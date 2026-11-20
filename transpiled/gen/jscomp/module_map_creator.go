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
// Port of com/google/javascript/jscomp/modules/ModuleMapCreator.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
)

// Java class: com.google.javascript.jscomp.modules.ModuleMapCreator
type ModuleMapCreator struct {
	compiler                            AbstractCompiler                            `gson:"compiler"`
	esModuleProcessor                   *EsModuleProcessor                          `gson:"esModuleProcessor"`
	closureModuleProcessor              *ClosureModuleProcessor                     `gson:"closureModuleProcessor"`
	nonEsModuleProcessor                *NonEsModuleProcessor                       `gson:"nonEsModuleProcessor"`
	unresolvedModules                   j2gort.OrderedMap[string, UnresolvedModule] `gson:"unresolvedModules"`
	unresolvedModulesByClosureNamespace j2gort.OrderedMap[string, UnresolvedModule] `gson:"unresolvedModulesByClosureNamespace"`
	moduleMetadataMap                   *ModuleMetadataMap                          `gson:"moduleMetadataMap"`
	idHash_                             atomic.Int32
}

func (this_ *ModuleMapCreator) BindSelf_(s any) {
	_ = s
}
func (this_ *ModuleMapCreator) IsJavaNil_() bool { return this_ == nil }
func (*ModuleMapCreator) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.modules.ModuleMapCreator", "ModuleMapCreator"
}
func (this_ *ModuleMapCreator) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*ModuleMapCreator) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.modules.ModuleMapCreator", "com.google.javascript.jscomp.modules.ModuleMapCreator", func(v_ any) bool {
		return j2gort.InstanceOf[*ModuleMapCreator](v_)
	})
}
func (*ModuleMapCreator) Implements_jscomp_CompilerPass() {}

var static_ModuleMapCreator_DOES_NOT_HAVE_EXPORT *DiagnosticType

func StaticGet_ModuleMapCreator_DOES_NOT_HAVE_EXPORT() *DiagnosticType {
	clinit_ModuleMapCreator()
	return static_ModuleMapCreator_DOES_NOT_HAVE_EXPORT
}
func StaticSet_ModuleMapCreator_DOES_NOT_HAVE_EXPORT(v *DiagnosticType) {
	clinit_ModuleMapCreator()
	static_ModuleMapCreator_DOES_NOT_HAVE_EXPORT = v
}

var static_ModuleMapCreator_DOES_NOT_HAVE_EXPORT_WITH_DETAILS *DiagnosticType

func StaticGet_ModuleMapCreator_DOES_NOT_HAVE_EXPORT_WITH_DETAILS() *DiagnosticType {
	clinit_ModuleMapCreator()
	return static_ModuleMapCreator_DOES_NOT_HAVE_EXPORT_WITH_DETAILS
}
func StaticSet_ModuleMapCreator_DOES_NOT_HAVE_EXPORT_WITH_DETAILS(v *DiagnosticType) {
	clinit_ModuleMapCreator()
	static_ModuleMapCreator_DOES_NOT_HAVE_EXPORT_WITH_DETAILS = v
}

var clinitDone_ModuleMapCreator uint32
var clinitState_ModuleMapCreator j2gort.ClinitState

func clinit_ModuleMapCreator() {
	if atomic.LoadUint32(&clinitDone_ModuleMapCreator) != 0 {
		return
	}
	clinitSlow_ModuleMapCreator()
}

func clinitSlow_ModuleMapCreator() {
	if !j2gort.EnterClinit(&clinitState_ModuleMapCreator, &clinitDone_ModuleMapCreator) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_ModuleMapCreator, &clinitDone_ModuleMapCreator)
	static_ModuleMapCreator_DOES_NOT_HAVE_EXPORT = DiagnosticType_Error("JSC_DOES_NOT_HAVE_EXPORT", "Requested module does not have an export \"{0}\".")
	static_ModuleMapCreator_DOES_NOT_HAVE_EXPORT_WITH_DETAILS = DiagnosticType_Error("JSC_DOES_NOT_HAVE_EXPORT_WITH_DETAILS", "Requested module does not have an export \"{0}\".{1}")
}

func (this_ *ModuleMapCreator) initFields_ModuleMapCreator() {
	_ = this_
}

