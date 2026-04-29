package signature_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/driftwatch/internal/signature"
)

func TestChecker_Valid(t *testing.T) {
	dir := t.TempDir()
	e := signature.Sign("svc-x", fields)
	if err := signature.Store(dir, e); err != nil {
		t.Fatalf("Store: %v", err)
	}

	var buf bytes.Buffer
	c := signature.NewCheckerWithWriter(dir, &buf)
	r := c.Check("svc-x", fields)

	if !r.Valid {
		t.Fatalf("expected valid result, got: %s", r.Message)
	}
	if !strings.Contains(buf.String(), "[OK]") {
		t.Fatalf("expected [OK] in output, got: %s", buf.String())
	}
}

func TestChecker_Mismatch(t *testing.T) {
	dir := t.TempDir()
	e := signature.Sign("svc-y", fields)
	if err := signature.Store(dir, e); err != nil {
		t.Fatalf("Store: %v", err)
	}

	changed := map[string]any{"replicas": "99", "image": "nginx:1.25", "port": "8080"}
	var buf bytes.Buffer
	c := signature.NewCheckerWithWriter(dir, &buf)
	r := c.Check("svc-y", changed)

	if r.Valid {
		t.Fatal("expected invalid result for changed fields")
	}
	if !strings.Contains(buf.String(), "[FAIL]") {
		t.Fatalf("expected [FAIL] in output, got: %s", buf.String())
	}
}

func TestChecker_NoStoredSignature(t *testing.T) {
	dir := t.TempDir()
	var buf bytes.Buffer
	c := signature.NewCheckerWithWriter(dir, &buf)
	r := c.Check("ghost", fields)

	if r.Valid {
		t.Fatal("expected invalid result when no signature stored")
	}
	if !strings.Contains(r.Message, "no stored signature") {
		t.Fatalf("unexpected message: %s", r.Message)
	}
}

func TestResult_String_OK(t *testing.T) {
	r := signature.Result{Service: "svc-z", Valid: true}
	if !strings.HasPrefix(r.String(), "[OK]") {
		t.Fatalf("expected [OK] prefix, got: %s", r.String())
	}
}

func TestResult_String_Fail(t *testing.T) {
	r := signature.Result{Service: "svc-z", Valid: false, Message: "signature mismatch"}
	if !strings.HasPrefix(r.String(), "[FAIL]") {
		t.Fatalf("expected [FAIL] prefix, got: %s", r.String())
	}
}
