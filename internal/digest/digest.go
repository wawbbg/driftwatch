// Package digest computes and compares deterministic hashes of service
// configuration maps, enabling quick equality checks before a full diff.
package digest

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"sort"
)

// Sum returns a stable SHA-256 hex digest of the provided map.
// Keys are sorted before serialisation so the result is deterministic
// regardless of the iteration order of the source map.
func Sum(m map[string]any) (string, error) {
	if m == nil {
		return emptyDigest(), nil
	}

	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	ordered := make([]any, 0, len(keys))
	for _, k := range keys {
		ordered = append(ordered, [2]any{k, m[k]})
	}

	b, err := json.Marshal(ordered)
	if err != nil {
		return "", fmt.Errorf("digest: marshal: %w", err)
	}

	h := sha256.Sum256(b)
	return hex.EncodeToString(h[:]), nil
}

// Equal reports whether two configuration maps produce the same digest.
// An error is returned only when serialisation of either map fails.
func Equal(a, b map[string]any) (bool, error) {
	da, err := Sum(a)
	if err != nil {
		return false, err
	}
	db, err := Sum(b)
	if err != nil {
		return false, err
	}
	return da == db, nil
}

// emptyDigest returns the SHA-256 hex digest of an empty JSON array.
func emptyDigest() string {
	h := sha256.Sum256([]byte("[]"))
	return hex.EncodeToString(h[:])
}
