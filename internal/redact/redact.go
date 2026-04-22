// Package redact provides utilities for masking sensitive fields
// in config maps before they are displayed or logged.
package redact

import "strings"

// DefaultSensitiveKeys is the list of field name substrings considered sensitive.
var DefaultSensitiveKeys = []string{
	"password",
	"secret",
	"token",
	"apikey",
	"api_key",
	"private",
	"credential",
}

const mask = "[REDACTED]"

// Redactor holds the set of sensitive key patterns.
type Redactor struct {
	keys []string
}

// New returns a Redactor using the provided key patterns.
// If keys is empty, DefaultSensitiveKeys is used.
func New(keys []string) *Redactor {
	if len(keys) == 0 {
		keys = DefaultSensitiveKeys
	}
	normalized := make([]string, len(keys))
	for i, k := range keys {
		normalized[i] = strings.ToLower(k)
	}
	return &Redactor{keys: normalized}
}

// IsSensitive reports whether the given field name matches any sensitive pattern.
func (r *Redactor) IsSensitive(field string) bool {
	lower := strings.ToLower(field)
	for _, k := range r.keys {
		if strings.Contains(lower, k) {
			return true
		}
	}
	return false
}

// Apply returns a copy of m with sensitive values replaced by the mask string.
func (r *Redactor) Apply(m map[string]any) map[string]any {
	out := make(map[string]any, len(m))
	for k, v := range m {
		if r.IsSensitive(k) {
			out[k] = mask
		} else {
			out[k] = v
		}
	}
	return out
}
