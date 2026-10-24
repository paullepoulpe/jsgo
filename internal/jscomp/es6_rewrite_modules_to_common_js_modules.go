// Port of com/google/javascript/jscomp/Es6RewriteModulesToCommonJsModules.java @ closure-compiler v20261021

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
	"github.com/paullepoulpe/jsgo/internal/parser"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// Es6RewriteModulesToCommonJsModules rewrites an ES6 module to a
// CommonJS-like module for the sake of per-file transpilation + bunlding
// (e.g. Closure Bundler). Output is not meant to be type checked.
type Es6RewriteModulesToCommonJsModules struct {
	compiler AbstractCompiler
}

const es6RewriteModulesToCommonJsModules_JSCOMP_DEFAULT_EXPORT = "$$default"
const es6RewriteModulesToCommonJsModules_MODULE = "$$module"
const es6RewriteModulesToCommonJsModules_EXPORTS = "$$exports"
const es6RewriteModulesToCommonJsModules_REQUIRE = "$$require"

// Java: Es6RewriteModulesToCommonJsModules(AbstractCompiler)
func NewEs6RewriteModulesToCommonJsModules(compiler AbstractCompiler) *Es6RewriteModulesToCommonJsModules {
	e := &Es6RewriteModulesToCommonJsModules{}
	e.compiler = compiler
	return e
}

// Java: Es6RewriteModulesToCommonJsModules.process(Node, Node)
func (e *Es6RewriteModulesToCommonJsModules) Process(externs rhino.Node, root rhino.Node) {
	for script := root.GetFirstChild(); !script.IsNil(); script = script.GetNext() {
		if Es6RewriteModules_IsEs6ModuleRoot(script) {
			NodeTraversal_Traverse(e.compiler, script, NewEs6RewriteModulesToCommonJsModules_Rewriter(e.compiler, script))
		}
	}
	// It is unusual to call this NodeUtil method instead of the TranspilationPasses one. This
	// pass is included in the {@code dependency_resolution} BUILD target and does not have access
	// to {@code TranspilationPasses}. Adding that dep produces a cycle in the BUILD dep graph.
	// Regular transpiler passes must use the {@code
	// TranspilationPasses.maybeMarkFeaturesAsTranspiledAway}
	NodeUtil_RemoveFeatureFromAllScripts(root, parser.FeatureSet_Feature_MODULES, e.compiler)
	GatherGetterAndSetterProperties_Update(e.compiler, externs, root)
}

// Es6RewriteModulesToCommonJsModules_LocalQName is the private static
// nested class Es6RewriteModulesToCommonJsModules.LocalQName.
type Es6RewriteModulesToCommonJsModules_LocalQName struct {
	qName string

	// Node to use for source information. All exported properties are
	// transpiled to ES5 getters so debuggers will automatically step into
	// these, even with source maps. When stepping in this node will be
	// displayed in the source map. In general it should either be the export
	// itself (e.g. export function or class) or the specific name being
	// exported (export specs, const, etc).
	nodeForSourceInfo rhino.Node
}

// Java: Es6RewriteModulesToCommonJsModules.LocalQName(String, Node)
func NewEs6RewriteModulesToCommonJsModules_LocalQName(qName string, nodeForSourceInfo rhino.Node) *Es6RewriteModulesToCommonJsModules_LocalQName {
	l := &Es6RewriteModulesToCommonJsModules_LocalQName{}
	l.qName = qName
	l.nodeForSourceInfo = nodeForSourceInfo
	return l
}

// Es6RewriteModulesToCommonJsModules_NormalizePath normalizes a registered
// or import path.
//
// Absolute import paths need to match the registered path exactly. Some
// ModuleLoader.ModulePaths will have a leading slash and some won't. So in
// order to have everything line up AND preserve schemes (if they exist)
// then just strip leading /.
//
// Additionally if any path contains a protocol it will be stripped only the
// path part will remain. This is done heuristically as we cannot use
// java.net.URL or java.nio.file.Path due to GWT. As a result of stripping
// this cross-domain imports are not compatible with this pass.
//
// Java: Es6RewriteModulesToCommonJsModules.normalizePath(String)
func Es6RewriteModulesToCommonJsModules_NormalizePath(path string) string {
	indexOfProtocol := jutil.StringIndexOf(path, "://")
	if indexOfProtocol > -1 {
		path = jutil.StringSubstringFrom(path, indexOfProtocol+3)
		indexOfSlash := jutil.StringIndexOfChar(path, '/')
		if indexOfSlash > -1 {
			path = jutil.StringSubstringFrom(path, indexOfSlash+1)
		}
	} else if strings.HasPrefix(path, "/") {
		path = jutil.StringSubstringFrom(path, 1)
	}
	return path
}

