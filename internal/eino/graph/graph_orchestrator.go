package graph

import (
	"context"
	"fmt"
	"io"
	"log"

	"github.com/UTC-Six/llm/internal/config"
	"github.com/UTC-Six/llm/internal/eino/components"
	"github.com/UTC-Six/llm/internal/types"
)

// GraphOrchestrator eino 框架的 Graph 编排器
// Graph 是复杂的有向图编排，支持分支、循环和条件判断
type GraphOrchestrator struct {
	config             config.Config
	chatModel          *components.ChatModel
	intentClassifier   *components.IntentClassifier
	contextManager     *components.ContextManager
	knowledgeRetriever *components.KnowledgeRetriever
	responseOptimizer  *components.ResponseOptimizer
}

// NewGraphOrchestrator 创建 Graph 编排器
func NewGraphOrchestrator(c config.Config) *GraphOrchestrator {
	return &GraphOrchestrator{
		config:             c,
		chatModel:          components.NewChatModel(c.LLM.GLM4),
		intentClassifier:   components.NewIntentClassifier(),
		contextManager:     components.NewContextManager(),
		knowledgeRetriever: components.NewKnowledgeRetriever(),
		responseOptimizer:  components.NewResponseOptimizer(),
	}
}

// ChatInvoke Graph 同步聊天
// 使用 eino 框架的 Graph 编排：用户输入 -> 意图识别 -> 分支处理 -> 上下文管理 -> 知识检索 -> 模型调用 -> 响应优化
func (g *GraphOrchestrator) ChatInvoke(ctx context.Context, req *types.ChatRequest) (*types.ChatResponse, error) {
	log.Printf("eino Graph 编排开始同步聊天: %s", req.Message)

	// Graph 节点1: 意图识别
	intentResult, err := g.intentClassifier.Invoke(ctx, req.Message)
	if err != nil {
		return nil, fmt.Errorf("eino Graph 意图识别失败: %v", err)
	}
	log.Printf("eino Graph 节点1 - 意图识别完成")

	// 解析意图结果
	intentData, ok := intentResult.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("eino Graph 意图结果格式错误")
	}
	intent, _ := intentData["intent"].(string)

	// Graph 节点2: 上下文管理
	contextInput := map[string]interface{}{
		"user_id":    "user_123",
		"session_id": "session_456",
		"message": components.ChatMessage{
			Role:    "user",
			Content: req.Message,
		},
	}

	_, err = g.contextManager.Invoke(ctx, contextInput)
	if err != nil {
		return nil, fmt.Errorf("eino Graph 上下文管理失败: %v", err)
	}
	log.Printf("eino Graph 节点2 - 上下文管理完成")

	// Graph 节点3: 根据意图选择分支处理
	var response *types.ChatResponse
	switch intent {
	case "greeting":
		response, err = g.handleGreetingBranch(ctx, req.Message, intentData)
	case "question":
		response, err = g.handleQuestionBranch(ctx, req.Message, intentData)
	case "complaint":
		response, err = g.handleComplaintBranch(ctx, req.Message, intentData)
	case "transfer":
		response, err = g.handleTransferBranch(ctx, req.Message, intentData)
	case "creative":
		response, err = g.handleCreativeBranch(ctx, req.Message, intentData)
	default:
		response, err = g.handleDefaultBranch(ctx, req.Message, intentData)
	}

	if err != nil {
		return nil, fmt.Errorf("eino Graph 分支处理失败: %v", err)
	}

	log.Printf("eino Graph 编排同步聊天完成，响应长度: %d", len(response.Data))
	return response, nil
}

