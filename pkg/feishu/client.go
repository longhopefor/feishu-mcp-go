package feishu

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/patrickmn/go-cache"
)

// Client 飞书API客户端
type Client struct {
	appID     string
	appSecret string
	baseURL   string
	client    *resty.Client
	cache     *cache.Cache
}

// AccessTokenResponse 获取访问令牌响应
type AccessTokenResponse struct {
	Code        int    `json:"code"`
	Msg         string `json:"msg"`
	AccessToken string `json:"tenant_access_token"`
	ExpireTime  int64  `json:"expire"`
}

// BaseResponse 基础响应结构
type BaseResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

// Document 文档相关数据结构
type Document struct {
	DocumentID string `json:"document_id"`
	Token      string `json:"token"`
	Title      string `json:"title"`
	OwnerID    string `json:"owner_id"`
	CreateTime string `json:"create_time"`
	UpdateTime string `json:"update_time"`
}

// CreateDocumentRequest 创建文档请求
type CreateDocumentRequest struct {
	FolderToken string `json:"folder_token"`
	Title       string `json:"title"`
}

// CreateDocumentResponse 创建文档响应
type CreateDocumentResponse struct {
	BaseResponse
	Data Document `json:"data"`
}

// DocumentInfo 文档信息
type DocumentInfo struct {
	Document Document `json:"document"`
}

// GetDocumentResponse 获取文档响应
type GetDocumentResponse struct {
	BaseResponse
	Data DocumentInfo `json:"data"`
}

// GetDocumentContentRequest 获取文档内容请求
type GetDocumentContentRequest struct {
	DocumentID string `json:"document_id"`
	Lang       int    `json:"lang,omitempty"` // 0: 中文, 1: 英文
}

// DocumentContent 文档内容
type DocumentContent struct {
	Content string `json:"content"`
}

// GetDocumentContentResponse 获取文档内容响应
type GetDocumentContentResponse struct {
	BaseResponse
	Data DocumentContent `json:"data"`
}

// Block 文档块结构
type Block struct {
	BlockID    string                 `json:"block_id"`
	BlockType  string                 `json:"block_type"`
	ParentID   string                 `json:"parent_id"`
	Children   []string               `json:"children"`
	Properties map[string]interface{} `json:"properties"`
}

// GetDocumentBlocksResponse 获取文档块响应
type GetDocumentBlocksResponse struct {
	BaseResponse
	Data struct {
		Blocks []Block `json:"blocks"`
	} `json:"data"`
}

// SearchDocumentsRequest 搜索文档请求
type SearchDocumentsRequest struct {
	SearchKey string `json:"search_key"`
	PageSize  int    `json:"page_size,omitempty"`
	PageToken string `json:"page_token,omitempty"`
}

// SearchDocument 搜索到的文档
type SearchDocument struct {
	DocumentID string `json:"document_id"`
	Token      string `json:"token"`
	Title      string `json:"title"`
	OwnerID    string `json:"owner_id"`
	CreateTime string `json:"create_time"`
	UpdateTime string `json:"update_time"`
	URL        string `json:"url"`
}

// SearchDocumentsResponse 搜索文档响应
type SearchDocumentsResponse struct {
	BaseResponse
	Data struct {
		Documents []SearchDocument `json:"documents"`
		PageToken string           `json:"page_token"`
		HasMore   bool             `json:"has_more"`
	} `json:"data"`
}

// NewClient 创建新的飞书客户端
func NewClient(appID, appSecret, baseURL string) (*Client, error) {
	if appID == "" || appSecret == "" {
		return nil, fmt.Errorf("应用ID和应用密钥不能为空")
	}

	if baseURL == "" {
		baseURL = "https://open.feishu.cn/open-apis"
	}

	client := resty.New().
		SetTimeout(30 * time.Second).
		SetRetryCount(3).
		SetRetryWaitTime(time.Second).
		SetRetryMaxWaitTime(5 * time.Second)

	// 创建缓存，默认过期时间1小时，清理间隔10分钟
	cache := cache.New(1*time.Hour, 10*time.Minute)

	return &Client{
		appID:     appID,
		appSecret: appSecret,
		baseURL:   baseURL,
		client:    client,
		cache:     cache,
	}, nil
}

// EnableDebugMode 启用调试模式
func (c *Client) EnableDebugMode(verbose bool) {
	monitor := NewRequestMonitor(true, verbose)
	c.client = monitor.EnableMonitoring(c.client)
}

