# Eino 框架使用指南

本文档详细介绍了如何在 LLM Demo 项目中使用 Eino 框架的三种编排方式：Chain、Graph 和 Workflow。

## 项目架构

```
internal/eino/
├── components/          # Eino 组件
│   └── glm_model.go    # GLM 模型组件
├── chain/              # Chain 编排
│   └── chain_service.go
├── graph/              # Graph 编排
│   └── graph_service.go
└── workflow/           # Workflow 编排
    └── workflow_service.go
```

## 1. Chain 编排（简单链式流程）

Chain 是最简单的编排方式，适用于线性流程，按顺序执行各个步骤。

### 特点
- **线性执行**：步骤按顺序执行，前一步的输出作为后一步的输入
- **简单易用**：适合简单的对话流程
- **低延迟**：执行效率高

### 实现原理

```go
// Chain 编排流程
用户输入 -> 预处理 -> 消息构建 -> GLM模型 -> 后处理 -> 响应
```

### 核心代码

```go
// 构建消息链
func (s *ChainService) buildMessageChain(userMessage string) []map[string]interface{} {
    // Chain步骤1: 系统提示
    systemMessage := map[string]interface{}{
        "role":    "system",
        "content": "你是一个智能助手，请用友好、专业的方式回答用户的问题。",
    }
    
    // Chain步骤2: 用户消息
    userMsg := map[string]interface{}{
        "role":    "user", 
        "content": userMessage,
    }
    
    // Chain步骤3: 构建消息链
    messages := []map[string]interface{}{
        systemMessage,
        userMsg,
    }
    
    return messages
}
```

### API 接口

- **同步调用**: `POST /api/v1/eino/chain/invoke`
- **流式调用**: `POST /api/v1/eino/chain/stream`

### 使用示例

```bash
# 同步调用
curl -X POST http://localhost:8888/api/v1/eino/chain/invoke \
  -H "Content-Type: application/json" \
  -d '{"message": "你好，请介绍一下自己"}'

# 流式调用
curl -X POST http://localhost:8888/api/v1/eino/chain/stream \
  -H "Content-Type: application/json" \
  -d '{"message": "请写一首关于春天的短诗"}'
```

## 2. Graph 编排（复杂有向图流程）

Graph 是复杂的有向图编排，支持分支、循环和条件判断，适合复杂的业务逻辑。

### 特点
- **分支处理**：根据条件选择不同的处理路径
- **意图识别**：自动识别用户意图并选择相应的处理策略
- **灵活性强**：支持复杂的业务逻辑

### 实现原理

```go
// Graph 编排流程
用户输入 -> 意图识别 -> 分支选择 -> GLM模型 -> 响应
                ↓
            [问候/问题/创作/默认]
```

### 核心代码

```go
// 意图识别节点
func (s *GraphService) identifyIntent(ctx context.Context, message string) (string, error) {
    message = strings.ToLower(message)
    
    if strings.Contains(message, "你好") || strings.Contains(message, "hello") {
        return "greeting", nil
    }
    
    if strings.Contains(message, "？") || strings.Contains(message, "?") {
        return "question", nil
    }
    
    if strings.Contains(message, "写") || strings.Contains(message, "创作") {
        return "creative", nil
    }
    
    return "default", nil
}

// 分支处理
switch intent {
case "greeting":
    response, err = s.handleGreeting(ctx, req.Message)
case "question":
    response, err = s.handleQuestion(ctx, req.Message)
case "creative":
    response, err = s.handleCreative(ctx, req.Message)
default:
    response, err = s.handleDefault(ctx, req.Message)
}
```

### API 接口

- **同步调用**: `POST /api/v1/eino/graph/invoke`
- **流式调用**: `POST /api/v1/eino/graph/stream`

### 使用示例

```bash
# 问候意图
curl -X POST http://localhost:8888/api/v1/eino/graph/invoke \
  -H "Content-Type: application/json" \
  -d '{"message": "你好，很高兴见到你"}'

# 问题意图
curl -X POST http://localhost:8888/api/v1/eino/graph/invoke \
  -H "Content-Type: application/json" \
  -d '{"message": "什么是人工智能？"}'

# 创作意图
curl -X POST http://localhost:8888/api/v1/eino/graph/invoke \
  -H "Content-Type: application/json" \
  -d '{"message": "请写一首关于春天的诗"}'
```

## 3. Workflow 编排（高级编排流程）

Workflow 是最高级的编排方式，支持复杂的数据流转和字段级映射，适合企业级应用。

### 特点
- **数据流转**：支持复杂的数据结构在节点间传递
- **字段映射**：支持字段级别的数据映射和转换
- **元数据管理**：完整的元数据跟踪和管理
- **可观测性**：详细的执行日志和监控

### 实现原理

```go
// Workflow 编排流程
数据输入 -> 预处理 -> 意图分析 -> 上下文构建 -> 消息构建 -> GLM模型 -> 后处理 -> 输出
    ↓         ↓         ↓          ↓          ↓         ↓         ↓
  用户输入   清理数据   关键词匹配   构建上下文   构建消息   模型调用   格式化输出
```

