package app

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"example.com/stradvision-project/cmd/client/config"
	"example.com/stradvision-project/pkg/kafka/producer"
	"example.com/stradvision-project/pkg/kube"
	"example.com/stradvision-project/pkg/logger"
)

type Application struct {
	k8sClient *kube.Client
	kp        *producer.KafkaProducer

	handler *Handler
}

func NewApplication(config *config.Config) (*Application, error) {
	// Kafka producer
	kp, err := producer.NewKafkaProducer(
		config.Kafka.Broker, config.Kafka.Topic,
		producer.WithTimeout(config.Kafka.Timeout),
		producer.WithRetry(config.Kafka.Retry),
		producer.WithRetryBackoff(config.Kafka.RetryBackoff),
		producer.WithFlushMaxMessages(config.Kafka.FlushMsg),
		producer.WithFlushFrequency(config.Kafka.FlushTime),
		producer.WithFlushBytes(config.Kafka.FlushByte),
		producer.WithErrorFunc(kafkaErrorHandler),
		producer.WithSuccessFunc(kafkaSuccessHandler),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create kafka producer: %w", err)
	}

	// kuberentes client
	handler := &Handler{kp: kp}
	kc, err := kube.NewClient(
		handler,
		kube.WithKubeConfig(config.Kube.Config),
		kube.WithResyncTime(config.Kube.Resync),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create kubernetes client: %w", err)
	}

	return &Application{
		k8sClient: kc,
		kp:        kp,
		handler:   handler,
	}, nil
}

func (app *Application) Run() {
	logger.Info("start application ...")

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go app.kp.Run()
	go app.k8sClient.Run()

	<-sigChan
	app.k8sClient.Close()
	logger.Info("closed kubernetes client")
	app.kp.Close()
	logger.Info("closed kafka producer")

	logger.Info("stop application ...")
}
