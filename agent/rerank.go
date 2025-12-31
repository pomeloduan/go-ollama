package agent

import (
	"go-ollama/ollama"
	"go-ollama/rule"
)

// 重排
type Reranker struct {
	ollama    *ollama.OllamaManager
	modelName string
	rule      *rule.RuleManager
}

func startReranker(ollama *ollama.OllamaManager, rule *rule.RuleManager) *Reranker {
	reranker := Reranker{
		ollama:    ollama,
		modelName: ollama.GetAvailableModelName("gemma"),
		rule:      rule,
	}
	return &reranker
}

// 调用LLM对比候选资料和特定文本
func (this *Reranker) RankCandidate(candidates string, text string, num int) string {
	message := this.rule.RerankMessage(candidates, text, num)
	return this.ollama.ChatWithoutContext(this.modelName, message)
}
