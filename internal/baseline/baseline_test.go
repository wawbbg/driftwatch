package baseline_test

import (
	"os"
	"testing"

	"github.com/yourorg/driftwatch/internal/baseline"
)

func TestStoreAndLoad(t *testing.T) {
	dir := t.TempDir()
	fields := map[string]interface{}{"replicas": 3, "image": "nginx:1.25"}

	if err := baseline.Store(dir, "web", fields); err != nil {
		t.Fatalf("Store: %v", err)
	}

	e, err := baseline.Load(dir, "web")
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if e.ServiceName != "web" {
		t.Errorf("service name: got %q, want %q", e.ServiceName, "web")
	}
	if e.Fields["image"] != "nginx:1.25" {
		t.Errorf("image field: got %v", e.Fields["image"])
	}
	if e.CapturedAt.IsZero() {
		t.Error("CapturedAt should not be zero")
	}
}

func TestLoad_NotFound(t *testing.T) {
	dir := t.TempDir()
	_, err := baseline.Load(dir, "missing")
	if err == nil {
		t.Fatal("expected error for missing baseline")
	}
}

func TestStore_BadDir(t *testing.T) {
	// Use a file as the directory to force failure.
	f, err := os.CreateTemp("", "bl")
	if err != nil {
		t.Fatal(err)
	}
	f.Close()
	defer os.Remove(f.Name())

	err = baseline.Store(f.Name(), "svc", map[string]interface{}{})
	if err == nil {
		t.Fatal("expected error writing into a file path as dir")
	}
}

func TestLoad_CorruptJSON(t *testing.T) {
	dir := t.TempDir()
	path := dir + "/bad.json"
	if err := os.WriteFile(path, []byte("not json{"), 0o644); err != nil {
		t.Fatal(err)
	}
	_, err := baseline.Load(dir, "bad")
	if err == nil {
		t.Fatal("expected unmarshal error")
	}
}
