package app

import (
	"encoding/json"

	"example.com/stradvision-project/pkg/kafka/producer"
	"example.com/stradvision-project/pkg/kube"
	"example.com/stradvision-project/pkg/logger"
	"go.uber.org/zap"
	v1 "k8s.io/api/events/v1"
)

type Handler struct {
	kp *producer.KafkaProducer
}

// OnAdd event handler
func (h *Handler) OnAdd(obj interface{}, _ bool) {
	object := obj.(*v1.Event)
	event := kube.ConvertEvent(object)

	jsonData, err := json.Marshal(event)
	if err != nil {
		logger.Error("[OnAdd] failed to marshal event object", zap.Error(err))
		return
	}

	h.kp.SendMessage("", jsonData)
	logger.Debug("[OnAdd] event object",
		zap.String("Kind", event.Regarding.Kind),
		zap.String("Namespace", event.Regarding.Namespace),
		zap.String("Name", event.Regarding.Name),
		zap.String("UID", event.Regarding.UID),
		zap.String("Reason", event.Reason),
	)
}

// OnUpdate event handler
func (h *Handler) OnUpdate(oldObj, newObj interface{}) {
	object := newObj.(*v1.Event)
	event := kube.ConvertEvent(object)

	jsonData, err := json.Marshal(event)
	if err != nil {
		logger.Error("[OnUpdate] failed to marshal event object", zap.Error(err))
		return
	}

	h.kp.SendMessage("", jsonData)
	logger.Debug("[OnUpdate] event object",
		zap.String("Kind", event.Regarding.Kind),
		zap.String("Namespace", event.Regarding.Namespace),
		zap.String("Name", event.Regarding.Name),
		zap.String("UID", event.Regarding.UID),
		zap.String("Reason", event.Reason),
	)
}

// OnDelete event handler
func (h *Handler) OnDelete(obj interface{}) {
	// delete의 경우는 데이터를 수집할 필요가 없으므로 구현하지 않음
}
