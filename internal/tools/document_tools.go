package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	"feishu-mcp-go/pkg/feishu"
	"feishu-mcp-go/pkg/logger"
)

// NewCreateDocumentTool 创建文档工具构造函数
func NewCreateDocumentTool(feishu *feishu.Client, logger logger.Logger) server.ServerTool {
	tool := mcp.NewTool("create_feishu_document",
		mcp.WithDescription("创建新的飞书文档"),
		mcp.WithString("title", mcp.Description("文档标题"), mcp.Required()),
		mcp.WithString("folderToken", mcp.Description("父文件夹的token，用于指定文档创建位置"), mcp.Required()))

	handler := func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return createDocument(feishu, logger, request)
	}
	return server.ServerTool{Tool: tool, Handler: handler}
}

func NewGetDocumentInfoTool(feishu *feishu.Client, logger logger.Logger) server.ServerTool {
	tool := mcp.NewTool("get_feishu_document_info",
		mcp.WithDescription("获取飞书文档的基本信息"),
		mcp.WithString("documentId", mcp.Description("文档ID，用于标识要获取信息的文档"), mcp.Required()))
	handler := func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return getDocumentInfo(feishu, logger, request)
	}
	return server.ServerTool{Tool: tool, Handler: handler}
}

func NewGetDocumentContentTool(feishu *feishu.Client, logger logger.Logger) server.ServerTool {
	tool := mcp.NewTool("get_feishu_document_content",
		mcp.WithDescription("获取飞书文档的纯文本内容"),
		mcp.WithString("documentId", mcp.Description("文档ID，用于标识要获取内容的文档"), mcp.Required()),
		mcp.WithNumber("lang", mcp.Description("语言类型，0表示中文，1表示英文，默认为0"), mcp.DefaultNumber(0)))
	handler := func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return getDocumentContent(feishu, logger, request)
	}
	return server.ServerTool{Tool: tool, Handler: handler}
}

func NewGetDocumentBlocksTool(feishu *feishu.Client, logger logger.Logger) server.ServerTool {
	tool := mcp.NewTool("get_feishu_document_blocks",
		mcp.WithDescription("获取飞书文档的块结构信息"),
		mcp.WithString("documentId", mcp.Description("文档ID，用于标识要获取块结构的文档"), mcp.Required()),
		mcp.WithNumber("pageSize", mcp.Description("每页返回的块数量，默认为50"), mcp.DefaultNumber(50)),
		mcp.WithString("pageToken", mcp.Description("分页令牌，用于获取下一页数据，可选")))
	handler := func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return getDocumentBlocks(feishu, logger, request)
	}
	return server.ServerTool{Tool: tool, Handler: handler}
}

func NewSearchDocumentsTool(feishu *feishu.Client, logger logger.Logger) server.ServerTool {
	tool := mcp.NewTool("search_feishu_documents",
		mcp.WithDescription("在飞书中搜索文档"),
		mcp.WithString("query", mcp.Description("搜索关键词，用于匹配文档标题或内容"), mcp.Required()),
		mcp.WithNumber("pageSize", mcp.Description("每页返回的文档数量，默认为10"), mcp.DefaultNumber(10)),
		mcp.WithString("pageToken", mcp.Description("分页令牌，用于获取下一页搜索结果，可选")))
	handler := func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return searchDocuments(feishu, logger, request)
	}
	return server.ServerTool{Tool: tool, Handler: handler}
}

// createDocument 创建文档实现
func createDocument(feishu *feishu.Client, logger logger.Logger, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	logger.Info("开始创建飞书文档", "name", request.Params.Name)

	// 解析参数
	var args map[string]interface{}

	// 安全地处理参数类型断言
	if request.Params.Arguments != nil {
		if argBytes, ok := request.Params.Arguments.([]byte); ok {
			if err := json.Unmarshal(argBytes, &args); err != nil {
				return nil, fmt.Errorf("解析参数失败: %w", err)
			}
		} else if argMap, ok := request.Params.Arguments.(map[string]interface{}); ok {
			args = argMap
		} else {
			// 尝试通过JSON序列化再反序列化
			if argBytes, err := json.Marshal(request.Params.Arguments); err != nil {
				return nil, fmt.Errorf("参数格式错误: %w", err)
			} else if err := json.Unmarshal(argBytes, &args); err != nil {
				return nil, fmt.Errorf("解析参数失败: %w", err)
			}
		}
	}

	title, ok := args["title"].(string)
	if !ok || title == "" {
		return nil, fmt.Errorf("title参数必须是非空字符串")
	}

	folderToken, ok := args["folderToken"].(string)
	if !ok || folderToken == "" {
		return nil, fmt.Errorf("folderToken参数必须是非空字符串")
	}

	// 调用飞书API创建文档
	document, err := feishu.CreateDocument(folderToken, title)
	if err != nil {
		logger.Error("创建文档失败", "error", err)
		return nil, fmt.Errorf("创建文档失败: %w", err)
	}

	logger.Info("文档创建成功", "documentId", document.DocumentID, "title", document.Title)

	// 构建成功响应
	result := map[string]interface{}{
		"success":    true,
		"message":    "文档创建成功",
		"documentId": document.DocumentID,
		"token":      document.Token,
		"title":      document.Title,
		"ownerId":    document.OwnerID,
		"createTime": document.CreateTime,
		"updateTime": document.UpdateTime,
	}

	content, _ := json.MarshalIndent(result, "", "  ")

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.TextContent{
				Type: "text",
				Text: string(content),
			},
		},
	}, nil
}

