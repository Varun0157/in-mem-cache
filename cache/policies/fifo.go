package policies

import (
	"container/list"
	"sync"

	"github.com/Varun0157/in-mem-cache/cache"
)

type fifoPolicy[K comparable] struct {
	mu     sync.RWMutex
	keys   *list.List
	keyMap map[K]*list.Element
}

// NewFIFO creates a new FIFO eviction policy.
func NewFIFO[K comparable]() cache.EvictionPolicy[K] {
	return &fifoPolicy[K]{
		keys:   list.New(),
		keyMap: make(map[K]*list.Element),
	}
}

func (p *fifoPolicy[K]) OnAdd(key K) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if _, exists := p.keyMap[key]; !exists {
		element := p.keys.PushBack(key)
		p.keyMap[key] = element
	}
}

func (p *fifoPolicy[K]) OnAccess(key K) {
	// No-op for FIFO
}

func (p *fifoPolicy[K]) OnRemove(key K) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if element, exists := p.keyMap[key]; exists {
		p.keys.Remove(element)
		delete(p.keyMap, key)
	}
}

func (p *fifoPolicy[K]) OnEvict() K {
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
