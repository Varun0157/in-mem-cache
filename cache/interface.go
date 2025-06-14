package cache

// Cacheable defines the public contract for any cache implementation.
// Both the core Cache and decorators will implement this interface.
type Cacheable[K comparable, V any] interface {
	Get(key K) (V, bool)
	Set(key K, value V)
	Delete(key K)
}