// Java constructor: ModuleMapCreator(com.google.javascript.jscomp.AbstractCompiler compiler, com.google.javascript.jscomp.modules.ModuleMetadataMap moduleMetadataMap)
func NewModuleMapCreator(compiler AbstractCompiler, moduleMetadataMap *ModuleMetadataMap) *ModuleMapCreator {
	clinit_ModuleMapCreator()
	this_ := &ModuleMapCreator{}
	InitModuleMapCreator(this_, compiler, moduleMetadataMap)
	return this_
}

func InitModuleMapCreator(this_ *ModuleMapCreator, compiler AbstractCompiler, moduleMetadataMap *ModuleMetadataMap) {
	_ = compiler
	_ = moduleMetadataMap
	this_.initFields_ModuleMapCreator()
	this_.compiler = compiler
	this_.moduleMetadataMap = moduleMetadataMap
	this_.esModuleProcessor = NewEsModuleProcessor(compiler)
	this_.closureModuleProcessor = NewClosureModuleProcessor(compiler)
	this_.nonEsModuleProcessor = NewNonEsModuleProcessor()
	this_.unresolvedModules = j2gort.NewOrderedMap[string, UnresolvedModule]()
	this_.unresolvedModulesByClosureNamespace = j2gort.NewOrderedMap[string, UnresolvedModule]()
}

// Java method: ModuleMapCreator.create()
func (this_ *ModuleMapCreator) priv_ModuleMapCreator_create_6782() *ModuleMap {
	this_.unresolvedModules.Clear()
	this_.unresolvedModulesByClosureNamespace.Clear()
	{
		cur_1 := j2gort.NewIterCursor(j2gort.CheckNotNull(this_.moduleMetadataMap.GetAllModuleMetadata()))
		for cur_1.HasNext() {
			rawElem_2 := cur_1.Next()
			moduleMetadata := j2gort.AsIface[ModuleMetadataMap_ModuleMetadata](rawElem_2)
			_ = moduleMetadata
			this_.priv_ModuleMapCreator_process_6784(moduleMetadata)
		}
	}
	return this_.priv_ModuleMapCreator_resolve_6783()
}

// Java method: ModuleMapCreator.resolve()
func (this_ *ModuleMapCreator) priv_ModuleMapCreator_resolve_6783() *ModuleMap {
	var requestResolver ModuleRequestResolver = NewModuleMapCreator_ModuleRequestResolverImpl(this_)
	_ = requestResolver
	var resolvedModules j2gort.OrderedMap[string, *Module] = j2gort.NewOrderedMap[string, *Module]()
	_ = resolvedModules
	var resolvedClosureModules j2gort.OrderedMap[string, *Module] = j2gort.NewOrderedMap[string, *Module]()
	_ = resolvedClosureModules
	for doFirst_3 := true; doFirst_3 || (!(j2gort.AsIface[j2gort.Ext_java_util_Set](resolvedModules.LiveKeySet()).ContainsAll(j2gort.AsIface[j2gort.Ext_java_util_Collection](this_.unresolvedModules.LiveKeySet())))); doFirst_3 = false {
		var toResolve j2gort.Ext_com_google_common_collect_ImmutableSet = j2gort.Ext_com_google_common_collect_Sets_Difference(j2gort.AsIface[j2gort.Ext_java_util_Set](this_.unresolvedModules.LiveKeySet()), j2gort.AsIface[j2gort.Ext_java_util_Set](resolvedModules.LiveKeySet())).ImmutableCopy()
		_ = toResolve
		{
			cur_4 := j2gort.NewIterCursor(j2gort.CheckNotNull(toResolve))
			for cur_4.HasNext() {
				key := cur_4.NextStr()
				_ = key
				var resolved *Module = this_.unresolvedModules.Get(key).GetUnresolvedModuleBase().Resolve_1(requestResolver)
				_ = resolved
				resolvedModules.Put(key, resolved)
				{
					cur_5 := j2gort.NewIterCursor(j2gort.CheckNotNull(resolved.Metadata().GoogNamespaces()))
					for cur_5.HasNext() {
						namespace := cur_5.NextStr()
						_ = namespace
						resolvedClosureModules.Put(namespace, resolved)
					}
				}
			}
		}
	}
	for doFirst_6 := true; doFirst_6 || (!(j2gort.AsIface[j2gort.Ext_java_util_Set](resolvedClosureModules.LiveKeySet()).ContainsAll(j2gort.AsIface[j2gort.Ext_java_util_Collection](this_.unresolvedModulesByClosureNamespace.LiveKeySet())))); doFirst_6 = false {
		var toResolve_2 j2gort.Ext_com_google_common_collect_ImmutableSet = j2gort.Ext_com_google_common_collect_Sets_Difference(j2gort.AsIface[j2gort.Ext_java_util_Set](this_.unresolvedModulesByClosureNamespace.LiveKeySet()), j2gort.AsIface[j2gort.Ext_java_util_Set](resolvedClosureModules.LiveKeySet())).ImmutableCopy()
		_ = toResolve_2
		{
			cur_7 := j2gort.NewIterCursor(j2gort.CheckNotNull(toResolve_2))
			for cur_7.HasNext() {
				namespace_2 := cur_7.NextStr()
				_ = namespace_2
				resolvedClosureModules.Put(namespace_2, this_.unresolvedModulesByClosureNamespace.Get(namespace_2).GetUnresolvedModuleBase().Resolve_1(requestResolver))
			}
		}
	}
	this_.unresolvedModules.Clear()
	this_.unresolvedModulesByClosureNamespace.Clear()
	return NewModuleMap(resolvedModules.ToImmutableMap(), resolvedClosureModules.ToImmutableMap())
}

