package pokecache

import (
	"testing"
	"time"
)

func TestEntryRemovedOnTime(t *testing.T) {
	duration := time.Millisecond * 20
	cache := NewCache(duration)
	entry := "entry"
	entryB := []byte(entry)
	url := "hi"
	cache.Add(url, entryB)

	time.Sleep(duration / 2)
	entryC, ok := cache.Get(url)
	if !ok {
		t.Errorf("entry not in cache")
		t.Fail()
	}

	output := string(entryC)
	if output != entry {
		t.Errorf("entry does not match")
		t.Fail()
	}

	time.Sleep(duration)
	_, ok = cache.Get(url)
	if ok {
		t.Errorf("entry not in cache")
		t.Fail()
	}
}
