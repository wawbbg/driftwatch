// Package prune removes stale snapshot and history entries that exceed a
// configurable retention window.
package prune

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

// Result summarises what was removed during a prune run.
type Result struct {
	Service  string
	Removed  int
	Retained int
}

func (r Result) String() string {
	return fmt.Sprintf("service=%s removed=%d retained=%d", r.Service, r.Removed, r.Retained)
}

// Pruner deletes files older than Retention inside a base directory.
type Pruner struct {
	Retention time.Duration
	w         io.Writer
}

// New returns a Pruner that writes log lines to w.
func New(retention time.Duration, w io.Writer) *Pruner {
	if w == nil {
		w = io.Discard
	}
	return &Pruner{Retention: retention, w: w}
}

// Run walks baseDir/<service>/ and removes any regular file whose modification
// time is older than p.Retention. It returns a Result describing the outcome.
func (p *Pruner) Run(baseDir, service string) (Result, error) {
	dir := filepath.Join(baseDir, service)
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return Result{Service: service}, nil
		}
		return Result{}, fmt.Errorf("prune: read dir %s: %w", dir, err)
	}

	cutoff := time.Now().Add(-p.Retention)
	var result Result
	result.Service = service

	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		info, err := e.Info()
		if err != nil {
			continue
		}
		path := filepath.Join(dir, e.Name())
		if info.ModTime().Before(cutoff) {
			if err := os.Remove(path); err != nil {
				fmt.Fprintf(p.w, "prune: remove %s: %v\n", path, err)
				continue
			}
			fmt.Fprintf(p.w, "pruned %s\n", path)
			result.Removed++
		} else {
			result.Retained++
		}
	}
	return result, nil
}