// ChatStream Graph 流式聊天
// 使用 eino 框架的 Graph 编排进行流式响应
func (g *GraphOrchestrator) ChatStream(ctx context.Context, req *types.StreamChatRequest, writer io.Writer) error {
	log.Printf("eino Graph 编排开始流式聊天: %s", req.Message)

	// Graph 节点1: 意图识别
	intentResult, err := g.intentClassifier.Invoke(ctx, req.Message)
	if err != nil {
		return fmt.Errorf("eino Graph 意图识别失败: %v", err)
	}

	// 解析意图结果
	intentData, ok := intentResult.(map[string]interface{})
	if !ok {
		return fmt.Errorf("eino Graph 意图结果格式错误")
	}
	intent, _ := intentData["intent"].(string)

	// Graph 节点2: 根据意图选择分支处理
	switch intent {
	case "greeting":
		err = g.handleGreetingStream(ctx, req.Message, writer)
	case "question":
		err = g.handleQuestionStream(ctx, req.Message, writer)
	case "complaint":
		err = g.handleComplaintStream(ctx, req.Message, writer)
	case "transfer":
		err = g.handleTransferStream(ctx, req.Message, writer)
	case "creative":
		err = g.handleCreativeStream(ctx, req.Message, writer)
	default:
		err = g.handleDefaultStream(ctx, req.Message, writer)
	}

	if err != nil {
		return fmt.Errorf("eino Graph 分支处理失败: %v", err)
	}

	log.Printf("eino Graph 编排流式聊天完成")
	return nil
}

// handleGreetingBranch 处理问候分支
func (g *GraphOrchestrator) handleGreetingBranch(ctx context.Context, message string, intentData map[string]interface{}) (*types.ChatResponse, error) {
	log.Printf("eino Graph 分支 - 处理问候意图")

	// 构建问候消息
	messages := []components.ChatMessage{
		{
			Role:    "system",
			Content: "你是一个友好的智能客服助手，请用温暖、热情的方式回应用户的问候。",
		},
		{
			Role:    "user",
			Content: message,
		},
	}

	// 调用模型
	result, err := g.chatModel.Invoke(ctx, messages)
	if err != nil {
		return nil, err
	}

	chatResp, ok := result.(*components.ChatResponse)
	if !ok {
		return nil, fmt.Errorf("模型结果格式错误")
	}

	// 响应优化
	optimizerInput := map[string]interface{}{
		"intent":   "greeting",
		"response": chatResp.Choices[0].Message.Content,
	}

	optimizerResult, err := g.responseOptimizer.Invoke(ctx, optimizerInput)
	if err != nil {
		return nil, err
	}

	optimizedData, ok := optimizerResult.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("优化结果格式错误")
	}

	optimizedResponse, _ := optimizedData["optimized_response"].(string)

	return &types.ChatResponse{
		Code:    200,
		Message: "success",
		Data:    optimizedResponse,
		Model:   g.config.LLM.GLM4.Model,
		Usage: types.Usage{
			PromptTokens:     chatResp.Usage.PromptTokens,
			CompletionTokens: chatResp.Usage.CompletionTokens,
			TotalTokens:      chatResp.Usage.TotalTokens,
		},
	}, nil
}

// handleQuestionBranch 处理问题分支
func (g *GraphOrchestrator) handleQuestionBranch(ctx context.Context, message string, intentData map[string]interface{}) (*types.ChatResponse, error) {
	log.Printf("eino Graph 分支 - 处理问题意图")

	// Graph 子节点: 知识检索
	knowledgeResult, err := g.knowledgeRetriever.Invoke(ctx, message)
	if err != nil {
		return nil, err
	}

	knowledgeData, ok := knowledgeResult.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("知识检索结果格式错误")
	}

	knowledgeResults, _ := knowledgeData["results"].([]string)
	knowledgeText := ""
	for _, result := range knowledgeResults {
		knowledgeText += result + "\n"
	}

	// 构建问题回答消息
	messages := []components.ChatMessage{
		{
			Role:    "system",
			Content: "你是一个专业的智能客服助手，请基于以下知识库信息准确回答用户的问题：\n" + knowledgeText,
		},
		{
			Role:    "user",
			Content: message,
		},
	}

	// 调用模型
	result, err := g.chatModel.Invoke(ctx, messages)
	if err != nil {
		return nil, err
	}

	chatResp, ok := result.(*components.ChatResponse)
	if !ok {
		return nil, fmt.Errorf("模型结果格式错误")
	}

	// 响应优化
	optimizerInput := map[string]interface{}{
		"intent":   "question",
		"response": chatResp.Choices[0].Message.Content,
	}

	optimizerResult, err := g.responseOptimizer.Invoke(ctx, optimizerInput)
	if err != nil {
		return nil, err
	}

	optimizedData, ok := optimizerResult.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("优化结果格式错误")
	}

	optimizedResponse, _ := optimizedData["optimized_response"].(string)

	return &types.ChatResponse{
		Code:    200,
		Message: "success",
		Data:    optimizedResponse,
		Model:   g.config.LLM.GLM4.Model,
		Usage: types.Usage{
			PromptTokens:     chatResp.Usage.PromptTokens,
			CompletionTokens: chatResp.Usage.CompletionTokens,
			TotalTokens:      chatResp.Usage.TotalTokens,
		},
	}, nil
}

