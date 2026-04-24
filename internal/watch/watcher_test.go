package watch_test

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/example/driftwatch/internal/watch"
)

func TestStart_DetectsModified(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "cfg.yaml")
	_ = os.WriteFile(p, []byte("v: 1"), 0o644)

	var mu sync.Mutex
	var events []watch.Event

	w := watch.NewWithWriter([]string{p}, 20*time.Millisecond, func(e watch.Event) error {
		mu.Lock()
		events = append(events, e)
		mu.Unlock()
		return nil
	}, &bytes.Buffer{})

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	go func() { done <- w.Start(ctx) }()

	time.Sleep(40 * time.Millisecond)
	_ = os.WriteFile(p, []byte("v: 2"), 0o644)
	time.Sleep(60 * time.Millisecond)
	cancel()
	<-done

	mu.Lock()
	defer mu.Unlock()
	if len(events) == 0 {
		t.Fatal("expected at least one event, got none")
	}
	if events[0].Type != watch.EventModified {
		t.Fatalf("expected modified, got %s", events[0].Type)
	}
}

func TestStart_DetectsCreated(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "new.yaml")

	var mu sync.Mutex
	var events []watch.Event

	w := watch.NewWithWriter([]string{p}, 20*time.Millisecond, func(e watch.Event) error {
		mu.Lock()
		events = append(events, e)
		mu.Unlock()
		return nil
	}, &bytes.Buffer{})

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	go func() { done <- w.Start(ctx) }()

	time.Sleep(40 * time.Millisecond)
	_ = os.WriteFile(p, []byte("v: 1"), 0o644)
	time.Sleep(60 * time.Millisecond)
	cancel()
	<-done

	mu.Lock()
	defer mu.Unlock()
	if len(events) == 0 {
		t.Fatal("expected created event")
	}
	if events[0].Type != watch.EventCreated {
		t.Fatalf("expected created, got %s", events[0].Type)
	}
}

func TestStart_CancelReturnsCtxErr(t *testing.T) {
	w := watch.New([]string{}, 50*time.Millisecond, func(watch.Event) error { return nil })
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	go func() { done <- w.Start(ctx) }()
	cancel()
	if err := <-done; err != context.Canceled {
		t.Fatalf("expected context.Canceled, got %v", err)
	}
}

func TestEvent_String(t *testing.T) {
	e := watch.Event{Path: "/etc/cfg.yaml", Type: watch.EventModified, Timestamp: time.Time{}}
	s := e.String()
	if s == "" {
		t.Fatal("expected non-empty string")
	}
}
