// Port of com/google/javascript/jscomp/GatherModuleMetadata.java @ closure-compiler v20261021

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

	"github.com/paullepoulpe/jsgo/internal/ast"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/parser"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// Go-only: this pass uses the typed AST view (internal/ast) to read the
// calls it records: VisitGoogCall takes the ast.Call, finds the root of its
// callee's property chain, and reads "exactly one argument, a string
// literal" (goog.provide('a.b')) with gatherModuleMetadata_onlyStringArg
// instead of HasTwoChildren() && GetLastChild().IsStringLit(). The source
// of an import or `export ... from` is read by role. Kept on rhino: the
// toggle checks, which read the call's parent of any kind.

// GatherModuleMetadata gathers metadata around modules that is useful for
// checking imports / requires and creates a ModuleMetadataMap.
type GatherModuleMetadata struct {
	// Map from module path to module. These modules represent files and thus
	// will contain all goog namespaces that are in the file. These are not
	// the same modules in modulesByGoogNamespace.
	modulesByPath *jutil.LinkedHashMap[string, *ModuleMetadataMap_ModuleMetadata]

	// Map from Closure namespace to module. These modules represent just the
	// single namespace and thus each module has only one goog namespace in
	// its ModuleMetadata#googNamespaces(). These are not the same modules in
	// modulesByPath.
	modulesByGoogNamespace *jutil.LinkedHashMap[string, *ModuleMetadataMap_ModuleMetadata]

	// The current module being traversed.
	currentModule *GatherModuleMetadata_ModuleMetadataBuilder

	// The module currentModule is nested under, if any. Modules are expected
	// to be at most two deep (a script and then a goog.loadModule call).
	parentModule *GatherModuleMetadata_ModuleMetadataBuilder

	// The call to goog.loadModule we are traversing.
	loadModuleCall rhino.Node

	compiler               AbstractCompiler
	processCommonJsModules bool
	moduleResolutionMode   ModuleLoader_ResolutionMode
}

var GatherModuleMetadata_MIXED_MODULE_TYPE = DiagnosticType_Error("JSC_MIXED_MODULE_TYPE", "A file cannot be both {0} and {1}.")

var GatherModuleMetadata_INVALID_NAMESPACE_OR_MODULE_ID = DiagnosticType_Error(
	"JSC_INVALID_NAMESPACE_OR_MODULE_ID",
	"Namespace and module ID must be a dot-separated sequence of legal property"+
		" identifiers and must only contain ASCII, 0-9, $, ., and _. Found ''{0}''")

var GatherModuleMetadata_INVALID_DECLARE_MODULE_ID_CALL = DiagnosticType_Error(
	"JSC_INVALID_DECLARE_NAMESPACE_CALL",
	"goog.declareModuleId parameter must be a string literal.")

var GatherModuleMetadata_DECLARE_MODULE_ID_OUTSIDE_ES6_MODULE = DiagnosticType_Error(
	"JSC_DECLARE_MODULE_NAMESPACE_OUTSIDE_ES6_MODULE",
	"goog.declareModuleId can only be called within ES6 modules.")

var GatherModuleMetadata_MULTIPLE_DECLARE_MODULE_NAMESPACE = DiagnosticType_Error(
	"JSC_MULTIPLE_DECLARE_MODULE_NAMESPACE",
	"goog.declareModuleId can only be called once per ES6 module.")

var GatherModuleMetadata_INVALID_REQUIRE_TYPE = DiagnosticType_Error(
	"JSC_INVALID_REQUIRE_TYPE", "Argument to goog.requireType must be a string.")

var GatherModuleMetadata_INVALID_REQUIRE_DYNAMIC = DiagnosticType_Error(
	"JSC_INVALID_REQUIRE_DYNAMIC", "Argument to goog.requireDynamic must be a string.")

var GatherModuleMetadata_INVALID_MAYBE_REQUIRE = DiagnosticType_Error(
	"JSC_INVALID_MAYBE_REQUIRE",
	"Argument to goog.maybeRequireFrameworkInternalOnlyDoNotCallOrElse must be a string.")

var GatherModuleMetadata_INVALID_SET_TEST_ONLY = DiagnosticType_Error(
	"JSC_INVALID_SET_TEST_ONLY",
	"Optional, single argument to goog.setTestOnly must be a string.")

