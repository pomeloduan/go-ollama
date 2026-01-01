package agent

import (
	"go-ollama/logger"
	"go-ollama/ollama"
	"go-ollama/rag"
	"go-ollama/rule"
	"sync"
)

// AgentManager Agent 管理器接口
// 创建/管理所有agent生命周期
type AgentManager interface {
	Chat(chat string) string
}

// agentManager Agent 管理器实现（包私有）
type agentManager struct {
	ollama        ollama.OllamaManager
	rag           rag.RagManager
	rule          rule.RuleManager
	coordinator   *Coordinator
	generalAgent  *Specialist
	specialistMap map[string]*Specialist
	reviewerMap   map[string]*Reviewer
	logger        logger.ErrorLogger
}

// rewriteScore 评审分数阈值，低于此分数将触发重写流程
const rewriteScore = 80

var (
	agentInstance *agentManager
	agentOnce     sync.Once
)

// newAgentManager 创建并初始化 Agent 管理器实例
// 初始化顺序：规则管理器 -> RAG 管理器 -> 协调者 -> 专家和评审者
// 参数 ollama: Ollama 管理器
// 参数 logger: 日志记录器
// 返回: agentManager 实例、error
func newAgentManager(ollama ollama.OllamaManager, logger logger.ErrorLogger) (*agentManager, error) {
	// 1 rule
	// 规则管理器需要先起
	ruleManager, err := rule.StartRuleManager()
	if err != nil {
		return nil, err
	}

	// 2 rag
	reranker := newReranker(ollama, ruleManager)
	ragMgr := rag.StartRagManager(reranker)

	// 3 coordinator
	coordinator := newCoordinator(ollama, ruleManager)

	// 4 specialist + reviewer
	general := newSpecialist(ollama, ragMgr, ruleManager.GetGeneralRule(), logger)
	specialistMap := make(map[string]*Specialist)
	reviewerMap := make(map[string]*Reviewer)
	for _, rule := range ruleManager.GetAllRules() {
		specialist := newSpecialist(ollama, ragMgr, rule, logger)
		specialistMap[rule.Name()] = specialist
		coordinator.addSpecialist(rule.Name(), rule.Introduction())
		if rule.NeedReviewer() {
			reviewer := newReviewer(ollama, rule, logger)
			reviewerMap[rule.Name()] = reviewer
		}
	}

	return &agentManager{
		ollama:        ollama,
		rag:           ragMgr,
		rule:          ruleManager,
		coordinator:   coordinator,
		generalAgent:  general,
		specialistMap: specialistMap,
		reviewerMap:   reviewerMap,
		logger:        logger,
	}, nil
}

// StartAgentManager 获取 Agent 管理器单例
// 初始化顺序：规则管理器 -> RAG 管理器 -> 协调者 -> 专家和评审者
// 返回初始化完成的 AgentManager 实例
// todo
// rag工程化
// function call
func StartAgentManager(ollama ollama.OllamaManager, logger logger.ErrorLogger) (AgentManager, error) {
	var err error
	agentOnce.Do(func() {
		agentInstance, err = newAgentManager(ollama, logger)
	})

	if err != nil {
		return nil, err
	}
	return agentInstance, nil
}

// Chat 处理用户输入的聊天请求，实现完整的 Agent 协作流程
// 流程：1. 协调者选择专家 2. 专家回答问题 3. 评审者评估 4. 低分重写
// 参数 chat: 用户输入的问题
// 返回: Agent 生成的回答
func (a *agentManager) Chat(chat string) string {
	// 1. 调用协调者选择最适合的专家
	name, err := a.coordinator.askForSpecialistName(chat)
	if err != nil {
		a.logger.LogError(err, "coordinator askForSpecialistName")
		// 如果协调者失败，使用通用专家
		name = ""
	}
	specialist, ok := a.specialistMap[name]
	// 如果没有匹配的专家，使用通用专家
	if !ok {
		specialist = a.generalAgent
	}

	// 2. 调用专家生成回答
	answer, err := specialist.chat(chat)
	if err != nil {
		a.logger.LogError(err, "specialist chat")
		return "抱歉，处理问题时出现错误，请稍后重试。"
	}
	
	// 3. 如果有评审者，进行质量评估
	reviewer, ok := a.reviewerMap[name]
	if ok {
		review := reviewer.review(chat, answer)
		// 4. 如果分数低于阈值，触发重写流程
		if review.Score < rewriteScore {
			// 获取 specialist 对应的规则并构建重写消息
			rule := specialist.getRule()
			message := rule.RewriteMessage(review.Review)
			rewrittenAnswer, err := specialist.chat(message)
			if err != nil {
				a.logger.LogError(err, "specialist rewrite")
				// 如果重写失败，返回原始答案
				return answer
			}
			answer = rewrittenAnswer
		}
	}
	return answer
}
