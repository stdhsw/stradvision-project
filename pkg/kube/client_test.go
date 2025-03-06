package kube

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	v1 "k8s.io/api/events/v1"
)

const (
	template = `{"index": {"_index": "%s" } }`
)

type testHandler struct{}

func (h *testHandler) OnAdd(obj interface{}, isInInitialList bool) {
	object := obj.(*v1.Event)

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

	//
	result := make([]byte, 0)
	meta := []byte(fmt.Sprintf(template, "event"))
	result = append(result, meta...)
	result = append(result, '\n')

	data, _ := json.Marshal(event)
	result = append(result, data...)
	result = append(result, '\n')

	fmt.Println("OnAdd")
	fmt.Println(string(result))
}

func (h *testHandler) OnUpdate(oldObj, newObj interface{}) {
	object := newObj.(*v1.Event)
	fmt.Println("OnUpdate")
	fmt.Println(object)
}

func (h *testHandler) OnDelete(obj interface{}) {
	_ = obj
}

func TestClient(t *testing.T) {
	// kubeconfig 파일 경로 설정
	home, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("failed to get user home dir: %v", err)
	}
	kubeConfig := filepath.Join(home, ".kube", "config")

	// 클라이언트 생성
	th := &testHandler{}
	client, err := NewClient(
		th,
		WithKubeConfig(kubeConfig),
	)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}
	client.Run()

	time.Sleep(60 * time.Second)
	client.Close()
}
