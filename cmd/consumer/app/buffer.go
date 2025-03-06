package app

import (
	"encoding/json"
	"fmt"

	"example.com/stradvision-project/pkg/kube"
	"example.com/stradvision-project/pkg/logger"
	"go.uber.org/zap"
)

const (
	template = `{"index": {"_index": "%s" } }`
)

func (app *Application) bufferDo(events []*kube.Event) error {
	// elasticsearch flush
	result := make([]byte, 0)

	for _, event := range events {
		data, err := json.Marshal(event)
		if err != nil {
			return fmt.Errorf("failed bufferDo marshal event: %w", err)
		}

		meta := []byte(fmt.Sprintf(template, app.index))
		result = append(result, meta...)
		result = append(result, '\n')
		result = append(result, data...)
		result = append(result, '\n')
	}

	result = append(result, '\n')
	logger.Debug("bufferDo", zap.String("result", string(result)))
	if err := app.ec.WriteBulk(app.index, result); err != nil {
		return fmt.Errorf("failed bufferDo send message: %w", err)
	}

	return nil
}

func (app *Application) bufferErrHandler(err error, events []*kube.Event) {
	// log error
	logger.Error("failed to flush events", zap.Error(err))

	// send to kafka dlq
	for _, event := range events {
		data, err := json.Marshal(event)
		if err != nil {
			logger.Error("failed bufferErrHandler marshal event", zap.Error(err))
			continue
		}

		app.dlpKp.SendMessage(app.index, data)
	}
}
