package agent

import (
	"go-ollama/logger"
	"go-ollama/ollama"
	"go-ollama/rule"
)

// Reviewer 评审者 Agent，负责评估专家生成答案的质量
// 通过评分和评价来指导答案的改进
type Reviewer struct {
	ollama    ollama.OllamaManager // Ollama 管理器
	modelName string               // 使用的模型名称
	rule      *rule.Rule          // 规则配置，包含评审相关的提示词
	chatCtx   *ollama.ChatContext // 对话上下文
	logger    logger.ErrorLogger  // 日志记录器
}

// newReviewer 创建并初始化评审者实例
func newReviewer(ollama ollama.OllamaManager, rule *rule.Rule, logger logger.ErrorLogger) *Reviewer {
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
func (r *Reviewer) prepareChat() {
	r.chatCtx = r.ollama.NewChat(r.modelName, r.rule.ReviewerSystemMessage())
}

// review 评审专家生成的答案
// 参数 question: 原始问题
// 参数 answer: 专家生成的答案
// 返回: ReviewResult，包含评分和评价文本
func (r *Reviewer) review(question string, answer string) rule.ReviewResult {
	// 延迟初始化
	if r.chatCtx == nil {
		r.prepareChat()
	}
	// 构建评审提示词
	message := r.rule.ReviewMessage(question, answer)
	// 调用 LLM 进行评审
	review, err := r.ollama.NextChat(r.chatCtx, message)
	if err != nil {
		// 如果评审失败，返回空结果
		r.logger.LogError(err, "review")
		return rule.ReviewResult{}
	}
	// 解析评审结果（提取分数和评价）
	return r.rule.ParseReview(review)
}
