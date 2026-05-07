// Package label provides key-value label management for driftwatch services.
// Labels can be attached to services and used for filtering, grouping, and reporting.
package label

import (
	"fmt"
	"sort"
	"strings"
)

// Set holds a collection of string key-value labels.
type Set map[string]string

// New returns an empty label Set.
func New() Set {
	return make(Set)
}

// FromPairs parses a slice of "key=value" strings into a Set.
// Malformed pairs (missing "=") are silently skipped.
func FromPairs(pairs []string) Set {
	s := New()
	for _, p := range pairs {
		k, v, ok := strings.Cut(p, "=")
		if !ok {
			continue
		}
		s[normalise(k)] = strings.TrimSpace(v)
	}
	return s
}

// Add sets a label, normalising the key to lowercase-trimmed form.
func (s Set) Add(key, value string) {
	s[normalise(key)] = strings.TrimSpace(value)
}

// Get returns the value for key and whether it was present.
func (s Set) Get(key string) (string, bool) {
	v, ok := s[normalise(key)]
	return v, ok
}

// Delete removes a label by key.
func (s Set) Delete(key string) {
	delete(s, normalise(key))
}

// Matches reports whether all labels in filter are present with matching values in s.
func (s Set) Matches(filter Set) bool {
	for k, v := range filter {
		if got, ok := s[k]; !ok || got != v {
			return false
		}
	}
	return true
}

// Pairs returns the labels as sorted "key=value" strings.
func (s Set) Pairs() []string {
	out := make([]string, 0, len(s))
	for k, v := range s {
		out = append(out, fmt.Sprintf("%s=%s", k, v))
	}
	sort.Strings(out)
	return out
}

func normalise(key string) string {
	return strings.ToLower(strings.TrimSpace(key))
}
