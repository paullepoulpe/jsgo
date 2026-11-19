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
// Port of com/google/javascript/jscomp/ShowByPathWarningsGuard.java (generated mechanically by j2go).

package jscomp

import (
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	"github.com/paullepoulpe/jsgo/transpiled/shims/gson"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.jscomp.ShowByPathWarningsGuard
type ShowByPathWarningsGuard struct {
	WarningsGuardBase
	warningsGuard *ByPathWarningsGuard `gson:"warningsGuard"`
}

func (this_ *ShowByPathWarningsGuard) BindSelf_(s any) {
	_ = s
}
func (this_ *ShowByPathWarningsGuard) IsJavaNil_() bool { return this_ == nil }
func (*ShowByPathWarningsGuard) JavaClassInfo_() (string, string) { return "com.google.javascript.jscomp.ShowByPathWarningsGuard", "ShowByPathWarningsGuard" }
func (this_ *ShowByPathWarningsGuard) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.WarningsGuardBase.IdentityHash_()
}
func (*ShowByPathWarningsGuard) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.ShowByPathWarningsGuard", "com.google.javascript.jscomp.ShowByPathWarningsGuard", func(v_ any) bool {
		return j2gort.InstanceOf[*ShowByPathWarningsGuard](v_)
	})
}
func (*ShowByPathWarningsGuard) Implements_jscomp_WarningsGuard() {}

func StaticGet_ShowByPathWarningsGuard_SerialVersionUID() int64 { return int64(1) }
func (this_ *ShowByPathWarningsGuard) initFields_ShowByPathWarningsGuard() {
	_ = this_
}

// Java constructor: ShowByPathWarningsGuard(java.lang.String checkWarningsOnlyForPath)
func NewShowByPathWarningsGuard_1_string(checkWarningsOnlyForPath string) *ShowByPathWarningsGuard {
	this_ := &ShowByPathWarningsGuard{}
	InitShowByPathWarningsGuard_1_string(this_, checkWarningsOnlyForPath)
	return this_
}

func InitShowByPathWarningsGuard_1_string(this_ *ShowByPathWarningsGuard, checkWarningsOnlyForPath string) {
	_ = checkWarningsOnlyForPath
	InitShowByPathWarningsGuard_2_string_ShowByPathWarningsGuard_ShowType(this_, checkWarningsOnlyForPath, StaticGet_ShowByPathWarningsGuard_ShowType_INCLUDE())
}

// Java constructor: ShowByPathWarningsGuard(java.lang.String[] checkWarningsOnlyForPath)
func NewShowByPathWarningsGuard_1_slice_string(checkWarningsOnlyForPath []string) *ShowByPathWarningsGuard {
	this_ := &ShowByPathWarningsGuard{}
	InitShowByPathWarningsGuard_1_slice_string(this_, checkWarningsOnlyForPath)
	return this_
}

func InitShowByPathWarningsGuard_1_slice_string(this_ *ShowByPathWarningsGuard, checkWarningsOnlyForPath []string) {
	_ = checkWarningsOnlyForPath
	InitShowByPathWarningsGuard_2_slice_string_ShowByPathWarningsGuard_ShowType(this_, checkWarningsOnlyForPath, StaticGet_ShowByPathWarningsGuard_ShowType_INCLUDE())
}

// Java constructor: ShowByPathWarningsGuard(java.lang.String path, com.google.javascript.jscomp.ShowByPathWarningsGuard.ShowType showType)
func NewShowByPathWarningsGuard_2_string_ShowByPathWarningsGuard_ShowType(path string, showType ShowByPathWarningsGuard_ShowType) *ShowByPathWarningsGuard {
	this_ := &ShowByPathWarningsGuard{}
	InitShowByPathWarningsGuard_2_string_ShowByPathWarningsGuard_ShowType(this_, path, showType)
	return this_
}

func InitShowByPathWarningsGuard_2_string_ShowByPathWarningsGuard_ShowType(this_ *ShowByPathWarningsGuard, path string, showType ShowByPathWarningsGuard_ShowType) {
	_ = path
	_ = showType
	InitShowByPathWarningsGuard_2_slice_string_ShowByPathWarningsGuard_ShowType(this_, []string{path}, showType)
}

