// Package audit records structured audit log entries for driftwatch.
//
// Each service maintains its own append-only JSONL audit file stored under
// a configurable directory. Entries capture the timestamp, service name,
// number of diffs detected, whether a policy error occurred, and an optional
// human-readable message.
//
// Usage:
//
//	// Record an audit event after a drift check.
//	err := audit.Record(".driftwatch/audit", "payments", 2, false, "")
//
//	// Retrieve the full audit history for a service.
//	entries, err := audit.List(".driftwatch/audit", "payments")
package audit
