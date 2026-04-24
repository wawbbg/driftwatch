package audit_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/user/driftwatch/internal/audit"
)

func TestRecordAndList(t *testing.T) {
	dir := t.TempDir()

	if err := audit.Record(dir, "svc-a", 3, false, "first run"); err != nil {
		t.Fatalf("Record: %v", err)
	}
	if err := audit.Record(dir, "svc-a", 0, false, "clean"); err != nil {
		t.Fatalf("Record: %v", err)
	}

	entries, err := audit.List(dir, "svc-a")
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
	if entries[0].DriftCount != 3 {
		t.Errorf("expected drift_count=3, got %d", entries[0].DriftCount)
	}
	if entries[1].Message != "clean" {
		t.Errorf("expected message=clean, got %q", entries[1].Message)
	}
}

func TestList_NotFound(t *testing.T) {
	dir := t.TempDir()
	entries, err := audit.List(dir, "nonexistent")
	if err != nil {
		t.Fatalf("expected no error for missing file, got %v", err)
	}
	if len(entries) != 0 {
		t.Errorf("expected empty slice, got %d entries", len(entries))
	}
}

func TestRecord_BadDir(t *testing.T) {
	// Use a file as the directory to force an error.
	f, err := os.CreateTemp("", "audit-baddir")
	if err != nil {
		t.Fatal(err)
	}
	f.Close()
	defer os.Remove(f.Name())

	err = audit.Record(f.Name(), "svc", 1, false, "")
	if err == nil {
		t.Fatal("expected error when dir is a file")
	}
}

func TestRecord_IsolatedByService(t *testing.T) {
	dir := t.TempDir()

	_ = audit.Record(dir, "alpha", 1, false, "")
	_ = audit.Record(dir, "beta", 2, true, "policy")

	alpha, _ := audit.List(dir, "alpha")
	beta, _ := audit.List(dir, "beta")

	if len(alpha) != 1 || alpha[0].Service != "alpha" {
		t.Errorf("unexpected alpha entries: %+v", alpha)
	}
	if len(beta) != 1 || !beta[0].PolicyError {
		t.Errorf("unexpected beta entries: %+v", beta)
	}
}

func TestRecord_FileCreated(t *testing.T) {
	dir := t.TempDir()
	_ = audit.Record(dir, "mysvc", 0, false, "ok")

	path := filepath.Join(dir, "mysvc.audit.jsonl")
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Errorf("expected audit file to be created at %s", path)
	}
}
