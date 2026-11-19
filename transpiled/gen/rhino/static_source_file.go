/*
 *
 * ***** BEGIN LICENSE BLOCK *****
 * Version: MPL 1.1
 *
 * The contents of this file are subject to the Mozilla Public License Version
 * 1.1 (the "License"); you may not use this file except in compliance with
 * the License. You may obtain a copy of the License at
 * http://www.mozilla.org/MPL/
 *
 * Software distributed under the License is distributed on an "AS IS" basis,
 * WITHOUT WARRANTY OF ANY KIND, either express or implied. See the License
 * for the specific language governing rights and limitations under the
 * License.
 *
 * The Original Code is Rhino code, released
 * May 6, 1999.
 *
 * The Initial Developer of the Original Code is
 * Netscape Communications Corporation.
 * Portions created by the Initial Developer are Copyright (C) 1997-1999
 * the Initial Developer. All Rights Reserved.
 *
 * Contributor(s):
 *   Nick Santos
 *
 * This Go translation is distributed under the MPL 1.1 only, as the
 * original's dual licence allows; the other licence it offers is not
 * used.
 *
 * ***** END LICENSE BLOCK ***** */
// Port of com/google/javascript/rhino/StaticSourceFile.java (generated mechanically by j2go).

package rhino

import (
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	"github.com/paullepoulpe/jsgo/transpiled/shims/gson"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java interface: com.google.javascript.rhino.StaticSourceFile
type StaticSourceFile interface {
	Implements_rhino_StaticSourceFile()
	IsJavaNil_() bool
	GetName() string
	GetKind() StaticSourceFile_SourceKind
	IsStrong() bool
	IsWeak() bool
	IsExtern() bool
	IsNonCode() bool
	IsTypeScriptSource() bool
	IsClosureUnawareCode() bool
	GetLineOffset(int32) int32
	GetLineOfOffset(int32) int32
	GetColumnOfOffset(int32) int32
}
func ToIface_StaticSourceFile[T interface {
	StaticSourceFile
	comparable
}](v T) StaticSourceFile {
	var z T
	if v == z {
		return nil
	}
	return v
}

func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.rhino.StaticSourceFile", "com.google.javascript.rhino.StaticSourceFile", func(v_ any) bool {
		return j2gort.InstanceOf[StaticSourceFile](v_)
	})
}
// Java default interface method: StaticSourceFile.isStrong()
func IfaceDefault_StaticSourceFile_IsStrong(this_ StaticSourceFile) bool {
	_ = this_
	return (this_.GetKind() == StaticGet_StaticSourceFile_SourceKind_STRONG())
}

// Java default interface method: StaticSourceFile.isWeak()
func IfaceDefault_StaticSourceFile_IsWeak(this_ StaticSourceFile) bool {
	_ = this_
	return (this_.GetKind() == StaticGet_StaticSourceFile_SourceKind_WEAK())
}

// Java default interface method: StaticSourceFile.isExtern()
func IfaceDefault_StaticSourceFile_IsExtern(this_ StaticSourceFile) bool {
	_ = this_
	return (this_.GetKind() == StaticGet_StaticSourceFile_SourceKind_EXTERN())
}

// Java default interface method: StaticSourceFile.isNonCode()
func IfaceDefault_StaticSourceFile_IsNonCode(this_ StaticSourceFile) bool {
	_ = this_
	return (this_.GetKind() == StaticGet_StaticSourceFile_SourceKind_NON_CODE())
}

// Java default interface method: StaticSourceFile.isTypeScriptSource()
func IfaceDefault_StaticSourceFile_IsTypeScriptSource(this_ StaticSourceFile) bool {
	_ = this_
	var sourceName string = this_.GetName()
	_ = sourceName
	return (javalang.StringEndsWith(j2gort.CheckNotNullString(sourceName), ".closure.js") || javalang.StringEndsWith(j2gort.CheckNotNullString(sourceName), ".tsx.cl.js"))
}

// Java enum: com.google.javascript.rhino.StaticSourceFile.SourceKind
type StaticSourceFile_SourceKind uint8

const (
	enum_StaticSourceFile_SourceKind_STRONG StaticSourceFile_SourceKind = 1
	enum_StaticSourceFile_SourceKind_WEAK StaticSourceFile_SourceKind = 2
	enum_StaticSourceFile_SourceKind_EXTERN StaticSourceFile_SourceKind = 3
	enum_StaticSourceFile_SourceKind_NON_CODE StaticSourceFile_SourceKind = 4
)

var enumBaseTable_StaticSourceFile_SourceKind = [...]j2gort.EnumBase{
	{Name_: "STRONG", Ordinal_: 0},
	{Name_: "WEAK", Ordinal_: 1},
	{Name_: "EXTERN", Ordinal_: 2},
	{Name_: "NON_CODE", Ordinal_: 3},
}
var enumValues_StaticSourceFile_SourceKind = []StaticSourceFile_SourceKind{1, 2, 3, 4}

