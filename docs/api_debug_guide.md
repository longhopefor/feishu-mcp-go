# 📋 飞书MCP Go API调试指南

## 🎯 概述

本指南将帮助你完成飞书MCP Go项目的API验证和调试工作。我们提供了完整的调试工具链，包括：

- **命令行调试工具** - 直接测试API调用
- **HTTP请求监控** - 实时查看请求和响应
- **健康检查系统** - 验证API连接状态
- **自动化测试脚本** - 批量验证所有API

## 🚀 快速开始

### 1. 设置配置文件

首先复制配置文件模板：

```bash
cp debug.env.example debug.env
```

编辑 `debug.env` 文件，填入你的飞书应用信息：

```bash
# 飞书应用配置
FEISHU_APP_ID=cli_xxxxxxxxxxxxxxxxx
FEISHU_APP_SECRET=xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
FEISHU_BASE_URL=https://open.feishu.cn/open-apis

# 测试用的文档和文件夹（可选）
TEST_FOLDER_TOKEN=fldbxxxxxxxxxxxxxxxxxxxxxxxxx
TEST_DOCUMENT_ID=docxxxxxxxxxxxxxxxxxxxxxxxxx
TEST_DOCUMENT_TITLE=API测试文档

# 调试配置
DEBUG_MODE=true
VERBOSE_MODE=true
```

### 2. 基础健康检查

运行基础健康检查来验证配置：

```bash
./scripts/debug.sh health
```

期望输出：
```
🏥 执行API健康检查...
1️⃣ 检查基础连接...
   ✅ 基础连接正常，状态码: 200
2️⃣ 检查Token获取...
   ✅ Token获取成功，长度: 3456
3️⃣ 检查API调用...
   ✅ API调用成功
🎉 健康检查完成!
```

### 3. 运行完整测试

运行所有API测试：

```bash
./scripts/debug.sh all
```

## 🔧 调试工具详解

### 命令行调试工具

直接使用Go程序进行调试：

```bash
# 构建调试工具
go build -o build/debug cmd/debug/main.go

# 基础用法
./build/debug -app-id="your_app_id" -app-secret="your_app_secret" -action=token

# 启用调试模式
./build/debug -app-id="your_app_id" -app-secret="your_app_secret" -action=health -debug -verbose
```

### 使用调试脚本

推荐使用便捷的调试脚本：

```bash
# 测试获取访问令牌
./scripts/debug.sh token

# 执行健康检查
./scripts/debug.sh health

# 验证API端点
./scripts/debug.sh validate

# 搜索文档
./scripts/debug.sh search "关键词"

# 运行所有测试
./scripts/debug.sh all
```

## 📊 支持的调试操作

### 1. 基础操作

| 操作 | 命令 | 说明 |
|------|------|------|
| 获取令牌 | `./scripts/debug.sh token` | 测试访问令牌获取 |
| 健康检查 | `./scripts/debug.sh health` | 执行系统健康检查 |
| 端点验证 | `./scripts/debug.sh validate` | 验证API端点可用性 |

### 2. 文档操作

| 操作 | 命令 | 说明 |
|------|------|------|
| 创建文档 | `./scripts/debug.sh create-doc` | 创建测试文档 |
| 获取文档 | `./scripts/debug.sh get-doc` | 获取文档基本信息 |
| 获取内容 | `./scripts/debug.sh get-content` | 获取文档纯文本内容 |
| 获取块 | `./scripts/debug.sh get-blocks` | 获取文档块结构 |
| 搜索文档 | `./scripts/debug.sh search` | 搜索文档 |

### 3. 综合测试

| 操作 | 命令 | 说明 |
|------|------|------|
| 全部测试 | `./scripts/debug.sh all` | 运行所有API测试 |

## 🔍 HTTP请求监控

### 启用监控

在调试模式下，你可以看到详细的HTTP请求信息：

```bash
./scripts/debug.sh token
```

监控输出示例：
```
🌐 HTTP请求监控
📍 URL: POST https://open.feishu.cn/open-apis/auth/v3/tenant_access_token/internal
🔑 Headers: map[Content-Type:[application/json]]
📄 Body: {"app_id":"cli_xxxxx","app_secret":"xxxxx"}
⏰ 时间: 2024-01-20 10:30:45
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

📥 HTTP响应监控
📊 状态码: 200
⏱️ 耗时: 234ms
📏 响应大小: 178 bytes
📄 响应内容: {"code":0,"msg":"ok","tenant_access_token":"t-xxxxx","expire":3600}
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
```

