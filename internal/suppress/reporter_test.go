package suppress_test

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/driftwatch/driftwatch/internal/suppress"
)

func TestWrite_NoSuppressions(t *testing.T) {
	dir := t.TempDir()
	var buf bytes.Buffer
	r := suppress.NewReporter(&buf)
	if err := r.Write(dir, "svc", now); err != nil {
		t.Fatalf("Write: %v", err)
	}
	if !strings.Contains(buf.String(), "no active suppressions") {
		t.Errorf("unexpected output: %q", buf.String())
	}
}

func TestWrite_WithSuppressions(t *testing.T) {
	dir := t.TempDir()
	e := suppress.Entry{
		Service:   "api",
		Field:     "replicas",
		Reason:    "planned scale",
		ExpiresAt: now.Add(48 * time.Hour),
	}
	_ = suppress.Store(dir, e)
	var buf bytes.Buffer
	r := suppress.NewReporter(&buf)
	if err := r.Write(dir, "api", now); err != nil {
		t.Fatalf("Write: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "replicas") {
		t.Errorf("expected field name in output: %q", out)
	}
	if !strings.Contains(out, "planned scale") {
		t.Errorf("expected reason in output: %q", out)
	}
}

func TestWrite_ExpiredNotShown(t *testing.T) {
	dir := t.TempDir()
	e := suppress.Entry{Service: "api", Field: "image", Reason: "old", ExpiresAt: past}
	_ = suppress.Store(dir, e)
	var buf bytes.Buffer
	r := suppress.NewReporter(&buf)
	_ = r.Write(dir, "api", now)
	if strings.Contains(buf.String(), "image") {
		t.Errorf("expired entry should not appear in output")
	}
}

func TestNewReporter_NilFallback(t *testing.T) {
	// Should not panic when nil writer is provided.
	r := suppress.NewReporter(nil)
	if r == nil {
		t.Error("expected non-nil reporter")
	}
}
