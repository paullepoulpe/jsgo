// Port of com/google/javascript/jscomp/RewriteGoogJsImports.java @ closure-compiler v20261015

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

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
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
	mode       RewriteGoogJsImports_Mode
	moduleMap  *ModuleMap
	compiler   AbstractCompiler
	googModule *Module
	// Module has value equality (AutoValue), hence the Object-keyed map.
	moduleReplacements *jutil.ObjectLinkedHashMap[*Module, *Module]
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

// Java: RewriteGoogJsImports.Mode.values()
func RewriteGoogJsImports_Mode_values() []RewriteGoogJsImports_Mode {
	return []RewriteGoogJsImports_Mode{RewriteGoogJsImports_Mode_LINT_ONLY, RewriteGoogJsImports_Mode_LINT_AND_REWRITE}
}

// Java: RewriteGoogJsImports.Mode.name()
func (m RewriteGoogJsImports_Mode) String() string {
	return [...]string{"LINT_ONLY", "LINT_AND_REWRITE"}[m]
}

const rewriteGoogJsImports_EXPECTED_BASE_PROVIDE = "goog"

// Java: RewriteGoogJsImports(AbstractCompiler, Mode, ModuleMap)
func NewRewriteGoogJsImports(compiler AbstractCompiler, mode RewriteGoogJsImports_Mode, moduleMap *ModuleMap) *RewriteGoogJsImports {
	jutil.CheckNotNull(moduleMap)
	r := &RewriteGoogJsImports{moduleReplacements: jutil.NewObjectLinkedHashMap[*Module, *Module]()}
	r.compiler = compiler
	r.mode = mode
	r.moduleMap = moduleMap
	return r
}

// rewriteGoogJsImports_getOrDefault is Map.getOrDefault(Object, V).
//
// Go-only: language-adapter — jutil.ObjectLinkedHashMap has Lookup.
func (r *RewriteGoogJsImports) replacementOrDefault(m *Module) *Module {
	if v, ok := r.moduleReplacements.Lookup(m); ok {
		return v
	}
	return m
}

// Java: RewriteGoogJsImports.changeModules()
func (r *RewriteGoogJsImports) ChangeModules() {
	resolvedModules := jutil.NewLinkedHashMap[string, *Module]()
	closureModules := jutil.NewLinkedHashMap[string, *Module]()

	for k, v := range r.moduleMap.GetModulesByPath().All() {
		newModule := r.replacementOrDefault(v)
		resolvedModules.Put(k, newModule)
	}

	for k, v := range r.moduleMap.GetModulesByClosureNamespace().All() {
		newModule := r.replacementOrDefault(v)
		closureModules.Put(k, newModule)
	}

	r.compiler.SetModuleMap(NewModuleMap(resolvedModules, closureModules))
	// Go-only: language-adapter — moduleReplacements.clear():
	// jutil.ObjectLinkedHashMap has no Clear, and the map is private.
	r.moduleReplacements = jutil.NewObjectLinkedHashMap[*Module, *Module]()
}

// RewriteGoogJsImports_ReferenceReplacer determines if a script imports the
// goog.js file and if so ensures that property accesses are valid and then
// detaches the import.
//
// It is an inner class: outer is the enclosing RewriteGoogJsImports.
type RewriteGoogJsImports_ReferenceReplacer struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	outer                  *RewriteGoogJsImports
	hasBadExport           bool
	googImportNode         rhino.Node
	globalizeAllReferences bool
}

