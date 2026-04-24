package metrics_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/yourorg/driftwatch/internal/metrics"
)

func TestNewRun_Fields(t *testing.T) {
	start := time.Now().Add(-50 * time.Millisecond)
	r := metrics.NewRun("payments", 10, 3, start)

	if r.Service != "payments" {
		t.Errorf("expected service=payments, got %s", r.Service)
	}
	if r.Total != 10 {
		t.Errorf("expected total=10, got %d", r.Total)
	}
	if r.Drifted != 3 {
		t.Errorf("expected drifted=3, got %d", r.Drifted)
	}
	if !r.HasDrift {
		t.Error("expected HasDrift=true")
	}
	if r.Duration <= 0 {
		t.Errorf("expected positive duration, got %f", r.Duration)
	}
}

func TestNewRun_NoDrift(t *testing.T) {
	r := metrics.NewRun("auth", 5, 0, time.Now())
	if r.HasDrift {
		t.Error("expected HasDrift=false when drifted=0")
	}
}

func TestRecord_WritesJSON(t *testing.T) {
	var buf bytes.Buffer
	col := metrics.NewWithWriter(&buf)

	r := metrics.NewRun("inventory", 8, 1, time.Now())
	if err := col.Record(r); err != nil {
		t.Fatalf("Record: %v", err)
	}

	line := strings.TrimSpace(buf.String())
	var got metrics.Run
	if err := json.Unmarshal([]byte(line), &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if got.Service != "inventory" {
		t.Errorf("expected service=inventory, got %s", got.Service)
	}
	if got.Drifted != 1 {
		t.Errorf("expected drifted=1, got %d", got.Drifted)
	}
}

func TestSummary_NoDrift(t *testing.T) {
	var buf bytes.Buffer
	col := metrics.NewWithWriter(&buf)
	r := metrics.NewRun("gateway", 6, 0, time.Now())
	col.Summary(r)

	out := buf.String()
	if !strings.Contains(out, "clean") {
		t.Errorf("expected 'clean' in output, got: %s", out)
	}
}

func TestSummary_WithDrift(t *testing.T) {
	var buf bytes.Buffer
	col := metrics.NewWithWriter(&buf)
	r := metrics.NewRun("gateway", 6, 2, time.Now())
	col.Summary(r)

	out := buf.String()
	if !strings.Contains(out, "DRIFT DETECTED") {
		t.Errorf("expected 'DRIFT DETECTED' in output, got: %s", out)
	}
}
