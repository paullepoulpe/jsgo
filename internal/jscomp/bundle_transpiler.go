// Port of com/google/javascript/jscomp/bundle/Transpiler.java @ closure-compiler v20261021

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

// Transpiler is a source transformer for lowering JS language versions. May
// also include a runtime that needs to be shipped with the final bundle.
//
// (transpile.Transpiler, in the same Go package, is Transpile_Transpiler.)
//
// Java: class bundle.Transpiler extends CompilerBasedTransformer
type Transpiler struct {
	CompilerBasedTransformerBase

	runtimeLibraryName string
}

// Java: Transpiler(CompilerBasedTransformer.CompilerSupplier, String)
func NewTranspiler(compilerSupplier CompilerBasedTransformer_CompilerSupplier,
	runtimeLibraryName string) *Transpiler {
	t := &Transpiler{}
	t.InitCompilerBasedTransformerBase(t, compilerSupplier)
	t.runtimeLibraryName = runtimeLibraryName
	return t
}

// Java: Transpiler.compilerSupplier()
func Transpiler_CompilerSupplier() CompilerBasedTransformer_CompilerSupplier {
	return NewCompilerBasedTransformer_CompilerSupplier()
}

// GetRuntime returns Optional.of(runtimeLibraryName).
//
// Java: Transpiler.getRuntime()
func (t *Transpiler) GetRuntime() *string {
	s := t.runtimeLibraryName
	return &s
}

// Java: Transpiler.getTranformationName()
func (t *Transpiler) GetTranformationName() string {
	return "Transpilation"
}

// Java: Transpiler.ES5_TRANSPILER
var Transpiler_ES5_TRANSPILER = NewTranspiler(NewCompilerBasedTransformer_CompilerSupplier(), "es6_runtime")

// Transpiler_ToEs5 returns the recommended transpiler.
//
// Java: Transpiler.toEs5()
func Transpiler_ToEs5() *Transpiler {
	return transpiler_TO_ES5
}

var transpiler_TO_ES5 = NewTranspiler(NewCompilerBasedTransformer_CompilerSupplier(), "es6_runtime")
