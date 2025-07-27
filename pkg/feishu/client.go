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
	BlockID   string                 `json:"block_id"`
	BlockType int                    `json:"block_type"` // 修正：API返回的是数字，不是字符串
	ParentID  string                 `json:"parent_id"`
	Children  []string               `json:"children,omitempty"`
	Page      map[string]interface{} `json:"page,omitempty"`    // 页面类型块的内容
	Text      map[string]interface{} `json:"text,omitempty"`    // 文本类型块的内容
	Code      map[string]interface{} `json:"code,omitempty"`    // 代码类型块的内容
	Heading   map[string]interface{} `json:"heading,omitempty"` // 标题类型块的内容
	List      map[string]interface{} `json:"list,omitempty"`    // 列表类型块的内容
}

// GetDocumentBlocksResponse 获取文档块响应
type GetDocumentBlocksResponse struct {
	BaseResponse
	Data struct {
		Items     []Block `json:"items"` // 修正：API返回的字段名是items，不是blocks
		HasMore   bool    `json:"has_more"`
		PageToken string  `json:"page_token"`
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

// FolderInfoResponse 文件夹信息响应
// FileInfo 文件信息结构体
type FileInfo struct {
	Name         string                 `json:"name"`
	ParentToken  string                 `json:"parent_token"`
	Token        string                 `json:"token"`
	Type         string                 `json:"type"`
	CreatedTime  string                 `json:"created_time"`
	ModifiedTime string                 `json:"modified_time"`
	OwnerID      string                 `json:"owner_id"`
	URL          string                 `json:"url"`
	ShortcutInfo map[string]interface{} `json:"shortcut_info,omitempty"`
}

// FolderInfoResponse 获取文件夹信息响应（修正为匹配飞书API实际响应格式）
type FolderInfoResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		Files   []FileInfo `json:"files"`
		HasMore bool       `json:"has_more"`
	} `json:"data"`
}

// DriveFile 云空间文件信息
type DriveFile struct {
	Token        string `json:"token"`
	Name         string `json:"name"`
	Type         string `json:"type"` // file, folder
	ParentToken  string `json:"parent_token"`
	URL          string `json:"url"`
	Size         int64  `json:"size,omitempty"`
	CreatedTime  string `json:"created_time"`
	ModifiedTime string `json:"modified_time"`
	OwnerID      string `json:"owner_id"`
	Creator      string `json:"creator"`
	Thumbnail    string `json:"thumbnail,omitempty"`
	MimeType     string `json:"mime_type,omitempty"`
}

// GetDriveFilesResponse 获取云空间文件列表响应
type GetDriveFilesResponse struct {
	BaseResponse
	Data struct {
		Files         []DriveFile `json:"files"`
		NextPageToken string      `json:"next_page_token"`
		HasMore       bool        `json:"has_more"`
	} `json:"data"`
}

// GetDriveMetaResponse 获取云空间目录元数据响应
type GetDriveMetaResponse struct {
	BaseResponse
	Data struct {
		Token        string `json:"token"`
		Name         string `json:"name"`
		Type         string `json:"type"`
		ParentToken  string `json:"parent_token"`
		URL          string `json:"url"`
		Size         int64  `json:"size,omitempty"`
		CreatedTime  string `json:"created_time"`
		ModifiedTime string `json:"modified_time"`
		OwnerID      string `json:"owner_id"`
		Creator      string `json:"creator"`
	} `json:"data"`
}

// RootFolderMeta 根文件夹元数据（修正为匹配实际API响应）
type RootFolderMeta struct {
	Token  string `json:"token"`
	ID     string `json:"id"`
	UserID string `json:"user_id"`
}

// GetRootFolderMetaResponse 获取根文件夹元数据响应
type GetRootFolderMetaResponse struct {
	BaseResponse
	Data RootFolderMeta `json:"data"`
}

// FolderDetailResponse 文件夹详细信息响应
type FolderDetailResponse struct {
	Code   int    `json:"code"`
	Msg    string `json:"msg"`
	Folder struct {
		FolderToken string `json:"folder_token"`
		FolderName  string `json:"folder_name"`
	} `json:"folder"`
}

// CreateFolderRequest 创建文件夹请求
type CreateFolderRequest struct {
	Name        string `json:"name"`
	ParentToken string `json:"parent_token,omitempty"`
}

// CreateFolderResponse 创建文件夹响应
type CreateFolderResponse struct {
	BaseResponse
	Data struct {
		Token string `json:"token"`
		Name  string `json:"name"`
		URL   string `json:"url"`
	} `json:"data"`
}

