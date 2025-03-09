package pokecache

import "time"

type Cache struct {
	entries  map[string]cacheEntry
	interval time.Duration
}

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

func NewCache(interval time.Duration) Cache {
	cache := Cache{
		entries:  make(map[string]cacheEntry),
		interval: interval,
	}
	go cache.reapLoop()
	return cache
}

func (c *Cache) Add(key string, val []byte) {
	c.entries[key] = cacheEntry{
		createdAt: time.Now(),
		val:       val,
	}
}

func (c Cache) Get(key string) ([]byte, bool) {
	entry, ok := c.entries[key]
	if !ok {
		return []byte{}, false
	}
	return entry.val, true
}

func (c *Cache) reapLoop() {
	for {
		for key, value := range c.entries {
			if time.Since(value.createdAt) > c.interval {
				delete(c.entries, key)
			}
		}
		time.Sleep(c.interval)
	}
}
