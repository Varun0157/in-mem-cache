package policies

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLIFOPolicy(t *testing.T) {
	policy := NewLIFO[string]()

	// Test OnAdd
	policy.OnAdd("a")
	policy.OnAdd("b")
	policy.OnAdd("c")

	// Test OnAccess - should not change order in LIFO
	policy.OnAccess("a")

	// Test OnEvict - should evict last item added
	evicted := policy.OnEvict()
	require.Equal(t, "c", evicted)

	// Test OnRemove
	policy.OnRemove("b")

	// Test OnEvict after removal
	evicted = policy.OnEvict()
	require.Equal(t, "a", evicted)
}

func TestLIFOPolicy_Empty(t *testing.T) {
	policy := NewLIFO[string]()

	// Test OnEvict on empty policy
	evicted := policy.OnEvict()
	require.Equal(t, "", evicted)
}

func TestLIFOPolicy_DuplicateAdd(t *testing.T) {
	policy := NewLIFO[string]()

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

func TestLIFOPolicy_AccessDoesNotChangeOrder(t *testing.T) {
	policy := NewLIFO[string]()

	// Add items
	policy.OnAdd("a")
	policy.OnAdd("b")
	policy.OnAdd("c")

	// Access first item multiple times
	policy.OnAccess("a")
	policy.OnAccess("a")
	policy.OnAccess("a")

	// Should still evict in LIFO order
	evicted := policy.OnEvict()
	require.Equal(t, "c", evicted)

	evicted = policy.OnEvict()
	require.Equal(t, "b", evicted)

	evicted = policy.OnEvict()
	require.Equal(t, "a", evicted)
}
