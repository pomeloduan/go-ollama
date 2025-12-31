package rule

import (
	"strconv"
	"strings"
)

// RuleManager 规则管理器
// 负责读取和解析 YAML 配置文件，管理所有 专家/评审者 Agent 的规则配置
type RuleManager struct {
	ruleMap map[string]*Rule // 规则名称到规则对象的映射
	config  *ChatConfig      // 完整配置对象
}

// StartRuleManager 启动规则管理器
// 从配置文件读取规则，创建规则对象
// 返回: RuleManager 实例、error
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

// GetGeneralRule 获取通用规则（用于通用专家）
// 返回一个空的 Rule 对象，表示没有特定配置
func (this *RuleManager) GetGeneralRule() *Rule {
	return &Rule{}
}

// GetAllRules 获取所有规则对象
// 返回: 规则对象数组
func (this *RuleManager) GetAllRules() []*Rule {
	var rules []*Rule
	for _, r := range this.ruleMap {
		rules = append(rules, r)
	}
	return rules
}

// RerankMessage 构建重排序提示词
// 替换模板中的占位符（{question}, {number}, {candidates}）
func (this *RuleManager) RerankMessage(candidates string, question string, number int) string {
	replacer := strings.NewReplacer(
		"{question}", question,
		"{number}", strconv.Itoa(number),
		"{candidates}", candidates,
	)
	return replacer.Replace(this.config.RerankMessage)
}

// CoordinatorMessage 构建协调者提示词
// 替换模板中的占位符（{question}）
func (this *RuleManager) CoordinatorMessage(question string) string {
	replacer := strings.NewReplacer(
		"{question}", question,
	)
	return replacer.Replace(this.config.CoordinatorMessage)
}

// CoordinatorSpecialistMessage 构建协调者专家信息提示词
// 替换模板中的占位符（{name}, {introduction}）
func (this *RuleManager) CoordinatorSpecialistMessage(name string, introduction string) string {
	replacer := strings.NewReplacer(
		"{name}", name,
		"{introduction}", introduction,
	)
	return replacer.Replace(this.config.CoordinatorSpecialistMessage)
}

// Rule 单个规则配置
// 包含一个专家 Agent 或评审者的所有配置信息
type Rule struct {
	name   string      // 规则名称（也是专家名称）
	config *RuleConfig // 规则配置对象
}

// Name 获取规则名称
func (this *Rule) Name() string {
	return this.name
}

// Introduction 获取专家介绍
// 用于协调者匹配问题
func (this *Rule) Introduction() string {
	if this.config == nil {
		return ""
	}
	return this.config.Introduction
}

// SystemMessage 获取系统提示词
// 定义 Agent 的角色和行为
func (this *Rule) SystemMessage() string {
	if this.config == nil {
		return ""
	}
	return this.config.SystemMessage
}

// NeedRag 判断是否需要 RAG
// 如果配置了源文件，则需要 RAG
func (this *Rule) NeedRag() bool {
	if this.config == nil {
		return false
	}
	return this.SourceFile() != ""
}

// SourceFile 获取 RAG 源文件路径
func (this *Rule) SourceFile() string {
	if this.config == nil {
		return ""
	}
	return this.config.SourceFile
}

// SourceMessage 构建包含检索文档的提示词
// 将检索到的文档和问题组合，替换模板中的占位符（{source}, {question}）
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

// NeedReviewer 判断是否需要评审者
// 如果配置了评审者系统提示词，则需要评审者
func (this *Rule) NeedReviewer() bool {
	if this.config == nil {
		return false
	}
	return this.ReviewerSystemMessage() != ""
}

// ReviewerSystemMessage 获取评审者系统提示词
func (this *Rule) ReviewerSystemMessage() string {
	if this.config == nil {
		return ""
	}
	return this.config.ReviewerSystemMessage
}

// ReviewMessage 构建评审提示词
// 将问题和答案组合，替换模板中的占位符（{question}, {answer}）
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

// RewriteMessage 构建重写提示词
// 将评审反馈组合到提示词中，替换模板中的占位符（{review}）
func (this *Rule) RewriteMessage(review string) string {
	if this.config == nil {
		return ""
	}
	replacer := strings.NewReplacer(
		"{review}", review,
	)
	return replacer.Replace(this.config.RewriteMessage)
}

// ParseReview 解析评审结果
// 从 LLM 返回的文本中提取分数和评价
// 期望格式：score: 分数\nreview: 评价
// 参数 text: LLM 返回的评审文本
// 返回: ReviewResult
func (this *Rule) ParseReview(text string) ReviewResult {
	var scoreString, review, isFormatted = parseKeyValueText(text, "score", "review")
	score, _ := strconv.Atoi(scoreString)
	if isFormatted {
		return ReviewResult{Score: score, Review: compactEmptyLines(review)}
	} else {
		return ReviewResult{}
	}
}

// ReviewResult 评审结果
type ReviewResult struct {
	Score  int    // 评分（0-100）
	Review string // 评价文本
}
