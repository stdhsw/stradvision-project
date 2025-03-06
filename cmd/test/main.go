package main

import (
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"example.com/stradvision-project/pkg/logger"
	"go.uber.org/zap"
)

const (
	AppName string = "test"

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
	logger.Info("test start ...")

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	i := 0
LOOP:
	for {
		select {
		case <-ticker.C:
			logger.Info("test ...", zap.Int("i", i))
			i++
		case <-sigChan:
			break LOOP
		}
	}

	logger.Info("test end ...")
}
