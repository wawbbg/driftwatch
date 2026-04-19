// Package schedule provides periodic drift-check scheduling.
package schedule

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"
)

// Runner executes a job function on a fixed interval.
type Runner struct {
	interval time.Duration
	job      func(ctx context.Context) error
	out      io.Writer
}

// New creates a Runner with the given interval and job.
func New(interval time.Duration, job func(ctx context.Context) error) *Runner {
	return &Runner{
		interval: interval,
		job:      job,
		out:      os.Stderr,
	}
}

// NewWithWriter creates a Runner with a custom writer for log output.
func NewWithWriter(interval time.Duration, job func(ctx context.Context) error, w io.Writer) *Runner {
	return &Runner{interval: interval, job: job, out: w}
}

// Start runs the job immediately and then on every interval tick until ctx is cancelled.
func (r *Runner) Start(ctx context.Context) error {
	if err := r.run(ctx); err != nil {
		fmt.Fprintf(r.out, "schedule: job error: %v\n", err)
	}

	ticker := time.NewTicker(r.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			if err := r.run(ctx); err != nil {
				fmt.Fprintf(r.out, "schedule: job error: %v\n", err)
			}
		}
	}
}

func (r *Runner) run(ctx context.Context) error {
	fmt.Fprintf(r.out, "schedule: running job at %s\n", time.Now().Format(time.RFC3339))
	return r.job(ctx)
}
