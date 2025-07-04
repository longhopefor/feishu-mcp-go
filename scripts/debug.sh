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
        if [ -z "$TEST_DOCUMENT_ID" ]; then
            echo "❌ 请在 debug.env 中设置 TEST_DOCUMENT_ID"
            exit 1
        fi
        echo "📄 测试获取文档信息: $TEST_DOCUMENT_ID"
        ./build/debug -app-id="$APP_ID" -app-secret="$APP_SECRET" -base-url="$BASE_URL" -action=get-doc -doc-id="$TEST_DOCUMENT_ID" $DEBUG_FLAG $VERBOSE_FLAG
        ;;
    "get-content")
        if [ -z "$TEST_DOCUMENT_ID" ]; then
            echo "❌ 请在 debug.env 中设置 TEST_DOCUMENT_ID"
            exit 1
        fi
        LANG="${1:-0}"
        echo "📖 测试获取文档内容: $TEST_DOCUMENT_ID (语言: $LANG)"
        ./build/debug -app-id="$APP_ID" -app-secret="$APP_SECRET" -base-url="$BASE_URL" -action=get-content -doc-id="$TEST_DOCUMENT_ID" -lang="$LANG" $DEBUG_FLAG $VERBOSE_FLAG
        ;;
    "get-blocks")
        if [ -z "$TEST_DOCUMENT_ID" ]; then
            echo "❌ 请在 debug.env 中设置 TEST_DOCUMENT_ID"
            exit 1
        fi
        echo "🧱 测试获取文档块: $TEST_DOCUMENT_ID"
        ./build/debug -app-id="$APP_ID" -app-secret="$APP_SECRET" -base-url="$BASE_URL" -action=get-blocks -doc-id="$TEST_DOCUMENT_ID" $DEBUG_FLAG $VERBOSE_FLAG
        ;;
    "search")
        SEARCH_KEY="${1:-test}"
        echo "🔍 测试搜索文档: $SEARCH_KEY"
        ./build/debug -app-id="$APP_ID" -app-secret="$APP_SECRET" -base-url="$BASE_URL" -action=search -search-key="$SEARCH_KEY" $DEBUG_FLAG $VERBOSE_FLAG
        ;;
    "folder-info"|"folder")
        echo "📁 测试获取文件夹信息..."
        ./build/debug -app-id="$APP_ID" -app-secret="$APP_SECRET" -base-url="$BASE_URL" -action=folder-info $DEBUG_FLAG $VERBOSE_FLAG
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
        
        echo "🎉 所有API测试完成!"
        ;;
    *)
        echo "飞书MCP Go API调试工具"
        echo ""
        echo "用法: $0 <action> [options]"
        echo ""
        echo "支持的操作:"
        echo "  token       - 测试获取访问令牌"
        echo "  health      - 执行健康检查"
        echo "  validate    - 验证API端点"
        echo "  create-doc  - 测试创建文档"
        echo "  get-doc     - 测试获取文档信息"
        echo "  get-content - 测试获取文档内容"
        echo "  get-blocks  - 测试获取文档块"
        echo "  search      - 测试搜索文档"
        echo "  folder-info - 测试获取文件夹信息"
        echo "  folder      - 测试获取文件夹信息(别名)"
        echo "  all         - 运行所有API测试"
        echo ""
        echo "配置文件: debug.env"
        echo "调试模式: $DEBUG_MODE"
        echo "详细模式: $VERBOSE_MODE"
        echo ""
        echo "示例:"
        echo "  $0 token"
        echo "  $0 health"
        echo "  $0 search 关键词"
        echo "  $0 all"
        exit 1
        ;;
esac 