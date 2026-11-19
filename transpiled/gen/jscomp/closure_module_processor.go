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
// Port of com/google/javascript/jscomp/modules/ClosureModuleProcessor.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.jscomp.modules.ClosureModuleProcessor
type ClosureModuleProcessor struct {
	compiler AbstractCompiler `gson:"compiler"`
	idHash_  atomic.Int32
}

func (this_ *ClosureModuleProcessor) BindSelf_(s any) {
	_ = s
}
func (this_ *ClosureModuleProcessor) IsJavaNil_() bool { return this_ == nil }
func (*ClosureModuleProcessor) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.modules.ClosureModuleProcessor", "ClosureModuleProcessor"
}
func (this_ *ClosureModuleProcessor) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*ClosureModuleProcessor) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.modules.ClosureModuleProcessor", "com.google.javascript.jscomp.modules.ClosureModuleProcessor", func(v_ any) bool {
		return j2gort.InstanceOf[*ClosureModuleProcessor](v_)
	})
}
func (*ClosureModuleProcessor) Implements_jscomp_ModuleMapCreator_ModuleProcessor() {}

func (this_ *ClosureModuleProcessor) initFields_ClosureModuleProcessor() {
	_ = this_
}

// Java constructor: ClosureModuleProcessor(com.google.javascript.jscomp.AbstractCompiler compiler)
func NewClosureModuleProcessor(compiler AbstractCompiler) *ClosureModuleProcessor {
	this_ := &ClosureModuleProcessor{}
	InitClosureModuleProcessor(this_, compiler)
	return this_
}

func InitClosureModuleProcessor(this_ *ClosureModuleProcessor, compiler AbstractCompiler) {
	_ = compiler
	this_.initFields_ClosureModuleProcessor()
	this_.compiler = compiler
}

// Java static method: ClosureModuleProcessor.mayBeAccidentalDefaultExport(java.lang.String importName, java.util.Map<java.lang.String,com.google.javascript.jscomp.modules.Binding> exports)
func ClosureModuleProcessor_MayBeAccidentalDefaultExport(importName string, exports j2gort.Ext_java_util_Map) bool {
	_ = importName
	_ = exports
	var defaultExport pkg_rhino.Node = j2gort.CastPtr[Binding](j2gort.MapGetStr(exports, "*exports*")).OriginatingExport().ExportNode()
	_ = defaultExport
	if !(defaultExport.MatchesName_1_string("exports") && defaultExport.GetParent().IsAssign()) {
		j2gort.CheckStateMsg(false, j2gort.BoxNode(defaultExport))
	}
	var exportedValue pkg_rhino.Node = defaultExport.GetNext()
	_ = exportedValue
	if !(exportedValue.IsObjectLit()) {
		return false
	}
	{
		var key pkg_rhino.Node = exportedValue.GetFirstChild()
		_ = key
		for ; key != 0; key = key.GetNext() {
			if key.IsStringKey() && (j2gort.CheckNotNullString(key.GetString()) == importName) {
				return true
			}
		}
	}
	return false
}

// Java method: ClosureModuleProcessor.process(com.google.javascript.jscomp.modules.ModuleMetadataMap.ModuleMetadata metadata, com.google.javascript.jscomp.deps.ModuleLoader.ModulePath path, com.google.javascript.rhino.Node script)
func (this_ *ClosureModuleProcessor) Process(metadata ModuleMetadataMap_ModuleMetadata, path *ModuleLoader_ModulePath, script pkg_rhino.Node) UnresolvedModule {
	_ = metadata
	_ = path
	_ = script
	if !(script.IsScript() || script.IsCall()) {
		j2gort.CheckArgumentFmt(false, "Unexpected module root %s", j2gort.BoxNode(script))
	}
	if !(script.IsCall() || (path != nil)) {
		j2gort.CheckArgumentStr(false, "Non goog.loadModules must have a path")
	}
	var moduleProcessingCallback *ClosureModuleProcessor_ModuleProcessingCallback = NewClosureModuleProcessor_ModuleProcessingCallback(metadata)
	_ = moduleProcessingCallback
	NodeTraversal_Traverse(this_.compiler, script, moduleProcessingCallback)
	return NewClosureModuleProcessor_UnresolvedGoogModule(metadata, script.GetSourceFileName(), path, moduleProcessingCallback.namespace.ToImmutableMap(), moduleProcessingCallback.requiresByLocalName.ToImmutableMap(), this_.compiler)
}