// handleComplaintBranch 处理投诉分支
func (g *GraphOrchestrator) handleComplaintBranch(ctx context.Context, message string, intentData map[string]interface{}) (*types.ChatResponse, error) {
	log.Printf("eino Graph 分支 - 处理投诉意图")

	// 构建投诉处理消息
	messages := []components.ChatMessage{
		{
			Role:    "system",
			Content: "你是一个专业的客服助手，请用诚恳、专业的方式处理用户的投诉，表达歉意并提供解决方案。",
		},
		{
			Role:    "user",
			Content: message,
		},
	}

	// 调用模型
	result, err := g.chatModel.Invoke(ctx, messages)
	if err != nil {
		return nil, err
	}

	chatResp, ok := result.(*components.ChatResponse)
	if !ok {
		return nil, fmt.Errorf("模型结果格式错误")
	}

	// 响应优化
	optimizerInput := map[string]interface{}{
		"intent":   "complaint",
		"response": chatResp.Choices[0].Message.Content,
	}

	optimizerResult, err := g.responseOptimizer.Invoke(ctx, optimizerInput)
	if err != nil {
		return nil, err
	}

	optimizedData, ok := optimizerResult.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("优化结果格式错误")
	}

	optimizedResponse, _ := optimizedData["optimized_response"].(string)

	return &types.ChatResponse{
		Code:    200,
		Message: "success",
		Data:    optimizedResponse,
		Model:   g.config.LLM.GLM4.Model,
		Usage: types.Usage{
			PromptTokens:     chatResp.Usage.PromptTokens,
			CompletionTokens: chatResp.Usage.CompletionTokens,
			TotalTokens:      chatResp.Usage.TotalTokens,
		},
	}, nil
}

// handleTransferBranch 处理转人工分支
func (g *GraphOrchestrator) handleTransferBranch(ctx context.Context, message string, intentData map[string]interface{}) (*types.ChatResponse, error) {
	log.Printf("eino Graph 分支 - 处理转人工意图")

	// 直接返回转人工响应
	response := "我理解您需要人工客服的帮助。正在为您转接到人工客服，请稍候...\n\n" +
		"在等待期间，您可以：\n" +
		"1. 描述您遇到的具体问题\n" +
		"2. 提供相关的订单号或账号信息\n" +
		"3. 说明您希望得到的帮助\n\n" +
		"我们的客服团队会尽快为您提供专业的服务。"

	return &types.ChatResponse{
		Code:    200,
		Message: "success",
		Data:    response,
		Model:   g.config.LLM.GLM4.Model,
		Usage: types.Usage{
			PromptTokens:     0,
			CompletionTokens: 0,
			TotalTokens:      0,
		},
	}, nil
}

// handleCreativeBranch 处理创作分支
func (g *GraphOrchestrator) handleCreativeBranch(ctx context.Context, message string, intentData map[string]interface{}) (*types.ChatResponse, error) {
	log.Printf("eino Graph 分支 - 处理创作意图")

	// 构建创作消息
	messages := []components.ChatMessage{
		{
			Role:    "system",
			Content: "你是一个富有创意的写作助手，请根据用户的要求创作生动、有趣的内容。",
		},
		{
			Role:    "user",
			Content: message,
		},
	}

	// 调用模型
	result, err := g.chatModel.Invoke(ctx, messages)
	if err != nil {
		return nil, err
	}

	chatResp, ok := result.(*components.ChatResponse)
	if !ok {
		return nil, fmt.Errorf("模型结果格式错误")
	}

	// 响应优化
	optimizerInput := map[string]interface{}{
		"intent":   "creative",
		"response": chatResp.Choices[0].Message.Content,
	}

	optimizerResult, err := g.responseOptimizer.Invoke(ctx, optimizerInput)
	if err != nil {
		return nil, err
	}

	optimizedData, ok := optimizerResult.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("优化结果格式错误")
	}

	optimizedResponse, _ := optimizedData["optimized_response"].(string)

	return &types.ChatResponse{
		Code:    200,
		Message: "success",
		Data:    optimizedResponse,
		Model:   g.config.LLM.GLM4.Model,
		Usage: types.Usage{
			PromptTokens:     chatResp.Usage.PromptTokens,
			CompletionTokens: chatResp.Usage.CompletionTokens,
			TotalTokens:      chatResp.Usage.TotalTokens,
		},
	}, nil
}

