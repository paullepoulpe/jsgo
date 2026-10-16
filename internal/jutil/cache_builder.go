package jutil

import (
	"strconv"
	"strings"
)

// Go-only: language-adapter — Guava CacheBuilder / CacheBuilderSpec /
// LoadingCache (LocalCache) for the size-bounded caches Closure builds.
//
// Entries are spread over Guava's segments (concurrencyLevel 4 by default,
// LocalCache.rehash of the key's hashCode) and each segment evicts its least
// recently accessed entry once it holds more than its share of maximumSize,
// which is how LocalCache evicts. Keys are compared with equals/hashCode.
// Time-based expiry, weak/soft references, weighers and statistics are not
// supported (asking for them panics).

// CacheBuilder is com.google.common.cache.CacheBuilder.
type CacheBuilder struct {
	initialCapacity  int
	concurrencyLevel int
	maximumSize      int64
	removalListener  func(key, value any)
}

// CacheBuilder_NewBuilder is CacheBuilder.newBuilder().
func CacheBuilder_NewBuilder() *CacheBuilder {
	return &CacheBuilder{initialCapacity: -1, concurrencyLevel: -1, maximumSize: -1}
}

// CacheBuilder_From is CacheBuilder.from(String spec) (CacheBuilderSpec).
func CacheBuilder_From(spec string) *CacheBuilder {
	b := CacheBuilder_NewBuilder()
	if spec == "" {
		return b
	}
	for _, keyValuePair := range strings.Split(spec, ",") {
		keyValuePair = strings.TrimSpace(keyValuePair)
		if keyValuePair == "" {
			continue
		}
		kv := strings.Split(keyValuePair, "=")
		for i := range kv {
			kv[i] = strings.TrimSpace(kv[i])
		}
		if len(kv) > 2 {
			panic(NewIllegalArgumentExceptionMsg("key-value pair " + keyValuePair + " with more than one equals sign"))
		}
		key := kv[0]
		var value *string
		if len(kv) == 2 {
			value = &kv[1]
		}
		parseLong := func() int64 {
			if value == nil || *value == "" {
				panic(NewIllegalArgumentExceptionMsg("value of key " + key + " omitted"))
			}
			v, err := strconv.ParseInt(*value, 10, 64)
			if err != nil {
				panic(NewIllegalArgumentExceptionMsg("key " + key + " value set to " + *value + ", must be integer"))
			}
			return v
		}
		switch key {
		case "initialCapacity":
			b.InitialCapacity(int(parseLong()))
		case "maximumSize":
			b.MaximumSize(parseLong())
		case "concurrencyLevel":
			b.ConcurrencyLevel(int(parseLong()))
		case "maximumWeight", "weakKeys", "softValues", "weakValues", "recordStats",
			"expireAfterAccess", "expireAfterWrite", "refreshAfterWrite", "refreshInterval":
			panic(NewUnsupportedOperationExceptionMsg("CacheBuilderSpec key not supported: " + key))
		default:
			panic(NewIllegalArgumentExceptionMsg("unknown key " + key))
		}
	}
	return b
}

// InitialCapacity is CacheBuilder.initialCapacity(int).
func (b *CacheBuilder) InitialCapacity(initialCapacity int) *CacheBuilder {
	CheckStateMsg1(b.initialCapacity == -1, "initial capacity was already set to %s", b.initialCapacity)
	CheckArgument(initialCapacity >= 0)
	b.initialCapacity = initialCapacity
	return b
}

// ConcurrencyLevel is CacheBuilder.concurrencyLevel(int).
func (b *CacheBuilder) ConcurrencyLevel(concurrencyLevel int) *CacheBuilder {
	CheckStateMsg1(b.concurrencyLevel == -1, "concurrency level was already set to %s", b.concurrencyLevel)
	CheckArgument(concurrencyLevel > 0)
	b.concurrencyLevel = concurrencyLevel
	return b
}

// MaximumSize is CacheBuilder.maximumSize(long).
func (b *CacheBuilder) MaximumSize(maximumSize int64) *CacheBuilder {
	CheckStateMsg1(b.maximumSize == -1, "maximum size was already set to %s", b.maximumSize)
	CheckArgumentMsg(maximumSize >= 0, "maximum size must not be negative")
	b.maximumSize = maximumSize
	return b
}

// RemovalListener is CacheBuilder.removalListener(RemovalListener); the
// listener receives the removed key and value.
func (b *CacheBuilder) RemovalListener(listener func(key, value any)) *CacheBuilder {
	b.removalListener = listener
	return b
}

// LoadingCache is com.google.common.cache.LoadingCache (LocalCache).
type LoadingCache[K any, V any] struct {
	loader       func(K) V
	listener     func(key, value any)
	segments     []*cacheSegment[K, V]
	segmentShift uint32
	segmentMask  int32
}

type cacheEntry[K any, V any] struct {
	key        K
	hash       int32
	value      V
	prev, next *cacheEntry[K, V] // access order
}

type cacheSegment[K any, V any] struct {
	maxSegmentWeight int64 // -1: unbounded
	table            map[int32][]*cacheEntry[K, V]
	count            int64
	head, tail       *cacheEntry[K, V] // head: least recently accessed
}