var GatherModuleMetadata_INVALID_NESTED_LOAD_MODULE = DiagnosticType_Error("JSC_INVALID_NESTED_LOAD_MODULE", "goog.loadModule cannot be nested.")

var GatherModuleMetadata_INVALID_TOGGLE_USAGE = DiagnosticType_Error("JSC_INVALID_TOGGLE_USAGE", "Invalid toggle usage: {0}")

var gatherModuleMetadata_GOOG_PROVIDE = rhino.IR_Getprop(rhino.IR_Name("goog"), "provide")

var gatherModuleMetadata_GOOG_MODULE = rhino.IR_Getprop(rhino.IR_Name("goog"), "module")

var gatherModuleMetadata_GOOG_MODULE_GET = rhino.IR_Getprop(rhino.IR_Getprop(rhino.IR_Name("goog"), "module"), "get")

var gatherModuleMetadata_GOOG_REQUIRE = rhino.IR_Getprop(rhino.IR_Name("goog"), "require")

var gatherModuleMetadata_GOOG_REQUIRE_TYPE = rhino.IR_Getprop(rhino.IR_Name("goog"), "requireType")

var gatherModuleMetadata_GOOG_MAYBE_REQUIRE = rhino.IR_Getprop(rhino.IR_Name("goog"), "maybeRequireFrameworkInternalOnlyDoNotCallOrElse")

var gatherModuleMetadata_GOOG_REQUIRE_DYNAMIC = rhino.IR_Getprop(rhino.IR_Name("goog"), "requireDynamic")

var gatherModuleMetadata_GOOG_SET_TEST_ONLY = rhino.IR_Getprop(rhino.IR_Name("goog"), "setTestOnly")

var gatherModuleMetadata_GOOG_MODULE_DECLARELEGACYNAMESPACE = rhino.IR_Getprop(gatherModuleMetadata_GOOG_MODULE.CloneTree(), "declareLegacyNamespace")

var gatherModuleMetadata_GOOG_DECLARE_MODULE_ID = rhino.IR_Getprop(rhino.IR_Name("goog"), "declareModuleId")

// TODO(johnplaisted): Remove once clients have migrated to declareModuleId
var gatherModuleMetadata_GOOG_MODULE_DECLARNAMESPACE = rhino.IR_Getprop(gatherModuleMetadata_GOOG_MODULE.CloneTree(), "declareNamespace")

const gatherModuleMetadata_TOGGLE_NAME_PREFIX = "TOGGLE_"

// Java: GatherModuleMetadata(AbstractCompiler, boolean, ModuleLoader.ResolutionMode)
func NewGatherModuleMetadata(
	compiler AbstractCompiler,
	processCommonJsModules bool,
	moduleResolutionMode ModuleLoader_ResolutionMode) *GatherModuleMetadata {
	g := &GatherModuleMetadata{
		modulesByPath:          jutil.NewLinkedHashMap[string, *ModuleMetadataMap_ModuleMetadata](),
		modulesByGoogNamespace: jutil.NewLinkedHashMap[string, *ModuleMetadataMap_ModuleMetadata](),
	}
	g.compiler = compiler
	g.processCommonJsModules = processCommonJsModules
	g.moduleResolutionMode = moduleResolutionMode
	return g
}

// GatherModuleMetadata_ModuleMetadataBuilder is the inner class
// GatherModuleMetadata.ModuleMetadataBuilder.
type GatherModuleMetadata_ModuleMetadataBuilder struct {
	outer *GatherModuleMetadata // the enclosing instance

	ambiguous               bool
	hasModuleBody           bool
	declaredModuleId        rhino.Node
	declaresLegacyNamespace rhino.Node
	metadataBuilder         *ModuleMetadataMap_ModuleMetadata_Builder
}

// Java: GatherModuleMetadata.ModuleMetadataBuilder(Node, ModulePath)
func NewGatherModuleMetadata_ModuleMetadataBuilder(outer *GatherModuleMetadata, rootNode rhino.Node, path *ModuleLoader_ModulePath) *GatherModuleMetadata_ModuleMetadataBuilder {
	b := &GatherModuleMetadata_ModuleMetadataBuilder{outer: outer}
	b.metadataBuilder = ModuleMetadataMap_ModuleMetadata_Builder_().
		Path(path).
		RootNode(rootNode).
		ModuleType(ModuleMetadataMap_ModuleType_SCRIPT).
		UsesClosure(false).
		IsTestOnly(false)
	return b
}

