/*
 * Copyright 2009 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/WarningLevel.java (generated mechanically by j2go).

package jscomp

import (
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	"github.com/paullepoulpe/jsgo/transpiled/shims/gson"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java enum: com.google.javascript.jscomp.WarningLevel
type WarningLevel uint8

const (
	enum_WarningLevel_QUIET   WarningLevel = 1
	enum_WarningLevel_DEFAULT WarningLevel = 2
	enum_WarningLevel_VERBOSE WarningLevel = 3
)

var enumBaseTable_WarningLevel = [...]j2gort.EnumBase{
	{Name_: "QUIET", Ordinal_: 0},
	{Name_: "DEFAULT", Ordinal_: 1},
	{Name_: "VERBOSE", Ordinal_: 2},
}
var enumValues_WarningLevel = []WarningLevel{1, 2, 3}

func StaticGet_WarningLevel_QUIET() WarningLevel {
	return enum_WarningLevel_QUIET
}
func StaticGet_WarningLevel_DEFAULT() WarningLevel {
	return enum_WarningLevel_DEFAULT
}
func StaticGet_WarningLevel_VERBOSE() WarningLevel {
	return enum_WarningLevel_VERBOSE
}

func WarningLevel_Values() []WarningLevel {
	return j2gort.CloneSlice(enumValues_WarningLevel)
}

func WarningLevel_ValueOf(name string) WarningLevel {
	if name == j2gort.NullString {
		j2gort.ThrowNPE()
	}
	for i_ := range enumBaseTable_WarningLevel {
		if enumBaseTable_WarningLevel[i_].Name_ == name {
			return WarningLevel(i_ + 1)
		}
	}
	panic(javalang.NewIllegalArgumentException("No enum constant com.google.javascript.jscomp.WarningLevel." + name))
}

// Java method: WarningLevel.setOptionsForWarningLevel(com.google.javascript.jscomp.CompilerOptions options)
func (this_ WarningLevel) SetOptionsForWarningLevel(options *CompilerOptions) {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	_ = options
	switch (this_).Ordinal() {
	case int32(0) /* WarningLevel.QUIET */ :
		WarningLevel_SilenceAllWarnings(options)
	case int32(1) /* WarningLevel.DEFAULT */ :
		WarningLevel_AddDefaultWarnings(options)
	case int32(2) /* WarningLevel.VERBOSE */ :
		WarningLevel_AddVerboseWarnings(options)
	}
}

// Java static method: WarningLevel.silenceAllWarnings(com.google.javascript.jscomp.CompilerOptions options)
func WarningLevel_SilenceAllWarnings(options *CompilerOptions) {
	_ = options
	options.AddWarningsGuard(NewShowByPathWarningsGuard_1_string("the_longest_path_that_cannot_be_expressed_as_a_string"))
	options.SetWarningLevel(StaticGet_DiagnosticGroups_MISSING_PROVIDE(), StaticGet_CheckLevel_OFF())
	options.SetCheckTypes(false)
	options.SetWarningLevel(StaticGet_DiagnosticGroups_CHECK_TYPES(), StaticGet_CheckLevel_OFF())
	options.SetWarningLevel(StaticGet_DiagnosticGroups_CHECK_USELESS_CODE(), StaticGet_CheckLevel_OFF())
	options.SetWarningLevel(StaticGet_DiagnosticGroups_MISSING_RETURN(), StaticGet_CheckLevel_OFF())
	options.SetWarningLevel(StaticGet_DiagnosticGroups_ACCESS_CONTROLS(), StaticGet_CheckLevel_OFF())
	options.SetWarningLevel(StaticGet_DiagnosticGroups_CONST(), StaticGet_CheckLevel_OFF())
	options.SetWarningLevel(StaticGet_DiagnosticGroups_CONSTANT_PROPERTY(), StaticGet_CheckLevel_OFF())
	options.SetCheckSuspiciousCode(false)
	options.SetWarningLevel(StaticGet_DiagnosticGroups_GLOBAL_THIS(), StaticGet_CheckLevel_OFF())
	options.SetWarningLevel(StaticGet_DiagnosticGroups_GLOBAL_THIS(), StaticGet_CheckLevel_OFF())
	options.SetWarningLevel(StaticGet_DiagnosticGroups_ES5_STRICT(), StaticGet_CheckLevel_OFF())
	options.SetWarningLevel(StaticGet_DiagnosticGroups_NON_STANDARD_JSDOC(), StaticGet_CheckLevel_OFF())
}

