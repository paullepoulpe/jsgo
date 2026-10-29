// Port of com/google/javascript/jscomp/bundle/CachedTransformer.java @ closure-compiler v20261021

/*
 * Copyright 2017 The Closure Compiler Authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package bundle

import (
	"strconv"
	"strings"
	"sync"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

type cachedSourceEntry struct {
	key        *Source
	hash       int32
	value      *Source
	prev, next *cachedSourceEntry
}

// CachedTransformer is a Transformer that caches output from a delegate
// transformer.
type CachedTransformer struct {
	delegate Source_Transformer
	maxSize  int // -1: unbounded
	mu       sync.Mutex
	table    map[int32][]*cachedSourceEntry
	count    int
	head     *cachedSourceEntry // least recently accessed
	tail     *cachedSourceEntry // most recently accessed
}

var _ Source_Transformer = (*CachedTransformer)(nil)

const cachedTransformer_DEFAULT_MAX_SIZE = 10000

// ParseCacheMaxSize parses a Guava CacheBuilderSpec string and returns the
// configured maximumSize, or -1 if unbounded.
func ParseCacheMaxSize(spec string) int {
	maxSize := int64(-1)
	if spec == "" {
		return -1
	}
	initialCapacity := int64(-1)
	concurrencyLevel := int64(-1)
	for rest := spec; rest != ""; {
		var keyValuePair string
		keyValuePair, rest, _ = strings.Cut(rest, ",")
		keyValuePair = strings.TrimSpace(keyValuePair)
		if keyValuePair == "" {
			continue
		}
		rawKey, rawVal, hasEq := strings.Cut(keyValuePair, "=")
		if hasEq && strings.ContainsRune(rawVal, '=') {
			panic(jutil.NewIllegalArgumentExceptionMsg("key-value pair " + keyValuePair + " with more than one equals sign"))
		}
		key := strings.TrimSpace(rawKey)
		value := strings.TrimSpace(rawVal)
		parseLong := func() int64 {
			if !hasEq || value == "" {
				panic(jutil.NewIllegalArgumentExceptionMsg("value of key " + key + " omitted"))
			}
			v, err := strconv.ParseInt(value, 10, 64)
			if err != nil {
				panic(jutil.NewIllegalArgumentExceptionMsg("key " + key + " value set to " + value + ", must be integer"))
			}
			return v
		}
		switch key {
		case "initialCapacity":
			jutil.CheckStateMsg1(initialCapacity == -1, "initial capacity was already set to %s", int(initialCapacity))
			v := parseLong()
			jutil.CheckArgument(v >= 0)
			initialCapacity = v
		case "maximumSize":
			jutil.CheckStateMsg1(maxSize == -1, "maximum size was already set to %s", maxSize)
			v := parseLong()
			jutil.CheckArgumentMsg(v >= 0, "maximum size must not be negative")
			maxSize = v
		case "concurrencyLevel":
			jutil.CheckStateMsg1(concurrencyLevel == -1, "concurrency level was already set to %s", int(concurrencyLevel))
			v := parseLong()
			jutil.CheckArgument(v > 0)
			concurrencyLevel = v
		case "maximumWeight", "weakKeys", "softValues", "weakValues", "recordStats",
			"expireAfterAccess", "expireAfterWrite", "refreshAfterWrite", "refreshInterval":
			panic(jutil.NewUnsupportedOperationExceptionMsg("CacheBuilderSpec key not supported: " + key))
		default:
			panic(jutil.NewIllegalArgumentExceptionMsg("unknown key " + key))
		}
	}
	return int(maxSize)
}

// NewCachedTransformerWithMaxSize constructs a CachedTransformer with the
// given maximum entry count (-1 for unbounded).
func NewCachedTransformerWithMaxSize(delegate Source_Transformer, maxSize int) *CachedTransformer {
	jutil.CheckNotNull(delegate)
	t := &CachedTransformer{
		delegate: delegate,
		maxSize:  maxSize,
	}
	if maxSize != 0 {
		t.table = make(map[int32][]*cachedSourceEntry)
	}
	return t
}

// Java: CachedTransformer(Source.Transformer, String)
func NewCachedTransformerWithSpec(delegate Source_Transformer, spec string) *CachedTransformer {
	return NewCachedTransformerWithMaxSize(delegate, ParseCacheMaxSize(spec))
}

// Java: CachedTransformer(Source.Transformer)
func NewCachedTransformer(delegate Source_Transformer) *CachedTransformer {
	return NewCachedTransformerWithMaxSize(delegate, cachedTransformer_DEFAULT_MAX_SIZE)
}

// Java: CachedTransformer.transform(Source)
func (t *CachedTransformer) Transform(input *Source) *Source {
	jutil.CheckNotNull(input)
	if t.maxSize == 0 {
		value := t.delegate.Transform(input)
		if value == nil {
			panic(jutil.NewRuntimeException("CacheLoader returned null."))
		}
		return value
	}

	hash := input.HashCode()
	t.mu.Lock()
	if e := t.findLocked(input, hash); e != nil {
		t.moveToTailLocked(e)
		val := e.value
		t.mu.Unlock()
		return val
	}
	t.mu.Unlock()

	value := t.delegate.Transform(input)
	if value == nil {
		panic(jutil.NewRuntimeException("CacheLoader returned null."))
	}

	t.mu.Lock()
	if e := t.findLocked(input, hash); e != nil {
		t.moveToTailLocked(e)
		val := e.value
		t.mu.Unlock()
		return val
	}
	e := &cachedSourceEntry{key: input, hash: hash, value: value}
	t.table[hash] = append(t.table[hash], e)
	t.count++
	t.linkLastLocked(e)
	if t.maxSize > 0 && t.count > t.maxSize {
		t.removeLocked(t.head)
	}
	t.mu.Unlock()
	return value
}

func (t *CachedTransformer) findLocked(key *Source, hash int32) *cachedSourceEntry {
	for _, e := range t.table[hash] {
		if e.key == key || e.key.Equals(key) {
			return e
		}
	}
	return nil
}

func (t *CachedTransformer) unlinkLocked(e *cachedSourceEntry) {
	if e.prev != nil {
		e.prev.next = e.next
	} else {
		t.head = e.next
	}
	if e.next != nil {
		e.next.prev = e.prev
	} else {
		t.tail = e.prev
	}
	e.prev, e.next = nil, nil
}

func (t *CachedTransformer) linkLastLocked(e *cachedSourceEntry) {
	e.prev = t.tail
	e.next = nil
	if t.tail != nil {
		t.tail.next = e
	} else {
		t.head = e
	}
	t.tail = e
}

func (t *CachedTransformer) moveToTailLocked(e *cachedSourceEntry) {
	if t.tail == e {
		return
	}
	t.unlinkLocked(e)
	t.linkLastLocked(e)
}

func (t *CachedTransformer) removeLocked(e *cachedSourceEntry) {
	t.unlinkLocked(e)
	bucket := t.table[e.hash]
	for i, x := range bucket {
		if x == e {
			t.table[e.hash] = append(bucket[:i:i], bucket[i+1:]...)
			break
		}
	}
	if len(t.table[e.hash]) == 0 {
		delete(t.table, e.hash)
	}
	t.count--
}
