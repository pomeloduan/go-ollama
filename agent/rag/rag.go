package rag

import (
	"go-ollama/logger"
	"sort"
	"strings"
)

type RagManager struct {
	chromem *ChromemManager
	gse     *GseManager
	logger  *logger.ErrorLogger
	chucks  []string
	rank    Rankable
}

type Rankable interface {
	RankCandidate(candidates string, text string, num int) string
}

const retrievalCount = 10
const rerankCount = 5

// todo rag context
func StartRag(rank Rankable, logger *logger.ErrorLogger) (*RagManager, error) {
	chromem, err := StartChromem()
	if err != nil {
		return nil, err
	}
	gse := StartGse()

	ragManager := RagManager{chromem: chromem, gse: gse, rank: rank, logger: logger}
	return &ragManager, nil
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
func (this *RagManager) PreprocessFromFile(filepath string) (chan ProgressInfo, error) {
	chucks, err := chucksFromTextFile(filepath)
	if err != nil {
		return nil, err
	}
	chProg := make(chan ProgressInfo)
	go func() {
		defer close(chProg)
		for i := 0; i < len(chucks); i++ {
			words := this.gse.SplitChineseWords(chucks[i])
			err = this.chromem.AddDocuments(i, words)

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
	this.chucks = chucks
	return chProg, nil
}

// 检索
func (this *RagManager) Query(text string) (chan string, error) {
	// 召回 向量相似
	indexArr, err := this.chromem.Query(text, retrievalCount)
	if err != nil {
		return nil, err
	}
	sort.Ints(indexArr)

	var textArr []string
	for i := 0; i < len(indexArr); i++ {
		index := indexArr[i]
		if i > 0 && index-indexArr[i-1] == 2 {
			textArr = append(textArr, this.chucks[index-1])
		}
		textArr = append(textArr, this.chucks[index])
	}
	// 重排 rank agent
	chRes := make(chan string)
	go func() {
		defer close(chRes)
		chRes <- this.rank.RankCandidate(strings.Join(textArr, "\n"), text, rerankCount)
	}()
	return chRes, nil
}
