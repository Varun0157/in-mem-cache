package policies

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLRUPolicy(t *testing.T) {
	policy := NewLRU[string]()

	// Test OnAdd
	policy.OnAdd("a")
	policy.OnAdd("b")
	policy.OnAdd("c")

	// Test OnAccess - should move accessed item to back
	policy.OnAccess("a")

	// Test OnEvict - should evict least recently used item
	evicted := policy.OnEvict()
	require.Equal(t, "b", evicted)

	// Test OnRemove
	policy.OnRemove("a")

	// Test OnEvict after removal
	evicted = policy.OnEvict()
	require.Equal(t, "c", evicted)
}

func TestLRUPolicy_Empty(t *testing.T) {
	policy := NewLRU[string]()

	// Test OnEvict on empty policy
	evicted := policy.OnEvict()
	require.Equal(t, "", evicted)
}

func TestLRUPolicy_DuplicateAdd(t *testing.T) {
	policy := NewLRU[string]()

	// Add same key multiple times
	policy.OnAdd("a")
	policy.OnAdd("a")
	policy.OnAdd("a")

	// Should only have one instance
	evicted := policy.OnEvict()
	require.Equal(t, "a", evicted)

	// Should be empty after
	evicted = policy.OnEvict()
	require.Equal(t, "", evicted)
}
