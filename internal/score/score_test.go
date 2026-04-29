package score_test

import {
	"bytes"
	"testing"

	"github.com/driftwatch/driftwatch/internal/diff"
	"github.com/driftwatch/driftwatch/internal/score"
)

func TestCompute_NoDrift(t *testing.T) {
	expected := map[string]any{"a": 1, "b": 2, "c": 3}
	r := score.Compute("svc", expected, nil)
	if r.Score != 100 {
		t.Fatalf("want 100, got %d", r.Score)
	}
	if r.Grade != score.GradeA {
		t.Fatalf("want A, got %s", r.Grade)
	}
	if r.Drifted != 0 || r.Total != 3 {
		t.Fatalf("unexpected counts: drifted=%d total=%d", r.Drifted, r.Total)
	}
}

func TestCompute_FullDrift(t *testing.T) {
	expected := map[string]any{"a": 1, "b": 2}
	diffs := []diff.Diff{
		{Field: "a", Expected: "1", Actual: "x"},
		{Field: "b", Expected: "2", Actual: "y"},
	}
	r := score.Compute("svc", expected, diffs)
	if r.Score != 0 {
		t.Fatalf("want 0, got %d", r.Score)
	}
	if r.Grade != score.GradeF {
		t.Fatalf("want F, got %s", r.Grade)
	}
}

func TestCompute_PartialDrift(t *testing.T) {
	expected := map[string]any{"a": 1, "b": 2, "c": 3, "d": 4}
	diffs := []diff.Diff{{Field: "a", Expected: "1", Actual: "z"}}
	r := score.Compute("svc", expected, diffs)
	// 3/4 healthy = 75
	if r.Score != 75 {
		t.Fatalf("want 75, got %d", r.Score)
	}
	if r.Grade != score.GradeB {
		t.Fatalf("want B, got %s", r.Grade)
	}
}

func TestCompute_EmptyExpected(t *testing.T) {
	r := score.Compute("svc", map[string]any{}, nil)
	if r.Score != 100 {
		t.Fatalf("want 100 for empty expected, got %d", r.Score)
	}
	if r.Grade != score.GradeA {
		t.Fatalf("want A, got %s", r.Grade)
	}
}

func TestResult_String(t *testing.T) {
	r := score.Result{Service: "api", Score: 80, Grade: score.GradeB, Drifted: 1, Total: 5}
	got := r.String()
	want := "api: score=80 grade=B drifted=1/5"
	if got != want {
		t.Fatalf("want %q, got %q", want, got)
	}
}

func TestReporter_Write(t *testing.T) {
	var buf bytes.Buffer
	rep := score.NewReporter(&buf)
	results := []score.Result{
		{Service: "svc-a", Score: 100, Grade: score.GradeA, Drifted: 0, Total: 4},
		{Service: "svc-b", Score: 50, Grade: score.GradeC, Drifted: 2, Total: 4},
	}
	rep.Write(results)
	out := buf.String()
	if out == "" {
		t.Fatal("expected non-empty output")
	}
	for _, res := range results {
		if !bytes.Contains(buf.Bytes(), []byte(res.Service)) {
			t.Errorf("output missing service %q", res.Service)
		}
	}
}

func TestNewReporter_NilFallback(t *testing.T) {
	// Should not panic when nil writer is provided.
	rep := score.NewReporter(nil)
	if rep == nil {
		t.Fatal("expected non-nil reporter")
	}
}
