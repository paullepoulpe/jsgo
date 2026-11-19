/*
 * Copyright 2016 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/PolymerPassErrors.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
)

// Java class: com.google.javascript.jscomp.PolymerPassErrors
type PolymerPassErrors struct {
	idHash_ atomic.Int32
}

func (this_ *PolymerPassErrors) BindSelf_(s any) {
	_ = s
}
func (this_ *PolymerPassErrors) IsJavaNil_() bool { return this_ == nil }
func (*PolymerPassErrors) JavaClassInfo_() (string, string) { return "com.google.javascript.jscomp.PolymerPassErrors", "PolymerPassErrors" }
func (this_ *PolymerPassErrors) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*PolymerPassErrors) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.PolymerPassErrors", "com.google.javascript.jscomp.PolymerPassErrors", func(v_ any) bool {
		return j2gort.InstanceOf[*PolymerPassErrors](v_)
	})
}

var static_PolymerPassErrors_POLYMER_DESCRIPTOR_NOT_VALID *DiagnosticType
func StaticGet_PolymerPassErrors_POLYMER_DESCRIPTOR_NOT_VALID() *DiagnosticType {
	clinit_PolymerPassErrors()
	return static_PolymerPassErrors_POLYMER_DESCRIPTOR_NOT_VALID
}
func StaticSet_PolymerPassErrors_POLYMER_DESCRIPTOR_NOT_VALID(v *DiagnosticType) {
	clinit_PolymerPassErrors()
	static_PolymerPassErrors_POLYMER_DESCRIPTOR_NOT_VALID = v
}
var static_PolymerPassErrors_POLYMER_INVALID_BEHAVIOR *DiagnosticType
func StaticGet_PolymerPassErrors_POLYMER_INVALID_BEHAVIOR() *DiagnosticType {
	clinit_PolymerPassErrors()
	return static_PolymerPassErrors_POLYMER_INVALID_BEHAVIOR
}
func StaticSet_PolymerPassErrors_POLYMER_INVALID_BEHAVIOR(v *DiagnosticType) {
	clinit_PolymerPassErrors()
	static_PolymerPassErrors_POLYMER_INVALID_BEHAVIOR = v
}
var static_PolymerPassErrors_POLYMER_MISSING_IS *DiagnosticType
func StaticGet_PolymerPassErrors_POLYMER_MISSING_IS() *DiagnosticType {
	clinit_PolymerPassErrors()
	return static_PolymerPassErrors_POLYMER_MISSING_IS
}
func StaticSet_PolymerPassErrors_POLYMER_MISSING_IS(v *DiagnosticType) {
	clinit_PolymerPassErrors()
	static_PolymerPassErrors_POLYMER_MISSING_IS = v
}
var static_PolymerPassErrors_POLYMER_UNEXPECTED_PARAMS *DiagnosticType
func StaticGet_PolymerPassErrors_POLYMER_UNEXPECTED_PARAMS() *DiagnosticType {
	clinit_PolymerPassErrors()
	return static_PolymerPassErrors_POLYMER_UNEXPECTED_PARAMS
}
func StaticSet_PolymerPassErrors_POLYMER_UNEXPECTED_PARAMS(v *DiagnosticType) {
	clinit_PolymerPassErrors()
	static_PolymerPassErrors_POLYMER_UNEXPECTED_PARAMS = v
}
var static_PolymerPassErrors_POLYMER_MISSING_EXTERNS *DiagnosticType
func StaticGet_PolymerPassErrors_POLYMER_MISSING_EXTERNS() *DiagnosticType {
	clinit_PolymerPassErrors()
	return static_PolymerPassErrors_POLYMER_MISSING_EXTERNS
}
func StaticSet_PolymerPassErrors_POLYMER_MISSING_EXTERNS(v *DiagnosticType) {
	clinit_PolymerPassErrors()
	static_PolymerPassErrors_POLYMER_MISSING_EXTERNS = v
}
var static_PolymerPassErrors_POLYMER_INVALID_PROPERTY *DiagnosticType
func StaticGet_PolymerPassErrors_POLYMER_INVALID_PROPERTY() *DiagnosticType {
	clinit_PolymerPassErrors()
	return static_PolymerPassErrors_POLYMER_INVALID_PROPERTY
}
func StaticSet_PolymerPassErrors_POLYMER_INVALID_PROPERTY(v *DiagnosticType) {
	clinit_PolymerPassErrors()
	static_PolymerPassErrors_POLYMER_INVALID_PROPERTY = v
}
var static_PolymerPassErrors_POLYMER_INVALID_EXTENDS *DiagnosticType
func StaticGet_PolymerPassErrors_POLYMER_INVALID_EXTENDS() *DiagnosticType {
	clinit_PolymerPassErrors()
	return static_PolymerPassErrors_POLYMER_INVALID_EXTENDS
}
func StaticSet_PolymerPassErrors_POLYMER_INVALID_EXTENDS(v *DiagnosticType) {
	clinit_PolymerPassErrors()
	static_PolymerPassErrors_POLYMER_INVALID_EXTENDS = v
}
var static_PolymerPassErrors_POLYMER_INVALID_BEHAVIOR_ARRAY *DiagnosticType
func StaticGet_PolymerPassErrors_POLYMER_INVALID_BEHAVIOR_ARRAY() *DiagnosticType {
	clinit_PolymerPassErrors()
	return static_PolymerPassErrors_POLYMER_INVALID_BEHAVIOR_ARRAY
}
func StaticSet_PolymerPassErrors_POLYMER_INVALID_BEHAVIOR_ARRAY(v *DiagnosticType) {
	clinit_PolymerPassErrors()
	static_PolymerPassErrors_POLYMER_INVALID_BEHAVIOR_ARRAY = v
}
var static_PolymerPassErrors_POLYMER_UNQUALIFIED_BEHAVIOR *DiagnosticType
func StaticGet_PolymerPassErrors_POLYMER_UNQUALIFIED_BEHAVIOR() *DiagnosticType {
	clinit_PolymerPassErrors()
	return static_PolymerPassErrors_POLYMER_UNQUALIFIED_BEHAVIOR
}
func StaticSet_PolymerPassErrors_POLYMER_UNQUALIFIED_BEHAVIOR(v *DiagnosticType) {
	clinit_PolymerPassErrors()
	static_PolymerPassErrors_POLYMER_UNQUALIFIED_BEHAVIOR = v
}
var static_PolymerPassErrors_POLYMER_UNANNOTATED_BEHAVIOR *DiagnosticType
func StaticGet_PolymerPassErrors_POLYMER_UNANNOTATED_BEHAVIOR() *DiagnosticType {
	clinit_PolymerPassErrors()
	return static_PolymerPassErrors_POLYMER_UNANNOTATED_BEHAVIOR
}
func StaticSet_PolymerPassErrors_POLYMER_UNANNOTATED_BEHAVIOR(v *DiagnosticType) {
	clinit_PolymerPassErrors()
	static_PolymerPassErrors_POLYMER_UNANNOTATED_BEHAVIOR = v
}
var static_PolymerPassErrors_POLYMER_CLASS_PROPERTIES_INVALID *DiagnosticType
func StaticGet_PolymerPassErrors_POLYMER_CLASS_PROPERTIES_INVALID() *DiagnosticType {
	clinit_PolymerPassErrors()
	return static_PolymerPassErrors_POLYMER_CLASS_PROPERTIES_INVALID
}
func StaticSet_PolymerPassErrors_POLYMER_CLASS_PROPERTIES_INVALID(v *DiagnosticType) {
	clinit_PolymerPassErrors()
	static_PolymerPassErrors_POLYMER_CLASS_PROPERTIES_INVALID = v
}
var static_PolymerPassErrors_POLYMER_CLASS_PROPERTIES_NOT_STATIC *DiagnosticType
func StaticGet_PolymerPassErrors_POLYMER_CLASS_PROPERTIES_NOT_STATIC() *DiagnosticType {
	clinit_PolymerPassErrors()
	return static_PolymerPassErrors_POLYMER_CLASS_PROPERTIES_NOT_STATIC
}
func StaticSet_PolymerPassErrors_POLYMER_CLASS_PROPERTIES_NOT_STATIC(v *DiagnosticType) {
	clinit_PolymerPassErrors()
	static_PolymerPassErrors_POLYMER_CLASS_PROPERTIES_NOT_STATIC = v
}
var static_PolymerPassErrors_POLYMER_CLASS_UNNAMED *DiagnosticType
func StaticGet_PolymerPassErrors_POLYMER_CLASS_UNNAMED() *DiagnosticType {
	clinit_PolymerPassErrors()
	return static_PolymerPassErrors_POLYMER_CLASS_UNNAMED
}
func StaticSet_PolymerPassErrors_POLYMER_CLASS_UNNAMED(v *DiagnosticType) {
	clinit_PolymerPassErrors()
	static_PolymerPassErrors_POLYMER_CLASS_UNNAMED = v
}
var static_PolymerPassErrors_POLYMER_MISPLACED_PROPERTY_JSDOC *DiagnosticType
func StaticGet_PolymerPassErrors_POLYMER_MISPLACED_PROPERTY_JSDOC() *DiagnosticType {
	clinit_PolymerPassErrors()
	return static_PolymerPassErrors_POLYMER_MISPLACED_PROPERTY_JSDOC
}
func StaticSet_PolymerPassErrors_POLYMER_MISPLACED_PROPERTY_JSDOC(v *DiagnosticType) {
	clinit_PolymerPassErrors()
	static_PolymerPassErrors_POLYMER_MISPLACED_PROPERTY_JSDOC = v
}
var static_PolymerPassErrors_POLYMER_UNPARSABLE_STRING *DiagnosticType
func StaticGet_PolymerPassErrors_POLYMER_UNPARSABLE_STRING() *DiagnosticType {
	clinit_PolymerPassErrors()
	return static_PolymerPassErrors_POLYMER_UNPARSABLE_STRING
}
func StaticSet_PolymerPassErrors_POLYMER_UNPARSABLE_STRING(v *DiagnosticType) {
	clinit_PolymerPassErrors()
	static_PolymerPassErrors_POLYMER_UNPARSABLE_STRING = v
}
var clinitDone_PolymerPassErrors uint32
var clinitState_PolymerPassErrors j2gort.ClinitState
func clinit_PolymerPassErrors() {
	if atomic.LoadUint32(&clinitDone_PolymerPassErrors) != 0 {
		return
	}
	clinitSlow_PolymerPassErrors()
}

func clinitSlow_PolymerPassErrors() {
	if !j2gort.EnterClinit(&clinitState_PolymerPassErrors, &clinitDone_PolymerPassErrors) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_PolymerPassErrors, &clinitDone_PolymerPassErrors)
	static_PolymerPassErrors_POLYMER_DESCRIPTOR_NOT_VALID = DiagnosticType_Warning("JSC_POLYMER_DESCRIPTOR_NOT_VALID", "The argument to Polymer() is not an obj lit or the Polymer 2 class does not have a static getter named 'config'. Ignoring this definition.")
	static_PolymerPassErrors_POLYMER_INVALID_BEHAVIOR = DiagnosticType_Error("JSC_POLYMER_INVALID_BEHAVIOR", "A Polymer behavior may not include an ''is'' property.")
	static_PolymerPassErrors_POLYMER_MISSING_IS = DiagnosticType_Error("JSC_POLYMER_MISSING_IS", "The class descriptor must include an ''is'' property.")
	static_PolymerPassErrors_POLYMER_UNEXPECTED_PARAMS = DiagnosticType_Error("JSC_POLYMER_UNEXPECTED_PARAMS", "The class definition has too many arguments.")
	static_PolymerPassErrors_POLYMER_MISSING_EXTERNS = DiagnosticType_Error("JSC_POLYMER_MISSING_EXTERNS", "Missing Polymer externs.")
	static_PolymerPassErrors_POLYMER_INVALID_PROPERTY = DiagnosticType_Error("JSC_POLYMER_INVALID_PROPERTY", "Polymer property has an invalid or missing type.")
	static_PolymerPassErrors_POLYMER_INVALID_EXTENDS = DiagnosticType_Error("JSC_POLYMER_INVALID_EXTENDS", "Cannot extend HTML element ''{0}''. The element is probably either misspelled, or needs to be added to the list of known elements.")
	static_PolymerPassErrors_POLYMER_INVALID_BEHAVIOR_ARRAY = DiagnosticType_Error("JSC_POLYMER_INVALID_BEHAVIOR_ARRAY", "The behaviors property must be an array literal.")
	static_PolymerPassErrors_POLYMER_UNQUALIFIED_BEHAVIOR = DiagnosticType_Error("JSC_POLYMER_UNQUALIFIED_BEHAVIOR", "Behaviors must be global names or qualified names that are declared as object literals or array literals of other valid Behaviors.")
	static_PolymerPassErrors_POLYMER_UNANNOTATED_BEHAVIOR = DiagnosticType_Error("JSC_POLYMER_UNANNOTATED_BEHAVIOR", "Behavior declarations must be annotated with @polymerBehavior.")
	static_PolymerPassErrors_POLYMER_CLASS_PROPERTIES_INVALID = DiagnosticType_Error("JSC_POLYMER_CLASS_PROPERTIES_INVALID", "The Polymer element class 'properties' getter does not return an object literal. Ignoring this definition.")
	static_PolymerPassErrors_POLYMER_CLASS_PROPERTIES_NOT_STATIC = DiagnosticType_Error("JSC_POLYMER_CLASS_PROPERTIES_NOT_STATIC", "The Polymer element class 'properties' getter is not declared static. Ignoring this definition.")
	static_PolymerPassErrors_POLYMER_CLASS_UNNAMED = DiagnosticType_Warning("JSC_POLYMER2_UNNAMED", "Unable to locate a valid name for the Polymer element class.Ignoring this definition.")
	static_PolymerPassErrors_POLYMER_MISPLACED_PROPERTY_JSDOC = DiagnosticType_Warning("JSC_POLYMER_MISPLACED_PROPERTY_JSDOC", "When a Polymer property is declared in the constructor, its JSDoc should only be in the constructor, not on the Polymer properties configuration.")
	static_PolymerPassErrors_POLYMER_UNPARSABLE_STRING = DiagnosticType_Error("JSC_POLYMER_UNPARSABLE_STRING", "The Polymer computed property or complex observer string could not be parsed.")
}

func (this_ *PolymerPassErrors) initFields_PolymerPassErrors() {
	_ = this_
}

// Java constructor: PolymerPassErrors()
func NewPolymerPassErrors() *PolymerPassErrors {
	clinit_PolymerPassErrors()
	this_ := &PolymerPassErrors{}
	InitPolymerPassErrors(this_)
	return this_
}

func InitPolymerPassErrors(this_ *PolymerPassErrors) {
	this_.initFields_PolymerPassErrors()
}

