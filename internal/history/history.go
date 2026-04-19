// Package history records drift check runs for trend analysis.
package history

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Entry represents a single drift check run result.
type Entry struct {
	Timestamp   time.Time `json:"timestamp"`
	ServiceName string    `json:"service_name"`
	DriftCount  int       `json:"drift_count"`
	HasDrift    bool      `json:"has_drift"`
}

// Record appends a new entry to the history file for the given service.
func Record(dir, serviceName string, driftCount int) error {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("history: create dir: %w", err)
	}

	path := filepath.Join(dir, serviceName+".json")

	var entries []Entry
	if data, err := os.ReadFile(path); err == nil {
		_ = json.Unmarshal(data, &entries)
	}

	entries = append(entries, Entry{
		Timestamp:   time.Now().UTC(),
		ServiceName: serviceName,
		DriftCount:  driftCount,
		HasDrift:    driftCount > 0,
	})

	data, err := json.MarshalIndent(entries, "", "  ")
	if err != nil {
		return fmt.Errorf("history: marshal: %w", err)
	}

	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("history: write: %w", err)
	}
	return nil
}

// List returns all recorded entries for the given service.
func List(dir, serviceName string) ([]Entry, error) {
	path := filepath.Join(dir, serviceName+".json")
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("history: read: %w", err)
	}

	var entries []Entry
	if err := json.Unmarshal(data, &entries); err != nil {
		return nil, fmt.Errorf("history: unmarshal: %w", err)
	}
	return entries, nil
}
