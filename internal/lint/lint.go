// Package lint validates service config maps against a set of rules,
// reporting fields that are empty, mistyped, or violate naming conventions.
package lint

import (
	"fmt"
	"strings"
)

// Severity indicates how serious a lint violation is.
type Severity string

const (
	SeverityError   Severity = "error"
	SeverityWarning Severity = "warning"
)

// Violation describes a single lint finding.
type Violation struct {
	Field    string
	Message  string
	Severity Severity
}

func (v Violation) String() string {
	return fmt.Sprintf("[%s] %s: %s", v.Severity, v.Field, v.Message)
}

// Rule is a function that inspects a config map and returns any violations.
type Rule func(cfg map[string]string) []Violation

// Linter holds a set of rules and runs them against config maps.
type Linter struct {
	rules []Rule
}

// New returns a Linter pre-loaded with the default rule set.
func New() *Linter {
	return &Linter{
		rules: []Rule{
			ruleNoEmptyValues,
			ruleKeyNaming,
		},
	}
}

// WithRule appends a custom rule to the linter.
func (l *Linter) WithRule(r Rule) *Linter {
	l.rules = append(l.rules, r)
	return l
}

// Run applies all rules to cfg and returns the combined violations.
func (l *Linter) Run(cfg map[string]string) []Violation {
	var out []Violation
	for _, r := range l.rules {
		out = append(out, r(cfg)...)
	}
	return out
}

// HasErrors returns true if any violation has severity error.
func HasErrors(vs []Violation) bool {
	for _, v := range vs {
		if v.Severity == SeverityError {
			return true
		}
	}
	return false
}

// ruleNoEmptyValues flags fields whose value is blank.
func ruleNoEmptyValues(cfg map[string]string) []Violation {
	var out []Violation
	for k, v := range cfg {
		if strings.TrimSpace(v) == "" {
			out = append(out, Violation{
				Field:    k,
				Message:  "value must not be empty",
				Severity: SeverityError,
			})
		}
	}
	return out
}

// ruleKeyNaming warns when a key contains uppercase letters.
func ruleKeyNaming(cfg map[string]string) []Violation {
	var out []Violation
	for k := range cfg {
		if k != strings.ToLower(k) {
			out = append(out, Violation{
				Field:    k,
				Message:  "key should be lowercase",
				Severity: SeverityWarning,
			})
		}
	}
	return out
}
