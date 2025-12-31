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
	domain string
	models []string
	logger *logger.ErrorLogger

	autogenChatId int

	// 数据统计
	TotalQCount   int
	TotalACount   int
	TotalDuration time.Duration
	TotalToken    int
}

func StartOllama(domain string, logger *logger.ErrorLogger) (*OllamaManager, error) {
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

	ollamaManager := OllamaManager{domain: domain, models: models, logger: logger}

	return &ollamaManager, nil
}

func (this *OllamaManager) GetAvailableModelName(modelName string) string {
	for i := 0; i < len(this.models); i++ {
		if strings.Contains(this.models[i], modelName) {
			return this.models[i]
		}
	}
	return ""
}

func (this *OllamaManager) GetDefaultEmbedModelName() string {
	return this.GetAvailableModelName("embed")
}

func (this *OllamaManager) GetDefaultLlmModelName() string {
	return this.GetAvailableModelName("deepseek")
}

// 无需上下文的对话（单次对话）
func (this *OllamaManager) ChatWithoutContext(modelName string, message string) string {
	this.logger.LogInfo("q#: " + message)
	this.TotalQCount++

	start := time.Now()
	response, err := sendChatRequest(this.domain, modelName, chatMessagesFromChatString(message))
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

	this.logger.LogInfo("a#: " + respMessage.Content)

	return respMessage.Content
}

// 新的对话
func (this *OllamaManager) NewChat(modelName string, systemMesssage string) *ChatContext {
	var chatId = this.autogenChatId
	this.autogenChatId++

	return newChat(modelName, chatId, systemMesssage)
}

// 对话
// todo 上下文优化 有限上下文
func (this *OllamaManager) NextChat(chatCtx *ChatContext, message string) string {
	this.logger.LogInfo("q" + strconv.Itoa(chatCtx.chatId) + ": " + message)
	this.TotalQCount++

	// 问题+历史记录
	chatCtx.addChatString(message)

	start := time.Now()
	response, err := sendChatRequest(this.domain, chatCtx.modelName, chatCtx.getMessages())
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

	this.logger.LogInfo("a" + strconv.Itoa(chatCtx.chatId) + ": " + respMessage.Content)

	// 保存历史记录
	chatCtx.addMessage(respMessage)

	return respMessage.Content
}
