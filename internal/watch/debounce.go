package watch

import (
	"sync"
	"time"
)

// Debouncer suppresses repeated invocations of a function within a quiet
// window, calling it only once after the activity has settled.
type Debouncer struct {
	mu      sync.Mutex
	delay   time.Duration
	timers  map[string]*time.Timer
	handler Handler
}

// NewDebouncer wraps handler so that rapid successive events for the same
// path are collapsed into a single call after delay has elapsed.
func NewDebouncer(delay time.Duration, handler Handler) *Debouncer {
	return &Debouncer{
		delay:   delay,
		timers:  make(map[string]*time.Timer),
		handler: handler,
	}
}

// Handle schedules handler to be called for e after the debounce delay.
// If another event for the same path arrives before the timer fires, the
// timer is reset and only the latest event is forwarded.
func (d *Debouncer) Handle(e Event) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	if t, ok := d.timers[e.Path]; ok {
		t.Stop()
	}

	d.timers[e.Path] = time.AfterFunc(d.delay, func() {
		d.mu.Lock()
		delete(d.timers, e.Path)
		d.mu.Unlock()
		_ = d.handler(e)
	})

	return nil
}
