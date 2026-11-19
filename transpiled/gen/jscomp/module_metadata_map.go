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
// Port of com/google/javascript/jscomp/modules/ModuleMetadataMap.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
	"github.com/paullepoulpe/jsgo/transpiled/shims/gson"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.jscomp.modules.ModuleMetadataMap
type ModuleMetadataMap struct {
	modulesByPath          j2gort.Ext_com_google_common_collect_ImmutableMap `gson:"modulesByPath"`
	modulesByGoogNamespace j2gort.Ext_com_google_common_collect_ImmutableMap `gson:"modulesByGoogNamespace"`
	moduleMetadata         j2gort.Ext_com_google_common_collect_ImmutableSet `gson:"moduleMetadata"`
	idHash_                atomic.Int32
}

func (this_ *ModuleMetadataMap) BindSelf_(s any) {
	_ = s
}
func (this_ *ModuleMetadataMap) IsJavaNil_() bool { return this_ == nil }
func (*ModuleMetadataMap) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.modules.ModuleMetadataMap", "ModuleMetadataMap"
}
func (this_ *ModuleMetadataMap) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*ModuleMetadataMap) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.modules.ModuleMetadataMap", "com.google.javascript.jscomp.modules.ModuleMetadataMap", func(v_ any) bool {
		return j2gort.InstanceOf[*ModuleMetadataMap](v_)
	})
}

func (this_ *ModuleMetadataMap) initFields_ModuleMetadataMap() {
	_ = this_
}

// Java constructor: ModuleMetadataMap(java.util.Map<java.lang.String,com.google.javascript.jscomp.modules.ModuleMetadataMap.ModuleMetadata> modulesByPath, java.util.Map<java.lang.String,com.google.javascript.jscomp.modules.ModuleMetadataMap.ModuleMetadata> modulesByGoogNamespace)
func NewModuleMetadataMap(modulesByPath j2gort.Ext_java_util_Map, modulesByGoogNamespace j2gort.Ext_java_util_Map) *ModuleMetadataMap {
	this_ := &ModuleMetadataMap{}
	InitModuleMetadataMap(this_, modulesByPath, modulesByGoogNamespace)
	return this_
}

func InitModuleMetadataMap(this_ *ModuleMetadataMap, modulesByPath j2gort.Ext_java_util_Map, modulesByGoogNamespace j2gort.Ext_java_util_Map) {
	_ = modulesByPath
	_ = modulesByGoogNamespace
	this_.initFields_ModuleMetadataMap()
	this_.modulesByPath = j2gort.Ext_com_google_common_collect_ImmutableMap_CopyOf_1_Ext_java_util_Map(modulesByPath)
	this_.modulesByGoogNamespace = j2gort.Ext_com_google_common_collect_ImmutableMap_CopyOf_1_Ext_java_util_Map(modulesByGoogNamespace)
	this_.moduleMetadata = j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableSet](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableSet_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableSet_Builder](j2gort.Ext_com_google_common_collect_ImmutableSet_S_Builder().AddAll_1_Ext_java_lang_Iterable(j2gort.AsIface[j2gort.Ext_java_util_Collection](modulesByPath.Values()))).AddAll_1_Ext_java_lang_Iterable(j2gort.AsIface[j2gort.Ext_java_util_Collection](modulesByGoogNamespace.Values()))).Build())
}

// Java method: ModuleMetadataMap.getModulesByPath()
func (this_ *ModuleMetadataMap) GetModulesByPath() j2gort.Ext_com_google_common_collect_ImmutableMap {
	return this_.modulesByPath
}

// Java method: ModuleMetadataMap.getModulesByGoogNamespace()
func (this_ *ModuleMetadataMap) GetModulesByGoogNamespace() j2gort.Ext_com_google_common_collect_ImmutableMap {
	return this_.modulesByGoogNamespace
}

// Java method: ModuleMetadataMap.getAllModuleMetadata()
func (this_ *ModuleMetadataMap) GetAllModuleMetadata() j2gort.Ext_com_google_common_collect_ImmutableSet {
	return this_.moduleMetadata
}

