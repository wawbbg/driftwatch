package report_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/driftwatch/internal/drift"
	"github.com/user/driftwatch/internal/report"
)

func diffs() []drift.Difference {
	return []drift.Difference{
		{Service: "api", Field: "replicas", Expected: "3", Actual: "1"},
		{Service: "api", Field: "image", Expected: "v2", Actual: "v1"},
	}
}

func TestWriteText_NoDrift(t *testing.T) {
	var buf bytes.Buffer
	w := report.NewWithWriter(&buf, report.FormatText)
	if err := w.Write(nil); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(buf.String(), "No drift detected") {
		t.Errorf("expected no-drift message, got: %s", buf.String())
	}
}

func TestWriteText_WithDrift(t *testing.T) {
	var buf bytes.Buffer
	w := report.NewWithWriter(&buf, report.FormatText)
	if err := w.Write(diffs()); err != nil {
		t.Fatal(err)
	}
	out := buf.String()
	if !strings.Contains(out, "2 difference") {
		t.Errorf("expected difference count, got: %s", out)
	}
	if !strings.Contains(out, "replicas") {
		t.Errorf("expected field name in output, got: %s", out)
	}
}

func TestWriteJSON_NoDrift(t *testing.T) {
	var buf bytes.Buffer
	w := report.NewWithWriter(&buf, report.FormatJSON)
	if err := w.Write(nil); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(buf.String(), `"drift":false`) {
		t.Errorf("expected drift:false, got: %s", buf.String())
	}
}

func TestWriteJSON_WithDrift(t *testing.T) {
	var buf bytes.Buffer
	w := report.NewWithWriter(&buf, report.FormatJSON)
	if err := w.Write(diffs()); err != nil {
		t.Fatal(err)
	}
	out := buf.String()
	if !strings.Contains(out, `"drift":true`) {
		t.Errorf("expected drift:true, got: %s", out)
	}
	if !strings.Contains(out, `"field":"replicas"`) {
		t.Errorf("expected field in JSON, got: %s", out)
	}
}
