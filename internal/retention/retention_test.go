package retention_test

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/driftwatch/driftwatch/internal/retention"
)

func writeFile(t *testing.T, dir, name string, modTime time.Time) string {
	t.Helper()
	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, []byte("data"), 0o644); err != nil {
		t.Fatalf("writeFile: %v", err)
	}
	if err := os.Chtimes(path, modTime, modTime); err != nil {
		t.Fatalf("chtimes: %v", err)
	}
	return path
}

func TestEnforce_RemovesOldFiles(t *testing.T) {
	dir := t.TempDir()
	now := time.Now()

	old := writeFile(t, dir, "old.json", now.Add(-48*time.Hour))
	recent := writeFile(t, dir, "recent.json", now.Add(-1*time.Hour))

	p := retention.Policy{MaxAge: 24 * time.Hour, BaseDir: dir}
	var buf bytes.Buffer
	e := retention.NewWithWriter(p, &buf)
	res := e.Enforce()

	if res.Removed != 1 {
		t.Errorf("expected 1 removed, got %d", res.Removed)
	}
	if res.Errors != 0 {
		t.Errorf("expected 0 errors, got %d", res.Errors)
	}
	if _, err := os.Stat(old); !os.IsNotExist(err) {
		t.Error("expected old file to be removed")
	}
	if _, err := os.Stat(recent); err != nil {
		t.Error("expected recent file to still exist")
	}
}

func TestEnforce_NothingToRemove(t *testing.T) {
	dir := t.TempDir()
	now := time.Now()
	writeFile(t, dir, "new.json", now.Add(-1*time.Minute))

	p := retention.Policy{MaxAge: 24 * time.Hour, BaseDir: dir}
	res := retention.New(p).Enforce()

	if res.Removed != 0 {
		t.Errorf("expected 0 removed, got %d", res.Removed)
	}
}

func TestEnforce_MissingDir(t *testing.T) {
	p := retention.Policy{MaxAge: time.Hour, BaseDir: "/nonexistent/driftwatch/retention"}
	res := retention.New(p).Enforce()

	if res.Removed != 0 || res.Errors != 0 {
		t.Errorf("expected zero result for missing dir, got %+v", res)
	}
}

func TestEnforce_SkipsSubdirectories(t *testing.T) {
	dir := t.TempDir()
	subdir := filepath.Join(dir, "subdir")
	if err := os.Mkdir(subdir, 0o755); err != nil {
		t.Fatal(err)
	}

	p := retention.Policy{MaxAge: time.Millisecond, BaseDir: dir}
	var buf bytes.Buffer
	e := retention.NewWithWriter(p, &buf)
	res := e.Enforce()

	if res.Removed != 0 {
		t.Errorf("expected subdirectory to be skipped, got removed=%d", res.Removed)
	}
}

func TestResult_String(t *testing.T) {
	r := retention.Result{Removed: 3, Errors: 1}
	if got := r.String(); got != "removed=3 errors=1" {
		t.Errorf("unexpected string: %q", got)
	}
}
