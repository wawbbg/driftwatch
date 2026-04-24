// Package watch implements lightweight polling-based file watching for
// driftwatch. It monitors one or more file paths at a configurable interval
// and invokes a Handler callback whenever a file is created, modified, or
// deleted. This allows driftwatch to automatically re-run drift detection
// when local config definitions change on disk without requiring an external
// file-notification library.
package watch
