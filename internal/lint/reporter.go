package lint

import (
	"fmt"
	"io"
	"os"
	"sort"
)

// Reporter writes lint violations to a writer in a human-readable format.
type Reporter struct {
	w io.Writer
}

// NewReporter returns a Reporter that writes to w.
// If w is nil it falls back to os.Stderr.
func NewReporter(w io.Writer) *Reporter {
	if w == nil {
		w = os.Stderr
	}
	return &Reporter{w: w}
}

// Write outputs all violations grouped by severity.
// It returns the number of error-level violations found.
func (r *Reporter) Write(service string, vs []Violation) int {
	if len(vs) == 0 {
		fmt.Fprintf(r.w, "lint: %s — no violations\n", service)
		return 0
	}

	// Sort: errors first, then warnings; within each group sort by field.
	sort.Slice(vs, func(i, j int) bool {
		if vs[i].Severity != vs[j].Severity {
			return vs[i].Severity == SeverityError
		}
		return vs[i].Field < vs[j].Field
	})

	fmt.Fprintf(r.w, "lint: %s — %d violation(s)\n", service, len(vs))
	for _, v := range vs {
		fmt.Fprintf(r.w, "  %s\n", v.String())
	}

	var errCount int
	for _, v := range vs {
		if v.Severity == SeverityError {
			errCount++
		}
	}
	return errCount
}