// Java: GatherModuleMetadata.ModuleMetadataBuilder.moduleType(ModuleType, NodeTraversal, Node)
func (b *GatherModuleMetadata_ModuleMetadataBuilder) ModuleType(type_ ModuleMetadataMap_ModuleType, t *NodeTraversal, n rhino.Node) {
	// checkNotNull(type): a Go enum value cannot be null.

	if b.metadataBuilder.GetModuleType() == type_ {
		return
	}

	if b.metadataBuilder.GetModuleType() == ModuleMetadataMap_ModuleType_SCRIPT {
		b.metadataBuilder.ModuleType(type_)
		return
	}

	b.ambiguous = true
	t.Report(n, GatherModuleMetadata_MIXED_MODULE_TYPE, b.metadataBuilder.GetModuleType().Description(), type_.Description())
}

// Java: GatherModuleMetadata.ModuleMetadataBuilder.recordDeclareModuleId(Node)
func (b *GatherModuleMetadata_ModuleMetadataBuilder) RecordDeclareModuleId(declaredModuleId rhino.Node) {
	b.declaredModuleId = declaredModuleId
}

// Java: GatherModuleMetadata.ModuleMetadataBuilder.recordDeclareLegacyNamespace(Node)
func (b *GatherModuleMetadata_ModuleMetadataBuilder) RecordDeclareLegacyNamespace(declaresLegacyNamespace rhino.Node) {
	b.declaresLegacyNamespace = declaresLegacyNamespace
}

// Java: GatherModuleMetadata.ModuleMetadataBuilder.isScript()
func (b *GatherModuleMetadata_ModuleMetadataBuilder) IsScript() bool {
	return b.metadataBuilder.GetModuleType() == ModuleMetadataMap_ModuleType_SCRIPT
}

// Java: GatherModuleMetadata.ModuleMetadataBuilder.build()
func (b *GatherModuleMetadata_ModuleMetadataBuilder) Build() *ModuleMetadataMap_ModuleMetadata {
	if !b.ambiguous {
		if b.hasModuleBody && b.metadataBuilder.GetModuleType() == ModuleMetadataMap_ModuleType_SCRIPT {
			// A script with no imports or exports, but has a module body, must be an ES module.
			b.metadataBuilder.ModuleType(ModuleMetadataMap_ModuleType_ES6_MODULE)
		}

		if !b.declaredModuleId.IsNil() && b.metadataBuilder.GetModuleType() != ModuleMetadataMap_ModuleType_ES6_MODULE {
			b.outer.compiler.Report(JSError_MakeWithNode(b.declaredModuleId, GatherModuleMetadata_DECLARE_MODULE_ID_OUTSIDE_ES6_MODULE))
		}

		if !b.declaresLegacyNamespace.IsNil() {
			if b.metadataBuilder.GetModuleType() == ModuleMetadataMap_ModuleType_GOOG_MODULE {
				b.metadataBuilder.ModuleType(ModuleMetadataMap_ModuleType_LEGACY_GOOG_MODULE)
			} else {
				b.outer.compiler.Report(
					JSError_MakeWithNode(b.declaresLegacyNamespace, ClosureCheckModule_DECLARE_LEGACY_NAMESPACE_IN_NON_MODULE))
			}
		}
	}

	return b.metadataBuilder.Build()
}

var gatherModuleMetadata_GOOG_LOADMODULE = rhino.QualifiedName_Of("goog.loadModule")

// GatherModuleMetadata_Finder traverses the AST and build a sets of
// ModuleMetadatas.
type GatherModuleMetadata_Finder struct {
	outer *GatherModuleMetadata // the enclosing instance

	// Store both names and vars. Strings alone is insufficient to determine
	// whether a name is actually a toggle module (since it could have been
	// shadowed, or may have been defined in a different file), but looking up
	// by only vars is much slower. This way we can do a fast name lookup,
	// followed by a slower var lookup only if the name is known to be a
	// toggle module name.
	toggleModuleNames map[string]struct{}
	toggleModules     map[*Var]struct{}
}

// Java: GatherModuleMetadata.Finder()
func NewGatherModuleMetadata_Finder(outer *GatherModuleMetadata) *GatherModuleMetadata_Finder {
	return &GatherModuleMetadata_Finder{
		outer:             outer,
		toggleModuleNames: make(map[string]struct{}),
		toggleModules:     make(map[*Var]struct{}),
	}
}