// Java method: ModuleMapCreator.process(com.google.javascript.jscomp.modules.ModuleMetadataMap.ModuleMetadata moduleMetadata)
func (this_ *ModuleMapCreator) priv_ModuleMapCreator_process_6784(moduleMetadata ModuleMetadataMap_ModuleMetadata) {
	_ = moduleMetadata
	var processor ModuleMapCreator_ModuleProcessor = nil
	_ = processor
	switch (moduleMetadata.ModuleType()).Ordinal() {
	case int32(0) /* ModuleType.ES6_MODULE */ :
		processor = ToIface_ModuleMapCreator_ModuleProcessor(this_.esModuleProcessor)
	case int32(1) /* ModuleType.GOOG_MODULE */, int32(2) /* ModuleType.LEGACY_GOOG_MODULE */ :
		if moduleMetadata.GoogNamespaces().Size() == int32(1) {
			processor = ToIface_ModuleMapCreator_ModuleProcessor(this_.closureModuleProcessor)
		} else {
			processor = ToIface_ModuleMapCreator_ModuleProcessor(this_.nonEsModuleProcessor)
		}
	default:
		processor = ToIface_ModuleMapCreator_ModuleProcessor(this_.nonEsModuleProcessor)
	}
	var module UnresolvedModule = processor.Process(moduleMetadata, moduleMetadata.Path(), moduleMetadata.RootNode())
	_ = module
	if moduleMetadata.Path() != nil {
		this_.unresolvedModules.Put(moduleMetadata.Path().ToModuleName(), module)
	}
	{
		cur_8 := j2gort.NewIterCursor(j2gort.CheckNotNull(moduleMetadata.GoogNamespaces()))
		for cur_8.HasNext() {
			namespace := cur_8.NextStr()
			_ = namespace
			this_.unresolvedModulesByClosureNamespace.Put(namespace, module)
		}
	}
}

// Java method: ModuleMapCreator.process(com.google.javascript.rhino.Node externs, com.google.javascript.rhino.Node root)
func (this_ *ModuleMapCreator) Process_2(externs pkg_rhino.Node, root pkg_rhino.Node) {
	_ = externs
	_ = root
	this_.compiler.SetModuleMap(this_.priv_ModuleMapCreator_create_6782())
}

// Java class: com.google.javascript.jscomp.modules.ModuleMapCreator.ModuleRequestResolverImpl
type ModuleMapCreator_ModuleRequestResolverImpl struct {
	outer_  *ModuleMapCreator
	idHash_ atomic.Int32
}

func (this_ *ModuleMapCreator_ModuleRequestResolverImpl) BindSelf_(s any) {
	_ = s
}
func (this_ *ModuleMapCreator_ModuleRequestResolverImpl) IsJavaNil_() bool { return this_ == nil }
func (*ModuleMapCreator_ModuleRequestResolverImpl) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.modules.ModuleMapCreator$ModuleRequestResolverImpl", "ModuleRequestResolverImpl"
}
func (this_ *ModuleMapCreator_ModuleRequestResolverImpl) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*ModuleMapCreator_ModuleRequestResolverImpl) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.modules.ModuleMapCreator.ModuleRequestResolverImpl", "com.google.javascript.jscomp.modules.ModuleMapCreator$ModuleRequestResolverImpl", func(v_ any) bool {
		return j2gort.InstanceOf[*ModuleMapCreator_ModuleRequestResolverImpl](v_)
	})
}
func (this_ *ModuleMapCreator_ModuleRequestResolverImpl) SetOuter_(o *ModuleMapCreator) {
	this_.outer_ = o
}
func (this_ *ModuleMapCreator_ModuleRequestResolverImpl) GetOuter_() *ModuleMapCreator {
	return this_.outer_
}
func (*ModuleMapCreator_ModuleRequestResolverImpl) Implements_jscomp_ModuleRequestResolver() {}

