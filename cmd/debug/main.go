package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"feishu-mcp-go/pkg/feishu"
	"feishu-mcp-go/pkg/logger"
)

var (
	appID     = flag.String("app-id", "", "飞书应用ID")
	appSecret = flag.String("app-secret", "", "飞书应用密钥")
	baseURL   = flag.String("base-url", "https://open.feishu.cn/open-apis", "飞书API基础URL")
	action    = flag.String("action", "", "要执行的操作：token, create-doc, get-doc, get-content, get-blocks, search, health, validate, folder-info, folder, wiki-spaces, wiki-nodes, wiki-content, wiki-meta, get-block, update-block, create-text, create-code, create-heading, create-list, delete-blocks, batch-create")

	// 文档相关参数
	docID       = flag.String("doc-id", "", "文档ID")
	folderToken = flag.String("folder-token", "", "文件夹Token")
	title       = flag.String("title", "", "文档标题")
	searchKey   = flag.String("search-key", "", "搜索关键词")
	lang        = flag.Int("lang", 0, "语言代码 (0:中文, 1:英文)")

	// 内容操作相关参数
	blockID    = flag.String("block-id", "", "块ID")
	parentID   = flag.String("parent-id", "", "父块ID")
	content    = flag.String("content", "", "内容文本")
	code       = flag.String("code", "", "代码内容")
	language   = flag.String("language", "text", "编程语言")
	level      = flag.Int("level", 1, "标题级别 (1-6)")
	listType   = flag.String("list-type", "bullet", "列表类型 (bullet, numbered)")
	items      = flag.String("items", "", "列表项，用逗号分隔")
	startIndex = flag.Int("start-index", 0, "起始索引")
	endIndex   = flag.Int("end-index", 0, "结束索引")

	// 知识库相关参数
	spaceID    = flag.String("space-id", "", "知识库空间ID")
	nodeToken  = flag.String("node-token", "", "知识库节点Token")
	parentNode = flag.String("parent-node", "", "父节点Token")
	pageSize   = flag.Int("page-size", 20, "分页大小")
	pageToken  = flag.String("page-token", "", "分页Token")

	// 调试相关参数
	verbose = flag.Bool("verbose", false, "启用详细日志输出")
	debug   = flag.Bool("debug", false, "启用调试模式")
)

