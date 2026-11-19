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
// Port of com/google/javascript/jscomp/CompilerOptionsPreprocessor.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.jscomp.CompilerOptionsPreprocessor
type CompilerOptionsPreprocessor struct {
	idHash_ atomic.Int32
}

func (this_ *CompilerOptionsPreprocessor) BindSelf_(s any) {
	_ = s
}
func (this_ *CompilerOptionsPreprocessor) IsJavaNil_() bool { return this_ == nil }
func (*CompilerOptionsPreprocessor) JavaClassInfo_() (string, string) { return "com.google.javascript.jscomp.CompilerOptionsPreprocessor", "CompilerOptionsPreprocessor" }
func (this_ *CompilerOptionsPreprocessor) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*CompilerOptionsPreprocessor) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.CompilerOptionsPreprocessor", "com.google.javascript.jscomp.CompilerOptionsPreprocessor", func(v_ any) bool {
		return j2gort.InstanceOf[*CompilerOptionsPreprocessor](v_)
	})
}

func (this_ *CompilerOptionsPreprocessor) initFields_CompilerOptionsPreprocessor() {
	_ = this_
}

// Java constructor: CompilerOptionsPreprocessor()
func NewCompilerOptionsPreprocessor() *CompilerOptionsPreprocessor {
	this_ := &CompilerOptionsPreprocessor{}
	InitCompilerOptionsPreprocessor(this_)
	return this_
}

func InitCompilerOptionsPreprocessor(this_ *CompilerOptionsPreprocessor) {
	this_.initFields_CompilerOptionsPreprocessor()
}

// Java static method: CompilerOptionsPreprocessor.preprocess(com.google.javascript.jscomp.CompilerOptions options)
func CompilerOptionsPreprocessor_Preprocess(options *CompilerOptions) {
	_ = options
	if (((options.GetInlineFunctionsLevel() == StaticGet_CompilerOptions_Reach_NONE()) && (options.GetMaxFunctionSizeAfterInlining() != int32(-1)))) {
		panic(NewCompilerOptionsPreprocessor_InvalidOptionsException("max_function_size_after_inlining has no effect if inlining is disabled."))
	}
}

// Java class: com.google.javascript.jscomp.CompilerOptionsPreprocessor.InvalidOptionsException
type CompilerOptionsPreprocessor_InvalidOptionsException struct {
	javalang.RuntimeException
	idHash_ atomic.Int32
}

func (this_ *CompilerOptionsPreprocessor_InvalidOptionsException) BindSelf_(s any) {
	_ = s
}
func (this_ *CompilerOptionsPreprocessor_InvalidOptionsException) IsJavaNil_() bool { return this_ == nil }
func (*CompilerOptionsPreprocessor_InvalidOptionsException) JavaClassInfo_() (string, string) { return "com.google.javascript.jscomp.CompilerOptionsPreprocessor$InvalidOptionsException", "InvalidOptionsException" }
func (this_ *CompilerOptionsPreprocessor_InvalidOptionsException) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.CompilerOptionsPreprocessor.InvalidOptionsException", "com.google.javascript.jscomp.CompilerOptionsPreprocessor$InvalidOptionsException", func(v_ any) bool {
		return j2gort.InstanceOf[*CompilerOptionsPreprocessor_InvalidOptionsException](v_)
	})
}

func (this_ *CompilerOptionsPreprocessor_InvalidOptionsException) initFields_CompilerOptionsPreprocessor_InvalidOptionsException() {
	_ = this_
}

// Java constructor: InvalidOptionsException(java.lang.String message)
func NewCompilerOptionsPreprocessor_InvalidOptionsException(message string) *CompilerOptionsPreprocessor_InvalidOptionsException {
	this_ := &CompilerOptionsPreprocessor_InvalidOptionsException{}
	InitCompilerOptionsPreprocessor_InvalidOptionsException(this_, message)
	return this_
}

func InitCompilerOptionsPreprocessor_InvalidOptionsException(this_ *CompilerOptionsPreprocessor_InvalidOptionsException, message string) {
	_ = message
	j2gort.InitUserThrowable(this_, "com.google.javascript.jscomp.CompilerOptionsPreprocessor.InvalidOptionsException", j2gort.BoxString(message))
	this_.initFields_CompilerOptionsPreprocessor_InvalidOptionsException()
}

