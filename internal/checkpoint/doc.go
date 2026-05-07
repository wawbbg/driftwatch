// Package checkpoint provides save/load/list operations for named
// drift-check checkpoints.
//
// A checkpoint captures a map of config fields for a named service at
// a specific point in time (e.g. "pre-deploy", "post-deploy"). It can
// later be loaded and compared against a live snapshot to detect
// changes introduced during a deployment window.
//
// Checkpoints are stored as JSON files under a configurable directory:
//
//	<dir>/<service>/<name>.json
package checkpoint
