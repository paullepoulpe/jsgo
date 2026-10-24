// Port of com/google/javascript/jscomp/transpile/TranspilerBuilder.java @ closure-compiler v20261021

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

// TranspilerBuilder is the basic Transpiler implementation for outputting
// ES5 code.
//
// Java: final class transpile.TranspilerBuilder
type TranspilerBuilder struct {
	transpiler Transpile_Transpiler
}

// TranspilerBuilder_ToEs5 returns a new TranspilerBuilder that transpiles
// down to ES5.
//
// Java: TranspilerBuilder.toEs5()
func TranspilerBuilder_ToEs5() *TranspilerBuilder {
	return transpilerBuilder_TO_ES5
}

var transpilerBuilder_TO_ES5 = NewTranspilerBuilder(
	NewBaseTranspiler(NewBaseTranspiler_CompilerSupplier(), "es6_runtime"))

// Java: TranspilerBuilder(Transpiler)
func NewTranspilerBuilder(transpiler Transpile_Transpiler) *TranspilerBuilder {
	b := &TranspilerBuilder{}
	b.transpiler = transpiler
	return b
}

// Caching returns a TranspilerBuilder with cached transpilations, using the
// default cache settings (maximum size of 10,000). Note that the builder
// itself is not changed.
//
// Java: TranspilerBuilder.caching()
func (b *TranspilerBuilder) Caching() *TranspilerBuilder {
	return b.CachingWithSpec(transpilerBuilder_DEFAULT_CACHE_SPEC)
}

// CachingWithSpec returns a TranspilerBuilder with cached transpilations,
// using the given cache spec. Note that the builder itself is not changed.
//
// Java: TranspilerBuilder.caching(String)
func (b *TranspilerBuilder) CachingWithSpec(spec string) *TranspilerBuilder {
	return b.CachingWithBuilder(jutil.CacheBuilder_From(spec))
}

// CachingWithBuilder returns a TranspilerBuilder with cached transpilations,
// using the given cache builder. Note that the builder itself is not
// changed.
//
// Java: TranspilerBuilder.caching(CacheBuilder<Object, ? super TranspileResult>)
func (b *TranspilerBuilder) CachingWithBuilder(builder *jutil.CacheBuilder) *TranspilerBuilder {
	return NewTranspilerBuilder(NewCachingTranspiler(b.transpiler, builder))
}

const transpilerBuilder_DEFAULT_CACHE_SPEC = "maximumSize=10000"

// Build returns the built Transpiler.
//
// Java: TranspilerBuilder.build()
func (b *TranspilerBuilder) Build() Transpile_Transpiler {
	return b.transpiler
}
