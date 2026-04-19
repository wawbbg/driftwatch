package filter_test

import (
	"testing"

	"github.com/yourorg/driftwatch/internal/filter"
)

func svc(name string, labels map[string]string) filter.Service {
	return filter.Service{Name: name, Labels: labels}
}

func TestMatch_NoFilter(t *testing.T) {
	s := svc("api", map[string]string{"env": "prod"})
	if !filter.Match(s, filter.Options{}) {
		t.Error("empty options should match any service")
	}
}

func TestMatch_ByName(t *testing.T) {
	s := svc("api", nil)
	if !filter.Match(s, filter.Options{Names: []string{"api"}}) {
		t.Error("expected match by name")
	}
	if filter.Match(s, filter.Options{Names: []string{"worker"}}) {
		t.Error("expected no match for different name")
	}
}

func TestMatch_ByLabel(t *testing.T) {
	s := svc("api", map[string]string{"env": "prod", "team": "platform"})
	opts := filter.Options{Labels: map[string]string{"env": "prod"}}
	if !filter.Match(s, opts) {
		t.Error("expected label match")
	}
	opts2 := filter.Options{Labels: map[string]string{"env": "staging"}}
	if filter.Match(s, opts2) {
		t.Error("expected no match for wrong label value")
	}
}

func TestMatch_NameAndLabel(t *testing.T) {
	s := svc("api", map[string]string{"env": "prod"})
	opts := filter.Options{Names: []string{"api"}, Labels: map[string]string{"env": "prod"}}
	if !filter.Match(s, opts) {
		t.Error("expected match on name+label")
	}
	opts2 := filter.Options{Names: []string{"worker"}, Labels: map[string]string{"env": "prod"}}
	if filter.Match(s, opts2) {
		t.Error("expected no match when name differs")
	}
}

func TestApply(t *testing.T) {
	services := []filter.Service{
		svc("api", map[string]string{"env": "prod"}),
		svc("worker", map[string]string{"env": "staging"}),
		svc("cron", map[string]string{"env": "prod"}),
	}
	got := filter.Apply(services, filter.Options{Labels: map[string]string{"env": "prod"}})
	if len(got) != 2 {
		t.Fatalf("expected 2 services, got %d", len(got))
	}
}

func TestApply_EmptyOpts(t *testing.T) {
	services := []filter.Service{svc("a", nil), svc("b", nil)}
	got := filter.Apply(services, filter.Options{})
	if len(got) != 2 {
		t.Error("empty opts should return all services")
	}
}
