// Port of com/google/javascript/jscomp/GatherModuleMetadata.java @ closure-compiler v20261104

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
	"github.com/paullepoulpe/jsgo/internal/strpool"
)

// GatherModuleMetadata gathers metadata around modules that is useful for
// checking imports / requires and creates a ModuleMetadataMap.
type GatherModuleMetadata struct {
	// Map from module path to module. These modules represent files and thus
	// will contain all goog namespaces that are in the file. These are not
	// the same modules in modulesByGoogNamespace.
	modulesByPath OrderedStringMap[*ModuleMetadataMap_ModuleMetadata]

	// Map from Closure namespace to module. These modules represent just the
	// single namespace and thus each module has only one goog namespace in
	// its ModuleMetadata#googNamespaces(). These are not the same modules in
	// modulesByPath.
	modulesByGoogNamespace OrderedStringMap[*ModuleMetadataMap_ModuleMetadata]

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

const gatherModuleMetadata_TOGGLE_NAME_PREFIX = "TOGGLE_"

var (
	gatherMetaGoogID       = strpool.Intern("goog")
	gatherMetaLoadModuleID = strpool.Intern("loadModule")
)

// Java: GatherModuleMetadata(AbstractCompiler, boolean, ModuleLoader.ResolutionMode)
func NewGatherModuleMetadata(
	compiler AbstractCompiler,
	processCommonJsModules bool,
	moduleResolutionMode ModuleLoader_ResolutionMode) *GatherModuleMetadata {
	return &GatherModuleMetadata{
		compiler:               compiler,
		processCommonJsModules: processCommonJsModules,
		moduleResolutionMode:   moduleResolutionMode,
	}
}

// gatherModuleMetadataBuilder is the inner class
// GatherModuleMetadata.ModuleMetadataBuilder.
type gatherModuleMetadataBuilder struct {
	compiler                AbstractCompiler
	metadata                *ModuleMetadataMap_ModuleMetadata
	declaredModuleId        rhino.Node
	declaresLegacyNamespace rhino.Node
	ambiguous               bool
	hasModuleBody           bool
}

func (b *gatherModuleMetadataBuilder) reset(compiler AbstractCompiler, rootNode rhino.Node, path *ModuleLoader_ModulePath) {
	*b = gatherModuleMetadataBuilder{
		compiler: compiler,
		metadata: &ModuleMetadataMap_ModuleMetadata{
			path:       path,
			rootNode:   rootNode,
			moduleType: ModuleMetadataMap_ModuleType_SCRIPT,
		},
	}
}

// Java: GatherModuleMetadata.ModuleMetadataBuilder(Node, ModulePath)
func newGatherModuleMetadataBuilder(compiler AbstractCompiler, rootNode rhino.Node, path *ModuleLoader_ModulePath) *gatherModuleMetadataBuilder {
	b := &gatherModuleMetadataBuilder{}
	b.reset(compiler, rootNode, path)
	return b
}

// Java: GatherModuleMetadata.ModuleMetadataBuilder.moduleType(ModuleType, NodeTraversal, Node)
func (b *gatherModuleMetadataBuilder) setModuleType(type_ ModuleMetadataMap_ModuleType, t *NodeTraversal, n rhino.Node) {
	if b.metadata.moduleType == type_ {
		return
	}
	if b.metadata.moduleType == ModuleMetadataMap_ModuleType_SCRIPT {
		b.metadata.moduleType = type_
		return
	}
	b.ambiguous = true
	t.Report(n, GatherModuleMetadata_MIXED_MODULE_TYPE, b.metadata.moduleType.Description(), type_.Description())
}

// Java: GatherModuleMetadata.ModuleMetadataBuilder.recordDeclareModuleId(Node)
func (b *gatherModuleMetadataBuilder) recordDeclareModuleId(declaredModuleId rhino.Node) {
	b.declaredModuleId = declaredModuleId
}

// Java: GatherModuleMetadata.ModuleMetadataBuilder.recordDeclareLegacyNamespace(Node)
func (b *gatherModuleMetadataBuilder) recordDeclareLegacyNamespace(declaresLegacyNamespace rhino.Node) {
	b.declaresLegacyNamespace = declaresLegacyNamespace
}

// Java: GatherModuleMetadata.ModuleMetadataBuilder.isScript()
func (b *gatherModuleMetadataBuilder) isScript() bool {
	return b.metadata.moduleType == ModuleMetadataMap_ModuleType_SCRIPT
}

// Java: GatherModuleMetadata.ModuleMetadataBuilder.build()
func (b *gatherModuleMetadataBuilder) build() *ModuleMetadataMap_ModuleMetadata {
	if !b.ambiguous {
		if b.hasModuleBody && b.metadata.moduleType == ModuleMetadataMap_ModuleType_SCRIPT {
			// A script with no imports or exports, but has a module body, must be an ES module.
			b.metadata.moduleType = ModuleMetadataMap_ModuleType_ES6_MODULE
		}

		if !b.declaredModuleId.IsNil() && b.metadata.moduleType != ModuleMetadataMap_ModuleType_ES6_MODULE {
			b.compiler.Report(JSError_MakeWithNode(b.declaredModuleId, GatherModuleMetadata_DECLARE_MODULE_ID_OUTSIDE_ES6_MODULE))
		}

		if !b.declaresLegacyNamespace.IsNil() {
			if b.metadata.moduleType == ModuleMetadataMap_ModuleType_GOOG_MODULE {
				b.metadata.moduleType = ModuleMetadataMap_ModuleType_LEGACY_GOOG_MODULE
			} else {
				b.compiler.Report(
					JSError_MakeWithNode(b.declaresLegacyNamespace, ClosureCheckModule_DECLARE_LEGACY_NAMESPACE_IN_NON_MODULE))
			}
		}
	}

	return b.metadata
}

var (
	gatherMetaVisitTokens = func() ast.TokenSet {
		var s ast.TokenSet
		s.Add(rhino.Token_SCRIPT)
		s.Add(rhino.Token_NAME)
		s.Add(rhino.Token_CALL)
		return s
	}()
	gatherMetaCommonJsVisitTokens = func() ast.TokenSet {
		s := gatherMetaVisitTokens
		s.Add(rhino.Token_GETPROP)
		s.Add(rhino.Token_GETELEM)
		return s
	}()
)

type gatherMetaOpKind uint8

const (
	gatherMetaOpAddNamespace gatherMetaOpKind = iota
	gatherMetaOpLeaveModule
)

type gatherMetaScriptOp struct {
	kind       gatherMetaOpKind
	eventIdx   int
	module     *ModuleMetadataMap_ModuleMetadata
	moduleType ModuleMetadataMap_ModuleType
	namespace  string
	node       rhino.Node
	sourceName string
	sameModDup bool
	sameModTyp ModuleMetadataMap_ModuleType
}

// moduleMetadataFinder traverses the AST and builds ModuleMetadatas.
//
// Java: GatherModuleMetadata.Finder
type moduleMetadataFinder struct {
	outer          *GatherModuleMetadata
	compiler       AbstractCompiler
	langInFeatures parser.FeatureSet

	currentModule  *gatherModuleMetadataBuilder
	parentModule   *gatherModuleMetadataBuilder
	rootSlot       gatherModuleMetadataBuilder
	nestedSlot     gatherModuleMetadataBuilder
	loadModuleCall rhino.Node

	onOp func(gatherMetaScriptOp)

	// Lazily allocated only if a $2etoggles goog.require is encountered.
	toggleModuleIDs map[strpool.ID]struct{}
	toggleModules   map[*Var]struct{}
}

var (
	_ NodeTraversal_Callback       = (*moduleMetadataFinder)(nil)
	_ nodeTraversal_visitTokens    = (*moduleMetadataFinder)(nil)
	_ nodeTraversal_onlyNameLeaves = (*moduleMetadataFinder)(nil)
)

func (f *moduleMetadataFinder) visitTokens() *ast.TokenSet {
	if f.outer != nil && f.outer.processCommonJsModules {
		return &gatherMetaCommonJsVisitTokens
	}
	return &gatherMetaVisitTokens
}

func (*moduleMetadataFinder) onlyNameLeaves() {}

// Java: GatherModuleMetadata.Finder.shouldTraverse(NodeTraversal, Node, Node)
func (f *moduleMetadataFinder) ShouldTraverse(t *NodeTraversal, n rhino.Node, parent rhino.Node) bool {
	switch n.GetToken() {
	case rhino.Token_SCRIPT:
		f.enterModule(t, n, t.GetInput().GetPath())
	case rhino.Token_IMPORT, rhino.Token_EXPORT:
		f.visitImportOrExport(t, n)
	case rhino.Token_CALL:
		if call, ok := ast.AsCall(n); ok {
			if gp, ok := call.Callee().AsGetProp(); ok && nodeStringID(gp.Node()) == gatherMetaLoadModuleID {
				if base, ok := gp.Object().AsName(); ok && nodeStringID(base.Node()) == gatherMetaGoogID {
					f.loadModuleCall = n
					f.enterModule(t, n, nil)
				}
			}
		}
	case rhino.Token_MODULE_BODY:
		f.currentModule.hasModuleBody = true
	case rhino.Token_DYNAMIC_IMPORT:
		f.visitDynamicImport(ast.MustDynamicImport(n))
	default:
	}

	return true
}

// Java: GatherModuleMetadata.Finder.visitImportOrExport(NodeTraversal, Node)
func (f *moduleMetadataFinder) visitImportOrExport(t *NodeTraversal, importOrExport rhino.Node) {
	jutil.CheckNotNull(f.currentModule)
	f.currentModule.setModuleType(ModuleMetadataMap_ModuleType_ES6_MODULE, t, importOrExport)
	var source ast.StringLit
	hasSource := false
	if imp, ok := ast.AsImport(importOrExport); ok {
		source, hasSource = imp.Source(), true
	} else if exp, ok := ast.AsExport(importOrExport); ok {
		source, hasSource = exp.Source() // export ... from
	}
	if hasSource {
		f.currentModule.metadata.es6ImportSpecifiers.Add(source.Value())
	}
}

// Java: GatherModuleMetadata.Finder.visitDynamicImport(Node)
func (f *moduleMetadataFinder) visitDynamicImport(dynamicImport ast.DynamicImport) {
	if specifier, ok := dynamicImport.Specifier().AsStringLit(); ok {
		f.currentModule.metadata.es6ImportSpecifiers.Add(specifier.Value())
	}
}

// Java: GatherModuleMetadata.Finder.enterModule(NodeTraversal, Node, ModulePath)
func (f *moduleMetadataFinder) enterModule(t *NodeTraversal, n rhino.Node, path *ModuleLoader_ModulePath) {
	var newModule *gatherModuleMetadataBuilder
	if f.currentModule == nil {
		f.rootSlot.reset(f.compiler, n, path)
		newModule = &f.rootSlot
	} else {
		if f.parentModule != nil {
			t.Report(n, GatherModuleMetadata_INVALID_NESTED_LOAD_MODULE)
			newModule = newGatherModuleMetadataBuilder(f.compiler, n, path)
		} else {
			f.nestedSlot.reset(f.compiler, n, path)
			newModule = &f.nestedSlot
		}
		f.parentModule = f.currentModule
	}
	f.currentModule = newModule
}

// Java: GatherModuleMetadata.Finder.leaveModule()
func (f *moduleMetadataFinder) leaveModule() {
	jutil.CheckNotNull(f.currentModule)
	module := f.currentModule.build()
	if f.onOp != nil {
		f.onOp(gatherMetaScriptOp{
			kind:   gatherMetaOpLeaveModule,
			module: module,
		})
	} else {
		f.outer.recordBuiltModule(module)
	}
	if f.parentModule != nil {
		f.parentModule.metadata.nestedModules = append(f.parentModule.metadata.nestedModules, module)
	}
	*f.currentModule = gatherModuleMetadataBuilder{}
	f.currentModule = f.parentModule
	f.parentModule = nil
}

func (g *GatherModuleMetadata) recordBuiltModule(module *ModuleMetadataMap_ModuleMetadata) {
	if module.Path() != nil {
		g.modulesByPath.Put(module.Path().String(), module)
	}
	for namespace := range module.GoogNamespaces().All() {
		g.modulesByGoogNamespace.Put(namespace, module)
	}
}

// Java: GatherModuleMetadata.Finder.visit(NodeTraversal, Node, Node)
func (f *moduleMetadataFinder) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	if f.outer.processCommonJsModules && f.currentModule != nil && f.currentModule.isScript() {
		// A common JS import (call to "require") does not force a module to be rewritten as
		// commonJS. Only an export statement.
		if ProcessCommonJSModules_IsCommonJsExport(t, n, f.outer.moduleResolutionMode) {
			f.currentModule.setModuleType(ModuleMetadataMap_ModuleType_COMMON_JS, t, n)
			return
		}
	}

	switch n.GetToken() {
	case rhino.Token_SCRIPT:
		f.leaveModule()
	case rhino.Token_NAME:
		f.visitName(t, n, parent)
	case rhino.Token_CALL:
		if f.loadModuleCall == n {
			f.leaveModule()
			f.loadModuleCall = rhino.NilNode
		} else {
			f.visitGoogCall(t, ast.MustCall(n), parent)
		}
	default:
	}
}

