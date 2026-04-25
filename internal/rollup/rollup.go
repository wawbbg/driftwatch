// Package rollup aggregates drift results across multiple services
// into a single summary suitable for reporting or alerting.
package rollup

import (
	"fmt"
	"io"
	"os"
	"sort"
	"text/tabwriter"

	"github.com/yourorg/driftwatch/internal/diff"
)

// ServiceResult holds the drift outcome for a single service.
type ServiceResult struct {
	Service string
	Diffs   []diff.Difference
}

// Report is a rolled-up view of drift across all checked services.
type Report struct {
	Results []ServiceResult
}

// New builds a Report from a map of service name to diffs.
func New(results map[string][]diff.Difference) *Report {
	r := &Report{}
	for svc, diffs := range results {
		r.Results = append(r.Results, ServiceResult{Service: svc, Diffs: diffs})
	}
	sort.Slice(r.Results, func(i, j int) bool {
		return r.Results[i].Service < r.Results[j].Service
	})
	return r
}

// TotalDrifted returns the number of services that have at least one diff.
func (r *Report) TotalDrifted() int {
	count := 0
	for _, res := range r.Results {
		if len(res.Diffs) > 0 {
			count++
		}
	}
	return count
}

// HasDrift reports whether any service has drifted.
func (r *Report) HasDrift() bool {
	return r.TotalDrifted() > 0
}

// Write prints a human-readable rollup table to w.
// Falls back to os.Stdout if w is nil.
func (r *Report) Write(w io.Writer) {
	if w == nil {
		w = os.Stdout
	}
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "SERVICE\tDRIFTED FIELDS\tSTATUS")
	for _, res := range r.Results {
		status := "ok"
		if len(res.Diffs) > 0 {
			status = "DRIFT"
		}
		fmt.Fprintf(tw, "%s\t%d\t%s\n", res.Service, len(res.Diffs), status)
	}
	tw.Flush()
}