func main() {
	flag.Parse()

	if *appID == "" || *appSecret == "" {
		fmt.Println("错误: 必须提供 app-id 和 app-secret")
		fmt.Println("使用方法:")
		fmt.Println("  go run cmd/debug/main.go -app-id=your_app_id -app-secret=your_app_secret -action=token")
		os.Exit(1)
	}

	// 创建日志器
	logger := logger.New("info")

	// 创建飞书客户端
	client, err := feishu.NewClient(*appID, *appSecret, *baseURL)
	if err != nil {
		log.Fatalf("创建飞书客户端失败: %v", err)
	}

	// 启用调试模式
	if *debug {
		client.EnableDebugMode(*verbose)
	}

	// 根据action执行不同的调试操作
	switch *action {
	case "token":
		debugGetToken(client, logger)
	case "create-doc":
		debugCreateDocument(client, logger)
	case "get-doc":
		debugGetDocument(client, logger)
	case "get-content":
		debugGetDocumentContent(client, logger)
	case "get-blocks":
		debugGetDocumentBlocks(client, logger)
	case "search":
		debugSearchDocuments(client, logger)
	case "health":
		debugHealthCheck(client, logger)
	case "validate":
		debugValidateAPI(client, logger)
	case "folder-info":
		debugGetFolderInfo(client, logger)
	case "folder":
		debugGetFolderInfo(client, logger)
	case "wiki-spaces":
		debugGetWikiSpaces(client, logger)
	case "wiki-nodes":
		debugGetWikiNodes(client, logger)
	case "wiki-content":
		debugGetWikiNodeContent(client, logger)
	case "wiki-meta":
		debugGetWikiNodeMeta(client, logger)
	// 内容操作相关的debug测试
	case "get-block":
		debugGetBlockContent(client, logger)
	case "update-block":
		debugUpdateBlockText(client, logger)
	case "create-text":
		debugCreateTextBlock(client, logger)
	case "create-code":
		debugCreateCodeBlock(client, logger)
	case "create-heading":
		debugCreateHeadingBlock(client, logger)
	case "create-list":
		debugCreateListBlock(client, logger)
	case "delete-blocks":
		debugDeleteBlocks(client, logger)
	case "batch-create":
		debugBatchCreateBlocks(client, logger)
	default:
		fmt.Println("支持的操作:")
		fmt.Println("文档操作:")
		fmt.Println("  token       - 测试获取访问令牌")
		fmt.Println("  create-doc  - 测试创建文档 (需要 -folder-token 和 -title)")
		fmt.Println("  get-doc     - 测试获取文档信息 (需要 -doc-id)")
		fmt.Println("  get-content - 测试获取文档内容 (需要 -doc-id, 可选 -lang)")
		fmt.Println("  get-blocks  - 测试获取文档块 (需要 -doc-id)")
		fmt.Println("  search      - 测试搜索文档 (需要 -search-key)")
		fmt.Println("内容操作:")
		fmt.Println("  get-block     - 测试获取块内容 (需要 -doc-id 和 -block-id)")
		fmt.Println("  update-block  - 测试更新块内容 (需要 -doc-id, -block-id 和 -content)")
		fmt.Println("  create-text   - 测试创建文本块 (需要 -doc-id 和 -content, 可选 -parent-id)")
		fmt.Println("  create-code   - 测试创建代码块 (需要 -doc-id 和 -code, 可选 -parent-id, -language)")
		fmt.Println("  create-heading- 测试创建标题块 (需要 -doc-id 和 -content, 可选 -parent-id, -level)")
		fmt.Println("  create-list   - 测试创建列表块 (需要 -doc-id 和 -items, 可选 -parent-id, -list-type)")
		fmt.Println("  delete-blocks - 测试删除块 (需要 -doc-id, -start-index, 可选 -end-index)")
		fmt.Println("  batch-create  - 测试批量创建块 (需要 -doc-id)")
		fmt.Println("文件夹操作:")
		fmt.Println("  folder-info - 测试获取文件夹信息")
		fmt.Println("  folder      - 测试获取文件夹信息")
		fmt.Println("知识库操作:")
		fmt.Println("  wiki-spaces - 测试获取知识库空间列表 (可选 -page-size, -page-token)")
		fmt.Println("  wiki-nodes  - 测试获取知识库节点列表 (需要 -space-id, 可选 -parent-node)")
		fmt.Println("  wiki-content- 测试获取知识库节点内容 (需要 -space-id 和 -node-token, 可选 -lang)")
		fmt.Println("  wiki-meta   - 测试获取知识库节点元信息 (需要 -space-id 和 -node-token)")
		fmt.Println("系统操作:")
		fmt.Println("  health      - 执行健康检查")
		fmt.Println("  validate    - 验证API端点")
		fmt.Println("")
		fmt.Println("调试选项:")
		fmt.Println("  -debug      - 启用调试模式")
		fmt.Println("  -verbose    - 启用详细日志输出")
	}
}

