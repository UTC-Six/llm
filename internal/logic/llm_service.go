package logic

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/UTC-Six/llm/internal/config"
	"github.com/UTC-Six/llm/internal/types"
)

// LLMService 大模型服务
type LLMService struct {
	config config.Config
}

// NewLLMService 创建大模型服务实例
func NewLLMService(c config.Config) *LLMService {
	return &LLMService{
		config: c,
	}
}

// ChatInvoke 同步调用大模型进行聊天
// 这是传统的请求-响应模式，适合需要完整响应的场景
func (s *LLMService) ChatInvoke(ctx context.Context, req *types.ChatRequest) (*types.ChatResponse, error) {
	log.Printf("开始处理聊天请求: %s", req.Message)

	// 获取模型配置
	modelConfig := s.config.LLM.GLM4
	if req.Model != "" {
		// 如果请求中指定了模型，使用指定的模型
		modelConfig.Model = req.Model
	}

	// 构建请求消息
	messages := []map[string]interface{}{
		{
			"role":    "user",
			"content": req.Message,
		},
	}

	// 构建请求参数
	requestData := map[string]interface{}{
		"model":       modelConfig.Model,
		"messages":    messages,
		"max_tokens":  modelConfig.MaxTokens,
		"temperature": modelConfig.Temperature,
		"top_p":       modelConfig.TopP,
		"stream":      false, // 非流式请求
	}

	log.Printf("发送请求到模型 %s，参数: %+v", modelConfig.Model, requestData)

	// 调用真实的GLM API
	response, err := s.callGLMAPI(ctx, requestData)
	if err != nil {
		log.Printf("调用GLM API失败: %v", err)
		return nil, fmt.Errorf("调用大模型API失败: %v", err)
	}

	// 构建响应
	resp := &types.ChatResponse{
		Code:    200,
		Message: "success",
		Data:    response.Choices[0].Message.Content,
		Model:   modelConfig.Model,
		Usage: types.Usage{
			PromptTokens:     response.Usage.PromptTokens,
			CompletionTokens: response.Usage.CompletionTokens,
			TotalTokens:      response.Usage.TotalTokens,
		},
	}

	log.Printf("聊天请求处理完成，返回内容长度: %d", len(response.Choices[0].Message.Content))
	return resp, nil
}

// ChatStream 流式调用大模型进行聊天
// 这是流式响应模式，适合需要实时显示生成内容的场景
func (s *LLMService) ChatStream(ctx context.Context, req *types.StreamChatRequest, writer io.Writer) error {
	log.Printf("开始处理流式聊天请求: %s", req.Message)

	// 获取模型配置
	modelConfig := s.config.LLM.GLM4
	if req.Model != "" {
		modelConfig.Model = req.Model
	}

	// 构建请求消息
	messages := []map[string]interface{}{
		{
			"role":    "user",
			"content": req.Message,
		},
	}

	// 构建请求参数
	requestData := map[string]interface{}{
		"model":       modelConfig.Model,
		"messages":    messages,
		"max_tokens":  modelConfig.MaxTokens,
		"temperature": modelConfig.Temperature,
		"top_p":       modelConfig.TopP,
		"stream":      true, // 流式请求
	}

	log.Printf("发送流式请求到模型 %s，参数: %+v", modelConfig.Model, requestData)

	// 调用真实的GLM流式API
	return s.callGLMStreamAPI(ctx, requestData, writer)
}

// callGLMAPI 调用真实的GLM API
func (s *LLMService) callGLMAPI(ctx context.Context, requestData map[string]interface{}) (*GLMResponse, error) {
	modelConfig := s.config.LLM.GLM4

	// 构建GLM API请求
	glmRequest := GLMRequest{
		Model:       requestData["model"].(string),
		Messages:    requestData["messages"].([]map[string]interface{}),
		MaxTokens:   requestData["max_tokens"].(int),
		Temperature: requestData["temperature"].(float64),
		TopP:        requestData["top_p"].(float64),
		Stream:      false,
	}

	// 序列化请求
	requestBody, err := json.Marshal(glmRequest)
	if err != nil {
		return nil, fmt.Errorf("序列化请求失败: %v", err)
	}

	// 创建HTTP请求
	req, err := http.NewRequestWithContext(ctx, "POST", modelConfig.URL+"chat/completions", bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, fmt.Errorf("创建HTTP请求失败: %v", err)
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+modelConfig.AppKey)

	// 发送请求
	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("发送HTTP请求失败: %v", err)
	}
	defer resp.Body.Close()

	// 检查响应状态
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("GLM API返回错误状态 %d: %s", resp.StatusCode, string(body))
	}

	// 解析响应
	var glmResp GLMResponse
	if err := json.NewDecoder(resp.Body).Decode(&glmResp); err != nil {
		return nil, fmt.Errorf("解析GLM响应失败: %v", err)
	}

	// 检查是否有错误
	if len(glmResp.Choices) == 0 {
		return nil, fmt.Errorf("GLM API返回空响应")
	}

	return &glmResp, nil
}

