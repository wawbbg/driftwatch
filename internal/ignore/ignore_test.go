package ignore_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yourorg/driftwatch/internal/ignore"
)

func writeIgnoreFile(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, "ignore.yaml")
	if err := os.WriteFile(p, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	return p
}

func TestLoadFromFile_Valid(t *testing.T) {
	p := writeIgnoreFile(t, `ignore:
  - service: api
    fields: [version, replicas]
  - service: "*"
    fields: [timestamp]
`)
	rules, err := ignore.LoadFromFile(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(rules.Ignore) != 2 {
		t.Fatalf("expected 2 rules, got %d", len(rules.Ignore))
	}
}

func TestLoadFromFile_NotFound(t *testing.T) {
	_, err := ignore.LoadFromFile("/nonexistent/path.yaml")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestMatchesService(t *testing.T) {
	p := writeIgnoreFile(t, `ignore:
  - service: api
    fields: [version]
`)
	rules, _ := ignore.LoadFromFile(p)
	if !rules.MatchesService("api") {
		t.Error("expected api to match")
	}
	if rules.MatchesService("worker") {
		t.Error("expected worker not to match")
	}
}

func TestIgnoredFields_Wildcard(t *testing.T) {
	p := writeIgnoreFile(t, `ignore:
  - service: "*"
    fields: [timestamp]
  - service: api
    fields: [version]
`)
	rules, _ := ignore.LoadFromFile(p)
	fields := rules.IgnoredFields("api")
	if !fields["timestamp"] {
		t.Error("expected timestamp to be ignored via wildcard")
	}
	if !fields["version"] {
		t.Error("expected version to be ignored for api")
	}
}

func TestIgnoredFields_NoMatch(t *testing.T) {
	p := writeIgnoreFile(t, `ignore:
  - service: api
    fields: [version]
`)
	rules, _ := ignore.LoadFromFile(p)
	fields := rules.IgnoredFields("worker")
	if len(fields) != 0 {
		t.Errorf("expected no ignored fields for worker, got %v", fields)
	}
}
