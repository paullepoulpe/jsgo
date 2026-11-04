// Port of com/google/javascript/jscomp/RewriteGoogJsImports.java @ closure-compiler v20261028

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
	"slices"
	"strings"

	"github.com/paullepoulpe/jsgo/internal/ast"
	"github.com/paullepoulpe/jsgo/internal/base"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
	"github.com/paullepoulpe/jsgo/internal/strpool"
)

// RewriteGoogJsImports looks for references to Closure's goog.js file and
// globalizes. The goog.js file is an ES6 module that forwards some symbols
// on goog from Closure's base.js file, like require.
//
// This pass scans the goog.js file to find what keys are exported. When
// rewriting imports only those in the goog.js file are globalized; if there
// is a reference to something not in the goog.js file it is rewritten in
// such a way to cause an error later at type checking.
//
// This is a separate pass that needs to run before Es6RewriteModules which
// will remove import statements.
//
// This pass enforces the following so that later compiler passes and regex
// based tools can correctly recognize references to these properties of
// goog (like goog.require):
//
//   - No other file is named goog.js
//   - The file is always imported as `import * as goog`
//   - No other module level variable is named `goog`
//   - `export from` is never used on goog.js
//
// Example:
//
// `import * as goog from 'path/to/closure/goog.js'; const myNamespace =
// goog.require('my.namespace');`
//
// Will be rewritten to:
//
// `import 'path/to/closure/goog.js'; const myNamespace =
// goog.require('my.namespace');`
type RewriteGoogJsImports struct {
	mode               RewriteGoogJsImports_Mode
	moduleMap          *ModuleMap
	compiler           AbstractCompiler
	googModule         *Module
	moduleReplacements map[*Module]*Module
}

var _ CompilerPass = (*RewriteGoogJsImports)(nil)

// Java: RewriteGoogJsImports.GOOG_JS_IMPORT_MUST_BE_GOOG_STAR
var RewriteGoogJsImports_GOOG_JS_IMPORT_MUST_BE_GOOG_STAR = DiagnosticType_Error(
	"JSC_GOOG_JS_IMPORT_MUST_BE_GOOG_STAR",
	"Closure''s goog.js file must be imported as `import * as goog`.")

// RewriteGoogJsImports_GOOG_JS_REEXPORTED: since many tools scan for
// "goog.require" ban re-exporting.
var RewriteGoogJsImports_GOOG_JS_REEXPORTED = DiagnosticType_Error("JSC_GOOG_JS_REEXPORTED", "Do not re-export from goog.js.")

// Java: RewriteGoogJsImports.CANNOT_NAME_FILE_GOOG
var RewriteGoogJsImports_CANNOT_NAME_FILE_GOOG = DiagnosticType_Error(
	"JSC_CANNOT_NAME_FILE_GOOG",
	"Do not name files goog.js, it is reserved for Closure Library.")

// Java: RewriteGoogJsImports.CANNOT_HAVE_MODULE_VAR_NAMED_GOOG
var RewriteGoogJsImports_CANNOT_HAVE_MODULE_VAR_NAMED_GOOG = DiagnosticType_Error(
	"JSC_CANNOT_HAVE_MODULE_VAR_NAMED_GOOG",
	"Module scoped variables named ''goog'' must come from importing Closure Library''s "+
		"goog.js file..")

// RewriteGoogJsImports_Mode is the possible traversal modes - either linting
// or linting+rewriting.
type RewriteGoogJsImports_Mode int32

const (
	// Lint only.
	RewriteGoogJsImports_Mode_LINT_ONLY RewriteGoogJsImports_Mode = iota
	// Lint and perform a rewrite on an entire compilation job
	RewriteGoogJsImports_Mode_LINT_AND_REWRITE
)

// Java: RewriteGoogJsImports.Mode.name()
func (m RewriteGoogJsImports_Mode) String() string {
	return [...]string{"LINT_ONLY", "LINT_AND_REWRITE"}[m]
}

const expectedGoogBaseProvide = "goog"

var rewriteGoogJsGoogID = strpool.Intern("goog")

// Java: RewriteGoogJsImports(AbstractCompiler, Mode, ModuleMap)
func NewRewriteGoogJsImports(compiler AbstractCompiler, mode RewriteGoogJsImports_Mode, moduleMap *ModuleMap) *RewriteGoogJsImports {
	jutil.CheckNotNull(moduleMap)
	return &RewriteGoogJsImports{
		mode:               mode,
		moduleMap:          moduleMap,
		compiler:           compiler,
		moduleReplacements: make(map[*Module]*Module),
	}
}

func (r *RewriteGoogJsImports) replacementOrDefault(m *Module) *Module {
	if v, ok := r.moduleReplacements[m]; ok {
		return v
	}
	return m
}

