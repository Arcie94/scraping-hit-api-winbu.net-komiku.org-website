package cache

import (
	"sync"
	"time"
)

// CacheItem represents a cached value with expiration
type cacheItem struct {
	Value      interface{}
	Expiration time.Time
}

// Cache provides in-memory caching with TTL support
type Cache struct {
	mu    sync.RWMutex
	items map[string]*cacheItem
	stats Stats
}

// Stats tracks cache performance metrics
type Stats struct {
	Hits   uint64
	Misses uint64
	Sets   uint64
}

// New creates a new Cache instance
func New() *Cache {
	c := &Cache{
		items: make(map[string]*cacheItem),
	}

	// Start background cleanup goroutine
	go c.cleanupExpired()

	return c
}

// Get retrieves a value from cache
func (c *Cache) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	item, found := c.items[key]
	if !found {
		c.stats.Misses++
		return nil, false
	}

	// Check if expired
	if time.Now().After(item.Expiration) {
		c.stats.Misses++
		return nil, false
	}

	c.stats.Hits++
	return item.Value, true
}

// Set stores a value in cache with TTL
func (c *Cache) Set(key string, value interface{}, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.items[key] = &cacheItem{
		Value:      value,
		Expiration: time.Now().Add(ttl),
	}
	c.stats.Sets++
}

// Delete removes a key from cache
func (c *Cache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.items, key)
}

// Clear removes all items from cache
func (c *Cache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.items = make(map[string]*cacheItem)
	// Reset stats
	c.stats = Stats{}
}

// Stats returns current cache statistics
func (c *Cache) GetStats() Stats {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.stats
}

// Size returns the number of items in cache
func (c *Cache) Size() int {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return len(c.items)
}

// cleanupExpired removes expired items periodically
func (c *Cache) cleanupExpired() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		c.mu.Lock()
		now := time.Now()

		for key, item := range c.items {
			if now.After(item.Expiration) {
				delete(c.items, key)
			}
		}

		c.mu.Unlock()
	}
}

// HitRate returns the cache hit ratio (0.0 to 1.0)
func (c *Cache) HitRate() float64 {
	c.mu.RLock()
	defer c.mu.RUnlock()

	total := c.stats.Hits + c.stats.Misses
	if total == 0 {
		return 0.0
	}

	return float64(c.stats.Hits) / float64(total)
}
