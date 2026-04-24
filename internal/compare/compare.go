// Package compare provides utilities for comparing service config
// snapshots across two points in time, producing a structured delta.
package compare

import (
	"fmt"
	"time"

	"github.com/driftwatch/internal/diff"
)

// Delta represents the result of comparing two snapshots for a service.
type Delta struct {
	Service   string
	Before    map[string]any
	After     map[string]any
	Diffs     []diff.Difference
	CapturedAt time.Time
}

// HasDrift reports whether any differences were found.
func (d Delta) HasDrift() bool {
	return len(d.Diffs) > 0
}

// String returns a human-readable summary of the delta.
func (d Delta) String() string {
	if !d.HasDrift() {
		return fmt.Sprintf("[%s] no drift detected", d.Service)
	}
	return fmt.Sprintf("[%s] %d field(s) changed", d.Service, len(d.Diffs))
}

// Snapshots compares two config maps for the named service and returns
// a Delta describing what changed between them.
func Snapshots(service string, before, after map[string]any) Delta {
	diffs := diff.Maps(before, after)
	return Delta{
		Service:    service,
		Before:     before,
		After:      after,
		Diffs:      diffs,
		CapturedAt: time.Now().UTC(),
	}
}
