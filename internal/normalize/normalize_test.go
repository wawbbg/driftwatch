package normalize

import (
	"testing"
)

func TestApply_LowercaseKeys(t *testing.T) {
	in := map[string]string{"HOST": "localhost", "Port": "8080"}
	out := Apply(in, Options{LowercaseKeys: true})
	if _, ok := out["host"]; !ok {
		t.Error("expected key 'host'")
	}
	if _, ok := out["port"]; !ok {
		t.Error("expected key 'port'")
	}
	if len(out) != 2 {
		t.Errorf("expected 2 keys, got %d", len(out))
	}
}

func TestApply_TrimValues(t *testing.T) {
	in := map[string]string{"key": "  value  "}
	out := Apply(in, Options{TrimValues: true})
	if got := out["key"]; got != "value" {
		t.Errorf("expected 'value', got %q", got)
	}
}

func TestApply_CollapseWhitespace(t *testing.T) {
	in := map[string]string{"desc": "hello   world\t!"}
	out := Apply(in, Options{CollapseWhitespace: true})
	if got := out["desc"]; got != "hello world !" {
		t.Errorf("unexpected value: %q", got)
	}
}

func TestApply_DoesNotMutateOriginal(t *testing.T) {
	in := map[string]string{"KEY": "  val  "}
	_ = Apply(in, DefaultOptions())
	if in["KEY"] != "  val  " {
		t.Error("original map was mutated")
	}
	if _, ok := in["key"]; ok {
		t.Error("original map gained a lowercase key")
	}
}

func TestApply_NilMap(t *testing.T) {
	out := Apply(nil, DefaultOptions())
	if out != nil {
		t.Errorf("expected nil, got %v", out)
	}
}

func TestApply_EmptyMap(t *testing.T) {
	out := Apply(map[string]string{}, DefaultOptions())
	if len(out) != 0 {
		t.Errorf("expected empty map, got %v", out)
	}
}

func TestDefaultOptions(t *testing.T) {
	opts := DefaultOptions()
	if !opts.LowercaseKeys {
		t.Error("expected LowercaseKeys to be true")
	}
	if !opts.TrimValues {
		t.Error("expected TrimValues to be true")
	}
	if opts.CollapseWhitespace {
		t.Error("expected CollapseWhitespace to be false")
	}
}

func TestApply_AllOptions(t *testing.T) {
	in := map[string]string{"MY_KEY": "  foo   bar  "}
	opts := Options{LowercaseKeys: true, TrimValues: true, CollapseWhitespace: true}
	out := Apply(in, opts)
	if got := out["my_key"]; got != "foo bar" {
		t.Errorf("expected 'foo bar', got %q", got)
	}
}
