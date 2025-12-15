package ollama

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"go-ollama/logger"
)

type OllamaClient struct {
	domain    string
	modelName string
	logger    *logger.ErrorLogger

	contextMap map[int]ChatContext

	chatId int

	// 数据统计
	TotalQCount   int
	TotalACount   int
	TotalDuration time.Duration
	TotalToken    int
}

type ChatContext struct {
	chatId  int
	history []ChatMessage
}

func StartOllamaClient(domain, defaultModel string, logger *logger.ErrorLogger) (*OllamaClient, error) {
	// 1. 检查 Ollama 服务是否运行
	_, err := http.Get(domain)
	if err != nil {
		return nil, fmt.Errorf("need ollama server")
	}

	// 2. 列出可用模型
	models, err := listModels(domain)
	if err != nil || len(models) == 0 {
		return nil, fmt.Errorf("no model")
	}

	// 3. 使用deepseek，或第一个模型
	modelName, getModel := findFirstContaining(models, defaultModel)
	if !getModel {
		modelName = models[0]
	}

	ollamaClient := OllamaClient{domain: domain, modelName: modelName, logger: logger}
	ollamaClient.contextMap = make(map[int]ChatContext)

	return &ollamaClient, nil
}

func (this *OllamaClient) NewChat(systemMesssage string) int {
	var chatId = this.chatId
	this.chatId++

	history := make([]ChatMessage, 1)
	history[0] = ChatMessage{Role: "system", Content: systemMesssage}
	chatContext := ChatContext{chatId: chatId, history: history}
	this.contextMap[chatId] = chatContext
	return chatId
}

func (this *OllamaClient) NextChat(message string, chatId int) string {
	this.logger.LogInfo("q: " + message)
	this.TotalQCount++

	chatContext, ok := this.contextMap[chatId]
	if !ok {
		this.logger.LogError(fmt.Errorf("no chat Id found"), "sendchat")
		return ""
	}

	allChat := append(chatContext.history, ChatMessage{Role: "user", Content: message})

	start := time.Now()
	response, err := sendChatRequest(this.domain, this.modelName, allChat)
	if err != nil {
		this.logger.LogError(fmt.Errorf("send chat err: %v", err), "sendchat")
		return ""
	}

	elapsed := time.Since(start)

	this.TotalACount++
	this.TotalDuration += elapsed
	this.TotalToken += response.EvalCount

	respMessage := response.Message

	this.logger.LogInfo("a: " + respMessage.Content)

	chatContext.history = append(allChat, respMessage)

	return respMessage.Content
}

// 查找第一个包含子字符串的元素
func findFirstContaining(arr []string, substr string) (string, bool) {
	for _, str := range arr {
		if strings.Contains(strings.ToLower(str), strings.ToLower(substr)) {
			return str, true
		}
	}
	return "", false
}