// Java static method: ModuleMetadataMap.emptyForTesting()
func ModuleMetadataMap_EmptyForTesting() *ModuleMetadataMap {
	return NewModuleMetadataMap(j2gort.AsIface[j2gort.Ext_java_util_Map](j2gort.Ext_com_google_common_collect_ImmutableMap_Of_0()), j2gort.AsIface[j2gort.Ext_java_util_Map](j2gort.Ext_com_google_common_collect_ImmutableMap_Of_0()))
}

// Java enum: com.google.javascript.jscomp.modules.ModuleMetadataMap.ModuleType
type ModuleMetadataMap_ModuleType uint8

const (
	enum_ModuleMetadataMap_ModuleType_ES6_MODULE         ModuleMetadataMap_ModuleType = 1
	enum_ModuleMetadataMap_ModuleType_GOOG_MODULE        ModuleMetadataMap_ModuleType = 2
	enum_ModuleMetadataMap_ModuleType_LEGACY_GOOG_MODULE ModuleMetadataMap_ModuleType = 3
	enum_ModuleMetadataMap_ModuleType_COMMON_JS          ModuleMetadataMap_ModuleType = 4
	enum_ModuleMetadataMap_ModuleType_GOOG_PROVIDE       ModuleMetadataMap_ModuleType = 5
	enum_ModuleMetadataMap_ModuleType_SCRIPT             ModuleMetadataMap_ModuleType = 6
)

var enumBaseTable_ModuleMetadataMap_ModuleType = [...]j2gort.EnumBase{
	{Name_: "ES6_MODULE", Ordinal_: 0},
	{Name_: "GOOG_MODULE", Ordinal_: 1},
	{Name_: "LEGACY_GOOG_MODULE", Ordinal_: 2},
	{Name_: "COMMON_JS", Ordinal_: 3},
	{Name_: "GOOG_PROVIDE", Ordinal_: 4},
	{Name_: "SCRIPT", Ordinal_: 5},
}
var enumValues_ModuleMetadataMap_ModuleType = []ModuleMetadataMap_ModuleType{1, 2, 3, 4, 5, 6}

var enumField_ModuleMetadataMap_ModuleType_description = [...]string{
	"an ES6 module",
	"a goog.module",
	"a goog.module",
	"a CommonJS module",
	"a script file that contains at least one goog.provide",
	"a script file that does not contain a goog.provide",
}

func (this_ ModuleMetadataMap_ModuleType) FieldGet_description() string {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumField_ModuleMetadataMap_ModuleType_description[this_-1]
}

func StaticGet_ModuleMetadataMap_ModuleType_ES6_MODULE() ModuleMetadataMap_ModuleType {
	return enum_ModuleMetadataMap_ModuleType_ES6_MODULE
}
func StaticGet_ModuleMetadataMap_ModuleType_GOOG_MODULE() ModuleMetadataMap_ModuleType {
	return enum_ModuleMetadataMap_ModuleType_GOOG_MODULE
}
func StaticGet_ModuleMetadataMap_ModuleType_LEGACY_GOOG_MODULE() ModuleMetadataMap_ModuleType {
	return enum_ModuleMetadataMap_ModuleType_LEGACY_GOOG_MODULE
}
func StaticGet_ModuleMetadataMap_ModuleType_COMMON_JS() ModuleMetadataMap_ModuleType {
	return enum_ModuleMetadataMap_ModuleType_COMMON_JS
}
func StaticGet_ModuleMetadataMap_ModuleType_GOOG_PROVIDE() ModuleMetadataMap_ModuleType {
	return enum_ModuleMetadataMap_ModuleType_GOOG_PROVIDE
}
func StaticGet_ModuleMetadataMap_ModuleType_SCRIPT() ModuleMetadataMap_ModuleType {
	return enum_ModuleMetadataMap_ModuleType_SCRIPT
}

func ModuleMetadataMap_ModuleType_Values() []ModuleMetadataMap_ModuleType {
	return j2gort.CloneSlice(enumValues_ModuleMetadataMap_ModuleType)
}

