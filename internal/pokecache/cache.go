package pokecache

import (
	"bytes"
	"sync"
	"time"
)

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

type Cache struct {
	mu       sync.Mutex
	entries  map[string]cacheEntry
	interval time.Duration

	done      chan struct{}
	closeOnce sync.Once
}

func NewCache(interval time.Duration) *Cache {
	c := &Cache{
		entries:  make(map[string]cacheEntry),
		interval: interval,
		done:     make(chan struct{}),
	}
	go c.reapLoop()
	return c
}

func (c *Cache) Add(key string, val []byte) {
	entry := cacheEntry{
		createdAt: time.Now(),
		val:       bytes.Clone(val),
	}

	c.mu.Lock()
	defer c.mu.Unlock()
	c.entries[key] = entry
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.Lock()
	entry, ok := c.entries[key]
	c.mu.Unlock()

	if !ok {
		return nil, false
	}
	return bytes.Clone(entry.val), true
}

func (c *Cache) Close() {
	c.closeOnce.Do(func() { close(c.done) })
}

func (c *Cache) reapLoop() {
	ticker := time.NewTicker(c.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			c.reap(time.Now().Add(-c.interval))
		case <-c.done:
			return
		}
	}
}

func (c *Cache) reap(cutoff time.Time) {
	c.mu.Lock()
	defer c.mu.Unlock()

	for key, entry := range c.entries {
		if entry.createdAt.Before(cutoff) {
			delete(c.entries, key)
		}
	}
}
