package app

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"example.com/stradvision-project/cmd/consumer/config"
	"example.com/stradvision-project/pkg/es"
	"example.com/stradvision-project/pkg/kafka/consumer"
	"example.com/stradvision-project/pkg/kafka/producer"
	"example.com/stradvision-project/pkg/kube"
	"example.com/stradvision-project/pkg/logger"
)

type Application struct {
	// kafka
	kc    *consumer.KafkaConsumer
	dlpKp *producer.KafkaProducer

	// elasticsearch
	index string
	ec    *es.Client

	// data buffer
	buf *kube.EventBuffer
}

func NewApplication(config *config.Config) (*Application, error) {
	app := &Application{}

	// elasticsearch client
	ec, err := es.NewElasticsearchClient(
		config.ElasticSearch.Addresses, config.ElasticSearch.User, config.ElasticSearch.Pass,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create elasticsearch client: %w", err)
	}
	app.ec = ec
	app.index = config.ElasticSearch.Index

	// kafka dlq producer
	kp, err := producer.NewKafkaProducer(
		config.Kafka.Broker, config.Kafka.DlqTopic,
		producer.WithTimeout(config.Kafka.Timeout),
		producer.WithRetry(config.Kafka.Retry),
		producer.WithRetryBackoff(config.Kafka.RetryBackoff),
		producer.WithFlushMaxMessages(config.Kafka.FlushMsg),
		producer.WithFlushFrequency(config.Kafka.FlushTime),
		producer.WithFlushBytes(config.Kafka.FlushByte),
		producer.WithErrorFunc(ProducerErrorHandler),
		producer.WithSuccessFunc(ProducerSuccessHandler),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create kafka producer dlq: %w", err)
	}
	app.dlpKp = kp

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
	logger.Info("start consumer application ...")

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go app.buf.Run()
	go app.dlpKp.Run()
	go app.kc.Run()

	<-sigChan
	app.buf.Close()
	logger.Info("close buffer ...")
	app.kc.Close()
	logger.Info("close consumer ...")
	app.dlpKp.Close()
	logger.Info("close dlq producer ...")

	logger.Info("stop consumer application ...")
}
