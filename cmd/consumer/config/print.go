package config

import (
	"strings"

	"example.com/stradvision-project/pkg/logger"
	"go.uber.org/zap"
)

func ShowConfig(config *Config) {
	logger.Debug("kafka consumer",
		zap.String("broker", strings.Join(config.Kafka.Broker, ",")),
		zap.String("groupID", config.Kafka.GroupID), zap.String("topic", config.Kafka.Topic),
		zap.String("rebalanceStrategy", config.Kafka.RebalanceStrategy),
	)

	logger.Debug("kafka producer",
		zap.String("DlqTopic", config.Kafka.DlqTopic),
	)

	logger.Debug("elasticsearch",
		zap.String("address", strings.Join(config.ElasticSearch.Addresses, ",")),
		zap.String("username", config.ElasticSearch.User),
		zap.String("password", config.ElasticSearch.Pass),
		zap.String("index", config.ElasticSearch.Index),
	)
}
