package dedupe

import (
	"testing"

	"github.com/example/driftwatch/internal/diff"
)

func diffs(pairs ...string) []diff.Difference {
	var out []diff.Difference
	for i := 0; i+1 < len(pairs); i += 2 {
		out = append(out, diff.Difference{Field: pairs[i], Want: pairs[i+1], Got: "actual"})
	}
	return out
}

func TestApply_FirstCallPassesThrough(t *testing.T) {
	d := New()
	in := diffs("port", "8080", "host", "localhost")
	out := d.Apply("svc", in)
	if len(out) != 2 {
		t.Fatalf("expected 2 diffs, got %d", len(out))
	}
}

func TestApply_DuplicateDropped(t *testing.T) {
	d := New()
	in := diffs("port", "8080")
	d.Apply("svc", in)
	out := d.Apply("svc", in)
	if len(out) != 0 {
		t.Fatalf("expected 0 diffs on second call, got %d", len(out))
	}
}

func TestApply_DifferentServiceNotDeduplicated(t *testing.T) {
	d := New()
	in := diffs("port", "8080")
	d.Apply("svc-a", in)
	out := d.Apply("svc-b", in)
	if len(out) != 1 {
		t.Fatalf("expected 1 diff for different service, got %d", len(out))
	}
}

func TestApply_PartialDuplicate(t *testing.T) {
	d := New()
	first := diffs("port", "8080")
	d.Apply("svc", first)

	second := []diff.Difference{
		{Field: "port", Want: "8080", Got: "actual"},  // duplicate
		{Field: "host", Want: "localhost", Got: "actual"}, // new
	}
	out := d.Apply("svc", second)
	if len(out) != 1 {
		t.Fatalf("expected 1 new diff, got %d", len(out))
	}
	if out[0].Field != "host" {
		t.Errorf("expected field 'host', got %q", out[0].Field)
	}
}

func TestReset_ClearsState(t *testing.T) {
	d := New()
	in := diffs("port", "8080")
	d.Apply("svc", in)
	d.Reset()
	out := d.Apply("svc", in)
	if len(out) != 1 {
		t.Fatalf("expected 1 diff after reset, got %d", len(out))
	}
}

func TestLen_TracksUniqueCount(t *testing.T) {
	d := New()
	if d.Len() != 0 {
		t.Fatalf("expected 0 initially, got %d", d.Len())
	}
	d.Apply("svc", diffs("port", "8080", "host", "localhost"))
	if d.Len() != 2 {
		t.Fatalf("expected 2, got %d", d.Len())
	}
	// duplicate should not increment
	d.Apply("svc", diffs("port", "8080"))
	if d.Len() != 2 {
		t.Fatalf("expected still 2, got %d", d.Len())
	}
}