func ModuleMetadataMap_ModuleType_ValueOf(name string) ModuleMetadataMap_ModuleType {
	if name == j2gort.NullString {
		j2gort.ThrowNPE()
	}
	for i_ := range enumBaseTable_ModuleMetadataMap_ModuleType {
		if enumBaseTable_ModuleMetadataMap_ModuleType[i_].Name_ == name {
			return ModuleMetadataMap_ModuleType(i_ + 1)
		}
	}
	panic(javalang.NewIllegalArgumentException("No enum constant com.google.javascript.jscomp.modules.ModuleMetadataMap.ModuleType." + name))
}

func (this_ ModuleMetadataMap_ModuleType) IsJavaNil_() bool                              { return this_ == 0 }
func (ModuleMetadataMap_ModuleType) HasIdentityEquals_()                                 {}
func (ModuleMetadataMap_ModuleType) Implements_j2gort_Ext_java_lang_Enum()               {}
func (ModuleMetadataMap_ModuleType) Implements_j2gort_Ext_java_lang_constant_Constable() {}
func (this_ ModuleMetadataMap_ModuleType) GetExt_java_lang_EnumBase() *j2gort.Ext_java_lang_EnumBase {
	return &j2gort.Ext_java_lang_EnumBase{Core_: this_.GetEnumBase()}
}
func (ModuleMetadataMap_ModuleType) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.modules.ModuleMetadataMap$ModuleType", "ModuleType"
}
func (this_ ModuleMetadataMap_ModuleType) GetEnumBase() *j2gort.EnumBase {
	if this_ == 0 {
		return nil
	}
	return &enumBaseTable_ModuleMetadataMap_ModuleType[this_-1]
}
func (this_ ModuleMetadataMap_ModuleType) Ordinal() int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return int32(this_ - 1)
}
func (this_ ModuleMetadataMap_ModuleType) EnumOrdinal_() int32 { return int32(this_) - 1 }
func (this_ ModuleMetadataMap_ModuleType) Name() string {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_ModuleMetadataMap_ModuleType[this_-1].Name_
}
func (this_ ModuleMetadataMap_ModuleType) String() string {
	if this_ == 0 {
		return "null"
	}
	return enumBaseTable_ModuleMetadataMap_ModuleType[this_-1].Name_
}
func (this_ ModuleMetadataMap_ModuleType) Equals(other_ any) bool {
	o_, ok_ := other_.(ModuleMetadataMap_ModuleType)
	return this_ != 0 && ok_ && this_ == o_
}
func (this_ ModuleMetadataMap_ModuleType) HashCode() int32 {
	if this_ == 0 {
		return 0
	}
	return enumBaseTable_ModuleMetadataMap_ModuleType[this_-1].HashCode()
}
func (this_ ModuleMetadataMap_ModuleType) CompareTo(other_ any) int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_ModuleMetadataMap_ModuleType[this_-1].CompareTo(other_)
}
func (this_ ModuleMetadataMap_ModuleType) WriteGson(out_ *gson.JsonWriter) {
	if this_ == 0 {
		out_.NullValue()
		return
	}
	out_.ValueString(enumBaseTable_ModuleMetadataMap_ModuleType[this_-1].Name_)
}
func init() {
	vals_ := []any{ModuleMetadataMap_ModuleType(1), ModuleMetadataMap_ModuleType(2), ModuleMetadataMap_ModuleType(3), ModuleMetadataMap_ModuleType(4), ModuleMetadataMap_ModuleType(5), ModuleMetadataMap_ModuleType(6)}
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.modules.ModuleMetadataMap.ModuleType", vals_)
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.modules.ModuleMetadataMap$ModuleType", vals_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.modules.ModuleMetadataMap.ModuleType", "com.google.javascript.jscomp.modules.ModuleMetadataMap$ModuleType", func(v_ any) bool {
		return j2gort.InstanceOf[ModuleMetadataMap_ModuleType](v_)
	})
}

