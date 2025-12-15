package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
)

const ollamaDomain = "http://localhost:11434"
const defaultModel = "deepseek"

// 查找第一个包含子字符串的元素
func findFirstContaining(arr []string, substr string) (string, bool) {
	for _, str := range arr {
		if strings.Contains(strings.ToLower(str), strings.ToLower(substr)) {
			return str, true
		}
	}
	return "", false
}

func main() {
	fmt.Println("--> Ollama Go demo")

	// 创建错误日志
	logger, err := NewErrorLogger("errors.log")
	if err != nil {
		log.Fatal(err)
	}
	defer logger.Close()

	ollama, err := StartOllamaClient(ollamaDomain, defaultModel, logger)
	if err != nil {
		logger.LogError(err, "lauching")
	}

	chatId := ollama.NewChat(`你是一位数学老师。你的任务是解答数学题。

# 行动格式:
你的回答必须严格遵循以下格式。首先是这是不是数学问题，然后是解答。
isMath: [这里是判断这是不是数学问题]
resolvation: [这里是你的解答]`)

	// 提问 / 回答
	processInput(ollama, chatId)

	fmt.Println("提问：", ollama.TotalQCount)
	fmt.Println("回答：", ollama.TotalACount)
	fmt.Printf("总用时：%f\n", ollama.TotalDuration.Seconds())
	fmt.Println("token使用：", ollama.TotalToken)
}

var logger *ErrorLogger

// 处理用户输入
func processInput(ollama *OllamaClient, chatId int) {
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
			var isMath, resolvation, canAnswer = parseKeyValuePairs(answer, "isMath", "resolvation")
			if canAnswer && isMath == "true" {
				fmt.Println(CompactEmptyLines(resolvation))
			} else {
				fmt.Println("不是数学问题呢")
			}
		}
	}
}

// 解析格式 "key0:text0 key1:text1" 的字符串
func parseKeyValuePairs(input, key0, key1 string) (string, string, bool) {
	pKey0 := strings.Index(input, key0+":")
	if pKey0 == -1 {
		return "", "", false
	}

	pKey1 := strings.Index(input, key1+":")
	if pKey1 == -1 {
		return "", "", false
	}

	var text0 = strings.TrimSpace(input[pKey0+len(key0)+1 : pKey1])
	var text1 = strings.TrimSpace(input[pKey1+len(key1)+1:])

	return text0, text1, true
}

// 压缩空行
func CompactEmptyLines(input string) string {
	re := regexp.MustCompile(`\n\s*\n`)

	output := re.ReplaceAllString(input, "\n")

	output = strings.TrimLeft(output, "\n")
	output = strings.TrimRight(output, "\n")

	return output
}
