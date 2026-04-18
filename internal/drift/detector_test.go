package drift_test

import (
	"testing"

	"github.com/driftwatch/driftwatch/internal/drift"
)

func TestDetect_NoDrift(t *testing.T) {
	expected := map[string]interface{}{"replicas": 3, "image": "nginx:1.25"}
	actual := map[string]interface{}{"replicas": 3, "image": "nginx:1.25"}

	result := drift.Detect("svc-a", expected, actual)
	if result.HasDrift() {
		t.Errorf("expected no drift, got %d difference(s)", len(result.Diffs))
	}
}

func TestDetect_ValueChanged(t *testing.T) {
	expected := map[string]interface{}{"replicas": 3}
	actual := map[string]interface{}{"replicas": 1}

	result := drift.Detect("svc-b", expected, actual)
	if !result.HasDrift() {
		t.Fatal("expected drift, got none")
	}
	if len(result.Diffs) != 1 {
		t.Fatalf("expected 1 diff, got %d", len(result.Diffs))
	}
	d := result.Diffs[0]
	if d.Field != "replicas" || d.Expected != 3 || d.Actual != 1 {
		t.Errorf("unexpected diff: %v", d)
	}
}

func TestDetect_MissingField(t *testing.T) {
	expected := map[string]interface{}{"timeout": "30s"}
	actual := map[string]interface{}{}

	result := drift.Detect("svc-c", expected, actual)
	if !result.HasDrift() {
		t.Fatal("expected drift for missing field")
	}
	if result.Diffs[0].Actual != nil {
		t.Errorf("expected nil actual for missing field, got %v", result.Diffs[0].Actual)
	}
}

func TestDetect_ServiceName(t *testing.T) {
	result := drift.Detect("my-service", map[string]interface{}{}, map[string]interface{}{})
	if result.ServiceName != "my-service" {
		t.Errorf("expected service name 'my-service', got %q", result.ServiceName)
	}
}

func TestDifference_String(t *testing.T) {
	d := drift.Difference{Field: "image", Expected: "nginx:1.25", Actual: "nginx:1.24"}
	s := d.String()
	if s == "" {
		t.Error("expected non-empty string representation")
	}
}