// Java: GatherModuleMetadata.Finder.shouldTraverse(NodeTraversal, Node, Node)
func (f *GatherModuleMetadata_Finder) ShouldTraverse(t *NodeTraversal, n rhino.Node, parent rhino.Node) bool {
	switch n.GetToken() {
	case rhino.Token_SCRIPT:
		f.EnterModule(t, n, t.GetInput().GetPath())
	case rhino.Token_IMPORT, rhino.Token_EXPORT:
		f.VisitImportOrExport(t, n)
	case rhino.Token_CALL:
		if n.IsCall() && gatherModuleMetadata_GOOG_LOADMODULE.Matches(n.GetFirstChild()) {
			f.outer.loadModuleCall = n
			f.EnterModule(t, n, nil)
		}
	case rhino.Token_MODULE_BODY:
		f.outer.currentModule.hasModuleBody = true
	case rhino.Token_DYNAMIC_IMPORT:
		f.VisitDynamicImport(ast.MustDynamicImport(n))
	default:
	}

	return true
}

// Java: GatherModuleMetadata.Finder.visitImportOrExport(NodeTraversal, Node)
func (f *GatherModuleMetadata_Finder) VisitImportOrExport(t *NodeTraversal, importOrExport rhino.Node) {
	jutil.CheckNotNull(f.outer.currentModule)
	f.outer.currentModule.ModuleType(ModuleMetadataMap_ModuleType_ES6_MODULE, t, importOrExport)
	var source ast.StringLit
	hasSource := false
	if imp, ok := ast.AsImport(importOrExport); ok {
		source, hasSource = imp.Source(), true
	} else if exp, ok := ast.AsExport(importOrExport); ok {
		source, hasSource = exp.Source() // export ... from
	}
	if hasSource {
		f.outer.currentModule.
			metadataBuilder.
			Es6ImportSpecifiersBuilder().
			Add(source.Value())
	}
}

// Java: GatherModuleMetadata.Finder.visitDynamicImport(Node)
func (f *GatherModuleMetadata_Finder) VisitDynamicImport(dynamicImport ast.DynamicImport) {
	if specifier, ok := dynamicImport.Specifier().AsStringLit(); ok {
		f.outer.currentModule.
			metadataBuilder.
			Es6ImportSpecifiersBuilder().
			Add(specifier.Value())
	}
}

// Java: GatherModuleMetadata.Finder.enterModule(NodeTraversal, Node, ModulePath)
func (f *GatherModuleMetadata_Finder) EnterModule(t *NodeTraversal, n rhino.Node, path *ModuleLoader_ModulePath) {
	newModule := NewGatherModuleMetadata_ModuleMetadataBuilder(f.outer, n, path)
	if f.outer.currentModule != nil {
		if f.outer.parentModule != nil {
			t.Report(n, GatherModuleMetadata_INVALID_NESTED_LOAD_MODULE)
		}
		f.outer.parentModule = f.outer.currentModule
	}
	f.outer.currentModule = newModule
}

// Java: GatherModuleMetadata.Finder.leaveModule()
func (f *GatherModuleMetadata_Finder) LeaveModule() {
	jutil.CheckNotNull(f.outer.currentModule)
	module := f.outer.currentModule.Build()
	if module.Path() != nil {
		f.outer.modulesByPath.Put(module.Path().String(), module)
	}
	for namespace := range module.GoogNamespaces().All() {
		f.outer.modulesByGoogNamespace.Put(namespace, module)
	}
	if f.outer.parentModule != nil {
		f.outer.parentModule.metadataBuilder.AddNestedModule(module)
	}
	f.outer.currentModule = f.outer.parentModule
	f.outer.parentModule = nil
}

// Java: GatherModuleMetadata.Finder.visit(NodeTraversal, Node, Node)
func (f *GatherModuleMetadata_Finder) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	if f.outer.processCommonJsModules && f.outer.currentModule != nil && f.outer.currentModule.IsScript() {
		// A common JS import (call to "require") does not force a module to be rewritten as
		// commonJS. Only an export statement.
		if ProcessCommonJSModules_IsCommonJsExport(t, n, f.outer.moduleResolutionMode) {
			f.outer.currentModule.ModuleType(ModuleMetadataMap_ModuleType_COMMON_JS, t, n)
			return
		}
	}

	switch n.GetToken() {
	case rhino.Token_SCRIPT:
		f.LeaveModule()
	case rhino.Token_NAME:
		f.VisitName(t, n)
	case rhino.Token_CALL:
		if f.outer.loadModuleCall == n {
			f.LeaveModule()
			f.outer.loadModuleCall = rhino.NilNode
		} else {
			f.VisitGoogCall(t, ast.MustCall(n))
		}
	default:
	}
}

