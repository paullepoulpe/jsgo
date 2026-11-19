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
// Port of com/google/javascript/jscomp/ModuleImportResolver.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.jscomp.ModuleImportResolver
type ModuleImportResolver struct {
	moduleMap         *ModuleMap                             `gson:"moduleMap"`
	nodeToScopeMapper j2gort.Ext_java_util_function_Function `gson:"nodeToScopeMapper"`
	registry          *pkg_rhino.JSTypeRegistry              `gson:"registry"`
	idHash_           atomic.Int32
}

func (this_ *ModuleImportResolver) BindSelf_(s any) {
	_ = s
}
func (this_ *ModuleImportResolver) IsJavaNil_() bool { return this_ == nil }
func (*ModuleImportResolver) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.ModuleImportResolver", "ModuleImportResolver"
}
func (this_ *ModuleImportResolver) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*ModuleImportResolver) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.ModuleImportResolver", "com.google.javascript.jscomp.ModuleImportResolver", func(v_ any) bool {
		return j2gort.InstanceOf[*ModuleImportResolver](v_)
	})
}

func StaticGet_ModuleImportResolver_GOOG() string { return "goog" }

var static_ModuleImportResolver_GOOG_DEPENDENCY_CALLS j2gort.Ext_com_google_common_collect_ImmutableSet

func StaticGet_ModuleImportResolver_GOOG_DEPENDENCY_CALLS() j2gort.Ext_com_google_common_collect_ImmutableSet {
	clinit_ModuleImportResolver()
	return static_ModuleImportResolver_GOOG_DEPENDENCY_CALLS
}
func StaticSet_ModuleImportResolver_GOOG_DEPENDENCY_CALLS(v j2gort.Ext_com_google_common_collect_ImmutableSet) {
	clinit_ModuleImportResolver()
	static_ModuleImportResolver_GOOG_DEPENDENCY_CALLS = v
}

var static_ModuleImportResolver_GOOG_MODULE_GET pkg_rhino.QualifiedName

func StaticGet_ModuleImportResolver_GOOG_MODULE_GET() pkg_rhino.QualifiedName {
	clinit_ModuleImportResolver()
	return static_ModuleImportResolver_GOOG_MODULE_GET
}
func StaticSet_ModuleImportResolver_GOOG_MODULE_GET(v pkg_rhino.QualifiedName) {
	clinit_ModuleImportResolver()
	static_ModuleImportResolver_GOOG_MODULE_GET = v
}

var clinitDone_ModuleImportResolver uint32
var clinitState_ModuleImportResolver j2gort.ClinitState

func clinit_ModuleImportResolver() {
	if atomic.LoadUint32(&clinitDone_ModuleImportResolver) != 0 {
		return
	}
	clinitSlow_ModuleImportResolver()
}

func clinitSlow_ModuleImportResolver() {
	if !j2gort.EnterClinit(&clinitState_ModuleImportResolver, &clinitDone_ModuleImportResolver) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_ModuleImportResolver, &clinitDone_ModuleImportResolver)
	static_ModuleImportResolver_GOOG_DEPENDENCY_CALLS = j2gort.Ext_com_google_common_collect_ImmutableSet_Of_4("require", "requireType", "forwardDeclare", "requireDynamic")
	static_ModuleImportResolver_GOOG_MODULE_GET = pkg_rhino.QualifiedName_Of("goog.module.get")
}

func (this_ *ModuleImportResolver) initFields_ModuleImportResolver() {
	_ = this_
}

// Java constructor: ModuleImportResolver(com.google.javascript.jscomp.modules.ModuleMap moduleMap, java.util.function.Function<com.google.javascript.rhino.Node,com.google.javascript.jscomp.TypedScope> nodeToScopeMapper, com.google.javascript.rhino.jstype.JSTypeRegistry registry)
func NewModuleImportResolver(moduleMap *ModuleMap, nodeToScopeMapper j2gort.Ext_java_util_function_Function, registry *pkg_rhino.JSTypeRegistry) *ModuleImportResolver {
	this_ := &ModuleImportResolver{}
	InitModuleImportResolver(this_, moduleMap, nodeToScopeMapper, registry)
	return this_
}

