package taskrunner

import (
	"time"
)

type OutputStore struct {
	cache *LRUCache
}

func NewOutputStore() *OutputStore {
	const (
		defaultMaxSize = 10000
		defaultTTL     = 1 * time.Hour
	)
	return &OutputStore{
		cache: NewLRUCache(defaultMaxSize, defaultTTL),
	}
}

func (r *OutputStore) Store(nodeKey string, outputs []map[string]any) {
	r.cache.Store(nodeKey, outputs)
}

func (r *OutputStore) Get(nodeKey string) ([]map[string]any, bool) {
	return r.cache.Get(nodeKey)
}