func debugGetToken(client *feishu.Client, logger logger.Logger) {
	fmt.Println("🔐 测试获取访问令牌...")

	// 直接获取访问令牌
	token, err := client.GetToken()
	if err != nil {
		fmt.Printf("❌ 获取令牌失败: %v\n", err)
		return
	}

	fmt.Printf("✅ 成功获取访问令牌！\n")
	fmt.Printf("🔑 令牌长度: %d 字符\n", len(token))
	fmt.Printf("🔑 令牌前缀: %s...\n", token[:min(len(token), 20)])

	// 通过调用一个简单的API来验证token有效性
	fmt.Println("🔍 验证令牌有效性...")
	err = client.HealthCheck()
	if err != nil {
		fmt.Printf("⚠️  令牌验证失败: %v\n", err)
		return
	}

	fmt.Println("✅ 令牌验证成功，可以正常访问飞书API！")
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func debugCreateDocument(client *feishu.Client, logger logger.Logger) {
	if *folderToken == "" || *title == "" {
		fmt.Println("❌ 创建文档需要 -folder-token 和 -title 参数")
		return
	}

	fmt.Printf("📝 测试创建文档: %s\n", *title)

	document, err := client.CreateDocument(*folderToken, *title)
	if err != nil {
		fmt.Printf("❌ 创建文档失败: %v\n", err)
		return
	}

	fmt.Println("✅ 文档创建成功!")
	printJSON("文档信息", document)
}

func debugGetDocument(client *feishu.Client, logger logger.Logger) {
	if *docID == "" {
		fmt.Println("❌ 获取文档信息需要 -doc-id 参数")
		return
	}

	fmt.Printf("📄 测试获取文档信息: %s\n", *docID)

	document, err := client.GetDocument(*docID)
	if err != nil {
		fmt.Printf("❌ 获取文档信息失败: %v\n", err)
		return
	}

	fmt.Println("✅ 获取文档信息成功!")
	printJSON("文档信息", document)
}

func debugGetDocumentContent(client *feishu.Client, logger logger.Logger) {
	if *docID == "" {
		fmt.Println("❌ 获取文档内容需要 -doc-id 参数")
		return
	}

	fmt.Printf("📖 测试获取文档内容: %s (语言: %d)\n", *docID, *lang)

	content, err := client.GetDocumentContent(*docID, *lang)
	if err != nil {
		fmt.Printf("❌ 获取文档内容失败: %v\n", err)
		return
	}

	fmt.Println("✅ 获取文档内容成功!")
	printJSON("文档内容", content)
}

func debugGetDocumentBlocks(client *feishu.Client, logger logger.Logger) {
	if *docID == "" {
		fmt.Println("❌ 获取文档块需要 -doc-id 参数")
		return
	}

	fmt.Printf("🧱 测试获取文档块: %s\n", *docID)

	blocks, err := client.GetDocumentBlocks(*docID)
	if err != nil {
		fmt.Printf("❌ 获取文档块失败: %v\n", err)
		return
	}

	fmt.Printf("✅ 获取文档块成功! 共 %d 个块\n", len(blocks))
	printJSON("文档块", blocks)
}

func debugSearchDocuments(client *feishu.Client, logger logger.Logger) {
	if *searchKey == "" {
		fmt.Println("❌ 搜索文档需要 -search-key 参数")
		return
	}

	fmt.Printf("🔍 测试搜索文档: %s\n", *searchKey)

	result, err := client.SearchDocuments(*searchKey, 10, "")
	if err != nil {
		fmt.Printf("❌ 搜索文档失败: %v\n", err)
		return
	}

	fmt.Printf("✅ 搜索文档成功! 共找到 %d 个文档\n", len(result.Data.Documents))
	printJSON("搜索结果", result)
}

func debugHealthCheck(client *feishu.Client, logger logger.Logger) {
	fmt.Println("🏥 执行健康检查...")

	err := client.HealthCheck()
	if err != nil {
		fmt.Printf("❌ 健康检查失败: %v\n", err)
		return
	}

	fmt.Println("✅ 健康检查完成!")
}

func debugValidateAPI(client *feishu.Client, logger logger.Logger) {
	fmt.Println("🔍 验证API端点...")

	validator := feishu.NewAPIValidator(client)
	validator.ValidateEndpoints()

	fmt.Println("✅ API端点验证完成!")
}

func debugGetFolderInfo(client *feishu.Client, logger logger.Logger) {
	fmt.Println("📁 测试获取文件夹信息...")

	// 从环境变量获取测试文件夹token
	folderToken := os.Getenv("TEST_FOLDER_TOKEN")
	if folderToken == "" {
		fmt.Println("⚠️ 未设置TEST_FOLDER_TOKEN环境变量，将尝试获取根目录信息")

		// 尝试获取根目录信息
		rootInfo, err := client.GetRootFolderInfo()
		if err != nil {
			fmt.Printf("❌ 获取根目录信息失败: %v\n", err)
			return
		}

		fmt.Printf("✅ 成功获取根目录信息！\n")
		fmt.Printf("📊 响应代码: %d\n", rootInfo.Code)
		fmt.Printf("📄 消息: %s\n", rootInfo.Msg)
		fmt.Printf("📁 文件夹数量: %d\n", len(rootInfo.Folders))

		for i, folder := range rootInfo.Folders {
			fmt.Printf("   %d. 文件夹名: %s\n", i+1, folder.FolderName)
			fmt.Printf("      Token: %s\n", folder.FolderToken)
		}
		return
	}

	// 获取指定文件夹信息
	folderDetail, err := client.GetFolderInfo(folderToken)
	if err != nil {
		fmt.Printf("❌ 获取文件夹信息失败: %v\n", err)
		return
	}

	fmt.Printf("✅ 成功获取文件夹信息！\n")
	fmt.Printf("📊 响应代码: %d\n", folderDetail.Code)
	fmt.Printf("📄 消息: %s\n", folderDetail.Msg)
	fmt.Printf("📁 文件夹名: %s\n", folderDetail.Folder.FolderName)
	fmt.Printf("🔑 文件夹Token: %s\n", folderDetail.Folder.FolderToken)

	// 尝试获取文件夹下的文件列表
	fmt.Println("\n📋 获取文件夹下的文件列表...")
	fileList, err := client.GetFolderFiles(folderToken, 10, "")
	if err != nil {
		fmt.Printf("❌ 获取文件列表失败: %v\n", err)
		return
	}

	fmt.Printf("✅ 成功获取文件列表！\n")
	fmt.Printf("📊 响应代码: %d\n", fileList.Code)
	fmt.Printf("📄 消息: %s\n", fileList.Msg)
	fmt.Printf("📁 文件数量: %d\n", len(fileList.Folders))

	for i, file := range fileList.Folders {
		fmt.Printf("   %d. 文件名: %s\n", i+1, file.FolderName)
		fmt.Printf("      Token: %s\n", file.FolderToken)
	}
}

func printJSON(title string, data interface{}) {
	fmt.Printf("\n📋 %s:\n", title)
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		fmt.Printf("JSON序列化失败: %v\n", err)
		return
	}
	fmt.Println(string(jsonData))
	fmt.Println()
}

