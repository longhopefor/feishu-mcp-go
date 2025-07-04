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

// ==================== 知识库工具构造函数 ====================

// NewGetWikiSpacesTool 获取知识库空间列表工具构造函数
func NewGetWikiSpacesTool(feishu *feishu.Client, logger logger.Logger) server.ServerTool {
	tool := mcp.NewTool("get_feishu_wiki_spaces",
		mcp.WithDescription("获取飞书知识库空间列表"))

	handler := func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return getWikiSpaces(feishu, logger, request)
	}
	return server.ServerTool{Tool: tool, Handler: handler}
}

// NewGetWikiNodesTool 获取知识库节点列表工具构造函数
func NewGetWikiNodesTool(feishu *feishu.Client, logger logger.Logger) server.ServerTool {
	tool := mcp.NewTool("get_feishu_wiki_nodes",
		mcp.WithDescription("获取飞书知识库空间下的节点列表"))

	handler := func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return getWikiNodes(feishu, logger, request)
	}
	return server.ServerTool{Tool: tool, Handler: handler}
}

// NewGetWikiNodeContentTool 获取知识库节点内容工具构造函数
func NewGetWikiNodeContentTool(feishu *feishu.Client, logger logger.Logger) server.ServerTool {
	tool := mcp.NewTool("get_feishu_wiki_node_content",
		mcp.WithDescription("获取飞书知识库节点的内容"))

	handler := func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return getWikiNodeContent(feishu, logger, request)
	}
	return server.ServerTool{Tool: tool, Handler: handler}
}

// NewGetWikiNodeMetaTool 获取知识库节点元信息工具构造函数
func NewGetWikiNodeMetaTool(feishu *feishu.Client, logger logger.Logger) server.ServerTool {
	tool := mcp.NewTool("get_feishu_wiki_node_meta",
		mcp.WithDescription("获取飞书知识库节点的元信息"))

	handler := func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return getWikiNodeMeta(feishu, logger, request)
	}
	return server.ServerTool{Tool: tool, Handler: handler}
}

// ==================== 知识库工具实现 ====================

