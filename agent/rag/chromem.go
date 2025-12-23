package rag

import (
	"context"
	"strconv"

	"github.com/philippgille/chromem-go"
)

// 向量数据库
type ChromemManager struct {
	db         *chromem.DB
	collection *chromem.Collection
}

const ollamaEmbedModelName = "nomic-embed-text-v2-moe"

func StartChromem() (*ChromemManager, error) {
	db := chromem.NewDB()

	collection, err := db.CreateCollection("knowledge", nil, chromem.NewEmbeddingFuncOllama(ollamaEmbedModelName, ""))
	if err != nil {
		return nil, err
	}

	chromemManager := ChromemManager{db: db, collection: collection}

	return &chromemManager, nil
}

// 向量化+存储
func (this *ChromemManager) AddDocuments(index int, content string) error {
	ctx := context.Background()
	return this.collection.AddDocument(ctx, chromem.Document{ID: strconv.Itoa(index), Content: content})
}

// 向量检索
func (this *ChromemManager) Query(text string, nResults int) ([]int, error) {
	ctx := context.Background()
	res, err := this.collection.Query(ctx, text, nResults, nil, nil)
	if err != nil {
		return nil, err
	}
	var indexArr []int
	for i := 0; i < len(res); i++ {
		index, _ := strconv.Atoi(res[i].ID)
		indexArr = append(indexArr, index)
	}
	return indexArr, nil
}
