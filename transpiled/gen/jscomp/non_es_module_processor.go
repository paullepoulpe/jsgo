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
// Port of com/google/javascript/jscomp/modules/NonEsModuleProcessor.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.jscomp.modules.NonEsModuleProcessor
type NonEsModuleProcessor struct {
	idHash_ atomic.Int32
}

func (this_ *NonEsModuleProcessor) BindSelf_(s any) {
	_ = s
}
func (this_ *NonEsModuleProcessor) IsJavaNil_() bool { return this_ == nil }
func (*NonEsModuleProcessor) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.modules.NonEsModuleProcessor", "NonEsModuleProcessor"
}
func (this_ *NonEsModuleProcessor) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*NonEsModuleProcessor) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.modules.NonEsModuleProcessor", "com.google.javascript.jscomp.modules.NonEsModuleProcessor", func(v_ any) bool {
		return j2gort.InstanceOf[*NonEsModuleProcessor](v_)
	})
}
func (*NonEsModuleProcessor) Implements_jscomp_ModuleMapCreator_ModuleProcessor() {}

func (this_ *NonEsModuleProcessor) initFields_NonEsModuleProcessor() {
	_ = this_
}

// Java constructor: NonEsModuleProcessor()
func NewNonEsModuleProcessor() *NonEsModuleProcessor {
	this_ := &NonEsModuleProcessor{}
	InitNonEsModuleProcessor(this_)
	return this_
}

func InitNonEsModuleProcessor(this_ *NonEsModuleProcessor) {
	this_.initFields_NonEsModuleProcessor()
}

// Java method: NonEsModuleProcessor.process(com.google.javascript.jscomp.modules.ModuleMetadataMap.ModuleMetadata metadata, com.google.javascript.jscomp.deps.ModuleLoader.ModulePath path, com.google.javascript.rhino.Node script)
func (this_ *NonEsModuleProcessor) Process(metadata ModuleMetadataMap_ModuleMetadata, path *ModuleLoader_ModulePath, script pkg_rhino.Node) UnresolvedModule {
	_ = metadata
	_ = path
	_ = script
	return NewNonEsModuleProcessor_NonEsModule(metadata, path, script)
}

// Java class: com.google.javascript.jscomp.modules.NonEsModuleProcessor.NonEsModule
type NonEsModuleProcessor_NonEsModule struct {
	UnresolvedModuleBase
	metadata   ModuleMetadataMap_ModuleMetadata `gson:"metadata"`
	path       *ModuleLoader_ModulePath         `gson:"path"`
	scriptNode pkg_rhino.Node                   `gson:"scriptNode"`
}

func (this_ *NonEsModuleProcessor_NonEsModule) BindSelf_(s any) {
	_ = s
	this_.UnresolvedModuleBase.BindSelf_(s)
}
func (this_ *NonEsModuleProcessor_NonEsModule) IsJavaNil_() bool { return this_ == nil }
func (*NonEsModuleProcessor_NonEsModule) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.modules.NonEsModuleProcessor$NonEsModule", "NonEsModule"
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.modules.NonEsModuleProcessor.NonEsModule", "com.google.javascript.jscomp.modules.NonEsModuleProcessor$NonEsModule", func(v_ any) bool {
		return j2gort.InstanceOf[*NonEsModuleProcessor_NonEsModule](v_)
	})
}
func (*NonEsModuleProcessor_NonEsModule) Implements_jscomp_UnresolvedModule() {}

func (this_ *NonEsModuleProcessor_NonEsModule) initFields_NonEsModuleProcessor_NonEsModule() {
	_ = this_
}

// Java constructor: NonEsModule(com.google.javascript.jscomp.modules.ModuleMetadataMap.ModuleMetadata metadata, com.google.javascript.jscomp.deps.ModuleLoader.ModulePath path, com.google.javascript.rhino.Node scriptNode)
func NewNonEsModuleProcessor_NonEsModule(metadata ModuleMetadataMap_ModuleMetadata, path *ModuleLoader_ModulePath, scriptNode pkg_rhino.Node) *NonEsModuleProcessor_NonEsModule {
	this_ := &NonEsModuleProcessor_NonEsModule{}
	this_.BindSelf_(this_)
	InitNonEsModuleProcessor_NonEsModule(this_, metadata, path, scriptNode)
	return this_
}

func InitNonEsModuleProcessor_NonEsModule(this_ *NonEsModuleProcessor_NonEsModule, metadata ModuleMetadataMap_ModuleMetadata, path *ModuleLoader_ModulePath, scriptNode pkg_rhino.Node) {
	_ = metadata
	_ = path
	_ = scriptNode
	InitUnresolvedModule(&this_.UnresolvedModuleBase)
	this_.initFields_NonEsModuleProcessor_NonEsModule()
	this_.metadata = metadata
	this_.path = path
	this_.scriptNode = scriptNode
}

