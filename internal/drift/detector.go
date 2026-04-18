package drift

import (
	"fmt"
)

// Difference represents a single detected drift between expected and actual config.
type Difference struct {
	Service  string
	Field    string
	Expected string
	Actual   string
}

// String returns a human-readable description of the difference.
func (d Difference) String() string {
	return fmt.Sprintf("[%s] %s: expected %q, got %q", d.Service, d.Field, d.Expected, d.Actual)
}

// Config holds a named set of key-value configuration fields.
type Config struct {
	Name   string
	Fields map[string]string
}

// Detect compares expected config against actual config and returns all differences.
func Detect(expected, actual Config) []Difference {
	var diffs []Difference
	service := expected.Name
	if service == "" {
		service = actual.Name
	}
	for key, expVal := range expected.Fields {
		actVal, ok := actual.Fields[key]
		if !ok {
			diffs = append(diffs, Difference{
				Service:  service,
				Field:    key,
				Expected: expVal,
				Actual:   "<missing>",
			})
			continue
		}
		if expVal != actVal {
			diffs = append(diffs, Difference{
				Service:  service,
				Field:    key,
				Expected: expVal,
				Actual:   actVal,
			})
		}
	}
	return diffs
}
