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
// Port of com/google/javascript/jscomp/modules/Export.java (generated mechanically by j2go).

package jscomp

import (
	"fmt"
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java record: com.google.javascript.jscomp.modules.Export
type Export struct {
	exportName       string                           `gson:"exportName"`
	moduleRequest    string                           `gson:"moduleRequest"`
	importName       string                           `gson:"importName"`
	localName        string                           `gson:"localName"`
	modulePath       *ModuleLoader_ModulePath         `gson:"modulePath"`
	moduleMetadata   ModuleMetadataMap_ModuleMetadata `gson:"moduleMetadata"`
	closureNamespace string                           `gson:"closureNamespace"`
	exportNode       pkg_rhino.Node                   `gson:"exportNode"`
	nameNode         pkg_rhino.Node                   `gson:"nameNode"`
	mutated          bool                             `gson:"mutated"`
}

func (this_ *Export) BindSelf_(s any) {
	_ = s
}
func (this_ *Export) IsJavaNil_() bool { return this_ == nil }
func (*Export) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.modules.Export", "Export"
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.modules.Export", "com.google.javascript.jscomp.modules.Export", func(v_ any) bool {
		return j2gort.InstanceOf[*Export](v_)
	})
}
func (*Export) Implements_j2gort_Ext_java_lang_Record() {}

func StaticGet_Export_DEFAULT_EXPORT_NAME() string { return "*default*" }
func StaticGet_Export_DEFAULT() string             { return "default" }
func StaticGet_Export_NAMESPACE() string           { return "*exports*" }
func (this_ *Export) initFields_Export() {
	_ = this_
}

// Java constructor: Export(java.lang.@org.jspecify.annotations.Nullable String exportName, java.lang.@org.jspecify.annotations.Nullable String moduleRequest, java.lang.@org.jspecify.annotations.Nullable String importName, java.lang.@org.jspecify.annotations.Nullable String localName, com.google.javascript.jscomp.deps.ModuleLoader.@org.jspecify.annotations.Nullable ModulePath modulePath, com.google.javascript.rhino.@org.jspecify.annotations.Nullable Node exportNode, com.google.javascript.rhino.@org.jspecify.annotations.Nullable Node nameNode, com.google.javascript.jscomp.modules.ModuleMetadataMap.ModuleMetadata moduleMetadata, java.lang.@org.jspecify.annotations.Nullable String closureNamespace, boolean mutated)
func NewExport(exportName string, moduleRequest string, importName string, localName string, modulePath *ModuleLoader_ModulePath, exportNode pkg_rhino.Node, nameNode pkg_rhino.Node, moduleMetadata ModuleMetadataMap_ModuleMetadata, closureNamespace string, mutated bool) *Export {
	this_ := &Export{}
	InitExport(this_, exportName, moduleRequest, importName, localName, modulePath, exportNode, nameNode, moduleMetadata, closureNamespace, mutated)
	return this_
}

func InitExport(this_ *Export, exportName string, moduleRequest string, importName string, localName string, modulePath *ModuleLoader_ModulePath, exportNode pkg_rhino.Node, nameNode pkg_rhino.Node, moduleMetadata ModuleMetadataMap_ModuleMetadata, closureNamespace string, mutated bool) {
	_ = exportName
	_ = moduleRequest
	_ = importName
	_ = localName
	_ = modulePath
	_ = exportNode
	_ = nameNode
	_ = moduleMetadata
	_ = closureNamespace
	_ = mutated
	this_.initFields_Export()
	j2gort.CheckNotNullIface(moduleMetadata)
	this_.exportName = exportName
	this_.moduleRequest = moduleRequest
	this_.importName = importName
	this_.localName = localName
	this_.modulePath = modulePath
	this_.exportNode = exportNode
	this_.nameNode = nameNode
	this_.moduleMetadata = moduleMetadata
	this_.closureNamespace = closureNamespace
	this_.mutated = mutated
}

// Java static method: Export.builder()
func Export_S_Builder() Export_Builder {
	return NewAutoBuilder_Export_Builder().Mutated(false)
}

// Java method: Export.toBuilder()
func (this_ *Export) ToBuilder() Export_Builder {
	return NewAutoBuilder_Export_Builder_1(this_)
}

// Java method: Export.mutatedCopy()
func (this_ *Export) MutatedCopy() *Export {
	return this_.ToBuilder().Mutated(true).AutoBuild()
}

