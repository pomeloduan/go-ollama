# Ollama Local Service Demo

A demonstration project for Ollama local service written in Go.

## Project Introduction

This project is a simple Go application that demonstrates how to interact with the Ollama local AI model service. Ollama is a tool for running large language models locally, and this project shows how to call these models through Go programs.

## Features

- Connect to local Ollama service
- Simple command-line interface
- Configurable model parameters
- Error handling and logging

## Prerequisites

Before running this project, please ensure:

1. **Install Go**: Version 1.21 or higher
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
info.log
error.log
```

## Related Links

- [Ollama Official Website](https://ollama.ai/)
- [Ollama GitHub](https://github.com/ollama/ollama)
- [Go Official Documentation](https://golang.org/doc/)

---

**Note**: This is a demonstration project suitable for learning and testing purposes.

# Ollama 本地服务 Demo

一个用 Go 语言编写的 Ollama 本地服务演示项目。

## 项目简介

本项目是一个简单的 Go 应用程序，演示了如何与 Ollama 本地 AI 模型服务进行交互。Ollama 是一个用于在本地运行大型语言模型的工具，本项目展示了如何通过 Go 程序调用这些模型。

## 功能特性

- 连接到本地 Ollama 服务
- 简单的命令行界面
- 可配置的模型参数
- 错误处理和日志记录

## 前提条件

在运行此项目之前，请确保：

1. **安装 Go**: 版本 1.21 或更高
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
info.log
error.log
```


## 相关链接

- [Ollama 官网](https://ollama.ai/)
- [Ollama GitHub](https://github.com/ollama/ollama)
- [Go 官方文档](https://golang.org/doc/)

---

**提示**: 这是一个演示项目，适用于学习和测试目的。
