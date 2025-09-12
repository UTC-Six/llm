# LLM Demo 集成指南

本文档详细说明如何将项目从模拟实现升级为真实的大模型API集成。

## 项目架构说明

### 当前实现
- **模拟模式**: 使用 `LLMService` 提供模拟的大模型响应
- **真实模式**: 使用 `RealLLMService` 展示如何集成真实的 eino 框架

### 核心组件

1. **配置管理** (`internal/config/config.go`)
   - 集中管理所有大模型相关配置
   - 支持多模型配置
   - 支持环境变量覆盖

2. **业务逻辑层** (`internal/logic/`)
   - `llm_service.go`: 模拟实现，用于演示和测试
   - `real_llm_service.go`: 真实实现，展示 eino 集成

3. **HTTP处理器** (`internal/handler/chat_handler.go`)
   - 统一的HTTP接口
   - 支持同步和流式两种模式
   - 完整的错误处理

## 集成真实大模型API

### 1. 安装 eino 框架

```bash
go get github.com/cloudwego/eino
go get github.com/cloudwego/eino-ext
```

### 2. 更新 go.mod

```go
require (
    github.com/cloudwego/eino v0.0.0-20241220000000-000000000000
    github.com/cloudwego/eino-ext v0.0.0-20241220000000-000000000000
)
```

### 3. 实现真实API调用

参考 `internal/logic/real_llm_service.go` 中的注释代码：

#### 同步调用实现

```go
// 1. 创建 eino 客户端
client, err := eino.NewClient(eino.WithAPIKey(modelConfig.AppKey))
if err != nil {
    return nil, fmt.Errorf("创建eino客户端失败: %v", err)
}

// 2. 构建请求
request := &eino.ChatCompletionRequest{
    Model:       modelConfig.Model,
    Messages:    messages,
    MaxTokens:   modelConfig.MaxTokens,
    Temperature: modelConfig.Temperature,
    TopP:        modelConfig.TopP,
    Stream:      false,
}

// 3. 发送请求
response, err := client.ChatCompletion(ctx, request)
if err != nil {
    return nil, fmt.Errorf("调用大模型API失败: %v", err)
}

// 4. 处理响应
content := response.Choices[0].Message.Content
usage := response.Usage
```

#### 流式调用实现

```go
// 1. 创建 eino 客户端
client, err := eino.NewClient(eino.WithAPIKey(modelConfig.AppKey))
if err != nil {
    return fmt.Errorf("创建eino客户端失败: %v", err)
}

// 2. 构建流式请求
request := &eino.ChatCompletionRequest{
    Model:       modelConfig.Model,
    Messages:    messages,
    MaxTokens:   modelConfig.MaxTokens,
    Temperature: modelConfig.Temperature,
    TopP:        modelConfig.TopP,
    Stream:      true,
}

// 3. 发送流式请求
stream, err := client.ChatCompletionStream(ctx, request)
if err != nil {
    return fmt.Errorf("调用大模型流式API失败: %v", err)
}
defer stream.Close()

// 4. 处理流式响应
for {
    chunk, err := stream.Recv()
    if err == io.EOF {
        break
    }
    if err != nil {
        return fmt.Errorf("接收流式数据失败: %v", err)
    }
    
    // 写入响应数据
    if len(chunk.Choices) > 0 {
        content := chunk.Choices[0].Delta.Content
        if content != "" {
            _, err := writer.Write([]byte(fmt.Sprintf("data: %s\n\n", content)))
            if err != nil {
                return fmt.Errorf("写入流式响应失败: %v", err)
            }
        }
    }
    
    // 检查是否完成
    if chunk.Choices[0].FinishReason == "stop" {
        break
    }
}
```

### 4. 切换服务实现

在 `internal/svc/service_context.go` 中：

```go
// 使用模拟服务
llmService := logic.NewLLMService(c)

// 或者使用真实服务
// llmService := logic.NewRealLLMService(c)
```

## 配置说明

### 环境变量配置

```bash
# 大模型配置
export GLM4_API_KEY="your-api-key-here"
export GLM4_MODEL="glm-4-0520"
export GLM4_MAX_TOKENS="4096"
export GLM4_TEMPERATURE="0.7"
export GLM4_TOP_P="0.9"
```

### 配置文件示例

```yaml
# etc/config.yaml
Name: llm-demo
Host: 0.0.0.0
Port: 8888
Mode: dev

LLM:
  GLM4:
    Name: "glm4.5"
    URL: "https://open.bigmodel.cn/api/paas/v4/"
    AppKey: "${GLM4_API_KEY}"
    Model: "${GLM4_MODEL}"
    MaxTokens: 4096
    Temperature: 0.7
    TopP: 0.9

Log:
  ServiceName: llm-demo
  Mode: console
  Level: info
  Encoding: json
```

