package components

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"
)

// IntentClassifier 意图分类器组件
// 这是 eino 框架中的工具组件，用于意图识别
type IntentClassifier struct {
	patterns map[string][]string
}

// NewIntentClassifier 创建意图分类器
func NewIntentClassifier() *IntentClassifier {
	return &IntentClassifier{
		patterns: map[string][]string{
			"greeting":  {"你好", "hello", "hi", "早上好", "下午好", "晚上好"},
			"question":  {"什么", "如何", "怎么", "为什么", "？", "?"},
			"complaint": {"投诉", "问题", "不满意", "差评", "退款"},
			"transfer":  {"人工", "客服", "转接", "人工客服"},
			"creative":  {"写", "创作", "诗", "故事", "文章"},
		},
	}
}

// Invoke 意图分类方法
func (ic *IntentClassifier) Invoke(ctx context.Context, input interface{}) (interface{}, error) {
	log.Printf("eino IntentClassifier 组件开始意图识别")

	// 获取用户输入
	userInput, ok := input.(string)
	if !ok {
		return nil, fmt.Errorf("输入格式错误，期望 string")
	}

	// 转换为小写进行匹配
	lowerInput := strings.ToLower(userInput)

	// 遍历模式进行匹配
	for intent, patterns := range ic.patterns {
		for _, pattern := range patterns {
			if strings.Contains(lowerInput, pattern) {
				log.Printf("eino IntentClassifier 识别意图: %s", intent)
				return map[string]interface{}{
					"intent":     intent,
					"confidence": 0.9,
					"input":      userInput,
				}, nil
			}
		}
	}

	// 默认意图
	log.Printf("eino IntentClassifier 识别为默认意图")
	return map[string]interface{}{
		"intent":     "default",
		"confidence": 0.5,
		"input":      userInput,
	}, nil
}

// ContextManager 上下文管理器组件
// 这是 eino 框架中的状态管理组件
type ContextManager struct {
	sessions map[string]*SessionContext
}

// SessionContext 会话上下文
type SessionContext struct {
	UserID    string                 `json:"user_id"`
	SessionID string                 `json:"session_id"`
	History   []ChatMessage          `json:"history"`
	Metadata  map[string]interface{} `json:"metadata"`
	CreatedAt time.Time              `json:"created_at"`
	UpdatedAt time.Time              `json:"updated_at"`
}

// NewContextManager 创建上下文管理器
func NewContextManager() *ContextManager {
	return &ContextManager{
		sessions: make(map[string]*SessionContext),
	}
}

// Invoke 上下文管理方法
func (cm *ContextManager) Invoke(ctx context.Context, input interface{}) (interface{}, error) {
	log.Printf("eino ContextManager 组件开始上下文管理")

	// 解析输入
	inputMap, ok := input.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("输入格式错误，期望 map[string]interface{}")
	}

	userID, _ := inputMap["user_id"].(string)
	sessionID, _ := inputMap["session_id"].(string)
	message, _ := inputMap["message"].(ChatMessage)

	// 获取或创建会话上下文
	sessionCtx, exists := cm.sessions[sessionID]
	if !exists {
		sessionCtx = &SessionContext{
			UserID:    userID,
			SessionID: sessionID,
			History:   make([]ChatMessage, 0),
			Metadata:  make(map[string]interface{}),
			CreatedAt: time.Now(),
		}
		cm.sessions[sessionID] = sessionCtx
	}

	// 更新上下文
	sessionCtx.History = append(sessionCtx.History, message)
	sessionCtx.UpdatedAt = time.Now()

	// 构建增强的上下文
	enhancedContext := map[string]interface{}{
		"session":    sessionCtx,
		"history":    sessionCtx.History,
		"user_id":    userID,
		"session_id": sessionID,
		"message":    message,
		"timestamp":  time.Now().Unix(),
	}

	log.Printf("eino ContextManager 上下文管理完成，历史消息数: %d", len(sessionCtx.History))
	return enhancedContext, nil
}

