package rule

import (
	"sync"
)

type Rule interface {
	Introduction() string
	SystemMessage() string
	ParseAnswer(text string) string
	SourceFile() string
	MessageFromSource(source string, question string) string
}

var ruleMap map[string]Rule
var muLock sync.Mutex

type GeneralRule struct {
}

func (this *GeneralRule) Introduction() string {
	return ""
}

func (this *GeneralRule) SystemMessage() string {
	return ""
}

func (this *GeneralRule) SourceFile() string {
	return ""
}

func (this *GeneralRule) MessageFromSource(source string, question string) string {
	return ""
}

func (this *GeneralRule) ParseAnswer(text string) string {
	return text
}

func GetRule(name string) Rule {
	return ruleMap[name]
}

func AllRuleNames() []string {
	names := []string{}
	for k := range ruleMap {
		names = append(names, k)
	}
	return names
}

func registerRule(name string, rule Rule) {
	muLock.Lock()
	defer muLock.Unlock()

	if ruleMap == nil {
		ruleMap = make(map[string]Rule)
	}
	ruleMap[name] = rule
}
