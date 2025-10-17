package main

import (
	"context"
	"fmt"
	"time"

	"github.com/cloudwego/eino-ext/components/tool/browseruse"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/components/tool/utils"
	"github.com/cloudwego/eino/schema"
	"github.com/zeromicro/go-zero/core/logx"
)

func main() {
	ctx := context.Background()
	bt, err := browseruse.NewBrowserUseTool(ctx, &browseruse.Config{})
	if err != nil {
		panic(err)
	}

	url := "https://www.bilibili.com"
	result, err := bt.Execute(&browseruse.Param{
		Action: browseruse.ActionGoToURL,
		URL:    &url,
	})
	if err != nil {
		panic(err)
	}
	logx.Infof("%+v", result)
	time.Sleep(10 * time.Second)
	bt.Cleanup()
}

type Game struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

type InputParams struct {
	Name string `json:"name" jsonschema:"description=the name of game"`
}

func GetGame(_ context.Context, params *InputParams) (string, error) {
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
		Name:  "get_name",
		Desc:  "get game url by name",
		Extra: nil,
		ParamsOneOf: schema.NewParamsOneOfByParams(
			map[string]*schema.ParameterInfo{
				"name": &schema.ParameterInfo{
					Type:     schema.String,
					Desc:     "get game`s url by name",
					Required: true,
				},
			},
		),
	}, GetGame)
}
