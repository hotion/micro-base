package log

import (
	"go.uber.org/zap"
)

// Log 实现go-kit日志接口
type Log struct {
	*zap.SugaredLogger
}

// NewLog 创建日志对象
func NewLog(logger *zap.SugaredLogger) *Log {
	return &Log{logger}
}

// Log 输出日志
func (l *Log) Log(keyvals ...interface{}) error {
	l.Error(keyvals...)
	return nil
}
