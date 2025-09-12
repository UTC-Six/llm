package chain

import (
	"context"
	"fmt"
	"io"
	"log"

	"github.com/UTC-Six/llm/internal/config"
	"github.com/UTC-Six/llm/internal/eino/components"
	"github.com/UTC-Six/llm/internal/types"
)

// ChainOrchestrator eino 框架的 Chain 编排器
// Chain 是最简单的编排方式，形成有向链式结构，数据单向流动
type ChainOrchestrator struct {
	config    config.Config
	chatModel *components.ChatModel
	optimizer *components.ResponseOptimizer
}

// NewChainOrchestrator 创建 Chain 编排器
func NewChainOrchestrator(c config.Config) *ChainOrchestrator {
	return &ChainOrchestrator{
		config:    c,
		chatModel: components.NewChatModel(c.LLM.GLM4),
		optimizer: components.NewResponseOptimizer(),
	}
}

// ChatInvoke Chain 同步聊天
// 使用 eino 框架的 Chain 编排：用户输入 -> 消息构建 -> 模型调用 -> 响应优化 -> 输出
func (co *ChainOrchestrator) ChatInvoke(ctx context.Context, req *types.ChatRequest) (*types.ChatResponse, error) {
	log.Printf("eino Chain 编排开始同步聊天: %s", req.Message)

	// Chain 步骤1: 构建消息
	messages := co.buildMessages(req.Message)
	log.Printf("eino Chain 步骤1 - 消息构建完成，消息数: %d", len(messages))

	// Chain 步骤2: 调用聊天模型
	modelResult, err := co.chatModel.Invoke(ctx, messages)
	if err != nil {
		return nil, fmt.Errorf("eino Chain 模型调用失败: %v", err)
	}
	log.Printf("eino Chain 步骤2 - 模型调用完成")

	// 解析模型结果
	chatResp, ok := modelResult.(*components.ChatResponse)
	if !ok {
		return nil, fmt.Errorf("eino Chain 模型结果格式错误")
	}

	// Chain 步骤3: 响应优化
	optimizerInput := map[string]interface{}{
		"intent":   "default",
		"response": chatResp.Choices[0].Message.Content,
	}

	optimizerResult, err := co.optimizer.Invoke(ctx, optimizerInput)
	if err != nil {
		return nil, fmt.Errorf("eino Chain 响应优化失败: %v", err)
	}
	log.Printf("eino Chain 步骤3 - 响应优化完成")

	// 解析优化结果
	optimizedData, ok := optimizerResult.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("eino Chain 优化结果格式错误")
	}

	optimizedResponse, _ := optimizedData["optimized_response"].(string)

	// 构建最终响应
	response := &types.ChatResponse{
		Code:    200,
		Message: "success",
		Data:    optimizedResponse,
		Model:   co.config.LLM.GLM4.Model,
		Usage: types.Usage{
			PromptTokens:     chatResp.Usage.PromptTokens,
			CompletionTokens: chatResp.Usage.CompletionTokens,
			TotalTokens:      chatResp.Usage.TotalTokens,
		},
	}

	log.Printf("eino Chain 编排同步聊天完成，响应长度: %d", len(optimizedResponse))
	return response, nil
}

// ChatStream Chain 流式聊天
// 使用 eino 框架的 Chain 编排进行流式响应
func (co *ChainOrchestrator) ChatStream(ctx context.Context, req *types.StreamChatRequest, writer io.Writer) error {
	log.Printf("eino Chain 编排开始流式聊天: %s", req.Message)

	// Chain 步骤1: 构建消息
	messages := co.buildMessages(req.Message)
	log.Printf("eino Chain 步骤1 - 消息构建完成，消息数: %d", len(messages))

	// Chain 步骤2: 流式调用聊天模型
	err := co.chatModel.Stream(ctx, messages, writer)
	if err != nil {
		return fmt.Errorf("eino Chain 流式模型调用失败: %v", err)
	}

	log.Printf("eino Chain 编排流式聊天完成")
	return nil
}

// buildMessages 构建消息链
// 这是 eino Chain 编排的核心：按顺序构建消息
func (co *ChainOrchestrator) buildMessages(userMessage string) []components.ChatMessage {
	// Chain 步骤1: 系统消息
	systemMessage := components.ChatMessage{
		Role:    "system",
		Content: "你是一个智能客服助手，请用友好、专业的方式回答用户的问题。",
	}

	// Chain 步骤2: 用户消息
	userMsg := components.ChatMessage{
		Role:    "user",
		Content: userMessage,
	}

	// Chain 步骤3: 构建消息链
	messages := []components.ChatMessage{
		systemMessage,
		userMsg,
	}

	return messages
}

// SimpleQAService 简单问答服务
// 展示 eino Chain 编排的简单应用场景
func (co *ChainOrchestrator) SimpleQAService(ctx context.Context, question string) (string, error) {
	log.Printf("eino Chain 简单问答服务开始: %s", question)

	// 构建问答消息
	messages := []components.ChatMessage{
		{
			Role:    "system",
			Content: "你是一个知识问答助手，请准确回答用户的问题。",
		},
		{
			Role:    "user",
			Content: question,
		},
	}

	// 调用模型
	result, err := co.chatModel.Invoke(ctx, messages)
	if err != nil {
		return "", fmt.Errorf("eino Chain 简单问答失败: %v", err)
	}

	// 解析结果
	chatResp, ok := result.(*components.ChatResponse)
	if !ok {
		return "", fmt.Errorf("eino Chain 问答结果格式错误")
	}

	answer := chatResp.Choices[0].Message.Content
	log.Printf("eino Chain 简单问答服务完成，答案长度: %d", len(answer))
	return answer, nil
}

// CreativeWritingService 创意写作服务
// 展示 eino Chain 编排的创意应用场景
func (co *ChainOrchestrator) CreativeWritingService(ctx context.Context, prompt string) (string, error) {
	log.Printf("eino Chain 创意写作服务开始: %s", prompt)

	// 构建创作消息
	messages := []components.ChatMessage{
		{
			Role:    "system",
			Content: "你是一个富有创意的写作助手，请根据用户的要求创作内容。",
		},
		{
			Role:    "user",
			Content: prompt,
		},
	}

	// 调用模型
	result, err := co.chatModel.Invoke(ctx, messages)
	if err != nil {
		return "", fmt.Errorf("eino Chain 创意写作失败: %v", err)
	}

	// 解析结果
	chatResp, ok := result.(*components.ChatResponse)
	if !ok {
		return "", fmt.Errorf("eino Chain 创作结果格式错误")
	}

	content := chatResp.Choices[0].Message.Content
	log.Printf("eino Chain 创意写作服务完成，内容长度: %d", len(content))
	return content, nil
}
