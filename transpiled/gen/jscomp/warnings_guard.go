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
// Port of com/google/javascript/jscomp/WarningsGuard.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	pkg_base "github.com/paullepoulpe/jsgo/transpiled/gen/base"
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	"github.com/paullepoulpe/jsgo/transpiled/shims/gson"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java polymorphic class interface: com.google.javascript.jscomp.WarningsGuard
type WarningsGuard interface {
	Implements_jscomp_WarningsGuard()
	IsJavaNil_() bool
	GetWarningsGuardBase() *WarningsGuardBase
	Level(*JSError) *CheckLevel
	MustRunChecks(*DiagnosticGroup) pkg_base.Tri
	GetPriority() int32
}
func ToIface_WarningsGuard[T interface {
	WarningsGuard
	comparable
}](v T) WarningsGuard {
	var z T
	if v == z {
		return nil
	}
	return v
}

func Upcast_WarningsGuard[T WarningsGuard](v T) WarningsGuard {
	var vi WarningsGuard = v
	if vi == nil || vi.IsJavaNil_() {
		return nil
	}
	return vi
}

// Java class: com.google.javascript.jscomp.WarningsGuard
type WarningsGuardBase struct {
	idHash_ atomic.Int32
}

func (this_ *WarningsGuardBase) BindSelf_(s any) {
	_ = s
}
func (this_ *WarningsGuardBase) IsJavaNil_() bool { return this_ == nil }
func (*WarningsGuardBase) JavaClassInfo_() (string, string) { return "com.google.javascript.jscomp.WarningsGuard", "WarningsGuard" }
func (this_ *WarningsGuardBase) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.WarningsGuard", "com.google.javascript.jscomp.WarningsGuard", func(v_ any) bool {
		return j2gort.InstanceOf[WarningsGuard](v_)
	})
}
func (b *WarningsGuardBase) GetWarningsGuardBase() *WarningsGuardBase { return b }
func (*WarningsGuardBase) Implements_jscomp_WarningsGuard() {}

func (this_ *WarningsGuardBase) initFields_WarningsGuard() {
	_ = this_
}

func InitWarningsGuard(this_ *WarningsGuardBase) {
	this_.initFields_WarningsGuard()
}

// Java method: WarningsGuard.level(com.google.javascript.jscomp.JSError error)
func (this_ *WarningsGuardBase) Level(error_ *JSError) *CheckLevel {
	_ = error_
	panic("abstract or native method invoked")
}

// Java method: WarningsGuard.mustRunChecks(com.google.javascript.jscomp.DiagnosticGroup group)
func (this_ *WarningsGuardBase) MustRunChecks(group *DiagnosticGroup) pkg_base.Tri {
	_ = group
	return pkg_base.StaticGet_Tri_UNKNOWN()
}

// Java method: WarningsGuard.getPriority()
func (this_ *WarningsGuardBase) GetPriority() int32 {
	return (StaticGet_WarningsGuard_Priority_DEFAULT()).FieldGet_value()
}

// Java enum: com.google.javascript.jscomp.WarningsGuard.Priority
type WarningsGuard_Priority uint8

const (
	enum_WarningsGuard_Priority_MAX WarningsGuard_Priority = 1
	enum_WarningsGuard_Priority_MIN WarningsGuard_Priority = 2
	enum_WarningsGuard_Priority_STRICT WarningsGuard_Priority = 3
	enum_WarningsGuard_Priority_DEFAULT WarningsGuard_Priority = 4
	enum_WarningsGuard_Priority_SUPPRESS_BY_ALLOWLIST WarningsGuard_Priority = 5
	enum_WarningsGuard_Priority_SUPPRESS_DOC WarningsGuard_Priority = 6
	enum_WarningsGuard_Priority_FILTER_BY_PATH WarningsGuard_Priority = 7
)

var enumBaseTable_WarningsGuard_Priority = [...]j2gort.EnumBase{
	{Name_: "MAX", Ordinal_: 0},
	{Name_: "MIN", Ordinal_: 1},
	{Name_: "STRICT", Ordinal_: 2},
	{Name_: "DEFAULT", Ordinal_: 3},
	{Name_: "SUPPRESS_BY_ALLOWLIST", Ordinal_: 4},
	{Name_: "SUPPRESS_DOC", Ordinal_: 5},
	{Name_: "FILTER_BY_PATH", Ordinal_: 6},
}
var enumValues_WarningsGuard_Priority = []WarningsGuard_Priority{1, 2, 3, 4, 5, 6, 7}

var enumField_WarningsGuard_Priority_value = [...]int32{
	int32(1),
	int32(100),
	int32(100),
	int32(50),
	int32(40),
	int32(20),
	int32(1),
}
func (this_ WarningsGuard_Priority) FieldGet_value() int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumField_WarningsGuard_Priority_value[this_-1]
}

