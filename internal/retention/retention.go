// Package retention enforces data retention limits on stored drift records,
// removing entries that exceed a configurable maximum age.
package retention

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

// Policy describes how long drift records should be kept.
type Policy struct {
	// MaxAge is the maximum age of a record before it is eligible for removal.
	MaxAge time.Duration
	// BaseDir is the root directory under which per-service record files live.
	BaseDir string
}

// Result summarises the outcome of an enforcement run.
type Result struct {
	Removed int
	Errors  int
}

func (r Result) String() string {
	return fmt.Sprintf("removed=%d errors=%d", r.Removed, r.Errors)
}

// Enforcer applies a retention Policy to stored records.
type Enforcer struct {
	policy Policy
	now    func() time.Time
	w      io.Writer
}

// New returns an Enforcer using the given Policy.
func New(p Policy) *Enforcer {
	return NewWithWriter(p, os.Stdout)
}

// NewWithWriter returns an Enforcer that writes log output to w.
func NewWithWriter(p Policy, w io.Writer) *Enforcer {
	return &Enforcer{policy: p, now: time.Now, w: w}
}

// Enforce scans BaseDir for files older than MaxAge and removes them.
// It returns a Result summarising what was done.
func (e *Enforcer) Enforce() Result {
	var res Result

	entries, err := os.ReadDir(e.policy.BaseDir)
	if err != nil {
		if os.IsNotExist(err) {
			return res
		}
		fmt.Fprintf(e.w, "retention: read dir error: %v\n", err)
		res.Errors++
		return res
	}

	cutoff := e.now().Add(-e.policy.MaxAge)

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		info, err := entry.Info()
		if err != nil {
			fmt.Fprintf(e.w, "retention: stat %s: %v\n", entry.Name(), err)
			res.Errors++
			continue
		}
		if info.ModTime().Before(cutoff) {
			path := filepath.Join(e.policy.BaseDir, entry.Name())
			if err := os.Remove(path); err != nil {
				fmt.Fprintf(e.w, "retention: remove %s: %v\n", path, err)
				res.Errors++
				continue
			}
			fmt.Fprintf(e.w, "retention: removed %s (age %s)\n", path, e.now().Sub(info.ModTime()).Round(time.Second))
			res.Removed++
		}
	}
	return res
}
