package app

import (
	"encoding/json"

	"example.com/stradvision-project/pkg/kube"
	"example.com/stradvision-project/pkg/logger"
	"go.uber.org/zap"
)

func (app *Application) ConsumerDo(data []byte) {
	event := &kube.Event{}
	if err := json.Unmarshal(data, event); err != nil {
		logger.Error("failed to consume unmarshal data", zap.Error(err))
		return
	}

	app.buf.AddEvent(event)
}

// Run application
func ConsumerErrorHandler(topic, msg string) {
	logger.Error("failed consumer error", zap.String("topic", topic), zap.String("msg", msg))
}
