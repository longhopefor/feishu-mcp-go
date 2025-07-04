# 飞书知识库API使用指南

## 概述

本指南介绍如何使用飞书MCP Go项目中的知识库相关API接口。知识库功能让你可以管理和访问飞书中的wiki内容，包括获取空间列表、节点信息、内容和元数据等。

## 支持的知识库工具

### 1. 获取知识库空间列表 (`get_feishu_wiki_spaces`)
获取用户有权限访问的知识库空间列表。

**参数：**
- `pageSize` (可选): 分页大小，默认20
- `pageToken` (可选): 分页令牌

**返回值：**
- 知识库空间列表
- 分页信息

### 2. 获取知识库节点列表 (`get_feishu_wiki_nodes`)
获取指定知识库空间中的节点列表。

**参数：**
- `spaceID` (必需): 知识库空间ID
- `pageSize` (可选): 分页大小，默认20
- `pageToken` (可选): 分页令牌
- `parentNode` (可选): 父节点Token，用于获取子节点

**返回值：**
- 知识库节点列表
- 分页信息

### 3. 获取知识库节点内容 (`get_feishu_wiki_node_content`)
获取指定知识库节点的内容。

**参数：**
- `spaceID` (必需): 知识库空间ID
- `nodeToken` (必需): 节点Token
- `lang` (可选): 语言代码，0为中文，1为英文

**返回值：**
- 节点内容
- 内容类型

### 4. 获取知识库节点元信息 (`get_feishu_wiki_node_meta`)
获取指定知识库节点的元信息。

**参数：**
- `spaceID` (必需): 知识库空间ID
- `nodeToken` (必需): 节点Token

**返回值：**
- 节点标题
- 节点类型
- 创建时间
- 最后修改时间
- 是否有子节点等信息

## 如何获取知识库参数

### 1. 获取知识库空间ID (Space ID)

#### 方法1：通过浏览器URL
1. 打开飞书，进入知识库
2. 在浏览器地址栏中，URL格式为：`https://xxx.feishu.cn/wiki/空间ID`
3. 复制URL中的空间ID部分

#### 方法2：通过API调用
```bash
# 获取知识库空间列表
./scripts/debug.sh wiki-spaces
```

### 2. 获取节点Token (Node Token)

#### 方法1：通过浏览器URL
1. 打开飞书知识库中的具体页面
2. 在浏览器地址栏中，URL格式为：`https://xxx.feishu.cn/wiki/空间ID?table=节点Token`
3. 复制URL中的节点Token部分

#### 方法2：通过API调用
```bash
# 获取知识库节点列表
./scripts/debug.sh wiki-nodes
```

## 配置指南

### 1. 配置debug.env文件

复制 `debug.env.example` 为 `debug.env`，并设置以下参数：

```bash
# 基础配置
FEISHU_APP_ID=your_app_id_here
FEISHU_APP_SECRET=your_app_secret_here

# 知识库测试配置
TEST_WIKI_SPACE_ID=your_wiki_space_id_here
TEST_WIKI_NODE_TOKEN=your_wiki_node_token_here
```

### 2. 权限配置

确保你的飞书应用具有以下权限：
- `wiki:wiki:read` - 读取知识库内容
- `wiki:wiki:readonly` - 只读访问知识库
- `wiki:node:read` - 读取知识库节点

在飞书开放平台中配置权限：
1. 进入飞书开放平台控制台
2. 选择你的应用
3. 在"权限管理"中添加知识库相关权限
4. 发布应用版本

## 使用示例

### 1. 获取知识库空间列表
```bash
# 获取默认的知识库空间列表
./scripts/debug.sh wiki-spaces

# 获取指定分页大小的空间列表
./scripts/debug.sh wiki-spaces 10

# 获取指定分页令牌的空间列表
./scripts/debug.sh wiki-spaces 20 "your_page_token"
```

