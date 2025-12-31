package rag

import (
	"go-ollama/rule"
	"sort"
	"strings"
)

// RagManager RAG 管理器，负责检索增强生成的完整流程
// 包括文本预处理、向量检索和结果重排序
type RagManager struct {
	chromem  *ChromemManager // 向量数据库管理器
	gse      *GseManager     // 中文分词管理器
	reranker Rerankable      // 重排序器接口

	autogenRagId int // 自动生成的 RAG 上下文 ID
}

// Rerankable 重排序器接口，用于对检索结果进行重排序
type Rerankable interface {
	RankCandidate(candidates string, text string, num int) string
}

// retrievalCount 向量检索返回的候选文档数量
const retrievalCount = 10
// rerankingCount 重排序后返回的最终文档数量
const rerankingCount = 5

func StartRag(reranker Rerankable) *RagManager {
	chromem := startChromem()
	gse := startGse()
	ragManager := RagManager{chromem: chromem, gse: gse, reranker: reranker}
	return &ragManager
}

// RAG 流程说明
// 预处理阶段：
// 1. 文本分块（chunking）：将文档按段落分割成多个块
// 2. 向量化存储（embedding）：对每个文本块进行向量化并存储到向量数据库
// 检索阶段：
// 3.1 召回（retrieval）：通过向量相似度检索相关文档
// 3.2 重排（reranking）：使用 LLM 对检索结果进行相关性重排序

// ProgressInfo 预处理进度信息，通过 channel 实时返回
type ProgressInfo struct {
	Current    int     // 当前处理的项数
	Total      int     // 总项数
	Percentage float32 // 完成百分比
	Err        error   // 错误信息
	Text       string  // 当前处理的文本内容
}

// PreprocessFromFile 从文件预处理知识库
// 包括文本分块、中文分词、向量化和存储
// 参数 filepath: 源文件路径
// 返回: RagContext、ProgressInfo channel、error
// 注意：ProgressInfo channel 需要调用者消费，否则会导致 goroutine 阻塞
func (this *RagManager) PreprocessFromFile(filepath string) (*RagContext, chan ProgressInfo, error) {
	var ragId = this.autogenRagId
	this.autogenRagId++

	chunks, err := chunksFromTextFile(filepath)
	if err != nil {
		return nil, nil, err
	}

	err = this.chromem.newCollection(ragId)
	if err != nil {
		return nil, nil, err
	}

	// 创建进度 channel，使用 goroutine 异步处理
	chProg := make(chan ProgressInfo)
	go func() {
		defer close(chProg)
		for i := 0; i < len(chunks); i++ {
			// 对中文文本进行分词，提升向量化效果
			words := this.gse.splitChineseWords(chunks[i])
			// 将文档添加到向量数据库（自动进行向量化）
			err = this.chromem.addDocuments(ragId, i, words)

			percentage := float32(i+1) / float32(len(chunks)) * 100
			// 发送进度信息
			chProg <- ProgressInfo{
				Current:    i + 1,
				Total:      len(chunks),
				Percentage: percentage,
				Err:        err,
				Text:       chunks[i],
			}
		}
	}()

	ragCtx := RagContext{ragId: ragId, chunks: chunks}
	return &ragCtx, chProg, nil
}

// Query 检索与问题相关的文档
// 流程：1. 向量相似度召回 2. 相邻块合并 3. LLM 重排序
// 参数 ragCtx: RAG 上下文，包含知识库信息
// 参数 text: 用户问题
// 参数 rule: 规则配置（当前未使用，保留用于扩展）
// 返回: string channel、error
// 注意：返回的 channel 需要调用者消费
func (this *RagManager) Query(ragCtx *RagContext, text string, rule *rule.Rule) (chan string, error) {
	// 1. 向量相似度召回：检索最相似的文档块索引
	indexArr, err := this.chromem.query(ragCtx.ragId, text, retrievalCount)
	if err != nil {
		return nil, err
	}
	// 对索引排序，便于后续处理
	sort.Ints(indexArr)

	// 2. 合并相邻的文档块，保持上下文连贯性
	var textArr []string
	for i := 0; i < len(indexArr); i++ {
		index := indexArr[i]
		// 如果当前块和前一个块只间隔一个块，将中间块也加入，保证上下文完整
		if i > 0 && index-indexArr[i-1] == 2 {
			textArr = append(textArr, ragCtx.chunks[index-1])
		}
		textArr = append(textArr, ragCtx.chunks[index])
	}
	
	// 3. 使用 LLM 对候选文档进行重排，选择最相关的文档
	chRes := make(chan string)
	go func() {
		defer close(chRes)
		chRes <- this.reranker.RankCandidate(strings.Join(textArr, "\n"), text, rerankingCount)
	}()
	return chRes, nil
}
