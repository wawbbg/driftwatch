// Package dedupe deduplicates drift differences so that the same
// field-level change is not reported more than once within a single
// driftwatch run, even when the same service is checked multiple times
// (e.g. via scheduled polling or parallel resolution).
//
// Usage:
//
//	d := dedupe.New()
//	unique := d.Apply("my-service", diffs)
//
// Call Reset between independent runs if a fresh deduplication window
// is required.
package dedupe