// Java: GatherModuleMetadata.Finder.isFromGoogImport(Var)
func (f *moduleMetadataFinder) isFromGoogImport(goog *Var) bool {
	// Because other tools are regex based we force importing this file as "import * as goog".
	star, ok := ast.AsImportStar(goog.GetNameNode())
	if !ok || nodeStringID(star.Node()) != gatherMetaGoogID {
		return false
	}
	imp, ok := ast.AsImport(star.Node().GetParent())
	if !ok {
		return false
	}
	_, hasDefault := imp.Default()
	return !hasDefault && strings.HasSuffix(imp.Source().Value(), "/goog.js")
}

// Java: GatherModuleMetadata.Finder.visitName(NodeTraversal, Node)
func (f *moduleMetadataFinder) visitName(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	id := nodeStringID(n)
	if f.toggleModuleIDs != nil {
		if _, ok := f.toggleModuleIDs[id]; ok {
			nameVar := t.GetScope().GetVarByID(uint32(id))
			if _, ok := f.toggleModules[nameVar]; ok {
				if gp, ok := ast.AsGetProp(parent); ok {
					f.addToggle(t, n, gp.Property())
				} else if !NodeUtil_IsNameDeclaration(parent) {
					t.Report(
						n,
						GatherModuleMetadata_INVALID_TOGGLE_USAGE,
						"toggle modules may not be used other than looking up properties")
				}
			}
		}
	}
	if id != gatherMetaGoogID {
		return
	}

	root := t.GetScope().GetVarByID(uint32(gatherMetaGoogID))
	if root != nil && !f.isFromGoogImport(root) {
		return
	}

	f.currentModule.metadata.usesClosure = true
}

