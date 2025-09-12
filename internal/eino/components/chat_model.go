package components

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
)

// ChatModel eino 框架的聊天模型组件
// 这是 eino 框架中的核心组件，用于与大模型进行交互
type ChatModel struct {
	config     config.GLM4Config
	httpClient *http.Client
}

// NewChatModel 创建 eino 聊天模型组件
func NewChatModel(cfg config.GLM4Config) *ChatModel {
	return &ChatModel{
		config: cfg,
		httpClient: &http.Client{
			Timeout: 60 * time.Second,
		},
	}
}

// ChatMessage 聊天消息结构
type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ChatRequest 聊天请求结构
type ChatRequest struct {
	Model       string        `json:"model"`
	Messages    []ChatMessage `json:"messages"`
	MaxTokens   int           `json:"max_tokens"`
	Temperature float64       `json:"temperature"`
	TopP        float64       `json:"top_p"`
	Stream      bool          `json:"stream"`
}

// ChatResponse 聊天响应结构
type ChatResponse struct {
	ID      string   `json:"id"`
	Object  string   `json:"object"`
	Created int64    `json:"created"`
	Model   string   `json:"model"`
	Choices []Choice `json:"choices"`
	Usage   Usage    `json:"usage"`
}

// Choice 选择结构
type Choice struct {
	Index        int         `json:"index"`
	Message      ChatMessage `json:"message"`
	FinishReason string      `json:"finish_reason"`
}

// Usage 使用情况
type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

// StreamResponse 流式响应结构
type StreamResponse struct {
	ID      string         `json:"id"`
	Object  string         `json:"object"`
	Created int64          `json:"created"`
	Model   string         `json:"model"`
	Choices []StreamChoice `json:"choices"`
	Usage   *Usage         `json:"usage,omitempty"`
}

// StreamChoice 流式选择结构
type StreamChoice struct {
	Index        int     `json:"index"`
	Delta        Delta   `json:"delta"`
	FinishReason *string `json:"finish_reason"`
}

// Delta 增量结构
type Delta struct {
	Role    *string `json:"role,omitempty"`
	Content *string `json:"content,omitempty"`
}

// Invoke 同步调用方法
// 这是 eino 框架中组件的标准接口方法
func (cm *ChatModel) Invoke(ctx context.Context, input interface{}) (interface{}, error) {
	log.Printf("eino ChatModel 组件同步调用开始")

	// 将输入转换为消息格式
	messages, ok := input.([]ChatMessage)
	if !ok {
		return nil, fmt.Errorf("输入格式错误，期望 []ChatMessage")
	}

	// 构建请求
	request := ChatRequest{
		Model:       cm.config.Model,
		Messages:    messages,
		MaxTokens:   cm.config.MaxTokens,
		Temperature: cm.config.Temperature,
		TopP:        cm.config.TopP,
		Stream:      false,
	}

	// 发送请求
	response, err := cm.callAPI(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("eino ChatModel 调用失败: %v", err)
	}

	log.Printf("eino ChatModel 组件同步调用完成")
	return response, nil
}

// Stream 流式调用方法
// 这是 eino 框架中组件的标准流式接口方法
func (cm *ChatModel) Stream(ctx context.Context, input interface{}, writer io.Writer) error {
	log.Printf("eino ChatModel 组件流式调用开始")

	// 将输入转换为消息格式
	messages, ok := input.([]ChatMessage)
	if !ok {
		return fmt.Errorf("输入格式错误，期望 []ChatMessage")
	}

	// 构建请求
	request := ChatRequest{
		Model:       cm.config.Model,
		Messages:    messages,
		MaxTokens:   cm.config.MaxTokens,
		Temperature: cm.config.Temperature,
		TopP:        cm.config.TopP,
		Stream:      true,
	}

	// 发送流式请求
	err := cm.callStreamAPI(ctx, request, writer)
	if err != nil {
		return fmt.Errorf("eino ChatModel 流式调用失败: %v", err)
	}

	log.Printf("eino ChatModel 组件流式调用完成")
	return nil
}

// callAPI 调用 API
func (cm *ChatModel) callAPI(ctx context.Context, request ChatRequest) (*ChatResponse, error) {
	// 序列化请求
	requestBody, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("序列化请求失败: %v", err)
	}

	// 创建 HTTP 请求
	req, err := http.NewRequestWithContext(ctx, "POST", cm.config.URL+"chat/completions", bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, fmt.Errorf("创建 HTTP 请求失败: %v", err)
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+cm.config.AppKey)

	// 发送请求
	resp, err := cm.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("发送 HTTP 请求失败: %v", err)
	}
	defer resp.Body.Close()

	// 检查响应状态
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API 返回错误状态 %d: %s", resp.StatusCode, string(body))
	}

	// 解析响应
	var chatResp ChatResponse
	if err := json.NewDecoder(resp.Body).Decode(&chatResp); err != nil {
		return nil, fmt.Errorf("解析响应失败: %v", err)
	}

	// 检查是否有错误
	if len(chatResp.Choices) == 0 {
		return nil, fmt.Errorf("API 返回空响应")
	}

	return &chatResp, nil
}

// callStreamAPI 调用流式 API
func (cm *ChatModel) callStreamAPI(ctx context.Context, request ChatRequest, writer io.Writer) error {
	// 序列化请求
	requestBody, err := json.Marshal(request)
	if err != nil {
		return fmt.Errorf("序列化请求失败: %v", err)
	}

	// 创建 HTTP 请求
	req, err := http.NewRequestWithContext(ctx, "POST", cm.config.URL+"chat/completions", bytes.NewBuffer(requestBody))
	if err != nil {
		return fmt.Errorf("创建 HTTP 请求失败: %v", err)
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+cm.config.AppKey)
	req.Header.Set("Accept", "text/event-stream")
	req.Header.Set("Cache-Control", "no-cache")

	// 发送请求
	client := &http.Client{Timeout: 0} // 流式请求不设置超时
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("发送 HTTP 请求失败: %v", err)
	}
	defer resp.Body.Close()

	// 检查响应状态
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API 返回错误状态 %d: %s", resp.StatusCode, string(body))
	}

	// 读取流式响应
	buffer := make([]byte, 1024)
	for {
		// 检查上下文是否被取消
		select {
		case <-ctx.Done():
			log.Printf("eino ChatModel 流式请求被取消")
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

			// 处理 SSE 格式的数据
			if bytes.HasPrefix(line, []byte("data: ")) {
				content := string(line[6:]) // 去掉 "data: " 前缀

				// 检查是否是结束标记
				if content == "[DONE]" {
					log.Printf("eino ChatModel 流式响应完成")
					return nil
				}

				// 解析 JSON 响应
				var streamResp StreamResponse
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

					log.Printf("eino ChatModel 发送流式数据块: %s", chunk)
				}
			}
		}
	}

	return nil
}
