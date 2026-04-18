package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/driftwatch/internal/config"
)

func writeTemp(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "driftwatch-*.yaml")
	if err != nil {
		t.Fatalf("creating temp file: %v", err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("writing temp file: %v", err)
	}
	f.Close()
	return f.Name()
}

func TestLoadFromFile_Valid(t *testing.T) {
	yaml := `
services:
  - name: api
    source_path: ./manifests/api.yaml
    deployed_at: k8s://default/api
    labels:
      env: prod
output_format: text
`
	path := writeTemp(t, yaml)
	cfg, err := config.LoadFromFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cfg.Services) != 1 {
		t.Fatalf("expected 1 service, got %d", len(cfg.Services))
	}
	if cfg.Services[0].Name != "api" {
		t.Errorf("expected name 'api', got %q", cfg.Services[0].Name)
	}
}

func TestLoadFromFile_MissingName(t *testing.T) {
	yaml := `
services:
  - source_path: ./manifests/api.yaml
    deployed_at: k8s://default/api
`
	path := writeTemp(t, yaml)
	_, err := config.LoadFromFile(path)
	if err == nil {
		t.Fatal("expected validation error, got nil")
	}
}

func TestLoadFromFile_NotFound(t *testing.T) {
	_, err := config.LoadFromFile(filepath.Join(t.TempDir(), "nonexistent.yaml"))
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestLoadFromFile_EmptyPath(t *testing.T) {
	_, err := config.LoadFromFile("")
	if err == nil {
		t.Fatal("expected error for empty path")
	}
}
