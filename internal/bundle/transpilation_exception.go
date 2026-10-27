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

package bundle

import (
	"strings"

	"github.com/paullepoulpe/jsgo/internal/jscomp"
	"github.com/paullepoulpe/jsgo/internal/jutil"
)

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

// NewTranspilationExceptionWithCause takes cause, an Exception.
//
// Java: TranspilationException(Exception)
func NewTranspilationExceptionWithCause(cause jutil.Throwable) *TranspilationException {
	return newTranspilationExceptionFromRoot(transpilationException_tryCastToTranspilationException(cause.GetCause()), cause)
}

// Java: TranspilationException(SourceExcerptProvider, ImmutableList<JSError>, ImmutableList<JSError>)
func NewTranspilationException(
	source jscomp.SourceExcerptProvider,
	errors []*jscomp.JSError,
	warnings []*jscomp.JSError) *TranspilationException {
	formatted := transpilationException_format(source, errors, warnings)
	return newTranspilationException(errors, warnings, &formatted, nil)
}

// Java: TranspilationException(TranspilationException, Exception)
func newTranspilationExceptionFromRoot(root *TranspilationException, cause jutil.Throwable) *TranspilationException {
	if root == nil {
		panic(jutil.NewNullPointerException())
	}
	// root's message is the formatted report, never null.
	message := root.GetMessage()
	return newTranspilationException(root.errors, root.warnings, &message, cause)
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
	for _, error := range errors {
		sb.WriteString("\n")
		sb.WriteString(formatter.FormatError(error))
	}
	for _, warning := range warnings {
		sb.WriteString("\n")
		sb.WriteString(formatter.FormatError(warning))
	}
	return sb.String()
}

// transpilationException_tryCastToTranspilationException may return nil.
//
// Java: TranspilationException.tryCastToTranspilationException(Throwable)
func transpilationException_tryCastToTranspilationException(t jutil.Throwable) *TranspilationException {
	if transpilationException, ok := t.(*TranspilationException); ok {
		return transpilationException
	}
	return nil
}
