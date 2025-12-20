package rag

import (
	"go-ollama/logger"
	"sort"
)

type RagManager struct {
	chromem *ChromemManager
	gse     *GseManager
	logger  *logger.ErrorLogger
	chucks  []string
}

func StartRag(logger *logger.ErrorLogger) (*RagManager, error) {
	chromem, err := StartChromem()
	if err != nil {
		return nil, err
	}
	gse := StartGse()

	ragManager := RagManager{chromem: chromem, gse: gse, logger: logger}
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
func (this *RagManager) Query(text string) ([]string, error) {
	indexArr, err := this.chromem.Query(text, 10)
	if err != nil {
		return nil, err
	}
	sort.Ints(indexArr)

	var resArr []string
	for i := 0; i < len(indexArr); i++ {
		index := indexArr[i]
		if i > 0 && index-indexArr[i-1] == 2 {
			resArr = append(resArr, this.chucks[index-1])
		}
		resArr = append(resArr, this.chucks[index])
	}
	return resArr, nil
}
