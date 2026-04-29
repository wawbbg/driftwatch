package signature

import (
	"fmt"
	"io"
	"os"
	"time"
)

// Result describes the outcome of a signature check for one service.
type Result struct {
	Service string
	Valid   bool
	StoredAt time.Time
	Message string
}

func (r Result) String() string {
	if r.Valid {
		return fmt.Sprintf("[OK]   %s (signed %s)", r.Service, r.StoredAt.Format(time.RFC3339))
	}
	return fmt.Sprintf("[FAIL] %s — %s", r.Service, r.Message)
}

// Checker verifies live fields against stored signatures.
type Checker struct {
	dir string
	out io.Writer
}

// NewChecker returns a Checker that reads signatures from dir.
func NewChecker(dir string) *Checker {
	return &Checker{dir: dir, out: os.Stdout}
}

// NewCheckerWithWriter returns a Checker writing output to w.
func NewCheckerWithWriter(dir string, w io.Writer) *Checker {
	return &Checker{dir: dir, out: w}
}

// Check verifies fields for service and writes a human-readable result line.
func (c *Checker) Check(service string, fields map[string]any) Result {
	e, err := Load(c.dir, service)
	if err != nil {
		r := Result{Service: service, Valid: false, Message: "no stored signature: " + err.Error()}
		fmt.Fprintln(c.out, r)
		return r
	}
	if !Verify(e, fields) {
		r := Result{Service: service, Valid: false, StoredAt: e.SignedAt, Message: "signature mismatch"}
		fmt.Fprintln(c.out, r)
		return r
	}
	r := Result{Service: service, Valid: true, StoredAt: e.SignedAt}
	fmt.Fprintln(c.out, r)
	return r
}
