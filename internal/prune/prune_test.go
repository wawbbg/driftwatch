package prune_test

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/driftwatch/driftwatch/internal/prune"
)

func writeFile(t *testing.T, dir, name string, modTime time.Time) {
	t.Helper()
	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, []byte(`{}`), 0o644); err != nil {
		t.Fatalf("writeFile: %v", err)
	}
	if err := os.Chtimes(path, modTime, modTime); err != nil {
		t.Fatalf("chtimes: %v", err)
	}
}

func TestRun_RemovesOldFiles(t *testing.T) {
	base := t.TempDir()
	svcDir := filepath.Join(base, "svc-a")
	os.MkdirAll(svcDir, 0o755)

	old := time.Now().Add(-48 * time.Hour)
	recent := time.Now().Add(-1 * time.Hour)

	writeFile(t, svcDir, "old.json", old)
	writeFile(t, svcDir, "recent.json", recent)

	p := prune.New(24*time.Hour, nil)
	res, err := p.Run(base, "svc-a")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Removed != 1 {
		t.Errorf("removed: got %d, want 1", res.Removed)
	}
	if res.Retained != 1 {
		t.Errorf("retained: got %d, want 1", res.Retained)
	}
	if _, err := os.Stat(filepath.Join(svcDir, "old.json")); !os.IsNotExist(err) {
		t.Error("expected old.json to be removed")
	}
}

func TestRun_ServiceDirMissing(t *testing.T) {
	base := t.TempDir()
	p := prune.New(24*time.Hour, nil)
	res, err := p.Run(base, "ghost-svc")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Removed != 0 || res.Retained != 0 {
		t.Errorf("expected zero counts for missing dir, got %+v", res)
	}
}

func TestRun_WritesLog(t *testing.T) {
	base := t.TempDir()
	svcDir := filepath.Join(base, "svc-b")
	os.MkdirAll(svcDir, 0o755)
	writeFile(t, svcDir, "stale.json", time.Now().Add(-72*time.Hour))

	var buf bytes.Buffer
	p := prune.New(24*time.Hour, &buf)
	p.Run(base, "svc-b") //nolint:errcheck

	if buf.Len() == 0 {
		t.Error("expected log output, got none")
	}
}

func TestResult_String(t *testing.T) {
	r := prune.Result{Service: "my-svc", Removed: 3, Retained: 7}
	got := r.String()
	want := "service=my-svc removed=3 retained=7"
	if got != want {
		t.Errorf("String(): got %q, want %q", got, want)
	}
}
