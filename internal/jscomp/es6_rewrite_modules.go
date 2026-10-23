// Port of com/google/javascript/jscomp/Es6RewriteModules.java @ closure-compiler v20261015

/*
 * Copyright 2014 The Closure Compiler Authors.
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
	"github.com/paullepoulpe/jsgo/internal/parser"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// Es6RewriteModules rewrites a ES6 module into a form that can be safely
// concatenated. Note that we treat a file as an ES6 module if it has at
// least one import or export statement.
//
// Also rewrites any goog.{require,requireType,forwardDeclare,goog.module.get}
// calls that are either in an ES module or of an ES module using
// goog.declareModuleId.
type Es6RewriteModules struct {
	compiler    AbstractCompiler
	astFactory  *AstFactory
	unknownType rhino.JSType

	preprocessorSymbolTable *PreprocessorSymbolTable

	// Local variable names that were goog.require'd to qualified name we need
	// to line.
	//
	// We need to inline all required names since there are certain
	// well-known Closure symbols (like goog.asserts) that later stages of the
	// compiler check for and cannot handle aliases.
	//
	// We use this to rewrite something like:
	//
	//	import {x} from '';
	//	const {assert} = goog.require('goog.asserts');
	//	assert(x);
	//
	// To:
	//
	//	import {x} from '';
	//	goog.asserts.assert(x);
	//
	// Because if we used an alias like below the assertion would not be
	// recognized:
	//
	//	import {x} from '';
	//	const {assert} = goog.asserts;
	//	assert(x);
	//
	// TODO(johnplaisted): This is actually incorrect if the require'd thing is
	// mutated. But we need it so that things like goog.asserts work. Mutated
	// closure symbols are a lot rarer than needing to use asserts and the
	// like. Until there's a better solution to finding aliases of well known
	// symbols we have to inline anything that is require'd.
	namesToInlineByAlias *jutil.LinkedHashMap[string, *ModuleRenaming_GlobalizedModuleName]

	typedefs *jutil.LinkedHashSet[string]

	moduleMetadataMap *ModuleMetadataMap
	moduleMap         *ModuleMap
	globalTypedScope  *TypedScope
	chunkOutputType   CompilerOptions_ChunkOutputType
}

var Es6RewriteModules_LHS_OF_GOOG_REQUIRE_MUST_BE_CONST = DiagnosticType_Error(
	"JSC_LHS_OF_GOOG_REQUIRE_MUST_BE_CONST",
	"The left side of a goog.require() or goog.requireType() "+
		"must use ''const'' (not ''let'' or ''var'')")

var Es6RewriteModules_REQUIRE_TYPE_FOR_ES6_SHOULD_BE_CONST = DiagnosticType_Error(
	"JSC_REQUIRE_TYPE_FOR_ES6_SHOULD_BE_CONST",
	"goog.requireType alias for ES6 module should be const.")

var Es6RewriteModules_FORWARD_DECLARE_FOR_ES6_SHOULD_BE_CONST = DiagnosticType_Error(
	"JSC_FORWARD_DECLARE_FOR_ES6_SHOULD_BE_CONST",
	"goog.forwardDeclare alias for ES6 module should be const.")

var Es6RewriteModules_SHOULD_IMPORT_ES6_MODULE = DiagnosticType_Warning(
	"JSC_SHOULD_IMPORT_ES6_MODULE",
	"ES6 modules should import other ES6 modules rather than goog.require them.")

// es6RewriteModules_DOT_SPLITTER_splitToList is
// Splitter.on(".").splitToList(name).
//
// Go-only: language-adapter — a Guava Splitter is a function here; "." is
// ASCII, so splitting the WTF-8 bytes splits the UTF-16 units alike.
func es6RewriteModules_DOT_SPLITTER_splitToList(name string) []string {
	return strings.Split(name, ".")
}

// es6RewriteModules_DOT_SPLITTER_limit2_splitToList is
// Splitter.on(".").limit(2).splitToList(name).
//
// Go-only: language-adapter — see es6RewriteModules_DOT_SPLITTER_splitToList.
func es6RewriteModules_DOT_SPLITTER_limit2_splitToList(name string) []string {
	return strings.SplitN(name, ".", 2)
}

// NewEs6RewriteModules creates a new Es6RewriteModules instance which can be
// used to rewrite ES6 modules to a concatenable form.
//
// Java: Es6RewriteModules(AbstractCompiler, ModuleMetadataMap, ModuleMap, PreprocessorSymbolTable, TypedScope)
func NewEs6RewriteModules(
	compiler AbstractCompiler,
	moduleMetadataMap *ModuleMetadataMap,
	moduleMap *ModuleMap,
	preprocessorSymbolTable *PreprocessorSymbolTable,
	globalTypedScope *TypedScope) *Es6RewriteModules {
	return NewEs6RewriteModulesWithChunkOutputType(
		compiler,
		moduleMetadataMap,
		moduleMap,
		preprocessorSymbolTable,
		globalTypedScope,
		CompilerOptions_ChunkOutputType_GLOBAL_NAMESPACE)
}

// NewEs6RewriteModulesWithChunkOutputType creates a new Es6RewriteModules
// instance which can be used to rewrite ES6 modules to a concatenable form.
//
// Java: Es6RewriteModules(AbstractCompiler, ModuleMetadataMap, ModuleMap, PreprocessorSymbolTable, TypedScope, ChunkOutputType)
func NewEs6RewriteModulesWithChunkOutputType(
	compiler AbstractCompiler,
	moduleMetadataMap *ModuleMetadataMap,
	moduleMap *ModuleMap,
	preprocessorSymbolTable *PreprocessorSymbolTable,
	globalTypedScope *TypedScope,
	chunkOutputType CompilerOptions_ChunkOutputType) *Es6RewriteModules {
	jutil.CheckNotNull(moduleMetadataMap)
	e := &Es6RewriteModules{}
	e.compiler = compiler
	e.astFactory = compiler.CreateAstFactory()
	e.moduleMetadataMap = moduleMetadataMap
	e.moduleMap = moduleMap
	e.preprocessorSymbolTable = preprocessorSymbolTable
	e.globalTypedScope = globalTypedScope
	e.unknownType = compiler.GetTypeRegistry().GetNativeType(rhino.JSTypeNative_UNKNOWN_TYPE)
	e.chunkOutputType = chunkOutputType
	return e
}

// Es6RewriteModules_IsEs6ModuleRoot returns whether or not the given script
// node represents an ES6 module file.
//
// Java: Es6RewriteModules.isEs6ModuleRoot(Node)
func Es6RewriteModules_IsEs6ModuleRoot(scriptNode *rhino.Node) bool {
	jutil.CheckArgumentMsg1(scriptNode.IsScript(), "%s", scriptNode)
	if scriptNode.GetBooleanProp(rhino.Node_GOOG_MODULE) {
		return false
	}
	return scriptNode.HasChildren() && scriptNode.GetFirstChild().IsModuleBody()
}

// Java: Es6RewriteModules.process(Node, Node)
func (e *Es6RewriteModules) Process(externs *rhino.Node, root *rhino.Node) {
	jutil.CheckArgumentMsg1(externs.IsRoot(), "%s", externs)
	jutil.CheckArgumentMsg1(root.IsRoot(), "%s", root)
	NodeTraversal_TraverseRoots(e.compiler, e, externs, root)
	// It is unusual to call this NodeUtil method instead of the TranspilationPasses one. This
	// pass is included in the {@code dependency_resolution} BUILD target and does not have access
	// to {@code TranspilationPasses}. Adding that dep produces a cycle in the BUILD dep graph.
	// Regular transpiler passes must use the {@code
	// TranspilationPasses.maybeMarkFeaturesAsTranspiledAway}
	NodeUtil_RemoveFeatureFromAllScripts(root, parser.FeatureSet_Feature_MODULES, e.compiler)
	// This pass may add getters properties on module objects.
	GatherGetterAndSetterProperties_Update(e.compiler, externs, root)
}

// Java: Es6RewriteModules.clearPerFileState()
func (e *Es6RewriteModules) ClearPerFileState() {
	e.typedefs = jutil.NewLinkedHashSet[string]()
	e.namesToInlineByAlias = jutil.NewLinkedHashMap[string, *ModuleRenaming_GlobalizedModuleName]()
}

var es6RewriteModules_GOOG_REQUIRE = rhino.QualifiedName_Of("goog.require")
var es6RewriteModules_GOOG_REQUIRETYPE = rhino.QualifiedName_Of("goog.requireType")
var es6RewriteModules_GOOG_MODULE_GET = rhino.QualifiedName_Of("goog.module.get")
var es6RewriteModules_GOOG_FORWARDDECLARE = rhino.QualifiedName_Of("goog.forwardDeclare")
var es6RewriteModules_GOOG_DECLAREMODULEID = rhino.QualifiedName_Of("goog.declareModuleId")

// Es6RewriteModules_RewriteRequiresForEs6Modules checks for goog.require,
// goog.requireType, goog.module.get and goog.forwardDeclare calls that are
// meant to import ES6 modules and rewrites them.
type Es6RewriteModules_RewriteRequiresForEs6Modules struct {
	NodeTraversal_AbstractPostOrderCallbackBase

	outer *Es6RewriteModules // the enclosing instance

	// An (s, old, new) entry indicates that occurrences of `old` in scope `s`
	// should be rewritten as `new`. This is used to rewrite namespaces that
	// appear in calls to goog.requireType and goog.forwardDeclare.
	renameTable *jutil.HashBasedTable[*rhino.Node, string, string]
}

// Java: Es6RewriteModules.RewriteRequiresForEs6Modules()
func NewEs6RewriteModules_RewriteRequiresForEs6Modules(outer *Es6RewriteModules) *Es6RewriteModules_RewriteRequiresForEs6Modules {
	return &Es6RewriteModules_RewriteRequiresForEs6Modules{outer: outer}
}

// Java: Es6RewriteModules.RewriteRequiresForEs6Modules.rewrite(Node)
func (r *Es6RewriteModules_RewriteRequiresForEs6Modules) Rewrite(scriptNode *rhino.Node) {
	r.renameTable = jutil.HashBasedTable_Create[*rhino.Node, string, string]()
	NodeTraversal_Traverse(r.outer.compiler, scriptNode, r)

	if !r.renameTable.IsEmpty() {
		NodeTraversal_Traverse(r.outer.compiler, scriptNode, NewEs6RenameTypeReferences(r.renameTable))
	}
}

// Java: Es6RewriteModules.RewriteRequiresForEs6Modules.visit(NodeTraversal, Node, Node)
func (r *Es6RewriteModules_RewriteRequiresForEs6Modules) Visit(t *NodeTraversal, n *rhino.Node, parent *rhino.Node) {
	e := r.outer
	if !n.IsCall() {
		return
	}

	callTarget := n.GetFirstChild()
	if !callTarget.IsGetProp() {
		return
	}

	isRequire := es6RewriteModules_GOOG_REQUIRE.Matches(callTarget)
	isRequireType := es6RewriteModules_GOOG_REQUIRETYPE.Matches(callTarget)
	isGet := es6RewriteModules_GOOG_MODULE_GET.Matches(callTarget)
	isForwardDeclare := es6RewriteModules_GOOG_FORWARDDECLARE.Matches(callTarget)

	if !isRequire && !isRequireType && !isGet && !isForwardDeclare {
		return
	}

	if !n.HasTwoChildren() || !n.GetLastChild().IsStringLit() {
		// Reported in CheckClosureImports, don't report here.

		return
	}

	name := n.GetLastChild().GetString()
	moduleMetadata := e.moduleMetadataMap.GetModulesByGoogNamespace().Get(name)

	if moduleMetadata == nil || !moduleMetadata.IsEs6Module() {
		return
	}

	// TODO(johnplaisted): Once we have an alternative to forwardDeclare / requireType that
	// doesn't require Closure Library warn about those too.
	// TODO(johnplaisted): Once we have import() support warn about goog.module.get.
	if isRequire {
		currentModuleMetadata :=
			e.moduleMetadataMap.GetModulesByPath().Get(t.GetInput().GetPath().String())
		if currentModuleMetadata != nil && currentModuleMetadata.IsEs6Module() {
			t.Report(n, Es6RewriteModules_SHOULD_IMPORT_ES6_MODULE)
		}
	}

	if isGet && t.InGlobalHoistScope() {
		return
	}

	statementNode := NodeUtil_GetEnclosingStatement(n)
	importHasAlias := NodeUtil_IsNameDeclaration(statementNode)
	if importHasAlias {
		if statementNode.GetFirstChild().IsDestructuringLhs() {
			if isForwardDeclare {
				// const {a, c:b} = goog.forwardDeclare('an.es6.namespace');
				t.Report(n, ClosurePrimitiveErrors_INVALID_DESTRUCTURING_FORWARD_DECLARE)
				return
			}
			if isRequireType {
				if !statementNode.IsConst() {
					t.Report(statementNode, Es6RewriteModules_REQUIRE_TYPE_FOR_ES6_SHOULD_BE_CONST)
					return
				}
				// const {a, c:b} = goog.requireType('an.es6.namespace');
				for child := statementNode.GetFirstFirstChild().GetFirstChild(); child != nil; child = child.GetNext() {
					jutil.CheckState(child.IsStringKey())
					jutil.CheckState(child.GetFirstChild().IsName())
					r.renameTable.Put(
						t.GetScopeRoot(),
						child.GetFirstChild().GetString(),
						ModuleRenaming_GetGlobalNameWithModuleMetadata(moduleMetadata, jutil.StringPtr(name)).
							Getprop(child.GetString()).
							Join())
				}
			} else {
				// Work around a bug in the type checker where destructuring can create
				// too many layers of aliases and confuse the type checker. b/112061124.

				// const {a, c:b} = goog.require('an.es6.namespace');
				// const a = module$es6.a;
				// const b = module$es6.c;
				for child := statementNode.GetFirstFirstChild().GetFirstChild(); child != nil; child = child.GetNext() {
					jutil.CheckState(child.IsStringKey())
					jutil.CheckState(child.GetFirstChild().IsName())
					globalName :=
						e.GetGlobalNameAndType(
							moduleMetadata, jutil.StringPtr(name) /* isFromMissingModuleOrProvide= */, false).
							Getprop(child.GetString())
					constNode :=
						e.astFactory.CreateSingleConstNameDeclaration(
							child.GetFirstChild().GetString(), globalName.ToQname(e.astFactory))
					constNode.SrcrefTree(child)
					constNode.InsertBefore(statementNode)
				}
			}
			statementNode.Detach()
			t.ReportCodeChange()
		} else {
			if isForwardDeclare || isRequireType {
				if !statementNode.IsConst() {
					diagnostic := Es6RewriteModules_REQUIRE_TYPE_FOR_ES6_SHOULD_BE_CONST
					if isForwardDeclare {
						diagnostic = Es6RewriteModules_FORWARD_DECLARE_FOR_ES6_SHOULD_BE_CONST
					}
					t.Report(statementNode, diagnostic)
					return
				}
				// const namespace = goog.forwardDeclare('an.es6.namespace');
				// const namespace = goog.requireType('an.es6.namespace');
				r.renameTable.Put(
					t.GetScopeRoot(),
					statementNode.GetFirstChild().GetString(),
					ModuleRenaming_GetGlobalNameWithModuleMetadata(moduleMetadata, jutil.StringPtr(name)).Join())
				statementNode.Detach()
				t.ReportCodeChange()
			} else {
				// const module = goog.require('an.es6.namespace');
				// const module = module$es6;
				n.ReplaceWith(
					e.astFactory.
						CreateName(
							ModuleRenaming_GetGlobalNameWithModuleMetadata(moduleMetadata, jutil.StringPtr(name)).GetRoot(), AstFactory_Type_(n)).
						SrcrefTree(n))
				t.ReportCodeChange()
			}
		}
	} else {
		if isForwardDeclare || isRequireType {
			// goog.forwardDeclare('an.es6.namespace')
			// goog.requireType('an.es6.namespace')
			r.renameTable.Put(
				t.GetScopeRoot(), name, ModuleRenaming_GetGlobalNameWithModuleMetadata(moduleMetadata, jutil.StringPtr(name)).Join())
			statementNode.Detach()
		} else {
			// goog.require('an.es6.namespace')
			if statementNode.IsExprResult() && statementNode.GetFirstChild() == n {
				statementNode.Detach()
			} else {
				n.ReplaceWith(
					e.astFactory.
						CreateName(
							ModuleRenaming_GetGlobalNameWithModuleMetadata(moduleMetadata, jutil.StringPtr(name)).GetRoot(), AstFactory_Type_(n)).
						SrcrefTree(n))
			}
		}
		t.ReportCodeChange()
	}
}