// Es6RewriteModulesToCommonJsModules_ModuleRequest is the private record
// Es6RewriteModulesToCommonJsModules.ModuleRequest(String specifier, String
// varName), compared by value in a LinkedHashSet.
type Es6RewriteModulesToCommonJsModules_ModuleRequest struct {
	specifier string
	varName   string
}

// Java: Es6RewriteModulesToCommonJsModules.ModuleRequest(String, String)
func NewEs6RewriteModulesToCommonJsModules_ModuleRequest(specifier string, varName string) *Es6RewriteModulesToCommonJsModules_ModuleRequest {
	// requireNonNull(specifier, "specifier"); requireNonNull(varName, "varName"): Go strings
	// cannot be null.
	return &Es6RewriteModulesToCommonJsModules_ModuleRequest{specifier: specifier, varName: varName}
}

// Specifier is the record component accessor specifier().
//
// Java: Es6RewriteModulesToCommonJsModules.ModuleRequest.specifier()
func (m *Es6RewriteModulesToCommonJsModules_ModuleRequest) Specifier() string { return m.specifier }

// VarName is the record component accessor varName().
//
// Java: Es6RewriteModulesToCommonJsModules.ModuleRequest.varName()
func (m *Es6RewriteModulesToCommonJsModules_ModuleRequest) VarName() string { return m.varName }

// Equals is the record's equals: both components equal.
//
// Java: Es6RewriteModulesToCommonJsModules.ModuleRequest.equals(Object)
func (m *Es6RewriteModulesToCommonJsModules_ModuleRequest) Equals(o any) bool {
	that, ok := o.(*Es6RewriteModulesToCommonJsModules_ModuleRequest)
	return ok && that != nil && m.specifier == that.specifier && m.varName == that.varName
}

// HashCode is the record's hashCode, derived from both components. The
// set that holds requests iterates in insertion order, so only its
// consistency with Equals matters.
//
// Java: Es6RewriteModulesToCommonJsModules.ModuleRequest.hashCode()
func (m *Es6RewriteModulesToCommonJsModules_ModuleRequest) HashCode() int32 {
	return 31*jutil.StringHashCode(m.specifier) + jutil.StringHashCode(m.varName)
}

// Java: Es6RewriteModulesToCommonJsModules.ModuleRequest.create(String, String)
func Es6RewriteModulesToCommonJsModules_ModuleRequest_Create(specifier string, varName string) *Es6RewriteModulesToCommonJsModules_ModuleRequest {
	return NewEs6RewriteModulesToCommonJsModules_ModuleRequest(specifier, varName)
}

// Es6RewriteModulesToCommonJsModules_Rewriter rewrites a single ES6 module
// into a CommonJS like module designed to be loaded in the compiler's module
// runtime.
type Es6RewriteModulesToCommonJsModules_Rewriter struct {
	NodeTraversal_AbstractPostOrderCallbackBase

	requireInsertSpot        rhino.Node
	script                   rhino.Node
	exportedNameToLocalQName *jutil.TreeMap[string, *Es6RewriteModulesToCommonJsModules_LocalQName]
	imports                  *jutil.LinkedHashSet[rhino.Node]
	importRequests           *jutil.ObjectLinkedHashSet[*Es6RewriteModulesToCommonJsModules_ModuleRequest]
	compiler                 AbstractCompiler
	modulePath               *ModuleLoader_ModulePath
}

