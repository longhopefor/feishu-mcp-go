package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	"feishu-mcp-go/pkg/feishu"
	"feishu-mcp-go/pkg/logger"
)

var (
	appID     = flag.String("app-id", "", "飞书应用ID")
	appSecret = flag.String("app-secret", "", "飞书应用密钥")
	baseURL   = flag.String("base-url", "https://open.feishu.cn/open-apis", "飞书API基础URL")
	action    = flag.String("action", "", "要执行的操作：token, create-doc, get-doc, get-content, get-blocks, search, health, validate, folder-info, folder, wiki-spaces, wiki-nodes, wiki-content, wiki-meta")

	// 文档相关参数
	docID       = flag.String("doc-id", "", "文档ID")
	folderToken = flag.String("folder-token", "", "文件夹Token")
	title       = flag.String("title", "", "文档标题")
	searchKey   = flag.String("search-key", "", "搜索关键词")
	lang        = flag.Int("lang", 0, "语言代码 (0:中文, 1:英文)")

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
	default:
		fmt.Println("支持的操作:")
		fmt.Println("文档操作:")
		fmt.Println("  token       - 测试获取访问令牌")
		fmt.Println("  create-doc  - 测试创建文档 (需要 -folder-token 和 -title)")
		fmt.Println("  get-doc     - 测试获取文档信息 (需要 -doc-id)")
		fmt.Println("  get-content - 测试获取文档内容 (需要 -doc-id, 可选 -lang)")
		fmt.Println("  get-blocks  - 测试获取文档块 (需要 -doc-id)")
		fmt.Println("  search      - 测试搜索文档 (需要 -search-key)")
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
