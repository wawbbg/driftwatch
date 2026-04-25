package trend_test

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/example/driftwatch/internal/trend"
)

func makeTrend(service string, samples int, avg float64, dir trend.Direction) trend.ServiceTrend {
	return trend.ServiceTrend{
		Service:   service,
		Samples:   samples,
		AvgDiffs:  avg,
		Direction: dir,
		LastSeen:  time.Date(2024, 1, 15, 9, 0, 0, 0, time.UTC),
	}
}

func TestWrite_NoTrends(t *testing.T) {
	var buf bytes.Buffer
	r := trend.NewReporter(&buf)
	if err := r.Write(nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "no trend data") {
		t.Fatalf("expected 'no trend data', got: %s", buf.String())
	}
}

func TestWrite_WithTrends(t *testing.T) {
	var buf bytes.Buffer
	r := trend.NewReporter(&buf)
	types := []trend.ServiceTrend{
		makeTrend("svc-z", 5, 2.0, trend.DirectionStable),
		makeTrend("svc-a", 3, 4.5, trend.DirectionWorsening),
	}
	if err := r.Write(types); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "SERVICE") {
		t.Error("expected header row")
	}
	if !strings.Contains(out, "svc-a") {
		t.Error("expected svc-a in output")
	}
	if !strings.Contains(out, "worsening") {
		t.Error("expected worsening direction")
	}
	// verify sorted order: svc-a before svc-z
	idxA := strings.Index(out, "svc-a")
	idxZ := strings.Index(out, "svc-z")
	if idxA > idxZ {
		t.Error("expected svc-a before svc-z (sorted)")
	}
}

func TestNewReporter_NilFallback(t *testing.T) {
	// should not panic when w is nil
	r := trend.NewReporter(nil)
	if r == nil {
		t.Fatal("expected non-nil reporter")
	}
}
