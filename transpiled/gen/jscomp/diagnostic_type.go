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
// Port of com/google/javascript/jscomp/DiagnosticType.java (generated mechanically by j2go).

package jscomp

import (
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.jscomp.DiagnosticType
type DiagnosticType struct {
	key string `gson:"key"`
	format string `gson:"format"`
	level *CheckLevel `gson:"level"`
}

func (this_ *DiagnosticType) BindSelf_(s any) {
	_ = s
}
func (this_ *DiagnosticType) IsJavaNil_() bool { return this_ == nil }
func (*DiagnosticType) JavaClassInfo_() (string, string) { return "com.google.javascript.jscomp.DiagnosticType", "DiagnosticType" }
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.DiagnosticType", "com.google.javascript.jscomp.DiagnosticType", func(v_ any) bool {
		return j2gort.InstanceOf[*DiagnosticType](v_)
	})
}

func StaticGet_DiagnosticType_SerialVersionUID() int64 { return int64(1) }
func (this_ *DiagnosticType) initFields_DiagnosticType() {
	_ = this_
	this_.key = j2gort.NullString
	this_.format = j2gort.NullString
}

// Java constructor: DiagnosticType(java.lang.String key, com.google.javascript.jscomp.CheckLevel level, java.lang.String format)
func NewDiagnosticType(key string, level *CheckLevel, format string) *DiagnosticType {
	this_ := &DiagnosticType{}
	InitDiagnosticType(this_, key, level, format)
	return this_
}

func InitDiagnosticType(this_ *DiagnosticType, key string, level *CheckLevel, format string) {
	_ = key
	_ = level
	_ = format
	this_.initFields_DiagnosticType()
	this_.key = key
	this_.level = level
	this_.format = format
}

// Java static method: DiagnosticType.error(java.lang.String name, java.lang.String descriptionFormat)
func DiagnosticType_Error(name string, descriptionFormat string) *DiagnosticType {
	_ = name
	_ = descriptionFormat
	return DiagnosticType_Make(name, StaticGet_CheckLevel_ERROR(), descriptionFormat)
}

// Java static method: DiagnosticType.warning(java.lang.String name, java.lang.String descriptionFormat)
func DiagnosticType_Warning(name string, descriptionFormat string) *DiagnosticType {
	_ = name
	_ = descriptionFormat
	return DiagnosticType_Make(name, StaticGet_CheckLevel_WARNING(), descriptionFormat)
}

// Java static method: DiagnosticType.disabled(java.lang.String name, java.lang.String descriptionFormat)
func DiagnosticType_Disabled(name string, descriptionFormat string) *DiagnosticType {
	_ = name
	_ = descriptionFormat
	return DiagnosticType_Make(name, StaticGet_CheckLevel_OFF(), descriptionFormat)
}

// Java static method: DiagnosticType.make(java.lang.String name, com.google.javascript.jscomp.CheckLevel level, java.lang.String descriptionFormat)
func DiagnosticType_Make(name string, level *CheckLevel, descriptionFormat string) *DiagnosticType {
	_ = name
	_ = level
	_ = descriptionFormat
	return NewDiagnosticType(name, level, descriptionFormat)
}

// Java method: DiagnosticType.format(java.lang.String[] arguments)
func (this_ *DiagnosticType) Format(arguments []string) string {
	_ = arguments
	return Platform_FormatMessage(this_.format, arguments)
}

// Java method: DiagnosticType.equals(java.lang.Object type)
func (this_ *DiagnosticType) Equals(arg_1 any) bool {
	type_ := arg_1
	_ = type_
	var diagnosticType *DiagnosticType
	_ = diagnosticType
	var ok_2 bool
	diagnosticType, ok_2 = j2gort.TryCast[*DiagnosticType](type_)
	return ok_2 && (func() bool { sr_3 := diagnosticType.key; se_4 := this_.key; return j2gort.CheckNotNullString(sr_3) == se_4 }())
}

// Java method: DiagnosticType.hashCode()
func (this_ *DiagnosticType) HashCode() int32 {
	return javalang.StringHashCode(j2gort.CheckNotNullString(this_.key))
}

// Java method: DiagnosticType.compareTo(com.google.javascript.jscomp.DiagnosticType diagnosticType)
func (this_ *DiagnosticType) CompareTo(arg_5 any) int32 {
	diagnosticType := j2gort.CastPtr[DiagnosticType](arg_5)
	_ = diagnosticType
	return j2gort.JStrCompareTo(this_.key, diagnosticType.key)
}

// Java method: DiagnosticType.toString()
func (this_ *DiagnosticType) String() string {
	return (j2gort.StrValString(this_.key) + ": " + j2gort.StrValString(this_.format))
}

