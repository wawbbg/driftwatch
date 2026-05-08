package classify_test

import (
	"testing"

	"github.com/example/driftwatch/internal/classify"
	"github.com/example/driftwatch/internal/diff"
)

func diffs() []diff.Difference {
	return []diff.Difference{
		{Field: "api_token", Want: "abc", Got: "xyz"},
		{Field: "host", Want: "prod.example.com", Got: "staging.example.com"},
		{Field: "replicas", Want: "3", Got: "2"},
		{Field: "log_level", Want: "info", Got: ""},
	}
}

func TestClassify_Critical(t *testing.T) {
	c := classify.New()
	r := c.Classify(diff.Difference{Field: "api_token", Want: "a", Got: "b"})
	if r.Severity != classify.SeverityCritical {
		t.Fatalf("expected critical, got %s", r.Severity)
	}
}

func TestClassify_High(t *testing.T) {
	c := classify.New()
	r := c.Classify(diff.Difference{Field: "db_host", Want: "a", Got: "b"})
	if r.Severity != classify.SeverityHigh {
		t.Fatalf("expected high, got %s", r.Severity)
	}
}

func TestClassify_Medium_MissingValue(t *testing.T) {
	c := classify.New()
	r := c.Classify(diff.Difference{Field: "log_level", Want: "info", Got: ""})
	if r.Severity != classify.SeverityMedium {
		t.Fatalf("expected medium, got %s", r.Severity)
	}
}

func TestClassify_Low(t *testing.T) {
	c := classify.New()
	r := c.Classify(diff.Difference{Field: "replicas", Want: "3", Got: "2"})
	if r.Severity != classify.SeverityLow {
		t.Fatalf("expected low, got %s", r.Severity)
	}
}

func TestApply_ReturnsAllResults(t *testing.T) {
	c := classify.New()
	res := c.Apply(diffs())
	if len(res) != 4 {
		t.Fatalf("expected 4 results, got %d", len(res))
	}
}

func TestHasCritical_True(t *testing.T) {
	c := classify.New()
	res := c.Apply(diffs())
	if !classify.HasCritical(res) {
		t.Fatal("expected HasCritical to be true")
	}
}

func TestHasCritical_False(t *testing.T) {
	res := []classify.Result{
		{Severity: classify.SeverityLow},
		{Severity: classify.SeverityMedium},
	}
	if classify.HasCritical(res) {
		t.Fatal("expected HasCritical to be false")
	}
}

func TestResult_String(t *testing.T) {
	r := classify.Result{
		Diff:     diff.Difference{Field: "host", Want: "a", Got: "b"},
		Severity: classify.SeverityHigh,
	}
	s := r.String()
	if s == "" {
		t.Fatal("expected non-empty string")
	}
}

func TestNewWithFields_CustomCritical(t *testing.T) {
	c := classify.NewWithFields([]string{"myfield"}, nil)
	r := c.Classify(diff.Difference{Field: "myfield", Want: "a", Got: "b"})
	if r.Severity != classify.SeverityCritical {
		t.Fatalf("expected critical, got %s", r.Severity)
	}
}
