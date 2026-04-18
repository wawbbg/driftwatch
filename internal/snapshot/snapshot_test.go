package snapshot_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/user/driftwatch/internal/drift"
	"github.com/user/driftwatch/internal/snapshot"
)

func TestNew(t *testing.T) {
	before := time.Now().UTC()
	diffs := []drift.Difference{{Field: "replicas", Expected: "3", Actual: "2"}}
	s := snapshot.New("my-service", diffs)
	if s.Service != "my-service" {
		t.Errorf("expected service %q, got %q", "my-service", s.Service)
	}
	if len(s.Diffs) != 1 {
		t.Fatalf("expected 1 diff, got %d", len(s.Diffs))
	}
	if s.Timestamp.Before(before) {
		t.Error("timestamp should not be before test start")
	}
}

func TestSaveAndLoad(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "snap.json")

	orig := snapshot.Snapshot{
		Timestamp: time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC),
		Service:   "api",
		Diffs: []drift.Difference{
			{Field: "image", Expected: "v1.2", Actual: "v1.1"},
		},
	}

	if err := snapshot.Save(path, orig); err != nil {
		t.Fatalf("Save: %v", err)
	}

	loaded, err := snapshot.Load(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}

	if loaded.Service != orig.Service {
		t.Errorf("service: want %q got %q", orig.Service, loaded.Service)
	}
	if len(loaded.Diffs) != len(orig.Diffs) {
		t.Fatalf("diffs len: want %d got %d", len(orig.Diffs), len(loaded.Diffs))
	}
	if loaded.Diffs[0].Field != orig.Diffs[0].Field {
		t.Errorf("diff field: want %q got %q", orig.Diffs[0].Field, loaded.Diffs[0].Field)
	}
}

func TestLoad_NotFound(t *testing.T) {
	_, err := snapshot.Load("/nonexistent/snap.json")
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestSave_BadPath(t *testing.T) {
	s := snapshot.New("svc", nil)
	err := snapshot.Save("/nonexistent/dir/snap.json", s)
	if err == nil {
		t.Error("expected error for bad path")
	}
	os.Remove("/nonexistent/dir/snap.json")
}
