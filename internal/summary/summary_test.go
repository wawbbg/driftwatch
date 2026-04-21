package summary_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/example/driftwatch/internal/drift"
	"github.com/example/driftwatch/internal/summary"
)

func diffs(pairs ...string) []drift.Difference {
	var out []drift.Difference
	for i := 0; i+1 < len(pairs); i += 2 {
		out = append(out, drift.Difference{
			Field:    pairs[i],
			Expected: pairs[i+1],
			Actual:   "actual-" + pairs[i],
		})
	}
	return out
}

func TestNew_Counts(t *testing.T) {
	input := map[string][]drift.Difference{
		"svc-a": diffs("replicas", "3"),
		"svc-b": {},
		"svc-c": diffs("image", "v1", "env", "prod"),
	}
	r := summary.New(input)
	if r.Total != 3 {
		t.Errorf("Total: got %d, want 3", r.Total)
	}
	if r.Drifted != 2 {
		t.Errorf("Drifted: got %d, want 2", r.Drifted)
	}
	if r.Clean != 1 {
		t.Errorf("Clean: got %d, want 1", r.Clean)
	}
}

func TestNew_Empty(t *testing.T) {
	r := summary.New(map[string][]drift.Difference{})
	if r.Total != 0 || r.Drifted != 0 || r.Clean != 0 {
		t.Errorf("expected all zeros for empty input, got %+v", r)
	}
}

func TestWrite_NoDrift(t *testing.T) {
	r := summary.New(map[string][]drift.Difference{
		"svc-a": {},
	})
	var buf bytes.Buffer
	summary.Write(&buf, r)
	out := buf.String()
	if !strings.Contains(out, "Clean            : 1") {
		t.Errorf("expected clean count in output, got:\n%s", out)
	}
	if strings.Contains(out, "Affected services") {
		t.Errorf("should not list affected services when none drifted")
	}
}

func TestWrite_WithDrift(t *testing.T) {
	r := summary.New(map[string][]drift.Difference{
		"svc-x": diffs("replicas", "2"),
	})
	var buf bytes.Buffer
	summary.Write(&buf, r)
	out := buf.String()
	if !strings.Contains(out, "svc-x") {
		t.Errorf("expected service name in output, got:\n%s", out)
	}
	if !strings.Contains(out, "replicas") {
		t.Errorf("expected field name in output, got:\n%s", out)
	}
}

func TestWrite_NilWriter(t *testing.T) {
	// should not panic when w is nil (falls back to os.Stdout)
	r := summary.New(map[string][]drift.Difference{})
	defer func() {
		if rec := recover(); rec != nil {
			t.Errorf("Write panicked with nil writer: %v", rec)
		}
	}()
	summary.Write(nil, r)
}
