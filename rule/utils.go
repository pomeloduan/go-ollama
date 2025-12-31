package rule

import (
	"regexp"
	"strings"
)

// parseKeyValueText 解析键值对格式的文本
// 从文本中提取两个键对应的值
// 格式示例："score: 75 review: 这是一段评价"
// 参数 input: 输入文本
// 参数 key0: 第一个键名
// 参数 key1: 第二个键名
// 返回: key0 的值、key1 的值、是否解析成功
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

// compactEmptyLines 压缩文本中的连续空行
// 将多个连续的空行压缩为一个，并去除首尾的换行符
// 参数 input: 输入文本
// 返回: 处理后的文本
func compactEmptyLines(input string) string {
	re := regexp.MustCompile(`\n\s*\n`)

	output := re.ReplaceAllString(input, "\n")

	output = strings.TrimLeft(output, "\n")
	output = strings.TrimRight(output, "\n")

	return output
}
