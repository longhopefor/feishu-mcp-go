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

// 为所有内容操作工具创建基础实现

func NewGetBlockContentTool(feishu *feishu.Client, logger logger.Logger) server.ServerTool {
	tool := mcp.NewTool("get_feishu_block_content",
		mcp.WithDescription("获取飞书文档中指定块的详细内容"),
		mcp.WithString("documentId", mcp.Description("文档ID，用于标识包含目标块的文档"), mcp.Required()),
		mcp.WithString("blockId", mcp.Description("块ID，用于标识要获取内容的具体块"), mcp.Required()))
	handler := func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return getBlockContent(feishu, logger, request)
	}
	return server.ServerTool{Tool: tool, Handler: handler}
}

func NewUpdateBlockTextTool(feishu *feishu.Client, logger logger.Logger) server.ServerTool {
	tool := mcp.NewTool("update_feishu_block_text",
		mcp.WithDescription("更新飞书文档中指定块的文本内容和样式"),
		mcp.WithString("documentId", mcp.Description("文档ID，用于标识包含目标块的文档"), mcp.Required()),
		mcp.WithString("blockId", mcp.Description("块ID，用于标识要更新的具体块"), mcp.Required()),
		mcp.WithString("content", mcp.Description("要更新的文本内容"), mcp.Required()),
		mcp.WithString("style", mcp.Description("文本样式，可选，如粗体、斜体等")))
	handler := func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return updateBlockText(feishu, logger, request)
	}
	return server.ServerTool{Tool: tool, Handler: handler}
}

func NewBatchCreateBlocksTool(feishu *feishu.Client, logger logger.Logger) server.ServerTool {
	tool := mcp.NewTool("batch_create_feishu_blocks",
		mcp.WithDescription("批量创建多个飞书文档块"),
		mcp.WithString("documentId", mcp.Description("文档ID，用于标识要添加块的文档"), mcp.Required()),
		mcp.WithString("parentId", mcp.Description("父块ID，用于指定新块的位置，可选")),
		mcp.WithString("blocks", mcp.Description("要创建的块列表，JSON格式字符串"), mcp.Required()))
	handler := func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return batchCreateBlocks(feishu, logger, request)
	}
	return server.ServerTool{Tool: tool, Handler: handler}
}

func NewCreateTextBlockTool(feishu *feishu.Client, logger logger.Logger) server.ServerTool {
	tool := mcp.NewTool("create_feishu_text_block",
		mcp.WithDescription("创建新的文本块"),
		mcp.WithString("documentId", mcp.Description("文档ID，用于标识要添加文本块的文档"), mcp.Required()),
		mcp.WithString("content", mcp.Description("文本块的内容"), mcp.Required()),
		mcp.WithString("parentId", mcp.Description("父块ID，用于指定文本块的位置，可选")),
		mcp.WithString("style", mcp.Description("文本样式，可选，如粗体、斜体等")))
	handler := func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return createTextBlock(feishu, logger, request)
	}
	return server.ServerTool{Tool: tool, Handler: handler}
}

func NewCreateCodeBlockTool(feishu *feishu.Client, logger logger.Logger) server.ServerTool {
	tool := mcp.NewTool("create_feishu_code_block",
		mcp.WithDescription("创建新的代码块"),
		mcp.WithString("documentId", mcp.Description("文档ID，用于标识要添加代码块的文档"), mcp.Required()),
		mcp.WithString("content", mcp.Description("代码块的内容"), mcp.Required()),
		mcp.WithString("language", mcp.Description("编程语言类型，如javascript、python等"), mcp.Required()),
		mcp.WithString("parentId", mcp.Description("父块ID，用于指定代码块的位置，可选")))
	handler := func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return createCodeBlock(feishu, logger, request)
	}
	return server.ServerTool{Tool: tool, Handler: handler}
}

func NewCreateHeadingBlockTool(feishu *feishu.Client, logger logger.Logger) server.ServerTool {
	tool := mcp.NewTool("create_feishu_heading_block",
		mcp.WithDescription("创建新的标题块"),
		mcp.WithString("documentId", mcp.Description("文档ID，用于标识要添加标题块的文档"), mcp.Required()),
		mcp.WithString("content", mcp.Description("标题内容"), mcp.Required()),
		mcp.WithNumber("level", mcp.Description("标题级别，1-3，1为最高级别"), mcp.Required()),
		mcp.WithString("parentId", mcp.Description("父块ID，用于指定标题块的位置，可选")))
	handler := func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return createHeadingBlock(feishu, logger, request)
	}
	return server.ServerTool{Tool: tool, Handler: handler}
}

