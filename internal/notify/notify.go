// Package notify provides webhook-based notification support for drift events.
package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Payload is the JSON body sent to a webhook endpoint.
type Payload struct {
	Service   string    `json:"service"`
	DriftCount int      `json:"drift_count"`
	Fields    []string  `json:"fields"`
	Timestamp time.Time `json:"timestamp"`
}

// Notifier sends drift payloads to a configured webhook URL.
type Notifier struct {
	webhookURL string
	client     *http.Client
}

// New returns a Notifier using the default HTTP client.
func New(webhookURL string) *Notifier {
	return NewWithClient(webhookURL, &http.Client{Timeout: 10 * time.Second})
}

// NewWithClient returns a Notifier using the provided HTTP client.
func NewWithClient(webhookURL string, client *http.Client) *Notifier {
	return &Notifier{webhookURL: webhookURL, client: client}
}

// Send posts a Payload to the webhook URL.
// It returns an error if the request fails or the server responds with a
// non-2xx status code.
func (n *Notifier) Send(p Payload) error {
	if p.Timestamp.IsZero() {
		p.Timestamp = time.Now().UTC()
	}

	body, err := json.Marshal(p)
	if err != nil {
		return fmt.Errorf("notify: marshal payload: %w", err)
	}

	resp, err := n.client.Post(n.webhookURL, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("notify: post webhook: %w", err)
	}
	defer resp.Body.Close()
	io.Copy(io.Discard, resp.Body)

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("notify: webhook returned status %d", resp.StatusCode)
	}
	return nil
}