// Java polymorphic class interface: com.google.javascript.jscomp.modules.ModuleMetadataMap.ModuleMetadata
type ModuleMetadataMap_ModuleMetadata interface {
	Implements_jscomp_ModuleMetadataMap_ModuleMetadata()
	IsJavaNil_() bool
	GetModuleMetadataMap_ModuleMetadataBase() *ModuleMetadataMap_ModuleMetadataBase
	ModuleType() ModuleMetadataMap_ModuleType
	IsEs6Module() bool
	IsGoogModule() bool
	IsNonLegacyGoogModule() bool
	IsLegacyGoogModule() bool
	IsGoogProvide() bool
	HasLegacyGoogNamespaces() bool
	IsCommonJs() bool
	IsNonProvideScript() bool
	IsModule() bool
	RootNode() pkg_rhino.Node
	UsesClosure() bool
	IsTestOnly() bool
	GoogNamespaces() j2gort.Ext_com_google_common_collect_ImmutableMultiset
	StronglyRequiredGoogNamespaces() j2gort.Ext_com_google_common_collect_ImmutableMultiset
	DynamicallyRequiredGoogNamespaces() j2gort.Ext_com_google_common_collect_ImmutableMultiset
	MaybeRequiredGoogNamespaces() j2gort.Ext_com_google_common_collect_ImmutableMultiset
	WeaklyRequiredGoogNamespaces() j2gort.Ext_com_google_common_collect_ImmutableMultiset
	Es6ImportSpecifiers() j2gort.Ext_com_google_common_collect_ImmutableMultiset
	NestedModules() j2gort.Ext_com_google_common_collect_ImmutableList
	ReadToggles() j2gort.Ext_com_google_common_collect_ImmutableMultiset
	Path() *ModuleLoader_ModulePath
	Equals(any) bool
	HashCode() int32
}

func ToIface_ModuleMetadataMap_ModuleMetadata[T interface {
	ModuleMetadataMap_ModuleMetadata
	comparable
}](v T) ModuleMetadataMap_ModuleMetadata {
	var z T
	if v == z {
		return nil
	}
	return v
}

// Java class: com.google.javascript.jscomp.modules.ModuleMetadataMap.ModuleMetadata
type ModuleMetadataMap_ModuleMetadataBase struct {
	self_ ModuleMetadataMap_ModuleMetadata
}

func (this_ *ModuleMetadataMap_ModuleMetadataBase) BindSelf_(s any) {
	_ = s
	this_.self_ = s.(ModuleMetadataMap_ModuleMetadata)
}
func (this_ *ModuleMetadataMap_ModuleMetadataBase) IsJavaNil_() bool { return this_ == nil }
func (*ModuleMetadataMap_ModuleMetadataBase) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.modules.ModuleMetadataMap$ModuleMetadata", "ModuleMetadata"
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.modules.ModuleMetadataMap.ModuleMetadata", "com.google.javascript.jscomp.modules.ModuleMetadataMap$ModuleMetadata", func(v_ any) bool {
		return j2gort.InstanceOf[ModuleMetadataMap_ModuleMetadata](v_)
	})
}
func (b *ModuleMetadataMap_ModuleMetadataBase) GetModuleMetadataMap_ModuleMetadataBase() *ModuleMetadataMap_ModuleMetadataBase {
	return b
}
func (*ModuleMetadataMap_ModuleMetadataBase) Implements_jscomp_ModuleMetadataMap_ModuleMetadata() {}

func (this_ *ModuleMetadataMap_ModuleMetadataBase) initFields_ModuleMetadataMap_ModuleMetadata() {
	_ = this_
}

func InitModuleMetadataMap_ModuleMetadata(this_ *ModuleMetadataMap_ModuleMetadataBase) {
	this_.initFields_ModuleMetadataMap_ModuleMetadata()
}

// Java method: ModuleMetadata.moduleType()
func (this_ *ModuleMetadataMap_ModuleMetadataBase) ModuleType() ModuleMetadataMap_ModuleType {
	panic("abstract or native method invoked")
}

// Java method: ModuleMetadata.isEs6Module()
func (this_ *ModuleMetadataMap_ModuleMetadataBase) IsEs6Module() bool {
	return (this_.self_.ModuleType() == StaticGet_ModuleMetadataMap_ModuleType_ES6_MODULE())
}

