# 飞书MCP Go项目工具列表

## 总览
本项目提供了25个MCP工具，覆盖了飞书API的各个方面。

## 工具分类

### 📄 文档管理工具 (5个)
1. **create_feishu_document** - 创建新的飞书文档
2. **get_feishu_document_info** - 获取飞书文档的基本信息
3. **get_feishu_document_content** - 获取飞书文档的纯文本内容
4. **get_feishu_document_blocks** - 获取飞书文档的块结构信息
5. **search_feishu_documents** - 在飞书中搜索文档

### 🔧 内容操作工具 (8个)
1. **get_feishu_block_content** - 获取飞书文档中指定块的详细内容
2. **update_feishu_block_text** - 更新飞书文档中指定块的文本内容和样式
3. **batch_create_feishu_blocks** - 批量创建多个飞书文档块
4. **create_feishu_text_block** - 创建新的文本块
5. **create_feishu_code_block** - 创建新的代码块
6. **create_feishu_heading_block** - 创建新的标题块
7. **create_feishu_list_block** - 创建新的列表块
8. **delete_feishu_document_blocks** - 删除飞书文档中的一个或多个连续块

### 📁 文件夹管理工具 (7个)
1. **get_feishu_root_folder_info** - 获取飞书云盘根文件夹信息
2. **get_feishu_root_folder_meta** - 获取云空间根文件夹元数据信息（使用新的API端点）
3. **get_feishu_folder_files** - 获取指定文件夹中的文件和子文件夹列表
4. **create_feishu_folder** - 在指定父文件夹中创建新文件夹
5. **get_feishu_drive_files_with_meta** - 获取云空间目录下所有文件的详细元数据信息
6. **get_feishu_drive_meta** - 获取云空间目录/文件的元数据信息
7. **get_all_feishu_drive_files** - 获取指定目录下所有文件（支持分页和数量限制）

### 📚 知识库工具 (4个)
1. **get_feishu_wiki_spaces** - 获取飞书知识库空间列表
2. **get_feishu_wiki_nodes** - 获取飞书知识库空间下的节点列表
3. **get_feishu_wiki_node_content** - 获取飞书知识库节点的内容
4. **get_feishu_wiki_node_meta** - 获取飞书知识库节点的元信息

### 🛠️ 工具功能 (2个)
1. **convert_feishu_wiki_to_document_id** - 将飞书Wiki链接转换为文档ID
2. **get_feishu_image_resource** - 获取飞书图片资源信息和内容

## 使用统计
- 总工具数量：25个
- 文档管理：5个
- 内容操作：8个
- 文件夹管理：7个
- 知识库：4个
- 工具功能：2个

## 更新日志
- 2024-12-25：新增 `get_feishu_root_folder_meta` 工具，使用新的API端点获取根文件夹元数据
- 2024-12-25：新增云空间相关的3个工具：`get_feishu_drive_files_with_meta`、`get_feishu_drive_meta`、`get_all_feishu_drive_files` 