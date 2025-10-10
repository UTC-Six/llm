package main

import (
	"context"

	"github.com/cloudwego/eino-ext/components/embedding/ark"
	"github.com/zeromicro/go-zero/core/logx"
)

var (
	BaseURL         = "https://ark.cn-beijing.volces.com/api/v3"
	ApiKey          = "39916561-ff4e-42c3-adc3-efb164dabe1b"
	ModelDouBaoSeed = "doubao-seed-1-6-250615"
	EmbedderModel   = "doubao-embedding-text-240715"
)

func main() {
	ctx := context.Background()
	embedder, err := ark.NewEmbedder(ctx, &ark.EmbeddingConfig{
		APIKey: ApiKey,
		Model:  EmbedderModel,
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
