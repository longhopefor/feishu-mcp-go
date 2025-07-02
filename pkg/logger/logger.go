package logger

import (
	"os"
	"strings"

	"github.com/sirupsen/logrus"
)

// Logger 日志接口
type Logger interface {
	Debug(msg string, args ...interface{})
	Info(msg string, args ...interface{})
	Warn(msg string, args ...interface{})
	Error(msg string, args ...interface{})
	Fatal(msg string, args ...interface{})
}

// LogrusLogger logrus实现的日志器
type LogrusLogger struct {
	logger *logrus.Logger
}

// New 创建新的日志器
func New(level string) Logger {
	logger := logrus.New()

	// 设置日志级别
	switch strings.ToLower(level) {
	case "debug":
		logger.SetLevel(logrus.DebugLevel)
	case "info":
		logger.SetLevel(logrus.InfoLevel)
	case "warn":
		logger.SetLevel(logrus.WarnLevel)
	case "error":
		logger.SetLevel(logrus.ErrorLevel)
	default:
		logger.SetLevel(logrus.InfoLevel)
	}

	// 设置日志格式
	logger.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})

	// 设置输出目标
	logger.SetOutput(os.Stdout)

	return &LogrusLogger{logger: logger}
}

// Debug 记录调试信息
func (l *LogrusLogger) Debug(msg string, args ...interface{}) {
	l.logWithFields(logrus.DebugLevel, msg, args...)
}

// Info 记录一般信息
func (l *LogrusLogger) Info(msg string, args ...interface{}) {
	l.logWithFields(logrus.InfoLevel, msg, args...)
}

// Warn 记录警告信息
func (l *LogrusLogger) Warn(msg string, args ...interface{}) {
	l.logWithFields(logrus.WarnLevel, msg, args...)
}

// Error 记录错误信息
func (l *LogrusLogger) Error(msg string, args ...interface{}) {
	l.logWithFields(logrus.ErrorLevel, msg, args...)
}

// Fatal 记录致命错误并退出程序
func (l *LogrusLogger) Fatal(msg string, args ...interface{}) {
	l.logWithFields(logrus.FatalLevel, msg, args...)
	os.Exit(1)
}

// logWithFields 记录日志并处理字段
func (l *LogrusLogger) logWithFields(level logrus.Level, msg string, args ...interface{}) {
	fields := logrus.Fields{}

	// 将参数转换为字段
	for i := 0; i < len(args); i += 2 {
		if i+1 < len(args) {
			key := args[i]
			value := args[i+1]
			if keyStr, ok := key.(string); ok {
				fields[keyStr] = value
			}
		}
	}

	entry := l.logger.WithFields(fields)

	switch level {
	case logrus.DebugLevel:
		entry.Debug(msg)
	case logrus.InfoLevel:
		entry.Info(msg)
	case logrus.WarnLevel:
		entry.Warn(msg)
	case logrus.ErrorLevel:
		entry.Error(msg)
	case logrus.FatalLevel:
		entry.Fatal(msg)
	}
}

// NewNoop 创建一个无操作的日志器（用于测试）
func NewNoop() Logger {
	return &NoopLogger{}
}

// NoopLogger 无操作日志器
type NoopLogger struct{}

func (n *NoopLogger) Debug(msg string, args ...interface{}) {}
func (n *NoopLogger) Info(msg string, args ...interface{})  {}
func (n *NoopLogger) Warn(msg string, args ...interface{})  {}
func (n *NoopLogger) Error(msg string, args ...interface{}) {}
func (n *NoopLogger) Fatal(msg string, args ...interface{}) {
	os.Exit(1)
}
