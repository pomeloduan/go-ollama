package agent

import (
	"go-ollama/ollama"
	"go-ollama/rule"
)

type Coordinator struct {
	ollama        *ollama.OllamaManager
	modelName     string
	specialistMap map[string]string
	rule          *rule.RuleManager
}

func startCoordinator(ollama *ollama.OllamaManager, rule *rule.RuleManager) *Coordinator {
	coordinator := Coordinator{
		ollama:        ollama,
		modelName:     ollama.GetAvailableModelName("deepseek"),
		specialistMap: make(map[string]string),
		rule:          rule,
	}
	return &coordinator
}

func (this *Coordinator) addSpecialist(name string, introduction string) {
	this.specialistMap[name] = introduction
}

func (this *Coordinator) askForSpecialist(chat string) string {
	message := this.rule.CoordinatorMessage(chat)
	for name, introduction := range this.specialistMap {
		message += this.rule.CoordinatorSpecialistMessage(name, introduction)
	}
	return this.ollama.ChatWithoutContext(this.modelName, message)
}