// callGLMStreamAPI 调用真实的GLM流式API
func (s *LLMService) callGLMStreamAPI(ctx context.Context, requestData map[string]interface{}, writer io.Writer) error {
	modelConfig := s.config.LLM.GLM4

	// 构建GLM API请求
	glmRequest := GLMRequest{
		Model:       requestData["model"].(string),
		Messages:    requestData["messages"].([]map[string]interface{}),
		MaxTokens:   requestData["max_tokens"].(int),
		Temperature: requestData["temperature"].(float64),
		TopP:        requestData["top_p"].(float64),
		Stream:      true,
	}

	// 序列化请求
	requestBody, err := json.Marshal(glmRequest)
	if err != nil {
		return fmt.Errorf("序列化请求失败: %v", err)
	}

	// 创建HTTP请求
	req, err := http.NewRequestWithContext(ctx, "POST", modelConfig.URL+"chat/completions", bytes.NewBuffer(requestBody))
	if err != nil {
		return fmt.Errorf("创建HTTP请求失败: %v", err)
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+modelConfig.AppKey)
	req.Header.Set("Accept", "text/event-stream")
	req.Header.Set("Cache-Control", "no-cache")

	// 发送请求
	client := &http.Client{Timeout: 0} // 流式请求不设置超时
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("发送HTTP请求失败: %v", err)
	}
	defer resp.Body.Close()

	// 检查响应状态
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("GLM API返回错误状态 %d: %s", resp.StatusCode, string(body))
	}

	// 读取流式响应
	buffer := make([]byte, 1024)
	for {
		// 检查上下文是否被取消
		select {
		case <-ctx.Done():
			log.Printf("流式请求被取消")
			return ctx.Err()
		default:
		}

		// 读取数据
		n, err := resp.Body.Read(buffer)
		if err != nil {
			if err == io.EOF {
				break
			}
			return fmt.Errorf("读取流式响应失败: %v", err)
		}

		// 处理数据
		data := string(buffer[:n])
		lines := bytes.Split([]byte(data), []byte("\n"))

		for _, line := range lines {
			lineStr := string(line)
			if lineStr == "" || lineStr == "\r" {
				continue
			}

			// 处理SSE格式的数据
			if bytes.HasPrefix(line, []byte("data: ")) {
				content := string(line[6:]) // 去掉 "data: " 前缀

				// 检查是否是结束标记
				if content == "[DONE]" {
					log.Printf("流式响应完成")
					return nil
				}

				// 解析JSON响应
				var streamResp GLMStreamResponse
				if err := json.Unmarshal([]byte(content), &streamResp); err != nil {
					log.Printf("解析流式响应失败: %v", err)
					continue
				}

				// 提取内容
				if len(streamResp.Choices) > 0 && streamResp.Choices[0].Delta.Content != nil && *streamResp.Choices[0].Delta.Content != "" {
					chunk := *streamResp.Choices[0].Delta.Content

					// 写入响应
					_, err := writer.Write([]byte(fmt.Sprintf("data: %s\n\n", chunk)))
					if err != nil {
						log.Printf("写入流式响应失败: %v", err)
						return err
					}

					log.Printf("发送流式数据块: %s", chunk)
				}
			}
		}
	}

	log.Printf("流式聊天请求处理完成")
	return nil
}

// GLMRequest GLM API请求结构
type GLMRequest struct {
	Model       string                   `json:"model"`
	Messages    []map[string]interface{} `json:"messages"`
	MaxTokens   int                      `json:"max_tokens"`
	Temperature float64                  `json:"temperature"`
	TopP        float64                  `json:"top_p"`
	Stream      bool                     `json:"stream"`
}

// GLMResponse GLM API响应结构
type GLMResponse struct {
	ID      string   `json:"id"`
	Object  string   `json:"object"`
	Created int64    `json:"created"`
	Model   string   `json:"model"`
	Choices []Choice `json:"choices"`
	Usage   GLMUsage `json:"usage"`
}

// Choice GLM API选择结构
type Choice struct {
	Index        int     `json:"index"`
	Message      Message `json:"message"`
	FinishReason string  `json:"finish_reason"`
}

// Message GLM API消息结构
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// GLMUsage GLM API使用情况
type GLMUsage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

// GLMStreamResponse GLM流式响应结构
type GLMStreamResponse struct {
	ID      string         `json:"id"`
	Object  string         `json:"object"`
	Created int64          `json:"created"`
	Model   string         `json:"model"`
	Choices []StreamChoice `json:"choices"`
	Usage   *GLMUsage      `json:"usage,omitempty"`
}

// StreamChoice GLM流式选择结构
type StreamChoice struct {
	Index        int     `json:"index"`
	Delta        Delta   `json:"delta"`
	FinishReason *string `json:"finish_reason"`
}

// Delta GLM流式增量结构
type Delta struct {
	Role    *string `json:"role,omitempty"`
	Content *string `json:"content,omitempty"`
}
