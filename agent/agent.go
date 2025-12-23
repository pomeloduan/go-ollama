package agent

import (
	"go-ollama/agent/rag"
	"go-ollama/agent/rule"
	"go-ollama/logger"
	"go-ollama/ollama"
	"sync"
)

type AgentManager struct {
	ollama          *ollama.OllamaManager
	rag             *rag.RagManager
	coordinateAgent *CoordinateAgent
	generalAgent    *SpecialistAgent
	specialistMap   map[string]*SpecialistAgent
	muLock          sync.Mutex
	logger          *logger.ErrorLogger
}

func StartAgentManager(ollama *ollama.OllamaManager, logger *logger.ErrorLogger) (*AgentManager, error) {
	rank := StartRankAgent(ollama)
	rag, err := rag.StartRag(rank, logger)
	if err != nil {
		return nil, err
	}

	coordinate := StartCoordinateAgent(ollama)

	general := StartSpecialistAgent(ollama, rag, &rule.GeneralRule{}, logger)
	specialistMap := make(map[string]*SpecialistAgent)
	for _, n := range rule.AllRuleNames() {
		rule := rule.GetRule(n)
		specialist := StartSpecialistAgent(ollama, rag, rule, logger)
		specialistMap[n] = specialist
		coordinate.AddSpecialist(n, rule.Introduction())
	}

	agentManager := AgentManager{
		ollama:          ollama,
		rag:             rag,
		coordinateAgent: coordinate,
		generalAgent:    general,
		specialistMap:   specialistMap,
		logger:          logger,
	}
	return &agentManager, nil
}

func (this *AgentManager) Chat(chat string) string {
	// call coordinate agent
	specialistName := this.coordinateAgent.AskForSpecialist(chat)
	specialist, ok := this.specialistMap[specialistName]
	if !ok {
		specialist = this.generalAgent
	}
	return specialist.Chat(chat)
}

// func (this AgentManager) StartAllSpecialist() error {
// 	for _, n := range rule.AllRuleNames() {
// 		specialist := StartSpecialistAgent(this.ollama, this.rag, rule.GetRule(n), this.logger)
// 		this.specialistMap[n] = specialist
// 	}
// 	return nil
// }

// func (this AgentManager) GetSpecialist(name string) SpecialistAgent {
// 	return this.specialistMap[name]
// }

// func (this AgentManager) GetCoordinateAgent() *CoordinateAgent {
// 	return this.coordinateAgent
// }
