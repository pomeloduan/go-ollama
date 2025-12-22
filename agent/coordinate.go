package agent

import "go-ollama/ollama"

type CoordinateAgent struct {
	ollama    *ollama.OllamaManager
	modelName string
}

func StartCoordinateAgent(ollama *ollama.OllamaManager) *CoordinateAgent {
	coordinateAgent := CoordinateAgent{ollama: ollama}
	coordinateAgent.modelName = ollama.GetAvailableModelName(coordinateAgent.DefaultModelName())
	return &coordinateAgent
}

func (this CoordinateAgent) DefaultModelName() string {
	return "deepseek"
}

func (this CoordinateAgent) Chat(chat string) string {
	return ""
}