// Java: Es6RewriteModulesToCommonJsModules.Rewriter(AbstractCompiler, Node)
func NewEs6RewriteModulesToCommonJsModules_Rewriter(compiler AbstractCompiler, script rhino.Node) *Es6RewriteModulesToCommonJsModules_Rewriter {
	r := &Es6RewriteModulesToCommonJsModules_Rewriter{}
	r.compiler = compiler
	r.script = script
	r.requireInsertSpot = rhino.NilNode
	// TreeMap because ES6 orders the export key using natural ordering.
	r.exportedNameToLocalQName = jutil.NewTreeMap[string, *Es6RewriteModulesToCommonJsModules_LocalQName](jutil.StringCompareTo)
	r.importRequests = jutil.NewObjectLinkedHashSet[*Es6RewriteModulesToCommonJsModules_ModuleRequest]()
	r.imports = jutil.NewLinkedHashSet[rhino.Node]()
	r.modulePath = compiler.GetInput(script.GetInputId()).GetPath()
	return r
}

// Java: Es6RewriteModulesToCommonJsModules.Rewriter.visit(NodeTraversal, Node, Node)
func (r *Es6RewriteModulesToCommonJsModules_Rewriter) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	switch n.GetToken() {
	case rhino.Token_IMPORT:
		r.VisitImport(t.GetInput().GetPath(), n)
	case rhino.Token_EXPORT:
		r.VisitExport(t, n, parent)
	case rhino.Token_SCRIPT:
		r.VisitScript(t, n)
	case rhino.Token_NAME:
		r.MaybeRenameImportedValue(t, n)
	default:
	}
}

// GetVarNameOfImport, given an import node, gets the name of the var to use
// for the imported module.
//
// Example: import {v} from './foo.js'; use(v); Can become:
//
//	const module$foo = require('./foo.js');
//	use(module$foo.v);
//
// This method would return "module$foo".
//
// Note that if there is a star import the name will be preserved.
//
// Example:
//
//	import defaultValue, * as foo from './foo.js';
//	use(defaultValue, foo.bar);
//
// Can become:
//
//	const foo = require('./foo.js'); use(foo.defaultValue, foo.bar);
//
// This makes debugging quite a bit easier as source maps are not great with
// renaming.
//
// Java: Es6RewriteModulesToCommonJsModules.Rewriter.getVarNameOfImport(Node)
func (r *Es6RewriteModulesToCommonJsModules_Rewriter) GetVarNameOfImport(importDecl rhino.Node) string {
	jutil.CheckState(importDecl.IsImport())
	if importDecl.GetSecondChild().IsImportStar() {
		return importDecl.GetSecondChild().GetString()
	}
	return r.GetVarNameOfImportWithImportRequest(importDecl.GetLastChild().GetString())
}

// Java: Es6RewriteModulesToCommonJsModules.Rewriter.getVarNameOfImport(String)
func (r *Es6RewriteModulesToCommonJsModules_Rewriter) GetVarNameOfImportWithImportRequest(importRequest string) string {
	return r.modulePath.ResolveModuleAsPath(importRequest).ToModuleName()
}

// GetNameOfImportedValue returns the qualified name to use to reference an
// imported value.
//
// Examples:
//
//   - If referencing an import spec like v in "import {v} from './foo.js'"
//     then this would return "module$foo.v".
//   - If referencing an import star like m in "import * as m from
//     './foo.js'" then this would return "m".
//   - If referencing an import default like d in "import d from './foo.js'"
//     then this would return "module$foo.default".
//
// Used to rename references to imported values within this module.
//
// Java: Es6RewriteModulesToCommonJsModules.Rewriter.getNameOfImportedValue(Node)
func (r *Es6RewriteModulesToCommonJsModules_Rewriter) GetNameOfImportedValue(nameNode rhino.Node) string {
	importDecl := nameNode

	for !importDecl.IsImport() {
		importDecl = importDecl.GetParent()
	}

	moduleName := r.GetVarNameOfImport(importDecl)

	if nameNode.GetParent().IsImportSpec() {
		return moduleName + "." + nameNode.GetParent().GetFirstChild().GetString()
	} else if nameNode.IsImportStar() {
		return moduleName
	} else {
		jutil.CheckState(nameNode.GetParent().IsImport())
		return moduleName + ".default"
	}
}