// Java: RewriteGoogJsImports.changeModules()
func (r *RewriteGoogJsImports) changeModules() {
	if len(r.moduleReplacements) == 0 {
		return
	}
	byPath := r.moduleMap.GetModulesByPath()
	byClosure := r.moduleMap.GetModulesByClosureNamespace()
	resolvedModules := NewOrderedStringMapWithCapacity[*Module](byPath.Size())
	closureModules := NewOrderedStringMapWithCapacity[*Module](byClosure.Size())

	for _, e := range byPath.Entries() {
		resolvedModules.Put(e.Key, r.replacementOrDefault(e.Value))
	}

	for _, e := range byClosure.Entries() {
		closureModules.Put(e.Key, r.replacementOrDefault(e.Value))
	}

	r.compiler.SetModuleMap(NewModuleMap(resolvedModules, closureModules))
	clear(r.moduleReplacements)
}

// googJsReferenceReplacer determines if a script imports the goog.js file and
// if so ensures that property accesses are valid and then detaches the import.
type googJsReferenceReplacer struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	outer                  *RewriteGoogJsImports
	hasBadExport           bool
	googImportNode         rhino.Node
	globalizeAllReferences bool
}

// Java: RewriteGoogJsImports.ReferenceReplacer(Node, Node, Module, boolean)
func (r *RewriteGoogJsImports) rewriteReferences(
	script rhino.Node, googImportNode rhino.Node, module *Module, globalizeAllReferences bool,
) {
	p := &googJsReferenceReplacer{
		outer:                  r,
		googImportNode:         googImportNode,
		globalizeAllReferences: globalizeAllReferences,
	}
	NodeTraversal_Traverse(r.compiler, script, p)

	imp := ast.MustImport(googImportNode)
	if bindings, ok := imp.Bindings(); ok {
		if star, ok := bindings.AsImportStar(); ok {
			starNode := star.Node()
			oldBound := module.BoundNamesMap()
			newBindings := NewOrderedStringMapWithCapacity[*Binding](oldBound.Size())
			var googBinding *Binding
			for _, e := range oldBound.Entries() {
				if e.Key == "goog" {
					googBinding = e.Value
					continue
				}
				newBindings.Put(e.Key, jutil.CheckNotNull(e.Value))
			}
			if p.hasBadExport {
				// We might be able to get rid of this case now. Originally this was
				// to cause a type error later in compilation. But now the
				// ModuleMapCreator should log an error for the bad export. However if
				// we're still running we might be in some tool that expects us to be
				// fault tolerant and still rewrite modules. So assume this is needed
				// in those tools.
				starNode.SetOriginalName("goog")
				starNode.SetString("$goog")
				newBindings.Put("$goog", jutil.CheckNotNull(googBinding))
			} else {
				starNode.ReplaceWith(rhino.IR_Empty())
			}
			newModule := module.WithBoundNamesMap(newBindings)
			r.moduleReplacements[module] = newModule
			r.compiler.ReportChangeToEnclosingScope(googImportNode)
		}
	}
}

// Java: RewriteGoogJsImports.ReferenceReplacer.maybeRewriteBadGoogJsImportRef(NodeTraversal, Node, Node)
func (p *googJsReferenceReplacer) maybeRewriteBadGoogJsImportRef(t *NodeTraversal, nameNode rhino.Node, parent rhino.Node) {
	if nodeStringID(nameNode) != rewriteGoogJsGoogID {
		return
	}
	gp, ok := ast.AsGetProp(parent)
	if !ok {
		return
	}

	googVar := t.GetScope().GetVarByID(uint32(rewriteGoogJsGoogID))
	if googVar == nil ||
		googVar.GetNameNode().IsNil() ||
		googVar.GetNameNode().GetParent() != p.googImportNode {
		return
	}

	if p.globalizeAllReferences || p.outer.googModule.NamespaceMap().ContainsKey(gp.Property()) {
		return
	}

	// Rewrite and keep the import so later type checking can report an error
	// that this property does not exist on the ES6 module.
	p.hasBadExport = true
	nameNode.SetOriginalName("goog")
	nameNode.SetString("$goog")
	t.ReportCodeChange()
}

var rewriteGoogJsRefReplacerVisitTokens = func() ast.TokenSet {
	var s ast.TokenSet
	s.Add(rhino.Token_NAME)
	return s
}()

var (
	_ nodeTraversal_visitTokens    = (*googJsReferenceReplacer)(nil)
	_ nodeTraversal_onlyNameLeaves = (*googJsReferenceReplacer)(nil)
)

func (*googJsReferenceReplacer) onlyNameLeaves() {}
func (*googJsReferenceReplacer) visitTokens() *ast.TokenSet {
	return &rewriteGoogJsRefReplacerVisitTokens
}

