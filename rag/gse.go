package rag

import (
	"strings"

	"github.com/go-ego/gse"
)

// 中文分词
type GseManager struct {
	seg gse.Segmenter
}

func startGse() *GseManager {
	var seg gse.Segmenter
	seg.LoadDict("zh_s")
	gseManager := GseManager{seg: seg}
	return &gseManager
}

func (this *GseManager) splitChineseWords(text string) string {
	cut := this.seg.Cut(text)
	return strings.Join(cut, " ")
}
