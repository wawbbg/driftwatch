package diff_test

import (
	"testing"

	"github.com/yourorg/driftwatch/internal/diff"
)

func TestMaps_NoDrift(t *testing.T) {
	expected := map[string]interface{}{"replicas": 3, "image": "nginx:1.25"}
	actual := map[string]interface{}{"replicas": 3, "image": "nginx:1.25"}

	results := diff.Maps("svc-a", expected, actual)
	if len(results) != 0 {
		t.Fatalf("expected no drift, got %d result(s)", len(results))
	}
}

func TestMaps_ValueChanged(t *testing.T) {
	expected := map[string]interface{}{"replicas": 3}
	actual := map[string]interface{}{"replicas": 5}

	results := diff.Maps("svc-b", expected, actual)
	if len(results) != 1 {
		t.Fatalf("expected 1 diff, got %d", len(results))
	}
	if results[0].Field != "replicas" {
		t.Errorf("unexpected field: %s", results[0].Field)
	}
	if results[0].Service != "svc-b" {
		t.Errorf("unexpected service: %s", results[0].Service)
	}
}

func TestMaps_MissingField(t *testing.T) {
	expected := map[string]interface{}{"image": "alpine:3.18", "port": 8080}
	actual := map[string]interface{}{"image": "alpine:3.18"}

	results := diff.Maps("svc-c", expected, actual)
	if len(results) != 1 {
		t.Fatalf("expected 1 diff, got %d", len(results))
	}
	if results[0].Actual != nil {
		t.Errorf("expected nil actual for missing field, got %v", results[0].Actual)
	}
}

func TestMaps_ExtraActualFieldIgnored(t *testing.T) {
	expected := map[string]interface{}{"image": "redis:7"}
	actual := map[string]interface{}{"image": "redis:7", "extra": "ignored"}

	results := diff.Maps("svc-d", expected, actual)
	if len(results) != 0 {
		t.Fatalf("extra actual fields should be ignored, got %d diff(s)", len(results))
	}
}

func TestHasDrift(t *testing.T) {
	if diff.HasDrift(nil) {
		t.Error("nil slice should not indicate drift")
	}
	results := []diff.Result{{Service: "x", Field: "f", Expected: 1, Actual: 2}}
	if !diff.HasDrift(results) {
		t.Error("non-empty slice should indicate drift")
	}
}

func TestResult_String(t *testing.T) {
	r := diff.Result{Service: "svc", Field: "replicas", Expected: 2, Actual: 4}
	s := r.String()
	if s == "" {
		t.Error("String() should not return empty string")
	}
}