func StaticGet_WarningsGuard_Priority_MAX() WarningsGuard_Priority {
	return enum_WarningsGuard_Priority_MAX
}
func StaticGet_WarningsGuard_Priority_MIN() WarningsGuard_Priority {
	return enum_WarningsGuard_Priority_MIN
}
func StaticGet_WarningsGuard_Priority_STRICT() WarningsGuard_Priority {
	return enum_WarningsGuard_Priority_STRICT
}
func StaticGet_WarningsGuard_Priority_DEFAULT() WarningsGuard_Priority {
	return enum_WarningsGuard_Priority_DEFAULT
}
func StaticGet_WarningsGuard_Priority_SUPPRESS_BY_ALLOWLIST() WarningsGuard_Priority {
	return enum_WarningsGuard_Priority_SUPPRESS_BY_ALLOWLIST
}
func StaticGet_WarningsGuard_Priority_SUPPRESS_DOC() WarningsGuard_Priority {
	return enum_WarningsGuard_Priority_SUPPRESS_DOC
}
func StaticGet_WarningsGuard_Priority_FILTER_BY_PATH() WarningsGuard_Priority {
	return enum_WarningsGuard_Priority_FILTER_BY_PATH
}

func WarningsGuard_Priority_Values() []WarningsGuard_Priority {
	return j2gort.CloneSlice(enumValues_WarningsGuard_Priority)
}

func WarningsGuard_Priority_ValueOf(name string) WarningsGuard_Priority {
	if name == j2gort.NullString {
		j2gort.ThrowNPE()
	}
	for i_ := range enumBaseTable_WarningsGuard_Priority {
		if enumBaseTable_WarningsGuard_Priority[i_].Name_ == name {
			return WarningsGuard_Priority(i_ + 1)
		}
	}
	panic(javalang.NewIllegalArgumentException("No enum constant com.google.javascript.jscomp.WarningsGuard.Priority." + name))
}

// Java method: Priority.getValue()
func (this_ WarningsGuard_Priority) GetValue() int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return this_.FieldGet_value()
}

func (this_ WarningsGuard_Priority) IsJavaNil_() bool { return this_ == 0 }
func (WarningsGuard_Priority) HasIdentityEquals_() {}
func (WarningsGuard_Priority) Implements_j2gort_Ext_java_lang_Enum() {}
func (WarningsGuard_Priority) Implements_j2gort_Ext_java_lang_constant_Constable() {}
func (this_ WarningsGuard_Priority) GetExt_java_lang_EnumBase() *j2gort.Ext_java_lang_EnumBase { return &j2gort.Ext_java_lang_EnumBase{Core_: this_.GetEnumBase()} }
func (WarningsGuard_Priority) JavaClassInfo_() (string, string) { return "com.google.javascript.jscomp.WarningsGuard$Priority", "Priority" }
func (this_ WarningsGuard_Priority) GetEnumBase() *j2gort.EnumBase {
	if this_ == 0 {
		return nil
	}
	return &enumBaseTable_WarningsGuard_Priority[this_-1]
}
func (this_ WarningsGuard_Priority) Ordinal() int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return int32(this_ - 1)
}
func (this_ WarningsGuard_Priority) EnumOrdinal_() int32 { return int32(this_) - 1 }
func (this_ WarningsGuard_Priority) Name() string {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_WarningsGuard_Priority[this_-1].Name_
}
func (this_ WarningsGuard_Priority) String() string {
	if this_ == 0 {
		return "null"
	}
	return enumBaseTable_WarningsGuard_Priority[this_-1].Name_
}
func (this_ WarningsGuard_Priority) Equals(other_ any) bool {
	o_, ok_ := other_.(WarningsGuard_Priority)
	return this_ != 0 && ok_ && this_ == o_
}
func (this_ WarningsGuard_Priority) HashCode() int32 {
	if this_ == 0 {
		return 0
	}
	return enumBaseTable_WarningsGuard_Priority[this_-1].HashCode()
}
func (this_ WarningsGuard_Priority) CompareTo(other_ any) int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_WarningsGuard_Priority[this_-1].CompareTo(other_)
}
func (this_ WarningsGuard_Priority) WriteGson(out_ *gson.JsonWriter) {
	if this_ == 0 {
		out_.NullValue()
		return
	}
	out_.ValueString(enumBaseTable_WarningsGuard_Priority[this_-1].Name_)
}
func init() {
	vals_ := []any{WarningsGuard_Priority(1), WarningsGuard_Priority(2), WarningsGuard_Priority(3), WarningsGuard_Priority(4), WarningsGuard_Priority(5), WarningsGuard_Priority(6), WarningsGuard_Priority(7)}
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.WarningsGuard.Priority", vals_)
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.WarningsGuard$Priority", vals_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.WarningsGuard.Priority", "com.google.javascript.jscomp.WarningsGuard$Priority", func(v_ any) bool {
		return j2gort.InstanceOf[WarningsGuard_Priority](v_)
	})
}

