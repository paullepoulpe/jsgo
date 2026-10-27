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

	"github.com/paullepoulpe/jsgo/internal/ast"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/parser"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// Es6RewriteModulesToCommonJsModules rewrites an ES6 module to a
// CommonJS-like module for the sake of per-file transpilation + bundling
// (e.g. Closure Bundler). Output is not meant to be type checked.
type Es6RewriteModulesToCommonJsModules struct {
	compiler AbstractCompiler
}

const (
	es6CommonJsDefaultExport = "$$default"
	es6CommonJsModule        = "$$module"
	es6CommonJsExports       = "$$exports"
	es6CommonJsRequire       = "$$require"
)

// Java: Es6RewriteModulesToCommonJsModules(AbstractCompiler)
func NewEs6RewriteModulesToCommonJsModules(compiler AbstractCompiler) *Es6RewriteModulesToCommonJsModules {
	return &Es6RewriteModulesToCommonJsModules{compiler: compiler}
}

// Java: Es6RewriteModulesToCommonJsModules.process(Node, Node)
func (e *Es6RewriteModulesToCommonJsModules) Process(externs rhino.Node, root rhino.Node) {
	for script := root.GetFirstChild(); !script.IsNil(); script = script.GetNext() {
		if Es6RewriteModules_IsEs6ModuleRoot(script) {
			NodeTraversal_Traverse(e.compiler, script, newEs6CommonJsModuleRewriter(e.compiler, script))
		}
	}
	NodeUtil_RemoveFeatureFromAllScripts(root, parser.FeatureSet_Feature_MODULES, e.compiler)
	GatherGetterAndSetterProperties_Update(e.compiler, externs, root)
}

type es6CommonJsLocalQName struct {
	qName             string
	nodeForSourceInfo rhino.Node
}

type es6CommonJsModuleRequest struct {
	specifier string
	varName   string
}

// Es6RewriteModulesToCommonJsModules_NormalizePath normalizes a registered
// or import path by stripping any protocol+host prefix or leading '/'.
//
// Java: Es6RewriteModulesToCommonJsModules.normalizePath(String)
func Es6RewriteModulesToCommonJsModules_NormalizePath(path string) string {
	if indexOfProtocol := strings.Index(path, "://"); indexOfProtocol >= 0 {
		path = path[indexOfProtocol+3:]
		if indexOfSlash := strings.IndexByte(path, '/'); indexOfSlash >= 0 {
			path = path[indexOfSlash+1:]
		}
	} else {
		path = strings.TrimPrefix(path, "/")
	}
	return path
}

// es6CommonJsModuleRewriter rewrites a single ES6 module into a CommonJS-like
// module designed to be loaded in the compiler's module runtime.
type es6CommonJsModuleRewriter struct {
	NodeTraversal_AbstractPostOrderCallbackBase

	requireInsertSpot        rhino.Node
	script                   rhino.Node
	exportedNameToLocalQName *jutil.TreeMap[string, es6CommonJsLocalQName]
	imports                  *jutil.LinkedHashSet[rhino.Node]
	importRequests           *jutil.LinkedHashSet[es6CommonJsModuleRequest]
	compiler                 AbstractCompiler
	modulePath               *ModuleLoader_ModulePath
}

func newEs6CommonJsModuleRewriter(compiler AbstractCompiler, script rhino.Node) *es6CommonJsModuleRewriter {
	return &es6CommonJsModuleRewriter{
		compiler: compiler,
		script:   script,
		// TreeMap because ES6 orders the export keys using natural ordering.
		exportedNameToLocalQName: jutil.NewTreeMap[string, es6CommonJsLocalQName](jutil.StringCompareTo),
		importRequests:           jutil.NewLinkedHashSet[es6CommonJsModuleRequest](),
		imports:                  jutil.NewLinkedHashSet[rhino.Node](),
		modulePath:               compiler.GetInput(script.GetInputId()).GetPath(),
	}
}

func (r *es6CommonJsModuleRewriter) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	switch n.GetToken() {
	case rhino.Token_IMPORT:
		r.visitImport(t.GetInput().GetPath(), ast.MustImport(n))
	case rhino.Token_EXPORT:
		r.visitExport(t, ast.MustExport(n), parent)
	case rhino.Token_SCRIPT:
		r.visitScript(t, n)
	case rhino.Token_NAME:
		r.maybeRenameImportedValue(t, ast.MustName(n))
	}
}

