# 飞书 MCP Go

[![Go Version](https://img.shields.io/badge/go-%3E%3D1.21-blue.svg)](https://golang.org/)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](https://opensource.org/licenses/MIT)

基于 [mcp-go](https://github.com/mark3labs/mcp-go) 框架的飞书文档操作 MCP 服务器，为 Cursor、Windsurf、Cline 等 AI 驱动的编码工具提供访问飞书文档的能力。

> 📋 **项目开发进度**: [查看详细开发计划](./todo_list.md) | [技术实现方案](./cursor_readme.md)  
> 🚀 **当前状态**: Phase 2 飞书API集成中 (25%完成)

## 🚀 特性

- **完全兼容 MCP 协议**：支持所有 AI 编程工具
- **17个核心工具**：涵盖文档管理、内容操作、文件夹管理等
- **高性能**：基于 Go 语言和 mcp-go 框架
- **易部署**：单一可执行文件，支持多平台
- **企业级特性**：Session管理、中间件系统、错误恢复

## 📦 支持的工具

### 📄 文档管理（5个工具）
- `create_feishu_document` - 创建飞书文档
- `get_feishu_document_info` - 获取文档信息  
- `get_feishu_document_content` - 获取文档内容
- `get_feishu_document_blocks` - 获取文档块结构
- `search_feishu_documents` - 搜索文档

### ✏️ 内容操作（8个工具）
- `get_feishu_block_content` - 获取块内容
- `update_feishu_block_text` - 更新块文本
- `batch_create_feishu_blocks` - 批量创建块
- `create_feishu_text_block` - 创建文本块
- `create_feishu_code_block` - 创建代码块
- `create_feishu_heading_block` - 创建标题块
- `create_feishu_list_block` - 创建列表块
- `delete_feishu_document_blocks` - 删除文档块

### 📁 文件夹管理（3个工具）
- `get_feishu_root_folder_info` - 获取根文件夹信息
- `get_feishu_folder_files` - 获取文件夹文件列表
- `create_feishu_folder` - 创建文件夹

### 🔧 工具功能（2个工具）
- `convert_feishu_wiki_to_document_id` - Wiki链接转文档ID
- `get_feishu_image_resource` - 获取图片资源

## 🛠️ 安装

### 方式1：使用 make（推荐）

```bash
# 克隆仓库
git clone https://github.com/yourname/feishu-mcp-go.git
cd feishu-mcp-go

# 下载依赖并构建
make build

# 运行
make run
```

### 方式2：手动构建

```bash
# 下载依赖
go mod download
go mod tidy

# 构建
go build -o build/feishu-mcp ./cmd/server

# 运行
./build/feishu-mcp
```

### 方式3：直接运行

```bash
go run ./cmd/server
```

## ⚙️ 配置

### 环境变量配置（推荐）

```bash
# 必需配置
export FEISHU_MCP_FEISHU_APP_ID="cli_xxxxxxxxxxxxxxxx"
export FEISHU_MCP_FEISHU_APP_SECRET="xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"

# 可选配置
export FEISHU_MCP_LOG_LEVEL="info"
export FEISHU_MCP_CACHE_ENABLED="true"
export FEISHU_MCP_SERVER_MODE="auto"
```

### 配置文件

复制 `configs/config.yaml` 并修改：

```yaml
feishu:
  app_id: "your_app_id"
  app_secret: "your_app_secret"
  base_url: "https://open.feishu.cn/open-apis"
  timeout: 30

log:
  level: "info"
  format: "json"

cache:
  enabled: true
  ttl: 3600

server:
  mode: "auto"
  port: 3333
```

### 命令行参数

```bash
./feishu-mcp \
  --feishu-app-id="cli_xxxxxxxxxxxxxxxx" \
  --feishu-app-secret="xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx" \
  --log-level="debug" \
  --stdio
```

## 🚀 使用方法

### 在 Cursor/Windsurf 中配置

1. 打开设置，找到 MCP 配置
2. 添加新的 MCP 服务器：

```json
{
  "mcpServers": {
    "feishu": {
      "command": "/path/to/feishu-mcp",
      "args": ["--stdio"],
      "env": {
        "FEISHU_MCP_FEISHU_APP_ID": "cli_xxxxxxxxxxxxxxxx",
        "FEISHU_MCP_FEISHU_APP_SECRET": "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
      }
    }
  }
}
```

### 命令行使用

```bash
# 显示帮助
./feishu-mcp --help

# 以 stdio 模式运行
./feishu-mcp --stdio

# 查看版本
./feishu-mcp version

# 开发模式（调试日志）
./feishu-mcp --log-level=debug
```

## 🏗️ 开发

### 项目结构

```
feishu-mcp-go/
├── cmd/server/              # 主程序入口
├── internal/
│   ├── config/             # 配置管理
│   └── tools/              # 工具实现
├── pkg/
│   ├── feishu/             # 飞书API客户端
│   └── logger/             # 日志管理
├── configs/                # 配置文件
├── build/                  # 构建输出
├── go.mod                  # Go模块定义
├── Makefile               # 构建脚本
└── README.md              # 项目文档
```

### Make 命令

```bash
make help        # 显示帮助
make deps        # 下载依赖
make build       # 构建二进制
make run         # 运行程序
make dev         # 开发模式运行
make test        # 运行测试
make lint        # 代码检查
make clean       # 清理构建文件
make install     # 安装到系统
make build-all   # 跨平台构建
```

### 添加新工具

1. 在 `internal/tools/` 下创建工具实现
2. 在 `registry.go` 中注册工具
3. 实现 `server.Tool` 接口

```go
type MyTool struct {
    feishu *feishu.Client
    logger logger.Logger
}

func (t *MyTool) Execute(ctx context.Context, request server.ToolRequest) (*server.ToolResponse, error) {
    // 实现工具逻辑
    return &server.ToolResponse{
        Content: []server.ToolResponseContent{
            {Type: "text", Text: "result"},
        },
    }, nil
}
```

## 📋 系统要求

- Go 1.21 或更高版本
- 有效的飞书应用 ID 和密钥
- 网络连接（访问飞书 API）

## 🤝 参与贡献

1. Fork 项目
2. 创建功能分支 (`git checkout -b feature/AmazingFeature`)
3. 提交变更 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 打开 Pull Request

## 📄 许可证

本项目基于 MIT 许可证开源。详见 [LICENSE](LICENSE) 文件。

## 🙏 致谢

- [mcp-go](https://github.com/mark3labs/mcp-go) - 优秀的 MCP Go 框架
- [Feishu-MCP](https://github.com/cso1z/Feishu-MCP) - TypeScript 版本参考实现

## 📞 支持

如有问题或建议，请：

1. 查看 [Issues](https://github.com/yourname/feishu-mcp-go/issues)
2. 创建新的 Issue
3. 参考项目文档 