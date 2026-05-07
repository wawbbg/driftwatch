package checkpoint

import (
	"fmt"
	"io"
	"os"
	"text/tabwriter"
)

// Reporter writes checkpoint listings to a writer.
type Reporter struct {
	w io.Writer
}

// NewReporter returns a Reporter that writes to w.
// If w is nil, os.Stdout is used.
func NewReporter(w io.Writer) *Reporter {
	if w == nil {
		w = os.Stdout
	}
	return &Reporter{w: w}
}

// Write prints all checkpoints for service stored in dir.
func (r *Reporter) Write(dir, service string) error {
	names, err := List(dir, service)
	if err != nil {
		return err
	}

	if len(names) == 0 {
		fmt.Fprintf(r.w, "no checkpoints found for %s\n", service)
		return nil
	}

	tw := tabwriter.NewWriter(r.w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "NAME\tSERVICE\tCREATED")

	for _, name := range names {
		e, err := Load(dir, service, name)
		if err != nil {
			fmt.Fprintf(tw, "%s\t%s\t(error: %v)\n", name, service, err)
			continue
		}
		fmt.Fprintf(tw, "%s\t%s\t%s\n", e.Name, e.Service, e.CreatedAt.Format("2006-01-02 15:04:05"))
	}

	return tw.Flush()
}