// Java: Es6RewriteModules.shouldTraverse(NodeTraversal, Node, Node)
func (e *Es6RewriteModules) ShouldTraverse(nodeTraversal *NodeTraversal, n *rhino.Node, parent *rhino.Node) bool {
	if n.IsScript() {
		// Trigger creation of the global scope before inserting any synthetic code.
		nodeTraversal.GetScope()
		NewEs6RewriteModules_RewriteRequiresForEs6Modules(e).Rewrite(n)
		if Es6RewriteModules_IsEs6ModuleRoot(n) {
			e.ClearPerFileState()
		} else {
			return false
		}
	}
	return true
}

// Java: Es6RewriteModules.visit(NodeTraversal, Node, Node)
func (e *Es6RewriteModules) Visit(t *NodeTraversal, n *rhino.Node, parent *rhino.Node) {
	if n.IsImport() {
		e.MaybeWarnExternModule(t, n, parent)
		e.VisitImport(t, n, parent)
	} else if n.IsExport() {
		e.MaybeWarnExternModule(t, n, parent)
		e.VisitExport(t, n, parent)
	} else if n.IsScript() {
		e.VisitScript(t, n)
	} else if n.IsCall() {
		if es6RewriteModules_GOOG_DECLAREMODULEID.Matches(n.GetFirstChild()) {
			n.GetParent().Detach()
		}
	} else if n.IsImportMeta() {
		if e.chunkOutputType != CompilerOptions_ChunkOutputType_ES_MODULES {
			// Since import.meta cannot be transpiled, we are only supporting it when the output format
			// is a module. The default output format type is just script.
			t.Report(
				n,
				TranspilationUtil_CANNOT_CONVERT,
				"import.meta. Use --chunk_output_type=ES_MODULES to allow passthrough support.")
		}
	}
}

