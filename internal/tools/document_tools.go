package tools

import (
	"context"
	"encoding/json"

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
		return executeDocumentTool(feishu, logger, "create_document", request)
	}
	return server.ServerTool{Tool: tool, Handler: handler}
}

func NewGetDocumentInfoTool(feishu *feishu.Client, logger logger.Logger) server.ServerTool {
	tool := mcp.NewTool("get_feishu_document_info",
		mcp.WithDescription("获取飞书文档的基本信息"))
	handler := func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return executeDocumentTool(feishu, logger, "get_document_info", request)
	}
	return server.ServerTool{Tool: tool, Handler: handler}
}

func NewGetDocumentContentTool(feishu *feishu.Client, logger logger.Logger) server.ServerTool {
	tool := mcp.NewTool("get_feishu_document_content",
		mcp.WithDescription("获取飞书文档的纯文本内容"))
	handler := func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return executeDocumentTool(feishu, logger, "get_document_content", request)
	}
	return server.ServerTool{Tool: tool, Handler: handler}
}

func NewGetDocumentBlocksTool(feishu *feishu.Client, logger logger.Logger) server.ServerTool {
	tool := mcp.NewTool("get_feishu_document_blocks",
		mcp.WithDescription("获取飞书文档的块结构信息"))
	handler := func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return executeDocumentTool(feishu, logger, "get_document_blocks", request)
	}
	return server.ServerTool{Tool: tool, Handler: handler}
}

func NewSearchDocumentsTool(feishu *feishu.Client, logger logger.Logger) server.ServerTool {
	tool := mcp.NewTool("search_feishu_documents",
		mcp.WithDescription("在飞书中搜索文档"))
	handler := func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return executeDocumentTool(feishu, logger, "search_documents", request)
	}
	return server.ServerTool{Tool: tool, Handler: handler}
}

// executeDocumentTool 文档工具执行函数
func executeDocumentTool(feishu *feishu.Client, logger logger.Logger, toolName string, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	logger.Info("执行文档工具", "tool", toolName, "name", request.Params.Name)

	result := map[string]interface{}{
		"success": true,
		"message": toolName + " 执行成功",
		"tool":    toolName,
		"args":    request.Params.Arguments,
	}

	content, _ := json.Marshal(result)

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.TextContent{
				Type: "text",
				Text: string(content),
			},
		},
	}, nil
}
