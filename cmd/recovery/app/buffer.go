package app

import (
	"encoding/json"

	"example.com/stradvision-project/pkg/kube"
	"example.com/stradvision-project/pkg/logger"
	"go.uber.org/zap"
)

func (app *Application) bufferDo(events []*kube.Event) (err error) {
	// storage flush
	for _, event := range events {
		var data []byte
		data, err = json.Marshal(event)
		if err != nil {
			continue
		}

		// 개행 추가
		data = append(data, '\n')
		if err = app.stg.WriteData(data); err != nil {
			continue
		}

		logger.Debug("storage flush",
			zap.String("event", event.Metadata.Name),
			zap.String("kind", event.Regarding.Kind),
			zap.String("namespace", event.Regarding.Namespace),
			zap.String("name", event.Regarding.Name),
		)
	}

	return err
}

func (app *Application) bufferErrHandler(err error, events []*kube.Event) {
	// log error
	logger.Error("failed to flush events", zap.Error(err))

	// log write
	for _, event := range events {
		logger.Error("failed storage flush", zap.String("event", event.Metadata.Name))
	}
}
