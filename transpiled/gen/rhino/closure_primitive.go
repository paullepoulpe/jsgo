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
 *   Bob Jervis
 *   Google Inc.
 *
 * This Go translation is distributed under the MPL 1.1 only, as the
 * original's dual licence allows; the other licence it offers is not
 * used.
 *
 * ***** END LICENSE BLOCK ***** */
// Port of com/google/javascript/rhino/ClosurePrimitive.java (generated mechanically by j2go).

package rhino

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	"github.com/paullepoulpe/jsgo/transpiled/shims/gson"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java enum: com.google.javascript.rhino.ClosurePrimitive
type ClosurePrimitive uint8

const (
	enum_ClosurePrimitive_ASSERTS_FAIL           ClosurePrimitive = 1
	enum_ClosurePrimitive_ASSERTS_MATCHES_RETURN ClosurePrimitive = 2
	enum_ClosurePrimitive_ASSERTS_TRUTHY         ClosurePrimitive = 3
)

var enumBaseTable_ClosurePrimitive = [...]j2gort.EnumBase{
	{Name_: "ASSERTS_FAIL", Ordinal_: 0},
	{Name_: "ASSERTS_MATCHES_RETURN", Ordinal_: 1},
	{Name_: "ASSERTS_TRUTHY", Ordinal_: 2},
}
var enumValues_ClosurePrimitive = []ClosurePrimitive{1, 2, 3}

func StaticGet_ClosurePrimitive_ASSERTS_FAIL() ClosurePrimitive {
	return enum_ClosurePrimitive_ASSERTS_FAIL
}
func StaticGet_ClosurePrimitive_ASSERTS_MATCHES_RETURN() ClosurePrimitive {
	return enum_ClosurePrimitive_ASSERTS_MATCHES_RETURN
}
func StaticGet_ClosurePrimitive_ASSERTS_TRUTHY() ClosurePrimitive {
	return enum_ClosurePrimitive_ASSERTS_TRUTHY
}

var static_ClosurePrimitive_ID_TO_ENUM j2gort.Ext_com_google_common_collect_ImmutableMap

func StaticGet_ClosurePrimitive_ID_TO_ENUM() j2gort.Ext_com_google_common_collect_ImmutableMap {
	clinit_ClosurePrimitive()
	return static_ClosurePrimitive_ID_TO_ENUM
}
func StaticSet_ClosurePrimitive_ID_TO_ENUM(v j2gort.Ext_com_google_common_collect_ImmutableMap) {
	clinit_ClosurePrimitive()
	static_ClosurePrimitive_ID_TO_ENUM = v
}

var clinitDone_ClosurePrimitive uint32
var clinitState_ClosurePrimitive j2gort.ClinitState

func clinit_ClosurePrimitive() {
	if atomic.LoadUint32(&clinitDone_ClosurePrimitive) != 0 {
		return
	}
	clinitSlow_ClosurePrimitive()
}

func clinitSlow_ClosurePrimitive() {
	if !j2gort.EnterClinit(&clinitState_ClosurePrimitive, &clinitDone_ClosurePrimitive) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_ClosurePrimitive, &clinitDone_ClosurePrimitive)
	static_ClosurePrimitive_ID_TO_ENUM = j2gort.Ext_com_google_common_collect_ImmutableMap_Of_6("asserts.fail", j2gort.BoxNullableVal(StaticGet_ClosurePrimitive_ASSERTS_FAIL()), "asserts.truthy", j2gort.BoxNullableVal(StaticGet_ClosurePrimitive_ASSERTS_TRUTHY()), "asserts.matchesReturn", j2gort.BoxNullableVal(StaticGet_ClosurePrimitive_ASSERTS_MATCHES_RETURN()))
}

func ClosurePrimitive_Values() []ClosurePrimitive {
	clinit_ClosurePrimitive()
	return j2gort.CloneSlice(enumValues_ClosurePrimitive)
}

