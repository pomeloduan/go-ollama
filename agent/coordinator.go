package agent

import (
	"go-ollama/ollama"
	"go-ollama/rule"
)

// Coordinator 协调者，负责分析问题并选择最合适的专家 Agent
type Coordinator struct {
	ollama        ollama.OllamaManager // Ollama 管理器，用于调用 LLM
	modelName     string               // 使用的模型名称
	specialistMap map[string]string    // 专家名称到介绍的映射
	rule          rule.RuleManager     // 规则管理器
}

// newCoordinator 创建并初始化协调者实例
func newCoordinator(ollama ollama.OllamaManager, rule rule.RuleManager) *Coordinator {
	coordinator := Coordinator{
		ollama:        ollama,
		modelName:     ollama.GetAvailableModelName("deepseek"),
		specialistMap: make(map[string]string),
		rule:          rule,
	}
	return &coordinator
}

// addSpecialist 注册专家到协调者的专家列表中
// 参数 name: 专家名称
// 参数 introduction: 专家介绍，用于匹配问题
func (this *Coordinator) addSpecialist(name string, introduction string) {
	this.specialistMap[name] = introduction
}

// askForSpecialistName 分析用户问题，选择最合适的专家来回答
// 使用 LLM 根据专家介绍和问题内容进行匹配
// 参数 chat: 用户输入的问题
// 返回: 匹配的专家名称，如果没有匹配则返回空字符串
func (this *Coordinator) askForSpecialistName(chat string) string {
	message := this.rule.CoordinatorMessage(chat)
	for name, introduction := range this.specialistMap {
		message += this.rule.CoordinatorSpecialistMessage(name, introduction)
	}
	return this.ollama.ChatWithoutContext(this.modelName, message)
}
