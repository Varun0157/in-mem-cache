package policies

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFIFOPolicy(t *testing.T) {
	policy := NewFIFO[string]()

	// Test OnAdd
	policy.OnAdd("a")
	policy.OnAdd("b")
	policy.OnAdd("c")

	// Test OnAccess - should not change order in FIFO
	policy.OnAccess("a")

	// Test OnEvict - should evict first item added
	evicted := policy.OnEvict()
	require.Equal(t, "a", evicted)

	// Test OnRemove
	policy.OnRemove("b")

	// Test OnEvict after removal
	evicted = policy.OnEvict()
	require.Equal(t, "c", evicted)
}

func TestFIFOPolicy_Empty(t *testing.T) {
	policy := NewFIFO[string]()

	// Test OnEvict on empty policy
	evicted := policy.OnEvict()
	require.Equal(t, "", evicted)
}

func TestFIFOPolicy_DuplicateAdd(t *testing.T) {
	policy := NewFIFO[string]()

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

func TestFIFOPolicy_AccessDoesNotChangeOrder(t *testing.T) {
	policy := NewFIFO[string]()

	// Add items
	policy.OnAdd("a")
	policy.OnAdd("b")
	policy.OnAdd("c")

	// Access middle item multiple times
	policy.OnAccess("b")
	policy.OnAccess("b")
	policy.OnAccess("b")

	// Should still evict in FIFO order
	evicted := policy.OnEvict()
	require.Equal(t, "a", evicted)

	evicted = policy.OnEvict()
	require.Equal(t, "b", evicted)

	evicted = policy.OnEvict()
	require.Equal(t, "c", evicted)
}
