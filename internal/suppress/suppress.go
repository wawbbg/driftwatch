// Package suppress provides a mechanism to silence known drift entries
// for a configurable duration, preventing repeated alerts for accepted deviations.
package suppress

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Entry represents a suppressed drift field for a service.
type Entry struct {
	Service   string    `json:"service"`
	Field     string    `json:"field"`
	Reason    string    `json:"reason"`
	ExpiresAt time.Time `json:"expires_at"`
}

// IsExpired reports whether the suppression window has passed.
func (e Entry) IsExpired(now time.Time) bool {
	return now.After(e.ExpiresAt)
}

// Store persists a suppression entry under dir/<service>/suppress.json.
func Store(dir string, e Entry) error {
	if e.Service == "" {
		return fmt.Errorf("suppress: service name required")
	}
	path := filepath.Join(dir, e.Service, "suppress.json")
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("suppress: mkdir: %w", err)
	}
	entries, _ := load(path)
	entries = append(entries, e)
	return write(path, entries)
}

// Active returns all non-expired suppressions for the given service.
func Active(dir, service string, now time.Time) ([]Entry, error) {
	path := filepath.Join(dir, service, "suppress.json")
	all, err := load(path)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var active []Entry
	for _, e := range all {
		if !e.IsExpired(now) {
			active = append(active, e)
		}
	}
	return active, nil
}

// IsSuppressed reports whether field is currently suppressed for service.
func IsSuppressed(dir, service, field string, now time.Time) (bool, error) {
	entries, err := Active(dir, service, now)
	if err != nil {
		return false, err
	}
	for _, e := range entries {
		if e.Field == field || e.Field == "*" {
			return true, nil
		}
	}
	return false, nil
}

func load(path string) ([]Entry, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var entries []Entry
	return entries, json.Unmarshal(data, &entries)
}

func write(path string, entries []Entry) error {
	data, err := json.MarshalIndent(entries, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}
