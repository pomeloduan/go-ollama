package rule

import (
	"strconv"
)

type RuleManager struct {
	ruleMap map[string]*Rule
}

func StartRuleManager() *RuleManager {
	// read file
	ruleConfigMap, err := readConfig("./rule/config.yml")
	if err != nil {

	}

	ruleMap := make(map[string]*Rule)
	for name, ruleConfig := range ruleConfigMap {
		ruleMap[name] = &Rule{name: name, config: ruleConfig}
	}
	return &RuleManager{ruleMap: ruleMap}
}

func (this *RuleManager) GetGeneralRule() *Rule {
	return &Rule{}
}

func (this *RuleManager) GetAllRules() []*Rule {
	var rules []*Rule
	for _, r := range this.ruleMap {
		rules = append(rules, r)
	}
	return rules
}

type Rule struct {
	name   string
	config RuleConfig
}

func (this *Rule) Name() string {
	return this.name
}

func (this *Rule) Introduction() string {
	return this.config.introduction
}

func (this *Rule) SystemMessage() string {
	return this.config.systemMessage
}

func (this *Rule) NeedRag() bool {
	return this.SourceFile() != ""
}

func (this *Rule) SourceFile() string {
	return this.config.sourceFile
}

func (this *Rule) MessageFromSource(source string, question string) string {
	return "请阅读以下文字，并优先根据这段内容回答之后的问题：\n" + source + "\n问题：" + question
}

func (this *Rule) NeedReviewer() bool {
	return this.ReviewerSystemMessage() != ""
}

func (this *Rule) ReviewerSystemMessage() string {
	return this.config.reviewerSystemMessage
}

func (this *Rule) ParseReview(text string) ReviewResult {
	var scoreString, review, isFormated = parseKeyValueText(text, "score", "review")
	score, _ := strconv.Atoi(scoreString)
	if isFormated {
		return ReviewResult{Score: score, Review: compactEmptyLines(review)}
	} else {
		return ReviewResult{}
	}
}

type ReviewResult struct {
	Score  int
	Review string
}