func InitModuleImportResolver(this_ *ModuleImportResolver, moduleMap *ModuleMap, nodeToScopeMapper j2gort.Ext_java_util_function_Function, registry *pkg_rhino.JSTypeRegistry) {
	_ = moduleMap
	_ = nodeToScopeMapper
	_ = registry
	this_.initFields_ModuleImportResolver()
	this_.moduleMap = moduleMap
	this_.nodeToScopeMapper = nodeToScopeMapper
	this_.registry = registry
}

// Java static method: ModuleImportResolver.isGoogModuleDependencyCall(com.google.javascript.rhino.Node value)
func ModuleImportResolver_IsGoogModuleDependencyCall(value pkg_rhino.Node) bool {
	_ = value
	if (((value == 0) || !(value.IsCall())) || !(value.HasTwoChildren())) || !(value.GetSecondChild().IsStringLit()) {
		return false
	}
	var callee pkg_rhino.Node = value.GetFirstChild()
	_ = callee
	if !(callee.IsGetProp()) {
		return false
	}
	var owner pkg_rhino.Node = callee.GetFirstChild()
	_ = owner
	return (((owner.IsName() && func() bool { sr_1 := owner.GetString(); se_2 := "goog"; return j2gort.CheckNotNullString(sr_1) == se_2 }()) && j2gort.SetContainsStr(StaticGet_ModuleImportResolver_GOOG_DEPENDENCY_CALLS(), callee.GetString())) || StaticGet_ModuleImportResolver_GOOG_MODULE_GET().Matches(callee))
}

// Java method: ModuleImportResolver.getClosureNamespaceTypeFromCall(com.google.javascript.rhino.Node googRequire)
func (this_ *ModuleImportResolver) GetClosureNamespaceTypeFromCall(googRequire pkg_rhino.Node) ScopedName {
	_ = googRequire
	if this_.moduleMap == nil {
		return nil
	}
	var moduleId string = googRequire.GetSecondChild().GetString()
	_ = moduleId
	return this_.priv_ModuleImportResolver_getScopedNameForClosureNamespace_3892(moduleId)
}

// Java method: ModuleImportResolver.getScopedNameForClosureNamespace(java.lang.String moduleId)
func (this_ *ModuleImportResolver) priv_ModuleImportResolver_getScopedNameForClosureNamespace_3892(moduleId string) ScopedName {
	_ = moduleId
	var module *Module = this_.moduleMap.GetClosureModule(moduleId)
	_ = module
	if module == nil {
		return nil
	}
	switch (module.Metadata().ModuleType()).Ordinal() {
	case int32(4) /* ModuleType.GOOG_PROVIDE */ :
		var provide pkg_rhino.Node = module.Metadata().RootNode()
		_ = provide
		if (provide != 0) && provide.IsScript() {
			return ScopedName_Of(moduleId, provide.GetGrandparent())
		} else {
			return nil
		}
	case int32(1) /* ModuleType.GOOG_MODULE */, int32(2) /* ModuleType.LEGACY_GOOG_MODULE */ :
		var scopeRoot pkg_rhino.Node = this_.priv_ModuleImportResolver_getGoogModuleScopeRoot_3893(module)
		_ = scopeRoot
		if scopeRoot != 0 {
			return ScopedName_Of("exports", scopeRoot)
		} else {
			return nil
		}
	case int32(0) /* ModuleType.ES6_MODULE */ :
		var moduleBody pkg_rhino.Node = module.Metadata().RootNode().GetFirstChild()
		_ = moduleBody
		return ScopedName_Of("*exports*", moduleBody)
	case int32(3) /* ModuleType.COMMON_JS */ :
		panic(j2gort.AsIface[javalang.IIllegalStateException](j2gort.ConstructShimThrowable("java.lang.IllegalStateException", "Type checking CommonJs modules not yet supported")))
	case int32(5) /* ModuleType.SCRIPT */ :
		panic(j2gort.AsIface[javalang.IIllegalStateException](j2gort.ConstructShimThrowable("java.lang.IllegalStateException", "Cannot import a name from a SCRIPT")))
	}
	panic(j2gort.AsIface[javalang.IAssertionError](j2gort.ConstructShimThrowable("java.lang.AssertionError")))
}

