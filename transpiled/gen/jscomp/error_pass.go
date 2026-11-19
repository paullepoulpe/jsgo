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
// Port of com/google/javascript/jscomp/ErrorPass.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
)

// Java class: com.google.javascript.jscomp.ErrorPass
type ErrorPass struct {
	compiler AbstractCompiler `gson:"compiler"`
	error_ *JSError `gson:"error"`
	idHash_ atomic.Int32
}

func (this_ *ErrorPass) BindSelf_(s any) {
	_ = s
}
func (this_ *ErrorPass) IsJavaNil_() bool { return this_ == nil }
func (*ErrorPass) JavaClassInfo_() (string, string) { return "com.google.javascript.jscomp.ErrorPass", "ErrorPass" }
func (this_ *ErrorPass) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*ErrorPass) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.ErrorPass", "com.google.javascript.jscomp.ErrorPass", func(v_ any) bool {
		return j2gort.InstanceOf[*ErrorPass](v_)
	})
}
func (*ErrorPass) Implements_jscomp_CompilerPass() {}

func (this_ *ErrorPass) initFields_ErrorPass() {
	_ = this_
}

// Java constructor: ErrorPass(com.google.javascript.jscomp.AbstractCompiler compiler, com.google.javascript.jscomp.DiagnosticType error)
func NewErrorPass_2_AbstractCompiler_DiagnosticType(compiler AbstractCompiler, error_ *DiagnosticType) *ErrorPass {
	this_ := &ErrorPass{}
	InitErrorPass_2_AbstractCompiler_DiagnosticType(this_, compiler, error_)
	return this_
}

func InitErrorPass_2_AbstractCompiler_DiagnosticType(this_ *ErrorPass, compiler AbstractCompiler, error_ *DiagnosticType) {
	_ = compiler
	_ = error_
	InitErrorPass_2_AbstractCompiler_JSError(this_, compiler, JSError_Make_2(error_, []string{}))
}

// Java constructor: ErrorPass(com.google.javascript.jscomp.AbstractCompiler compiler, com.google.javascript.jscomp.JSError error)
func NewErrorPass_2_AbstractCompiler_JSError(compiler AbstractCompiler, error_ *JSError) *ErrorPass {
	this_ := &ErrorPass{}
	InitErrorPass_2_AbstractCompiler_JSError(this_, compiler, error_)
	return this_
}

func InitErrorPass_2_AbstractCompiler_JSError(this_ *ErrorPass, compiler AbstractCompiler, error_ *JSError) {
	_ = compiler
	_ = error_
	this_.initFields_ErrorPass()
	this_.compiler = compiler
	this_.error_ = error_
}

// Java method: ErrorPass.process(com.google.javascript.rhino.Node externs, com.google.javascript.rhino.Node root)
func (this_ *ErrorPass) Process_2(externs pkg_rhino.Node, root pkg_rhino.Node) {
	_ = externs
	_ = root
	this_.compiler.Report_1(this_.error_)
}

