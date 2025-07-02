# 飞书 MCP Go 版本 - 详细技术实现方案

## 1. 架构设计深度分析

### 1.1 整体架构图

```
┌─────────────────────────────────────────────────────────────┐
│                    MCP Client (Cursor/Cline)                │
└─────────────────────┬───────────────────────────────────────┘
                      │ MCP Protocol (JSON-RPC)
                      │ Transport: stdio/HTTP+SSE
┌─────────────────────▼───────────────────────────────────────┐
│                  MCP Server (Go)                            │
├─────────────────────────────────────────────────────────────┤
│  Transport Layer                                            │
│  ├── StdIO Handler                                          │
│  └── HTTP/SSE Handler                                       │
├─────────────────────────────────────────────────────────────┤
│  Protocol Layer                                             │
│  ├── Message Parser/Serializer                             │
│  ├── Tool Registry & Dispatcher                            │
│  └── Resource Manager                                       │
├─────────────────────────────────────────────────────────────┤
│  Business Layer                                             │
│  ├── Document Handlers                                      │
│  ├── Block Handlers                                         │
│  ├── Folder Handlers                                        │
│  └── Utility Handlers                                       │
├─────────────────────────────────────────────────────────────┤
│  Integration Layer                                          │
│  ├── Auth Manager (Token Management)                       │
│  ├── API Client (HTTP Client Pool)                         │
│  ├── Cache Layer (Memory + Persistence)                    │
│  └── Rate Limiter                                          │
└─────────────────────┬───────────────────────────────────────┘
                      │ HTTPS/REST API
                      │
┌─────────────────────▼───────────────────────────────────────┐
│              Feishu Open Platform APIs                     │
└─────────────────────────────────────────────────────────────┘
```

### 1.2 核心接口设计

#### MCP Server 核心接口

```go
// MCP 服务器主接口
type MCPServer interface {
    Start(ctx context.Context) error
    Stop() error
    RegisterTool(tool Tool) error
    RegisterResource(resource Resource) error
}

// 工具处理接口
type Tool interface {
    Name() string
    Description() string
    InputSchema() map[string]interface{}
    Execute(ctx context.Context, args map[string]interface{}) (*ToolResult, error)
}

// 传输层接口
type Transport interface {
    Listen(ctx context.Context) error
    Send(message *MCPMessage) error
    Receive() (*MCPMessage, error)
    Close() error
}
```

#### 飞书 API 客户端接口

```go
// 飞书 API 客户端接口
type FeishuClient interface {
    // 认证相关
    GetAppAccessToken(ctx context.Context) (*TokenResponse, error)
    RefreshToken(ctx context.Context) error
    
    // 文档操作
    CreateDocument(ctx context.Context, req *CreateDocumentRequest) (*Document, error)
    GetDocument(ctx context.Context, docID string) (*Document, error)
    UpdateDocument(ctx context.Context, docID string, req *UpdateDocumentRequest) error
    DeleteDocument(ctx context.Context, docID string) error
    
    // 块操作
    GetBlocks(ctx context.Context, docID string) (*BlocksResponse, error)
    CreateBlocks(ctx context.Context, docID string, req *CreateBlocksRequest) (*CreateBlocksResponse, error)
    UpdateBlock(ctx context.Context, docID, blockID string, req *UpdateBlockRequest) error
    DeleteBlocks(ctx context.Context, docID string, req *DeleteBlocksRequest) error
    
    // 文件夹操作
    GetRootFolder(ctx context.Context) (*Folder, error)
    GetFolderFiles(ctx context.Context, folderToken string) (*FolderFilesResponse, error)
    CreateFolder(ctx context.Context, req *CreateFolderRequest) (*Folder, error)
}
```

## 2. 关键组件实现细节

### 2.1 MCP 协议实现

#### 消息类型定义

```go
// MCP 消息类型
type MCPMessage struct {
    JSONRPC string      `json:"jsonrpc"`
    ID      interface{} `json:"id,omitempty"`
    Method  string      `json:"method,omitempty"`
    Params  interface{} `json:"params,omitempty"`
    Result  interface{} `json:"result,omitempty"`
    Error   *MCPError   `json:"error,omitempty"`
}

// MCP 错误类型
type MCPError struct {
    Code    int         `json:"code"`
    Message string      `json:"message"`
    Data    interface{} `json:"data,omitempty"`
}

// 工具调用请求
type ToolCallRequest struct {
    Name      string                 `json:"name"`
    Arguments map[string]interface{} `json:"arguments"`
}

// 工具调用结果
type ToolResult struct {
    Content []ContentBlock `json:"content"`
    IsError bool          `json:"isError,omitempty"`
}

// 内容块
type ContentBlock struct {
    Type string      `json:"type"`
    Text string      `json:"text,omitempty"`
    Data interface{} `json:"data,omitempty"`
}
```

