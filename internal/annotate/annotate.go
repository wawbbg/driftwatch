// Package annotate attaches free-form key/value annotations to drift results,
// allowing operators to add context (e.g. ticket IDs, owner info) to detected
// differences before they are persisted or reported.
package annotate

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Annotation holds a single key/value pair and the time it was recorded.
type Annotation struct {
	Key       string    `json:"key"`
	Value     string    `json:"value"`
	RecordedAt time.Time `json:"recorded_at"`
}

// Set writes an annotation for the given service to the annotations directory.
// Annotations are stored as a JSON array; each call appends to the list.
func Set(dir, service, key, value string) error {
	if service == "" {
		return fmt.Errorf("annotate: service name must not be empty")
	}
	if key == "" {
		return fmt.Errorf("annotate: key must not be empty")
	}
	key = strings.ToLower(strings.TrimSpace(key))

	path := filepath.Join(dir, service+".json")
	anns, err := readAll(path)
	if err != nil {
		return err
	}

	anns = append(anns, Annotation{
		Key:       key,
		Value:     value,
		RecordedAt: time.Now().UTC(),
	})
	return writeAll(path, anns)
}

// Get returns all annotations for the given service.
// Returns an empty slice (and no error) when none exist.
func Get(dir, service string) ([]Annotation, error) {
	if service == "" {
		return nil, fmt.Errorf("annotate: service name must not be empty")
	}
	path := filepath.Join(dir, service+".json")
	return readAll(path)
}

// Delete removes all annotations whose key matches the given key for a service.
func Delete(dir, service, key string) error {
	if service == "" {
		return fmt.Errorf("annotate: service name must not be empty")
	}
	key = strings.ToLower(strings.TrimSpace(key))
	path := filepath.Join(dir, service+".json")
	anns, err := readAll(path)
	if err != nil {
		return err
	}
	filtered := anns[:0]
	for _, a := range anns {
		if a.Key != key {
			filtered = append(filtered, a)
		}
	}
	return writeAll(path, filtered)
}

func readAll(path string) ([]Annotation, error) {
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return []Annotation{}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("annotate: read %s: %w", path, err)
	}
	var anns []Annotation
	if err := json.Unmarshal(data, &anns); err != nil {
		return nil, fmt.Errorf("annotate: parse %s: %w", path, err)
	}
	return anns, nil
}

func writeAll(path string, anns []Annotation) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("annotate: mkdir: %w", err)
	}
	data, err := json.MarshalIndent(anns, "", "  ")
	if err != nil {
		return fmt.Errorf("annotate: marshal: %w", err)
	}
	return os.WriteFile(path, data, 0o644)
}
