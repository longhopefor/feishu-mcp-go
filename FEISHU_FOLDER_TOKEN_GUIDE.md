# 飞书 Folder Token 获取指南

## 🎯 什么是 Folder Token？

Folder Token（文件夹令牌）是飞书云空间中每个文件夹的唯一标识符，用于通过API访问和操作特定文件夹。当您需要在某个文件夹中创建文档、获取文件列表时，都需要提供对应的folder token。

## 📋 获取方法

### 方法1: 通过浏览器URL获取（最简单）

1. **打开飞书云空间**
   - 登录飞书客户端或网页版
   - 进入"云空间"模块

2. **导航到目标文件夹**
   - 找到您想要操作的文件夹
   - 双击进入该文件夹

3. **从URL中提取Token**
   - 查看浏览器地址栏
   - URL格式通常为：`https://xxx.feishu.cn/drive/folder/{folder_token}`
   - 复制`{folder_token}`部分

**示例URL分析**：
```
https://bytedance.feishu.cn/drive/folder/FldcnH5V8loCfBde4HHWc4Mygnb
                                        ↑________________________↑
                                              这部分就是folder_token
```

### 方法2: 通过API获取（推荐用于程序）

我们的项目中包含了获取文件夹信息的工具，您可以通过以下方式使用：

#### 2.1 获取根目录信息
```bash
# 运行我们项目的调试工具
./scripts/debug.sh folder-info
```

#### 2.2 通过文件夹名称搜索（未来功能）
目前我们项目中暂未实现文件夹搜索功能，但这是一个很有用的功能，我们可以添加。

### 方法3: 通过飞书开放平台测试（适合开发者）

1. **访问飞书开放平台**
   - 打开：https://open.feishu.cn/
   - 登录您的开发者账号

2. **进入API调试台**
   - 选择"工具与资源" > "API调试台"
   - 找到"云文档"相关的API

3. **使用文件夹列表API**
   - 选择"获取文件夹下的文档列表"API
   - 从根目录开始逐级查找您需要的文件夹

## 🔧 在我们项目中使用 Folder Token

### 配置调试环境
```bash
# 1. 复制配置文件
cp debug.env.example debug.env

# 2. 编辑配置文件，填入您获取的folder token
vim debug.env
```

在`debug.env`中设置：
```bash
# 测试用的文档和文件夹
TEST_FOLDER_TOKEN=您获取的folder_token_在这里
TEST_DOCUMENT_ID=your_test_document_id_here
TEST_DOCUMENT_TITLE=API测试文档
```

### 测试文件夹访问
```bash
# 测试创建文档到指定文件夹
./scripts/debug.sh create-doc
```

## ⚠️ 重要注意事项

### 1. 权限要求
- 您的飞书应用必须有访问对应文件夹的权限
- 如果是企业自建应用，需要管理员授权相应的云文档权限
- 个人文件夹和共享文件夹的权限处理可能不同

### 2. Token格式特征
飞书的folder token通常具有以下特征：
- 长度约为27个字符
- 由字母和数字组成
- 格式类似：`FldcnH5V8loCfBde4HHWc4Mygnb`

### 3. 常见错误
- **folder not found**: folder token不存在或无权限访问
- **invalid token**: folder token格式错误
- **permission denied**: 应用没有访问该文件夹的权限

## 🛠️ 权限配置指南

### 1. 飞书开放平台权限设置
1. 登录 https://open.feishu.cn/
2. 进入您的应用管理页面
3. 点击"权限管理"
4. 搜索并开通以下权限：
   - `drive:drive` (查看云空间文件)
   - `docx:document` (创建、编辑云文档)
   - `drive:drive:readonly` (只读访问云空间)

### 2. 应用发布
开通权限后，需要：
1. 创建应用版本
2. 申请线上发布
3. 等待审核通过

## 🎯 实用示例

### 获取个人根目录的folder token
```bash
# 方法1: 直接查看个人云空间URL
# 进入飞书 > 云空间 > 我的空间
# URL通常为: https://xxx.feishu.cn/drive/folder/{root_folder_token}

# 方法2: 使用我们的调试工具
./scripts/debug.sh health  # 验证基础连接
```

### 创建测试文档到指定文件夹
```bash
# 设置好TEST_FOLDER_TOKEN后
./scripts/debug.sh create-doc
```

## 📚 API参考

### 相关API端点
- 获取文件夹信息：`/drive/v1/folders/{folder_token}`
- 获取文件夹文件列表：`/drive/v1/files?folder_token={folder_token}`
- 创建文档到文件夹：`/docx/v1/documents` (Body中指定folder_token)

### 请求示例
```bash
# 获取文件夹信息
curl -X GET "https://open.feishu.cn/open-apis/drive/v1/folders/{folder_token}" \
  -H "Authorization: Bearer {your_access_token}"
```

## 🔍 故障排除

### 常见问题和解决方案

1. **Q: 无法访问文件夹**
   - A: 检查应用是否有相应权限，确认权限已发布上线

2. **Q: folder token格式看起来不对**
   - A: 确保复制的是完整的token，没有包含额外的字符

3. **Q: 在共享文件夹中无法创建文档**
   - A: 检查您在该共享文件夹中的权限级别，需要编辑权限

### 调试建议
1. 先测试简单的API调用，如健康检查
2. 使用个人文件夹进行初始测试
3. 检查应用权限配置是否完整
4. 查看详细的错误日志进行诊断

## 💡 最佳实践

1. **测试环境**：先在测试文件夹中验证功能
2. **权限最小化**：只申请必要的权限
3. **错误处理**：在代码中添加完善的错误处理逻辑
4. **缓存机制**：频繁使用的folder token可以考虑缓存

---

通过以上方法，您应该能够成功获取到需要的飞书folder token。如果仍有问题，请查看我们项目的调试日志，或者检查飞书开放平台的API文档。 