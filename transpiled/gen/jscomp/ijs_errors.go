/*
 * Copyright 2019 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/ijs/IjsErrors.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
)

// Java class: com.google.javascript.jscomp.ijs.IjsErrors
type IjsErrors struct {
	idHash_ atomic.Int32
}

func (this_ *IjsErrors) BindSelf_(s any) {
	_ = s
}
func (this_ *IjsErrors) IsJavaNil_() bool { return this_ == nil }
func (*IjsErrors) JavaClassInfo_() (string, string) { return "com.google.javascript.jscomp.ijs.IjsErrors", "IjsErrors" }
func (this_ *IjsErrors) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*IjsErrors) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.ijs.IjsErrors", "com.google.javascript.jscomp.ijs.IjsErrors", func(v_ any) bool {
		return j2gort.InstanceOf[*IjsErrors](v_)
	})
}

var static_IjsErrors_BAD_IJS_FILE_NAME *DiagnosticType
func StaticGet_IjsErrors_BAD_IJS_FILE_NAME() *DiagnosticType {
	clinit_IjsErrors()
	return static_IjsErrors_BAD_IJS_FILE_NAME
}
func StaticSet_IjsErrors_BAD_IJS_FILE_NAME(v *DiagnosticType) {
	clinit_IjsErrors()
	static_IjsErrors_BAD_IJS_FILE_NAME = v
}
var static_IjsErrors_CONFLICTING_IJS_FILE *DiagnosticType
func StaticGet_IjsErrors_CONFLICTING_IJS_FILE() *DiagnosticType {
	clinit_IjsErrors()
	return static_IjsErrors_CONFLICTING_IJS_FILE
}
func StaticSet_IjsErrors_CONFLICTING_IJS_FILE(v *DiagnosticType) {
	clinit_IjsErrors()
	static_IjsErrors_CONFLICTING_IJS_FILE = v
}
var clinitDone_IjsErrors uint32
var clinitState_IjsErrors j2gort.ClinitState
func clinit_IjsErrors() {
	if atomic.LoadUint32(&clinitDone_IjsErrors) != 0 {
		return
	}
	clinitSlow_IjsErrors()
}

func clinitSlow_IjsErrors() {
	if !j2gort.EnterClinit(&clinitState_IjsErrors, &clinitDone_IjsErrors) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_IjsErrors, &clinitDone_IjsErrors)
	static_IjsErrors_BAD_IJS_FILE_NAME = DiagnosticType_Error("JSC_BAD_IJS_FILE_NAME", "--ijs file name must end in .i.js: {0}")
	static_IjsErrors_CONFLICTING_IJS_FILE = DiagnosticType_Warning("JSC_CONFLICTING_IJS_FILE", "Found --ijs file in a conflict with --js or --weakdep file. This usually means the same file appears both in the srcs and the transitive deps of a build target.\nFile 1: {0}\nFile 2: {1}")
}

func (this_ *IjsErrors) initFields_IjsErrors() {
	_ = this_
}

// Java constructor: IjsErrors()
func NewIjsErrors() *IjsErrors {
	clinit_IjsErrors()
	this_ := &IjsErrors{}
	InitIjsErrors(this_)
	return this_
}

func InitIjsErrors(this_ *IjsErrors) {
	this_.initFields_IjsErrors()
}