func NewCreateListBlockTool(feishu *feishu.Client, logger logger.Logger) server.ServerTool {
	tool := mcp.NewTool("create_feishu_list_block",
		mcp.WithDescription("创建新的列表块"),
		mcp.WithString("documentId", mcp.Description("文档ID，用于标识要添加列表块的文档"), mcp.Required()),
		mcp.WithString("content", mcp.Description("列表项内容"), mcp.Required()),
		mcp.WithString("type", mcp.Description("列表类型，bullet表示无序列表，ordered表示有序列表"), mcp.Required()),
		mcp.WithString("parentId", mcp.Description("父块ID，用于指定列表块的位置，可选")))
	handler := func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return createListBlock(feishu, logger, request)
	}
	return server.ServerTool{Tool: tool, Handler: handler}
}

func NewDeleteBlocksTool(feishu *feishu.Client, logger logger.Logger) server.ServerTool {
	tool := mcp.NewTool("delete_feishu_document_blocks",
		mcp.WithDescription("删除飞书文档中的一个或多个连续块"),
		mcp.WithString("documentId", mcp.Description("文档ID，用于标识包含要删除块的文档"), mcp.Required()),
		mcp.WithString("startIndex", mcp.Description("删除起始块的索引位置"), mcp.Required()),
		mcp.WithString("endIndex", mcp.Description("删除结束块的索引位置，包含此位置"), mcp.Required()))
	handler := func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return deleteBlocks(feishu, logger, request)
	}
	return server.ServerTool{Tool: tool, Handler: handler}
}

// getBlockContent 获取块内容实现
func getBlockContent(feishu *feishu.Client, logger logger.Logger, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	logger.Info("开始获取块内容", "name", request.Params.Name)

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

	blockId, ok := args["blockId"].(string)
	if !ok || blockId == "" {
		return nil, fmt.Errorf("blockId参数必须是非空字符串")
	}

	// 调用飞书API获取块内容
	blockContent, err := feishu.GetBlockContent(documentId, blockId)
	if err != nil {
		logger.Error("获取块内容失败", "error", err)
		return nil, fmt.Errorf("获取块内容失败: %w", err)
	}

	logger.Info("获取块内容成功", "blockId", blockContent.BlockID)

	// 根据块类型构建内容
	var content interface{}
	switch blockContent.BlockType {
	case 1: // 页面块
		content = blockContent.Page
	case 2: // 文本块
		content = blockContent.Text
	case 3: // 代码块
		content = blockContent.Code
	case 4: // 标题块
		content = blockContent.Heading
	case 5: // 列表块
		content = blockContent.List
	default:
		content = map[string]interface{}{
			"blockType": blockContent.BlockType,
			"note":      "未知块类型",
		}
	}

	// 构建成功响应
	result := map[string]interface{}{
		"success":   true,
		"message":   "获取块内容成功",
		"blockId":   blockContent.BlockID,
		"blockType": blockContent.BlockType,
		"parentId":  blockContent.ParentID,
		"content":   content,
	}

	resultContent, _ := json.MarshalIndent(result, "", "  ")

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.TextContent{
				Type: "text",
				Text: string(resultContent),
			},
		},
	}, nil
}

// updateBlockText 更新块文本实现
func updateBlockText(feishu *feishu.Client, logger logger.Logger, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	logger.Info("开始更新块文本", "name", request.Params.Name)

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

	blockId, ok := args["blockId"].(string)
	if !ok || blockId == "" {
		return nil, fmt.Errorf("blockId参数必须是非空字符串")
	}

	// 获取内容参数，可以是字符串
	var content map[string]interface{}
	if contentStr, ok := args["content"].(string); ok {
		// 使用正确的API格式 - 使用text字段传递字符串内容
		content = map[string]interface{}{
			"text": contentStr,
		}
	} else if contentMap, ok := args["content"].(map[string]interface{}); ok {
		content = contentMap
	} else {
		return nil, fmt.Errorf("content参数必须是字符串或对象")
	}

	// 调用飞书API更新块文本
	err := feishu.UpdateBlockText(documentId, blockId, content)
	if err != nil {
		logger.Error("更新块文本失败", "error", err)
		return nil, fmt.Errorf("更新块文本失败: %w", err)
	}

	logger.Info("更新块文本成功", "blockId", blockId)

	// 构建成功响应
	result := map[string]interface{}{
		"success":    true,
		"message":    "更新块文本成功",
		"blockId":    blockId,
		"documentId": documentId,
	}

	resultContent, _ := json.MarshalIndent(result, "", "  ")

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.TextContent{
				Type: "text",
				Text: string(resultContent),
			},
		},
	}, nil
}

