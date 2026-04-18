// Package snapshot provides functionality to save and load drift snapshots
// for comparing service state across multiple runs.
package snapshot

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/user/driftwatch/internal/drift"
)

// Snapshot represents a recorded state of drift results at a point in time.
type Snapshot struct {
	Timestamp time.Time        `json:"timestamp"`
	Service   string           `json:"service"`
	Diffs     []drift.Difference `json:"diffs"`
}

// Save writes a snapshot to the given file path as JSON.
func Save(path string, s Snapshot) error {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("snapshot: create %q: %w", path, err)
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	if err := enc.Encode(s); err != nil {
		return fmt.Errorf("snapshot: encode: %w", err)
	}
	return nil
}

// Load reads a snapshot from the given file path.
func Load(path string) (Snapshot, error) {
	var s Snapshot
	f, err := os.Open(path)
	if err != nil {
		return s, fmt.Errorf("snapshot: open %q: %w", path, err)
	}
	defer f.Close()

	if err := json.NewDecoder(f).Decode(&s); err != nil {
		return s, fmt.Errorf("snapshot: decode: %w", err)
	}
	return s, nil
}

// New creates a new Snapshot for the given service and diffs.
func New(service string, diffs []drift.Difference) Snapshot {
	return Snapshot{
		Timestamp: time.Now().UTC(),
		Service:   service,
		Diffs:     diffs,
	}
}
