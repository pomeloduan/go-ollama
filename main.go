package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"go-ollama/agent"
	"go-ollama/logger"
	"go-ollama/ollama"
)

const ollamaDomain = "http://localhost:11434"

// todo mcp func call
func main() {
	fmt.Println("--> Ollama Local Service Demo")

	// 创建错误日志
	logger, err := logger.NewErrorLogger("info.log")
	if err != nil {
		log.Fatal(err)
	}
	defer logger.Close()

	// rule := getMathRule()

	ollama, err := ollama.StartOllama(ollamaDomain, logger)
	if err != nil {
		logger.LogError(err, "lauching")
		return
	}

	agent, err := agent.StartAgentManager(ollama, logger)
	if err != nil {
		logger.LogError(err, "lauching")
		return
	}

	// chatId := newChat(ollama, rule, logger)

	// 提问 / 回答
	processInput(agent, logger)

	fmt.Println("提问：", ollama.TotalQCount)
	fmt.Println("回答：", ollama.TotalACount)
	fmt.Printf("总用时：%f\n", ollama.TotalDuration.Seconds())
	fmt.Println("token使用：", ollama.TotalToken)
}

// 处理用户输入
func processInput(agent *agent.AgentManager, logger *logger.ErrorLogger) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("等待提问 -->")
	fmt.Println("  q - 退出程序")

	// test
	// agent.Chat("1+2+3+...+100的值是多少？")
	// agent.Chat("请你以猫为主题，写一首诗。")
	// agent.Chat("在《哈利波特》中，为了通过下棋关卡，哈利代替了什么棋子，罗恩代替了什么棋子？")
	// agent.Chat("今天天气怎么样")

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
			answer := agent.Chat(input)
			fmt.Println(answer)
			fmt.Println("等待提问 -->")
		}

	}
}
