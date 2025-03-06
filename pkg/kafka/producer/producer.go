package producer

import (
	"fmt"
	"time"

	"github.com/IBM/sarama"
)

type KafkaProducer struct {
	producer sarama.AsyncProducer
	topic    string
	closeCh  chan struct{}

	errFunc     func(ts time.Time, topic string, partition int32, err error)
	successFunc func(ts time.Time, topic string, partition int32)
}

// NewKafkaProducer KafkaProducer 생성
// brokers: Kafka 브로커 주소
// topic: Kafka 토픽
// opt: producer 설정
func NewKafkaProducer(brokers []string, topic string, opts ...Option) (*KafkaProducer, error) {
	pConfig := fromOptions(opts)

	producer, err := sarama.NewAsyncProducer(brokers, pConfig.config)
	if err != nil {
		return nil, fmt.Errorf("failed to create producer: %w", err)
	}

	kp := &KafkaProducer{
		producer:    producer,
		topic:       topic,
		closeCh:     make(chan struct{}),
		errFunc:     pConfig.errFunc,
		successFunc: pConfig.successFunc,
	}

	return kp, nil
}

// run KafkaProducer 결과 처리
func (kp *KafkaProducer) Run() {
	for {
		select {
		case <-kp.closeCh:
			return
		case err := <-kp.producer.Errors():
			kp.errFunc(err.Msg.Timestamp, err.Msg.Topic, err.Msg.Partition, err.Err)
		case success := <-kp.producer.Successes():
			kp.successFunc(success.Timestamp, success.Topic, success.Partition)
		}
	}
}

// SendMessage 메시지 전송
func (kp *KafkaProducer) SendMessage(key string, data []byte) {
	msg := &sarama.ProducerMessage{
		Topic: kp.topic,
		Key:   sarama.StringEncoder(key),
		Value: sarama.ByteEncoder(data),
	}
	kp.producer.Input() <- msg
}

// Close KafkaProducer 종료
func (kp *KafkaProducer) Close() {
	close(kp.closeCh)
	kp.producer.AsyncClose()
}
