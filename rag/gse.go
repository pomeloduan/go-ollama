package rag

import (
	"strings"
	"sync"

	"github.com/go-ego/gse"
)

// GseManager 中文分词管理器
// 使用 GSE 库对中文文本进行分词，提升向量化的准确性
type GseManager struct {
	seg gse.Segmenter // GSE 分词器实例
}

// newGseManager 创建并初始化中文分词管理器
// 加载中文词典（zh_s 表示简体中文）
func newGseManager() *GseManager {
	var seg gse.Segmenter
	seg.LoadDict("zh_s")
	return &GseManager{seg: seg}
}

// splitChineseWords 对中文文本进行分词
// 将分词结果用空格连接，便于后续向量化处理
// 参数 text: 待分词的中文文本
// 返回: 分词后的文本（词语间用空格分隔）
func (this *GseManager) splitChineseWords(text string) string {
	cut := this.seg.Cut(text)
	return strings.Join(cut, " ")
}