### 2. 获取知识库节点列表
```bash
# 获取指定空间的节点列表
./scripts/debug.sh wiki-nodes

# 获取指定分页大小的节点列表
./scripts/debug.sh wiki-nodes 10

# 获取指定父节点下的子节点
./scripts/debug.sh wiki-nodes 20 "parent_node_token"
```

### 3. 获取知识库节点内容
```bash
# 获取节点内容（中文）
./scripts/debug.sh wiki-content

# 获取节点内容（英文）
./scripts/debug.sh wiki-content 1
```

### 4. 获取知识库节点元信息
```bash
# 获取节点元信息
./scripts/debug.sh wiki-meta
```

### 5. 运行所有测试
```bash
# 运行包含知识库在内的所有API测试
./scripts/debug.sh all
```

## API端点说明

### 知识库空间相关端点
- `/wiki/v2/spaces` - 获取知识库空间列表
- `/wiki/v1/spaces` - 备用端点

### 知识库节点相关端点
- `/wiki/v2/spaces/{space_id}/nodes` - 获取节点列表
- `/wiki/v2/spaces/{space_id}/nodes/{node_token}` - 获取节点信息
- `/wiki/v2/spaces/{space_id}/nodes/{node_token}/content` - 获取节点内容

## 错误处理

### 常见错误及解决方案

1. **403 Forbidden**
   - 检查应用权限配置
   - 确认应用已发布并获得授权

2. **404 Not Found**
   - 检查空间ID和节点Token是否正确
   - 确认知识库内容是否存在

3. **400 Bad Request**
   - 检查请求参数格式
   - 确认必需参数是否提供

4. **401 Unauthorized**
   - 检查应用ID和密钥是否正确
   - 确认访问令牌是否有效

## 性能优化建议

1. **分页处理**
   - 使用合适的分页大小（建议10-50）
   - 通过pageToken进行分页遍历

2. **缓存策略**
   - 对元信息进行适当缓存
   - 注意内容更新时的缓存失效

3. **批量操作**
   - 避免频繁的单个请求
   - 使用分页一次性获取多个节点

## 完整工作流程示例

```bash
# 1. 检查基础连接
./scripts/debug.sh token
./scripts/debug.sh health

# 2. 获取知识库空间列表
./scripts/debug.sh wiki-spaces

# 3. 设置配置文件中的TEST_WIKI_SPACE_ID

# 4. 获取节点列表
./scripts/debug.sh wiki-nodes

# 5. 设置配置文件中的TEST_WIKI_NODE_TOKEN

# 6. 获取节点内容和元信息
./scripts/debug.sh wiki-content
./scripts/debug.sh wiki-meta

# 7. 运行完整测试
./scripts/debug.sh all
```

## 故障排除

### 调试模式
启用调试模式以获取详细的请求和响应信息：

```bash
# 在debug.env中设置
DEBUG_MODE=true
VERBOSE_MODE=true
```

### 日志分析
检查日志输出中的以下信息：
- HTTP状态码
- 错误响应详情
- 请求参数验证结果
- API端点尝试顺序

### 常见问题检查清单
- [ ] 飞书应用ID和密钥是否正确
- [ ] 知识库权限是否已配置
- [ ] 应用是否已发布
- [ ] 知识库空间ID是否有效
- [ ] 节点Token是否存在
- [ ] 网络连接是否正常

## 最佳实践

1. **权限最小化原则**
   - 只申请必要的知识库权限
   - 定期审查权限使用情况

2. **错误处理**
   - 实现适当的重试机制
   - 提供友好的错误信息

3. **数据保护**
   - 不要在日志中记录敏感的Token信息
   - 安全存储配置文件

4. **性能监控**
   - 监控API调用频率
   - 跟踪响应时间和错误率

## 总结

通过本指南，你已经了解了如何使用飞书知识库API。这些工具可以帮助你构建强大的知识管理应用，实现对飞书知识库内容的程序化访问和管理。

如果遇到问题，请参考故障排除部分，或查看项目的完整API文档。 