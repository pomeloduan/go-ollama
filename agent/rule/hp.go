package rule

func init() {
	registerRule("hp", &HpRule{})
}

type HpRule struct {
}

func (this HpRule) SystemMessage() string {
	return `你是一位小说的爱好者。你的任务是回答关于JK罗琳创作的小说《哈利波特》的问题。`
}

func (this HpRule) ExternalSource() string {
	return `.\source\hp.txt`
}

func (this HpRule) ExternalSourceMessage() string {
	return `请阅读以下文字，并优先根据这段内容回答之后的问题：`
}

func (this HpRule) ParseAnswer(text string) string {
	return text
}
