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
// Port of com/google/javascript/jscomp/CssRenamingMap.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	"github.com/paullepoulpe/jsgo/transpiled/shims/gson"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java interface: com.google.javascript.jscomp.CssRenamingMap
type CssRenamingMap interface {
	Implements_jscomp_CssRenamingMap()
	IsJavaNil_() bool
	RenamingMap
	GetStyle() CssRenamingMap_Style
}
func ToIface_CssRenamingMap[T interface {
	CssRenamingMap
	comparable
}](v T) CssRenamingMap {
	var z T
	if v == z {
		return nil
	}
	return v
}

func Upcast_CssRenamingMap[T CssRenamingMap](v T) CssRenamingMap {
	var vi CssRenamingMap = v
	if vi == nil || vi.IsJavaNil_() {
		return nil
	}
	return vi
}

func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.CssRenamingMap", "com.google.javascript.jscomp.CssRenamingMap", func(v_ any) bool {
		return j2gort.InstanceOf[CssRenamingMap](v_)
	})
}
// Java enum: com.google.javascript.jscomp.CssRenamingMap.Style
type CssRenamingMap_Style uint8

const (
	enum_CssRenamingMap_Style_BY_WHOLE CssRenamingMap_Style = 1
	enum_CssRenamingMap_Style_BY_PART CssRenamingMap_Style = 2
)

var enumBaseTable_CssRenamingMap_Style = [...]j2gort.EnumBase{
	{Name_: "BY_WHOLE", Ordinal_: 0},
	{Name_: "BY_PART", Ordinal_: 1},
}
var enumValues_CssRenamingMap_Style = []CssRenamingMap_Style{1, 2}

func StaticGet_CssRenamingMap_Style_BY_WHOLE() CssRenamingMap_Style {
	return enum_CssRenamingMap_Style_BY_WHOLE
}
func StaticGet_CssRenamingMap_Style_BY_PART() CssRenamingMap_Style {
	return enum_CssRenamingMap_Style_BY_PART
}

func CssRenamingMap_Style_Values() []CssRenamingMap_Style {
	return j2gort.CloneSlice(enumValues_CssRenamingMap_Style)
}

func CssRenamingMap_Style_ValueOf(name string) CssRenamingMap_Style {
	if name == j2gort.NullString {
		j2gort.ThrowNPE()
	}
	for i_ := range enumBaseTable_CssRenamingMap_Style {
		if enumBaseTable_CssRenamingMap_Style[i_].Name_ == name {
			return CssRenamingMap_Style(i_ + 1)
		}
	}
	panic(javalang.NewIllegalArgumentException("No enum constant com.google.javascript.jscomp.CssRenamingMap.Style." + name))
}

func (this_ CssRenamingMap_Style) IsJavaNil_() bool { return this_ == 0 }
func (CssRenamingMap_Style) HasIdentityEquals_() {}
func (CssRenamingMap_Style) Implements_j2gort_Ext_java_lang_Enum() {}
func (CssRenamingMap_Style) Implements_j2gort_Ext_java_lang_constant_Constable() {}
func (this_ CssRenamingMap_Style) GetExt_java_lang_EnumBase() *j2gort.Ext_java_lang_EnumBase { return &j2gort.Ext_java_lang_EnumBase{Core_: this_.GetEnumBase()} }
func (CssRenamingMap_Style) JavaClassInfo_() (string, string) { return "com.google.javascript.jscomp.CssRenamingMap$Style", "Style" }
func (this_ CssRenamingMap_Style) GetEnumBase() *j2gort.EnumBase {
	if this_ == 0 {
		return nil
	}
	return &enumBaseTable_CssRenamingMap_Style[this_-1]
}
func (this_ CssRenamingMap_Style) Ordinal() int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return int32(this_ - 1)
}
func (this_ CssRenamingMap_Style) EnumOrdinal_() int32 { return int32(this_) - 1 }
func (this_ CssRenamingMap_Style) Name() string {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_CssRenamingMap_Style[this_-1].Name_
}
func (this_ CssRenamingMap_Style) String() string {
	if this_ == 0 {
		return "null"
	}
	return enumBaseTable_CssRenamingMap_Style[this_-1].Name_
}
func (this_ CssRenamingMap_Style) Equals(other_ any) bool {
	o_, ok_ := other_.(CssRenamingMap_Style)
	return this_ != 0 && ok_ && this_ == o_
}
func (this_ CssRenamingMap_Style) HashCode() int32 {
	if this_ == 0 {
		return 0
	}
	return enumBaseTable_CssRenamingMap_Style[this_-1].HashCode()
}
func (this_ CssRenamingMap_Style) CompareTo(other_ any) int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_CssRenamingMap_Style[this_-1].CompareTo(other_)
}
func (this_ CssRenamingMap_Style) WriteGson(out_ *gson.JsonWriter) {
	if this_ == 0 {
		out_.NullValue()
		return
	}
	out_.ValueString(enumBaseTable_CssRenamingMap_Style[this_-1].Name_)
}
func init() {
	vals_ := []any{CssRenamingMap_Style(1), CssRenamingMap_Style(2)}
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.CssRenamingMap.Style", vals_)
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.CssRenamingMap$Style", vals_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.CssRenamingMap.Style", "com.google.javascript.jscomp.CssRenamingMap$Style", func(v_ any) bool {
		return j2gort.InstanceOf[CssRenamingMap_Style](v_)
	})
}

