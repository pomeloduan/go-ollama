package rule

import "strconv"

func init() {
	registerRule("poet", &PoetRule{})
}

type PoetRule struct {
}

func (this *PoetRule) Introduction() string {
	return "擅于创作诗歌，涉及诗歌相关都可以来问。"
}

func (this *PoetRule) SystemMessage() string {
	return "你是一位诗歌的爱好者。你的任务是用王维的风格创作诗歌"
}

func (this *PoetRule) SourceFile() string {
	return ""
}

func (this *PoetRule) MessageFromSource(source string, question string) string {
	return ""
}

func (this *PoetRule) ReviewerSystemMessage() string {
	return `你是一位诗歌的爱好者。你的任务是给诗歌评分，你会收到一段写作要求，和一首诗歌作品，请你从诗歌是否符合要求，以及诗歌的语言、结构和审美等综合评价
# 行动格式:
你的回答必须严格遵循以下格式。首先是诗歌的分数，然后是评价，90以下需要尽量给出不够好的点和修改意见。
score: [这里是诗歌的分数，使用0-100回答]
review: [这里是你的评价]`
}

func (this *PoetRule) ParseReview(text string) ReviewResult {
	var scoreString, review, isFormated = parseKeyValueText(text, "score", "review")
	score, _ := strconv.Atoi(scoreString)
	if isFormated {
		return ReviewResult{Score: score, Review: compactEmptyLines(review)}
	} else {
		return ReviewResult{}
	}
}

func (this *PoetRule) ParseAnswer(text string) string {
	return text
}
