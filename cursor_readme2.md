# 飞书 MCP Go 版本项目规划 (基于 mcp-go 库)

## 项目概述

基于优秀的开源库 [mcp-go](https://github.com/mark3labs/mcp-go) (6.1k stars)，快速构建飞书文档操作的 MCP 服务器，为 Cursor、Windsurf、Cline 等 AI 驱动的编码工具提供访问飞书文档的能力。

**核心优势**：
- 🏗️ **基于成熟库**：使用经过验证的 mcp-go 框架，减少 80% 基础开发工作
- ⚡ **快速开发**：专注业务逻辑，预计开发周期从 10 周缩短至 4-5 周  
- 🔧 **功能丰富**：开箱即用的 Session 管理、中间件、Hooks 等高级特性
- 🎯 **稳定可靠**：基于 6.1k stars 的成熟开源项目，质量有保障

参考项目：https://github.com/cso1z/Feishu-MCP

## mcp-go 库分析

### 🎯 核心优势

**1. 完整的 MCP 协议实现**
- 完全符合 MCP 1.0 规范
- 支持 stdio、HTTP/SSE、streamable-HTTP 多种传输层
- 自动处理协议握手和消息序列化

**2. 高级功能特性**
- **Session 管理**：支持多客户端会话，per-session 工具定制
- **中间件系统**：支持请求中间件和错误恢复中间件
- **Hooks 机制**：请求生命周期监控，便于遥测和日志
- **工具过滤**：基于会话的工具权限控制

**3. 开发者友好**
- 清晰的 API 设计和丰富的示例
- 完善的文档和类型定义
- 内置测试工具 (mcptest)

### 🏗️ 重新设计的技术架构

基于 mcp-go 库的新架构：

```
┌─────────────────────────────────────────────────────────────┐
│                    MCP Client (Cursor/Cline)                │
└─────────────────────┬───────────────────────────────────────┘
                      │ MCP Protocol (JSON-RPC)
┌─────────────────────▼───────────────────────────────────────┐
│                  mcp-go Framework                           │
│  ✅ Protocol Layer    ✅ Transport Layer                    │
│  ✅ Session Management ✅ Middleware System                │
│  ✅ Tool Registry     ✅ Hook System                       │
└─────────────────────┬───────────────────────────────────────┘
                      │ 我们只需要实现这一层 ⬇️
┌─────────────────────▼───────────────────────────────────────┐
│              Feishu MCP Business Logic                     │
├─────────────────────────────────────────────────────────────┤
│  Feishu Tools (17个工具实现)                                │
│  ├── Document Tools (5个)                                  │
│  ├── Block Tools (8个)                                     │
│  ├── Folder Tools (3个)                                    │
│  └── Utility Tools (2个)                                   │
├─────────────────────────────────────────────────────────────┤
│  Feishu API Client                                         │
│  ├── Authentication & Token Management                     │
│  ├── HTTP Client with Retry & Rate Limiting               │
│  └── Request/Response Models                              │
└─────────────────────┬───────────────────────────────────────┘
                      │ HTTPS/REST API
┌─────────────────────▼───────────────────────────────────────┐
│              Feishu Open Platform APIs                     │
└─────────────────────────────────────────────────────────────┘
```

## 🚀 优化后的项目结构

```
feishu-mcp-go/
├── cmd/
│   └── server/
│       └── main.go              # 程序入口，使用 mcp-go 框架
├── internal/
│   ├── config/
│   │   └── config.go           # 配置管理
│   ├── feishu/
│   │   ├── client.go           # 飞书API客户端
│   │   ├── auth.go             # 认证和Token管理
│   │   └── models.go           # 数据模型
│   ├── tools/                  # 飞书工具实现（主要开发内容）
│   │   ├── document/
│   │   │   ├── create.go
│   │   │   ├── get.go
│   │   │   ├── search.go
│   │   │   └── ...
│   │   ├── block/
│   │   │   ├── create.go
│   │   │   ├── batch.go
│   │   │   └── ...
│   │   ├── folder/
│   │   └── utils/
│   └── middleware/
│       ├── auth.go             # 飞书认证中间件
│       ├── logging.go          # 日志中间件
│       └── recovery.go         # 错误恢复中间件
├── pkg/
│   └── errors/
│       └── errors.go           # 自定义错误类型
├── go.mod                      # 依赖 mcp-go
├── go.sum
├── Makefile
├── Dockerfile
└── README.md
```

## 🎯 简化后的核心依赖

```go
// go.mod - 大幅简化的依赖
module github.com/your-org/feishu-mcp-go

go 1.21

require (
    // 核心MCP框架 - 这个就够了！
    github.com/mark3labs/mcp-go v0.32.0
    
    // 飞书API客户端
    github.com/go-resty/resty/v2 v2.7.0
    
    // 配置和工具
    github.com/spf13/viper v1.16.0
    github.com/spf13/cobra v1.7.0
    
    // 缓存（可选）
    github.com/patrickmn/go-cache v2.1.0+incompatible
    
    // 测试
    github.com/stretchr/testify v1.8.4
)
```

## 🛠️ 核心实现代码示例

### 主程序入口

```go
// cmd/server/main.go
package main

import (
    "context"
    "log"
    
    "github.com/mark3labs/mcp-go/mcp"
    "github.com/mark3labs/mcp-go/server"
    
    "your-org/feishu-mcp-go/internal/tools"
    "your-org/feishu-mcp-go/internal/config"
)

func main() {
    // 加载配置
    cfg := config.Load()
    
    // 创建 MCP 服务器
    s := server.NewMCPServer(
        "Feishu MCP",
        "1.0.0",
        server.WithToolCapabilities(true),
        server.WithResourceCapabilities(true),
        server.WithRecovery(),
    )
    
    // 注册飞书工具
    tools.RegisterAll(s, cfg)
    
    // 启动服务器
    if err := s.Serve(context.Background()); err != nil {
        log.Fatal(err)
    }
}
```

### 工具注册

```go
// internal/tools/register.go
package tools

import (
    "github.com/mark3labs/mcp-go/server"
    "your-org/feishu-mcp-go/internal/config"
)

func RegisterAll(s *server.MCPServer, cfg *config.Config) {
    // 创建飞书客户端
    client := feishu.NewClient(cfg.Feishu)
    
    // 注册文档工具
    RegisterDocumentTools(s, client)
    RegisterBlockTools(s, client)
    RegisterFolderTools(s, client)
    RegisterUtilityTools(s, client)
}
```

### 工具实现示例

```go
// internal/tools/document/create.go
package document

import (
    "context"
    
    "github.com/mark3labs/mcp-go/mcp"
    "github.com/mark3labs/mcp-go/server"
)

func RegisterCreateTool(s *server.MCPServer, client *feishu.Client) {
    tool := mcp.NewTool(
        "create_feishu_document",
        mcp.WithDescription("创建新的飞书文档"),
        mcp.WithInputSchema(map[string]any{
            "type": "object",
            "properties": map[string]any{
                "title": map[string]any{
                    "type": "string",
                    "description": "文档标题",
                },
                "folderToken": map[string]any{
                    "type": "string", 
                    "description": "文件夹Token",
                },
            },
            "required": []string{"title", "folderToken"},
        }),
    )
    
    s.AddTool(tool, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
        // 解析参数
        title := req.Params.Arguments["title"].(string)
        folderToken := req.Params.Arguments["folderToken"].(string)
        
        // 调用飞书API
        doc, err := client.CreateDocument(ctx, title, folderToken)
        if err != nil {
            return mcp.NewToolResultError(err.Error()), nil
        }
        
        // 返回结果
        return mcp.NewToolResultText(
            fmt.Sprintf("文档创建成功！\n文档ID: %s\n访问链接: %s", 
                doc.DocumentID, doc.URL),
        ), nil
    })
}
```

## ⚡ 大幅缩短的开发计划

### 新的 4 周开发计划

**第 1 周：项目搭建和飞书集成**
- 基于 mcp-go 创建项目骨架
- 实现飞书 API 客户端和认证
- 完成 2-3 个基础工具验证架构

**第 2 周：核心工具实现**  
- 实现所有 17 个飞书工具
- 完善错误处理和参数验证
- 实现批量操作优化

**第 3 周：高级功能和优化**
- 添加缓存机制
- 实现中间件（日志、监控、限流）
- 性能优化和并发控制

**第 4 周：测试和发布**
- 完善单元测试和集成测试
- 编写文档和使用指南
- CI/CD 和发布准备

### 开发效率提升对比

| 方面 | 原方案 | 基于 mcp-go | 提升 |
|------|--------|-------------|------|
| **开发周期** | 10 周 | 4 周 | 🚀 **60% 缩短** |
| **代码量** | ~5000 行 | ~2000 行 | 🎯 **60% 减少** |
| **核心依赖** | 15+ 个 | 5 个 | ✅ **70% 简化** |
| **MCP 协议** | 完全自实现 | 开箱即用 | 🏗️ **零开发** |
| **传输层** | 需要实现 | 已完成 | ⚡ **即用** |
| **Session 管理** | 需要设计 | 已完成 | 🔧 **企业级** |

## 🎁 额外获得的高级特性

使用 mcp-go 库，我们免费获得了很多企业级特性：

### 1. 会话管理
```go
// 支持多客户端会话
session := &MySession{ID: "user-123"}
s.RegisterSession(ctx, session)

// 向特定客户端发送通知
s.SendNotificationToSpecificClient(session.ID, "update", data)
```

### 2. 中间件系统
```go
// 添加认证中间件
s.AddTool(tool, middleware.WithAuth(handler))

// 添加日志中间件
s.WithToolHandlerMiddleware(middleware.Logging)
```

### 3. 工具过滤
```go
// 基于用户权限过滤工具
s.WithToolFilter(func(ctx context.Context, tools []mcp.Tool) []mcp.Tool {
    session := server.ClientSessionFromContext(ctx)
    if session.IsAdmin() {
        return tools // 管理员看到所有工具
    }
    return filterUserTools(tools) // 普通用户只看部分工具
})
```

### 4. 监控和遥测
```go
// 请求生命周期监控
s.WithHooks(&server.Hooks{
    OnToolCall: func(ctx context.Context, req mcp.CallToolRequest) {
        // 记录工具调用
        log.Printf("Tool called: %s", req.Params.Name)
    },
})
```

## 🏆 最终技术方案总结

**选择 mcp-go 库的决定是明智的**：

1. **开发效率提升 60%**：从 10 周缩短到 4 周
2. **代码质量保证**：基于 6.1k stars 的成熟项目  
3. **功能更丰富**：获得企业级的会话管理和中间件系统
4. **维护成本降低**：协议层由开源社区维护
5. **扩展性更强**：标准化的架构便于后续功能扩展

我们的核心工作重点转为：
- ✅ **70% 时间**：实现 17 个飞书工具的业务逻辑
- ✅ **20% 时间**：飞书 API 客户端和认证
- ✅ **10% 时间**：配置、部署和文档

这是一个非常好的技术决策！🎉 