// Java polymorphic class interface: com.google.javascript.jscomp.CssRenamingMap.ByPart
type CssRenamingMap_ByPart interface {
	Implements_jscomp_CssRenamingMap_ByPart()
	IsJavaNil_() bool
	GetCssRenamingMap_ByPartBase() *CssRenamingMap_ByPartBase
	CssRenamingMap
}
func ToIface_CssRenamingMap_ByPart[T interface {
	CssRenamingMap_ByPart
	comparable
}](v T) CssRenamingMap_ByPart {
	var z T
	if v == z {
		return nil
	}
	return v
}

func Upcast_CssRenamingMap_ByPart[T CssRenamingMap_ByPart](v T) CssRenamingMap_ByPart {
	var vi CssRenamingMap_ByPart = v
	if vi == nil || vi.IsJavaNil_() {
		return nil
	}
	return vi
}

// Java class: com.google.javascript.jscomp.CssRenamingMap.ByPart
type CssRenamingMap_ByPartBase struct {
	idHash_ atomic.Int32
}

func (this_ *CssRenamingMap_ByPartBase) BindSelf_(s any) {
	_ = s
}
func (this_ *CssRenamingMap_ByPartBase) IsJavaNil_() bool { return this_ == nil }
func (*CssRenamingMap_ByPartBase) JavaClassInfo_() (string, string) { return "com.google.javascript.jscomp.CssRenamingMap$ByPart", "ByPart" }
func (this_ *CssRenamingMap_ByPartBase) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.CssRenamingMap.ByPart", "com.google.javascript.jscomp.CssRenamingMap$ByPart", func(v_ any) bool {
		return j2gort.InstanceOf[CssRenamingMap_ByPart](v_)
	})
}
func (b *CssRenamingMap_ByPartBase) GetCssRenamingMap_ByPartBase() *CssRenamingMap_ByPartBase { return b }
func (*CssRenamingMap_ByPartBase) Implements_jscomp_CssRenamingMap_ByPart() {}
func (*CssRenamingMap_ByPartBase) Implements_jscomp_CssRenamingMap() {}
func (*CssRenamingMap_ByPartBase) Implements_jscomp_RenamingMap() {}

func (this_ *CssRenamingMap_ByPartBase) initFields_CssRenamingMap_ByPart() {
	_ = this_
}

func InitCssRenamingMap_ByPart(this_ *CssRenamingMap_ByPartBase) {
	this_.initFields_CssRenamingMap_ByPart()
}

// Java method: ByPart.get(java.lang.String value)
func (this_ *CssRenamingMap_ByPartBase) Get(value string) string {
	_ = value
	panic("abstract or native method invoked")
}

// Java method: ByPart.getStyle()
func (this_ *CssRenamingMap_ByPartBase) GetStyle() CssRenamingMap_Style {
	return StaticGet_CssRenamingMap_Style_BY_PART()
}

// Java polymorphic class interface: com.google.javascript.jscomp.CssRenamingMap.ByWhole
type CssRenamingMap_ByWhole interface {
	Implements_jscomp_CssRenamingMap_ByWhole()
	IsJavaNil_() bool
	GetCssRenamingMap_ByWholeBase() *CssRenamingMap_ByWholeBase
	CssRenamingMap
}
func ToIface_CssRenamingMap_ByWhole[T interface {
	CssRenamingMap_ByWhole
	comparable
}](v T) CssRenamingMap_ByWhole {
	var z T
	if v == z {
		return nil
	}
	return v
}

func Upcast_CssRenamingMap_ByWhole[T CssRenamingMap_ByWhole](v T) CssRenamingMap_ByWhole {
	var vi CssRenamingMap_ByWhole = v
	if vi == nil || vi.IsJavaNil_() {
		return nil
	}
	return vi
}

// Java class: com.google.javascript.jscomp.CssRenamingMap.ByWhole
type CssRenamingMap_ByWholeBase struct {
	idHash_ atomic.Int32
}

func (this_ *CssRenamingMap_ByWholeBase) BindSelf_(s any) {
	_ = s
}
func (this_ *CssRenamingMap_ByWholeBase) IsJavaNil_() bool { return this_ == nil }
func (*CssRenamingMap_ByWholeBase) JavaClassInfo_() (string, string) { return "com.google.javascript.jscomp.CssRenamingMap$ByWhole", "ByWhole" }
func (this_ *CssRenamingMap_ByWholeBase) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.CssRenamingMap.ByWhole", "com.google.javascript.jscomp.CssRenamingMap$ByWhole", func(v_ any) bool {
		return j2gort.InstanceOf[CssRenamingMap_ByWhole](v_)
	})
}
func (b *CssRenamingMap_ByWholeBase) GetCssRenamingMap_ByWholeBase() *CssRenamingMap_ByWholeBase { return b }
func (*CssRenamingMap_ByWholeBase) Implements_jscomp_CssRenamingMap_ByWhole() {}
func (*CssRenamingMap_ByWholeBase) Implements_jscomp_CssRenamingMap() {}
func (*CssRenamingMap_ByWholeBase) Implements_jscomp_RenamingMap() {}

func (this_ *CssRenamingMap_ByWholeBase) initFields_CssRenamingMap_ByWhole() {
	_ = this_
}

func InitCssRenamingMap_ByWhole(this_ *CssRenamingMap_ByWholeBase) {
	this_.initFields_CssRenamingMap_ByWhole()
}

// Java method: ByWhole.get(java.lang.String value)
func (this_ *CssRenamingMap_ByWholeBase) Get(value string) string {
	_ = value
	panic("abstract or native method invoked")
}

// Java method: ByWhole.getStyle()
func (this_ *CssRenamingMap_ByWholeBase) GetStyle() CssRenamingMap_Style {
	return StaticGet_CssRenamingMap_Style_BY_WHOLE()
}

