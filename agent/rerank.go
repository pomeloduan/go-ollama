package agent

import (
	"go-ollama/ollama"
	"go-ollama/rule"
)

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

func (this *Reranker) RankCandidate(candidates string, text string, num int) string {
	message := this.rule.RerankMessage(text, num, candidates)
	return this.ollama.ChatWithoutContext(this.modelName, message)
}