// Java class: com.google.javascript.jscomp.modules.ClosureModuleProcessor.UnresolvedGoogModule
type ClosureModuleProcessor_UnresolvedGoogModule struct {
	UnresolvedModuleBase
	metadata            ModuleMetadataMap_ModuleMetadata                  `gson:"metadata"`
	srcFileName         string                                            `gson:"srcFileName"`
	path                *ModuleLoader_ModulePath                          `gson:"path"`
	namespace           j2gort.Ext_com_google_common_collect_ImmutableMap `gson:"namespace"`
	requiresByLocalName j2gort.Ext_com_google_common_collect_ImmutableMap `gson:"requiresByLocalName"`
	compiler            AbstractCompiler                                  `gson:"compiler"`
	resolved            *Module                                           `gson:"resolved"`
}

func (this_ *ClosureModuleProcessor_UnresolvedGoogModule) BindSelf_(s any) {
	_ = s
	this_.UnresolvedModuleBase.BindSelf_(s)
}
func (this_ *ClosureModuleProcessor_UnresolvedGoogModule) IsJavaNil_() bool { return this_ == nil }
func (*ClosureModuleProcessor_UnresolvedGoogModule) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.modules.ClosureModuleProcessor$UnresolvedGoogModule", "UnresolvedGoogModule"
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.modules.ClosureModuleProcessor.UnresolvedGoogModule", "com.google.javascript.jscomp.modules.ClosureModuleProcessor$UnresolvedGoogModule", func(v_ any) bool {
		return j2gort.InstanceOf[*ClosureModuleProcessor_UnresolvedGoogModule](v_)
	})
}
func (*ClosureModuleProcessor_UnresolvedGoogModule) Implements_jscomp_UnresolvedModule() {}

func (this_ *ClosureModuleProcessor_UnresolvedGoogModule) initFields_ClosureModuleProcessor_UnresolvedGoogModule() {
	_ = this_
	this_.srcFileName = j2gort.NullString

}

// Java constructor: UnresolvedGoogModule(com.google.javascript.jscomp.modules.ModuleMetadataMap.ModuleMetadata metadata, java.lang.String srcFileName, com.google.javascript.jscomp.deps.ModuleLoader.ModulePath path, com.google.common.collect.ImmutableMap<java.lang.String,com.google.javascript.jscomp.modules.Binding> namespace, com.google.common.collect.ImmutableMap<java.lang.String,com.google.javascript.jscomp.modules.ClosureRequireProcessor.Require> requires, com.google.javascript.jscomp.AbstractCompiler compiler)
func NewClosureModuleProcessor_UnresolvedGoogModule(metadata ModuleMetadataMap_ModuleMetadata, srcFileName string, path *ModuleLoader_ModulePath, namespace j2gort.Ext_com_google_common_collect_ImmutableMap, requires j2gort.Ext_com_google_common_collect_ImmutableMap, compiler AbstractCompiler) *ClosureModuleProcessor_UnresolvedGoogModule {
	this_ := &ClosureModuleProcessor_UnresolvedGoogModule{}
	this_.BindSelf_(this_)
	InitClosureModuleProcessor_UnresolvedGoogModule(this_, metadata, srcFileName, path, namespace, requires, compiler)
	return this_
}

func InitClosureModuleProcessor_UnresolvedGoogModule(this_ *ClosureModuleProcessor_UnresolvedGoogModule, metadata ModuleMetadataMap_ModuleMetadata, srcFileName string, path *ModuleLoader_ModulePath, namespace j2gort.Ext_com_google_common_collect_ImmutableMap, requires j2gort.Ext_com_google_common_collect_ImmutableMap, compiler AbstractCompiler) {
	_ = metadata
	_ = srcFileName
	_ = path
	_ = namespace
	_ = requires
	_ = compiler
	InitUnresolvedModule(&this_.UnresolvedModuleBase)
	this_.initFields_ClosureModuleProcessor_UnresolvedGoogModule()
	this_.metadata = metadata
	this_.srcFileName = srcFileName
	this_.path = path
	this_.namespace = namespace
	this_.requiresByLocalName = requires
	this_.compiler = compiler
}

