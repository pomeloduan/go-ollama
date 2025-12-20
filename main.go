package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"go-ollama/logger"
	"go-ollama/ollama"
	"go-ollama/rag"
	"go-ollama/rule"
)

const ollamaDomain = "http://localhost:11434"

var ragManager *rag.RagManager

func main() {
	fmt.Println("--> Ollama Local Service Demo")

	// 创建错误日志
	logger, err := logger.NewErrorLogger("errors.log")
	if err != nil {
		log.Fatal(err)
	}
	defer logger.Close()

	rule := getMathRule()

	ollama, err := ollama.StartOllama(ollamaDomain, rule.DefaultModel(), logger)
	if err != nil {
		logger.LogError(err, "lauching")
	}

	chatId := newChat(ollama, rule, logger)

	// 提问 / 回答
	processInput(chatId, ollama, rule, logger)

	fmt.Println("提问：", ollama.TotalQCount)
	fmt.Println("回答：", ollama.TotalACount)
	fmt.Printf("总用时：%f\n", ollama.TotalDuration.Seconds())
	fmt.Println("token使用：", ollama.TotalToken)
}

func getHpRule() rule.Rule {
	return rule.UseRule("hp")
}

func getMathRule() rule.Rule {
	return rule.UseRule("math")
}

func newChat(ollama *ollama.OllamaManager, rule rule.Rule, logger *logger.ErrorLogger) int {
	if rule.ExternalSource() != "" {
		if ragManager == nil {
			var err error
			ragManager, err = rag.StartRag(logger)
			if err != nil {
				logger.LogError(err, "rag start")
			}
		}
		chProg, err := ragManager.PreprocessFromFile(rule.ExternalSource())
		if err != nil {
			logger.LogError(err, "rag preprocess")
		} else {
			fmt.Println("知识库预处理，需要一些时间")
			errCount := 0
			for prog := range chProg {
				if prog.Err != nil {
					logger.LogError(prog.Err, "rag preprocess", prog.Text)
					errCount++
				}
				fmt.Printf("\r进度：%.1f%% 第%d项，共%d项", prog.Percentage, prog.Current, prog.Total)
			}

			if errCount > 0 {
				fmt.Println(" 预处理错误" + strconv.Itoa(errCount) + "项")
			} else {
				fmt.Println()
			}
		}
	}
	return ollama.NewChat(rule.SystemMessage())
}

// 处理用户输入
func processInput(chatId int, ollama *ollama.OllamaManager, rule rule.Rule, logger *logger.ErrorLogger) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("等待提问 -->")
	fmt.Println("  h - 哈利波特小说相关问题")
	fmt.Println("  m - 数学问题")
	fmt.Println("  q - 退出程序")

	for {
		input, err := reader.ReadString('\n')
		if err != nil {
			logger.LogError(fmt.Errorf("read error: %v", err), "input")
			continue
		}

		input = strings.TrimSpace(input)

		switch input {
		case "q":
			return
		case "h":
			rule = getHpRule()
			chatId = newChat(ollama, rule, logger)
			fmt.Println("等待提问 -->")
			//q := "下棋的过程中，哈利代替了什么棋子，罗恩代替了什么棋子？"
			continue
		case "m":
			rule = getMathRule()
			chatId = newChat(ollama, rule, logger)
			fmt.Println("等待提问 -->")
			continue
		case "":
			continue

		default:
			if rule.ExternalSource() != "" {
				sources, err := ragManager.Query(input)
				if err != nil {
					logger.LogError(err, "rag query")
					continue
				}
				ollama.AddSystemMessage(chatId, "请阅读以下文字，并优先根据这段内容回答之后的问题：\n"+strings.Join(sources, "\n"))
			}
			var answer = ollama.NextChat(chatId, input)
			fmt.Println(rule.ParseAnswer(answer))
			fmt.Println("等待提问 -->")
		}

	}
}
