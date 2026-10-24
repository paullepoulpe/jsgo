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

package jscomp

import "github.com/paullepoulpe/jsgo/internal/jutil"

// CachingTranspiler is a transpiler implementation that delegates to a
// lower-level implementation on a cache miss. Passed a CacheBuilder to
// allow specifying maximum size and other requirements externally.
type CachingTranspiler struct {
	cache   *jutil.LoadingCache[*cachingTranspiler_Key, *TranspileResult]
	runtime func() string
}

var _ Transpile_Transpiler = (*CachingTranspiler)(nil)

// Java: CachingTranspiler(Transpiler, CacheBuilder<Object, ? super TranspileResult>)
func NewCachingTranspiler(
	delegate Transpile_Transpiler, builder *jutil.CacheBuilder) *CachingTranspiler {
	jutil.CheckNotNull(delegate)
	t := &CachingTranspiler{}
	t.cache =
		jutil.CacheBuilder_Build(builder,
			func(key *cachingTranspiler_Key) *TranspileResult {
				return delegate.Transpile(key.path, key.code)
			})
	t.runtime = source_Suppliers_Memoize(delegate.Runtime)
	return t
}

// Java: CachingTranspiler.transpile(URI, String)
func (t *CachingTranspiler) Transpile(path *jutil.URI, code string) (result *TranspileResult) {
	jutil.Try(func() {
		result = t.cache.GetUnchecked(newCachingTranspiler_Key(path, code))
	}, func(e jutil.IUncheckedExecutionException) {
		if _, ok := e.GetCause().(ITranspilationException); ok {
			// If transpilation fails due to a parse error we can get an UncheckedExecutionException.
			// This is because BaseTranspiler wraps the parse error as an TranspilationException.
			// TODO(joeltine): This might better as a checked exception?
			panic(NewTranspilationExceptionWithCause(e))
		} else {
			panic(e)
		}
	})
	return result
}

// Java: CachingTranspiler.runtime()
func (t *CachingTranspiler) Runtime() string {
	return t.runtime()
}

// cachingTranspiler_Key is CachingTranspiler.Key.
type cachingTranspiler_Key struct {
	path *jutil.URI
	code string
}

// Java: CachingTranspiler.Key(URI, String)
func newCachingTranspiler_Key(path *jutil.URI, code string) *cachingTranspiler_Key {
	return &cachingTranspiler_Key{path: jutil.CheckNotNull(path), code: code}
}

// Java: CachingTranspiler.Key.equals(Object)
func (k *cachingTranspiler_Key) Equals(other any) bool {
	key, ok := other.(*cachingTranspiler_Key)
	return ok && key != nil && key.path.Equals(k.path) && key.code == k.code
}

// Java: CachingTranspiler.Key.hashCode()
func (k *cachingTranspiler_Key) HashCode() int32 {
	// Objects.hash(path, code)
	h := int32(1)
	h = 31*h + k.path.HashCode()
	h = 31*h + jutil.StringHashCode(k.code)
	return h
}
