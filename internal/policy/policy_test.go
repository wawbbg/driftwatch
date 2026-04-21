package policy_test

import (
	"strings"
	"testing"

	"github.com/example/driftwatch/internal/drift"
	"github.com/example/driftwatch/internal/policy"
)

func diffs() []drift.Difference {
	return []drift.Difference{
		{Service: "svc", Field: "replicas", Expected: "3", Actual: "1"},
		{Service: "svc", Field: "image", Expected: "v2", Actual: "v1"},
		{Service: "svc", Field: "env", Expected: "prod", Actual: "staging"},
	}
}

func TestEvaluate_NoRules(t *testing.T) {
	p := &policy.Policy{}
	violations := p.Evaluate(diffs())
	if len(violations) != 0 {
		t.Fatalf("expected 0 violations, got %d", len(violations))
	}
}

func TestEvaluate_SpecificField(t *testing.T) {
	p := &policy.Policy{
		Rules: []policy.Rule{
			{Field: "replicas", Level: policy.LevelWarn},
		},
	}
	violations := p.Evaluate(diffs())
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(violations))
	}
	if violations[0].Level != policy.LevelWarn {
		t.Errorf("expected warn, got %s", violations[0].Level)
	}
}

func TestEvaluate_Wildcard(t *testing.T) {
	p := &policy.Policy{
		Rules: []policy.Rule{
			{Field: "*", Level: policy.LevelError},
		},
	}
	violations := p.Evaluate(diffs())
	if len(violations) != 3 {
		t.Fatalf("expected 3 violations, got %d", len(violations))
	}
}

func TestEvaluate_CaseInsensitive(t *testing.T) {
	p := &policy.Policy{
		Rules: []policy.Rule{
			{Field: "IMAGE", Level: policy.LevelError},
		},
	}
	violations := p.Evaluate(diffs())
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(violations))
	}
}

func TestHasErrors_True(t *testing.T) {
	violations := []policy.Violation{
		{Level: policy.LevelWarn},
		{Level: policy.LevelError},
	}
	if !policy.HasErrors(violations) {
		t.Error("expected HasErrors to return true")
	}
}

func TestHasErrors_False(t *testing.T) {
	violations := []policy.Violation{
		{Level: policy.LevelWarn},
	}
	if policy.HasErrors(violations) {
		t.Error("expected HasErrors to return false")
	}
}

func TestViolation_String(t *testing.T) {
	v := policy.Violation{
		Diff:  drift.Difference{Service: "svc", Field: "replicas", Expected: "3", Actual: "1"},
		Level: policy.LevelError,
	}
	s := v.String()
	if !strings.Contains(s, "ERROR") {
		t.Errorf("expected ERROR in string, got: %s", s)
	}
}
