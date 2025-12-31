package agent

import (
	"go-ollama/logger"
	"go-ollama/ollama"
	"go-ollama/rag"
	"go-ollama/rule"
)

// agent 管理器
// 创建/管理所有agent生命周期
type AgentManager struct {
	ollama        *ollama.OllamaManager
	rag           *rag.RagManager
	rule          *rule.RuleManager
	coordinator   *Coordinator
	generalAgent  *Specialist
	specialistMap map[string]*Specialist
	reviewerMap   map[string]*Reviewer
	logger        *logger.ErrorLogger
}

// rewriteScore 评审分数阈值，低于此分数将触发重写流程
const rewriteScore = 80

// StartAgentManager 启动 Agent 管理器，初始化所有组件
// 初始化顺序：规则管理器 -> RAG 管理器 -> 协调者 -> 专家和评审者
// 返回初始化完成的 AgentManager 实例
// todo
// rag工程化
// function call
func StartAgentManager(ollama *ollama.OllamaManager, logger *logger.ErrorLogger) (*AgentManager, error) {
	// 1 rule
	// 规则管理器需要先起
	ruleManager, err := rule.StartRuleManager()
	if err != nil {
		return nil, err
	}

	// 2 rag
	reranker := startReranker(ollama, ruleManager)
	rag := rag.StartRag(reranker)

	// 3 coordinator
	coordinator := startCoordinator(ollama, ruleManager)

	// 4 specialist + reviewer
	general := startSpecialist(ollama, rag, ruleManager.GetGeneralRule(), logger)
	specialistMap := make(map[string]*Specialist)
	reviewerMap := make(map[string]*Reviewer)
	for _, rule := range ruleManager.GetAllRules() {
		specialist := startSpecialist(ollama, rag, rule, logger)
		specialistMap[rule.Name()] = specialist
		coordinator.addSpecialist(rule.Name(), rule.Introduction())
		if rule.NeedReviewer() {
			reviewer := startReviewer(ollama, rule, logger)
			reviewerMap[rule.Name()] = reviewer
		}
	}

	agentManager := AgentManager{
		ollama:        ollama,
		rag:           rag,
		coordinator:   coordinator,
		generalAgent:  general,
		specialistMap: specialistMap,
		reviewerMap:   reviewerMap,
		logger:        logger,
	}
	return &agentManager, nil
}

// Chat 处理用户输入的聊天请求，实现完整的 Agent 协作流程
// 流程：1. 协调者选择专家 2. 专家回答问题 3. 评审者评估 4. 低分重写
// 参数 chat: 用户输入的问题
// 返回: Agent 生成的回答
func (this *AgentManager) Chat(chat string) string {
	// 1. 调用协调者选择最适合的专家
	name := this.coordinator.askForSpecialistName(chat)
	specialist, ok := this.specialistMap[name]
	// 如果没有匹配的专家，使用通用专家
	if !ok {
		specialist = this.generalAgent
	}

	// 2. 调用专家生成回答
	answer := specialist.chat(chat)
	
	// 3. 如果有评审者，进行质量评估
	reviewer, ok := this.reviewerMap[name]
	if ok {
		review := reviewer.review(chat, answer)
		// 4. 如果分数低于阈值，触发重写流程
		if review.Score < rewriteScore {
			message := specialist.rule.RewriteMessage(review.Review)
			answer = specialist.chat(message)
		}
	}
	return answer
}
