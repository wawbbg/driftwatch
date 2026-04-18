package fetcher

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// ServiceState holds the live key-value config state fetched from a service.
type ServiceState struct {
	Name   string
	Fields map[string]string
}

// Fetcher retrieves live configuration state from a remote service endpoint.
type Fetcher struct {
	client *http.Client
}

// New returns a Fetcher with a default timeout.
func New() *Fetcher {
	return &Fetcher{
		client: &http.Client{Timeout: 10 * time.Second},
	}
}

// NewWithClient returns a Fetcher using the provided HTTP client (useful for testing).
func NewWithClient(c *http.Client) *Fetcher {
	return &Fetcher{client: c}
}

// Fetch calls the given URL and decodes the JSON response into a flat string map.
// The returned ServiceState will have Name set to the provided name argument.
func (f *Fetcher) Fetch(name, url string) (*ServiceState, error) {
	resp, err := f.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("fetcher: GET %s: %w", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("fetcher: unexpected status %d for %s", resp.StatusCode, url)
	}

	var raw map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return nil, fmt.Errorf("fetcher: decode %s: %w", url, err)
	}

	fields := make(map[string]string, len(raw))
	for k, v := range raw {
		fields[k] = fmt.Sprintf("%v", v)
	}

	return &ServiceState{Name: name, Fields: fields}, nil
}
