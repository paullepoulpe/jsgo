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
// Port of com/google/javascript/jscomp/ijs/CheckTypeSummaryWarningsGuard.java (generated mechanically by j2go).

package jscomp

import (
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.jscomp.ijs.CheckTypeSummaryWarningsGuard
type CheckTypeSummaryWarningsGuard struct {
	WarningsGuardBase
	level *CheckLevel `gson:"level"`
}

func (this_ *CheckTypeSummaryWarningsGuard) BindSelf_(s any) {
	_ = s
}
func (this_ *CheckTypeSummaryWarningsGuard) IsJavaNil_() bool { return this_ == nil }
func (*CheckTypeSummaryWarningsGuard) JavaClassInfo_() (string, string) { return "com.google.javascript.jscomp.ijs.CheckTypeSummaryWarningsGuard", "CheckTypeSummaryWarningsGuard" }
func (this_ *CheckTypeSummaryWarningsGuard) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.WarningsGuardBase.IdentityHash_()
}
func (*CheckTypeSummaryWarningsGuard) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.ijs.CheckTypeSummaryWarningsGuard", "com.google.javascript.jscomp.ijs.CheckTypeSummaryWarningsGuard", func(v_ any) bool {
		return j2gort.InstanceOf[*CheckTypeSummaryWarningsGuard](v_)
	})
}
func (*CheckTypeSummaryWarningsGuard) Implements_jscomp_WarningsGuard() {}

func (this_ *CheckTypeSummaryWarningsGuard) initFields_CheckTypeSummaryWarningsGuard() {
	_ = this_
}

// Java constructor: CheckTypeSummaryWarningsGuard(com.google.javascript.jscomp.CheckLevel level)
func NewCheckTypeSummaryWarningsGuard(level *CheckLevel) *CheckTypeSummaryWarningsGuard {
	this_ := &CheckTypeSummaryWarningsGuard{}
	InitCheckTypeSummaryWarningsGuard(this_, level)
	return this_
}

func InitCheckTypeSummaryWarningsGuard(this_ *CheckTypeSummaryWarningsGuard, level *CheckLevel) {
	_ = level
	InitWarningsGuard(&this_.WarningsGuardBase)
	this_.initFields_CheckTypeSummaryWarningsGuard()
	this_.level = level
}

// Java method: CheckTypeSummaryWarningsGuard.level(com.google.javascript.jscomp.JSError error)
func (this_ *CheckTypeSummaryWarningsGuard) Level(error_ *JSError) *CheckLevel {
	_ = error_
	j2gort.CheckNotNullPtr(error_)
	if (!(this_.priv_CheckTypeSummaryWarningsGuard_shouldDemoteIfTypeSummary_6391(error_))) {
		return nil
	}
	if (this_.priv_CheckTypeSummaryWarningsGuard_inTypeSummary_6390(error_)) {
		return this_.level
	}
	return nil
}

// Java method: CheckTypeSummaryWarningsGuard.getPriority()
func (this_ *CheckTypeSummaryWarningsGuard) GetPriority() int32 {
	return StaticGet_WarningsGuard_Priority_SUPPRESS_BY_ALLOWLIST().GetValue()
}

// Java method: CheckTypeSummaryWarningsGuard.inTypeSummary(com.google.javascript.jscomp.JSError error)
func (this_ *CheckTypeSummaryWarningsGuard) priv_CheckTypeSummaryWarningsGuard_inTypeSummary_6390(error_ *JSError) bool {
	_ = error_
	return ((error_.SourceName() != j2gort.NullString) && javalang.StringEndsWith(j2gort.CheckNotNullString(error_.SourceName()), ".i.js"))
}

// Java method: CheckTypeSummaryWarningsGuard.shouldDemoteIfTypeSummary(com.google.javascript.jscomp.JSError error)
func (this_ *CheckTypeSummaryWarningsGuard) priv_CheckTypeSummaryWarningsGuard_shouldDemoteIfTypeSummary_6391(error_ *JSError) bool {
	_ = error_
	return !(error_.Type().Equals(j2gort.PtrToAny(StaticGet_ClosurePrimitiveErrors_DUPLICATE_NAMESPACE_AND_MODULE())))
}