// KnowledgeRetriever 知识检索器组件
// 这是 eino 框架中的知识库组件
type KnowledgeRetriever struct {
	knowledge map[string][]string
}

// NewKnowledgeRetriever 创建知识检索器
func NewKnowledgeRetriever() *KnowledgeRetriever {
	return &KnowledgeRetriever{
		knowledge: map[string][]string{
			"产品信息": {
				"我们提供多种AI产品和服务",
				"包括聊天机器人、智能客服、内容生成等",
				"所有产品都基于先进的大语言模型技术",
			},
			"技术支持": {
				"我们提供7x24小时技术支持",
				"可以通过在线客服、邮件、电话联系我们",
				"技术支持团队由经验丰富的工程师组成",
			},
			"价格信息": {
				"我们提供灵活的定价方案",
				"包括免费试用、按量计费、包年套餐等",
				"具体价格请联系销售团队获取详细报价",
			},
			"常见问题": {
				"如何开始使用我们的产品？",
				"如何联系技术支持？",
				"如何升级或降级套餐？",
			},
		},
	}
}

// Invoke 知识检索方法
func (kr *KnowledgeRetriever) Invoke(ctx context.Context, input interface{}) (interface{}, error) {
	log.Printf("eino KnowledgeRetriever 组件开始知识检索")

	// 获取查询内容
	query, ok := input.(string)
	if !ok {
		return nil, fmt.Errorf("输入格式错误，期望 string")
	}

	// 简单的关键词匹配
	var results []string
	lowerQuery := strings.ToLower(query)

	for category, knowledge := range kr.knowledge {
		if strings.Contains(lowerQuery, strings.ToLower(category)) {
			results = append(results, knowledge...)
		}
	}

	// 如果没有找到匹配的知识，返回通用信息
	if len(results) == 0 {
		results = []string{
			"很抱歉，我没有找到相关的信息",
			"您可以尝试重新描述您的问题",
			"或者联系我们的技术支持团队获取帮助",
		}
	}

	log.Printf("eino KnowledgeRetriever 检索到 %d 条知识", len(results))
	return map[string]interface{}{
		"query":   query,
		"results": results,
		"count":   len(results),
	}, nil
}

// ResponseOptimizer 响应优化器组件
// 这是 eino 框架中的后处理组件
type ResponseOptimizer struct {
	templates map[string]string
}

// NewResponseOptimizer 创建响应优化器
func NewResponseOptimizer() *ResponseOptimizer {
	return &ResponseOptimizer{
		templates: map[string]string{
			"greeting":  "您好！我是智能客服助手，很高兴为您服务。",
			"question":  "根据您的问题，我为您提供以下信息：",
			"complaint": "非常抱歉给您带来不便，我会认真处理您的问题。",
			"transfer":  "我理解您需要人工客服的帮助，正在为您转接...",
			"default":   "感谢您的咨询，我会尽力为您提供帮助。",
		},
	}
}

// Invoke 响应优化方法
func (ro *ResponseOptimizer) Invoke(ctx context.Context, input interface{}) (interface{}, error) {
	log.Printf("eino ResponseOptimizer 组件开始响应优化")

	// 解析输入
	inputMap, ok := input.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("输入格式错误，期望 map[string]interface{}")
	}

	intent, _ := inputMap["intent"].(string)
	response, _ := inputMap["response"].(string)

	// 根据意图选择模板
	template, exists := ro.templates[intent]
	if !exists {
		template = ro.templates["default"]
	}

	// 优化响应
	optimizedResponse := fmt.Sprintf("%s\n\n%s", template, response)

	// 添加结尾
	optimizedResponse += "\n\n如果您还有其他问题，请随时告诉我！"

	log.Printf("eino ResponseOptimizer 响应优化完成")
	return map[string]interface{}{
		"original_response":  response,
		"optimized_response": optimizedResponse,
		"intent":             intent,
		"template":           template,
	}, nil
}
