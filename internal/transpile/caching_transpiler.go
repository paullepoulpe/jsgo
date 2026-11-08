// Port of com/google/javascript/jscomp/transpile/CachingTranspiler.java @ closure-compiler v20261104

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

// CachingTranspiler is a transpiler implementation that delegates to a
// lower-level implementation on a cache miss.
type CachingTranspiler struct {
	delegate Transpiler
	cache    jutil.LRUCache[cacheKey, *TranspileResult]
	runtime  func() string
}

var _ Transpiler = (*CachingTranspiler)(nil)

// NewCachingTranspiler constructs a CachingTranspiler with the given maximum
// entry count (-1 for unbounded).
func NewCachingTranspiler(delegate Transpiler, maxSize int) *CachingTranspiler {
	jutil.CheckNotNull(delegate)
	t := &CachingTranspiler{
		delegate: delegate,
		runtime:  sync.OnceValue(delegate.Runtime),
	}
	t.cache.Init(maxSize)
	return t
}

// Java: CachingTranspiler.transpile(URI, String)
func (t *CachingTranspiler) Transpile(path string, code string) *TranspileResult {
	return t.cache.GetOrLoad(cacheKey{path: path, code: code}, func() *TranspileResult {
		value := t.delegate.Transpile(path, code)
		if value == nil {
			panic(jutil.NewRuntimeException("CacheLoader returned null."))
		}
		return value
	})
}

// Java: CachingTranspiler.runtime()
func (t *CachingTranspiler) Runtime() string {
	return t.runtime()
}
