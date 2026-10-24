// Port of com/google/javascript/jscomp/ModuleImportResolver.java @ closure-compiler v20261021

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

package jscomp

import (
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// ModuleImportResolver resolves module requires into TypedVars.
//
// Currently this only supports goog.modules, but can be extended for ES
// modules.
type ModuleImportResolver struct {
	moduleMap         *ModuleMap
	nodeToScopeMapper func(rhino.Node) *TypedScope
	registry          *rhino.JSTypeRegistry
}

const moduleImportResolver_GOOG = "goog"

var moduleImportResolver_GOOG_DEPENDENCY_CALLS = jutil.NewLinkedHashSet(
	"require", "requireType", "forwardDeclare", "requireDynamic")

var moduleImportResolver_GOOG_MODULE_GET = rhino.QualifiedName_Of("goog.module.get")

// Java: ModuleImportResolver(ModuleMap, Function<Node, TypedScope>, JSTypeRegistry)
func NewModuleImportResolver(
	moduleMap *ModuleMap, nodeToScopeMapper func(rhino.Node) *TypedScope, registry *rhino.JSTypeRegistry) *ModuleImportResolver {
	r := &ModuleImportResolver{}
	r.moduleMap = moduleMap
	r.nodeToScopeMapper = nodeToScopeMapper
	r.registry = registry
	return r
}

// ModuleImportResolver_IsGoogModuleDependencyCall returns whether this is a
// CALL node for goog.require, goog.requireType, goog.requireDynamic,
// goog.forwardDeclare, or goog.module.get.
//
// This method does not verify that the call is actually in a valid
// location. For example, this method does not verify that goog.require
// calls are at the top-level. That is left to the caller.
//
// Java: ModuleImportResolver.isGoogModuleDependencyCall(Node)
func ModuleImportResolver_IsGoogModuleDependencyCall(value rhino.Node) bool {
	if value.IsNil() ||
		!value.IsCall() ||
		!value.HasTwoChildren() ||
		!value.GetSecondChild().IsStringLit() {
		return false
	}
	callee := value.GetFirstChild()
	if !callee.IsGetProp() {
		return false
	}
	owner := callee.GetFirstChild()
	return (owner.IsName() &&
		owner.GetString() == moduleImportResolver_GOOG &&
		moduleImportResolver_GOOG_DEPENDENCY_CALLS.Contains(callee.GetString())) ||
		moduleImportResolver_GOOG_MODULE_GET.Matches(callee)
}

// GetClosureNamespaceTypeFromCall attempts to look up the type of a Closure
// namespace from a require call
//
// This returns null if the given ModuleMap is null, if the required module
// does not exist, or if support is missing for the type of required Module.
// Currently only requires of goog.modules, goog.provides, and ES module with
// goog.declareModuleId are supported.
//
// googRequire is a CALL node representing some kind of Closure require.
//
// Java: ModuleImportResolver.getClosureNamespaceTypeFromCall(Node)
func (r *ModuleImportResolver) GetClosureNamespaceTypeFromCall(googRequire rhino.Node) ScopedName {
	if r.moduleMap == nil {
		// TODO(b/124919359): make sure all tests have generated a ModuleMap
		return nil
	}
	moduleId := googRequire.GetSecondChild().GetString()
	return r.GetScopedNameForClosureNamespace(moduleId)
}

// GetScopedNameForClosureNamespace attempts to look up the type of a
// Closure namespace from a require call
//
// This returns null if the given ModuleMap is null, if the required module
// does not exist, or if support is missing for the type of required Module.
// Currently only requires of goog.modules, goog.provides, and ES module with
// goog.declareModuleId are supported.
//
// moduleId is a Closure namespace, such as "foo.bar"
//
// Java: ModuleImportResolver.getScopedNameForClosureNamespace(String)
func (r *ModuleImportResolver) GetScopedNameForClosureNamespace(moduleId string) ScopedName {
	module := r.moduleMap.GetClosureModule(moduleId)
	if module == nil {
		return nil
	}
	switch module.Metadata().ModuleType() {
	case ModuleMetadataMap_ModuleType_GOOG_PROVIDE:
		// Expect this to be a global variable
		provide := module.Metadata().RootNode()
		if !provide.IsNil() && provide.IsScript() {
			return ScopedName_Of(moduleId, provide.GetGrandparent())
		} else {
			// Unknown module requires default to 'goog provides', but we don't want to type them.
			return nil
		}
	case ModuleMetadataMap_ModuleType_GOOG_MODULE, ModuleMetadataMap_ModuleType_LEGACY_GOOG_MODULE:
		// TODO(b/124919359): Fix getGoogModuleScopeRoot to never return null.
		scopeRoot := r.GetGoogModuleScopeRoot(module)
		if !scopeRoot.IsNil() {
			return ScopedName_Of("exports", scopeRoot)
		}
		return nil
	case ModuleMetadataMap_ModuleType_ES6_MODULE:
		moduleBody := module.Metadata().RootNode().GetFirstChild() // SCRIPT -> MODULE_BODY
		return ScopedName_Of(Export_NAMESPACE, moduleBody)
	case ModuleMetadataMap_ModuleType_COMMON_JS:
		panic(jutil.NewIllegalStateExceptionMsg("Type checking CommonJs modules not yet supported"))
	case ModuleMetadataMap_ModuleType_SCRIPT:
		panic(jutil.NewIllegalStateExceptionMsg("Cannot import a name from a SCRIPT"))
	}
	panic(jutil.NewAssertionErrorNoMessage())
}

// GetGoogModuleScopeRoot returns the corresponding scope root Node from a
// goog.module.
//
// Java: ModuleImportResolver.getGoogModuleScopeRoot(Module)
func (r *ModuleImportResolver) GetGoogModuleScopeRoot(module *Module) rhino.Node {
	jutil.CheckArgumentMsg1(module.Metadata().IsGoogModule(), "%s", module.Metadata())
	scriptNode := module.Metadata().RootNode()

	if scriptNode.IsScript() &&
		scriptNode.HasOneChild() &&
		scriptNode.GetOnlyChild().IsModuleBody() {
		// The module root node should be a SCRIPT, whose first child is a MODULE_BODY.
		// The map is keyed off a MODULE_BODY node for a goog.module,
		// which is the only child of our SCRIPT node.
		return scriptNode.GetOnlyChild()
	} else if scriptNode.IsCall() {
		// This is a goog.loadModule call, and the scope is keyed off the FUNCTION node's BLOCK in:
		//   goog.loadModule(function(exports)
		functionLiteral := scriptNode.GetSecondChild()
		return NodeUtil_GetFunctionBody(functionLiteral)
	}
	// TODO(b/124919359): this case should not happen, but is triggering on goog.require calls in
	// rewritten modules with preserveClosurePrimitives enabled.
	return rhino.NilNode
}

// DeclareEsModuleImports declares/updates the type of all bindings imported
// into the ES module scope
//
// Returns a map from local nodes to ScopedNames for which nodeToScopeMapper
// couldn't find a scope, despite the original module existing. This is
// expected to happen for circular references if not all module scopes are
// created and the caller should handle declaring these names later, e.g. in
// TypedScopeCreator.
//
// Java: ModuleImportResolver.declareEsModuleImports(Module, TypedScope, CompilerInput)
func (r *ModuleImportResolver) DeclareEsModuleImports(
	module *Module, scope *TypedScope, moduleInput *CompilerInput) *jutil.LinkedHashMap[rhino.Node, ScopedName] {
	jutil.CheckArgumentMsg1(module.Metadata().IsEs6Module(), "%s", module)
	jutil.CheckArgumentMsg1(scope.IsModuleScope(), "%s", scope)
	missingNames := jutil.NewImmutableMapBuilder[rhino.Node, ScopedName]()
	for localName, binding := range module.BoundNames().All() {
		if !binding.IsCreatedByEsImport() {
			continue
		}
		// ES imports fall into two categories:
		//  - namespace imports. These correspond to an object type containing all named exports.
		//  - named imports. These always correspond, eventually, to a name local to a module.
		//    Note that we include imports of an `export default` in this case and map them to a
		//    pseudo-variable named *default*.
		export := r.GetScopedNameFromEsBinding(binding)
		modScope := r.nodeToScopeMapper(export.GetScopeRoot())
		if modScope == nil {
			jutil.CheckStateMsg1(binding.SourceNode().GetString() == localName, "%s", binding.SourceNode())
			missingNames.Put(binding.SourceNode(), export)
			continue
		}

		originalVar := modScope.GetVar(export.GetName())
		importType := originalVar.GetType()
		scope.Declare(
			localName,
			binding.SourceNode(),
			importType,
			moduleInput,
			/* inferred= */ originalVar.IsTypeInferred())

		// Non-namespace imports may be typedefs; if so, propagate the typedef prop onto the
		// export and import bindings, if not already there.
		if !binding.IsModuleNamespace() && binding.SourceNode().GetTypedefTypeProp() == nil {
			typedefType := originalVar.GetNameNode().GetTypedefTypeProp()
			if typedefType != nil {
				binding.SourceNode().SetTypedefTypeProp(typedefType)
				r.registry.DeclareType(scope, localName, typedefType)
			}
		}
	}
	return missingNames.BuildOrThrow()
}

// UpdateEsModuleNamespaceType declares or updates the type of properties
// representing exported names from ES module
//
// When the given object type does not have existing properties
// corresponding to exported names, this method adds new properties to the
// object type. If the object type already has properties, this method will
// ignore declared properties and update the type of inferred properties.
//
// The additional properties will be inferred (instead of declared) if and
// only if TypedVar#isTypeInferred() is true for the original exported name.
//
// We create this type to support 'import *' and goog.requires of this
// module. Note: we could lazily initialize this type if always creating it
// hurts performance.
//
// namespace is an object type which may already have properties
// representing exported names; scope is the scope rooted at the given
// module.
//
// Java: ModuleImportResolver.updateEsModuleNamespaceType(ObjectType, Module, TypedScope)
func (r *ModuleImportResolver) UpdateEsModuleNamespaceType(namespace rhino.ObjectType, module *Module, scope *TypedScope) {
	jutil.CheckArgumentMsg1(module.Metadata().IsEs6Module(), "%s", module)
	jutil.CheckArgumentMsg1(scope.IsModuleScope(), "%s", scope)

	for exportKey, binding := range module.Namespace().All() {
		if namespace.IsPropertyTypeDeclared(exportKey) {
			// Cannot change the type of a declared property after it is added to the ObjectType.
			continue
		}

		bindingSourceNode := binding.SourceNode() // e.g. 'x' in `export let x;` or `export {x};`
		export := r.GetScopedNameFromEsBinding(binding)
		var originalScope *TypedScope
		if export.GetScopeRoot() == scope.GetRootNode() {
			originalScope = scope
		} else {
			originalScope = r.nodeToScopeMapper(export.GetScopeRoot())
		}
		if originalScope == nil {
			// Exporting an import from an invalid module load or early reference.
			namespace.DefineInferredProperty(
				exportKey, r.registry.GetNativeType(rhino.JSTypeNative_UNKNOWN_TYPE), bindingSourceNode)
			r.UpdateAstForExport(
				bindingSourceNode,
				r.registry.GetNativeType(rhino.JSTypeNative_UNKNOWN_TYPE),
				/* typedefType= */ nil)
			continue
		}

		originalName, _ := originalScope.GetSlot(export.GetName()).(*TypedVar)
		exportType := originalName.GetType()
		if exportType == nil {
			exportType = r.registry.GetNativeType(rhino.JSTypeNative_NO_TYPE)
		}
		if originalName.IsTypeInferred() {
			// NB: this method may be either adding a new inferred property or updating the type of an
			// existing inferred property.
			namespace.DefineInferredProperty(exportKey, exportType, bindingSourceNode)
		} else {
			namespace.DefineDeclaredProperty(exportKey, exportType, bindingSourceNode)
		}

		r.UpdateAstForExport(
			bindingSourceNode, exportType, originalName.GetNameNode().GetTypedefTypeProp())
	}
}

// Java: ModuleImportResolver.updateAstForExport(Node, JSType, JSType)
func (r *ModuleImportResolver) UpdateAstForExport(
	bindingSourceNode rhino.Node, exportType rhino.JSType, typedefType rhino.JSType) {
	bindingSourceNode.SetJSType(exportType)
	bindingSourceNode.SetTypedefTypeProp(typedefType)
	if bindingSourceNode.GetParent().IsExportSpec() {
		// given  `export {x as y} from './a/b'` update the y NAME node
		bindingSourceNode.GetNext().SetJSType(exportType)
	}
}

// GetScopedNameFromEsBinding: given a Binding from an ES module, return the
// name and scope of the bound name.
//
// Java: ModuleImportResolver.getScopedNameFromEsBinding(Binding)
func (r *ModuleImportResolver) GetScopedNameFromEsBinding(binding *Binding) ScopedName {
	// NB: If the original export was an `export default` then the local name is *default*.
	// We've already declared a dummy variable named `*default*` in the scope.
	//
	// Go-only: language-adapter — Java's String name is null when the originating export has no
	// local name (a goog.module export, or a name "exported" by a script); see
	// moduleImportResolver_scopedNameOf.
	var name *string
	if binding.IsModuleNamespace() {
		namespace := Export_NAMESPACE
		name = &namespace
	} else {
		name = binding.BoundName()
	}
	var originalMetadata *ModuleMetadataMap_ModuleMetadata
	if binding.IsModuleNamespace() {
		originalMetadata = binding.Metadata()
	} else {
		originalMetadata = binding.OriginatingExport().ModuleMetadata()
	}
	switch originalMetadata.ModuleType() {
	case ModuleMetadataMap_ModuleType_ES6_MODULE:
		scriptNode := originalMetadata.RootNode()
		// Imports of nonexistent modules have a null 'root node'. Imports of names from scripts are
		// meaningless.
		jutil.CheckStateMsg1(scriptNode.IsNil() || scriptNode.IsScript(), "%s", scriptNode)
		var scopeRoot rhino.Node
		if !scriptNode.IsNil() {
			scopeRoot = scriptNode.GetOnlyChild()
		}
		return moduleImportResolver_scopedNameOf(name, scopeRoot)
	case ModuleMetadataMap_ModuleType_GOOG_MODULE, ModuleMetadataMap_ModuleType_GOOG_PROVIDE, ModuleMetadataMap_ModuleType_LEGACY_GOOG_MODULE:
		closureModuleObject :=
			r.GetScopedNameForClosureNamespace(*binding.ClosureNamespace())
		if closureModuleObject == nil {
			// This is an error, but one that should be reported elsewhere. assume the module is
			// a goog.provide for legacy compatibility.
			return ScopedName_Of(*binding.ClosureNamespace(), rhino.NilNode)
		}
		if binding.IsModuleNamespace() {
			return closureModuleObject
		}
		return ScopedName_Of(
			closureModuleObject.GetName()+"."+modules_nullableString(binding.OriginatingExport().ExportName()),
			closureModuleObject.GetScopeRoot())
	case ModuleMetadataMap_ModuleType_SCRIPT:
		// Importing SCRIPTs should not allow you to look up names in scope.
		// we also don't really support CommonJs.
		return moduleImportResolver_scopedNameOf(name, rhino.NilNode)
	case ModuleMetadataMap_ModuleType_COMMON_JS:
		panic(jutil.NewIllegalStateExceptionMsg("Typechecking CommonJS modules is not supported"))
	}
	// Go-only: language-adapter — an exhaustive Java switch expression
	// throws MatchException for an enum constant it does not know.
	panic(jutil.NewIllegalStateException())
}

// moduleImportResolver_scopedNameOf is ScopedName.of(name, scopeRoot) for
// getScopedNameFromEsBinding's nullable name.
//
// Go-only: language-adapter — ScopedName's name is a Go string, and Java
// passes a null name here for a name imported from a script (the Export that
// NonEsModuleProcessor resolves has no local name). The scope root is then
// null as well, and nothing reads the name of such a ScopedName (its scope is
// looked up first and is missing), so the null name becomes "".
func moduleImportResolver_scopedNameOf(name *string, scopeRoot rhino.Node) ScopedName {
	if name == nil {
		return ScopedName_Of("", scopeRoot)
	}
	return ScopedName_Of(*name, scopeRoot)
}

// ModuleImportResolver_GetModuleFromScopeRoot returns the Module
// corresponding to this scope root, or null if not a module root.
//
// Java: ModuleImportResolver.getModuleFromScopeRoot(ModuleMap, CompilerInputProvider, Node)
func ModuleImportResolver_GetModuleFromScopeRoot(
	moduleMap *ModuleMap, inputProvider CompilerInputProvider, moduleBody rhino.Node) *Module {
	if ModuleImportResolver_IsGoogModuleBody(moduleBody) {
		googModuleCall := moduleBody.GetFirstChild()
		namespace := googModuleCall.GetFirstChild().GetSecondChild().GetString()
		return moduleMap.GetClosureModule(namespace)
	} else if moduleBody.IsModuleBody() {
		scriptNode := moduleBody.GetParent()
		input := jutil.CheckNotNull(inputProvider.GetInput(scriptNode.GetInputId()))
		module := moduleMap.GetModuleByPath(input.GetPath())
		// TODO(b/131418081): Also cover CommonJS modules.
		jutil.CheckStateMsg(
			module.Metadata().IsEs6Module(),
			"Typechecking of non-goog- and non-es-modules not supported")
		return module
	}
	return nil
}

// Java: ModuleImportResolver.isGoogModuleBody(Node)
func ModuleImportResolver_IsGoogModuleBody(moduleBody rhino.Node) bool {
	if moduleBody.IsModuleBody() {
		return moduleBody.GetParent().GetBooleanProp(rhino.Node_GOOG_MODULE)
	} else if moduleBody.IsBlock() {
		return moduleBody.GetParent().IsFunction() &&
			NodeUtil_IsBundledGoogModuleCall(moduleBody.GetGrandparent())
	}
	return false
}
