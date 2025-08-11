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
		mcp.WithDescription("获取指定文件夹中的文件和子文件夹列表"),
		mcp.WithString("folderToken", mcp.Description("文件夹token，用于标识要获取文件列表的文件夹"), mcp.Required()),
		mcp.WithNumber("pageSize", mcp.Description("每页返回的文件数量，默认为50"), mcp.DefaultNumber(50)),
		mcp.WithString("pageToken", mcp.Description("分页令牌，用于获取下一页数据，可选")))
	handler := func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return getFolderFiles(feishu, logger, request)
	}
	return server.ServerTool{Tool: tool, Handler: handler}
}

func NewCreateFolderTool(feishu *feishu.Client, logger logger.Logger) server.ServerTool {
	tool := mcp.NewTool("create_feishu_folder",
		mcp.WithDescription("在指定父文件夹中创建新文件夹"),
		mcp.WithString("name", mcp.Description("新文件夹的名称"), mcp.Required()),
		mcp.WithString("parentToken", mcp.Description("父文件夹的token，用于指定新文件夹的创建位置"), mcp.Required()))
	handler := func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return createFolder(feishu, logger, request)
	}
	return server.ServerTool{Tool: tool, Handler: handler}
}

func NewGetDriveFilesWithMetaTool(feishu *feishu.Client, logger logger.Logger) server.ServerTool {
	tool := mcp.NewTool("get_feishu_drive_files_with_meta",
		mcp.WithDescription("获取云空间目录下所有文件的详细元数据信息"),
		mcp.WithString("folderToken", mcp.Description("文件夹token，用于标识要获取元数据的文件夹"), mcp.Required()),
		mcp.WithNumber("pageSize", mcp.Description("每页返回的文件数量，默认为50"), mcp.DefaultNumber(50)),
		mcp.WithString("pageToken", mcp.Description("分页令牌，用于获取下一页数据，可选")))
	handler := func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return getDriveFilesWithMeta(feishu, logger, request)
	}
	return server.ServerTool{Tool: tool, Handler: handler}
}

func NewGetDriveMetaTool(feishu *feishu.Client, logger logger.Logger) server.ServerTool {
	tool := mcp.NewTool("get_feishu_drive_meta",
		mcp.WithDescription("获取云空间目录/文件的元数据信息"),
		mcp.WithString("requestDoc", mcp.Description("请求的文档token，用于获取指定文档的元数据"), mcp.Required()))
	handler := func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return getDriveMeta(feishu, logger, request)
	}
	return server.ServerTool{Tool: tool, Handler: handler}
}

func NewGetAllDriveFilesTool(feishu *feishu.Client, logger logger.Logger) server.ServerTool {
	tool := mcp.NewTool("get_all_feishu_drive_files",
		mcp.WithDescription("获取指定目录下所有文件（支持分页和数量限制）"),
		mcp.WithString("folderToken", mcp.Description("文件夹token，用于标识要获取所有文件的文件夹"), mcp.Required()),
		mcp.WithNumber("pageSize", mcp.Description("每页返回的文件数量，默认为50"), mcp.DefaultNumber(50)),
		mcp.WithString("pageToken", mcp.Description("分页令牌，用于获取下一页数据，可选")))
	handler := func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return getAllDriveFiles(feishu, logger, request)
	}
	return server.ServerTool{Tool: tool, Handler: handler}
}

