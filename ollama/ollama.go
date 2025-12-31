package ollama

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"go-ollama/logger"
)

// OllamaManager Ollama 服务管理器
// 负责与本地 Ollama 服务通信，管理模型和对话上下文
type OllamaManager struct {
	domain string            // Ollama 服务地址
	models []string          // 可用的模型列表
	logger *logger.ErrorLogger // 日志记录器

	autogenChatId int // 自动生成的对话 ID，用于区分不同的对话上下文

	// 数据统计
	TotalQCount   int           // 总问题数
	TotalACount   int           // 总回答数
	TotalDuration time.Duration // 总响应时间
	TotalToken    int           // 总 token 使用量
}

// StartOllamaManager 初始化 Ollama 管理器
// 检查服务是否运行，获取可用模型列表
// 参数 domain: Ollama 服务地址
// 参数 logger: 日志记录器
// 返回: OllamaManager 实例、error
func StartOllamaManager(domain string, logger *logger.ErrorLogger) (*OllamaManager, error) {
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

// GetAvailableModelName 获取包含指定关键词的可用模型名称
// 用于模糊匹配模型名称（如 "deepseek" 会匹配 "deepseek-chat"）
// 参数 modelName: 模型名称关键词
// 返回: 完整的模型名称，如果未找到则返回空字符串
func (this *OllamaManager) GetAvailableModelName(modelName string) string {
	for i := 0; i < len(this.models); i++ {
		if strings.Contains(this.models[i], modelName) {
			return this.models[i]
		}
	}
	return ""
}

// GetDefaultEmbedModelName 获取默认的嵌入模型名称
// 用于文档向量化
func (this *OllamaManager) GetDefaultEmbedModelName() string {
	return this.GetAvailableModelName("embed")
}

// GetDefaultLlmModelName 获取默认的 LLM 模型名称
// 用于文本生成
func (this *OllamaManager) GetDefaultLlmModelName() string {
	return this.GetAvailableModelName("deepseek")
}

// ChatWithoutContext 单次对话，不维护上下文
// 适用于不需要历史对话的场景（如协调者选择专家、重排序等）
// 参数 modelName: 模型名称
// 参数 message: 用户消息
// 返回: LLM 生成的回答
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

// NewChat 创建新的对话上下文
// 为每个 Agent 创建独立的对话上下文，用于维护多轮对话历史
// 参数 modelName: 模型名称
// 参数 systemMessage: 系统提示词
// 返回: ChatContext 实例
func (this *OllamaManager) NewChat(modelName string, systemMessage string) *ChatContext {
	var chatId = this.autogenChatId
	this.autogenChatId++

	return newChat(modelName, chatId, systemMessage)
}

// NextChat 继续进行对话，维护上下文
// 将新的消息添加到历史记录，调用 LLM 生成回答，并保存回答到历史
// 参数 chatCtx: 对话上下文
// 参数 message: 用户消息
// 返回: LLM 生成的回答
// todo 上下文优化：实现有限上下文窗口，避免历史记录过长导致 token 超限
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
