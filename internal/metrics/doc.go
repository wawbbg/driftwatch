// Package metrics provides lightweight instrumentation for driftwatch runs.
//
// Each drift detection pass can be recorded as a Run, which captures the
// service name, timestamp, field counts, drift status, and elapsed time.
//
// Runs can be emitted as newline-delimited JSON for downstream ingestion
// (e.g. log aggregators, dashboards) or printed as a human-readable summary.
//
// Example:
//
//	col := metrics.New()
//	start := time.Now()
//	// ... run detection ...
//	run := metrics.NewRun("payments", 12, 2, start)
//	col.Summary(run)
//	_ = col.Record(run)
package metrics
