package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"go-ollama/logger"
	"go-ollama/ollama"
	"go-ollama/rule"
)

const ollamaDomain = "http://localhost:11434"

func main() {
	fmt.Println("--> Ollama Local Service Demo")

	// 创建错误日志
	logger, err := logger.NewErrorLogger("errors.log")
	if err != nil {
		log.Fatal(err)
	}
	defer logger.Close()

	rule := rule.UseRule("math")

	ollama, err := ollama.StartOllamaClient(ollamaDomain, rule.DefaultModel(), logger)
	if err != nil {
		logger.LogError(err, "lauching")
	}

	chatId := ollama.NewChat(rule.SystemMessage())

	// 提问 / 回答
	processInput(ollama, chatId, rule, logger)

	fmt.Println("提问：", ollama.TotalQCount)
	fmt.Println("回答：", ollama.TotalACount)
	fmt.Printf("总用时：%f\n", ollama.TotalDuration.Seconds())
	fmt.Println("token使用：", ollama.TotalToken)
}

// 处理用户输入
func processInput(ollama *ollama.OllamaClient, chatId int, rule rule.Rule, logger *logger.ErrorLogger) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("等待提问 -->")
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

		case "":
			continue

		default:
			var answer = ollama.NextChat(input, chatId)
			fmt.Println(rule.ParseAnswer(answer))
		}
	}
}
