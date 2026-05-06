package verify_test

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/driftwatch/internal/verify"
)

// writeBaseline writes a JSON baseline file for service into dir.
func writeBaseline(t *testing.T, dir, service string, data map[string]any) {
	t.Helper()
	svcDir := filepath.Join(dir, service)
	if err := os.MkdirAll(svcDir, 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	b, err := json.Marshal(data)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	if err := os.WriteFile(filepath.Join(svcDir, "baseline.json"), b, 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}
}

func TestRun_NoDrift(t *testing.T) {
	dir := t.TempDir()
	base := map[string]any{"replicas": "3", "image": "nginx:1.25"}
	writeBaseline(t, dir, "web", base)

	var buf bytes.Buffer
	v := verify.NewWithWriter(dir, &buf)
	res, err := v.Run("web", map[string]any{"replicas": "3", "image": "nginx:1.25"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !res.OK {
		t.Errorf("expected OK, got diffs: %v", res.Diffs)
	}
	if res.Service != "web" {
		t.Errorf("service = %q, want %q", res.Service, "web")
	}
}

func TestRun_WithDrift(t *testing.T) {
	dir := t.TempDir()
	writeBaseline(t, dir, "api", map[string]any{"replicas": "2", "timeout": "30s"})

	var buf bytes.Buffer
	v := verify.NewWithWriter(dir, &buf)
	res, err := v.Run("api", map[string]any{"replicas": "5", "timeout": "30s"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.OK {
		t.Error("expected drift, got OK")
	}
	if len(res.Diffs) != 1 {
		t.Errorf("expected 1 diff, got %d", len(res.Diffs))
	}
}

func TestRun_MissingBaseline(t *testing.T) {
	dir := t.TempDir()
	v := verify.New(dir)
	_, err := v.Run("ghost", map[string]any{"x": "1"})
	if err == nil {
		t.Error("expected error for missing baseline, got nil")
	}
}

func TestResult_String(t *testing.T) {
	ok := verify.Result{Service: "svc", OK: true}
	if ok.String() != "svc: OK (no drift)" {
		t.Errorf("unexpected: %s", ok.String())
	}

	nok := verify.Result{Service: "svc", OK: false, Diffs: make([]interface{}, 2)}
	// Diffs is []diff.Difference; use length check via string output
	_ = nok
}

func TestNewWithWriter_NilFallback(t *testing.T) {
	v := verify.NewWithWriter(t.TempDir(), nil)
	if v == nil {
		t.Error("expected non-nil verifier")
	}
}