// ==================== 知识库调试函数 ====================

func debugGetWikiSpaces(client *feishu.Client, logger logger.Logger) {
	fmt.Println("📚 测试获取知识库空间列表...")

	response, err := client.GetWikiSpaces(*pageSize, *pageToken)
	if err != nil {
		fmt.Printf("❌ 获取知识库空间列表失败: %v\n", err)
		return
	}

	fmt.Printf("✅ 获取知识库空间列表成功! 共找到 %d 个空间\n", len(response.Data.Items))
	fmt.Printf("📊 分页信息: hasMore=%t, pageToken=%s\n", response.Data.HasMore, response.Data.PageToken)
	printJSON("知识库空间列表", response)
}

func debugGetWikiNodes(client *feishu.Client, logger logger.Logger) {
	if *spaceID == "" {
		fmt.Println("❌ 获取知识库节点列表需要 -space-id 参数")
		return
	}

	fmt.Printf("📁 测试获取知识库节点列表: 空间ID=%s\n", *spaceID)
	if *parentNode != "" {
		fmt.Printf("📂 父节点: %s\n", *parentNode)
	}

	response, err := client.GetWikiSpaceNodes(*spaceID, *pageSize, *pageToken, *parentNode)
	if err != nil {
		fmt.Printf("❌ 获取知识库节点列表失败: %v\n", err)
		return
	}

	fmt.Printf("✅ 获取知识库节点列表成功! 共找到 %d 个节点\n", len(response.Data.Items))
	fmt.Printf("📊 分页信息: hasMore=%t, pageToken=%s\n", response.Data.HasMore, response.Data.PageToken)
	printJSON("知识库节点列表", response)
}

