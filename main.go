package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/UTC-Six/llm/internal/config"
	"github.com/UTC-Six/llm/internal/handler"
	"github.com/UTC-Six/llm/internal/svc"
)

// 配置文件路径
var configFile = flag.String("f", "etc/config.yaml", "配置文件路径")

func main() {
	flag.Parse()

	// 加载配置文件
	var c config.Config
	if err := loadConfig(*configFile, &c); err != nil {
		log.Fatalf("加载配置文件失败: %v", err)
	}

	log.Printf("LLM Demo 服务启动中...")

	// 创建服务上下文
	ctx := svc.NewServiceContext(c)

	// 创建HTTP服务器
	server := &http.Server{
		Addr:    ":8888",
		Handler: setupRoutes(ctx),
	}

	// 启动服务器
	go func() {
		log.Printf("LLM Demo 服务启动成功，监听地址: %s", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("服务器启动失败: %v", err)
		}
	}()

	// 等待中断信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Printf("正在关闭服务器...")
}

// loadConfig 加载配置文件
func loadConfig(configFile string, c *config.Config) error {
	// 这里简化处理，直接设置默认值
	// 在实际项目中应该使用YAML解析库
	c.Host = "0.0.0.0"
	c.Port = 8888
	c.Mode = "dev"
	c.LLM.GLM4.Name = "glm4.5"
	c.LLM.GLM4.URL = "https://open.bigmodel.cn/api/paas/v4/"
	c.LLM.GLM4.AppKey = "2bfeecd60bf44969acf97f24ce734fa4.ytUPBofnk8rSM3Ag"
	c.LLM.GLM4.Model = "glm-4-0520"
	c.LLM.GLM4.MaxTokens = 4096
	c.LLM.GLM4.Temperature = 0.7
	c.LLM.GLM4.TopP = 0.9
	return nil
}

// setupRoutes 设置路由
func setupRoutes(ctx *svc.ServiceContext) http.Handler {
	mux := http.NewServeMux()

	// 创建处理器
	chatHandler := handler.NewChatHandler(ctx)

	// 注册聊天相关路由
	mux.HandleFunc("/api/v1/chat/invoke", chatHandler.ChatInvoke)
	mux.HandleFunc("/api/v1/chat/stream", chatHandler.ChatStream)

	// 注册健康检查路由
	mux.HandleFunc("/api/v1/health/check", healthCheckHandler)

	// 注册根路径，显示API文档
	mux.HandleFunc("/", apiDocHandler)

	log.Printf("路由注册完成")
	return mux
}

// healthCheckHandler 健康检查处理器
func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ok","message":"LLM Demo 服务运行正常"}`))
}

// apiDocHandler API文档处理器
func apiDocHandler(w http.ResponseWriter, r *http.Request) {
	html := `
<!DOCTYPE html>
<html>
<head>
    <title>LLM Demo API</title>
    <meta charset="utf-8">
    <style>
        body { font-family: Arial, sans-serif; margin: 40px; }
        .endpoint { background: #f5f5f5; padding: 15px; margin: 10px 0; border-radius: 5px; }
        .method { color: #007bff; font-weight: bold; }
        .path { color: #28a745; font-family: monospace; }
        .description { color: #666; margin-top: 5px; }
    </style>
</head>
<body>
    <h1>LLM Demo API 文档</h1>
    <p>这是一个基于 eino 框架的大模型交互示例项目</p>
    
    <h2>API 接口</h2>
    
    <div class="endpoint">
        <div><span class="method">POST</span> <span class="path">/api/v1/chat/invoke</span></div>
        <div class="description">同步聊天接口 - 传统的请求-响应模式</div>
        <p>请求体: {"message": "你好", "model": "glm-4-0520"}</p>
    </div>
    
    <div class="endpoint">
        <div><span class="method">POST</span> <span class="path">/api/v1/chat/stream</span></div>
        <div class="description">流式聊天接口 - 实时流式响应模式</div>
        <p>请求体: {"message": "你好", "model": "glm-4-0520"}</p>
    </div>
    
    <div class="endpoint">
        <div><span class="method">GET</span> <span class="path">/api/v1/health/check</span></div>
        <div class="description">健康检查接口</div>
    </div>
    
    <h2>使用说明</h2>
    <ul>
        <li><strong>同步模式</strong>: 发送请求后等待完整响应，适合需要完整内容的场景</li>
        <li><strong>流式模式</strong>: 实时接收生成的内容片段，适合需要实时显示的场景</li>
        <li>支持指定不同的模型，如果不指定则使用默认配置的模型</li>
        <li>所有接口都支持跨域访问</li>
    </ul>
    
    <h2>测试示例</h2>
    <p>使用 curl 测试同步接口：</p>
    <pre>curl -X POST http://localhost:8888/api/v1/chat/invoke \
  -H "Content-Type: application/json" \
  -d '{"message": "你好，请介绍一下自己"}'</pre>
    
    <p>使用 curl 测试流式接口：</p>
    <pre>curl -X POST http://localhost:8888/api/v1/chat/stream \
  -H "Content-Type: application/json" \
  -d '{"message": "你好，请介绍一下自己"}'</pre>
</body>
</html>
`
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(html))
}
