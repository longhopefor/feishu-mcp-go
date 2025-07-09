# 飞书MCP本地使用指南

## 🎯 完整使用流程

本指南将带您从零开始在本地部署和使用飞书MCP服务器，并与AI编程工具（如Cursor、Windsurf、Cline）集成。

---

## 📋 前提条件

### 系统要求
- **Go 1.21+** - [下载安装](https://golang.org/dl/)
- **Git** - 用于克隆代码
- **飞书开发者账号** - 用于创建应用获取凭证

### 支持的AI工具
- ✅ **Cursor** - VS Code分支的AI编程工具
- ✅ **Windsurf** - Codeium的AI编程工具  
- ✅ **Cline** (原Claude Dev) - VS Code扩展
- ✅ **其他支持MCP协议的工具**

---

## 🔑 第一步：获取飞书应用凭证

### 1.1 创建飞书应用

1. 访问 [飞书开放平台](https://open.feishu.cn/)
2. 登录并进入 **开发者后台**
3. 点击 **创建应用** → 选择 **自建应用**
4. 填写应用信息：
   - **应用名称**: 如 "MCP文档助手"
   - **应用描述**: 描述用途
   - **应用类型**: 选择 **企业自建应用**

### 1.2 配置应用权限

在应用管理页面，进入 **权限管理**，添加以下权限：

**📄 文档权限**:
- `docx:document` - 查看、编辑、创建新版云文档
- `docx:document:readonly` - 查看云文档  
- `drive:file` - 查看、下载、分享云空间中的文件

**🔍 搜索权限**:
- `search:message` - 搜索消息
- `search:data_source` - 获取数据源信息

**📁 云空间权限**:
- `drive:drive` - 获取云空间信息
- `drive:file` - 文件操作权限

### 1.3 获取凭证信息

在应用的 **凭证与基础信息** 页面，复制：
- **App ID** (cli_开头)
- **App Secret** (长字符串)

⚠️ **重要**: 妥善保管这些凭证，不要泄露！

### 1.4 发布应用

1. 进入 **版本管理与发布**
2. 创建版本并提交审核
3. 审核通过后，在企业内发布应用

---

## 💻 第二步：安装和配置项目

### 2.1 克隆项目

```bash
# 克隆代码仓库
git clone https://github.com/longhopefor/feishu-mcp-go.git
cd feishu-mcp-go

# 确保Go模块正常
go mod download
go mod tidy
```

### 2.2 配置环境变量（推荐方式）

创建 `.env` 文件：

```bash
# 创建环境配置文件
cat > .env << 'EOF'
# 飞书应用配置（必需）
export FEISHU_MCP_FEISHU_APP_ID="cli_你的应用ID"
export FEISHU_MCP_FEISHU_APP_SECRET="你的应用密钥"

# 可选配置
export FEISHU_MCP_LOG_LEVEL="info"
export FEISHU_MCP_CACHE_ENABLED="true"
export FEISHU_MCP_SERVER_MODE="stdio"
EOF

# 加载环境变量
source .env
```

### 2.3 配置文件方式（备选）

复制并编辑配置文件：

```bash
# 复制配置文件
cp configs/config.yaml configs/my-config.yaml

# 编辑配置文件
cat > configs/my-config.yaml << 'EOF'
feishu:
  app_id: "cli_你的应用ID"
  app_secret: "你的应用密钥"
  base_url: "https://open.feishu.cn/open-apis"
  timeout: 30

log:
  level: "info"
  format: "json"

cache:
  enabled: true
  ttl: 3600

server:
  mode: "stdio"  # 重要：必须设置为stdio模式
  port: 3333
EOF
```

---

## 🚀 第三步：构建和测试

### 3.1 构建项目

```bash
# 使用Makefile构建（推荐）
make build

# 或手动构建
go build -o build/feishu-mcp ./cmd/server/main.go
```

### 3.2 测试连接

```bash
# 测试API连接
source .env  # 加载环境变量
go run ./cmd/debug/main.go -app-id="$FEISHU_MCP_FEISHU_APP_ID" -app-secret="$FEISHU_MCP_FEISHU_APP_SECRET" -action=token

# 应该看到类似输出：
# ✅ 成功获取访问令牌！
# ✅ 令牌验证成功，可以正常访问飞书API！

# 测试MCP协议连接
./test_mcp_connection.sh

# 应该看到输出包含：
# ✅ initialize响应 - 协议版本和服务器能力
# ✅ tools/list响应 - 21个可用工具列表
# ✅ ping响应成功
```

### 3.3 测试工具功能

```bash
# 运行全面测试
chmod +x comprehensive_api_test.sh
./comprehensive_api_test.sh

# 测试特定功能
go run ./cmd/debug/main.go -app-id="$FEISHU_MCP_FEISHU_APP_ID" -app-secret="$FEISHU_MCP_FEISHU_APP_SECRET" -action=health
```

---

## 🔌 第四步：与AI工具集成

### 4.1 Cursor 集成

1. **打开Cursor设置**：
   - 按 `Cmd/Ctrl + ,` 打开设置
   - 搜索 "MCP" 或进入 Extensions → MCP

2. **配置MCP服务器**：
   ```json
   {
     "mcpServers": {
       "feishu": {
         "command": "/Users/你的用户名/path/to/feishu-mcp-go/build/feishu-mcp",
         "args": ["--stdio"],
         "env": {
           "FEISHU_MCP_FEISHU_APP_ID": "cli_你的应用ID",
           "FEISHU_MCP_FEISHU_APP_SECRET": "你的应用密钥",
           "FEISHU_MCP_LOG_LEVEL": "info"
         }
       }
     }
   }
   ```

### 4.2 Windsurf 集成

1. **打开设置文件**：
   - 找到Windsurf的MCP配置文件（通常在用户配置目录）

2. **添加配置**：
   ```json
   {
     "mcpServers": {
       "feishu-mcp": {
         "command": "/path/to/feishu-mcp-go/build/feishu-mcp",
         "args": ["--stdio"],
         "env": {
           "FEISHU_MCP_FEISHU_APP_ID": "cli_你的应用ID",
           "FEISHU_MCP_FEISHU_APP_SECRET": "你的应用密钥"
         }
       }
     }
   }
   ```

### 4.3 Cline (Claude Dev) 集成

1. **安装Cline扩展**（在VS Code中）
2. **配置MCP服务器**：
   - 打开Cline扩展设置
   - 添加MCP服务器配置
   - 使用与Cursor相同的配置格式

### 4.4 配置路径说明

⚠️ **重要**: 确保路径正确！

```bash
# 获取绝对路径
pwd  # 显示当前目录完整路径
# 例如：/Users/wanglong/mcp/code-demo

# 在配置中使用完整路径：
# "/Users/wanglong/mcp/code-demo/build/feishu-mcp"
```

---

## 🎮 第五步：实际使用示例

### 5.1 重启AI工具

配置完成后，**重启AI工具**让MCP配置生效。

### 5.2 验证集成

在AI工具中，应该能看到飞书MCP工具可用。试试以下对话：

```
你好！我想测试飞书MCP功能，请帮我：
1. 获取我的飞书文档列表
2. 创建一个新文档
```

### 5.3 常用功能示例

**📝 创建文档**:
```
请在我的飞书云盘中创建一个名为"项目需求文档"的新文档
```

**📋 获取文档内容**:
```
请获取文档ID为"xxx"的内容，并总结要点
```

**✏️ 编辑文档**:
```
请在现有文档中添加一个标题"项目背景"，并在下面添加项目背景的详细描述
```

**🔍 搜索文档**:
```
请搜索包含"技术方案"关键词的飞书文档
```

---

## 🛠️ 故障排除

### 常见问题

**❌ 问题1**: "获取令牌失败"
```bash
# 解决方案：
1. 检查APP_ID和APP_SECRET是否正确
2. 确认应用已发布且权限充足
3. 检查网络连接

# 调试命令：
go run ./cmd/debug/main.go -app-id="你的ID" -app-secret="你的密钥" -action=token -debug
```

**❌ 问题2**: "AI工具找不到MCP服务器"
```bash
# 解决方案：
1. 检查可执行文件路径是否正确
2. 确认文件有执行权限：chmod +x build/feishu-mcp
3. 使用绝对路径而不是相对路径
4. 重启AI工具
```

**❌ 问题3**: "权限不足"
```bash
# 解决方案：
1. 检查飞书应用权限配置
2. 确认应用已正确发布
3. 检查用户是否在应用可见范围内
```

**❌ 问题4**: "JSON RPC Parse error"
```bash
# 原因：
直接向MCP服务器发送随意输入会导致解析错误，这是正常现象

# 解决方案：
1. 使用正确的MCP协议初始化序列
2. 运行测试脚本：./test_mcp_connection.sh
3. 通过AI工具连接，而不是直接手动输入
```

### 调试模式

启用详细日志：

```bash
# 方式1：环境变量
export FEISHU_MCP_LOG_LEVEL="debug"

# 方式2：命令行参数
./build/feishu-mcp --log-level=debug --stdio

# 方式3：使用debug工具
go run ./cmd/debug/main.go -debug -verbose -action=health
```

### 测试脚本

运行完整测试确保功能正常：

```bash
# 加载环境变量
source .env

# 运行全面测试
./comprehensive_api_test.sh

# 检查特定功能
go run ./cmd/debug/main.go -action=health -debug
```

---

## 📚 高级用法

### 自定义配置

```bash
# 使用自定义配置文件
./build/feishu-mcp --config=configs/my-config.yaml --stdio

# 覆盖特定参数
./build/feishu-mcp --feishu-app-id="cli_xxx" --log-level=debug --stdio
```

### 开发模式

```bash
# 开发模式运行（实时重载）
make dev

# 或直接运行
go run ./cmd/server/main.go --log-level=debug --stdio
```

### 性能监控

```bash
# 启用性能日志
export FEISHU_MCP_LOG_LEVEL="debug"
./build/feishu-mcp --stdio 2>&1 | tee feishu-mcp.log

# 查看日志
tail -f feishu-mcp.log
```

---

## 🎉 恭喜！

如果您完成了以上步骤，现在应该可以：

✅ **在AI工具中使用飞书MCP功能**  
✅ **创建、编辑、搜索飞书文档**  
✅ **管理文件夹和内容块**  
✅ **获得强大的文档操作能力**  

享受AI驱动的飞书文档操作体验吧！

---

## 📞 获取帮助

- **项目仓库**: [GitHub](https://github.com/longhopefor/feishu-mcp-go)
- **问题反馈**: [Issues](https://github.com/longhopefor/feishu-mcp-go/issues)
- **API文档**: [飞书开放平台](https://open.feishu.cn/document/)
- **MCP协议**: [Model Context Protocol](https://modelcontextprotocol.io/)

**🚀 开始您的AI驱动文档操作之旅！** 