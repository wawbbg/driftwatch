package lint_test

import (
	"testing"

	"github.com/driftwatch/driftwatch/internal/lint"
)

func TestRun_NoViolations(t *testing.T) {
	l := lint.New()
	cfg := map[string]string{
		"host": "localhost",
		"port": "8080",
	}
	vs := l.Run(cfg)
	if len(vs) != 0 {
		t.Fatalf("expected no violations, got %d: %v", len(vs), vs)
	}
}

func TestRun_EmptyValue(t *testing.T) {
	l := lint.New()
	cfg := map[string]string{"host": ""}
	vs := l.Run(cfg)
	if len(vs) == 0 {
		t.Fatal("expected at least one violation for empty value")
	}
	if vs[0].Severity != lint.SeverityError {
		t.Errorf("expected error severity, got %s", vs[0].Severity)
	}
}

func TestRun_UppercaseKey(t *testing.T) {
	l := lint.New()
	cfg := map[string]string{"HOST": "localhost"}
	vs := l.Run(cfg)
	var found bool
	for _, v := range vs {
		if v.Severity == lint.SeverityWarning && v.Field == "HOST" {
			found = true
		}
	}
	if !found {
		t.Error("expected warning for uppercase key")
	}
}

func TestHasErrors_True(t *testing.T) {
	vs := []lint.Violation{{Field: "x", Message: "bad", Severity: lint.SeverityError}}
	if !lint.HasErrors(vs) {
		t.Error("expected HasErrors to return true")
	}
}

func TestHasErrors_False(t *testing.T) {
	vs := []lint.Violation{{Field: "x", Message: "warn", Severity: lint.SeverityWarning}}
	if lint.HasErrors(vs) {
		t.Error("expected HasErrors to return false")
	}
}

func TestViolation_String(t *testing.T) {
	v := lint.Violation{Field: "port", Message: "must not be empty", Severity: lint.SeverityError}
	s := v.String()
	if s == "" {
		t.Error("expected non-empty string representation")
	}
}

func TestWithRule_CustomRule(t *testing.T) {
	customRule := func(cfg map[string]string) []lint.Violation {
		var out []lint.Violation
		if _, ok := cfg["required_key"]; !ok {
			out = append(out, lint.Violation{
				Field:    "required_key",
				Message:  "missing required field",
				Severity: lint.SeverityError,
			})
		}
		return out
	}
	l := lint.New().WithRule(customRule)
	cfg := map[string]string{"host": "localhost"}
	vs := l.Run(cfg)
	var found bool
	for _, v := range vs {
		if v.Field == "required_key" {
			found = true
		}
	}
	if !found {
		t.Error("expected custom rule violation for missing required_key")
	}
}
