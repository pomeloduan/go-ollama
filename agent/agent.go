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

var coordinateAgent Agent
var agentMap map[string]Agent
var muLock sync.Mutex
var once sync.Once

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

func GetCoordinateAgent() Agent {
	once.Do(func() {
		coordinateAgent = &CoordinateAgent{}
	})
	return coordinateAgent
}
