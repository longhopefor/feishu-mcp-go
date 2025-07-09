#!/bin/bash

# 飞书MCP快速开始脚本
# 自动化安装和配置流程

set -e

echo "🚀 === 飞书MCP快速开始脚本 ==="
echo "本脚本将帮助您快速配置和启动飞书MCP服务器"
echo ""

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 日志函数
log_info() {
    echo -e "${BLUE}ℹ️  $1${NC}"
}

log_success() {
    echo -e "${GREEN}✅ $1${NC}"
}

log_warning() {
    echo -e "${YELLOW}⚠️  $1${NC}"
}

log_error() {
    echo -e "${RED}❌ $1${NC}"
}

# 检查系统依赖
check_dependencies() {
    log_info "检查系统依赖..."
    
    # 检查Go
    if ! command -v go &> /dev/null; then
        log_error "Go未安装或不在PATH中"
        log_info "请访问 https://golang.org/dl/ 下载安装Go 1.21+"
        exit 1
    fi
    
    GO_VERSION=$(go version | grep -o 'go[0-9]\+\.[0-9]\+' | sed 's/go//')
    log_success "Go版本: $GO_VERSION"
    
    # 检查Git
    if ! command -v git &> /dev/null; then
        log_error "Git未安装"
        exit 1
    fi
    log_success "Git已安装"
    
    # 检查make
    if ! command -v make &> /dev/null; then
        log_warning "Make未安装，将使用go build"
    else
        log_success "Make已安装"
    fi
}

# 获取飞书配置
get_feishu_config() {
    log_info "配置飞书应用凭证"
    echo ""
    
    # 检查是否已有配置
    if [ -f ".env" ]; then
        log_warning "发现现有的.env配置文件"
        read -p "是否使用现有配置？(y/n): " use_existing
        if [[ $use_existing =~ ^[Yy]$ ]]; then
            source .env
            return
        fi
    fi
    
    echo "请准备您的飞书应用凭证："
    echo "1. 访问 https://open.feishu.cn/"
    echo "2. 创建企业自建应用"
    echo "3. 配置必要权限（docx:document, drive:file等）"
    echo "4. 获取App ID和App Secret"
    echo ""
    
    # 获取App ID
    while true; do
        read -p "请输入飞书App ID (cli_开头): " APP_ID
        if [[ $APP_ID =~ ^cli_ ]]; then
            break
        else
            log_error "App ID格式不正确，应该以'cli_'开头"
        fi
    done
    
    # 获取App Secret
    while true; do
        read -s -p "请输入飞书App Secret: " APP_SECRET
        echo ""
        if [ ${#APP_SECRET} -gt 20 ]; then
            break
        else
            log_error "App Secret长度不足，请检查是否正确"
        fi
    done
    
    # 创建.env文件
    cat > .env << EOF
# 飞书应用配置
export FEISHU_MCP_FEISHU_APP_ID="$APP_ID"
export FEISHU_MCP_FEISHU_APP_SECRET="$APP_SECRET"

# 可选配置
export FEISHU_MCP_LOG_LEVEL="info"
export FEISHU_MCP_CACHE_ENABLED="true"
export FEISHU_MCP_SERVER_MODE="stdio"
EOF
    
    log_success "配置文件已创建: .env"
    
    # 加载环境变量
    source .env
}

# 构建项目
build_project() {
    log_info "构建项目..."
    
    # 下载依赖
    go mod download
    go mod tidy
    
    # 构建
    if command -v make &> /dev/null; then
        make build
    else
        mkdir -p build
        go build -o build/feishu-mcp ./cmd/server/main.go
    fi
    
    # 检查构建结果
    if [ -f "build/feishu-mcp" ]; then
        chmod +x build/feishu-mcp
        log_success "项目构建成功"
    else
        log_error "项目构建失败"
        exit 1
    fi
}

# 测试连接
test_connection() {
    log_info "测试飞书API连接..."
    
    # 确保环境变量已加载
    source .env
    
    # 测试令牌获取
    if go run ./cmd/debug/main.go -app-id="$FEISHU_MCP_FEISHU_APP_ID" -app-secret="$FEISHU_MCP_FEISHU_APP_SECRET" -action=token >/dev/null 2>&1; then
        log_success "飞书API连接测试成功"
    else
        log_error "飞书API连接失败，请检查配置"
        log_info "尝试运行调试命令："
        echo "go run ./cmd/debug/main.go -app-id=\"\$FEISHU_MCP_FEISHU_APP_ID\" -app-secret=\"\$FEISHU_MCP_FEISHU_APP_SECRET\" -action=token -debug"
        exit 1
    fi
}

# 生成MCP配置
generate_mcp_config() {
    log_info "生成MCP配置示例..."
    
    # 获取当前绝对路径
    CURRENT_PATH=$(pwd)
    
    # 生成Cursor配置
    cat > cursor-mcp-config.json << EOF
{
  "mcpServers": {
    "feishu": {
      "command": "$CURRENT_PATH/build/feishu-mcp",
      "args": ["--stdio"],
      "env": {
        "FEISHU_MCP_FEISHU_APP_ID": "$FEISHU_MCP_FEISHU_APP_ID",
        "FEISHU_MCP_FEISHU_APP_SECRET": "$FEISHU_MCP_FEISHU_APP_SECRET",
        "FEISHU_MCP_LOG_LEVEL": "info"
      }
    }
  }
}
EOF
    
    log_success "MCP配置已生成: cursor-mcp-config.json"
}

# 显示使用说明
show_usage_instructions() {
    echo ""
    log_success "🎉 快速开始完成！"
    echo ""
    
    log_info "下一步操作："
    echo "1. 配置AI工具（Cursor/Windsurf/Cline）："
    echo "   - 复制 cursor-mcp-config.json 中的配置"
    echo "   - 添加到AI工具的MCP设置中"
    echo "   - 重启AI工具"
    echo ""
    
    echo "2. 测试功能："
    echo "   - 在AI工具中说：'请帮我创建一个飞书文档'"
    echo "   - 或运行全面测试：./comprehensive_api_test.sh"
    echo ""
    
    echo "3. 常用命令："
    echo "   - 手动启动服务器：./build/feishu-mcp --stdio"
    echo "   - 调试模式：./build/feishu-mcp --log-level=debug --stdio"
    echo "   - 测试连接：source .env && go run ./cmd/debug/main.go -action=health"
    echo ""
    
    log_info "详细使用说明请查看：LOCAL_USAGE_GUIDE.md"
    
    echo ""
    log_warning "重要提醒："
    echo "- 确保飞书应用已发布且权限充足"
    echo "- 妥善保管.env文件中的凭证信息"
    echo "- 使用绝对路径配置MCP服务器"
}

# 主执行流程
main() {
    echo "开始快速安装和配置..."
    echo ""
    
    # 1. 检查依赖
    check_dependencies
    echo ""
    
    # 2. 获取配置
    get_feishu_config
    echo ""
    
    # 3. 构建项目
    build_project
    echo ""
    
    # 4. 测试连接
    test_connection
    echo ""
    
    # 5. 生成配置
    generate_mcp_config
    echo ""
    
    # 6. 显示说明
    show_usage_instructions
}

# 错误处理
trap 'log_error "脚本执行失败，请检查错误信息"; exit 1' ERR

# 确认执行
read -p "是否开始快速配置飞书MCP？(y/n): " -n 1 -r
echo ""
if [[ $REPLY =~ ^[Yy]$ ]]; then
    main
else
    echo "取消配置"
    exit 0
fi 