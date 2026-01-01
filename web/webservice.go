package web

import (
	"encoding/json"
	"fmt"
	"net/http"

	"go-ollama/agent"
	"go-ollama/ollama"
)

// ChatRequest 聊天请求结构
type ChatRequest struct {
	Message string `json:"message"`
}

// ChatResponse 聊天响应结构
type ChatResponse struct {
	Answer string `json:"answer"`
	Error  string `json:"error,omitempty"`
}

// StatsResponse 统计信息响应结构
type StatsResponse struct {
	QuestionCount int     `json:"question_count"`
	AnswerCount   int     `json:"answer_count"`
	TotalDuration float64 `json:"total_duration"`
	TotalToken    int     `json:"total_token"`
}

// WebService Web服务，包含所有需要的依赖
type WebService struct {
	agentMgr  agent.AgentManager
	ollamaMgr ollama.OllamaManager
}

// NewWebService 创建Web服务实例
// 参数 agentMgr: Agent管理器
// 参数 ollamaMgr: Ollama管理器
// 返回: WebService实例
func NewWebService(agentMgr agent.AgentManager, ollamaMgr ollama.OllamaManager) *WebService {
	return &WebService{
		agentMgr:  agentMgr,
		ollamaMgr: ollamaMgr,
	}
}

// HandleIndex 处理首页请求，返回HTML页面
func (ws *WebService) HandleIndex(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, IndexHTML)
}

// HandleChat 处理聊天API请求
func (ws *WebService) HandleChat(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req ChatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response := ChatResponse{Error: "无效的请求格式"}
		json.NewEncoder(w).Encode(response)
		return
	}

	if req.Message == "" {
		response := ChatResponse{Error: "消息不能为空"}
		json.NewEncoder(w).Encode(response)
		return
	}

	// 调用Agent处理问题
	answer := ws.agentMgr.Chat(req.Message)

	response := ChatResponse{Answer: answer}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(response)
}

// HandleStats 处理统计信息API请求
func (ws *WebService) HandleStats(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	stats := StatsResponse{
		QuestionCount: ws.ollamaMgr.GetTotalQCount(),
		AnswerCount:   ws.ollamaMgr.GetTotalACount(),
		TotalDuration: ws.ollamaMgr.GetTotalDuration().Seconds(),
		TotalToken:    ws.ollamaMgr.GetTotalToken(),
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(stats)
}

// RegisterRoutes 注册所有HTTP路由
// 参数 mux: HTTP多路复用器，如果为nil则使用默认的http.DefaultServeMux
func (ws *WebService) RegisterRoutes(mux *http.ServeMux) {
	if mux == nil {
		mux = http.DefaultServeMux
	}
	
	mux.HandleFunc("/", ws.HandleIndex)
	mux.HandleFunc("/api/chat", ws.HandleChat)
	mux.HandleFunc("/api/stats", ws.HandleStats)
}