// getDocumentInfo 获取文档信息实现
func getDocumentInfo(feishu *feishu.Client, logger logger.Logger, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	logger.Info("开始获取文档信息", "name", request.Params.Name)

	// 解析参数
	var args map[string]interface{}

	// 安全地处理参数类型断言
	if request.Params.Arguments != nil {
		if argBytes, ok := request.Params.Arguments.([]byte); ok {
			if err := json.Unmarshal(argBytes, &args); err != nil {
				return nil, fmt.Errorf("解析参数失败: %w", err)
			}
		} else if argMap, ok := request.Params.Arguments.(map[string]interface{}); ok {
			args = argMap
		} else {
			// 尝试通过JSON序列化再反序列化
			if argBytes, err := json.Marshal(request.Params.Arguments); err != nil {
				return nil, fmt.Errorf("参数格式错误: %w", err)
			} else if err := json.Unmarshal(argBytes, &args); err != nil {
				return nil, fmt.Errorf("解析参数失败: %w", err)
			}
		}
	}

	documentId, ok := args["documentId"].(string)
	if !ok || documentId == "" {
		return nil, fmt.Errorf("documentId参数必须是非空字符串")
	}

	// 调用飞书API获取文档信息
	document, err := feishu.GetDocument(documentId)
	if err != nil {
		logger.Error("获取文档信息失败", "error", err)
		return nil, fmt.Errorf("获取文档信息失败: %w", err)
	}

	logger.Info("获取文档信息成功", "documentId", document.DocumentID)

	// 构建成功响应
	result := map[string]interface{}{
		"success":    true,
		"message":    "获取文档信息成功",
		"documentId": document.DocumentID,
		"token":      document.Token,
		"title":      document.Title,
		"ownerId":    document.OwnerID,
		"createTime": document.CreateTime,
		"updateTime": document.UpdateTime,
	}

	content, _ := json.MarshalIndent(result, "", "  ")

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.TextContent{
				Type: "text",
				Text: string(content),
			},
		},
	}, nil
}

// getDocumentContent 获取文档内容实现
func getDocumentContent(feishu *feishu.Client, logger logger.Logger, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	logger.Info("开始获取文档内容", "name", request.Params.Name)

	// 解析参数
	var args map[string]interface{}

	// 安全地处理参数类型断言
	if request.Params.Arguments != nil {
		if argBytes, ok := request.Params.Arguments.([]byte); ok {
			if err := json.Unmarshal(argBytes, &args); err != nil {
				return nil, fmt.Errorf("解析参数失败: %w", err)
			}
		} else if argMap, ok := request.Params.Arguments.(map[string]interface{}); ok {
			args = argMap
		} else {
			// 尝试通过JSON序列化再反序列化
			if argBytes, err := json.Marshal(request.Params.Arguments); err != nil {
				return nil, fmt.Errorf("参数格式错误: %w", err)
			} else if err := json.Unmarshal(argBytes, &args); err != nil {
				return nil, fmt.Errorf("解析参数失败: %w", err)
			}
		}
	}

	documentId, ok := args["documentId"].(string)
	if !ok || documentId == "" {
		return nil, fmt.Errorf("documentId参数必须是非空字符串")
	}

	// 获取语言参数，默认为0（中文）
	lang := 0
	if langVal, ok := args["lang"]; ok {
		if langFloat, ok := langVal.(float64); ok {
			lang = int(langFloat)
		}
	}

	// 调用飞书API获取文档内容
	content, err := feishu.GetDocumentContent(documentId, lang)
	if err != nil {
		logger.Error("获取文档内容失败", "error", err)
		return nil, fmt.Errorf("获取文档内容失败: %w", err)
	}

	logger.Info("获取文档内容成功", "documentId", documentId, "lang", lang)

	// 构建成功响应
	result := map[string]interface{}{
		"success":    true,
		"message":    "获取文档内容成功",
		"documentId": documentId,
		"lang":       lang,
		"content":    content.Content,
	}

	resultJson, _ := json.MarshalIndent(result, "", "  ")

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.TextContent{
				Type: "text",
				Text: string(resultJson),
			},
		},
	}, nil
}

