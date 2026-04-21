package policy_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/example/driftwatch/internal/policy"
)

func writePolicy(t *testing.T, p policy.Policy) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "policy.json")
	f, err := os.Create(path)
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	defer f.Close()
	if err := json.NewEncoder(f).Encode(p); err != nil {
		t.Fatalf("encode: %v", err)
	}
	return path
}

func TestLoadFromFile_Valid(t *testing.T) {
	p := policy.Policy{
		Rules: []policy.Rule{
			{Field: "replicas", Level: policy.LevelWarn},
			{Field: "image", Level: policy.LevelError},
		},
	}
	path := writePolicy(t, p)
	loaded, err := policy.LoadFromFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(loaded.Rules) != 2 {
		t.Errorf("expected 2 rules, got %d", len(loaded.Rules))
	}
}

func TestLoadFromFile_EmptyPath(t *testing.T) {
	p, err := policy.LoadFromFile("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(p.Rules) != 0 {
		t.Errorf("expected empty policy")
	}
}

func TestLoadFromFile_NotFound(t *testing.T) {
	_, err := policy.LoadFromFile("/nonexistent/policy.json")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestLoadFromFile_UnsupportedFormat(t *testing.T) {
	_, err := policy.LoadFromFile("policy.yaml")
	if err == nil {
		t.Fatal("expected error for unsupported format")
	}
}

func TestLoadFromFile_InvalidLevel(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "policy.json")
	content := `{"rules":[{"field":"replicas","level":"critical"}]}`
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}
	_, err := policy.LoadFromFile(path)
	if err == nil {
		t.Fatal("expected validation error for unknown level")
	}
}

func TestLoadFromFile_EmptyField(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "policy.json")
	content := `{"rules":[{"field":"","level":"warn"}]}`
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}
	_, err := policy.LoadFromFile(path)
	if err == nil {
		t.Fatal("expected validation error for empty field")
	}
}
