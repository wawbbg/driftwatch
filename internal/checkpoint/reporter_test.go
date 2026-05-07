package checkpoint_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/driftwatch/driftwatch/internal/checkpoint"
)

func TestWrite_NoCheckpoints(t *testing.T) {
	dir := t.TempDir()
	var buf bytes.Buffer
	r := checkpoint.NewReporter(&buf)

	if err := r.Write(dir, "empty-svc"); err != nil {
		t.Fatalf("Write: %v", err)
	}
	if !strings.Contains(buf.String(), "no checkpoints") {
		t.Errorf("expected 'no checkpoints' message, got: %q", buf.String())
	}
}

func TestWrite_WithCheckpoints(t *testing.T) {
	dir := t.TempDir()
	_ = checkpoint.Save(dir, "svc-x", "v1", map[string]string{"img": "alpine"})
	_ = checkpoint.Save(dir, "svc-x", "v2", map[string]string{"img": "alpine:3.18"})

	var buf bytes.Buffer
	r := checkpoint.NewReporter(&buf)

	if err := r.Write(dir, "svc-x"); err != nil {
		t.Fatalf("Write: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "v1") {
		t.Errorf("expected v1 in output, got: %q", out)
	}
	if !strings.Contains(out, "v2") {
		t.Errorf("expected v2 in output, got: %q", out)
	}
	if !strings.Contains(out, "NAME") {
		t.Errorf("expected header NAME in output, got: %q", out)
	}
}

func TestNewReporter_NilFallback(t *testing.T) {
	// Should not panic when nil writer is passed.
	r := checkpoint.NewReporter(nil)
	if r == nil {
		t.Fatal("expected non-nil reporter")
	}
}
