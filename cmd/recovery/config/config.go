package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

const (
	// kafka 설정 환경변수
	EnvKafkaBroker  string = "KAFKA_BROKER"
	EnvKafkaGroupID string = "KAFKA_GROUP_ID"
	EnvKafkaTopic   string = "KAFKA_TOPIC"

	// storage 설정 환경변수
	EnvStorageName         string = "STORAGE_NAME"
	EnvStoragePath         string = "STORAGE_PATH"
	EnvStorageMaxFileSize  string = "STORAGE_MAX_FILE_SIZE"
	EnvStorageMaxFileCount string = "STORAGE_MAX_FILE_COUNT"
)

type Config struct {
	Kafka struct {
		Broker  []string `yaml:"broker"`  // 필수
		GroupID string   `yaml:"groupID"` // 필수
		Topic   string   `yaml:"topic"`   // 필수

		RebalanceStrategy string `yaml:"rebalanceStrategy"`
	} `yaml:"kafka"`

	Storage struct {
		Name string `yaml:"name"` // 필수
		Path string `yaml:"path"` // 필수

		MaxFileSize  int `yaml:"maxFileSize"`
		MaxFileCount int `yaml:"maxFileCount"`
	} `yaml:"storage"`
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

	// Storage
	if config.Storage.Name == "" {
		return fmt.Errorf("config storage name required")
	}
	if config.Storage.Path == "" {
		return fmt.Errorf("config storage path required")
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
	if env := os.Getenv(EnvKafkaGroupID); env != "" {
		config.Kafka.GroupID = env
	}
	if env := os.Getenv(EnvKafkaTopic); env != "" {
		config.Kafka.Topic = env
	}

	// Storage
	if env := os.Getenv(EnvStorageName); env != "" {
		config.Storage.Name = env
	}
	if env := os.Getenv(EnvStoragePath); env != "" {
		config.Storage.Path = env
	}
	if env := os.Getenv(EnvStorageMaxFileSize); env != "" {
		if value, err := strconv.Atoi(env); err == nil {
			config.Storage.MaxFileSize = value
		}
	}
	if env := os.Getenv(EnvStorageMaxFileCount); env != "" {
		if value, err := strconv.Atoi(env); err == nil {
			config.Storage.MaxFileCount = value
		}
	}
}