func (this_ *Export) ExportName() string                               { return this_.exportName }
func (this_ *Export) ModuleRequest() string                            { return this_.moduleRequest }
func (this_ *Export) ImportName() string                               { return this_.importName }
func (this_ *Export) LocalName() string                                { return this_.localName }
func (this_ *Export) ModulePath() *ModuleLoader_ModulePath             { return this_.modulePath }
func (this_ *Export) ExportNode() pkg_rhino.Node                       { return this_.exportNode }
func (this_ *Export) NameNode() pkg_rhino.Node                         { return this_.nameNode }
func (this_ *Export) ModuleMetadata() ModuleMetadataMap_ModuleMetadata { return this_.moduleMetadata }
func (this_ *Export) ClosureNamespace() string                         { return this_.closureNamespace }
func (this_ *Export) Mutated() bool                                    { return this_.mutated }
func (this_ *Export) Equals(other_ any) bool {
	o_, ok_ := other_.(*Export)
	if !ok_ || o_ == nil {
		return false
	}
	if this_ == o_ {
		return true
	}
	if this_.exportName != o_.exportName {
		return false
	}
	if this_.moduleRequest != o_.moduleRequest {
		return false
	}
	if this_.importName != o_.importName {
		return false
	}
	if this_.localName != o_.localName {
		return false
	}
	if this_.modulePath != o_.modulePath {
		return false
	}
	if this_.exportNode != o_.exportNode {
		return false
	}
	if this_.nameNode != o_.nameNode {
		return false
	}
	if !javalang.ObjectsEquals(this_.moduleMetadata, o_.moduleMetadata) {
		return false
	}
	if this_.closureNamespace != o_.closureNamespace {
		return false
	}
	if this_.mutated != o_.mutated {
		return false
	}
	return true
}
func (this_ *Export) HashCode() int32 {
	var h_ int32 = 0
	{
		hs_1 := this_.exportName
		h_ = 31 * h_
		if hs_1 != j2gort.NullString {
			h_ += javalang.StringHashCode(hs_1)
		}
	}
	{
		hs_2 := this_.moduleRequest
		h_ = 31 * h_
		if hs_2 != j2gort.NullString {
			h_ += javalang.StringHashCode(hs_2)
		}
	}
	{
		hs_3 := this_.importName
		h_ = 31 * h_
		if hs_3 != j2gort.NullString {
			h_ += javalang.StringHashCode(hs_3)
		}
	}
	{
		hs_4 := this_.localName
		h_ = 31 * h_
		if hs_4 != j2gort.NullString {
			h_ += javalang.StringHashCode(hs_4)
		}
	}
	{
		hp_5 := this_.modulePath
		h_ = 31 * h_
		if hp_5 != nil {
			h_ += j2gort.IdentityHash(hp_5)
		}
	}
	h_ = 31*h_ + int32(this_.exportNode)
	h_ = 31*h_ + int32(this_.nameNode)
	{
		hi_6 := this_.moduleMetadata
		h_ = 31 * h_
		if hi_6 != nil && !hi_6.IsJavaNil_() {
			h_ += hi_6.GetModuleMetadataMap_ModuleMetadataBase().HashCode()
		}
	}
	{
		hs_7 := this_.closureNamespace
		h_ = 31 * h_
		if hs_7 != j2gort.NullString {
			h_ += javalang.StringHashCode(hs_7)
		}
	}
	if this_.mutated {
		h_ = 31*h_ + 1231
	} else {
		h_ = 31*h_ + 1237
	}
	return h_
}
func (this_ *Export) String() string {
	_ = fmt.Sprint
	return "Export[" + "exportName=" + j2gort.StrVal(j2gort.BoxString(this_.exportName)) + ", " + "moduleRequest=" + j2gort.StrVal(j2gort.BoxString(this_.moduleRequest)) + ", " + "importName=" + j2gort.StrVal(j2gort.BoxString(this_.importName)) + ", " + "localName=" + j2gort.StrVal(j2gort.BoxString(this_.localName)) + ", " + "modulePath=" + j2gort.StrVal(j2gort.PtrToAny(this_.modulePath)) + ", " + "exportNode=" + j2gort.StrVal(j2gort.BoxNode(this_.exportNode)) + ", " + "nameNode=" + j2gort.StrVal(j2gort.BoxNode(this_.nameNode)) + ", " + "moduleMetadata=" + j2gort.StrVal(this_.moduleMetadata) + ", " + "closureNamespace=" + j2gort.StrVal(j2gort.BoxString(this_.closureNamespace)) + ", " + "mutated=" + j2gort.StrVal(j2gort.ValueOfBoolean(this_.mutated)) + "]"
}

