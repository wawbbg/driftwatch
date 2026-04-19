package alert_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/example/driftwatch/internal/alert"
	"github.com/example/driftwatch/internal/drift"
)

func diffs(pairs ...string) []drift.Difference {
	var out []drift.Difference
	for i := 0; i+2 < len(pairs); i += 3 {
		out = append(out, drift.Difference{
			Field:    pairs[i],
			Expected: pairs[i+1],
			Actual:   pairs[i+2],
		})
	}
	return out
}

func TestNotify_NoDiffs(t *testing.T) {
	var buf bytes.Buffer
	n := alert.NewWithWriter(&buf)
	count := n.Notify("svc", nil)
	if count != 0 {
		t.Fatalf("expected 0 alerts, got %d", count)
	}
	if buf.Len() != 0 {
		t.Fatalf("expected no output, got %q", buf.String())
	}
}

func TestNotify_WithDiffs(t *testing.T) {
	var buf bytes.Buffer
	n := alert.NewWithWriter(&buf)
	ds := diffs("replicas", "3", "1", "image", "v2", "v1")
	count := n.Notify("api", ds)
	if count != 2 {
		t.Fatalf("expected 2 alerts, got %d", count)
	}
	out := buf.String()
	if !strings.Contains(out, "api") {
		t.Errorf("expected service name in output")
	}
	if !strings.Contains(out, "WARN") {
		t.Errorf("expected WARN level in output")
	}
}

func TestNotify_MissingExpected_IsError(t *testing.T) {
	var buf bytes.Buffer
	n := alert.NewWithWriter(&buf)
	ds := []drift.Difference{{Field: "port", Expected: "", Actual: "9090"}}
	n.Notify("gateway", ds)
	if !strings.Contains(buf.String(), "ERROR") {
		t.Errorf("expected ERROR level when Expected is empty")
	}
}

func TestAlert_String(t *testing.T) {
	a := alert.Alert{Service: "db", Level: alert.LevelWarn, Message: "field changed"}
	s := a.String()
	if !strings.HasPrefix(s, "[WARN]") {
		t.Errorf("unexpected format: %q", s)
	}
}
