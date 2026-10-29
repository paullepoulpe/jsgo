// Port of com/google/javascript/jscomp/bundle/Transpiler.java @ closure-compiler v20261021
// Port of com/google/javascript/jscomp/bundle/TranspilationException.java @ closure-compiler v20261021

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

// Package bundle implements per-source compilation, transpilation, caching, and
// coverage instrumentation transformers (com.google.javascript.jscomp.bundle).
package bundle

import (
	"strings"

	"github.com/paullepoulpe/jsgo/internal/jscomp"
	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// Transpiler is a source transformer for lowering JS language versions. May
// also include a runtime that needs to be shipped with the final bundle.
//
// Java: class bundle.Transpiler extends CompilerBasedTransformer
type Transpiler struct {
	CompilerBasedTransformer
}

var _ Source_Transformer = (*Transpiler)(nil)

// Java: Transpiler(CompilerBasedTransformer.CompilerSupplier, String)
func NewTranspiler(compilerSupplier CompilerBasedTransformer_CompilerSupplier,
	runtimeLibraryName string) *Transpiler {
	return &Transpiler{
		CompilerBasedTransformer: CompilerBasedTransformer{
			compilerSupplier:   jutil.CheckNotNull(compilerSupplier),
			transformationName: "Transpilation",
			runtimeLibraryName: runtimeLibraryName,
			hasRuntime:         true,
		},
	}
}

// Java: Transpiler.compilerSupplier()
func Transpiler_CompilerSupplier() CompilerBasedTransformer_CompilerSupplier {
	return NewCompilerBasedTransformer_CompilerSupplier()
}

// TranspilationException is an unchecked exception thrown when
// transpilation fails due to one or more errors in the input script.
type TranspilationException struct {
	jutil.RuntimeException
	errors   []*jscomp.JSError
	warnings []*jscomp.JSError
}

func (*TranspilationException) transpilationException() {}

// ITranspilationException matches TranspilationException.
type ITranspilationException interface {
	jutil.IRuntimeException
	transpilationException()
}

// Java: TranspilationException(SourceExcerptProvider, ImmutableList<JSError>, ImmutableList<JSError>)
func NewTranspilationException(
	source jscomp.SourceExcerptProvider,
	errors []*jscomp.JSError,
	warnings []*jscomp.JSError) *TranspilationException {
	formatted := transpilationException_format(source, errors, warnings)
	return newTranspilationException(errors, warnings, &formatted, nil)
}

// Java: TranspilationException(ImmutableList<JSError>, ImmutableList<JSError>, String, Exception)
func newTranspilationException(
	errors []*jscomp.JSError,
	warnings []*jscomp.JSError,
	formatted *string,
	cause jutil.Throwable) *TranspilationException {
	e := &TranspilationException{errors: errors, warnings: warnings}
	jutil.InitThrowable(e, "com.google.javascript.jscomp.bundle.TranspilationException", formatted, cause)
	return e
}

// Java: TranspilationException.errors()
func (e *TranspilationException) Errors() []*jscomp.JSError {
	return e.errors
}

// Java: TranspilationException.warnings()
func (e *TranspilationException) Warnings() []*jscomp.JSError {
	return e.warnings
}

// Java: TranspilationException.format(SourceExcerptProvider, ImmutableList<JSError>, ImmutableList<JSError>)
func transpilationException_format(
	source jscomp.SourceExcerptProvider,
	errors []*jscomp.JSError,
	warnings []*jscomp.JSError) string {
	var sb strings.Builder
	sb.WriteString("Transpilation failed:\n")
	var formatter jscomp.MessageFormatter
	if source != nil {
		formatter = jscomp.ErrorFormat_SINGLELINE.ToFormatter(source, false)
	} else {
		formatter = jscomp.ErrorFormat_SOURCELESS.ToFormatter(source, false)
	}
	for _, err := range errors {
		sb.WriteByte('\n')
		sb.WriteString(formatter.FormatError(err))
	}
	for _, warning := range warnings {
		sb.WriteByte('\n')
		sb.WriteString(formatter.FormatError(warning))
	}
	return sb.String()
}
