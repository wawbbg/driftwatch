// Package tag provides utilities for tagging drift results with
// metadata such as environment, team, or severity labels.
package tag

import (
	"fmt"
	"sort"
	"strings"
)

// Tag represents a key-value metadata label.
type Tag struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// String returns the tag in "key=value" format.
func (t Tag) String() string {
	return fmt.Sprintf("%s=%s", t.Key, t.Value)
}

// Set holds a collection of unique tags keyed by name.
type Set struct {
	tags map[string]string
}

// New creates an empty tag Set.
func New() *Set {
	return &Set{tags: make(map[string]string)}
}

// Add inserts or overwrites a tag with the given key and value.
func (s *Set) Add(key, value string) {
	s.tags[strings.ToLower(strings.TrimSpace(key))] = strings.TrimSpace(value)
}

// Parse parses a slice of "key=value" strings into the Set.
// Entries that do not contain "=" are silently skipped.
func (s *Set) Parse(pairs []string) {
	for _, p := range pairs {
		parts := strings.SplitN(p, "=", 2)
		if len(parts) != 2 {
			continue
		}
		s.Add(parts[0], parts[1])
	}
}

// Get returns the value for a key and whether it was found.
func (s *Set) Get(key string) (string, bool) {
	v, ok := s.tags[strings.ToLower(strings.TrimSpace(key))]
	return v, ok
}

// All returns all tags as a sorted slice.
func (s *Set) All() []Tag {
	out := make([]Tag, 0, len(s.tags))
	for k, v := range s.tags {
		out = append(out, Tag{Key: k, Value: v})
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Key < out[j].Key })
	return out
}

// Len returns the number of tags in the set.
func (s *Set) Len() int { return len(s.tags) }
