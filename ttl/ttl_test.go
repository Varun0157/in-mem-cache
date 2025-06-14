package ttl_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/Varun0157/cache-library/cache"
	"github.com/Varun0157/cache-library/cache/policies"
	"github.com/Varun0157/cache-library/ttl"
)

func TestTTLCache_Expiration(t *testing.T) {
	// 1. Setup: Create a core cache and wrap it with the TTL decorator.
	core := cache.New[string, string](10, policies.NewLRU[string]())
	ttlCache := ttl.NewCache(core)

	// 2. Act: Set a key with a very short TTL.
	ttlCache.SetWithTTL("a", "alpha", 50*time.Millisecond)

	// 3. Assert: Immediately after, the key should exist.
	val, found := ttlCache.Get("a")
	require.True(t, found)
	require.Equal(t, "alpha", val)

	// 4. Act: Wait for the TTL to expire.
	time.Sleep(100 * time.Millisecond)

	// 5. Assert: Now the key should be gone.
	_, found = ttlCache.Get("a")
	require.False(t, found, "Key 'a' should have expired and been evicted")
}

func TestTTLCache_NoExpiration(t *testing.T) {
	core := cache.New[string, string](10, policies.NewLRU[string]())
	ttlCache := ttl.NewCache(core)

	// Set a key with no TTL (or using the standard Set method).
	ttlCache.Set("b", "beta")
	time.Sleep(50 * time.Millisecond) // Wait a bit

	// The key should still exist.
	_, found := ttlCache.Get("b")
	require.True(t, found, "Key 'b' without TTL should not expire")
}

func TestTTLCache_UpdateTTL(t *testing.T) {
	core := cache.New[string, string](10, policies.NewLRU[string]())
	ttlCache := ttl.NewCache(core)

	// Set a key with a short TTL
	ttlCache.SetWithTTL("c", "gamma", 50*time.Millisecond)

	// Update the same key with a longer TTL
	ttlCache.SetWithTTL("c", "gamma", 200*time.Millisecond)

	// Wait for the original TTL to expire
	time.Sleep(100 * time.Millisecond)

	// The key should still exist because we updated the TTL
	val, found := ttlCache.Get("c")
	require.True(t, found)
	require.Equal(t, "gamma", val)

	// Wait for the new TTL to expire
	time.Sleep(150 * time.Millisecond)

	// Now the key should be gone
	_, found = ttlCache.Get("c")
	require.False(t, found, "Key 'c' should have expired after the updated TTL")
}

func TestTTLCache_Delete(t *testing.T) {
	core := cache.New[string, string](10, policies.NewLRU[string]())
	ttlCache := ttl.NewCache(core)

	// Set a key with TTL
	ttlCache.SetWithTTL("d", "delta", 100*time.Millisecond)

	// Delete it immediately
	ttlCache.Delete("d")

	// The key should be gone
	_, found := ttlCache.Get("d")
	require.False(t, found, "Key 'd' should be deleted")

	// Wait for what would have been the TTL
	time.Sleep(150 * time.Millisecond)

	// Still should be gone
	_, found = ttlCache.Get("d")
	require.False(t, found, "Key 'd' should still be deleted after TTL would have expired")
}
