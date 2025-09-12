package svc

import (
	"github.com/UTC-Six/llm/internal/config"
	"github.com/UTC-Six/llm/internal/logic"
)

// ServiceContext 服务上下文
type ServiceContext struct {
	Config     config.Config
	LLMService *logic.LLMService
}

// NewServiceContext 创建服务上下文
func NewServiceContext(c config.Config) *ServiceContext {
	// 初始化大模型服务
	llmService := logic.NewLLMService(c)

	return &ServiceContext{
		Config:     c,
		LLMService: llmService,
	}
}
