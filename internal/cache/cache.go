package cache

import (
	"context"
	"sync"
)

type cacheItem struct {
	value any
	count int
}

type LFUCache struct {
	mu       sync.RWMutex
	items    map[string]*cacheItem
	capacity int
}

func NewLFUCache(capacity int) *LFUCache {
	return &LFUCache{
		items:    make(map[string]*cacheItem),
		capacity: capacity,
	}
}

func (c *LFUCache) Get(ctx context.Context, key string) (any, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	item, ok := c.items[key]
	if !ok {
		return nil, false
	}
	item.count++
	return item.value, true
}

func (c *LFUCache) Set(ctx context.Context, key string, value any) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if item, ok := c.items[key]; ok {
		item.value = value
		item.count++
		return
	}

	if len(c.items) >= c.capacity {
		c.removeLFU()
	}

	c.items[key] = &cacheItem{
		value: value,
		count: 1,
	}
}

func (c *LFUCache) Delete(ctx context.Context, key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.items, key)
}

func (c *LFUCache) DeletePrefix(ctx context.Context, prefix string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	for k := range c.items {
		if startsWith(k, prefix) {
			delete(c.items, k)
		}
	}
}

func (c *LFUCache) removeLFU() {
	var lfuKey string
	var minCount = int(^uint(0) >> 1) 

	for k, v := range c.items {
		if v.count < minCount {
			minCount = v.count
			lfuKey = k
		}
	}
	if lfuKey != "" {
		delete(c.items, lfuKey)
	}
}

func startsWith(s, prefix string) bool {
	if len(s) < len(prefix) {
		return false
	}
	return s[:len(prefix)] == prefix
}