func (this_ *ModuleMapCreator_ModuleRequestResolverImpl) initFields_ModuleMapCreator_ModuleRequestResolverImpl() {
	_ = this_
}

// Java constructor: ModuleRequestResolverImpl()
func NewModuleMapCreator_ModuleRequestResolverImpl(outer_ *ModuleMapCreator) *ModuleMapCreator_ModuleRequestResolverImpl {
	this_ := &ModuleMapCreator_ModuleRequestResolverImpl{}
	this_.outer_ = outer_
	InitModuleMapCreator_ModuleRequestResolverImpl(this_)
	return this_
}

func InitModuleMapCreator_ModuleRequestResolverImpl(this_ *ModuleMapCreator_ModuleRequestResolverImpl) {
	this_.initFields_ModuleMapCreator_ModuleRequestResolverImpl()
}

// Java method: ModuleRequestResolverImpl.getFallbackForMissingNonClosureModule(com.google.javascript.jscomp.deps.ModuleLoader.ModulePath path)
func (this_ *ModuleMapCreator_ModuleRequestResolverImpl) priv_ModuleRequestResolverImpl_getFallbackForMissingNonClosureModule_6785(path *ModuleLoader_ModulePath) UnresolvedModule {
	_ = path
	var metadata ModuleMetadataMap_ModuleMetadata = ModuleMetadataMap_ModuleMetadata_S_Builder().RootNode(pkg_rhino.Node(0)).Path(path).ModuleType_1(StaticGet_ModuleMetadataMap_ModuleType_ES6_MODULE()).IsTestOnly(false).UsesClosure(false).Build()
	_ = metadata
	return NewModuleMapCreator__Anon_123(this_, path, metadata)
}

// Java method: ModuleRequestResolverImpl.getFallbackForMissingClosureModule(java.lang.String namespace)
func (this_ *ModuleMapCreator_ModuleRequestResolverImpl) priv_ModuleRequestResolverImpl_getFallbackForMissingClosureModule_6786(namespace string) UnresolvedModule {
	_ = namespace
	return this_.outer_.nonEsModuleProcessor.Process(ModuleMetadataMap_ModuleMetadata_S_Builder().GetModuleMetadataMap_ModuleMetadata_BuilderBase().AddGoogNamespace(namespace).IsTestOnly(false).ModuleType_1(StaticGet_ModuleMetadataMap_ModuleType_GOOG_PROVIDE()).Path(nil).RootNode(pkg_rhino.Node(0)).UsesClosure(true).Build(), nil, pkg_rhino.Node(0))
}

// Java method: ModuleRequestResolverImpl.resolve(com.google.javascript.jscomp.modules.Import i)
func (this_ *ModuleMapCreator_ModuleRequestResolverImpl) Resolve_1_Import(i *Import) UnresolvedModule {
	_ = i
	if i.ModulePath() == nil {
		return this_.priv_ModuleRequestResolverImpl_resolveForClosure_6789(i.ModuleRequest())
	}
	return this_.priv_ModuleRequestResolverImpl_resolve_6790(i.ModuleRequest(), i.ModulePath(), i.ImportNode())
}

// Java method: ModuleRequestResolverImpl.resolve(com.google.javascript.jscomp.modules.Export e)
func (this_ *ModuleMapCreator_ModuleRequestResolverImpl) Resolve_1_Export(e *Export) UnresolvedModule {
	_ = e
	return this_.priv_ModuleRequestResolverImpl_resolve_6790(e.ModuleRequest(), e.ModulePath(), e.ExportNode())
}

