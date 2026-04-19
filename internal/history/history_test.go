package history_test

import (
	"os"
	"testing"

	"github.com/driftwatch/driftwatch/internal/history"
)

func TestRecordAndList(t *testing.T) {
	dir := t.TempDir()

	if err := history.Record(dir, "svc-a", 3); err != nil {
		t.Fatalf("Record: %v", err)
	}
	if err := history.Record(dir, "svc-a", 0); err != nil {
		t.Fatalf("Record: %v", err)
	}

	entries, err := history.List(dir, "svc-a")
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
	if entries[0].DriftCount != 3 || !entries[0].HasDrift {
		t.Errorf("first entry: unexpected values %+v", entries[0])
	}
	if entries[1].DriftCount != 0 || entries[1].HasDrift {
		t.Errorf("second entry: unexpected values %+v", entries[1])
	}
	if entries[0].ServiceName != "svc-a" {
		t.Errorf("expected service name svc-a, got %s", entries[0].ServiceName)
	}
}

func TestList_NotFound(t *testing.T) {
	dir := t.TempDir()
	entries, err := history.List(dir, "missing")
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if entries != nil {
		t.Errorf("expected nil entries, got %v", entries)
	}
}

func TestRecord_BadDir(t *testing.T) {
	// Use a file as the directory to force an error.
	f, err := os.CreateTemp("", "notadir")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(f.Name())
	f.Close()

	err = history.Record(f.Name(), "svc", 1)
	if err == nil {
		t.Fatal("expected error when dir is a file")
	}
}

func TestRecord_IsolatedByService(t *testing.T) {
	dir := t.TempDir()
	_ = history.Record(dir, "svc-a", 1)
	_ = history.Record(dir, "svc-b", 5)

	a, _ := history.List(dir, "svc-a")
	b, _ := history.List(dir, "svc-b")

	if len(a) != 1 || a[0].DriftCount != 1 {
		t.Errorf("svc-a isolation failed: %+v", a)
	}
	if len(b) != 1 || b[0].DriftCount != 5 {
		t.Errorf("svc-b isolation failed: %+v", b)
	}
}