// Java method: ModuleMetadata.isGoogModule()
func (this_ *ModuleMetadataMap_ModuleMetadataBase) IsGoogModule() bool {
	return (this_.IsNonLegacyGoogModule() || this_.IsLegacyGoogModule())
}

// Java method: ModuleMetadata.isNonLegacyGoogModule()
func (this_ *ModuleMetadataMap_ModuleMetadataBase) IsNonLegacyGoogModule() bool {
	return (this_.self_.ModuleType() == StaticGet_ModuleMetadataMap_ModuleType_GOOG_MODULE())
}

// Java method: ModuleMetadata.isLegacyGoogModule()
func (this_ *ModuleMetadataMap_ModuleMetadataBase) IsLegacyGoogModule() bool {
	return (this_.self_.ModuleType() == StaticGet_ModuleMetadataMap_ModuleType_LEGACY_GOOG_MODULE())
}

// Java method: ModuleMetadata.isGoogProvide()
func (this_ *ModuleMetadataMap_ModuleMetadataBase) IsGoogProvide() bool {
	return (this_.self_.ModuleType() == StaticGet_ModuleMetadataMap_ModuleType_GOOG_PROVIDE())
}

// Java method: ModuleMetadata.hasLegacyGoogNamespaces()
func (this_ *ModuleMetadataMap_ModuleMetadataBase) HasLegacyGoogNamespaces() bool {
	return (this_.IsGoogProvide() || this_.IsLegacyGoogModule())
}

// Java method: ModuleMetadata.isCommonJs()
func (this_ *ModuleMetadataMap_ModuleMetadataBase) IsCommonJs() bool {
	return (this_.self_.ModuleType() == StaticGet_ModuleMetadataMap_ModuleType_COMMON_JS())
}

// Java method: ModuleMetadata.isNonProvideScript()
func (this_ *ModuleMetadataMap_ModuleMetadataBase) IsNonProvideScript() bool {
	return (this_.self_.ModuleType() == StaticGet_ModuleMetadataMap_ModuleType_SCRIPT())
}

// Java method: ModuleMetadata.isModule()
func (this_ *ModuleMetadataMap_ModuleMetadataBase) IsModule() bool {
	switch (this_.self_.ModuleType()).Ordinal() {
	case int32(4) /* ModuleType.GOOG_PROVIDE */, int32(5) /* ModuleType.SCRIPT */ :
		return false
	case int32(3) /* ModuleType.COMMON_JS */, int32(0) /* ModuleType.ES6_MODULE */, int32(1) /* ModuleType.GOOG_MODULE */, int32(2) /* ModuleType.LEGACY_GOOG_MODULE */ :
		return true
	}
	panic("unreachable switch expression")
}

// Java method: ModuleMetadata.rootNode()
func (this_ *ModuleMetadataMap_ModuleMetadataBase) RootNode() pkg_rhino.Node {
	panic("abstract or native method invoked")
}

// Java method: ModuleMetadata.usesClosure()
func (this_ *ModuleMetadataMap_ModuleMetadataBase) UsesClosure() bool {
	panic("abstract or native method invoked")
}

// Java method: ModuleMetadata.isTestOnly()
func (this_ *ModuleMetadataMap_ModuleMetadataBase) IsTestOnly() bool {
	panic("abstract or native method invoked")
}

// Java method: ModuleMetadata.googNamespaces()
func (this_ *ModuleMetadataMap_ModuleMetadataBase) GoogNamespaces() j2gort.Ext_com_google_common_collect_ImmutableMultiset {
	panic("abstract or native method invoked")
}

// Java method: ModuleMetadata.stronglyRequiredGoogNamespaces()
func (this_ *ModuleMetadataMap_ModuleMetadataBase) StronglyRequiredGoogNamespaces() j2gort.Ext_com_google_common_collect_ImmutableMultiset {
	panic("abstract or native method invoked")
}

// Java method: ModuleMetadata.dynamicallyRequiredGoogNamespaces()
func (this_ *ModuleMetadataMap_ModuleMetadataBase) DynamicallyRequiredGoogNamespaces() j2gort.Ext_com_google_common_collect_ImmutableMultiset {
	panic("abstract or native method invoked")
}

