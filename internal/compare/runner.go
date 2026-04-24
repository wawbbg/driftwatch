package compare

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/driftwatch/internal/snapshot"
)

// Runner loads two named snapshots for a service and compares them.
type Runner struct {
	dir    string
	writer io.Writer
}

// NewRunner creates a Runner that reads snapshots from dir.
func NewRunner(dir string) *Runner {
	return &Runner{dir: dir, writer: os.Stdout}
}

// NewRunnerWithWriter creates a Runner with a custom writer for output.
func NewRunnerWithWriter(dir string, w io.Writer) *Runner {
	return &Runner{dir: dir, writer: w}
}

// Run loads the snapshots identified by beforeTag and afterTag for the
// given service, compares them, and returns the resulting Delta.
func (r *Runner) Run(_ context.Context, service, beforeTag, afterTag string) (Delta, error) {
	before, err := snapshot.Load(r.dir, service+"."+beforeTag)
	if err != nil {
		return Delta{}, fmt.Errorf("compare: load before snapshot %q: %w", beforeTag, err)
	}

	after, err := snapshot.Load(r.dir, service+"."+afterTag)
	if err != nil {
		return Delta{}, fmt.Errorf("compare: load after snapshot %q: %w", afterTag, err)
	}

	delta := Snapshots(service, before, after)

	fmt.Fprintln(r.writer, delta.String())
	return delta, nil
}
