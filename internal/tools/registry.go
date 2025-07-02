package tools

import (
	"fmt"

	"github.com/mark3labs/mcp-go/server"

	"feishu-mcp-go/internal/config"
	"feishu-mcp-go/pkg/feishu"
	"feishu-mcp-go/pkg/logger"
)

// ToolRegistry 工具注册器
type ToolRegistry struct {
	server *server.MCPServer
	feishu *feishu.Client
	logger logger.Logger
	config *config.Config
}

// NewRegistry 创建新的工具注册器
func NewRegistry(s *server.MCPServer, cfg *config.Config, log logger.Logger) (*ToolRegistry, error) {
	// 创建飞书客户端
	feishuClient, err := feishu.NewClient(cfg.Feishu.AppID, cfg.Feishu.AppSecret, cfg.Feishu.BaseURL)
	if err != nil {
		return nil, fmt.Errorf("创建飞书客户端失败: %w", err)
	}

	return &ToolRegistry{
		server: s,
		feishu: feishuClient,
		logger: log,
		config: cfg,
	}, nil
}

// RegisterAll 注册所有飞书工具
func RegisterAll(s *server.MCPServer, cfg *config.Config, log logger.Logger) error {
	registry, err := NewRegistry(s, cfg, log)
	if err != nil {
		return err
	}

	// 注册文档管理工具
	if err := registry.registerDocumentTools(); err != nil {
		return fmt.Errorf("注册文档管理工具失败: %w", err)
	}

	// 注册内容操作工具
	if err := registry.registerContentTools(); err != nil {
		return fmt.Errorf("注册内容操作工具失败: %w", err)
	}

	// 注册文件夹管理工具
	if err := registry.registerFolderTools(); err != nil {
		return fmt.Errorf("注册文件夹管理工具失败: %w", err)
	}

	// 注册工具功能
	if err := registry.registerUtilityTools(); err != nil {
		return fmt.Errorf("注册工具功能失败: %w", err)
	}

	log.Info("所有飞书工具注册完成", "count", registry.getToolCount())
	return nil
}

// registerDocumentTools 注册文档管理工具
func (r *ToolRegistry) registerDocumentTools() error {
	tools := []server.ServerTool{
		NewCreateDocumentTool(r.feishu, r.logger),
		NewGetDocumentInfoTool(r.feishu, r.logger),
		NewGetDocumentContentTool(r.feishu, r.logger),
		NewGetDocumentBlocksTool(r.feishu, r.logger),
		NewSearchDocumentsTool(r.feishu, r.logger),
	}

	r.server.AddTools(tools...)
	for _, tool := range tools {
		r.logger.Debug("注册文档管理工具", "name", tool.Tool.Name)
	}

	return nil
}

// registerContentTools 注册内容操作工具
func (r *ToolRegistry) registerContentTools() error {
	tools := []server.ServerTool{
		NewGetBlockContentTool(r.feishu, r.logger),
		NewUpdateBlockTextTool(r.feishu, r.logger),
		NewBatchCreateBlocksTool(r.feishu, r.logger),
		NewCreateTextBlockTool(r.feishu, r.logger),
		NewCreateCodeBlockTool(r.feishu, r.logger),
		NewCreateHeadingBlockTool(r.feishu, r.logger),
		NewCreateListBlockTool(r.feishu, r.logger),
		NewDeleteBlocksTool(r.feishu, r.logger),
	}

	r.server.AddTools(tools...)
	for _, tool := range tools {
		r.logger.Debug("注册内容操作工具", "name", tool.Tool.Name)
	}

	return nil
}

// registerFolderTools 注册文件夹管理工具
func (r *ToolRegistry) registerFolderTools() error {
	tools := []server.ServerTool{
		NewGetRootFolderInfoTool(r.feishu, r.logger),
		NewGetFolderFilesTool(r.feishu, r.logger),
		NewCreateFolderTool(r.feishu, r.logger),
	}

	r.server.AddTools(tools...)
	for _, tool := range tools {
		r.logger.Debug("注册文件夹管理工具", "name", tool.Tool.Name)
	}

	return nil
}

// registerUtilityTools 注册工具功能
func (r *ToolRegistry) registerUtilityTools() error {
	tools := []server.ServerTool{
		NewConvertWikiTool(r.feishu, r.logger),
		NewGetImageResourceTool(r.feishu, r.logger),
	}

	r.server.AddTools(tools...)
	for _, tool := range tools {
		r.logger.Debug("注册工具功能", "name", tool.Tool.Name)
	}

	return nil
}

// getToolCount 获取已注册工具数量
func (r *ToolRegistry) getToolCount() int {
	return 17 // 总共17个工具
}
