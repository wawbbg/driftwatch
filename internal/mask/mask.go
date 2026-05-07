// Package mask provides field-level masking for config maps,
// replacing sensitive values with a fixed placeholder before display or export.
package mask

import "strings"

const Placeholder = "***"

// Masker holds the set of field keys that should be masked.
type Masker struct {
	keys map[string]struct{}
}

// defaultKeys are masked unless overridden.
var defaultKeys = []string{"password", "secret", "token", "apikey", "api_key", "private_key"}

// New returns a Masker using the default sensitive keys.
func New() *Masker {
	return NewWithKeys(defaultKeys)
}

// NewWithKeys returns a Masker that masks exactly the provided keys (case-insensitive).
func NewWithKeys(keys []string) *Masker {
	m := &Masker{keys: make(map[string]struct{}, len(keys))}
	for _, k := range keys {
		m.keys[strings.ToLower(k)] = struct{}{}
	}
	return m
}

// IsSensitive reports whether key should be masked.
func (m *Masker) IsSensitive(key string) bool {
	_, ok := m.keys[strings.ToLower(key)]
	return ok
}

// Apply returns a shallow copy of cfg with sensitive values replaced by Placeholder.
// The original map is never mutated.
func (m *Masker) Apply(cfg map[string]string) map[string]string {
	out := make(map[string]string, len(cfg))
	for k, v := range cfg {
		if m.IsSensitive(k) {
			out[k] = Placeholder
		} else {
			out[k] = v
		}
	}
	return out
}

// Keys returns the sorted list of keys that will be masked.
func (m *Masker) Keys() []string {
	out := make([]string, 0, len(m.keys))
	for k := range m.keys {
		out = append(out, k)
	}
	return out
}
