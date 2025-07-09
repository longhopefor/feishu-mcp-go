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

// 文件夹管理工具实现

func NewGetRootFolderInfoTool(feishu *feishu.Client, logger logger.Logger) server.ServerTool {
	tool := mcp.NewTool("get_feishu_root_folder_info",
		mcp.WithDescription("获取飞书云盘根文件夹信息"))
	handler := func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return getRootFolderInfo(feishu, logger, request)
	}
	return server.ServerTool{Tool: tool, Handler: handler}
}

func NewGetFolderFilesTool(feishu *feishu.Client, logger logger.Logger) server.ServerTool {
	tool := mcp.NewTool("get_feishu_folder_files",
		mcp.WithDescription("获取指定文件夹中的文件和子文件夹列表"))
	handler := func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return getFolderFiles(feishu, logger, request)
	}
	return server.ServerTool{Tool: tool, Handler: handler}
}

func NewCreateFolderTool(feishu *feishu.Client, logger logger.Logger) server.ServerTool {
	tool := mcp.NewTool("create_feishu_folder",
		mcp.WithDescription("在指定父文件夹中创建新文件夹"))
	handler := func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return createFolder(feishu, logger, request)
	}
	return server.ServerTool{Tool: tool, Handler: handler}
}

// getRootFolderInfo 获取根文件夹信息实现
func getRootFolderInfo(feishu *feishu.Client, logger logger.Logger, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	logger.Info("开始获取根文件夹信息", "name", request.Params.Name)

	// 调用飞书API获取根文件夹信息
	folderInfo, err := feishu.GetRootFolderInfo()
	if err != nil {
		logger.Error("获取根文件夹信息失败", "error", err)
		return nil, fmt.Errorf("获取根文件夹信息失败: %w", err)
	}

	logger.Info("获取根文件夹信息成功", "folders_count", len(folderInfo.Folders))

	// 构建成功响应
	result := map[string]interface{}{
		"success":       true,
		"message":       "获取根文件夹信息成功",
		"folders_count": len(folderInfo.Folders),
		"folders":       folderInfo.Folders,
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

// getFolderFiles 获取文件夹文件列表实现
func getFolderFiles(feishu *feishu.Client, logger logger.Logger, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	logger.Info("开始获取文件夹文件列表", "name", request.Params.Name)

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

	folderToken, ok := args["folderToken"].(string)
	if !ok || folderToken == "" {
		return nil, fmt.Errorf("folderToken参数必须是非空字符串")
	}

	// 获取可选参数
	pageSize := 10 // 默认值
	if ps, ok := args["pageSize"].(float64); ok {
		pageSize = int(ps)
	}

	pageToken := ""
	if pt, ok := args["pageToken"].(string); ok {
		pageToken = pt
	}

	// 调用飞书API获取文件夹文件列表
	folderFiles, err := feishu.GetFolderFiles(folderToken, pageSize, pageToken)
	if err != nil {
		logger.Error("获取文件夹文件列表失败", "error", err)
		return nil, fmt.Errorf("获取文件夹文件列表失败: %w", err)
	}

	logger.Info("获取文件夹文件列表成功", "folderToken", folderToken, "files_count", len(folderFiles.Folders))

	// 构建成功响应
	result := map[string]interface{}{
		"success":     true,
		"message":     "获取文件夹文件列表成功",
		"folderToken": folderToken,
		"files_count": len(folderFiles.Folders),
		"files":       folderFiles.Folders,
		"pageSize":    pageSize,
		"pageToken":   pageToken,
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

// createFolder 创建文件夹实现
func createFolder(feishu *feishu.Client, logger logger.Logger, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	logger.Info("开始创建文件夹", "name", request.Params.Name)

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

	name, ok := args["name"].(string)
	if !ok || name == "" {
		return nil, fmt.Errorf("name参数必须是非空字符串")
	}

	// 获取可选的父文件夹token
	parentToken, _ := args["parentToken"].(string)

	// 调用飞书API创建文件夹
	folderResp, err := feishu.CreateFolder(name, parentToken)
	if err != nil {
		logger.Error("创建文件夹失败", "error", err)
		return nil, fmt.Errorf("创建文件夹失败: %w", err)
	}

	logger.Info("创建文件夹成功", "name", name, "token", folderResp.Data.Token)

	// 构建成功响应
	result := map[string]interface{}{
		"success":     true,
		"message":     "创建文件夹成功",
		"name":        name,
		"token":       folderResp.Data.Token,
		"url":         folderResp.Data.URL,
		"parentToken": parentToken,
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
