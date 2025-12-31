package rule

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

// RuleConfig 单个规则的配置结构
// 对应 YAML 配置文件中 rules 下的单个规则
type RuleConfig struct {
	Introduction          string `yaml:"introduction"`           // 专家介绍，用于协调者匹配
	SystemMessage         string `yaml:"system_message"`         // 专家系统提示词
	SourceFile            string `yaml:"source_file"`            // RAG 源文件路径
	SourceMessage         string `yaml:"source_message"`         // RAG 检索文档的提示词模板
	ReviewerSystemMessage string `yaml:"reviewer_system_message"`// 评审者系统提示词
	ReviewMessage         string `yaml:"review_message"`         // 评审提示词模板
	RewriteMessage        string `yaml:"rewrite_message"`        // 重写提示词模板
}

// ChatConfig 完整的配置结构
// 对应整个 YAML 配置文件
type ChatConfig struct {
	Rules map[string]RuleConfig `yaml:"rules"` // 规则字典，key 是规则名称

	// 全局配置
	RerankMessage                string `yaml:"rerank_message"`                 // 重排提示词模板
	CoordinatorMessage           string `yaml:"coordinator_message"`            // 协调者提示词模板
	CoordinatorSpecialistMessage string `yaml:"coordinator_specialist_message"` // 协调者专家信息提示词模板
}

// readConfig 从 YAML 文件读取配置
// 参数 filepath: 配置文件路径
// 返回: ChatConfig 实例、error
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
