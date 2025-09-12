#!/bin/bash

# Eino 框架编排 API 测试脚本

BASE_URL="http://localhost:8888"

echo "=== Eino 框架编排 API 测试脚本 ==="
echo

# 检查服务是否运行
echo "1. 检查服务健康状态..."
curl -s "$BASE_URL/api/v1/health/check" | jq .
echo
echo

# 测试传统聊天接口
echo "2. 测试传统聊天接口..."
echo "发送消息: 你好，请简单介绍一下自己"
curl -X POST "$BASE_URL/api/v1/chat/invoke" \
  -H "Content-Type: application/json" \
  -d '{"message": "你好，请简单介绍一下自己"}' | jq .
echo
echo

# 测试 Chain 编排
echo "3. 测试 Chain 编排（简单链式流程）..."
echo "Chain 同步调用:"
curl -X POST "$BASE_URL/api/v1/eino/chain/invoke" \
  -H "Content-Type: application/json" \
  -d '{"message": "你好，请介绍一下自己"}' | jq .
echo
echo

echo "Chain 流式调用:"
curl -X POST "$BASE_URL/api/v1/eino/chain/stream" \
  -H "Content-Type: application/json" \
  -d '{"message": "请写一首关于春天的短诗"}'
echo
echo

# 测试 Graph 编排
echo "4. 测试 Graph 编排（复杂有向图流程）..."
echo "Graph 同步调用 - 问候意图:"
curl -X POST "$BASE_URL/api/v1/eino/graph/invoke" \
  -H "Content-Type: application/json" \
  -d '{"message": "你好，很高兴见到你"}' | jq .
echo
echo

echo "Graph 同步调用 - 问题意图:"
curl -X POST "$BASE_URL/api/v1/eino/graph/invoke" \
  -H "Content-Type: application/json" \
  -d '{"message": "什么是人工智能？"}' | jq .
echo
echo

echo "Graph 同步调用 - 创作意图:"
curl -X POST "$BASE_URL/api/v1/eino/graph/invoke" \
  -H "Content-Type: application/json" \
  -d '{"message": "请写一首关于春天的诗"}' | jq .
echo
echo

echo "Graph 流式调用:"
curl -X POST "$BASE_URL/api/v1/eino/graph/stream" \
  -H "Content-Type: application/json" \
  -d '{"message": "请写一首关于春天的短诗"}'
echo
echo

# 测试 Workflow 编排
echo "5. 测试 Workflow 编排（高级编排流程）..."
echo "Workflow 同步调用:"
curl -X POST "$BASE_URL/api/v1/eino/workflow/invoke" \
  -H "Content-Type: application/json" \
  -d '{"message": "你好，请介绍一下自己"}' | jq .
echo
echo

echo "Workflow 流式调用:"
curl -X POST "$BASE_URL/api/v1/eino/workflow/stream" \
  -H "Content-Type: application/json" \
  -d '{"message": "请写一首关于春天的短诗"}'
echo
echo

# 性能对比测试
echo "6. 性能对比测试..."
echo "测试不同编排方式的响应时间..."

echo "传统聊天接口:"
time curl -s -X POST "$BASE_URL/api/v1/chat/invoke" \
  -H "Content-Type: application/json" \
  -d '{"message": "你好"}' > /dev/null

echo "Chain 编排:"
time curl -s -X POST "$BASE_URL/api/v1/eino/chain/invoke" \
  -H "Content-Type: application/json" \
  -d '{"message": "你好"}' > /dev/null

echo "Graph 编排:"
time curl -s -X POST "$BASE_URL/api/v1/eino/graph/invoke" \
  -H "Content-Type: application/json" \
  -d '{"message": "你好"}' > /dev/null

echo "Workflow 编排:"
time curl -s -X POST "$BASE_URL/api/v1/eino/workflow/invoke" \
  -H "Content-Type: application/json" \
  -d '{"message": "你好"}' > /dev/null

echo
echo "=== 测试完成 ==="
echo
echo "总结："
echo "- Chain 编排：简单链式流程，适合线性对话"
echo "- Graph 编排：复杂有向图流程，支持意图识别和分支处理"
echo "- Workflow 编排：高级编排流程，支持复杂数据流转和字段级映射"
echo "- 所有编排方式都支持同步和流式两种调用模式"
