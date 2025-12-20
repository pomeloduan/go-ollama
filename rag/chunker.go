package rag

import (
	"bufio"
	"io"
	"os"
	"strings"
)

const minParaCount = 2
const minCharCount = 100

// 文本分块（按段落）
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

func chucksFromParagraphs(paragraphs []string) []string {
	var chucks []string
	var chuck string
	var paraCount int

	for i := 0; i < len(paragraphs); i++ {
		para := paragraphs[i]
		chuck += para
		paraCount++

		// 组合一个chuck
		if paraCount >= minParaCount && len(chuck) >= minCharCount {
			chucks = append(chucks, chuck)

			// chuck = para
			// paraCount = 1

			chuck = ""
			paraCount = 0
		}
		// else 内容较少，继续加入内容
	}

	if paraCount > 1 && len(chuck) > 0 {
		chucks = append(chucks, chuck)
	}

	return chucks
}