// Java method: UnresolvedGoogModule.resolveExport(com.google.javascript.jscomp.modules.ModuleRequestResolver moduleRequestResolver, java.lang.String exportName)
func (this_ *ClosureModuleProcessor_UnresolvedGoogModule) ResolveExport_2(moduleRequestResolver ModuleRequestResolver, exportName string) *ResolveExportResult {
	_ = moduleRequestResolver
	_ = exportName
	if j2gort.MapContainsKeyStr(this_.namespace, exportName) {
		return ResolveExportResult_Of(j2gort.CastPtr[Binding](j2gort.MapGetStr(this_.namespace, exportName)))
	}
	return StaticGet_ResolveExportResult_NOT_FOUND()
}

// Java method: UnresolvedGoogModule.resolveExport(com.google.javascript.jscomp.modules.ModuleRequestResolver moduleRequestResolver, java.lang.@org.jspecify.annotations.Nullable String moduleSpecifier, java.lang.String exportName, java.util.Set<com.google.javascript.jscomp.modules.ExportTrace> resolveSet, java.util.Set<com.google.javascript.jscomp.modules.UnresolvedModule> exportStarSet)
func (this_ *ClosureModuleProcessor_UnresolvedGoogModule) ResolveExport_5(moduleRequestResolver ModuleRequestResolver, moduleSpecifier string, exportName string, resolveSet j2gort.Ext_java_util_Set, exportStarSet j2gort.Ext_java_util_Set) *ResolveExportResult {
	_ = moduleRequestResolver
	_ = moduleSpecifier
	_ = exportName
	_ = resolveSet
	_ = exportStarSet
	return this_.ResolveExport_2(moduleRequestResolver, exportName)
}

// Java method: UnresolvedGoogModule.resolve(com.google.javascript.jscomp.modules.ModuleRequestResolver moduleRequestResolver, java.lang.@org.jspecify.annotations.Nullable String moduleSpecifier)
func (this_ *ClosureModuleProcessor_UnresolvedGoogModule) Resolve_2(moduleRequestResolver ModuleRequestResolver, moduleSpecifier string) *Module {
	_ = moduleRequestResolver
	_ = moduleSpecifier
	if this_.resolved == nil {
		var boundNames j2gort.Ext_java_util_Map = j2gort.NewExt_java_util_LinkedHashMap_1_Ext_java_util_Map(this_.GetAllResolvedImports(moduleRequestResolver))
		_ = boundNames
		this_.resolved = Module_S_Builder().Path(this_.path).Metadata(this_.metadata).Namespace(this_.namespace).BoundNames(j2gort.Ext_com_google_common_collect_ImmutableMap_CopyOf_1_Ext_java_util_Map(boundNames)).LocalNameToLocalExport(j2gort.Ext_com_google_common_collect_ImmutableMap_Of_0()).ClosureNamespace(j2gort.CastString(j2gort.Ext_com_google_common_collect_Iterables_GetOnlyElement_1(j2gort.AsIface[j2gort.Ext_java_lang_Iterable](this_.metadata.GoogNamespaces())))).Build()
	}
	return this_.resolved
}

// Java method: UnresolvedGoogModule.getAllResolvedImports(com.google.javascript.jscomp.modules.ModuleRequestResolver moduleRequestResolver)
func (this_ *ClosureModuleProcessor_UnresolvedGoogModule) GetAllResolvedImports(moduleRequestResolver ModuleRequestResolver) j2gort.Ext_java_util_Map {
	_ = moduleRequestResolver
	var imports j2gort.Ext_java_util_Map = j2gort.NewExt_java_util_LinkedHashMap()
	_ = imports
	{
		cur_1 := j2gort.NewMapKeyStrCursor(this_.requiresByLocalName)
		for cur_1.HasNext() {
			name := cur_1.NextStr()
			_ = name
			var b *ResolveExportResult = this_.ResolveImport(moduleRequestResolver, name)
			_ = b
			if b.Resolved() {
				j2gort.MapPutStr(imports, name, j2gort.PtrToAny(b.GetBinding()))
			}
		}
	}
	return imports
}

