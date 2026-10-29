// Port of com/google/javascript/jscomp/transpile/CachingTranspiler.java @ closure-compiler v20261021

/*
 * Copyright 2016 The Closure Compiler Authors.
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

package transpile

import (
	"sync"

	"github.com/paullepoulpe/jsgo/internal/bundle"
	"github.com/paullepoulpe/jsgo/internal/jutil"
)

type cachingTranspilerEntry struct {
	path       *jutil.URI
	code       string
	hash       int32
	value      *TranspileResult
	prev, next *cachingTranspilerEntry
}

// CachingTranspiler is a transpiler implementation that delegates to a
// lower-level implementation on a cache miss.
type CachingTranspiler struct {
	delegate Transpiler
	maxSize  int // -1: unbounded
	mu       sync.Mutex
	table    map[int32][]*cachingTranspilerEntry
	count    int
	head     *cachingTranspilerEntry // least recently accessed
	tail     *cachingTranspilerEntry // most recently accessed
	runtime  func() string
}

var _ Transpiler = (*CachingTranspiler)(nil)

// NewCachingTranspiler constructs a CachingTranspiler with the given maximum
// entry count (-1 for unbounded).
func NewCachingTranspiler(delegate Transpiler, maxSize int) *CachingTranspiler {
	jutil.CheckNotNull(delegate)
	t := &CachingTranspiler{
		delegate: delegate,
		maxSize:  maxSize,
		runtime:  bundle.Source_Suppliers_Memoize(delegate.Runtime),
	}
	if maxSize != 0 {
		t.table = make(map[int32][]*cachingTranspilerEntry)
	}
	return t
}

// Java: CachingTranspiler.transpile(URI, String)
func (t *CachingTranspiler) Transpile(path *jutil.URI, code string) *TranspileResult {
	jutil.CheckNotNull(path)
	if t.maxSize == 0 {
		value := t.delegate.Transpile(path, code)
		if value == nil {
			panic(jutil.NewRuntimeException("CacheLoader returned null."))
		}
		return value
	}

	hash := 31*(31+path.HashCode()) + jutil.StringHashCode(code)
	t.mu.Lock()
	if e := t.findLocked(path, code, hash); e != nil {
		t.moveToTailLocked(e)
		val := e.value
		t.mu.Unlock()
		return val
	}
	t.mu.Unlock()

	value := t.delegate.Transpile(path, code)
	if value == nil {
		panic(jutil.NewRuntimeException("CacheLoader returned null."))
	}

	t.mu.Lock()
	if e := t.findLocked(path, code, hash); e != nil {
		t.moveToTailLocked(e)
		val := e.value
		t.mu.Unlock()
		return val
	}
	e := &cachingTranspilerEntry{path: path, code: code, hash: hash, value: value}
	t.table[hash] = append(t.table[hash], e)
	t.count++
	t.linkLastLocked(e)
	if t.maxSize > 0 && t.count > t.maxSize {
		t.removeLocked(t.head)
	}
	t.mu.Unlock()
	return value
}

// Java: CachingTranspiler.runtime()
func (t *CachingTranspiler) Runtime() string {
	return t.runtime()
}

func (t *CachingTranspiler) findLocked(path *jutil.URI, code string, hash int32) *cachingTranspilerEntry {
	for _, e := range t.table[hash] {
		if e.code == code && (e.path == path || e.path.Equals(path)) {
			return e
		}
	}
	return nil
}

func (t *CachingTranspiler) unlinkLocked(e *cachingTranspilerEntry) {
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

func (t *CachingTranspiler) linkLastLocked(e *cachingTranspilerEntry) {
	e.prev = t.tail
	e.next = nil
	if t.tail != nil {
		t.tail.next = e
	} else {
		t.head = e
	}
	t.tail = e
}

func (t *CachingTranspiler) moveToTailLocked(e *cachingTranspilerEntry) {
	if t.tail == e {
		return
	}
	t.unlinkLocked(e)
	t.linkLastLocked(e)
}

func (t *CachingTranspiler) removeLocked(e *cachingTranspilerEntry) {
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