func ClosurePrimitive_ValueOf(name string) ClosurePrimitive {
	clinit_ClosurePrimitive()
	if name == j2gort.NullString {
		j2gort.ThrowNPE()
	}
	for i_ := range enumBaseTable_ClosurePrimitive {
		if enumBaseTable_ClosurePrimitive[i_].Name_ == name {
			return ClosurePrimitive(i_ + 1)
		}
	}
	panic(javalang.NewIllegalArgumentException("No enum constant com.google.javascript.rhino.ClosurePrimitive." + name))
}

// Java static method: ClosurePrimitive.fromStringId(java.lang.@org.jspecify.annotations.Nullable String id)
func ClosurePrimitive_FromStringId(id string) ClosurePrimitive {
	clinit_ClosurePrimitive()
	_ = id
	if id == j2gort.NullString {
		return ClosurePrimitive(0)
	}
	var primitive ClosurePrimitive = j2gort.Cast[ClosurePrimitive](j2gort.MapGetStr(StaticGet_ClosurePrimitive_ID_TO_ENUM(), id))
	_ = primitive
	if primitive == 0 {
		panic(j2gort.AsIface[javalang.IIllegalArgumentException](j2gort.ConstructShimThrowable("java.lang.IllegalArgumentException", "String id "+j2gort.StrValString(id)+" does not match any ClosurePrimitive")))
	}
	return primitive
}

func (this_ ClosurePrimitive) IsJavaNil_() bool                              { return this_ == 0 }
func (ClosurePrimitive) HasIdentityEquals_()                                 {}
func (ClosurePrimitive) Implements_j2gort_Ext_java_lang_Enum()               {}
func (ClosurePrimitive) Implements_j2gort_Ext_java_lang_constant_Constable() {}
func (this_ ClosurePrimitive) GetExt_java_lang_EnumBase() *j2gort.Ext_java_lang_EnumBase {
	return &j2gort.Ext_java_lang_EnumBase{Core_: this_.GetEnumBase()}
}
func (ClosurePrimitive) JavaClassInfo_() (string, string) {
	return "com.google.javascript.rhino.ClosurePrimitive", "ClosurePrimitive"
}
func (this_ ClosurePrimitive) GetEnumBase() *j2gort.EnumBase {
	if this_ == 0 {
		return nil
	}
	return &enumBaseTable_ClosurePrimitive[this_-1]
}
func (this_ ClosurePrimitive) Ordinal() int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return int32(this_ - 1)
}
func (this_ ClosurePrimitive) EnumOrdinal_() int32 { return int32(this_) - 1 }
func (this_ ClosurePrimitive) Name() string {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_ClosurePrimitive[this_-1].Name_
}
func (this_ ClosurePrimitive) String() string {
	if this_ == 0 {
		return "null"
	}
	return enumBaseTable_ClosurePrimitive[this_-1].Name_
}
func (this_ ClosurePrimitive) Equals(other_ any) bool {
	o_, ok_ := other_.(ClosurePrimitive)
	return this_ != 0 && ok_ && this_ == o_
}
func (this_ ClosurePrimitive) HashCode() int32 {
	if this_ == 0 {
		return 0
	}
	return enumBaseTable_ClosurePrimitive[this_-1].HashCode()
}
func (this_ ClosurePrimitive) CompareTo(other_ any) int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_ClosurePrimitive[this_-1].CompareTo(other_)
}
func (this_ ClosurePrimitive) WriteGson(out_ *gson.JsonWriter) {
	if this_ == 0 {
		out_.NullValue()
		return
	}
	out_.ValueString(enumBaseTable_ClosurePrimitive[this_-1].Name_)
}
func init() {
	vals_ := []any{ClosurePrimitive(1), ClosurePrimitive(2), ClosurePrimitive(3)}
	j2gort.RegisterEnumConstants("com.google.javascript.rhino.ClosurePrimitive", vals_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.rhino.ClosurePrimitive", "com.google.javascript.rhino.ClosurePrimitive", func(v_ any) bool {
		return j2gort.InstanceOf[ClosurePrimitive](v_)
	})
}