// Java: GatherModuleMetadata.Finder.isFromGoogImport(Var)
func (f *GatherModuleMetadata_Finder) IsFromGoogImport(goog *Var) bool {
	nameNode := goog.GetNameNode()

	// Because other tools are regex based we force importing this file as "import * as goog".
	return !nameNode.IsNil() &&
		nameNode.IsImportStar() &&
		nameNode.GetString() == "goog" &&
		nameNode.GetParent().GetFirstChild().IsEmpty() &&
		strings.HasSuffix(nameNode.GetParent().GetLastChild().GetString(), "/goog.js")
}

// Java: GatherModuleMetadata.Finder.visitName(NodeTraversal, Node)
func (f *GatherModuleMetadata_Finder) VisitName(t *NodeTraversal, n rhino.Node) {
	name := n.GetString()
	if _, ok := f.toggleModuleNames[name]; ok {
		nameVar := t.GetScope().GetVar(name)
		if _, ok := f.toggleModules[nameVar]; ok {
			parent := n.GetParent()
			if parent.IsGetProp() {
				f.AddToggle(t, n, parent.GetString())
			} else if !NodeUtil_IsNameDeclaration(parent) {
				t.Report(
					n,
					GatherModuleMetadata_INVALID_TOGGLE_USAGE,
					"toggle modules may not be used other than looking up properties")
			}
		}
	}
	if "goog" != name {
		return
	}

	root := t.GetScope().GetVar("goog")
	if root != nil && !f.IsFromGoogImport(root) {
		return
	}

	f.outer.currentModule.metadataBuilder.UsesClosure(true)
}