// MaybeGetNameOfImportedValue: nameNode is any variable name that is
// potentially from an import statement. Returns the qualified name to use to
// reference an imported value if the given node is an imported name or nil
// if the value is not imported or if it is in the import statement itself.
//
// Java: Es6RewriteModulesToCommonJsModules.Rewriter.maybeGetNameOfImportedValue(Scope, Node)
func (r *Es6RewriteModulesToCommonJsModules_Rewriter) MaybeGetNameOfImportedValue(s *Scope, nameNode rhino.Node) *string {
	jutil.CheckState(nameNode.IsName())
	var_ := s.GetVar(nameNode.GetString())

	if var_ != nil &&
		// variables added implicitly to the scope, like arguments, have a null name node
		!var_.GetNameNode().IsNil() &&
		NodeUtil_IsImportedName(var_.GetNameNode()) &&
		nameNode != var_.GetNameNode() {
		return jutil.StringPtr(r.GetNameOfImportedValue(var_.GetNameNode()))
	}

	return nil
}

// MaybeRenameImportedValue renames the given name node if it is an imported
// value.
//
// Java: Es6RewriteModulesToCommonJsModules.Rewriter.maybeRenameImportedValue(NodeTraversal, Node)
func (r *Es6RewriteModulesToCommonJsModules_Rewriter) MaybeRenameImportedValue(t *NodeTraversal, n rhino.Node) {
	jutil.CheckState(n.IsName())
	parent := n.GetParent()

	if parent.IsExport() ||
		parent.IsExportSpec() ||
		parent.IsImport() ||
		parent.IsImportSpec() {
		return
	}

	qName := r.MaybeGetNameOfImportedValue(t.GetScope(), n)

	if qName != nil {
		n.ReplaceWith(NodeUtil_NewQName(r.compiler, *qName))
		t.ReportCodeChange()
	}
}

// Java: Es6RewriteModulesToCommonJsModules.Rewriter.visitScript(NodeTraversal, Node)
func (r *Es6RewriteModulesToCommonJsModules_Rewriter) VisitScript(t *NodeTraversal, script rhino.Node) {
	jutil.CheckState(r.script == script)
	moduleNode := script.GetFirstChild()
	jutil.CheckState(moduleNode.IsModuleBody())
	moduleNode.Detach()
	script.AddChildrenToFront(moduleNode.RemoveChildren())

	// Order here is important. We want the end result to be:
	//  $jscomp.registerAndLoadModule(function($$require, $$exports, $$module) {
	//   // First to ensure circular deps can see exports of this module before we require them,
	//   // and also so that temporal deadzone is respected.
	//   //<export def>
	//   // Second so the module definition can reference imported modules, and so any require'd
	//   // modules are loaded.
	//   //<requires>
	//   // And finally last is the actual module definition.
	//   //<module def>
	//  }, /* <module path> */, [/* <deps> */]);
	// As a result the calls below are in *inverse* order to what we want above so they can keep
	// adding to the front of the script.
	r.AddRequireCalls()
	r.AddExportDef()
	r.RegisterAndLoadModule(t)
}

// AddRequireCalls adds one call to require per imported module.
//
// Java: Es6RewriteModulesToCommonJsModules.Rewriter.addRequireCalls()
func (r *Es6RewriteModulesToCommonJsModules_Rewriter) AddRequireCalls() {
	if !r.importRequests.IsEmpty() {
		for importDecl := range r.imports.All() {
			importDecl.Detach()
		}

		importedNames := jutil.NewLinkedHashSet[string]()

		for request := range r.importRequests.All() {
			varName := request.VarName()
			if importedNames.Add(varName) {
				requireCall := rhino.IR_Call(rhino.IR_Name(es6RewriteModulesToCommonJsModules_REQUIRE), rhino.IR_String(request.Specifier()))
				requireCall.PutBooleanProp(rhino.Node_FREE_CALL, true)
				decl := rhino.IR_VarWithValue(rhino.IR_Name(varName), requireCall)
				decl.SrcrefTreeIfMissing(r.script)
				if r.requireInsertSpot.IsNil() {
					r.script.AddChildToFront(decl)
				} else {
					decl.InsertAfter(r.requireInsertSpot)
				}
				r.requireInsertSpot = decl
			}
		}
	}
}