// Java polymorphic class interface: com.google.javascript.jscomp.modules.Export.Builder
type Export_Builder interface {
	Implements_jscomp_Export_Builder()
	IsJavaNil_() bool
	GetExport_BuilderBase() *Export_BuilderBase
	ExportName(string) Export_Builder
	ModuleRequest(string) Export_Builder
	ImportName(string) Export_Builder
	LocalName(string) Export_Builder
	ModulePath(*ModuleLoader_ModulePath) Export_Builder
	ExportNode(pkg_rhino.Node) Export_Builder
	NameNode(pkg_rhino.Node) Export_Builder
	ModuleMetadata(ModuleMetadataMap_ModuleMetadata) Export_Builder
	ClosureNamespace(string) Export_Builder
	Mutated(bool) Export_Builder
	AutoBuild() *Export
	Build() *Export
}

func ToIface_Export_Builder[T interface {
	Export_Builder
	comparable
}](v T) Export_Builder {
	var z T
	if v == z {
		return nil
	}
	return v
}

func Upcast_Export_Builder[T Export_Builder](v T) Export_Builder {
	var vi Export_Builder = v
	if vi == nil || vi.IsJavaNil_() {
		return nil
	}
	return vi
}

// Java class: com.google.javascript.jscomp.modules.Export.Builder
type Export_BuilderBase struct {
	self_   Export_Builder
	idHash_ atomic.Int32
}

func (this_ *Export_BuilderBase) BindSelf_(s any) {
	_ = s
	this_.self_ = s.(Export_Builder)
}
func (this_ *Export_BuilderBase) IsJavaNil_() bool { return this_ == nil }
func (*Export_BuilderBase) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.modules.Export$Builder", "Builder"
}
func (this_ *Export_BuilderBase) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.modules.Export.Builder", "com.google.javascript.jscomp.modules.Export$Builder", func(v_ any) bool {
		return j2gort.InstanceOf[Export_Builder](v_)
	})
}
func (b *Export_BuilderBase) GetExport_BuilderBase() *Export_BuilderBase { return b }
func (*Export_BuilderBase) Implements_jscomp_Export_Builder()            {}

func (this_ *Export_BuilderBase) initFields_Export_Builder() {
	_ = this_
}

func InitExport_Builder(this_ *Export_BuilderBase) {
	this_.initFields_Export_Builder()
}

// Java method: Builder.exportName(java.lang.@org.jspecify.annotations.Nullable String value)
func (this_ *Export_BuilderBase) ExportName(value string) Export_Builder {
	_ = value
	panic("abstract or native method invoked")
}

// Java method: Builder.moduleRequest(java.lang.@org.jspecify.annotations.Nullable String value)
func (this_ *Export_BuilderBase) ModuleRequest(value string) Export_Builder {
	_ = value
	panic("abstract or native method invoked")
}

// Java method: Builder.importName(java.lang.@org.jspecify.annotations.Nullable String value)
func (this_ *Export_BuilderBase) ImportName(value string) Export_Builder {
	_ = value
	panic("abstract or native method invoked")
}

// Java method: Builder.localName(java.lang.@org.jspecify.annotations.Nullable String value)
func (this_ *Export_BuilderBase) LocalName(value string) Export_Builder {
	_ = value
	panic("abstract or native method invoked")
}

// Java method: Builder.modulePath(com.google.javascript.jscomp.deps.ModuleLoader.@org.jspecify.annotations.Nullable ModulePath value)
func (this_ *Export_BuilderBase) ModulePath(value *ModuleLoader_ModulePath) Export_Builder {
	_ = value
	panic("abstract or native method invoked")
}

// Java method: Builder.exportNode(com.google.javascript.rhino.@org.jspecify.annotations.Nullable Node value)
func (this_ *Export_BuilderBase) ExportNode(value pkg_rhino.Node) Export_Builder {
	_ = value
	panic("abstract or native method invoked")
}

// Java method: Builder.nameNode(com.google.javascript.rhino.@org.jspecify.annotations.Nullable Node value)
func (this_ *Export_BuilderBase) NameNode(value pkg_rhino.Node) Export_Builder {
	_ = value
	panic("abstract or native method invoked")
}

// Java method: Builder.moduleMetadata(com.google.javascript.jscomp.modules.ModuleMetadataMap.ModuleMetadata value)
func (this_ *Export_BuilderBase) ModuleMetadata(value ModuleMetadataMap_ModuleMetadata) Export_Builder {
	_ = value
	panic("abstract or native method invoked")
}

