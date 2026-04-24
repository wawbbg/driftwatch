package compare_test

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/driftwatch/internal/compare"
)

func writeSnap(t *testing.T, dir, name string, data map[string]any) {
	t.Helper()
	b, err := json.Marshal(data)
	if err != nil {
		t.Fatalf("marshal snapshot: %v", err)
	}
	if err := os.WriteFile(filepath.Join(dir, name+".json"), b, 0o644); err != nil {
		t.Fatalf("write snapshot: %v", err)
	}
}

func TestRunner_NoDrift(t *testing.T) {
	dir := t.TempDir()
	data := map[string]any{"replicas": "2"}
	writeSnap(t, dir, "svc.v1", data)
	writeSnap(t, dir, "svc.v2", data)

	var buf bytes.Buffer
	r := compare.NewRunnerWithWriter(dir, &buf)
	delta, err := r.Run(context.Background(), "svc", "v1", "v2")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if delta.HasDrift() {
		t.Error("expected no drift")
	}
}

func TestRunner_WithDrift(t *testing.T) {
	dir := t.TempDir()
	writeSnap(t, dir, "svc.v1", map[string]any{"replicas": "2"})
	writeSnap(t, dir, "svc.v2", map[string]any{"replicas": "4"})

	var buf bytes.Buffer
	r := compare.NewRunnerWithWriter(dir, &buf)
	delta, err := r.Run(context.Background(), "svc", "v1", "v2")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !delta.HasDrift() {
		t.Error("expected drift")
	}
}

func TestRunner_MissingSnapshot(t *testing.T) {
	dir := t.TempDir()
	r := compare.NewRunner(dir)
	_, err := r.Run(context.Background(), "svc", "v1", "v2")
	if err == nil {
		t.Fatal("expected error for missing snapshot")
	}
}