### 2.2 飞书 API 客户端实现

#### HTTP 客户端配置

```go
type FeishuClientConfig struct {
    AppID       string
    AppSecret   string
    BaseURL     string
    Timeout     time.Duration
    RetryCount  int
    RateLimit   int
    CacheEnabled bool
    CacheTTL    time.Duration
}

type feishuClient struct {
    config     *FeishuClientConfig
    httpClient *resty.Client
    tokenCache TokenCache
    rateLimiter RateLimiter
    logger     Logger
}

func NewFeishuClient(config *FeishuClientConfig) FeishuClient {
    client := resty.New().
        SetTimeout(config.Timeout).
        SetRetryCount(config.RetryCount).
        SetBaseURL(config.BaseURL).
        SetHeader("Content-Type", "application/json")
    
    return &feishuClient{
        config:      config,
        httpClient:  client,
        tokenCache:  NewTokenCache(config.CacheTTL),
        rateLimiter: NewRateLimiter(config.RateLimit),
        logger:      NewLogger(),
    }
}
```

#### Token 管理实现

```go
type TokenCache interface {
    Get(key string) (*TokenInfo, bool)
    Set(key string, token *TokenInfo, ttl time.Duration)
    Delete(key string)
}

type TokenInfo struct {
    AccessToken string
    ExpiresAt   time.Time
    TokenType   string
}

type tokenManager struct {
    client  *feishuClient
    cache   TokenCache
    mutex   sync.RWMutex
}

func (tm *tokenManager) GetValidToken(ctx context.Context) (string, error) {
    tm.mutex.RLock()
    if token, exists := tm.cache.Get("app_access_token"); exists {
        if time.Now().Before(token.ExpiresAt.Add(-5 * time.Minute)) {
            tm.mutex.RUnlock()
            return token.AccessToken, nil
        }
    }
    tm.mutex.RUnlock()
    
    tm.mutex.Lock()
    defer tm.mutex.Unlock()
    
    // 重新检查避免并发问题
    if token, exists := tm.cache.Get("app_access_token"); exists {
        if time.Now().Before(token.ExpiresAt.Add(-5 * time.Minute)) {
            return token.AccessToken, nil
        }
    }
    
    return tm.refreshToken(ctx)
}
```

## 3. 项目结构和开发规划

### 3.1 详细目录结构

```
feishu-mcp-go/
├── cmd/
│   └── server/
│       └── main.go              # 程序入口点
├── internal/
│   ├── config/
│   │   ├── config.go           # 配置结构和加载
│   │   └── validation.go       # 配置验证
│   ├── mcp/
│   │   ├── server.go           # MCP服务器实现
│   │   ├── transport/
│   │   │   ├── stdio.go        # Stdio传输层
│   │   │   └── http.go         # HTTP/SSE传输层
│   │   ├── protocol/
│   │   │   ├── handler.go      # 协议处理器
│   │   │   ├── message.go      # 消息类型定义
│   │   │   └── validator.go    # 消息验证
│   │   └── registry/
│   │       ├── tool.go         # 工具注册管理
│   │       └── resource.go     # 资源注册管理
│   ├── feishu/
│   │   ├── client.go           # 飞书API客户端
│   │   ├── auth/
│   │   │   ├── token.go        # Token管理
│   │   │   └── cache.go        # Token缓存
│   │   ├── api/
│   │   │   ├── document.go     # 文档API
│   │   │   ├── block.go        # 块API
│   │   │   ├── folder.go       # 文件夹API
│   │   │   └── utils.go        # 工具API
│   │   └── models/
│   │       ├── document.go     # 文档数据模型
│   │       ├── block.go        # 块数据模型
│   │       ├── folder.go       # 文件夹数据模型
│   │       └── common.go       # 通用数据模型
│   ├── handlers/
│   │   ├── base.go             # 基础处理器
│   │   ├── document/
│   │   │   ├── create.go       # 创建文档
│   │   │   ├── get.go          # 获取文档
│   │   │   ├── search.go       # 搜索文档
│   │   │   └── info.go         # 文档信息
│   │   ├── block/
│   │   │   ├── create.go       # 创建块
│   │   │   ├── batch.go        # 批量操作
│   │   │   ├── update.go       # 更新块
│   │   │   ├── delete.go       # 删除块
│   │   │   └── get.go          # 获取块
│   │   ├── folder/
│   │   │   ├── create.go       # 创建文件夹
│   │   │   ├── get.go          # 获取文件夹
│   │   │   └── root.go         # 根文件夹
│   │   └── utils/
│   │       ├── convert.go      # Wiki转换
│   │       └── image.go        # 图片资源
│   ├── cache/
│   │   ├── interface.go        # 缓存接口
│   │   ├── memory.go           # 内存缓存
│   │   └── redis.go            # Redis缓存
│   ├── middleware/
│   │   ├── ratelimit.go        # 速率限制
│   │   ├── metrics.go          # 指标收集
│   │   └── recovery.go         # 错误恢复
│   └── utils/
│       ├── logger.go           # 日志工具
│       ├── validator.go        # 参数验证
│       └── converter.go        # 数据转换
├── pkg/
│   ├── logger/
│   │   ├── interface.go        # 日志接口
│   │   ├── logrus.go          # Logrus实现
│   │   └── zap.go             # Zap实现
│   └── errors/
│       └── errors.go          # 错误类型定义
├── api/
│   └── openapi.yaml           # OpenAPI规范
├── configs/
│   ├── config.yaml            # 默认配置
│   ├── config.prod.yaml       # 生产配置
│   └── config.dev.yaml        # 开发配置
├── scripts/
│   ├── build.sh               # 构建脚本
│   ├── test.sh                # 测试脚本
│   └── deploy.sh              # 部署脚本
├── test/
│   ├── integration/           # 集成测试
│   ├── fixtures/              # 测试数据
│   └── mocks/                 # Mock对象
├── docs/                      # 文档目录
├── .github/
│   └── workflows/
│       ├── ci.yml             # CI流程
│       └── release.yml        # 发布流程
├── go.mod
├── go.sum
├── Makefile
├── Dockerfile
├── docker-compose.yml
├── .golangci.yml             # Lint配置
├── .env.example              # 环境变量示例
└── README.md
```

