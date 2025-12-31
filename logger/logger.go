package logger

import (
	"fmt"
	"os"
	"sync"
	"time"
)

// ErrorLogger 错误日志记录器接口
type ErrorLogger interface {
	LogError(err error, context ...string) error
	LogInfo(info string) error
	Close() error
}

// errorLogger 错误日志记录器实现（包私有）
type errorLogger struct {
	file *os.File
}

var (
	loggerInstance *errorLogger
	loggerOnce     sync.Once
)

// newErrorLogger 创建并初始化错误日志记录器实例
// 参数 filename: 日志文件路径
// 返回: errorLogger 实例、error
func newErrorLogger(filename string) (*errorLogger, error) {
	// 打开或创建日志文件（追加模式）
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("无法打开日志文件: %v", err)
	}

	return &errorLogger{file: file}, nil
}

// StartErrorLogger 获取错误日志记录器单例
func StartErrorLogger(filename string) (ErrorLogger, error) {
	var err error
	loggerOnce.Do(func() {
		loggerInstance, err = newErrorLogger(filename)
	})

	if err != nil {
		return nil, err
	}
	return loggerInstance, nil
}

// LogError 记录错误到文件
func (e *errorLogger) LogError(err error, context ...string) error {
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
	if _, writeErr := e.file.WriteString(logEntry); writeErr != nil {
		return fmt.Errorf("写入日志失败: %v", writeErr)
	}

	// 确保立即写入磁盘
	return e.file.Sync()
}

// LogInfo 记录信息到文件
func (e *errorLogger) LogInfo(info string) error {
	timestamp := time.Now().Format("2006-01-02 15:04:05")

	// 构建日志条目
	logEntry := fmt.Sprintf("[info][%s] %v", timestamp, info)

	logEntry += "\n"

	// 写入文件
	if _, writeErr := e.file.WriteString(logEntry); writeErr != nil {
		return fmt.Errorf("写入日志失败: %v", writeErr)
	}

	// 确保立即写入磁盘
	return e.file.Sync()
}

// Close 关闭日志文件
func (e *errorLogger) Close() error {
	if e.file != nil {
		return e.file.Close()
	}
	return nil
}
