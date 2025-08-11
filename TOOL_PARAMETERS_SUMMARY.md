# 飞书MCP工具参数描述总结

本文档总结了所有飞书MCP工具的参数描述，方便大模型在调用时了解每个工具需要的参数。

## 文档管理工具 (5个)

### 1. create_feishu_document - 创建新的飞书文档
- **title** (string, 必需): 文档标题
- **folderToken** (string, 必需): 父文件夹的token，用于指定文档创建位置

### 2. get_feishu_document_info - 获取飞书文档的基本信息
- **documentId** (string, 必需): 文档ID，用于标识要获取信息的文档

### 3. get_feishu_document_content - 获取飞书文档的纯文本内容
- **documentId** (string, 必需): 文档ID，用于标识要获取内容的文档
- **lang** (number, 可选): 语言类型，0表示中文，1表示英文，默认为0

### 4. get_feishu_document_blocks - 获取飞书文档的块结构信息
- **documentId** (string, 必需): 文档ID，用于标识要获取块结构的文档
- **pageSize** (number, 可选): 每页返回的块数量，默认为50
- **pageToken** (string, 可选): 分页令牌，用于获取下一页数据

### 5. search_feishu_documents - 在飞书中搜索文档
- **query** (string, 必需): 搜索关键词，用于匹配文档标题或内容
- **pageSize** (number, 可选): 每页返回的文档数量，默认为10
- **pageToken** (string, 可选): 分页令牌，用于获取下一页搜索结果

## 内容操作工具 (8个)

### 6. get_feishu_block_content - 获取飞书文档中指定块的详细内容
- **documentId** (string, 必需): 文档ID，用于标识包含目标块的文档
- **blockId** (string, 必需): 块ID，用于标识要获取内容的具体块

### 7. update_feishu_block_text - 更新飞书文档中指定块的文本内容和样式
- **documentId** (string, 必需): 文档ID，用于标识包含目标块的文档
- **blockId** (string, 必需): 块ID，用于标识要更新的具体块
- **content** (string, 必需): 要更新的文本内容
- **style** (string, 可选): 文本样式，可选，如粗体、斜体等

### 8. batch_create_feishu_blocks - 批量创建多个飞书文档块
- **documentId** (string, 必需): 文档ID，用于标识要添加块的文档
- **parentId** (string, 可选): 父块ID，用于指定新块的位置
- **blocks** (string, 必需): 要创建的块列表，JSON格式字符串

### 9. create_feishu_text_block - 创建新的文本块
- **documentId** (string, 必需): 文档ID，用于标识要添加文本块的文档
- **content** (string, 必需): 文本块的内容
- **parentId** (string, 可选): 父块ID，用于指定文本块的位置
- **style** (string, 可选): 文本样式，可选，如粗体、斜体等

### 10. create_feishu_code_block - 创建新的代码块
- **documentId** (string, 必需): 文档ID，用于标识要添加代码块的文档
- **content** (string, 必需): 代码块的内容
- **language** (string, 必需): 编程语言类型，如javascript、python等
- **parentId** (string, 可选): 父块ID，用于指定代码块的位置

### 11. create_feishu_heading_block - 创建新的标题块
- **documentId** (string, 必需): 文档ID，用于标识要添加标题块的文档
- **content** (string, 必需): 标题内容
- **level** (number, 必需): 标题级别，1-3，1为最高级别
- **parentId** (string, 可选): 父块ID，用于指定标题块的位置

### 12. create_feishu_list_block - 创建新的列表块
- **documentId** (string, 必需): 文档ID，用于标识要添加列表块的文档
- **content** (string, 必需): 列表项内容
- **type** (string, 必需): 列表类型，bullet表示无序列表，ordered表示有序列表
- **parentId** (string, 可选): 父块ID，用于指定列表块的位置

### 13. delete_feishu_document_blocks - 删除飞书文档中的一个或多个连续块
- **documentId** (string, 必需): 文档ID，用于标识包含要删除块的文档
- **startIndex** (string, 必需): 删除起始块的索引位置
- **endIndex** (string, 必需): 删除结束块的索引位置，包含此位置

## 文件夹管理工具 (7个)

### 14. get_feishu_root_folder_info - 获取飞书云盘根文件夹信息
- 无参数

### 15. get_feishu_folder_files - 获取指定文件夹中的文件和子文件夹列表
- **folderToken** (string, 必需): 文件夹token，用于标识要获取文件列表的文件夹
- **pageSize** (number, 可选): 每页返回的文件数量，默认为50
- **pageToken** (string, 可选): 分页令牌，用于获取下一页数据

### 16. create_feishu_folder - 在指定父文件夹中创建新文件夹
- **name** (string, 必需): 新文件夹的名称
- **parentToken** (string, 必需): 父文件夹的token，用于指定新文件夹的创建位置

### 17. get_feishu_drive_files_with_meta - 获取云空间目录下所有文件的详细元数据信息
- **folderToken** (string, 必需): 文件夹token，用于标识要获取元数据的文件夹
- **pageSize** (number, 可选): 每页返回的文件数量，默认为50
- **pageToken** (string, 可选): 分页令牌，用于获取下一页数据

### 18. get_feishu_drive_meta - 获取云空间目录/文件的元数据信息
- **requestDoc** (string, 必需): 请求的文档token，用于获取指定文档的元数据

### 19. get_all_feishu_drive_files - 获取指定目录下所有文件（支持分页和数量限制）
- **folderToken** (string, 必需): 文件夹token，用于标识要获取所有文件的文件夹
- **pageSize** (number, 可选): 每页返回的文件数量，默认为50
- **pageToken** (string, 可选): 分页令牌，用于获取下一页数据

### 20. get_feishu_root_folder_meta - 获取云空间根文件夹元数据信息（使用新的API端点）
- 无参数

## 知识库工具 (4个)

### 21. get_feishu_wiki_spaces - 获取飞书知识库空间列表
- **pageSize** (number, 可选): 每页返回的知识库空间数量，默认为50
- **pageToken** (string, 可选): 分页令牌，用于获取下一页数据

### 22. get_feishu_wiki_nodes - 获取飞书知识库空间下的节点列表
- **spaceId** (string, 必需): 知识库空间ID，用于标识要获取节点的知识库空间
- **pageSize** (number, 可选): 每页返回的节点数量，默认为50
- **pageToken** (string, 可选): 分页令牌，用于获取下一页数据

### 23. get_feishu_wiki_node_content - 获取飞书知识库节点的内容
- **token** (string, 必需): 知识库节点token，用于标识要获取内容的节点

### 24. get_feishu_wiki_node_meta - 获取飞书知识库节点的元信息
- **token** (string, 必需): 知识库节点token，用于标识要获取元信息的节点

## 工具功能 (2个)

### 25. convert_feishu_wiki_to_document_id - 将飞书Wiki链接转换为文档ID
- **wikiUrl** (string, 必需): 飞书Wiki的URL链接，需要转换为文档ID

### 26. get_feishu_image_resource - 获取飞书图片资源信息和内容
- **token** (string, 必需): 图片资源的token，用于标识要获取的图片资源

## 参数类型说明

- **string**: 字符串类型参数
- **number**: 数字类型参数
- **必需**: 必须提供的参数
- **可选**: 可以选择性提供的参数，通常有默认值

## 使用建议

1. 大模型在调用工具时，应该根据工具描述和参数说明来提供正确的参数
2. 对于必需参数，必须提供有效值
3. 对于可选参数，如果不提供会使用默认值
4. 分页相关的参数（pageToken）通常用于获取大量数据时的分页处理
5. 所有ID和token参数都需要从飞书API的响应中获取