### 3.2 核心依赖规划

```go
// go.mod
module github.com/your-org/feishu-mcp-go

go 1.21

require (
    // Web框架和HTTP
    github.com/gin-gonic/gin v1.9.1
    github.com/go-resty/resty/v2 v2.7.0
    
    // 配置管理
    github.com/spf13/viper v1.16.0
    github.com/spf13/cobra v1.7.0
    
    // 日志
    github.com/sirupsen/logrus v1.9.3
    go.uber.org/zap v1.24.0
    
    // JSON处理
    github.com/tidwall/gjson v1.14.4
    github.com/json-iterator/go v1.1.12
    
    // 缓存
    github.com/allegro/bigcache/v3 v3.1.0
    github.com/patrickmn/go-cache v2.1.0+incompatible
    github.com/go-redis/redis/v8 v8.11.5
    
    // 监控和指标
    github.com/prometheus/client_golang v1.16.0
    
    // 测试
    github.com/stretchr/testify v1.8.4
    github.com/golang/mock v1.6.0
    
    // 工具库
    github.com/google/uuid v1.3.0
    golang.org/x/time v0.3.0
    golang.org/x/sync v0.3.0
)
```

## 4. 实现优先级和开发路径

### 4.1 第一阶段：基础架构 (Week 1-2)

**目标**: 搭建完整的项目骨架和基础功能

**任务清单**:
1. **项目初始化**
   - 创建标准Go项目结构
   - 配置Go modules和基础依赖
   - 设置开发工具链 (golangci-lint, pre-commit hooks)

2. **配置管理系统**
   - 实现多环境配置支持 (dev/prod/test)
   - 环境变量和配置文件融合
   - 配置验证和默认值设置

3. **日志系统**
   - 结构化日志实现
   - 日志级别和格式配置
   - 性能监控日志

4. **MCP协议基础实现**
   - 基本消息类型定义
   - Stdio传输层实现
   - 简单的协议处理框架

**验收标准**:
- [ ] 项目可以正常启动并监听stdio
- [ ] 配置系统工作正常
- [ ] 基础的MCP握手协议可以工作
- [ ] 日志输出格式正确

### 4.2 第二阶段：飞书API集成 (Week 3-4)

**目标**: 完成飞书API客户端和认证系统

**任务清单**:
1. **飞书API客户端**
   - HTTP客户端封装 (基于resty)
   - 请求/响应模型定义
   - 错误处理和重试机制

2. **认证和Token管理**
   - App Access Token获取和刷新
   - Token缓存机制
   - 并发安全的Token管理

3. **基础API实现**
   - 文档基础操作 (创建、获取、删除)
   - 简单的块操作
   - 文件夹基础操作

4. **缓存系统**
   - 内存缓存实现
   - 缓存键策略和过期管理
   - 缓存预热机制

**验收标准**:
- [ ] 可以成功获取飞书App Access Token
- [ ] 基础的文档操作API工作正常
- [ ] 缓存系统有效减少API调用
- [ ] 错误处理机制完善

### 4.3 第三阶段：核心工具实现 (Week 5-7)

**目标**: 实现所有核心MCP工具

**任务清单**:
1. **文档管理工具**
   - create_feishu_document
   - get_feishu_document_info
   - get_feishu_document_content
   - search_feishu_documents

