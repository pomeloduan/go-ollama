package rule

import (
	"regexp"
	"strings"
)

func init() {
	registerRule("math", &MathRule{})
}

type MathRule struct {
}

func (this *MathRule) DefaultModel() string {
	return "deepseek"
}

func (this *MathRule) SystemMessage() string {
	return `你是一位数学老师。你的任务是解答数学题。

# 行动格式:
你的回答必须严格遵循以下格式。首先是这是不是数学问题，然后是解答。
isMath: [这里是判断这是不是数学问题]
resolvation: [这里是你的解答]`
}

func (this *MathRule) ParseAnswer(text string) string {
	var isMath, resolvation, formatedAnswer = parseKeyValueText(text, "isMath", "resolvation")
	if formatedAnswer && isMath == "true" {
		return compactEmptyLines(resolvation)
	} else {
		return "不是数学问题"
	}
}

// 解析格式 "key0:text0 key1:text1" 的字符串
func parseKeyValueText(input, key0, key1 string) (string, string, bool) {
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
func compactEmptyLines(input string) string {
	re := regexp.MustCompile(`\n\s*\n`)

	output := re.ReplaceAllString(input, "\n")

	output = strings.TrimLeft(output, "\n")
	output = strings.TrimRight(output, "\n")

	return output
}
