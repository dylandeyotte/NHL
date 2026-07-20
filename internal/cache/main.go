package cache

import (
	"errors"
	"net/url"
	"sync"
	"time"
)

type Cache struct {
	entries map[string]cacheEntry
	mu      *sync.Mutex
}

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

func NewCache(interval time.Duration) Cache {
	cache := Cache{
		entries: make(map[string]cacheEntry),
		mu:      &sync.Mutex{},
	}
	go cache.killLoop(interval)

	return cache
}

func (c *Cache) Add(key string, val []byte) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if !isValidURL(key) {
		return errors.New("Bad URL")
	}

	c.entries[key] = cacheEntry{
		createdAt: time.Now(),
		val:       val,
	}
	return nil
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	entry, ok := c.entries[key]

	return entry.val, ok
}

func (c *Cache) killLoop(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		c.mu.Lock()

		for key, entry := range c.entries {
			duration := time.Since(entry.createdAt)
			if duration > interval {
				delete(c.entries, key)
			}
		}
		c.mu.Unlock()
	}
}

func isValidURL(str string) bool {
	url, err := url.ParseRequestURI(str)
	if err != nil {
		return false
	}

	if url.Scheme != "http" && url.Scheme != "https" {
		return false
	}

	return url.Host != ""
}
