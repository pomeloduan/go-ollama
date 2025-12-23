package rag

import (
	"sort"
	"strings"
)

type RagManager struct {
	chromem *ChromemManager
	gse     *GseManager
	rerank  Rerankable

	autogenRagId int
}

type Rerankable interface {
	RankCandidate(candidates string, text string, num int) string
}

const retrievalCount = 10
const rerankCount = 5

func StartRag(rerank Rerankable) *RagManager {
	chromem := StartChromem()
	gse := StartGse()
	ragManager := RagManager{chromem: chromem, gse: gse, rerank: rerank}
	return &ragManager
}

// 预处理
// 1 文本分块
// 2 向量化 存储
// 检索
// 3.1 召回
// 3.2 重排

type ProgressInfo struct {
	Current    int
	Total      int
	Percentage float32
	Err        error
	Text       string
}

// 预处理
func (this *RagManager) PreprocessFromFile(filepath string) (*RagContext, chan ProgressInfo, error) {
	var ragId = this.autogenRagId
	this.autogenRagId++

	chucks, err := chucksFromTextFile(filepath)
	if err != nil {
		return nil, nil, err
	}

	err = this.chromem.newCollection(ragId)
	if err != nil {
		return nil, nil, err
	}

	chProg := make(chan ProgressInfo)
	go func() {
		defer close(chProg)
		for i := 0; i < len(chucks); i++ {
			words := this.gse.SplitChineseWords(chucks[i])
			err = this.chromem.addDocuments(ragId, i, words)

			percentage := float32(i+1) / float32(len(chucks)) * 100
			// 发送进度信息
			chProg <- ProgressInfo{
				Current:    i + 1,
				Total:      len(chucks),
				Percentage: percentage,
				Err:        err,
				Text:       chucks[i],
			}
		}
	}()

	ragCtx := RagContext{ragId: ragId, chucks: chucks}
	return &ragCtx, chProg, nil
}

// 检索
func (this *RagManager) Query(ragCtx *RagContext, text string) (chan string, error) {
	// 召回 向量相似
	indexArr, err := this.chromem.query(ragCtx.ragId, text, retrievalCount)
	if err != nil {
		return nil, err
	}
	sort.Ints(indexArr)

	var textArr []string
	for i := 0; i < len(indexArr); i++ {
		index := indexArr[i]
		if i > 0 && index-indexArr[i-1] == 2 {
			textArr = append(textArr, ragCtx.chucks[index-1])
		}
		textArr = append(textArr, ragCtx.chucks[index])
	}
	// 重排 rank agent
	chRes := make(chan string)
	go func() {
		defer close(chRes)
		chRes <- this.rerank.RankCandidate(strings.Join(textArr, "\n"), text, rerankCount)
	}()
	return chRes, nil
}
