package workflow

import (
	"context"
	"fmt"
	"io"
	"log"
	"strings"
	"time"

	"github.com/UTC-Six/llm/internal/config"
	"github.com/UTC-Six/llm/internal/eino/components"
	"github.com/UTC-Six/llm/internal/types"
)

// WorkflowOrchestrator eino 框架的 Workflow 编排器
// Workflow 是最高级的编排方式，支持复杂的数据流转和字段级映射
type WorkflowOrchestrator struct {
	config             config.Config
	chatModel          *components.ChatModel
	intentClassifier   *components.IntentClassifier
	contextManager     *components.ContextManager
	knowledgeRetriever *components.KnowledgeRetriever
	responseOptimizer  *components.ResponseOptimizer
}

// WorkflowData Workflow 数据流结构
// 这是 eino Workflow 编排的核心数据结构，支持字段级映射
type WorkflowData struct {
	// 输入数据
	UserInput string `json:"user_input"`
	UserID    string `json:"user_id"`
	SessionID string `json:"session_id"`
	Timestamp int64  `json:"timestamp"`

	// 意图分析数据
	Intent     string  `json:"intent"`
	Confidence float64 `json:"confidence"`

	// 上下文数据
	Context map[string]interface{}   `json:"context"`
	History []components.ChatMessage `json:"history"`

	// 知识检索数据
	Knowledge     map[string]interface{} `json:"knowledge"`
	KnowledgeText string                 `json:"knowledge_text"`

	// 消息构建数据
	Messages []components.ChatMessage `json:"messages"`

	// 模型响应数据
	ModelResponse *components.ChatResponse `json:"model_response"`
	RawResponse   string                   `json:"raw_response"`

	// 优化数据
	OptimizedResponse string `json:"optimized_response"`
	Template          string `json:"template"`

	// 元数据
	Metadata       map[string]interface{} `json:"metadata"`
	ProcessingTime int64                  `json:"processing_time"`
	Status         string                 `json:"status"`
}

// NewWorkflowOrchestrator 创建 Workflow 编排器
func NewWorkflowOrchestrator(c config.Config) *WorkflowOrchestrator {
	return &WorkflowOrchestrator{
		config:             c,
		chatModel:          components.NewChatModel(c.LLM.GLM4),
		intentClassifier:   components.NewIntentClassifier(),
		contextManager:     components.NewContextManager(),
		knowledgeRetriever: components.NewKnowledgeRetriever(),
		responseOptimizer:  components.NewResponseOptimizer(),
	}
}

// ChatInvoke Workflow 同步聊天
// 使用 eino 框架的 Workflow 编排：数据输入 -> 预处理 -> 意图分析 -> 上下文管理 -> 知识检索 -> 消息构建 -> 模型调用 -> 响应优化 -> 输出
func (wo *WorkflowOrchestrator) ChatInvoke(ctx context.Context, req *types.ChatRequest) (*types.ChatResponse, error) {
	log.Printf("eino Workflow 编排开始同步聊天: %s", req.Message)

	// 创建 Workflow 数据流
	data := &WorkflowData{
		UserInput: req.Message,
		UserID:    "user_123",
		SessionID: "session_456",
		Timestamp: time.Now().Unix(),
		Context:   make(map[string]interface{}),
		Knowledge: make(map[string]interface{}),
		Metadata:  make(map[string]interface{}),
		Status:    "started",
	}

	// Workflow 节点1: 数据预处理
	err := wo.preprocessData(ctx, data)
	if err != nil {
		return nil, fmt.Errorf("eino Workflow 数据预处理失败: %v", err)
	}
	log.Printf("eino Workflow 节点1 - 数据预处理完成")

	// Workflow 节点2: 意图分析
	err = wo.analyzeIntent(ctx, data)
	if err != nil {
		return nil, fmt.Errorf("eino Workflow 意图分析失败: %v", err)
	}
	log.Printf("eino Workflow 节点2 - 意图分析完成: %s", data.Intent)

	// Workflow 节点3: 上下文管理
	err = wo.manageContext(ctx, data)
	if err != nil {
		return nil, fmt.Errorf("eino Workflow 上下文管理失败: %v", err)
	}
	log.Printf("eino Workflow 节点3 - 上下文管理完成")

	// Workflow 节点4: 知识检索（根据意图决定是否执行）
	if data.Intent == "question" {
		err = wo.retrieveKnowledge(ctx, data)
		if err != nil {
			return nil, fmt.Errorf("eino Workflow 知识检索失败: %v", err)
		}
		log.Printf("eino Workflow 节点4 - 知识检索完成")
	}

	// Workflow 节点5: 消息构建
	err = wo.buildMessages(ctx, data)
	if err != nil {
		return nil, fmt.Errorf("eino Workflow 消息构建失败: %v", err)
	}
	log.Printf("eino Workflow 节点5 - 消息构建完成")

	// Workflow 节点6: 模型调用
	err = wo.callModel(ctx, data)
	if err != nil {
		return nil, fmt.Errorf("eino Workflow 模型调用失败: %v", err)
	}
	log.Printf("eino Workflow 节点6 - 模型调用完成")

	// Workflow 节点7: 响应优化
	err = wo.optimizeResponse(ctx, data)
	if err != nil {
		return nil, fmt.Errorf("eino Workflow 响应优化失败: %v", err)
	}
	log.Printf("eino Workflow 节点7 - 响应优化完成")

	// Workflow 节点8: 后处理
	err = wo.postprocess(ctx, data)
	if err != nil {
		return nil, fmt.Errorf("eino Workflow 后处理失败: %v", err)
	}
	log.Printf("eino Workflow 节点8 - 后处理完成")

	// 构建最终响应
	response := &types.ChatResponse{
		Code:    200,
		Message: "success",
		Data:    data.OptimizedResponse,
		Model:   wo.config.LLM.GLM4.Model,
		Usage: types.Usage{
			PromptTokens:     data.ModelResponse.Usage.PromptTokens,
			CompletionTokens: data.ModelResponse.Usage.CompletionTokens,
			TotalTokens:      data.ModelResponse.Usage.TotalTokens,
		},
	}

	log.Printf("eino Workflow 编排同步聊天完成，响应长度: %d", len(data.OptimizedResponse))
	return response, nil
}

