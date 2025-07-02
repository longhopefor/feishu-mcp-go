package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/mark3labs/mcp-go/server"
	"github.com/spf13/cobra"

	"feishu-mcp-go/internal/config"
	"feishu-mcp-go/internal/tools"
	"feishu-mcp-go/pkg/logger"
)

var (
	version   = "dev"
	buildTime = "unknown"
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "feishu-mcp",
		Short: "Feishu MCP Server - 为AI工具提供飞书文档操作能力",
		Long: `Feishu MCP Server 基于 Model Context Protocol (MCP) 协议，
为 Cursor、Windsurf、Cline 等 AI 驱动的编码工具提供访问飞书文档的能力。

支持的功能：
- 文档管理：创建、获取、搜索文档
- 内容操作：创建和编辑文本块、代码块、标题块等
- 文件夹管理：获取文件夹信息、创建文件夹
- 工具功能：Wiki链接转换、图片资源获取`,
		Run: runServer,
	}

	// 添加命令行参数
	rootCmd.Flags().String("config", "", "配置文件路径")
	rootCmd.Flags().String("feishu-app-id", "", "飞书应用 ID")
	rootCmd.Flags().String("feishu-app-secret", "", "飞书应用密钥")
	rootCmd.Flags().String("log-level", "info", "日志级别 (debug, info, warn, error)")
	rootCmd.Flags().Bool("stdio", false, "使用 stdio 模式运行")
	rootCmd.Flags().Bool("version", false, "显示版本信息")

	// 添加版本命令
	var versionCmd = &cobra.Command{
		Use:   "version",
		Short: "显示版本信息",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("Feishu MCP Server\n")
			fmt.Printf("Version: %s\n", version)
			fmt.Printf("Build Time: %s\n", buildTime)
			fmt.Printf("Based on: github.com/mark3labs/mcp-go\n")
		},
	}
	rootCmd.AddCommand(versionCmd)

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func runServer(cmd *cobra.Command, args []string) {
	// 处理版本标志
	if showVersion, _ := cmd.Flags().GetBool("version"); showVersion {
		fmt.Printf("Feishu MCP Server %s (built %s)\n", version, buildTime)
		return
	}

	// 加载配置
	cfg, err := config.Load(cmd)
	if err != nil {
		log.Fatalf("配置加载失败: %v", err)
	}

	// 初始化日志
	logger := logger.New(cfg.Log.Level)
	logger.Info("启动 Feishu MCP Server",
		"version", version,
		"mode", getRunMode(cmd))

	// 验证飞书配置
	if err := cfg.Feishu.Validate(); err != nil {
		logger.Fatal("飞书配置验证失败", "error", err)
	}

	// 创建 MCP 服务器
	s := server.NewMCPServer(
		"Feishu MCP",
		version,
		server.WithToolCapabilities(true),
		server.WithResourceCapabilities(false, false), // subscribe=false, listChanged=false
		server.WithRecovery(),                         // 启用错误恢复中间件
	)

	// 注册飞书工具
	if err := tools.RegisterAll(s, cfg, logger); err != nil {
		logger.Fatal("工具注册失败", "error", err)
	}

	logger.Info("所有飞书工具注册完成")

	// 启动服务器
	logger.Info("Feishu MCP Server 启动中...")
	if err := server.ServeStdio(s); err != nil {
		logger.Fatal("服务器启动失败", "error", err)
	}

	logger.Info("Feishu MCP Server 已安全关闭")
}

func getRunMode(cmd *cobra.Command) string {
	if stdio, _ := cmd.Flags().GetBool("stdio"); stdio {
		return "stdio"
	}
	return "auto"
}

func init() {
	flag.Parse()
}
