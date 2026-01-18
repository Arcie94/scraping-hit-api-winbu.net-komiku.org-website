package cache

import (
	"sync"
	"time"
)

// CacheItem represents a cached item with expiration
type CacheItem struct {
	Data      interface{}
	ExpiresAt time.Time
}

// Cache is a simple in-memory cache with TTL
type Cache struct {
	items sync.Map
}

// NewCache creates a new cache instance
func NewCache() *Cache {
	c := &Cache{}
	// Start cleanup goroutine
	go c.cleanup()
	return c
}

// Set stores an item in cache with TTL
func (c *Cache) Set(key string, value interface{}, ttl time.Duration) {
	c.items.Store(key, CacheItem{
		Data:      value,
		ExpiresAt: time.Now().Add(ttl),
	})
}

// Get retrieves an item from cache
func (c *Cache) Get(key string) (interface{}, bool) {
	val, ok := c.items.Load(key)
	if !ok {
		return nil, false
	}

	item := val.(CacheItem)

	// Check if expired
	if time.Now().After(item.ExpiresAt) {
		c.items.Delete(key)
		return nil, false
	}

	return item.Data, true
}

// Delete removes an item from cache
func (c *Cache) Delete(key string) {
	c.items.Delete(key)
}

// Clear removes all items from cache
func (c *Cache) Clear() {
	c.items = sync.Map{}
}

// cleanup runs periodically to remove expired items
func (c *Cache) cleanup() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		now := time.Now()
		c.items.Range(func(key, value interface{}) bool {
			item := value.(CacheItem)
			if now.After(item.ExpiresAt) {
				c.items.Delete(key)
			}
			return true
		})
	}
}
