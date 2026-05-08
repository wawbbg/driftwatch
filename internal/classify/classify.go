// Package classify assigns severity levels to detected drift differences.
package classify

import (
	"strings"

	"github.com/example/driftwatch/internal/diff"
)

// Severity represents the importance of a drift difference.
type Severity string

const (
	SeverityLow      Severity = "low"
	SeverityMedium   Severity = "medium"
	SeverityHigh     Severity = "high"
	SeverityCritical Severity = "critical"
)

// Result holds a classified difference.
type Result struct {
	Diff     diff.Difference
	Severity Severity
}

// String returns a human-readable representation.
func (r Result) String() string {
	return string(r.Severity) + ": " + r.Diff.String()
}

// Classifier assigns severities to differences based on field rules.
type Classifier struct {
	criticalFields []string
	highFields     []string
}

// New returns a Classifier with sensible defaults.
func New() *Classifier {
	return &Classifier{
		criticalFields: []string{"password", "secret", "token", "key"},
		highFields:     []string{"host", "port", "endpoint", "url", "database"},
	}
}

// NewWithFields returns a Classifier with explicit field lists.
func NewWithFields(critical, high []string) *Classifier {
	return &Classifier{criticalFields: critical, highFields: high}
}

// Classify assigns a Severity to a single difference.
func (c *Classifier) Classify(d diff.Difference) Result {
	field := strings.ToLower(d.Field)
	for _, f := range c.criticalFields {
		if strings.Contains(field, f) {
			return Result{Diff: d, Severity: SeverityCritical}
		}
	}
	for _, f := range c.highFields {
		if strings.Contains(field, f) {
			return Result{Diff: d, Severity: SeverityHigh}
		}
	}
	if d.Got == "" {
		return Result{Diff: d, Severity: SeverityMedium}
	}
	return Result{Diff: d, Severity: SeverityLow}
}

// Apply classifies a slice of differences and returns Results.
func (c *Classifier) Apply(diffs []diff.Difference) []Result {
	out := make([]Result, 0, len(diffs))
	for _, d := range diffs {
		out = append(out, c.Classify(d))
	}
	return out
}

// HasCritical reports whether any result carries critical severity.
func HasCritical(results []Result) bool {
	for _, r := range results {
		if r.Severity == SeverityCritical {
			return true
		}
	}
	return false
}
