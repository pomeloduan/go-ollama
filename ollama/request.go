package ollama

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// ChatRequest Ollama API 聊天请求结构
type ChatRequest struct {
	Model    string        `json:"model"`    // 模型名称
	Messages []ChatMessage `json:"messages"` // 消息列表
	Stream   bool          `json:"stream"`   // 是否流式输出（当前未使用）
}

// ChatMessage 对话消息结构
type ChatMessage struct {
	Role    string `json:"role"`    // 角色：system/user/assistant
	Content string `json:"content"` // 消息内容
}

// ChatResponse Ollama API 聊天响应结构
type ChatResponse struct {
	Model              string      `json:"model"`                         // 使用的模型
	CreatedAt          string      `json:"created_at"`                    // 创建时间
	Message            ChatMessage `json:"message"`                       // 返回的消息
	Done               bool        `json:"done"`                          // 是否完成
	DoneReason         string      `json:"done_reason"`                   // 完成原因
	TotalDuration      int64       `json:"total_duration,omitempty"`      // 总耗时（纳秒）
	LoadDuration       int64       `json:"load_duration,omitempty"`       // 模型加载耗时
	PromptEvalCount    int         `json:"prompt_eval_count,omitempty"`   // 提示词 token 数
	PromptEvalDuration int64       `json:"prompt_eval_duration,omitempty"`// 提示词评估耗时
	EvalCount          int         `json:"eval_count,omitempty"`          // 生成 token 数
	EvalDuration       int64       `json:"eval_duration,omitempty"`       // 生成耗时
}

// listModels 列出本地 Ollama 服务可用的所有模型
// 参数 domain: Ollama 服务地址
// 返回: 模型名称数组、error
func listModels(domain string) ([]string, error) {
	resp, err := http.Get(domain + "/api/tags")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	models := []string{}
	if modelsList, ok := result["models"].([]interface{}); ok {
		for _, m := range modelsList {
			if modelMap, ok := m.(map[string]interface{}); ok {
				if name, ok := modelMap["name"].(string); ok {
					models = append(models, name)
				}
			}
		}
	}

	return models, nil
}

// sendChatRequest 发送聊天请求到 Ollama API
// 参数 domain: Ollama 服务地址
// 参数 model: 模型名称
// 参数 messages: 消息列表
// 返回: ChatResponse、error
func sendChatRequest(domain string, model string, messages []ChatMessage) (*ChatResponse, error) {
	requestData := ChatRequest{
		Model:    model,
		Messages: messages,
		Stream:   false,
	}

	jsonData, err := json.Marshal(requestData)
	if err != nil {
		return nil, fmt.Errorf("json error: %v", err)
	}

	client := &http.Client{
		Timeout: 180 * time.Second,
	}

	resp, err := client.Post(
		domain+"/api/chat",
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return nil, fmt.Errorf("http request error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("api error: %s - %s", resp.Status, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read resp error: %v", err)
	}

	var chatResp ChatResponse
	if err := json.Unmarshal(body, &chatResp); err != nil {
		return nil, fmt.Errorf("json error: %v", err)
	}

	return &chatResp, nil
}