// Java: Es6RewriteModules.maybeWarnExternModule(NodeTraversal, Node, Node)
func (e *Es6RewriteModules) MaybeWarnExternModule(t *NodeTraversal, n *rhino.Node, parent *rhino.Node) {
	jutil.CheckState(parent.IsModuleBody())
	if parent.IsFromExterns() && !NodeUtil_IsFromTypeSummary(parent.GetParent()) {
		t.Report(n, TranspilationUtil_CANNOT_CONVERT_YET, "ES6 modules in externs")
	}
}

// Java: Es6RewriteModules.visitImport(NodeTraversal, Node, Node)
func (e *Es6RewriteModules) VisitImport(t *NodeTraversal, importDecl *rhino.Node, parent *rhino.Node) {
	jutil.CheckArgumentMsg1(parent.IsModuleBody(), "%s", parent)
	importName := importDecl.GetLastChild().GetString()
	isNamespaceImport := strings.HasPrefix(importName, "goog:")
	if isNamespaceImport {
		// Allow importing Closure namespace objects (e.g. from goog.provide or goog.module) as
		//   import ... from 'goog:my.ns.Object'.
		namespace := jutil.StringSubstringFrom(importName, jutil.StringLength("goog:"))
		m := e.moduleMetadataMap.GetModulesByGoogNamespace().Get(namespace)

		if m == nil {
			t.Report(importDecl, ClosurePrimitiveErrors_MISSING_MODULE_OR_PROVIDE, namespace)
		} else {
			jutil.CheckState(m.IsEs6Module() || m.IsGoogModule() || m.IsGoogProvide())
		}
	} else {
		modulePath :=
			t.GetInput().
				GetPath().
				ResolveJsModule(
					importName,
					processCommonJSModules_sourceFileName(importDecl),
					importDecl.GetLineno(),
					importDecl.GetCharno())
		if modulePath == nil {
			// The module loader issues an error
			// Fall back to assuming the module is a file path
			modulePath = t.GetInput().GetPath().ResolveModuleAsPath(importName)
		}

		e.MaybeAddImportedFileReferenceToSymbolTable(importDecl.GetLastChild(), modulePath.String())
		// TODO(johnplaisted): Use ModuleMetadata to ensure the path required is CommonJs or ES6 and
		// if not give a better error.
	}

	for child := importDecl.GetFirstChild(); child != nil; child = child.GetNext() {
		if child.IsImportSpecs() {
			for grandChild := child.GetFirstChild(); grandChild != nil; grandChild = grandChild.GetNext() {
				e.MaybeAddAliasToSymbolTable(grandChild.GetFirstChild(), jutil.ToJavaString(t.GetSourceName()))
				jutil.CheckState(grandChild.HasTwoChildren())
			}
		} else if child.IsImportStar() {
			// import * as ns from "mod"
			e.MaybeAddAliasToSymbolTable(child, jutil.ToJavaString(t.GetSourceName()))
		}
	}

	importDecl.Detach()
	t.ReportCodeChange()
}

