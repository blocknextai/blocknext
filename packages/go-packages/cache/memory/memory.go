// Package memory provides an in-process cache.Service for single-instance
// deployments that do not want to run Redis. Everything lives in this process,
// so a second instance would not see these entries.
package memory

import (
	"context"
	"path"
	"strconv"
	"sync"
	"time"

	"github.com/blocknextai/go-packages/cache"
)

// cleanupInterval bounds how long an expired entry can linger before the
// janitor reclaims it; reads never observe an expired entry regardless.
const cleanupInterval = time.Minute

type entry struct {
	value     string
	expiresAt time.Time
}

func (e entry) expired(now time.Time) bool {
	return !e.expiresAt.IsZero() && now.After(e.expiresAt)
}

type provider struct {
	mu      sync.Mutex
	entries map[string]entry

	stop     chan struct{}
	stopOnce sync.Once
}

// New returns an in-process cache. The returned service starts a janitor
// goroutine that is released by Close.
func New() cache.Service {
	p := &provider{
		entries: make(map[string]entry),
		stop:    make(chan struct{}),
	}

	go p.janitor()

	return p
}

func (p *provider) janitor() {
	ticker := time.NewTicker(cleanupInterval)
	defer ticker.Stop()

	for {
		select {
		case <-p.stop:
			return
		case <-ticker.C:
			p.deleteExpired()
		}
	}
}

func (p *provider) deleteExpired() {
	now := time.Now()

	p.mu.Lock()
	defer p.mu.Unlock()

	for key, item := range p.entries {
		if item.expired(now) {
			delete(p.entries, key)
		}
	}
}

// get reads an entry, treating an expired one as absent. The caller holds the lock.
func (p *provider) get(key string) (entry, bool) {
	item, ok := p.entries[key]
	if !ok {
		return entry{}, false
	}
	if item.expired(time.Now()) {
		delete(p.entries, key)
		return entry{}, false
	}
	return item, true
}

func (p *provider) Ping(_ context.Context) error {
	return nil
}

// Get mirrors the Redis provider: a missing key is an empty string, not an error.
func (p *provider) Get(_ context.Context, key string) (string, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	item, ok := p.get(key)
	if !ok {
		return "", nil
	}
	return item.value, nil
}

func (p *provider) Set(_ context.Context, key string, value string, expiration time.Duration) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.entries[key] = entry{value: value, expiresAt: expiresAt(expiration)}
	return nil
}

func (p *provider) Incr(ctx context.Context, key string) (int64, error) {
	return p.addTo(key, 1)
}

func (p *provider) Decr(ctx context.Context, key string) (int64, error) {
	return p.addTo(key, -1)
}

// addTo applies delta to a counter, keeping the entry's existing expiry as
// Redis INCR/DECR do.
func (p *provider) addTo(key string, delta int64) (int64, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	var current int64
	expiresAt := time.Time{}

	if item, ok := p.get(key); ok {
		parsed, err := strconv.ParseInt(item.value, 10, 64)
		if err != nil {
			return 0, cache.ErrNotAnInteger
		}
		current = parsed
		expiresAt = item.expiresAt
	}

	current += delta
	p.entries[key] = entry{value: strconv.FormatInt(current, 10), expiresAt: expiresAt}

	return current, nil
}

func (p *provider) Delete(_ context.Context, key string) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	delete(p.entries, key)
	return nil
}

func (p *provider) GetAndDelete(_ context.Context, key string) (string, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	item, ok := p.get(key)
	if !ok {
		return "", nil
	}
	delete(p.entries, key)

	return item.value, nil
}

func (p *provider) Exists(_ context.Context, key string) (bool, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	_, ok := p.get(key)
	return ok, nil
}

// Keys matches the glob syntax Redis uses, which path.Match implements for the
// patterns this project relies on ("prefix:*").
func (p *provider) Keys(_ context.Context, pattern string) ([]string, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	now := time.Now()
	keys := make([]string, 0, len(p.entries))

	for key, item := range p.entries {
		if item.expired(now) {
			continue
		}
		matched, err := path.Match(pattern, key)
		if err != nil {
			return nil, err
		}
		if matched {
			keys = append(keys, key)
		}
	}

	return keys, nil
}

func (p *provider) Expire(_ context.Context, key string, expiration time.Duration) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	item, ok := p.get(key)
	if !ok {
		return nil
	}
	item.expiresAt = expiresAt(expiration)
	p.entries[key] = item

	return nil
}

// AcquireSemaphoreAtomic counts holders in a single key, refreshing the TTL on
// the first acquisition exactly as the Redis script does.
func (p *provider) AcquireSemaphoreAtomic(_ context.Context, key string, maxCount int64, ttl time.Duration) (bool, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	var current int64
	if item, ok := p.get(key); ok {
		parsed, err := strconv.ParseInt(item.value, 10, 64)
		if err != nil {
			return false, cache.ErrNotAnInteger
		}
		current = parsed
	}

	if current >= maxCount {
		return false, nil
	}

	current++
	item := entry{value: strconv.FormatInt(current, 10)}
	if current == 1 {
		item.expiresAt = expiresAt(ttl)
	} else if existing, ok := p.entries[key]; ok {
		item.expiresAt = existing.expiresAt
	}
	p.entries[key] = item

	return true, nil
}

// ReleaseSemaphoreAtomic drops a holder and removes the key once it empties.
func (p *provider) ReleaseSemaphoreAtomic(_ context.Context, key string) (int64, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	item, ok := p.get(key)
	if !ok {
		return 0, nil
	}

	current, err := strconv.ParseInt(item.value, 10, 64)
	if err != nil {
		return 0, cache.ErrNotAnInteger
	}

	current--
	if current <= 0 {
		delete(p.entries, key)
		return 0, nil
	}

	item.value = strconv.FormatInt(current, 10)
	p.entries[key] = item

	return current, nil
}

func (p *provider) Close() error {
	p.stopOnce.Do(func() { close(p.stop) })
	return nil
}

func expiresAt(expiration time.Duration) time.Time {
	if expiration <= 0 {
		return time.Time{}
	}
	return time.Now().Add(expiration)
}
