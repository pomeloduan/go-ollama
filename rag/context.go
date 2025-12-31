package rag

// RagContext RAG 上下文，存储知识库的相关信息
type RagContext struct {
	ragId  int      // RAG 上下文 ID，对应向量数据库中的集合 ID
	chucks []string // 原始文本块数组，用于根据索引检索完整文本
}