## 错误处理

### 常见错误类型

1. **API认证错误**
   - 检查 AppKey 是否正确
   - 检查 API 权限

2. **网络错误**
   - 实现重试机制
   - 设置超时时间

3. **模型错误**
   - 检查模型名称
   - 检查参数范围

### 错误处理示例

```go
func (s *RealLLMService) ChatInvoke(ctx context.Context, req *types.ChatRequest) (*types.ChatResponse, error) {
    // 设置超时
    ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
    defer cancel()
    
    // 实现重试机制
    var lastErr error
    for i := 0; i < 3; i++ {
        resp, err := s.callAPI(ctx, req)
        if err == nil {
            return resp, nil
        }
        lastErr = err
        
        // 指数退避
        time.Sleep(time.Duration(i+1) * time.Second)
    }
    
    return nil, fmt.Errorf("重试3次后仍然失败: %v", lastErr)
}
```

## 性能优化

### 1. 连接池

```go
// 创建带连接池的客户端
client, err := eino.NewClient(
    eino.WithAPIKey(modelConfig.AppKey),
    eino.WithHTTPClient(&http.Client{
        Transport: &http.Transport{
            MaxIdleConns:        100,
            MaxIdleConnsPerHost: 10,
            IdleConnTimeout:     90 * time.Second,
        },
    }),
)
```

### 2. 缓存机制

```go
// 实现响应缓存
type CacheService struct {
    cache map[string]*types.ChatResponse
    mutex sync.RWMutex
}

func (c *CacheService) Get(key string) (*types.ChatResponse, bool) {
    c.mutex.RLock()
    defer c.mutex.RUnlock()
    resp, ok := c.cache[key]
    return resp, ok
}
```

### 3. 限流控制

```go
// 实现限流
type RateLimiter struct {
    limiter *rate.Limiter
}

func (r *RateLimiter) Allow() bool {
    return r.limiter.Allow()
}
```

## 监控和日志

### 1. 指标监控

```go
// 添加监控指标
var (
    requestTotal = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "llm_requests_total",
            Help: "Total number of LLM requests",
        },
        []string{"model", "status"},
    )
    
    requestDuration = prometheus.NewHistogramVec(
        prometheus.HistogramOpts{
            Name: "llm_request_duration_seconds",
            Help: "Duration of LLM requests",
        },
        []string{"model"},
    )
)
```

### 2. 结构化日志

```go
// 使用结构化日志
log.WithFields(log.Fields{
    "model":     modelConfig.Model,
    "user_id":   userID,
    "request_id": requestID,
    "tokens":    usage.TotalTokens,
}).Info("LLM request completed")
```

## 测试

### 1. 单元测试

```go
func TestLLMService_ChatInvoke(t *testing.T) {
    // 测试模拟实现
    service := logic.NewLLMService(testConfig)
    resp, err := service.ChatInvoke(context.Background(), &types.ChatRequest{
        Message: "Hello",
    })
    
    assert.NoError(t, err)
    assert.NotEmpty(t, resp.Data)
}
```

### 2. 集成测试

```go
func TestRealLLMService_Integration(t *testing.T) {
    if testing.Short() {
        t.Skip("跳过集成测试")
    }
    
    // 测试真实API
    service := logic.NewRealLLMService(realConfig)
    resp, err := service.ChatInvoke(context.Background(), &types.ChatRequest{
        Message: "Hello",
    })
    
    assert.NoError(t, err)
    assert.NotEmpty(t, resp.Data)
}
```

## 部署

### 1. Docker 部署

```dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o llm-demo main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/llm-demo .
COPY --from=builder /app/etc ./etc
CMD ["./llm-demo"]
```

### 2. Kubernetes 部署

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: llm-demo
spec:
  replicas: 3
  selector:
    matchLabels:
      app: llm-demo
  template:
    metadata:
      labels:
        app: llm-demo
    spec:
      containers:
      - name: llm-demo
        image: llm-demo:latest
        ports:
        - containerPort: 8888
        env:
        - name: GLM4_API_KEY
          valueFrom:
            secretKeyRef:
              name: llm-secrets
              key: api-key
```

## 总结

本指南提供了从模拟实现到真实大模型API集成的完整路径。通过 eino 框架，可以轻松集成各种大模型服务，实现高性能、可扩展的AI应用。

关键要点：
1. 使用 eino 框架简化大模型集成
2. 实现完整的错误处理和重试机制
3. 添加监控和日志记录
4. 进行充分的测试
5. 考虑性能和安全性
