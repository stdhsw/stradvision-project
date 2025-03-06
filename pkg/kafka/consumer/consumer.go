package consumer

import (
	"context"
	"fmt"

	"github.com/IBM/sarama"
)

type KafkaConsumer struct {
	cg    sarama.ConsumerGroup
	topic string

	handler sarama.ConsumerGroupHandler
	ctx     context.Context
	cancel  context.CancelFunc

	errFunc func(topic, msg string)
}

func NewKafkaConsumer(brokers []string, groupID, topic string, opts ...Option) (*KafkaConsumer, error) {
	cConfig := fromOptions(opts)

	consumerGroup, err := sarama.NewConsumerGroup(brokers, groupID, cConfig.config)
	if err != nil {
		return nil, fmt.Errorf("failed to create consumer group: %w", err)
	}

	kc := &KafkaConsumer{
		cg:    consumerGroup,
		topic: topic,
		handler: consumerGroupHandler{
			doFunc: cConfig.doFunc,
		},
		errFunc: cConfig.errFunc,
	}
	kc.ctx, kc.cancel = context.WithCancel(context.Background())

	return kc, nil
}

func (kc *KafkaConsumer) Run() {
	for {
		if err := kc.cg.Consume(kc.ctx, []string{kc.topic}, kc.handler); err != nil {
			kc.errFunc(kc.topic, fmt.Errorf("failed to consume: %w", err).Error())
			return
		}
	}
}

func (kc *KafkaConsumer) Close() {
	kc.cancel()
	kc.cg.Close()
}
