package types

// ChatRequest 聊天请求结构
type ChatRequest struct {
	Message string `json:"message"`         // 用户输入的消息
	Model   string `json:"model,omitempty"` // 可选指定模型，默认使用配置中的模型
}

// ChatResponse 聊天响应结构
type ChatResponse struct {
	Code    int    `json:"code"`    // 响应码
	Message string `json:"message"` // 响应消息
	Data    string `json:"data"`    // 大模型返回的内容
	Model   string `json:"model"`   // 使用的模型名称
	Usage   Usage  `json:"usage"`   // Token使用情况
}

// StreamChatRequest 流式聊天请求结构
type StreamChatRequest struct {
	Message string `json:"message"`         // 用户输入的消息
	Model   string `json:"model,omitempty"` // 可选指定模型
}

// StreamChatResponse 流式聊天响应结构
type StreamChatResponse struct {
	Code    int    `json:"code"`            // 响应码
	Message string `json:"message"`         // 响应消息
	Data    string `json:"data"`            // 流式返回的内容片段
	Model   string `json:"model"`           // 使用的模型名称
	Done    bool   `json:"done"`            // 是否完成
	Usage   *Usage `json:"usage,omitempty"` // Token使用情况（仅在完成时返回）
}

// Usage Token使用情况
type Usage struct {
	PromptTokens     int `json:"promptTokens"`     // 输入token数
	CompletionTokens int `json:"completionTokens"` // 输出token数
	TotalTokens      int `json:"totalTokens"`      // 总token数
}

// ErrorResponse 错误响应结构
type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