// Java: RewriteGoogJsImports.ReferenceReplacer.visit(NodeTraversal, Node, Node)
func (p *googJsReferenceReplacer) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	if n.IsName() {
		p.maybeRewriteBadGoogJsImportRef(t, n, parent)
	}
}

// googJsFindReexports finds instances where goog is reexported. Not an
// exhaustive search; does not find alias, e.g. `export const x = goog;`.
type googJsFindReexports struct {
	NodeTraversal_AbstractPreOrderCallbackBase
	hasGoogImport bool
}

// Java: RewriteGoogJsImports.FindReexports(boolean)
func newGoogJsFindReexports(hasGoogImport bool) *googJsFindReexports {
	return &googJsFindReexports{hasGoogImport: hasGoogImport}
}

// Java: RewriteGoogJsImports.FindReexports.checkIfForwardingExport(NodeTraversal, Node)
func (f *googJsFindReexports) checkIfForwardingExport(t *NodeTraversal, export rhino.Node) {
	if source, ok := ast.MustExport(export).Source(); ok && strings.HasSuffix(source.Value(), "/goog.js") {
		t.Report(export, RewriteGoogJsImports_GOOG_JS_REEXPORTED)
	}
}

// Java: RewriteGoogJsImports.FindReexports.checkIfNameFowardedExport(NodeTraversal, Node, Node)
func (f *googJsFindReexports) checkIfNameForwardedExport(t *NodeTraversal, nameNode rhino.Node, parent rhino.Node) {
	if f.hasGoogImport && nodeStringID(nameNode) == rewriteGoogJsGoogID {
		if spec, ok := ast.AsExportSpec(parent); (ok && spec.Local().Is(nameNode)) || parent.IsExport() {
			t.Report(nameNode, RewriteGoogJsImports_GOOG_JS_REEXPORTED)
		}
	}
}

// Java: RewriteGoogJsImports.FindReexports.shouldTraverse(NodeTraversal, Node, Node)
func (f *googJsFindReexports) ShouldTraverse(t *NodeTraversal, n rhino.Node, parent rhino.Node) bool {
	switch n.GetToken() {
	case rhino.Token_ROOT, rhino.Token_SCRIPT, rhino.Token_MODULE_BODY, rhino.Token_EXPORT_SPECS, rhino.Token_EXPORT_SPEC:
		return true
	case rhino.Token_EXPORT:
		f.checkIfForwardingExport(t, n)
		return true
	case rhino.Token_NAME:
		// Visit names in export specs and export defaults.
		f.checkIfNameForwardedExport(t, n, parent)
		return false
	default:
		return false
	}
}

func optNodeSourceFileName(n rhino.Node) base.OptString {
	if s, ok := n.SourceFileName(); ok {
		return base.SomeString(s)
	}
	return base.NoString()
}

// Java: RewriteGoogJsImports.findGoogImportNode(Node)
func (r *RewriteGoogJsImports) findGoogImportNode(scriptRoot rhino.Node) rhino.Node {
	// Cannot use the module map here - information is lost about the
	// imports. The "bound names" could be from transitive imports, but we
	// lose the original import.
	valid := true
	var googImportNode rhino.Node

	moduleBody := ast.MustModuleBody(ast.MustScript(scriptRoot).MustFirstItem().Node())
	for stmt := range moduleBody.Statements() {
		imp, ok := stmt.AsImport()
		if !ok || !strings.HasSuffix(imp.Source().Value(), "/goog.js") {
			continue
		}
		isGoogStar := false
		if _, hasDefault := imp.Default(); !hasDefault {
			if bindings, hasBindings := imp.Bindings(); hasBindings {
				if star, isStar := bindings.AsImportStar(); isStar && star.Ident() == "goog" {
					isGoogStar = true
				}
			}
		}
		child := stmt.Node()
		if isGoogStar {
			googImportNode = child
		} else {
			valid = false
			r.compiler.Report(
				JSError_MakeWithOptSourceLocation(
					optNodeSourceFileName(scriptRoot),
					child.GetLineno(),
					child.GetCharno(),
					RewriteGoogJsImports_GOOG_JS_IMPORT_MUST_BE_GOOG_STAR))
		}
	}

	if !valid {
		return rhino.NilNode
	} else if !googImportNode.IsNil() {
		return googImportNode
	}

	moduleScope :=
		NewSyntacticScopeCreator(r.compiler).
			CreateScopeOfScope(moduleBody.Node(), Scope_CreateGlobalScope(scriptRoot))
	googVar := moduleScope.GetVarByID(uint32(rewriteGoogJsGoogID))

	if googVar != nil && !googVar.GetNameNode().IsNil() {
		r.compiler.Report(
			JSError_MakeWithOptSourceLocation(
				optNodeSourceFileName(scriptRoot),
				googVar.GetNameNode().GetLineno(),
				googVar.GetNameNode().GetCharno(),
				RewriteGoogJsImports_CANNOT_HAVE_MODULE_VAR_NAMED_GOOG))
	}

	return rhino.NilNode
}

