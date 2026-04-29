package lint_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/driftwatch/driftwatch/internal/lint"
)

func TestWrite_NoViolations(t *testing.T) {
	var buf bytes.Buffer
	r := lint.NewReporter(&buf)
	count := r.Write("svc-a", nil)
	if count != 0 {
		t.Errorf("expected 0 errors, got %d", count)
	}
	if !strings.Contains(buf.String(), "no violations") {
		t.Errorf("expected 'no violations' in output, got: %s", buf.String())
	}
}

func TestWrite_WithViolations(t *testing.T) {
	var buf bytes.Buffer
	r := lint.NewReporter(&buf)
	vs := []lint.Violation{
		{Field: "port", Message: "must not be empty", Severity: lint.SeverityError},
		{Field: "HOST", Message: "key should be lowercase", Severity: lint.SeverityWarning},
	}
	count := r.Write("svc-b", vs)
	if count != 1 {
		t.Errorf("expected 1 error, got %d", count)
	}
	out := buf.String()
	if !strings.Contains(out, "svc-b") {
		t.Error("expected service name in output")
	}
	if !strings.Contains(out, "[error]") {
		t.Error("expected error severity in output")
	}
	if !strings.Contains(out, "[warning]") {
		t.Error("expected warning severity in output")
	}
}

func TestWrite_ErrorsFirst(t *testing.T) {
	var buf bytes.Buffer
	r := lint.NewReporter(&buf)
	vs := []lint.Violation{
		{Field: "z_warn", Message: "warn", Severity: lint.SeverityWarning},
		{Field: "a_err", Message: "err", Severity: lint.SeverityError},
	}
	r.Write("svc-c", vs)
	out := buf.String()
	errIdx := strings.Index(out, "[error]")
	warnIdx := strings.Index(out, "[warning]")
	if errIdx > warnIdx {
		t.Error("expected errors to appear before warnings in output")
	}
}

func TestNewReporter_NilFallback(t *testing.T) {
	// Should not panic when writer is nil; falls back to os.Stderr.
	r := lint.NewReporter(nil)
	if r == nil {
		t.Fatal("expected non-nil reporter")
	}
}
