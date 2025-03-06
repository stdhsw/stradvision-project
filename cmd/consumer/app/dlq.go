package app

import (
	"time"

	"example.com/stradvision-project/pkg/logger"
	"go.uber.org/zap"
)

// ProducerErrorHandler
func ProducerErrorHandler(ts time.Time, topic string, partition int32, err error) {
	logger.Error("failed dlq producer error", zap.Time("ts", ts), zap.String("topic", topic), zap.Int32("partition", partition), zap.Error(err))
}

// ProducerSuccessHandler
func ProducerSuccessHandler(ts time.Time, topic string, partition int32) {
	logger.Info("success dlq producer", zap.Time("ts", ts), zap.String("topic", topic), zap.Int32("partition", partition))
}
