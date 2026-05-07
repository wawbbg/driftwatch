package replay

import (
	"fmt"
	"io"
	"os"
	"time"
)

// Reporter writes a human-readable replay result to a writer.
type Reporter struct {
	out io.Writer
}

// NewReporter returns a Reporter that writes to w.
// If w is nil, os.Stdout is used.
func NewReporter(w io.Writer) *Reporter {
	if w == nil {
		w = os.Stdout
	}
	return &Reporter{out: w}
}

// Write formats and writes the replay result.
func (r *Reporter) Write(res Result) {
	fmt.Fprintf(r.out, "=== Replay: %s ===\n", res.Service)
	fmt.Fprintf(r.out, "Window : %s → %s\n",
		res.From.Format(time.RFC3339),
		res.To.Format(time.RFC3339))
	fmt.Fprintf(r.out, "Entries: %d\n", len(res.Entries))

	if len(res.Entries) == 0 {
		fmt.Fprintln(r.out, "No drift recorded in this window.")
		return
	}

	for _, e := range res.Entries {
		fmt.Fprintf(r.out, "\n[%s]\n", e.Timestamp.Format(time.RFC3339))
		if len(e.Diffs) == 0 {
			fmt.Fprintln(r.out, "  (no diffs)")
			continue
		}
		for _, d := range e.Diffs {
			fmt.Fprintf(r.out, "  - %s\n", d)
		}
	}
}