// Java method: ModuleRequestResolverImpl.resolveForClosure(java.lang.String namespace)
func (this_ *ModuleMapCreator_ModuleRequestResolverImpl) priv_ModuleRequestResolverImpl_resolveForClosure_6789(namespace string) UnresolvedModule {
	_ = namespace
	{

		bound_9 := this_
		return j2gort.OrderedMapComputeIfAbsent(&this_.outer_.unresolvedModulesByClosureNamespace, namespace, nil, func(mrp0_ string) UnresolvedModule {
			return bound_9.priv_ModuleRequestResolverImpl_getFallbackForMissingClosureModule_6786(mrp0_)
		})
	}

}

// Java method: ModuleRequestResolverImpl.resolve(java.lang.String moduleRequest, com.google.javascript.jscomp.deps.ModuleLoader.ModulePath modulePath, com.google.javascript.rhino.Node forLineInfo)
func (this_ *ModuleMapCreator_ModuleRequestResolverImpl) priv_ModuleRequestResolverImpl_resolve_6790(moduleRequest string, modulePath *ModuleLoader_ModulePath, forLineInfo pkg_rhino.Node) UnresolvedModule {
	_ = moduleRequest
	_ = modulePath
	_ = forLineInfo
	if GoogEsImports_IsGoogImportSpecifier(moduleRequest) {
		var namespace string = GoogEsImports_GetClosureIdFromGoogImportSpecifier(moduleRequest)
		_ = namespace
		return this_.priv_ModuleRequestResolverImpl_resolveForClosure_6789(namespace)
	}
	var requestedPath *ModuleLoader_ModulePath = modulePath.ResolveJsModule(moduleRequest, forLineInfo.GetSourceFileName(), forLineInfo.GetLineno(), forLineInfo.GetCharno())
	_ = requestedPath
	if requestedPath == nil {
		requestedPath = modulePath.ResolveModuleAsPath(moduleRequest)
		if !(this_.outer_.unresolvedModules.ContainsKey(requestedPath.ToModuleName())) {
			var module UnresolvedModule = this_.priv_ModuleRequestResolverImpl_getFallbackForMissingNonClosureModule_6785(requestedPath)
			_ = module
			this_.outer_.unresolvedModules.Put(requestedPath.ToModuleName(), module)
			return module
		}
	}
	return this_.outer_.unresolvedModules.Get(requestedPath.ToModuleName())
}

// Java interface: com.google.javascript.jscomp.modules.ModuleMapCreator.ModuleProcessor
type ModuleMapCreator_ModuleProcessor interface {
	Implements_jscomp_ModuleMapCreator_ModuleProcessor()
	IsJavaNil_() bool
	Process(ModuleMetadataMap_ModuleMetadata, *ModuleLoader_ModulePath, pkg_rhino.Node) UnresolvedModule
}

func ToIface_ModuleMapCreator_ModuleProcessor[T interface {
	ModuleMapCreator_ModuleProcessor
	comparable
}](v T) ModuleMapCreator_ModuleProcessor {
	var z T
	if v == z {
		return nil
	}
	return v
}

func Upcast_ModuleMapCreator_ModuleProcessor[T ModuleMapCreator_ModuleProcessor](v T) ModuleMapCreator_ModuleProcessor {
	var vi ModuleMapCreator_ModuleProcessor = v
	if vi == nil || vi.IsJavaNil_() {
		return nil
	}
	return vi
}

func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.modules.ModuleMapCreator.ModuleProcessor", "com.google.javascript.jscomp.modules.ModuleMapCreator$ModuleProcessor", func(v_ any) bool {
		return j2gort.InstanceOf[ModuleMapCreator_ModuleProcessor](v_)
	})
}

// Java class:
type ModuleMapCreator__Anon_123 struct {
	UnresolvedModuleBase
	outer_         *ModuleMapCreator_ModuleRequestResolverImpl
	cap_0_path     *ModuleLoader_ModulePath
	cap_1_metadata ModuleMetadataMap_ModuleMetadata
}

func (this_ *ModuleMapCreator__Anon_123) BindSelf_(s any) {
	_ = s
	this_.UnresolvedModuleBase.BindSelf_(s)
}
func (this_ *ModuleMapCreator__Anon_123) IsJavaNil_() bool { return this_ == nil }
func (*ModuleMapCreator__Anon_123) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.modules.ModuleMapCreator$ModuleRequestResolverImpl$1", ""
}
func (this_ *ModuleMapCreator__Anon_123) SetOuter_(o *ModuleMapCreator_ModuleRequestResolverImpl) {
	this_.outer_ = o
}
func (this_ *ModuleMapCreator__Anon_123) GetOuter_() *ModuleMapCreator_ModuleRequestResolverImpl {
	return this_.outer_
}
func (*ModuleMapCreator__Anon_123) Implements_jscomp_UnresolvedModule() {}

