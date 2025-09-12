#!/bin/bash

# Eino 框架智能客服系统演示脚本

BASE_URL="http://localhost:8888"

echo "=========================================="
echo "    Eino 框架智能客服系统演示"
echo "=========================================="
echo

# 检查服务是否运行
echo "1. 检查服务健康状态..."
curl -s "$BASE_URL/api/v1/health/check" | jq .
echo
echo

# 测试 Chain 编排
echo "2. 测试 Chain 编排（简单链式流程）..."
echo "Chain 编排特点：线性执行，数据单向流动，适合简单问答"
echo

echo "2.1 Chain 同步调用 - 简单问答："
curl -X POST "$BASE_URL/api/v1/eino/chain/invoke" \
  -H "Content-Type: application/json" \
  -d '{"message": "你好，请介绍一下自己"}' | jq .
echo
echo

echo "2.2 Chain 流式调用 - 创意写作："
curl -X POST "$BASE_URL/api/v1/eino/chain/stream" \
  -H "Content-Type: application/json" \
  -d '{"message": "请写一首关于春天的短诗"}'
echo
echo

# 测试 Graph 编排
echo "3. 测试 Graph 编排（复杂有向图流程）..."
echo "Graph 编排特点：支持分支、循环和条件判断，适合复杂业务逻辑"
echo

echo "3.1 Graph 同步调用 - 问候意图："
curl -X POST "$BASE_URL/api/v1/eino/graph/invoke" \
  -H "Content-Type: application/json" \
  -d '{"message": "你好，很高兴见到你"}' | jq .
echo
echo

echo "3.2 Graph 同步调用 - 问题意图："
curl -X POST "$BASE_URL/api/v1/eino/graph/invoke" \
  -H "Content-Type: application/json" \
  -d '{"message": "什么是人工智能？"}' | jq .
echo
echo

echo "3.3 Graph 同步调用 - 投诉意图："
curl -X POST "$BASE_URL/api/v1/eino/graph/invoke" \
  -H "Content-Type: application/json" \
  -d '{"message": "我对你们的服务不满意，要投诉"}' | jq .
echo
echo

echo "3.4 Graph 同步调用 - 转人工意图："
curl -X POST "$BASE_URL/api/v1/eino/graph/invoke" \
  -H "Content-Type: application/json" \
  -d '{"message": "请转接人工客服"}' | jq .
echo
echo

echo "3.5 Graph 同步调用 - 创作意图："
curl -X POST "$BASE_URL/api/v1/eino/graph/invoke" \
  -H "Content-Type: application/json" \
  -d '{"message": "请写一首关于春天的诗"}' | jq .
echo
echo

echo "3.6 Graph 流式调用 - 创作流式响应："
curl -X POST "$BASE_URL/api/v1/eino/graph/stream" \
  -H "Content-Type: application/json" \
  -d '{"message": "请写一首关于春天的短诗"}'
echo
echo

# 测试 Workflow 编排
echo "4. 测试 Workflow 编排（高级编排流程）..."
echo "Workflow 编排特点：支持复杂数据流转和字段级映射，适合企业级应用"
echo

echo "4.1 Workflow 同步调用 - 复杂工作流："
curl -X POST "$BASE_URL/api/v1/eino/workflow/invoke" \
  -H "Content-Type: application/json" \
  -d '{"message": "你好，请介绍一下自己"}' | jq .
echo
echo

echo "4.2 Workflow 同步调用 - 问题处理工作流："
curl -X POST "$BASE_URL/api/v1/eino/workflow/invoke" \
  -H "Content-Type: application/json" \
  -d '{"message": "什么是人工智能？"}' | jq .
echo
echo

echo "4.3 Workflow 流式调用 - 复杂工作流流式响应："
curl -X POST "$BASE_URL/api/v1/eino/workflow/stream" \
  -H "Content-Type: application/json" \
  -d '{"message": "请写一首关于春天的短诗"}'
echo
echo

# 性能对比测试
echo "5. 性能对比测试..."
echo "测试不同编排方式的响应时间："
echo

echo "5.1 Chain 编排响应时间："
time curl -s -X POST "$BASE_URL/api/v1/eino/chain/invoke" \
  -H "Content-Type: application/json" \
  -d '{"message": "你好"}' > /dev/null

echo "5.2 Graph 编排响应时间："
time curl -s -X POST "$BASE_URL/api/v1/eino/graph/invoke" \
  -H "Content-Type: application/json" \
  -d '{"message": "你好"}' > /dev/null

echo "5.3 Workflow 编排响应时间："
time curl -s -X POST "$BASE_URL/api/v1/eino/workflow/invoke" \
  -H "Content-Type: application/json" \
  -d '{"message": "你好"}' > /dev/null

echo
echo "=========================================="
echo "           演示总结"
echo "=========================================="
echo
echo "✅ Chain 编排："
echo "   - 简单链式流程，线性执行"
echo "   - 适合简单问答和创意写作"
echo "   - 执行效率高，延迟低"
echo
echo "✅ Graph 编排："
echo "   - 复杂有向图流程，支持分支和循环"
echo "   - 支持意图识别和智能路由"
echo "   - 适合智能客服和复杂业务逻辑"
echo
echo "✅ Workflow 编排："
echo "   - 高级编排流程，支持字段级映射"
echo "   - 支持复杂数据流转和状态管理"
echo "   - 适合企业级应用和高可观测性要求"
echo
echo "🎯 Eino 框架核心功能："
echo "   - 组件化设计：可复用的组件构建复杂应用"
echo "   - 编排能力：三种不同复杂度的编排方式"
echo "   - 数据流转：支持复杂数据在节点间的流转和映射"
echo "   - 流式处理：自动处理流式数据的拼接和合并"
echo "   - 状态管理：全局状态管理和会话维护"
echo "   - 可观测性：完整的日志记录和监控"
echo
echo "🚀 这是一个真正基于 eino 框架的中等复杂度项目！"
echo "=========================================="