// Java: GatherModuleMetadata.Finder.visitGoogCall(NodeTraversal, Node)
func (f *moduleMetadataFinder) visitGoogCall(t *NodeTraversal, call ast.Call, parent rhino.Node) {
	callee, ok := call.Callee().AsGetProp()
	if !ok {
		return
	}

	// Walk the GetProp chain to its root object, capturing depth (1 = goog.X, 2 = goog.Y.X).
	var innerPropNode ast.GetProp
	depth := 1
	firstProp := callee.Object()
	for {
		inner, ok := firstProp.AsGetProp()
		if !ok {
			break
		}
		innerPropNode = inner
		firstProp = inner.Object()
		depth++
	}

	if name, ok := firstProp.AsName(); !ok || nodeStringID(name.Node()) != gatherMetaGoogID {
		return
	}

	prop := callee.Property()
	parentProp := ""
	if depth > 1 {
		parentProp = innerPropNode.Property()
	}

	n := call.Node()
	root := t.GetScope().GetVarByID(uint32(gatherMetaGoogID))
	if root != nil {
		scope := root.GetScope().(*Scope)
		if scope.IsModuleScope() {
			if !f.isFromGoogImport(root) {
				return
			}
		} else if root.IsLocal() {
			return
		}
	}

	// If goog is defined in this script then it does not use Closure. If this is a bundle with
	// base.js in it, then it doesn't need base.js again.
	if root == nil ||
		NodeUtil_GetEnclosingScript(root.GetNameNode()) != NodeUtil_GetEnclosingScript(n) {
		f.currentModule.metadata.usesClosure = true
	}

	switch {
	case depth == 1:
		switch prop {
		case "provide":
			f.currentModule.setModuleType(ModuleMetadataMap_ModuleType_GOOG_PROVIDE, t, n)
			if namespace, ok := closureOnlyStringArg(call); ok {
				f.addNamespace(f.currentModule, ModuleMetadataMap_ModuleType_GOOG_PROVIDE, namespace, t, n)
			} else {
				t.Report(n, ClosureRewriteModule_INVALID_PROVIDE_NAMESPACE)
				f.currentModule.metadata.usesClosure = false
			}
		case "module":
			f.currentModule.setModuleType(ModuleMetadataMap_ModuleType_GOOG_MODULE, t, n)
			if namespace, ok := closureOnlyStringArg(call); ok {
				f.addNamespace(f.currentModule, ModuleMetadataMap_ModuleType_GOOG_MODULE, namespace, t, n)
			} else {
				t.Report(n, ClosureRewriteModule_INVALID_MODULE_ID_ARG)
				f.currentModule.metadata.usesClosure = false
			}
		case "declareModuleId":
			f.visitDeclareModuleId(t, call, n)
		case "require":
			if namespace, ok := closureOnlyStringArg(call); ok {
				f.currentModule.metadata.stronglyRequiredGoogNamespaces.Add(namespace)
				if strings.HasSuffix(namespace, "$2etoggles") {
					f.recordToggleRequire(t, n, parent)
				}
			} else {
				t.Report(n, ClosurePrimitiveErrors_INVALID_REQUIRE_NAMESPACE)
			}
		case "requireType":
			if namespace, ok := closureOnlyStringArg(call); ok {
				f.currentModule.metadata.weaklyRequiredGoogNamespaces.Add(namespace)
			} else {
				t.Report(n, GatherModuleMetadata_INVALID_REQUIRE_TYPE)
			}
		case "maybeRequireFrameworkInternalOnlyDoNotCallOrElse":
			if namespace, ok := closureOnlyStringArg(call); ok {
				f.currentModule.metadata.maybeRequiredGoogNamespaces.Add(namespace)
			} else {
				t.Report(n, GatherModuleMetadata_INVALID_MAYBE_REQUIRE)
			}
		case "setTestOnly":
			if _, isString := closureOnlyStringArg(call); !call.HasArgs() || isString {
				f.currentModule.metadata.isTestOnly = true
			} else {
				t.Report(n, GatherModuleMetadata_INVALID_SET_TEST_ONLY)
			}
		case "requireDynamic":
			if namespace, ok := closureOnlyStringArg(call); ok {
				f.currentModule.metadata.dynamicallyRequiredGoogNamespaces.Add(namespace)
			} else {
				t.Report(n, GatherModuleMetadata_INVALID_REQUIRE_DYNAMIC)
			}
		}
	case depth == 2 && parentProp == "module":
		switch prop {
		case "get":
			// Look for a module named ???$2etoggles and a getprop of .TOGGLE_??? on the call's return
			if namespace, ok := closureOnlyStringArg(call); !ok || !strings.HasSuffix(namespace, "$2etoggles") {
				return
			}
			gp, ok := ast.AsGetProp(parent)
			if !ok {
				t.Report(
					n,
					GatherModuleMetadata_INVALID_TOGGLE_USAGE,
					"goog.module.get of toggles module must immediately look up a single toggle")
				return
			}
			f.addToggle(t, gp.Node(), gp.Property())
		case "declareLegacyNamespace":
			f.currentModule.recordDeclareLegacyNamespace(n)
		case "declareNamespace":
			// TODO(johnplaisted): Remove once clients have migrated to declareModuleId
			f.visitDeclareModuleId(t, call, n)
		}
	}
}

