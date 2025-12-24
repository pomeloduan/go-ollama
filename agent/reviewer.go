package agent

import (
	"go-ollama/logger"
	"go-ollama/ollama"
	"go-ollama/rule"
)

type Reviewer struct {
	ollama    *ollama.OllamaManager
	modelName string
	rule      *rule.Rule
	chatCtx   *ollama.ChatContext
	logger    *logger.ErrorLogger
}

func startReviewer(ollama *ollama.OllamaManager, rule *rule.Rule, logger *logger.ErrorLogger) *Reviewer {
	reviewer := Reviewer{
		ollama:    ollama,
		modelName: ollama.GetAvailableModelName("gemma"),
		rule:      rule,
		logger:    logger,
	}
	return &reviewer
}

func (this *Reviewer) prepareChat() {
	this.chatCtx = this.ollama.NewChat(this.modelName, this.rule.ReviewerSystemMessage())
}

func (this *Reviewer) review(question string, answer string) rule.ReviewResult {
	if this.chatCtx == nil {
		this.prepareChat()
	}
	message := "要求：\n" + question + "\n作品：\n" + answer
	review := this.ollama.NextChat(this.chatCtx, message)
	return this.rule.ParseReview(review)
}
