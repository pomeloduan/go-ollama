package agent

import (
	"go-ollama/ollama"
	"strconv"
)

type Reranker struct {
	ollama    *ollama.OllamaManager
	modelName string
}

func startReranker(ollama *ollama.OllamaManager) *Reranker {
	reranker := Reranker{ollama: ollama, modelName: ollama.GetAvailableModelName("gemma")}
	return &reranker
}

func (this *Reranker) RankCandidate(candidates string, text string, num int) string {
	message := "话题：" + text +
		"\n请从以下许多段文字中，先每一段都和话题进行比较，给出一个相关性评分，然后选择相关性最高的" + strconv.Itoa(num) +
		"段，最后仅回复相关性最高的" + strconv.Itoa(num) +
		"段文字原文，不需要回复原因和分数：\n" + candidates
	return this.ollama.ChatWithoutContext(this.modelName, message)
}
