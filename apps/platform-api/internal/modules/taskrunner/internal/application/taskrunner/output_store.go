package taskrunner

import (
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
)

type OutputStore struct {
	cache   *LRUCache
	mu      sync.RWMutex
	indexes map[uuid.UUID]map[string][]int
}

func scopedKey(taskID uuid.UUID, key string) string {
	var builder strings.Builder
	builder.WriteString(taskID.String())
	builder.WriteString("|")
	builder.WriteString(key)
	return builder.String()
}

func NewOutputStore() *OutputStore {
	const (
		defaultMaxSize = 10000
		defaultTTL     = 1 * time.Hour
	)
	return &OutputStore{
		cache:   NewLRUCache(defaultMaxSize, defaultTTL),
		indexes: make(map[uuid.UUID]map[string][]int),
	}
}

func (r *OutputStore) Store(taskID uuid.UUID, nodeKey string, outputs []map[string]any) {
	r.cache.Store(scopedKey(taskID, nodeKey), outputs)
}

func (r *OutputStore) Get(taskID uuid.UUID, nodeKey string) ([]map[string]any, bool) {
	return r.cache.Get(scopedKey(taskID, nodeKey))
}

func (r *OutputStore) StoreBranch(taskID uuid.UUID, nodeKey, handle string, outputs []map[string]any, absolute []int) {
	branchKey := BuildBranchKey(nodeKey, handle)
	r.cache.Store(scopedKey(taskID, branchKey), outputs)

	r.mu.Lock()
	defer r.mu.Unlock()
	if r.indexes[taskID] == nil {
		r.indexes[taskID] = make(map[string][]int)
	}
	r.indexes[taskID][branchKey] = absolute
}

func (r *OutputStore) Indexes(taskID uuid.UUID, branchKey string) ([]int, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	absolute, ok := r.indexes[taskID][branchKey]
	return absolute, ok
}

func (r *OutputStore) Release(taskID uuid.UUID) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.indexes, taskID)
}
