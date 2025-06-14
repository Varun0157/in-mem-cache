# In-Memory Cache Library

A thread-safe, extensible, in-memory caching library in Go that supports multiple eviction policies.

## Features

- Thread-safe operations
- Generic key-value storage
- Multiple eviction policies:
  - LRU (Least Recently Used)
  - FIFO (First In, First Out)
  - LIFO (Last In, First Out)
- Extensible design for custom eviction policies
- High performance with minimal allocations

## Installation

```bash
go get github.com/Varun0157/cache-library
```

## Usage

### Basic Usage

```go
import (
    "github.com/Varun0157/cache-library/cache"
    "github.com/Varun0157/cache-library/cache/policies"
)

// Create a new cache with LRU policy
lruCache := cache.New[string, int](100, policies.NewLRU[string]())

// Set values
lruCache.Set("key1", 1)
lruCache.Set("key2", 2)

// Get values
if val, found := lruCache.Get("key1"); found {
    fmt.Println("Value:", val)
}

// Delete values
lruCache.Delete("key1")
```

### Using Different Eviction Policies

```go
// FIFO Cache
fifoCache := cache.New[string, int](100, policies.NewFIFO[string]())

// LIFO Cache
lifoCache := cache.New[string, int](100, policies.NewLIFO[string]())
```

### Creating Custom Eviction Policies

You can create custom eviction policies by implementing the `EvictionPolicy` interface:

```go
type MyPolicy[K comparable] struct {
    // Your custom fields
}

func (p *MyPolicy[K]) OnAdd(key K) {
    // Implementation
}

func (p *MyPolicy[K]) OnAccess(key K) {
    // Implementation
}

func (p *MyPolicy[K]) OnRemove(key K) {
    // Implementation
}

func (p *MyPolicy[K]) OnEvict() K {
    // Implementation
    return keyToEvict
}

// Usage
myCache := cache.New[string, int](100, &MyPolicy[string]{})
```

## Performance

The library is designed for high performance with minimal allocations. The cache operations are thread-safe and use a combination of a map for O(1) lookups and a linked list for maintaining the eviction order.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the LICENSE file for details. 
