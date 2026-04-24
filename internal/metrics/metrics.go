// Package metrics tracks drift detection run statistics.
package metrics

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"
)

// Run holds statistics for a single drift detection run.
type Run struct {
	Service    string    `json:"service"`
	Timestamp  time.Time `json:"timestamp"`
	Total      int       `json:"total_fields"`
	Drifted    int       `json:"drifted_fields"`
	Duration   float64   `json:"duration_ms"`
	HasDrift   bool      `json:"has_drift"`
}

// Collector accumulates run metrics and writes summaries.
type Collector struct {
	out io.Writer
}

// New returns a Collector that writes to stdout.
func New() *Collector {
	return NewWithWriter(os.Stdout)
}

// NewWithWriter returns a Collector that writes to w.
func NewWithWriter(w io.Writer) *Collector {
	return &Collector{out: w}
}

// Record writes a single run metric as a JSON line.
func (c *Collector) Record(r Run) error {
	b, err := json.Marshal(r)
	if err != nil {
		return fmt.Errorf("metrics: marshal: %w", err)
	}
	_, err = fmt.Fprintf(c.out, "%s\n", b)
	return err
}

// Summary prints a human-readable summary of a run to the collector's writer.
func (c *Collector) Summary(r Run) {
	status := "clean"
	if r.HasDrift {
		status = "DRIFT DETECTED"
	}
	fmt.Fprintf(c.out, "[metrics] service=%s status=%s drifted=%d/%d duration=%.2fms\n",
		r.Service, status, r.Drifted, r.Total, r.Duration)
}

// NewRun is a convenience constructor for a Run.
func NewRun(service string, total, drifted int, start time.Time) Run {
	dur := float64(time.Since(start).Microseconds()) / 1000.0
	return Run{
		Service:   service,
		Timestamp: time.Now().UTC(),
		Total:     total,
		Drifted:   drifted,
		Duration:  dur,
		HasDrift:  drifted > 0,
	}
}
