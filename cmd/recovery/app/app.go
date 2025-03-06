package app

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"example.com/stradvision-project/cmd/recovery/config"
	"example.com/stradvision-project/pkg/kafka/consumer"
	"example.com/stradvision-project/pkg/kube"
	"example.com/stradvision-project/pkg/logger"
	"example.com/stradvision-project/pkg/storage"
)

type Application struct {
	// kafka
	kc *consumer.KafkaConsumer

	// storage
	stg *storage.Handler

	// data buffer
	buf *kube.EventBuffer
}

func NewApplication(config *config.Config) (*Application, error) {
	app := &Application{}

	// storage handler
	stg, err := storage.NewHandler(
		config.Storage.Name, config.Storage.Path,
		storage.WithMaxFileSize(config.Storage.MaxFileSize),
		storage.WithMaxFileCount(config.Storage.MaxFileCount),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create storage handler: %w", err)
	}
	app.stg = stg

	// kafka consumer
	kc, err := consumer.NewKafkaConsumer(
		config.Kafka.Broker, config.Kafka.GroupID, config.Kafka.Topic,
		consumer.WithErrFunc(ConsumerErrorHandler),
		consumer.WithDoFunc(app.ConsumerDo),
		consumer.WithBalanceStrategy(config.Kafka.RebalanceStrategy),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create kafka consumer: %w", err)
	}
	app.kc = kc

	// data buffer
	buf, err := kube.NewEventBuffer(app.bufferDo, app.bufferErrHandler)
	if err != nil {
		return nil, fmt.Errorf("failed to create event buffer: %w", err)
	}
	app.buf = buf

	return app, nil
}

func (app *Application) Run() {
	logger.Info("start recovery application ...	")

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go app.buf.Run()
	go app.kc.Run()

	<-sigChan
	app.buf.Close()
	logger.Info("close buffer ...")
	app.kc.Close()
	logger.Info("close consumer ...")

	logger.Info("stop recovery application ...")
}
