#!/bin/bash

# LLM Demo API 测试脚本

BASE_URL="http://localhost:8888"

echo "=== LLM Demo API 测试脚本 ==="
echo

# 检查服务是否运行
echo "1. 检查服务健康状态..."
curl -s "$BASE_URL/api/v1/health/check" | jq .
echo
echo

# 测试同步聊天接口
echo "2. 测试同步聊天接口..."
echo "发送消息: 你好，请简单介绍一下自己"
curl -X POST "$BASE_URL/api/v1/chat/invoke" \
  -H "Content-Type: application/json" \
  -d '{"message": "你好，请简单介绍一下自己"}' | jq .
echo
echo

# 测试流式聊天接口
echo "3. 测试流式聊天接口..."
echo "发送消息: 请写一首关于春天的短诗"
echo "流式响应:"
curl -X POST "$BASE_URL/api/v1/chat/stream" \
  -H "Content-Type: application/json" \
  -d '{"message": "请写一首关于春天的短诗"}'
echo
echo

# 测试指定模型
echo "4. 测试指定模型..."
echo "发送消息: 什么是人工智能？"
curl -X POST "$BASE_URL/api/v1/chat/invoke" \
  -H "Content-Type: application/json" \
  -d '{"message": "什么是人工智能？", "model": "glm-4-0520"}' | jq .
echo
echo

echo "=== 测试完成 ==="
