// Package summary provides aggregation of drift detection results
// across multiple services into a concise summary report.
package summary

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/example/driftwatch/internal/drift"
)

// Result holds aggregated drift statistics for a run.
type Result struct {
	Timestamp   time.Time
	Total       int
	Drifted     int
	Clean       int
	ServiceDrift map[string][]drift.Difference
}

// New builds a Result from a map of service name to detected differences.
func New(diffs map[string][]drift.Difference) Result {
	r := Result{
		Timestamp:    time.Now().UTC(),
		Total:        len(diffs),
		ServiceDrift: make(map[string][]drift.Difference, len(diffs)),
	}
	for svc, d := range diffs {
		r.ServiceDrift[svc] = d
		if len(d) > 0 {
			r.Drifted++
		} else {
			r.Clean++
		}
	}
	return r
}

// Write prints a human-readable summary to w.
func Write(w io.Writer, r Result) {
	if w == nil {
		w = os.Stdout
	}
	fmt.Fprintf(w, "Drift Summary — %s\n", r.Timestamp.Format(time.RFC3339))
	fmt.Fprintf(w, "  Services checked : %d\n", r.Total)
	fmt.Fprintf(w, "  Clean            : %d\n", r.Clean)
	fmt.Fprintf(w, "  Drifted          : %d\n", r.Drifted)
	if r.Drifted == 0 {
		return
	}
	fmt.Fprintln(w, "  Affected services:")
	for svc, diffs := range r.ServiceDrift {
		if len(diffs) == 0 {
			continue
		}
		fmt.Fprintf(w, "    • %s (%d field(s))\n", svc, len(diffs))
		for _, d := range diffs {
			fmt.Fprintf(w, "        - %s\n", d)
		}
	}
}