// ChatStream Workflow 流式聊天
// 使用 eino 框架的 Workflow 编排进行流式响应
func (wo *WorkflowOrchestrator) ChatStream(ctx context.Context, req *types.StreamChatRequest, writer io.Writer) error {
	log.Printf("eino Workflow 编排开始流式聊天: %s", req.Message)

	// 创建 Workflow 数据流
	data := &WorkflowData{
		UserInput: req.Message,
		UserID:    "user_123",
		SessionID: "session_456",
		Timestamp: time.Now().Unix(),
		Context:   make(map[string]interface{}),
		Knowledge: make(map[string]interface{}),
		Metadata:  make(map[string]interface{}),
		Status:    "started",
	}

	// Workflow 节点1-5: 预处理、意图分析、上下文管理、知识检索、消息构建
	err := wo.preprocessData(ctx, data)
	if err != nil {
		return fmt.Errorf("eino Workflow 数据预处理失败: %v", err)
	}

	err = wo.analyzeIntent(ctx, data)
	if err != nil {
		return fmt.Errorf("eino Workflow 意图分析失败: %v", err)
	}

	err = wo.manageContext(ctx, data)
	if err != nil {
		return fmt.Errorf("eino Workflow 上下文管理失败: %v", err)
	}

	if data.Intent == "question" {
		err = wo.retrieveKnowledge(ctx, data)
		if err != nil {
			return fmt.Errorf("eino Workflow 知识检索失败: %v", err)
		}
	}

	err = wo.buildMessages(ctx, data)
	if err != nil {
		return fmt.Errorf("eino Workflow 消息构建失败: %v", err)
	}

	// Workflow 节点6: 流式模型调用
	err = wo.callModelStream(ctx, data, writer)
	if err != nil {
		return fmt.Errorf("eino Workflow 流式模型调用失败: %v", err)
	}

	log.Printf("eino Workflow 编排流式聊天完成")
	return nil
}

// preprocessData 数据预处理节点
func (wo *WorkflowOrchestrator) preprocessData(ctx context.Context, data *WorkflowData) error {
	log.Printf("eino Workflow 节点1 - 开始数据预处理")

	// 清理用户输入
	data.UserInput = strings.TrimSpace(data.UserInput)

	// 添加预处理元数据
	data.Metadata["preprocessing_time"] = time.Now().Unix()
	data.Metadata["input_length"] = len(data.UserInput)
	data.Metadata["language"] = "zh-CN"

	// 更新状态
	data.Status = "preprocessed"

	log.Printf("eino Workflow 节点1 - 数据预处理完成，输入长度: %d", len(data.UserInput))
	return nil
}

