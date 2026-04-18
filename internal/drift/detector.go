package drift

import (
	"fmt"
	"reflect"
)

// Field represents a single configuration field with its value.
type Field struct {
	Name  string
	Value interface{}
}

// Difference describes a detected drift between expected and actual values.
type Difference struct {
	Field    string
	Expected interface{}
	Actual   interface{}
}

func (d Difference) String() string {
	return fmt.Sprintf("field %q: expected %v, got %v", d.Field, d.Expected, d.Actual)
}

// Result holds the outcome of a drift detection run.
type Result struct {
	ServiceName string
	Diffs       []Difference
}

// HasDrift returns true if any differences were found.
func (r Result) HasDrift() bool {
	return len(r.Diffs) > 0
}

// Detect compares expected fields against actual fields and returns a Result.
// Fields present in expected but missing from actual are reported as drift.
func Detect(serviceName string, expected, actual map[string]interface{}) Result {
	result := Result{ServiceName: serviceName}

	for key, expVal := range expected {
		actVal, ok := actual[key]
		if !ok {
			result.Diffs = append(result.Diffs, Difference{
				Field:    key,
				Expected: expVal,
				Actual:   nil,
			})
			continue
		}
		if !reflect.DeepEqual(expVal, actVal) {
			result.Diffs = append(result.Diffs, Difference{
				Field:    key,
				Expected: expVal,
				Actual:   actVal,
			})
		}
	}

	return result
}