func (f *moduleMetadataFinder) visitDeclareModuleId(t *NodeTraversal, call ast.Call, n rhino.Node) {
	if !f.currentModule.declaredModuleId.IsNil() {
		t.Report(n, GatherModuleMetadata_MULTIPLE_DECLARE_MODULE_NAMESPACE)
	}
	if namespace, ok := closureOnlyStringArg(call); ok {
		f.currentModule.recordDeclareModuleId(n)
		f.addNamespace(f.currentModule, ModuleMetadataMap_ModuleType_GOOG_MODULE, namespace, t, n)
	} else {
		t.Report(n, GatherModuleMetadata_INVALID_DECLARE_MODULE_ID_CALL)
	}
}

func (f *moduleMetadataFinder) recordToggleRequire(t *NodeTraversal, n rhino.Node, callParent rhino.Node) {
	// Track imports of *.toggles.ts, which are rewritten to $2etoggles.
	if dlhs, ok := ast.AsDestructuringLhs(callParent); ok {
		// const {TOGGLE_foo} = goog.require('foo$2etoggles');
		if objPat, ok := dlhs.Pattern().AsObjectPattern(); ok {
			for prop := range objPat.Props() {
				if key, ok := prop.AsStringKey(); ok {
					f.addToggle(t, n, key.Key())
				} else {
					t.Report(n, GatherModuleMetadata_INVALID_TOGGLE_USAGE, "must be destructured with string keys")
				}
			}
		} else if arrPat, ok := dlhs.Pattern().AsArrayPattern(); ok {
			for range arrPat.Elements() {
				t.Report(n, GatherModuleMetadata_INVALID_TOGGLE_USAGE, "must be destructured with string keys")
			}
		}
	} else if nameNode, ok := ast.AsName(callParent); ok {
		// const fooToggles = goog.require('foo$2etoggles');
		nameID := nodeStringID(nameNode.Node())
		nameVar := t.GetScope().GetVarByID(uint32(nameID))
		if f.toggleModules == nil {
			f.toggleModules = make(map[*Var]struct{})
			f.toggleModuleIDs = make(map[strpool.ID]struct{})
		}
		f.toggleModules[nameVar] = struct{}{}
		f.toggleModuleIDs[nameID] = struct{}{}
	} else if f.currentModule.metadata.moduleType != ModuleMetadataMap_ModuleType_GOOG_PROVIDE {
		// Side-effect toggle-module imports are not allowed, since they don't actually do
		// anything. We allow it in `goog.provide()` files because there's no other way to
		// import, and since toggle modules don't declare a legacy namespace, it's unusable
		// without a `goog.module.get()` (so we can catch the toggle use there, instead).
		t.Report(n, GatherModuleMetadata_INVALID_TOGGLE_USAGE, "import must be assigned")
	}
}