// Java method: ModuleMetadata.maybeRequiredGoogNamespaces()
func (this_ *ModuleMetadataMap_ModuleMetadataBase) MaybeRequiredGoogNamespaces() j2gort.Ext_com_google_common_collect_ImmutableMultiset {
	panic("abstract or native method invoked")
}

// Java method: ModuleMetadata.weaklyRequiredGoogNamespaces()
func (this_ *ModuleMetadataMap_ModuleMetadataBase) WeaklyRequiredGoogNamespaces() j2gort.Ext_com_google_common_collect_ImmutableMultiset {
	panic("abstract or native method invoked")
}

// Java method: ModuleMetadata.es6ImportSpecifiers()
func (this_ *ModuleMetadataMap_ModuleMetadataBase) Es6ImportSpecifiers() j2gort.Ext_com_google_common_collect_ImmutableMultiset {
	panic("abstract or native method invoked")
}

// Java method: ModuleMetadata.nestedModules()
func (this_ *ModuleMetadataMap_ModuleMetadataBase) NestedModules() j2gort.Ext_com_google_common_collect_ImmutableList {
	panic("abstract or native method invoked")
}

// Java method: ModuleMetadata.readToggles()
func (this_ *ModuleMetadataMap_ModuleMetadataBase) ReadToggles() j2gort.Ext_com_google_common_collect_ImmutableMultiset {
	panic("abstract or native method invoked")
}

// Java method: ModuleMetadata.path()
func (this_ *ModuleMetadataMap_ModuleMetadataBase) Path() *ModuleLoader_ModulePath {
	panic("abstract or native method invoked")
}

// Java static method: ModuleMetadata.builder()
func ModuleMetadataMap_ModuleMetadata_S_Builder() ModuleMetadataMap_ModuleMetadata_Builder {
	return NewAutoValue_ModuleMetadataMap_ModuleMetadata_Builder()
}

// Java method: ModuleMetadata.equals(java.lang.Object other)
func (this_ *ModuleMetadataMap_ModuleMetadataBase) Equals(arg_1 any) bool {
	other := arg_1
	_ = other
	return j2gort.RefEquals(this_.self_, other)
}

// Java method: ModuleMetadata.hashCode()
func (this_ *ModuleMetadataMap_ModuleMetadataBase) HashCode() int32 {
	return javalang.IdentityHashCode(this_.self_)
}

// Java polymorphic class interface: com.google.javascript.jscomp.modules.ModuleMetadataMap.ModuleMetadata.Builder
type ModuleMetadataMap_ModuleMetadata_Builder interface {
	Implements_jscomp_ModuleMetadataMap_ModuleMetadata_Builder()
	IsJavaNil_() bool
	GetModuleMetadataMap_ModuleMetadata_BuilderBase() *ModuleMetadataMap_ModuleMetadata_BuilderBase
	Build() ModuleMetadataMap_ModuleMetadata
	GoogNamespacesBuilder() *j2gort.Ext_com_google_common_collect_ImmutableMultiset_Builder
	AddGoogNamespace(string) ModuleMetadataMap_ModuleMetadata_Builder
	StronglyRequiredGoogNamespacesBuilder() *j2gort.Ext_com_google_common_collect_ImmutableMultiset_Builder
	DynamicallyRequiredGoogNamespacesBuilder() *j2gort.Ext_com_google_common_collect_ImmutableMultiset_Builder
	MaybeRequiredGoogNamespacesBuilder() *j2gort.Ext_com_google_common_collect_ImmutableMultiset_Builder
	WeaklyRequiredGoogNamespacesBuilder() *j2gort.Ext_com_google_common_collect_ImmutableMultiset_Builder
	Es6ImportSpecifiersBuilder() *j2gort.Ext_com_google_common_collect_ImmutableMultiset_Builder
	NestedModulesBuilder() *j2gort.Ext_com_google_common_collect_ImmutableList_Builder
	ReadTogglesBuilder() *j2gort.Ext_com_google_common_collect_ImmutableMultiset_Builder
	Path(*ModuleLoader_ModulePath) ModuleMetadataMap_ModuleMetadata_Builder
	UsesClosure(bool) ModuleMetadataMap_ModuleMetadata_Builder
	IsTestOnly(bool) ModuleMetadataMap_ModuleMetadata_Builder
	ModuleType_0() ModuleMetadataMap_ModuleType
	ModuleType_1(ModuleMetadataMap_ModuleType) ModuleMetadataMap_ModuleMetadata_Builder
	RootNode(pkg_rhino.Node) ModuleMetadataMap_ModuleMetadata_Builder
}

