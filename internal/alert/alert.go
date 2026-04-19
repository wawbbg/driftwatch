// Package alert provides simple alerting hooks for drift detection results.
package alert

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/example/driftwatch/internal/drift"
)

// Level represents the severity of an alert.
type Level string

const (
	LevelWarn  Level = "warn"
	LevelError Level = "error"
)

// Alert holds a formatted drift alert message.
type Alert struct {
	Service string
	Level   Level
	Message string
}

// String returns a human-readable representation of the alert.
func (a Alert) String() string {
	return fmt.Sprintf("[%s] %s: %s", strings.ToUpper(string(a.Level)), a.Service, a.Message)
}

// Notifier sends alerts to a writer.
type Notifier struct {
	w io.Writer
}

// New returns a Notifier writing to stderr.
func New() *Notifier {
	return &Notifier{w: os.Stderr}
}

// NewWithWriter returns a Notifier writing to w.
func NewWithWriter(w io.Writer) *Notifier {
	return &Notifier{w: w}
}

// Notify emits an alert for each drift difference.
// Returns the number of alerts emitted.
func (n *Notifier) Notify(service string, diffs []drift.Difference) int {
	count := 0
	for _, d := range diffs {
		lvl := LevelWarn
		if d.Expected == "" {
			lvl = LevelError
		}
		a := Alert{
			Service: service,
			Level:   lvl,
			Message: d.String(),
		}
		fmt.Fprintln(n.w, a.String())
		count++
	}
	return count
}
