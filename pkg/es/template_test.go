package es

import (
	"fmt"
	"testing"
)

// json yaml to go struct
type test struct {
	Kind     string `json:"kind" yaml:"kind"`
	Metadata struct {
		Name              string `json:"name" yaml:"name"`
		Namespace         string `json:"namespace" yaml:"namespace"`
		CreationTimestamp string `json:"creationTimestamp" yaml:"creationTimestamp"`
	} `json:"metadata" yaml:"metadata"`
	Message string `json:"message" yaml:"message"`
	Reason  string `json:"reason" yaml:"reason"`
}

func TestTemplate(t *testing.T) {
	index := "test"
	docs := test{
		Kind: "test",
		Metadata: struct {
			Name              string `json:"name" yaml:"name"`
			Namespace         string `json:"namespace" yaml:"namespace"`
			CreationTimestamp string `json:"creationTimestamp" yaml:"creationTimestamp"`
		}{
			Name:              "test",
			Namespace:         "test",
			CreationTimestamp: "2021-09-01T00:00:00Z",
		},
		Message: "test message",
		Reason:  "create",
	}

	result, err := ConvertTemplate(index, docs)
	if err != nil {
		t.Errorf("ConvertTemplate() error = %v, want %v", err, nil)
	}
	if len(result) == 0 {
		t.Errorf("ConvertTemplate() = %v, want %v", len(result), 0)
	}

	strResult := string(result)
	fmt.Println(strResult)
}