func (r *es6CommonJsModuleRewriter) varNameOfImport(importDecl ast.Import) string {
	if bindings, ok := importDecl.Bindings(); ok {
		if star, isStar := bindings.AsImportStar(); isStar {
			return star.Ident()
		}
	}
	return r.varNameOfImportRequest(importDecl.Source().Value())
}

func (r *es6CommonJsModuleRewriter) varNameOfImportRequest(importRequest string) string {
	return r.modulePath.ResolveModuleAsPath(importRequest).ToModuleName()
}

func (r *es6CommonJsModuleRewriter) nameOfImportedValue(nameNode rhino.Node) string {
	curr := nameNode
	for !curr.IsImport() {
		curr = curr.GetParent()
	}
	moduleName := r.varNameOfImport(ast.MustImport(curr))

	if spec, ok := ast.AsImportSpec(nameNode.GetParent()); ok {
		return moduleName + "." + spec.Imported().Ident()
	}
	if _, ok := ast.AsImportStar(nameNode); ok {
		return moduleName
	}
	jutil.CheckState(nameNode.GetParent().IsImport())
	return moduleName + ".default"
}

func (r *es6CommonJsModuleRewriter) maybeGetNameOfImportedValue(s *Scope, nameNode ast.Name) (string, bool) {
	v := s.GetVar(nameNode.Ident())
	if v != nil &&
		!v.GetNameNode().IsNil() &&
		NodeUtil_IsImportedName(v.GetNameNode()) &&
		nameNode.Node() != v.GetNameNode() {
		return r.nameOfImportedValue(v.GetNameNode()), true
	}
	return "", false
}

func (r *es6CommonJsModuleRewriter) maybeRenameImportedValue(t *NodeTraversal, nm ast.Name) {
	parent := nm.Node().GetParent()
	if parent.IsExport() ||
		parent.IsExportSpec() ||
		parent.IsImport() ||
		parent.IsImportSpec() {
		return
	}

	if qName, ok := r.maybeGetNameOfImportedValue(t.GetScope(), nm); ok {
		nm.Node().ReplaceWith(NodeUtil_NewQName(r.compiler, qName))
		t.ReportCodeChange()
	}
}

func (r *es6CommonJsModuleRewriter) visitScript(t *NodeTraversal, script rhino.Node) {
	jutil.CheckState(r.script == script)
	moduleBody := ast.MustModuleBody(script.GetFirstChild())
	moduleNode := moduleBody.Node()
	moduleNode.Detach()
	script.AddChildrenToFront(moduleNode.RemoveChildren())

	// Prepend sections in reverse order so the resulting body has:
	//   1) export definitions
	//   2) require calls
	//   3) module body statements
	r.addRequireCalls()
	r.addExportDef()
	r.registerAndLoadModule(t)
}

