package suppress

import (
	"fmt"
	"io"
	"os"
	"text/tabwriter"
	"time"
)

// Reporter writes suppression listings in a human-readable table.
type Reporter struct {
	w io.Writer
}

// NewReporter returns a Reporter writing to w, falling back to os.Stdout.
func NewReporter(w io.Writer) *Reporter {
	if w == nil {
		w = os.Stdout
	}
	return &Reporter{w: w}
}

// Write outputs active suppressions for the given service to the reporter's writer.
func (r *Reporter) Write(dir, service string, now time.Time) error {
	entries, err := Active(dir, service, now)
	if err != nil {
		return fmt.Errorf("suppress reporter: %w", err)
	}
	if len(entries) == 0 {
		fmt.Fprintf(r.w, "no active suppressions for %s\n", service)
		return nil
	}
	tw := tabwriter.NewWriter(r.w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "FIELD\tREASON\tEXPIRES")
	for _, e := range entries {
		fmt.Fprintf(tw, "%s\t%s\t%s\n", e.Field, e.Reason, e.ExpiresAt.Format(time.RFC3339))
	}
	return tw.Flush()
}
