package rollup_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/yourorg/driftwatch/internal/diff"
	"github.com/yourorg/driftwatch/internal/rollup"
)

func diffs(fields ...string) []diff.Difference {
	var out []diff.Difference
	for _, f := range fields {
		out = append(out, diff.Difference{Field: f, Expected: "a", Actual: "b"})
	}
	return out
}

func TestNew_SortsByService(t *testing.T) {
	r := rollup.New(map[string][]diff.Difference{
		"zebra": diffs("port"),
		"alpha": diffs(),
		"mango": diffs("image", "replicas"),
	})
	if len(r.Results) != 3 {
		t.Fatalf("expected 3 results, got %d", len(r.Results))
	}
	if r.Results[0].Service != "alpha" {
		t.Errorf("expected first service to be alpha, got %s", r.Results[0].Service)
	}
	if r.Results[2].Service != "zebra" {
		t.Errorf("expected last service to be zebra, got %s", r.Results[2].Service)
	}
}

func TestTotalDrifted(t *testing.T) {
	r := rollup.New(map[string][]diff.Difference{
		"svc-a": diffs("port"),
		"svc-b": diffs(),
		"svc-c": diffs("image"),
	})
	if got := r.TotalDrifted(); got != 2 {
		t.Errorf("expected 2 drifted, got %d", got)
	}
}

func TestHasDrift_False(t *testing.T) {
	r := rollup.New(map[string][]diff.Difference{
		"svc-a": diffs(),
		"svc-b": diffs(),
	})
	if r.HasDrift() {
		t.Error("expected no drift")
	}
}

func TestHasDrift_True(t *testing.T) {
	r := rollup.New(map[string][]diff.Difference{
		"svc-a": diffs("replicas"),
	})
	if !r.HasDrift() {
		t.Error("expected drift")
	}
}

func TestWrite_ContainsHeaders(t *testing.T) {
	r := rollup.New(map[string][]diff.Difference{
		"api": diffs("port"),
		"worker": diffs(),
	})
	var buf bytes.Buffer
	r.Write(&buf)
	out := buf.String()
	for _, want := range []string{"SERVICE", "DRIFTED FIELDS", "STATUS", "api", "DRIFT", "worker", "ok"} {
		if !strings.Contains(out, want) {
			t.Errorf("output missing %q\n%s", want, out)
		}
	}
}

func TestWrite_NilWriterDoesNotPanic(t *testing.T) {
	r := rollup.New(map[string][]diff.Difference{})
	defer func() {
		if rec := recover(); rec != nil {
			t.Errorf("Write panicked: %v", rec)
		}
	}()
	r.Write(nil)
}