// createTextBlock 创建文本块实现
func createTextBlock(feishu *feishu.Client, logger logger.Logger, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	logger.Info("开始创建文本块", "name", request.Params.Name)

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

	text, ok := args["text"].(string)
	if !ok || text == "" {
		return nil, fmt.Errorf("text参数必须是非空字符串")
	}

	// 获取可选的父块ID
	parentId, _ := args["parentId"].(string)

	// 获取可选的样式参数
	var style map[string]interface{}
	if styleArg, ok := args["style"].(map[string]interface{}); ok {
		style = styleArg
	}

	// 调用飞书API创建文本块
	blockId, err := feishu.CreateTextBlock(documentId, parentId, text, style)
	if err != nil {
		logger.Error("创建文本块失败", "error", err)
		return nil, fmt.Errorf("创建文本块失败: %w", err)
	}

	logger.Info("创建文本块成功", "blockId", blockId)

	// 构建成功响应
	result := map[string]interface{}{
		"success":    true,
		"message":    "创建文本块成功",
		"blockId":    blockId,
		"documentId": documentId,
		"text":       text,
	}

	resultContent, _ := json.MarshalIndent(result, "", "  ")

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.TextContent{
				Type: "text",
				Text: string(resultContent),
			},
		},
	}, nil
}

// createCodeBlock 创建代码块实现
func createCodeBlock(feishu *feishu.Client, logger logger.Logger, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	logger.Info("开始创建代码块", "name", request.Params.Name)

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

	code, ok := args["code"].(string)
	if !ok || code == "" {
		return nil, fmt.Errorf("code参数必须是非空字符串")
	}

	// 获取可选的父块ID
	parentId, _ := args["parentId"].(string)

	// 获取可选的编程语言参数，默认为 "text"
	language, _ := args["language"].(string)
	if language == "" {
		language = "text"
	}

	// 调用飞书API创建代码块
	blockId, err := feishu.CreateCodeBlock(documentId, parentId, code, language)
	if err != nil {
		logger.Error("创建代码块失败", "error", err)
		return nil, fmt.Errorf("创建代码块失败: %w", err)
	}

	logger.Info("创建代码块成功", "blockId", blockId)

	// 构建成功响应
	result := map[string]interface{}{
		"success":    true,
		"message":    "创建代码块成功",
		"blockId":    blockId,
		"documentId": documentId,
		"code":       code,
		"language":   language,
	}

	resultContent, _ := json.MarshalIndent(result, "", "  ")

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.TextContent{
				Type: "text",
				Text: string(resultContent),
			},
		},
	}, nil
}

// createHeadingBlock 创建标题块实现
func createHeadingBlock(feishu *feishu.Client, logger logger.Logger, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	logger.Info("开始创建标题块", "name", request.Params.Name)

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

	text, ok := args["text"].(string)
	if !ok || text == "" {
		return nil, fmt.Errorf("text参数必须是非空字符串")
	}

	// 获取可选的父块ID
	parentId, _ := args["parentId"].(string)

	// 获取标题级别参数，默认为 1
	level := 1
	if levelArg, ok := args["level"].(float64); ok {
		level = int(levelArg)
	} else if levelArg, ok := args["level"].(int); ok {
		level = levelArg
	}

	// 验证标题级别范围
	if level < 1 || level > 6 {
		return nil, fmt.Errorf("level参数必须在1-6范围内")
	}

	// 调用飞书API创建标题块
	blockId, err := feishu.CreateHeadingBlock(documentId, parentId, text, level)
	if err != nil {
		logger.Error("创建标题块失败", "error", err)
		return nil, fmt.Errorf("创建标题块失败: %w", err)
	}

	logger.Info("创建标题块成功", "blockId", blockId)

	// 构建成功响应
	result := map[string]interface{}{
		"success":    true,
		"message":    "创建标题块成功",
		"blockId":    blockId,
		"documentId": documentId,
		"text":       text,
		"level":      level,
	}

	resultContent, _ := json.MarshalIndent(result, "", "  ")

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.TextContent{
				Type: "text",
				Text: string(resultContent),
			},
		},
	}, nil
}

// createListBlock 创建列表块实现
func createListBlock(feishu *feishu.Client, logger logger.Logger, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	logger.Info("开始创建列表块", "name", request.Params.Name)

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

	// 获取可选的父块ID
	parentId, _ := args["parentId"].(string)

	// 获取列表项参数
	var items []string
	if itemsArg, ok := args["items"].([]interface{}); ok {
		for _, item := range itemsArg {
			if itemStr, ok := item.(string); ok {
				items = append(items, itemStr)
			}
		}
	} else if itemsArg, ok := args["items"].([]string); ok {
		items = itemsArg
	}

	if len(items) == 0 {
		return nil, fmt.Errorf("items参数必须是非空字符串数组")
	}

	// 获取列表类型参数，默认为 "bullet"
	listType, _ := args["listType"].(string)
	if listType == "" {
		listType = "bullet"
	}

	// 验证列表类型
	if listType != "bullet" && listType != "numbered" {
		return nil, fmt.Errorf("listType参数必须是 'bullet' 或 'numbered'")
	}

	// 调用飞书API创建列表块
	blockId, err := feishu.CreateListBlock(documentId, parentId, items, listType)
	if err != nil {
		logger.Error("创建列表块失败", "error", err)
		return nil, fmt.Errorf("创建列表块失败: %w", err)
	}

	logger.Info("创建列表块成功", "blockId", blockId)

	// 构建成功响应
	result := map[string]interface{}{
		"success":    true,
		"message":    "创建列表块成功",
		"blockId":    blockId,
		"documentId": documentId,
		"items":      items,
		"listType":   listType,
	}

	resultContent, _ := json.MarshalIndent(result, "", "  ")

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.TextContent{
				Type: "text",
				Text: string(resultContent),
			},
		},
	}, nil
}

