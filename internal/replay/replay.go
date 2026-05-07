// Package replay provides functionality for replaying historical drift
// records to reconstruct the state of a service at a point in time.
package replay

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/driftwatch/internal/history"
)

// Entry represents a single replayed drift event.
type Entry struct {
	Service   string
	Timestamp time.Time
	Diffs     []string
}

// Result holds the outcome of a replay operation.
type Result struct {
	Service string
	Entries []Entry
	From    time.Time
	To      time.Time
}

// HasDrift reports whether any entry in the result contains diffs.
func (r Result) HasDrift() bool {
	for _, e := range r.Entries {
		if len(e.Diffs) > 0 {
			return true
		}
	}
	return false
}

// Replayer loads and filters history records for a service.
type Replayer struct {
	dir string
	out io.Writer
}

// New returns a Replayer that reads from dir and writes to os.Stdout.
func New(dir string) *Replayer {
	return NewWithWriter(dir, os.Stdout)
}

// NewWithWriter returns a Replayer with a custom writer.
func NewWithWriter(dir string, w io.Writer) *Replayer {
	if w == nil {
		w = os.Stdout
	}
	return &Replayer{dir: dir, out: w}
}

// Run replays all history records for service between from and to (inclusive).
func (r *Replayer) Run(service string, from, to time.Time) (Result, error) {
	records, err := history.List(r.dir, service)
	if err != nil {
		return Result{}, fmt.Errorf("replay: list history for %q: %w", service, err)
	}

	result := Result{Service: service, From: from, To: to}
	for _, rec := range records {
		if rec.Timestamp.Before(from) || rec.Timestamp.After(to) {
			continue
		}
		result.Entries = append(result.Entries, Entry{
			Service:   rec.Service,
			Timestamp: rec.Timestamp,
			Diffs:     rec.Diffs,
		})
	}

	fmt.Fprintf(r.out, "replay: %s — %d entries between %s and %s\n",
		service, len(result.Entries),
		from.Format(time.RFC3339), to.Format(time.RFC3339))

	return result, nil
}
