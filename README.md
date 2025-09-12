# LLM Demo - 大模型交互示例项目

这是一个基于 [eino](https://www.cloudwego.io/zh/docs/eino/overview/eino_open_source/) 框架和 go-zero 的大模型交互示例项目，展示了如何与大模型进行 Invoke 和 Stream 两种方式的交互。

**✅ 已集成真实的 GLM4.5 API 调用**

## 项目特性

- 🚀 基于 go-zero 框架构建，提供高性能的 HTTP 服务
- 🤖 **真实 GLM4.5 大模型集成** - 已集成真实的智谱清言 API
- 📡 提供两种交互模式：同步调用（Invoke）和流式调用（Stream）
- 🔧 完整的配置管理和日志系统
- 📚 详细的代码注释和文档
- ⚡ 生产就绪的 API 调用实现

## 项目结构

```
llm-demo/
├── api/                    # API 定义文件
│   └── chat.api           # 聊天相关接口定义
├── etc/                   # 配置文件目录
│   └── config.yaml        # 应用配置文件
├── internal/              # 内部代码
│   ├── config/            # 配置结构定义
│   │   └── config.go
│   ├── handler/           # HTTP 处理器
│   │   └── chat_handler.go
│   ├── logic/             # 业务逻辑层
│   │   └── llm_service.go
│   ├── svc/               # 服务上下文
│   │   └── service_context.go
│   └── types/             # 类型定义
│       └── types.go
├── main.go                # 程序入口
├── go.mod                 # Go 模块定义
└── README.md              # 项目说明
```

## 快速开始

### 1. 安装依赖

```bash
go mod tidy
```

### 2. 配置模型参数

编辑 `etc/config.yaml` 文件，配置你的大模型参数：

```yaml
LLM:
  GLM4:
    Name: "glm4.5"
    URL: "https://open.bigmodel.cn/api/paas/v4/"
    AppKey: "your-app-key-here"
    Model: "glm-4-0520"
    MaxTokens: 4096
    Temperature: 0.7
    TopP: 0.9
```

### 3. 运行服务

```bash
go run main.go
```

服务将在 `http://localhost:8888` 启动。

### 4. 访问 API 文档

打开浏览器访问 `http://localhost:8888` 查看 API 文档。

## API 接口

### 同步聊天接口

**POST** `/api/v1/chat/invoke`

传统的请求-响应模式，适合需要完整响应的场景。

**请求体：**
```json
{
  "message": "你好，请介绍一下自己",
  "model": "glm-4-0520"  // 可选，不指定则使用默认模型
}
```

**响应：**
```json
{
  "code": 200,
  "message": "success",
  "data": "你好！我是一个AI助手...",
  "model": "glm-4-0520",
  "usage": {
    "promptTokens": 50,
    "completionTokens": 30,
    "totalTokens": 80
  }
}
```

### 流式聊天接口

**POST** `/api/v1/chat/stream`

实时流式响应模式，适合需要实时显示生成内容的场景。

**请求体：**
```json
{
  "message": "你好，请介绍一下自己",
  "model": "glm-4-0520"  // 可选
}
```

**响应：** Server-Sent Events (SSE) 格式
```
data: 你好！
data: 我是
data: 一个
data: AI助手...
data: {"code":200,"message":"success","data":"","model":"glm-4-0520","done":true,"usage":{"promptTokens":50,"completionTokens":30,"totalTokens":80}}
```

## 测试示例

### 使用 curl 测试同步接口

```bash
curl -X POST http://localhost:8888/api/v1/chat/invoke \
  -H "Content-Type: application/json" \
  -d '{"message": "你好，请介绍一下自己"}'
```

### 使用 curl 测试流式接口

```bash
curl -X POST http://localhost:8888/api/v1/chat/stream \
  -H "Content-Type: application/json" \
  -d '{"message": "你好，请介绍一下自己"}'
```

### 健康检查

```bash
curl http://localhost:8888/api/v1/health/check
```

## 测试结果

### 同步调用测试
```bash
$ curl -X POST http://localhost:8888/api/v1/chat/invoke \
  -H "Content-Type: application/json" \
  -d '{"message": "你好，请简单介绍一下自己"}'

{
  "code": 200,
  "message": "success",
  "data": "你好，我是智谱清言，是智谱 AI 公司训练的语言模型...",
  "model": "glm-4-0520",
  "usage": {
    "promptTokens": 11,
    "completionTokens": 58,
    "totalTokens": 69
  }
}
```

### 流式调用测试
```bash
$ curl -X POST http://localhost:8888/api/v1/chat/stream \
  -H "Content-Type: application/json" \
  -d '{"message": "请写一首关于春天的短诗"}'

data: 春风
data: 拂
data: 面
data: 花
data: 千
data: 树
...
```

## 核心实现说明

### 1. 配置管理

项目使用 go-zero 的配置系统，支持 YAML 格式的配置文件。所有大模型相关的配置都集中在 `etc/config.yaml` 中。

### 2. 服务架构

- **Handler 层**：处理 HTTP 请求和响应
- **Logic 层**：实现具体的业务逻辑
- **Service 层**：管理服务上下文和依赖注入

### 3. 两种交互模式

#### Invoke 模式（同步）
- 发送完整请求到模型
- 等待模型返回完整响应
- 适合需要完整内容的场景
- 实现简单，易于调试

#### Stream 模式（流式）
- 建立流式连接
- 实时接收模型生成的内容片段
- 适合需要实时显示的场景
- 用户体验更好，响应更快

### 4. 错误处理

项目实现了完整的错误处理机制：
- 参数验证
- 业务逻辑错误处理
- HTTP 状态码规范
- 详细的错误日志

### 5. 日志系统

使用 go-zero 的 logx 进行日志管理：
- 结构化日志输出
- 不同级别的日志记录
- 请求链路追踪

## 扩展说明

### 集成真实的大模型 API

当前项目使用模拟数据，要集成真实的大模型 API，需要：

1. 在 `internal/logic/llm_service.go` 中实现真实的 API 调用
2. 处理 API 认证和错误响应
3. 实现流式响应的解析
4. 添加重试和熔断机制

### 添加更多模型支持

可以通过以下方式扩展：

1. 在配置文件中添加更多模型配置
2. 在 `LLMService` 中添加模型选择逻辑
3. 实现不同模型的适配器模式

### 数据持久化

可以添加数据库支持：
1. 集成 go-zero 的数据库组件
2. 创建聊天记录表
3. 实现聊天历史查询功能

## 技术栈

- **框架**：go-zero
- **大模型框架**：eino
- **语言**：Go 1.21+
- **配置**：YAML
- **日志**：go-zero logx

## 许可证

MIT License