// handleDefaultBranch 处理默认分支
func (g *GraphOrchestrator) handleDefaultBranch(ctx context.Context, message string, intentData map[string]interface{}) (*types.ChatResponse, error) {
	log.Printf("eino Graph 分支 - 处理默认意图")

	// 构建默认消息
	messages := []components.ChatMessage{
		{
			Role:    "system",
			Content: "你是一个智能客服助手，请根据用户的需求提供合适的帮助。",
		},
		{
			Role:    "user",
			Content: message,
		},
	}

	// 调用模型
	result, err := g.chatModel.Invoke(ctx, messages)
	if err != nil {
		return nil, err
	}

	chatResp, ok := result.(*components.ChatResponse)
	if !ok {
		return nil, fmt.Errorf("模型结果格式错误")
	}

	// 响应优化
	optimizerInput := map[string]interface{}{
		"intent":   "default",
		"response": chatResp.Choices[0].Message.Content,
	}

	optimizerResult, err := g.responseOptimizer.Invoke(ctx, optimizerInput)
	if err != nil {
		return nil, err
	}

	optimizedData, ok := optimizerResult.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("优化结果格式错误")
	}

	optimizedResponse, _ := optimizedData["optimized_response"].(string)

	return &types.ChatResponse{
		Code:    200,
		Message: "success",
		Data:    optimizedResponse,
		Model:   g.config.LLM.GLM4.Model,
		Usage: types.Usage{
			PromptTokens:     chatResp.Usage.PromptTokens,
			CompletionTokens: chatResp.Usage.CompletionTokens,
			TotalTokens:      chatResp.Usage.TotalTokens,
		},
	}, nil
}

// 流式处理方法的实现
func (g *GraphOrchestrator) handleGreetingStream(ctx context.Context, message string, writer io.Writer) error {
	messages := []components.ChatMessage{
		{
			Role:    "system",
			Content: "你是一个友好的智能客服助手，请用温暖、热情的方式回应用户的问候。",
		},
		{
			Role:    "user",
			Content: message,
		},
	}
	return g.chatModel.Stream(ctx, messages, writer)
}

func (g *GraphOrchestrator) handleQuestionStream(ctx context.Context, message string, writer io.Writer) error {
	messages := []components.ChatMessage{
		{
			Role:    "system",
			Content: "你是一个专业的智能客服助手，请准确回答用户的问题。",
		},
		{
			Role:    "user",
			Content: message,
		},
	}
	return g.chatModel.Stream(ctx, messages, writer)
}

func (g *GraphOrchestrator) handleComplaintStream(ctx context.Context, message string, writer io.Writer) error {
	messages := []components.ChatMessage{
		{
			Role:    "system",
			Content: "你是一个专业的客服助手，请用诚恳、专业的方式处理用户的投诉。",
		},
		{
			Role:    "user",
			Content: message,
		},
	}
	return g.chatModel.Stream(ctx, messages, writer)
}

func (g *GraphOrchestrator) handleTransferStream(ctx context.Context, message string, writer io.Writer) error {
	response := "正在为您转接到人工客服，请稍候..."
	_, err := writer.Write([]byte(fmt.Sprintf("data: %s\n\n", response)))
	return err
}

func (g *GraphOrchestrator) handleCreativeStream(ctx context.Context, message string, writer io.Writer) error {
	messages := []components.ChatMessage{
		{
			Role:    "system",
			Content: "你是一个富有创意的写作助手，请根据用户的要求创作生动、有趣的内容。",
		},
		{
			Role:    "user",
			Content: message,
		},
	}
	return g.chatModel.Stream(ctx, messages, writer)
}

func (g *GraphOrchestrator) handleDefaultStream(ctx context.Context, message string, writer io.Writer) error {
	messages := []components.ChatMessage{
		{
			Role:    "system",
			Content: "你是一个智能客服助手，请根据用户的需求提供合适的帮助。",
		},
		{
			Role:    "user",
			Content: message,
		},
	}
	return g.chatModel.Stream(ctx, messages, writer)
}