func debugGetWikiNodeContent(client *feishu.Client, logger logger.Logger) {
	if *spaceID == "" || *nodeToken == "" {
		fmt.Println("❌ 获取知识库节点内容需要 -space-id 和 -node-token 参数")
		return
	}

	fmt.Printf("📖 测试获取知识库节点内容: 空间ID=%s, 节点Token=%s, 语言=%d\n", *spaceID, *nodeToken, *lang)

	content, err := client.GetWikiNodeContent(*spaceID, *nodeToken, *lang)
	if err != nil {
		fmt.Printf("❌ 获取知识库节点内容失败: %v\n", err)
		return
	}

	fmt.Println("✅ 获取知识库节点内容成功!")
	fmt.Printf("📄 内容类型: %s\n", content.Type)
	fmt.Printf("📄 内容长度: %d 字符\n", len(content.Content))
	printJSON("知识库节点内容", content)
}

func debugGetWikiNodeMeta(client *feishu.Client, logger logger.Logger) {
	if *spaceID == "" || *nodeToken == "" {
		fmt.Println("❌ 获取知识库节点元信息需要 -space-id 和 -node-token 参数")
		return
	}

	fmt.Printf("ℹ️ 测试获取知识库节点元信息: 空间ID=%s, 节点Token=%s\n", *spaceID, *nodeToken)

	meta, err := client.GetWikiNodeMeta(*spaceID, *nodeToken)
	if err != nil {
		fmt.Printf("❌ 获取知识库节点元信息失败: %v\n", err)
		return
	}

	fmt.Println("✅ 获取知识库节点元信息成功!")
	fmt.Printf("📝 节点标题: %s\n", meta.Title)
	fmt.Printf("📝 节点类型: %s\n", meta.NodeType)
	fmt.Printf("📝 对象类型: %s\n", meta.ObjType)
	fmt.Printf("📝 是否有子节点: %t\n", meta.HasChild)
	fmt.Printf("📅 节点创建时间: %s\n", meta.NodeCreateTime)
	fmt.Printf("📅 对象编辑时间: %s\n", meta.ObjEditTime)
	fmt.Printf("👤 创建者: %s\n", meta.Creator)
	fmt.Printf("👤 拥有者: %s\n", meta.Owner)
	printJSON("知识库节点元信息", meta)
}

// ==================== 内容操作调试函数 ====================

func debugGetBlockContent(client *feishu.Client, logger logger.Logger) {
	if *docID == "" || *blockID == "" {
		fmt.Println("❌ 获取块内容需要 -doc-id 和 -block-id 参数")
		return
	}

	fmt.Printf("🧱 测试获取块内容: 文档ID=%s, 块ID=%s\n", *docID, *blockID)

	blockContent, err := client.GetBlockContent(*docID, *blockID)
	if err != nil {
		fmt.Printf("❌ 获取块内容失败: %v\n", err)
		return
	}

	fmt.Println("✅ 获取块内容成功!")
	fmt.Printf("🆔 块ID: %s\n", blockContent.BlockID)
	fmt.Printf("📝 块类型: %d\n", blockContent.BlockType)
	fmt.Printf("👨‍👩‍👧‍👦 父块ID: %s\n", blockContent.ParentID)

	// 根据块类型显示相应的内容
	switch blockContent.BlockType {
	case 2: // 文本块
		if blockContent.Text != nil {
			fmt.Printf("📄 文本内容: \n")
			printJSON("文本块内容", blockContent.Text)
		}
	case 3: // 代码块
		if blockContent.Code != nil {
			fmt.Printf("💻 代码内容: \n")
			printJSON("代码块内容", blockContent.Code)
		}
	case 4: // 标题块
		if blockContent.Heading != nil {
			fmt.Printf("📊 标题内容: \n")
			printJSON("标题块内容", blockContent.Heading)
		}
	case 5: // 列表块
		if blockContent.List != nil {
			fmt.Printf("📋 列表内容: \n")
			printJSON("列表块内容", blockContent.List)
		}
	case 1: // 页面块
		if blockContent.Page != nil {
			fmt.Printf("📜 页面内容: \n")
			printJSON("页面块内容", blockContent.Page)
		}
	default:
		fmt.Printf("❓ 未知块类型: %d\n", blockContent.BlockType)
		printJSON("完整块内容", blockContent)
	}
}

