package cache

import (
	"L0-wb/internal/models"
	"container/list"
	"sync"
)

type cacheItem struct {
	key   string
	value *models.Order
}

type Cache interface {
	Set(key string, order *models.Order)
	Get(key string) (*models.Order, bool)
	Close()
}

type lruCache struct {
	capacity int
	items    map[string]*list.Element
	queue    *list.List
	mutex    sync.RWMutex
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		items:    make(map[string]*list.Element),
		queue:    list.New(),
	}
}

func (c *lruCache) Set(key string, order *models.Order) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if elem, exists := c.items[key]; exists {
		c.queue.MoveToFront(elem)
		elem.Value.(*cacheItem).value = order
		return
	}

	// Add new item
	item := &cacheItem{key: key, value: order}
	elem := c.queue.PushFront(item)
	c.items[key] = elem

	// Evict if over capacity
	if c.queue.Len() > c.capacity {
		c.evictOldest()
	}
}

func (c *lruCache) Get(key string) (*models.Order, bool) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	if elem, exists := c.items[key]; exists {
		c.queue.MoveToFront(elem)
		return elem.Value.(*cacheItem).value, true
	}
	return nil, false
}

func (c *lruCache) evictOldest() {
	if elem := c.queue.Back(); elem != nil {
		c.queue.Remove(elem)
		item := elem.Value.(*cacheItem)
		delete(c.items, item.key)
	}
}

func (c *lruCache) Close() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.items = nil
	c.queue = nil
}