// Java: Es6RewriteModules.visitExport(NodeTraversal, Node, Node)
func (e *Es6RewriteModules) VisitExport(t *NodeTraversal, export *rhino.Node, parent *rhino.Node) {
	jutil.CheckArgumentMsg1(parent.IsModuleBody(), "%s", parent)
	if export.GetBooleanProp(rhino.Node_EXPORT_DEFAULT) {
		// export default

		// If the thing being exported is a class or function that has a name,
		// extract it from the export statement, so that it can be referenced
		// from within the module.
		//
		//   export default class X {} -> class X {}; ... moduleName.default = X;
		//   export default function X() {} -> function X() {}; ... moduleName.default = X;
		//
		// Otherwise, create a local variable for it and export that.
		//
		//   export default 'someExpression'
		//     ->
		//   var $jscompDefaultExport = 'someExpression';
		//   ...
		//   moduleName.default = $jscompDefaultExport;
		child := export.GetFirstChild()
		var name *string

		if child.IsFunction() || child.IsClass() {
			name = NodeUtil_GetName(child)
		}

		if name != nil {
			decl := child.Detach()
			export.ReplaceWith(decl)
		} else {
			var_ :=
				e.astFactory.CreateSingleVarNameDeclarationWithValue(
					ModuleRenaming_DEFAULT_EXPORT_VAR_PREFIX, export.RemoveFirstChild())
			var_.SetJSDocInfo(child.GetJSDocInfo())
			child.SetJSDocInfo(nil)
			var_.SrcrefTreeIfMissing(export)
			export.ReplaceWith(var_)
		}
		t.ReportCodeChange()
	} else if export.GetBooleanProp(rhino.Node_EXPORT_ALL_FROM) ||
		export.HasTwoChildren() ||
		export.GetFirstChild().IsExportSpecs() {
		//   export * from 'moduleIdentifier';
		//   export {x, y as z} from 'moduleIdentifier';
		//   export {Foo};
		export.Detach()
		t.ReportCodeChange()
	} else {
		e.VisitExportDeclaration(t, export)
	}
}

// Java: Es6RewriteModules.visitExportNameDeclaration(Node)
func (e *Es6RewriteModules) VisitExportNameDeclaration(declaration *rhino.Node) {
	if declaration.GetJSDocInfo() != nil && declaration.GetJSDocInfo().HasTypedefType() {
		//    export var Foo;
		//    export let {a, b:[c,d]} = {};
		NodeUtil_VisitLhsNodesInNode(declaration, func(lhs *rhino.Node) { e.typedefs.Add(lhs.GetString()) })
	}
}

// Java: Es6RewriteModules.visitExportDeclaration(NodeTraversal, Node)
func (e *Es6RewriteModules) VisitExportDeclaration(t *NodeTraversal, export *rhino.Node) {
	//    export var Foo;
	//    export function Foo() {}
	// etc.
	declaration := export.GetFirstChild()

	if NodeUtil_IsNameDeclaration(declaration) {
		e.VisitExportNameDeclaration(declaration)
	}

	export.ReplaceWith(declaration.Detach())
	t.ReportCodeChange()
}

// Java: Es6RewriteModules.inlineModuleToGlobalScope(Node)
func (e *Es6RewriteModules) InlineModuleToGlobalScope(moduleNode *rhino.Node) {
	jutil.CheckState(moduleNode.IsModuleBody())
	scriptNode := moduleNode.GetParent()
	moduleNode.Detach()
	scriptNode.AddChildrenToFront(moduleNode.RemoveChildren())
}

// Java: Es6RewriteModules.visitScript(NodeTraversal, Node)
func (e *Es6RewriteModules) VisitScript(t *NodeTraversal, script *rhino.Node) {
	moduleBody := script.GetFirstChild()
	// TypedScopeCreator sets the module object type on the MODULE_BODY during type checking.
	moduleObjectType := AstFactory_Type_(moduleBody)
	e.InlineModuleToGlobalScope(moduleBody)

	ClosureRewriteModule_CheckAndSetStrictModeDirective(t, script)

	thisModule := e.moduleMap.GetModuleByPath(t.GetInput().GetPath())
	qualifiedName :=
		ModuleRenaming_GetGlobalNameWithModuleMetadata(thisModule.Metadata() /* googNamespace= */, nil)
	jutil.CheckStateMsg1(qualifiedName.IsSimple(), "Unexpected qualified name %s", qualifiedName)
	moduleName := qualifiedName.GetRoot()

	moduleVar := e.CreateExportsObject(moduleName, t, script, moduleObjectType)

	// Rename vars to not conflict in global scope.
	NodeTraversal_Traverse(e.compiler, script, NewEs6RewriteModules_RenameGlobalVars(e, thisModule))

	// Rename the exports object to something we can reference later.
	moduleVar.GetFirstChild().SetString(moduleName)
	moduleVar.MakeNonIndexableRecursive()
	e.DeclareGlobalVariable(moduleVar.GetFirstChild(), t)

	// rewriteRequires is here (rather than being part of the main visit() method, because we only
	// want to rewrite the requires if this is an ES6 module. Note that we also want to do this
	// AFTER renaming all module scoped vars in the event that something that is goog.require'd is
	// a global, unqualified name (e.g. if "goog.provide('foo')" exists, we don't want to rewrite
	// "const foo = goog.require('foo')" to "const foo = foo". If we rewrite our module scoped names
	// first then we'll rewrite to "const foo$module$fudge = goog.require('foo')", then to
	// "const foo$module$fudge = foo".
	e.RewriteRequires(script)

	t.ReportCodeChange()
}

