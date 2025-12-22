package agent

type CoordinateAgent struct {
}

func DefaultModelName() string {
	return "deepseek"
}

func Chat(chat string) string
