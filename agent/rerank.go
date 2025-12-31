package agent

import (
	"go-ollama/ollama"
	"go-ollama/rule"
)

// Reranker 重排器，用于对 RAG 检索结果进行重排
// 使用 LLM 评估检索到的候选文档与问题的相关性，选择最相关的文档
type Reranker struct {
	ollama    ollama.OllamaManager // Ollama 管理器
	modelName string               // 使用的模型名称
	rule      rule.RuleManager     // 规则管理器，包含重排提示词模板
}

// newReranker 创建并初始化重排序器实例
func newReranker(ollama ollama.OllamaManager, rule rule.RuleManager) *Reranker {
	reranker := Reranker{
		ollama:    ollama,
		modelName: ollama.GetAvailableModelName("gemma"),
		rule:      rule,
	}
	return &reranker
}

// RankCandidate 对候选文档进行重排序
// 使用 LLM 评估每个候选文档与问题的相关性，返回最相关的文档
// 参数 candidates: 候选文档文本，多个文档用换行分隔
// 参数 text: 用户问题
// 参数 num: 返回的文档数量
// 返回: 重排序后的文档文本
func (this *Reranker) RankCandidate(candidates string, text string, num int) string {
	message := this.rule.RerankMessage(candidates, text, num)
	return this.ollama.ChatWithoutContext(this.modelName, message)
}
