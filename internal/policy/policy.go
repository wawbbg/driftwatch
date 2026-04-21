// Package policy evaluates drift results against user-defined severity rules.
package policy

import (
	"fmt"
	"strings"

	"github.com/example/driftwatch/internal/drift"
)

// Level represents the severity of a policy violation.
type Level string

const (
	LevelWarn  Level = "warn"
	LevelError Level = "error"
	LevelNone  Level = "none"
)

// Rule defines a policy rule that matches fields and assigns a severity level.
type Rule struct {
	Field string `yaml:"field" json:"field"`
	Level Level  `yaml:"level" json:"level"`
}

// Policy holds a collection of rules used to evaluate drift differences.
type Policy struct {
	Rules []Rule `yaml:"rules" json:"rules"`
}

// Violation represents a drift difference that matched a policy rule.
type Violation struct {
	Diff  drift.Difference
	Level Level
}

// String returns a human-readable representation of the violation.
func (v Violation) String() string {
	return fmt.Sprintf("[%s] %s", strings.ToUpper(string(v.Level)), v.Diff)
}

// Evaluate checks a slice of differences against the policy rules and returns
// any violations found. Differences that match no rule are skipped.
func (p *Policy) Evaluate(diffs []drift.Difference) []Violation {
	var violations []Violation
	for _, d := range diffs {
		level := p.levelFor(d.Field)
		if level == LevelNone {
			continue
		}
		violations = append(violations, Violation{Diff: d, Level: level})
	}
	return violations
}

// levelFor returns the Level for the given field name by scanning rules in order.
// Returns LevelNone when no rule matches.
func (p *Policy) levelFor(field string) Level {
	for _, r := range p.Rules {
		if r.Field == "*" || strings.EqualFold(r.Field, field) {
			return r.Level
		}
	}
	return LevelNone
}

// HasErrors reports whether any of the violations are at error level.
func HasErrors(violations []Violation) bool {
	for _, v := range violations {
		if v.Level == LevelError {
			return true
		}
	}
	return false
}
