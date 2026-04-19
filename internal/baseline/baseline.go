// Package baseline manages storing and comparing baseline configs
// for drift detection across runs.
package baseline

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Entry represents a stored baseline for a single service.
type Entry struct {
	ServiceName string                 `json:"service_name"`
	CapturedAt  time.Time              `json:"captured_at"`
	Fields      map[string]interface{} `json:"fields"`
}

// Store persists a baseline entry to a JSON file under dir.
func Store(dir, serviceName string, fields map[string]interface{}) error {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("baseline: mkdir %s: %w", dir, err)
	}
	e := Entry{
		ServiceName: serviceName,
		CapturedAt:  time.Now().UTC(),
		Fields:      fields,
	}
	data, err := json.MarshalIndent(e, "", "  ")
	if err != nil {
		return fmt.Errorf("baseline: marshal: %w", err)
	}
	path := filepath.Join(dir, serviceName+".json")
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("baseline: write %s: %w", path, err)
	}
	return nil
}

// Load reads a baseline entry for the given service from dir.
func Load(dir, serviceName string) (*Entry, error) {
	path := filepath.Join(dir, serviceName+".json")
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, fmt.Errorf("baseline: no baseline for service %q", serviceName)
		}
		return nil, fmt.Errorf("baseline: read %s: %w", path, err)
	}
	var e Entry
	if err := json.Unmarshal(data, &e); err != nil {
		return nil, fmt.Errorf("baseline: unmarshal: %w", err)
	}
	return &e, nil
}
