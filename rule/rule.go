package rule

import (
	"regexp"
	"strings"
	"sync"
)

type Rule interface {
	DefaultModel() string
	SystemMessage() string
	ParseAnswer(text string) string
	ExternalSource() string
	ExternalSourceMessage() string
}

func UseRule(name string) Rule {
	return ruleMap[name]
}

var ruleMap map[string]Rule
var muLock sync.Mutex

func registerRule(name string, rule Rule) {
	muLock.Lock()
	defer muLock.Unlock()

	if ruleMap == nil {
		ruleMap = make(map[string]Rule)
	}
	ruleMap[name] = rule
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
