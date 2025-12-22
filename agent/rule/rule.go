package rule

import (
	"sync"
)

type Rule interface {
	SystemMessage() string
	ParseAnswer(text string) string
	ExternalSource() string
	ExternalSourceMessage() string
}

var ruleMap map[string]Rule
var muLock sync.Mutex

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
