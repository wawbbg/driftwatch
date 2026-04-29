// Package score computes a numeric drift health score for a service
// based on the number and severity of diffs detected.
package score

import (
	"fmt"
	"io"
	"os"

	"github.com/driftwatch/driftwatch/internal/diff"
)

// Grade represents a letter-grade bucket for a drift score.
type Grade string

const (
	GradeA Grade = "A" // 90–100: healthy
	GradeB Grade = "B" // 75–89
	GradeC Grade = "C" // 50–74
	GradeD Grade = "D" // 25–49
	GradeF Grade = "F" // 0–24: severe drift
)

// Result holds the computed score and grade for a single service.
type Result struct {
	Service string
	Score   int
	Grade   Grade
	Drifted int
	Total   int
}

// String returns a human-readable summary of the result.
func (r Result) String() string {
	return fmt.Sprintf("%s: score=%d grade=%s drifted=%d/%d",
		r.Service, r.Score, r.Grade, r.Drifted, r.Total)
}

// Compute calculates a drift health score given the full set of expected
// keys and the diffs detected. Score is 0–100; higher is healthier.
func Compute(service string, expected map[string]any, diffs []diff.Diff) Result {
	total := len(expected)
	drifted := len(diffs)

	var score int
	if total == 0 {
		score = 100
	} else {
		healthy := total - drifted
		if healthy < 0 {
			healthy = 0
		}
		score = (healthy * 100) / total
	}

	return Result{
		Service: service,
		Score:   score,
		Grade:   gradeFor(score),
		Drifted: drifted,
		Total:   total,
	}
}

// gradeFor maps a numeric score to a letter grade.
func gradeFor(s int) Grade {
	switch {
	case s >= 90:
		return GradeA
	case s >= 75:
		return GradeB
	case s >= 50:
		return GradeC
	case s >= 25:
		return GradeD
	default:
		return GradeF
	}
}

// Reporter writes score results to an io.Writer.
type Reporter struct{ w io.Writer }

// NewReporter returns a Reporter that writes to w, falling back to os.Stdout.
func NewReporter(w io.Writer) *Reporter {
	if w == nil {
		w = os.Stdout
	}
	return &Reporter{w: w}
}

// Write prints each result on its own line.
func (r *Reporter) Write(results []Result) {
	for _, res := range results {
		fmt.Fprintln(r.w, res.String())
	}
}
