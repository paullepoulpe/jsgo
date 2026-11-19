/*
 * Copyright 2020 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/base/Tri.java (generated mechanically by j2go).

package base

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	"github.com/paullepoulpe/jsgo/transpiled/shims/gson"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java enum: com.google.javascript.jscomp.base.Tri
type Tri uint8

const (
	enum_Tri_FALSE Tri = 1
	enum_Tri_UNKNOWN Tri = 2
	enum_Tri_TRUE Tri = 3
)

var enumBaseTable_Tri = [...]j2gort.EnumBase{
	{Name_: "FALSE", Ordinal_: 0},
	{Name_: "UNKNOWN", Ordinal_: 1},
	{Name_: "TRUE", Ordinal_: 2},
}
var enumValues_Tri = []Tri{1, 2, 3}

var enumField_Tri_value = [...]int32{
	int32(-1),
	int32(0),
	int32(1),
}
func (this_ Tri) FieldGet_value() int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumField_Tri_value[this_-1]
}

func StaticGet_Tri_FALSE() Tri {
	return enum_Tri_FALSE
}
func StaticGet_Tri_UNKNOWN() Tri {
	return enum_Tri_UNKNOWN
}
func StaticGet_Tri_TRUE() Tri {
	return enum_Tri_TRUE
}

var static_Tri_VALUES []Tri
func StaticGet_Tri_VALUES() []Tri {
	clinit_Tri()
	return static_Tri_VALUES
}
func StaticSet_Tri_VALUES(v []Tri) {
	clinit_Tri()
	static_Tri_VALUES = v
}
var clinitDone_Tri uint32
var clinitState_Tri j2gort.ClinitState
func clinit_Tri() {
	if atomic.LoadUint32(&clinitDone_Tri) != 0 {
		return
	}
	clinitSlow_Tri()
}

func clinitSlow_Tri() {
	if !j2gort.EnterClinit(&clinitState_Tri, &clinitDone_Tri) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_Tri, &clinitDone_Tri)
	static_Tri_VALUES = Tri_Values()
}

func Tri_Values() []Tri {
	clinit_Tri()
	return j2gort.CloneSlice(enumValues_Tri)
}

func Tri_ValueOf(name string) Tri {
	clinit_Tri()
	if name == j2gort.NullString {
		j2gort.ThrowNPE()
	}
	for i_ := range enumBaseTable_Tri {
		if enumBaseTable_Tri[i_].Name_ == name {
			return Tri(i_ + 1)
		}
	}
	panic(javalang.NewIllegalArgumentException("No enum constant com.google.javascript.jscomp.base.Tri." + name))
}

// Java method: Tri.or(com.google.javascript.jscomp.base.Tri x)
func (this_ Tri) Or(x Tri) Tri {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	_ = x
	if ((this_.FieldGet_value() > x.FieldGet_value())) {
		return this_
	} else {
		return x
	}
}

// Java method: Tri.and(com.google.javascript.jscomp.base.Tri x)
func (this_ Tri) And(x Tri) Tri {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	_ = x
	if ((this_.FieldGet_value() < x.FieldGet_value())) {
		return this_
	} else {
		return x
	}
}

// Java method: Tri.xor(com.google.javascript.jscomp.base.Tri x)
func (this_ Tri) Xor(x Tri) Tri {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	_ = x
	return Tri_ForInt((-(this_.FieldGet_value()) * x.FieldGet_value()))
}

// Java method: Tri.not()
func (this_ Tri) Not() Tri {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return Tri_ForInt(-(this_.FieldGet_value()))
}

// Java method: Tri.toBoolean(boolean x)
func (this_ Tri) ToBoolean(x bool) bool {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	_ = x
	switch (this_).Ordinal() {
	case int32(0) /* Tri.FALSE */:
		return false
	case int32(1) /* Tri.UNKNOWN */:
		return x
	case int32(2) /* Tri.TRUE */:
		return true
	}
	panic("unreachable switch expression")
}

// Java static method: Tri.forBoolean(boolean x)
func Tri_ForBoolean(x bool) Tri {
	clinit_Tri()
	_ = x
	if x {
		return StaticGet_Tri_TRUE()
	} else {
		return StaticGet_Tri_FALSE()
	}
}

// Java method: Tri.toString()
func (this_ Tri) String() string {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return javalang.StringToLowerCaseASCII((this_).Name())
}

// Java static method: Tri.forInt(int x)
func Tri_ForInt(x int32) Tri {
	clinit_Tri()
	_ = x
	return StaticGet_Tri_VALUES()[(x + int32(1))]
}

func (this_ Tri) IsJavaNil_() bool { return this_ == 0 }
func (Tri) HasIdentityEquals_() {}
func (Tri) Implements_j2gort_Ext_java_lang_Enum() {}
func (Tri) Implements_j2gort_Ext_java_lang_constant_Constable() {}
func (this_ Tri) GetExt_java_lang_EnumBase() *j2gort.Ext_java_lang_EnumBase { return &j2gort.Ext_java_lang_EnumBase{Core_: this_.GetEnumBase()} }
func (Tri) JavaClassInfo_() (string, string) { return "com.google.javascript.jscomp.base.Tri", "Tri" }
func (this_ Tri) GetEnumBase() *j2gort.EnumBase {
	if this_ == 0 {
		return nil
	}
	return &enumBaseTable_Tri[this_-1]
}
func (this_ Tri) Ordinal() int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return int32(this_ - 1)
}
func (this_ Tri) EnumOrdinal_() int32 { return int32(this_) - 1 }
func (this_ Tri) Name() string {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_Tri[this_-1].Name_
}
func (this_ Tri) Equals(other_ any) bool {
	o_, ok_ := other_.(Tri)
	return this_ != 0 && ok_ && this_ == o_
}
func (this_ Tri) HashCode() int32 {
	if this_ == 0 {
		return 0
	}
	return enumBaseTable_Tri[this_-1].HashCode()
}
func (this_ Tri) CompareTo(other_ any) int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_Tri[this_-1].CompareTo(other_)
}
func (this_ Tri) WriteGson(out_ *gson.JsonWriter) {
	if this_ == 0 {
		out_.NullValue()
		return
	}
	out_.ValueString(enumBaseTable_Tri[this_-1].Name_)
}
func init() {
	vals_ := []any{Tri(1), Tri(2), Tri(3)}
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.base.Tri", vals_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.base.Tri", "com.google.javascript.jscomp.base.Tri", func(v_ any) bool {
		return j2gort.InstanceOf[Tri](v_)
	})
}

