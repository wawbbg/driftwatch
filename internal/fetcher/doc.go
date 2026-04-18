// Package fetcher provides functionality for retrieving live configuration
// state from running services over HTTP.
//
// A service is expected to expose a JSON endpoint that returns a flat
// key-value map of its current configuration. The Fetcher normalises all
// values to strings so they can be compared directly against the expected
// values defined in driftwatch.yaml by the drift detector.
//
// Basic usage:
//
//	f := fetcher.New()
//	state, err := f.Fetch("my-service", "http://my-service/config")
//	if err != nil {
//		log.Fatal(err)
//	}
//	// pass state.Fields to drift.Detect
package fetcher