// Java: RewriteGoogJsImports.ReferenceReplacer(Node, Node, Module, boolean)
func (r *RewriteGoogJsImports) NewReferenceReplacer(
	script rhino.Node, googImportNode rhino.Node, module *Module, globalizeAllReferences bool) *RewriteGoogJsImports_ReferenceReplacer {
	p := &RewriteGoogJsImports_ReferenceReplacer{outer: r}
	p.googImportNode = googImportNode
	p.globalizeAllReferences = globalizeAllReferences
	NodeTraversal_Traverse(r.compiler, script, p)

	if googImportNode.GetSecondChild().IsImportStar() {
		newBindings := jutil.NewLinkedHashMapFrom(module.BoundNames())
		newBindings.Remove("goog")
		if p.hasBadExport {
			// We might be able to get rid of this case now. Originally this was
			// to cause a type error later in compilation. But now the
			// ModuleMapCreator should log an error for the bad export. However if
			// we're still running we might be in some tool that expects us to be
			// fault tolerant and still rewrite modules. So assume this is needed
			// in those tools.
			googImportNode.GetSecondChild().SetOriginalName(jutil.StringPtr("goog"))
			googImportNode.GetSecondChild().SetString("$goog")
			newBindings.Put("$goog", module.BoundNames().Get("goog"))
		} else {
			googImportNode.GetSecondChild().ReplaceWith(rhino.IR_Empty())
		}
		newModule := module.ToBuilder().BoundNames(closureModuleProcessor_immutableMapCopyOf(newBindings)).Build()
		r.moduleReplacements.Put(module, newModule)
		r.compiler.ReportChangeToEnclosingScope(googImportNode)
	}
	return p
}

// Java: RewriteGoogJsImports.ReferenceReplacer.maybeRewriteBadGoogJsImportRef(NodeTraversal, Node, Node)
func (p *RewriteGoogJsImports_ReferenceReplacer) MaybeRewriteBadGoogJsImportRef(t *NodeTraversal, nameNode rhino.Node, parent rhino.Node) {
	if !parent.IsGetProp() || nameNode.GetString() != "goog" {
		return
	}

	var_ := t.GetScope().GetVar(nameNode.GetString())

	if var_ == nil ||
		var_.GetNameNode().IsNil() ||
		var_.GetNameNode().GetParent() != p.googImportNode {
		return
	}

	if p.globalizeAllReferences || p.outer.googModule.Namespace().ContainsKey(parent.GetString()) {
		return
	}

	// Rewrite and keep the import so later type checking can report an error
	// that this property does not exist on the ES6 module.
	p.hasBadExport = true
	nameNode.SetOriginalName(jutil.StringPtr("goog"))
	nameNode.SetString("$goog")
	t.ReportCodeChange()

}

// Java: RewriteGoogJsImports.ReferenceReplacer.visit(NodeTraversal, Node, Node)
func (p *RewriteGoogJsImports_ReferenceReplacer) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	switch n.GetToken() {
	case rhino.Token_NAME:
		p.MaybeRewriteBadGoogJsImportRef(t, n, parent)
	default:
	}
}

// RewriteGoogJsImports_FindReexports finds instances where goog is
// reexported. Not an exhaustive search; does not find alias, e.g. `export
// const x = goog;`. Not to say that people should export using aliases; they
// shouldn't. Just that we're catching the 99% case and that should be an
// indication to people not to do bad things.
type RewriteGoogJsImports_FindReexports struct {
	NodeTraversal_AbstractPreOrderCallbackBase
	hasGoogImport bool
}

// Java: RewriteGoogJsImports.FindReexports(boolean)
func NewRewriteGoogJsImports_FindReexports(hasGoogImport bool) *RewriteGoogJsImports_FindReexports {
	return &RewriteGoogJsImports_FindReexports{hasGoogImport: hasGoogImport}
}

// Java: RewriteGoogJsImports.FindReexports.checkIfForwardingExport(NodeTraversal, Node)
func (f *RewriteGoogJsImports_FindReexports) CheckIfForwardingExport(t *NodeTraversal, export rhino.Node) {
	if export.HasTwoChildren() && strings.HasSuffix(export.GetLastChild().GetString(), "/goog.js") {
		t.Report(export, RewriteGoogJsImports_GOOG_JS_REEXPORTED)
	}
}