// CacheBuilder_Build is CacheBuilder.build(CacheLoader).
func CacheBuilder_Build[K any, V any](b *CacheBuilder, loader func(K) V) *LoadingCache[K, V] {
	concurrencyLevel := b.concurrencyLevel
	if concurrencyLevel == -1 {
		concurrencyLevel = 4
	}
	concurrencyLevel = min(concurrencyLevel, 1<<16)
	maxWeight := b.maximumSize
	evictsBySize := maxWeight >= 0

	segmentShift := 0
	segmentCount := 1
	for segmentCount < concurrencyLevel && (!evictsBySize || int64(segmentCount)*20 <= maxWeight) {
		segmentShift++
		segmentCount <<= 1
	}
	c := &LoadingCache[K, V]{
		loader:       loader,
		listener:     b.removalListener,
		segmentShift: uint32(32 - segmentShift),
		segmentMask:  int32(segmentCount - 1),
	}
	if evictsBySize {
		// Ensure sum of segment max weights = overall max weights
		maxSegmentWeight := maxWeight/int64(segmentCount) + 1
		remainder := maxWeight % int64(segmentCount)
		for i := 0; i < segmentCount; i++ {
			if int64(i) == remainder {
				maxSegmentWeight--
			}
			c.segments = append(c.segments, &cacheSegment[K, V]{maxSegmentWeight: maxSegmentWeight, table: map[int32][]*cacheEntry[K, V]{}})
		}
	} else {
		for i := 0; i < segmentCount; i++ {
			c.segments = append(c.segments, &cacheSegment[K, V]{maxSegmentWeight: -1, table: map[int32][]*cacheEntry[K, V]{}})
		}
	}
	return c
}

// cacheRehash is LocalCache.rehash(int).
func cacheRehash(h int32) int32 {
	h += (h << 15) ^ -12931 // 0xffffcd7d
	h ^= int32(uint32(h) >> 10)
	h += h << 3
	h ^= int32(uint32(h) >> 6)
	h += (h << 2) + (h << 14)
	return h ^ int32(uint32(h)>>16)
}

func (c *LoadingCache[K, V]) segmentFor(hash int32) *cacheSegment[K, V] {
	return c.segments[int32(uint32(hash)>>c.segmentShift)&c.segmentMask]
}

func (s *cacheSegment[K, V]) find(key K, hash int32) *cacheEntry[K, V] {
	for _, e := range s.table[hash] {
		if ObjectsEquals(e.key, key) {
			return e
		}
	}
	return nil
}

func (s *cacheSegment[K, V]) unlink(e *cacheEntry[K, V]) {
	if e.prev != nil {
		e.prev.next = e.next
	} else {
		s.head = e.next
	}
	if e.next != nil {
		e.next.prev = e.prev
	} else {
		s.tail = e.prev
	}
	e.prev, e.next = nil, nil
}

func (s *cacheSegment[K, V]) linkLast(e *cacheEntry[K, V]) {
	e.prev = s.tail
	e.next = nil
	if s.tail != nil {
		s.tail.next = e
	} else {
		s.head = e
	}
	s.tail = e
}

func (s *cacheSegment[K, V]) remove(e *cacheEntry[K, V]) {
	s.unlink(e)
	bucket := s.table[e.hash]
	for i, x := range bucket {
		if x == e {
			s.table[e.hash] = append(bucket[:i:i], bucket[i+1:]...)
			break
		}
	}
	if len(s.table[e.hash]) == 0 {
		delete(s.table, e.hash)
	}
	s.count--
}

// GetUnchecked is LoadingCache.getUnchecked(K): the cached value, or the
// loader's value (stored, possibly evicting the least recently used entry).
// Exceptions thrown by the loader are wrapped like Guava's.
func (c *LoadingCache[K, V]) GetUnchecked(key K) V {
	CheckNotNull(any(key))
	hash := cacheRehash(ObjectHashCode(key))
	s := c.segmentFor(hash)
	if e := s.find(key, hash); e != nil {
		// recordRead: move to the most recently accessed end.
		s.unlink(e)
		s.linkLast(e)
		return e.value
	}
	var value V
	func() {
		defer func() {
			if r := recover(); r != nil {
				t := AsThrowable(r)
				if err, ok := t.(IError); ok {
					panic(NewExecutionError(err))
				}
				panic(NewUncheckedExecutionException(t))
			}
		}()
		value = c.loader(key)
	}()
	if IsNil(any(value)) {
		panic(NewCacheLoader_InvalidCacheLoadException("CacheLoader returned null for key " + ToJavaString(key) + "."))
	}
	e := &cacheEntry[K, V]{key: key, hash: hash, value: value}
	s.table[hash] = append(s.table[hash], e)
	s.count++
	s.linkLast(e)
	c.evictEntries(s, e)
	return value
}

// evictEntries is LocalCache.Segment.evictEntries.
func (c *LoadingCache[K, V]) evictEntries(s *cacheSegment[K, V], newest *cacheEntry[K, V]) {
	if s.maxSegmentWeight < 0 {
		return
	}
	if 1 > s.maxSegmentWeight {
		s.remove(newest)
		c.notify(newest)
	}
	for s.count > s.maxSegmentWeight {
		e := s.head
		s.remove(e)
		c.notify(e)
	}
}

func (c *LoadingCache[K, V]) notify(e *cacheEntry[K, V]) {
	if c.listener != nil {
		c.listener(e.key, e.value)
	}
}

// Size is Cache.size().
func (c *LoadingCache[K, V]) Size() int64 {
	var n int64
	for _, s := range c.segments {
		n += s.count
	}
	return n
}

// GetIfPresent is Cache.getIfPresent(Object); ok reports presence.
func (c *LoadingCache[K, V]) GetIfPresent(key K) (V, bool) {
	hash := cacheRehash(ObjectHashCode(key))
	s := c.segmentFor(hash)
	if e := s.find(key, hash); e != nil {
		s.unlink(e)
		s.linkLast(e)
		return e.value, true
	}
	var zero V
	return zero, false
}
