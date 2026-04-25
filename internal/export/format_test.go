package export_test

import (
	"testing"

	"github.com/driftwatch/internal/export"
)

func TestParseFormat_Valid(t *testing.T) {
	cases := []struct {
		input string
		want  export.Format
	}{
		{"csv", export.FormatCSV},
		{"CSV", export.FormatCSV},
		{"ndjson", export.FormatNDJSON},
		{"NDJSON", export.FormatNDJSON},
		{"jsonl", export.FormatNDJSON},
	}
	for _, tc := range cases {
		t.Run(tc.input, func(t *testing.T) {
			got, err := export.ParseFormat(tc.input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tc.want {
				t.Errorf("got %q, want %q", got, tc.want)
			}
		})
	}
}

func TestParseFormat_Invalid(t *testing.T) {
	_, err := export.ParseFormat("yaml")
	if err == nil {
		t.Fatal("expected error for unsupported format")
	}
}

func TestSupported(t *testing.T) {
	list := export.Supported()
	if len(list) == 0 {
		t.Fatal("expected at least one supported format")
	}
}