// Java method: ModuleImportResolver.getGoogModuleScopeRoot(com.google.javascript.jscomp.modules.@org.jspecify.annotations.Nullable Module module)
func (this_ *ModuleImportResolver) priv_ModuleImportResolver_getGoogModuleScopeRoot_3893(module *Module) pkg_rhino.Node {
	_ = module
	j2gort.CheckArgumentMsg(module.Metadata().GetModuleMetadataMap_ModuleMetadataBase().IsGoogModule(), module.Metadata())
	var scriptNode pkg_rhino.Node = module.Metadata().RootNode()
	_ = scriptNode
	if (scriptNode.IsScript() && scriptNode.HasOneChild()) && scriptNode.GetOnlyChild().IsModuleBody() {
		return scriptNode.GetOnlyChild()
	} else if scriptNode.IsCall() {
		var functionLiteral pkg_rhino.Node = scriptNode.GetSecondChild()
		_ = functionLiteral
		return NodeUtil_GetFunctionBody(functionLiteral)
	}
	return pkg_rhino.Node(0)
}

// Java method: ModuleImportResolver.declareEsModuleImports(com.google.javascript.jscomp.modules.Module module, com.google.javascript.jscomp.TypedScope scope, com.google.javascript.jscomp.CompilerInput moduleInput)
func (this_ *ModuleImportResolver) DeclareEsModuleImports(module *Module, scope *TypedScope, moduleInput *CompilerInput) j2gort.Ext_com_google_common_collect_ImmutableMap {
	_ = module
	_ = scope
	_ = moduleInput
	if !(module.Metadata().GetModuleMetadataMap_ModuleMetadataBase().IsEs6Module()) {
		j2gort.CheckArgumentMsg(false, j2gort.PtrToAny(module))
	}
	if !(scope.IsModuleScope()) {
		j2gort.CheckArgumentMsg(false, j2gort.PtrToAny(scope))
	}
	var missingNames j2gort.Ext_com_google_common_collect_ImmutableMap_Builder = j2gort.Ext_com_google_common_collect_ImmutableMap_S_Builder()
	_ = missingNames
	{
		cur_3 := j2gort.NewMapEntryCursor(module.BoundNames())
		for cur_3.NextEntry() {
			var binding *Binding = j2gort.CastPtr[Binding](cur_3.GetValue())
			_ = binding
			var localName string = cur_3.KeyStr()
			_ = localName
			if !(binding.IsCreatedByEsImport()) {
				continue
			}
			var export ScopedName = this_.priv_ModuleImportResolver_getScopedNameFromEsBinding_3897(binding)
			_ = export
			var modScope *TypedScope = j2gort.CastPtr[TypedScope](this_.nodeToScopeMapper.Apply(j2gort.BoxNode(export.GetScopedNameBase().GetScopeRoot())))
			_ = modScope
			if modScope == nil {
				j2gort.CheckStateMsg((j2gort.CheckNotNullString(binding.SourceNode().GetString()) == localName), j2gort.BoxNode(binding.SourceNode()))
				missingNames.Put_2(j2gort.BoxNode(binding.SourceNode()), export)
				continue
			}
			var originalVar *TypedVar = j2gort.CastPtr[TypedVar](modScope.GetVar(export.GetScopedNameBase().GetName()))
			_ = originalVar
			var importType pkg_rhino.JSType = originalVar.GetType()
			_ = importType
			scope.Declare(localName, binding.SourceNode(), importType, moduleInput, originalVar.IsTypeInferred())
			if !(binding.IsModuleNamespace()) && j2gort.IfaceEquals[pkg_rhino.JSType](binding.SourceNode().GetTypedefTypeProp(), nil) {
				var typedefType pkg_rhino.JSType = originalVar.GetNameNode().GetTypedefTypeProp()
				_ = typedefType
				if typedefType != nil && !typedefType.IsJavaNil_() {
					binding.SourceNode().SetTypedefTypeProp(typedefType)
					this_.registry.DeclareType(pkg_rhino.ToIface_StaticScope(scope), localName, typedefType)
				}
			}
		}
	}
	return j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap](missingNames.BuildOrThrow())
}

