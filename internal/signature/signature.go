// Package signature computes and verifies a canonical signature for a
// service's live configuration, allowing callers to detect whether a config
// has changed since it was last signed.
package signature

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/user/driftwatch/internal/digest"
)

// Entry holds the signature for a single service at a point in time.
type Entry struct {
	Service   string    `json:"service"`
	Hex       string    `json:"hex"`
	SignedAt  time.Time `json:"signed_at"`
}

// Sign computes the hex signature of fields and returns an Entry.
func Sign(service string, fields map[string]any) Entry {
	raw := digest.Sum(fields)
	return Entry{
		Service:  service,
		Hex:      hex.EncodeToString(raw[:]),
		SignedAt: time.Now().UTC(),
	}
}

// Verify returns true when fields produce the same hex as e.Hex.
func Verify(e Entry, fields map[string]any) bool {
	return digest.Equal(fields, decodeHex(e.Hex))
}

// Store persists an Entry to dir/<service>.sig.json.
func Store(dir string, e Entry) error {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("signature: mkdir %s: %w", dir, err)
	}
	path := filepath.Join(dir, e.Service+".sig.json")
	b, err := json.MarshalIndent(e, "", "  ")
	if err != nil {
		return fmt.Errorf("signature: marshal: %w", err)
	}
	return os.WriteFile(path, b, 0o644)
}

// Load reads the stored Entry for service from dir.
func Load(dir, service string) (Entry, error) {
	path := filepath.Join(dir, service+".sig.json")
	b, err := os.ReadFile(path)
	if err != nil {
		return Entry{}, fmt.Errorf("signature: read %s: %w", path, err)
	}
	var e Entry
	if err := json.Unmarshal(b, &e); err != nil {
		return Entry{}, fmt.Errorf("signature: unmarshal: %w", err)
	}
	return e, nil
}

func decodeHex(h string) [32]byte {
	var out [32]byte
	b, _ := hex.DecodeString(h)
	copy(out[:], b)
	return out
}
