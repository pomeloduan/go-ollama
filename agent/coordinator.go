package agent

import (
	"go-ollama/ollama"
	"go-ollama/rule"
)

// 协调者
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

// 添加专家的名称/介绍
func (this *Coordinator) addSpecialist(name string, introduction string) {
	this.specialistMap[name] = introduction
}

// 协调者选择主题相关的专家回答问题
func (this *Coordinator) askForSpecialistName(chat string) string {
	message := this.rule.CoordinatorMessage(chat)
	for name, introduction := range this.specialistMap {
		message += this.rule.CoordinatorSpecialistMessage(name, introduction)
	}
	return this.ollama.ChatWithoutContext(this.modelName, message)
}
