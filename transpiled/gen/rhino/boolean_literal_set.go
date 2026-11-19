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
// Port of com/google/javascript/rhino/jstype/BooleanLiteralSet.java (generated mechanically by j2go).

package rhino

import (
	pkg_base "github.com/paullepoulpe/jsgo/transpiled/gen/base"
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	"github.com/paullepoulpe/jsgo/transpiled/shims/gson"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java enum: com.google.javascript.rhino.jstype.BooleanLiteralSet
type BooleanLiteralSet uint8

const (
	enum_BooleanLiteralSet_EMPTY BooleanLiteralSet = 1
	enum_BooleanLiteralSet_TRUE  BooleanLiteralSet = 2
	enum_BooleanLiteralSet_FALSE BooleanLiteralSet = 3
	enum_BooleanLiteralSet_BOTH  BooleanLiteralSet = 4
)

var enumBaseTable_BooleanLiteralSet = [...]j2gort.EnumBase{
	{Name_: "EMPTY", Ordinal_: 0},
	{Name_: "TRUE", Ordinal_: 1},
	{Name_: "FALSE", Ordinal_: 2},
	{Name_: "BOTH", Ordinal_: 3},
}
var enumValues_BooleanLiteralSet = []BooleanLiteralSet{1, 2, 3, 4}

func StaticGet_BooleanLiteralSet_EMPTY() BooleanLiteralSet {
	return enum_BooleanLiteralSet_EMPTY
}
func StaticGet_BooleanLiteralSet_TRUE() BooleanLiteralSet {
	return enum_BooleanLiteralSet_TRUE
}
func StaticGet_BooleanLiteralSet_FALSE() BooleanLiteralSet {
	return enum_BooleanLiteralSet_FALSE
}
func StaticGet_BooleanLiteralSet_BOTH() BooleanLiteralSet {
	return enum_BooleanLiteralSet_BOTH
}

func BooleanLiteralSet_Values() []BooleanLiteralSet {
	return j2gort.CloneSlice(enumValues_BooleanLiteralSet)
}

func BooleanLiteralSet_ValueOf(name string) BooleanLiteralSet {
	if name == j2gort.NullString {
		j2gort.ThrowNPE()
	}
	for i_ := range enumBaseTable_BooleanLiteralSet {
		if enumBaseTable_BooleanLiteralSet[i_].Name_ == name {
			return BooleanLiteralSet(i_ + 1)
		}
	}
	panic(javalang.NewIllegalArgumentException("No enum constant com.google.javascript.rhino.jstype.BooleanLiteralSet." + name))
}

// Java method: BooleanLiteralSet.fromOrdinal(int ordinal)
func (this_ BooleanLiteralSet) priv_BooleanLiteralSet_fromOrdinal_8824(ordinal int32) BooleanLiteralSet {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	_ = ordinal
	switch ordinal {
	case int32(0):
		return StaticGet_BooleanLiteralSet_EMPTY()
	case int32(1):
		return StaticGet_BooleanLiteralSet_TRUE()
	case int32(2):
		return StaticGet_BooleanLiteralSet_FALSE()
	case int32(3):
		return StaticGet_BooleanLiteralSet_BOTH()
	default:
		panic(j2gort.AsIface[javalang.IIllegalArgumentException](j2gort.ConstructShimThrowable("java.lang.IllegalArgumentException", "Ordinal: "+j2gort.StrValInt32(ordinal))))
	}
	panic("unreachable switch expression")
}

// Java method: BooleanLiteralSet.intersection(com.google.javascript.rhino.jstype.BooleanLiteralSet that)
func (this_ BooleanLiteralSet) Intersection(that BooleanLiteralSet) BooleanLiteralSet {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	_ = that
	return this_.priv_BooleanLiteralSet_fromOrdinal_8824(((this_).Ordinal() & (that).Ordinal()))
}

// Java method: BooleanLiteralSet.union(com.google.javascript.rhino.jstype.BooleanLiteralSet that)
func (this_ BooleanLiteralSet) Union(that BooleanLiteralSet) BooleanLiteralSet {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	_ = that
	return this_.priv_BooleanLiteralSet_fromOrdinal_8824(((this_).Ordinal() | (that).Ordinal()))
}

// Java method: BooleanLiteralSet.contains(boolean literalValue)
func (this_ BooleanLiteralSet) Contains(literalValue bool) bool {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	_ = literalValue
	switch (this_).Ordinal() {
	case int32(0):
		return false
	case int32(1):
		return literalValue
	case int32(2):
		return !(literalValue)
	case int32(3):
		return true
	default:
		panic(j2gort.AsIface[javalang.IIndexOutOfBoundsException](j2gort.ConstructShimThrowable("java.lang.IndexOutOfBoundsException", "Ordinal: "+j2gort.StrValInt32((this_).Ordinal()))))
	}
	panic("unreachable switch expression")
}

// Java static method: BooleanLiteralSet.get(boolean literalValue)
func BooleanLiteralSet_Get(literalValue bool) BooleanLiteralSet {
	_ = literalValue
	if literalValue {
		return StaticGet_BooleanLiteralSet_TRUE()
	} else {
		return StaticGet_BooleanLiteralSet_FALSE()
	}
}

// Java method: BooleanLiteralSet.toTri()
func (this_ BooleanLiteralSet) ToTri() pkg_base.Tri {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	if this_ == StaticGet_BooleanLiteralSet_TRUE() {
		return pkg_base.StaticGet_Tri_TRUE()
	} else if this_ == StaticGet_BooleanLiteralSet_FALSE() {
		return pkg_base.StaticGet_Tri_FALSE()
	}
	return pkg_base.StaticGet_Tri_UNKNOWN()
}

func (this_ BooleanLiteralSet) IsJavaNil_() bool                              { return this_ == 0 }
func (BooleanLiteralSet) HasIdentityEquals_()                                 {}
func (BooleanLiteralSet) Implements_j2gort_Ext_java_lang_Enum()               {}
func (BooleanLiteralSet) Implements_j2gort_Ext_java_lang_constant_Constable() {}
func (this_ BooleanLiteralSet) GetExt_java_lang_EnumBase() *j2gort.Ext_java_lang_EnumBase {
	return &j2gort.Ext_java_lang_EnumBase{Core_: this_.GetEnumBase()}
}
func (BooleanLiteralSet) JavaClassInfo_() (string, string) {
	return "com.google.javascript.rhino.jstype.BooleanLiteralSet", "BooleanLiteralSet"
}
func (this_ BooleanLiteralSet) GetEnumBase() *j2gort.EnumBase {
	if this_ == 0 {
		return nil
	}
	return &enumBaseTable_BooleanLiteralSet[this_-1]
}
func (this_ BooleanLiteralSet) Ordinal() int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return int32(this_ - 1)
}
func (this_ BooleanLiteralSet) EnumOrdinal_() int32 { return int32(this_) - 1 }
func (this_ BooleanLiteralSet) Name() string {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_BooleanLiteralSet[this_-1].Name_
}
func (this_ BooleanLiteralSet) String() string {
	if this_ == 0 {
		return "null"
	}
	return enumBaseTable_BooleanLiteralSet[this_-1].Name_
}
func (this_ BooleanLiteralSet) Equals(other_ any) bool {
	o_, ok_ := other_.(BooleanLiteralSet)
	return this_ != 0 && ok_ && this_ == o_
}
func (this_ BooleanLiteralSet) HashCode() int32 {
	if this_ == 0 {
		return 0
	}
	return enumBaseTable_BooleanLiteralSet[this_-1].HashCode()
}
func (this_ BooleanLiteralSet) CompareTo(other_ any) int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_BooleanLiteralSet[this_-1].CompareTo(other_)
}
func (this_ BooleanLiteralSet) WriteGson(out_ *gson.JsonWriter) {
	if this_ == 0 {
		out_.NullValue()
		return
	}
	out_.ValueString(enumBaseTable_BooleanLiteralSet[this_-1].Name_)
}
func init() {
	vals_ := []any{BooleanLiteralSet(1), BooleanLiteralSet(2), BooleanLiteralSet(3), BooleanLiteralSet(4)}
	j2gort.RegisterEnumConstants("com.google.javascript.rhino.jstype.BooleanLiteralSet", vals_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.rhino.jstype.BooleanLiteralSet", "com.google.javascript.rhino.jstype.BooleanLiteralSet", func(v_ any) bool {
		return j2gort.InstanceOf[BooleanLiteralSet](v_)
	})
}