// addToggle records a toggle usage (either a destructured import or a
// property lookup on a module).
//
// Java: GatherModuleMetadata.Finder.addToggle(NodeTraversal, Node, String)
func (f *moduleMetadataFinder) addToggle(t *NodeTraversal, n rhino.Node, name string) {
	if strings.HasPrefix(name, gatherModuleMetadata_TOGGLE_NAME_PREFIX) {
		toggleName := name[len(gatherModuleMetadata_TOGGLE_NAME_PREFIX):]
		f.currentModule.metadata.readToggles.Add(toggleName)
	} else {
		t.Report(n, GatherModuleMetadata_INVALID_TOGGLE_USAGE, "all toggle names must start with `TOGGLE_`")
	}
}

// addNamespace adds the namespaces to the module and checks if the given
// Closure namespace is a duplicate or not.
//
// Java: GatherModuleMetadata.Finder.addNamespace(ModuleMetadataBuilder, ModuleType, String, NodeTraversal, Node)
func (f *moduleMetadataFinder) addNamespace(
	module *gatherModuleMetadataBuilder,
	moduleType ModuleMetadataMap_ModuleType,
	namespace string,
	t *NodeTraversal,
	n rhino.Node) {
	if moduleType == ModuleMetadataMap_ModuleType_GOOG_PROVIDE ||
		moduleType == ModuleMetadataMap_ModuleType_LEGACY_GOOG_MODULE {
		if !NodeUtil_IsValidQualifiedNameFeatureSet(f.langInFeatures, namespace) {
			f.compiler.Report(JSError_MakeWithNode(n, GatherModuleMetadata_INVALID_NAMESPACE_OR_MODULE_ID, namespace))
		}
	}
	if moduleType == ModuleMetadataMap_ModuleType_GOOG_MODULE ||
		moduleType == ModuleMetadataMap_ModuleType_LEGACY_GOOG_MODULE {
		// non-legacy goog.modules don't technically need to be valid qualified names
		if !GatherModuleMetadata_IsValidModuleId(namespace) {
			f.compiler.Report(JSError_MakeWithNode(n, GatherModuleMetadata_INVALID_NAMESPACE_OR_MODULE_ID, namespace))
		}
	}

	sameModDup := module.metadata.googNamespaces.Contains(namespace)
	sameModTyp := module.metadata.moduleType
	sameModSrc := "null"
	if sameModDup {
		if src, ok := t.SourceNameOk(); ok {
			sameModSrc = src
		}
	}
	f.currentModule.metadata.googNamespaces.Add(namespace)
	if f.onOp != nil {
		f.onOp(gatherMetaScriptOp{
			kind:       gatherMetaOpAddNamespace,
			moduleType: moduleType,
			namespace:  namespace,
			node:       n,
			sourceName: sameModSrc,
			sameModDup: sameModDup,
			sameModTyp: sameModTyp,
		})
		return
	}
	g := f.outer
	g.checkDuplicateNamespace(moduleType, namespace, n, sameModDup, sameModTyp, sameModSrc)
}