### 核心代码

```go
// Workflow数据流结构
type WorkflowData struct {
    UserInput     string                 `json:"user_input"`
    Intent        string                 `json:"intent"`
    Context       map[string]interface{} `json:"context"`
    Messages      []map[string]interface{} `json:"messages"`
    Response      string                 `json:"response"`
    Metadata      map[string]interface{} `json:"metadata"`
}

// 数据预处理
func (s *WorkflowService) preprocessData(ctx context.Context, data *WorkflowData) error {
    data.UserInput = strings.TrimSpace(data.UserInput)
    data.Context["timestamp"] = time.Now().Unix()
    data.Context["user_id"] = "user_123"
    data.Context["session_id"] = "session_456"
    return nil
}

// 意图分析
func (s *WorkflowService) analyzeIntent(ctx context.Context, data *WorkflowData) error {
    input := strings.ToLower(data.UserInput)
    
    if strings.Contains(input, "你好") || strings.Contains(input, "hello") {
        data.Intent = "greeting"
    } else if strings.Contains(input, "？") || strings.Contains(input, "?") {
        data.Intent = "question"
    } else if strings.Contains(input, "写") || strings.Contains(input, "创作") {
        data.Intent = "creative"
    } else {
        data.Intent = "general"
    }
    
    data.Context["intent"] = data.Intent
    data.Context["confidence"] = 0.9
    return nil
}
```

### API 接口

- **同步调用**: `POST /api/v1/eino/workflow/invoke`
- **流式调用**: `POST /api/v1/eino/workflow/stream`

### 使用示例

```bash
# 同步调用
curl -X POST http://localhost:8888/api/v1/eino/workflow/invoke \
  -H "Content-Type: application/json" \
  -d '{"message": "你好，请介绍一下自己"}'

# 流式调用
curl -X POST http://localhost:8888/api/v1/eino/workflow/stream \
  -H "Content-Type: application/json" \
  -d '{"message": "请写一首关于春天的短诗"}'
```

## 三种编排方式对比

| 特性 | Chain | Graph | Workflow |
|------|-------|-------|----------|
| **复杂度** | 简单 | 中等 | 复杂 |
| **适用场景** | 线性对话 | 条件分支 | 企业级应用 |
| **数据流转** | 简单 | 中等 | 复杂 |
| **可扩展性** | 低 | 中 | 高 |
| **性能** | 高 | 中 | 中 |
| **维护成本** | 低 | 中 | 高 |

## 选择建议

### 使用 Chain 的场景
- 简单的问答对话
- 线性处理流程
- 快速原型开发
- 对性能要求高的场景

### 使用 Graph 的场景
- 需要意图识别的对话
- 有分支逻辑的业务
- 需要条件判断的流程
- 中等复杂度的应用

### 使用 Workflow 的场景
- 企业级应用
- 复杂的数据处理
- 需要完整监控和日志
- 高可扩展性要求

## 扩展开发

### 添加新的意图类型

在 Graph 编排中添加新的意图：

```go
func (s *GraphService) identifyIntent(ctx context.Context, message string) (string, error) {
    message = strings.ToLower(message)
    
    // 添加新的意图识别
    if strings.Contains(message, "翻译") || strings.Contains(message, "translate") {
        return "translation", nil
    }
    
    // 现有意图...
    return "default", nil
}

// 添加对应的处理方法
func (s *GraphService) handleTranslation(ctx context.Context, message string) (*types.ChatResponse, error) {
    // 翻译处理逻辑
}
```

### 添加新的数据字段

在 Workflow 编排中添加新的数据字段：

```go
type WorkflowData struct {
    UserInput     string                 `json:"user_input"`
    Intent        string                 `json:"intent"`
    Context       map[string]interface{} `json:"context"`
    Messages      []map[string]interface{} `json:"messages"`
    Response      string                 `json:"response"`
    Metadata      map[string]interface{} `json:"metadata"`
    
    // 添加新字段
    Language      string                 `json:"language"`
    Emotion       string                 `json:"emotion"`
    Priority      int                    `json:"priority"`
}
```

## 监控和调试

### 日志记录

每个编排方式都有详细的日志记录：

```go
log.Printf("Chain步骤1 - 预处理完成: %s", processedInput)
log.Printf("Graph节点1 - 意图识别完成: %s", intent)
log.Printf("Workflow节点1 - 数据预处理完成")
```

### 性能监控

可以通过元数据跟踪性能：

```go
data.Metadata["processing_time"] = time.Now().Unix() - data.Context["timestamp"].(int64)
data.Metadata["response_length"] = len(data.Response)
```

## 总结

Eino 框架提供了三种不同复杂度的编排方式，可以根据具体需求选择合适的编排方式：

1. **Chain**：适合简单的线性流程
2. **Graph**：适合有分支逻辑的复杂流程
3. **Workflow**：适合企业级的高复杂度应用

每种方式都支持同步和流式两种调用模式，可以根据实际场景选择最合适的组合。
