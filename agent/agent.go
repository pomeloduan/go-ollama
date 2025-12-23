package agent

import (
	"go-ollama/agent/rag"
	"go-ollama/agent/rule"
	"go-ollama/logger"
	"go-ollama/ollama"
	"sync"
)

type AgentManager struct {
	ollama        *ollama.OllamaManager
	rag           *rag.RagManager
	coordinator   *Coordinator
	generalAgent  *Specialist
	specialistMap map[string]*Specialist
	muLock        sync.Mutex
	logger        *logger.ErrorLogger
}

func StartAgentManager(ollama *ollama.OllamaManager, logger *logger.ErrorLogger) (*AgentManager, error) {
	rank := startReranker(ollama)
	rag := rag.StartRag(rank)

	coordinator := startCoordinator(ollama)

	general := startSpecialist(ollama, rag, &rule.GeneralRule{}, logger)
	specialistMap := make(map[string]*Specialist)
	for _, n := range rule.AllRuleNames() {
		rule := rule.GetRule(n)
		specialist := startSpecialist(ollama, rag, rule, logger)
		specialistMap[n] = specialist
		coordinator.addSpecialist(n, rule.Introduction())
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
	specialistName := this.coordinator.askForSpecialist(chat)
	specialist, ok := this.specialistMap[specialistName]
	if !ok {
		specialist = this.generalAgent
	}
	return specialist.chat(chat)
}