// deleteBlocks 删除文档块实现
func deleteBlocks(feishu *feishu.Client, logger logger.Logger, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	logger.Info("开始删除文档块", "name", request.Params.Name)

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

	// 获取起始索引
	startIndex := 0
	if startIndexArg, ok := args["startIndex"].(float64); ok {
		startIndex = int(startIndexArg)
	} else if startIndexArg, ok := args["startIndex"].(int); ok {
		startIndex = startIndexArg
	}

	// 获取结束索引，如果未提供，则默认为startIndex
	endIndex := startIndex
	if endIndexArg, ok := args["endIndex"].(float64); ok {
		endIndex = int(endIndexArg)
	} else if endIndexArg, ok := args["endIndex"].(int); ok {
		endIndex = endIndexArg
	}

	// 验证索引范围
	if startIndex < 0 {
		return nil, fmt.Errorf("startIndex参数不能为负数")
	}
	if endIndex < startIndex {
		return nil, fmt.Errorf("endIndex参数不能小于startIndex")
	}

	// 调用飞书API删除块
	err := feishu.DeleteBlocks(documentId, startIndex, endIndex)
	if err != nil {
		logger.Error("删除文档块失败", "error", err)
		return nil, fmt.Errorf("删除文档块失败: %w", err)
	}

	blocksCount := endIndex - startIndex + 1
	logger.Info("删除文档块成功", "blocksCount", blocksCount)

	// 构建成功响应
	result := map[string]interface{}{
		"success":     true,
		"message":     "删除文档块成功",
		"documentId":  documentId,
		"startIndex":  startIndex,
		"endIndex":    endIndex,
		"blocksCount": blocksCount,
	}

	resultContent, _ := json.MarshalIndent(result, "", "  ")

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.TextContent{
				Type: "text",
				Text: string(resultContent),
			},
		},
	}, nil
}

// batchCreateBlocks 批量创建块实现
func batchCreateBlocks(feishu *feishu.Client, logger logger.Logger, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	logger.Info("开始批量创建块", "name", request.Params.Name)

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

	// 获取块列表参数
	blocksArg, ok := args["blocks"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("blocks参数必须是数组")
	}

	if len(blocksArg) == 0 {
		return nil, fmt.Errorf("blocks参数不能为空")
	}

	// 由于无法直接使用feishu.CreateBlockRequest类型，我们将使用多个单独的CreateBlock调用
	var blockIds []string
	for i, blockArg := range blocksArg {
		blockMap, ok := blockArg.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("块 %d 必须是对象", i)
		}

		blockType, ok := blockMap["blockType"].(string)
		if !ok || blockType == "" {
			return nil, fmt.Errorf("块 %d 的blockType参数必须是非空字符串", i)
		}

		parentId, _ := blockMap["parentId"].(string)

		// 获取内容参数
		var content map[string]interface{}
		if contentArg, ok := blockMap["content"].(map[string]interface{}); ok {
			content = contentArg
		} else {
			return nil, fmt.Errorf("块 %d 的content参数必须是对象", i)
		}

		// 调用CreateBlock方法来创建单个块
		blockId, err := feishu.CreateBlock(documentId, blockType, parentId, content)
		if err != nil {
			logger.Error("创建块失败", "blockIndex", i, "error", err)
			return nil, fmt.Errorf("创建块 %d 失败: %w", i, err)
		}

		blockIds = append(blockIds, blockId)
	}

	logger.Info("批量创建块成功", "blockCount", len(blockIds))

	// 构建成功响应
	result := map[string]interface{}{
		"success":    true,
		"message":    "批量创建块成功",
		"documentId": documentId,
		"blockIds":   blockIds,
		"blockCount": len(blockIds),
	}

	resultContent, _ := json.MarshalIndent(result, "", "  ")

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.TextContent{
				Type: "text",
				Text: string(resultContent),
			},
		},
	}, nil
}

// executeBaseTool 基础工具执行函数 (占位符，将被具体实现替换)
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
