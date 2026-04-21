// Package diff provides low-level map comparison utilities used by the
// drift detector to identify field-level discrepancies between a service's
// expected configuration and its live state.
//
// Usage:
//
//	results := diff.Maps("my-service", expectedMap, actualMap)
//	if diff.HasDrift(results) {
//		for _, r := range results {
//			fmt.Println(r)
//		}
//	}
//
// Maps only inspects keys present in the expected map; additional keys in
// the actual map are silently ignored, keeping comparisons focused on
// declared intent rather than runtime metadata.
package diff
