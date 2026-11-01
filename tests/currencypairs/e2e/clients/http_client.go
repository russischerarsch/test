package clients

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type httpClient struct {
	client  *http.Client
	baseURL string
}

func NewHttpClient(url string) *httpClient {
	return &httpClient{
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
		baseURL: url,
	}
}

type Response[T any] struct {
	Code int
	Data *T
	Body []byte
}

func Do[T any](ctx context.Context, c *httpClient, path, method string, reqBody interface{}) (*Response[T], error) {
	var body io.Reader
	if reqBody != nil {
		data, err := json.Marshal(reqBody)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal input, err %w", err)
		}
		body = bytes.NewReader(data)
	}

	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create a request, err %w", err)
	}
	if reqBody != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request, err %w", err)
	}
	defer resp.Body.Close()

	r := &Response[T]{Code: resp.StatusCode}
	if r.Code != http.StatusOK && r.Code != http.StatusCreated {
		r.Body, err = json.Marshal(resp)
		if err != nil {
			return r, fmt.Errorf("failed to marshal response body, err %w", err)
		}
		return r, nil
	}
	if err := json.NewDecoder(resp.Body).Decode(&r.Data); err != nil {
		return r, fmt.Errorf("failed to decode response data to struct, err %w", err)
	}

	return r, nil
}
