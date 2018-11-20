// LRU cache with a global ttl. Each cache is created with a ttl which is updated during Put.
// Every ttl/4 interval cache iterates over all entries are removes those whose ttl has expired.
// Based code copied from https://github.com/golang/build/blob/master/internal/lru/cache.go

package cache

import (
	"container/list"
	"goprizm/goprof"
	"sync"
	"time"
)

// LRU is an LRU cache, safe for concurrent access.
type LRU struct {
	maxEntries int

	mu    sync.Mutex
	ll    *list.List
	cache map[interface{}]*list.Element
}

// *entry is the type stored in each *list.Element.
type entry struct {
	key, value interface{}
	updatedAt  time.Time
}

// New returns a new cache with the provided maximum items.
// to disable cleanup set ttl to 0
func NewLRU(maxEntries int, ttl time.Duration) *LRU {
	lru := &LRU{
		maxEntries: maxEntries,
		ll:         list.New(),
		cache:      make(map[interface{}]*list.Element),
	}

	// run cleanup go routine if ttl is greater than 0
	if ttl != 0 {
		go lru.cleanup(ttl)
	}

	return lru
}

func (lru *LRU) cleanup(ttl time.Duration) {
	ticker := time.NewTicker(ttl / 4)
	for t := range ticker.C {
		limit := t.Add(-ttl)
		lru.mu.Lock()
		for key, v := range lru.cache {
			if v.Value.(entry).updatedAt.Before(limit) {
				lru.ll.Remove(v)
				delete(lru.cache, key)
			}
		}
		lru.mu.Unlock()
	}
}

// Add adds the provided key and value to the cache, evicting
// an old item if necessary.
func (lru *LRU) Add(key, value interface{}) {
	lru.mu.Lock()
	defer lru.mu.Unlock()

	// Already in cache?
	if ee, ok := lru.cache[key]; ok {
		lru.ll.MoveToFront(ee)
		ee.Value = entry{key, value, time.Now()}
		return
	}

	// Add to cache if not present
	ele := lru.ll.PushFront(entry{key, value, time.Now()})
	lru.cache[key] = ele

	if lru.ll.Len() > lru.maxEntries {
		lru.removeOldest()
	}
}

// Get fetches the key's value from the cache.
// Get moves the key to front of list.
// The ok result will be true if the item was found.
func (lru *LRU) Get(key interface{}) (value interface{}, ok bool) {
	lru.mu.Lock()
	defer lru.mu.Unlock()
	if ele, hit := lru.cache[key]; hit {
		lru.ll.MoveToFront(ele)
		return ele.Value.(entry).value, true
	}
	return
}

// Peek fetches the key's value from the cache.
// The ok result will be true if the item was found.
func (lru *LRU) Peek(key interface{}) (value interface{}, ok bool) {
	lru.mu.Lock()
	defer lru.mu.Unlock()
	if ele, hit := lru.cache[key]; hit {
		return ele.Value.(entry).value, true
	}
	return
}

// note: must hold lru.mu
func (lru *LRU) removeOldest() (key, value interface{}) {
	ele := lru.ll.Back()
	if ele == nil {
		return
	}
	lru.ll.Remove(ele)
	ent := ele.Value.(entry)
	delete(lru.cache, ent.key)
	return ent.key, ent.value

}

// Remove removes a key from the cache. The return value is the value that was
// removed or nil if the key was not present. The value's EvictionNotifier is
// not run by Remove.
func (lru *LRU) Remove(key interface{}) interface{} {
	lru.mu.Lock()
	defer lru.mu.Unlock()
	if ele, found := lru.cache[key]; found {
		_, value := lru.remove(ele)
		return value
	}
	return nil
}

// note: must hold c.mu
func (lru *LRU) remove(ele *list.Element) (key, value interface{}) {
	lru.ll.Remove(ele)
	ent := ele.Value.(entry)
	delete(lru.cache, ent.key)
	return ent.key, ent.value
}

// Len returns the number of items in the cache.
func (lru *LRU) Len() int {
	lru.mu.Lock()
	defer lru.mu.Unlock()
	return lru.ll.Len()
}

func (lru *LRU) Size() (size uint64) {
	it := lru.NewIterator()
	for {
		key, val, ok := it.GetAndAdvance()
		if !ok {
			break
		}

		size = size + goprof.Sizeof(key)
		size = size + goprof.Sizeof(val)
	}

	return size
}

// Iterator is an iterator through the list. The iterator points to a nil
// element at the end of the list.
type Iterator struct {
	lru     *LRU
	this    *list.Element
	forward bool
}

// NewIterator returns a new iterator for the LRU.
func (lru *LRU) NewIterator() *Iterator {
	lru.mu.Lock()
	defer lru.mu.Unlock()
	return &Iterator{lru: lru, this: lru.ll.Front(), forward: true}
}

// NewReverseIterator returns a new reverse iterator for the LRU.
func (lru *LRU) NewReverseIterator() *Iterator {
	lru.mu.Lock()
	defer lru.mu.Unlock()
	return &Iterator{lru: lru, this: lru.ll.Back(), forward: false}
}

// GetAndAdvance returns key, value, true if the current entry is valid and advances the
// iterator. Otherwise it returns nil, nil, false.
func (i *Iterator) GetAndAdvance() (interface{}, interface{}, bool) {
	i.lru.mu.Lock()
	defer i.lru.mu.Unlock()
	if i.this == nil {
		return nil, nil, false
	}
	ent := i.this.Value.(entry)
	if i.forward {
		i.this = i.this.Next()
	} else {
		i.this = i.this.Prev()
	}
	return ent.key, ent.value, true
}
