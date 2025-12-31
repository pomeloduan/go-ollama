package ollama

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"go-ollama/logger"
)

// OllamaManager Ollama 服务管理器接口
// 负责与本地 Ollama 服务通信，管理模型和对话上下文
type OllamaManager interface {
	GetAvailableModelName(modelName string) string
	GetDefaultEmbedModelName() string
	GetDefaultLlmModelName() string
	ChatWithoutContext(modelName string, message string) (string, error)
	NewChat(modelName string, systemMessage string) *ChatContext
	NextChat(chatCtx *ChatContext, message string) (string, error)
	// 统计信息
	GetTotalQCount() int
	GetTotalACount() int
	GetTotalDuration() time.Duration
	GetTotalToken() int
}

// ollamaManager Ollama 服务管理器实现（包私有）
type ollamaManager struct {
	domain string            // Ollama 服务地址
	models []string          // 可用的模型列表
	logger logger.ErrorLogger // 日志记录器

	mu            sync.RWMutex // 保护并发访问的读写锁
	autogenChatId int          // 自动生成的对话 ID，用于区分不同的对话上下文

	// 数据统计（需要并发保护）
	totalQCount   int           // 总问题数
	totalACount   int           // 总回答数
	totalDuration time.Duration // 总响应时间
	totalToken    int           // 总 token 使用量
}

var (
	ollamaInstance *ollamaManager
	ollamaOnce     sync.Once
)

// newOllamaManager 创建并初始化 Ollama 管理器实例
// 参数 domain: Ollama 服务地址
// 参数 logger: 日志记录器
// 返回: ollamaManager 实例、error
func newOllamaManager(domain string, logger logger.ErrorLogger) (*ollamaManager, error) {
	// 1. 检查 Ollama 服务是否运行
	resp, err := http.Get(domain)
	if err != nil {
		return nil, fmt.Errorf("need ollama server: %w", err)
	}
	defer resp.Body.Close()

	// 2. 列出可用模型
	models, err := listModels(domain)
	if err != nil || len(models) == 0 {
		return nil, fmt.Errorf("no model")
	}

	return &ollamaManager{
		domain: domain,
		models: models,
		logger: logger,
	}, nil
}

// StartOllamaManager 获取 Ollama 管理器单例
// 检查服务是否运行，获取可用模型列表
// 参数 domain: Ollama 服务地址
// 参数 logger: 日志记录器
// 返回: OllamaManager 实例、error
func StartOllamaManager(domain string, logger logger.ErrorLogger) (OllamaManager, error) {
	var err error
	ollamaOnce.Do(func() {
		ollamaInstance, err = newOllamaManager(domain, logger)
	})

	if err != nil {
		return nil, err
	}
	return ollamaInstance, nil
}

// GetAvailableModelName 获取包含指定关键词的可用模型名称
// 用于模糊匹配模型名称（如 "deepseek" 会匹配 "deepseek-chat"）
// 参数 modelName: 模型名称关键词
// 返回: 完整的模型名称，如果未找到则返回空字符串
func (o *ollamaManager) GetAvailableModelName(modelName string) string {
	for i := 0; i < len(o.models); i++ {
		if strings.Contains(o.models[i], modelName) {
			return o.models[i]
		}
	}
	return ""
}

// GetDefaultEmbedModelName 获取默认的嵌入模型名称
// 用于文档向量化
func (o *ollamaManager) GetDefaultEmbedModelName() string {
	return o.GetAvailableModelName("embed")
}

// GetDefaultLlmModelName 获取默认的 LLM 模型名称
// 用于文本生成
func (o *ollamaManager) GetDefaultLlmModelName() string {
	return o.GetAvailableModelName("deepseek")
}

// ChatWithoutContext 单次对话，不维护上下文
// 适用于不需要历史对话的场景（如协调者选择专家、重排序等）
// 参数 modelName: 模型名称
// 参数 message: 用户消息
// 返回: LLM 生成的回答、error
func (o *ollamaManager) ChatWithoutContext(modelName string, message string) (string, error) {
	o.logger.LogInfo("q#: " + message)
	
	o.mu.Lock()
	o.totalQCount++
	o.mu.Unlock()

	start := time.Now()
	response, err := sendChatRequest(o.domain, modelName, chatMessagesFromChatString(message))
	if err != nil {
		o.logger.LogError(fmt.Errorf("send chat err: %v", err), "sendchat")
		return "", fmt.Errorf("chat request failed: %w", err)
	}

	// 统计
	elapsed := time.Since(start)
	respMessage := response.Message

	o.mu.Lock()
	defer o.mu.Unlock()
	o.totalACount++
	o.totalDuration += elapsed
	o.totalToken += response.EvalCount

	o.logger.LogInfo("a#: " + respMessage.Content)

	return respMessage.Content, nil
}

// NewChat 创建新的对话上下文
// 为每个 Agent 创建独立的对话上下文，用于维护多轮对话历史
// 参数 modelName: 模型名称
// 参数 systemMessage: 系统提示词
// 返回: ChatContext 实例
func (o *ollamaManager) NewChat(modelName string, systemMessage string) *ChatContext {
	o.mu.Lock()
	defer o.mu.Unlock()
	chatId := o.autogenChatId
	o.autogenChatId++

	return newChat(modelName, chatId, systemMessage)
}

// NextChat 继续进行对话，维护上下文
// 将新的消息添加到历史记录，调用 LLM 生成回答，并保存回答到历史
// 参数 chatCtx: 对话上下文
// 参数 message: 用户消息
// 返回: LLM 生成的回答、error
// todo 上下文优化：实现有限上下文窗口，避免历史记录过长导致 token 超限
func (o *ollamaManager) NextChat(chatCtx *ChatContext, message string) (string, error) {
	o.logger.LogInfo("q" + strconv.Itoa(chatCtx.chatId) + ": " + message)
	
	o.mu.Lock()
	o.totalQCount++
	o.mu.Unlock()

	// 问题+历史记录
	chatCtx.addChatString(message)

	start := time.Now()
	response, err := sendChatRequest(o.domain, chatCtx.modelName, chatCtx.getMessages())
	if err != nil {
		o.logger.LogError(fmt.Errorf("send chat err: %v", err), "sendchat")
		return "", fmt.Errorf("chat request failed: %w", err)
	}

	// 统计
	elapsed := time.Since(start)
	respMessage := response.Message

	o.mu.Lock()
	defer o.mu.Unlock()
	o.totalACount++
	o.totalDuration += elapsed
	o.totalToken += response.EvalCount

	o.logger.LogInfo("a" + strconv.Itoa(chatCtx.chatId) + ": " + respMessage.Content)

	// 保存历史记录
	chatCtx.addMessage(respMessage)

	return respMessage.Content, nil
}

// GetTotalQCount 获取总问题数
func (o *ollamaManager) GetTotalQCount() int {
	o.mu.RLock()
	defer o.mu.RUnlock()
	return o.totalQCount
}

// GetTotalACount 获取总回答数
func (o *ollamaManager) GetTotalACount() int {
	o.mu.RLock()
	defer o.mu.RUnlock()
	return o.totalACount
}

// GetTotalDuration 获取总响应时间
func (o *ollamaManager) GetTotalDuration() time.Duration {
	o.mu.RLock()
	defer o.mu.RUnlock()
	return o.totalDuration
}

// GetTotalToken 获取总 token 使用量
func (o *ollamaManager) GetTotalToken() int {
	o.mu.RLock()
	defer o.mu.RUnlock()
	return o.totalToken
}
