// Port of com/google/javascript/jscomp/transpile/Transpiler.java @ closure-compiler v20261021

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

// Transpile_Transpiler is the common interface for a transpiler.
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
//
// Go-only: language-adapter — transpile.Transpiler and bundle.Transpiler
// are both in the Go package jscomp, so this one carries its subpackage
// name.
type Transpile_Transpiler interface {
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

// Transpile_Transpiler_NULL is the null implementation that does no
// transpilation at all.
//
// Java: Transpiler.NULL
var Transpile_Transpiler_NULL Transpile_Transpiler = transpile_Transpiler_null{}
