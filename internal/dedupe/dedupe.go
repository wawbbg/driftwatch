// Package dedupe provides deduplication of drift differences across
// multiple service checks within a single run.
package dedupe

import (
	"fmt"
	"sync"

	"github.com/example/driftwatch/internal/diff"
)

// key uniquely identifies a drift difference by service, field, and values.
type key struct {
	service string
	field   string
	want    string
	got     string
}

// Deduper tracks seen differences and filters duplicates.
type Deduper struct {
	mu   sync.Mutex
	seen map[key]struct{}
}

// New returns an initialised Deduper.
func New() *Deduper {
	return &Deduper{seen: make(map[key]struct{})}
}

// Apply returns only those diffs in ds that have not been seen before
// for the given service name. Duplicate entries are silently dropped.
func (d *Deduper) Apply(service string, ds []diff.Difference) []diff.Difference {
	d.mu.Lock()
	defer d.mu.Unlock()

	out := ds[:0:0]
	for _, entry := range ds {
		k := key{
			service: service,
			field:   entry.Field,
			want:    fmt.Sprintf("%v", entry.Want),
			got:     fmt.Sprintf("%v", entry.Got),
		}
		if _, exists := d.seen[k]; exists {
			continue
		}
		d.seen[k] = struct{}{}
		out = append(out, entry)
	}
	return out
}

// Reset clears all previously seen differences, allowing subsequent calls
// to Apply to treat every difference as new.
func (d *Deduper) Reset() {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.seen = make(map[key]struct{})
}

// Len returns the number of unique differences recorded so far.
func (d *Deduper) Len() int {
	d.mu.Lock()
	defer d.mu.Unlock()
	return len(d.seen)
}
