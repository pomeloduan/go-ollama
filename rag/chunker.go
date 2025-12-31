package rag

import (
	"bufio"
	"io"
	"os"
	"strings"
)

// minParaCount 每个文本块最少包含的段落数
const minParaCount = 2
// minCharCount 每个文本块最少的字符数
const minCharCount = 100

// chucksFromTextFile 从文本文件读取并分块
// 按段落读取文件，然后组合成合适大小的文本块
// 参数 filePath: 文件路径
// 返回: 文本块数组、error
func chucksFromTextFile(filePath string) ([]string, error) {
	paragraphs, err := readParagraphs(filePath)
	if err != nil {
		return nil, err
	}
	chucks := chucksFromParagraphs(paragraphs)
	return chucks, nil
}

// func chucksFromText(text string) ([]string, error) {
// 	paragraphs, err := readParagraphs(filePath)
// 	if err != nil {
// 		return nil, err
// 	}
// 	chucks := chucksFromParagraphs(paragraphs)
// 	return chucks, nil
// }

// readParagraphs 从文件读取段落
// 按空行分割段落，连续的非空行组成一个段落
// 参数 filePath: 文件路径
// 返回: 段落数组、error
func readParagraphs(filePath string) ([]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	var paragraphs []string
	var builder strings.Builder

	for {
		line, err := reader.ReadString('\n')

		// 去除换行符
		line = strings.TrimSuffix(line, "\n")
		line = strings.TrimSuffix(line, "\r")

		// 检查是否为空行
		if strings.TrimSpace(line) == "" {
			if builder.Len() > 0 {
				paragraphs = append(paragraphs, builder.String())
				builder.Reset()
			}
		} else {
			if builder.Len() > 0 {
				builder.WriteString("\n")
			}
			builder.WriteString(line)
		}

		// 处理文件结束
		if err != nil {
			if err == io.EOF {
				if builder.Len() > 0 {
					paragraphs = append(paragraphs, builder.String())
				}
				break
			}
			return nil, err
		}
	}

	return paragraphs, nil
}

// chucksFromParagraphs 将段落组合成文本块
// 策略：至少包含 minParaCount 个段落且字符数 >= minCharCount 才形成一个块
// 这样可以保证每个块有足够的上下文信息，同时避免块过大
// 参数 paragraphs: 段落数组
// 返回: 文本块数组
func chucksFromParagraphs(paragraphs []string) []string {
	var chucks []string
	var chuck string
	var paraCount int

	for i := 0; i < len(paragraphs); i++ {
		para := paragraphs[i]
		chuck += para
		paraCount++

		// 当达到最小段落数和最小字符数时，形成一个文本块
		if paraCount >= minParaCount && len(chuck) >= minCharCount {
			chucks = append(chucks, chuck)

			// 重置，开始下一个块
			// 注意：这里采用重置策略而不是滑动窗口
			// 可以根据需要修改为：chuck = para; paraCount = 1（保留最后一段）

			chuck = ""
			paraCount = 0
		}
		// 内容较少时，继续累积到下一个块
	}

	// 处理剩余的段落（至少 2 段且有内容）
	if paraCount > 1 && len(chuck) > 0 {
		chucks = append(chucks, chuck)
	}

	return chucks
}