// Java: Es6RewriteModules.createExportsObject(String, NodeTraversal, Node, AstFactory.Type)
func (e *Es6RewriteModules) CreateExportsObject(
	moduleName string, t *NodeTraversal, script *rhino.Node, moduleObjectType AstFactory_Type) *rhino.Node {
	moduleObject := e.astFactory.CreateObjectLitWithType(moduleObjectType)
	// Going to get renamed by RenameGlobalVars, so the name we choose here doesn't matter as long
	// as it doesn't collide with an existing variable. (We can't use `moduleName` since then
	// RenameGlobalVars will rename all references to `moduleName` incorrectly). We'll fix the name
	// in visitScript after the global renaming to ensure it has a name that is deterministic from
	// the path.
	//
	// So after this method we'll have:
	// var $jscomp$tmp$exports$module$name = {};
	// module$name.exportName = localName;
	//
	// After RenameGlobalVars:
	// var $jscomp$tmp$exports$module$nameglobalized = {};
	// module$name.exportName = localName$globalized;
	//
	// After visitScript:
	// var module$name = {};
	// module$name.exportName = localName$globalized;
	moduleVar :=
		e.astFactory.CreateSingleVarNameDeclarationWithValue("$jscomp$tmp$exports$module$name", moduleObject)
	moduleVar.GetFirstChild().PutBooleanProp(rhino.Node_MODULE_EXPORT, true)
	// TODO(b/144593112): Stop adding JSDoc when this pass moves to always be after typechecking.
	infoBuilder := rhino.JSDocInfo_Builder_()
	infoBuilder.RecordConstancy()
	moduleVar.SetJSDocInfo(infoBuilder.Build())
	moduleVar.GetFirstChild().SetDeclaredConstantVar(true)
	script.AddChildToBack(moduleVar.SrcrefTreeIfMissing(script))

	thisModule := e.moduleMap.GetModuleByPath(t.GetInput().GetPath())

	for exportedName, binding := range thisModule.Namespace().All() {
		nodeForSourceInfo := binding.SourceNode()
		mutated := binding.IsMutated()
		boundVariableQualifiedName := ModuleRenaming_GetGlobalNameWithBinding(binding)
		jutil.CheckStateMsg1(
			boundVariableQualifiedName.IsSimple(),
			"unexpected qualified name: %s",
			boundVariableQualifiedName)
		boundVariableName := boundVariableQualifiedName.GetRoot()

		getProp :=
			e.astFactory.CreateGetPropWithoutColor(
				e.astFactory.CreateName(moduleName, moduleObjectType), exportedName)
		getProp.PutBooleanProp(rhino.Node_MODULE_EXPORT, true)

		if e.typedefs.Contains(exportedName) {
			// /** @typedef {foo} */
			// moduleName.foo;
			builder := rhino.JSDocInfo_Builder_().ParseDocumentation()
			typeExpr :=
				rhino.NewJSTypeExpression(
					e.astFactory.CreateString(exportedName).Srcref(nodeForSourceInfo),
					jutil.ToJavaString(script.GetSourceFileName()))
			builder.RecordTypedef(typeExpr)
			info := builder.Build()
			getProp.SetJSDocInfo(info)
			exprResult := e.astFactory.ExprResult(getProp).SrcrefTreeIfMissing(nodeForSourceInfo)
			script.AddChildToBack(exprResult)
		} else if mutated {
			globalExportName := e.astFactory.CreateName(boundVariableName, AstFactory_Type_(getProp))
			e.AddGetterExport(script, nodeForSourceInfo, moduleObject, exportedName, globalExportName)
			NodeUtil_AddFeatureToScript(t.GetCurrentScript(), parser.FeatureSet_Feature_GETTER, e.compiler)
		} else {
			// Avoid the extra complexity of using getters when the property isn't mutated.
			// exports.foo = foo;
			assign :=
				e.astFactory.CreateAssign(
					getProp, e.astFactory.CreateName(boundVariableName, AstFactory_Type_(getProp)))
			// TODO(b/144593112): Stop adding JSDoc when this pass moves to always be after typechecking
			builder := rhino.JSDocInfo_Builder_().ParseDocumentation()
			builder.RecordConstancy()
			info := builder.Build()
			assign.SetJSDocInfo(info)
			script.AddChildToBack(e.astFactory.ExprResult(assign).SrcrefTreeIfMissing(nodeForSourceInfo))
		}
	}

	return moduleVar
}

// Java: Es6RewriteModules.addGetterExport(Node, Node, Node, String, Node)
func (e *Es6RewriteModules) AddGetterExport(
	script *rhino.Node, forSourceInfo *rhino.Node, objLit *rhino.Node, exportedName string, value *rhino.Node) {
	getter := e.astFactory.CreateGetterDef(exportedName, value)
	getter.PutBooleanProp(rhino.Node_MODULE_EXPORT, true)
	objLit.AddChildToBack(getter)

	if !e.astFactory.IsAddingTypes() {
		// TODO(b/143904518): Remove this code when this pass is permanently moved after type checking
		// Type checker doesn't infer getters so mark the return as unknown.
		// { /** @return {?} */ get foo() { return foo; } }
		builder := rhino.JSDocInfo_Builder_().ParseDocumentation()
		builder.RecordReturnType(
			rhino.NewJSTypeExpression(
				rhino.NewNode(rhino.Token_QMARK).Srcref(forSourceInfo), jutil.ToJavaString(script.GetSourceFileName())))
		getter.SetJSDocInfo(builder.Build())
	} else {
		// For a property typed as number, synthesize a type `function(): number`.
		getter.SetJSType(e.compiler.GetTypeRegistry().CreateFunctionType(value.GetJSType()))
	}

	getter.SrcrefTreeIfMissing(forSourceInfo)
	e.compiler.ReportChangeToEnclosingScope(getter.GetFirstChild().GetLastChild())
	e.compiler.ReportChangeToEnclosingScope(getter)
}

