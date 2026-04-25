package tag_test

import (
	"testing"

	"github.com/example/driftwatch/internal/tag"
)

func TestAdd_And_Get(t *testing.T) {
	s := tag.New()
	s.Add("env", "production")

	v, ok := s.Get("env")
	if !ok {
		t.Fatal("expected tag 'env' to exist")
	}
	if v != "production" {
		t.Fatalf("expected 'production', got %q", v)
	}
}

func TestAdd_NormalisesKey(t *testing.T) {
	s := tag.New()
	s.Add("  ENV  ", "staging")

	_, ok := s.Get("env")
	if !ok {
		t.Fatal("expected normalised key 'env' to exist")
	}
}

func TestGet_Missing(t *testing.T) {
	s := tag.New()
	_, ok := s.Get("missing")
	if ok {
		t.Fatal("expected missing key to return false")
	}
}

func TestParse_ValidPairs(t *testing.T) {
	s := tag.New()
	s.Parse([]string{"team=platform", "env=prod", "tier=backend"})

	if s.Len() != 3 {
		t.Fatalf("expected 3 tags, got %d", s.Len())
	}

	v, _ := s.Get("team")
	if v != "platform" {
		t.Fatalf("expected 'platform', got %q", v)
	}
}

func TestParse_SkipsMalformed(t *testing.T) {
	s := tag.New()
	s.Parse([]string{"valid=yes", "noequals", "another=ok"})

	if s.Len() != 2 {
		t.Fatalf("expected 2 tags, got %d", s.Len())
	}
}

func TestAll_SortedByKey(t *testing.T) {
	s := tag.New()
	s.Parse([]string{"z=last", "a=first", "m=middle"})

	all := s.All()
	if len(all) != 3 {
		t.Fatalf("expected 3 tags, got %d", len(all))
	}
	if all[0].Key != "a" || all[1].Key != "m" || all[2].Key != "z" {
		t.Fatalf("tags not sorted: %v", all)
	}
}

func TestTag_String(t *testing.T) {
	tg := tag.Tag{Key: "env", Value: "prod"}
	if tg.String() != "env=prod" {
		t.Fatalf("unexpected String(): %q", tg.String())
	}
}

func TestAdd_Overwrites(t *testing.T) {
	s := tag.New()
	s.Add("env", "dev")
	s.Add("env", "prod")

	v, _ := s.Get("env")
	if v != "prod" {
		t.Fatalf("expected overwritten value 'prod', got %q", v)
	}
	if s.Len() != 1 {
		t.Fatalf("expected 1 tag after overwrite, got %d", s.Len())
	}
}