// Java: RewriteGoogJsImports.FindReexports.checkIfNameFowardedExport(NodeTraversal, Node, Node)
func (f *RewriteGoogJsImports_FindReexports) CheckIfNameFowardedExport(t *NodeTraversal, nameNode rhino.Node, parent rhino.Node) {
	if f.hasGoogImport && nameNode.GetString() == "goog" {
		if (parent.IsExportSpec() && parent.GetFirstChild() == nameNode) || parent.IsExport() {
			t.Report(nameNode, RewriteGoogJsImports_GOOG_JS_REEXPORTED)
		}
	}
}

// Java: RewriteGoogJsImports.FindReexports.shouldTraverse(NodeTraversal, Node, Node)
func (f *RewriteGoogJsImports_FindReexports) ShouldTraverse(t *NodeTraversal, n rhino.Node, parent rhino.Node) bool {
	switch n.GetToken() {
	case rhino.Token_ROOT, rhino.Token_SCRIPT, rhino.Token_MODULE_BODY, rhino.Token_EXPORT_SPECS, rhino.Token_EXPORT_SPEC:
		return true
	case rhino.Token_EXPORT:
		f.CheckIfForwardingExport(t, n)
		return true
	case rhino.Token_NAME:
		// Visit names in export specs and export defaults.
		f.CheckIfNameFowardedExport(t, n, parent)
		return false
	default:
		return false
	}
}

// Java: RewriteGoogJsImports.findGoogImportNode(Node)
func (r *RewriteGoogJsImports) FindGoogImportNode(scriptRoot rhino.Node) rhino.Node {
	// Cannot use the module map here - information is lost about the
	// imports. The "bound names" could be from transitive imports, but we
	// lose the original import.
	valid := true
	var googImportNode rhino.Node

	for child := scriptRoot.GetFirstFirstChild(); !child.IsNil(); child = child.GetNext() {
		if child.IsImport() && strings.HasSuffix(child.GetLastChild().GetString(), "/goog.js") {
			if child.GetFirstChild().IsEmpty() &&
				child.GetSecondChild().IsImportStar() &&
				child.GetSecondChild().GetString() == "goog" {
				googImportNode = child
			} else {
				valid = false
				r.compiler.Report(
					JSError_MakeWithSourceLocation(
						scriptRoot.GetSourceFileName(),
						child.GetLineno(),
						child.GetCharno(),
						RewriteGoogJsImports_GOOG_JS_IMPORT_MUST_BE_GOOG_STAR))
			}
		}
	}

	if !valid {
		return rhino.NilNode
	} else if !googImportNode.IsNil() {
		return googImportNode
	}

	moduleScope :=
		NewSyntacticScopeCreator(r.compiler).
			CreateScopeOfScope(scriptRoot.GetFirstChild(), Scope_CreateGlobalScope(scriptRoot))
	googVar := moduleScope.GetVar("goog")

	if googVar != nil && !googVar.GetNameNode().IsNil() {
		r.compiler.Report(
			JSError_MakeWithSourceLocation(
				scriptRoot.GetSourceFileName(),
				googVar.GetNameNode().GetLineno(),
				googVar.GetNameNode().GetCharno(),
				RewriteGoogJsImports_CANNOT_HAVE_MODULE_VAR_NAMED_GOOG))
	}

	return rhino.NilNode
}

// Java: RewriteGoogJsImports.rewriteImports(Node)
func (r *RewriteGoogJsImports) RewriteImports(scriptRoot rhino.Node) {
	if Es6RewriteModules_IsEs6ModuleRoot(scriptRoot) {
		googImportNode := r.FindGoogImportNode(scriptRoot)
		NodeTraversal_Traverse(r.compiler, scriptRoot, NewRewriteGoogJsImports_FindReexports(!googImportNode.IsNil()))
		module := r.moduleMap.GetModuleByPath(r.compiler.GetInput(scriptRoot.GetInputId()).GetPath())
		jutil.CheckNotNull(module)

		if !googImportNode.IsNil() && r.mode == RewriteGoogJsImports_Mode_LINT_AND_REWRITE {
			// If googModule is null then goog.js was not part of the input. Try
			// to be fault tolerant and just assume that everything exported is
			// on the global goog.
			_ = r.NewReferenceReplacer(
				scriptRoot,
				googImportNode,
				module,
				/* globalizeAllReferences= */ r.googModule == nil)
		}
	}
}

