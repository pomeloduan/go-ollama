package agent

import (
	"go-ollama/logger"
	"go-ollama/ollama"
	"sync"
)

type Agent interface {
	DefaultModelName() string
	Chat(chat string) string
}

type AgentManager struct {
	ollama          *ollama.OllamaManager
	coordinateAgent *CoordinateAgent
	rankAgent       *RankAgent
	specialistMap   map[string]Agent
	muLock          sync.Mutex
}

func StartAgentManager(ollama *ollama.OllamaManager, logger *logger.ErrorLogger) (*AgentManager, error) {
	agentManager := AgentManager{
		ollama:          ollama,
		coordinateAgent: StartCoordinateAgent(ollama),
		rankAgent:       StartRankAgent(ollama),
		specialistMap:   make(map[string]Agent),
	}
	return &agentManager, nil
}

func (this AgentManager) registerSpecialist(name string, specialist Agent) {
	this.specialistMap[name] = specialist
}

func (this AgentManager) StartAllSpecialist(ollama *ollama.OllamaManager, logger *logger.ErrorLogger) error {
	return nil
}

func (this AgentManager) GetSpecialist(name string) Agent {
	return this.specialistMap[name]
}

func (this AgentManager) GetCoordinateAgent() *CoordinateAgent {
	return this.coordinateAgent
}

func (this AgentManager) GetRankAgent() *RankAgent {
	return this.rankAgent
}
