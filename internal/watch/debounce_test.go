package watch_test

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/example/driftwatch/internal/watch"
)

func TestDebouncer_CollapsesBurstIntoOne(t *testing.T) {
	var count int64
	d := watch.NewDebouncer(50*time.Millisecond, func(e watch.Event) error {
		atomic.AddInt64(&count, 1)
		return nil
	})

	e := watch.Event{Path: "/cfg.yaml", Type: watch.EventModified, Timestamp: time.Now()}
	for i := 0; i < 10; i++ {
		_ = d.Handle(e)
		time.Sleep(5 * time.Millisecond)
	}

	time.Sleep(120 * time.Millisecond)

	if got := atomic.LoadInt64(&count); got != 1 {
		t.Fatalf("expected 1 invocation, got %d", got)
	}
}

func TestDebouncer_SeparatePathsAreIndependent(t *testing.T) {
	var mu sync.Mutex
	paths := map[string]int{}

	d := watch.NewDebouncer(30*time.Millisecond, func(e watch.Event) error {
		mu.Lock()
		paths[e.Path]++
		mu.Unlock()
		return nil
	})

	_ = d.Handle(watch.Event{Path: "/a.yaml", Type: watch.EventModified, Timestamp: time.Now()})
	_ = d.Handle(watch.Event{Path: "/b.yaml", Type: watch.EventModified, Timestamp: time.Now()})

	time.Sleep(100 * time.Millisecond)

	mu.Lock()
	defer mu.Unlock()
	if paths["/a.yaml"] != 1 {
		t.Errorf("expected 1 call for /a.yaml, got %d", paths["/a.yaml"])
	}
	if paths["/b.yaml"] != 1 {
		t.Errorf("expected 1 call for /b.yaml, got %d", paths["/b.yaml"])
	}
}