// RegisterAndLoadModule wraps the entire current module definition in a
// $jscomp.registerAndLoadModule function.
//
// Java: Es6RewriteModulesToCommonJsModules.Rewriter.registerAndLoadModule(NodeTraversal)
func (r *Es6RewriteModulesToCommonJsModules_Rewriter) RegisterAndLoadModule(t *NodeTraversal) {
	block := rhino.IR_Block()
	block.AddChildrenToFront(r.script.RemoveChildren())

	// TODO(b/282006497): Maybe mark this function for strict mode?
	// NOTE: One might be tempted to add `'use strict';` here, but that causes problems.
	// Optimizations and transpilations are not written to look for this statement,
	// and are likely to either remove it or move other statements ahead of it,
	// making it ineffective.
	// The "right" way to mark the method for strict mode would be to apply the
	// USE_STRICT node property to its body, but at the moment that will have no
	// effect because 1) the code printer ignores it and 2) TypedAst doesn't
	// serialize it.

	moduleFunction :=
		rhino.IR_Function(
			rhino.IR_Name(""),
			rhino.IR_ParamList(
				rhino.IR_Name(es6RewriteModulesToCommonJsModules_REQUIRE),
				rhino.IR_Name(es6RewriteModulesToCommonJsModules_EXPORTS),
				rhino.IR_Name(es6RewriteModulesToCommonJsModules_MODULE)),
			block)

	shallowDeps := rhino.NewNode(rhino.Token_ARRAYLIT)

	for request := range r.importRequests.All() {
		shallowDeps.AddChildToBack(rhino.IR_String(request.Specifier()))
	}

	exprResult :=
		rhino.IR_ExprResult(
			rhino.IR_Call(
				rhino.IR_Getprop(rhino.IR_Name("$jscomp"), "registerAndLoadModule"),
				moduleFunction,
				// Resolving this path enables removing module roots from this path.
				rhino.IR_String(
					Es6RewriteModulesToCommonJsModules_NormalizePath(
						r.compiler.GetModuleLoader().Resolve(t.GetInput().GetName()).String())),
				shallowDeps))

	r.script.AddChildToBack(exprResult.SrcrefTreeIfMissing(r.script))

	r.compiler.ReportChangeToChangeScope(r.script)
	r.compiler.ReportChangeToChangeScope(moduleFunction)
	t.ReportCodeChange()
}

// AddExportDef adds exports to the exports object using
// Object.defineProperties.
//
// Java: Es6RewriteModulesToCommonJsModules.Rewriter.addExportDef()
func (r *Es6RewriteModulesToCommonJsModules_Rewriter) AddExportDef() {
	if !r.exportedNameToLocalQName.IsEmpty() {
		definePropertiesLit := rhino.IR_Objectlit()

		for key, value := range r.exportedNameToLocalQName.All() {
			r.AddExport(definePropertiesLit, key, value)
		}

		r.script.AddChildToFront(
			rhino.IR_ExprResult(
				rhino.IR_Call(
					NodeUtil_NewQName(r.compiler, "Object.defineProperties"),
					rhino.IR_Name(es6RewriteModulesToCommonJsModules_EXPORTS),
					definePropertiesLit)).
				SrcrefTreeIfMissing(r.script))
	}
}

// AddExport adds an ES5 getter to the given object literal to use an an
// export.
//
// Java: Es6RewriteModulesToCommonJsModules.Rewriter.addExport(Node, String, LocalQName)
func (r *Es6RewriteModulesToCommonJsModules_Rewriter) AddExport(definePropertiesLit rhino.Node, exportedName string, localQName *Es6RewriteModulesToCommonJsModules_LocalQName) {
	exportedValue := NodeUtil_NewQName(r.compiler, localQName.qName)
	getterFunction :=
		rhino.IR_Function(rhino.IR_Name(""), rhino.IR_ParamList(), rhino.IR_BlockOfNode(rhino.IR_ReturnNodeWithExpr(exportedValue)))
	getterFunction.SrcrefTree(localQName.nodeForSourceInfo)

	objLit :=
		rhino.IR_Objectlit(
			rhino.IR_StringKeyWithValue("enumerable", rhino.IR_TrueNode()), rhino.IR_StringKeyWithValue("get", getterFunction))
	definePropertiesLit.AddChildToBack(rhino.IR_StringKeyWithValue(exportedName, objLit))

	r.compiler.ReportChangeToChangeScope(getterFunction)
}

