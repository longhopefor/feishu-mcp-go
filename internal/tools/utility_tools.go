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

// 工具功能实现

func NewConvertWikiTool(feishu *feishu.Client, logger logger.Logger) server.ServerTool {
	tool := mcp.NewTool("convert_feishu_wiki_to_document_id",
		mcp.WithDescription("将飞书Wiki链接转换为文档ID"))
	handler := func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return convertWiki(feishu, logger, request)
	}
	return server.ServerTool{Tool: tool, Handler: handler}
}

func NewGetImageResourceTool(feishu *feishu.Client, logger logger.Logger) server.ServerTool {
	tool := mcp.NewTool("get_feishu_image_resource",
		mcp.WithDescription("获取飞书图片资源信息和内容"))
	handler := func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return getImageResource(feishu, logger, request)
	}
	return server.ServerTool{Tool: tool, Handler: handler}
}

// convertWiki 转换Wiki链接为文档ID实现
func convertWiki(feishu *feishu.Client, logger logger.Logger, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	logger.Info("开始转换Wiki链接", "name", request.Params.Name)

	// 解析参数
	var args map[string]interface{}
	if request.Params.Arguments != nil {
		if argBytes, ok := request.Params.Arguments.([]byte); ok {
			if err := json.Unmarshal(argBytes, &args); err != nil {
				return nil, fmt.Errorf("解析参数失败: %w", err)
			}
		} else if argMap, ok := request.Params.Arguments.(map[string]interface{}); ok {
			args = argMap
		} else {
			if argBytes, err := json.Marshal(request.Params.Arguments); err != nil {
				return nil, fmt.Errorf("参数格式错误: %w", err)
			} else if err := json.Unmarshal(argBytes, &args); err != nil {
				return nil, fmt.Errorf("解析参数失败: %w", err)
			}
		}
	}

	objToken, ok := args["objToken"].(string)
	if !ok || objToken == "" {
		return nil, fmt.Errorf("objToken参数必须是非空字符串")
	}

	objType, ok := args["objType"].(string)
	if !ok || objType == "" {
		objType = "doc" // 默认类型
	}

	// 调用飞书API转换Wiki
	wikiResp, err := feishu.ConvertWiki(objToken, objType)
	if err != nil {
		logger.Error("转换Wiki失败", "error", err)
		return nil, fmt.Errorf("转换Wiki失败: %w", err)
	}

	logger.Info("转换Wiki成功", "objToken", objToken, "documentID", wikiResp.Data.DocumentID)

	// 构建成功响应
	result := map[string]interface{}{
		"success":    true,
		"message":    "转换Wiki链接成功",
		"objToken":   objToken,
		"objType":    objType,
		"documentID": wikiResp.Data.DocumentID,
		"url":        wikiResp.Data.URL,
		"token":      wikiResp.Data.Token,
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

// getImageResource 获取图片资源实现
func getImageResource(feishu *feishu.Client, logger logger.Logger, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	logger.Info("开始获取图片资源", "name", request.Params.Name)

	// 解析参数
	var args map[string]interface{}
	if request.Params.Arguments != nil {
		if argBytes, ok := request.Params.Arguments.([]byte); ok {
			if err := json.Unmarshal(argBytes, &args); err != nil {
				return nil, fmt.Errorf("解析参数失败: %w", err)
			}
		} else if argMap, ok := request.Params.Arguments.(map[string]interface{}); ok {
			args = argMap
		} else {
			if argBytes, err := json.Marshal(request.Params.Arguments); err != nil {
				return nil, fmt.Errorf("参数格式错误: %w", err)
			} else if err := json.Unmarshal(argBytes, &args); err != nil {
				return nil, fmt.Errorf("解析参数失败: %w", err)
			}
		}
	}

	imageKey, ok := args["imageKey"].(string)
	if !ok || imageKey == "" {
		return nil, fmt.Errorf("imageKey参数必须是非空字符串")
	}

	// 调用飞书API获取图片资源
	imageResp, err := feishu.GetImageResource(imageKey)
	if err != nil {
		logger.Error("获取图片资源失败", "error", err)
		return nil, fmt.Errorf("获取图片资源失败: %w", err)
	}

	logger.Info("获取图片资源成功", "imageKey", imageKey, "contentType", imageResp.Data.ContentType)

	// 构建成功响应
	result := map[string]interface{}{
		"success":     true,
		"message":     "获取图片资源成功",
		"imageKey":    imageKey,
		"url":         imageResp.Data.URL,
		"contentType": imageResp.Data.ContentType,
		"size":        imageResp.Data.Size,
		"base64Data":  imageResp.Data.Base64Data,
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
