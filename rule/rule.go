package rule

import "sync"

type Rule interface {
	DefaultModel() string
	SystemMessage() string
	ParseAnswer(text string) string
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
