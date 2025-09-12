package config

// Config 应用配置结构
type Config struct {
	Host string    `json:"host,omitempty"`
	Port int       `json:"port,omitempty"`
	Mode string    `json:"mode,omitempty"`
	LLM  LLMConfig `json:"llm,omitempty"`
	Log  LogConfig `json:"log,omitempty"`
}

// LLMConfig 大模型配置
type LLMConfig struct {
	GLM4 GLM4Config `json:"glm4,omitempty"`
}

// GLM4Config GLM4.5 模型配置
type GLM4Config struct {
	Name        string  `json:"name"`        // 模型名称
	URL         string  `json:"url"`         // API 地址
	AppKey      string  `json:"appKey"`      // 应用密钥
	Model       string  `json:"model"`       // 具体模型名称
	MaxTokens   int     `json:"maxTokens"`   // 最大token数
	Temperature float64 `json:"temperature"` // 温度参数
	TopP        float64 `json:"topP"`        // TopP参数
}

// LogConfig 日志配置
type LogConfig struct {
	ServiceName string `json:"serviceName"`
	Mode        string `json:"mode"`
	Level       string `json:"level"`
	Encoding    string `json:"encoding"`
}
