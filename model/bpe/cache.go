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

// Fresh create a fresh `Cache` with the same configuration
func (c *Cache) Fresh() {
	c = NewCache(c.Capacity)
}

// Clear clears the cache
func (c *Cache) Clear() {
	// NOTE:Should we just make a new map instead of use `for` loop
	// Ref. https://stackoverflow.com/questions/13812121
	for k := range c.cmap {
		delete(c.cmap, k)
	}
}

// GetValues returns slices of values associated with input keys
func (c *Cache) GetValues(keys []string) []Word {
	c.mux.Lock() // Lock so only one goroutine at a time can access
	defer c.mux.Unlock()

	var res []Word

	for _, k := range keys {
		res = append(res, c.cmap[k])
	}

	return res
}

func (c *Cache) SetValues(values []CacheItem) {

	// Before trying to acquire a write lock, we check if we are already at
	// capacity with a read handler.
	c.mux.Lock()
	defer c.mux.Unlock()

	if len(c.cmap) == c.Capacity {
		return
	}

	for _, v := range values {
		if len(c.cmap) == c.Capacity {
			break
		}

		c.cmap[v.Key] = v.Value
	}
}
