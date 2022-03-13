package log

import (
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func init() {
	config := zap.NewProductionConfig()
	config.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	config.OutputPaths = []string{"stdout", "/tmp/coordinator.log"}
	
	logger, err := config.Build()
	if err != nil {
		panic(err)
	}

	zap.ReplaceGlobals(logger)
}

func Info(msg string, fields ...zapcore.Field) {
	zap.L().Info(msg, fields...)
}

func Infof(msg string, args ...interface{}) {
	zap.L().Info(fmt.Sprintf(msg, args...))
}

func Warn(msg string, fields ...zapcore.Field) {
	zap.L().Warn(msg, fields...)
}

func Warnf(msg string, args ...interface{}) {
	zap.L().Warn(fmt.Sprintf(msg, args...))
}

func Error(msg string, fields ...zapcore.Field) {
	zap.L().Error(msg, fields...)
}

func Errorf(msg string, args ...interface{}) {
	zap.L().Error(fmt.Sprintf(msg, args...))
}
