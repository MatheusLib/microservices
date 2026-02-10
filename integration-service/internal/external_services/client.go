package external_services

import (
	"context"
	"net/http"
	"time"
)

type Client interface {
	Ping(ctx context.Context) error
}

type httpClient struct {
	baseURL string
	client  *http.Client
}

func NewHTTPClient(baseURL string) Client {
	return &httpClient{
		baseURL: baseURL,
		client: &http.Client{Timeout: 3 * time.Second},
	}
}

func (c *httpClient) Ping(ctx context.Context) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL, nil)
	if err != nil {
		return err
	}
	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	_ = resp.Body.Close()
	return nil
}
