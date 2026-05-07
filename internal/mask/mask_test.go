package mask_test

import (
	"testing"

	"github.com/yourorg/driftwatch/internal/mask"
)

func TestIsSensitive_DefaultKeys(t *testing.T) {
	m := mask.New()
	sensitive := []string{"password", "Password", "TOKEN", "secret", "apikey", "api_key", "private_key"}
	for _, k := range sensitive {
		if !m.IsSensitive(k) {
			t.Errorf("expected %q to be sensitive", k)
		}
	}
}

func TestIsSensitive_SafeKey(t *testing.T) {
	m := mask.New()
	if m.IsSensitive("host") {
		t.Error("expected 'host' to be safe")
	}
}

func TestApply_MasksSensitiveFields(t *testing.T) {
	m := mask.New()
	input := map[string]string{
		"host":     "localhost",
		"password": "s3cr3t",
		"token":    "abc123",
	}
	out := m.Apply(input)
	if out["host"] != "localhost" {
		t.Errorf("host should be unchanged, got %q", out["host"])
	}
	if out["password"] != mask.Placeholder {
		t.Errorf("password should be masked, got %q", out["password"])
	}
	if out["token"] != mask.Placeholder {
		t.Errorf("token should be masked, got %q", out["token"])
	}
}

func TestApply_DoesNotMutateOriginal(t *testing.T) {
	m := mask.New()
	input := map[string]string{"password": "real"}
	_ = m.Apply(input)
	if input["password"] != "real" {
		t.Error("Apply must not mutate the original map")
	}
}

func TestApply_EmptyMap(t *testing.T) {
	m := mask.New()
	out := m.Apply(map[string]string{})
	if len(out) != 0 {
		t.Errorf("expected empty map, got %v", out)
	}
}

func TestNewWithKeys_CustomKeys(t *testing.T) {
	m := mask.NewWithKeys([]string{"db_pass", "cert"})
	if !m.IsSensitive("db_pass") {
		t.Error("expected db_pass to be sensitive")
	}
	if m.IsSensitive("password") {
		t.Error("default key 'password' should not be sensitive with custom keys")
	}
}

func TestKeys_ReturnsAll(t *testing.T) {
	m := mask.NewWithKeys([]string{"alpha", "beta", "gamma"})
	keys := m.Keys()
	if len(keys) != 3 {
		t.Errorf("expected 3 keys, got %d", len(keys))
	}
}
