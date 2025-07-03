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
