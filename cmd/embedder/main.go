package main

import (
	"context"
	"os"

	"github.com/cloudwego/eino-ext/components/embedding/ark"
	"github.com/joho/godotenv"
	"github.com/zeromicro/go-zero/core/logx"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}
	ctx := context.Background()
	embedder, err := ark.NewEmbedder(ctx, &ark.EmbeddingConfig{
		APIKey: os.Getenv("API_KEY"),
		Model:  os.Getenv("EMBEDDER_MODEL"),
	})
	if err != nil {
		logx.Errorf("ark.NewEmbedder failed, err=%v", err)
		return
	}
	input := []string{"你好, 泥豪", "Hello carey"}

	embeddings, err := embedder.EmbedStrings(ctx, input)
	if err != nil {
		logx.Errorf("embedder.EmbedStrings failed, err=%v, input=%+v", err, input)
		return
	}
	for _, embedding := range embeddings {
		logx.Infof("embedding=%f", embedding)
	}
}
