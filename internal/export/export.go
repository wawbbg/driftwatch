// Package export writes drift results to external formats (CSV, NDJSON).
package export

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/driftwatch/internal/diff"
)

// Format represents a supported export format.
type Format string

const (
	FormatCSV    Format = "csv"
	FormatNDJSON Format = "ndjson"
)

// Record is a single exportable drift entry.
type Record struct {
	Service   string    `json:"service"`
	Field     string    `json:"field"`
	Expected  string    `json:"expected"`
	Actual    string    `json:"actual"`
	Timestamp time.Time `json:"timestamp"`
}

// Exporter writes drift records to a writer in a chosen format.
type Exporter struct {
	w   io.Writer
	fmt Format
}

// New returns an Exporter that writes to w using the given format.
func New(w io.Writer, f Format) *Exporter {
	return &Exporter{w: w, fmt: f}
}

// Write converts diffs for a named service into export records and writes them.
func (e *Exporter) Write(service string, diffs []diff.Difference) error {
	records := make([]Record, len(diffs))
	now := time.Now().UTC()
	for i, d := range diffs {
		records[i] = Record{
			Service:   service,
			Field:     d.Field,
			Expected:  fmt.Sprintf("%v", d.Expected),
			Actual:    fmt.Sprintf("%v", d.Actual),
			Timestamp: now,
		}
	}
	switch e.fmt {
	case FormatCSV:
		return writeCSV(e.w, records)
	case FormatNDJSON:
		return writeNDJSON(e.w, records)
	default:
		return fmt.Errorf("export: unsupported format %q", e.fmt)
	}
}

func writeCSV(w io.Writer, records []Record) error {
	cw := csv.NewWriter(w)
	_ = cw.Write([]string{"service", "field", "expected", "actual", "timestamp"})
	for _, r := range records {
		_ = cw.Write([]string{r.Service, r.Field, r.Expected, r.Actual, r.Timestamp.Format(time.RFC3339)})
	}
	cw.Flush()
	return cw.Error()
}

func writeNDJSON(w io.Writer, records []Record) error {
	enc := json.NewEncoder(w)
	for _, r := range records {
		if err := enc.Encode(r); err != nil {
			return err
		}
	}
	return nil
}
