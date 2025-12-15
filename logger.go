package main

import (
	"fmt"
	"os"
	"time"
)

// ErrorLogger 错误日志记录器
type ErrorLogger struct {
	file *os.File
}

// NewErrorLogger 创建新的错误日志记录器
func NewErrorLogger(filename string) (*ErrorLogger, error) {
	// 打开或创建日志文件（追加模式）
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("无法打开日志文件: %v", err)
	}

	return &ErrorLogger{file: file}, nil
}

// LogError 记录错误到文件
func (this *ErrorLogger) LogError(err error, context ...string) error {
	timestamp := time.Now().Format("2006-01-02 15:04:05")

	// 构建日志条目
	logEntry := fmt.Sprintf("[ERROR][%s] %v", timestamp, err)

	// 添加上下文信息
	if len(context) > 0 {
		logEntry += fmt.Sprintf(" | 上下文: %s", context[0])
		for i := 1; i < len(context); i++ {
			logEntry += fmt.Sprintf(", %s", context[i])
		}
	}

	logEntry += "\n"

	// 写入文件
	if _, writeErr := this.file.WriteString(logEntry); writeErr != nil {
		return fmt.Errorf("写入日志失败: %v", writeErr)
	}

	// 确保立即写入磁盘
	return this.file.Sync()
}

func (this *ErrorLogger) LogInfo(info string) error {
	timestamp := time.Now().Format("2006-01-02 15:04:05")

	// 构建日志条目
	logEntry := fmt.Sprintf("[info][%s] %v", timestamp, info)

	logEntry += "\n"

	// 写入文件
	if _, writeErr := this.file.WriteString(logEntry); writeErr != nil {
		return fmt.Errorf("写入日志失败: %v", writeErr)
	}

	// 确保立即写入磁盘
	return this.file.Sync()
}

// Close 关闭日志文件
func (this *ErrorLogger) Close() error {
	if this.file != nil {
		return this.file.Close()
	}
	return nil
}
