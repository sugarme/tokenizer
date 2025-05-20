package bpe

import (
	"sync"
	"testing"
)

func TestNewCache(t *testing.T) {
	capacity := 10
	cache := NewCache(capacity)

	if cache.Capacity != capacity {
		t.Errorf("Expected capacity %d, got %d", capacity, cache.Capacity)
	}

	if len(cache.cmap) != 0 {
		t.Errorf("Expected empty map, got map with %d items", len(cache.cmap))
	}
}

func TestCache_Clear(t *testing.T) {
	cache := NewCache(5)
	word1 := NewWord()
	word2 := NewWord()
	cache.SetValues([]CacheItem{
		{Key: "test1", Value: *word1},
		{Key: "test2", Value: *word2},
	})

	cache.Clear()

	if len(cache.cmap) != 0 {
		t.Errorf("Expected empty map after Clear(), got map with %d items", len(cache.cmap))
	}
}

func TestCache_GetValues(t *testing.T) {
	cache := NewCache(5)
	word1 := NewWord()
	word2 := NewWord()
	word3 := NewWord()
	testItems := []CacheItem{
		{Key: "test1", Value: *word1},
		{Key: "test2", Value: *word2},
		{Key: "test3", Value: *word3},
	}
	cache.SetValues(testItems)

	tests := []struct {
		name     string
		keys     []string
		expected []Word
	}{
		{
			name:     "Get existing values",
			keys:     []string{"test1", "test2"},
			expected: []Word{*word1, *word2},
		},
		{
			name:     "Get non-existing value",
			keys:     []string{"nonexistent"},
			expected: []Word{{}},
		},
		{
			name:     "Get mixed existing and non-existing values",
			keys:     []string{"test1", "nonexistent", "test3"},
			expected: []Word{*word1, {}, *word3},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			values := cache.GetValues(tt.keys)
			if len(values) != len(tt.expected) {
				t.Errorf("Expected %d values, got %d", len(tt.expected), len(values))
				return
			}
			for i, v := range values {
				if len(v.Symbols) != len(tt.expected[i].Symbols) {
					t.Errorf("Expected value %v at index %d, got %v", tt.expected[i], i, v)
				}
			}
		})
	}
}

func TestCache_SetValues(t *testing.T) {
	tests := []struct {
		name           string
		capacity       int
		initialItems   []CacheItem
		itemsToAdd     []CacheItem
		expectedLength int
	}{
		{
			name:     "Add items within capacity",
			capacity: 5,
			itemsToAdd: []CacheItem{
				{Key: "test1", Value: *NewWord()},
				{Key: "test2", Value: *NewWord()},
			},
			expectedLength: 2,
		},
		{
			name:     "Add items exceeding capacity",
			capacity: 2,
			itemsToAdd: []CacheItem{
				{Key: "test1", Value: *NewWord()},
				{Key: "test2", Value: *NewWord()},
				{Key: "test3", Value: *NewWord()},
			},
			expectedLength: 2,
		},
		{
			name:     "Add items to full cache",
			capacity: 2,
			initialItems: []CacheItem{
				{Key: "test1", Value: *NewWord()},
				{Key: "test2", Value: *NewWord()},
			},
			itemsToAdd: []CacheItem{
				{Key: "test3", Value: *NewWord()},
			},
			expectedLength: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cache := NewCache(tt.capacity)
			if len(tt.initialItems) > 0 {
				cache.SetValues(tt.initialItems)
			}
			cache.SetValues(tt.itemsToAdd)

			if len(cache.cmap) != tt.expectedLength {
				t.Errorf("Expected %d items in cache, got %d", tt.expectedLength, len(cache.cmap))
			}
		})
	}
}

func TestCache_GetSize(t *testing.T) {
	cache := NewCache(5)
	word1 := NewWord()
	word2 := NewWord()
	cache.SetValues([]CacheItem{
		{Key: "test1", Value: *word1},
		{Key: "test2", Value: *word2},
	})

	size := cache.GetSize()
	if size != 2 {
		t.Errorf("Expected size 2, got %d", size)
	}
}

func TestCache_IsFull(t *testing.T) {
	cache := NewCache(2)
	word1 := NewWord()
	word2 := NewWord()
	cache.SetValues([]CacheItem{
		{Key: "test1", Value: *word1},
		{Key: "test2", Value: *word2},
	})

	full := cache.IsFull()
	if !full {
		t.Errorf("Expected cache to be full, but it is not")
	}
}

func TestCache_ConcurrentReads(t *testing.T) {
	t.Parallel()

	cache := NewCache(100)
	word := NewWord()

	// Fill cache with some initial data
	for i := 0; i < 50; i++ {
		cache.SetValues([]CacheItem{
			{Key: "test" + string(rune(i)), Value: *word},
		})
	}

	var wg sync.WaitGroup
	readers := 100
	iterations := 1000

	wg.Add(readers)
	for i := 0; i < readers; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < iterations; j++ {
				// Generate random keys to read
				keys := make([]string, 5)
				for k := 0; k < 5; k++ {
					keys[k] = "test" + string(rune(j%50))
				}
				_ = cache.GetValues(keys)
			}
		}()
	}

	wg.Wait()
}

func TestCache_ConcurrentWrites(t *testing.T) {
	t.Parallel()

	cache := NewCache(100)
	word := NewWord()

	var wg sync.WaitGroup
	writers := 10
	iterations := 100

	wg.Add(writers)
	for i := 0; i < writers; i++ {
		go func(writerID int) {
			defer wg.Done()
			for j := 0; j < iterations; j++ {
				key := "test" + string(rune(writerID)) + string(rune(j))
				cache.SetValues([]CacheItem{
					{Key: key, Value: *word},
				})
			}
		}(i)
	}

	wg.Wait()

	// Verify that the cache size is within capacity
	if len(cache.cmap) > cache.Capacity {
		t.Errorf("Cache size %d exceeds capacity %d", len(cache.cmap), cache.Capacity)
	}
}

func TestCache_ConcurrentReadsAndWrites(t *testing.T) {
	t.Parallel()

	cache := NewCache(100)
	word := NewWord()

	// Fill cache with some initial data
	for i := 0; i < 50; i++ {
		cache.SetValues([]CacheItem{
			{Key: "test" + string(rune(i)), Value: *word},
		})
	}

	var wg sync.WaitGroup
	readers := 50
	writers := 10
	iterations := 100

	// Start readers
	wg.Add(readers)
	for i := 0; i < readers; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < iterations; j++ {
				keys := make([]string, 5)
				for k := 0; k < 5; k++ {
					keys[k] = "test" + string(rune(j%50))
				}
				_ = cache.GetValues(keys)
			}
		}()
	}

	// Start writers
	wg.Add(writers)
	for i := 0; i < writers; i++ {
		go func(writerID int) {
			defer wg.Done()
			for j := 0; j < iterations; j++ {
				key := "test" + string(rune(writerID)) + string(rune(j))
				cache.SetValues([]CacheItem{
					{Key: key, Value: *word},
				})
			}
		}(i)
	}

	wg.Wait()

	// Verify that the cache size is within capacity
	if len(cache.cmap) > cache.Capacity {
		t.Errorf("Cache size %d exceeds capacity %d", len(cache.cmap), cache.Capacity)
	}
}
