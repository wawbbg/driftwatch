// Package audit provides a structured audit log for drift detection events.
// Each audit entry records when a check ran, which service was checked,
// how many diffs were found, and whether policy errors were present.
package audit

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Entry represents a single audit log record.
type Entry struct {
	Timestamp   time.Time `json:"timestamp"`
	Service     string    `json:"service"`
	DriftCount  int       `json:"drift_count"`
	PolicyError bool      `json:"policy_error"`
	Message     string    `json:"message,omitempty"`
}

// Record appends an audit entry for the given service to the audit log file
// located at dir/<service>.audit.jsonl. The file is created if it does not exist.
func Record(dir, service string, driftCount int, policyError bool, message string) error {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("audit: create dir: %w", err)
	}

	entry := Entry{
		Timestamp:   time.Now().UTC(),
		Service:     service,
		DriftCount:  driftCount,
		PolicyError: policyError,
		Message:     message,
	}

	data, err := json.Marshal(entry)
	if err != nil {
		return fmt.Errorf("audit: marshal entry: %w", err)
	}

	path := filepath.Join(dir, service+".audit.jsonl")
	f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return fmt.Errorf("audit: open file: %w", err)
	}
	defer f.Close()

	_, err = fmt.Fprintf(f, "%s\n", data)
	return err
}

// List returns all audit entries recorded for the given service.
// Returns an empty slice if no audit file exists.
func List(dir, service string) ([]Entry, error) {
	path := filepath.Join(dir, service+".audit.jsonl")
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return []Entry{}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("audit: read file: %w", err)
	}

	var entries []Entry
	for _, line := range splitLines(data) {
		if len(line) == 0 {
			continue
		}
		var e Entry
		if err := json.Unmarshal(line, &e); err != nil {
			return nil, fmt.Errorf("audit: parse entry: %w", err)
		}
		entries = append(entries, e)
	}
	return entries, nil
}

func splitLines(data []byte) [][]byte {
	var lines [][]byte
	start := 0
	for i, b := range data {
		if b == '\n' {
			lines = append(lines, data[start:i])
			start = i + 1
		}
	}
	if start < len(data) {
		lines = append(lines, data[start:])
	}
	return lines
}