// Java: GatherModuleMetadata.Finder.visitGoogCall(NodeTraversal, Node)
func (f *GatherModuleMetadata_Finder) VisitGoogCall(t *NodeTraversal, call ast.Call) {
	n := call.Node()
	callee, ok := call.Callee().AsGetProp()
	if !ok || !callee.Node().IsQualifiedName() {
		return
	}
	getprop := callee.Node()

	// The root of the property chain: goog in goog.module.get.
	firstProp := callee.Object()
	for {
		inner, ok := firstProp.AsGetProp()
		if !ok {
			break
		}
		firstProp = inner.Object()
	}

	if name, ok := firstProp.AsName(); !ok || name.Ident() != "goog" {
		return
	}

	root := t.GetScope().GetVar("goog")

	// If this is a locally defined variable it can't be the global "goog", so exit early.
	if root != nil && root.IsLocal() && !root.GetScope().(*Scope).IsModuleScope() {
		return
	}

	// If this is a module-level variable but wasn't created by importing goog.js, return.
	if root != nil && root.GetScope().(*Scope).IsModuleScope() && !f.IsFromGoogImport(root) {
		return
	}

	// If goog is defined in this script then it does not use Closure. If this is a bundle with
	// base.js in it, then it doesn't need base.js again.
	if root == nil ||
		NodeUtil_GetEnclosingScript(root.GetNameNode()) != NodeUtil_GetEnclosingScript(n) {
		f.outer.currentModule.metadataBuilder.UsesClosure(true)
	}

	if getprop.MatchesQualifiedNameNode(gatherModuleMetadata_GOOG_PROVIDE) {
		f.outer.currentModule.ModuleType(ModuleMetadataMap_ModuleType_GOOG_PROVIDE, t, n)
		if namespace, ok := gatherModuleMetadata_onlyStringArg(call); ok {
			f.AddNamespace(f.outer.currentModule, ModuleMetadataMap_ModuleType_GOOG_PROVIDE, namespace, t, n)
		} else {
			t.Report(n, ClosureRewriteModule_INVALID_PROVIDE_NAMESPACE)
			f.outer.currentModule.metadataBuilder.UsesClosure(false)
		}
	} else if getprop.MatchesQualifiedNameNode(gatherModuleMetadata_GOOG_MODULE) {
		f.outer.currentModule.ModuleType(ModuleMetadataMap_ModuleType_GOOG_MODULE, t, n)
		if namespace, ok := gatherModuleMetadata_onlyStringArg(call); ok {
			f.AddNamespace(f.outer.currentModule, ModuleMetadataMap_ModuleType_GOOG_MODULE, namespace, t, n)
		} else {
			t.Report(n, ClosureRewriteModule_INVALID_MODULE_ID_ARG)
			f.outer.currentModule.metadataBuilder.UsesClosure(false)
		}
	} else if getprop.MatchesQualifiedNameNode(gatherModuleMetadata_GOOG_MODULE_GET) {
		// Look for a module named ???$2etoggles and a getprop of .TOGGLE_??? on the call's return
		if namespace, ok := gatherModuleMetadata_onlyStringArg(call); !ok || !strings.HasSuffix(namespace, "$2etoggles") {
			return // only do anything with toggle namespaces
		}
		parent := n.GetParent()
		if !parent.IsGetProp() {
			t.Report(
				n,
				GatherModuleMetadata_INVALID_TOGGLE_USAGE,
				"goog.module.get of toggles module must immediately look up a single toggle")
			return
		}
		f.AddToggle(t, parent, parent.GetString())
	} else if getprop.MatchesQualifiedNameNode(gatherModuleMetadata_GOOG_MODULE_DECLARELEGACYNAMESPACE) {
		f.outer.currentModule.RecordDeclareLegacyNamespace(n)
	} else if getprop.MatchesQualifiedNameNode(gatherModuleMetadata_GOOG_DECLARE_MODULE_ID) ||
		getprop.MatchesQualifiedNameNode(gatherModuleMetadata_GOOG_MODULE_DECLARNAMESPACE) {
		if !f.outer.currentModule.declaredModuleId.IsNil() {
			t.Report(n, GatherModuleMetadata_MULTIPLE_DECLARE_MODULE_NAMESPACE)
		}
		if namespace, ok := gatherModuleMetadata_onlyStringArg(call); ok {
			f.outer.currentModule.RecordDeclareModuleId(n)
			f.AddNamespace(f.outer.currentModule, ModuleMetadataMap_ModuleType_GOOG_MODULE, namespace, t, n)
		} else {
			t.Report(n, GatherModuleMetadata_INVALID_DECLARE_MODULE_ID_CALL)
		}
	} else if getprop.MatchesQualifiedNameNode(gatherModuleMetadata_GOOG_REQUIRE) {
		if namespace, ok := gatherModuleMetadata_onlyStringArg(call); ok {
			f.outer.currentModule.metadataBuilder.StronglyRequiredGoogNamespacesBuilder().Add(namespace)
			if strings.HasSuffix(namespace, "$2etoggles") {
				// Track imports of *.toggles.ts, which are rewritten to $2etoggles.
				callParent := n.GetParent()
				lhs := callParent.GetFirstChild()
				if callParent.IsDestructuringLhs() {
					// const {TOGGLE_foo} = goog.require('foo$2etoggles');
					for key := range lhs.Children() {
						if key.IsStringKey() {
							f.AddToggle(t, n, key.GetString())
						} else {
							t.Report(n, GatherModuleMetadata_INVALID_TOGGLE_USAGE, "must be destructured with string keys")
						}
					}
				} else if callParent.IsName() {
					// const fooToggles = goog.require('foo$2etoggles');
					name := callParent.GetString()
					nameVar := t.GetScope().GetVar(name)
					f.toggleModules[nameVar] = struct{}{}
					f.toggleModuleNames[name] = struct{}{}
				} else if f.outer.currentModule.metadataBuilder.GetModuleType() != ModuleMetadataMap_ModuleType_GOOG_PROVIDE {
					// Side-effect toggle-module imports are not allowed, since they don't actually do
					// anything.  We allow it in `goog.provide()` files because there's no other way to
					// import, and since toggle modules don't declare a legacy namespace, it's unusable
					// without a `goog.module.get()` (so we can catch the toggle use there, instead).
					t.Report(n, GatherModuleMetadata_INVALID_TOGGLE_USAGE, "import must be assigned")
				}
			}
		} else {
			t.Report(n, ClosurePrimitiveErrors_INVALID_REQUIRE_NAMESPACE)
		}
	} else if getprop.MatchesQualifiedNameNode(gatherModuleMetadata_GOOG_REQUIRE_TYPE) {
		if namespace, ok := gatherModuleMetadata_onlyStringArg(call); ok {
			f.outer.currentModule.
				metadataBuilder.
				WeaklyRequiredGoogNamespacesBuilder().
				Add(namespace)
		} else {
			t.Report(n, GatherModuleMetadata_INVALID_REQUIRE_TYPE)
		}
	} else if getprop.MatchesQualifiedNameNode(gatherModuleMetadata_GOOG_MAYBE_REQUIRE) {
		if namespace, ok := gatherModuleMetadata_onlyStringArg(call); ok {
			f.outer.currentModule.
				metadataBuilder.
				MaybeRequiredGoogNamespacesBuilder().
				Add(namespace)
		} else {
			t.Report(n, GatherModuleMetadata_INVALID_MAYBE_REQUIRE)
		}
	} else if getprop.MatchesQualifiedNameNode(gatherModuleMetadata_GOOG_SET_TEST_ONLY) {
		if _, isString := gatherModuleMetadata_onlyStringArg(call); !call.HasArgs() || isString {
			f.outer.currentModule.metadataBuilder.IsTestOnly(true)
		} else {
			t.Report(n, GatherModuleMetadata_INVALID_SET_TEST_ONLY)
		}
	} else if getprop.MatchesQualifiedNameNode(gatherModuleMetadata_GOOG_REQUIRE_DYNAMIC) {
		if namespace, ok := gatherModuleMetadata_onlyStringArg(call); ok {
			f.outer.currentModule.
				metadataBuilder.
				DynamicallyRequiredGoogNamespacesBuilder().
				Add(namespace)
		} else {
			t.Report(n, GatherModuleMetadata_INVALID_REQUIRE_DYNAMIC)
		}
	}
}

