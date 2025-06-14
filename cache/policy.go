package cache

// EvictionPolicy defines the public interface for a cache eviction strategy.
// Any custom policy must implement this interface.
type EvictionPolicy[K comparable] interface {
	// OnAdd is called when a new key is added to the cache.
	OnAdd(key K)

	// OnAccess is called when a key is accessed (e.g., via Get).
	OnAccess(key K)

	// OnRemove is called when a key is explicitly removed from the cache.
	OnRemove(key K)

	// OnEvict is called to determine which key should be evicted when the cache is full.
	// It should return the key to be removed.
	OnEvict() K
}
