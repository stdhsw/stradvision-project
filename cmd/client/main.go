package main

import (
	"os"
	"strconv"

	"example.com/stradvision-project/cmd/client/app"
	"example.com/stradvision-project/cmd/client/config"
	"example.com/stradvision-project/pkg/logger"
	"go.uber.org/zap"
)

const (
	AppName string = "client"

	EnvLogLevel    string = "LOG_LEVEL"
	EnvLogSize     string = "LOG_SIZE"
	EnvLogAge      string = "LOG_AGE"
	EnvLogBack     string = "LOG_BACK"
	EnvLogCompress string = "LOG_COMPRESS"

	DefaultConfigPath string = "/etc/stradvision/config.yaml"
)

func init() {
	logLevel := os.Getenv(EnvLogLevel)
	logSize, _ := strconv.Atoi(os.Getenv(EnvLogSize))
	logAge, _ := strconv.Atoi(os.Getenv(EnvLogAge))
	logBack, _ := strconv.Atoi(os.Getenv(EnvLogBack))
	logCompress, _ := strconv.ParseBool(os.Getenv(EnvLogCompress))

	logger.InitLogger(AppName,
		logger.WithLogLevel(logLevel),
		logger.WithLogMaxSize(logSize),
		logger.WithLogMaxAge(logAge),
		logger.WithLogMaxBackups(logBack),
		logger.WithLogCompress(logCompress),
	)
}

func main() {
	cfg, err := config.LoadConfig(DefaultConfigPath)
	if err != nil {
		logger.Panic("failed to load config", zap.Error(err))
	}

	app, err := app.NewApplication(cfg)
	if err != nil {
		logger.Panic("failed to create application", zap.Error(err))
	}
	app.Run()
}
