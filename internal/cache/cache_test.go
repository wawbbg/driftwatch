package cache_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/yourorg/driftwatch/internal/cache"
)

func TestGet_MissOnEmpty(t *testing.T) {
	c := cache.New("")
	_, ok := c.Get("svc-a")
	if ok {
		t.Fatal("expected miss on empty cache")
	}
}

func TestSetAndGet_InMemory(t *testing.T) {
	c := cache.New("")
	val := map[string]any{"replicas": 3}

	if err := c.Set("svc-a", val, time.Minute); err != nil {
		t.Fatalf("Set: %v", err)
	}

	e, ok := c.Get("svc-a")
	if !ok {
		t.Fatal("expected hit after Set")
	}
	if e.Value["replicas"] != 3 {
		t.Errorf("got %v, want 3", e.Value["replicas"])
	}
}

func TestGet_ExpiredEntry(t *testing.T) {
	c := cache.New("")
	if err := c.Set("svc-b", map[string]any{"env": "prod"}, -time.Second); err != nil {
		t.Fatalf("Set: %v", err)
	}
	_, ok := c.Get("svc-b")
	if ok {
		t.Fatal("expected miss for expired entry")
	}
}

func TestSetAndGet_Persisted(t *testing.T) {
	dir := t.TempDir()
	c := cache.New(dir)
	val := map[string]any{"image": "nginx:1.25"}

	if err := c.Set("svc-c", val, time.Hour); err != nil {
		t.Fatalf("Set: %v", err)
	}

	// New cache instance reads from disk.
	c2 := cache.New(dir)
	e, ok := c2.Get("svc-c")
	if !ok {
		t.Fatal("expected hit from disk cache")
	}
	if e.Value["image"] != "nginx:1.25" {
		t.Errorf("got %v, want nginx:1.25", e.Value["image"])
	}
}

func TestSet_BadDir(t *testing.T) {
	c := cache.New(filepath.Join(t.TempDir(), "no", "such", "\x00"))
	err := c.Set("svc-d", map[string]any{}, time.Minute)
	if err == nil {
		t.Fatal("expected error for bad dir")
	}
}

func TestExpired_FileOnDisk_ReturnsMiss(t *testing.T) {
	dir := t.TempDir()
	c := cache.New(dir)

	// Write an already-expired entry.
	if err := c.Set("svc-e", map[string]any{"ok": true}, -time.Second); err != nil {
		t.Fatalf("Set: %v", err)
	}

	// Confirm file exists on disk.
	if _, err := os.Stat(filepath.Join(dir, "svc-e.json")); err != nil {
		t.Fatalf("expected file on disk: %v", err)
	}

	c2 := cache.New(dir)
	_, ok := c2.Get("svc-e")
	if ok {
		t.Fatal("expected miss for on-disk expired entry")
	}
}
