// Port of com/google/javascript/jscomp/FindModuleDependencies.java @ closure-compiler v20261015

/*
 * Copyright 2017 The Closure Compiler Authors.
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

// FindModuleDependencies finds and updates any direct dependencies of an
// input. Used to walk the dependency graph and support a strict depth-first
// dependency ordering. Marks an input as providing its module name.
//
// Discovers dependencies from:
//
//   - goog.require calls
//   - ES6 import statements
//   - CommonJS require statements
//   - goog.requireDynamic calls
//
// The order of dependency references is preserved so that a deterministic
// depth-first ordering can be achieved.
type FindModuleDependencies struct {
	compiler                AbstractCompiler
	supportsEs6Modules      bool
	supportsCommonJsModules bool
	moduleType              CompilerInput_ModuleType
	dynamicImportScope      *Scope
	inputPathByWebpackId    *jutil.LinkedHashMap[string, string]
}

var findModuleDependencies_GOOG_MODULE = rhino.QualifiedName_Of("goog.module")
var findModuleDependencies_GOOG_PROVIDE = rhino.QualifiedName_Of("goog.provide")

// Java: FindModuleDependencies(AbstractCompiler, boolean, boolean, ImmutableMap<String, String>)
func NewFindModuleDependencies(
	compiler AbstractCompiler,
	supportsEs6Modules bool,
	supportsCommonJsModules bool,
	inputPathByWebpackId *jutil.LinkedHashMap[string, string]) *FindModuleDependencies {
	f := &FindModuleDependencies{}
	f.moduleType = CompilerInput_ModuleType_NONE
	f.dynamicImportScope = nil
	f.compiler = compiler
	f.supportsEs6Modules = supportsEs6Modules
	f.supportsCommonJsModules = supportsCommonJsModules
	f.inputPathByWebpackId = inputPathByWebpackId
	return f
}

// FindModuleDependencies_FindGoogProvideOrGoogModule is the private static
// nested class FindModuleDependencies.FindGoogProvideOrGoogModule.
type FindModuleDependencies_FindGoogProvideOrGoogModule struct {
	NodeTraversal_AbstractPreOrderCallbackBase

	found bool
}

// Java: FindModuleDependencies.FindGoogProvideOrGoogModule()
func NewFindModuleDependencies_FindGoogProvideOrGoogModule() *FindModuleDependencies_FindGoogProvideOrGoogModule {
	return &FindModuleDependencies_FindGoogProvideOrGoogModule{}
}

// Java: FindModuleDependencies.FindGoogProvideOrGoogModule.isFound()
func (f *FindModuleDependencies_FindGoogProvideOrGoogModule) IsFound() bool {
	return f.found
}

// Java: FindModuleDependencies.FindGoogProvideOrGoogModule.shouldTraverse(NodeTraversal, Node, Node)
func (f *FindModuleDependencies_FindGoogProvideOrGoogModule) ShouldTraverse(nodeTraversal *NodeTraversal, n *rhino.Node, parent *rhino.Node) bool {
	if f.found {
		return false
	}
	// Shallow traversal, since we don't need to inspect within functions or expressions.
	if NodeUtil_IsShallowStatementTree(parent) {
		if n.IsExprResult() {
			maybeGetProp := n.GetFirstFirstChild()
			if maybeGetProp != nil &&
				(findModuleDependencies_GOOG_PROVIDE.Matches(maybeGetProp) || findModuleDependencies_GOOG_MODULE.Matches(maybeGetProp)) {
				f.found = true
				return false
			}
		}
		return true
	}
	return false
}

// Java: FindModuleDependencies.process(Node)
func (f *FindModuleDependencies) Process(root *rhino.Node) {
	jutil.CheckArgument(root.IsScript())
	if Es6RewriteModules_IsEs6ModuleRoot(root) {
		f.moduleType = CompilerInput_ModuleType_ES6
	}
	input := f.compiler.GetInput(root.GetInputId())

	// The "goog" namespace isn't always specifically required.
	// The deps parser will pick up any access to a `goog.foo()` call
	// and add "goog" as a dependency. If "goog" is a dependency of the
	// file we add it here to the ordered requires so that it's always
	// first.
	if findModuleDependencies_containsRequire(input.GetRequires(), DependencyInfo_Require_BASE) {
		input.AddOrderedRequire(DependencyInfo_Require_BASE)
	}

	sc := NewSyntacticScopeCreator(f.compiler)
	b := NodeTraversal_Builder{
		compiler:     f.compiler,
		callback:     f,
		scopeCreator: sc,
	}
	b.Traverse(root)
	f.dynamicImportScope = nil

	if f.moduleType == CompilerInput_ModuleType_ES6 {
		f.ConvertToEs6Module(root, true)
	} else if f.moduleType == CompilerInput_ModuleType_NONE &&
		f.inputPathByWebpackId != nil &&
		f.inputPathByWebpackId.ContainsValue(input.GetPath().String()) {
		f.moduleType = CompilerInput_ModuleType_IMPORTED_SCRIPT
	}

	input.AddProvide(input.GetPath().ToModuleName())
	input.SetJsModuleType(f.moduleType)
	input.SetHasFullParseDependencyInfo(true)
}

// findModuleDependencies_containsRequire is
// ImmutableList.contains(Object) on the requires (Require is an AutoValue
// compared with equals).
//
// Go-only: language-adapter — List.contains on a Go slice.
func findModuleDependencies_containsRequire(requires []*DependencyInfo_Require, require *DependencyInfo_Require) bool {
	for _, r := range requires {
		if jutil.ObjectsEquals(r, require) {
			return true
		}
	}
	return false
}

// Java: FindModuleDependencies.shouldTraverse(NodeTraversal, Node, Node)
func (f *FindModuleDependencies) ShouldTraverse(t *NodeTraversal, n *rhino.Node, parent *rhino.Node) bool {
	if f.supportsCommonJsModules &&
		n.IsFunction() &&
		ProcessCommonJSModules_IsCommonJsDynamicImportCallback(
			n, f.compiler.GetOptions().GetModuleResolutionMode()) {
		if f.dynamicImportScope == nil {
			f.dynamicImportScope = t.GetScope()
		}
	}

	return true
}

// Java: FindModuleDependencies.visit(NodeTraversal, Node, Node)
func (f *FindModuleDependencies) Visit(t *NodeTraversal, n *rhino.Node, parent *rhino.Node) {
	resolutionMode := f.compiler.GetOptions().GetModuleResolutionMode()
	if NodeUtil_IsShallowStatementTree(parent) {
		if n.IsExprResult() {
			maybeGetProp := n.GetFirstFirstChild()
			if maybeGetProp != nil &&
				(findModuleDependencies_GOOG_PROVIDE.Matches(maybeGetProp) || findModuleDependencies_GOOG_MODULE.Matches(maybeGetProp)) {
				f.moduleType = CompilerInput_ModuleType_GOOG
				return
			}
		}
	}

	// goog.requireDynamic()
	if NodeUtil_IsGoogRequireDynamicCall(n) {
		t.GetInput().AddRequireDynamicImports(n.GetLastChild().GetString())
	}

	if f.supportsEs6Modules && n.IsExport() {
		f.moduleType = CompilerInput_ModuleType_ES6
		if n.GetBooleanProp(rhino.Node_EXPORT_DEFAULT) {
			// export default
		} else if n.HasTwoChildren() {
			// export * from 'moduleIdentifier';
			// export {x, y as z} from 'moduleIdentifier';
			f.AddEs6ModuleImportToGraph(t, n)
		}
	} else if f.supportsEs6Modules && n.IsImport() {
		f.moduleType = CompilerInput_ModuleType_ES6
		f.AddEs6ModuleImportToGraph(t, n)
	} else if f.supportsEs6Modules &&
		n.GetToken() == rhino.Token_DYNAMIC_IMPORT &&
		n.GetFirstChild().IsString() {
		path := n.GetFirstChild().GetString()
		modulePath := t.GetInput().GetPath().ResolveJsModuleSilently(path)
		if modulePath != nil {
			t.GetInput().AddDynamicRequire(modulePath.ToModuleName())
		}
	} else if f.supportsCommonJsModules {
		if f.moduleType != CompilerInput_ModuleType_GOOG &&
			ProcessCommonJSModules_IsCommonJsExport(t, n, resolutionMode) {
			f.moduleType = CompilerInput_ModuleType_COMMONJS
		} else if ProcessCommonJSModules_IsCommonJsImport(n, resolutionMode) {
			path := ProcessCommonJSModules_GetCommonJsImportPath(n, resolutionMode)

			modulePath :=
				t.GetInput().
					GetPath().
					ResolveJsModule(path, processCommonJSModules_sourceFileName(n), n.GetLineno(), n.GetCharno())

			if modulePath != nil {
				if f.dynamicImportScope != nil ||
					ProcessCommonJSModules_IsCommonJsDynamicImportCallback(
						NodeUtil_GetEnclosingFunction(n), resolutionMode) {
					t.GetInput().AddDynamicRequire(modulePath.ToModuleName())
				} else {
					t.GetInput().AddOrderedRequire(DependencyInfo_Require_CommonJs(modulePath.ToModuleName(), path))
				}
			}
		}

		// TODO(ChadKillingsworth) add require.ensure support
	}

	if parent != nil &&
		(parent.IsExprResult() || !t.InGlobalHoistScope()) &&
		NodeUtil_IsGoogRequireCall(n) &&
		n.GetSecondChild() != nil &&
		n.GetSecondChild().IsStringLit() {
		namespace := n.GetSecondChild().GetString()
		if strings.HasPrefix(namespace, "goog.") {
			t.GetInput().AddOrderedRequire(DependencyInfo_Require_BASE)
		}
		t.GetInput().AddOrderedRequire(DependencyInfo_Require_GoogRequireSymbol(namespace))
	}
}

// Java: FindModuleDependencies.enterScope(NodeTraversal)
func (f *FindModuleDependencies) EnterScope(t *NodeTraversal) {}

// Java: FindModuleDependencies.exitScope(NodeTraversal)
func (f *FindModuleDependencies) ExitScope(t *NodeTraversal) {
	if f.dynamicImportScope != nil && t.GetScopeRoot() == f.dynamicImportScope.GetRootNode() && t.GetScope() == f.dynamicImportScope {
		f.dynamicImportScope = nil
	}
}

// AddEs6ModuleImportToGraph adds an es6 module from an import node (import
// or export statement) to the graph.
//
// Java: FindModuleDependencies.addEs6ModuleImportToGraph(NodeTraversal, Node)
func (f *FindModuleDependencies) AddEs6ModuleImportToGraph(t *NodeTraversal, n *rhino.Node) {
	moduleName := f.GetEs6ModuleNameFromImportNode(t, n)
	if strings.HasPrefix(moduleName, "goog.") {
		t.GetInput().AddOrderedRequire(DependencyInfo_Require_BASE)
	}
	t.GetInput().AddOrderedRequire(DependencyInfo_Require_Es6Import(moduleName, n.GetLastChild().GetString()))
}

// GetEs6ModuleNameFromImportNode gets the module name from an import node
// (import or export statement).
//
// Java: FindModuleDependencies.getEs6ModuleNameFromImportNode(NodeTraversal, Node)
func (f *FindModuleDependencies) GetEs6ModuleNameFromImportNode(t *NodeTraversal, n *rhino.Node) string {
	importName := n.GetLastChild().GetString()
	isNamespaceImport := strings.HasPrefix(importName, "goog:")
	if isNamespaceImport {
		// Allow importing Closure namespace objects (e.g. from goog.provide or goog.module) as
		//   import ... from 'goog:my.ns.Object'.
		// These are rewritten to plain namespace object accesses.
		return jutil.StringSubstringFrom(importName, jutil.StringLength("goog:"))
	} else {
		modulePath :=
			t.GetInput().
				GetPath().
				ResolveJsModule(importName, processCommonJSModules_sourceFileName(n), n.GetLineno(), n.GetCharno())
		if modulePath == nil {
			// The module loader issues an error
			// Fall back to assuming the module is a file path
			modulePath = t.GetInput().GetPath().ResolveModuleAsPath(importName)
		}
		return modulePath.ToModuleName()
	}
}

// Java: FindModuleDependencies.convertToEs6Module(Node, boolean)
func (f *FindModuleDependencies) ConvertToEs6Module(root *rhino.Node, skipGoogProvideModuleCheck bool) bool {
	if Es6RewriteModules_IsEs6ModuleRoot(root) {
		return true
	}
	if !skipGoogProvideModuleCheck {
		finder := NewFindModuleDependencies_FindGoogProvideOrGoogModule()
		NodeTraversal_Traverse(f.compiler, root, finder)
		if finder.IsFound() {
			return false
		}
	}
	moduleNode := rhino.NewNode(rhino.Token_MODULE_BODY).Srcref(root)
	moduleNode.AddChildrenToBack(root.RemoveChildren())
	root.AddChildToBack(moduleNode)
	return true
}