// WikiSpace 知识库空间结构
type WikiSpace struct {
	SpaceID     string `json:"space_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	SpaceType   string `json:"space_type"`
	Visibility  string `json:"visibility"`
}

// GetWikiSpacesResponse 获取知识库空间响应
type GetWikiSpacesResponse struct {
	BaseResponse
	Data struct {
		Items     []WikiSpace `json:"items"` // 修正：API返回的字段名是items，不是spaces
		PageToken string      `json:"page_token"`
		HasMore   bool        `json:"has_more"`
	} `json:"data"`
}

// WikiNode 知识库节点结构
type WikiNode struct {
	SpaceID    string `json:"space_id"`
	NodeToken  string `json:"node_token"`
	ObjToken   string `json:"obj_token"`
	ObjType    string `json:"obj_type"`
	ParentNode string `json:"parent_node_token"`
	NodeType   string `json:"node_type"`
	OriginNode string `json:"origin_node_token"`
	Title      string `json:"title"`
	HasChild   bool   `json:"has_child"`
}

// GetWikiNodesResponse 获取知识库节点响应
type GetWikiNodesResponse struct {
	BaseResponse
	Data struct {
		Items     []WikiNode `json:"items"`
		PageToken string     `json:"page_token"`
		HasMore   bool       `json:"has_more"`
	} `json:"data"`
}

// WikiNodeContent 知识库节点内容
type WikiNodeContent struct {
	Content string `json:"content"`
	Type    string `json:"type"`
}

// GetWikiNodeContentResponse 获取知识库节点内容响应
type GetWikiNodeContentResponse struct {
	BaseResponse
	Data WikiNodeContent `json:"data"`
}

// WikiNodeMeta 知识库节点元信息
type WikiNodeMeta struct {
	SpaceID        string `json:"space_id"`
	NodeToken      string `json:"node_token"`
	ObjToken       string `json:"obj_token"`
	ObjType        string `json:"obj_type"`
	ParentNode     string `json:"parent_node_token"`
	NodeType       string `json:"node_type"`
	OriginNode     string `json:"origin_node_token"`
	OriginSpaceID  string `json:"origin_space_id"`
	Title          string `json:"title"`
	HasChild       bool   `json:"has_child"`
	NodeCreateTime string `json:"node_create_time"` // 修正：API实际字段名
	NodeCreator    string `json:"node_creator"`     // 修正：API实际字段名
	Creator        string `json:"creator"`
	Owner          string `json:"owner"`
	ObjCreateTime  string `json:"obj_create_time"` // 新增：对象创建时间
	ObjEditTime    string `json:"obj_edit_time"`   // 新增：对象编辑时间
}

// GetWikiNodeMetaResponse 获取知识库节点元信息响应
type GetWikiNodeMetaResponse struct {
	BaseResponse
	Data struct {
		Node WikiNodeMeta `json:"node"` // 修正：API返回的数据在data.node中
	} `json:"data"`
}

// BlockContent 块内容
type BlockContent struct {
	BlockID   string                 `json:"block_id"`
	BlockType int                    `json:"block_type"` // 修正：API返回的是数字
	ParentID  string                 `json:"parent_id"`  // 添加父块ID
	Children  []string               `json:"children,omitempty"`
	Page      map[string]interface{} `json:"page,omitempty"`    // 页面类型块的内容
	Text      map[string]interface{} `json:"text,omitempty"`    // 文本类型块的内容
	Code      map[string]interface{} `json:"code,omitempty"`    // 代码类型块的内容
	Heading   map[string]interface{} `json:"heading,omitempty"` // 标题类型块的内容
	List      map[string]interface{} `json:"list,omitempty"`    // 列表类型块的内容
}

// GetBlockContentResponse 获取块内容响应
type GetBlockContentResponse struct {
	BaseResponse
	Data struct {
		Block BlockContent `json:"block"` // 修正：API返回的数据在data.block中
	} `json:"data"`
}

// CreateBlockRequest 创建块请求
type CreateBlockRequest struct {
	BlockType string                 `json:"block_type"`
	ParentID  string                 `json:"parent_id,omitempty"`
	Content   map[string]interface{} `json:"content"`
}

// CreateBlockResponse 创建块响应
type CreateBlockResponse struct {
	BaseResponse
	Data struct {
		BlockID string `json:"block_id"`
	} `json:"data"`
}

// BatchCreateBlocksRequest 批量创建块请求
type BatchCreateBlocksRequest struct {
	Blocks []CreateBlockRequest `json:"blocks"`
}

// BatchCreateBlocksResponse 批量创建块响应
type BatchCreateBlocksResponse struct {
	BaseResponse
	Data struct {
		BlockIDs []string `json:"block_ids"`
	} `json:"data"`
}

// UpdateBlockRequest 更新块请求
type UpdateBlockRequest struct {
	Content map[string]interface{} `json:"content"`
}

// UpdateBlockResponse 更新块响应
type UpdateBlockResponse struct {
	BaseResponse
	Data struct {
		BlockID string `json:"block_id"`
	} `json:"data"`
}

// DeleteBlocksRequest 删除块请求
type DeleteBlocksRequest struct {
	StartIndex int `json:"start_index"`
	EndIndex   int `json:"end_index"`
}

// DeleteBlocksResponse 删除块响应
type DeleteBlocksResponse struct {
	BaseResponse
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
	// 尝试多个可能的API端点
	endpoints := []string{
		fmt.Sprintf("/docx/v1/documents/%s", documentID),   // 标准文档端点
		fmt.Sprintf("/drive/v1/files/%s", documentID),      // 云文档端点
		fmt.Sprintf("/drive/v1/files/%s/meta", documentID), // 文档元信息端点
	}

	// 获取访问令牌
	token, err := c.getAccessToken()
	if err != nil {
		return nil, err
	}

	var lastErr error

	// 尝试不同的端点
	for _, endpoint := range endpoints {
		resp, err := c.client.R().
			SetHeader("Authorization", "Bearer "+token).
			Get(c.baseURL + endpoint)

		if err != nil {
			lastErr = fmt.Errorf("获取文档请求失败 (端点: %s): %w", endpoint, err)
			continue
		}

		// 如果状态码是200，尝试解析响应
		if resp.StatusCode() == 200 {
			var getResp GetDocumentResponse
			if err := json.Unmarshal(resp.Body(), &getResp); err != nil {
				lastErr = fmt.Errorf("解析获取文档响应失败 (端点: %s): %w", endpoint, err)
				continue
			}

			if getResp.Code == 0 {
				// 成功
				return &getResp.Data.Document, nil
			} else {
				lastErr = fmt.Errorf("获取文档失败 (端点: %s, code: %d): %s", endpoint, getResp.Code, getResp.Msg)
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

	return nil, fmt.Errorf("获取文档失败: 所有端点都无法访问")
}

// GetDocumentContent 获取文档内容
func (c *Client) GetDocumentContent(documentID string, lang int) (*DocumentContent, error) {
	// 尝试多个可能的API端点
	endpoints := []string{
		fmt.Sprintf("/docx/v1/documents/%s/raw_content", documentID), // 原始内容端点
		fmt.Sprintf("/drive/v1/files/%s/content", documentID),        // 云文档内容端点
		fmt.Sprintf("/docx/v1/documents/%s/content", documentID),     // 原端点（备用）
	}

	// 获取访问令牌
	token, err := c.getAccessToken()
	if err != nil {
		return nil, err
	}

	var lastErr error

	// 尝试不同的端点
	for _, endpoint := range endpoints {
		// 构建查询参数
		reqURL := c.baseURL + endpoint
		if lang != 0 {
			reqURL += fmt.Sprintf("?lang=%d", lang)
		}

		resp, err := c.client.R().
			SetHeader("Authorization", "Bearer "+token).
			Get(reqURL)

		if err != nil {
			lastErr = fmt.Errorf("获取文档内容请求失败 (端点: %s): %w", endpoint, err)
			continue
		}

		// 如果状态码是200，尝试解析响应
		if resp.StatusCode() == 200 {
			var contentResp GetDocumentContentResponse
			if err := json.Unmarshal(resp.Body(), &contentResp); err != nil {
				lastErr = fmt.Errorf("解析文档内容响应失败 (端点: %s): %w", endpoint, err)
				continue
			}

			if contentResp.Code == 0 {
				// 成功
				return &contentResp.Data, nil
			} else {
				lastErr = fmt.Errorf("获取文档内容失败 (端点: %s, code: %d): %s", endpoint, contentResp.Code, contentResp.Msg)
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

	return nil, fmt.Errorf("获取文档内容失败: 所有端点都无法访问")
}

// GetDocumentBlocks 获取文档块结构
func (c *Client) GetDocumentBlocks(documentID string) ([]Block, error) {
	// 尝试多个可能的API端点
	endpoints := []string{
		fmt.Sprintf("/docx/v1/documents/%s/blocks", documentID),     // 标准块结构端点
		fmt.Sprintf("/drive/v1/files/%s/blocks", documentID),        // 云文档块端点
		fmt.Sprintf("/docx/v1/documents/%s/raw_blocks", documentID), // 原始块端点
	}

	// 获取访问令牌
	token, err := c.getAccessToken()
	if err != nil {
		return nil, err
	}

	var lastErr error

	// 尝试不同的端点
	for _, endpoint := range endpoints {
		resp, err := c.client.R().
			SetHeader("Authorization", "Bearer "+token).
			Get(c.baseURL + endpoint)

		if err != nil {
			lastErr = fmt.Errorf("获取文档块结构请求失败 (端点: %s): %w", endpoint, err)
			continue
		}

		// 如果状态码是200，尝试解析响应
		if resp.StatusCode() == 200 {
			var blocksResp GetDocumentBlocksResponse
			if err := json.Unmarshal(resp.Body(), &blocksResp); err != nil {
				lastErr = fmt.Errorf("解析文档块结构响应失败 (端点: %s): %w", endpoint, err)
				continue
			}

			if blocksResp.Code == 0 {
				// 成功
				return blocksResp.Data.Items, nil
			} else {
				lastErr = fmt.Errorf("获取文档块结构失败 (端点: %s, code: %d): %s", endpoint, blocksResp.Code, blocksResp.Msg)
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

	return nil, fmt.Errorf("获取文档块结构失败: 所有端点都无法访问")
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

// GetRootFolderInfo 获取根目录信息（修正为匹配飞书API实际响应格式）
func (c *Client) GetRootFolderInfo() (*FolderInfoResponse, error) {
	endpoint := "/drive/v1/files"

	// 获取访问令牌
	token, err := c.getAccessToken()
	if err != nil {
		return nil, err
	}

	resp, err := c.client.R().
		SetHeader("Authorization", "Bearer "+token).
		SetQueryParam("folder_token", ""). // 空值表示根目录
		SetQueryParam("page_size", "50").  // 增加页面大小以获取更多文件
		Get(c.baseURL + endpoint)

	if err != nil {
		return nil, fmt.Errorf("获取根目录信息失败: %w", err)
	}

	var folderResp FolderInfoResponse
	if err := json.Unmarshal(resp.Body(), &folderResp); err != nil {
		return nil, fmt.Errorf("解析根目录信息响应失败: %w", err)
	}

	if folderResp.Code != 0 {
		return nil, fmt.Errorf("获取根目录信息失败: %s", folderResp.Msg)
	}

	return &folderResp, nil
}

// GetRootFolderMeta 获取根文件夹元数据（使用新的API端点）
func (c *Client) GetRootFolderMeta() (*GetRootFolderMetaResponse, error) {
	endpoint := "/drive/explorer/v2/root_folder/meta"

	// 获取访问令牌
	token, err := c.getAccessToken()
	if err != nil {
		return nil, fmt.Errorf("获取访问令牌失败: %w", err)
	}

	resp, err := c.client.R().
		SetHeader("Authorization", "Bearer "+token).
		Get(c.baseURL + endpoint)

	if err != nil {
		return nil, fmt.Errorf("获取根文件夹元数据失败: %w", err)
	}

	var metaResp GetRootFolderMetaResponse
	if err := json.Unmarshal(resp.Body(), &metaResp); err != nil {
		return nil, fmt.Errorf("解析根文件夹元数据响应失败: %w", err)
	}

	if metaResp.Code != 0 {
		return nil, fmt.Errorf("获取根文件夹元数据失败 (code: %d): %s", metaResp.Code, metaResp.Msg)
	}

	return &metaResp, nil
}

// GetFolderInfo 获取指定文件夹信息
func (c *Client) GetFolderInfo(folderToken string) (*FolderDetailResponse, error) {
	endpoint := fmt.Sprintf("/drive/v1/folders/%s", folderToken)

	// 获取访问令牌
	token, err := c.getAccessToken()
	if err != nil {
		return nil, err
	}

	resp, err := c.client.R().
		SetHeader("Authorization", "Bearer "+token).
		Get(c.baseURL + endpoint)

	if err != nil {
		return nil, fmt.Errorf("获取文件夹信息失败: %w", err)
	}

	var folderResp FolderDetailResponse
	if err := json.Unmarshal(resp.Body(), &folderResp); err != nil {
		return nil, fmt.Errorf("解析文件夹信息响应失败: %w", err)
	}

	if folderResp.Code != 0 {
		return nil, fmt.Errorf("获取文件夹信息失败: %s", folderResp.Msg)
	}

	return &folderResp, nil
}

// GetFolderFiles 获取文件夹下的文件列表
func (c *Client) GetFolderFiles(folderToken string, pageSize int, pageToken string) (*FolderInfoResponse, error) {
	endpoint := "/drive/v1/files"

	// 获取访问令牌
	token, err := c.getAccessToken()
	if err != nil {
		return nil, err
	}

	req := c.client.R().
		SetHeader("Authorization", "Bearer "+token).
		SetQueryParam("folder_token", folderToken)

	if pageSize > 0 {
		req.SetQueryParam("page_size", fmt.Sprintf("%d", pageSize))
	}

	if pageToken != "" {
		req.SetQueryParam("page_token", pageToken)
	}

	resp, err := req.Get(c.baseURL + endpoint)

	if err != nil {
		return nil, fmt.Errorf("获取文件夹文件列表失败: %w", err)
	}

	var folderResp FolderInfoResponse
	if err := json.Unmarshal(resp.Body(), &folderResp); err != nil {
		return nil, fmt.Errorf("解析文件夹文件列表响应失败: %w", err)
	}

	if folderResp.Code != 0 {
		return nil, fmt.Errorf("获取文件夹文件列表失败: %s", folderResp.Msg)
	}

	return &folderResp, nil
}

// GetDriveFilesWithMeta 获取云空间目录下所有文件的详细信息
func (c *Client) GetDriveFilesWithMeta(folderToken string, pageSize int, pageToken string) (*GetDriveFilesResponse, error) {
	endpoint := "/drive/v1/files"

	// 获取访问令牌
	token, err := c.getAccessToken()
	if err != nil {
		return nil, fmt.Errorf("获取访问令牌失败: %w", err)
	}

	req := c.client.R().
		SetHeader("Authorization", "Bearer "+token)

	// 设置查询参数
	if folderToken != "" {
		req.SetQueryParam("folder_token", folderToken)
	}
	if pageSize > 0 {
		req.SetQueryParam("page_size", fmt.Sprintf("%d", pageSize))
	} else {
		req.SetQueryParam("page_size", "50") // 默认分页大小
	}
	if pageToken != "" {
		req.SetQueryParam("page_token", pageToken)
	}

	// 设置返回字段，获取文件的详细元数据
	req.SetQueryParam("fields", "name,type,parent_token,url,size,created_time,modified_time,owner_id,creator")

	resp, err := req.Get(c.baseURL + endpoint)
	if err != nil {
		return nil, fmt.Errorf("获取云空间文件列表失败: %w", err)
	}

	var filesResp GetDriveFilesResponse
	if err := json.Unmarshal(resp.Body(), &filesResp); err != nil {
		return nil, fmt.Errorf("解析云空间文件列表响应失败: %w", err)
	}

	if filesResp.Code != 0 {
		return nil, fmt.Errorf("获取云空间文件列表失败 (code: %d): %s", filesResp.Code, filesResp.Msg)
	}

	return &filesResp, nil
}

// GetDriveMeta 获取云空间目录/文件的元数据信息
func (c *Client) GetDriveMeta(fileToken string) (*GetDriveMetaResponse, error) {
	endpoint := fmt.Sprintf("/drive/v1/files/%s/meta", fileToken)

	// 获取访问令牌
	token, err := c.getAccessToken()
	if err != nil {
		return nil, fmt.Errorf("获取访问令牌失败: %w", err)
	}

	resp, err := c.client.R().
		SetHeader("Authorization", "Bearer "+token).
		Get(c.baseURL + endpoint)

	if err != nil {
		return nil, fmt.Errorf("获取云空间文件元数据失败: %w", err)
	}

	var metaResp GetDriveMetaResponse
	if err := json.Unmarshal(resp.Body(), &metaResp); err != nil {
		return nil, fmt.Errorf("解析云空间文件元数据响应失败: %w", err)
	}

	if metaResp.Code != 0 {
		return nil, fmt.Errorf("获取云空间文件元数据失败 (code: %d): %s", metaResp.Code, metaResp.Msg)
	}

	return &metaResp, nil
}

// GetAllDriveFiles 获取指定目录下所有文件（递归获取，支持分页）
func (c *Client) GetAllDriveFiles(folderToken string, maxFiles int) ([]DriveFile, error) {
	var allFiles []DriveFile
	pageToken := ""
	pageSize := 50

	if maxFiles > 0 && maxFiles < pageSize {
		pageSize = maxFiles
	}

	for {
		resp, err := c.GetDriveFilesWithMeta(folderToken, pageSize, pageToken)
		if err != nil {
			return nil, fmt.Errorf("获取文件列表失败: %w", err)
		}

		allFiles = append(allFiles, resp.Data.Files...)

		// 检查是否达到最大文件数限制
		if maxFiles > 0 && len(allFiles) >= maxFiles {
			if len(allFiles) > maxFiles {
				allFiles = allFiles[:maxFiles]
			}
			break
		}

		// 检查是否还有更多数据
		if !resp.Data.HasMore || resp.Data.NextPageToken == "" {
			break
		}

		pageToken = resp.Data.NextPageToken
	}

	return allFiles, nil
}

// CreateFolder 创建文件夹
func (c *Client) CreateFolder(name, parentToken string) (*CreateFolderResponse, error) {
	endpoint := "/drive/v1/folders"

	// 获取访问令牌
	token, err := c.getAccessToken()
	if err != nil {
		return nil, err
	}

	// 构建请求体
	requestBody := CreateFolderRequest{
		Name:        name,
		ParentToken: parentToken,
	}

	resp, err := c.client.R().
		SetHeader("Authorization", "Bearer "+token).
		SetHeader("Content-Type", "application/json").
		SetBody(requestBody).
		Post(c.baseURL + endpoint)

	if err != nil {
		return nil, fmt.Errorf("创建文件夹请求失败: %w", err)
	}

	if err := c.CheckResponse(resp); err != nil {
		return nil, err
	}

	var createResp CreateFolderResponse
	if err := json.Unmarshal(resp.Body(), &createResp); err != nil {
		return nil, fmt.Errorf("解析创建文件夹响应失败: %w", err)
	}

	if createResp.Code != 0 {
		return nil, fmt.Errorf("创建文件夹失败: %s", createResp.Msg)
	}

	return &createResp, nil
}

// ==================== 知识库相关方法 ====================

// GetWikiSpaces 获取知识库空间列表
func (c *Client) GetWikiSpaces(pageSize int, pageToken string) (*GetWikiSpacesResponse, error) {
	// 知识库空间API端点
	endpoints := []string{
		"/wiki/v2/spaces",  // 标准知识库空间端点
		"/wiki/v1/spaces",  // 备用端点
		"/drive/v1/spaces", // 云文档空间端点
	}

	// 获取访问令牌
	token, err := c.getAccessToken()
	if err != nil {
		return nil, err
	}

	// 构建查询参数
	queryParams := map[string]string{}
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
			lastErr = fmt.Errorf("获取知识库空间请求失败 (端点: %s): %w", endpoint, err)
			continue
		}

		// 如果状态码是200，尝试解析响应
		if resp.StatusCode() == 200 {
			var spacesResp GetWikiSpacesResponse
			if err := json.Unmarshal(resp.Body(), &spacesResp); err != nil {
				lastErr = fmt.Errorf("解析知识库空间响应失败 (端点: %s): %w", endpoint, err)
				continue
			}

			if spacesResp.Code == 0 {
				// 成功
				return &spacesResp, nil
			} else {
				lastErr = fmt.Errorf("获取知识库空间失败 (端点: %s, code: %d): %s", endpoint, spacesResp.Code, spacesResp.Msg)
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

	return nil, fmt.Errorf("获取知识库空间失败: 所有端点都无法访问")
}

// GetWikiSpaceNodes 获取知识库空间下的节点列表
func (c *Client) GetWikiSpaceNodes(spaceID string, pageSize int, pageToken string, parentNodeToken string) (*GetWikiNodesResponse, error) {
	// 知识库节点API端点
	endpoints := []string{
		fmt.Sprintf("/wiki/v2/spaces/%s/nodes", spaceID),  // 标准知识库节点端点
		fmt.Sprintf("/wiki/v1/spaces/%s/nodes", spaceID),  // 备用端点
		fmt.Sprintf("/drive/v1/spaces/%s/nodes", spaceID), // 云文档节点端点
	}

	// 获取访问令牌
	token, err := c.getAccessToken()
	if err != nil {
		return nil, err
	}

	// 构建查询参数
	queryParams := map[string]string{}
	if pageSize > 0 {
		queryParams["page_size"] = fmt.Sprintf("%d", pageSize)
	}
	if pageToken != "" {
		queryParams["page_token"] = pageToken
	}
	if parentNodeToken != "" {
		queryParams["parent_node_token"] = parentNodeToken
	}

	var lastErr error

	// 尝试不同的端点
	for _, endpoint := range endpoints {
		resp, err := c.client.R().
			SetHeader("Authorization", "Bearer "+token).
			SetQueryParams(queryParams).
			Get(c.baseURL + endpoint)

		if err != nil {
			lastErr = fmt.Errorf("获取知识库节点请求失败 (端点: %s): %w", endpoint, err)
			continue
		}

		// 如果状态码是200，尝试解析响应
		if resp.StatusCode() == 200 {
			var nodesResp GetWikiNodesResponse
			if err := json.Unmarshal(resp.Body(), &nodesResp); err != nil {
				lastErr = fmt.Errorf("解析知识库节点响应失败 (端点: %s): %w", endpoint, err)
				continue
			}

			if nodesResp.Code == 0 {
				// 成功
				return &nodesResp, nil
			} else {
				lastErr = fmt.Errorf("获取知识库节点失败 (端点: %s, code: %d): %s", endpoint, nodesResp.Code, nodesResp.Msg)
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

	return nil, fmt.Errorf("获取知识库节点失败: 所有端点都无法访问")
}

// GetWikiNodeContent 获取知识库节点内容
func (c *Client) GetWikiNodeContent(spaceID, nodeToken string, lang int) (*WikiNodeContent, error) {
	// 知识库节点内容API端点
	endpoints := []string{
		fmt.Sprintf("/wiki/v2/spaces/%s/nodes/%s/content", spaceID, nodeToken),     // 标准知识库内容端点
		fmt.Sprintf("/wiki/v1/spaces/%s/nodes/%s/content", spaceID, nodeToken),     // 备用端点v1
		fmt.Sprintf("/drive/v1/spaces/%s/nodes/%s/content", spaceID, nodeToken),    // 云文档内容端点
		fmt.Sprintf("/wiki/v2/spaces/%s/nodes/%s/raw_content", spaceID, nodeToken), // 原始内容端点
	}

	// 获取访问令牌
	token, err := c.getAccessToken()
	if err != nil {
		return nil, err
	}

	// 构建查询参数
	queryParams := map[string]string{}
	if lang > 0 {
		queryParams["lang"] = fmt.Sprintf("%d", lang)
	}

	var lastErr error

	// 尝试不同的端点
	for _, endpoint := range endpoints {
		resp, err := c.client.R().
			SetHeader("Authorization", "Bearer "+token).
			SetQueryParams(queryParams).
			Get(c.baseURL + endpoint)

		if err != nil {
			lastErr = fmt.Errorf("获取知识库节点内容请求失败 (端点: %s): %w", endpoint, err)
			continue
		}

		// 如果状态码是200，尝试解析响应
		if resp.StatusCode() == 200 {
			var contentResp GetWikiNodeContentResponse
			if err := json.Unmarshal(resp.Body(), &contentResp); err != nil {
				lastErr = fmt.Errorf("解析知识库节点内容响应失败 (端点: %s): %w", endpoint, err)
				continue
			}

			if contentResp.Code == 0 {
				// 成功
				return &contentResp.Data, nil
			} else {
				lastErr = fmt.Errorf("获取知识库节点内容失败 (端点: %s, code: %d): %s", endpoint, contentResp.Code, contentResp.Msg)
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

	return nil, fmt.Errorf("获取知识库节点内容失败: 所有端点都无法访问")
}

// GetWikiNodeMeta 获取知识库节点元信息
func (c *Client) GetWikiNodeMeta(spaceID, nodeToken string) (*WikiNodeMeta, error) {
	// 知识库节点元信息API端点
	endpoints := []string{
		fmt.Sprintf("/wiki/v2/spaces/%s/nodes/%s", spaceID, nodeToken),      // 标准知识库元信息端点
		fmt.Sprintf("/wiki/v1/spaces/%s/nodes/%s", spaceID, nodeToken),      // 备用端点v1
		fmt.Sprintf("/drive/v1/spaces/%s/nodes/%s", spaceID, nodeToken),     // 云文档元信息端点
		fmt.Sprintf("/wiki/v2/spaces/%s/nodes/%s/meta", spaceID, nodeToken), // 元信息端点
	}

	// 获取访问令牌
	token, err := c.getAccessToken()
	if err != nil {
		return nil, err
	}

	var lastErr error

	// 尝试不同的端点
	for _, endpoint := range endpoints {
		resp, err := c.client.R().
			SetHeader("Authorization", "Bearer "+token).
			Get(c.baseURL + endpoint)

		if err != nil {
			lastErr = fmt.Errorf("获取知识库节点元信息请求失败 (端点: %s): %w", endpoint, err)
			continue
		}

		// 如果状态码是200，尝试解析响应
		if resp.StatusCode() == 200 {
			var metaResp GetWikiNodeMetaResponse
			if err := json.Unmarshal(resp.Body(), &metaResp); err != nil {
				lastErr = fmt.Errorf("解析知识库节点元信息响应失败 (端点: %s): %w", endpoint, err)
				continue
			}

			if metaResp.Code == 0 {
				// 成功
				return &metaResp.Data.Node, nil
			} else {
				lastErr = fmt.Errorf("获取知识库节点元信息失败 (端点: %s, code: %d): %s", endpoint, metaResp.Code, metaResp.Msg)
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

	return nil, fmt.Errorf("获取知识库节点元信息失败: 所有端点都无法访问")
}

// GetBlockContent 获取特定块的详细内容
func (c *Client) GetBlockContent(documentID, blockID string) (*BlockContent, error) {
	token, err := c.getAccessToken()
	if err != nil {
		return nil, err
	}

	endpoint := fmt.Sprintf("/docx/v1/documents/%s/blocks/%s", documentID, blockID)
	resp, err := c.client.R().
		SetHeader("Authorization", "Bearer "+token).
		Get(c.baseURL + endpoint)

	if err != nil {
		return nil, fmt.Errorf("获取块内容请求失败: %w", err)
	}

	if err := c.CheckResponse(resp); err != nil {
		return nil, err
	}

	var blockContentResp GetBlockContentResponse
	if err := json.Unmarshal(resp.Body(), &blockContentResp); err != nil {
		return nil, fmt.Errorf("解析块内容响应失败: %w", err)
	}

	if blockContentResp.Code != 0 {
		return nil, fmt.Errorf("获取块内容失败: %s", blockContentResp.Msg)
	}

	return &blockContentResp.Data.Block, nil
}

// UpdateBlockText 更新块的文本内容
func (c *Client) UpdateBlockText(documentID, blockID string, content map[string]interface{}) error {
	token, err := c.getAccessToken()
	if err != nil {
		return err
	}

	endpoint := fmt.Sprintf("/docx/v1/documents/%s/blocks/%s", documentID, blockID)

	// 从content参数中提取文本内容
	var textContent string
	var hasValidContent bool

	// 尝试从不同的字段中提取文本内容
	if text, exists := content["text"]; exists {
		// 如果text是字符串，直接使用
		if textStr, ok := text.(string); ok {
			textContent = textStr
			hasValidContent = true
		} else if textMap, ok := text.(map[string]interface{}); ok {
			// 如果text是map，尝试从elements中提取内容
			if elements, ok := textMap["elements"].([]interface{}); ok && len(elements) > 0 {
				if element, ok := elements[0].(map[string]interface{}); ok {
					if textRun, ok := element["text_run"].(map[string]interface{}); ok {
						if content, ok := textRun["content"].(string); ok {
							textContent = content
							hasValidContent = true
						}
					}
				}
			}
		}
	}

	// 如果没有找到text字段，尝试content字段
	if !hasValidContent {
		if text, exists := content["content"]; exists {
			if textStr, ok := text.(string); ok {
				textContent = textStr
				hasValidContent = true
			} else {
				return fmt.Errorf("content字段必须是字符串类型")
			}
		}
	}

	// 如果没有找到有效的文本内容，返回错误
	if !hasValidContent {
		return fmt.Errorf("content参数中必须包含有效的text或content字段")
	}

	// 使用update_text_elements字段来更新块内容
	requestBody := map[string]interface{}{
		"update_text_elements": map[string]interface{}{
			"elements": []interface{}{
				map[string]interface{}{
					"text_run": map[string]interface{}{
						"content": textContent,
						"text_element_style": map[string]interface{}{
							"bold":          false,
							"inline_code":   false,
							"italic":        false,
							"strikethrough": false,
							"underline":     false,
						},
					},
				},
			},
			"style": map[string]interface{}{
				"align":  1,
				"folded": false,
			},
		},
	}

	resp, err := c.client.R().
		SetHeader("Authorization", "Bearer "+token).
		SetHeader("Content-Type", "application/json").
		SetBody(requestBody).
		Patch(c.baseURL + endpoint)

	if err != nil {
		return fmt.Errorf("更新块内容请求失败: %w", err)
	}

	if err := c.CheckResponse(resp); err != nil {
		return err
	}

	var updateResp UpdateBlockResponse
	if err := json.Unmarshal(resp.Body(), &updateResp); err != nil {
		return fmt.Errorf("解析更新响应失败: %w", err)
	}

	if updateResp.Code != 0 {
		return fmt.Errorf("更新块内容失败: %s", updateResp.Msg)
	}

	return nil
}

// CreateBlock 创建单个块
func (c *Client) CreateBlock(documentID string, blockType, parentID string, content map[string]interface{}) (string, error) {
	token, err := c.getAccessToken()
	if err != nil {
		return "", err
	}

	// 尝试使用不同的API端点格式，在父块下创建子块
	endpoint := fmt.Sprintf("/docx/v1/documents/%s/blocks/%s/children", documentID, parentID)

	// 创建最简单的文本块请求
	blockTypeInt := 2 // 默认文本块类型
	if blockType == "2" {
		blockTypeInt = 2
	}

	createReq := map[string]interface{}{
		"block_type": blockTypeInt,
		"text":       content["text"],
	}

	resp, err := c.client.R().
		SetHeader("Authorization", "Bearer "+token).
		SetHeader("Content-Type", "application/json").
		SetBody(createReq).
		Post(c.baseURL + endpoint)

	if err != nil {
		return "", fmt.Errorf("创建块请求失败: %w", err)
	}

	if err := c.CheckResponse(resp); err != nil {
		return "", err
	}

	var createResp CreateBlockResponse
	if err := json.Unmarshal(resp.Body(), &createResp); err != nil {
		return "", fmt.Errorf("解析创建响应失败: %w", err)
	}

	if createResp.Code != 0 {
		return "", fmt.Errorf("创建块失败: %s", createResp.Msg)
	}

	return createResp.Data.BlockID, nil
}

// BatchCreateBlocks 批量创建多个块
func (c *Client) BatchCreateBlocks(documentID string, blocks []CreateBlockRequest) ([]string, error) {
	token, err := c.getAccessToken()
	if err != nil {
		return nil, err
	}

	endpoint := fmt.Sprintf("/docx/v1/documents/%s/blocks/batch_create", documentID)
	batchReq := BatchCreateBlocksRequest{
		Blocks: blocks,
	}

	resp, err := c.client.R().
		SetHeader("Authorization", "Bearer "+token).
		SetHeader("Content-Type", "application/json").
		SetBody(batchReq).
		Post(c.baseURL + endpoint)

	if err != nil {
		return nil, fmt.Errorf("批量创建块请求失败: %w", err)
	}

	if err := c.CheckResponse(resp); err != nil {
		return nil, err
	}

	var batchResp BatchCreateBlocksResponse
	if err := json.Unmarshal(resp.Body(), &batchResp); err != nil {
		return nil, fmt.Errorf("解析批量创建响应失败: %w", err)
	}

	if batchResp.Code != 0 {
		return nil, fmt.Errorf("批量创建块失败: %s", batchResp.Msg)
	}

	return batchResp.Data.BlockIDs, nil
}

// CreateTextBlock 创建文本块
func (c *Client) CreateTextBlock(documentID, parentID, text string, style map[string]interface{}) (string, error) {
	// 使用完整的文本块结构
	textContent := map[string]interface{}{
		"elements": []interface{}{
			map[string]interface{}{
				"text_run": map[string]interface{}{
					"content": text,
					"text_element_style": map[string]interface{}{
						"bold":          false,
						"inline_code":   false,
						"italic":        false,
						"strikethrough": false,
						"underline":     false,
					},
				},
			},
		},
		"style": map[string]interface{}{
			"align":  1,
			"folded": false,
		},
	}

	// 如果提供了样式参数，合并到文本样式中
	if style != nil {
		if textStyle, ok := textContent["style"].(map[string]interface{}); ok {
			for k, v := range style {
				textStyle[k] = v
			}
		}
	}

	content := map[string]interface{}{
		"text": textContent,
	}

	return c.CreateBlock(documentID, "2", parentID, content)
}

// CreateCodeBlock 创建代码块
func (c *Client) CreateCodeBlock(documentID, parentID, code, language string) (string, error) {
	content := map[string]interface{}{
		"code": code,
	}
	if language != "" {
		content["language"] = language
	}

	return c.CreateBlock(documentID, "code", parentID, content)
}

// CreateHeadingBlock 创建标题块
func (c *Client) CreateHeadingBlock(documentID, parentID, text string, level int) (string, error) {
	content := map[string]interface{}{
		"text":  text,
		"level": level,
	}

	return c.CreateBlock(documentID, "heading", parentID, content)
}

// CreateListBlock 创建列表块
func (c *Client) CreateListBlock(documentID, parentID string, items []string, listType string) (string, error) {
	content := map[string]interface{}{
		"items": items,
		"type":  listType, // "ordered" 或 "unordered"
	}

	return c.CreateBlock(documentID, "list", parentID, content)
}

// DeleteBlocks 删除文档块
func (c *Client) DeleteBlocks(documentID string, startIndex, endIndex int) error {
	token, err := c.getAccessToken()
	if err != nil {
		return err
	}

	endpoint := fmt.Sprintf("/docx/v1/documents/%s/blocks/batch_delete", documentID)
	deleteReq := DeleteBlocksRequest{
		StartIndex: startIndex,
		EndIndex:   endIndex,
	}

	resp, err := c.client.R().
		SetHeader("Authorization", "Bearer "+token).
		SetHeader("Content-Type", "application/json").
		SetBody(deleteReq).
		Delete(c.baseURL + endpoint)

	if err != nil {
		return fmt.Errorf("删除块请求失败: %w", err)
	}

	if err := c.CheckResponse(resp); err != nil {
		return err
	}

	var deleteResp DeleteBlocksResponse
	if err := json.Unmarshal(resp.Body(), &deleteResp); err != nil {
		return fmt.Errorf("解析删除响应失败: %w", err)
	}

	if deleteResp.Code != 0 {
		return fmt.Errorf("删除块失败: %s", deleteResp.Msg)
	}

	return nil
}

// ConvertWiki API 相关结构
type ConvertWikiRequest struct {
	ObjToken string `json:"obj_token"`
	ObjType  string `json:"obj_type"`
}

type ConvertWikiResponse struct {
	Data *ConvertWikiData `json:"data"`
	Code int              `json:"code"`
	Msg  string           `json:"msg"`
}

type ConvertWikiData struct {
	DocumentID string `json:"document_id"`
	URL        string `json:"url"`
	Token      string `json:"token"`
}

// GetImageResource API 相关结构
type GetImageResourceRequest struct {
	ImageKey string `json:"image_key"`
}

type GetImageResourceResponse struct {
	Data *ImageResourceData `json:"data"`
	Code int                `json:"code"`
	Msg  string             `json:"msg"`
}

type ImageResourceData struct {
	ImageKey    string `json:"image_key"`
	URL         string `json:"url"`
	ContentType string `json:"content_type"`
	Size        int64  `json:"size"`
	Base64Data  string `json:"base64_data,omitempty"`
}

// ConvertWiki 将飞书Wiki链接转换为文档ID
func (c *Client) ConvertWiki(objToken string, objType string) (*ConvertWikiResponse, error) {
	// 使用wiki API端点
	endpoint := fmt.Sprintf("/wiki/v2/spaces/%s/nodes/%s", objToken, objToken)

	resp, err := c.Get(endpoint)
	if err != nil {
		return nil, fmt.Errorf("转换Wiki失败: %w", err)
	}

	var response ConvertWikiResponse
	if err := json.Unmarshal(resp.Body(), &response); err != nil {
		return nil, fmt.Errorf("解析转换Wiki响应失败: %w", err)
	}

	if response.Code != 0 {
		return nil, fmt.Errorf("转换Wiki失败: code=%d, msg=%s", response.Code, response.Msg)
	}

	return &response, nil
}

// GetImageResource 获取飞书图片资源
func (c *Client) GetImageResource(imageKey string) (*GetImageResourceResponse, error) {
	endpoint := fmt.Sprintf("/im/v1/images/%s", imageKey)

	resp, err := c.Get(endpoint)
	if err != nil {
		return nil, fmt.Errorf("获取图片资源失败: %w", err)
	}

	var response GetImageResourceResponse
	if err := json.Unmarshal(resp.Body(), &response); err != nil {
		return nil, fmt.Errorf("解析图片资源响应失败: %w", err)
	}

	if response.Code != 0 {
		return nil, fmt.Errorf("获取图片资源失败: code=%d, msg=%s", response.Code, response.Msg)
	}

	return &response, nil
}