func (g *GatherModuleMetadata) checkDuplicateNamespace(
	moduleType ModuleMetadataMap_ModuleType,
	namespace string,
	n rhino.Node,
	sameModDup bool,
	sameModTyp ModuleMetadataMap_ModuleType,
	sameModSrc string,
) {
	var existingType ModuleMetadataMap_ModuleType
	hasExistingType := false
	existingFileSource := "null"
	if sameModDup {
		existingType = sameModTyp
		hasExistingType = true
		existingFileSource = sameModSrc
	} else if existingModule := g.modulesByGoogNamespace.Get(namespace); existingModule != nil {
		existingType = existingModule.ModuleType()
		hasExistingType = true
		if src, ok := existingModule.RootNode().SourceFileName(); ok {
			existingFileSource = src
		}
	}
	if hasExistingType {
		switch existingType {
		case ModuleMetadataMap_ModuleType_ES6_MODULE, ModuleMetadataMap_ModuleType_GOOG_MODULE, ModuleMetadataMap_ModuleType_LEGACY_GOOG_MODULE:
			diagnostic := ClosurePrimitiveErrors_DUPLICATE_MODULE
			if moduleType == ModuleMetadataMap_ModuleType_GOOG_PROVIDE {
				diagnostic = ClosurePrimitiveErrors_DUPLICATE_NAMESPACE_AND_MODULE
			}
			g.compiler.Report(JSError_MakeWithNode(n, diagnostic, namespace, existingFileSource))
			return
		case ModuleMetadataMap_ModuleType_GOOG_PROVIDE:
			diagnostic := ClosurePrimitiveErrors_DUPLICATE_NAMESPACE_AND_MODULE
			if moduleType == ModuleMetadataMap_ModuleType_GOOG_PROVIDE {
				diagnostic = ClosurePrimitiveErrors_DUPLICATE_NAMESPACE
			}
			g.compiler.Report(JSError_MakeWithNode(n, diagnostic, namespace, existingFileSource))
			return
		case ModuleMetadataMap_ModuleType_COMMON_JS, ModuleMetadataMap_ModuleType_SCRIPT:
		}
		panic(jutil.NewIllegalStateExceptionMsg("Unexpected module type: " + existingType.String()))
	}
}