// Java static method: WarningLevel.addDefaultWarnings(com.google.javascript.jscomp.CompilerOptions options)
func WarningLevel_AddDefaultWarnings(options *CompilerOptions) {
	_ = options
	options.SetCheckSuspiciousCode(true)
	options.SetWarningLevel(StaticGet_DiagnosticGroups_NON_STANDARD_JSDOC(), StaticGet_CheckLevel_OFF())
}

// Java static method: WarningLevel.addVerboseWarnings(com.google.javascript.jscomp.CompilerOptions options)
func WarningLevel_AddVerboseWarnings(options *CompilerOptions) {
	_ = options
	WarningLevel_AddDefaultWarnings(options)
	options.SetCheckSuspiciousCode(true)
	options.SetWarningLevel(StaticGet_DiagnosticGroups_GLOBAL_THIS(), StaticGet_CheckLevel_WARNING())
	options.SetCheckSymbols(true)
	options.SetCheckTypes(true)
	options.SetWarningLevel(StaticGet_DiagnosticGroups_MISSING_PROPERTIES(), StaticGet_CheckLevel_WARNING())
	options.SetWarningLevel(StaticGet_DiagnosticGroups_DEPRECATED(), StaticGet_CheckLevel_WARNING())
	options.SetWarningLevel(StaticGet_DiagnosticGroups_ES5_STRICT(), StaticGet_CheckLevel_WARNING())
	options.SetWarningLevel(StaticGet_DiagnosticGroups_VISIBILITY(), StaticGet_CheckLevel_WARNING())
	options.SetWarningLevel(StaticGet_DiagnosticGroups_CONST(), StaticGet_CheckLevel_WARNING())
	options.SetWarningLevel(StaticGet_DiagnosticGroups_CHECK_REGEXP(), StaticGet_CheckLevel_WARNING())
	options.SetWarningLevel(StaticGet_DiagnosticGroups_STRICT_MODULE_DEP_CHECK(), StaticGet_CheckLevel_WARNING())
	options.SetWarningLevel(StaticGet_DiagnosticGroups_MISSING_RETURN(), StaticGet_CheckLevel_WARNING())
	options.SetWarningLevel(StaticGet_DiagnosticGroups_NON_STANDARD_JSDOC(), StaticGet_CheckLevel_WARNING())
}

func (this_ WarningLevel) IsJavaNil_() bool                              { return this_ == 0 }
func (WarningLevel) HasIdentityEquals_()                                 {}
func (WarningLevel) Implements_j2gort_Ext_java_lang_Enum()               {}
func (WarningLevel) Implements_j2gort_Ext_java_lang_constant_Constable() {}
func (this_ WarningLevel) GetExt_java_lang_EnumBase() *j2gort.Ext_java_lang_EnumBase {
	return &j2gort.Ext_java_lang_EnumBase{Core_: this_.GetEnumBase()}
}
func (WarningLevel) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.WarningLevel", "WarningLevel"
}
func (this_ WarningLevel) GetEnumBase() *j2gort.EnumBase {
	if this_ == 0 {
		return nil
	}
	return &enumBaseTable_WarningLevel[this_-1]
}
func (this_ WarningLevel) Ordinal() int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return int32(this_ - 1)
}
func (this_ WarningLevel) EnumOrdinal_() int32 { return int32(this_) - 1 }
func (this_ WarningLevel) Name() string {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_WarningLevel[this_-1].Name_
}
func (this_ WarningLevel) String() string {
	if this_ == 0 {
		return "null"
	}
	return enumBaseTable_WarningLevel[this_-1].Name_
}
func (this_ WarningLevel) Equals(other_ any) bool {
	o_, ok_ := other_.(WarningLevel)
	return this_ != 0 && ok_ && this_ == o_
}
func (this_ WarningLevel) HashCode() int32 {
	if this_ == 0 {
		return 0
	}
	return enumBaseTable_WarningLevel[this_-1].HashCode()
}
func (this_ WarningLevel) CompareTo(other_ any) int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_WarningLevel[this_-1].CompareTo(other_)
}
func (this_ WarningLevel) WriteGson(out_ *gson.JsonWriter) {
	if this_ == 0 {
		out_.NullValue()
		return
	}
	out_.ValueString(enumBaseTable_WarningLevel[this_-1].Name_)
}
func init() {
	vals_ := []any{WarningLevel(1), WarningLevel(2), WarningLevel(3)}
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.WarningLevel", vals_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.WarningLevel", "com.google.javascript.jscomp.WarningLevel", func(v_ any) bool {
		return j2gort.InstanceOf[WarningLevel](v_)
	})
}
