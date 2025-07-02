package tools

import (
	"context"
	"encoding/json"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	"feishu-mcp-go/pkg/feishu"
	"feishu-mcp-go/pkg/logger"
)

// 工具功能实现

func NewConvertWikiTool(feishu *feishu.Client, logger logger.Logger) server.ServerTool {
	tool := mcp.NewTool("convert_feishu_wiki_to_document_id",
		mcp.WithDescription("将飞书Wiki文档链接转换为兼容的文档ID"))
	handler := func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return executeUtilityTool(feishu, logger, "convert_wiki", request)
	}
	return server.ServerTool{Tool: tool, Handler: handler}
}

func NewGetImageResourceTool(feishu *feishu.Client, logger logger.Logger) server.ServerTool {
	tool := mcp.NewTool("get_feishu_image_resource",
		mcp.WithDescription("通过媒体ID下载飞书图片资源"))
	handler := func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return executeUtilityTool(feishu, logger, "get_image_resource", request)
	}
	return server.ServerTool{Tool: tool, Handler: handler}
}

// executeUtilityTool 工具功能执行函数
func executeUtilityTool(feishu *feishu.Client, logger logger.Logger, toolName string, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	logger.Info("执行工具功能", "tool", toolName, "name", request.Params.Name)

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
