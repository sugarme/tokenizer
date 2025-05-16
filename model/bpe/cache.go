package bpe

import (
	"sync"
)

// Cache is a map with read-write mutex included
// to hold map of `word` strings
// E.g. https://tour.golang.org/concurrency/9
// NOTE: can we you sync.Map struct instead???
type Cache struct {
	mux sync.RWMutex
	// cmap     map[interface{}]interface{}
	cmap     map[string]Word
	Capacity int
}

type CacheItem struct {
	// Key   interface{}
	// Value interface{}
	Key   string
	Value Word // `word` string
}

// NewCache create an empty Cache with a specified capacity
func NewCache(capacity int) *Cache {
	return &Cache{
		// cmap:     make(map[interface{}]interface{}, capacity),
		cmap:     make(map[string]Word, capacity),
		Capacity: capacity,
	}
}

// Clear clears the cache
func (c *Cache) Clear() {
	c.mux.Lock()
	defer c.mux.Unlock()

	// Create a new map instead of using delete loop
	c.cmap = make(map[string]Word, c.Capacity)
}

// Get returns the value associated with the given key
func (c *Cache) Get(key string) (Word, bool) {
	c.mux.RLock()
	defer c.mux.RUnlock()
	word, ok := c.cmap[key]

	return word, ok
}

// GetValues returns slices of values associated with input keys
func (c *Cache) GetValues(keys []string) []Word {
	c.mux.RLock() // Use read lock for concurrent reads
	defer c.mux.RUnlock()

	var res []Word
	res = make([]Word, len(keys)) // Pre-allocate slice for better performance

	for i, k := range keys {
		res[i] = c.cmap[k]
	}

	return res
}

// SetValues sets values in the cache, respecting capacity limits
func (c *Cache) SetValues(values []CacheItem) {
	c.mux.Lock()
	defer c.mux.Unlock()

	// Check if we're already at capacity
	if len(c.cmap) >= c.Capacity {
		return
	}

	// Calculate how many items we can add
	remaining := c.Capacity - len(c.cmap)
	if remaining <= 0 {
		return
	}

	// Add items up to capacity
	for i, v := range values {
		if i >= remaining {
			break
		}
		c.cmap[v.Key] = v.Value
	}
}

// GetSize returns the current number of items in the cache
func (c *Cache) GetSize() int {
	c.mux.RLock()
	defer c.mux.RUnlock()
	return len(c.cmap)
}

// IsFull returns true if the cache has reached its capacity
func (c *Cache) IsFull() bool {
	c.mux.RLock()
	defer c.mux.RUnlock()
	return len(c.cmap) >= c.Capacity
}