### 错误诊断

当API调用失败时，监控器会提供详细的错误信息和调试建议：

```
🚨 API错误响应
📍 URL: https://open.feishu.cn/open-apis/docx/v1/documents/invalid_id
📊 状态码: 404
📄 错误内容: {"code":404,"msg":"document not found"}
💡 调试建议: 检查API端点是否正确，资源是否存在
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
```

## 🛠️ 常见问题诊断

### 1. Token获取失败

**错误信息**: `Token获取失败: 获取访问令牌失败: invalid app_id or app_secret`

**解决方案**:
1. 检查 `debug.env` 中的 `FEISHU_APP_ID` 和 `FEISHU_APP_SECRET`
2. 确认应用已在飞书开发者控制台激活
3. 检查网络连接和防火墙设置

### 2. 权限不足

**错误信息**: `API调用失败 (code: 403): permission denied`

**解决方案**:
1. 检查应用权限配置
2. 确认已申请相应的API权限
3. 检查用户是否有相应的文档访问权限

### 3. 文档不存在

**错误信息**: `API调用失败 (code: 404): document not found`

**解决方案**:
1. 检查 `TEST_DOCUMENT_ID` 是否正确
2. 确认文档存在且用户有访问权限
3. 尝试创建新的测试文档

### 4. 网络连接问题

**错误信息**: `HTTP请求失败: connection timeout`

**解决方案**:
1. 检查网络连接
2. 确认防火墙和代理设置
3. 尝试不同的 `FEISHU_BASE_URL`

## 🎨 自定义调试

### 直接使用Go程序

```bash
# 启用详细模式
go run cmd/debug/main.go \
  -app-id="your_app_id" \
  -app-secret="your_app_secret" \
  -action=get-doc \
  -doc-id="your_doc_id" \
  -debug \
  -verbose

# 测试不同环境
go run cmd/debug/main.go \
  -app-id="your_app_id" \
  -app-secret="your_app_secret" \
  -base-url="https://open.feishu.cn/open-apis" \
  -action=health
```

### 集成到CI/CD

在自动化流程中使用：

```yaml
# .github/workflows/api-test.yml
name: API Tests
on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v2
    - uses: actions/setup-go@v2
      with:
        go-version: 1.21
    - name: Run API Tests
      env:
        FEISHU_APP_ID: ${{ secrets.FEISHU_APP_ID }}
        FEISHU_APP_SECRET: ${{ secrets.FEISHU_APP_SECRET }}
      run: |
        echo "FEISHU_APP_ID=$FEISHU_APP_ID" > debug.env
        echo "FEISHU_APP_SECRET=$FEISHU_APP_SECRET" >> debug.env
        echo "DEBUG_MODE=true" >> debug.env
        echo "VERBOSE_MODE=false" >> debug.env
        chmod +x scripts/debug.sh
        ./scripts/debug.sh health
```

## 📈 性能监控

### 监控指标

调试工具会自动收集以下指标：

- **响应时间**: 每个API调用的耗时
- **成功率**: API调用的成功/失败比例
- **错误类型**: 具体的错误代码和消息
- **请求大小**: 请求和响应的数据大小

### 优化建议

根据监控结果，你可以：

1. **优化请求频率**: 如果遇到429错误，增加请求间隔
2. **缓存Token**: 重用有效的访问令牌
3. **批量操作**: 合并多个小请求为批量请求
4. **错误重试**: 对临时错误实施重试策略

## 🔗 相关链接

- [飞书开发者文档](https://open.feishu.cn/document/)
- [飞书API权限申请](https://open.feishu.cn/document/home/introduction-to-application-permissions)
- [项目GitHub仓库](https://github.com/longhopefor/feishu-mcp-go)

## 🤝 获取帮助

如果遇到问题，请：

1. 查看本指南的常见问题部分
2. 检查项目的 `cursor_readme.md` 文件
3. 在GitHub仓库提交Issue
4. 联系项目维护者

祝你调试顺利！🎉 