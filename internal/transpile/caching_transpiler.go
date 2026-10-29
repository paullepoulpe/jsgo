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

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

type cacheKey struct {
	path, code string
}

type cachingTranspilerNode struct {
	key        cacheKey
	result     *TranspileResult
	prev, next *cachingTranspilerNode
}

// CachingTranspiler is a transpiler implementation that delegates to a
// lower-level implementation on a cache miss.
type CachingTranspiler struct {
	delegate Transpiler
	maxSize  int // -1: unbounded
	mu       sync.Mutex
	cache    map[cacheKey]*cachingTranspilerNode
	head     *cachingTranspilerNode // least recently accessed
	tail     *cachingTranspilerNode // most recently accessed
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
		runtime:  sync.OnceValue(delegate.Runtime),
	}
	if maxSize != 0 {
		t.cache = make(map[cacheKey]*cachingTranspilerNode)
	}
	return t
}

// Java: CachingTranspiler.transpile(URI, String)
func (t *CachingTranspiler) Transpile(path string, code string) *TranspileResult {
	if t.maxSize == 0 {
		value := t.delegate.Transpile(path, code)
		if value == nil {
			panic(jutil.NewRuntimeException("CacheLoader returned null."))
		}
		return value
	}

	key := cacheKey{path: path, code: code}
	t.mu.Lock()
	if node, ok := t.cache[key]; ok {
		if t.maxSize > 0 {
			t.moveToTailLocked(node)
		}
		res := node.result
		t.mu.Unlock()
		return res
	}
	t.mu.Unlock()

	value := t.delegate.Transpile(path, code)
	if value == nil {
		panic(jutil.NewRuntimeException("CacheLoader returned null."))
	}

	t.mu.Lock()
	if node, ok := t.cache[key]; ok {
		if t.maxSize > 0 {
			t.moveToTailLocked(node)
		}
		res := node.result
		t.mu.Unlock()
		return res
	}
	node := &cachingTranspilerNode{key: key, result: value}
	t.cache[key] = node
	if t.maxSize > 0 {
		t.linkLastLocked(node)
		if len(t.cache) > t.maxSize {
			oldest := t.head
			t.unlinkLocked(oldest)
			delete(t.cache, oldest.key)
		}
	}
	t.mu.Unlock()
	return value
}

// Java: CachingTranspiler.runtime()
func (t *CachingTranspiler) Runtime() string {
	return t.runtime()
}

func (t *CachingTranspiler) unlinkLocked(n *cachingTranspilerNode) {
	if n.prev != nil {
		n.prev.next = n.next
	} else {
		t.head = n.next
	}
	if n.next != nil {
		n.next.prev = n.prev
	} else {
		t.tail = n.prev
	}
	n.prev, n.next = nil, nil
}

func (t *CachingTranspiler) linkLastLocked(n *cachingTranspilerNode) {
	n.prev = t.tail
	n.next = nil
	if t.tail != nil {
		t.tail.next = n
	} else {
		t.head = n
	}
	t.tail = n
}

func (t *CachingTranspiler) moveToTailLocked(n *cachingTranspilerNode) {
	if t.tail == n {
		return
	}
	t.unlinkLocked(n)
	t.linkLastLocked(n)
}
