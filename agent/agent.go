package agent

import (
	"go-ollama/logger"
	"go-ollama/ollama"
	"go-ollama/rag"
	"go-ollama/rule"
)

type AgentManager struct {
	ollama        *ollama.OllamaManager
	rag           *rag.RagManager
	rule          *rule.RuleManager
	coordinator   *Coordinator
	generalAgent  *Specialist
	specialistMap map[string]*Specialist
	reviewerMap   map[string]*Reviewer
	logger        *logger.ErrorLogger
}

func StartAgentManager(ollama *ollama.OllamaManager, logger *logger.ErrorLogger) (*AgentManager, error) {
	rank := startReranker(ollama)
	rag := rag.StartRag(rank)
	ruleManager := rule.StartRuleManager()

	coordinator := startCoordinator(ollama)

	general := startSpecialist(ollama, rag, ruleManager.GetGeneralRule(), logger)
	specialistMap := make(map[string]*Specialist)
	reviewerMap := make(map[string]*Reviewer)
	for _, rule := range ruleManager.GetAllRules() {
		specialist := startSpecialist(ollama, rag, rule, logger)
		specialistMap[rule.Name()] = specialist
		coordinator.addSpecialist(rule.Name(), rule.Introduction())
		if rule.NeedReviewer() {
			reviewer := startReviewer(ollama, rule, logger)
			reviewerMap[rule.Name()] = reviewer
		}
	}

	agentManager := AgentManager{
		ollama:        ollama,
		rag:           rag,
		coordinator:   coordinator,
		generalAgent:  general,
		specialistMap: specialistMap,
		logger:        logger,
	}
	return &agentManager, nil
}

func (this *AgentManager) Chat(chat string) string {
	// call coordinator
	name := this.coordinator.askForSpecialist(chat)
	specialist, ok := this.specialistMap[name]
	if !ok {
		specialist = this.generalAgent
	}
	answer := specialist.chat(chat)
	reviewer, ok := this.reviewerMap[name]
	if ok {
		review := reviewer.review(chat, answer)
		if review.Score < 80 {
			answer = specialist.chat("请参考以下评价重新写作：\n" + review.Review)
		}
	}
	return answer
}
