package config

import (
	"strings"

	"example.com/stradvision-project/pkg/logger"
	"go.uber.org/zap"
)

func ShowConfig(config *Config) {
	logger.Debug("kafka",
		zap.String("broker", strings.Join(config.Kafka.Broker, ",")),
		zap.String("groupID", config.Kafka.GroupID), zap.String("topic", config.Kafka.Topic),
		zap.String("rebalanceStrategy", config.Kafka.RebalanceStrategy),
	)

	logger.Debug("storage",
		zap.String("name", config.Storage.Name), zap.String("path", config.Storage.Path),
		zap.Int("maxFileSize", config.Storage.MaxFileSize), zap.Int("maxFileCount", config.Storage.MaxFileCount),
	)
}
