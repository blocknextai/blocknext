package taskrunner

import (
	"container/list"
	"sync"
	"time"
)

type cacheEntry struct {
	key        string
	value      []map[string]any
	expiration time.Time
}

type LRUCache struct {
	maxSize     int
	ttl         time.Duration
	cache       map[string]*list.Element
	lruList     *list.List
	mu          sync.RWMutex
	stopCleanup chan struct{}
}

func NewLRUCache(maxSize int, ttl time.Duration) *LRUCache {
	if maxSize <= 0 {
		maxSize = 10000
	}
	if ttl <= 0 {
		ttl = 1 * time.Hour
	}

	cache := &LRUCache{
		maxSize:     maxSize,
		ttl:         ttl,
		cache:       make(map[string]*list.Element),
		lruList:     list.New(),
		stopCleanup: make(chan struct{}),
	}

	go cache.cleanupExpiredEntries()

	return cache
}

func (c *LRUCache) Store(key string, value []map[string]any) {
	expiration := time.Now().UTC().Add(c.ttl)

	c.mu.Lock()
	defer c.mu.Unlock()

	if element, exists := c.cache[key]; exists {
		entry, ok := element.Value.(*cacheEntry)
		if !ok {
			c.removeElement(element)
		} else {
			entry.value = value
			entry.expiration = expiration
			c.lruList.MoveToFront(element)
			return
		}
	}

	if c.lruList.Len() >= c.maxSize {
		c.evictLRU()
	}

	entry := &cacheEntry{
		key:        key,
		value:      value,
		expiration: expiration,
	}
	element := c.lruList.PushFront(entry)
	c.cache[key] = element
}

func (c *LRUCache) Get(key string) ([]map[string]any, bool) {
	utcNow := time.Now().UTC()

	c.mu.Lock()
	defer c.mu.Unlock()

	element, exists := c.cache[key]
	if !exists {
		return nil, false
	}

	entry, ok := element.Value.(*cacheEntry)
	if !ok {
		c.removeElement(element)
		return nil, false
	}

	if utcNow.After(entry.expiration) {
		c.removeElement(element)
		return nil, false
	}

	c.lruList.MoveToFront(element)
	return entry.value, true
}

func (c *LRUCache) evictLRU() {
	element := c.lruList.Back()
	if element != nil {
		c.removeElement(element)
	}
}

func (c *LRUCache) removeElement(element *list.Element) {
	c.lruList.Remove(element)
	entry, ok := element.Value.(*cacheEntry)
	if ok {
		delete(c.cache, entry.key)
	}
}

func (c *LRUCache) cleanupExpiredEntries() {
	cleanupInterval := max(c.ttl/2, 1*time.Minute)

	ticker := time.NewTicker(cleanupInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			c.performCleanup()
		case <-c.stopCleanup:
			return
		}
	}
}

func (c *LRUCache) performCleanup() {
	c.mu.Lock()
	defer c.mu.Unlock()

	utcNow := time.Now().UTC()
	var toRemove []*list.Element

	for element := c.lruList.Front(); element != nil; element = element.Next() {
		entry, ok := element.Value.(*cacheEntry)
		if !ok {
			toRemove = append(toRemove, element)
			continue
		}
		if utcNow.After(entry.expiration) {
			toRemove = append(toRemove, element)
		}
	}

	for _, element := range toRemove {
		c.removeElement(element)
	}
}

func (c *LRUCache) Close() {
	close(c.stopCleanup)
}

func (c *LRUCache) Len() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.lruList.Len()
}

func (c *LRUCache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.cache = make(map[string]*list.Element)
	c.lruList = list.New()
}
