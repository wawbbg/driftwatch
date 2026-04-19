// Package baseline provides utilities for capturing and retrieving
// baseline service configurations used by driftwatch to compare
// against live state.
//
// A baseline is a point-in-time snapshot of a service's expected
// field values. It is stored as a JSON file on disk (one per service)
// and loaded during drift detection runs to establish the reference
// state.
//
// Typical usage:
//
//	err := baseline.Store(".driftwatch/baselines", "payments", liveFields)
//
//	entry, err := baseline.Load(".driftwatch/baselines", "payments")
package baseline
