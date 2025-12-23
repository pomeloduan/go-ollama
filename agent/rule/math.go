package rule

func init() {
	registerRule("math", &MathRule{})
}

type MathRule struct {
}

func (this *MathRule) Introduction() string {
	return "擅于解答数学问题，涉及代数、几何、概率等数学相关都可以来问。"
}

func (this *MathRule) SystemMessage() string {
	return `你是一位数学老师。你的任务是解答数学题。

# 行动格式:
你的回答必须严格遵循以下格式。首先是这是不是数学问题，然后是解答。
isMath: [这里是判断这是不是数学问题，使用true/false回答]
resolvation: [这里是你的解答]`
}

func (this *MathRule) SourceFile() string {
	return ""
}

func (this *MathRule) MessageFromSource(source string, question string) string {
	return ""
}

func (this *MathRule) ParseAnswer(text string) string {
	var isMath, resolvation, formatedAnswer = parseKeyValueText(text, "isMath", "resolvation")
	if formatedAnswer && isMath == "true" {
		return compactEmptyLines(resolvation)
	} else {
		return "不是数学问题"
	}
}
