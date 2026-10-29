// Port of com/google/javascript/jscomp/transpile/Transpiler.java @ closure-compiler v20261021
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

// Package transpile implements the standalone per-file ES transpiler API
// (Transpiler, BaseTranspiler, CachingTranspiler, TranspileResult) and
// ClosureBundler.
package transpile

import (
	"github.com/paullepoulpe/jsgo/internal/bundle"
	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// Transpiler is the common interface for a transpiler.
//
// There are a number of considerations when implementing this interface.
// Specifically,
//  1. Which compiler and options to use, including language mode and any
//     specific Compiler subclass. This is handled by BaseTranspiler
//     accepting a CompilerSupplier.
//  2. Whether or not to generate external or embedded source maps. This is
//     handled by returning a TranspileResult, which is able to return any
//     combination of embedded or external source map.
//  3. Specification of a sourceURL, handling of goog.module, or other
//     postprocessing, such as wrapping the script in eval. This is left to
//     other collaborators.
//  4. Caching. This is handled by a CachingTranspiler that decorates an
//     existing Transpiler with caching behavior.
type Transpiler interface {
	// Transpile transforms the given chunk of code. The input should be an
	// entire file worth of code.
	Transpile(path *jutil.URI, code string) *TranspileResult

	// Runtime returns any necessary runtime code as a string. This should
	// include everything that could possibly be required at runtime,
	// regardless of whether it's actually used by any of the code that will
	// be transpiled.
	Runtime() string
}

type transpile_Transpiler_null struct{}

func (transpile_Transpiler_null) Transpile(path *jutil.URI, code string) *TranspileResult {
	return NewTranspileResult(path, code, code, "")
}

func (transpile_Transpiler_null) Runtime() string { return "" }

// Transpiler_NULL is the null implementation that does no transpilation at all.
//
// Java: Transpiler.NULL
var Transpiler_NULL Transpiler = transpile_Transpiler_null{}

// TranspilerBuilder is the basic Transpiler implementation for outputting
// ES5 code.
//
// Java: final class transpile.TranspilerBuilder
type TranspilerBuilder struct {
	transpiler Transpiler
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
func NewTranspilerBuilder(transpiler Transpiler) *TranspilerBuilder {
	return &TranspilerBuilder{transpiler: transpiler}
}

const transpilerBuilder_DEFAULT_MAX_SIZE = 10000

// Caching returns a TranspilerBuilder with cached transpilations, using the
// default cache settings (maximum size of 10,000). Note that the builder
// itself is not changed.
//
// Java: TranspilerBuilder.caching()
func (b *TranspilerBuilder) Caching() *TranspilerBuilder {
	return NewTranspilerBuilder(NewCachingTranspiler(b.transpiler, transpilerBuilder_DEFAULT_MAX_SIZE))
}

// CachingWithSpec returns a TranspilerBuilder with cached transpilations,
// using the given cache spec. Note that the builder itself is not changed.
//
// Java: TranspilerBuilder.caching(String)
func (b *TranspilerBuilder) CachingWithSpec(spec string) *TranspilerBuilder {
	return NewTranspilerBuilder(NewCachingTranspiler(b.transpiler, bundle.ParseCacheMaxSize(spec)))
}

// Build returns the built Transpiler.
//
// Java: TranspilerBuilder.build()
func (b *TranspilerBuilder) Build() Transpiler {
	return b.transpiler
}
