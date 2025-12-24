[中文](README.md) | English

While most AI Agent implementations are in Python, practical examples using Go are relatively scarce. This project is a complete, locally runnable AI Agent system built with Go. It serves as:
1. A record of my personal learning journey
2. A reference for AI development within the Go community
3. A space to exchange ideas with others interested in building AI Agents and learning Go language

# Ollama Local Service Demo

This is a Go sample project that demonstrates how to call a local Ollama AI model with minimal code.

This project builds an agent system with the following features:

- **Prompt Management**: Optimizes dialogue instructions for the model
- **Context Maintenance**: Maintains coherence across multi-turn conversations
- **RAG Integration**: Incorporates Retrieval-Augmented Generation technology to enhance the accuracy of professional question answering
- **Multi-Agent Collaboration**: Employs a centralized architecture where a **Coordinator** allocates tasks and **Specialist** agents handle domain-specific problems
- **Reviewer Mode**: Adheres to a Generate → Review → Refine workflow to ensure continuous optimization of output quality
- **Agent Configuration**: Uses YAML configuration files to dynamically define and generate **Specialist** agents and **Reviewer** agents

```
. . . . . . . . . . . . . . . . . . .
.
. Agent Manager
.                 │
.          ┌──────▼──────┐
.          │ Coordinator │
.          │ ---         │
.          │ Analysis    │
.          │ Matching    │
.          └──────┬──────┘
.                 │
.      ┌──────────┴─────────┐
       │                    │
┌──────▼──────┐       ┌─────▼─────┐
│ Specialist  │       │ Reviewer  │
│ ---         │◄──┐   │ ---       │
│ - hp        │   └──►│ Scoring   │       ┌─────────────┐
│ - math      │-------│ Feedback  │------►│ Rag Manager │
│ - poet      │       └─────┬─────┘       │ ---         │
│ ...         │             │             │ Chucking    │
└──────┬──────┘             │             │     ↓       │
       │                    │             │ Embedding + │
       └──────────┬─────────┘             │     ↓     +-│--→ Gse Segmentation
                  │                       │ Retrieval + │    Chromem Vector DB
         ┌────────▼────────┐              │     ↓       │    ┌──────────┐
         │ Ollama Manager  │              │ Reranking +-│---►│ Reranker │
         │ ---             │              └─────────────┘    └─────┬────┘
         │ Ollama Service  │◄──────────────────────────────────────┘
         │ Local LLM       │
         └─────────────────┘

```

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
   - Download at least one language model, for example:
     ```bash
     ollama pull llama2
     ollama pull mistral
     ```
   - Download at least one embedding model

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
```

## Related Links

- [Ollama Official Website](https://ollama.ai/)
- [Ollama GitHub](https://github.com/ollama/ollama)
- [Go Official Documentation](https://golang.org/doc/)

**Note**: This is a demonstration project suitable for learning and testing purposes.
