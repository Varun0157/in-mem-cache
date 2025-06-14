package cache

import (
	"log"
	"sync"
)

// Cache is a thread-safe, generic, in-memory cache.
type Cache[K comparable, V any] struct {
	capacity int
	policy   EvictionPolicy[K]

	mu      sync.RWMutex
	storage map[K]V
}

// New creates a new Cache with a given capacity and eviction policy.
func New[K comparable, V any](capacity int, policy EvictionPolicy[K]) *Cache[K, V] {
	if capacity <= 0 {
		log.Println("Cache capacity must be greater than 0, defaulting to 1")
		capacity = 1
	}
	return &Cache[K, V]{
		capacity: capacity,
		policy:   policy,
		storage:  make(map[K]V, capacity),
	}
}

// Set adds or updates a value in the cache.
func (c *Cache[K, V]) Set(key K, value V) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Check if the key already exists
	if _, ok := c.storage[key]; ok {
		// Update the value directly
		c.storage[key] = value
		// Notify the policy of the access
		c.policy.OnAccess(key)
		return
	}

	// Check if the cache is at capacity BEFORE adding
	if len(c.storage) >= c.capacity {
		// Ask the policy for the key to evict
		keyToEvict := c.policy.OnEvict()
		// Remove the evicted key from storage
		delete(c.storage, keyToEvict)
	}

	// Add the new key-value pair to storage
	c.storage[key] = value
	// Notify the policy that a new key was added
	c.policy.OnAdd(key)
}

// Get retrieves a value from the cache.
func (c *Cache[K, V]) Get(key K) (V, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	// Look up the value directly in the storage map
	value, ok := c.storage[key]
	if !ok {
		var zeroV V
		return zeroV, false
	}

	// If found, notify the policy of the access
	c.policy.OnAccess(key)

	// Return the value
	return value, true
}

// Delete removes a value from the cache.
func (c *Cache[K, V]) Delete(key K) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Check if the key exists before trying to delete
	if _, ok := c.storage[key]; !ok {
		return
	}

	// Delete from the storage map
	delete(c.storage, key)

	// Notify the policy of the removal
	c.policy.OnRemove(key)
}

// Static assertion to ensure *Cache satisfies the Cacheable interface.
var _ Cacheable[any, any] = (*Cache[any, any])(nil)