// Java method: ModuleImportResolver.updateEsModuleNamespaceType(com.google.javascript.rhino.jstype.ObjectType namespace, com.google.javascript.jscomp.modules.Module module, com.google.javascript.jscomp.TypedScope scope)
func (this_ *ModuleImportResolver) UpdateEsModuleNamespaceType(namespace pkg_rhino.ObjectType, module *Module, scope *TypedScope) {
	_ = namespace
	_ = module
	_ = scope
	if !(module.Metadata().GetModuleMetadataMap_ModuleMetadataBase().IsEs6Module()) {
		j2gort.CheckArgumentMsg(false, j2gort.PtrToAny(module))
	}
	if !(scope.IsModuleScope()) {
		j2gort.CheckArgumentMsg(false, j2gort.PtrToAny(scope))
	}
	{
		cur_4 := j2gort.NewMapEntryCursor(module.Namespace())
		for cur_4.NextEntry() {
			var exportKey string = cur_4.KeyStr()
			_ = exportKey
			if namespace.GetObjectTypeBase().IsPropertyTypeDeclared_1_string(exportKey) {
				continue
			}
			var binding *Binding = j2gort.CastPtr[Binding](cur_4.GetValue())
			_ = binding
			var bindingSourceNode pkg_rhino.Node = binding.SourceNode()
			_ = bindingSourceNode
			var export ScopedName = this_.priv_ModuleImportResolver_getScopedNameFromEsBinding_3897(binding)
			_ = export
			var originalScope *TypedScope
			_ = originalScope
			if export.GetScopedNameBase().GetScopeRoot() == scope.GetRootNode_0() {
				originalScope = scope
			} else {
				originalScope = j2gort.CastPtr[TypedScope](this_.nodeToScopeMapper.Apply(j2gort.BoxNode(export.GetScopedNameBase().GetScopeRoot())))
			}
			if originalScope == nil {
				namespace.GetObjectTypeBase().DefineInferredProperty_3_string_JSType_Node(exportKey, this_.registry.GetNativeType(pkg_rhino.StaticGet_JSTypeNative_UNKNOWN_TYPE()), bindingSourceNode)
				this_.priv_ModuleImportResolver_updateAstForExport_3896(bindingSourceNode, this_.registry.GetNativeType(pkg_rhino.StaticGet_JSTypeNative_UNKNOWN_TYPE()), nil)
				continue
			}
			var originalName *TypedVar = j2gort.CastPtr[TypedVar](originalScope.GetSlot(export.GetScopedNameBase().GetName()))
			_ = originalName
			var exportType pkg_rhino.JSType = originalName.GetType()
			_ = exportType
			if exportType == nil || exportType.IsJavaNil_() {
				exportType = this_.registry.GetNativeType(pkg_rhino.StaticGet_JSTypeNative_NO_TYPE())
			}
			if originalName.IsTypeInferred() {
				namespace.GetObjectTypeBase().DefineInferredProperty_3_string_JSType_Node(exportKey, exportType, bindingSourceNode)
			} else {
				namespace.GetObjectTypeBase().DefineDeclaredProperty_3_string_JSType_Node(exportKey, exportType, bindingSourceNode)
			}
			this_.priv_ModuleImportResolver_updateAstForExport_3896(bindingSourceNode, exportType, originalName.GetNameNode().GetTypedefTypeProp())
		}
	}
}

// Java method: ModuleImportResolver.updateAstForExport(com.google.javascript.rhino.Node bindingSourceNode, com.google.javascript.rhino.jstype.JSType exportType, com.google.javascript.rhino.jstype.@org.jspecify.annotations.Nullable JSType typedefType)
func (this_ *ModuleImportResolver) priv_ModuleImportResolver_updateAstForExport_3896(bindingSourceNode pkg_rhino.Node, exportType pkg_rhino.JSType, typedefType pkg_rhino.JSType) {
	_ = bindingSourceNode
	_ = exportType
	_ = typedefType
	bindingSourceNode.SetJSType(exportType)
	bindingSourceNode.SetTypedefTypeProp(typedefType)
	if bindingSourceNode.GetParent().IsExportSpec() {
		bindingSourceNode.GetNext().SetJSType(exportType)
	}
}