// analyzeIntent 意图分析节点
func (wo *WorkflowOrchestrator) analyzeIntent(ctx context.Context, data *WorkflowData) error {
	log.Printf("eino Workflow 节点2 - 开始意图分析")

	// 调用意图分类器
	result, err := wo.intentClassifier.Invoke(ctx, data.UserInput)
	if err != nil {
		return err
	}

	// 解析结果
	intentData, ok := result.(map[string]interface{})
	if !ok {
		return fmt.Errorf("意图分析结果格式错误")
	}

	// 字段级映射
	data.Intent = intentData["intent"].(string)
	data.Confidence = intentData["confidence"].(float64)

	// 添加意图分析元数据
	data.Metadata["intent_analysis_time"] = time.Now().Unix()
	data.Metadata["intent_confidence"] = data.Confidence

	// 更新状态
	data.Status = "intent_analyzed"

	log.Printf("eino Workflow 节点2 - 意图分析完成: %s (置信度: %.2f)", data.Intent, data.Confidence)
	return nil
}

// manageContext 上下文管理节点
func (wo *WorkflowOrchestrator) manageContext(ctx context.Context, data *WorkflowData) error {
	log.Printf("eino Workflow 节点3 - 开始上下文管理")

	// 构建上下文管理输入
	contextInput := map[string]interface{}{
		"user_id":    data.UserID,
		"session_id": data.SessionID,
		"message": components.ChatMessage{
			Role:    "user",
			Content: data.UserInput,
		},
	}

	// 调用上下文管理器
	result, err := wo.contextManager.Invoke(ctx, contextInput)
	if err != nil {
		return err
	}

	// 解析结果
	contextData, ok := result.(map[string]interface{})
	if !ok {
		return fmt.Errorf("上下文管理结果格式错误")
	}

	// 字段级映射
	data.Context = contextData
	if session, ok := contextData["session"].(*components.SessionContext); ok {
		data.History = session.History
	}

	// 添加上下文管理元数据
	data.Metadata["context_management_time"] = time.Now().Unix()
	data.Metadata["history_length"] = len(data.History)

	// 更新状态
	data.Status = "context_managed"

	log.Printf("eino Workflow 节点3 - 上下文管理完成，历史消息数: %d", len(data.History))
	return nil
}

// retrieveKnowledge 知识检索节点
func (wo *WorkflowOrchestrator) retrieveKnowledge(ctx context.Context, data *WorkflowData) error {
	log.Printf("eino Workflow 节点4 - 开始知识检索")

	// 调用知识检索器
	result, err := wo.knowledgeRetriever.Invoke(ctx, data.UserInput)
	if err != nil {
		return err
	}

	// 解析结果
	knowledgeData, ok := result.(map[string]interface{})
	if !ok {
		return fmt.Errorf("知识检索结果格式错误")
	}

	// 字段级映射
	data.Knowledge = knowledgeData
	if results, ok := knowledgeData["results"].([]string); ok {
		data.KnowledgeText = ""
		for _, result := range results {
			data.KnowledgeText += result + "\n"
		}
	}

	// 添加知识检索元数据
	data.Metadata["knowledge_retrieval_time"] = time.Now().Unix()
	data.Metadata["knowledge_count"] = knowledgeData["count"]

	// 更新状态
	data.Status = "knowledge_retrieved"

	log.Printf("eino Workflow 节点4 - 知识检索完成，知识条数: %v", knowledgeData["count"])
	return nil
}

// buildMessages 消息构建节点
func (wo *WorkflowOrchestrator) buildMessages(ctx context.Context, data *WorkflowData) error {
	log.Printf("eino Workflow 节点5 - 开始消息构建")

	// 根据意图和上下文构建系统消息
	var systemContent string
	switch data.Intent {
	case "greeting":
		systemContent = "你是一个友好的智能客服助手，请用温暖、热情的方式回应用户的问候。"
	case "question":
		if data.KnowledgeText != "" {
			systemContent = fmt.Sprintf("你是一个专业的智能客服助手，请基于以下知识库信息准确回答用户的问题：\n%s", data.KnowledgeText)
		} else {
			systemContent = "你是一个专业的智能客服助手，请准确回答用户的问题。"
		}
	case "complaint":
		systemContent = "你是一个专业的客服助手，请用诚恳、专业的方式处理用户的投诉，表达歉意并提供解决方案。"
	case "transfer":
		systemContent = "你是一个智能客服助手，请礼貌地告知用户正在转接人工客服。"
	case "creative":
		systemContent = "你是一个富有创意的写作助手，请根据用户的要求创作生动、有趣的内容。"
	default:
		systemContent = "你是一个智能客服助手，请根据用户的需求提供合适的帮助。"
	}

	// 构建消息数组
	data.Messages = []components.ChatMessage{
		{
			Role:    "system",
			Content: systemContent,
		},
		{
			Role:    "user",
			Content: data.UserInput,
		},
	}

	// 添加历史消息（如果有）
	if len(data.History) > 0 {
		// 只添加最近的几条历史消息
		recentHistory := data.History
		if len(recentHistory) > 5 {
			recentHistory = recentHistory[len(recentHistory)-5:]
		}
		data.Messages = append(recentHistory, data.Messages...)
	}

	// 添加消息构建元数据
	data.Metadata["message_building_time"] = time.Now().Unix()
	data.Metadata["message_count"] = len(data.Messages)

	// 更新状态
	data.Status = "messages_built"

	log.Printf("eino Workflow 节点5 - 消息构建完成，消息数: %d", len(data.Messages))
	return nil
}

