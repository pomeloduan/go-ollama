package agent

import (
	"fmt"
	"go-ollama/agent/rag"
	"go-ollama/agent/rule"
	"go-ollama/logger"
	"go-ollama/ollama"
	"strconv"
)

type Specialist struct {
	ollama    *ollama.OllamaManager
	rag       *rag.RagManager
	modelName string
	rule      rule.Rule
	chatCtx   *ollama.ChatContext
	ragCtx    *rag.RagContext
	logger    *logger.ErrorLogger
}

func startSpecialist(ollama *ollama.OllamaManager, rag *rag.RagManager, rule rule.Rule, logger *logger.ErrorLogger) *Specialist {
	specialist := Specialist{
		ollama:    ollama,
		rag:       rag,
		modelName: ollama.GetAvailableModelName("deepseek"),
		rule:      rule,
		logger:    logger,
	}
	return &specialist
}

func (this *Specialist) prepareChat() {
	if this.rule.SourceFile() != "" {
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
	this.chatCtx = this.ollama.NewChat(this.modelName, this.rule.SystemMessage())
}

func (this *Specialist) chat(chat string) string {
	if this.chatCtx == nil {
		this.prepareChat()
	}
	if this.rule.SourceFile() != "" {
		chSource, err := this.rag.Query(this.ragCtx, chat)
		if err != nil {
			this.logger.LogError(err, "rag query")
		}
		source := ""
		for s := range chSource {
			source += s
		}
		chat = this.rule.MessageFromSource(source, chat)
	}
	var answer = this.ollama.NextChat(this.chatCtx, chat)
	return this.rule.ParseAnswer(answer)
}