// getWikiSpaces 获取知识库空间列表实现
func getWikiSpaces(feishu *feishu.Client, logger logger.Logger, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	logger.Info("开始获取知识库空间列表", "name", request.Params.Name)

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

	// 解析可选参数
	pageSize := 20 // 默认页大小
	if ps, ok := args["pageSize"]; ok {
		if psFloat, ok := ps.(float64); ok {
			pageSize = int(psFloat)
		}
	}

	pageToken := ""
	if pt, ok := args["pageToken"]; ok {
		if ptStr, ok := pt.(string); ok {
			pageToken = ptStr
		}
	}

	// 调用飞书API获取知识库空间列表
	response, err := feishu.GetWikiSpaces(pageSize, pageToken)
	if err != nil {
		logger.Error("获取知识库空间列表失败", "error", err)
		return nil, fmt.Errorf("获取知识库空间列表失败: %w", err)
	}

	logger.Info("获取知识库空间列表成功", "count", len(response.Data.Items))

	// 构建成功响应
	result := map[string]interface{}{
		"success":   true,
		"message":   "获取知识库空间列表成功",
		"spaces":    response.Data.Items, // 修正：使用Items字段
		"pageToken": response.Data.PageToken,
		"hasMore":   response.Data.HasMore,
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

// getWikiNodes 获取知识库节点列表实现
func getWikiNodes(feishu *feishu.Client, logger logger.Logger, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	logger.Info("开始获取知识库节点列表", "name", request.Params.Name)

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

	// 必需参数：空间ID
	spaceId, ok := args["spaceId"].(string)
	if !ok || spaceId == "" {
		return nil, fmt.Errorf("spaceId参数必须是非空字符串")
	}

	// 解析可选参数
	pageSize := 20 // 默认页大小
	if ps, ok := args["pageSize"]; ok {
		if psFloat, ok := ps.(float64); ok {
			pageSize = int(psFloat)
		}
	}

	pageToken := ""
	if pt, ok := args["pageToken"]; ok {
		if ptStr, ok := pt.(string); ok {
			pageToken = ptStr
		}
	}

	parentNodeToken := ""
	if pnt, ok := args["parentNodeToken"]; ok {
		if pntStr, ok := pnt.(string); ok {
			parentNodeToken = pntStr
		}
	}

	// 调用飞书API获取知识库节点列表
	response, err := feishu.GetWikiSpaceNodes(spaceId, pageSize, pageToken, parentNodeToken)
	if err != nil {
		logger.Error("获取知识库节点列表失败", "error", err)
		return nil, fmt.Errorf("获取知识库节点列表失败: %w", err)
	}

	logger.Info("获取知识库节点列表成功", "spaceId", spaceId, "count", len(response.Data.Items))

	// 构建成功响应
	result := map[string]interface{}{
		"success":   true,
		"message":   "获取知识库节点列表成功",
		"spaceId":   spaceId,
		"nodes":     response.Data.Items,
		"pageToken": response.Data.PageToken,
		"hasMore":   response.Data.HasMore,
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

// getWikiNodeContent 获取知识库节点内容实现
func getWikiNodeContent(feishu *feishu.Client, logger logger.Logger, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	logger.Info("开始获取知识库节点内容", "name", request.Params.Name)

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

	// 必需参数：空间ID和节点Token
	spaceId, ok := args["spaceId"].(string)
	if !ok || spaceId == "" {
		return nil, fmt.Errorf("spaceId参数必须是非空字符串")
	}

	nodeToken, ok := args["nodeToken"].(string)
	if !ok || nodeToken == "" {
		return nil, fmt.Errorf("nodeToken参数必须是非空字符串")
	}

	// 解析可选参数
	lang := 0 // 默认语言 0: 中文
	if l, ok := args["lang"]; ok {
		if lFloat, ok := l.(float64); ok {
			lang = int(lFloat)
		}
	}

	// 调用飞书API获取知识库节点内容
	content, err := feishu.GetWikiNodeContent(spaceId, nodeToken, lang)
	if err != nil {
		logger.Error("获取知识库节点内容失败", "error", err)
		return nil, fmt.Errorf("获取知识库节点内容失败: %w", err)
	}

	logger.Info("获取知识库节点内容成功", "spaceId", spaceId, "nodeToken", nodeToken)

	// 构建成功响应
	result := map[string]interface{}{
		"success":   true,
		"message":   "获取知识库节点内容成功",
		"spaceId":   spaceId,
		"nodeToken": nodeToken,
		"content":   content.Content,
		"type":      content.Type,
	}

	contentBytes, _ := json.MarshalIndent(result, "", "  ")

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.TextContent{
				Type: "text",
				Text: string(contentBytes),
			},
		},
	}, nil
}

// getWikiNodeMeta 获取知识库节点元信息实现
func getWikiNodeMeta(feishu *feishu.Client, logger logger.Logger, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	logger.Info("开始获取知识库节点元信息", "name", request.Params.Name)

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

	// 必需参数：空间ID和节点Token
	spaceId, ok := args["spaceId"].(string)
	if !ok || spaceId == "" {
		return nil, fmt.Errorf("spaceId参数必须是非空字符串")
	}

	nodeToken, ok := args["nodeToken"].(string)
	if !ok || nodeToken == "" {
		return nil, fmt.Errorf("nodeToken参数必须是非空字符串")
	}

	// 调用飞书API获取知识库节点元信息
	meta, err := feishu.GetWikiNodeMeta(spaceId, nodeToken)
	if err != nil {
		logger.Error("获取知识库节点元信息失败", "error", err)
		return nil, fmt.Errorf("获取知识库节点元信息失败: %w", err)
	}

	logger.Info("获取知识库节点元信息成功", "spaceId", spaceId, "nodeToken", nodeToken, "title", meta.Title)

	// 构建成功响应
	result := map[string]interface{}{
		"success":        true,
		"message":        "获取知识库节点元信息成功",
		"spaceId":        meta.SpaceID,
		"nodeToken":      meta.NodeToken,
		"objToken":       meta.ObjToken,
		"objType":        meta.ObjType,
		"parentNode":     meta.ParentNode,
		"nodeType":       meta.NodeType,
		"originNode":     meta.OriginNode,
		"originSpaceId":  meta.OriginSpaceID,
		"title":          meta.Title,
		"hasChild":       meta.HasChild,
		"nodeCreateTime": meta.NodeCreateTime,
		"nodeCreator":    meta.NodeCreator,
		"creator":        meta.Creator,
		"owner":          meta.Owner,
		"objCreateTime":  meta.ObjCreateTime,
		"objEditTime":    meta.ObjEditTime,
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