func NewGetRootFolderMetaTool(feishu *feishu.Client, logger logger.Logger) server.ServerTool {
	tool := mcp.NewTool("get_feishu_root_folder_meta",
		mcp.WithDescription("获取云空间根文件夹元数据信息（使用新的API端点）"))
	handler := func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return getRootFolderMeta(feishu, logger, request)
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

	logger.Info("获取根文件夹信息成功", "files_count", len(folderInfo.Data.Files))

	// 直接返回原始结构体内容，保证与飞书官方API一致
	resultContent, _ := json.MarshalIndent(folderInfo, "", "  ")

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

	logger.Info("获取文件夹文件列表成功", "folderToken", folderToken, "files_count", len(folderFiles.Data.Files))

	// 构建成功响应
	result := map[string]interface{}{
		"success":     true,
		"message":     "获取文件夹文件列表成功",
		"folderToken": folderToken,
		"files_count": len(folderFiles.Data.Files),
		"files":       folderFiles.Data.Files,
		"has_more":    folderFiles.Data.HasMore,
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

// getDriveFilesWithMeta 获取云空间目录下所有文件的详细元数据信息实现
func getDriveFilesWithMeta(feishu *feishu.Client, logger logger.Logger, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	logger.Info("开始获取云空间文件详细元数据", "name", request.Params.Name)

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

	// 获取参数
	folderToken := ""
	if ft, ok := args["folderToken"].(string); ok {
		folderToken = ft
	}

	pageSize := 50 // 默认值
	if ps, ok := args["pageSize"].(float64); ok {
		pageSize = int(ps)
	}

	pageToken := ""
	if pt, ok := args["pageToken"].(string); ok {
		pageToken = pt
	}

	// 调用飞书API获取文件详细信息
	filesResp, err := feishu.GetDriveFilesWithMeta(folderToken, pageSize, pageToken)
	if err != nil {
		logger.Error("获取云空间文件详细元数据失败", "error", err)
		return nil, fmt.Errorf("获取云空间文件详细元数据失败: %w", err)
	}

	logger.Info("获取云空间文件详细元数据成功", "folderToken", folderToken, "files_count", len(filesResp.Data.Files))

	// 构建成功响应
	result := map[string]interface{}{
		"success":       true,
		"message":       "获取云空间文件详细元数据成功",
		"folderToken":   folderToken,
		"files_count":   len(filesResp.Data.Files),
		"files":         filesResp.Data.Files,
		"pageSize":      pageSize,
		"pageToken":     pageToken,
		"nextPageToken": filesResp.Data.NextPageToken,
		"hasMore":       filesResp.Data.HasMore,
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

// getDriveMeta 获取云空间目录/文件的元数据信息实现
func getDriveMeta(feishu *feishu.Client, logger logger.Logger, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	logger.Info("开始获取云空间文件元数据", "name", request.Params.Name)

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

	fileToken, ok := args["fileToken"].(string)
	if !ok || fileToken == "" {
		return nil, fmt.Errorf("fileToken参数必须是非空字符串")
	}

	// 调用飞书API获取文件元数据
	metaResp, err := feishu.GetDriveMeta(fileToken)
	if err != nil {
		logger.Error("获取云空间文件元数据失败", "error", err)
		return nil, fmt.Errorf("获取云空间文件元数据失败: %w", err)
	}

	logger.Info("获取云空间文件元数据成功", "fileToken", fileToken, "name", metaResp.Data.Name)

	// 构建成功响应
	result := map[string]interface{}{
		"success":   true,
		"message":   "获取云空间文件元数据成功",
		"fileToken": fileToken,
		"meta":      metaResp.Data,
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

// getAllDriveFiles 获取指定目录下所有文件（支持分页和数量限制）实现
func getAllDriveFiles(feishu *feishu.Client, logger logger.Logger, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	logger.Info("开始获取目录下所有文件", "name", request.Params.Name)

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

	// 获取参数
	folderToken := ""
	if ft, ok := args["folderToken"].(string); ok {
		folderToken = ft
	}

	maxFiles := 0 // 默认无限制
	if mf, ok := args["maxFiles"].(float64); ok {
		maxFiles = int(mf)
	}

	// 调用飞书API获取所有文件
	allFiles, err := feishu.GetAllDriveFiles(folderToken, maxFiles)
	if err != nil {
		logger.Error("获取目录下所有文件失败", "error", err)
		return nil, fmt.Errorf("获取目录下所有文件失败: %w", err)
	}

	logger.Info("获取目录下所有文件成功", "folderToken", folderToken, "files_count", len(allFiles))

	// 构建成功响应
	result := map[string]interface{}{
		"success":     true,
		"message":     "获取目录下所有文件成功",
		"folderToken": folderToken,
		"files_count": len(allFiles),
		"files":       allFiles,
		"maxFiles":    maxFiles,
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

// getRootFolderMeta 获取根文件夹元数据实现（使用新的API端点）
func getRootFolderMeta(feishu *feishu.Client, logger logger.Logger, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	logger.Info("开始获取根文件夹元数据", "name", request.Params.Name)

	// 调用飞书API获取根文件夹元数据
	metaResp, err := feishu.GetRootFolderMeta()
	if err != nil {
		logger.Error("获取根文件夹元数据失败", "error", err)
		return nil, fmt.Errorf("获取根文件夹元数据失败: %w", err)
	}

	logger.Info("获取根文件夹元数据成功", "token", metaResp.Data.Token, "id", metaResp.Data.ID)

	// 直接返回原始结构体内容，保证与飞书官方API一致
	resultContent, _ := json.MarshalIndent(metaResp, "", "  ")

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.TextContent{
				Type: "text",
				Text: string(resultContent),
			},
		},
	}, nil
}
