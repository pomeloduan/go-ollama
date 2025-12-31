package agent

import (
	"fmt"
	"go-ollama/logger"
	"go-ollama/ollama"
	"go-ollama/rag"
	"go-ollama/rule"
	"strconv"
)

// Specialist 专家 Agent，负责处理特定领域的问题
// 支持 RAG（检索增强生成）来提升回答的准确性
type Specialist struct {
	ollama    *ollama.OllamaManager // Ollama 管理器
	rag       *rag.RagManager       // RAG 管理器，用于检索外部知识
	modelName string                // 使用的 LLM 模型名称
	rule      *rule.Rule            // 规则配置
	chatCtx   *ollama.ChatContext   // 对话上下文，维护多轮对话历史
	ragCtx    *rag.RagContext       // RAG 上下文，存储知识库信息
	logger    *logger.ErrorLogger   // 日志记录器
}

// startSpecialist 创建并初始化专家实例
// 参数 rag: RAG 管理器
// 参数 rule: 专家规则配置
func startSpecialist(ollama *ollama.OllamaManager, rag *rag.RagManager, rule *rule.Rule, logger *logger.ErrorLogger) *Specialist {
	specialist := Specialist{
		ollama:    ollama,
		rag:       rag,
		modelName: ollama.GetAvailableModelName("deepseek"),
		rule:      rule,
		logger:    logger,
	}
	return &specialist
}

// prepareChat 初始化对话环境
// 如果需要 RAG，会预处理知识库（文本分块、向量化、存储）
// 创建对话上下文，设置系统提示词
func (this *Specialist) prepareChat() {
	if this.rule.NeedRag() {
		// 导入外部知识库，进行预处理
		ragCtx, chProg, err := this.rag.PreprocessFromFile(this.rule.SourceFile())
		if err != nil {
			this.logger.LogError(err, "rag preprocess")
		} else {
			this.ragCtx = ragCtx
			fmt.Println("需要导入外部知识库，请稍等...")
			errCount := 0
			for p := range chProg {
				if p.Err != nil {
					this.logger.LogError(p.Err, "rag preprocess", p.Text)
					errCount++
				}
				fmt.Printf("\r进度：%.1f%% 第%d项，共%d项", p.Percentage, p.Current, p.Total)
			}
			if errCount > 0 {
				fmt.Println(" 预处理错误" + strconv.Itoa(errCount) + "项")
			} else {
				fmt.Println()
			}
		}
	}
	// 创建对话上下文，设置系统提示词
	this.chatCtx = this.ollama.NewChat(this.modelName, this.rule.SystemMessage())
}

// chat 处理用户问题并生成回答
// 如果配置了 RAG，会先检索相关文档，然后将检索结果和问题一起发送给 LLM
// 参数 chat: 用户输入的问题
// 返回: 专家生成的回答
func (this *Specialist) chat(chat string) string {
	// 延迟初始化，首次调用时准备对话环境
	if this.chatCtx == nil {
		this.prepareChat()
	}
	
	// 如果需要 RAG，检索相关文档并增强问题
	if this.rule.NeedRag() {
		chSource, err := this.rag.Query(this.ragCtx, chat, this.rule)
		if err != nil {
			this.logger.LogError(err, "rag query")
		}
		// 从 channel 中读取检索结果
		source := ""
		for s := range chSource {
			source += s
		}
		// 将检索到的文档和问题组合成新的提示词
		chat = this.rule.SourceMessage(source, chat)
	}
	// 调用 LLM 生成回答，维护对话上下文
	return this.ollama.NextChat(this.chatCtx, chat)
}
