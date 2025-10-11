package main

import (
	"context"
	"os"

	"github.com/cloudwego/eino-ext/components/embedding/ark"
	"github.com/cloudwego/eino-ext/components/retriever/milvus"
	"github.com/joho/godotenv"
	cli "github.com/milvus-io/milvus-sdk-go/v2/client"
	"github.com/zeromicro/go-zero/core/logx"
)

var MilvusCli cli.Client

func initClient(ctx context.Context) {
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}
	client, err := cli.NewClient(ctx, cli.Config{
		Address: os.Getenv("MILVUS_ADDRESS"),
		DBName:  os.Getenv("DB_NAME"),
	})
	if err != nil {
		panic(err)
	}
	MilvusCli = client
}

func main() {
	ctx := context.Background()
	initClient(ctx)

	embedder, err := ark.NewEmbedder(ctx, &ark.EmbeddingConfig{
		APIKey: os.Getenv("API_KEY"),
		Model:  os.Getenv("EMBEDDER_MODEL"),
	})
	if err != nil {
		panic(err)
	}
	retriever, err := milvus.NewRetriever(ctx, &milvus.RetrieverConfig{
		Client:       MilvusCli,
		Collection:   "test",
		VectorField:  "vector",
		OutputFields: []string{"id", "content", "metadata"},
		TopK:         2,
		Embedding:    embedder,
	})
	if err != nil {
		panic(err)
	}
	results, err := retriever.Retrieve(ctx, "原神")
	if err != nil {
		panic(err)
	}
	for _, result := range results {
		logx.Infof("result=%+v", result.MetaData)
	}
}
