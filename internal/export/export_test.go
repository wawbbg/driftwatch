package export_test

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"strings"
	"testing"

	"github.com/driftwatch/internal/diff"
	"github.com/driftwatch/internal/export"
)

var testDiffs = []diff.Difference{
	{Field: "replicas", Expected: "3", Actual: "1"},
	{Field: "image", Expected: "nginx:1.25", Actual: "nginx:1.20"},
}

func TestWrite_CSV_NoDrift(t *testing.T) {
	var buf bytes.Buffer
	e := export.New(&buf, export.FormatCSV)
	if err := e.Write("svc", nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	r := csv.NewReader(&buf)
	rows, _ := r.ReadAll()
	if len(rows) != 1 {
		t.Fatalf("expected header only, got %d rows", len(rows))
	}
}

func TestWrite_CSV_WithDrift(t *testing.T) {
	var buf bytes.Buffer
	e := export.New(&buf, export.FormatCSV)
	if err := e.Write("api", testDiffs); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	r := csv.NewReader(&buf)
	rows, _ := r.ReadAll()
	// header + 2 data rows
	if len(rows) != 3 {
		t.Fatalf("expected 3 rows, got %d", len(rows))
	}
	if rows[1][0] != "api" {
		t.Errorf("expected service 'api', got %q", rows[1][0])
	}
	if rows[1][1] != "replicas" {
		t.Errorf("expected field 'replicas', got %q", rows[1][1])
	}
}

func TestWrite_NDJSON_WithDrift(t *testing.T) {
	var buf bytes.Buffer
	e := export.New(&buf, export.FormatNDJSON)
	if err := e.Write("worker", testDiffs); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 2 {
		t.Fatalf("expected 2 lines, got %d", len(lines))
	}
	var rec map[string]interface{}
	if err := json.Unmarshal([]byte(lines[0]), &rec); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if rec["service"] != "worker" {
		t.Errorf("expected service 'worker', got %v", rec["service"])
	}
}

func TestWrite_UnknownFormat(t *testing.T) {
	var buf bytes.Buffer
	e := export.New(&buf, export.Format("xml"))
	if err := e.Write("svc", testDiffs); err == nil {
		t.Fatal("expected error for unknown format")
	}
}
