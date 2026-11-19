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
// Port of com/google/javascript/jscomp/modules/Module.java (generated mechanically by j2go).

package jscomp

import (
	"fmt"
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java record: com.google.javascript.jscomp.modules.Module
type Module struct {
	metadata               ModuleMetadataMap_ModuleMetadata                  `gson:"metadata"`
	path                   *ModuleLoader_ModulePath                          `gson:"path"`
	namespace              j2gort.Ext_com_google_common_collect_ImmutableMap `gson:"namespace"`
	boundNames             j2gort.Ext_com_google_common_collect_ImmutableMap `gson:"boundNames"`
	localNameToLocalExport j2gort.Ext_com_google_common_collect_ImmutableMap `gson:"localNameToLocalExport"`
	closureNamespace       string                                            `gson:"closureNamespace"`
}

func (this_ *Module) BindSelf_(s any) {
	_ = s
}
func (this_ *Module) IsJavaNil_() bool { return this_ == nil }
func (*Module) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.modules.Module", "Module"
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.modules.Module", "com.google.javascript.jscomp.modules.Module", func(v_ any) bool {
		return j2gort.InstanceOf[*Module](v_)
	})
}
func (*Module) Implements_j2gort_Ext_java_lang_Record() {}

func (this_ *Module) initFields_Module() {
	_ = this_
}

// Java constructor: Module(com.google.javascript.jscomp.modules.ModuleMetadataMap.ModuleMetadata metadata, com.google.javascript.jscomp.deps.ModuleLoader.@org.jspecify.annotations.Nullable ModulePath path, com.google.common.collect.ImmutableMap<java.lang.String,com.google.javascript.jscomp.modules.Binding> namespace, com.google.common.collect.ImmutableMap<java.lang.String,com.google.javascript.jscomp.modules.Binding> boundNames, com.google.common.collect.ImmutableMap<java.lang.String,com.google.javascript.jscomp.modules.Export> localNameToLocalExport, java.lang.@org.jspecify.annotations.Nullable String closureNamespace)
func NewModule(metadata ModuleMetadataMap_ModuleMetadata, path *ModuleLoader_ModulePath, namespace j2gort.Ext_com_google_common_collect_ImmutableMap, boundNames j2gort.Ext_com_google_common_collect_ImmutableMap, localNameToLocalExport j2gort.Ext_com_google_common_collect_ImmutableMap, closureNamespace string) *Module {
	this_ := &Module{}
	InitModule(this_, metadata, path, namespace, boundNames, localNameToLocalExport, closureNamespace)
	return this_
}

func InitModule(this_ *Module, metadata ModuleMetadataMap_ModuleMetadata, path *ModuleLoader_ModulePath, namespace j2gort.Ext_com_google_common_collect_ImmutableMap, boundNames j2gort.Ext_com_google_common_collect_ImmutableMap, localNameToLocalExport j2gort.Ext_com_google_common_collect_ImmutableMap, closureNamespace string) {
	_ = metadata
	_ = path
	_ = namespace
	_ = boundNames
	_ = localNameToLocalExport
	_ = closureNamespace
	this_.initFields_Module()
	j2gort.CheckNotNullIface(metadata)
	j2gort.CheckNotNull(namespace)
	j2gort.CheckNotNull(boundNames)
	j2gort.CheckNotNull(localNameToLocalExport)
	this_.metadata = metadata
	this_.path = path
	this_.namespace = namespace
	this_.boundNames = boundNames
	this_.localNameToLocalExport = localNameToLocalExport
	this_.closureNamespace = closureNamespace
}

// Java static method: Module.builder()
func Module_S_Builder() Module_Builder {
	return NewAutoBuilder_Module_Builder()
}

// Java method: Module.toBuilder()
func (this_ *Module) ToBuilder() Module_Builder {
	return NewAutoBuilder_Module_Builder_1(this_)
}

func (this_ *Module) Metadata() ModuleMetadataMap_ModuleMetadata { return this_.metadata }
func (this_ *Module) Path() *ModuleLoader_ModulePath             { return this_.path }
func (this_ *Module) Namespace() j2gort.Ext_com_google_common_collect_ImmutableMap {
	return this_.namespace
}
func (this_ *Module) BoundNames() j2gort.Ext_com_google_common_collect_ImmutableMap {
	return this_.boundNames
}
func (this_ *Module) LocalNameToLocalExport() j2gort.Ext_com_google_common_collect_ImmutableMap {
	return this_.localNameToLocalExport
}
func (this_ *Module) ClosureNamespace() string { return this_.closureNamespace }
func (this_ *Module) Equals(other_ any) bool {
	o_, ok_ := other_.(*Module)
	if !ok_ || o_ == nil {
		return false
	}
	if this_ == o_ {
		return true
	}
	if !javalang.ObjectsEquals(this_.metadata, o_.metadata) {
		return false
	}
	if this_.path != o_.path {
		return false
	}
	if !javalang.ObjectsEquals(this_.namespace, o_.namespace) {
		return false
	}
	if !javalang.ObjectsEquals(this_.boundNames, o_.boundNames) {
		return false
	}
	if !javalang.ObjectsEquals(this_.localNameToLocalExport, o_.localNameToLocalExport) {
		return false
	}
	if this_.closureNamespace != o_.closureNamespace {
		return false
	}
	return true
}
func (this_ *Module) HashCode() int32 {
	var h_ int32 = 0
	{
		hi_1 := this_.metadata
		h_ = 31 * h_
		if hi_1 != nil && !hi_1.IsJavaNil_() {
			h_ += hi_1.HashCode()
		}
	}
	{
		hp_2 := this_.path
		h_ = 31 * h_
		if hp_2 != nil {
			h_ += j2gort.IdentityHash(hp_2)
		}
	}
	h_ = 31*h_ + j2gort.ObjectHash(this_.namespace)
	h_ = 31*h_ + j2gort.ObjectHash(this_.boundNames)
	h_ = 31*h_ + j2gort.ObjectHash(this_.localNameToLocalExport)
	{
		hs_3 := this_.closureNamespace
		h_ = 31 * h_
		if hs_3 != j2gort.NullString {
			h_ += javalang.StringHashCode(hs_3)
		}
	}
	return h_
}
func (this_ *Module) String() string {
	_ = fmt.Sprint
	return "Module[" + "metadata=" + j2gort.StrVal(this_.metadata) + ", " + "path=" + j2gort.StrVal(j2gort.PtrToAny(this_.path)) + ", " + "namespace=" + j2gort.StrVal(this_.namespace) + ", " + "boundNames=" + j2gort.StrVal(this_.boundNames) + ", " + "localNameToLocalExport=" + j2gort.StrVal(this_.localNameToLocalExport) + ", " + "closureNamespace=" + j2gort.StrVal(j2gort.BoxString(this_.closureNamespace)) + "]"
}

