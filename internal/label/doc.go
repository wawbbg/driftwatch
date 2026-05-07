// Package label provides a lightweight key-value label Set for driftwatch
// services. Labels are normalised (lowercase, trimmed) on insertion so that
// comparisons are case-insensitive and whitespace-tolerant.
//
// Typical usage:
//
//	ls := label.FromPairs([]string{"env=prod", "team=platform"})
//	if ls.Matches(filter) {
//		// service matches filter criteria
//	}
package label
