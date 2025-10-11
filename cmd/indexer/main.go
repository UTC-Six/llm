package main

import (
	"context"
	"os"

	"github.com/cloudwego/eino-ext/components/embedding/ark"
	"github.com/cloudwego/eino-ext/components/indexer/milvus"
	"github.com/cloudwego/eino/schema"
	"github.com/joho/godotenv"
	cli "github.com/milvus-io/milvus-sdk-go/v2/client"
	"github.com/milvus-io/milvus-sdk-go/v2/entity"
	"github.com/zeromicro/go-zero/core/logx"
)

var MilvusCli cli.Client

func initClient(ctx context.Context) {
	client, err := cli.NewClient(ctx, cli.Config{
		Address: "localhost:19530",
		DBName:  "llm",
	})
	if err != nil {
		panic(err)
	}
	MilvusCli = client
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}
	ctx := context.Background()
	initClient(ctx)

	embedder, err := ark.NewEmbedder(ctx, &ark.EmbeddingConfig{
		APIKey: os.Getenv("API_KEY"),
		Model:  os.Getenv("EMBEDDER_MODEL"),
	})
	if err != nil {
		panic(err)
	}
	var (
		collection = "test"
		fields     = []*entity.Field{
			{
				Name:       "id",
				PrimaryKey: true,
				DataType:   entity.FieldTypeVarChar,
				TypeParams: map[string]string{"max_length": "256"},
			},
			{
				Name:       "vector",
				DataType:   entity.FieldTypeBinaryVector,
				TypeParams: map[string]string{"dim": "81920"},
			},
			{
				Name:       "content",
				DataType:   entity.FieldTypeVarChar,
				TypeParams: map[string]string{"max_length": "8192"},
			},
			{
				Name:       "metadata",
				DataType:   entity.FieldTypeJSON,
				TypeParams: map[string]string{"max_length": "8192"},
			},
		}
	)

	indexer, err := milvus.NewIndexer(ctx, &milvus.IndexerConfig{
		Client:     MilvusCli,
		Collection: collection,
		Fields:     fields,
		Embedding:  embedder,
	})
	if err != nil {
		panic(err)
	}

	docs := []*schema.Document{
		{
			ID:       "1",
			Content:  "你说的对，但是原神是一款二次元开放大世界游戏",
			MetaData: map[string]any{"author": "木乔"},
		},
		{
			ID:       "2",
			Content:  "你说的对，但是原神是一款二次元开放大世界游戏",
			MetaData: map[string]any{"author": "鹰角"},
		},
	}

	ids, err := indexer.Store(ctx, docs)
	if err != nil {
		panic(err)
	}
	logx.Infof("ids=%+v", ids)
}