// Java constructor: ShowByPathWarningsGuard(java.lang.String[] paths, com.google.javascript.jscomp.ShowByPathWarningsGuard.ShowType showType)
func NewShowByPathWarningsGuard_2_slice_string_ShowByPathWarningsGuard_ShowType(paths []string, showType ShowByPathWarningsGuard_ShowType) *ShowByPathWarningsGuard {
	this_ := &ShowByPathWarningsGuard{}
	InitShowByPathWarningsGuard_2_slice_string_ShowByPathWarningsGuard_ShowType(this_, paths, showType)
	return this_
}

func InitShowByPathWarningsGuard_2_slice_string_ShowByPathWarningsGuard_ShowType(this_ *ShowByPathWarningsGuard, paths []string, showType ShowByPathWarningsGuard_ShowType) {
	_ = paths
	_ = showType
	InitWarningsGuard(&this_.WarningsGuardBase)
	this_.initFields_ShowByPathWarningsGuard()
	if !((paths != nil)) { j2gort.CheckArgument(false) }
	if !((showType != 0)) { j2gort.CheckArgument(false) }
	var pathList j2gort.Ext_com_google_common_collect_ImmutableList = j2gort.Ext_com_google_common_collect_ImmutableList_CopyOf_1_slice_any(func() []any { arr_1 := paths; if arr_1 == nil { return nil }; out_2 := make([]any, len(arr_1)); for i_, el_3 := range arr_1 { out_2[i_] = j2gort.BoxString(el_3) }; return out_2 }())
	_ = pathList
	if ((showType == StaticGet_ShowByPathWarningsGuard_ShowType_INCLUDE())) {
		this_.warningsGuard = ByPathWarningsGuard_ExceptPath(j2gort.AsIface[j2gort.Ext_java_util_List](pathList), StaticGet_CheckLevel_OFF())
	} else {
		this_.warningsGuard = ByPathWarningsGuard_ForPath(j2gort.AsIface[j2gort.Ext_java_util_List](pathList), StaticGet_CheckLevel_OFF())
	}
}

// Java method: ShowByPathWarningsGuard.level(com.google.javascript.jscomp.JSError error)
func (this_ *ShowByPathWarningsGuard) Level(error_ *JSError) *CheckLevel {
	_ = error_
	return this_.warningsGuard.Level(error_)
}

// Java method: ShowByPathWarningsGuard.getPriority()
func (this_ *ShowByPathWarningsGuard) GetPriority() int32 {
	return this_.warningsGuard.GetPriority()
}

// Java enum: com.google.javascript.jscomp.ShowByPathWarningsGuard.ShowType
type ShowByPathWarningsGuard_ShowType uint8

const (
	enum_ShowByPathWarningsGuard_ShowType_INCLUDE ShowByPathWarningsGuard_ShowType = 1
	enum_ShowByPathWarningsGuard_ShowType_EXCLUDE ShowByPathWarningsGuard_ShowType = 2
)

var enumBaseTable_ShowByPathWarningsGuard_ShowType = [...]j2gort.EnumBase{
	{Name_: "INCLUDE", Ordinal_: 0},
	{Name_: "EXCLUDE", Ordinal_: 1},
}
var enumValues_ShowByPathWarningsGuard_ShowType = []ShowByPathWarningsGuard_ShowType{1, 2}

func StaticGet_ShowByPathWarningsGuard_ShowType_INCLUDE() ShowByPathWarningsGuard_ShowType {
	return enum_ShowByPathWarningsGuard_ShowType_INCLUDE
}
func StaticGet_ShowByPathWarningsGuard_ShowType_EXCLUDE() ShowByPathWarningsGuard_ShowType {
	return enum_ShowByPathWarningsGuard_ShowType_EXCLUDE
}

