package agent

import (
	"go-ollama/logger"
	"go-ollama/ollama"
	"sync"
)

type Agent interface {
	DefaultModelName() string
	SendChat(chat string) string
}

var agentMap map[string]Agent
var muLock sync.Mutex

func registerAgent(name string, agent Agent) {
	muLock.Lock()
	defer muLock.Unlock()

	if agentMap == nil {
		agentMap = make(map[string]Agent)
	}
	agentMap[name] = agent
}

func StartAllAgent(ollama *ollama.OllamaManager, logger *logger.ErrorLogger) error {
	return nil
}

func GetAgent(name string) Agent {
	return agentMap[name]
}
