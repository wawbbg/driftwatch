package compare_test

import (
	"strings"
	"testing"

	"github.com/driftwatch/internal/compare"
)

func TestSnapshots_NoDrift(t *testing.T) {
	before := map[string]any{"replicas": "3", "image": "nginx:1.25"}
	after := map[string]any{"replicas": "3", "image": "nginx:1.25"}

	d := compare.Snapshots("svc-a", before, after)

	if d.HasDrift() {
		t.Fatalf("expected no drift, got %d diffs", len(d.Diffs))
	}
	if d.Service != "svc-a" {
		t.Errorf("unexpected service name: %s", d.Service)
	}
	if d.CapturedAt.IsZero() {
		t.Error("CapturedAt should not be zero")
	}
}

func TestSnapshots_ValueChanged(t *testing.T) {
	before := map[string]any{"replicas": "3"}
	after := map[string]any{"replicas": "5"}

	d := compare.Snapshots("svc-b", before, after)

	if !d.HasDrift() {
		t.Fatal("expected drift but none found")
	}
	if len(d.Diffs) != 1 {
		t.Fatalf("expected 1 diff, got %d", len(d.Diffs))
	}
}

func TestSnapshots_MissingField(t *testing.T) {
	before := map[string]any{"replicas": "3", "timeout": "30s"}
	after := map[string]any{"replicas": "3"}

	d := compare.Snapshots("svc-c", before, after)

	if !d.HasDrift() {
		t.Fatal("expected drift for missing field")
	}
}

func TestDelta_String_NoDrift(t *testing.T) {
	d := compare.Snapshots("svc-x", map[string]any{"k": "v"}, map[string]any{"k": "v"})
	if !strings.Contains(d.String(), "no drift") {
		t.Errorf("unexpected string: %s", d.String())
	}
}

func TestDelta_String_WithDrift(t *testing.T) {
	d := compare.Snapshots("svc-y", map[string]any{"k": "v1"}, map[string]any{"k": "v2"})
	if !strings.Contains(d.String(), "changed") {
		t.Errorf("unexpected string: %s", d.String())
	}
}
