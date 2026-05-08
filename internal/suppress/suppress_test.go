package suppress_test

import (
	"os"
	"testing"
	"time"

	"github.com/driftwatch/driftwatch/internal/suppress"
)

var (
	now    = time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC)
	future = now.Add(24 * time.Hour)
	past   = now.Add(-1 * time.Hour)
)

func TestStoreAndActive(t *testing.T) {
	dir := t.TempDir()
	e := suppress.Entry{
		Service:   "api",
		Field:     "replicas",
		Reason:    "scaling event",
		ExpiresAt: future,
	}
	if err := suppress.Store(dir, e); err != nil {
		t.Fatalf("Store: %v", err)
	}
	entries, err := suppress.Active(dir, "api", now)
	if err != nil {
		t.Fatalf("Active: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 active entry, got %d", len(entries))
	}
	if entries[0].Field != "replicas" {
		t.Errorf("unexpected field: %s", entries[0].Field)
	}
}

func TestActive_ExpiredEntryExcluded(t *testing.T) {
	dir := t.TempDir()
	e := suppress.Entry{Service: "api", Field: "image", ExpiresAt: past}
	_ = suppress.Store(dir, e)
	entries, err := suppress.Active(dir, "api", now)
	if err != nil {
		t.Fatalf("Active: %v", err)
	}
	if len(entries) != 0 {
		t.Errorf("expected no active entries, got %d", len(entries))
	}
}

func TestActive_NotFound(t *testing.T) {
	dir := t.TempDir()
	entries, err := suppress.Active(dir, "missing", now)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if entries != nil {
		t.Errorf("expected nil, got %v", entries)
	}
}

func TestIsSuppressed_True(t *testing.T) {
	dir := t.TempDir()
	e := suppress.Entry{Service: "svc", Field: "timeout", ExpiresAt: future}
	_ = suppress.Store(dir, e)
	ok, err := suppress.IsSuppressed(dir, "svc", "timeout", now)
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Error("expected field to be suppressed")
	}
}

func TestIsSuppressed_Wildcard(t *testing.T) {
	dir := t.TempDir()
	e := suppress.Entry{Service: "svc", Field: "*", ExpiresAt: future}
	_ = suppress.Store(dir, e)
	ok, err := suppress.IsSuppressed(dir, "svc", "any-field", now)
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Error("expected wildcard suppression to match")
	}
}

func TestStore_EmptyService(t *testing.T) {
	dir := t.TempDir()
	err := suppress.Store(dir, suppress.Entry{Field: "x", ExpiresAt: future})
	if err == nil {
		t.Error("expected error for empty service")
	}
}

func TestStore_BadDir(t *testing.T) {
	// Use a file as the directory to force a failure.
	f, _ := os.CreateTemp("", "suppress")
	_ = f.Close()
	defer os.Remove(f.Name())
	err := suppress.Store(f.Name(), suppress.Entry{Service: "svc", Field: "x", ExpiresAt: future})
	if err == nil {
		t.Error("expected error when dir is a file")
	}
}
