package policies

import (
	"container/list"
	"sync"

	"github.com/Varun0157/in-mem-cache/cache"
)

type lruPolicy[K comparable] struct {
	mu     sync.RWMutex
	keys   *list.List
	keyMap map[K]*list.Element
}

// NewLRU creates a new LRU eviction policy.
func NewLRU[K comparable]() cache.EvictionPolicy[K] {
	return &lruPolicy[K]{
		keys:   list.New(),
		keyMap: make(map[K]*list.Element),
	}
}

func (p *lruPolicy[K]) OnAdd(key K) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if element, exists := p.keyMap[key]; exists {
		p.keys.MoveToBack(element)
	} else {
		element := p.keys.PushBack(key)
		p.keyMap[key] = element
	}
}

func (p *lruPolicy[K]) OnAccess(key K) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if element, exists := p.keyMap[key]; exists {
		p.keys.MoveToBack(element)
	}
}

func (p *lruPolicy[K]) OnRemove(key K) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if element, exists := p.keyMap[key]; exists {
		p.keys.Remove(element)
		delete(p.keyMap, key)
	}
}

func (p *lruPolicy[K]) OnEvict() K {
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
