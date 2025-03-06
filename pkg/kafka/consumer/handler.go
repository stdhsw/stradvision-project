package consumer

import "github.com/IBM/sarama"

type consumerGroupHandler struct {
	doFunc func([]byte)
}

func (consumerGroupHandler) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (consumerGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }
func (h consumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		h.doFunc(msg.Value)
		session.MarkMessage(msg, "")
	}
	return nil
}