// Java: Es6RewriteModules.rewriteRequires(Node)
func (e *Es6RewriteModules) RewriteRequires(script *rhino.Node) {
	NodeTraversal_Builder_().
		SetCompiler(e.compiler).
		SetCallbackFunc(
			func(t *NodeTraversal, n *rhino.Node, parent *rhino.Node) {
				if n.IsCall() {
					fn := n.GetFirstChild()
					if es6RewriteModules_GOOG_REQUIRE.Matches(fn) || es6RewriteModules_GOOG_REQUIRETYPE.Matches(fn) {
						// TODO(tjgq): This will rewrite both type references and code references. For
						// goog.requireType, the latter are potentially broken because the symbols aren't
						// guaranteed to be available at run time. A separate pass needs to be added to
						// detect these incorrect uses of goog.requireType.
						e.VisitRequireOrGet(t, n, parent /* isRequire= */, true)
					} else if es6RewriteModules_GOOG_MODULE_GET.Matches(fn) {
						e.VisitGoogModuleGet(t, n, parent)
					}
				}
			}).
		Traverse(script)
	NodeTraversal_Builder_().
		SetCompiler(e.compiler).
		SetCallbackFunc(
			func(t *NodeTraversal, n *rhino.Node, parent *rhino.Node) {
				info := n.GetJSDocInfo()
				if info != nil {
					for it := info.TypeExprIter(); it.Next(); {
						e.InlineAliasedTypes(t, it.Node())
					}
				}

				if n.IsName() && e.namesToInlineByAlias.ContainsKey(n.GetString()) {
					v := t.GetScope().GetVar(n.GetString())
					if v == nil || v.GetNameNode() != n {
						replacementName := e.namesToInlineByAlias.Get(n.GetString())
						replacement := replacementName.ToQname(e.astFactory).SrcrefTree(n)
						n.ReplaceWith(replacement)
					}
				}
			}).
		Traverse(script)
}

// Java: Es6RewriteModules.inlineAliasedTypes(NodeTraversal, Node)
func (e *Es6RewriteModules) InlineAliasedTypes(t *NodeTraversal, typeNode *rhino.Node) {
	if typeNode.IsStringLit() {
		name := typeNode.GetString()
		split := es6RewriteModules_DOT_SPLITTER_limit2_splitToList(name)

		// We've already removed the alias.
		if t.GetScope().GetVar(split[0]) == nil {
			replacement := e.namesToInlineByAlias.Get(split[0])
			if replacement != nil {
				rest := ""
				if len(split) == 2 {
					rest = "." + split[1]
				}
				typeNode.SetOriginalName(jutil.StringPtr(name))
				typeNode.SetString(replacement.AliasName().Join() + rest)
				t.ReportCodeChange()
			}
		}
	}
	for child := typeNode.GetFirstChild(); child != nil; child = child.GetNext() {
		e.InlineAliasedTypes(t, child)
	}
}

// Java: Es6RewriteModules.visitGoogModuleGet(NodeTraversal, Node, Node)
func (e *Es6RewriteModules) VisitGoogModuleGet(t *NodeTraversal, getCall *rhino.Node, parent *rhino.Node) {
	if !getCall.HasTwoChildren() || !getCall.GetLastChild().IsStringLit() {
		t.Report(getCall, ClosurePrimitiveErrors_INVALID_GET_NAMESPACE)
		return
	}

	// Module has already been turned into a script at this point.
	if t.InGlobalHoistScope() {
		t.Report(getCall, ClosurePrimitiveErrors_MODULE_USES_GOOG_MODULE_GET)
		return
	}

	e.VisitRequireOrGet(t, getCall, parent /* isRequire= */, false)
}

// GetFallbackMetadataForNamespace gets some made-up metadata for the given
// Closure namespace.
//
// This is used when the namespace is not part of the input so that this
// pass can be fault tolerant and still rewrite to something. Some tools
// don't care about rewriting correctly and just want the type information
// of this module (e.g. clutz).
//
// Java: Es6RewriteModules.getFallbackMetadataForNamespace(String)
func (e *Es6RewriteModules) GetFallbackMetadataForNamespace(namespace string) *ModuleMetadataMap_ModuleMetadata {
	// Assume a provide'd file to be consistent with goog.module rewriting.
	builder :=
		ModuleMetadataMap_ModuleMetadata_Builder_().
			ModuleType(ModuleMetadataMap_ModuleType_GOOG_PROVIDE).
			UsesClosure(true).
			IsTestOnly(false)
	builder.GoogNamespacesBuilder().Add(namespace)
	return builder.Build()
}

// Java: Es6RewriteModules.visitRequireOrGet(NodeTraversal, Node, Node, boolean)
func (e *Es6RewriteModules) VisitRequireOrGet(
	t *NodeTraversal, requireCall *rhino.Node, parent *rhino.Node, isRequire bool) {
	if !requireCall.HasTwoChildren() || !requireCall.GetLastChild().IsStringLit() {
		t.Report(requireCall, ClosurePrimitiveErrors_INVALID_REQUIRE_NAMESPACE)
		return
	}

	// Module has already been turned into a script at this point.
	if isRequire && !t.GetScope().IsGlobal() {
		t.Report(requireCall, ClosurePrimitiveErrors_INVALID_CLOSURE_CALL_SCOPE_ERROR)
		return
	}

	namespace := requireCall.GetLastChild().GetString()

	isStoredInDeclaration := NodeUtil_IsDeclaration(parent.GetParent())

	if isStoredInDeclaration && !parent.GetParent().IsConst() {
		e.compiler.Report(JSError_MakeWithNode(parent.GetParent(), Es6RewriteModules_LHS_OF_GOOG_REQUIRE_MUST_BE_CONST))
	}

	m := e.moduleMetadataMap.GetModulesByGoogNamespace().Get(namespace)

	isFromFallbackMetadata := m == nil
	if isFromFallbackMetadata {
		t.Report(requireCall, ClosurePrimitiveErrors_MISSING_MODULE_OR_PROVIDE, namespace)
		m = e.GetFallbackMetadataForNamespace(namespace)
	}

	if isStoredInDeclaration {
		if isRequire {
			var toDetach *rhino.Node

			if parent.IsDestructuringLhs() {
				jutil.CheckState(parent.GetFirstChild().IsObjectPattern())
				toDetach = parent.GetParent()
				for child := parent.GetFirstFirstChild(); child != nil; child = child.GetNext() {
					jutil.CheckStateMsg1(child.IsStringKey() && child.GetFirstChild().IsName(), "%s", child)
					rep :=
						e.GetGlobalNameAndType(m, jutil.StringPtr(namespace), isFromFallbackMetadata).
							Getprop(child.GetString())
					e.namesToInlineByAlias.Put(child.GetFirstChild().GetString(), rep)
				}
			} else if parent.IsName() {
				alias := e.GetGlobalNameAndType(m, jutil.StringPtr(namespace), isFromFallbackMetadata)
				e.namesToInlineByAlias.Put(parent.GetString(), alias)

				toDetach = parent.GetParent()
			} else {
				jutil.CheckState(parent.IsExprResult())
				toDetach = parent
			}
			toDetach.Detach()
		} else {
			name := e.GetGlobalNameAndType(m, jutil.StringPtr(namespace), isFromFallbackMetadata)
			replacement := name.ToQname(e.astFactory).SrcrefTree(requireCall)
			requireCall.ReplaceWith(replacement)
		}
	} else {
		jutil.CheckState(requireCall.GetParent().IsExprResult())
		requireCall.GetParent().Detach()
	}
}