// Java polymorphic class interface: com.google.javascript.jscomp.modules.Module.Builder
type Module_Builder interface {
	Implements_jscomp_Module_Builder()
	IsJavaNil_() bool
	GetModule_BuilderBase() *Module_BuilderBase
	Metadata(ModuleMetadataMap_ModuleMetadata) Module_Builder
	Path(*ModuleLoader_ModulePath) Module_Builder
	Namespace(j2gort.Ext_com_google_common_collect_ImmutableMap) Module_Builder
	BoundNames(j2gort.Ext_com_google_common_collect_ImmutableMap) Module_Builder
	LocalNameToLocalExport(j2gort.Ext_com_google_common_collect_ImmutableMap) Module_Builder
	ClosureNamespace(string) Module_Builder
	Build() *Module
}

func ToIface_Module_Builder[T interface {
	Module_Builder
	comparable
}](v T) Module_Builder {
	var z T
	if v == z {
		return nil
	}
	return v
}

// Java class: com.google.javascript.jscomp.modules.Module.Builder
type Module_BuilderBase struct {
	idHash_ atomic.Int32
}

func (this_ *Module_BuilderBase) BindSelf_(s any) {
	_ = s
}
func (this_ *Module_BuilderBase) IsJavaNil_() bool { return this_ == nil }
func (*Module_BuilderBase) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.modules.Module$Builder", "Builder"
}
func (this_ *Module_BuilderBase) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.modules.Module.Builder", "com.google.javascript.jscomp.modules.Module$Builder", func(v_ any) bool {
		return j2gort.InstanceOf[Module_Builder](v_)
	})
}
func (b *Module_BuilderBase) GetModule_BuilderBase() *Module_BuilderBase { return b }
func (*Module_BuilderBase) Implements_jscomp_Module_Builder()            {}

func (this_ *Module_BuilderBase) initFields_Module_Builder() {
	_ = this_
}

func InitModule_Builder(this_ *Module_BuilderBase) {
	this_.initFields_Module_Builder()
}

// Java method: Builder.metadata(com.google.javascript.jscomp.modules.ModuleMetadataMap.ModuleMetadata value)
func (this_ *Module_BuilderBase) Metadata(value ModuleMetadataMap_ModuleMetadata) Module_Builder {
	_ = value
	panic("abstract or native method invoked")
}

// Java method: Builder.path(com.google.javascript.jscomp.deps.ModuleLoader.@org.jspecify.annotations.Nullable ModulePath value)
func (this_ *Module_BuilderBase) Path(value *ModuleLoader_ModulePath) Module_Builder {
	_ = value
	panic("abstract or native method invoked")
}

// Java method: Builder.namespace(com.google.common.collect.ImmutableMap<java.lang.String,com.google.javascript.jscomp.modules.Binding> value)
func (this_ *Module_BuilderBase) Namespace(value j2gort.Ext_com_google_common_collect_ImmutableMap) Module_Builder {
	_ = value
	panic("abstract or native method invoked")
}

// Java method: Builder.boundNames(com.google.common.collect.ImmutableMap<java.lang.String,com.google.javascript.jscomp.modules.Binding> value)
func (this_ *Module_BuilderBase) BoundNames(value j2gort.Ext_com_google_common_collect_ImmutableMap) Module_Builder {
	_ = value
	panic("abstract or native method invoked")
}

// Java method: Builder.localNameToLocalExport(com.google.common.collect.ImmutableMap<java.lang.String,com.google.javascript.jscomp.modules.Export> value)
func (this_ *Module_BuilderBase) LocalNameToLocalExport(value j2gort.Ext_com_google_common_collect_ImmutableMap) Module_Builder {
	_ = value
	panic("abstract or native method invoked")
}

// Java method: Builder.closureNamespace(java.lang.@org.jspecify.annotations.Nullable String value)
func (this_ *Module_BuilderBase) ClosureNamespace(value string) Module_Builder {
	_ = value
	panic("abstract or native method invoked")
}

// Java method: Builder.build()
func (this_ *Module_BuilderBase) Build() *Module {
	panic("abstract or native method invoked")
}