type gatherMetaScriptItem struct {
	scopeRoot rhino.Node
	script    rhino.Node
	lazyScope *lazyGlobalScope
}

type gatherMetaWorkerState struct {
	rc     Compiler
	sc     *SyntacticScopeCreator
	finder moduleMetadataFinder
	curIdx int
}

func (g *GatherModuleMetadata) processParallel(externs rhino.Node, root rhino.Node) bool {
	compiler := g.compiler
	if compiler == nil || compiler.IsDebugLoggingEnabled() {
		return false
	}
	extCount, okExt := 0, true
	if !externs.IsNil() {
		extCount, okExt = countRootScripts(externs)
	}
	rootCount, okRoot := 0, true
	if !root.IsNil() {
		rootCount, okRoot = countRootScripts(root)
	}
	if !okExt || !okRoot {
		return false
	}
	total := extCount + rootCount
	numWorkers := parallelThreadCount(compiler, total)
	if numWorkers <= 1 {
		return false
	}

	prewarmCompilerForParallelReads(compiler)
	langInFeatures := compiler.GetOptions().GetLanguageIn().ToFeatureSet()
	var lazyExt, lazyRoot *lazyGlobalScope
	if extCount > 0 {
		lazyExt = newLazyGlobalScope(compiler, externs)
	}
	if rootCount > 0 {
		lazyRoot = newLazyGlobalScope(compiler, root)
	}
	items := make([]gatherMetaScriptItem, 0, total)
	if extCount > 0 {
		for s := externs.GetFirstChild(); !s.IsNil(); s = s.GetNext() {
			items = append(items, gatherMetaScriptItem{scopeRoot: externs, script: s, lazyScope: lazyExt})
		}
	}
	if rootCount > 0 {
		for s := root.GetFirstChild(); !s.IsNil(); s = s.GetNext() {
			items = append(items, gatherMetaScriptItem{scopeRoot: root, script: s, lazyScope: lazyRoot})
		}
	}

	eventsByItem := make([][]recordedCompilerEvent, len(items))
	opsByItem := make([][]gatherMetaScriptOp, len(items))

	panicsByItem, workerPanic := collectParallelWorkersDeferredPins(
		len(items),
		numWorkers,
		func(_ int, w *gatherMetaWorkerState) {
			initRecordingCompiler(&w.rc, compiler)
			w.sc = NewSyntacticScopeCreator(&w.rc)
			w.finder = moduleMetadataFinder{
				outer:          g,
				compiler:       &w.rc,
				langInFeatures: langInFeatures,
			}
			w.finder.onOp = func(op gatherMetaScriptOp) {
				op.eventIdx = len(w.rc.parallelRecorder.events)
				opsByItem[w.curIdx] = append(opsByItem[w.curIdx], op)
			}
		},
		func(w *gatherMetaWorkerState, idx int) {
			resetRecordingCompiler(&w.rc, compiler)
			w.curIdx = idx
			w.finder.currentModule = nil
			w.finder.parentModule = nil
			w.finder.loadModuleCall = rhino.NilNode
			if len(w.finder.toggleModuleIDs) > 0 {
				clear(w.finder.toggleModuleIDs)
				clear(w.finder.toggleModules)
			}
			item := items[idx]
			wt := acquireNodeTraversalDirect(&w.rc, &w.finder, w.sc, false)
			defer releaseNodeTraversal(wt)
			func() {
				defer func() {
					if r := recover(); r != nil {
						wt.ThrowUnexpectedException(jutil.AsThrowable(r))
					}
				}()
				wt.lazyGlobalScope = item.lazyScope
				wt.InitTraversal(item.scopeRoot)
				wt.currentNode = item.scopeRoot
				wt.PushScope(item.scopeRoot)
				wt.TraverseBranch(item.script, item.scopeRoot)
				wt.PopScope()
			}()
		},
		func(w *gatherMetaWorkerState, idx int) {
			eventsByItem[idx] = w.rc.takeEvents()
		},
		nil,
		false,
	)

	for i := range items {
		evs := eventsByItem[i]
		evPos := 0
		for _, op := range opsByItem[i] {
			if op.eventIdx > evPos && evPos < len(evs) {
				end := min(op.eventIdx, len(evs))
				replayCompilerEvents(compiler, evs[evPos:end])
				evPos = end
			}
			switch op.kind {
			case gatherMetaOpAddNamespace:
				g.checkDuplicateNamespace(op.moduleType, op.namespace, op.node, op.sameModDup, op.sameModTyp, op.sourceName)
			case gatherMetaOpLeaveModule:
				g.recordBuiltModule(op.module)
			}
		}
		if evPos < len(evs) {
			replayCompilerEvents(compiler, evs[evPos:])
		}
		if p := panicsByItem[i]; p != nil {
			rethrowParallelPanic(compiler, p)
		}
	}
	if workerPanic != nil {
		rethrowParallelPanic(compiler, workerPanic)
	}
	return true
}

