While most AI Agent implementations are in Python, practical examples using Go are relatively scarce. This project is a complete, locally runnable AI Agent system built with Go. It serves as:
- A record of my personal learning journey
- A reference for AI development within the Go community
- A space to exchange ideas with others interested in building AI Agents and learning Go language

# Ollama Local Service Demo

This is a Go sample project that demonstrates how to call a local Ollama AI model with minimal code.

It implements three core functions:

- **Prompt Management**: Optimize instructions for model interaction  
- **Context Handling**: Maintain coherence across multi-turn conversations  
- **RAG Integration**: Enhance response quality through retrieval-augmented generation

rag效果（需要导入《哈利波特与魔法石》16章内容）
```
哈利和罗恩认为谁想去偷魔法石？
哈利和罗恩认为是斯内普想去偷魔法石。  
哈利在讲述时明确表示：“所以我现在只能等着斯内普去偷魔法石。”罗恩和赫敏虽然没有直接重复这一观点，但从对话中可以推断他们也认同斯内普的意图。
等待提问 -->
下棋的过程中，哈利代替了什么棋子，罗恩代替了什么棋子？
根据第二段文字，哈利代替了**主教**的位置，罗恩代替了**骑士**。  
赫敏则代替了城堡的位置，但问题仅询问哈利和罗恩。
等待提问 -->
```

## Features

- Connect to local Ollama service
- Simple command-line interface
- Configurable model parameters
- Error handling and logging

## Prerequisites

Before running this project, please ensure:

1. **Install Go**: Version 1.23 or higher
   ```bash
   # Check Go version
   go version
   ```

2. **Install and run Ollama**:
   - Visit the [Ollama official website](https://ollama.ai/) to download and install
   - Start the Ollama service
   - Download at least one model, for example:
     ```bash
     ollama pull llama2
     ollama pull mistral
     ```

## Important Notes

1. **Ensure the Ollama service is running**:
   ```bash
   # Check Ollama service status
   curl http://localhost:11434/api/tags
   ```

2. **Model availability**: Ensure the required models are downloaded and available

3. **Performance considerations**: Response times may vary depending on model size and hardware configuration

## Common Issues

1. **Connection failure**:
   - Check if the Ollama service is running: `ollama serve`
   - Confirm that port 11434 is available

2. **Model not found**:
   - List installed models: `ollama list`
   - Download the required model: `ollama pull <model-name>`

3. **Insufficient memory**:
   - Try using smaller models
   - Close other memory-intensive applications

## Log Viewing

```
error.log
```

## Related Links

- [Ollama Official Website](https://ollama.ai/)
- [Ollama GitHub](https://github.com/ollama/ollama)
- [Go Official Documentation](https://golang.org/doc/)

**Note**: This is a demonstration project suitable for learning and testing purposes.

---

当前AI Agent的实现多以Python为主，Go语言的相关实践较为缺乏。本项目旨在使用Go语言完整实现一套可本地运行的AI Agent系统，希望能够：
- 作为个人学习总结
- 为Go社区提供AI开发的参考示例
- 与其他关注AI agent开发和Go语言的小伙伴交流

# Ollama 本地服务 Demo

这是一个 Go 语言演示项目，展示了如何用最少代码调用本地 Ollama AI 模型。

项目实现三个基本功能：

- **提示词管理**：优化与模型的对话指令
- **上下文维护**：保持多轮对话连贯性
- **RAG 集成**：通过检索增强生成提升回答质量

## 功能特性

- 连接到本地 Ollama 服务
- 简单的命令行界面
- 可配置的模型参数
- 错误处理和日志记录

## 前提条件

在运行此项目之前，请确保：

1. **安装 Go**: 版本 1.23 或更高
   ```bash
   # 检查 Go 版本
   go version
   ```

2. **安装并运行 Ollama**: 
   - 访问 [Ollama 官网](https://ollama.ai/) 下载并安装
   - 启动 Ollama 服务
   - 下载至少一个模型，例如：
     ```bash
     ollama pull llama2
     ollama pull mistral
     ```

## 注意事项

1. **确保 Ollama 服务正在运行**：
   ```bash
   # 检查 Ollama 服务状态
   curl http://localhost:11434/api/tags
   ```

2. **模型可用性**：确保所需的模型已下载并可用

3. **性能考虑**：根据模型大小和硬件配置，响应时间可能有所不同

## 常见问题

1. **连接失败**：
   - 检查 Ollama 服务是否运行：`ollama serve`
   - 确认端口 11434 是否可用

2. **模型未找到**：
   - 列出已安装的模型：`ollama list`
   - 下载所需模型：`ollama pull <model-name>`

3. **内存不足**：
   - 尝试使用较小的模型
   - 关闭其他占用内存的应用程序

## 日志查看

```
error.log
```


## 相关链接

- [Ollama 官网](https://ollama.ai/)
- [Ollama GitHub](https://github.com/ollama/ollama)
- [Go 官方文档](https://golang.org/doc/)

**提示**: 这是一个演示项目，适用于学习和测试目的。