// Java method: UnresolvedGoogModule.resolveImport(com.google.javascript.jscomp.modules.ModuleRequestResolver moduleRequestResolver, java.lang.String name)
func (this_ *ClosureModuleProcessor_UnresolvedGoogModule) ResolveImport(moduleRequestResolver ModuleRequestResolver, name string) *ResolveExportResult {
	_ = moduleRequestResolver
	_ = name
	var require *ClosureRequireProcessor_Require = j2gort.CastPtr[ClosureRequireProcessor_Require](j2gort.MapGetStr(this_.requiresByLocalName, name))
	_ = require
	var importRecord *Import = require.ImportRecord()
	_ = importRecord
	var requested UnresolvedModule = moduleRequestResolver.Resolve_1_Import(importRecord)
	_ = requested
	if requested == nil || requested.IsJavaNil_() {
		return StaticGet_ResolveExportResult_ERROR()
	} else {
		var iifeCond_1 bool
		{

			sr_2 := importRecord.ImportName()
			se_3 := "*exports*"
			iifeCond_1 = j2gort.CheckNotNullString(sr_2) == se_3
		}
		if iifeCond_1 {
			return ResolveExportResult_Of(Binding_From_4(requested.Metadata(), importRecord.NameNode(), importRecord.ModuleRequest(), require.CreatedBy()))
		} else {
			var result *ResolveExportResult = requested.ResolveExport_5(moduleRequestResolver, importRecord.ModuleRequest(), importRecord.ImportName(), j2gort.NewExt_java_util_LinkedHashSet(), j2gort.NewExt_java_util_LinkedHashSet())
			_ = result
			if !(result.Found()) && !(result.HadError()) {
				this_.priv_UnresolvedGoogModule_reportInvalidDestructuringRequire_6670(requested, importRecord)
				return StaticGet_ResolveExportResult_ERROR()
			}
			var forSourceInfo pkg_rhino.Node
			_ = forSourceInfo
			if importRecord.NameNode() == 0 {
				forSourceInfo = importRecord.ImportNode()
			} else {
				forSourceInfo = importRecord.NameNode()
			}
			return result.Copy(forSourceInfo, require.CreatedBy())
		}
	}
}

// Java method: UnresolvedGoogModule.metadata()
func (this_ *ClosureModuleProcessor_UnresolvedGoogModule) Metadata() ModuleMetadataMap_ModuleMetadata {
	return this_.metadata
}

// Java method: UnresolvedGoogModule.getExportedNames(com.google.javascript.jscomp.modules.ModuleRequestResolver moduleRequestResolver)
func (this_ *ClosureModuleProcessor_UnresolvedGoogModule) GetExportedNames_1(moduleRequestResolver ModuleRequestResolver) j2gort.Ext_com_google_common_collect_ImmutableSet {
	_ = moduleRequestResolver
	panic(j2gort.AsIface[javalang.IUnsupportedOperationException](j2gort.ConstructShimThrowable("java.lang.UnsupportedOperationException")))
}

// Java method: UnresolvedGoogModule.getExportedNames(com.google.javascript.jscomp.modules.ModuleRequestResolver moduleRequestResolver, java.util.Set<com.google.javascript.jscomp.modules.UnresolvedModule> visited)
func (this_ *ClosureModuleProcessor_UnresolvedGoogModule) GetExportedNames_2(moduleRequestResolver ModuleRequestResolver, visited j2gort.Ext_java_util_Set) j2gort.Ext_com_google_common_collect_ImmutableSet {
	_ = moduleRequestResolver
	_ = visited
	panic(j2gort.AsIface[javalang.IUnsupportedOperationException](j2gort.ConstructShimThrowable("java.lang.UnsupportedOperationException")))
}

// Java method: UnresolvedGoogModule.reset()
func (this_ *ClosureModuleProcessor_UnresolvedGoogModule) Reset() {
	this_.resolved = nil
}

