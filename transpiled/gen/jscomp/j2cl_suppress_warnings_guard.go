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
// Port of com/google/javascript/jscomp/J2clSuppressWarningsGuard.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
)

// Java class: com.google.javascript.jscomp.J2clSuppressWarningsGuard
type J2clSuppressWarningsGuard struct {
	WarningsGuardBase
}

func (this_ *J2clSuppressWarningsGuard) BindSelf_(s any) {
	_ = s
}
func (this_ *J2clSuppressWarningsGuard) IsJavaNil_() bool { return this_ == nil }
func (*J2clSuppressWarningsGuard) JavaClassInfo_() (string, string) { return "com.google.javascript.jscomp.J2clSuppressWarningsGuard", "J2clSuppressWarningsGuard" }
func (this_ *J2clSuppressWarningsGuard) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.WarningsGuardBase.IdentityHash_()
}
func (*J2clSuppressWarningsGuard) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.J2clSuppressWarningsGuard", "com.google.javascript.jscomp.J2clSuppressWarningsGuard", func(v_ any) bool {
		return j2gort.InstanceOf[*J2clSuppressWarningsGuard](v_)
	})
}
func (*J2clSuppressWarningsGuard) Implements_jscomp_WarningsGuard() {}

var static_J2clSuppressWarningsGuard_DEFAULT_J2CL_SUPRRESIONS *DiagnosticGroup
func StaticGet_J2clSuppressWarningsGuard_DEFAULT_J2CL_SUPRRESIONS() *DiagnosticGroup {
	clinit_J2clSuppressWarningsGuard()
	return static_J2clSuppressWarningsGuard_DEFAULT_J2CL_SUPRRESIONS
}
func StaticSet_J2clSuppressWarningsGuard_DEFAULT_J2CL_SUPRRESIONS(v *DiagnosticGroup) {
	clinit_J2clSuppressWarningsGuard()
	static_J2clSuppressWarningsGuard_DEFAULT_J2CL_SUPRRESIONS = v
}
var clinitDone_J2clSuppressWarningsGuard uint32
var clinitState_J2clSuppressWarningsGuard j2gort.ClinitState
func clinit_J2clSuppressWarningsGuard() {
	if atomic.LoadUint32(&clinitDone_J2clSuppressWarningsGuard) != 0 {
		return
	}
	clinitSlow_J2clSuppressWarningsGuard()
}

func clinitSlow_J2clSuppressWarningsGuard() {
	if !j2gort.EnterClinit(&clinitState_J2clSuppressWarningsGuard, &clinitDone_J2clSuppressWarningsGuard) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_J2clSuppressWarningsGuard, &clinitDone_J2clSuppressWarningsGuard)
	static_J2clSuppressWarningsGuard_DEFAULT_J2CL_SUPRRESIONS = NewDiagnosticGroup_2_string_DiagnosticGroup("j2clIncomaptible", []*DiagnosticGroup{StaticGet_DiagnosticGroups_CHECK_STATIC_OVERRIDES(), StaticGet_DiagnosticGroups_CHECK_USELESS_CODE(), StaticGet_DiagnosticGroups_CONST(), StaticGet_DiagnosticGroups_EXTRA_REQUIRE(), StaticGet_DiagnosticGroups_INVALID_CASTS(), StaticGet_DiagnosticGroups_LATE_PROVIDE(), StaticGet_DiagnosticGroups_MISSING_OVERRIDE(), StaticGet_DiagnosticGroups_MISSING_REQUIRE(), StaticGet_DiagnosticGroups_STRICT_MODULE_DEP_CHECK(), StaticGet_DiagnosticGroups_SUSPICIOUS_CODE(), StaticGet_DiagnosticGroups_UNUSED_LOCAL_VARIABLE(), StaticGet_DiagnosticGroups_STRICT_MISSING_PROPERTIES(), DiagnosticGroups_ForName("transitionalSuspiciousCodeWarnings")})
}

func (this_ *J2clSuppressWarningsGuard) initFields_J2clSuppressWarningsGuard() {
	_ = this_
}

// Java constructor: J2clSuppressWarningsGuard()
func NewJ2clSuppressWarningsGuard() *J2clSuppressWarningsGuard {
	clinit_J2clSuppressWarningsGuard()
	this_ := &J2clSuppressWarningsGuard{}
	InitJ2clSuppressWarningsGuard(this_)
	return this_
}

func InitJ2clSuppressWarningsGuard(this_ *J2clSuppressWarningsGuard) {
	InitWarningsGuard(&this_.WarningsGuardBase)
	this_.initFields_J2clSuppressWarningsGuard()
}

// Java method: J2clSuppressWarningsGuard.level(com.google.javascript.jscomp.JSError error)
func (this_ *J2clSuppressWarningsGuard) Level(error_ *JSError) *CheckLevel {
	_ = error_
	if (!(J2clSourceUtils_IsJ2clSource_1_string(error_.SourceName()))) {
		return nil
	}
	if ((StaticGet_DiagnosticGroups_DEPRECATED().Matches_1_JSError(error_) && (j2gort.JStrIndexOf(error_.Description(), "forwardDeclare") >= 0))) {
		return StaticGet_CheckLevel_OFF()
	}
	if StaticGet_J2clSuppressWarningsGuard_DEFAULT_J2CL_SUPRRESIONS().Matches_1_JSError(error_) {
		return StaticGet_CheckLevel_OFF()
	} else {
		return nil
	}
}

// Java method: J2clSuppressWarningsGuard.getPriority()
func (this_ *J2clSuppressWarningsGuard) GetPriority() int32 {
	return StaticGet_WarningsGuard_Priority_MAX().GetValue()
}

