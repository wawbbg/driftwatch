// Package verify checks that a live service configuration matches
// a stored baseline, emitting a structured result for each field.
package verify

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/driftwatch/internal/baseline"
	"github.com/driftwatch/internal/diff"
)

// Result holds the outcome of a single verification run.
type Result struct {
	Service   string
	CheckedAt time.Time
	Diffs     []diff.Difference
	OK        bool
}

// String returns a human-readable summary of the result.
func (r Result) String() string {
	if r.OK {
		return fmt.Sprintf("%s: OK (no drift)", r.Service)
	}
	return fmt.Sprintf("%s: DRIFT detected (%d field(s))", r.Service, len(r.Diffs))
}

// Verifier compares live config against a stored baseline.
type Verifier struct {
	dir string
	out io.Writer
}

// New returns a Verifier that reads baselines from dir.
func New(dir string) *Verifier {
	return NewWithWriter(dir, os.Stdout)
}

// NewWithWriter returns a Verifier with a custom writer for log output.
func NewWithWriter(dir string, w io.Writer) *Verifier {
	if w == nil {
		w = os.Stdout
	}
	return &Verifier{dir: dir, out: w}
}

// Run loads the baseline for service and compares it against live.
// It returns a Result and any I/O error encountered.
func (v *Verifier) Run(service string, live map[string]any) (Result, error) {
	stored, err := baseline.Load(v.dir, service)
	if err != nil {
		return Result{}, fmt.Errorf("verify: load baseline for %q: %w", service, err)
	}

	diffs := diff.Maps(stored, live)
	res := Result{
		Service:   service,
		CheckedAt: time.Now().UTC(),
		Diffs:     diffs,
		OK:        !diff.HasDrift(diffs),
	}

	fmt.Fprintln(v.out, res.String())
	return res, nil
}
