package throttle_test

import (
	"bytes"
	"context"
	"testing"
	"time"

	"github.com/example/driftwatch/internal/throttle"
)

func TestWait_FirstCallImmediate(t *testing.T) {
	th := throttle.New(100 * time.Millisecond)
	start := time.Now()
	if err := th.Wait(context.Background()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if elapsed := time.Since(start); elapsed > 20*time.Millisecond {
		t.Errorf("first call should be immediate, took %s", elapsed)
	}
}

func TestWait_SecondCallThrottled(t *testing.T) {
	var buf bytes.Buffer
	th := throttle.NewWithWriter(80*time.Millisecond, &buf)

	_ = th.Wait(context.Background())

	start := time.Now()
	_ = th.Wait(context.Background())
	elapsed := time.Since(start)

	if elapsed < 60*time.Millisecond {
		t.Errorf("second call should have been throttled, elapsed=%s", elapsed)
	}
	if buf.Len() == 0 {
		t.Error("expected throttle log message, got none")
	}
}

func TestWait_CancelledContext(t *testing.T) {
	th := throttle.New(500 * time.Millisecond)
	_ = th.Wait(context.Background()) // prime the timer

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
	defer cancel()

	err := th.Wait(ctx)
	if err == nil {
		t.Fatal("expected context error, got nil")
	}
}

func TestReset_AllowsImmediateCall(t *testing.T) {
	th := throttle.New(500 * time.Millisecond)
	_ = th.Wait(context.Background())
	th.Reset()

	start := time.Now()
	_ = th.Wait(context.Background())
	if elapsed := time.Since(start); elapsed > 20*time.Millisecond {
		t.Errorf("call after Reset should be immediate, took %s", elapsed)
	}
}

func TestWait_ConcurrentSafe(t *testing.T) {
	th := throttle.New(5 * time.Millisecond)
	ctx := context.Background()
	for i := 0; i < 5; i++ {
		go func() { _ = th.Wait(ctx) }()
	}
	// no race detector hit is the success condition
	time.Sleep(60 * time.Millisecond)
}