func ToIface_ModuleMetadataMap_ModuleMetadata_Builder[T interface {
	ModuleMetadataMap_ModuleMetadata_Builder
	comparable
}](v T) ModuleMetadataMap_ModuleMetadata_Builder {
	var z T
	if v == z {
		return nil
	}
	return v
}

// Java class: com.google.javascript.jscomp.modules.ModuleMetadataMap.ModuleMetadata.Builder
type ModuleMetadataMap_ModuleMetadata_BuilderBase struct {
	self_   ModuleMetadataMap_ModuleMetadata_Builder
	idHash_ atomic.Int32
}

func (this_ *ModuleMetadataMap_ModuleMetadata_BuilderBase) BindSelf_(s any) {
	_ = s
	this_.self_ = s.(ModuleMetadataMap_ModuleMetadata_Builder)
}
func (this_ *ModuleMetadataMap_ModuleMetadata_BuilderBase) IsJavaNil_() bool { return this_ == nil }
func (*ModuleMetadataMap_ModuleMetadata_BuilderBase) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.modules.ModuleMetadataMap$ModuleMetadata$Builder", "Builder"
}
func (this_ *ModuleMetadataMap_ModuleMetadata_BuilderBase) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.modules.ModuleMetadataMap.ModuleMetadata.Builder", "com.google.javascript.jscomp.modules.ModuleMetadataMap$ModuleMetadata$Builder", func(v_ any) bool {
		return j2gort.InstanceOf[ModuleMetadataMap_ModuleMetadata_Builder](v_)
	})
}
func (b *ModuleMetadataMap_ModuleMetadata_BuilderBase) GetModuleMetadataMap_ModuleMetadata_BuilderBase() *ModuleMetadataMap_ModuleMetadata_BuilderBase {
	return b
}
func (*ModuleMetadataMap_ModuleMetadata_BuilderBase) Implements_jscomp_ModuleMetadataMap_ModuleMetadata_Builder() {
}

func (this_ *ModuleMetadataMap_ModuleMetadata_BuilderBase) initFields_ModuleMetadataMap_ModuleMetadata_Builder() {
	_ = this_
}

func InitModuleMetadataMap_ModuleMetadata_Builder(this_ *ModuleMetadataMap_ModuleMetadata_BuilderBase) {
	this_.initFields_ModuleMetadataMap_ModuleMetadata_Builder()
}

// Java method: Builder.build()
func (this_ *ModuleMetadataMap_ModuleMetadata_BuilderBase) Build() ModuleMetadataMap_ModuleMetadata {
	panic("abstract or native method invoked")
}

// Java method: Builder.googNamespacesBuilder()
func (this_ *ModuleMetadataMap_ModuleMetadata_BuilderBase) GoogNamespacesBuilder() *j2gort.Ext_com_google_common_collect_ImmutableMultiset_Builder {
	panic("abstract or native method invoked")
}

// Java method: Builder.addGoogNamespace(java.lang.String namespace)
func (this_ *ModuleMetadataMap_ModuleMetadata_BuilderBase) AddGoogNamespace(namespace string) ModuleMetadataMap_ModuleMetadata_Builder {
	_ = namespace
	this_.self_.GoogNamespacesBuilder().Add_1_any(j2gort.BoxString(namespace))
	return this_.self_
}

