package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/cloudwego/eino-ext/components/model/ark"
	"github.com/cloudwego/eino/callbacks"
	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/components/tool/utils"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
	callbackHelpers "github.com/cloudwego/eino/utils/callbacks"
	"github.com/joho/godotenv"
)

/*
// lambda demo code

func main() {
	ctx := context.Background()
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}

	chatModel, err := ark.NewChatModel(ctx, &ark.ChatModelConfig{
		APIKey: os.Getenv("API_KEY"),
		Model:  os.Getenv("CHAT_MODEL_NAME"),
	})
	if err != nil {
		panic(err)
	}

	lambda := compose.InvokableLambda(func(ctx context.Context, input string) (output []*schema.Message, err error) {
		return []*schema.Message{
			{
				Role:    schema.User,
				Content: fmt.Sprintf("%s 回答结尾加上desuwa", input),
			},
		}, nil
	})

	chain := compose.NewChain[string, *schema.Message]()
	r, err := chain.AppendLambda(lambda).AppendChatModel(chatModel).Compile(ctx)
	if err != nil {
		panic(err)
	}
	resp, err := r.Invoke(ctx, "你好，你可以告诉我你的名字吗？")
	if err != nil {
		panic(err)
	}
	logx.Infof("result=%s", resp.Content)
}*/

type Game struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

type InputParams struct {
	Name string `json:"name" jsonschema:"description=the name of game"`
}

func GetGame(_ context.Context, params *InputParams) (output string, err error) {
	GameSet := []*Game{
		{
			Name: "原神",
			Url:  "https://yuanshen.com/tool",
		},
		{
			Name: "鸣潮",
			Url:  "https://mingchao.com/tool",
		},
		{
			Name: "明日方舟",
			Url:  "https://mingrifangzhou.com/tool",
		},
	}

	for _, game := range GameSet {
		if game.Name == params.Name {
			return game.Url, nil
		}
	}
	return "", fmt.Errorf("name=%s`s game not found", params.Name)
}
func CreateTool() tool.InvokableTool {
	return utils.NewTool(&schema.ToolInfo{
		Name: "get_name",
		Desc: "get game url by name",
		ParamsOneOf: schema.NewParamsOneOfByParams(
			map[string]*schema.ParameterInfo{
				"name": {
					Type:     schema.String,
					Desc:     "get game`s url by name",
					Required: true,
				},
			},
		),
	}, GetGame)
}

// chain compose demo with callback、lambda、tool、chat_model
func main() {
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}

	ctx := context.Background()

	modelHandler := &callbackHelpers.ModelCallbackHandler{
		OnEnd: func(ctx context.Context, runInfo *callbacks.RunInfo, output *model.CallbackOutput) context.Context {
			fmt.Println("模型的思考过程为：")
			fmt.Println(output.Message.ReasoningContent)
			return ctx
		},
	}

	toolHandler := &callbackHelpers.ToolCallbackHandler{
		OnStart: func(ctx context.Context, info *callbacks.RunInfo, input *tool.CallbackInput) context.Context {
			fmt.Printf("调用工具前的输入参数为：%v \n", input.ArgumentsInJSON)
			return ctx
		},
		OnEnd: func(ctx context.Context, info *callbacks.RunInfo, output *tool.CallbackOutput) context.Context {
			fmt.Printf("工具执行完成，输出内容为：%v \n", output.Response)
			return ctx
		},
	}

	handler := callbackHelpers.NewHandlerHelper().ChatModel(modelHandler).Tool(toolHandler).Handler()
	timeout := 30 * time.Second
	chatModel, err := ark.NewChatModel(ctx, &ark.ChatModelConfig{
		APIKey:  os.Getenv("API_KEY"),
		Model:   os.Getenv("CHAT_MODEL_NAME"),
		Timeout: &timeout,
	})
	if err != nil {
		panic(err)
	}

	// 绑定工具
	getNameTool := CreateTool()
	info, err := getNameTool.Info(ctx)
	if err != nil {
		panic(err)
	}

	fmt.Printf("=============%+v \n", info)
	infos := []*schema.ToolInfo{
		info,
	}

	err = chatModel.BindTools(infos)
	if err != nil {
		panic(err)
	}
	chain := compose.NewChain[[]*schema.Message, []*schema.Message]()
	toolNode, err := compose.NewToolNode(ctx, &compose.ToolsNodeConfig{
		Tools: []tool.BaseTool{
			getNameTool,
		},
	})
	if err != nil {
		panic(err)
	}

	chain.AppendChatModel(chatModel, compose.WithNodeName("chat_model")).AppendToolsNode(toolNode, compose.WithNodeName("tools"))
	r, err := chain.Compile(ctx)
	if err != nil {
		panic(err)
	}
	resp, err := r.Invoke(ctx, []*schema.Message{
		{
			Role:    schema.User,
			Content: "请告诉我原神、鸣潮、明日的 URL 是啥？",
		},
	}, compose.WithCallbacks(handler))
	if err != nil {
		panic(err)
	}
	for _, result := range resp {
		fmt.Printf("%+v", result)
	}
}
