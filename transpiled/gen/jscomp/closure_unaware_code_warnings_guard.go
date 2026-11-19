/*
 * Copyright 2024 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/ClosureUnawareCodeWarningsGuard.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
)

// Java class: com.google.javascript.jscomp.ClosureUnawareCodeWarningsGuard
type ClosureUnawareCodeWarningsGuard struct {
	WarningsGuardBase
}

func (this_ *ClosureUnawareCodeWarningsGuard) BindSelf_(s any) {
	_ = s
}
func (this_ *ClosureUnawareCodeWarningsGuard) IsJavaNil_() bool { return this_ == nil }
func (*ClosureUnawareCodeWarningsGuard) JavaClassInfo_() (string, string) { return "com.google.javascript.jscomp.ClosureUnawareCodeWarningsGuard", "ClosureUnawareCodeWarningsGuard" }
func (this_ *ClosureUnawareCodeWarningsGuard) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.WarningsGuardBase.IdentityHash_()
}
func (*ClosureUnawareCodeWarningsGuard) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.ClosureUnawareCodeWarningsGuard", "com.google.javascript.jscomp.ClosureUnawareCodeWarningsGuard", func(v_ any) bool {
		return j2gort.InstanceOf[*ClosureUnawareCodeWarningsGuard](v_)
	})
}
func (*ClosureUnawareCodeWarningsGuard) Implements_jscomp_WarningsGuard() {}

var static_ClosureUnawareCodeWarningsGuard_DEFAULT_CLOSURE_UNAWARE_CODE_SUPPRESSIONS *DiagnosticGroup
func StaticGet_ClosureUnawareCodeWarningsGuard_DEFAULT_CLOSURE_UNAWARE_CODE_SUPPRESSIONS() *DiagnosticGroup {
	clinit_ClosureUnawareCodeWarningsGuard()
	return static_ClosureUnawareCodeWarningsGuard_DEFAULT_CLOSURE_UNAWARE_CODE_SUPPRESSIONS
}
func StaticSet_ClosureUnawareCodeWarningsGuard_DEFAULT_CLOSURE_UNAWARE_CODE_SUPPRESSIONS(v *DiagnosticGroup) {
	clinit_ClosureUnawareCodeWarningsGuard()
	static_ClosureUnawareCodeWarningsGuard_DEFAULT_CLOSURE_UNAWARE_CODE_SUPPRESSIONS = v
}
var clinitDone_ClosureUnawareCodeWarningsGuard uint32
var clinitState_ClosureUnawareCodeWarningsGuard j2gort.ClinitState
func clinit_ClosureUnawareCodeWarningsGuard() {
	if atomic.LoadUint32(&clinitDone_ClosureUnawareCodeWarningsGuard) != 0 {
		return
	}
	clinitSlow_ClosureUnawareCodeWarningsGuard()
}

func clinitSlow_ClosureUnawareCodeWarningsGuard() {
	if !j2gort.EnterClinit(&clinitState_ClosureUnawareCodeWarningsGuard, &clinitDone_ClosureUnawareCodeWarningsGuard) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_ClosureUnawareCodeWarningsGuard, &clinitDone_ClosureUnawareCodeWarningsGuard)
	static_ClosureUnawareCodeWarningsGuard_DEFAULT_CLOSURE_UNAWARE_CODE_SUPPRESSIONS = NewDiagnosticGroup_2_string_DiagnosticType("closureUnawareCodeJSDocIncompatible", []*DiagnosticType{StaticGet_RhinoErrorReporter_TYPE_PARSE_ERROR(), StaticGet_RhinoErrorReporter_UNRECOGNIZED_TYPE_ERROR(), StaticGet_RhinoErrorReporter_JSDOC_MISSING_BRACES_WARNING(), StaticGet_RhinoErrorReporter_JSDOC_MISSING_TYPE_WARNING(), StaticGet_RhinoErrorReporter_UNNECESSARY_ESCAPE(), StaticGet_RhinoErrorReporter_BAD_JSDOC_ANNOTATION(), StaticGet_RhinoErrorReporter_UNSUPPORTED_BOUNDED_GENERIC_TYPES(), StaticGet_RhinoErrorReporter_BOUNDED_GENERIC_TYPE_ERROR()})
}

func (this_ *ClosureUnawareCodeWarningsGuard) initFields_ClosureUnawareCodeWarningsGuard() {
	_ = this_
}

// Java constructor: ClosureUnawareCodeWarningsGuard()
func NewClosureUnawareCodeWarningsGuard() *ClosureUnawareCodeWarningsGuard {
	clinit_ClosureUnawareCodeWarningsGuard()
	this_ := &ClosureUnawareCodeWarningsGuard{}
	InitClosureUnawareCodeWarningsGuard(this_)
	return this_
}

func InitClosureUnawareCodeWarningsGuard(this_ *ClosureUnawareCodeWarningsGuard) {
	InitWarningsGuard(&this_.WarningsGuardBase)
	this_.initFields_ClosureUnawareCodeWarningsGuard()
}

// Java method: ClosureUnawareCodeWarningsGuard.level(com.google.javascript.jscomp.JSError error)
func (this_ *ClosureUnawareCodeWarningsGuard) Level(error_ *JSError) *CheckLevel {
	_ = error_
	if (!(StaticGet_ClosureUnawareCodeWarningsGuard_DEFAULT_CLOSURE_UNAWARE_CODE_SUPPRESSIONS().Matches_1_JSError(error_))) {
		return nil
	}
	var node pkg_rhino.Node = error_.Node()
	_ = node
	if ((node == 0)) {
		return nil
	}
	if (node.IsClosureUnawareCode()) {
		return StaticGet_CheckLevel_OFF()
	}
	return nil
}

// Java method: ClosureUnawareCodeWarningsGuard.getPriority()
func (this_ *ClosureUnawareCodeWarningsGuard) GetPriority() int32 {
	return StaticGet_WarningsGuard_Priority_MAX().GetValue()
}