func debugUpdateBlockText(client *feishu.Client, logger logger.Logger) {
	if *docID == "" || *blockID == "" || *content == "" {
		fmt.Println("❌ 更新块内容需要 -doc-id, -block-id 和 -content 参数")
		return
	}

	fmt.Printf("✏️ 测试更新块内容: 文档ID=%s, 块ID=%s\n", *docID, *blockID)
	fmt.Printf("📝 新内容: %s\n", *content)

	// 使用与获取到的块结构相同的格式
	contentMap := map[string]interface{}{
		"text": map[string]interface{}{
			"elements": []interface{}{
				map[string]interface{}{
					"text_run": map[string]interface{}{
						"content": *content,
						"text_element_style": map[string]interface{}{
							"bold":          false,
							"inline_code":   false,
							"italic":        false,
							"strikethrough": false,
							"underline":     false,
						},
					},
				},
			},
			"style": map[string]interface{}{
				"align":  1,
				"folded": false,
			},
		},
	}

	err := client.UpdateBlockText(*docID, *blockID, contentMap)
	if err != nil {
		fmt.Printf("❌ 更新块内容失败: %v\n", err)
		return
	}

	fmt.Println("✅ 更新块内容成功!")
}

func debugCreateTextBlock(client *feishu.Client, logger logger.Logger) {
	if *docID == "" || *content == "" {
		fmt.Println("❌ 创建文本块需要 -doc-id 和 -content 参数")
		return
	}

	fmt.Printf("📝 测试创建文本块: 文档ID=%s\n", *docID)
	fmt.Printf("📄 文本内容: %s\n", *content)
	if *parentID != "" {
		fmt.Printf("👨‍👩‍👧‍👦 父块ID: %s\n", *parentID)
	}

	blockID, err := client.CreateTextBlock(*docID, *parentID, *content, nil)
	if err != nil {
		fmt.Printf("❌ 创建文本块失败: %v\n", err)
		return
	}

	fmt.Printf("✅ 创建文本块成功! 块ID: %s\n", blockID)
}

func debugCreateCodeBlock(client *feishu.Client, logger logger.Logger) {
	if *docID == "" || *code == "" {
		fmt.Println("❌ 创建代码块需要 -doc-id 和 -code 参数")
		return
	}

	fmt.Printf("💻 测试创建代码块: 文档ID=%s\n", *docID)
	fmt.Printf("📄 代码内容: %s\n", *code)
	fmt.Printf("🔤 编程语言: %s\n", *language)
	if *parentID != "" {
		fmt.Printf("👨‍👩‍👧‍👦 父块ID: %s\n", *parentID)
	}

	blockID, err := client.CreateCodeBlock(*docID, *parentID, *code, *language)
	if err != nil {
		fmt.Printf("❌ 创建代码块失败: %v\n", err)
		return
	}

	fmt.Printf("✅ 创建代码块成功! 块ID: %s\n", blockID)
}

func debugCreateHeadingBlock(client *feishu.Client, logger logger.Logger) {
	if *docID == "" || *content == "" {
		fmt.Println("❌ 创建标题块需要 -doc-id 和 -content 参数")
		return
	}

	fmt.Printf("📊 测试创建标题块: 文档ID=%s\n", *docID)
	fmt.Printf("📄 标题内容: %s\n", *content)
	fmt.Printf("🔢 标题级别: %d\n", *level)
	if *parentID != "" {
		fmt.Printf("👨‍👩‍👧‍👦 父块ID: %s\n", *parentID)
	}

	blockID, err := client.CreateHeadingBlock(*docID, *parentID, *content, *level)
	if err != nil {
		fmt.Printf("❌ 创建标题块失败: %v\n", err)
		return
	}

	fmt.Printf("✅ 创建标题块成功! 块ID: %s\n", blockID)
}

