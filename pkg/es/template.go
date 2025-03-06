package es

import (
	"fmt"

	"encoding/json"
)

const (
	template = `{"index": {"_index": "%s"}}\n`
)

func ConvertTemplate(index string, doc interface{}) ([]byte, error) {
	data, err := json.Marshal(doc)
	if err != nil {
		return nil, err
	}

	result := make([]byte, 0)
	meta := []byte(fmt.Sprintf(template, index))
	result = append(result, meta...)
	result = append(result, data...)
	result = append(result, '\n')

	return result, nil
}

func ConvertTemplates(index string, docs []interface{}) ([]byte, error) {
	result := make([]byte, 0)
	for _, doc := range docs {
		data, err := ConvertTemplate(index, doc)
		if err != nil {
			return nil, err
		}
		result = append(result, data...)
	}

	return result, nil
}
