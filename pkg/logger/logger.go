package logger

import (
	"fmt"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func InitLogger(appName string, opts ...Option) error {
	c := fromOptions(appName, opts...)
	writer = zap.New(
		zapcore.NewCore(
			c.encoder,
			zapcore.NewMultiWriteSyncer(append([]zapcore.WriteSyncer{zapcore.AddSync(os.Stdout)}, zapcore.AddSync(&c.logger))...),
			c.level,
		),
	)

	if writer == nil {
		return fmt.Errorf("failed to create logger")
	}

	return nil
}