func (r *es6CommonJsModuleRewriter) addRequireCalls() {
	if r.importRequests.IsEmpty() {
		return
	}
	for importDecl := range r.imports.All() {
		importDecl.Detach()
	}

	importedNames := jutil.NewLinkedHashSet[string]()
	for request := range r.importRequests.All() {
		if importedNames.Add(request.varName) {
			requireCall := rhino.IR_Call(rhino.IR_Name(es6CommonJsRequire), rhino.IR_String(request.specifier))
			ast.MustCall(requireCall).SetFreeCall(true)
			decl := rhino.IR_VarWithValue(rhino.IR_Name(request.varName), requireCall)
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

func (r *es6CommonJsModuleRewriter) registerAndLoadModule(t *NodeTraversal) {
	block := rhino.IR_Block()
	block.AddChildrenToFront(r.script.RemoveChildren())

	moduleFunction :=
		rhino.IR_Function(
			rhino.IR_Name(""),
			rhino.IR_ParamList(
				rhino.IR_Name(es6CommonJsRequire),
				rhino.IR_Name(es6CommonJsExports),
				rhino.IR_Name(es6CommonJsModule)),
			block)

	shallowDeps := rhino.NewNode(rhino.Token_ARRAYLIT)
	for request := range r.importRequests.All() {
		shallowDeps.AddChildToBack(rhino.IR_String(request.specifier))
	}

	exprResult :=
		rhino.IR_ExprResult(
			rhino.IR_Call(
				rhino.IR_Getprop(rhino.IR_Name("$jscomp"), "registerAndLoadModule"),
				moduleFunction,
				rhino.IR_String(
					Es6RewriteModulesToCommonJsModules_NormalizePath(
						r.compiler.GetModuleLoader().Resolve(t.GetInput().GetName()).String())),
				shallowDeps))

	r.script.AddChildToBack(exprResult.SrcrefTreeIfMissing(r.script))

	r.compiler.ReportChangeToChangeScope(r.script)
	r.compiler.ReportChangeToChangeScope(moduleFunction)
	t.ReportCodeChange()
}

func (r *es6CommonJsModuleRewriter) addExportDef() {
	if r.exportedNameToLocalQName.IsEmpty() {
		return
	}
	definePropertiesLit := rhino.IR_Objectlit()
	for key, value := range r.exportedNameToLocalQName.All() {
		r.addExport(definePropertiesLit, key, value)
	}

	r.script.AddChildToFront(
		rhino.IR_ExprResult(
			rhino.IR_Call(
				NodeUtil_NewQName(r.compiler, "Object.defineProperties"),
				rhino.IR_Name(es6CommonJsExports),
				definePropertiesLit)).
			SrcrefTreeIfMissing(r.script))
}

func (r *es6CommonJsModuleRewriter) addExport(definePropertiesLit rhino.Node, exportedName string, localQName es6CommonJsLocalQName) {
	exportedValue := NodeUtil_NewQName(r.compiler, localQName.qName)
	getterFunction :=
		rhino.IR_Function(rhino.IR_Name(""), rhino.IR_ParamList(), rhino.IR_BlockOfNode(rhino.IR_ReturnNodeWithExpr(exportedValue)))
	getterFunction.SrcrefTree(localQName.nodeForSourceInfo)

	objLit :=
		rhino.IR_Objectlit(
			rhino.IR_StringKeyWithValue("enumerable", rhino.IR_TrueNode()),
			rhino.IR_StringKeyWithValue("get", getterFunction))
	definePropertiesLit.AddChildToBack(rhino.IR_StringKeyWithValue(exportedName, objLit))

	r.compiler.ReportChangeToChangeScope(getterFunction)
}

func (r *es6CommonJsModuleRewriter) visitImport(path *ModuleLoader_ModulePath, importDecl ast.Import) {
	specifier := importDecl.Source().Value()
	if strings.Contains(specifier, "://") {
		r.compiler.Report(
			JSError_MakeWithNode(
				importDecl.Node(), TranspilationUtil_CANNOT_CONVERT, "Module requests with protocols."))
	}

	r.importRequests.Add(es6CommonJsModuleRequest{
		specifier: Es6RewriteModulesToCommonJsModules_NormalizePath(path.ResolveModuleAsPath(specifier).String()),
		varName:   r.varNameOfImport(importDecl),
	})
	r.imports.Add(importDecl.Node())
}

func (r *es6CommonJsModuleRewriter) visitExportDefault(t *NodeTraversal, export ast.Export) {
	exportNode := export.Node()
	child := export.Body().Node()
	var localName string
	hasDeclName := false

	if fn, ok := ast.AsFunction(child); ok {
		if nm, named := fn.Name(); named {
			localName = nm.Ident()
			hasDeclName = true
		}
	} else if cls, ok := ast.AsClass(child); ok {
		if nm, named := cls.Name(); named {
			localName = nm.Ident()
			hasDeclName = true
		}
	}

	if hasDeclName {
		exportNode.ReplaceWith(child.Detach())
	} else {
		localName = es6CommonJsDefaultExport
		constDecl := rhino.IR_ConstNode(rhino.IR_Name(localName), exportNode.RemoveFirstChild())
		exportNode.ReplaceWith(constDecl.SrcrefTreeIfMissing(exportNode))
		NodeUtil_AddFeatureToScript(t.GetCurrentScript(), parser.FeatureSet_Feature_CONST_DECLARATIONS, r.compiler)
	}

	r.exportedNameToLocalQName.Put("default", es6CommonJsLocalQName{qName: localName, nodeForSourceInfo: exportNode})
	t.ReportCodeChange()
}

func (r *es6CommonJsModuleRewriter) visitExportFrom(t *NodeTraversal, export ast.Export, source ast.StringLit, parent rhino.Node) {
	exportNode := export.Node()
	importNode := rhino.IR_ImportNode(rhino.IR_Empty(), rhino.IR_Empty(), source.Node().CloneNode())
	importNode.Srcref(exportNode)
	importNode.InsertBefore(exportNode)
	r.Visit(t, importNode, parent)

	moduleName := r.varNameOfImportRequest(source.Value())
	specs := ast.MustExportSpecs(export.Body().Node())
	for exportSpec := range specs.Specs() {
		r.exportedNameToLocalQName.Put(
			exportSpec.Exported().Ident(),
			es6CommonJsLocalQName{
				qName:             moduleName + "." + exportSpec.Local().Ident(),
				nodeForSourceInfo: exportSpec.Node(),
			})
	}

	exportNode.Detach()
	t.ReportCodeChange()
}

func (r *es6CommonJsModuleRewriter) visitExportSpecs(t *NodeTraversal, export ast.Export, specs ast.ExportSpecs) {
	for exportSpec := range specs.Specs() {
		localNode := exportSpec.Local()
		localName := localNode.Ident()
		v := t.GetScope().GetVar(localName)
		if v != nil && NodeUtil_IsImportedName(v.GetNameNode()) {
			importedQName, ok := r.maybeGetNameOfImportedValue(t.GetScope(), localNode)
			if !ok {
				panic(jutil.NewNullPointerException())
			}
			localName = importedQName
		}
		r.exportedNameToLocalQName.Put(
			exportSpec.Exported().Ident(),
			es6CommonJsLocalQName{qName: localName, nodeForSourceInfo: exportSpec.Node()})
	}
	export.Node().Detach()
	t.ReportCodeChange()
}

func (r *es6CommonJsModuleRewriter) addExportedName(lhs rhino.Node) {
	nm := ast.MustName(lhs)
	name := nm.Ident()
	r.exportedNameToLocalQName.Put(name, es6CommonJsLocalQName{qName: name, nodeForSourceInfo: lhs})
}

func (r *es6CommonJsModuleRewriter) visitExportDeclaration(t *NodeTraversal, export ast.Export) {
	exportNode := export.Node()
	declaration := export.Body().Node()

	if _, ok := ast.AsNameDeclaration(declaration); ok {
		NodeUtil_VisitLhsNodesInNode(declaration, r.addExportedName)
	} else if fn, ok := ast.AsFunction(declaration); ok {
		name := fn.NameNode().Ident()
		r.exportedNameToLocalQName.Put(name, es6CommonJsLocalQName{qName: name, nodeForSourceInfo: exportNode})
	} else {
		cls := ast.MustClass(declaration)
		nm, _ := cls.Name()
		name := nm.Ident()
		r.exportedNameToLocalQName.Put(name, es6CommonJsLocalQName{qName: name, nodeForSourceInfo: exportNode})
	}

	exportNode.ReplaceWith(declaration.Detach())
	t.ReportCodeChange()
}

func (r *es6CommonJsModuleRewriter) visitExportStar(t *NodeTraversal, export ast.Export, parent rhino.Node) {
	exportNode := export.Node()
	source, ok := export.Source()
	jutil.CheckState(ok)

	importNode := rhino.IR_ImportNode(rhino.IR_Empty(), rhino.IR_Empty(), source.Node().CloneNode())
	importNode.Srcref(exportNode)
	importNode.InsertBefore(exportNode)
	r.Visit(t, importNode, parent)

	moduleName := r.varNameOfImportRequest(source.Value())
	exportNode.ReplaceWith(
		rhino.IR_ExprResult(
			rhino.IR_Call(rhino.IR_Getprop(rhino.IR_Name(es6CommonJsModule), "exportAllFrom"), rhino.IR_Name(moduleName))).
			SrcrefTree(exportNode))

	t.ReportCodeChange()
}

func (r *es6CommonJsModuleRewriter) visitExport(t *NodeTraversal, export ast.Export, parent rhino.Node) {
	if export.IsDefault() {
		r.visitExportDefault(t, export)
	} else if export.IsAllFrom() {
		r.visitExportStar(t, export, parent)
	} else if source, hasSource := export.Source(); hasSource {
		r.visitExportFrom(t, export, source, parent)
	} else if specs, isSpecs := export.Body().AsExportSpecs(); isSpecs {
		r.visitExportSpecs(t, export, specs)
	} else {
		r.visitExportDeclaration(t, export)
	}
}
