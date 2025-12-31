package rag

import (
	"context"
	"strconv"

	"github.com/philippgille/chromem-go"
)

// ChromemManager 向量数据库管理器
// 使用 Chromem 库管理向量集合，支持文档的向量化和相似度检索
type ChromemManager struct {
	db            *chromem.DB                      // Chromem 数据库实例
	collectionMap map[int]*chromem.Collection      // RAG ID 到向量集合的映射
}

// ollamaEmbedModelName Ollama 嵌入模型名称，用于文档向量化
const ollamaEmbedModelName = "nomic-embed-text-v2-moe"

// startChromem 创建并初始化向量数据库管理器
func startChromem() *ChromemManager {
	return &ChromemManager{db: chromem.NewDB(), collectionMap: make(map[int]*chromem.Collection)}
}

// newCollection 为指定的 RAG 上下文创建新的向量集合
// 参数 ragId: RAG 上下文 ID
// 返回: error
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

// addDocuments 添加文档到向量集合
// 会自动调用 Ollama 进行向量化，然后存储到向量数据库
// 参数 ragId: RAG 上下文 ID
// 参数 index: 文档索引（用作文档 ID）
// 参数 content: 文档内容
// 返回: error
func (this *ChromemManager) addDocuments(ragId int, index int, content string) error {
	ctx := context.Background()
	collection := this.collectionMap[ragId]
	return collection.AddDocument(ctx, chromem.Document{ID: strconv.Itoa(index), Content: content})
}

// query 向量相似度检索
// 将查询文本向量化，然后检索最相似的文档
// 参数 ragId: RAG 上下文 ID
// 参数 text: 查询文本
// 参数 nResults: 返回的文档数量
// 返回: 文档索引数组（按相似度排序）、error
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
