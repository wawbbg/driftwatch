package classify

import (
	"fmt"
	"io"
	"os"
	"sort"
)

// Reporter writes classified results to a writer.
type Reporter struct {
	w io.Writer
}

// NewReporter returns a Reporter writing to w; falls back to os.Stdout if nil.
func NewReporter(w io.Writer) *Reporter {
	if w == nil {
		w = os.Stdout
	}
	return &Reporter{w: w}
}

// Write outputs classified results grouped by severity.
func (r *Reporter) Write(service string, results []Result) {
	if len(results) == 0 {
		fmt.Fprintf(r.w, "[classify] %s: no classified drift\n", service)
		return
	}

	// Sort: critical first, then high, medium, low.
	order := map[Severity]int{
		SeverityCritical: 0,
		SeverityHigh:     1,
		SeverityMedium:   2,
		SeverityLow:      3,
	}
	sorted := make([]Result, len(results))
	copy(sorted, results)
	sort.Slice(sorted, func(i, j int) bool {
		return order[sorted[i].Severity] < order[sorted[j].Severity]
	})

	fmt.Fprintf(r.w, "[classify] %s: %d classified difference(s)\n", service, len(sorted))
	for _, res := range sorted {
		fmt.Fprintf(r.w, "  [%s] field=%q want=%q got=%q\n",
			res.Severity, res.Diff.Field, res.Diff.Want, res.Diff.Got)
	}
}