// getAccessToken 获取访问令牌
func (c *Client) getAccessToken() (string, error) {
	cacheKey := "access_token"

	// 尝试从缓存获取
	if token, found := c.cache.Get(cacheKey); found {
		return token.(string), nil
	}

	// 请求新的访问令牌
	resp, err := c.client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(map[string]string{
			"app_id":     c.appID,
			"app_secret": c.appSecret,
		}).
		Post(c.baseURL + "/auth/v3/tenant_access_token/internal")

	if err != nil {
		return "", fmt.Errorf("请求访问令牌失败: %w", err)
	}

	var tokenResp AccessTokenResponse
	if err := json.Unmarshal(resp.Body(), &tokenResp); err != nil {
		return "", fmt.Errorf("解析访问令牌响应失败: %w", err)
	}

	if tokenResp.Code != 0 {
		return "", fmt.Errorf("获取访问令牌失败: %s", tokenResp.Msg)
	}

	// 缓存令牌（提前5分钟过期）
	expireDuration := time.Duration(tokenResp.ExpireTime-300) * time.Second
	c.cache.Set(cacheKey, tokenResp.AccessToken, expireDuration)

	return tokenResp.AccessToken, nil
}

// GetToken 获取访问令牌（公共方法，用于调试）
func (c *Client) GetToken() (string, error) {
	return c.getAccessToken()
}

// Get 发送GET请求
func (c *Client) Get(endpoint string) (*resty.Response, error) {
	return c.request("GET", endpoint, nil)
}

// Post 发送POST请求
func (c *Client) Post(endpoint string, body interface{}) (*resty.Response, error) {
	return c.request("POST", endpoint, body)
}

// Put 发送PUT请求
func (c *Client) Put(endpoint string, body interface{}) (*resty.Response, error) {
	return c.request("PUT", endpoint, body)
}

// Delete 发送DELETE请求
func (c *Client) Delete(endpoint string) (*resty.Response, error) {
	return c.request("DELETE", endpoint, nil)
}

// request 发送HTTP请求
func (c *Client) request(method, endpoint string, body interface{}) (*resty.Response, error) {
	// 获取访问令牌
	token, err := c.getAccessToken()
	if err != nil {
		return nil, err
	}

	// 构建请求
	req := c.client.R().
		SetHeader("Authorization", "Bearer "+token).
		SetHeader("Content-Type", "application/json")

	if body != nil {
		req.SetBody(body)
	}

	// 发送请求
	var resp *resty.Response
	switch method {
	case "GET":
		resp, err = req.Get(c.baseURL + endpoint)
	case "POST":
		resp, err = req.Post(c.baseURL + endpoint)
	case "PUT":
		resp, err = req.Put(c.baseURL + endpoint)
	case "DELETE":
		resp, err = req.Delete(c.baseURL + endpoint)
	default:
		return nil, fmt.Errorf("不支持的HTTP方法: %s", method)
	}

	if err != nil {
		return nil, fmt.Errorf("HTTP请求失败: %w", err)
	}

	return resp, nil
}

// CheckResponse 检查响应是否成功
func (c *Client) CheckResponse(resp *resty.Response) error {
	var baseResp BaseResponse
	if err := json.Unmarshal(resp.Body(), &baseResp); err != nil {
		return fmt.Errorf("解析响应失败: %w", err)
	}

	if baseResp.Code != 0 {
		return fmt.Errorf("API调用失败 (code: %d): %s", baseResp.Code, baseResp.Msg)
	}

	return nil
}

// GetBaseURL 获取基础URL
func (c *Client) GetBaseURL() string {
	return c.baseURL
}

// SetTimeout 设置请求超时时间
func (c *Client) SetTimeout(timeout time.Duration) {
	c.client.SetTimeout(timeout)
}

// SetRetry 设置重试配置
func (c *Client) SetRetry(count int, waitTime, maxWaitTime time.Duration) {
	c.client.SetRetryCount(count).
		SetRetryWaitTime(waitTime).
		SetRetryMaxWaitTime(maxWaitTime)
}

// CreateDocument 创建文档
func (c *Client) CreateDocument(folderToken, title string) (*Document, error) {
	req := CreateDocumentRequest{
		FolderToken: folderToken,
		Title:       title,
	}

	resp, err := c.Post("/docx/v1/documents", req)
	if err != nil {
		return nil, fmt.Errorf("创建文档请求失败: %w", err)
	}

	var createResp CreateDocumentResponse
	if err := json.Unmarshal(resp.Body(), &createResp); err != nil {
		return nil, fmt.Errorf("解析创建文档响应失败: %w", err)
	}

	if createResp.Code != 0 {
		return nil, fmt.Errorf("创建文档失败 (code: %d): %s", createResp.Code, createResp.Msg)
	}

	return &createResp.Data, nil
}

