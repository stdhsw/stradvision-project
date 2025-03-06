package es

import (
	"bytes"
	"context"
	"crypto/tls"
	"fmt"
	"net/http"

	"github.com/elastic/go-elasticsearch/v8"
)

type Client struct {
	es *elasticsearch.Client
}

func NewElasticsearchClient(addrs []string, user, pass string) (*Client, error) {
	config := elasticsearch.Config{
		Addresses: addrs,
		Username:  user,
		Password:  pass,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true, // TLS 인증서 검증 비활성화
			},
		},
	}

	es, err := elasticsearch.NewClient(config)
	if err != nil {
		return nil, err
	}

	return &Client{
		es: es,
	}, nil
}

func (c *Client) WriteBulk(index string, data []byte) error {
	buf := bytes.NewBuffer(data)
	res, err := c.es.Bulk(buf, c.es.Bulk.WithContext(context.Background()))
	if err != nil {
		return fmt.Errorf("failed to send elasticsearch bulk request: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("elasticsearch bulk request failed: %s", res.String())
	}

	return nil
}