// Java: Es6RewriteModulesToCommonJsModules.Rewriter.visitImport(ModulePath, Node)
func (r *Es6RewriteModulesToCommonJsModules_Rewriter) VisitImport(path *ModuleLoader_ModulePath, importDecl rhino.Node) {
	if strings.Contains(importDecl.GetLastChild().GetString(), "://") {
		r.compiler.Report(
			JSError_MakeWithNode(
				importDecl, TranspilationUtil_CANNOT_CONVERT, "Module requests with protocols."))
	}

	// Normalize the import path according to the module resolution scheme so that bundles are
	// compatible with the compiler's module loader options.
	r.importRequests.Add(
		Es6RewriteModulesToCommonJsModules_ModuleRequest_Create(
			Es6RewriteModulesToCommonJsModules_NormalizePath(
				path.ResolveModuleAsPath(importDecl.GetLastChild().GetString()).String()),
			r.GetVarNameOfImport(importDecl)))
	r.imports.Add(importDecl)
}

// Java: Es6RewriteModulesToCommonJsModules.Rewriter.visitExportDefault(NodeTraversal, Node)
func (r *Es6RewriteModulesToCommonJsModules_Rewriter) VisitExportDefault(t *NodeTraversal, export rhino.Node) {
	child := export.GetFirstChild()
	var name *string

	if child.IsFunction() || child.IsClass() {
		name = NodeUtil_GetName(child)
	}

	if name != nil {
		decl := child.Detach()
		export.ReplaceWith(decl)
	} else {
		name = jutil.StringPtr(es6RewriteModulesToCommonJsModules_JSCOMP_DEFAULT_EXPORT)
		// Default exports are constant in more ways than one. Not only can they not be
		// overwritten but they also act like a const for temporal dead-zone purposes.
		var_ := rhino.IR_ConstNode(rhino.IR_Name(*name), export.RemoveFirstChild())
		export.ReplaceWith(var_.SrcrefTreeIfMissing(export))
		NodeUtil_AddFeatureToScript(t.GetCurrentScript(), parser.FeatureSet_Feature_CONST_DECLARATIONS, r.compiler)
	}

	r.exportedNameToLocalQName.Put("default", NewEs6RewriteModulesToCommonJsModules_LocalQName(*name, export))
	t.ReportCodeChange()
}

// Java: Es6RewriteModulesToCommonJsModules.Rewriter.visitExportFrom(NodeTraversal, Node, Node)
func (r *Es6RewriteModulesToCommonJsModules_Rewriter) VisitExportFrom(t *NodeTraversal, export rhino.Node, parent rhino.Node) {
	//   export {x, y as z} from 'moduleIdentifier';
	moduleIdentifier := export.GetLastChild()
	importNode := rhino.IR_ImportNode(rhino.IR_Empty(), rhino.IR_Empty(), moduleIdentifier.CloneNode())
	importNode.Srcref(export)
	importNode.InsertBefore(export)
	r.Visit(t, importNode, parent)

	moduleName := r.GetVarNameOfImportWithImportRequest(moduleIdentifier.GetString())

	for exportSpec := export.GetFirstFirstChild(); !exportSpec.IsNil(); exportSpec = exportSpec.GetNext() {
		r.exportedNameToLocalQName.Put(
			exportSpec.GetLastChild().GetString(),
			NewEs6RewriteModulesToCommonJsModules_LocalQName(moduleName+"."+exportSpec.GetFirstChild().GetString(), exportSpec))
	}

	export.Detach()
	t.ReportCodeChange()
}

// Java: Es6RewriteModulesToCommonJsModules.Rewriter.visitExportSpecs(NodeTraversal, Node)
func (r *Es6RewriteModulesToCommonJsModules_Rewriter) VisitExportSpecs(t *NodeTraversal, export rhino.Node) {
	//     export {Foo};
	for exportSpec := export.GetFirstFirstChild(); !exportSpec.IsNil(); exportSpec = exportSpec.GetNext() {
		localName := exportSpec.GetFirstChild().GetString()
		var_ := t.GetScope().GetVar(localName)
		if var_ != nil && NodeUtil_IsImportedName(var_.GetNameNode()) {
			localName = *jutil.CheckNotNull(r.MaybeGetNameOfImportedValue(t.GetScope(), exportSpec.GetFirstChild()))
		}
		r.exportedNameToLocalQName.Put(
			exportSpec.GetLastChild().GetString(), NewEs6RewriteModulesToCommonJsModules_LocalQName(localName, exportSpec))
	}
	export.Detach()
	t.ReportCodeChange()
}

