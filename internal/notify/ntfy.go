package notify

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type Ntfy struct {
	baseURL string
	topic   string
	http    *http.Client
}

func NewNtfy(baseURL, topic string) *Ntfy {
	return &Ntfy{
		baseURL: baseURL,
		topic:   topic,
		http:    &http.Client{Timeout: 10 * time.Second},
	}
}

func (n *Ntfy) Send(ctx context.Context, title, message string) error {
	url := n.baseURL + "/" + n.topic

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, strings.NewReader(message))
	if err != nil {
		return fmt.Errorf("build request: %w", err)
	}
	req.Header.Set("Title", title)

	resp, err := n.http.Do(req)
	if err != nil {
		return fmt.Errorf("send notification: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status: %s", resp.Status)
	}
	return nil
}