// GetDocument 获取文档信息
func (c *Client) GetDocument(documentID string) (*Document, error) {
	endpoint := fmt.Sprintf("/docx/v1/documents/%s", documentID)

	resp, err := c.Get(endpoint)
	if err != nil {
		return nil, fmt.Errorf("获取文档请求失败: %w", err)
	}

	var getResp GetDocumentResponse
	if err := json.Unmarshal(resp.Body(), &getResp); err != nil {
		return nil, fmt.Errorf("解析获取文档响应失败: %w", err)
	}

	if getResp.Code != 0 {
		return nil, fmt.Errorf("获取文档失败 (code: %d): %s", getResp.Code, getResp.Msg)
	}

	return &getResp.Data.Document, nil
}

// GetDocumentContent 获取文档内容
func (c *Client) GetDocumentContent(documentID string, lang int) (*DocumentContent, error) {
	endpoint := fmt.Sprintf("/docx/v1/documents/%s/content", documentID)

	// 构建查询参数
	reqURL := c.baseURL + endpoint
	if lang != 0 {
		reqURL += fmt.Sprintf("?lang=%d", lang)
	}

	// 获取访问令牌
	token, err := c.getAccessToken()
	if err != nil {
		return nil, err
	}

	resp, err := c.client.R().
		SetHeader("Authorization", "Bearer "+token).
		Get(reqURL)

	if err != nil {
		return nil, fmt.Errorf("获取文档内容请求失败: %w", err)
	}

	var contentResp GetDocumentContentResponse
	if err := json.Unmarshal(resp.Body(), &contentResp); err != nil {
		return nil, fmt.Errorf("解析文档内容响应失败: %w", err)
	}

	if contentResp.Code != 0 {
		return nil, fmt.Errorf("获取文档内容失败 (code: %d): %s", contentResp.Code, contentResp.Msg)
	}

	return &contentResp.Data, nil
}

// GetDocumentBlocks 获取文档块结构
func (c *Client) GetDocumentBlocks(documentID string) ([]Block, error) {
	endpoint := fmt.Sprintf("/docx/v1/documents/%s/blocks", documentID)

	resp, err := c.Get(endpoint)
	if err != nil {
		return nil, fmt.Errorf("获取文档块结构请求失败: %w", err)
	}

	var blocksResp GetDocumentBlocksResponse
	if err := json.Unmarshal(resp.Body(), &blocksResp); err != nil {
		return nil, fmt.Errorf("解析文档块结构响应失败: %w", err)
	}

	if blocksResp.Code != 0 {
		return nil, fmt.Errorf("获取文档块结构失败 (code: %d): %s", blocksResp.Code, blocksResp.Msg)
	}

	return blocksResp.Data.Blocks, nil
}

// SearchDocuments 搜索文档
func (c *Client) SearchDocuments(searchKey string, pageSize int, pageToken string) (*SearchDocumentsResponse, error) {
	// 使用正确的飞书搜索API端点 - 尝试多个可能的端点
	endpoints := []string{
		"/drive/v1/files/search", // 飞书云文档搜索API
		"/search/v2/documents",   // 可能的搜索端点
		"/search/v1/documents",   // 备用搜索端点
	}

	// 获取访问令牌
	token, err := c.getAccessToken()
	if err != nil {
		return nil, err
	}

	// 构建查询参数
	queryParams := map[string]string{
		"search_key": searchKey,
	}

	if pageSize > 0 {
		queryParams["page_size"] = fmt.Sprintf("%d", pageSize)
	}

	if pageToken != "" {
		queryParams["page_token"] = pageToken
	}

	var lastErr error

	// 尝试不同的端点
	for _, endpoint := range endpoints {
		resp, err := c.client.R().
			SetHeader("Authorization", "Bearer "+token).
			SetQueryParams(queryParams).
			Get(c.baseURL + endpoint)

		if err != nil {
			lastErr = fmt.Errorf("搜索文档请求失败 (端点: %s): %w", endpoint, err)
			continue
		}

		// 如果状态码是200，尝试解析响应
		if resp.StatusCode() == 200 {
			var searchResp SearchDocumentsResponse
			if err := json.Unmarshal(resp.Body(), &searchResp); err != nil {
				lastErr = fmt.Errorf("解析搜索文档响应失败 (端点: %s): %w", endpoint, err)
				continue
			}

			if searchResp.Code == 0 {
				// 成功
				return &searchResp, nil
			} else {
				lastErr = fmt.Errorf("搜索文档失败 (端点: %s, code: %d): %s", endpoint, searchResp.Code, searchResp.Msg)
				continue
			}
		} else {
			lastErr = fmt.Errorf("API请求失败 (端点: %s, 状态码: %d): %s", endpoint, resp.StatusCode(), resp.String())
			continue
		}
	}

	// 如果所有端点都失败，返回最后一个错误
	if lastErr != nil {
		return nil, lastErr
	}

	return nil, fmt.Errorf("搜索文档失败: 所有端点都无法访问")
}
