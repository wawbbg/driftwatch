// Package cache provides a simple in-memory and file-backed cache for
// storing fetched service configurations to reduce redundant HTTP requests.
package cache

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// Entry holds a cached value along with its expiry time.
type Entry struct {
	Value     map[string]any `json:"value"`
	FetchedAt time.Time      `json:"fetched_at"`
	TTL       time.Duration  `json:"ttl"`
}

// Expired reports whether the cache entry has passed its TTL.
func (e Entry) Expired() bool {
	return time.Since(e.FetchedAt) > e.TTL
}

// Cache is a thread-safe store for service config snapshots.
type Cache struct {
	mu      sync.RWMutex
	entries map[string]Entry
	dir     string
}

// New returns a Cache that persists entries under dir.
// Pass an empty string to use a purely in-memory cache.
func New(dir string) *Cache {
	return &Cache{
		entries: make(map[string]Entry),
		dir:     dir,
	}
}

// Set stores value for key with the given TTL, writing through to disk when a
// directory was provided.
func (c *Cache) Set(key string, value map[string]any, ttl time.Duration) error {
	e := Entry{Value: value, FetchedAt: time.Now(), TTL: ttl}
	c.mu.Lock()
	c.entries[key] = e
	c.mu.Unlock()

	if c.dir != "" {
		return c.persist(key, e)
	}
	return nil
}

// Get returns the cached entry for key. ok is false when the key is absent or
// the entry has expired.
func (c *Cache) Get(key string) (Entry, bool) {
	c.mu.RLock()
	e, ok := c.entries[key]
	c.mu.RUnlock()

	if ok && !e.Expired() {
		return e, true
	}

	if c.dir != "" {
		if loaded, err := c.load(key); err == nil && !loaded.Expired() {
			c.mu.Lock()
			c.entries[key] = loaded
			c.mu.Unlock()
			return loaded, true
		}
	}
	return Entry{}, false
}

func (c *Cache) persist(key string, e Entry) error {
	if err := os.MkdirAll(c.dir, 0o755); err != nil {
		return fmt.Errorf("cache: mkdir %s: %w", c.dir, err)
	}
	data, err := json.Marshal(e)
	if err != nil {
		return fmt.Errorf("cache: marshal: %w", err)
	}
	return os.WriteFile(filepath.Join(c.dir, key+".json"), data, 0o644)
}

func (c *Cache) load(key string) (Entry, error) {
	data, err := os.ReadFile(filepath.Join(c.dir, key+".json"))
	if err != nil {
		return Entry{}, err
	}
	var e Entry
	if err := json.Unmarshal(data, &e); err != nil {
		return Entry{}, fmt.Errorf("cache: unmarshal: %w", err)
	}
	return e, nil
}
