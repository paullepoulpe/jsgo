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
// Port of com/google/javascript/jscomp/TranspilationUtil.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.jscomp.TranspilationUtil
type TranspilationUtil struct {
	idHash_ atomic.Int32
}

func (this_ *TranspilationUtil) BindSelf_(s any) {
	_ = s
}
func (this_ *TranspilationUtil) IsJavaNil_() bool { return this_ == nil }
func (*TranspilationUtil) JavaClassInfo_() (string, string) { return "com.google.javascript.jscomp.TranspilationUtil", "TranspilationUtil" }
func (this_ *TranspilationUtil) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*TranspilationUtil) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.TranspilationUtil", "com.google.javascript.jscomp.TranspilationUtil", func(v_ any) bool {
		return j2gort.InstanceOf[*TranspilationUtil](v_)
	})
}

var static_TranspilationUtil_CANNOT_CONVERT *DiagnosticType
func StaticGet_TranspilationUtil_CANNOT_CONVERT() *DiagnosticType {
	clinit_TranspilationUtil()
	return static_TranspilationUtil_CANNOT_CONVERT
}
func StaticSet_TranspilationUtil_CANNOT_CONVERT(v *DiagnosticType) {
	clinit_TranspilationUtil()
	static_TranspilationUtil_CANNOT_CONVERT = v
}
var static_TranspilationUtil_CANNOT_CONVERT_YET *DiagnosticType
func StaticGet_TranspilationUtil_CANNOT_CONVERT_YET() *DiagnosticType {
	clinit_TranspilationUtil()
	return static_TranspilationUtil_CANNOT_CONVERT_YET
}
func StaticSet_TranspilationUtil_CANNOT_CONVERT_YET(v *DiagnosticType) {
	clinit_TranspilationUtil()
	static_TranspilationUtil_CANNOT_CONVERT_YET = v
}
var clinitDone_TranspilationUtil uint32
var clinitState_TranspilationUtil j2gort.ClinitState
func clinit_TranspilationUtil() {
	if atomic.LoadUint32(&clinitDone_TranspilationUtil) != 0 {
		return
	}
	clinitSlow_TranspilationUtil()
}

func clinitSlow_TranspilationUtil() {
	if !j2gort.EnterClinit(&clinitState_TranspilationUtil, &clinitDone_TranspilationUtil) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_TranspilationUtil, &clinitDone_TranspilationUtil)
	static_TranspilationUtil_CANNOT_CONVERT = DiagnosticType_Error("JSC_CANNOT_CONVERT", "This code cannot be transpiled. {0}")
	static_TranspilationUtil_CANNOT_CONVERT_YET = DiagnosticType_Error("JSC_CANNOT_CONVERT_YET", "Transpilation of ''{0}'' is not yet implemented.")
}

func (this_ *TranspilationUtil) initFields_TranspilationUtil() {
	_ = this_
}

// Java constructor: TranspilationUtil()
func NewTranspilationUtil() *TranspilationUtil {
	clinit_TranspilationUtil()
	this_ := &TranspilationUtil{}
	InitTranspilationUtil(this_)
	return this_
}

func InitTranspilationUtil(this_ *TranspilationUtil) {
	this_.initFields_TranspilationUtil()
}

// Java static method: TranspilationUtil.cannotConvert(com.google.javascript.jscomp.AbstractCompiler compiler, com.google.javascript.rhino.Node n, java.lang.String message)
func TranspilationUtil_CannotConvert(compiler AbstractCompiler, n pkg_rhino.Node, message string) {
	clinit_TranspilationUtil()
	_ = compiler
	_ = n
	_ = message
	compiler.Report_1(JSError_Make_3(n, StaticGet_TranspilationUtil_CANNOT_CONVERT(), []string{message}))
}

// Java static method: TranspilationUtil.cannotConvertYet(com.google.javascript.jscomp.AbstractCompiler compiler, com.google.javascript.rhino.Node n, java.lang.String feature)
func TranspilationUtil_CannotConvertYet(compiler AbstractCompiler, n pkg_rhino.Node, feature string) {
	clinit_TranspilationUtil()
	_ = compiler
	_ = n
	_ = feature
	compiler.Report_1(JSError_Make_3(n, StaticGet_TranspilationUtil_CANNOT_CONVERT_YET(), []string{feature}))
}

// Java static method: TranspilationUtil.preloadTranspilationRuntimeFunction(com.google.javascript.jscomp.AbstractCompiler compiler, java.lang.String function)
func TranspilationUtil_PreloadTranspilationRuntimeFunction(compiler AbstractCompiler, function string) {
	clinit_TranspilationUtil()
	_ = compiler
	_ = function
	compiler.GetRuntimeJsLibManager().EnsureLibraryInjected(("es6/util/" + j2gort.StrValString(javalang.StringToLowerCaseASCII(function))), false)
}

