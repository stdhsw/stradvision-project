package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

const (
	// Kafka 설정 환경변수
	EnvKafkaBroker string = "KAFKA_BROKER"
	// kafka consumer
	EnvKafkaGoupID     string = "KAFKA_GROUP_ID"
	EnvKafkaTopic      string = "KAFKA_TOPIC"
	EnvKafkaRebanlance string = "KAFKA_REBALANCE"
	// kafka producer
	EnvKafkaDlqTopic     string = "KAFKA_DLQ_TOPIC"
	EnvKafkaTimeout      string = "KAFKA_TIMEOUT"
	EnvKafkaRetry        string = "KAFKA_RETRY"
	EnvKafkaRetryBackoff string = "KAFKA_RETRY_BACKOFF"
	EnvKafkaFlushMsg     string = "KAFKA_FLUSH_MSG"
	EnvKafkaFlushSec     string = "KAFKA_FLUSH_SEC"
	EnvKafkaFlushByte    string = "KAFKA_FLUSH_BYTE"

	// Storage ElasticSearch 설정 환경변수
	EnvElasticAddress string = "ELASTIC_ADDRESS"
	EnvElasticUser    string = "ELASTIC_USER"
	EnvElasticPass    string = "ELASTIC_PASS"
	EnvElasticIndex   string = "ELASTIC_INDEX"
)

type Config struct {
	Kafka struct {
		// kafka 공통 설정
		Broker []string `yaml:"broker"` // 필수

		// Consumer 그룹 설정
		GroupID           string `yaml:"groupID"` // 필수
		Topic             string `yaml:"topic"`   // 필수
		RebalanceStrategy string `yaml:"rebalanceStrategy"`

		// Dead Letter Queue 설정
		DlqTopic     string        `yaml:"dlqTopic"` // 필수
		Timeout      time.Duration `yaml:"timeout"`
		Retry        int           `yaml:"retry"`
		RetryBackoff time.Duration `yaml:"retryBackoff"`
		FlushMsg     int           `yaml:"flushMsg"`
		FlushTime    time.Duration `yaml:"flushTime"`
		FlushByte    int           `yaml:"flushByte"`
	} `yaml:"kafka"`

	ElasticSearch struct {
		Addresses []string `yaml:"addresses"` // 필수
		User      string   `yaml:"user"`      // 필수
		Pass      string   `yaml:"pass"`      // 필수
		Index     string   `yaml:"index"`     // 필수
	} `yaml:"elasticsearch"`
}

// LoadConfig 설정 파일을 읽어서 Config 구조체로 반환
func LoadConfig(fileName string) (*Config, error) {
	config := &Config{}
	if err := readFile(fileName, config); err != nil {
		return nil, fmt.Errorf("failed config read file: %w", err)
	}
	readEnv(config)
	ShowConfig(config)

	return config, checkConfig(config)
}

// checkConfig 필수 설정 값이 있는지 확인
func checkConfig(config *Config) error {
	// Kafka
	if config.Kafka.Broker == nil || len(config.Kafka.Broker) == 0 {
		return fmt.Errorf("config kafka broker required")
	}
	if config.Kafka.GroupID == "" {
		return fmt.Errorf("config kafka groupID required")
	}
	if config.Kafka.Topic == "" {
		return fmt.Errorf("config kafka topic required")
	}

	// ElasticSearch
	if config.ElasticSearch.Addresses == nil || len(config.ElasticSearch.Addresses) == 0 {
		return fmt.Errorf("config elasticsearch addresses required")
	}

	return nil
}

// readFile 설정 파일을 읽어서 Config 구조체로 반환
func readFile(filename string, config *Config) error {
	file, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	if err := yaml.Unmarshal(file, config); err != nil {
		return err
	}

	return nil
}

// readEnv 환경변수를 읽어서 Config 구조체에 설정
func readEnv(config *Config) {
	// Kafka
	if env := os.Getenv(EnvKafkaBroker); env != "" {
		brokers := strings.Split(env, ",")
		config.Kafka.Broker = brokers
	}
	// Consumer 그룹 설정
	if env := os.Getenv(EnvKafkaGoupID); env != "" {
		config.Kafka.GroupID = env
	}
	if env := os.Getenv(EnvKafkaTopic); env != "" {
		config.Kafka.Topic = env
	}
	if env := os.Getenv(EnvKafkaRebanlance); env != "" {
		config.Kafka.RebalanceStrategy = env
	}
	// Dead Letter Queue 설정
	if env := os.Getenv(EnvKafkaDlqTopic); env != "" {
		config.Kafka.DlqTopic = env
	}
	if env := os.Getenv(EnvKafkaTimeout); env != "" {
		if value, err := time.ParseDuration(env); err == nil {
			config.Kafka.Timeout = value
		}
	}
	if env := os.Getenv(EnvKafkaRetry); env != "" {
		if value, err := strconv.Atoi(env); err == nil {
			config.Kafka.Retry = value
		}
	}
	if env := os.Getenv(EnvKafkaRetryBackoff); env != "" {
		if value, err := time.ParseDuration(env); err == nil {
			config.Kafka.RetryBackoff = value
		}
	}
	if env := os.Getenv(EnvKafkaFlushMsg); env != "" {
		if value, err := strconv.Atoi(env); err == nil {
			config.Kafka.FlushMsg = value
		}
	}
	if env := os.Getenv(EnvKafkaFlushSec); env != "" {
		if value, err := time.ParseDuration(env); err == nil {
			config.Kafka.FlushTime = value
		}
	}
	if env := os.Getenv(EnvKafkaFlushByte); env != "" {
		if value, err := strconv.Atoi(env); err == nil {
			config.Kafka.FlushByte = value
		}
	}

	// Elasticsearch
	if env := os.Getenv(EnvElasticAddress); env != "" {
		addresses := strings.Split(env, ",")
		config.ElasticSearch.Addresses = addresses
	}
	if env := os.Getenv(EnvElasticUser); env != "" {
		config.ElasticSearch.User = env
	}
	if env := os.Getenv(EnvElasticPass); env != "" {
		config.ElasticSearch.Pass = env
	}
}
