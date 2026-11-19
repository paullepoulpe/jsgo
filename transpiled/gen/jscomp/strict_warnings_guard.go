/*
 * Copyright 2008 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/StrictWarningsGuard.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
)

// Java class: com.google.javascript.jscomp.StrictWarningsGuard
type StrictWarningsGuard struct {
	WarningsGuardBase
}

func (this_ *StrictWarningsGuard) BindSelf_(s any) {
	_ = s
}
func (this_ *StrictWarningsGuard) IsJavaNil_() bool { return this_ == nil }
func (*StrictWarningsGuard) JavaClassInfo_() (string, string) { return "com.google.javascript.jscomp.StrictWarningsGuard", "StrictWarningsGuard" }
func (this_ *StrictWarningsGuard) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.WarningsGuardBase.IdentityHash_()
}
func (*StrictWarningsGuard) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.StrictWarningsGuard", "com.google.javascript.jscomp.StrictWarningsGuard", func(v_ any) bool {
		return j2gort.InstanceOf[*StrictWarningsGuard](v_)
	})
}
func (*StrictWarningsGuard) Implements_jscomp_WarningsGuard() {}

func StaticGet_StrictWarningsGuard_SerialVersionUID() int64 { return int64(1) }
var static_StrictWarningsGuard_UNRAISABLE_WARNING *DiagnosticType
func StaticGet_StrictWarningsGuard_UNRAISABLE_WARNING() *DiagnosticType {
	clinit_StrictWarningsGuard()
	return static_StrictWarningsGuard_UNRAISABLE_WARNING
}
func StaticSet_StrictWarningsGuard_UNRAISABLE_WARNING(v *DiagnosticType) {
	clinit_StrictWarningsGuard()
	static_StrictWarningsGuard_UNRAISABLE_WARNING = v
}
var clinitDone_StrictWarningsGuard uint32
var clinitState_StrictWarningsGuard j2gort.ClinitState
func clinit_StrictWarningsGuard() {
	if atomic.LoadUint32(&clinitDone_StrictWarningsGuard) != 0 {
		return
	}
	clinitSlow_StrictWarningsGuard()
}

func clinitSlow_StrictWarningsGuard() {
	if !j2gort.EnterClinit(&clinitState_StrictWarningsGuard, &clinitDone_StrictWarningsGuard) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_StrictWarningsGuard, &clinitDone_StrictWarningsGuard)
	static_StrictWarningsGuard_UNRAISABLE_WARNING = DiagnosticType_Warning("JSC_UNRAISABLE_WARNING", "{0}")
}

func (this_ *StrictWarningsGuard) initFields_StrictWarningsGuard() {
	_ = this_
}

// Java constructor: StrictWarningsGuard()
func NewStrictWarningsGuard() *StrictWarningsGuard {
	clinit_StrictWarningsGuard()
	this_ := &StrictWarningsGuard{}
	InitStrictWarningsGuard(this_)
	return this_
}

func InitStrictWarningsGuard(this_ *StrictWarningsGuard) {
	InitWarningsGuard(&this_.WarningsGuardBase)
	this_.initFields_StrictWarningsGuard()
}

// Java method: StrictWarningsGuard.level(com.google.javascript.jscomp.JSError error)
func (this_ *StrictWarningsGuard) Level(error_ *JSError) *CheckLevel {
	_ = error_
	if ((error_.Type() == StaticGet_StrictWarningsGuard_UNRAISABLE_WARNING())) {
		return nil
	}
	if error_.DefaultLevel().IsOn() {
		return StaticGet_CheckLevel_ERROR()
	} else {
		return nil
	}
}

// Java method: StrictWarningsGuard.getPriority()
func (this_ *StrictWarningsGuard) GetPriority() int32 {
	return (StaticGet_WarningsGuard_Priority_STRICT()).FieldGet_value()
}