func debugCreateListBlock(client *feishu.Client, logger logger.Logger) {
	if *docID == "" || *items == "" {
		fmt.Println("❌ 创建列表块需要 -doc-id 和 -items 参数")
		return
	}

	// 解析逗号分隔的列表项
	testItems := strings.Split(*items, ",")
	for i, item := range testItems {
		testItems[i] = strings.TrimSpace(item)
	}

	fmt.Printf("📋 测试创建列表块: 文档ID=%s\n", *docID)
	fmt.Printf("📝 列表类型: %s\n", *listType)
	fmt.Printf("📄 列表项数量: %d\n", len(testItems))
	if *parentID != "" {
		fmt.Printf("👨‍👩‍👧‍👦 父块ID: %s\n", *parentID)
	}

	blockID, err := client.CreateListBlock(*docID, *parentID, testItems, *listType)
	if err != nil {
		fmt.Printf("❌ 创建列表块失败: %v\n", err)
		return
	}

	fmt.Printf("✅ 创建列表块成功! 块ID: %s\n", blockID)
}

func debugDeleteBlocks(client *feishu.Client, logger logger.Logger) {
	if *docID == "" {
		fmt.Println("❌ 删除块需要 -doc-id 和 -start-index 参数")
		return
	}

	endIdx := *endIndex
	if endIdx == 0 {
		endIdx = *startIndex
	}

	fmt.Printf("🗑️ 测试删除块: 文档ID=%s\n", *docID)
	fmt.Printf("📍 起始索引: %d\n", *startIndex)
	fmt.Printf("📍 结束索引: %d\n", endIdx)

	err := client.DeleteBlocks(*docID, *startIndex, endIdx)
	if err != nil {
		fmt.Printf("❌ 删除块失败: %v\n", err)
		return
	}

	blocksCount := endIdx - *startIndex + 1
	fmt.Printf("✅ 删除块成功! 删除了 %d 个块\n", blocksCount)
}

func debugBatchCreateBlocks(client *feishu.Client, logger logger.Logger) {
	if *docID == "" {
		fmt.Println("❌ 批量创建块需要 -doc-id 参数")
		return
	}

	fmt.Printf("🏗️ 测试批量创建块: 文档ID=%s\n", *docID)

	// 创建示例块定义
	blocks := []map[string]interface{}{
		{
			"blockType": "text",
			"content": map[string]interface{}{
				"text": "这是第一个文本块",
			},
		},
		{
			"blockType": "text",
			"content": map[string]interface{}{
				"text": "这是第二个文本块",
			},
		},
		{
			"blockType": "code",
			"content": map[string]interface{}{
				"code":     "fmt.Println(\"Hello, World!\")",
				"language": "go",
			},
		},
	}

	fmt.Printf("📦 将创建 %d 个块\n", len(blocks))

	// 由于我们的实现是循环调用单个CreateBlock，这里模拟批量创建
	var blockIDs []string
	for i, block := range blocks {
		blockType := block["blockType"].(string)
		content := block["content"].(map[string]interface{})

		blockID, err := client.CreateBlock(*docID, blockType, "", content)
		if err != nil {
			fmt.Printf("❌ 创建块 %d 失败: %v\n", i+1, err)
			return
		}

		blockIDs = append(blockIDs, blockID)
		fmt.Printf("✅ 创建块 %d 成功，ID: %s\n", i+1, blockID)
	}

	fmt.Printf("🎉 批量创建块完成! 共创建 %d 个块\n", len(blockIDs))
	printJSON("创建的块ID列表", blockIDs)
}
