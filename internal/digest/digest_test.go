package digest_test

import (
	"testing"

	"github.com/user/driftwatch/internal/digest"
)

func TestSum_NilMap(t *testing.T) {
	d, err := digest.Sum(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if d == "" {
		t.Fatal("expected non-empty digest for nil map")
	}
}

func TestSum_EmptyMap(t *testing.T) {
	nil1, _ := digest.Sum(nil)
	empty, _ := digest.Sum(map[string]any{})
	if nil1 != empty {
		t.Errorf("nil and empty map should produce same digest: %q vs %q", nil1, empty)
	}
}

func TestSum_Deterministic(t *testing.T) {
	m := map[string]any{"replicas": 3, "image": "nginx:latest", "port": 8080}
	d1, err := digest.Sum(m)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	d2, err := digest.Sum(m)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if d1 != d2 {
		t.Errorf("non-deterministic digest: %q vs %q", d1, d2)
	}
}

func TestSum_OrderIndependent(t *testing.T) {
	a := map[string]any{"x": 1, "y": 2}
	b := map[string]any{"y": 2, "x": 1}
	da, _ := digest.Sum(a)
	db, _ := digest.Sum(b)
	if da != db {
		t.Errorf("digest should be key-order independent: %q vs %q", da, db)
	}
}

func TestSum_DifferentValues(t *testing.T) {
	a := map[string]any{"replicas": 1}
	b := map[string]any{"replicas": 2}
	da, _ := digest.Sum(a)
	db, _ := digest.Sum(b)
	if da == db {
		t.Error("expected different digests for different values")
	}
}

func TestEqual_SameMaps(t *testing.T) {
	m := map[string]any{"env": "prod", "replicas": 5}
	ok, err := digest.Equal(m, m)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok {
		t.Error("expected Equal to return true for identical maps")
	}
}

func TestEqual_DifferentMaps(t *testing.T) {
	a := map[string]any{"env": "prod"}
	b := map[string]any{"env": "staging"}
	ok, err := digest.Equal(a, b)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ok {
		t.Error("expected Equal to return false for different maps")
	}
}

func TestEqual_NilVsEmpty(t *testing.T) {
	ok, err := digest.Equal(nil, map[string]any{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok {
		t.Error("nil and empty map should be considered equal")
	}
}