// Java: RewriteGoogJsImports.rewriteImports(Node)
func (r *RewriteGoogJsImports) rewriteImports(scriptRoot rhino.Node) {
	if Es6RewriteModules_IsEs6ModuleRoot(scriptRoot) {
		googImportNode := r.findGoogImportNode(scriptRoot)
		NodeTraversal_Traverse(r.compiler, scriptRoot, newGoogJsFindReexports(!googImportNode.IsNil()))
		module := r.moduleMap.GetModuleByPath(r.compiler.GetInput(scriptRoot.GetInputId()).GetPath())
		jutil.CheckNotNull(module)

		if !googImportNode.IsNil() && r.mode == RewriteGoogJsImports_Mode_LINT_AND_REWRITE {
			// If googModule is null then goog.js was not part of the input. Try
			// to be fault tolerant and just assume that everything exported is
			// on the global goog.
			r.rewriteReferences(
				scriptRoot,
				googImportNode,
				module,
				/* globalizeAllReferences= */ r.googModule == nil)
		}
	}
}

// Java: RewriteGoogJsImports.findGoogJsScriptNode(Node)
func (r *RewriteGoogJsImports) findGoogJsScriptNode(root rhino.Node) rhino.Node {
	if root.IsNil() {
		return rhino.NilNode
	}
	var expectedGoogPath *ModuleLoader_ModulePath

	// Find Closure's base.js file. goog.js should be right next to it.
	for item := range ast.MustRoot(root).Items() {
		script := item.Node()
		inp := r.compiler.GetInput(script.GetInputId())
		if slices.Contains(inp.GetProvides(), expectedGoogBaseProvide) {
			// Use resolveModuleAsPath as if it is not part of the input we don't
			// want to report an error.
			expectedGoogPath = inp.GetPath().ResolveModuleAsPath("./goog.js")
			break
		}
	}

	if expectedGoogPath == nil {
		return rhino.NilNode
	}

	var googScriptNode rhino.Node
	for item := range ast.MustRoot(root).Items() {
		script := item.Node()
		if r.compiler.
			GetInput(script.GetInputId()).
			GetPath().
			EqualsIgnoreLeadingSlash(expectedGoogPath) {
			googScriptNode = script
		} else if srcFileName, ok := script.SourceFileName(); ok && strings.HasSuffix(srcFileName, "/goog.js") {
			// Ban the name goog.js as input except for Closure's goog.js file.
			// This simplifies a lot of logic if the only file that is allowed
			// to be named goog.js is Closure's.
			r.compiler.Report(JSError_MakeAtSourceLocation(srcFileName, -1, -1, RewriteGoogJsImports_CANNOT_NAME_FILE_GOOG))
		}
	}
	return googScriptNode
}

// Java: RewriteGoogJsImports.process(Node, Node)
func (r *RewriteGoogJsImports) Process(externs rhino.Node, root rhino.Node) {
	r.googModule = nil

	googJsScriptNode := r.findGoogJsScriptNode(root)
	if googJsScriptNode.IsNil() {
		// Potentially in externs if library level type checking.
		googJsScriptNode = r.findGoogJsScriptNode(externs)
	}

	if r.mode == RewriteGoogJsImports_Mode_LINT_AND_REWRITE {
		if !googJsScriptNode.IsNil() {
			googJsPath := r.compiler.GetInput(googJsScriptNode.GetInputId()).GetPath()
			r.googModule = r.moduleMap.GetModuleByPath(googJsPath)
			jutil.CheckNotNull(r.googModule)
			isFromGoog := func(b *Binding) bool { return b.OriginatingExport().ModulePath() == googJsPath }
			jutil.CheckStateMsg(
				allBindingsMatch(r.googModule.BoundNamesMap(), isFromGoog),
				"goog.js should never import anything")
			jutil.CheckStateMsg(
				!r.googModule.NamespaceMap().ContainsKey("default"),
				"goog.js should never have a default export.")
			jutil.CheckStateMsg(
				allBindingsMatch(r.googModule.NamespaceMap(), isFromGoog),
				"goog.js should never export from anything.")
		}
	} else {
		jutil.CheckState(r.mode == RewriteGoogJsImports_Mode_LINT_ONLY)
	}

	if !root.IsNil() {
		for item := range ast.MustRoot(root).Items() {
			r.rewriteImports(item.Node())
		}
	}

	r.changeModules()
}

// allBindingsMatch is map.values().stream().allMatch(pred).
func allBindingsMatch(m *OrderedStringMap[*Binding], pred func(*Binding) bool) bool {
	for _, e := range m.Entries() {
		if !pred(e.Value) {
			return false
		}
	}
	return true
}
