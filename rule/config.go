package rule

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

type RuleConfig struct {
	introduction          string `yaml:"introduction"`
	systemMessage         string `yaml:"system_message"`
	sourceFile            string `yaml:"source_file"`
	reviewerSystemMessage string `yaml:"reviewer_system_message"`
}

type RulesConfig struct {
	rules map[string]RuleConfig `yaml:"rules"`
}

func readConfig(filepath string) (map[string]RuleConfig, error) {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("读取配置文件失败: %v", err)
	}

	var config RulesConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("解析YAML失败: %v", err)
	}

	return config.rules, nil
}
