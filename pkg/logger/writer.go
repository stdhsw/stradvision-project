package logger

import "go.uber.org/zap"

var (
	writer *zap.Logger
)

func Debug(msg string, fields ...zap.Field) {
	writer.Debug(msg, fields...)
}

func Info(msg string, fields ...zap.Field) {
	writer.Info(msg, fields...)
}

func Warn(msg string, fields ...zap.Field) {
	writer.Warn(msg, fields...)
}

func Error(msg string, fields ...zap.Field) {
	writer.Error(msg, fields...)
}

func Fatal(msg string, fields ...zap.Field) {
	writer.Fatal(msg, fields...)
}

func Panic(msg string, fields ...zap.Field) {
	writer.Panic(msg, fields...)
}
