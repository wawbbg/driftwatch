// Package diff provides utilities for comparing two arbitrary
// key-value maps and returning a list of named differences.
package diff

import "fmt"

// Result holds a single field-level difference between an expected
// and an actual value for a named service.
type Result struct {
	Service  string
	Field    string
	Expected interface{}
	Actual   interface{}
}

// String returns a human-readable representation of the difference.
func (r Result) String() string {
	return fmt.Sprintf("[%s] %s: expected %v, got %v",
		r.Service, r.Field, r.Expected, r.Actual)
}

// Maps compares two flat string-keyed maps and returns all differences.
// Fields present in expected but missing in actual are reported with
// Actual set to nil. Fields present only in actual are ignored.
func Maps(service string, expected, actual map[string]interface{}) []Result {
	var results []Result

	for key, expVal := range expected {
		actVal, ok := actual[key]
		if !ok {
			results = append(results, Result{
				Service:  service,
				Field:    key,
				Expected: expVal,
				Actual:   nil,
			})
			continue
		}
		if fmt.Sprintf("%v", expVal) != fmt.Sprintf("%v", actVal) {
			results = append(results, Result{
				Service:  service,
				Field:    key,
				Expected: expVal,
				Actual:   actVal,
			})
		}
	}

	return results
}

// HasDrift returns true when the provided slice contains at least one Result.
func HasDrift(results []Result) bool {
	return len(results) > 0
}