// GetGlobalNameAndType looks up information about the globalized name and
// type of a given module: metadata (required) is the metadata for the
// module or provide being imported; googNamespace is optional;
// isFromMissingModuleOrProvide whether the metadata is synthesized fallback
// metadata.
//
// Java: Es6RewriteModules.getGlobalNameAndType(ModuleMetadata, String, boolean)
func (e *Es6RewriteModules) GetGlobalNameAndType(
	metadata *ModuleMetadataMap_ModuleMetadata,
	googNamespace *string,
	isFromMissingModuleOrProvide bool) *ModuleRenaming_GlobalizedModuleName {
	if isFromMissingModuleOrProvide {
		// The missing namespace presumably does not have a corresponding type defined in the scope.
		// Use the unknownType instead of asking ModuleRenaming to look up the type.
		globalName := ModuleRenaming_GetGlobalNameWithModuleMetadata(metadata, googNamespace)
		return ModuleRenaming_GlobalizedModuleName_Create(globalName, e.unknownType)
	}
	return ModuleRenaming_GlobalizedModuleName_CreateWithModuleMetadata(metadata, googNamespace, e.globalTypedScope)
}

// Es6RewriteModules_RenameGlobalVars traverses a node tree and
//
//  1. Appends a suffix to all global variable names defined in this module.
//  2. Changes references to imported values to access the exported variable.
type Es6RewriteModules_RenameGlobalVars struct {
	NodeTraversal_AbstractPostOrderCallbackBase

	outer *Es6RewriteModules // the enclosing instance

	thisModule *Module
}

// Java: Es6RewriteModules.RenameGlobalVars(Module)
func NewEs6RewriteModules_RenameGlobalVars(outer *Es6RewriteModules, thisModule *Module) *Es6RewriteModules_RenameGlobalVars {
	r := &Es6RewriteModules_RenameGlobalVars{outer: outer}
	r.thisModule = thisModule
	return r
}

// Java: Es6RewriteModules.RenameGlobalVars.visit(NodeTraversal, Node, Node)
func (r *Es6RewriteModules_RenameGlobalVars) Visit(t *NodeTraversal, n *rhino.Node, parent *rhino.Node) {
	info := n.GetJSDocInfo()
	if info != nil {
		for it := info.TypeExprIter(); it.Next(); {
			r.FixTypeNode(t, it.Node())
		}
	}

	if n.IsName() {
		name := n.GetString()

		var_ := t.GetScope().GetVar(name)
		if var_ != nil && var_.IsGlobal() {
			// Avoid polluting the global namespace.
			newName :=
				ModuleRenaming_GetGlobalNameOfEsModuleLocalVariable(r.thisModule.Metadata(), name).
					Join()
			n.SetString(newName)
			n.SetOriginalName(jutil.StringPtr(name))
			t.ReportCodeChangeWithNode(n)
			if NodeUtil_IsDeclarationLValue(n) {
				r.outer.DeclareGlobalVariable(n, t)
			}
		} else if var_ == nil && r.thisModule.BoundNames().ContainsKey(name) {
			// Imports have been detached, so they won't show up in scope. Thus if we have a variable
			// not in scope that shares the name of an import it is the import.
			r.outer.MaybeAddAliasToSymbolTable(n, jutil.ToJavaString(t.GetSourceName()))
			binding := r.thisModule.BoundNames().Get(name)

			replacement := r.Replace(n, binding)

			// `n.x()` may become `foo()`
			if replacement.IsName() &&
				parent.IsCall() &&
				parent.GetFirstChild() == n &&
				parent.GetBooleanProp(rhino.Node_FREE_CALL) {
				parent.PutBooleanProp(rhino.Node_FREE_CALL, true)
			}

			if NodeUtil_IsDeclarationLValue(n) {
				r.outer.DeclareGlobalVariable(n, t)
			}
			t.ReportCodeChange()
		}
	}
}

// Replace replaces the reference to a given binding. See
// ModuleRenaming#getGlobalNameForJsDoc(ModuleMap, Binding, List) for a JS
// Doc version.
//
// For example:
//
//	// bar
//	export let baz = {qux: 0};
//
//	// foo
//	import * as bar from 'bar';
//	export {bar};
//
//	import * as foo from 'foo';
//	use(foo.bar.baz.qux);
//
// Should call this method with the binding and node for foo. In this
// example any of these properties could also be modules. This method will
// replace as much as the GETPROP as it can with module exported variables.
// Meaning in the above example this would return something like
// "baz$$module$bar.qux", whereas if this method were called for just
// "foo.bar" it would return "module$bar", as it refers to a module object
// itself.
//
// n is the node to replace; binding the binding nameNode is a reference to.
//
// Java: Es6RewriteModules.RenameGlobalVars.replace(Node, Binding)
func (r *Es6RewriteModules_RenameGlobalVars) Replace(n *rhino.Node, binding *Binding) *rhino.Node {
	e := r.outer
	jutil.CheckState(n.IsName())

	for binding.IsModuleNamespace() &&
		binding.Metadata().IsEs6Module() &&
		n.GetParent().IsGetProp() {
		propertyName := n.GetParent().GetString()
		m := e.moduleMap.GetModuleByPath(binding.Metadata().Path())
		if m.Namespace().ContainsKey(propertyName) {
			binding = m.Namespace().Get(propertyName)
			n = n.GetParent()
		} else {
			// This means someone referenced an invalid export on a module object. This should be an
			// error, so just rewrite and let the type checker complain later. It isn't a super clear
			// error, but we're working on type checking modules soon.
			break
		}
	}

	globalName := ModuleRenaming_GetGlobalNameWithBinding(binding)
	var newNode *rhino.Node
	if !globalName.IsSimple() {
		root := globalName.GetRoot()
		// we might encounter a name not in the global scope when requiring a missing symbol.
		if e.globalTypedScope != nil && e.globalTypedScope.HasSlot(root) {
			newNode = e.astFactory.CreateQNameUsingJSTypeInfo(e.globalTypedScope, globalName.Join())
		} else {
			newNode = e.astFactory.CreateQNameWithUnknownType(globalName.Join())
		}
	} else {
		// Because this pass does not update the global scope with injected names, t.getScope()
		// will not contain a declaration for this global name. Fortunately, we already have the
		// JSType on the existing node to pass to AstFactory.
		newNode = e.astFactory.CreateName(globalName.GetRoot(), AstFactory_Type_(n))
	}

	// For kythe: the new node only represents the last name it replaced, not all the names.
	// e.g. if we rewrite `a.b.c.d.e` to `x.d.e`, then `x` should map to `c`, not `a.b.c`.
	n.ReplaceWith(newNode)
	newNode.SrcrefTree(n)
	newNode.SetOriginalName(jutil.StringPtr(n.GetString()))
	return newNode
}

