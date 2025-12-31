package agent

import (
	"go-ollama/logger"
	"go-ollama/ollama"
	"go-ollama/rag"
	"go-ollama/rule"
)

// agent 管理器
// 创建/管理所有agent生命周期
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

const rewriteScore = 80

// todo
// rag工程化
// function call
func StartAgentManager(ollama *ollama.OllamaManager, logger *logger.ErrorLogger) (*AgentManager, error) {
	// 1 rule
	// 规则管理器需要先起
	ruleManager, err := rule.StartRuleManager()
	if err != nil {
		return nil, err
	}

	// 2 rag
	reranker := startReranker(ollama, ruleManager)
	rag := rag.StartRag(reranker)

	// 3 coordinator
	coordinator := startCoordinator(ollama, ruleManager)

	// 4 specialist + reviewer
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
		reviewerMap:   reviewerMap,
		logger:        logger,
	}
	return &agentManager, nil
}

func (this *AgentManager) Chat(chat string) string {
	// call coordinator
	name := this.coordinator.askForSpecialistName(chat)
	specialist, ok := this.specialistMap[name]
	if !ok {
		specialist = this.generalAgent
	}

	// call specialist
	answer := specialist.chat(chat)
	// review
	reviewer, ok := this.reviewerMap[name]
	if ok {
		review := reviewer.review(chat, answer)
		if review.Score < rewriteScore {
			message := specialist.rule.RewriteMessage(review.Review)
			answer = specialist.chat(message)
		}
	}
	return answer
}
