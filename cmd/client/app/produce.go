package app

import (
	"time"

	"example.com/stradvision-project/pkg/logger"
	"go.uber.org/zap"
)

func kafkaErrorHandler(ts time.Time, topic string, partition int32, err error) {
	logger.Error("failed kafka send message", zap.Time("ts", ts), zap.String("topic", topic), zap.Int32("partition", partition), zap.Error(err))
}

func kafkaSuccessHandler(ts time.Time, topic string, partition int32) {
	logger.Debug("success kafka send message", zap.Time("ts", ts), zap.String("topic", topic), zap.Int32("partition", partition))
}