// Java method: UnresolvedGoogModule.reportInvalidDestructuringRequire(com.google.javascript.jscomp.modules.UnresolvedModule requested, com.google.javascript.jscomp.modules.Import importRecord)
func (this_ *ClosureModuleProcessor_UnresolvedGoogModule) priv_UnresolvedGoogModule_reportInvalidDestructuringRequire_6670(requested UnresolvedModule, importRecord *Import) {
	_ = requested
	_ = importRecord
	var additionalInfo string = ""
	_ = additionalInfo
	var unresolvedGoogModule *ClosureModuleProcessor_UnresolvedGoogModule
	_ = unresolvedGoogModule
	var ok_4 bool
	unresolvedGoogModule, ok_4 = j2gort.TryCast[*ClosureModuleProcessor_UnresolvedGoogModule](requested)
	if ok_4 {
		var exports j2gort.Ext_com_google_common_collect_ImmutableMap = unresolvedGoogModule.namespace
		_ = exports
		if j2gort.MapContainsKeyStr(exports, "*exports*") {
			additionalInfo = j2gort.Ext_com_google_common_base_Strings_LenientFormat("\nThe goog.module \"%s\" cannot be destructured as it contains a default export, not named exports. See %s.", []any{j2gort.BoxString(importRecord.ModuleRequest()), "https://github.com/google/closure-library/wiki/goog.module%3A-an-ES6-module-like-alternative-to-goog.provide#destructuring-imports"})
			if ClosureModuleProcessor_MayBeAccidentalDefaultExport(importRecord.ImportName(), j2gort.AsIface[j2gort.Ext_java_util_Map](exports)) {
				additionalInfo = (j2gort.StrValString(additionalInfo) + j2gort.StrValString(j2gort.Ext_com_google_common_base_Strings_LenientFormat("\nEither use a non-destructuring require or rewrite the goog.module \"%s\" to support destructuring requires. For example, consider replacing\n  exports = {%s: <value>[, ...]};\nwith individual named export assignments like\n  exports.%s = <value>;\n", []any{j2gort.BoxString(importRecord.ModuleRequest()), j2gort.BoxString(importRecord.ImportName()), j2gort.BoxString(importRecord.ImportName())})))
			}
		}
	}
	this_.compiler.Report_1(JSError_Make_5(this_.srcFileName, importRecord.ImportNode().GetLineno(), importRecord.ImportNode().GetCharno(), StaticGet_ModuleMapCreator_DOES_NOT_HAVE_EXPORT_WITH_DETAILS(), []string{importRecord.ImportName(), additionalInfo}))
}

// Java class: com.google.javascript.jscomp.modules.ClosureModuleProcessor.ModuleProcessingCallback
type ClosureModuleProcessor_ModuleProcessingCallback struct {
	NodeTraversal_AbstractPreOrderCallbackBase
	metadata              ModuleMetadataMap_ModuleMetadata                            `gson:"metadata"`
	closureNamespace      string                                                      `gson:"closureNamespace"`
	namespace             j2gort.OrderedMap[string, *Binding]                         `gson:"namespace"`
	requiresByLocalName   j2gort.OrderedMap[string, *ClosureRequireProcessor_Require] `gson:"requiresByLocalName"`
	seenExportsAssignment bool                                                        `gson:"seenExportsAssignment"`
}

func (this_ *ClosureModuleProcessor_ModuleProcessingCallback) BindSelf_(s any) {
	_ = s
}
func (this_ *ClosureModuleProcessor_ModuleProcessingCallback) IsJavaNil_() bool { return this_ == nil }
func (*ClosureModuleProcessor_ModuleProcessingCallback) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.modules.ClosureModuleProcessor$ModuleProcessingCallback", "ModuleProcessingCallback"
}
func (this_ *ClosureModuleProcessor_ModuleProcessingCallback) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.NodeTraversal_AbstractPreOrderCallbackBase.IdentityHash_()
}
func (*ClosureModuleProcessor_ModuleProcessingCallback) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.modules.ClosureModuleProcessor.ModuleProcessingCallback", "com.google.javascript.jscomp.modules.ClosureModuleProcessor$ModuleProcessingCallback", func(v_ any) bool {
		return j2gort.InstanceOf[*ClosureModuleProcessor_ModuleProcessingCallback](v_)
	})
}
func (*ClosureModuleProcessor_ModuleProcessingCallback) Implements_jscomp_NodeTraversal_AbstractPreOrderCallback() {
}
func (*ClosureModuleProcessor_ModuleProcessingCallback) Implements_jscomp_NodeTraversal_Callback() {}

