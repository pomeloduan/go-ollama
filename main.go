package main

import (
	"fmt"
	"log"
	"net/http"

	"go-ollama/agent"
	"go-ollama/logger"
	"go-ollama/ollama"
	"go-ollama/web"
)

// ollamaDomain Ollama 服务的地址，默认本地 11434 端口
const ollamaDomain = "http://localhost:11434"

// serverAddr Web服务器监听地址
const serverAddr = ":8080"

// test
// agent.Chat("1+2+3+...+100的值是多少？")
// agent.Chat("请你以猫为主题，写一首诗。")
// agent.Chat("在《哈利波特》中，为了通过下棋关卡，哈利代替了什么棋子，罗恩代替了什么棋子？")
// agent.Chat("今天天气怎么样")

// main 程序入口函数
// 初始化日志、Ollama 连接和 Agent 管理器，然后启动Web服务器
// todo mcp func call
func main() {
	fmt.Println("--> Ollama Local Service Demo")
	fmt.Println("正在初始化...")

	// 创建错误日志
	errorLog, err := logger.NewErrorLogger("info.log")
	if err != nil {
		log.Fatal(err)
	}
	defer errorLog.Close()

	// 连接本地模型
	ollamaMgr, err := ollama.StartOllamaManager(ollamaDomain, errorLog)
	if err != nil {
		errorLog.LogError(err, "launching")
		return
	}

	// 启动agent
	agentMgr, err := agent.StartAgentManager(ollamaMgr, errorLog)
	if err != nil {
		errorLog.LogError(err, "launching")
		return
	}

	// 创建Web服务并注册路由
	webService := web.NewWebService(agentMgr, ollamaMgr)
	webService.RegisterRoutes(nil)

	// 启动Web服务器
	fmt.Printf("Web服务已启动，请访问: http://localhost%s\n", serverAddr)
	fmt.Println("按 Ctrl+C 停止服务")
	if err := http.ListenAndServe(serverAddr, nil); err != nil {
		errorLog.LogError(err, "http server")
		log.Fatal(err)
	}
}
