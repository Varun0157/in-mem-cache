package main

import (
	"fmt"
	"time"

	"github.com/Varun0157/in-mem-cache/cache"
	"github.com/Varun0157/in-mem-cache/cache/policies"
	"github.com/Varun0157/in-mem-cache/ttl"
)

// SimplePolicy is a custom eviction policy that implements the cache.EvictionPolicy interface.
type SimplePolicy[K comparable] struct {
	keys []K
}

func NewSimplePolicy[K comparable]() cache.EvictionPolicy[K] {
	return &SimplePolicy[K]{}
}

func (p *SimplePolicy[K]) OnAdd(key K) {
	p.keys = append(p.keys, key)
}

func (p *SimplePolicy[K]) OnAccess(key K) {
	// No-op
}

func (p *SimplePolicy[K]) OnRemove(key K) {
	for i, k := range p.keys {
		if k == key {
			p.keys = append(p.keys[:i], p.keys[i+1:]...)
			break
		}
	}
}

func (p *SimplePolicy[K]) OnEvict() K {
	if len(p.keys) == 0 {
		var zero K
		return zero
	}
	keyToEvict := p.keys[0]
	p.keys = p.keys[1:]
	return keyToEvict
}

func main() {
	fmt.Println("--- Demonstrating LRU Cache ---")
	lruCache := cache.New[string, int](2, policies.NewLRU[string]())
	lruCache.Set("a", 1)
	lruCache.Set("b", 2)
	lruCache.Set("c", 3) // Evicts "a"

	if _, found := lruCache.Get("a"); !found {
		fmt.Println("LRU Cache correctly evicted 'a'")
	}

	if val, found := lruCache.Get("b"); found {
		fmt.Printf("LRU Cache contains 'b' with value: %d\n", val)
	}

	if val, found := lruCache.Get("c"); found {
		fmt.Printf("LRU Cache contains 'c' with value: %d\n", val)
	}

	fmt.Println("\n--- Demonstrating FIFO Cache ---")
	fifoCache := cache.New[string, int](2, policies.NewFIFO[string]())
	fifoCache.Set("a", 1)
	fifoCache.Set("b", 2)
	fifoCache.Set("c", 3) // Evicts "a"

	if _, found := fifoCache.Get("a"); !found {
		fmt.Println("FIFO Cache correctly evicted 'a'")
	}

	if val, found := fifoCache.Get("b"); found {
		fmt.Printf("FIFO Cache contains 'b' with value: %d\n", val)
	}

	if val, found := fifoCache.Get("c"); found {
		fmt.Printf("FIFO Cache contains 'c' with value: %d\n", val)
	}

	fmt.Println("\n--- Demonstrating Custom Policy ---")
	customCache := cache.New[string, int](2, NewSimplePolicy[string]())
	customCache.Set("a", 1)
	customCache.Set("b", 2)
	customCache.Set("c", 3) // Evicts "a"

	if _, found := customCache.Get("a"); !found {
		fmt.Println("Custom Cache correctly evicted 'a'")
	}

	if val, found := customCache.Get("b"); found {
		fmt.Printf("Custom Cache contains 'b' with value: %d\n", val)
	}

	if val, found := customCache.Get("c"); found {
		fmt.Printf("Custom Cache contains 'c' with value: %d\n", val)
	}

	fmt.Println("\n--- Demonstrating TTL Decorator Cache ---")
	// 1. Create the core cache first.
	coreCache := cache.New[string, int](10, policies.NewLRU[string]())

	// 2. Wrap it with the TTL decorator.
	ttlCache := ttl.NewCache[string, int](coreCache)

	// 3. Set a value that expires in 100ms.
	ttlCache.SetWithTTL("transient", 123, 100*time.Millisecond)
	fmt.Println("Set 'transient' key with 100ms TTL.")

	// 4. Set a value with no expiration.
	ttlCache.Set("permanent", 456)
	fmt.Println("Set 'permanent' key with no TTL.")

	// 5. Verify both values exist initially.
	if val, found := ttlCache.Get("transient"); found {
		fmt.Printf("'transient' key exists with value: %d\n", val)
	}
	if val, found := ttlCache.Get("permanent"); found {
		fmt.Printf("'permanent' key exists with value: %d\n", val)
	}

	// 6. Wait for the TTL to expire.
	fmt.Println("Waiting for TTL to expire...")
	time.Sleep(150 * time.Millisecond)

	// 7. Check values after expiration.
	if _, found := ttlCache.Get("transient"); !found {
		fmt.Println("'transient' key has correctly expired.")
	}
	if val, found := ttlCache.Get("permanent"); found {
		fmt.Printf("'permanent' key still exists with value: %d\n", val)
	}
}
