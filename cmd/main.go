package main

import (
	"context"
	"os"

	"github.com/cloudwego/eino-ext/components/model/ark"
	"github.com/cloudwego/eino/components/prompt"
	"github.com/cloudwego/eino/schema"
	"github.com/joho/godotenv"
	"github.com/zeromicro/go-zero/core/logx"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}
	ctx := context.Background()
	model, err := ark.NewChatModel(ctx, &ark.ChatModelConfig{
		APIKey: os.Getenv("API_KEY"),
		Model:  os.Getenv("CHAT_MODEL_NAME"),
	})

	template := prompt.FromMessages(schema.FString, schema.SystemMessage("你是一个{role}"), &schema.Message{
		Role:    schema.User,
		Content: "请帮帮我，施瓦罗先生，帮我解决{task}",
	})

	params := map[string]interface{}{
		"role": "高中数学老师",
		"task": "斐波那契数列的算法问题",
	}
	msgs, err := template.Format(ctx, params)
	if err != nil {
		logx.Errorf("template.Format failed, err=%v, params=%+v", err, params)
		return
	}

	/*	input := []*schema.Message{
		schema.SystemMessage("你是一个高中可爱美少女"),
		schema.UserMessage("你好呀"),
	}*/
	resp, err := model.Generate(ctx, msgs)
	if err != nil {
		logx.Errorf("model.Generate failed, err=%v, input=%+v", err, msgs)
		return
	}
	logx.Infof("resp=%s", resp.Content)
	/*reader, err := model.Stream(ctx, input)
	if err != nil {
		logx.Errorf("model.Stream failed, err=%v, input=%+v", err, input)
		return
	}
	defer reader.Close()
	for {
		chunk, err := reader.Recv()
		if err != nil {
			logx.Errorf("reader.Read failed, err=%v", err)
			return
		}
		//logx.Info("chunk=%s", chunk.Content)
		print(chunk.Content)
		//logx.Infof("---------------chunk start------------------------------")
		//logx.Infof("chunk=%+v", chunk)
		//logx.Infof("---------------chunk end------------------------------")

	}*/
}
