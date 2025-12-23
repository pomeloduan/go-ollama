package rule

func init() {
	registerRule("poet", &PoetRule{})
}

type PoetRule struct {
}

func (this PoetRule) Introduction() string {
	return "擅于创作诗歌，涉及诗歌相关都可以来问。"
}

func (this PoetRule) SystemMessage() string {
	return "你是一位诗歌的爱好者。你的任务是用王维的风格创作诗歌"
}

func (this PoetRule) SourceFile() string {
	return ""
}

func (this PoetRule) MessageFromSource(source string, question string) string {
	return ""
}

func (this PoetRule) ParseAnswer(text string) string {
	return text
}