// callModel 模型调用节点
func (wo *WorkflowOrchestrator) callModel(ctx context.Context, data *WorkflowData) error {
	log.Printf("eino Workflow 节点6 - 开始模型调用")

	// 调用聊天模型
	result, err := wo.chatModel.Invoke(ctx, data.Messages)
	if err != nil {
		return err
	}

	// 解析结果
	chatResp, ok := result.(*components.ChatResponse)
	if !ok {
		return fmt.Errorf("模型调用结果格式错误")
	}

	// 字段级映射
	data.ModelResponse = chatResp
	data.RawResponse = chatResp.Choices[0].Message.Content

	// 添加模型调用元数据
	data.Metadata["model_call_time"] = time.Now().Unix()
	data.Metadata["model_id"] = chatResp.ID
	data.Metadata["prompt_tokens"] = chatResp.Usage.PromptTokens
	data.Metadata["completion_tokens"] = chatResp.Usage.CompletionTokens
	data.Metadata["total_tokens"] = chatResp.Usage.TotalTokens

	// 更新状态
	data.Status = "model_called"

	log.Printf("eino Workflow 节点6 - 模型调用完成，响应长度: %d", len(data.RawResponse))
	return nil
}

// callModelStream 流式模型调用节点
func (wo *WorkflowOrchestrator) callModelStream(ctx context.Context, data *WorkflowData, writer io.Writer) error {
	log.Printf("eino Workflow 节点6 - 开始流式模型调用")

	// 调用流式聊天模型
	err := wo.chatModel.Stream(ctx, data.Messages, writer)
	if err != nil {
		return err
	}

	// 添加流式模型调用元数据
	data.Metadata["stream_model_call_time"] = time.Now().Unix()

	// 更新状态
	data.Status = "stream_model_called"

	log.Printf("eino Workflow 节点6 - 流式模型调用完成")
	return nil
}

// optimizeResponse 响应优化节点
func (wo *WorkflowOrchestrator) optimizeResponse(ctx context.Context, data *WorkflowData) error {
	log.Printf("eino Workflow 节点7 - 开始响应优化")

	// 构建优化器输入
	optimizerInput := map[string]interface{}{
		"intent":   data.Intent,
		"response": data.RawResponse,
	}

	// 调用响应优化器
	result, err := wo.responseOptimizer.Invoke(ctx, optimizerInput)
	if err != nil {
		return err
	}

	// 解析结果
	optimizedData, ok := result.(map[string]interface{})
	if !ok {
		return fmt.Errorf("响应优化结果格式错误")
	}

	// 字段级映射
	data.OptimizedResponse = optimizedData["optimized_response"].(string)
	data.Template = optimizedData["template"].(string)

	// 添加响应优化元数据
	data.Metadata["response_optimization_time"] = time.Now().Unix()
	data.Metadata["optimization_template"] = data.Template

	// 更新状态
	data.Status = "response_optimized"

	log.Printf("eino Workflow 节点7 - 响应优化完成，优化后长度: %d", len(data.OptimizedResponse))
	return nil
}

// postprocess 后处理节点
func (wo *WorkflowOrchestrator) postprocess(ctx context.Context, data *WorkflowData) error {
	log.Printf("eino Workflow 节点8 - 开始后处理")

	// 计算总处理时间
	data.ProcessingTime = time.Now().Unix() - data.Timestamp

	// 添加后处理元数据
	data.Metadata["postprocessing_time"] = time.Now().Unix()
	data.Metadata["total_processing_time"] = data.ProcessingTime
	data.Metadata["workflow_status"] = "completed"

	// 更新状态
	data.Status = "completed"

	log.Printf("eino Workflow 节点8 - 后处理完成，总处理时间: %d 秒", data.ProcessingTime)
	return nil
}
