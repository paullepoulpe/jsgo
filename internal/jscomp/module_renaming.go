// Port of com/google/javascript/jscomp/ModuleRenaming.java @ closure-compiler v20261021

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

package jscomp

import (
	"strings"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// ModuleRenaming is the centralized location for determining how to rename
// modules. (A final class with a private constructor and only static
// members: its members are package functions.)

// ModuleRenaming_DEFAULT_EXPORT_VAR_PREFIX is the name of the temporary
// variable created for the default export before globalization of the
// module.
const ModuleRenaming_DEFAULT_EXPORT_VAR_PREFIX = "$jscompDefaultExport"

// ModuleRenaming_GetGlobalNameOfEsModuleLocalVariable returns the global
// name of a variable declared in an ES module.
//
// Java: ModuleRenaming.getGlobalNameOfEsModuleLocalVariable(ModuleMetadata, String)
func ModuleRenaming_GetGlobalNameOfEsModuleLocalVariable(
	moduleMetadata *ModuleMetadataMap_ModuleMetadata, variableName string) rhino.QualifiedName {
	return rhino.QualifiedName_Of(
		variableName + "$$" + ModuleRenaming_GetGlobalNameWithModuleMetadata(moduleMetadata /* googNamespace= */, nil).Join())
}

// ModuleRenaming_GetGlobalNameOfAnonymousDefaultExport returns the global
// name of the anonymous default export for the given module.
//
// Java: ModuleRenaming.getGlobalNameOfAnonymousDefaultExport(ModuleMetadata)
func ModuleRenaming_GetGlobalNameOfAnonymousDefaultExport(moduleMetadata *ModuleMetadataMap_ModuleMetadata) rhino.QualifiedName {
	return ModuleRenaming_GetGlobalNameOfEsModuleLocalVariable(moduleMetadata, ModuleRenaming_DEFAULT_EXPORT_VAR_PREFIX)
}

// ModuleRenaming_GetGlobalNameWithModuleMetadata: moduleMetadata is the
// metadata of the module to get the global name of; googNamespace the
// Closure namespace that is being referenced fromEsModule this module, if
// any. Returns the global, qualified name to rewrite any references to this
// module to.
//
// Java: ModuleRenaming.getGlobalName(ModuleMetadata, String)
func ModuleRenaming_GetGlobalNameWithModuleMetadata(
	moduleMetadata *ModuleMetadataMap_ModuleMetadata, googNamespace *string) rhino.QualifiedName {
	return ModuleRenaming_GlobalizedModuleName_CreateWithModuleMetadata(moduleMetadata, googNamespace, nil).AliasName()
}

// ModuleRenaming_GetGlobalName returns the post-transpilation, globalized
// name of the export.
//
// Java: ModuleRenaming.getGlobalName(Export)
func ModuleRenaming_GetGlobalName(export *Export) rhino.QualifiedName {
	if export.ModuleMetadata().IsEs6Module() {
		if *export.LocalName() == Export_DEFAULT_EXPORT_NAME {
			return ModuleRenaming_GetGlobalNameOfAnonymousDefaultExport(export.ModuleMetadata())
		}
		return ModuleRenaming_GetGlobalNameOfEsModuleLocalVariable(export.ModuleMetadata(), *export.LocalName())
	}
	return ModuleRenaming_GetGlobalNameWithModuleMetadata(export.ModuleMetadata(), export.ClosureNamespace()).
		Getprop(*export.ExportName())
}

// ModuleRenaming_GetGlobalNameWithBinding returns the post-transpilation,
// globalized name of the binding.
//
// Java: ModuleRenaming.getGlobalName(Binding)
func ModuleRenaming_GetGlobalNameWithBinding(binding *Binding) rhino.QualifiedName {
	if binding.IsModuleNamespace() {
		return ModuleRenaming_GetGlobalNameWithModuleMetadata(binding.Metadata(), binding.ClosureNamespace())
	}
	return ModuleRenaming_GetGlobalName(binding.OriginatingExport())
}

// Java: ModuleRenaming.getNameRootType(String, TypedScope)
func ModuleRenaming_GetNameRootType(qname string, globalTypedScope *TypedScope) rhino.JSType {
	if globalTypedScope == nil {
		return nil
	}
	root := NodeUtil_GetRootOfQualifiedNameString(qname)
	return jutil.CheckNotNullMsg1(globalTypedScope.GetVar(root), "missing var for %s", root).GetType()
}

// ModuleRenaming_GetGlobalNameForJsDoc returns the globalized name of a
// reference to a binding in JS Doc.
//
// For example:
//
//	// bar
//	export class Bar {}
//
//	// foo
//	import * as bar from 'bar';
//	export {bar};
//
//	import * as foo from 'foo';
//	let /** !foo.bar.Bar *\/ b;
//
// Should call this method with the binding for foo and a list ("bar",
// "Bar"). In this example any of these properties could also be modules.
// This method will replace as much as the GETPROP as it can with module
// exported variables. Meaning in the above example this would return
// something like "baz$$module$bar", whereas if this method were called for
// just "foo.bar" it would return "module$bar", as it refers to a module
// object itself.
//
// Java: ModuleRenaming.getGlobalNameForJsDoc(ModuleMap, Binding, List<String>)
func ModuleRenaming_GetGlobalNameForJsDoc(
	moduleMap *ModuleMap, binding *Binding, propertyChain []string) string {
	prop := 0

	for binding.IsModuleNamespace() &&
		binding.Metadata().IsEs6Module() &&
		prop < len(propertyChain) {
		propertyName := propertyChain[prop]
		m := moduleMap.GetModuleByPath(binding.Metadata().Path())
		if m.Namespace().ContainsKey(propertyName) {
			binding = m.Namespace().Get(propertyName)
		} else {
			// This means someone referenced an invalid export on a module object. This should be an
			// error, so just rewrite and let the type checker complain later. It isn't a super clear
			// error, but we're working on type checking modules soon.
			break
		}
		prop++
	}

	globalName := ModuleRenaming_GetGlobalNameWithBinding(binding)

	for prop < len(propertyChain) {
		globalName = globalName.Getprop(propertyChain[prop])
		prop++
	}

	return globalName.Join()
}

// ModuleRenaming_GlobalizedModuleName stores a fully qualified globalized
// name of a module or provide + the type of its root node.
//
// For goog.provides and legacy goog.modules, the fully qualified name is
// just the namespace of the module "a.b.c.d". For ES modules, it's
// deterministically generated based on the module path, e.g.
// "module$path$to$my$dir". For non-legacy goog.modules, it's
// deterministically generated based on the module name, e.g.
// $module$exports$a$b$c$d".
//
// Used because in some cases, the root of the qualified name will not be
// present in any scopes until a subsequent run of ClosureRewriteModule. It's
// useful to store the type of Closure module exports
// `module$exports$my$goog$module` separately.
//
// rootNameType is the type of the root of `aliasName`, as it may not always
// exist in the scope yet.
//
// Java: record ModuleRenaming.GlobalizedModuleName(QualifiedName aliasName, JSType rootNameType)
type ModuleRenaming_GlobalizedModuleName struct {
	aliasName    rhino.QualifiedName
	rootNameType rhino.JSType
}

// Java: ModuleRenaming.GlobalizedModuleName(QualifiedName, JSType)
func NewModuleRenaming_GlobalizedModuleName(aliasName rhino.QualifiedName, rootNameType rhino.JSType) *ModuleRenaming_GlobalizedModuleName {
	jutil.CheckNotNullMsg1(aliasName, "%s", "aliasName") // requireNonNull(aliasName, "aliasName")
	return &ModuleRenaming_GlobalizedModuleName{aliasName: aliasName, rootNameType: rootNameType}
}

// AliasName is the record component accessor aliasName().
//
// Java: ModuleRenaming.GlobalizedModuleName.aliasName()
func (g *ModuleRenaming_GlobalizedModuleName) AliasName() rhino.QualifiedName { return g.aliasName }

// RootNameType is the record component accessor rootNameType().
//
// Java: ModuleRenaming.GlobalizedModuleName.rootNameType()
func (g *ModuleRenaming_GlobalizedModuleName) RootNameType() rhino.JSType { return g.rootNameType }

// ToQname creates a GETPROP chain with type information representing this
// name
//
// Java: ModuleRenaming.GlobalizedModuleName.toQname(AstFactory)
func (g *ModuleRenaming_GlobalizedModuleName) ToQname(astFactory *AstFactory) rhino.Node {
	rootName := astFactory.CreateName(g.AliasName().GetRoot(), AstFactory_Type_OfJSType(g.RootNameType()))
	if g.AliasName().IsSimple() {
		return rootName
	}
	return astFactory.CreateGetPropsWithoutColors(
		rootName, g.AliasName().Components()[1:])
}

// Getprop returns a copy of this name with the given property appended to
// the AliasName()
//
// Java: ModuleRenaming.GlobalizedModuleName.getprop(String)
func (g *ModuleRenaming_GlobalizedModuleName) Getprop(property string) *ModuleRenaming_GlobalizedModuleName {
	jutil.CheckArgument(property != "" && !strings.Contains(property, "."))
	return ModuleRenaming_GlobalizedModuleName_Create(g.AliasName().Getprop(property), g.RootNameType())
}

// Java: ModuleRenaming.GlobalizedModuleName.create(QualifiedName, JSType)
func ModuleRenaming_GlobalizedModuleName_Create(aliasName rhino.QualifiedName, rootNameType rhino.JSType) *ModuleRenaming_GlobalizedModuleName {
	return NewModuleRenaming_GlobalizedModuleName(aliasName, rootNameType)
}

// ModuleRenaming_GlobalizedModuleName_CreateWithModuleMetadata:
// moduleMetadata is the metadata of the module to get the global name of;
// googNamespace the Closure namespace that is being referenced fromEsModule
// this module, if any; globalTypedScope a global scope expected to contain
// the types of goog.provided names and rewritten ES6 module names. Returns
// the global, qualified name to rewrite any references to this module to,
// along with the type of the root of the module if globalTypedScope is not
// nil.
//
// Java: ModuleRenaming.GlobalizedModuleName.create(ModuleMetadata, String, TypedScope)
func ModuleRenaming_GlobalizedModuleName_CreateWithModuleMetadata(
	moduleMetadata *ModuleMetadataMap_ModuleMetadata,
	googNamespace *string,
	globalTypedScope *TypedScope) *ModuleRenaming_GlobalizedModuleName {
	jutil.CheckState(googNamespace == nil || moduleMetadata.GoogNamespaces().Contains(*googNamespace))
	switch moduleMetadata.ModuleType() {
	case ModuleMetadataMap_ModuleType_GOOG_MODULE:
		// The exported type is stored on the MODULE_BODY node.
		moduleBody := moduleMetadata.RootNode().GetFirstChild()
		return NewModuleRenaming_GlobalizedModuleName(
			rhino.QualifiedName_Of(ClosureRewriteModule_GetBinaryModuleNamespace(*googNamespace)),
			moduleBody.GetJSType())
	case ModuleMetadataMap_ModuleType_GOOG_PROVIDE, ModuleMetadataMap_ModuleType_LEGACY_GOOG_MODULE:
		return NewModuleRenaming_GlobalizedModuleName(
			rhino.QualifiedName_Of(*googNamespace), ModuleRenaming_GetNameRootType(*googNamespace, globalTypedScope))
	case ModuleMetadataMap_ModuleType_ES6_MODULE, ModuleMetadataMap_ModuleType_COMMON_JS:
		return NewModuleRenaming_GlobalizedModuleName(
			rhino.QualifiedName_Of(moduleMetadata.Path().ToModuleName()),
			ModuleRenaming_GetNameRootType(moduleMetadata.Path().ToModuleName(), globalTypedScope))
	case ModuleMetadataMap_ModuleType_SCRIPT:
	}
	panic(jutil.NewIllegalStateExceptionMsg("Unexpected module type: " + moduleMetadata.ModuleType().String()))
}
