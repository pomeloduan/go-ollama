package ollama

type ChatContext struct {
	modelName     string
	chatId        int
	systemMessage ChatMessage
	history       []ChatMessage
}

func newChat(modelName string, chatId int, systemMesssage string) *ChatContext {
	return &ChatContext{modelName: modelName, chatId: chatId, systemMessage: ChatMessage{Role: "system", Content: systemMesssage}}
}

func (this ChatContext) addMessage(messsage ChatMessage) {
	this.history = append(this.history, messsage)
}

func (this ChatContext) addChatString(content string) {
	this.history = append(this.history, ChatMessage{Role: "user", Content: content})
}

func (this ChatContext) getMessages() []ChatMessage {
	messages := make([]ChatMessage, 1)
	messages[0] = this.systemMessage
	messages = append(messages, this.history...)
	return messages
}
