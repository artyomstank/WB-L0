package cache

import (
	"L0-wb/internal/models"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCache_SetAndGet(t *testing.T) {
	c := NewCache(10)
	order := &models.Order{OrderUID: "test-123"}

	// Test Set and Get
	c.Set("test-123", order)
	got, exists := c.Get("test-123")

	assert.True(t, exists)
	assert.Equal(t, order, got)

	// Test non-existent key
	_, exists = c.Get("non-existent")
	assert.False(t, exists)
}

func TestCache_Eviction(t *testing.T) {
	// Create cache with size 2
	c := NewCache(2)

	// Add three orders to trigger eviction
	order1 := &models.Order{OrderUID: "1"}
	order2 := &models.Order{OrderUID: "2"}
	order3 := &models.Order{OrderUID: "3"}

	c.Set("1", order1)
	c.Set("2", order2)

	// Check both orders are in cache
	_, exists1 := c.Get("1")
	require.True(t, exists1)
	_, exists2 := c.Get("2")
	require.True(t, exists2)

	// Add third order to trigger eviction
	c.Set("3", order3)

	// Verify least recently used order (order1) was evicted
	_, exists1 = c.Get("1")
	assert.False(t, exists1, "order1 should be evicted")
	_, exists3 := c.Get("3")
	assert.True(t, exists3, "order3 should be in cache")
	_, exists2 = c.Get("2")
	assert.True(t, exists2, "order2 should still be in cache")
}

func TestCache_TTL(t *testing.T) {
	c := NewCache(10)
	order := &models.Order{OrderUID: "test-123"}

	c.Set("test-123", order)
	time.Sleep(150 * time.Millisecond)

	// Item should still exist after 150ms
	_, exists := c.Get("test-123")
	assert.True(t, exists)
}

func TestCache_Close(t *testing.T) {
	c := NewCache(10)
	c.Close() // Should not panic
}
