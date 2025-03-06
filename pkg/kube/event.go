package kube

import (
	"fmt"
	"time"

	v1 "k8s.io/api/events/v1"
)

/* Example
{
  "metadata": {
    "name": "recovery-6f86dc46dc-288hv.182a25008603f337",
    "namespace": "stradvision",
    "uid": "ea24c536-a48c-4047-a0ef-a44eb092f1ae",
    "resourceVersion": "209584",
    "creationTimestamp": "2025-03-06T07:08:10Z",
    "managedFields": [
      {
        "manager": "kube-scheduler",
        "operation": "Update",
        "apiVersion": "v1",
        "time": "2025-03-06T07:08:10Z",
        "fieldsType": "FieldsV1",
        "fieldsV1": {
          "f:count": {},
          "f:firstTimestamp": {},
          "f:involvedObject": {},
          "f:lastTimestamp": {},
          "f:message": {},
          "f:reason": {},
          "f:reportingComponent": {},
          "f:source": {
            "f:component": {}
          },
          "f:type": {}
        }
      }
    ]
  },
  "eventTime": null,
  "reportingController": "default-scheduler",
  "reason": "Scheduled",
  "regarding": {
    "kind": "Pod",
    "namespace": "stradvision",
    "name": "recovery-6f86dc46dc-288hv",
    "uid": "829030cc-8a1c-4d8f-8cc0-67d3b9c493c4",
    "apiVersion": "v1",
    "resourceVersion": "209579"
  },
  "note": "Successfully assigned stradvision/recovery-6f86dc46dc-288hv to docker-desktop",
  "type": "Normal",
  "deprecatedSource": {
    "component": "default-scheduler"
  },
  "deprecatedFirstTimestamp": "2025-03-06T07:08:10Z",
  "deprecatedLastTimestamp": "2025-03-06T07:08:10Z",
  "deprecatedCount": 1
}
*/

const (
	DefaultFlushMaxCount = 100
	DefaultFlushMaxTime  = 5 * time.Second
)

type Event struct {
	Metadata struct {
		Name              string    `json:"name"`
		Namespace         string    `json:"namespace"`
		UID               string    `json:"uid"`
		ResourceVersion   string    `json:"resourceVersion"`
		CreationTimestamp time.Time `json:"creationTimestamp"`
	} `json:"metadata"`

	EventTime            time.Time `json:"eventTime"`
	RetportingController string    `json:"reportingController"`
	Reason               string    `json:"reason"`

	Regarding struct {
		Kind            string `json:"kind"`
		Namespace       string `json:"namespace"`
		Name            string `json:"name"`
		UID             string `json:"uid"`
		ApiVersion      string `json:"apiVersion"`
		ResourceVersion string `json:"resourceVersion"`
	} `json:"regarding"`

	Note string `json:"note"`
	Type string `json:"type"`

	DeprecatedFirstTimestamp time.Time `json:"deprecatedFirstTimestamp"`
	DeprecatedLastTimestamp  time.Time `json:"deprecatedLastTimestamp"`
	DeprecatedCount          int       `json:"deprecatedCount"`
}

type EventBuffer struct {
	EventChan chan *Event
	Events    []*Event
	closeChan chan struct{}

	DoFunc  func([]*Event) error
	ErrFunc func(error, []*Event)
}

func NewEventBuffer(
	doFunc func([]*Event) error,
	errFunc func(error, []*Event),
) (*EventBuffer, error) {
	if doFunc == nil {
		return nil, fmt.Errorf("doFunc is nil")
	}
	if errFunc == nil {
		return nil, fmt.Errorf("errFunc is nil")
	}

	buffer := &EventBuffer{
		EventChan: make(chan *Event),
		Events:    make([]*Event, 0),
		closeChan: make(chan struct{}),
		DoFunc:    doFunc,
		ErrFunc:   errFunc,
	}

	return buffer, nil
}

func (eb *EventBuffer) AddEvent(event *Event) {
	eb.EventChan <- event
}

func (eb *EventBuffer) Run() {
	ticker := time.NewTicker(DefaultFlushMaxTime)
	defer ticker.Stop()

	for {
		select {
		case <-eb.closeChan:
			return
		case event := <-eb.EventChan:
			eb.Events = append(eb.Events, event)
			if len(eb.Events) >= DefaultFlushMaxCount {
				if err := eb.DoFunc(eb.Events); err != nil {
					eb.ErrFunc(err, eb.Events)
				}
				eb.Events = make([]*Event, 0)
			}
		case <-ticker.C:
			if len(eb.Events) > 0 {
				if err := eb.DoFunc(eb.Events); err != nil {
					eb.ErrFunc(err, eb.Events)
				}
				eb.Events = make([]*Event, 0)
			}
		}
	}
}

func (eb *EventBuffer) Close() {
	close(eb.closeChan)
	close(eb.EventChan)
}

func ConvertEvent(object *v1.Event) *Event {
	event := &Event{}
	event.Metadata.Name = object.Name
	event.Metadata.Namespace = object.Namespace
	event.Metadata.UID = string(object.UID)
	event.Metadata.ResourceVersion = object.ResourceVersion
	event.Metadata.CreationTimestamp = object.CreationTimestamp.Time

	event.EventTime = object.EventTime.Time
	event.RetportingController = object.ReportingController
	event.Reason = object.Reason

	event.Regarding.Kind = object.Regarding.Kind
	event.Regarding.Namespace = object.Regarding.Namespace
	event.Regarding.Name = object.Regarding.Name
	event.Regarding.UID = string(object.Regarding.UID)
	event.Regarding.ApiVersion = object.Regarding.APIVersion
	event.Regarding.ResourceVersion = object.Regarding.ResourceVersion

	event.Note = object.Note
	event.Type = object.Type

	event.DeprecatedFirstTimestamp = object.DeprecatedFirstTimestamp.Time
	event.DeprecatedLastTimestamp = object.DeprecatedLastTimestamp.Time
	event.DeprecatedCount = int(object.DeprecatedCount)

	return event
}
