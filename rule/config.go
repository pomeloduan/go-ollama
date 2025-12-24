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
	ReviewerSystemMessage string `yaml:"reviewer_system_message"`
}

type RulesConfig struct {
	Rules map[string]RuleConfig `yaml:"rules"`
}

func readConfig(filepath string) (map[string]RuleConfig, error) {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("读取配置文件失败: %v", err)
	}

	var config RulesConfig
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, fmt.Errorf("解析YAML失败: %v", err)
	}

	return config.Rules, nil
}
