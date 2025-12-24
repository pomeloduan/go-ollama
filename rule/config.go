package rule

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

type RuleConfig struct {
	Introduction          string `yaml:"introduction"`
	SystemMessage         string `yaml:"system_message"`
	SourceFile            string `yaml:"source_file"`
	SourceMessage         string `yaml:"source_message"`
	ReviewerSystemMessage string `yaml:"reviewer_system_message"`
	ReviewMessage         string `yaml:"review_message"`
	RewriteMessage        string `yaml:"rewrite_message"`
}

type ChatConfig struct {
	Rules map[string]RuleConfig `yaml:"rules"`

	RerankMessage                string `yaml:"rerank_message"`
	CoordinatorMessage           string `yaml:"coordinator_message"`
	CoordinatorSpecialistMessage string `yaml:"coordinator_specialist_message"`
}

func readConfig(filepath string) (*ChatConfig, error) {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("读取配置文件失败: %v", err)
	}

	var config ChatConfig
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, fmt.Errorf("解析YAML失败: %v", err)
	}

	return &config, nil
}
