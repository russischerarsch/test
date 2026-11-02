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
			Timeout: 30 * time.Second,
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
	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read body response, err %w", err)
	}
	r := &Response[T]{Code: resp.StatusCode, Body: raw}
	if r.Code != http.StatusOK && r.Code != http.StatusCreated {
		return r, nil
	}
	if err := json.Unmarshal(raw, &r.Data); err != nil {
		return r, fmt.Errorf("failed to decode response data to struct, err %w", err)
	}

	return r, nil
}

func Post[T any](ctx context.Context, c *httpClient, path, method string, reqBody interface{}) (*Response[T], error) {
	return Do[T](ctx, c, path, method, reqBody)
}
func Get[T any](ctx context.Context, c *httpClient, path, method string, reqBody interface{}) (*Response[T], error) {
	return Do[T](ctx, c, path, method, reqBody)
}
