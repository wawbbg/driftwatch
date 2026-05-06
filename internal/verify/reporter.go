package verify

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"
)

// jsonResult is the serialisable form of a Result.
type jsonResult struct {
	Service   string   `json:"service"`
	CheckedAt string   `json:"checked_at"`
	OK        bool     `json:"ok"`
	DriftCount int     `json:"drift_count"`
	Fields    []string `json:"drifted_fields,omitempty"`
}

// Reporter writes verification results to an io.Writer.
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

// WriteText writes a plain-text summary of results.
func (r *Reporter) WriteText(results []Result) {
	if len(results) == 0 {
		fmt.Fprintln(r.out, "no verification results")
		return
	}
	for _, res := range results {
		fmt.Fprintln(r.out, res.String())
	}
}

// WriteJSON writes results as a JSON array.
func (r *Reporter) WriteJSON(results []Result) error {
	out := make([]jsonResult, 0, len(results))
	for _, res := range results {
		fields := make([]string, 0, len(res.Diffs))
		for _, d := range res.Diffs {
			fields = append(fields, d.Field)
		}
		out = append(out, jsonResult{
			Service:    res.Service,
			CheckedAt:  res.CheckedAt.Format(time.RFC3339),
			OK:         res.OK,
			DriftCount: len(res.Diffs),
			Fields:     fields,
		})
	}
	enc := json.NewEncoder(r.out)
	enc.SetIndent("", "  ")
	return enc.Encode(out)
}
