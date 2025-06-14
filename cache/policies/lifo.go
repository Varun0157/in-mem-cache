package policies

import (
	"container/list"
	"sync"

	"github.com/Varun0157/cache-library/cache"
)

type lifoPolicy[K comparable] struct {
	mu     sync.RWMutex
	keys   *list.List
	keyMap map[K]*list.Element
}

// NewLIFO creates a new LIFO eviction policy.
func NewLIFO[K comparable]() cache.EvictionPolicy[K] {
	return &lifoPolicy[K]{
		keys:   list.New(),
		keyMap: make(map[K]*list.Element),
	}
}

func (p *lifoPolicy[K]) OnAdd(key K) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if _, exists := p.keyMap[key]; !exists {
		element := p.keys.PushFront(key)
		p.keyMap[key] = element
	}
}

func (p *lifoPolicy[K]) OnAccess(key K) {
	// No-op for LIFO
}

func (p *lifoPolicy[K]) OnRemove(key K) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if element, exists := p.keyMap[key]; exists {
		p.keys.Remove(element)
		delete(p.keyMap, key)
	}
}

func (p *lifoPolicy[K]) OnEvict() K {
	p.mu.Lock()
	defer p.mu.Unlock()

	element := p.keys.Front()
	if element == nil {
		var zero K
		return zero
	}
	key := element.Value.(K)
	p.keys.Remove(element)
	delete(p.keyMap, key)
	return key
}
