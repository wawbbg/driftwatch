package label_test

import (
	"testing"

	"github.com/driftwatch/driftwatch/internal/label"
)

func TestNew_Empty(t *testing.T) {
	s := label.New()
	if len(s) != 0 {
		t.Fatalf("expected empty set, got %d entries", len(s))
	}
}

func TestAdd_And_Get(t *testing.T) {
	s := label.New()
	s.Add("Env", "prod")
	v, ok := s.Get("env")
	if !ok || v != "prod" {
		t.Fatalf("expected prod, got %q ok=%v", v, ok)
	}
}

func TestGet_Missing(t *testing.T) {
	s := label.New()
	_, ok := s.Get("missing")
	if ok {
		t.Fatal("expected missing key to return false")
	}
}

func TestDelete_RemovesKey(t *testing.T) {
	s := label.New()
	s.Add("team", "platform")
	s.Delete("team")
	_, ok := s.Get("team")
	if ok {
		t.Fatal("expected key to be deleted")
	}
}

func TestFromPairs_Valid(t *testing.T) {
	s := label.FromPairs([]string{"env=prod", "team=platform"})
	if v, _ := s.Get("env"); v != "prod" {
		t.Errorf("env: want prod, got %q", v)
	}
	if v, _ := s.Get("team"); v != "platform" {
		t.Errorf("team: want platform, got %q", v)
	}
}

func TestFromPairs_SkipsMalformed(t *testing.T) {
	s := label.FromPairs([]string{"nodivider", "ok=yes"})
	if len(s) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(s))
	}
}

func TestMatches_AllPresent(t *testing.T) {
	s := label.FromPairs([]string{"env=prod", "team=platform", "region=us-east"})
	filter := label.FromPairs([]string{"env=prod", "team=platform"})
	if !s.Matches(filter) {
		t.Fatal("expected match")
	}
}

func TestMatches_ValueMismatch(t *testing.T) {
	s := label.FromPairs([]string{"env=staging"})
	filter := label.FromPairs([]string{"env=prod"})
	if s.Matches(filter) {
		t.Fatal("expected no match on value mismatch")
	}
}

func TestMatches_EmptyFilter(t *testing.T) {
	s := label.FromPairs([]string{"env=prod"})
	if !s.Matches(label.New()) {
		t.Fatal("empty filter should always match")
	}
}

func TestPairs_Sorted(t *testing.T) {
	s := label.FromPairs([]string{"z=last", "a=first"})
	pairs := s.Pairs()
	if len(pairs) != 2 {
		t.Fatalf("expected 2 pairs, got %d", len(pairs))
	}
	if pairs[0] != "a=first" || pairs[1] != "z=last" {
		t.Errorf("unexpected order: %v", pairs)
	}
}
