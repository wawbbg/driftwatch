package export

import (
	"fmt"
	"strings"
)

// ParseFormat parses a format string into a Format constant.
// It is case-insensitive and returns an error for unknown values.
func ParseFormat(s string) (Format, error) {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "csv":
		return FormatCSV, nil
	case "ndjson", "jsonl":
		return FormatNDJSON, nil
	default:
		return "", fmt.Errorf("export: unknown format %q (supported: csv, ndjson)", s)
	}
}

// Supported returns the list of supported format names.
func Supported() []string {
	return []string{"csv", "ndjson"}
}
