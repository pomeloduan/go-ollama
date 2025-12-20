package ollama

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"go-ollama/logger"
)

type OllamaManager struct {
	domain    string
	modelName string
	logger    *logger.ErrorLogger

	contextMap map[int]*ChatContext

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

func StartOllama(domain, defaultModel string, logger *logger.ErrorLogger) (*OllamaManager, error) {
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

	ollamaManager := OllamaManager{domain: domain, modelName: modelName, logger: logger}
	ollamaManager.contextMap = make(map[int]*ChatContext)

	return &ollamaManager, nil
}

func (this *OllamaManager) GetModelName() string {
	return this.modelName
}

// 新的对话
func (this *OllamaManager) NewChat(systemMesssage string) int {
	var chatId = this.chatId
	this.chatId++

	// 初始化历史记录
	history := make([]ChatMessage, 1)
	history[0] = ChatMessage{Role: "system", Content: systemMesssage}
	chatContext := ChatContext{chatId: chatId, history: history}
	this.contextMap[chatId] = &chatContext
	return chatId
}

// 添加系统信息
func (this *OllamaManager) AddSystemMessage(chatId int, systemMesssage string) {

	chatContext, ok := this.contextMap[chatId]
	if !ok {
		this.logger.LogError(fmt.Errorf("no chat Id found"), "add system")
		return
	}
	chatContext.history = append(chatContext.history, ChatMessage{Role: "system", Content: systemMesssage})
	this.logger.LogInfo("q" + strconv.Itoa(chatId) + ": " + systemMesssage)
}

// 对话
func (this *OllamaManager) NextChat(chatId int, message string) string {
	this.logger.LogInfo("q" + strconv.Itoa(chatId) + ": " + message)
	this.TotalQCount++

	// 问题+历史记录
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

	// 统计
	elapsed := time.Since(start)

	this.TotalACount++
	this.TotalDuration += elapsed
	this.TotalToken += response.EvalCount

	respMessage := response.Message

	this.logger.LogInfo("a" + strconv.Itoa(chatId) + ": " + respMessage.Content)

	// 保存历史记录
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