// Java: Es6RewriteModulesToCommonJsModules.Rewriter.visitExportNameDeclaration(Node)
func (r *Es6RewriteModulesToCommonJsModules_Rewriter) VisitExportNameDeclaration(declaration rhino.Node) {
	//    export var Foo;
	//    export let {a, b:[c,d]} = {};
	NodeUtil_VisitLhsNodesInNode(declaration, r.AddExportedName)
}

// Java: Es6RewriteModulesToCommonJsModules.Rewriter.addExportedName(Node)
func (r *Es6RewriteModulesToCommonJsModules_Rewriter) AddExportedName(lhs rhino.Node) {
	jutil.CheckState(lhs.IsName())
	name := lhs.GetString()
	r.exportedNameToLocalQName.Put(name, NewEs6RewriteModulesToCommonJsModules_LocalQName(name, lhs))
}

// Java: Es6RewriteModulesToCommonJsModules.Rewriter.visitExportDeclaration(NodeTraversal, Node)
func (r *Es6RewriteModulesToCommonJsModules_Rewriter) VisitExportDeclaration(t *NodeTraversal, export rhino.Node) {
	//    export var Foo;
	//    export function Foo() {}
	// etc.
	declaration := export.GetFirstChild()

	if NodeUtil_IsNameDeclaration(declaration) {
		r.VisitExportNameDeclaration(declaration)
	} else {
		jutil.CheckState(declaration.IsFunction() || declaration.IsClass())
		name := declaration.GetFirstChild().GetString()
		r.exportedNameToLocalQName.Put(name, NewEs6RewriteModulesToCommonJsModules_LocalQName(name, export))
	}

	export.ReplaceWith(declaration.Detach())
	t.ReportCodeChange()
}

// Java: Es6RewriteModulesToCommonJsModules.Rewriter.visitExportStar(NodeTraversal, Node, Node)
func (r *Es6RewriteModulesToCommonJsModules_Rewriter) VisitExportStar(t *NodeTraversal, export rhino.Node, parent rhino.Node) {
	//   export * from 'moduleIdentifier';
	moduleIdentifier := export.GetLastChild()

	// Make an "import 'spec'" from this export node and then visit it to rewrite to a require().
	importNode := rhino.IR_ImportNode(rhino.IR_Empty(), rhino.IR_Empty(), moduleIdentifier.CloneNode())
	importNode.Srcref(export)
	importNode.InsertBefore(export)
	r.Visit(t, importNode, parent)

	moduleName := r.GetVarNameOfImportWithImportRequest(moduleIdentifier.GetString())
	export.ReplaceWith(
		rhino.IR_ExprResult(
			rhino.IR_Call(rhino.IR_Getprop(rhino.IR_Name("$$module"), "exportAllFrom"), rhino.IR_Name(moduleName))).
			SrcrefTree(export))

	t.ReportCodeChange()
}

// Java: Es6RewriteModulesToCommonJsModules.Rewriter.visitExport(NodeTraversal, Node, Node)
func (r *Es6RewriteModulesToCommonJsModules_Rewriter) VisitExport(t *NodeTraversal, export rhino.Node, parent rhino.Node) {
	if export.GetBooleanProp(rhino.Node_EXPORT_DEFAULT) {
		r.VisitExportDefault(t, export)
	} else if export.GetBooleanProp(rhino.Node_EXPORT_ALL_FROM) {
		r.VisitExportStar(t, export, parent)
	} else if export.HasTwoChildren() {
		r.VisitExportFrom(t, export, parent)
	} else {
		if export.GetFirstChild().IsExportSpecs() {
			r.VisitExportSpecs(t, export)
		} else {
			r.VisitExportDeclaration(t, export)
		}
	}
}