// getDocumentBlocks 获取文档块结构实现
func getDocumentBlocks(feishu *feishu.Client, logger logger.Logger, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	logger.Info("开始获取文档块结构", "name", request.Params.Name)

	// 解析参数
	var args map[string]interface{}

	// 安全地处理参数类型断言
	if request.Params.Arguments != nil {
		if argBytes, ok := request.Params.Arguments.([]byte); ok {
			if err := json.Unmarshal(argBytes, &args); err != nil {
				return nil, fmt.Errorf("解析参数失败: %w", err)
			}
		} else if argMap, ok := request.Params.Arguments.(map[string]interface{}); ok {
			args = argMap
		} else {
			// 尝试通过JSON序列化再反序列化
			if argBytes, err := json.Marshal(request.Params.Arguments); err != nil {
				return nil, fmt.Errorf("参数格式错误: %w", err)
			} else if err := json.Unmarshal(argBytes, &args); err != nil {
				return nil, fmt.Errorf("解析参数失败: %w", err)
			}
		}
	}

	documentId, ok := args["documentId"].(string)
	if !ok || documentId == "" {
		return nil, fmt.Errorf("documentId参数必须是非空字符串")
	}

	// 调用飞书API获取文档块结构
	blocks, err := feishu.GetDocumentBlocks(documentId)
	if err != nil {
		logger.Error("获取文档块结构失败", "error", err)
		return nil, fmt.Errorf("获取文档块结构失败: %w", err)
	}

	logger.Info("获取文档块结构成功", "documentId", documentId, "blocksCount", len(blocks))

	// 构建成功响应
	result := map[string]interface{}{
		"success":     true,
		"message":     "获取文档块结构成功",
		"documentId":  documentId,
		"blocksCount": len(blocks),
		"blocks":      blocks,
	}

	resultJson, _ := json.MarshalIndent(result, "", "  ")

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.TextContent{
				Type: "text",
				Text: string(resultJson),
			},
		},
	}, nil
}

// searchDocuments 搜索文档实现
func searchDocuments(feishu *feishu.Client, logger logger.Logger, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	logger.Info("开始搜索文档", "name", request.Params.Name)

	// 解析参数
	var args map[string]interface{}

	// 安全地处理参数类型断言
	if request.Params.Arguments != nil {
		if argBytes, ok := request.Params.Arguments.([]byte); ok {
			if err := json.Unmarshal(argBytes, &args); err != nil {
				return nil, fmt.Errorf("解析参数失败: %w", err)
			}
		} else if argMap, ok := request.Params.Arguments.(map[string]interface{}); ok {
			args = argMap
		} else {
			// 尝试通过JSON序列化再反序列化
			if argBytes, err := json.Marshal(request.Params.Arguments); err != nil {
				return nil, fmt.Errorf("参数格式错误: %w", err)
			} else if err := json.Unmarshal(argBytes, &args); err != nil {
				return nil, fmt.Errorf("解析参数失败: %w", err)
			}
		}
	}

	searchKey, ok := args["searchKey"].(string)
	if !ok || searchKey == "" {
		return nil, fmt.Errorf("searchKey参数必须是非空字符串")
	}

	// 获取可选的分页参数
	pageSize := 10 // 默认分页大小
	if pageSizeVal, ok := args["pageSize"]; ok {
		if pageSizeFloat, ok := pageSizeVal.(float64); ok {
			pageSize = int(pageSizeFloat)
		}
	}

	pageToken := ""
	if pageTokenVal, ok := args["pageToken"]; ok {
		if pageTokenStr, ok := pageTokenVal.(string); ok {
			pageToken = pageTokenStr
		}
	}

	// 调用飞书API搜索文档
	searchResult, err := feishu.SearchDocuments(searchKey, pageSize, pageToken)
	if err != nil {
		logger.Error("搜索文档失败", "error", err)
		return nil, fmt.Errorf("搜索文档失败: %w", err)
	}

	logger.Info("搜索文档成功", "searchKey", searchKey, "documentsCount", len(searchResult.Data.Documents))

	// 构建成功响应
	result := map[string]interface{}{
		"success":        true,
		"message":        "搜索文档成功",
		"searchKey":      searchKey,
		"documentsCount": len(searchResult.Data.Documents),
		"documents":      searchResult.Data.Documents,
		"pageToken":      searchResult.Data.PageToken,
		"hasMore":        searchResult.Data.HasMore,
	}

	resultJson, _ := json.MarshalIndent(result, "", "  ")

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.TextContent{
				Type: "text",
				Text: string(resultJson),
			},
		},
	}, nil
}