// Java method: NonEsModule.resolveExport(com.google.javascript.jscomp.modules.ModuleRequestResolver moduleRequestResolver, java.lang.String exportName)
func (this_ *NonEsModuleProcessor_NonEsModule) ResolveExport_2(moduleRequestResolver ModuleRequestResolver, exportName string) *ResolveExportResult {
	_ = moduleRequestResolver
	_ = exportName
	panic(j2gort.AsIface[javalang.IUnsupportedOperationException](j2gort.ConstructShimThrowable("java.lang.UnsupportedOperationException")))
}

// Java method: NonEsModule.resolveExport(com.google.javascript.jscomp.modules.ModuleRequestResolver moduleRequestResolver, java.lang.@org.jspecify.annotations.Nullable String moduleSpecifier, java.lang.String exportName, java.util.Set<com.google.javascript.jscomp.modules.ExportTrace> resolveSet, java.util.Set<com.google.javascript.jscomp.modules.UnresolvedModule> exportStarSet)
func (this_ *NonEsModuleProcessor_NonEsModule) ResolveExport_5(moduleRequestResolver ModuleRequestResolver, moduleSpecifier string, exportName string, resolveSet j2gort.Ext_java_util_Set, exportStarSet j2gort.Ext_java_util_Set) *ResolveExportResult {
	_ = moduleRequestResolver
	_ = moduleSpecifier
	_ = exportName
	_ = resolveSet
	_ = exportStarSet
	var namespace string = j2gort.NullString
	_ = namespace
	if (moduleSpecifier != j2gort.NullString) && GoogEsImports_IsGoogImportSpecifier(moduleSpecifier) {
		namespace = GoogEsImports_GetClosureIdFromGoogImportSpecifier(moduleSpecifier)
	} else if this_.metadata.IsGoogProvide() {
		namespace = moduleSpecifier
	}
	if this_.metadata.IsCommonJs() {
		if !(j2gort.CheckNotNullString(exportName) == "default") {
			return StaticGet_ResolveExportResult_NOT_FOUND()
		}
	}
	return ResolveExportResult_Of(Binding_From_2(Export_S_Builder().ExportName(exportName).ModuleMetadata(this_.metadata).ModulePath(this_.path).ClosureNamespace(namespace).Build(), this_.scriptNode))
}

// Java method: NonEsModule.reset()
func (this_ *NonEsModuleProcessor_NonEsModule) Reset() {
}

// Java method: NonEsModule.resolve(com.google.javascript.jscomp.modules.ModuleRequestResolver moduleRequestResolver, java.lang.@org.jspecify.annotations.Nullable String moduleSpecifier)
func (this_ *NonEsModuleProcessor_NonEsModule) Resolve_2(moduleRequestResolver ModuleRequestResolver, moduleSpecifier string) *Module {
	_ = moduleRequestResolver
	_ = moduleSpecifier
	var namespace string = j2gort.NullString
	_ = namespace
	if (moduleSpecifier != j2gort.NullString) && GoogEsImports_IsGoogImportSpecifier(moduleSpecifier) {
		namespace = GoogEsImports_GetClosureIdFromGoogImportSpecifier(moduleSpecifier)
	}
	return Module_S_Builder().Path(this_.path).Metadata(this_.metadata).Namespace(j2gort.Ext_com_google_common_collect_ImmutableMap_Of_0()).BoundNames(j2gort.Ext_com_google_common_collect_ImmutableMap_Of_0()).LocalNameToLocalExport(j2gort.Ext_com_google_common_collect_ImmutableMap_Of_0()).ClosureNamespace(namespace).Build()
}

// Java method: NonEsModule.metadata()
func (this_ *NonEsModuleProcessor_NonEsModule) Metadata() ModuleMetadataMap_ModuleMetadata {
	return this_.metadata
}

// Java method: NonEsModule.getExportedNames(com.google.javascript.jscomp.modules.ModuleRequestResolver moduleRequestResolver)
func (this_ *NonEsModuleProcessor_NonEsModule) GetExportedNames_1(moduleRequestResolver ModuleRequestResolver) j2gort.Ext_com_google_common_collect_ImmutableSet {
	_ = moduleRequestResolver
	panic(j2gort.AsIface[javalang.IUnsupportedOperationException](j2gort.ConstructShimThrowable("java.lang.UnsupportedOperationException")))
}

// Java method: NonEsModule.getExportedNames(com.google.javascript.jscomp.modules.ModuleRequestResolver moduleRequestResolver, java.util.Set<com.google.javascript.jscomp.modules.UnresolvedModule> visited)
func (this_ *NonEsModuleProcessor_NonEsModule) GetExportedNames_2(moduleRequestResolver ModuleRequestResolver, visited j2gort.Ext_java_util_Set) j2gort.Ext_com_google_common_collect_ImmutableSet {
	_ = moduleRequestResolver
	_ = visited
	panic(j2gort.AsIface[javalang.IUnsupportedOperationException](j2gort.ConstructShimThrowable("java.lang.UnsupportedOperationException")))
}
