package classify_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/example/driftwatch/internal/classify"
	"github.com/example/driftwatch/internal/diff"
)

func TestWrite_NoResults(t *testing.T) {
	var buf bytes.Buffer
	r := classify.NewReporter(&buf)
	r.Write("svc-a", nil)
	if !strings.Contains(buf.String(), "no classified drift") {
		t.Fatalf("unexpected output: %s", buf.String())
	}
}

func TestWrite_WithResults(t *testing.T) {
	var buf bytes.Buffer
	r := classify.NewReporter(&buf)
	c := classify.New()
	res := c.Apply([]diff.Difference{
		{Field: "api_key", Want: "old", Got: "new"},
		{Field: "replicas", Want: "3", Got: "2"},
	})
	r.Write("svc-b", res)
	out := buf.String()
	if !strings.Contains(out, "svc-b") {
		t.Errorf("expected service name in output")
	}
	if !strings.Contains(out, "critical") {
		t.Errorf("expected critical severity in output")
	}
	if !strings.Contains(out, "low") {
		t.Errorf("expected low severity in output")
	}
}

func TestWrite_SortOrder(t *testing.T) {
	var buf bytes.Buffer
	r := classify.NewReporter(&buf)
	res := []classify.Result{
		{Diff: diff.Difference{Field: "x"}, Severity: classify.SeverityLow},
		{Diff: diff.Difference{Field: "y"}, Severity: classify.SeverityCritical},
	}
	r.Write("svc-c", res)
	out := buf.String()
	critIdx := strings.Index(out, "critical")
	lowIdx := strings.Index(out, "low")
	if critIdx > lowIdx {
		t.Error("expected critical to appear before low")
	}
}

func TestNewReporter_NilFallback(t *testing.T) {
	// Should not panic when nil is passed.
	r := classify.NewReporter(nil)
	if r == nil {
		t.Fatal("expected non-nil reporter")
	}
}
