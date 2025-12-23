package agent

import "go-ollama/ollama"

type Coordinator struct {
	ollama        *ollama.OllamaManager
	modelName     string
	specialistMap map[string]string
}

func startCoordinator(ollama *ollama.OllamaManager) *Coordinator {
	coordinator := Coordinator{
		ollama:        ollama,
		modelName:     ollama.GetAvailableModelName("gemma"),
		specialistMap: make(map[string]string),
	}
	return &coordinator
}

func (this *Coordinator) addSpecialist(name string, introduction string) {
	this.specialistMap[name] = introduction
}

func (this *Coordinator) askForSpecialist(chat string) string {
	message := "现在有一个问题需要寻求专家的帮助，有下面几位专家：\n"
	for name, introduction := range this.specialistMap {
		message += "专家名字：" + name + " 专家介绍：" + introduction + "\n"
	}
	message += "问题是：" + chat + "\n请选择与问题相关的适合解答问题的专家，回复专家名字，或者你认为没有专家能够解答，回复NA"
	return this.ollama.ChatWithoutContext(this.modelName, message)
}
