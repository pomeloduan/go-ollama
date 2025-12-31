package agent

import (
	"go-ollama/logger"
	"go-ollama/ollama"
	"go-ollama/rule"
)

// Reviewer 评审者 Agent，负责评估专家生成答案的质量
// 通过评分和评价来指导答案的改进
type Reviewer struct {
	ollama    *ollama.OllamaManager // Ollama 管理器
	modelName string                // 使用的模型名称
	rule      *rule.Rule            // 规则配置，包含评审相关的提示词
	chatCtx   *ollama.ChatContext   // 对话上下文
	logger    *logger.ErrorLogger   // 日志记录器
}

// newReviewer 创建并初始化评审者实例
func newReviewer(ollama *ollama.OllamaManager, rule *rule.Rule, logger *logger.ErrorLogger) *Reviewer {
	reviewer := Reviewer{
		ollama:    ollama,
		modelName: ollama.GetAvailableModelName("gemma"),
		rule:      rule,
		logger:    logger,
	}
	return &reviewer
}

// prepareChat 初始化评审者的对话上下文
// 设置评审者的系统提示词
func (this *Reviewer) prepareChat() {
	this.chatCtx = this.ollama.NewChat(this.modelName, this.rule.ReviewerSystemMessage())
}

// review 评审专家生成的答案
// 参数 question: 原始问题
// 参数 answer: 专家生成的答案
// 返回: ReviewResult，包含评分和评价文本
func (this *Reviewer) review(question string, answer string) rule.ReviewResult {
	// 延迟初始化
	if this.chatCtx == nil {
		this.prepareChat()
	}
	// 构建评审提示词
	message := this.rule.ReviewMessage(question, answer)
	// 调用 LLM 进行评审
	review := this.ollama.NextChat(this.chatCtx, message)
	// 解析评审结果（提取分数和评价）
	return this.rule.ParseReview(review)
}
