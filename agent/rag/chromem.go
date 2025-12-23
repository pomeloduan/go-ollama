package rag

import (
	"context"
	"strconv"

	"github.com/philippgille/chromem-go"
)

// 向量数据库
type ChromemManager struct {
	db            *chromem.DB
	collectionMap map[int]*chromem.Collection
}

const ollamaEmbedModelName = "nomic-embed-text-v2-moe"

func StartChromem() *ChromemManager {
	return &ChromemManager{db: chromem.NewDB(), collectionMap: make(map[int]*chromem.Collection)}
}

func (this *ChromemManager) newCollection(ragId int) error {
	collection, err := this.db.CreateCollection(
		"rag-"+strconv.Itoa(ragId),
		nil,
		chromem.NewEmbeddingFuncOllama(ollamaEmbedModelName, ""))
	if err != nil {
		return err
	}
	this.collectionMap[ragId] = collection
	return nil
}

// 向量化+存储
func (this *ChromemManager) addDocuments(ragId int, index int, content string) error {
	ctx := context.Background()
	collection := this.collectionMap[ragId]
	return collection.AddDocument(ctx, chromem.Document{ID: strconv.Itoa(index), Content: content})
}

// 向量检索
func (this *ChromemManager) query(ragId int, text string, nResults int) ([]int, error) {
	ctx := context.Background()
	collection := this.collectionMap[ragId]
	res, err := collection.Query(ctx, text, nResults, nil, nil)
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
