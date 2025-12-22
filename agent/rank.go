package agent

import (
	"go-ollama/ollama"
)

type RankAgent struct {
	ollama    *ollama.OllamaManager
	modelName string
}

func (this RankAgent) DefaultModelName() string {
	return "gemma"
}

func Start(ollama *ollama.OllamaManager) *RankAgent {
	rankAgent := RankAgent{ollama: ollama}
	rankAgent.modelName = ollama.GetAvailableModelName(rankAgent.DefaultModelName())
	return &rankAgent
}

func (this RankAgent) RankCandidate(candidates string, text string) string {
	chat := "话题：" + text + "请从以下许多段文字中，先每一段都和话题进行比较，给出一个相关性评分，然后选择相关性最高的3段，最后仅回复相关性最高的3段文字原文，不需要回复原因和分数：" + candidates
	return this.ollama.ChatWithoutContext(this.modelName, chat)
}
