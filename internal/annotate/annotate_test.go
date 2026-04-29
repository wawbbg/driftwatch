package annotate_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/driftwatch/driftwatch/internal/annotate"
)

func TestSetAndGet(t *testing.T) {
	dir := t.TempDir()

	if err := annotate.Set(dir, "svc-a", "ticket", "INFRA-1"); err != nil {
		t.Fatalf("Set: %v", err)
	}
	if err := annotate.Set(dir, "svc-a", "owner", "alice"); err != nil {
		t.Fatalf("Set: %v", err)
	}

	anns, err := annotate.Get(dir, "svc-a")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if len(anns) != 2 {
		t.Fatalf("expected 2 annotations, got %d", len(anns))
	}
	if anns[0].Key != "ticket" || anns[0].Value != "INFRA-1" {
		t.Errorf("unexpected first annotation: %+v", anns[0])
	}
	if anns[1].Key != "owner" || anns[1].Value != "alice" {
		t.Errorf("unexpected second annotation: %+v", anns[1])
	}
}

func TestGet_NotFound(t *testing.T) {
	dir := t.TempDir()
	anns, err := annotate.Get(dir, "ghost")
	if err != nil {
		t.Fatalf("expected no error for missing service, got %v", err)
	}
	if len(anns) != 0 {
		t.Errorf("expected empty slice, got %v", anns)
	}
}

func TestDelete_RemovesMatchingKey(t *testing.T) {
	dir := t.TempDir()
	_ = annotate.Set(dir, "svc-b", "ticket", "INFRA-2")
	_ = annotate.Set(dir, "svc-b", "owner", "bob")

	if err := annotate.Delete(dir, "svc-b", "ticket"); err != nil {
		t.Fatalf("Delete: %v", err)
	}

	anns, _ := annotate.Get(dir, "svc-b")
	if len(anns) != 1 {
		t.Fatalf("expected 1 annotation after delete, got %d", len(anns))
	}
	if anns[0].Key != "owner" {
		t.Errorf("expected remaining annotation key=owner, got %q", anns[0].Key)
	}
}

func TestSet_NormalisesKey(t *testing.T) {
	dir := t.TempDir()
	_ = annotate.Set(dir, "svc-c", "  TICKET  ", "X-99")
	anns, _ := annotate.Get(dir, "svc-c")
	if len(anns) == 0 {
		t.Fatal("expected annotation")
	}
	if anns[0].Key != "ticket" {
		t.Errorf("expected normalised key 'ticket', got %q", anns[0].Key)
	}
}

func TestSet_EmptyService(t *testing.T) {
	dir := t.TempDir()
	if err := annotate.Set(dir, "", "k", "v"); err == nil {
		t.Error("expected error for empty service name")
	}
}

func TestSet_BadDir(t *testing.T) {
	// Use a file as the directory to trigger a mkdir error.
	f, _ := os.CreateTemp("", "annotate-*")
	_ = f.Close()
	defer os.Remove(f.Name())

	badDir := filepath.Join(f.Name(), "sub")
	if err := annotate.Set(badDir, "svc", "k", "v"); err == nil {
		t.Error("expected error when dir is unwritable")
	}
}

func TestAnnotation_RecordedAt(t *testing.T) {
	dir := t.TempDir()
	_ = annotate.Set(dir, "svc-d", "env", "prod")
	anns, _ := annotate.Get(dir, "svc-d")
	if anns[0].RecordedAt.IsZero() {
		t.Error("expected RecordedAt to be set")
	}
}