func (this_ *ClosureModuleProcessor_ModuleProcessingCallback) initFields_ClosureModuleProcessor_ModuleProcessingCallback() {
	_ = this_
	this_.closureNamespace = j2gort.NullString
}

// Java constructor: ModuleProcessingCallback(com.google.javascript.jscomp.modules.ModuleMetadataMap.ModuleMetadata metadata)
func NewClosureModuleProcessor_ModuleProcessingCallback(metadata ModuleMetadataMap_ModuleMetadata) *ClosureModuleProcessor_ModuleProcessingCallback {
	this_ := &ClosureModuleProcessor_ModuleProcessingCallback{}
	InitClosureModuleProcessor_ModuleProcessingCallback(this_, metadata)
	return this_
}

func InitClosureModuleProcessor_ModuleProcessingCallback(this_ *ClosureModuleProcessor_ModuleProcessingCallback, metadata ModuleMetadataMap_ModuleMetadata) {
	_ = metadata
	InitNodeTraversal_AbstractPreOrderCallback(&this_.NodeTraversal_AbstractPreOrderCallbackBase)
	this_.initFields_ClosureModuleProcessor_ModuleProcessingCallback()
	this_.metadata = metadata
	this_.namespace = j2gort.NewOrderedMap[string, *Binding]()
	this_.requiresByLocalName = j2gort.NewOrderedMap[string, *ClosureRequireProcessor_Require]()
	this_.closureNamespace = j2gort.CastString(j2gort.Ext_com_google_common_collect_Iterables_GetOnlyElement_1(j2gort.AsIface[j2gort.Ext_java_lang_Iterable](metadata.GoogNamespaces())))
	this_.seenExportsAssignment = false
}

// Java method: ModuleProcessingCallback.shouldTraverse(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *ClosureModuleProcessor_ModuleProcessingCallback) ShouldTraverse_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) bool {
	_ = t
	_ = n
	_ = parent
	switch n.TokenSwitchOrd_() {
	case int32(123) /* Token.MODULE_BODY */, int32(94) /* Token.SCRIPT */, int32(28) /* Token.CALL */, int32(90) /* Token.BLOCK */ :
		return true
	case int32(73) /* Token.FUNCTION */ :
		return (parent.IsCall() && (parent == this_.metadata.RootNode()))
	case int32(93) /* Token.EXPR_RESULT */ :
		var expr pkg_rhino.Node = n.GetFirstChild()
		_ = expr
		if expr.IsAssign() {
			this_.priv_ModuleProcessingCallback_maybeInitializeExports_6671(expr)
		} else if expr.IsGetProp() {
			this_.priv_ModuleProcessingCallback_maybeInitializeExportsStub_6672(expr)
		}
		return false
	case int32(97) /* Token.CONST */, int32(84) /* Token.VAR */, int32(111) /* Token.LET */ :
		this_.priv_ModuleProcessingCallback_maybeInitializeRequire_6676(n)
		return false
	default:
		return false
	}
	panic("unreachable switch expression")
}

// Java method: ModuleProcessingCallback.maybeInitializeExports(com.google.javascript.rhino.Node assignment)
func (this_ *ClosureModuleProcessor_ModuleProcessingCallback) priv_ModuleProcessingCallback_maybeInitializeExports_6671(assignment pkg_rhino.Node) {
	_ = assignment
	var lhs pkg_rhino.Node = assignment.GetFirstChild()
	_ = lhs
	var rhs pkg_rhino.Node = assignment.GetSecondChild()
	_ = rhs
	if lhs.IsName() && (j2gort.CheckNotNullString(lhs.GetString()) == "exports") {
		if NodeUtil_IsNamedExportsLiteral(rhs) {
			this_.priv_ModuleProcessingCallback_initializeNamedExportsLiteral_6674(rhs)
		} else {
			this_.seenExportsAssignment = true
			this_.priv_ModuleProcessingCallback_markExportsAssignmentInNamespace_6673(lhs)
		}
	} else if (lhs.IsGetProp() && lhs.GetFirstChild().IsName()) && (j2gort.CheckNotNullString(lhs.GetFirstChild().GetString()) == "exports") {
		var exportedId string = lhs.GetString()
		_ = exportedId
		this_.priv_ModuleProcessingCallback_addPropertyExport_6675(exportedId, lhs)
	}
}

