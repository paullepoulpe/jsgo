// Port of com/google/javascript/jscomp/bundle/CachedTransformer.java @ closure-compiler v20261015

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

package jscomp

import "github.com/paullepoulpe/jsgo/internal/jutil"

// CachedTransformer is a Transformer that caches output from a delegate
// transformer.
type CachedTransformer struct {
	cache *jutil.LoadingCache[*Source, *Source]
}

var _ Source_Transformer = (*CachedTransformer)(nil)

const cachedTransformer_DEFAULT_CACHE_SPEC = "maximumSize=10000"

// Java: CachedTransformer(Source.Transformer, CacheBuilder<? super Source, ? super Source>)
func NewCachedTransformerWithBuilder(delegate Source_Transformer, builder *jutil.CacheBuilder) *CachedTransformer {
	return &CachedTransformer{cache: jutil.CacheBuilder_Build(builder, delegate.Transform)}
}

// Java: CachedTransformer(Source.Transformer, String)
func NewCachedTransformerWithSpec(delegate Source_Transformer, spec string) *CachedTransformer {
	return NewCachedTransformerWithBuilder(delegate, jutil.CacheBuilder_From(spec))
}

// Java: CachedTransformer(Source.Transformer)
func NewCachedTransformer(delegate Source_Transformer) *CachedTransformer {
	return NewCachedTransformerWithSpec(delegate, cachedTransformer_DEFAULT_CACHE_SPEC)
}

// Java: CachedTransformer.transform(Source)
func (t *CachedTransformer) Transform(input *Source) *Source {
	return t.cache.GetUnchecked(input)
}
