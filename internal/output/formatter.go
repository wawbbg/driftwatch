package output

import (
	"fmt"
	"io"
	"os"
	"strings"
)

// Format represents the output format for drift results.
type Format string

const (
	FormatText Format = "text"
	FormatJSON Format = "json"
)

// ParseFormat parses a string into a Format, returning an error if unrecognized.
func ParseFormat(s string) (Format, error) {
	switch strings.ToLower(s) {
	case "text", "":
		return FormatText, nil
	case "json":
		return FormatJSON, nil
	default:
		return "", fmt.Errorf("unknown output format %q: must be \"text\" or \"json\"", s)
	}
}

// ExitCode returns the appropriate process exit code based on whether drift
// was detected. Callers can use this to signal drift to CI pipelines.
func ExitCode(driftDetected bool) int {
	if driftDetected {
		return 1
	}
	return 0
}

// Printer writes a simple labelled line to w, falling back to os.Stdout.
func Printer(w io.Writer) func(format string, args ...any) {
	if w == nil {
		w = os.Stdout
	}
	return func(format string, args ...any) {
		fmt.Fprintf(w, format+"\n", args...)
	}
}