// gatherModuleMetadata_onlyStringArg returns the value of call's argument
// if it has exactly one and it is a string literal: goog.provide('a.b').
//
// Go-only: Java's n.hasTwoChildren() && n.getLastChild().isStringLit().
func gatherModuleMetadata_onlyStringArg(call ast.Call) (string, bool) {
	arg, ok := call.OnlyArg()
	if !ok {
		return "", false
	}
	s, ok := arg.AsStringLit()
	if !ok {
		return "", false
	}
	return s.Value(), true
}

// AddToggle records a toggle usage (either a destructured import or a
// property lookup on a module).
//
// Java: GatherModuleMetadata.Finder.addToggle(NodeTraversal, Node, String)
func (f *GatherModuleMetadata_Finder) AddToggle(t *NodeTraversal, n rhino.Node, name string) {
	if strings.HasPrefix(name, gatherModuleMetadata_TOGGLE_NAME_PREFIX) {
		toggleName := jutil.StringSubstringFrom(name, jutil.StringLength(gatherModuleMetadata_TOGGLE_NAME_PREFIX))
		f.outer.currentModule.metadataBuilder.ReadTogglesBuilder().Add(toggleName)
	} else {
		t.Report(n, GatherModuleMetadata_INVALID_TOGGLE_USAGE, "all toggle names must start with `TOGGLE_`")
	}
}

