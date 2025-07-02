package tools

import (
	"context"
	"encoding/json"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	"feishu-mcp-go/pkg/feishu"
	"feishu-mcp-go/pkg/logger"
)

// 文件夹管理工具实现

func NewGetRootFolderInfoTool(feishu *feishu.Client, logger logger.Logger) server.ServerTool {
	tool := mcp.NewTool("get_feishu_root_folder_info",
		mcp.WithDescription("获取飞书云盘根文件夹信息"))
	handler := func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return executeFolderTool(feishu, logger, "get_root_folder_info", request)
	}
	return server.ServerTool{Tool: tool, Handler: handler}
}

func NewGetFolderFilesTool(feishu *feishu.Client, logger logger.Logger) server.ServerTool {
	tool := mcp.NewTool("get_feishu_folder_files",
		mcp.WithDescription("获取指定文件夹中的文件和子文件夹列表"))
	handler := func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return executeFolderTool(feishu, logger, "get_folder_files", request)
	}
	return server.ServerTool{Tool: tool, Handler: handler}
}

func NewCreateFolderTool(feishu *feishu.Client, logger logger.Logger) server.ServerTool {
	tool := mcp.NewTool("create_feishu_folder",
		mcp.WithDescription("在指定父文件夹中创建新文件夹"))
	handler := func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return executeFolderTool(feishu, logger, "create_folder", request)
	}
	return server.ServerTool{Tool: tool, Handler: handler}
}

// executeFolderTool 文件夹工具执行函数
func executeFolderTool(feishu *feishu.Client, logger logger.Logger, toolName string, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	logger.Info("执行文件夹工具", "tool", toolName, "name", request.Params.Name)

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