func StaticGet_StaticSourceFile_SourceKind_STRONG() StaticSourceFile_SourceKind {
	return enum_StaticSourceFile_SourceKind_STRONG
}
func StaticGet_StaticSourceFile_SourceKind_WEAK() StaticSourceFile_SourceKind {
	return enum_StaticSourceFile_SourceKind_WEAK
}
func StaticGet_StaticSourceFile_SourceKind_EXTERN() StaticSourceFile_SourceKind {
	return enum_StaticSourceFile_SourceKind_EXTERN
}
func StaticGet_StaticSourceFile_SourceKind_NON_CODE() StaticSourceFile_SourceKind {
	return enum_StaticSourceFile_SourceKind_NON_CODE
}

func StaticSourceFile_SourceKind_Values() []StaticSourceFile_SourceKind {
	return j2gort.CloneSlice(enumValues_StaticSourceFile_SourceKind)
}

func StaticSourceFile_SourceKind_ValueOf(name string) StaticSourceFile_SourceKind {
	if name == j2gort.NullString {
		j2gort.ThrowNPE()
	}
	for i_ := range enumBaseTable_StaticSourceFile_SourceKind {
		if enumBaseTable_StaticSourceFile_SourceKind[i_].Name_ == name {
			return StaticSourceFile_SourceKind(i_ + 1)
		}
	}
	panic(javalang.NewIllegalArgumentException("No enum constant com.google.javascript.rhino.StaticSourceFile.SourceKind." + name))
}

func (this_ StaticSourceFile_SourceKind) IsJavaNil_() bool { return this_ == 0 }
func (StaticSourceFile_SourceKind) HasIdentityEquals_() {}
func (StaticSourceFile_SourceKind) Implements_j2gort_Ext_java_lang_Enum() {}
func (StaticSourceFile_SourceKind) Implements_j2gort_Ext_java_lang_constant_Constable() {}
func (this_ StaticSourceFile_SourceKind) GetExt_java_lang_EnumBase() *j2gort.Ext_java_lang_EnumBase { return &j2gort.Ext_java_lang_EnumBase{Core_: this_.GetEnumBase()} }
func (StaticSourceFile_SourceKind) JavaClassInfo_() (string, string) { return "com.google.javascript.rhino.StaticSourceFile$SourceKind", "SourceKind" }
func (this_ StaticSourceFile_SourceKind) GetEnumBase() *j2gort.EnumBase {
	if this_ == 0 {
		return nil
	}
	return &enumBaseTable_StaticSourceFile_SourceKind[this_-1]
}
func (this_ StaticSourceFile_SourceKind) Ordinal() int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return int32(this_ - 1)
}
func (this_ StaticSourceFile_SourceKind) EnumOrdinal_() int32 { return int32(this_) - 1 }
func (this_ StaticSourceFile_SourceKind) Name() string {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_StaticSourceFile_SourceKind[this_-1].Name_
}
func (this_ StaticSourceFile_SourceKind) String() string {
	if this_ == 0 {
		return "null"
	}
	return enumBaseTable_StaticSourceFile_SourceKind[this_-1].Name_
}
func (this_ StaticSourceFile_SourceKind) Equals(other_ any) bool {
	o_, ok_ := other_.(StaticSourceFile_SourceKind)
	return this_ != 0 && ok_ && this_ == o_
}
func (this_ StaticSourceFile_SourceKind) HashCode() int32 {
	if this_ == 0 {
		return 0
	}
	return enumBaseTable_StaticSourceFile_SourceKind[this_-1].HashCode()
}
func (this_ StaticSourceFile_SourceKind) CompareTo(other_ any) int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_StaticSourceFile_SourceKind[this_-1].CompareTo(other_)
}
func (this_ StaticSourceFile_SourceKind) WriteGson(out_ *gson.JsonWriter) {
	if this_ == 0 {
		out_.NullValue()
		return
	}
	out_.ValueString(enumBaseTable_StaticSourceFile_SourceKind[this_-1].Name_)
}
func init() {
	vals_ := []any{StaticSourceFile_SourceKind(1), StaticSourceFile_SourceKind(2), StaticSourceFile_SourceKind(3), StaticSourceFile_SourceKind(4)}
	j2gort.RegisterEnumConstants("com.google.javascript.rhino.StaticSourceFile.SourceKind", vals_)
	j2gort.RegisterEnumConstants("com.google.javascript.rhino.StaticSourceFile$SourceKind", vals_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.rhino.StaticSourceFile.SourceKind", "com.google.javascript.rhino.StaticSourceFile$SourceKind", func(v_ any) bool {
		return j2gort.InstanceOf[StaticSourceFile_SourceKind](v_)
	})
}

