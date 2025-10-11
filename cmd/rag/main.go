package main

import (
	"context"
	"fmt"
	"os"

	"github.com/cloudwego/eino-ext/components/document/transformer/splitter/markdown"
	"github.com/cloudwego/eino-ext/components/embedding/ark"
	"github.com/cloudwego/eino-ext/components/indexer/milvus"
	milvusRetriever "github.com/cloudwego/eino-ext/components/retriever/milvus"
	"github.com/cloudwego/eino/schema"
	"github.com/joho/godotenv"
	cli "github.com/milvus-io/milvus-sdk-go/v2/client"
	"github.com/milvus-io/milvus-sdk-go/v2/entity"
	"github.com/zeromicro/go-zero/core/logx"
)

var MilvusCli cli.Client

func initMilvusClient(ctx context.Context) {
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
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}

	ctx := context.Background()

	embedder, err := ark.NewEmbedder(ctx, &ark.EmbeddingConfig{
		APIKey: os.Getenv("API_KEY"),
		Model:  os.Getenv("EMBEDDER_MODEL")})
	if err != nil {
		panic(err)
	}

	initMilvusClient(ctx)
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

	retriever, err := milvusRetriever.NewRetriever(ctx, &milvusRetriever.RetrieverConfig{
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
			ID:      "doc",
			Content: string(bs),
		},
	}

	results, err := splitter.Transform(ctx, docs)
	if err != nil {
		panic(err)
	}

	for i, doc := range results {
		doc.ID = fmt.Sprintf("%s_%d", doc.ID, i+1)
	}

	_, err = indexer.Store(ctx, results)
	if err != nil {
		panic(err)
	}

	res, err := retriever.Retrieve(ctx, "v1.1.0 - 真实 GLM API 集成")
	if err != nil {
		panic(err)
	}
	for _, result := range res {
		logx.Infof("id=%s, content=%s, metadata=%v", result.ID, result.Content, result.MetaData)
	}
}
