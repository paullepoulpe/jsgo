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
// Port of com/google/javascript/jscomp/DiagnosticGroupWarningsGuard.java (generated mechanically by j2go).

package jscomp

import (
	pkg_base "github.com/paullepoulpe/jsgo/transpiled/gen/base"
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
)

// Java class: com.google.javascript.jscomp.DiagnosticGroupWarningsGuard
type DiagnosticGroupWarningsGuard struct {
	WarningsGuardBase
	group *DiagnosticGroup `gson:"group"`
	level *CheckLevel `gson:"level"`
}

func (this_ *DiagnosticGroupWarningsGuard) BindSelf_(s any) {
	_ = s
}
func (this_ *DiagnosticGroupWarningsGuard) IsJavaNil_() bool { return this_ == nil }
func (*DiagnosticGroupWarningsGuard) JavaClassInfo_() (string, string) { return "com.google.javascript.jscomp.DiagnosticGroupWarningsGuard", "DiagnosticGroupWarningsGuard" }
func (this_ *DiagnosticGroupWarningsGuard) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.WarningsGuardBase.IdentityHash_()
}
func (*DiagnosticGroupWarningsGuard) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.DiagnosticGroupWarningsGuard", "com.google.javascript.jscomp.DiagnosticGroupWarningsGuard", func(v_ any) bool {
		return j2gort.InstanceOf[*DiagnosticGroupWarningsGuard](v_)
	})
}
func (*DiagnosticGroupWarningsGuard) Implements_jscomp_WarningsGuard() {}

func StaticGet_DiagnosticGroupWarningsGuard_SerialVersionUID() int64 { return int64(1) }
func (this_ *DiagnosticGroupWarningsGuard) initFields_DiagnosticGroupWarningsGuard() {
	_ = this_
}

// Java constructor: DiagnosticGroupWarningsGuard(com.google.javascript.jscomp.DiagnosticGroup group, com.google.javascript.jscomp.CheckLevel level)
func NewDiagnosticGroupWarningsGuard(group *DiagnosticGroup, level *CheckLevel) *DiagnosticGroupWarningsGuard {
	this_ := &DiagnosticGroupWarningsGuard{}
	InitDiagnosticGroupWarningsGuard(this_, group, level)
	return this_
}

func InitDiagnosticGroupWarningsGuard(this_ *DiagnosticGroupWarningsGuard, group *DiagnosticGroup, level *CheckLevel) {
	_ = group
	_ = level
	InitWarningsGuard(&this_.WarningsGuardBase)
	this_.initFields_DiagnosticGroupWarningsGuard()
	this_.group = group
	this_.level = level
}

// Java method: DiagnosticGroupWarningsGuard.level(com.google.javascript.jscomp.JSError error)
func (this_ *DiagnosticGroupWarningsGuard) Level(error_ *JSError) *CheckLevel {
	_ = error_
	if this_.group.Matches_1_JSError(error_) {
		return this_.level
	} else {
		return nil
	}
}

// Java method: DiagnosticGroupWarningsGuard.mustRunChecks(com.google.javascript.jscomp.DiagnosticGroup otherGroup)
func (this_ *DiagnosticGroupWarningsGuard) MustRunChecks(otherGroup *DiagnosticGroup) pkg_base.Tri {
	_ = otherGroup
	if (this_.level.IsOn()) {
		if otherGroup.GetTypes().Stream().AnyMatch(func() j2gort.Ext_java_util_function_Predicate { bound_1 := this_.group; return j2gort.NewLambda_Ext_java_util_function_Predicate(func(mrp0_ any) bool { return bound_1.Matches_1_DiagnosticType(j2gort.CastPtr[DiagnosticType](mrp0_)) }) }()) {
			return pkg_base.StaticGet_Tri_TRUE()
		} else {
			return pkg_base.StaticGet_Tri_UNKNOWN()
		}
	} else {
		if otherGroup.GetTypes().Stream().AllMatch(func() j2gort.Ext_java_util_function_Predicate { bound_2 := this_.group; return j2gort.NewLambda_Ext_java_util_function_Predicate(func(mrp0_ any) bool { return bound_2.Matches_1_DiagnosticType(j2gort.CastPtr[DiagnosticType](mrp0_)) }) }()) {
			return pkg_base.StaticGet_Tri_FALSE()
		} else {
			return pkg_base.StaticGet_Tri_UNKNOWN()
		}
	}
}

// Java method: DiagnosticGroupWarningsGuard.toString()
func (this_ *DiagnosticGroupWarningsGuard) String() string {
	return (j2gort.StrVal(this_.group) + "(" + j2gort.StrVal(this_.level) + ")")
}