2. **内容操作工具**
   - create_feishu_text_block
   - create_feishu_code_block  
   - create_feishu_heading_block
   - create_feishu_list_block
   - update_feishu_block_text
   - delete_feishu_document_blocks

3. **批量操作优化**
   - batch_create_feishu_blocks
   - 批量处理逻辑
   - 请求合并优化

4. **文件夹管理工具**
   - get_feishu_root_folder_info
   - get_feishu_folder_files
   - create_feishu_folder

**验收标准**:
- [ ] 所有核心工具都能正常工作
- [ ] 工具参数验证完善
- [ ] 错误信息用户友好
- [ ] 批量操作性能优化有效

### 4.4 第四阶段：高级功能和优化 (Week 8-9) 

**目标**: 完善高级功能和性能优化

**任务清单**:
1. **工具功能完善**
   - convert_feishu_wiki_to_document_id
   - get_feishu_image_resource
   - get_feishu_block_content
   - get_feishu_document_blocks

2. **性能优化**
   - 连接池优化
   - 并发控制
   - 请求去重
   - 响应压缩

3. **HTTP/SSE传输层**
   - RESTful API实现
   - Server-Sent Events支持
   - WebSocket支持 (可选)

4. **监控和指标**
   - Prometheus指标集成
   - 健康检查端点
   - 性能监控

**验收标准**:
- [ ] 所有工具功能完整
- [ ] HTTP模式正常工作
- [ ] 性能指标达标
- [ ] 监控系统完善

### 4.5 第五阶段：测试和发布 (Week 10)

**目标**: 完善测试和准备发布

**任务清单**:
1. **测试完善**
   - 单元测试覆盖率 > 80%
   - 集成测试
   - 端到端测试
   - 性能测试

2. **文档编写**
   - API文档
   - 使用指南  
   - 部署文档
   - 故障排查指南

3. **CI/CD流程**
   - GitHub Actions配置
   - 自动化测试
   - 自动化构建和发布
   - Docker镜像构建

4. **发布准备**
   - 版本管理
   - 发布说明
   - 兼容性测试

**验收标准**:
- [ ] 测试覆盖率达标
- [ ] 文档完整
- [ ] CI/CD流程正常
- [ ] 可以成功发布

## 5. 关键技术决策和理由

### 5.1 架构模式选择

**选择**: 分层架构 + 依赖注入
**理由**: 
- 清晰的职责分离
- 便于测试和维护
- 支持功能扩展

### 5.2 并发模型

**选择**: Goroutine + Channel + sync包
**理由**:
- Go原生并发优势
- 高性能和低资源消耗
- 简单的并发控制

### 5.3 缓存策略

**选择**: 多层缓存 (L1内存 + L2可选Redis)
**理由**:
- 平衡性能和资源使用
- 灵活的缓存策略
- 支持集群部署

### 5.4 错误处理

**选择**: 结构化错误 + 统一错误处理
**理由**:
- 用户友好的错误信息
- 便于调试和监控
- 标准化的错误格式

## 6. 性能目标和监控指标

### 6.1 性能目标

- **响应时间**: 95%的请求在2秒内完成
- **并发处理**: 支持100个并发请求
- **内存使用**: 正常运行占用 < 100MB
- **CPU使用**: 正常负载下 < 10%
- **缓存命中率**: > 80%

### 6.2 监控指标

```go
// 关键指标定义
type Metrics struct {
    // 请求指标
    RequestsTotal     prometheus.CounterVec   // 总请求数
    RequestDuration   prometheus.HistogramVec // 请求耗时
    RequestsInFlight  prometheus.GaugeVec     // 进行中请求数
    
    // 错误指标  
    ErrorsTotal       prometheus.CounterVec   // 错误总数
    ErrorRate         prometheus.GaugeVec     // 错误率
    
    // 资源指标
    MemoryUsage       prometheus.GaugeVec     // 内存使用
    CPUUsage          prometheus.GaugeVec     // CPU使用
    GoroutineCount    prometheus.GaugeVec     // Goroutine数量
    
    // 缓存指标
    CacheHits         prometheus.CounterVec   // 缓存命中
    CacheMisses       prometheus.CounterVec   // 缓存未命中
    CacheSize         prometheus.GaugeVec     // 缓存大小
    
    // 飞书API指标
    FeishuAPICall     prometheus.CounterVec   // 飞书API调用数
    FeishuAPILatency  prometheus.HistogramVec // 飞书API延迟
    TokenRefreshCount prometheus.CounterVec   // Token刷新次数
}
```

---

这个详细的技术实现方案为开发Go版本的飞书MCP提供了全面的指导，包含了架构设计、开发规划、性能目标等各个方面的考虑。可以作为项目开发的蓝图和检查清单使用。 