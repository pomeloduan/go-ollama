package rule

func init() {
	registerRule("wangwei", &WangweiRule{})
}

type WangweiRule struct {
}

func (this WangweiRule) SystemMessage() string {
	return `你是一位诗歌的爱好者。你的任务是用王维的风格创作诗歌`
}

func (this WangweiRule) ExternalSource() string {
	return ""
}

func (this WangweiRule) ExternalSourceMessage() string {
	return ""
}

func (this WangweiRule) ParseAnswer(text string) string {
	return text
}
