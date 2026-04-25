// Package throttle provides a simple rate-limiter that prevents drift checks
// from hammering remote endpoints when many services are configured.
package throttle

import (
	"context"
	"fmt"
	"io"
	"os"
	"sync"
	"time"
)

// Throttle enforces a minimum gap between successive calls across all workers.
type Throttle struct {
	mu       sync.Mutex
	last     time.Time
	interval time.Duration
	w        io.Writer
}

// New returns a Throttle that allows at most one call per interval.
func New(interval time.Duration) *Throttle {
	return NewWithWriter(interval, os.Stderr)
}

// NewWithWriter returns a Throttle that logs wait events to w.
func NewWithWriter(interval time.Duration, w io.Writer) *Throttle {
	return &Throttle{interval: interval, w: w}
}

// Wait blocks until the throttle window has elapsed since the last call,
// or until ctx is cancelled. It returns ctx.Err() on cancellation.
func (t *Throttle) Wait(ctx context.Context) error {
	t.mu.Lock()
	now := time.Now()
	wait := t.interval - now.Sub(t.last)
	if wait > 0 {
		t.last = now.Add(wait)
	} else {
		t.last = now
		wait = 0
	}
	t.mu.Unlock()

	if wait <= 0 {
		return nil
	}

	fmt.Fprintf(t.w, "throttle: waiting %s before next check\n", wait.Round(time.Millisecond))

	select {
	case <-time.After(wait):
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// Reset clears the last-call timestamp so the next Wait returns immediately.
func (t *Throttle) Reset() {
	t.mu.Lock()
	t.last = time.Time{}
	t.mu.Unlock()
}
