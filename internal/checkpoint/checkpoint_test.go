package checkpoint_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/driftwatch/driftwatch/internal/checkpoint"
)

func TestSaveAndLoad(t *testing.T) {
	dir := t.TempDir()
	fields := map[string]string{"replicas": "3", "image": "nginx:1.25"}

	if err := checkpoint.Save(dir, "svc-a", "pre-deploy", fields); err != nil {
		t.Fatalf("Save: %v", err)
	}

	e, err := checkpoint.Load(dir, "svc-a", "pre-deploy")
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if e.Service != "svc-a" {
		t.Errorf("service = %q, want svc-a", e.Service)
	}
	if e.Fields["replicas"] != "3" {
		t.Errorf("replicas = %q, want 3", e.Fields["replicas"])
	}
	if e.CreatedAt.IsZero() {
		t.Error("CreatedAt should not be zero")
	}
}

func TestLoad_NotFound(t *testing.T) {
	dir := t.TempDir()
	_, err := checkpoint.Load(dir, "svc-a", "missing")
	if err == nil {
		t.Fatal("expected error for missing checkpoint")
	}
}

func TestSave_BadDir(t *testing.T) {
	err := checkpoint.Save("/proc/nonexistent/bad", "svc", "cp", nil)
	if err == nil {
		t.Fatal("expected error for bad dir")
	}
}

func TestSave_EmptyService(t *testing.T) {
	dir := t.TempDir()
	err := checkpoint.Save(dir, "", "cp", nil)
	if err == nil {
		t.Fatal("expected error for empty service")
	}
}

func TestList(t *testing.T) {
	dir := t.TempDir()
	for _, name := range []string{"alpha", "beta", "gamma"} {
		if err := checkpoint.Save(dir, "svc-b", name, map[string]string{"k": "v"}); err != nil {
			t.Fatalf("Save %s: %v", name, err)
		}
	}

	names, err := checkpoint.List(dir, "svc-b")
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(names) != 3 {
		t.Errorf("got %d names, want 3", len(names))
	}
}

func TestList_NotFound(t *testing.T) {
	dir := t.TempDir()
	names, err := checkpoint.List(dir, "no-such-service")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(names) != 0 {
		t.Errorf("expected empty list, got %v", names)
	}
}

func TestList_SkipsNonJSON(t *testing.T) {
	dir := t.TempDir()
	svcDir := filepath.Join(dir, "svc-c")
	_ = os.MkdirAll(svcDir, 0o755)
	_ = os.WriteFile(filepath.Join(svcDir, "notes.txt"), []byte("hi"), 0o644)
	_ = checkpoint.Save(dir, "svc-c", "real", map[string]string{})

	names, err := checkpoint.List(dir, "svc-c")
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(names) != 1 || names[0] != "real" {
		t.Errorf("got %v, want [real]", names)
	}
}
