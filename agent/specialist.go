package agent

import "go-ollama/ollama"

type SpecialistAgent struct {
	ollama    *ollama.OllamaManager
	modelName string
}