func (this_ *ModuleMapCreator__Anon_123) initFields_ModuleMapCreator__Anon_123() {
	_ = this_
}

func NewModuleMapCreator__Anon_123(outer_ *ModuleMapCreator_ModuleRequestResolverImpl, cap_p_path *ModuleLoader_ModulePath, cap_p_metadata ModuleMetadataMap_ModuleMetadata) *ModuleMapCreator__Anon_123 {
	this_ := &ModuleMapCreator__Anon_123{}
	this_.BindSelf_(this_)
	this_.outer_ = outer_
	this_.cap_0_path = cap_p_path
	this_.cap_1_metadata = cap_p_metadata
	InitUnresolvedModule(&this_.UnresolvedModuleBase)
	this_.initFields_ModuleMapCreator__Anon_123()
	return this_
}

// Java method: .reset()
func (this_ *ModuleMapCreator__Anon_123) Reset() {
}

// Java method: .resolve(com.google.javascript.jscomp.modules.ModuleRequestResolver moduleRequestResolver, java.lang.@org.jspecify.annotations.Nullable String moduleSpecifier)
func (this_ *ModuleMapCreator__Anon_123) Resolve_2(moduleRequestResolver ModuleRequestResolver, moduleSpecifier string) *Module {
	_ = moduleRequestResolver
	_ = moduleSpecifier
	return Module_S_Builder().BoundNames(j2gort.Ext_com_google_common_collect_ImmutableMap_Of_0()).Namespace(j2gort.Ext_com_google_common_collect_ImmutableMap_Of_0()).LocalNameToLocalExport(j2gort.Ext_com_google_common_collect_ImmutableMap_Of_0()).Path(this_.cap_0_path).Metadata(this_.cap_1_metadata).Build()
}

// Java method: .metadata()
func (this_ *ModuleMapCreator__Anon_123) Metadata() ModuleMetadataMap_ModuleMetadata {
	return this_.cap_1_metadata
}

// Java method: .getExportedNames(com.google.javascript.jscomp.modules.ModuleRequestResolver moduleRequestResolver)
func (this_ *ModuleMapCreator__Anon_123) GetExportedNames_1(moduleRequestResolver ModuleRequestResolver) j2gort.Ext_com_google_common_collect_ImmutableSet {
	_ = moduleRequestResolver
	return j2gort.Ext_com_google_common_collect_ImmutableSet_Of_0()
}

// Java method: .getExportedNames(com.google.javascript.jscomp.modules.ModuleRequestResolver moduleRequestResolver, java.util.Set<com.google.javascript.jscomp.modules.UnresolvedModule> visited)
func (this_ *ModuleMapCreator__Anon_123) GetExportedNames_2(moduleRequestResolver ModuleRequestResolver, visited j2gort.Ext_java_util_Set) j2gort.Ext_com_google_common_collect_ImmutableSet {
	_ = moduleRequestResolver
	_ = visited
	return j2gort.Ext_com_google_common_collect_ImmutableSet_Of_0()
}

// Java method: .resolveExport(com.google.javascript.jscomp.modules.ModuleRequestResolver moduleRequestResolver, java.lang.@org.jspecify.annotations.Nullable String moduleSpecifier, java.lang.String exportName, java.util.Set<com.google.javascript.jscomp.modules.ExportTrace> resolveSet, java.util.Set<com.google.javascript.jscomp.modules.UnresolvedModule> exportStarSet)
func (this_ *ModuleMapCreator__Anon_123) ResolveExport_5(moduleRequestResolver ModuleRequestResolver, moduleSpecifier string, exportName string, resolveSet j2gort.Ext_java_util_Set, exportStarSet j2gort.Ext_java_util_Set) *ResolveExportResult {
	_ = moduleRequestResolver
	_ = moduleSpecifier
	_ = exportName
	_ = resolveSet
	_ = exportStarSet
	return ResolveExportResult_Of(Binding_From_2(Export_S_Builder().LocalName(exportName).ModuleMetadata(this_.cap_1_metadata).ModulePath(this_.cap_0_path).ClosureNamespace(j2gort.NullString).GetExport_BuilderBase().Build(), pkg_rhino.Node(0)))
}