// Java: GatherModuleMetadata.process(Node, Node)
func (g *GatherModuleMetadata) Process(externs rhino.Node, root rhino.Node) {
	if !g.processParallel(externs, root) {
		langInFeatures := g.compiler.GetOptions().GetLanguageIn().ToFeatureSet()
		sc := NewSyntacticScopeCreator(g.compiler)
		finder := moduleMetadataFinder{outer: g, compiler: g.compiler, langInFeatures: langInFeatures}
		(NodeTraversal_Builder{compiler: g.compiler, scopeCreator: sc, callback: &finder}).Traverse(externs)
		finder = moduleMetadataFinder{outer: g, compiler: g.compiler, langInFeatures: langInFeatures}
		(NodeTraversal_Builder{compiler: g.compiler, scopeCreator: sc, callback: &finder}).Traverse(root)
	}
	g.compiler.SetModuleMetadataMap(newModuleMetadataMapDirect(g.modulesByPath, g.modulesByGoogNamespace))
}

// Must match closure/base.js's goog.VALID_MODULE_RE_ & also validates that
// dotted segments are non-empty.
//
// Java: GatherModuleMetadata.isValidModuleId(String)
func GatherModuleMetadata_IsValidModuleId(id string) bool {
	n := len(id)
	if n == 0 || id[n-1] == '.' {
		return false
	}
	c0 := id[0]
	if !((c0 >= 'a' && c0 <= 'z') || (c0 >= 'A' && c0 <= 'Z') || c0 == '_' || c0 == '$') {
		return false
	}
	for i := 1; i < n; i++ {
		c := id[i]
		if c == '.' {
			if id[i-1] == '.' {
				return false
			}
			continue
		}
		if !((c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') || c == '_' || c == '$') {
			return false
		}
	}
	return true
}