// FixTypeNode replaces type name references. Change short names to fully
// qualified names with namespace prefixes. Eg: {Foo} becomes
// {module$test.Foo}.
//
// Java: Es6RewriteModules.RenameGlobalVars.fixTypeNode(NodeTraversal, Node)
func (r *Es6RewriteModules_RenameGlobalVars) FixTypeNode(t *NodeTraversal, typeNode *rhino.Node) {
	e := r.outer
	if typeNode.IsStringLit() {
		thisModule := e.moduleMap.GetModuleByPath(t.GetInput().GetPath())
		name := typeNode.GetString()
		splitted := es6RewriteModules_DOT_SPLITTER_splitToList(name)
		baseName := splitted[0]
		rest := ""
		if len(splitted) > 1 {
			rest = jutil.StringSubstringFrom(name, jutil.StringLength(baseName))
		}
		var_ := t.GetScope().GetVar(baseName)
		if var_ != nil && var_.IsGlobal() {
			r.MaybeSetNewName(
				t,
				typeNode,
				name,
				ModuleRenaming_GetGlobalNameOfEsModuleLocalVariable(thisModule.Metadata(), baseName).
					Join()+
					rest)
		} else if var_ == nil && thisModule.BoundNames().ContainsKey(baseName) {
			// Imports have been detached, so they won't show up in scope. Thus if we have a variable
			// not in scope that shares the name of an import it is the import.
			binding := thisModule.BoundNames().Get(baseName)
			globalName :=
				ModuleRenaming_GetGlobalNameForJsDoc(
					e.moduleMap, binding, splitted[1:len(splitted)])
			r.MaybeSetNewName(t, typeNode, name, globalName)

			if e.preprocessorSymbolTable != nil {
				// Jsdoc type node is a single STRING node that spans the whole type. For example
				// STRING node "bar.Foo". ES6 import rewrite replaces only "module"
				// part of the type: "bar.Foo" => "module$full$path$bar$Foo". We have to record
				// "bar" as alias.
				onlyBaseName := rhino.Node_NewString(baseName).Srcref(typeNode)
				onlyBaseName.SetLength(jutil.StringLength(baseName))
				e.MaybeAddAliasToSymbolTable(onlyBaseName, jutil.ToJavaString(t.GetSourceName()))
			}
		}

		typeNode.SetOriginalName(jutil.StringPtr(name))
	}

	for child := typeNode.GetFirstChild(); child != nil; child = child.GetNext() {
		r.FixTypeNode(t, child)
	}
}

// Java: Es6RewriteModules.RenameGlobalVars.maybeSetNewName(NodeTraversal, Node, String, String)
func (r *Es6RewriteModules_RenameGlobalVars) MaybeSetNewName(t *NodeTraversal, node *rhino.Node, name string, newName string) {
	if name != newName {
		node.SetString(newName)
		node.SetOriginalName(jutil.StringPtr(name))
		t.ReportCodeChange()
	}
}

// MaybeAddAliasToSymbolTable adds alias nodes to the symbol table as they
// going to be removed by rewriter. Example aliases:
//
//	import * as foo from './foo';
//	import {doBar} from './bar';
//
//	console.log(doBar);
//
// n is the alias node. In the example above alias nodes are foo, doBar, and
// doBar. module is the name of the module currently being processed.
//
// Java: Es6RewriteModules.maybeAddAliasToSymbolTable(Node, String)
func (e *Es6RewriteModules) MaybeAddAliasToSymbolTable(n *rhino.Node, module string) {
	if e.preprocessorSymbolTable == nil {
		return
	}
	n.PutBooleanProp(rhino.Node_MODULE_ALIAS, true)
	// Alias can be used in js types. Types have node type STRING and not NAME so we have to
	// use their name as string.
	var nodeName string
	if n.IsStringLit() || n.IsImportStar() {
		nodeName = n.GetString()
	} else {
		nodeName = jutil.ToJavaString(e.preprocessorSymbolTable.GetQualifiedName(n))
	}
	// We need to include module as part of the name because aliases are local to current module.
	// Aliases with the same name from different module should be completely different entities.
	name := "alias_" + module + "_" + nodeName
	e.preprocessorSymbolTable.AddReferenceWithName(n, name)
}

// MaybeAddImportedFileReferenceToSymbolTable adds reference to a file that
// current module imports. Example:
//
//	import * as qux from '../some/file.js';
//
// Will add a reference to file.js on the string node `'../some/file.js'`.
//
// importStringNode is the string node from the import statement that
// references imported file. In the example above it is the
// '../some/file.js' STRING node. importedFilePath is the absolute path to
// the imported file. In the example above it can be
// myproject/folder/some/file.js
//
// Java: Es6RewriteModules.maybeAddImportedFileReferenceToSymbolTable(Node, String)
func (e *Es6RewriteModules) MaybeAddImportedFileReferenceToSymbolTable(
	importStringNode *rhino.Node, importedFilePath string) {
	if e.preprocessorSymbolTable == nil {
		return
	}

	// If this if the first import that mentions importedFilePath then we need to create a SCRIPT
	// node for the imported file.
	if e.preprocessorSymbolTable.GetSlot(importedFilePath) == nil {
		scriptNode := e.compiler.GetScriptNode(importedFilePath)
		if scriptNode != nil {
			e.preprocessorSymbolTable.AddReferenceWithName(scriptNode, importedFilePath)
		}
	}

	e.preprocessorSymbolTable.AddReferenceWithName(importStringNode, importedFilePath)
}

// Java: Es6RewriteModules.declareGlobalVariable(Node, NodeTraversal)
func (e *Es6RewriteModules) DeclareGlobalVariable(n *rhino.Node, t *NodeTraversal) {
	jutil.CheckState(n.IsName())
	if !e.astFactory.IsAddingTypes() {
		return
	}
	jutil.CheckNotNull(e.globalTypedScope)

	name := n.GetString()
	if e.globalTypedScope.HasOwnSlot(name) {
		t.Report(t.GetCurrentScript(), ClosureRewriteModule_ILLEGAL_MODULE_RENAMING_CONFLICT, name)
	} else {
		type_ := jutil.CheckNotNull(n.GetJSType())
		e.globalTypedScope.Declare(name, n, type_, t.GetInput(), false)
	}
}