// Java: RewriteGoogJsImports.findGoogJsScriptNode(Node)
func (r *RewriteGoogJsImports) FindGoogJsScriptNode(root rhino.Node) rhino.Node {
	var expectedGoogPath *ModuleLoader_ModulePath

	// Find Closure's base.js file. goog.js should be right next to it.
	for script := root.GetFirstChild(); !script.IsNil(); script = script.GetNext() {
		provides := r.compiler.GetInput(script.GetInputId()).GetProvides()
		if slices.Contains(provides, rewriteGoogJsImports_EXPECTED_BASE_PROVIDE) {
			// Use resolveModuleAsPath as if it is not part of the input we don't
			// want to report an error.
			expectedGoogPath =
				r.compiler.GetInput(script.GetInputId()).GetPath().ResolveModuleAsPath("./goog.js")
			break
		}
	}

	if expectedGoogPath != nil {
		var googScriptNode rhino.Node

		for script := root.GetFirstChild(); !script.IsNil(); script = script.GetNext() {
			if r.compiler.
				GetInput(script.GetInputId()).
				GetPath().
				EqualsIgnoreLeadingSlash(expectedGoogPath) {
				googScriptNode = script
			} else if strings.HasSuffix(*script.GetSourceFileName(), "/goog.js") {
				// Ban the name goog.js as input except for Closure's goog.js file.
				// This simplifies a lot of logic if the only file that is allowed
				// to be named goog.js is Closure's.
				r.compiler.Report(JSError_MakeWithSourceLocation(script.GetSourceFileName(), -1, -1, RewriteGoogJsImports_CANNOT_NAME_FILE_GOOG))
			}
		}

		return googScriptNode
	}

	return rhino.NilNode
}

// Java: RewriteGoogJsImports.process(Node, Node)
func (r *RewriteGoogJsImports) Process(externs rhino.Node, root rhino.Node) {
	r.googModule = nil

	googJsScriptNode := r.FindGoogJsScriptNode(root)

	if googJsScriptNode.IsNil() {
		// Potentially in externs if library level type checking.
		googJsScriptNode = r.FindGoogJsScriptNode(externs)
	}

	if r.mode == RewriteGoogJsImports_Mode_LINT_AND_REWRITE {
		if !googJsScriptNode.IsNil() {
			googJsPath := r.compiler.GetInput(googJsScriptNode.GetInputId()).GetPath()
			r.googModule = r.moduleMap.GetModuleByPath(googJsPath)
			jutil.CheckNotNull(r.googModule)
			isFromGoog := func(b *Binding) bool { return b.OriginatingExport().ModulePath() == googJsPath }
			jutil.CheckStateMsg(
				rewriteGoogJsImports_allMatch(r.googModule.BoundNames(), isFromGoog),
				"goog.js should never import anything")
			jutil.CheckStateMsg(
				!r.googModule.Namespace().ContainsKey("default"),
				"goog.js should never have a default export.")
			jutil.CheckStateMsg(
				rewriteGoogJsImports_allMatch(r.googModule.Namespace(), isFromGoog),
				"goog.js should never export from anything.")
		}
	} else {
		jutil.CheckState(r.mode == RewriteGoogJsImports_Mode_LINT_ONLY)
	}

	for script := root.GetFirstChild(); !script.IsNil(); script = script.GetNext() {
		r.RewriteImports(script)
	}

	r.ChangeModules()
}

// rewriteGoogJsImports_allMatch is map.values().stream().allMatch(pred).
//
// Go-only: language-adapter — Java streams.
func rewriteGoogJsImports_allMatch(m *jutil.LinkedHashMap[string, *Binding], pred func(*Binding) bool) bool {
	for b := range m.Values() {
		if !pred(b) {
			return false
		}
	}
	return true
}
