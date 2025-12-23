package agent

import "go-ollama/ollama"

type CoordinateAgent struct {
	ollama        *ollama.OllamaManager
	modelName     string
	specialistMap map[string]string
}

func StartCoordinateAgent(ollama *ollama.OllamaManager) *CoordinateAgent {
	coordinate := CoordinateAgent{
		ollama:        ollama,
		modelName:     ollama.GetAvailableModelName("deepseek"),
		specialistMap: make(map[string]string),
	}
	return &coordinate
}

func (this *CoordinateAgent) AddSpecialist(name string, introduction string) {
	this.specialistMap[name] = introduction
}
func (this *CoordinateAgent) AskForSpecialist(chat string) string {
	message := "现在有一个问题需要寻求专家的帮助，有下面几位专家：\n"
	for name, introduction := range this.specialistMap {
		message += "专家名字：" + name + " 专家介绍：" + introduction + "\n"
	}
	message += "问题是：" + chat + "\n请选择与问题相关的适合解答问题的专家，回复专家名字，或者你认为没有专家能够解答，回复NA"
	return this.ollama.ChatWithoutContext(this.modelName, message)
}
