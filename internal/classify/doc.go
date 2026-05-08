// Package classify assigns severity levels (low, medium, high, critical) to
// detected configuration drift differences.
//
// Severities are determined by matching field names against configurable lists:
//
//	- critical: fields that contain sensitive substrings such as "token" or "secret"
//	- high:     fields that affect service connectivity such as "host" or "port"
//	- medium:   differences where the actual value is missing entirely
//	- low:      all other differences
//
// Usage:
//
//	c := classify.New()
//	results := c.Apply(diffs)
//	if classify.HasCritical(results) {
//		// escalate
//	}
package classify
