package ttl

import (
	"sync"
	"time"

	"github.com/Varun0157/in-mem-cache/cache"
)

// ttlEntry stores the expiration time for a key.
type ttlEntry struct {
	expiresAt time.Time
}

// Cache is a decorator that adds TTL (Time-To-Live) functionality
// to any underlying cache that satisfies the cache.Cacheable interface.
type Cache[K comparable, V any] struct {
	// The underlying cache to store the actual key-value pairs.
	coreCache cache.Cacheable[K, V]

	mu          sync.RWMutex
	expirations map[K]ttlEntry // Stores only the expiration data
}

// NewCache creates a new TTL-enabled cache decorator.
// It wraps a core cache instance (like the one from your 'cache' package).
func NewCache[K comparable, V any](core cache.Cacheable[K, V]) *Cache[K, V] {
	return &Cache[K, V]{
		coreCache:   core,
		expirations: make(map[K]ttlEntry),
	}
}

// SetWithTTL adds a key-value pair to the cache with a specific TTL.
func (c *Cache[K, V]) SetWithTTL(key K, value V, ttl time.Duration) {
	c.coreCache.Set(key, value) // Set the value in the core cache

	c.mu.Lock()
	defer c.mu.Unlock()
	if ttl > 0 {
		c.expirations[key] = ttlEntry{expiresAt: time.Now().Add(ttl)}
	} else {
		// If TTL is zero or negative, it means no expiration.
		// We can remove it from our tracking map.
		delete(c.expirations, key)
	}
}

// Set is required to satisfy the cache.Cacheable interface.
// It sets a value with no expiration.
func (c *Cache[K, V]) Set(key K, value V) {
	c.SetWithTTL(key, value, 0)
}

// Get retrieves a value. It first checks for expiration.
func (c *Cache[K, V]) Get(key K) (V, bool) {
	c.mu.RLock()
	entry, hasExpiration := c.expirations[key]
	c.mu.RUnlock()

	// Check if the item has an expiration time and if it has passed.
	if hasExpiration && time.Now().After(entry.expiresAt) {
		// Item has expired. Delete it from both caches.
		c.Delete(key)
		var zeroV V
		return zeroV, false
	}

	// If not expired (or no expiration was set), get it from the core cache.
	return c.coreCache.Get(key)
}

// Delete removes a key from both the TTL tracker and the core cache.
func (c *Cache[K, V]) Delete(key K) {
	c.mu.Lock()
	delete(c.expirations, key)
	c.mu.Unlock()

	c.coreCache.Delete(key)
}

// Static assertion to ensure *ttl.Cache satisfies the cache.Cacheable interface.
var _ cache.Cacheable[any, any] = (*Cache[any, any])(nil)
