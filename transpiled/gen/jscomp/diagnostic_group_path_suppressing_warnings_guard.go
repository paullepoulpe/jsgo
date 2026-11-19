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
// Port of com/google/javascript/jscomp/DiagnosticGroupPathSuppressingWarningsGuard.java (generated mechanically by j2go).

package jscomp

import (
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
)

// Java class: com.google.javascript.jscomp.DiagnosticGroupPathSuppressingWarningsGuard
type DiagnosticGroupPathSuppressingWarningsGuard struct {
	WarningsGuardBase
	group *DiagnosticGroup `gson:"group"`
	part string `gson:"part"`
}

func (this_ *DiagnosticGroupPathSuppressingWarningsGuard) BindSelf_(s any) {
	_ = s
}
func (this_ *DiagnosticGroupPathSuppressingWarningsGuard) IsJavaNil_() bool { return this_ == nil }
func (*DiagnosticGroupPathSuppressingWarningsGuard) JavaClassInfo_() (string, string) { return "com.google.javascript.jscomp.DiagnosticGroupPathSuppressingWarningsGuard", "DiagnosticGroupPathSuppressingWarningsGuard" }
func (this_ *DiagnosticGroupPathSuppressingWarningsGuard) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.WarningsGuardBase.IdentityHash_()
}
func (*DiagnosticGroupPathSuppressingWarningsGuard) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.DiagnosticGroupPathSuppressingWarningsGuard", "com.google.javascript.jscomp.DiagnosticGroupPathSuppressingWarningsGuard", func(v_ any) bool {
		return j2gort.InstanceOf[*DiagnosticGroupPathSuppressingWarningsGuard](v_)
	})
}
func (*DiagnosticGroupPathSuppressingWarningsGuard) Implements_jscomp_WarningsGuard() {}

func (this_ *DiagnosticGroupPathSuppressingWarningsGuard) initFields_DiagnosticGroupPathSuppressingWarningsGuard() {
	_ = this_
	this_.part = j2gort.NullString
}

// Java constructor: DiagnosticGroupPathSuppressingWarningsGuard(com.google.javascript.jscomp.DiagnosticGroup group, java.lang.String part)
func NewDiagnosticGroupPathSuppressingWarningsGuard(group *DiagnosticGroup, part string) *DiagnosticGroupPathSuppressingWarningsGuard {
	this_ := &DiagnosticGroupPathSuppressingWarningsGuard{}
	InitDiagnosticGroupPathSuppressingWarningsGuard(this_, group, part)
	return this_
}

func InitDiagnosticGroupPathSuppressingWarningsGuard(this_ *DiagnosticGroupPathSuppressingWarningsGuard, group *DiagnosticGroup, part string) {
	_ = group
	_ = part
	InitWarningsGuard(&this_.WarningsGuardBase)
	this_.initFields_DiagnosticGroupPathSuppressingWarningsGuard()
	this_.group = group
	this_.part = part
}

// Java method: DiagnosticGroupPathSuppressingWarningsGuard.level(com.google.javascript.jscomp.JSError error)
func (this_ *DiagnosticGroupPathSuppressingWarningsGuard) Level(error_ *JSError) *CheckLevel {
	_ = error_
	if (((error_.SourceName() == j2gort.NullString) || !((j2gort.JStrIndexOf(error_.SourceName(), j2gort.CheckNotNullString(this_.part)) >= 0)))) {
		return nil
	}
	if this_.group.Matches_1_JSError(error_) {
		return StaticGet_CheckLevel_OFF()
	} else {
		return nil
	}
}

// Java method: DiagnosticGroupPathSuppressingWarningsGuard.toString()
func (this_ *DiagnosticGroupPathSuppressingWarningsGuard) String() string {
	return (j2gort.StrVal(this_.group) + "(" + j2gort.StrValString(this_.part) + ")")
}

