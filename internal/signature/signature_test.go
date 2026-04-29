package signature_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/user/driftwatch/internal/signature"
)

var fields = map[string]any{
	"replicas": "3",
	"image":    "nginx:1.25",
	"port":     "8080",
}

func TestSign_HexLength(t *testing.T) {
	e := signature.Sign("svc-a", fields)
	if len(e.Hex) != 64 {
		t.Fatalf("expected 64-char hex, got %d", len(e.Hex))
	}
	if e.Service != "svc-a" {
		t.Fatalf("unexpected service: %s", e.Service)
	}
}

func TestVerify_Match(t *testing.T) {
	e := signature.Sign("svc-a", fields)
	if !signature.Verify(e, fields) {
		t.Fatal("expected Verify to return true for unchanged fields")
	}
}

func TestVerify_Mismatch(t *testing.T) {
	e := signature.Sign("svc-a", fields)
	modified := map[string]any{
		"replicas": "5",
		"image":    "nginx:1.25",
		"port":     "8080",
	}
	if signature.Verify(e, modified) {
		t.Fatal("expected Verify to return false for changed fields")
	}
}

func TestStoreAndLoad(t *testing.T) {
	dir := t.TempDir()
	e := signature.Sign("svc-b", fields)

	if err := signature.Store(dir, e); err != nil {
		t.Fatalf("Store: %v", err)
	}

	got, err := signature.Load(dir, "svc-b")
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if got.Hex != e.Hex {
		t.Fatalf("hex mismatch: want %s got %s", e.Hex, got.Hex)
	}
}

func TestLoad_NotFound(t *testing.T) {
	dir := t.TempDir()
	_, err := signature.Load(dir, "missing")
	if err == nil {
		t.Fatal("expected error for missing signature file")
	}
}

func TestStore_BadDir(t *testing.T) {
	// Use a file as the directory to force an error.
	f, _ := os.CreateTemp("", "sig")
	f.Close()
	defer os.Remove(f.Name())

	e := signature.Sign("svc-c", fields)
	err := signature.Store(filepath.Join(f.Name(), "sub"), e)
	if err == nil {
		t.Fatal("expected error when dir cannot be created")
	}
}
