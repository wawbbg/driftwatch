// Package normalize provides utilities for normalizing config map keys and
// values before comparison, ensuring consistent drift detection regardless
// of minor formatting differences.
package normalize

import (
	"strings"
	"unicode"
)

// Options controls how normalization is applied.
type Options struct {
	// LowercaseKeys converts all map keys to lowercase before comparison.
	LowercaseKeys bool
	// TrimValues strips leading and trailing whitespace from string values.
	TrimValues bool
	// CollapseWhitespace replaces runs of whitespace in values with a single space.
	CollapseWhitespace bool
}

// DefaultOptions returns a sensible default normalization configuration.
func DefaultOptions() Options {
	return Options{
		LowercaseKeys:      true,
		TrimValues:         true,
		CollapseWhitespace: false,
	}
}

// Apply returns a new map with keys and values normalized according to opts.
// The original map is never mutated.
func Apply(m map[string]string, opts Options) map[string]string {
	if m == nil {
		return nil
	}
	out := make(map[string]string, len(m))
	for k, v := range m {
		if opts.LowercaseKeys {
			k = strings.ToLower(k)
		}
		if opts.TrimValues {
			v = strings.TrimSpace(v)
		}
		if opts.CollapseWhitespace {
			v = collapseWS(v)
		}
		out[k] = v
	}
	return out
}

// collapseWS replaces runs of whitespace characters with a single space.
func collapseWS(s string) string {
	var b strings.Builder
	inSpace := false
	for _, r := range s {
		if unicode.IsSpace(r) {
			if !inSpace {
				b.WriteRune(' ')
				inSpace = true
			}
		} else {
			b.WriteRune(r)
			inSpace = false
		}
	}
	return b.String()
}
