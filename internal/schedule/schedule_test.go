package schedule_test

import (
	"bytes"
	"context"
	"errors"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/example/driftwatch/internal/schedule"
)

func TestStart_RunsImmediately(t *testing.T) {
	var count int32
	var buf bytes.Buffer

	job := func(ctx context.Context) error {
		atomic.AddInt32(&count, 1)
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	r := schedule.NewWithWriter(1*time.Hour, job, &buf)
	r.Start(ctx) // nolint: errcheck

	if atomic.LoadInt32(&count) < 1 {
		t.Fatal("expected job to run at least once immediately")
	}
}

func TestStart_RunsOnInterval(t *testing.T) {
	var count int32
	var buf bytes.Buffer

	job := func(ctx context.Context) error {
		atomic.AddInt32(&count, 1)
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 180*time.Millisecond)
	defer cancel()

	r := schedule.NewWithWriter(60*time.Millisecond, job, &buf)
	r.Start(ctx) // nolint: errcheck

	if atomic.LoadInt32(&count) < 2 {
		t.Fatalf("expected at least 2 runs, got %d", atomic.LoadInt32(&count))
	}
}

func TestStart_JobErrorLogged(t *testing.T) {
	var buf bytes.Buffer

	job := func(ctx context.Context) error {
		return errors.New("boom")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Millisecond)
	defer cancel()

	r := schedule.NewWithWriter(1*time.Hour, job, &buf)
	r.Start(ctx) // nolint: errcheck

	if !strings.Contains(buf.String(), "boom") {
		t.Errorf("expected error to be logged, got: %s", buf.String())
	}
}

func TestStart_CancelReturnsCtxErr(t *testing.T) {
	var buf bytes.Buffer

	job := func(ctx context.Context) error { return nil }

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	r := schedule.NewWithWriter(10*time.Millisecond, job, &buf)
	err := r.Start(ctx)

	if !errors.Is(err, context.Canceled) {
		t.Errorf("expected context.Canceled, got %v", err)
	}
}
