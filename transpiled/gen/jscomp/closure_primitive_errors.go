/*
 * Copyright 2018 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/ClosurePrimitiveErrors.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
)

// Java class: com.google.javascript.jscomp.ClosurePrimitiveErrors
type ClosurePrimitiveErrors struct {
	idHash_ atomic.Int32
}

func (this_ *ClosurePrimitiveErrors) BindSelf_(s any) {
	_ = s
}
func (this_ *ClosurePrimitiveErrors) IsJavaNil_() bool { return this_ == nil }
func (*ClosurePrimitiveErrors) JavaClassInfo_() (string, string) { return "com.google.javascript.jscomp.ClosurePrimitiveErrors", "ClosurePrimitiveErrors" }
func (this_ *ClosurePrimitiveErrors) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*ClosurePrimitiveErrors) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.ClosurePrimitiveErrors", "com.google.javascript.jscomp.ClosurePrimitiveErrors", func(v_ any) bool {
		return j2gort.InstanceOf[*ClosurePrimitiveErrors](v_)
	})
}

var static_ClosurePrimitiveErrors_DUPLICATE_MODULE *DiagnosticType
func StaticGet_ClosurePrimitiveErrors_DUPLICATE_MODULE() *DiagnosticType {
	clinit_ClosurePrimitiveErrors()
	return static_ClosurePrimitiveErrors_DUPLICATE_MODULE
}
func StaticSet_ClosurePrimitiveErrors_DUPLICATE_MODULE(v *DiagnosticType) {
	clinit_ClosurePrimitiveErrors()
	static_ClosurePrimitiveErrors_DUPLICATE_MODULE = v
}
var static_ClosurePrimitiveErrors_DUPLICATE_NAMESPACE *DiagnosticType
func StaticGet_ClosurePrimitiveErrors_DUPLICATE_NAMESPACE() *DiagnosticType {
	clinit_ClosurePrimitiveErrors()
	return static_ClosurePrimitiveErrors_DUPLICATE_NAMESPACE
}
func StaticSet_ClosurePrimitiveErrors_DUPLICATE_NAMESPACE(v *DiagnosticType) {
	clinit_ClosurePrimitiveErrors()
	static_ClosurePrimitiveErrors_DUPLICATE_NAMESPACE = v
}
var static_ClosurePrimitiveErrors_DUPLICATE_NAMESPACE_AND_MODULE *DiagnosticType
func StaticGet_ClosurePrimitiveErrors_DUPLICATE_NAMESPACE_AND_MODULE() *DiagnosticType {
	clinit_ClosurePrimitiveErrors()
	return static_ClosurePrimitiveErrors_DUPLICATE_NAMESPACE_AND_MODULE
}
func StaticSet_ClosurePrimitiveErrors_DUPLICATE_NAMESPACE_AND_MODULE(v *DiagnosticType) {
	clinit_ClosurePrimitiveErrors()
	static_ClosurePrimitiveErrors_DUPLICATE_NAMESPACE_AND_MODULE = v
}
var static_ClosurePrimitiveErrors_INVALID_DESTRUCTURING_FORWARD_DECLARE *DiagnosticType
func StaticGet_ClosurePrimitiveErrors_INVALID_DESTRUCTURING_FORWARD_DECLARE() *DiagnosticType {
	clinit_ClosurePrimitiveErrors()
	return static_ClosurePrimitiveErrors_INVALID_DESTRUCTURING_FORWARD_DECLARE
}
func StaticSet_ClosurePrimitiveErrors_INVALID_DESTRUCTURING_FORWARD_DECLARE(v *DiagnosticType) {
	clinit_ClosurePrimitiveErrors()
	static_ClosurePrimitiveErrors_INVALID_DESTRUCTURING_FORWARD_DECLARE = v
}
var static_ClosurePrimitiveErrors_MODULE_USES_GOOG_MODULE_GET *DiagnosticType
func StaticGet_ClosurePrimitiveErrors_MODULE_USES_GOOG_MODULE_GET() *DiagnosticType {
	clinit_ClosurePrimitiveErrors()
	return static_ClosurePrimitiveErrors_MODULE_USES_GOOG_MODULE_GET
}
func StaticSet_ClosurePrimitiveErrors_MODULE_USES_GOOG_MODULE_GET(v *DiagnosticType) {
	clinit_ClosurePrimitiveErrors()
	static_ClosurePrimitiveErrors_MODULE_USES_GOOG_MODULE_GET = v
}
var static_ClosurePrimitiveErrors_INVALID_FORWARD_DECLARE_NAMESPACE *DiagnosticType
func StaticGet_ClosurePrimitiveErrors_INVALID_FORWARD_DECLARE_NAMESPACE() *DiagnosticType {
	clinit_ClosurePrimitiveErrors()
	return static_ClosurePrimitiveErrors_INVALID_FORWARD_DECLARE_NAMESPACE
}
func StaticSet_ClosurePrimitiveErrors_INVALID_FORWARD_DECLARE_NAMESPACE(v *DiagnosticType) {
	clinit_ClosurePrimitiveErrors()
	static_ClosurePrimitiveErrors_INVALID_FORWARD_DECLARE_NAMESPACE = v
}
var static_ClosurePrimitiveErrors_INVALID_GET_NAMESPACE *DiagnosticType
func StaticGet_ClosurePrimitiveErrors_INVALID_GET_NAMESPACE() *DiagnosticType {
	clinit_ClosurePrimitiveErrors()
	return static_ClosurePrimitiveErrors_INVALID_GET_NAMESPACE
}
func StaticSet_ClosurePrimitiveErrors_INVALID_GET_NAMESPACE(v *DiagnosticType) {
	clinit_ClosurePrimitiveErrors()
	static_ClosurePrimitiveErrors_INVALID_GET_NAMESPACE = v
}
var static_ClosurePrimitiveErrors_INVALID_REQUIRE_DYNAMIC *DiagnosticType
func StaticGet_ClosurePrimitiveErrors_INVALID_REQUIRE_DYNAMIC() *DiagnosticType {
	clinit_ClosurePrimitiveErrors()
	return static_ClosurePrimitiveErrors_INVALID_REQUIRE_DYNAMIC
}
func StaticSet_ClosurePrimitiveErrors_INVALID_REQUIRE_DYNAMIC(v *DiagnosticType) {
	clinit_ClosurePrimitiveErrors()
	static_ClosurePrimitiveErrors_INVALID_REQUIRE_DYNAMIC = v
}
var static_ClosurePrimitiveErrors_INVALID_REQUIRE_NAMESPACE *DiagnosticType
func StaticGet_ClosurePrimitiveErrors_INVALID_REQUIRE_NAMESPACE() *DiagnosticType {
	clinit_ClosurePrimitiveErrors()
	return static_ClosurePrimitiveErrors_INVALID_REQUIRE_NAMESPACE
}
func StaticSet_ClosurePrimitiveErrors_INVALID_REQUIRE_NAMESPACE(v *DiagnosticType) {
	clinit_ClosurePrimitiveErrors()
	static_ClosurePrimitiveErrors_INVALID_REQUIRE_NAMESPACE = v
}
var static_ClosurePrimitiveErrors_INVALID_REQUIRE_TYPE_NAMESPACE *DiagnosticType
func StaticGet_ClosurePrimitiveErrors_INVALID_REQUIRE_TYPE_NAMESPACE() *DiagnosticType {
	clinit_ClosurePrimitiveErrors()
	return static_ClosurePrimitiveErrors_INVALID_REQUIRE_TYPE_NAMESPACE
}
func StaticSet_ClosurePrimitiveErrors_INVALID_REQUIRE_TYPE_NAMESPACE(v *DiagnosticType) {
	clinit_ClosurePrimitiveErrors()
	static_ClosurePrimitiveErrors_INVALID_REQUIRE_TYPE_NAMESPACE = v
}
var static_ClosurePrimitiveErrors_MISSING_MODULE_OR_PROVIDE *DiagnosticType
func StaticGet_ClosurePrimitiveErrors_MISSING_MODULE_OR_PROVIDE() *DiagnosticType {
	clinit_ClosurePrimitiveErrors()
	return static_ClosurePrimitiveErrors_MISSING_MODULE_OR_PROVIDE
}
func StaticSet_ClosurePrimitiveErrors_MISSING_MODULE_OR_PROVIDE(v *DiagnosticType) {
	clinit_ClosurePrimitiveErrors()
	static_ClosurePrimitiveErrors_MISSING_MODULE_OR_PROVIDE = v
}
var static_ClosurePrimitiveErrors_MISSING_MODULE_OR_PROVIDE_FOR_FORWARD_DECLARE *DiagnosticType
func StaticGet_ClosurePrimitiveErrors_MISSING_MODULE_OR_PROVIDE_FOR_FORWARD_DECLARE() *DiagnosticType {
	clinit_ClosurePrimitiveErrors()
	return static_ClosurePrimitiveErrors_MISSING_MODULE_OR_PROVIDE_FOR_FORWARD_DECLARE
}
func StaticSet_ClosurePrimitiveErrors_MISSING_MODULE_OR_PROVIDE_FOR_FORWARD_DECLARE(v *DiagnosticType) {
	clinit_ClosurePrimitiveErrors()
	static_ClosurePrimitiveErrors_MISSING_MODULE_OR_PROVIDE_FOR_FORWARD_DECLARE = v
}
var static_ClosurePrimitiveErrors_INVALID_GET_CALL_SCOPE *DiagnosticType
func StaticGet_ClosurePrimitiveErrors_INVALID_GET_CALL_SCOPE() *DiagnosticType {
	clinit_ClosurePrimitiveErrors()
	return static_ClosurePrimitiveErrors_INVALID_GET_CALL_SCOPE
}
func StaticSet_ClosurePrimitiveErrors_INVALID_GET_CALL_SCOPE(v *DiagnosticType) {
	clinit_ClosurePrimitiveErrors()
	static_ClosurePrimitiveErrors_INVALID_GET_CALL_SCOPE = v
}
var static_ClosurePrimitiveErrors_INVALID_CLOSURE_CALL_SCOPE_ERROR *DiagnosticType
func StaticGet_ClosurePrimitiveErrors_INVALID_CLOSURE_CALL_SCOPE_ERROR() *DiagnosticType {
	clinit_ClosurePrimitiveErrors()
	return static_ClosurePrimitiveErrors_INVALID_CLOSURE_CALL_SCOPE_ERROR
}
func StaticSet_ClosurePrimitiveErrors_INVALID_CLOSURE_CALL_SCOPE_ERROR(v *DiagnosticType) {
	clinit_ClosurePrimitiveErrors()
	static_ClosurePrimitiveErrors_INVALID_CLOSURE_CALL_SCOPE_ERROR = v
}
var static_ClosurePrimitiveErrors_TOO_MANY_ARGUMENTS_ERROR *DiagnosticType
func StaticGet_ClosurePrimitiveErrors_TOO_MANY_ARGUMENTS_ERROR() *DiagnosticType {
	clinit_ClosurePrimitiveErrors()
	return static_ClosurePrimitiveErrors_TOO_MANY_ARGUMENTS_ERROR
}
func StaticSet_ClosurePrimitiveErrors_TOO_MANY_ARGUMENTS_ERROR(v *DiagnosticType) {
	clinit_ClosurePrimitiveErrors()
	static_ClosurePrimitiveErrors_TOO_MANY_ARGUMENTS_ERROR = v
}
var static_ClosurePrimitiveErrors_NULL_ARGUMENT_ERROR *DiagnosticType
func StaticGet_ClosurePrimitiveErrors_NULL_ARGUMENT_ERROR() *DiagnosticType {
	clinit_ClosurePrimitiveErrors()
	return static_ClosurePrimitiveErrors_NULL_ARGUMENT_ERROR
}
func StaticSet_ClosurePrimitiveErrors_NULL_ARGUMENT_ERROR(v *DiagnosticType) {
	clinit_ClosurePrimitiveErrors()
	static_ClosurePrimitiveErrors_NULL_ARGUMENT_ERROR = v
}
var static_ClosurePrimitiveErrors_INVALID_ARGUMENT_ERROR *DiagnosticType
func StaticGet_ClosurePrimitiveErrors_INVALID_ARGUMENT_ERROR() *DiagnosticType {
	clinit_ClosurePrimitiveErrors()
	return static_ClosurePrimitiveErrors_INVALID_ARGUMENT_ERROR
}
func StaticSet_ClosurePrimitiveErrors_INVALID_ARGUMENT_ERROR(v *DiagnosticType) {
	clinit_ClosurePrimitiveErrors()
	static_ClosurePrimitiveErrors_INVALID_ARGUMENT_ERROR = v
}
var static_ClosurePrimitiveErrors_GOOG_MODULE_GET_OF_WEAK_MODULE *DiagnosticType
func StaticGet_ClosurePrimitiveErrors_GOOG_MODULE_GET_OF_WEAK_MODULE() *DiagnosticType {
	clinit_ClosurePrimitiveErrors()
	return static_ClosurePrimitiveErrors_GOOG_MODULE_GET_OF_WEAK_MODULE
}
func StaticSet_ClosurePrimitiveErrors_GOOG_MODULE_GET_OF_WEAK_MODULE(v *DiagnosticType) {
	clinit_ClosurePrimitiveErrors()
	static_ClosurePrimitiveErrors_GOOG_MODULE_GET_OF_WEAK_MODULE = v
}
var clinitDone_ClosurePrimitiveErrors uint32
var clinitState_ClosurePrimitiveErrors j2gort.ClinitState
func clinit_ClosurePrimitiveErrors() {
	if atomic.LoadUint32(&clinitDone_ClosurePrimitiveErrors) != 0 {
		return
	}
	clinitSlow_ClosurePrimitiveErrors()
}

func clinitSlow_ClosurePrimitiveErrors() {
	if !j2gort.EnterClinit(&clinitState_ClosurePrimitiveErrors, &clinitDone_ClosurePrimitiveErrors) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_ClosurePrimitiveErrors, &clinitDone_ClosurePrimitiveErrors)
	static_ClosurePrimitiveErrors_DUPLICATE_MODULE = DiagnosticType_Error("JSC_DUPLICATE_MODULE", "Duplicate module: {0}. Module already declared in other file {1}")
	static_ClosurePrimitiveErrors_DUPLICATE_NAMESPACE = DiagnosticType_Error("JSC_DUPLICATE_NAMESPACE", "Duplicate namespace: {0}. Namespace already provided in other file {1}")
	static_ClosurePrimitiveErrors_DUPLICATE_NAMESPACE_AND_MODULE = DiagnosticType_Error("JSC_DUPLICATE_NAMESPACE_AND_MODULE", "Duplicate module: {0}. Module already goog.provided in other file {1}")
	static_ClosurePrimitiveErrors_INVALID_DESTRUCTURING_FORWARD_DECLARE = DiagnosticType_Error("JSC_INVALID_DESTRUCTURING_FORWARD_DECLARE", "Cannot destructure a forward-declared type")
	static_ClosurePrimitiveErrors_MODULE_USES_GOOG_MODULE_GET = DiagnosticType_Error("JSC_MODULE_USES_GOOG_MODULE_GET", "It's illegal to use a 'goog.module.get' at the module top-level. Did you mean to use goog.require instead?")
	static_ClosurePrimitiveErrors_INVALID_FORWARD_DECLARE_NAMESPACE = DiagnosticType_Error("JSC_GOOG_MODULE_INVALID_FORWARD_DECLARE_NAMESPACE", "goog.forwardDeclare parameter must be a string literal.")
	static_ClosurePrimitiveErrors_INVALID_GET_NAMESPACE = DiagnosticType_Error("JSC_GOOG_MODULE_INVALID_GET_NAMESPACE", "goog.module.get parameter must be a string literal.")
	static_ClosurePrimitiveErrors_INVALID_REQUIRE_DYNAMIC = DiagnosticType_Error("JSC_GOOG_MODULE_INVALID_REQUIRE_DYNAMIC_NAMESPACE", "goog.requireDynamic parameter must be a string literal.")
	static_ClosurePrimitiveErrors_INVALID_REQUIRE_NAMESPACE = DiagnosticType_Error("JSC_GOOG_MODULE_INVALID_REQUIRE_NAMESPACE", "goog.require parameter must be a string literal.")
	static_ClosurePrimitiveErrors_INVALID_REQUIRE_TYPE_NAMESPACE = DiagnosticType_Error("JSC_GOOG_MODULE_INVALID_REQUIRE_TYPE_NAMESPACE", "goog.requireType parameter must be a string literal.")
	static_ClosurePrimitiveErrors_MISSING_MODULE_OR_PROVIDE = DiagnosticType_Error("JSC_MISSING_MODULE_OR_PROVIDE", "Required namespace \"{0}\" never defined.")
	static_ClosurePrimitiveErrors_MISSING_MODULE_OR_PROVIDE_FOR_FORWARD_DECLARE = DiagnosticType_Error("JSC_MISSING_MODULE_OR_PROVIDE_FOR_FORWARD_DECLARE", "Required namespace \"{0}\" never defined.")
	static_ClosurePrimitiveErrors_INVALID_GET_CALL_SCOPE = DiagnosticType_Error("JSC_GOOG_MODULE_INVALID_GET_CALL_SCOPE", "goog.module.get values should not be assigned to global variables.")
	static_ClosurePrimitiveErrors_INVALID_CLOSURE_CALL_SCOPE_ERROR = DiagnosticType_Error("JSC_INVALID_CLOSURE_CALL_ERROR", "Closure primitive methods (goog.provide, goog.require, goog.define, etc) must be called at file scope.")
	static_ClosurePrimitiveErrors_TOO_MANY_ARGUMENTS_ERROR = DiagnosticType_Error("JSC_TOO_MANY_ARGUMENTS_ERROR", "method \"{0}\" called with more than one argument")
	static_ClosurePrimitiveErrors_NULL_ARGUMENT_ERROR = DiagnosticType_Error("JSC_NULL_ARGUMENT_ERROR", "method \"{0}\" called without an argument")
	static_ClosurePrimitiveErrors_INVALID_ARGUMENT_ERROR = DiagnosticType_Error("JSC_INVALID_ARGUMENT_ERROR", "method \"{0}\" called with invalid argument")
	static_ClosurePrimitiveErrors_GOOG_MODULE_GET_OF_WEAK_MODULE = DiagnosticType_Error("JSC_GOOG_MODULE_GET_OF_WEAK_MODULE", "goog.module.get of \"{0}\" references a namespace that is only imported by goog.requireType.")
}

func (this_ *ClosurePrimitiveErrors) initFields_ClosurePrimitiveErrors() {
	_ = this_
}

// Java constructor: ClosurePrimitiveErrors()
func NewClosurePrimitiveErrors() *ClosurePrimitiveErrors {
	clinit_ClosurePrimitiveErrors()
	this_ := &ClosurePrimitiveErrors{}
	InitClosurePrimitiveErrors(this_)
	return this_
}

func InitClosurePrimitiveErrors(this_ *ClosurePrimitiveErrors) {
	this_.initFields_ClosurePrimitiveErrors()
}