// Java method: ModuleImportResolver.getScopedNameFromEsBinding(com.google.javascript.jscomp.modules.Binding binding)
func (this_ *ModuleImportResolver) priv_ModuleImportResolver_getScopedNameFromEsBinding_3897(binding *Binding) ScopedName {
	_ = binding
	var name string
	_ = name
	if binding.IsModuleNamespace() {
		name = "*exports*"
	} else {
		name = binding.BoundName()
	}
	var originalMetadata ModuleMetadataMap_ModuleMetadata
	_ = originalMetadata
	if binding.IsModuleNamespace() {
		originalMetadata = binding.Metadata()
	} else {
		originalMetadata = binding.OriginatingExport().ModuleMetadata()
	}
	switch (originalMetadata.ModuleType()).Ordinal() {
	case int32(0) /* ModuleType.ES6_MODULE */ :
		var scriptNode pkg_rhino.Node = originalMetadata.RootNode()
		_ = scriptNode
		if !((scriptNode == 0) || scriptNode.IsScript()) {
			j2gort.CheckStateMsg(false, j2gort.BoxNode(scriptNode))
		}
		return ScopedName_Of(name, func() pkg_rhino.Node {
			if scriptNode != 0 {
				return scriptNode.GetOnlyChild()
			}
			return pkg_rhino.Node(0)
		}())
	case int32(1) /* ModuleType.GOOG_MODULE */, int32(4) /* ModuleType.GOOG_PROVIDE */, int32(2) /* ModuleType.LEGACY_GOOG_MODULE */ :
		var closureModuleObject ScopedName = this_.priv_ModuleImportResolver_getScopedNameForClosureNamespace_3892(binding.ClosureNamespace())
		_ = closureModuleObject
		if closureModuleObject == nil || closureModuleObject.IsJavaNil_() {
			return ScopedName_Of(binding.ClosureNamespace(), pkg_rhino.Node(0))
		}
		if binding.IsModuleNamespace() {
			return closureModuleObject
		}
		return ScopedName_Of((j2gort.StrValString(closureModuleObject.GetScopedNameBase().GetName()) + "." + j2gort.StrValString(binding.OriginatingExport().ExportName())), closureModuleObject.GetScopedNameBase().GetScopeRoot())
	case int32(5) /* ModuleType.SCRIPT */ :
		return ScopedName_Of(name, pkg_rhino.Node(0))
	case int32(3) /* ModuleType.COMMON_JS */ :
		panic(j2gort.AsIface[javalang.IIllegalStateException](j2gort.ConstructShimThrowable("java.lang.IllegalStateException", "Typechecking CommonJS modules is not supported")))
	}
	panic("unreachable switch expression")
}

// Java static method: ModuleImportResolver.getModuleFromScopeRoot(com.google.javascript.jscomp.modules.ModuleMap moduleMap, com.google.javascript.jscomp.CompilerInputProvider inputProvider, com.google.javascript.rhino.Node moduleBody)
func ModuleImportResolver_GetModuleFromScopeRoot(moduleMap *ModuleMap, inputProvider CompilerInputProvider, moduleBody pkg_rhino.Node) *Module {
	_ = moduleMap
	_ = inputProvider
	_ = moduleBody
	if ModuleImportResolver_IsGoogModuleBody(moduleBody) {
		var googModuleCall pkg_rhino.Node = moduleBody.GetFirstChild()
		_ = googModuleCall
		var namespace string = googModuleCall.GetFirstChild().GetSecondChild().GetString()
		_ = namespace
		return moduleMap.GetClosureModule(namespace)
	} else if moduleBody.IsModuleBody() {
		var scriptNode pkg_rhino.Node = moduleBody.GetParent()
		_ = scriptNode
		var input *CompilerInput = j2gort.CheckNotNullPtr(inputProvider.GetInput(scriptNode.GetInputId()))
		_ = input
		var module *Module = moduleMap.GetModule_1_ModuleLoader_ModulePath(input.GetPath())
		_ = module
		if !(module.Metadata().GetModuleMetadataMap_ModuleMetadataBase().IsEs6Module()) {
			j2gort.CheckStateStr(false, "Typechecking of non-goog- and non-es-modules not supported")
		}
		return module
	}
	return nil
}

// Java static method: ModuleImportResolver.isGoogModuleBody(com.google.javascript.rhino.Node moduleBody)
func ModuleImportResolver_IsGoogModuleBody(moduleBody pkg_rhino.Node) bool {
	_ = moduleBody
	if moduleBody.IsModuleBody() {
		return moduleBody.GetParent().GetBooleanProp(pkg_rhino.StaticGet_Node_GOOG_MODULE())
	} else if moduleBody.IsBlock() {
		return (moduleBody.GetParent().IsFunction() && NodeUtil_IsBundledGoogModuleCall(moduleBody.GetGrandparent()))
	}
	return false
}