// AddNamespace adds the namespaces to the module and checks if the given
// Closure namespace is a duplicate or not.
//
// Java: GatherModuleMetadata.Finder.addNamespace(ModuleMetadataBuilder, ModuleType, String, NodeTraversal, Node)
func (f *GatherModuleMetadata_Finder) AddNamespace(
	module *GatherModuleMetadata_ModuleMetadataBuilder,
	moduleType ModuleMetadataMap_ModuleType,
	namespace string,
	t *NodeTraversal,
	n rhino.Node) {
	if moduleType == ModuleMetadataMap_ModuleType_GOOG_PROVIDE ||
		moduleType == ModuleMetadataMap_ModuleType_LEGACY_GOOG_MODULE {
		if !NodeUtil_IsValidQualifiedNameFeatureSet(
			f.outer.compiler.GetOptions().GetLanguageIn().ToFeatureSet(), namespace) {
			f.outer.compiler.Report(JSError_MakeWithNode(n, GatherModuleMetadata_INVALID_NAMESPACE_OR_MODULE_ID, namespace))
		}
	}
	if moduleType == ModuleMetadataMap_ModuleType_GOOG_MODULE ||
		moduleType == ModuleMetadataMap_ModuleType_LEGACY_GOOG_MODULE {
		// non-legacy goog.modules don't technically need to be valid qualified names
		if !GatherModuleMetadata_IsValidModuleId(namespace) {
			f.outer.compiler.Report(JSError_MakeWithNode(n, GatherModuleMetadata_INVALID_NAMESPACE_OR_MODULE_ID, namespace))
		}
	}

	// existingType is null in Java until a duplicate is found.
	var existingType ModuleMetadataMap_ModuleType
	hasExistingType := false
	var existingFileSource *string
	if module.metadataBuilder.HasGoogNamespace(namespace) {
		existingType = module.metadataBuilder.GetModuleType()
		hasExistingType = true
		existingFileSource = t.GetSourceName()
	} else {
		existingModule := f.outer.modulesByGoogNamespace.Get(namespace)
		if existingModule != nil {
			existingType = existingModule.ModuleType()
			hasExistingType = true
			existingFileSource = existingModule.RootNode().GetSourceFileName()
		}
	}
	f.outer.currentModule.metadataBuilder.GoogNamespacesBuilder().Add(namespace)
	if hasExistingType {
		switch existingType {
		case ModuleMetadataMap_ModuleType_ES6_MODULE, ModuleMetadataMap_ModuleType_GOOG_MODULE, ModuleMetadataMap_ModuleType_LEGACY_GOOG_MODULE:
			diagnostic := ClosurePrimitiveErrors_DUPLICATE_MODULE
			if moduleType == ModuleMetadataMap_ModuleType_GOOG_PROVIDE {
				diagnostic = ClosurePrimitiveErrors_DUPLICATE_NAMESPACE_AND_MODULE
			}
			t.Report(n, diagnostic, namespace, jutil.ToJavaString(existingFileSource))
			return
		case ModuleMetadataMap_ModuleType_GOOG_PROVIDE:
			diagnostic := ClosurePrimitiveErrors_DUPLICATE_NAMESPACE_AND_MODULE
			if moduleType == ModuleMetadataMap_ModuleType_GOOG_PROVIDE {
				diagnostic = ClosurePrimitiveErrors_DUPLICATE_NAMESPACE
			}
			t.Report(n, diagnostic, namespace, jutil.ToJavaString(existingFileSource))
			return
		case ModuleMetadataMap_ModuleType_COMMON_JS, ModuleMetadataMap_ModuleType_SCRIPT:
		}
		panic(jutil.NewIllegalStateExceptionMsg("Unexpected module type: " + existingType.String()))
	}
}

// Java: GatherModuleMetadata.process(Node, Node)
func (g *GatherModuleMetadata) Process(externs rhino.Node, root rhino.Node) {
	sc := NewSyntacticScopeCreator(g.compiler)
	NodeTraversal_Builder_().SetCompiler(g.compiler).SetScopeCreator(sc).SetCallback(NewGatherModuleMetadata_Finder(g)).Traverse(externs)
	NodeTraversal_Builder_().SetCompiler(g.compiler).SetScopeCreator(sc).SetCallback(NewGatherModuleMetadata_Finder(g)).Traverse(root)
	g.compiler.SetModuleMetadataMap(NewModuleMetadataMap(g.modulesByPath, g.modulesByGoogNamespace))
}

// Must match closure/base.js's goog.VALID_MODULE_RE_ & also validates that
// dotted segments are non-empty.
//
// Java: GatherModuleMetadata.isValidModuleId(String)
func GatherModuleMetadata_IsValidModuleId(id string) bool {
	for _, segment := range gatherModuleMetadata_DOT_SPLITTER(id) {
		if jutil.StringLength(segment) == 0 {
			return false
		}
		for _, c := range jutil.StringToUTF16(segment) {
			if !parser.Identifiers_IsIdentifierPart(c) {
				return false
			}
		}
	}
	return gatherModuleMetadata_NAMESPACE_SEGMENT_REGEX.Matcher(id).Matches()
}

var gatherModuleMetadata_NAMESPACE_SEGMENT_REGEX = jutil.Pattern_Compile("^[a-zA-Z_$][a-zA-Z0-9_$.]*$")

// gatherModuleMetadata_DOT_SPLITTER is Splitter.on('.').split(id): every
// piece between dots, empty ones included.
//
// Go-only: language-adapter — a Guava Splitter is a function here; '.' is
// ASCII, so splitting the WTF-8 bytes splits the UTF-16 units alike.
func gatherModuleMetadata_DOT_SPLITTER(id string) []string {
	return strings.Split(id, ".")
}
