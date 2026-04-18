package report

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/user/driftwatch/internal/drift"
)

// Format represents the output format for reports.
type Format string

const (
	FormatText Format = "text"
	FormatJSON  Format = "json"
)

// Writer writes drift reports to an output stream.
type Writer struct {
	out    io.Writer
	format Format
}

// New creates a new report Writer with the given format.
func New(format Format) *Writer {
	return &Writer{out: os.Stdout, format: format}
}

// NewWithWriter creates a new report Writer with a custom io.Writer.
func NewWithWriter(out io.Writer, format Format) *Writer {
	return &Writer{out: out, format: format}
}

// Write outputs the drift results.
func (w *Writer) Write(results []drift.Difference) error {
	switch w.format {
	case FormatJSON:
		return w.writeJSON(results)
	default:
		return w.writeText(results)
	}
}

func (w *Writer) writeText(results []drift.Difference) error {
	if len(results) == 0 {
		_, err := fmt.Fprintln(w.out, "No drift detected.")
		return err
	}
	fmt.Fprintf(w.out, "Drift detected (%d difference(s)):\n", len(results))
	fmt.Fprintln(w.out, strings.Repeat("-", 40))
	for _, d := range results {
		fmt.Fprintln(w.out, d.String())
	}
	return nil
}

func (w *Writer) writeJSON(results []drift.Difference) error {
	if len(results) == 0 {
		_, err := fmt.Fprintln(w.out, `{"drift":false,"differences":[]}`)
		return err
	}
	fmt.Fprintf(w.out, `{"drift":true,"differences":[`)
	for i, d := range results {
		if i > 0 {
			fmt.Fprint(w.out, ",")
		}
		fmt.Fprintf(w.out, `{"service":%q,"field":%q,"expected":%q,"actual":%q}`,
			d.Service, d.Field, d.Expected, d.Actual)
	}
	_, err := fmt.Fprintln(w.out, `]}`)
	return err
}
