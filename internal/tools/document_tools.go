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
		mcp.WithDescription("创建新的飞书文档"))

	handler := func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return createDocument(feishu, logger, request)
	}
	return server.ServerTool{Tool: tool, Handler: handler}
}

func NewGetDocumentInfoTool(feishu *feishu.Client, logger logger.Logger) server.ServerTool {
	tool := mcp.NewTool("get_feishu_document_info",
		mcp.WithDescription("获取飞书文档的基本信息"))
	handler := func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return getDocumentInfo(feishu, logger, request)
	}
	return server.ServerTool{Tool: tool, Handler: handler}
}

func NewGetDocumentContentTool(feishu *feishu.Client, logger logger.Logger) server.ServerTool {
	tool := mcp.NewTool("get_feishu_document_content",
		mcp.WithDescription("获取飞书文档的纯文本内容"))
	handler := func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return getDocumentContent(feishu, logger, request)
	}
	return server.ServerTool{Tool: tool, Handler: handler}
}

func NewGetDocumentBlocksTool(feishu *feishu.Client, logger logger.Logger) server.ServerTool {
	tool := mcp.NewTool("get_feishu_document_blocks",
		mcp.WithDescription("获取飞书文档的块结构信息"))
	handler := func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return getDocumentBlocks(feishu, logger, request)
	}
	return server.ServerTool{Tool: tool, Handler: handler}
}

func NewSearchDocumentsTool(feishu *feishu.Client, logger logger.Logger) server.ServerTool {
	tool := mcp.NewTool("search_feishu_documents",
		mcp.WithDescription("在飞书中搜索文档"))
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

// getDocumentContent 获取文档内容实现（占位符）
func getDocumentContent(feishu *feishu.Client, logger logger.Logger, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	logger.Info("获取文档内容功能待实现", "name", request.Params.Name)

	result := map[string]interface{}{
		"success": false,
		"message": "get_feishu_document_content 功能待实现",
		"status":  "TODO",
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

// getDocumentBlocks 获取文档块结构实现（占位符）
func getDocumentBlocks(feishu *feishu.Client, logger logger.Logger, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	logger.Info("获取文档块结构功能待实现", "name", request.Params.Name)

	result := map[string]interface{}{
		"success": false,
		"message": "get_feishu_document_blocks 功能待实现",
		"status":  "TODO",
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

// searchDocuments 搜索文档实现（占位符）
func searchDocuments(feishu *feishu.Client, logger logger.Logger, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	logger.Info("搜索文档功能待实现", "name", request.Params.Name)

	result := map[string]interface{}{
		"success": false,
		"message": "search_feishu_documents 功能待实现",
		"status":  "TODO",
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
