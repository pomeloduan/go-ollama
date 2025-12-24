package rule

import (
	"strconv"
	"strings"
)

type RuleManager struct {
	ruleMap map[string]*Rule
	config  *ChatConfig
}

func StartRuleManager() (*RuleManager, error) {
	// read file
	config, err := readConfig("./rule/config.yml")
	if err != nil {
		return nil, err
	}

	ruleMap := make(map[string]*Rule)
	for name, ruleCfg := range config.Rules {
		ruleMap[name] = &Rule{name: name, config: &ruleCfg}
	}
	return &RuleManager{ruleMap: ruleMap, config: config}, nil
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

func (this *RuleManager) RerankMessage(question string, number int, candidates string) string {
	replacer := strings.NewReplacer(
		"{question}", question,
		"{number}", strconv.Itoa(number),
		"{candidates}", candidates,
	)
	return replacer.Replace(this.config.RerankMessage)
}

func (this *RuleManager) CoordinatorMessage(question string) string {
	replacer := strings.NewReplacer(
		"{question}", question,
	)
	return replacer.Replace(this.config.CoordinatorMessage)
}

func (this *RuleManager) CoordinatorSpecialistMessage(name string, introduction string) string {
	replacer := strings.NewReplacer(
		"{name}", name,
		"{introduction}", introduction,
	)
	return replacer.Replace(this.config.CoordinatorSpecialistMessage)
}

type Rule struct {
	name   string
	config *RuleConfig
}

func (this *Rule) Name() string {
	return this.name
}

func (this *Rule) Introduction() string {
	if this.config == nil {
		return ""
	}
	return this.config.Introduction
}

func (this *Rule) SystemMessage() string {
	if this.config == nil {
		return ""
	}
	return this.config.SystemMessage
}

func (this *Rule) NeedRag() bool {
	if this.config == nil {
		return false
	}
	return this.SourceFile() != ""
}

func (this *Rule) SourceFile() string {
	if this.config == nil {
		return ""
	}
	return this.config.SourceFile
}

func (this *Rule) SourceMessage(source string, question string) string {
	if this.config == nil {
		return ""
	}
	replacer := strings.NewReplacer(
		"{source}", source,
		"{question}", question,
	)
	return replacer.Replace(this.config.SourceMessage)
}

func (this *Rule) NeedReviewer() bool {
	if this.config == nil {
		return false
	}
	return this.ReviewerSystemMessage() != ""
}

func (this *Rule) ReviewerSystemMessage() string {
	if this.config == nil {
		return ""
	}
	return this.config.ReviewerSystemMessage
}

func (this *Rule) ReviewMessage(question string, answer string) string {
	if this.config == nil {
		return ""
	}
	replacer := strings.NewReplacer(
		"{question}", question,
		"{answer}", answer,
	)
	return replacer.Replace(this.config.ReviewMessage)
}

func (this *Rule) RewriteMessage(review string) string {
	if this.config == nil {
		return ""
	}
	replacer := strings.NewReplacer(
		"{review}", review,
	)
	return replacer.Replace(this.config.RewriteMessage)
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
