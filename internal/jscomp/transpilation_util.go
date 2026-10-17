// Port of com/google/javascript/jscomp/TranspilationUtil.java @ closure-compiler v20261015

/*
 * Copyright 2014 The Closure Compiler Authors.
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

import (
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// TranspilationUtil holds util functions for transpilation passes.
//
// Java: TranspilationUtil() is private (prevents instantiation); the class
// has only static members.
type TranspilationUtil struct{}

// Java: TranspilationUtil.CANNOT_CONVERT
var TranspilationUtil_CANNOT_CONVERT = DiagnosticType_Error("JSC_CANNOT_CONVERT", "This code cannot be transpiled. {0}")

// TODO(tbreisacher): Remove this once we have implemented transpilation for all the features
// we intend to support.
//
// Java: TranspilationUtil.CANNOT_CONVERT_YET
var TranspilationUtil_CANNOT_CONVERT_YET = DiagnosticType_Error(
	"JSC_CANNOT_CONVERT_YET", "Transpilation of ''{0}'' is not yet implemented.")

// Java: TranspilationUtil.cannotConvert(AbstractCompiler, Node, String)
func TranspilationUtil_CannotConvert(compiler AbstractCompiler, n *rhino.Node, message string) {
	compiler.Report(JSError_MakeWithNode(n, TranspilationUtil_CANNOT_CONVERT, message))
}

// TranspilationUtil_CannotConvertYet warns the user that the given feature
// cannot be transpiled because the transpilation is not yet implemented. A
// call to this method is essentially a "TODO(tbreisacher): Implement
// feature" comment.
//
// Java: TranspilationUtil.cannotConvertYet(AbstractCompiler, Node, String)
func TranspilationUtil_CannotConvertYet(compiler AbstractCompiler, n *rhino.Node, feature string) {
	compiler.Report(JSError_MakeWithNode(n, TranspilationUtil_CANNOT_CONVERT_YET, feature))
}

// Java: TranspilationUtil.preloadTranspilationRuntimeFunction(AbstractCompiler, String)
func TranspilationUtil_PreloadTranspilationRuntimeFunction(compiler AbstractCompiler, function string) {
	compiler.
		GetRuntimeJsLibManager().
		EnsureLibraryInjected("es6/util/"+jutil.StringToLowerCaseRoot(function), false)
}
