package log

import (
	"errors"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const stdout = "stdout"

func New(logFormat, logLevel string) (*zap.Logger, error) {
	var zapLogger *zap.Logger

	zapLogLevel, err := getLogLevel(logLevel)
	if err != nil {
		return nil, err
	}

	if logFormat == "json" {
		zapConfig := zap.NewProductionConfig()
		zapConfig.Level = zap.NewAtomicLevelAt(zapLogLevel)
		zapConfig.ErrorOutputPaths = []string{stdout}
		zapConfig.OutputPaths = []string{stdout}
		zapLogger, err = zapConfig.Build()
		if err != nil {
			return nil, err
		}
	} else {
		zapConfig := zap.NewDevelopmentConfig()
		zapConfig.Level = zap.NewAtomicLevelAt(zapLogLevel)
		zapConfig.ErrorOutputPaths = []string{stdout}
		zapConfig.OutputPaths = []string{stdout}
		zapLogger, err = zapConfig.Build()
		if err != nil {
			return nil, err
		}
	}

	return zapLogger, nil
}

var (
	ErrUnrecognizedLogLevel = errors.New("unrecognized log level")
)

func getLogLevel(level string) (zapcore.Level, error) {
	switch level {
	case "debug", "DEBUG":
		return zapcore.DebugLevel, nil
	case "info", "INFO", "": // Make the zero value useful.
		return zapcore.InfoLevel, nil
	case "warn", "WARN":
		return zapcore.WarnLevel, nil
	case "error", "ERROR":
		return zapcore.ErrorLevel, nil
	case "dpanic", "DPANIC":
		return zapcore.DPanicLevel, nil
	case "panic", "PANIC":
		return zapcore.PanicLevel, nil
	case "fatal", "FATAL":
		return zapcore.FatalLevel, nil
	default:
		return zapcore.ErrorLevel, ErrUnrecognizedLogLevel
	}
}

func Error(logger *zap.Logger, err error) error {
	logger.WithOptions(zap.AddCallerSkip(1)).Error(err.Error())
	return err
}
