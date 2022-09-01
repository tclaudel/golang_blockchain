package badger

import (
	badger "github.com/dgraph-io/badger/v3"
	"go.uber.org/zap"
)

const badgerLog = "[BADGER] : "

type Logger struct {
	logger       *zap.Logger
	badgerLogger bool
}

func (b Logger) Errorf(s string, i ...interface{}) {
	b.logger.Error(s, zap.Any("i", i), zap.String("repository", "badger"))
}

func (b Logger) Warningf(s string, i ...interface{}) {
	b.logger.Warn(s, zap.Any("i", i), zap.String("repository", "badger"))
}

func (b Logger) Infof(s string, i ...interface{}) {
	b.logger.Info(s, zap.Any("i", i), zap.String("repository", "badger"))
}

func (b Logger) Debugf(s string, i ...interface{}) {
	b.logger.Debug(s, zap.Any("i", i), zap.String("repository", "badger"))
}

func NewBadgerLogger(logger *zap.Logger) badger.Logger {
	return &Logger{
		logger: logger,
	}
}
