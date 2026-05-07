// Package checkpoint records and retrieves named drift-check
// checkpoints so that runs can be resumed or compared across time.
package checkpoint

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Entry is a single saved checkpoint.
type Entry struct {
	Service   string            `json:"service"`
	Name      string            `json:"name"`
	Fields    map[string]string `json:"fields"`
	CreatedAt time.Time         `json:"created_at"`
}

// Save writes a checkpoint entry to dir/<service>/<name>.json.
func Save(dir, service, name string, fields map[string]string) error {
	if dir == "" {
		return errors.New("checkpoint: dir must not be empty")
	}
	if service == "" {
		return errors.New("checkpoint: service must not be empty")
	}
	if name == "" {
		return errors.New("checkpoint: name must not be empty")
	}

	dest := filepath.Join(dir, service)
	if err := os.MkdirAll(dest, 0o755); err != nil {
		return fmt.Errorf("checkpoint: mkdir %s: %w", dest, err)
	}

	e := Entry{
		Service:   service,
		Name:      name,
		Fields:    fields,
		CreatedAt: time.Now().UTC(),
	}

	b, err := json.MarshalIndent(e, "", "  ")
	if err != nil {
		return fmt.Errorf("checkpoint: marshal: %w", err)
	}

	path := filepath.Join(dest, name+".json")
	return os.WriteFile(path, b, 0o644)
}

// Load reads a checkpoint entry from dir/<service>/<name>.json.
func Load(dir, service, name string) (Entry, error) {
	path := filepath.Join(dir, service, name+".json")
	b, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return Entry{}, fmt.Errorf("checkpoint: %s/%s not found", service, name)
		}
		return Entry{}, fmt.Errorf("checkpoint: read %s: %w", path, err)
	}

	var e Entry
	if err := json.Unmarshal(b, &e); err != nil {
		return Entry{}, fmt.Errorf("checkpoint: unmarshal: %w", err)
	}
	return e, nil
}

// List returns all checkpoint names stored for service under dir.
func List(dir, service string) ([]string, error) {
	path := filepath.Join(dir, service)
	entries, err := os.ReadDir(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, nil
		}
		return nil, fmt.Errorf("checkpoint: readdir %s: %w", path, err)
	}

	var names []string
	for _, e := range entries {
		if !e.IsDir() && filepath.Ext(e.Name()) == ".json" {
			names = append(names, e.Name()[:len(e.Name())-5])
		}
	}
	return names, nil
}
