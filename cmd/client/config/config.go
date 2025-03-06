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
	// Kubernetes 설정 환경변수
	EnvKubeConfig string = "KUBECONFIG"
	EnvResyncTime string = "RESYNC_TIME"

	// Kafka 설정 환경변수
	EnvKafkaBroker       string = "KAFKA_BROKER"
	EnvKafkaTopic        string = "KAFKA_TOPIC"
	EnvKafkaTimeout      string = "KAFKA_TIMEOUT"
	EnvKafkaRetry        string = "KAFKA_RETRY"
	EnvKafkaRetryBackoff string = "KAFKA_RETRY_BACKOFF"
	EnvKafkaFlushMsg     string = "KAFKA_FLUSH_MSG"
	EnvKafkaFlushSec     string = "KAFKA_FLUSH_SEC"
	EnvKafkaFlushByte    string = "KAFKA_FLUSH_BYTE"
)

type Config struct {
	Kube struct {
		Config string        `yaml:"config"` // 없으면 in-cluster 자동 설정
		Resync time.Duration `yaml:"resync"`
	} `yaml:"kube"`

	Kafka struct {
		Broker       []string      `yaml:"broker"` // 필수
		Topic        string        `yaml:"topic"`  // 필수
		Timeout      time.Duration `yaml:"timeout"`
		Retry        int           `yaml:"retry"`
		RetryBackoff time.Duration `yaml:"retryBackoff"`
		FlushMsg     int           `yaml:"flushMsg"`
		FlushTime    time.Duration `yaml:"flushTime"`
		FlushByte    int           `yaml:"flushByte"`
	} `yaml:"kafka"`
}

// LoadConfig 설정 파일을 읽어서 Config 구조체로 반환
func LoadConfig(filename string) (*Config, error) {
	config := &Config{}
	if err := readFile(filename, config); err != nil {
		return nil, fmt.Errorf("failed config read file: %w", err)
	}
	readEnv(config)
	ShowConfig(config)

	return config, checkConfig(config)
}

// checkConfig 필수 설정 값이 있는지 확인
func checkConfig(config *Config) error {
	if config.Kafka.Broker == nil || len(config.Kafka.Broker) == 0 {
		return fmt.Errorf("config kafka broker required")
	}
	if config.Kafka.Topic == "" {
		return fmt.Errorf("config kafka topic required")
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
	// kubernetes client 설정
	if env := os.Getenv(EnvKubeConfig); env != "" {
		config.Kube.Config = env
	}
	if env := os.Getenv(EnvResyncTime); env != "" {
		if value, err := time.ParseDuration(env); err == nil {
			config.Kube.Resync = value
		}
	}

	// kafka producer 설정
	if env := os.Getenv(EnvKafkaBroker); env != "" {
		brokers := strings.Split(env, ",")
		config.Kafka.Broker = brokers
	}
	if env := os.Getenv(EnvKafkaTopic); env != "" {
		config.Kafka.Topic = env
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
}
