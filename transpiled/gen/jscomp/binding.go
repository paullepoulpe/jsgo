/*
 * Copyright 2018 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/modules/Binding.java (generated mechanically by j2go).

package jscomp

import (
	"fmt"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
	"github.com/paullepoulpe/jsgo/transpiled/shims/gson"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java record: com.google.javascript.jscomp.modules.Binding
type Binding struct {
	metadata          ModuleMetadataMap_ModuleMetadata `gson:"metadata"`
	originatingExport *Export                          `gson:"originatingExport"`
	closureNamespace  string                           `gson:"closureNamespace"`
	sourceNode        pkg_rhino.Node                   `gson:"sourceNode"`
	isModuleNamespace bool                             `gson:"isModuleNamespace"`
	createdBy         Binding_CreatedBy                `gson:"createdBy"`
}

func (this_ *Binding) BindSelf_(s any) {
	_ = s
}
func (this_ *Binding) IsJavaNil_() bool { return this_ == nil }
func (*Binding) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.modules.Binding", "Binding"
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.modules.Binding", "com.google.javascript.jscomp.modules.Binding", func(v_ any) bool {
		return j2gort.InstanceOf[*Binding](v_)
	})
}
func (*Binding) Implements_j2gort_Ext_java_lang_Record() {}

func (this_ *Binding) initFields_Binding() {
	_ = this_
}

// Java constructor: Binding(com.google.javascript.jscomp.modules.ModuleMetadataMap.ModuleMetadata metadata, com.google.javascript.rhino.@org.jspecify.annotations.Nullable Node sourceNode, com.google.javascript.jscomp.modules.@org.jspecify.annotations.Nullable Export originatingExport, boolean isModuleNamespace, java.lang.@org.jspecify.annotations.Nullable String closureNamespace, com.google.javascript.jscomp.modules.Binding.CreatedBy createdBy)
func NewBinding(metadata ModuleMetadataMap_ModuleMetadata, sourceNode pkg_rhino.Node, originatingExport *Export, isModuleNamespace bool, closureNamespace string, createdBy Binding_CreatedBy) *Binding {
	this_ := &Binding{}
	InitBinding(this_, metadata, sourceNode, originatingExport, isModuleNamespace, closureNamespace, createdBy)
	return this_
}

func InitBinding(this_ *Binding, metadata ModuleMetadataMap_ModuleMetadata, sourceNode pkg_rhino.Node, originatingExport *Export, isModuleNamespace bool, closureNamespace string, createdBy Binding_CreatedBy) {
	_ = metadata
	_ = sourceNode
	_ = originatingExport
	_ = isModuleNamespace
	_ = closureNamespace
	_ = createdBy
	this_.initFields_Binding()
	j2gort.CheckNotNullIface(metadata)
	j2gort.CheckNotNullEnum(createdBy)
	this_.metadata = metadata
	this_.sourceNode = sourceNode
	this_.originatingExport = originatingExport
	this_.isModuleNamespace = isModuleNamespace
	this_.closureNamespace = closureNamespace
	this_.createdBy = createdBy
}

// Java static method: Binding.from(com.google.javascript.jscomp.modules.Export boundExport, com.google.javascript.rhino.Node sourceNode)
func Binding_From_2(boundExport *Export, sourceNode pkg_rhino.Node) *Binding {
	_ = boundExport
	_ = sourceNode
	return NewBinding(boundExport.ModuleMetadata(), sourceNode, boundExport, false, boundExport.ClosureNamespace(), StaticGet_Binding_CreatedBy_EXPORT())
}

// Java static method: Binding.from(com.google.javascript.jscomp.modules.ModuleMetadataMap.ModuleMetadata metadata, com.google.javascript.rhino.Node sourceNode, java.lang.String closureNamespace, com.google.javascript.jscomp.modules.Binding.CreatedBy createdBy)
func Binding_From_4(metadata ModuleMetadataMap_ModuleMetadata, sourceNode pkg_rhino.Node, closureNamespace string, createdBy Binding_CreatedBy) *Binding {
	_ = metadata
	_ = sourceNode
	_ = closureNamespace
	_ = createdBy
	if !(createdBy.IsClosureImport()) {
		j2gort.CheckArgumentFmt(false, "Expected goog.require(Type) or goog.forwardDeclare, got %s", j2gort.BoxNullableVal(createdBy))
	}
	return NewBinding(metadata, sourceNode, nil, true, closureNamespace, createdBy)
}

// Java static method: Binding.from(com.google.javascript.jscomp.modules.ModuleMetadataMap.ModuleMetadata metadataOfBoundModule, java.lang.@org.jspecify.annotations.Nullable String closureNamespace, com.google.javascript.rhino.Node sourceNode)
func Binding_From_3(metadataOfBoundModule ModuleMetadataMap_ModuleMetadata, closureNamespace string, sourceNode pkg_rhino.Node) *Binding {
	_ = metadataOfBoundModule
	_ = closureNamespace
	_ = sourceNode
	return NewBinding(metadataOfBoundModule, sourceNode, nil, true, closureNamespace, StaticGet_Binding_CreatedBy_IMPORT())
}

// Java method: Binding.copy(com.google.javascript.rhino.Node sourceNode, com.google.javascript.jscomp.modules.Binding.CreatedBy createdBy)
func (this_ *Binding) Copy(sourceNode pkg_rhino.Node, createdBy Binding_CreatedBy) *Binding {
	_ = sourceNode
	_ = createdBy
	j2gort.CheckNotNullNode(sourceNode)
	return NewBinding(this_.Metadata(), sourceNode, this_.OriginatingExport(), this_.IsModuleNamespace(), this_.ClosureNamespace(), createdBy)
}

// Java method: Binding.boundName()
func (this_ *Binding) BoundName() string {
	if !(!(this_.IsModuleNamespace())) {
		j2gort.CheckState(false)
	}
	return this_.OriginatingExport().LocalName()
}

// Java method: Binding.isMutated()
func (this_ *Binding) IsMutated() bool {
	return (!(this_.IsModuleNamespace()) && this_.OriginatingExport().Mutated())
}

// Java method: Binding.isCreatedByEsImport()
func (this_ *Binding) IsCreatedByEsImport() bool {
	return func() bool {
		re_1 := this_.CreatedBy()
		ae_2 := StaticGet_Binding_CreatedBy_IMPORT()
		if re_1 == 0 {
			j2gort.ThrowNPE()
		}
		return re_1 == ae_2
	}()
}

// Java method: Binding.isCreatedByEsExport()
func (this_ *Binding) IsCreatedByEsExport() bool {
	return func() bool {
		re_3 := this_.CreatedBy()
		ae_4 := StaticGet_Binding_CreatedBy_EXPORT()
		if re_3 == 0 {
			j2gort.ThrowNPE()
		}
		return re_3 == ae_4
	}()
}

// Java method: Binding.isSomeImport()
func (this_ *Binding) IsSomeImport() bool {
	return !(func() bool {
		re_5 := this_.CreatedBy()
		ae_6 := StaticGet_Binding_CreatedBy_EXPORT()
		if re_5 == 0 {
			j2gort.ThrowNPE()
		}
		return re_5 == ae_6
	}())
}

func (this_ *Binding) Metadata() ModuleMetadataMap_ModuleMetadata { return this_.metadata }
func (this_ *Binding) SourceNode() pkg_rhino.Node                 { return this_.sourceNode }
func (this_ *Binding) OriginatingExport() *Export                 { return this_.originatingExport }
func (this_ *Binding) IsModuleNamespace() bool                    { return this_.isModuleNamespace }
func (this_ *Binding) ClosureNamespace() string                   { return this_.closureNamespace }
func (this_ *Binding) CreatedBy() Binding_CreatedBy               { return this_.createdBy }
func (this_ *Binding) Equals(other_ any) bool {
	o_, ok_ := other_.(*Binding)
	if !ok_ || o_ == nil {
		return false
	}
	if this_ == o_ {
		return true
	}
	if !javalang.ObjectsEquals(this_.metadata, o_.metadata) {
		return false
	}
	if this_.sourceNode != o_.sourceNode {
		return false
	}
	if !javalang.ObjectsEquals(j2gort.PtrToAny(this_.originatingExport), j2gort.PtrToAny(o_.originatingExport)) {
		return false
	}
	if this_.isModuleNamespace != o_.isModuleNamespace {
		return false
	}
	if this_.closureNamespace != o_.closureNamespace {
		return false
	}
	if this_.createdBy != o_.createdBy {
		return false
	}
	return true
}
func (this_ *Binding) HashCode() int32 {
	var h_ int32 = 0
	{
		hi_7 := this_.metadata
		h_ = 31 * h_
		if hi_7 != nil && !hi_7.IsJavaNil_() {
			h_ += hi_7.HashCode()
		}
	}
	h_ = 31*h_ + int32(this_.sourceNode)
	{
		hp_8 := this_.originatingExport
		h_ = 31 * h_
		if hp_8 != nil {
			h_ += hp_8.HashCode()
		}
	}
	if this_.isModuleNamespace {
		h_ = 31*h_ + 1231
	} else {
		h_ = 31*h_ + 1237
	}
	{
		hs_9 := this_.closureNamespace
		h_ = 31 * h_
		if hs_9 != j2gort.NullString {
			h_ += javalang.StringHashCode(hs_9)
		}
	}
	h_ = 31*h_ + (this_.createdBy).HashCode()
	return h_
}
func (this_ *Binding) String() string {
	_ = fmt.Sprint
	return "Binding[" + "metadata=" + j2gort.StrVal(this_.metadata) + ", " + "sourceNode=" + j2gort.StrVal(j2gort.BoxNode(this_.sourceNode)) + ", " + "originatingExport=" + j2gort.StrVal(j2gort.PtrToAny(this_.originatingExport)) + ", " + "isModuleNamespace=" + j2gort.StrVal(j2gort.ValueOfBoolean(this_.isModuleNamespace)) + ", " + "closureNamespace=" + j2gort.StrVal(j2gort.BoxString(this_.closureNamespace)) + ", " + "createdBy=" + j2gort.StrVal(j2gort.BoxNullableVal(this_.createdBy)) + "]"
}

// Java enum: com.google.javascript.jscomp.modules.Binding.CreatedBy
type Binding_CreatedBy uint8

const (
	enum_Binding_CreatedBy_EXPORT               Binding_CreatedBy = 1
	enum_Binding_CreatedBy_IMPORT               Binding_CreatedBy = 2
	enum_Binding_CreatedBy_GOOG_REQUIRE         Binding_CreatedBy = 3
	enum_Binding_CreatedBy_GOOG_REQUIRE_TYPE    Binding_CreatedBy = 4
	enum_Binding_CreatedBy_GOOG_FORWARD_DECLARE Binding_CreatedBy = 5
)

var enumBaseTable_Binding_CreatedBy = [...]j2gort.EnumBase{
	{Name_: "EXPORT", Ordinal_: 0},
	{Name_: "IMPORT", Ordinal_: 1},
	{Name_: "GOOG_REQUIRE", Ordinal_: 2},
	{Name_: "GOOG_REQUIRE_TYPE", Ordinal_: 3},
	{Name_: "GOOG_FORWARD_DECLARE", Ordinal_: 4},
}
var enumValues_Binding_CreatedBy = []Binding_CreatedBy{1, 2, 3, 4, 5}

func StaticGet_Binding_CreatedBy_EXPORT() Binding_CreatedBy {
	return enum_Binding_CreatedBy_EXPORT
}
func StaticGet_Binding_CreatedBy_IMPORT() Binding_CreatedBy {
	return enum_Binding_CreatedBy_IMPORT
}
func StaticGet_Binding_CreatedBy_GOOG_REQUIRE() Binding_CreatedBy {
	return enum_Binding_CreatedBy_GOOG_REQUIRE
}
func StaticGet_Binding_CreatedBy_GOOG_REQUIRE_TYPE() Binding_CreatedBy {
	return enum_Binding_CreatedBy_GOOG_REQUIRE_TYPE
}
func StaticGet_Binding_CreatedBy_GOOG_FORWARD_DECLARE() Binding_CreatedBy {
	return enum_Binding_CreatedBy_GOOG_FORWARD_DECLARE
}

func Binding_CreatedBy_Values() []Binding_CreatedBy {
	return j2gort.CloneSlice(enumValues_Binding_CreatedBy)
}

func Binding_CreatedBy_ValueOf(name string) Binding_CreatedBy {
	if name == j2gort.NullString {
		j2gort.ThrowNPE()
	}
	for i_ := range enumBaseTable_Binding_CreatedBy {
		if enumBaseTable_Binding_CreatedBy[i_].Name_ == name {
			return Binding_CreatedBy(i_ + 1)
		}
	}
	panic(javalang.NewIllegalArgumentException("No enum constant com.google.javascript.jscomp.modules.Binding.CreatedBy." + name))
}

// Java method: CreatedBy.isClosureImport()
func (this_ Binding_CreatedBy) IsClosureImport() bool {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return ((func() bool {
		re_7 := this_
		ae_8 := StaticGet_Binding_CreatedBy_GOOG_REQUIRE()
		if re_7 == 0 {
			j2gort.ThrowNPE()
		}
		return re_7 == ae_8
	}() || func() bool {
		re_9 := this_
		ae_10 := StaticGet_Binding_CreatedBy_GOOG_REQUIRE_TYPE()
		if re_9 == 0 {
			j2gort.ThrowNPE()
		}
		return re_9 == ae_10
	}()) || func() bool {
		re_11 := this_
		ae_12 := StaticGet_Binding_CreatedBy_GOOG_FORWARD_DECLARE()
		if re_11 == 0 {
			j2gort.ThrowNPE()
		}
		return re_11 == ae_12
	}())
}

func (this_ Binding_CreatedBy) IsJavaNil_() bool                              { return this_ == 0 }
func (Binding_CreatedBy) HasIdentityEquals_()                                 {}
func (Binding_CreatedBy) Implements_j2gort_Ext_java_lang_Enum()               {}
func (Binding_CreatedBy) Implements_j2gort_Ext_java_lang_constant_Constable() {}
func (this_ Binding_CreatedBy) GetExt_java_lang_EnumBase() *j2gort.Ext_java_lang_EnumBase {
	return &j2gort.Ext_java_lang_EnumBase{Core_: this_.GetEnumBase()}
}
func (Binding_CreatedBy) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.modules.Binding$CreatedBy", "CreatedBy"
}
func (this_ Binding_CreatedBy) GetEnumBase() *j2gort.EnumBase {
	if this_ == 0 {
		return nil
	}
	return &enumBaseTable_Binding_CreatedBy[this_-1]
}
func (this_ Binding_CreatedBy) Ordinal() int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return int32(this_ - 1)
}
func (this_ Binding_CreatedBy) EnumOrdinal_() int32 { return int32(this_) - 1 }
func (this_ Binding_CreatedBy) Name() string {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_Binding_CreatedBy[this_-1].Name_
}
func (this_ Binding_CreatedBy) String() string {
	if this_ == 0 {
		return "null"
	}
	return enumBaseTable_Binding_CreatedBy[this_-1].Name_
}
func (this_ Binding_CreatedBy) Equals(other_ any) bool {
	o_, ok_ := other_.(Binding_CreatedBy)
	return this_ != 0 && ok_ && this_ == o_
}
func (this_ Binding_CreatedBy) HashCode() int32 {
	if this_ == 0 {
		return 0
	}
	return enumBaseTable_Binding_CreatedBy[this_-1].HashCode()
}
func (this_ Binding_CreatedBy) CompareTo(other_ any) int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_Binding_CreatedBy[this_-1].CompareTo(other_)
}
func (this_ Binding_CreatedBy) WriteGson(out_ *gson.JsonWriter) {
	if this_ == 0 {
		out_.NullValue()
		return
	}
	out_.ValueString(enumBaseTable_Binding_CreatedBy[this_-1].Name_)
}
func init() {
	vals_ := []any{Binding_CreatedBy(1), Binding_CreatedBy(2), Binding_CreatedBy(3), Binding_CreatedBy(4), Binding_CreatedBy(5)}
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.modules.Binding.CreatedBy", vals_)
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.modules.Binding$CreatedBy", vals_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.modules.Binding.CreatedBy", "com.google.javascript.jscomp.modules.Binding$CreatedBy", func(v_ any) bool {
		return j2gort.InstanceOf[Binding_CreatedBy](v_)
	})
}