// Java method: Builder.closureNamespace(java.lang.@org.jspecify.annotations.Nullable String value)
func (this_ *Export_BuilderBase) ClosureNamespace(value string) Export_Builder {
	_ = value
	panic("abstract or native method invoked")
}

// Java method: Builder.mutated(boolean value)
func (this_ *Export_BuilderBase) Mutated(value bool) Export_Builder {
	_ = value
	panic("abstract or native method invoked")
}

// Java method: Builder.autoBuild()
func (this_ *Export_BuilderBase) AutoBuild() *Export {
	panic("abstract or native method invoked")
}

// Java method: Builder.build()
func (this_ *Export_BuilderBase) Build() *Export {
	var e *Export = this_.self_.AutoBuild()
	_ = e
	if e.ModuleMetadata().GetModuleMetadataMap_ModuleMetadataBase().IsEs6Module() {
		this_.priv_Builder_validateEsModule_6774(e)
	} else if e.ModuleMetadata().GetModuleMetadataMap_ModuleMetadataBase().IsGoogModule() {
		Export_Builder_ValidateGoogModule(e)
	} else {
		Export_Builder_ValidateOtherModule(e)
	}
	return e
}

// Java method: Builder.validateEsModule(com.google.javascript.jscomp.modules.Export e)
func (this_ *Export_BuilderBase) priv_Builder_validateEsModule_6774(e *Export) {
	_ = e
	if !(e.ClosureNamespace() == j2gort.NullString) {
		j2gort.CheckState(false)
	}
	{
		var iifeCond_1 bool
		{

			sr_1 := "*"
			se_2 := e.ImportName()
			iifeCond_1 = j2gort.CheckNotNullString(sr_1) == se_2
		}
		if !(!(iifeCond_1) || (((e.ModuleRequest() != j2gort.NullString) && (e.ExportName() == j2gort.NullString)) && (e.LocalName() == j2gort.NullString))) {
			j2gort.CheckStateStr(false, "Star exports should not have exported / local names.")
		}
	}
	if !((e.LocalName() == j2gort.NullString) || (e.ModuleRequest() == j2gort.NullString)) {
		j2gort.CheckStateStr(false, "Local exports should not have module requests.")
	}
	if !((e.ModuleRequest() == j2gort.NullString) || (e.LocalName() == j2gort.NullString)) {
		j2gort.CheckStateStr(false, "Reexports should not have local names.")
	}
	if !((e.ModuleRequest() == j2gort.NullString) || (e.ImportName() != j2gort.NullString)) {
		j2gort.CheckStateStr(false, "Reexports should have import names.")
	}
	if !((e.ImportName() == j2gort.NullString) || (e.ModuleRequest() != j2gort.NullString)) {
		j2gort.CheckStateStr(false, "Exports with an import name should be a reexport.")
	}
}

// Java static method: Builder.validateGoogModule(com.google.javascript.jscomp.modules.Export e)
func Export_Builder_ValidateGoogModule(e *Export) {
	_ = e
	if !(e.ClosureNamespace() != j2gort.NullString) {
		j2gort.CheckStateStr(false, "Exports should be associated with a namespace")
	}
	if !(e.ExportName() != j2gort.NullString) {
		j2gort.CheckStateStr(false, "Exports should be named")
	}
	if !(e.ExportNode() != 0) {
		j2gort.CheckStateStr(false, "Exports should have a node")
	}
	if !(e.LocalName() == j2gort.NullString) {
		j2gort.CheckStateStr(false, "goog.module Exports don't set a localName")
	}
	if !(e.ModuleRequest() == j2gort.NullString) {
		j2gort.CheckStateStr(false, "goog modules cannot export from other modules")
	}
}

// Java static method: Builder.validateOtherModule(com.google.javascript.jscomp.modules.Export e)
func Export_Builder_ValidateOtherModule(e *Export) {
	_ = e
	j2gort.CheckNotNullString(e.ExportName())
	if !(e.ExportNode() == 0) {
		j2gort.CheckState(false)
	}
	if !(e.LocalName() == j2gort.NullString) {
		j2gort.CheckState(false)
	}
	if !(e.ModuleRequest() == j2gort.NullString) {
		j2gort.CheckState(false)
	}
	if !(e.ImportName() == j2gort.NullString) {
		j2gort.CheckState(false)
	}
	if !(e.NameNode() == 0) {
		j2gort.CheckState(false)
	}
}
