#!/bin/bash

# 飞书MCP Go API调试脚本
# 用法: ./scripts/debug.sh <action> [options]

set -e

# 检查是否有debug.env文件
if [ ! -f "debug.env" ]; then
    echo "❌ 找不到 debug.env 文件"
    echo "请复制 debug.env.example 为 debug.env，并填入真实的配置信息"
    exit 1
fi

# 加载环境变量
source debug.env

# 检查必要的环境变量
if [ -z "$FEISHU_APP_ID" ] || [ -z "$FEISHU_APP_SECRET" ]; then
    echo "❌ 请在 debug.env 中设置 FEISHU_APP_ID 和 FEISHU_APP_SECRET"
    exit 1
fi

# 构建调试工具
echo "🔧 构建调试工具..."
go build -o build/debug cmd/debug/main.go

# 设置默认参数
APP_ID="${FEISHU_APP_ID}"
APP_SECRET="${FEISHU_APP_SECRET}"
BASE_URL="${FEISHU_BASE_URL:-https://open.feishu.cn/open-apis}"
DEBUG_FLAG=""
VERBOSE_FLAG=""

# 设置调试标志
if [ "$DEBUG_MODE" = "true" ]; then
    DEBUG_FLAG="-debug"
fi

if [ "$VERBOSE_MODE" = "true" ]; then
    VERBOSE_FLAG="-verbose"
fi

# 获取第一个参数作为action
ACTION="$1"
shift || true

