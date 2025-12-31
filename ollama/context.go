package ollama

// ChatContext 对话上下文，维护多轮对话的历史记录
type ChatContext struct {
	modelName     string        // 使用的模型名称
	chatId        int           // 对话 ID，用于区分不同的对话
	systemMessage ChatMessage   // 系统提示词
	history       []ChatMessage // 对话历史（用户消息和助手回答）
}

// newChat 创建新的对话上下文
// 参数 modelName: 模型名称
// 参数 chatId: 对话 ID
// 参数 systemMessage: 系统提示词
// 返回: ChatContext 实例
func newChat(modelName string, chatId int, systemMessage string) *ChatContext {
	return &ChatContext{modelName: modelName, chatId: chatId, systemMessage: ChatMessage{Role: "system", Content: systemMessage}}
}

// chatMessagesFromChatString 将字符串转换为单次对话的消息数组
// 用于无需上下文的对话场景
func chatMessagesFromChatString(content string) []ChatMessage {
	messages := make([]ChatMessage, 1)
	messages[0] = ChatMessage{Role: "user", Content: content}
	return messages
}

// addMessage 添加消息到对话历史
// 参数 message: 要添加的消息（通常是助手的回答）
func (c *ChatContext) addMessage(message ChatMessage) {
	c.history = append(c.history, message)
}

// addChatString 添加用户消息到对话历史
// 参数 content: 用户消息内容
func (c *ChatContext) addChatString(content string) {
	c.history = append(c.history, ChatMessage{Role: "user", Content: content})
}

// getMessages 获取完整的消息列表，用于发送给 LLM
// 格式：系统消息 + 对话历史
// 返回: 消息数组，第一个是系统消息，后面是对话历史
func (c *ChatContext) getMessages() []ChatMessage {
	messages := make([]ChatMessage, 1)
	messages[0] = c.systemMessage
	messages = append(messages, c.history...)
	return messages
}
