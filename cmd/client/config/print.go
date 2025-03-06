package config

import (
	"example.com/stradvision-project/pkg/logger"
	"go.uber.org/zap"
)

func ShowConfig(config *Config) {
	logger.Debug("kubernetes",
		zap.String("config", config.Kube.Config),
		zap.Duration("resync", config.Kube.Resync),
	)

	logger.Debug("kafka",
		zap.Strings("broker", config.Kafka.Broker),
		zap.String("topic", config.Kafka.Topic),
	)
}
