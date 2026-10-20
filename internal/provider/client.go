package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Client struct {
	baseURL string
	http    *http.Client
}

func NewClient(baseURL string) *Client {
	return &Client{
		baseURL: baseURL,
		http:    &http.Client{Timeout: 10 * time.Second},
	}
}

func (c *Client) GetQuote(ctx context.Context, symbol, currency string) (Quote, error) {
	url := fmt.Sprintf("%s/price/%s/%s", c.baseURL, symbol, currency)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return Quote{}, fmt.Errorf("build request: %w", err)
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return Quote{}, fmt.Errorf("call API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return Quote{}, fmt.Errorf("unexpected status: %s", resp.Status)
	}

	var q Quote
	if err := json.NewDecoder(resp.Body).Decode(&q); err != nil {
		return Quote{}, fmt.Errorf("decode response: %w", err)
	}
	if q.Price <= 0 {
		return Quote{}, fmt.Errorf("invalid price %v", q.Price)
	}
	return q, nil
}
