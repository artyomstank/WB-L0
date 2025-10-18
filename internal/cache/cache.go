package cache

import (
	"L0-wb/internal/models"
	"sync"
	"time"
)

// Cache defines interface for caching operations
type Cache interface {
	Set(key string, order *models.Order)
	Get(key string) (*models.Order, bool)
	Close()
}

const (
	defaultMaxSize       = 1000
	defaultCleanupPeriod = 5 * time.Minute
	defaultTTL           = 30 * time.Minute
)

type cacheEntry struct {
	order      *models.Order
	lastAccess time.Time
}

// InMemoryCache implements Cache interface
type InMemoryCache struct {
	data        map[string]cacheEntry
	mu          sync.RWMutex
	maxSize     int
	ttl         time.Duration
	stopCleanup chan struct{}
}

func NewCache(maxSize int) Cache {
	if maxSize <= 0 {
		maxSize = defaultMaxSize
	}

	c := &InMemoryCache{
		data:        make(map[string]cacheEntry),
		maxSize:     maxSize,
		ttl:         defaultTTL,
		stopCleanup: make(chan struct{}),
	}

	go c.startCleanup()
	return c
}

func (c *InMemoryCache) startCleanup() {
	ticker := time.NewTicker(defaultCleanupPeriod)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			c.cleanup()
		case <-c.stopCleanup:
			return
		}
	}
}

func (c *InMemoryCache) cleanup() {
	c.mu.Lock()
	defer c.mu.Unlock()

	now := time.Now()
	for key, entry := range c.data {
		if now.Sub(entry.lastAccess) > c.ttl {
			delete(c.data, key)
		}
	}
}

func (c *InMemoryCache) Set(key string, order *models.Order) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if len(c.data) >= c.maxSize {
		c.evictOldest()
	}

	c.data[key] = cacheEntry{
		order:      order,
		lastAccess: time.Now(),
	}
}

func (c *InMemoryCache) Get(key string) (*models.Order, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	entry, exists := c.data[key]
	if !exists {
		return nil, false
	}

	entry.lastAccess = time.Now()
	c.data[key] = entry

	return entry.order, true
}

func (c *InMemoryCache) evictOldest() {
	var oldestKey string
	var oldestAccess time.Time
	first := true

	for key, entry := range c.data {
		if first || entry.lastAccess.Before(oldestAccess) {
			oldestKey = key
			oldestAccess = entry.lastAccess
			first = false
		}
	}

	if oldestKey != "" {
		delete(c.data, oldestKey)
	}
}

func (c *InMemoryCache) Close() {
	close(c.stopCleanup)
}