# 根据action执行不同的操作
case "$ACTION" in
    "token")
        echo "🔐 测试获取访问令牌..."
        ./build/debug -app-id="$APP_ID" -app-secret="$APP_SECRET" -base-url="$BASE_URL" -action=token $DEBUG_FLAG $VERBOSE_FLAG
        ;;
    "health")
        echo "🏥 执行健康检查..."
        ./build/debug -app-id="$APP_ID" -app-secret="$APP_SECRET" -base-url="$BASE_URL" -action=health $DEBUG_FLAG $VERBOSE_FLAG
        ;;
    "validate")
        echo "🔍 验证API端点..."
        ./build/debug -app-id="$APP_ID" -app-secret="$APP_SECRET" -base-url="$BASE_URL" -action=validate $DEBUG_FLAG $VERBOSE_FLAG
        ;;
    "create-doc")
        if [ -z "$TEST_FOLDER_TOKEN" ]; then
            echo "❌ 请在 debug.env 中设置 TEST_FOLDER_TOKEN"
            exit 1
        fi
        TITLE="${TEST_DOCUMENT_TITLE:-API测试文档}"
        echo "📝 测试创建文档: $TITLE"
        ./build/debug -app-id="$APP_ID" -app-secret="$APP_SECRET" -base-url="$BASE_URL" -action=create-doc -folder-token="$TEST_FOLDER_TOKEN" -title="$TITLE" $DEBUG_FLAG $VERBOSE_FLAG
        ;;
    "get-doc")
        # 支持通过命令行参数传递文档ID，如果没有则使用环境变量
        DOC_ID="${1:-$TEST_DOCUMENT_ID}"
        if [ -z "$DOC_ID" ]; then
            echo "❌ 请提供 doc_id 参数或在 debug.env 中设置 TEST_DOCUMENT_ID"
            echo "用法: $0 get-doc <doc_id>"
            exit 1
        fi
        echo "📄 测试获取文档信息: $DOC_ID"
        ./build/debug -app-id="$APP_ID" -app-secret="$APP_SECRET" -base-url="$BASE_URL" -action=get-doc -doc-id="$DOC_ID" $DEBUG_FLAG $VERBOSE_FLAG
        ;;
    "get-content")
        # 支持通过命令行参数传递文档ID，如果没有则使用环境变量
        DOC_ID="${1:-$TEST_DOCUMENT_ID}"
        if [ -z "$DOC_ID" ]; then
            echo "❌ 请提供 doc_id 参数或在 debug.env 中设置 TEST_DOCUMENT_ID"
            echo "用法: $0 get-content <doc_id> [lang]"
            exit 1
        fi
        LANG="${2:-0}"
        echo "📖 测试获取文档内容: $DOC_ID (语言: $LANG)"
        ./build/debug -app-id="$APP_ID" -app-secret="$APP_SECRET" -base-url="$BASE_URL" -action=get-content -doc-id="$DOC_ID" -lang="$LANG" $DEBUG_FLAG $VERBOSE_FLAG
        ;;
    "get-blocks")
        # 支持通过命令行参数传递文档ID，如果没有则使用环境变量
        DOC_ID="${1:-$TEST_DOCUMENT_ID}"
        if [ -z "$DOC_ID" ]; then
            echo "❌ 请提供 doc_id 参数或在 debug.env 中设置 TEST_DOCUMENT_ID"
            echo "用法: $0 get-blocks <doc_id>"
            exit 1
        fi
        echo "🧱 测试获取文档块: $DOC_ID"
        ./build/debug -app-id="$APP_ID" -app-secret="$APP_SECRET" -base-url="$BASE_URL" -action=get-blocks -doc-id="$DOC_ID" $DEBUG_FLAG $VERBOSE_FLAG
        ;;
    "search")
        SEARCH_KEY="${1:-test}"
        echo "🔍 测试搜索文档: $SEARCH_KEY"
        ./build/debug -app-id="$APP_ID" -app-secret="$APP_SECRET" -base-url="$BASE_URL" -action=search -search-key="$SEARCH_KEY" $DEBUG_FLAG $VERBOSE_FLAG
        ;;
    "get-block-content")
        # 支持通过命令行参数传递文档ID和块ID，如果没有则使用环境变量
        DOC_ID="${1:-$TEST_DOCUMENT_ID}"
        BLOCK_ID="${2:-$TEST_BLOCK_ID}"
        if [ -z "$DOC_ID" ] || [ -z "$BLOCK_ID" ]; then
            echo "❌ 请提供 doc_id 和 block_id 参数或在 debug.env 中设置 TEST_DOCUMENT_ID 和 TEST_BLOCK_ID"
            echo "用法: $0 get-block-content <doc_id> <block_id>"
            exit 1
        fi
        echo "🧱 测试获取块内容: 文档=$DOC_ID, 块=$BLOCK_ID"
        ./build/debug -app-id="$APP_ID" -app-secret="$APP_SECRET" -base-url="$BASE_URL" -action=get-block -doc-id="$DOC_ID" -block-id="$BLOCK_ID" $DEBUG_FLAG $VERBOSE_FLAG
        ;;
    "update-block-text")
        # 支持通过命令行参数传递文档ID、块ID和内容
        DOC_ID="${1:-$TEST_DOCUMENT_ID}"
        BLOCK_ID="${2:-$TEST_BLOCK_ID}"
        CONTENT="${3:-更新的测试内容}"
        if [ -z "$DOC_ID" ] || [ -z "$BLOCK_ID" ]; then
            echo "❌ 请提供 doc_id 和 block_id 参数或在 debug.env 中设置 TEST_DOCUMENT_ID 和 TEST_BLOCK_ID"
            echo "用法: $0 update-block-text <doc_id> <block_id> [content]"
            exit 1
        fi
        echo "✏️ 测试更新块文本: 文档=$DOC_ID, 块=$BLOCK_ID"
        ./build/debug -app-id="$APP_ID" -app-secret="$APP_SECRET" -base-url="$BASE_URL" -action=update-block -doc-id="$DOC_ID" -block-id="$BLOCK_ID" -content="$CONTENT" $DEBUG_FLAG $VERBOSE_FLAG
        ;;
    "create-text-block")
        # 支持通过命令行参数传递文档ID、父块ID和文本内容
        DOC_ID="${1:-$TEST_DOCUMENT_ID}"
        PARENT_ID="${2:-$TEST_PARENT_BLOCK_ID}"
        CONTENT="${3:-这是一个测试文本块}"
        if [ -z "$DOC_ID" ]; then
            echo "❌ 请提供 doc_id 参数或在 debug.env 中设置 TEST_DOCUMENT_ID"
            echo "用法: $0 create-text-block <doc_id> [parent_id] [content]"
            exit 1
        fi
        echo "📝 测试创建文本块: 文档=$DOC_ID, 内容=$CONTENT"
        ./build/debug -app-id="$APP_ID" -app-secret="$APP_SECRET" -base-url="$BASE_URL" -action=create-text -doc-id="$DOC_ID" -parent-id="$PARENT_ID" -content="$CONTENT" $DEBUG_FLAG $VERBOSE_FLAG
        ;;
    "create-code-block")
        # 支持通过命令行参数传递文档ID、父块ID、代码内容和编程语言
        DOC_ID="${1:-$TEST_DOCUMENT_ID}"
        PARENT_ID="${2:-$TEST_PARENT_BLOCK_ID}"
        CODE="${3:-console.log('Hello, World!');}"
        LANGUAGE="${4:-javascript}"
        if [ -z "$DOC_ID" ]; then
            echo "❌ 请提供 doc_id 参数或在 debug.env 中设置 TEST_DOCUMENT_ID"
            echo "用法: $0 create-code-block <doc_id> [parent_id] [code] [language]"
            exit 1
        fi
        echo "💻 测试创建代码块: 文档=$DOC_ID, 语言=$LANGUAGE"
        ./build/debug -app-id="$APP_ID" -app-secret="$APP_SECRET" -base-url="$BASE_URL" -action=create-code -doc-id="$DOC_ID" -parent-id="$PARENT_ID" -code="$CODE" -language="$LANGUAGE" $DEBUG_FLAG $VERBOSE_FLAG
        ;;
    "create-heading-block")
        # 支持通过命令行参数传递文档ID、父块ID、标题文本和级别
        DOC_ID="${1:-$TEST_DOCUMENT_ID}"
        PARENT_ID="${2:-$TEST_PARENT_BLOCK_ID}"
        CONTENT="${3:-测试标题}"
        LEVEL="${4:-1}"
        if [ -z "$DOC_ID" ]; then
            echo "❌ 请提供 doc_id 参数或在 debug.env 中设置 TEST_DOCUMENT_ID"
            echo "用法: $0 create-heading-block <doc_id> [parent_id] [content] [level]"
            exit 1
        fi
        echo "📋 测试创建标题块: 文档=$DOC_ID, H$LEVEL - $CONTENT"
        ./build/debug -app-id="$APP_ID" -app-secret="$APP_SECRET" -base-url="$BASE_URL" -action=create-heading -doc-id="$DOC_ID" -parent-id="$PARENT_ID" -content="$CONTENT" -level="$LEVEL" $DEBUG_FLAG $VERBOSE_FLAG
        ;;
    "create-list-block")
        # 支持通过命令行参数传递文档ID、父块ID、列表项和列表类型
        DOC_ID="${1:-$TEST_DOCUMENT_ID}"
        PARENT_ID="${2:-$TEST_PARENT_BLOCK_ID}"
        ITEMS="${3:-项目1,项目2,项目3}"
        LIST_TYPE="${4:-bullet}"
        if [ -z "$DOC_ID" ]; then
            echo "❌ 请提供 doc_id 参数或在 debug.env 中设置 TEST_DOCUMENT_ID"
            echo "用法: $0 create-list-block <doc_id> [parent_id] [items] [list_type]"
            exit 1
        fi
        echo "📃 测试创建列表块: 文档=$DOC_ID, $LIST_TYPE 类型"
        ./build/debug -app-id="$APP_ID" -app-secret="$APP_SECRET" -base-url="$BASE_URL" -action=create-list -doc-id="$DOC_ID" -parent-id="$PARENT_ID" -items="$ITEMS" -list-type="$LIST_TYPE" $DEBUG_FLAG $VERBOSE_FLAG
        ;;
    "batch-create-blocks")
        # 支持通过命令行参数传递文档ID
        DOC_ID="${1:-$TEST_DOCUMENT_ID}"
        PARENT_ID="${2:-$TEST_PARENT_BLOCK_ID}"
        if [ -z "$DOC_ID" ]; then
            echo "❌ 请提供 doc_id 参数或在 debug.env 中设置 TEST_DOCUMENT_ID"
            echo "用法: $0 batch-create-blocks <doc_id> [parent_id]"
            exit 1
        fi
        echo "📦 测试批量创建块: 文档=$DOC_ID"
        ./build/debug -app-id="$APP_ID" -app-secret="$APP_SECRET" -base-url="$BASE_URL" -action=batch-create -doc-id="$DOC_ID" -parent-id="$PARENT_ID" $DEBUG_FLAG $VERBOSE_FLAG
        ;;
    "delete-blocks")
        # 支持通过命令行参数传递文档ID、起始和结束索引
        DOC_ID="${1:-$TEST_DOCUMENT_ID}"
        START_INDEX="${2:-0}"
        END_INDEX="${3:-1}"
        if [ -z "$DOC_ID" ]; then
            echo "❌ 请提供 doc_id 参数或在 debug.env 中设置 TEST_DOCUMENT_ID"
            echo "用法: $0 delete-blocks <doc_id> [start_index] [end_index]"
            exit 1
        fi
        echo "🗑️ 测试删除文档块: $DOC_ID (索引 $START_INDEX-$END_INDEX)"
        ./build/debug -app-id="$APP_ID" -app-secret="$APP_SECRET" -base-url="$BASE_URL" -action=delete-blocks -doc-id="$DOC_ID" -start-index="$START_INDEX" -end-index="$END_INDEX" $DEBUG_FLAG $VERBOSE_FLAG
        ;;
    "folder-info"|"folder")
        echo "📁 测试获取文件夹信息..."
        ./build/debug -app-id="$APP_ID" -app-secret="$APP_SECRET" -base-url="$BASE_URL" -action=folder-info $DEBUG_FLAG $VERBOSE_FLAG
        ;;
    "wiki-spaces")
        echo "📚 测试获取知识库空间列表..."
        PAGE_SIZE="${1:-20}"
        PAGE_TOKEN="${2:-}"
        ARGS=""
        if [ -n "$PAGE_TOKEN" ]; then
            ARGS="-page-token=\"$PAGE_TOKEN\""
        fi
        ./build/debug -app-id="$APP_ID" -app-secret="$APP_SECRET" -base-url="$BASE_URL" -action=wiki-spaces -page-size="$PAGE_SIZE" $ARGS $DEBUG_FLAG $VERBOSE_FLAG
        ;;
    "wiki-nodes")
        # 支持通过命令行参数传递 space_id，如果没有则使用环境变量
        SPACE_ID="${1:-$TEST_WIKI_SPACE_ID}"
        if [ -z "$SPACE_ID" ]; then
            echo "❌ 请提供 space_id 参数或在 debug.env 中设置 TEST_WIKI_SPACE_ID"
            echo "用法: $0 wiki-nodes <space_id> [page_size] [parent_node]"
            exit 1
        fi
        echo "📁 测试获取知识库节点列表: 空间ID=$SPACE_ID"
        PAGE_SIZE="${2:-20}"
        PARENT_NODE="${3:-}"
        ARGS=""
        if [ -n "$PARENT_NODE" ]; then
            ARGS="-parent-node=\"$PARENT_NODE\""
        fi
        ./build/debug -app-id="$APP_ID" -app-secret="$APP_SECRET" -base-url="$BASE_URL" -action=wiki-nodes -space-id="$SPACE_ID" -page-size="$PAGE_SIZE" $ARGS $DEBUG_FLAG $VERBOSE_FLAG
        ;;
    "wiki-content")
        # 支持通过命令行参数传递参数，如果没有则使用环境变量
        SPACE_ID="${1:-$TEST_WIKI_SPACE_ID}"
        NODE_TOKEN="${2:-$TEST_WIKI_NODE_TOKEN}"
        if [ -z "$SPACE_ID" ] || [ -z "$NODE_TOKEN" ]; then
            echo "❌ 请提供 space_id 和 node_token 参数或在 debug.env 中设置 TEST_WIKI_SPACE_ID 和 TEST_WIKI_NODE_TOKEN"
            echo "用法: $0 wiki-content <space_id> <node_token> [lang]"
            exit 1
        fi
        LANG="${3:-0}"
        echo "📖 测试获取知识库节点内容: $SPACE_ID / $NODE_TOKEN (语言: $LANG)"
        ./build/debug -app-id="$APP_ID" -app-secret="$APP_SECRET" -base-url="$BASE_URL" -action=wiki-content -space-id="$SPACE_ID" -node-token="$NODE_TOKEN" -lang="$LANG" $DEBUG_FLAG $VERBOSE_FLAG
        ;;
    "wiki-meta")
        # 支持通过命令行参数传递参数，如果没有则使用环境变量
        SPACE_ID="${1:-$TEST_WIKI_SPACE_ID}"
        NODE_TOKEN="${2:-$TEST_WIKI_NODE_TOKEN}"
        if [ -z "$SPACE_ID" ] || [ -z "$NODE_TOKEN" ]; then
            echo "❌ 请提供 space_id 和 node_token 参数或在 debug.env 中设置 TEST_WIKI_SPACE_ID 和 TEST_WIKI_NODE_TOKEN"
            echo "用法: $0 wiki-meta <space_id> <node_token>"
            exit 1
        fi
        echo "ℹ️ 测试获取知识库节点元信息: $SPACE_ID / $NODE_TOKEN"
        ./build/debug -app-id="$APP_ID" -app-secret="$APP_SECRET" -base-url="$BASE_URL" -action=wiki-meta -space-id="$SPACE_ID" -node-token="$NODE_TOKEN" $DEBUG_FLAG $VERBOSE_FLAG
        ;;
    "all")
        echo "🚀 运行所有API测试..."
        echo ""
        
        # 基础检查
        echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
        echo "🔐 1. 测试获取访问令牌"
        echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
        ./build/debug -app-id="$APP_ID" -app-secret="$APP_SECRET" -base-url="$BASE_URL" -action=token $DEBUG_FLAG $VERBOSE_FLAG
        echo ""
        
        echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
        echo "🏥 2. 执行健康检查"
        echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
        ./build/debug -app-id="$APP_ID" -app-secret="$APP_SECRET" -base-url="$BASE_URL" -action=health $DEBUG_FLAG $VERBOSE_FLAG
        echo ""
        
        echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
        echo "🔍 3. 验证API端点"
        echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
        ./build/debug -app-id="$APP_ID" -app-secret="$APP_SECRET" -base-url="$BASE_URL" -action=validate $DEBUG_FLAG $VERBOSE_FLAG
        echo ""
        
        # 文档操作测试
        if [ -n "$TEST_DOCUMENT_ID" ]; then
            echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
            echo "📄 4. 测试获取文档信息"
            echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
            ./build/debug -app-id="$APP_ID" -app-secret="$APP_SECRET" -base-url="$BASE_URL" -action=get-doc -doc-id="$TEST_DOCUMENT_ID" $DEBUG_FLAG $VERBOSE_FLAG
            echo ""
            
            echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
            echo "📖 5. 测试获取文档内容"
            echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
            ./build/debug -app-id="$APP_ID" -app-secret="$APP_SECRET" -base-url="$BASE_URL" -action=get-content -doc-id="$TEST_DOCUMENT_ID" -lang=0 $DEBUG_FLAG $VERBOSE_FLAG
            echo ""
            
            echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
            echo "🧱 6. 测试获取文档块"
            echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
            ./build/debug -app-id="$APP_ID" -app-secret="$APP_SECRET" -base-url="$BASE_URL" -action=get-blocks -doc-id="$TEST_DOCUMENT_ID" $DEBUG_FLAG $VERBOSE_FLAG
            echo ""
        else
            echo "⚠️ 跳过文档操作测试，请在 debug.env 中设置 TEST_DOCUMENT_ID"
        fi
        
        echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
        echo "🔍 7. 测试搜索文档"
        echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
        ./build/debug -app-id="$APP_ID" -app-secret="$APP_SECRET" -base-url="$BASE_URL" -action=search -search-key="test" $DEBUG_FLAG $VERBOSE_FLAG
        echo ""
        
        echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
        echo "📁 8. 测试获取文件夹信息"
        echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
        ./build/debug -app-id="$APP_ID" -app-secret="$APP_SECRET" -base-url="$BASE_URL" -action=folder-info $DEBUG_FLAG $VERBOSE_FLAG
        echo ""
        
        # 内容操作工具测试
        if [ -n "$TEST_DOCUMENT_ID" ] && [ -n "$TEST_BLOCK_ID" ]; then
            echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
            echo "🧱 9. 测试获取块内容"
            echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
            ./build/debug -app-id="$APP_ID" -app-secret="$APP_SECRET" -base-url="$BASE_URL" -action=get-block -doc-id="$TEST_DOCUMENT_ID" -block-id="$TEST_BLOCK_ID" $DEBUG_FLAG $VERBOSE_FLAG
            echo ""
            
            echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
            echo "✏️ 10. 测试更新块文本"
            echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
            ./build/debug -app-id="$APP_ID" -app-secret="$APP_SECRET" -base-url="$BASE_URL" -action=update-block -doc-id="$TEST_DOCUMENT_ID" -block-id="$TEST_BLOCK_ID" -content="自动测试更新的内容" $DEBUG_FLAG $VERBOSE_FLAG
            echo ""
        else
            echo "⚠️ 跳过块内容和更新测试，请在 debug.env 中设置 TEST_DOCUMENT_ID 和 TEST_BLOCK_ID"
        fi
        
        if [ -n "$TEST_DOCUMENT_ID" ]; then
            echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
            echo "📝 11. 测试创建文本块"
            echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
            ./build/debug -app-id="$APP_ID" -app-secret="$APP_SECRET" -base-url="$BASE_URL" -action=create-text -doc-id="$TEST_DOCUMENT_ID" -parent-id="$TEST_PARENT_BLOCK_ID" -content="自动测试创建的文本块" $DEBUG_FLAG $VERBOSE_FLAG
            echo ""
            
            echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
            echo "💻 12. 测试创建代码块"
            echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
            ./build/debug -app-id="$APP_ID" -app-secret="$APP_SECRET" -base-url="$BASE_URL" -action=create-code -doc-id="$TEST_DOCUMENT_ID" -parent-id="$TEST_PARENT_BLOCK_ID" -code="console.log('自动测试代码块');" -language="javascript" $DEBUG_FLAG $VERBOSE_FLAG
            echo ""
            
            echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
            echo "📋 13. 测试创建标题块"
            echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
            ./build/debug -app-id="$APP_ID" -app-secret="$APP_SECRET" -base-url="$BASE_URL" -action=create-heading -doc-id="$TEST_DOCUMENT_ID" -parent-id="$TEST_PARENT_BLOCK_ID" -content="自动测试标题" -level="2" $DEBUG_FLAG $VERBOSE_FLAG
            echo ""
            
            echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
            echo "📃 14. 测试创建列表块"
            echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
            ./build/debug -app-id="$APP_ID" -app-secret="$APP_SECRET" -base-url="$BASE_URL" -action=create-list -doc-id="$TEST_DOCUMENT_ID" -parent-id="$TEST_PARENT_BLOCK_ID" -items="测试项目1,测试项目2,测试项目3" -list-type="bullet" $DEBUG_FLAG $VERBOSE_FLAG
            echo ""
            
            echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
            echo "📦 15. 测试批量创建块"
            echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
            ./build/debug -app-id="$APP_ID" -app-secret="$APP_SECRET" -base-url="$BASE_URL" -action=batch-create -doc-id="$TEST_DOCUMENT_ID" -parent-id="$TEST_PARENT_BLOCK_ID" $DEBUG_FLAG $VERBOSE_FLAG
            echo ""
        else
            echo "⚠️ 跳过块创建测试，请在 debug.env 中设置 TEST_DOCUMENT_ID"
        fi
        
        if [ -n "$TEST_DOCUMENT_ID" ]; then
            echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
            echo "🗑️ 16. 测试删除文档块 (谨慎操作)"
            echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
            echo "⚠️ 注意：此操作会删除文档中的块，请确保测试文档可以被修改"
            ./build/debug -app-id="$APP_ID" -app-secret="$APP_SECRET" -base-url="$BASE_URL" -action=delete-blocks -doc-id="$TEST_DOCUMENT_ID" -start-index="0" -end-index="0" $DEBUG_FLAG $VERBOSE_FLAG
            echo ""
        fi
        
        # 知识库操作测试
        echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
        echo "📚 17. 测试获取知识库空间列表"
        echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
        ./build/debug -app-id="$APP_ID" -app-secret="$APP_SECRET" -base-url="$BASE_URL" -action=wiki-spaces -page-size=20 $DEBUG_FLAG $VERBOSE_FLAG
        echo ""
        
        if [ -n "$TEST_WIKI_SPACE_ID" ]; then
            echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
            echo "📁 18. 测试获取知识库节点列表"
            echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
            ./build/debug -app-id="$APP_ID" -app-secret="$APP_SECRET" -base-url="$BASE_URL" -action=wiki-nodes -space-id="$TEST_WIKI_SPACE_ID" -page-size=20 $DEBUG_FLAG $VERBOSE_FLAG
            echo ""
            
            if [ -n "$TEST_WIKI_NODE_TOKEN" ]; then
                echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
                echo "📖 19. 测试获取知识库节点内容"
                echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
                ./build/debug -app-id="$APP_ID" -app-secret="$APP_SECRET" -base-url="$BASE_URL" -action=wiki-content -space-id="$TEST_WIKI_SPACE_ID" -node-token="$TEST_WIKI_NODE_TOKEN" -lang=0 $DEBUG_FLAG $VERBOSE_FLAG
                echo ""
                
                echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
                echo "ℹ️ 20. 测试获取知识库节点元信息"
                echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
                ./build/debug -app-id="$APP_ID" -app-secret="$APP_SECRET" -base-url="$BASE_URL" -action=wiki-meta -space-id="$TEST_WIKI_SPACE_ID" -node-token="$TEST_WIKI_NODE_TOKEN" $DEBUG_FLAG $VERBOSE_FLAG
                echo ""
            else
                echo "⚠️ 跳过知识库节点内容测试，请在 debug.env 中设置 TEST_WIKI_NODE_TOKEN"
            fi
        else
            echo "⚠️ 跳过知识库节点测试，请在 debug.env 中设置 TEST_WIKI_SPACE_ID"
        fi
        
        echo "🎉 所有API测试完成!"
        ;;
    *)
        echo "飞书MCP Go API调试工具"
        echo ""
        echo "用法: $0 <action> [options]"
        echo ""
        echo "支持的操作:"
        echo "基础操作:"
        echo "  token       - 测试获取访问令牌"
        echo "  health      - 执行健康检查"
        echo "  validate    - 验证API端点"
        echo "文档操作:"
        echo "  create-doc  - 测试创建文档"
        echo "  get-doc     - 测试获取文档信息 <doc_id>"
        echo "  get-content - 测试获取文档内容 <doc_id> [lang]"
        echo "  get-blocks  - 测试获取文档块 <doc_id>"
        echo "  search      - 测试搜索文档 <search_key>"
        echo "内容操作:"
        echo "  get-block-content     - 测试获取块内容 <doc_id> <block_id>"
        echo "  update-block-text     - 测试更新块文本 <doc_id> <block_id> [content]"
        echo "  create-text-block     - 测试创建文本块 <doc_id> [parent_id] [content]"
        echo "  create-code-block     - 测试创建代码块 <doc_id> [parent_id] [code] [language]"
        echo "  create-heading-block  - 测试创建标题块 <doc_id> [parent_id] [content] [level]"
        echo "  create-list-block     - 测试创建列表块 <doc_id> [parent_id] [items] [list_type]"
        echo "  batch-create-blocks   - 测试批量创建块 <doc_id> [parent_id]"
        echo "  delete-blocks         - 测试删除文档块 <doc_id> [start_index] [end_index]"
        echo "文件夹操作:"
        echo "  folder-info - 测试获取文件夹信息"
        echo "  folder      - 测试获取文件夹信息(别名)"
        echo "知识库操作:"
        echo "  wiki-spaces - 测试获取知识库空间列表 [page_size] [page_token]"
        echo "  wiki-nodes  - 测试获取知识库节点列表 <space_id> [page_size] [parent_node]"
        echo "  wiki-content- 测试获取知识库节点内容 <space_id> <node_token> [lang]"
        echo "  wiki-meta   - 测试获取知识库节点元信息 <space_id> <node_token>"
        echo "综合测试:"
        echo "  all         - 运行所有API测试"
        echo ""
        echo "配置文件: debug.env"
        echo "调试模式: $DEBUG_MODE"
        echo "详细模式: $VERBOSE_MODE"
        echo ""
        echo "示例:"
        echo "基础操作:"
        echo "  $0 token"
        echo "  $0 health"
        echo "文档操作:"
        echo "  $0 get-content BNTbdcURCoyPcHx4Cgzcz9xGn9e"
        echo "  $0 search 关键词"
        echo "内容操作:"
        echo "  $0 get-block-content doc_id_123 block_id_456"
        echo "  $0 update-block-text doc_id_123 block_id_456 '新的文本内容'"
        echo "  $0 create-text-block doc_id_123 parent_id_789 '这是新的文本块'"
        echo "  $0 create-code-block doc_id_123 parent_id_789 'console.log(\"Hello\");' javascript"
        echo "  $0 create-heading-block doc_id_123 parent_id_789 '标题文本' 2"
        echo "  $0 create-list-block doc_id_123 parent_id_789 '项目1,项目2,项目3' bullet"
        echo "知识库操作:"
        echo "  $0 wiki-spaces"
        echo "  $0 wiki-nodes 7523019799962943492"
        echo "  $0 wiki-content 7523019799962943492 token123"
        echo "综合测试:"
        echo "  $0 all"
        exit 1
        ;;
esac 