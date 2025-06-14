package cache_test

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/Varun0157/in-mem-cache/cache"
	"github.com/Varun0157/in-mem-cache/cache/policies"
)

func TestCache_With_LRU_Policy(t *testing.T) {
	lruCache := cache.New[string, int](2, policies.NewLRU[string]())

	lruCache.Set("a", 1)
	lruCache.Set("b", 2)
	lruCache.Set("c", 3) // Evicts "a"

	_, found := lruCache.Get("a")
	require.False(t, found)

	val, found := lruCache.Get("b")
	require.True(t, found)
	require.Equal(t, 2, val)

	val, found = lruCache.Get("c")
	require.True(t, found)
	require.Equal(t, 3, val)
}

func TestCache_With_FIFO_Policy(t *testing.T) {
	fifoCache := cache.New[string, int](2, policies.NewFIFO[string]())

	fifoCache.Set("a", 1)
	fifoCache.Set("b", 2)
	fifoCache.Set("c", 3) // Evicts "a"

	_, found := fifoCache.Get("a")
	require.False(t, found)

	val, found := fifoCache.Get("b")
	require.True(t, found)
	require.Equal(t, 2, val)

	val, found = fifoCache.Get("c")
	require.True(t, found)
	require.Equal(t, 3, val)
}

func TestCache_With_LIFO_Policy(t *testing.T) {
	lifoCache := cache.New[string, int](2, policies.NewLIFO[string]())

	lifoCache.Set("a", 1)
	lifoCache.Set("b", 2)
	lifoCache.Set("c", 3) // Evicts "b"

	val, found := lifoCache.Get("a")
	require.True(t, found)
	require.Equal(t, 1, val)

	_, found = lifoCache.Get("b")
	require.False(t, found)

	val, found = lifoCache.Get("c")
	require.True(t, found)
	require.Equal(t, 3, val)
}

func TestCache_ConcurrentAccess(t *testing.T) {
	c := cache.New[int, int](100, policies.NewLRU[int]())
	var wg sync.WaitGroup
	numGoroutines := 10
	numOperations := 1000

	for i := range numGoroutines {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := range numOperations {
				key := id*numOperations + j
				c.Set(key, key)
				val, found := c.Get(key)
				if found {
					require.Equal(t, key, val)
				}
			}
		}(i)
	}

	wg.Wait()
}

func BenchmarkCacheSet(b *testing.B) {
	c := cache.New[int, int](1024, policies.NewLRU[int]())
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		c.Set(i, i)
	}
}

func BenchmarkCacheGet(b *testing.B) {
	c := cache.New[int, int](1024, policies.NewLRU[int]())
	for i := range 1000 {
		c.Set(i, i)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c.Get(i % 1000)
	}
}

func TestCache_ZeroCapacity(t *testing.T) {
	// Should default to capacity 1
	c := cache.New[string, int](0, policies.NewLRU[string]())

	c.Set("a", 1)
	c.Set("b", 2) // Should evict "a"

	_, found := c.Get("a")
	require.False(t, found)

	val, found := c.Get("b")
	require.True(t, found)
	require.Equal(t, 2, val)
}

func TestCache_UpdateExistingKey(t *testing.T) {
	c := cache.New[string, int](2, policies.NewLRU[string]())

	c.Set("a", 1)
	c.Set("a", 2) // Update existing key

	val, found := c.Get("a")
	require.True(t, found)
	require.Equal(t, 2, val)
}

func TestCache_Delete(t *testing.T) {
	c := cache.New[string, int](2, policies.NewLRU[string]())

	c.Set("a", 1)
	c.Set("b", 2)

	// Delete existing key
	c.Delete("a")
	_, found := c.Get("a")
	require.False(t, found)

	// Delete non-existent key should not panic
	c.Delete("c")

	// Other key should still exist
	val, found := c.Get("b")
	require.True(t, found)
	require.Equal(t, 2, val)
}

func TestCache_ComplexValueType(t *testing.T) {
	type ComplexValue struct {
		ID    int
		Data  []string
		Extra map[string]interface{}
	}

	c := cache.New[string, ComplexValue](2, policies.NewLRU[string]())

	val1 := ComplexValue{
		ID:    1,
		Data:  []string{"test1", "test2"},
		Extra: map[string]interface{}{"key": 123},
	}

	val2 := ComplexValue{
		ID:    2,
		Data:  []string{"test3"},
		Extra: map[string]interface{}{"key": "value"},
	}

	c.Set("a", val1)
	c.Set("b", val2)

	// Test retrieval
	retrieved, found := c.Get("a")
	require.True(t, found)
	require.Equal(t, val1.ID, retrieved.ID)
	require.Equal(t, val1.Data, retrieved.Data)
	require.Equal(t, val1.Extra, retrieved.Extra)
}

func TestCache_ConcurrentDelete(t *testing.T) {
	c := cache.New[int, int](100, policies.NewLRU[int]())
	var wg sync.WaitGroup
	numGoroutines := 10
	numOperations := 1000

	// First populate the cache
	for i := 0; i < numGoroutines*numOperations; i++ {
		c.Set(i, i)
	}

	for i := range numGoroutines {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := range numOperations {
				key := id*numOperations + j
				c.Delete(key)
				_, found := c.Get(key)
				require.False(t, found)
			}
		}(i)
	}

	wg.Wait()
}

func TestCache_ConcurrentMixedOperations(t *testing.T) {
	c := cache.New[int, int](100, policies.NewLRU[int]())
	var wg sync.WaitGroup
	numGoroutines := 10
	numOperations := 1000

	for i := range numGoroutines {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := range numOperations {
				key := id*numOperations + j

				// Mix of operations
				switch j % 3 {
				case 0:
					c.Set(key, key)
				case 1:
					c.Get(key)
				case 2:
					c.Delete(key)
				}
			}
		}(i)
	}

	wg.Wait()
}
