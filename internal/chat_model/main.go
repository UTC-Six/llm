package main

import (
	"context"

	"github.com/cloudwego/eino-ext/components/model/openai"
	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/schema"
	"github.com/zeromicro/go-zero/core/logx"
)

/*
编排范式:
 - Chain: 链式有向无环图
 - Graph: 有向图或有向无环图
 - Workflow: 有字段映射能力的有向无环图
*/

func main() {
	ctx := context.Background()
	cm, err := openai.NewChatModel(ctx, &openai.ChatModelConfig{
		APIKey:  "2bfeecd60bf44969acf97f24ce734fa4.ytUPBofnk8rSM3Ag",
		BaseURL: "https://open.bigmodel.cn/api/paas/v4/",
		Model:   "glm-4.5",
	})
	if err != nil {
		logx.Errorf("模型初始化失败: %v", err)
		return
	}

	resp, err := cm.Generate(ctx, []*schema.Message{
		{
			Role:    "user",
			Content: "你好，请介绍一下自己",
		},
	}, model.WithTemperature(0.8))
	logx.Infof("响应结果，result=%+v, err=%v", resp, err)
}
