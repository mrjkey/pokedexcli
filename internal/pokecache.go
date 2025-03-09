package pokecache

import "time"

type Cache map[string]cacheEntry

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}
