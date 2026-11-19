/*
 * Copyright 2015 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/TypeMatchingStrategy.java (generated mechanically by j2go).

package jscomp

import (
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
	"github.com/paullepoulpe/jsgo/transpiled/shims/gson"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java enum: com.google.javascript.jscomp.TypeMatchingStrategy
type TypeMatchingStrategy uint8

const (
	enum_TypeMatchingStrategy_LOOSE TypeMatchingStrategy = 1
	enum_TypeMatchingStrategy_STRICT_NULLABILITY TypeMatchingStrategy = 2
	enum_TypeMatchingStrategy_SUBTYPES TypeMatchingStrategy = 3
	enum_TypeMatchingStrategy_EXACT TypeMatchingStrategy = 4
)

var enumBaseTable_TypeMatchingStrategy = [...]j2gort.EnumBase{
	{Name_: "LOOSE", Ordinal_: 0},
	{Name_: "STRICT_NULLABILITY", Ordinal_: 1},
	{Name_: "SUBTYPES", Ordinal_: 2},
	{Name_: "EXACT", Ordinal_: 3},
}
var enumValues_TypeMatchingStrategy = []TypeMatchingStrategy{1, 2, 3, 4}

var enumField_TypeMatchingStrategy_allowSubtypes = [...]bool{
	true,
	true,
	true,
	false,
}
func (this_ TypeMatchingStrategy) FieldGet_allowSubtypes() bool {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumField_TypeMatchingStrategy_allowSubtypes[this_-1]
}

var enumField_TypeMatchingStrategy_ignoreNullability = [...]bool{
	true,
	false,
	false,
	false,
}
func (this_ TypeMatchingStrategy) FieldGet_ignoreNullability() bool {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumField_TypeMatchingStrategy_ignoreNullability[this_-1]
}

var enumField_TypeMatchingStrategy_allowLooseMatches = [...]bool{
	true,
	true,
	false,
	false,
}
func (this_ TypeMatchingStrategy) FieldGet_allowLooseMatches() bool {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumField_TypeMatchingStrategy_allowLooseMatches[this_-1]
}

func StaticGet_TypeMatchingStrategy_LOOSE() TypeMatchingStrategy {
	return enum_TypeMatchingStrategy_LOOSE
}
func StaticGet_TypeMatchingStrategy_STRICT_NULLABILITY() TypeMatchingStrategy {
	return enum_TypeMatchingStrategy_STRICT_NULLABILITY
}
func StaticGet_TypeMatchingStrategy_SUBTYPES() TypeMatchingStrategy {
	return enum_TypeMatchingStrategy_SUBTYPES
}
func StaticGet_TypeMatchingStrategy_EXACT() TypeMatchingStrategy {
	return enum_TypeMatchingStrategy_EXACT
}

func TypeMatchingStrategy_Values() []TypeMatchingStrategy {
	return j2gort.CloneSlice(enumValues_TypeMatchingStrategy)
}

func TypeMatchingStrategy_ValueOf(name string) TypeMatchingStrategy {
	if name == j2gort.NullString {
		j2gort.ThrowNPE()
	}
	for i_ := range enumBaseTable_TypeMatchingStrategy {
		if enumBaseTable_TypeMatchingStrategy[i_].Name_ == name {
			return TypeMatchingStrategy(i_ + 1)
		}
	}
	panic(javalang.NewIllegalArgumentException("No enum constant com.google.javascript.jscomp.TypeMatchingStrategy." + name))
}

// Java method: TypeMatchingStrategy.match(com.google.javascript.rhino.jstype.JSType templateType, com.google.javascript.rhino.jstype.JSType type)
func (this_ TypeMatchingStrategy) Match(templateType pkg_rhino.JSType, type_ pkg_rhino.JSType) TypeMatchingStrategy_MatchResult {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	_ = templateType
	_ = type_
	if (templateType.IsUnknownType()) {
		return StaticGet_TypeMatchingStrategy_MatchResult_MATCH()
	}
	if ((((type_ == nil || type_.IsJavaNil_()) || type_.IsUnknownType()) || type_.IsAllType())) {
		if this_.FieldGet_allowLooseMatches() {
			return StaticGet_TypeMatchingStrategy_MatchResult_LOOSE_MATCH()
		} else {
			return StaticGet_TypeMatchingStrategy_MatchResult_NO_MATCH()
		}
	}
	if (this_.FieldGet_allowSubtypes()) {
		if (this_.FieldGet_ignoreNullability()) {
			type_ = type_.RestrictByNotNullOrUndefined()
		}
		if (type_.IsSubtypeOf_1(templateType)) {
			return StaticGet_TypeMatchingStrategy_MatchResult_MATCH()
		}
	}
	var nullableMismatch bool = (templateType.IsNullable() != type_.IsNullable())
	_ = nullableMismatch
	var voidableMismatch bool = (templateType.IsVoidable() != type_.IsVoidable())
	_ = voidableMismatch
	if ((!(this_.FieldGet_ignoreNullability()) && ((nullableMismatch || voidableMismatch)))) {
		return StaticGet_TypeMatchingStrategy_MatchResult_NO_MATCH()
	}
	if type_.Equals(templateType) {
		return StaticGet_TypeMatchingStrategy_MatchResult_MATCH()
	} else {
		return StaticGet_TypeMatchingStrategy_MatchResult_NO_MATCH()
	}
}

func (this_ TypeMatchingStrategy) IsJavaNil_() bool { return this_ == 0 }
func (TypeMatchingStrategy) HasIdentityEquals_() {}
func (TypeMatchingStrategy) Implements_j2gort_Ext_java_lang_Enum() {}
func (TypeMatchingStrategy) Implements_j2gort_Ext_java_lang_constant_Constable() {}
func (this_ TypeMatchingStrategy) GetExt_java_lang_EnumBase() *j2gort.Ext_java_lang_EnumBase { return &j2gort.Ext_java_lang_EnumBase{Core_: this_.GetEnumBase()} }
func (TypeMatchingStrategy) JavaClassInfo_() (string, string) { return "com.google.javascript.jscomp.TypeMatchingStrategy", "TypeMatchingStrategy" }
func (this_ TypeMatchingStrategy) GetEnumBase() *j2gort.EnumBase {
	if this_ == 0 {
		return nil
	}
	return &enumBaseTable_TypeMatchingStrategy[this_-1]
}
func (this_ TypeMatchingStrategy) Ordinal() int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return int32(this_ - 1)
}
func (this_ TypeMatchingStrategy) EnumOrdinal_() int32 { return int32(this_) - 1 }
func (this_ TypeMatchingStrategy) Name() string {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_TypeMatchingStrategy[this_-1].Name_
}
func (this_ TypeMatchingStrategy) String() string {
	if this_ == 0 {
		return "null"
	}
	return enumBaseTable_TypeMatchingStrategy[this_-1].Name_
}
func (this_ TypeMatchingStrategy) Equals(other_ any) bool {
	o_, ok_ := other_.(TypeMatchingStrategy)
	return this_ != 0 && ok_ && this_ == o_
}
func (this_ TypeMatchingStrategy) HashCode() int32 {
	if this_ == 0 {
		return 0
	}
	return enumBaseTable_TypeMatchingStrategy[this_-1].HashCode()
}
func (this_ TypeMatchingStrategy) CompareTo(other_ any) int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_TypeMatchingStrategy[this_-1].CompareTo(other_)
}
func (this_ TypeMatchingStrategy) WriteGson(out_ *gson.JsonWriter) {
	if this_ == 0 {
		out_.NullValue()
		return
	}
	out_.ValueString(enumBaseTable_TypeMatchingStrategy[this_-1].Name_)
}
func init() {
	vals_ := []any{TypeMatchingStrategy(1), TypeMatchingStrategy(2), TypeMatchingStrategy(3), TypeMatchingStrategy(4)}
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.TypeMatchingStrategy", vals_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.TypeMatchingStrategy", "com.google.javascript.jscomp.TypeMatchingStrategy", func(v_ any) bool {
		return j2gort.InstanceOf[TypeMatchingStrategy](v_)
	})
}

// Java enum: com.google.javascript.jscomp.TypeMatchingStrategy.MatchResult
type TypeMatchingStrategy_MatchResult uint8

const (
	enum_TypeMatchingStrategy_MatchResult_MATCH TypeMatchingStrategy_MatchResult = 1
	enum_TypeMatchingStrategy_MatchResult_NO_MATCH TypeMatchingStrategy_MatchResult = 2
	enum_TypeMatchingStrategy_MatchResult_LOOSE_MATCH TypeMatchingStrategy_MatchResult = 3
)

var enumBaseTable_TypeMatchingStrategy_MatchResult = [...]j2gort.EnumBase{
	{Name_: "MATCH", Ordinal_: 0},
	{Name_: "NO_MATCH", Ordinal_: 1},
	{Name_: "LOOSE_MATCH", Ordinal_: 2},
}
var enumValues_TypeMatchingStrategy_MatchResult = []TypeMatchingStrategy_MatchResult{1, 2, 3}

func StaticGet_TypeMatchingStrategy_MatchResult_MATCH() TypeMatchingStrategy_MatchResult {
	return enum_TypeMatchingStrategy_MatchResult_MATCH
}
func StaticGet_TypeMatchingStrategy_MatchResult_NO_MATCH() TypeMatchingStrategy_MatchResult {
	return enum_TypeMatchingStrategy_MatchResult_NO_MATCH
}
func StaticGet_TypeMatchingStrategy_MatchResult_LOOSE_MATCH() TypeMatchingStrategy_MatchResult {
	return enum_TypeMatchingStrategy_MatchResult_LOOSE_MATCH
}

func TypeMatchingStrategy_MatchResult_Values() []TypeMatchingStrategy_MatchResult {
	return j2gort.CloneSlice(enumValues_TypeMatchingStrategy_MatchResult)
}

func TypeMatchingStrategy_MatchResult_ValueOf(name string) TypeMatchingStrategy_MatchResult {
	if name == j2gort.NullString {
		j2gort.ThrowNPE()
	}
	for i_ := range enumBaseTable_TypeMatchingStrategy_MatchResult {
		if enumBaseTable_TypeMatchingStrategy_MatchResult[i_].Name_ == name {
			return TypeMatchingStrategy_MatchResult(i_ + 1)
		}
	}
	panic(javalang.NewIllegalArgumentException("No enum constant com.google.javascript.jscomp.TypeMatchingStrategy.MatchResult." + name))
}

// Java method: MatchResult.isMatch()
func (this_ TypeMatchingStrategy_MatchResult) IsMatch() bool {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return (this_ != StaticGet_TypeMatchingStrategy_MatchResult_NO_MATCH())
}

// Java method: MatchResult.isLooseMatch()
func (this_ TypeMatchingStrategy_MatchResult) IsLooseMatch() bool {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return (this_ == StaticGet_TypeMatchingStrategy_MatchResult_LOOSE_MATCH())
}

func (this_ TypeMatchingStrategy_MatchResult) IsJavaNil_() bool { return this_ == 0 }
func (TypeMatchingStrategy_MatchResult) HasIdentityEquals_() {}
func (TypeMatchingStrategy_MatchResult) Implements_j2gort_Ext_java_lang_Enum() {}
func (TypeMatchingStrategy_MatchResult) Implements_j2gort_Ext_java_lang_constant_Constable() {}
func (this_ TypeMatchingStrategy_MatchResult) GetExt_java_lang_EnumBase() *j2gort.Ext_java_lang_EnumBase { return &j2gort.Ext_java_lang_EnumBase{Core_: this_.GetEnumBase()} }
func (TypeMatchingStrategy_MatchResult) JavaClassInfo_() (string, string) { return "com.google.javascript.jscomp.TypeMatchingStrategy$MatchResult", "MatchResult" }
func (this_ TypeMatchingStrategy_MatchResult) GetEnumBase() *j2gort.EnumBase {
	if this_ == 0 {
		return nil
	}
	return &enumBaseTable_TypeMatchingStrategy_MatchResult[this_-1]
}
func (this_ TypeMatchingStrategy_MatchResult) Ordinal() int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return int32(this_ - 1)
}
func (this_ TypeMatchingStrategy_MatchResult) EnumOrdinal_() int32 { return int32(this_) - 1 }
func (this_ TypeMatchingStrategy_MatchResult) Name() string {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_TypeMatchingStrategy_MatchResult[this_-1].Name_
}
func (this_ TypeMatchingStrategy_MatchResult) String() string {
	if this_ == 0 {
		return "null"
	}
	return enumBaseTable_TypeMatchingStrategy_MatchResult[this_-1].Name_
}
func (this_ TypeMatchingStrategy_MatchResult) Equals(other_ any) bool {
	o_, ok_ := other_.(TypeMatchingStrategy_MatchResult)
	return this_ != 0 && ok_ && this_ == o_
}
func (this_ TypeMatchingStrategy_MatchResult) HashCode() int32 {
	if this_ == 0 {
		return 0
	}
	return enumBaseTable_TypeMatchingStrategy_MatchResult[this_-1].HashCode()
}
func (this_ TypeMatchingStrategy_MatchResult) CompareTo(other_ any) int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_TypeMatchingStrategy_MatchResult[this_-1].CompareTo(other_)
}
func (this_ TypeMatchingStrategy_MatchResult) WriteGson(out_ *gson.JsonWriter) {
	if this_ == 0 {
		out_.NullValue()
		return
	}
	out_.ValueString(enumBaseTable_TypeMatchingStrategy_MatchResult[this_-1].Name_)
}
func init() {
	vals_ := []any{TypeMatchingStrategy_MatchResult(1), TypeMatchingStrategy_MatchResult(2), TypeMatchingStrategy_MatchResult(3)}
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.TypeMatchingStrategy.MatchResult", vals_)
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.TypeMatchingStrategy$MatchResult", vals_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.TypeMatchingStrategy.MatchResult", "com.google.javascript.jscomp.TypeMatchingStrategy$MatchResult", func(v_ any) bool {
		return j2gort.InstanceOf[TypeMatchingStrategy_MatchResult](v_)
	})
}

