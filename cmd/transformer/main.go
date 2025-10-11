package main

import (
	"context"
	"os"

	"github.com/cloudwego/eino-ext/components/document/transformer/splitter/markdown"
	"github.com/cloudwego/eino/schema"
	"github.com/zeromicro/go-zero/core/logx"
)

func main() {
	ctx := context.Background()
	splitter, err := markdown.NewHeaderSplitter(ctx, &markdown.HeaderConfig{
		Headers: map[string]string{
			"#":   "H1",
			"##":  "H2",
			"###": "H3",
		},
		TrimHeaders: false,
	})
	if err != nil {
		panic(err)
	}

	content, err := os.OpenFile("./cmd/transformer/document.md", os.O_CREATE|os.O_RDWR, 0755)
	if err != nil {
		panic(err)
	}
	defer content.Close()
	bs, err := os.ReadFile("./cmd/transformer/document.md")
	if err != nil {
		panic(err)
	}

	docs := []*schema.Document{
		{
			ID:      "doc1",
			Content: string(bs),
		},
	}

	results, err := splitter.Transform(ctx, docs)
	if err != nil {
		panic(err)
	}

	for _, doc := range results {
		logx.Infof("content=%+v", doc.Content)
		logx.Info("标题层级")
		for k, v := range doc.MetaData {
			if k == "h1" || k == "h2" || k == "h3" {
				logx.Infof("%s=%v", k, v)
			}
		}
	}
}
