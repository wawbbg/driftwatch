// Package watch provides file-system watching for config files,
// triggering drift detection when source definitions change on disk.
package watch

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"time"
)

// EventType describes the kind of file-system event observed.
type EventType string

const (
	EventModified EventType = "modified"
	EventCreated  EventType = "created"
	EventDeleted  EventType = "deleted"
)

// Event represents a single file-system change.
type Event struct {
	Path      string
	Type      EventType
	Timestamp time.Time
}

// Handler is called whenever a watched file changes.
type Handler func(e Event) error

// Watcher polls a set of file paths for changes.
type Watcher struct {
	paths    []string
	interval time.Duration
	handler  Handler
	logger   *log.Logger
}

// New creates a Watcher that polls paths every interval.
func New(paths []string, interval time.Duration, handler Handler) *Watcher {
	return NewWithWriter(paths, interval, handler, os.Stderr)
}

// NewWithWriter creates a Watcher with a custom log writer.
func NewWithWriter(paths []string, interval time.Duration, handler Handler, w io.Writer) *Watcher {
	return &Watcher{
		paths:    paths,
		interval: interval,
		handler:  handler,
		logger:   log.New(w, "[watch] ", log.LstdFlags),
	}
}

// Start begins polling and blocks until ctx is cancelled.
func (w *Watcher) Start(ctx context.Context) error {
	modTimes := make(map[string]time.Time)

	// Seed initial mod times so we don't fire on first tick.
	for _, p := range w.paths {
		if fi, err := os.Stat(p); err == nil {
			modTimes[p] = fi.ModTime()
		}
	}

	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			w.poll(modTimes)
		}
	}
}

func (w *Watcher) poll(modTimes map[string]time.Time) {
	for _, p := range w.paths {
		fi, err := os.Stat(p)
		if err != nil {
			if os.IsNotExist(err) {
				if _, seen := modTimes[p]; seen {
					delete(modTimes, p)
					w.fire(Event{Path: p, Type: EventDeleted, Timestamp: time.Now()})
				}
			}
			continue
		}

		prev, seen := modTimes[p]
		switch {
		case !seen:
			modTimes[p] = fi.ModTime()
			w.fire(Event{Path: p, Type: EventCreated, Timestamp: fi.ModTime()})
		case fi.ModTime().After(prev):
			modTimes[p] = fi.ModTime()
			w.fire(Event{Path: p, Type: EventModified, Timestamp: fi.ModTime()})
		}
	}
}

func (w *Watcher) fire(e Event) {
	if err := w.handler(e); err != nil {
		w.logger.Printf("handler error for %s (%s): %v", e.Path, e.Type, err)
	}
}

// String returns a human-readable description of the event.
func (e Event) String() string {
	return fmt.Sprintf("%s %s at %s", e.Type, e.Path, e.Timestamp.Format(time.RFC3339))
}
