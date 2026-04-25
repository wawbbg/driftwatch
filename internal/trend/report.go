package trend

import (
	"fmt"
	"io"
	"os"
	"sort"
	"text/tabwriter"
)

// Reporter writes a formatted trend summary to a writer.
type Reporter struct {
	w io.Writer
}

// NewReporter returns a Reporter that writes to w. If w is nil, os.Stdout is
// used.
func NewReporter(w io.Writer) *Reporter {
	if w == nil {
		w = os.Stdout
	}
	return &Reporter{w: w}
}

// Write renders a slice of ServiceTrend values as a human-readable table.
func (r *Reporter) Write(trends []ServiceTrend) error {
	if len(trends) == 0 {
		_, err := fmt.Fprintln(r.w, "no trend data available")
		return err
	}

	sorted := make([]ServiceTrend, len(trends))
	copy(sorted, trends)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Service < sorted[j].Service
	})

	tw := tabwriter.NewWriter(r.w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "SERVICE\tSAMPLES\tAVG DIFFS\tDIRECTION\tLAST SEEN")
	for _, t := range sorted {
		fmt.Fprintf(tw, "%s\t%d\t%.1f\t%s\t%s\n",
			t.Service,
			t.Samples,
			t.AvgDiffs,
			t.Direction,
			t.LastSeen.Format("2006-01-02 15:04:05"),
		)
	}
	return tw.Flush()
}
