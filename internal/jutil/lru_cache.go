package jutil

import "sync"

type lruEntry[K any, V any] struct {
	key        K
	value      V
	hash       int32
	prev, next *lruEntry[K, V]
}

type lruList[K any, V any] struct {
	head *lruEntry[K, V] // least recently accessed
	tail *lruEntry[K, V] // most recently accessed
}

func (l *lruList[K, V]) linkLast(e *lruEntry[K, V]) {
	e.prev = l.tail
	e.next = nil
	if l.tail != nil {
		l.tail.next = e
	} else {
		l.head = e
	}
	l.tail = e
}

func (l *lruList[K, V]) unlink(e *lruEntry[K, V]) {
	if e.prev != nil {
		e.prev.next = e.next
	} else {
		l.head = e.next
	}
	if e.next != nil {
		e.next.prev = e.prev
	} else {
		l.tail = e.prev
	}
	e.prev, e.next = nil, nil
}

func (l *lruList[K, V]) moveToTail(e *lruEntry[K, V]) {
	if l.tail == e {
		return
	}
	l.unlink(e)
	l.linkLast(e)
}

// LRUCache is a concurrency-safe cache keyed by Go == with optional LRU eviction
// (maxSize > 0 bounded, maxSize < 0 unbounded, maxSize == 0 disabled).
type LRUCache[K comparable, V any] struct {
	maxSize int
	mu      sync.Mutex
	entries map[K]*lruEntry[K, V]
	order   lruList[K, V]
}

// Init initializes c in place with the given maximum entry count.
func (c *LRUCache[K, V]) Init(maxSize int) {
	c.maxSize = maxSize
	if maxSize != 0 {
		c.entries = make(map[K]*lruEntry[K, V])
	}
}

// GetOrLoad returns the cached value for key or computes it via load outside
// the cache lock, promoting hits and evicting the least recently used entry
// when maxSize > 0 is exceeded.
func (c *LRUCache[K, V]) GetOrLoad(key K, load func() V) V {
	if c.maxSize == 0 {
		return load()
	}

	c.mu.Lock()
	if e, ok := c.entries[key]; ok {
		if c.maxSize > 0 {
			c.order.moveToTail(e)
		}
		v := e.value
		c.mu.Unlock()
		return v
	}
	c.mu.Unlock()

	value := load()

	c.mu.Lock()
	if e, ok := c.entries[key]; ok {
		if c.maxSize > 0 {
			c.order.moveToTail(e)
		}
		v := e.value
		c.mu.Unlock()
		return v
	}
	e := &lruEntry[K, V]{key: key, value: value}
	c.entries[key] = e
	if c.maxSize > 0 {
		c.order.linkLast(e)
		if len(c.entries) > c.maxSize {
			oldest := c.order.head
			c.order.unlink(oldest)
			delete(c.entries, oldest.key)
		}
	}
	c.mu.Unlock()
	return value
}

// HashableKey is implemented by comparable reference keys that override Java
// equals(Object) and hashCode().
type HashableKey interface {
	comparable
	hashable
}

// ObjectLRUCache is a concurrency-safe LRU cache keyed by Java HashCode/Equals.
type ObjectLRUCache[K HashableKey, V any] struct {
	maxSize int
	mu      sync.Mutex
	table   map[int32][]*lruEntry[K, V]
	count   int
	order   lruList[K, V]
}

// Init initializes c in place with the given maximum entry count
// (-1 for unbounded, 0 for disabled).
func (c *ObjectLRUCache[K, V]) Init(maxSize int) {
	c.maxSize = maxSize
	if maxSize != 0 {
		c.table = make(map[int32][]*lruEntry[K, V])
	}
}

// GetOrLoad returns the cached value for key or computes it via load outside
// the cache lock.
func (c *ObjectLRUCache[K, V]) GetOrLoad(key K, load func() V) V {
	if c.maxSize == 0 {
		return load()
	}

	hash := key.HashCode()
	c.mu.Lock()
	if e := c.findLocked(key, hash); e != nil {
		if c.maxSize > 0 {
			c.order.moveToTail(e)
		}
		v := e.value
		c.mu.Unlock()
		return v
	}
	c.mu.Unlock()

	value := load()

	c.mu.Lock()
	if e := c.findLocked(key, hash); e != nil {
		if c.maxSize > 0 {
			c.order.moveToTail(e)
		}
		v := e.value
		c.mu.Unlock()
		return v
	}
	e := &lruEntry[K, V]{key: key, value: value, hash: hash}
	c.table[hash] = append(c.table[hash], e)
	c.count++
	if c.maxSize > 0 {
		c.order.linkLast(e)
		if c.count > c.maxSize {
			c.evictOldestLocked()
		}
	}
	c.mu.Unlock()
	return value
}

func (c *ObjectLRUCache[K, V]) findLocked(key K, hash int32) *lruEntry[K, V] {
	for _, e := range c.table[hash] {
		if e.key == key || e.key.Equals(key) {
			return e
		}
	}
	return nil
}

func (c *ObjectLRUCache[K, V]) evictOldestLocked() {
	oldest := c.order.head
	c.order.unlink(oldest)
	bucket := c.table[oldest.hash]
	for i, x := range bucket {
		if x == oldest {
			c.table[oldest.hash] = append(bucket[:i:i], bucket[i+1:]...)
			break
		}
	}
	if len(c.table[oldest.hash]) == 0 {
		delete(c.table, oldest.hash)
	}
	c.count--
}