// Java method: Builder.stronglyRequiredGoogNamespacesBuilder()
func (this_ *ModuleMetadataMap_ModuleMetadata_BuilderBase) StronglyRequiredGoogNamespacesBuilder() *j2gort.Ext_com_google_common_collect_ImmutableMultiset_Builder {
	panic("abstract or native method invoked")
}

// Java method: Builder.dynamicallyRequiredGoogNamespacesBuilder()
func (this_ *ModuleMetadataMap_ModuleMetadata_BuilderBase) DynamicallyRequiredGoogNamespacesBuilder() *j2gort.Ext_com_google_common_collect_ImmutableMultiset_Builder {
	panic("abstract or native method invoked")
}

// Java method: Builder.maybeRequiredGoogNamespacesBuilder()
func (this_ *ModuleMetadataMap_ModuleMetadata_BuilderBase) MaybeRequiredGoogNamespacesBuilder() *j2gort.Ext_com_google_common_collect_ImmutableMultiset_Builder {
	panic("abstract or native method invoked")
}

// Java method: Builder.weaklyRequiredGoogNamespacesBuilder()
func (this_ *ModuleMetadataMap_ModuleMetadata_BuilderBase) WeaklyRequiredGoogNamespacesBuilder() *j2gort.Ext_com_google_common_collect_ImmutableMultiset_Builder {
	panic("abstract or native method invoked")
}

// Java method: Builder.es6ImportSpecifiersBuilder()
func (this_ *ModuleMetadataMap_ModuleMetadata_BuilderBase) Es6ImportSpecifiersBuilder() *j2gort.Ext_com_google_common_collect_ImmutableMultiset_Builder {
	panic("abstract or native method invoked")
}

// Java method: Builder.nestedModulesBuilder()
func (this_ *ModuleMetadataMap_ModuleMetadata_BuilderBase) NestedModulesBuilder() *j2gort.Ext_com_google_common_collect_ImmutableList_Builder {
	panic("abstract or native method invoked")
}

// Java method: Builder.readTogglesBuilder()
func (this_ *ModuleMetadataMap_ModuleMetadata_BuilderBase) ReadTogglesBuilder() *j2gort.Ext_com_google_common_collect_ImmutableMultiset_Builder {
	panic("abstract or native method invoked")
}

// Java method: Builder.path(com.google.javascript.jscomp.deps.ModuleLoader.@org.jspecify.annotations.Nullable ModulePath value)
func (this_ *ModuleMetadataMap_ModuleMetadata_BuilderBase) Path(value *ModuleLoader_ModulePath) ModuleMetadataMap_ModuleMetadata_Builder {
	_ = value
	panic("abstract or native method invoked")
}

// Java method: Builder.usesClosure(boolean value)
func (this_ *ModuleMetadataMap_ModuleMetadata_BuilderBase) UsesClosure(value bool) ModuleMetadataMap_ModuleMetadata_Builder {
	_ = value
	panic("abstract or native method invoked")
}

// Java method: Builder.isTestOnly(boolean value)
func (this_ *ModuleMetadataMap_ModuleMetadata_BuilderBase) IsTestOnly(value bool) ModuleMetadataMap_ModuleMetadata_Builder {
	_ = value
	panic("abstract or native method invoked")
}

// Java method: Builder.moduleType()
func (this_ *ModuleMetadataMap_ModuleMetadata_BuilderBase) ModuleType_0() ModuleMetadataMap_ModuleType {
	panic("abstract or native method invoked")
}

// Java method: Builder.moduleType(com.google.javascript.jscomp.modules.ModuleMetadataMap.ModuleType value)
func (this_ *ModuleMetadataMap_ModuleMetadata_BuilderBase) ModuleType_1(value ModuleMetadataMap_ModuleType) ModuleMetadataMap_ModuleMetadata_Builder {
	_ = value
	panic("abstract or native method invoked")
}

// Java method: Builder.rootNode(com.google.javascript.rhino.@org.jspecify.annotations.Nullable Node root)
func (this_ *ModuleMetadataMap_ModuleMetadata_BuilderBase) RootNode(root pkg_rhino.Node) ModuleMetadataMap_ModuleMetadata_Builder {
	_ = root
	panic("abstract or native method invoked")
}
