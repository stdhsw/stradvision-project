package logger

import (
	"os"
	"testing"

	"go.uber.org/zap"
)

func Test_logger(t *testing.T) {
	nowPath, _ := os.Getwd()

	if err := InitLogger(
		"stradvision",
		WithPath(nowPath),
		WithLogLevel("debug"),
		WithLogLocalTime(false),
		WithLogCompress(true),
		WitchEncoder("console"),
	); err != nil {
		t.Error(err)
	}

	writer.Info("info message")
	writer.Debug("debug message", zap.String("key", "value"))
	writer.Warn("warn message", zap.String("key", "value"))
	writer.Error("error message", zap.String("key", "value"))
}
