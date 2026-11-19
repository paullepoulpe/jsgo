/*
 * Copyright 2019 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/modules/UnresolvedModule.java (generated mechanically by j2go).

package jscomp

import (
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java polymorphic class interface: com.google.javascript.jscomp.modules.UnresolvedModule
type UnresolvedModule interface {
	Implements_jscomp_UnresolvedModule()
	IsJavaNil_() bool
	GetUnresolvedModuleBase() *UnresolvedModuleBase
	Reset()
	Resolve_1(ModuleRequestResolver) *Module
	Resolve_2(ModuleRequestResolver, string) *Module
	Metadata() ModuleMetadataMap_ModuleMetadata
	GetExportedNames_1(ModuleRequestResolver) j2gort.Ext_com_google_common_collect_ImmutableSet
	GetExportedNames_2(ModuleRequestResolver, j2gort.Ext_java_util_Set) j2gort.Ext_com_google_common_collect_ImmutableSet
	ResolveExport_2(ModuleRequestResolver, string) *ResolveExportResult
	ResolveExport_5(ModuleRequestResolver, string, string, j2gort.Ext_java_util_Set, j2gort.Ext_java_util_Set) *ResolveExportResult
	Equals(any) bool
	HashCode() int32
}

func ToIface_UnresolvedModule[T interface {
	UnresolvedModule
	comparable
}](v T) UnresolvedModule {
	var z T
	if v == z {
		return nil
	}
	return v
}

// Java class: com.google.javascript.jscomp.modules.UnresolvedModule
type UnresolvedModuleBase struct {
	self_ UnresolvedModule
}

func (this_ *UnresolvedModuleBase) BindSelf_(s any) {
	_ = s
	this_.self_ = s.(UnresolvedModule)
}
func (this_ *UnresolvedModuleBase) IsJavaNil_() bool { return this_ == nil }
func (*UnresolvedModuleBase) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.modules.UnresolvedModule", "UnresolvedModule"
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.modules.UnresolvedModule", "com.google.javascript.jscomp.modules.UnresolvedModule", func(v_ any) bool {
		return j2gort.InstanceOf[UnresolvedModule](v_)
	})
}
func (b *UnresolvedModuleBase) GetUnresolvedModuleBase() *UnresolvedModuleBase { return b }
func (*UnresolvedModuleBase) Implements_jscomp_UnresolvedModule()              {}

func (this_ *UnresolvedModuleBase) initFields_UnresolvedModule() {
	_ = this_
}

func InitUnresolvedModule(this_ *UnresolvedModuleBase) {
	this_.initFields_UnresolvedModule()
}

// Java method: UnresolvedModule.reset()
func (this_ *UnresolvedModuleBase) Reset() {
	panic("abstract or native method invoked")
}

// Java method: UnresolvedModule.resolve(com.google.javascript.jscomp.modules.ModuleRequestResolver moduleRequestResolver)
func (this_ *UnresolvedModuleBase) Resolve_1(moduleRequestResolver ModuleRequestResolver) *Module {
	_ = moduleRequestResolver
	return this_.self_.Resolve_2(moduleRequestResolver, j2gort.NullString)
}

// Java method: UnresolvedModule.resolve(com.google.javascript.jscomp.modules.ModuleRequestResolver moduleRequestResolver, java.lang.@org.jspecify.annotations.Nullable String moduleSpecifier)
func (this_ *UnresolvedModuleBase) Resolve_2(moduleRequestResolver ModuleRequestResolver, moduleSpecifier string) *Module {
	_ = moduleRequestResolver
	_ = moduleSpecifier
	panic("abstract or native method invoked")
}

// Java method: UnresolvedModule.metadata()
func (this_ *UnresolvedModuleBase) Metadata() ModuleMetadataMap_ModuleMetadata {
	panic("abstract or native method invoked")
}

// Java method: UnresolvedModule.getExportedNames(com.google.javascript.jscomp.modules.ModuleRequestResolver moduleRequestResolver)
func (this_ *UnresolvedModuleBase) GetExportedNames_1(moduleRequestResolver ModuleRequestResolver) j2gort.Ext_com_google_common_collect_ImmutableSet {
	_ = moduleRequestResolver
	panic("abstract or native method invoked")
}

// Java method: UnresolvedModule.getExportedNames(com.google.javascript.jscomp.modules.ModuleRequestResolver moduleRequestResolver, java.util.Set<com.google.javascript.jscomp.modules.UnresolvedModule> visited)
func (this_ *UnresolvedModuleBase) GetExportedNames_2(moduleRequestResolver ModuleRequestResolver, visited j2gort.Ext_java_util_Set) j2gort.Ext_com_google_common_collect_ImmutableSet {
	_ = moduleRequestResolver
	_ = visited
	panic("abstract or native method invoked")
}

// Java method: UnresolvedModule.resolveExport(com.google.javascript.jscomp.modules.ModuleRequestResolver moduleRequestResolver, java.lang.String exportName)
func (this_ *UnresolvedModuleBase) ResolveExport_2(moduleRequestResolver ModuleRequestResolver, exportName string) *ResolveExportResult {
	_ = moduleRequestResolver
	_ = exportName
	return this_.self_.ResolveExport_5(moduleRequestResolver, j2gort.NullString, exportName, j2gort.NewExt_java_util_LinkedHashSet(), j2gort.NewExt_java_util_LinkedHashSet())
}

// Java method: UnresolvedModule.resolveExport(com.google.javascript.jscomp.modules.ModuleRequestResolver moduleRequestResolver, java.lang.@org.jspecify.annotations.Nullable String moduleSpecifier, java.lang.String exportName, java.util.Set<com.google.javascript.jscomp.modules.ExportTrace> resolveSet, java.util.Set<com.google.javascript.jscomp.modules.UnresolvedModule> exportStarSet)
func (this_ *UnresolvedModuleBase) ResolveExport_5(moduleRequestResolver ModuleRequestResolver, moduleSpecifier string, exportName string, resolveSet j2gort.Ext_java_util_Set, exportStarSet j2gort.Ext_java_util_Set) *ResolveExportResult {
	_ = moduleRequestResolver
	_ = moduleSpecifier
	_ = exportName
	_ = resolveSet
	_ = exportStarSet
	panic("abstract or native method invoked")
}

// Java method: UnresolvedModule.equals(java.lang.Object other)
func (this_ *UnresolvedModuleBase) Equals(arg_1 any) bool {
	other := arg_1
	_ = other
	return j2gort.RefEquals(this_.self_, other)
}

// Java method: UnresolvedModule.hashCode()
func (this_ *UnresolvedModuleBase) HashCode() int32 {
	return javalang.IdentityHashCode(this_.self_)
}
