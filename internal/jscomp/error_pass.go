// Port of com/google/javascript/jscomp/ErrorPass.java @ closure-compiler v20260415

/*
 * Copyright 2010 The Closure Compiler Authors.
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
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// ErrorPass is a compiler pass that just reports an error.
type ErrorPass struct {
	compiler AbstractCompiler
	error_   *JSError
}

// Java: ErrorPass(AbstractCompiler, DiagnosticType)
func NewErrorPass(compiler AbstractCompiler, error_ *DiagnosticType) *ErrorPass {
	return NewErrorPassWithJSError(compiler, JSError_Make(error_))
}

// Java: ErrorPass(AbstractCompiler, JSError)
func NewErrorPassWithJSError(compiler AbstractCompiler, error_ *JSError) *ErrorPass {
	return &ErrorPass{
		compiler: compiler,
		error_:   error_,
	}
}

// Java: ErrorPass.process(Node, Node)
func (p *ErrorPass) Process(externs *rhino.Node, root *rhino.Node) {
	p.compiler.Report(p.error_)
}
