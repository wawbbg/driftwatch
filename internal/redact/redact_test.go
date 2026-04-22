package redact_test

import (
	"testing"

	"github.com/yourorg/driftwatch/internal/redact"
)

func TestIsSensitive_DefaultKeys(t *testing.T) {
	r := redact.New(nil)
	cases := []struct {
		field     string
		expected  bool
	}{
		{"db_password", true},
		{"API_KEY", true},
		{"auth_token", true},
		{"SECRET", true},
		{"replica_count", false},
		{"port", false},
		{"image", false},
	}
	for _, tc := range cases {
		t.Run(tc.field, func(t *testing.T) {
			if got := r.IsSensitive(tc.field); got != tc.expected {
				t.Errorf("IsSensitive(%q) = %v, want %v", tc.field, got, tc.expected)
			}
		})
	}
}

func TestIsSensitive_CustomKeys(t *testing.T) {
	r := redact.New([]string{"pin", "ssn"})
	if !r.IsSensitive("user_pin") {
		t.Error("expected user_pin to be sensitive")
	}
	if r.IsSensitive("password") {
		t.Error("expected password NOT to be sensitive with custom keys")
	}
}

func TestApply_RedactsSensitiveFields(t *testing.T) {
	r := redact.New(nil)
	input := map[string]any{
		"image":       "nginx:latest",
		"db_password": "s3cr3t",
		"port":        "8080",
		"api_key":     "abc123",
	}
	out := r.Apply(input)

	if out["image"] != "nginx:latest" {
		t.Errorf("image should be unchanged, got %v", out["image"])
	}
	if out["port"] != "8080" {
		t.Errorf("port should be unchanged, got %v", out["port"])
	}
	if out["db_password"] != "[REDACTED]" {
		t.Errorf("db_password should be redacted, got %v", out["db_password"])
	}
	if out["api_key"] != "[REDACTED]" {
		t.Errorf("api_key should be redacted, got %v", out["api_key"])
	}
}

func TestApply_DoesNotMutateOriginal(t *testing.T) {
	r := redact.New(nil)
	input := map[string]any{
		"token": "original-value",
	}
	_ = r.Apply(input)
	if input["token"] != "original-value" {
		t.Error("Apply must not mutate the original map")
	}
}

func TestApply_EmptyMap(t *testing.T) {
	r := redact.New(nil)
	out := r.Apply(map[string]any{})
	if len(out) != 0 {
		t.Errorf("expected empty map, got %v", out)
	}
}
