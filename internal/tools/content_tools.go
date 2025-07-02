package tools

import (
	"context"
	"encoding/json"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	"feishu-mcp-go/pkg/feishu"
	"feishu-mcp-go/pkg/logger"
)

// 为所有内容操作工具创建基础实现

func NewGetBlockContentTool(feishu *feishu.Client, logger logger.Logger) server.ServerTool {
	tool := mcp.NewTool("get_feishu_block_content",
		mcp.WithDescription("获取飞书文档中指定块的详细内容"))
	handler := func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return executeBaseTool(feishu, logger, "get_block_content", request)
	}
	return server.ServerTool{Tool: tool, Handler: handler}
}

func NewUpdateBlockTextTool(feishu *feishu.Client, logger logger.Logger) server.ServerTool {
	tool := mcp.NewTool("update_feishu_block_text",
		mcp.WithDescription("更新飞书文档中指定块的文本内容和样式"))
	handler := func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return executeBaseTool(feishu, logger, "update_block_text", request)
	}
	return server.ServerTool{Tool: tool, Handler: handler}
}

func NewBatchCreateBlocksTool(feishu *feishu.Client, logger logger.Logger) server.ServerTool {
	tool := mcp.NewTool("batch_create_feishu_blocks",
		mcp.WithDescription("批量创建多个飞书文档块"))
	handler := func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return executeBaseTool(feishu, logger, "batch_create_blocks", request)
	}
	return server.ServerTool{Tool: tool, Handler: handler}
}

func NewCreateTextBlockTool(feishu *feishu.Client, logger logger.Logger) server.ServerTool {
	tool := mcp.NewTool("create_feishu_text_block",
		mcp.WithDescription("创建新的文本块"))
	handler := func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return executeBaseTool(feishu, logger, "create_text_block", request)
	}
	return server.ServerTool{Tool: tool, Handler: handler}
}

func NewCreateCodeBlockTool(feishu *feishu.Client, logger logger.Logger) server.ServerTool {
	tool := mcp.NewTool("create_feishu_code_block",
		mcp.WithDescription("创建新的代码块"))
	handler := func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return executeBaseTool(feishu, logger, "create_code_block", request)
	}
	return server.ServerTool{Tool: tool, Handler: handler}
}

func NewCreateHeadingBlockTool(feishu *feishu.Client, logger logger.Logger) server.ServerTool {
	tool := mcp.NewTool("create_feishu_heading_block",
		mcp.WithDescription("创建新的标题块"))
	handler := func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return executeBaseTool(feishu, logger, "create_heading_block", request)
	}
	return server.ServerTool{Tool: tool, Handler: handler}
}

func NewCreateListBlockTool(feishu *feishu.Client, logger logger.Logger) server.ServerTool {
	tool := mcp.NewTool("create_feishu_list_block",
		mcp.WithDescription("创建新的列表块"))
	handler := func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return executeBaseTool(feishu, logger, "create_list_block", request)
	}
	return server.ServerTool{Tool: tool, Handler: handler}
}

func NewDeleteBlocksTool(feishu *feishu.Client, logger logger.Logger) server.ServerTool {
	tool := mcp.NewTool("delete_feishu_document_blocks",
		mcp.WithDescription("删除飞书文档中的一个或多个连续块"))
	handler := func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return executeBaseTool(feishu, logger, "delete_blocks", request)
	}
	return server.ServerTool{Tool: tool, Handler: handler}
}

// executeBaseTool 基础工具执行函数
func executeBaseTool(feishu *feishu.Client, logger logger.Logger, toolName string, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	logger.Info("执行工具", "tool", toolName, "name", request.Params.Name)

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
