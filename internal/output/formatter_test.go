package output

import (
	"bytes"
	"strings"
	"testing"
)

func TestParseFormat_Valid(t *testing.T) {
	cases := []struct {
		input    string
		want     Format
	}{
		{"text", FormatText},
		{"TEXT", FormatText},
		{"", FormatText},
		{"json", FormatJSON},
		{"JSON", FormatJSON},
	}
	for _, tc := range cases {
		got, err := ParseFormat(tc.input)
		if err != nil {
			t.Errorf("ParseFormat(%q) unexpected error: %v", tc.input, err)
		}
		if got != tc.want {
			t.Errorf("ParseFormat(%q) = %q; want %q", tc.input, got, tc.want)
		}
	}
}

func TestParseFormat_Invalid(t *testing.T) {
	_, err := ParseFormat("xml")
	if err == nil {
		t.Fatal("expected error for unknown format, got nil")
	}
	if !strings.Contains(err.Error(), "xml") {
		t.Errorf("error message should mention the bad value, got: %v", err)
	}
}

func TestExitCode(t *testing.T) {
	if ExitCode(true) != 1 {
		t.Error("expected exit code 1 when drift detected")
	}
	if ExitCode(false) != 0 {
		t.Error("expected exit code 0 when no drift")
	}
}

func TestPrinter(t *testing.T) {
	var buf bytes.Buffer
	print := Printer(&buf)
	print("hello %s", "world")
	got := buf.String()
	if got != "hello world\n" {
		t.Errorf("Printer output = %q; want %q", got, "hello world\n")
	}
}

func TestPrinter_NilFallback(t *testing.T) {
	// Should not panic when writer is nil (falls back to stdout).
	print := Printer(nil)
	print("no panic")
}