// Java method: ModuleProcessingCallback.maybeInitializeExportsStub(com.google.javascript.rhino.Node qname)
func (this_ *ClosureModuleProcessor_ModuleProcessingCallback) priv_ModuleProcessingCallback_maybeInitializeExportsStub_6672(qname pkg_rhino.Node) {
	_ = qname
	var owner pkg_rhino.Node = qname.GetFirstChild()
	_ = owner
	if owner.IsName() && (j2gort.CheckNotNullString(owner.GetString()) == "exports") {
		this_.priv_ModuleProcessingCallback_addPropertyExport_6675(qname.GetString(), qname)
	}
}

// Java method: ModuleProcessingCallback.markExportsAssignmentInNamespace(com.google.javascript.rhino.Node exportsNode)
func (this_ *ClosureModuleProcessor_ModuleProcessingCallback) priv_ModuleProcessingCallback_markExportsAssignmentInNamespace_6673(exportsNode pkg_rhino.Node) {
	_ = exportsNode
	this_.namespace.Put("*exports*", Binding_From_2(Export_S_Builder().ExportName("*exports*").ExportNode(exportsNode).ModuleMetadata(this_.metadata).ClosureNamespace(this_.closureNamespace).ModulePath(this_.metadata.Path()).GetExport_BuilderBase().Build(), exportsNode))
}

// Java method: ModuleProcessingCallback.initializeNamedExportsLiteral(com.google.javascript.rhino.Node objectLit)
func (this_ *ClosureModuleProcessor_ModuleProcessingCallback) priv_ModuleProcessingCallback_initializeNamedExportsLiteral_6674(objectLit pkg_rhino.Node) {
	_ = objectLit
	{
		var key pkg_rhino.Node = objectLit.GetFirstChild()
		_ = key
		for ; key != 0; key = key.GetNext() {
			this_.priv_ModuleProcessingCallback_addPropertyExport_6675(key.GetString(), key)
		}
	}
}

// Java method: ModuleProcessingCallback.addPropertyExport(java.lang.String exportedId, com.google.javascript.rhino.Node propNode)
func (this_ *ClosureModuleProcessor_ModuleProcessingCallback) priv_ModuleProcessingCallback_addPropertyExport_6675(exportedId string, propNode pkg_rhino.Node) {
	_ = exportedId
	_ = propNode
	if this_.seenExportsAssignment {
		return
	} else if this_.namespace.ContainsKey(exportedId) {
		return
	}
	this_.namespace.Put(exportedId, Binding_From_2(Export_S_Builder().ExportName(exportedId).ExportNode(propNode).ModuleMetadata(this_.metadata).ClosureNamespace(this_.closureNamespace).ModulePath(this_.metadata.Path()).GetExport_BuilderBase().Build(), propNode))
}

// Java method: ModuleProcessingCallback.maybeInitializeRequire(com.google.javascript.rhino.Node nameDeclaration)
func (this_ *ClosureModuleProcessor_ModuleProcessingCallback) priv_ModuleProcessingCallback_maybeInitializeRequire_6676(nameDeclaration pkg_rhino.Node) {
	_ = nameDeclaration
	{
		cur_5 := j2gort.NewIterCursor(j2gort.CheckNotNull(ClosureRequireProcessor_GetAllRequires(nameDeclaration)))
		for cur_5.HasNext() {
			rawElem_6 := cur_5.Next()
			require := j2gort.CastPtr[ClosureRequireProcessor_Require](rawElem_6)
			_ = require
			this_.requiresByLocalName.PutIfAbsent(require.LocalName(), require)
		}
	}
}
func (*ClosureModuleProcessor_ModuleProcessingCallback) NoOpVisit_()                 {}
func (*ClosureModuleProcessor_ModuleProcessingCallback) CanReleaseSyntacticScopes_() {}
