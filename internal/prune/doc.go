// Package prune provides retention-based cleanup for driftwatch data
// directories. It removes snapshot, history, and audit files that are older
// than a configured duration, keeping storage bounded over long-running
// deployments.
//
// Usage:
//
//	p := prune.New(7*24*time.Hour, os.Stdout)
//	result, err := p.Run("/var/lib/driftwatch", "my-service")
package prune