func ShowByPathWarningsGuard_ShowType_Values() []ShowByPathWarningsGuard_ShowType {
	return j2gort.CloneSlice(enumValues_ShowByPathWarningsGuard_ShowType)
}

func ShowByPathWarningsGuard_ShowType_ValueOf(name string) ShowByPathWarningsGuard_ShowType {
	if name == j2gort.NullString {
		j2gort.ThrowNPE()
	}
	for i_ := range enumBaseTable_ShowByPathWarningsGuard_ShowType {
		if enumBaseTable_ShowByPathWarningsGuard_ShowType[i_].Name_ == name {
			return ShowByPathWarningsGuard_ShowType(i_ + 1)
		}
	}
	panic(javalang.NewIllegalArgumentException("No enum constant com.google.javascript.jscomp.ShowByPathWarningsGuard.ShowType." + name))
}

func (this_ ShowByPathWarningsGuard_ShowType) IsJavaNil_() bool { return this_ == 0 }
func (ShowByPathWarningsGuard_ShowType) HasIdentityEquals_() {}
func (ShowByPathWarningsGuard_ShowType) Implements_j2gort_Ext_java_lang_Enum() {}
func (ShowByPathWarningsGuard_ShowType) Implements_j2gort_Ext_java_lang_constant_Constable() {}
func (this_ ShowByPathWarningsGuard_ShowType) GetExt_java_lang_EnumBase() *j2gort.Ext_java_lang_EnumBase { return &j2gort.Ext_java_lang_EnumBase{Core_: this_.GetEnumBase()} }
func (ShowByPathWarningsGuard_ShowType) JavaClassInfo_() (string, string) { return "com.google.javascript.jscomp.ShowByPathWarningsGuard$ShowType", "ShowType" }
func (this_ ShowByPathWarningsGuard_ShowType) GetEnumBase() *j2gort.EnumBase {
	if this_ == 0 {
		return nil
	}
	return &enumBaseTable_ShowByPathWarningsGuard_ShowType[this_-1]
}
func (this_ ShowByPathWarningsGuard_ShowType) Ordinal() int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return int32(this_ - 1)
}
func (this_ ShowByPathWarningsGuard_ShowType) EnumOrdinal_() int32 { return int32(this_) - 1 }
func (this_ ShowByPathWarningsGuard_ShowType) Name() string {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_ShowByPathWarningsGuard_ShowType[this_-1].Name_
}
func (this_ ShowByPathWarningsGuard_ShowType) String() string {
	if this_ == 0 {
		return "null"
	}
	return enumBaseTable_ShowByPathWarningsGuard_ShowType[this_-1].Name_
}
func (this_ ShowByPathWarningsGuard_ShowType) Equals(other_ any) bool {
	o_, ok_ := other_.(ShowByPathWarningsGuard_ShowType)
	return this_ != 0 && ok_ && this_ == o_
}
func (this_ ShowByPathWarningsGuard_ShowType) HashCode() int32 {
	if this_ == 0 {
		return 0
	}
	return enumBaseTable_ShowByPathWarningsGuard_ShowType[this_-1].HashCode()
}
func (this_ ShowByPathWarningsGuard_ShowType) CompareTo(other_ any) int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_ShowByPathWarningsGuard_ShowType[this_-1].CompareTo(other_)
}
func (this_ ShowByPathWarningsGuard_ShowType) WriteGson(out_ *gson.JsonWriter) {
	if this_ == 0 {
		out_.NullValue()
		return
	}
	out_.ValueString(enumBaseTable_ShowByPathWarningsGuard_ShowType[this_-1].Name_)
}
func init() {
	vals_ := []any{ShowByPathWarningsGuard_ShowType(1), ShowByPathWarningsGuard_ShowType(2)}
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.ShowByPathWarningsGuard.ShowType", vals_)
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.ShowByPathWarningsGuard$ShowType", vals_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.ShowByPathWarningsGuard.ShowType", "com.google.javascript.jscomp.ShowByPathWarningsGuard$ShowType", func(v_ any) bool {
		return j2gort.InstanceOf[ShowByPathWarningsGuard_ShowType](v_)
	})
}

