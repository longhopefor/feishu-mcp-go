#!/bin/bash

# 测试MCP连接的脚本
echo "=== 测试 Feishu MCP Server 连接 ==="

# 设置环境变量
export FEISHU_MCP_FEISHU_APP_ID="test"
export FEISHU_MCP_FEISHU_APP_SECRET="test"

# 创建临时文件来测试JSON RPC通信
cat > mcp_test_input.json << 'EOF'
{"jsonrpc": "2.0", "id": 1, "method": "initialize", "params": {"protocolVersion": "2024-11-05", "capabilities": {"tools": {}}, "clientInfo": {"name": "test-client", "version": "1.0.0"}}}
{"jsonrpc": "2.0", "id": 2, "method": "tools/list", "params": {}}
{"jsonrpc": "2.0", "id": 3, "method": "ping", "params": {}}
EOF

echo "测试输入JSON文件内容:"
cat mcp_test_input.json
echo ""
echo "=== 启动MCP服务器并发送测试请求 ==="

# 启动服务器并发送测试输入
./build/feishu-mcp --stdio < mcp_test_input.json 2>&1 &
SERVER_PID=$!
sleep 3
kill $SERVER_PID 2>/dev/null || echo "服务器已停止"
wait $SERVER_PID 2>/dev/null || echo "测试完成"

# 清理临时文件
rm -f mcp_test_input.json

echo ""
echo "=== 测试结束 ===" 