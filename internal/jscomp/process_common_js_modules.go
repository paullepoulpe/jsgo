// Port of com/google/javascript/jscomp/ProcessCommonJSModules.java @ closure-compiler v20261021

/*
 * Copyright 2011 The Closure Compiler Authors.
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
	"strconv"
	"strings"

	"github.com/paullepoulpe/jsgo/internal/ast"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// Go-only: this pass uses the typed AST view (internal/ast) where it reads
// the shapes of the module code it recognizes: the callee and arguments of
// require(...), __webpack_require__(...), require.ensure(...) and
// Promise.all(...) calls, the IIFE and UMD factory wrappers (the statement,
// the call, its function and `.call` receiver), the parameters of a
// webpack module shim, the `module['exports']` element accesses, the
// assignments of exports and their object literals, the test, then and else
// of a UMD `if` or `?:`, the destructuring of a require, and the
// declarations it splits. The port no longer follows the Java line by line.
//
// Java quirks kept: the require.ensure warning counts the callee in the
// "found {0}" argument count, and the webpack module shim reads the
// argument of the `module` parameter by child index, one place earlier for
// a free call.
//
// Left on the rhino API: the per-node visits (they match qualified names
// on any node), the JSDoc type rewriting, and the name rewriting of
// RewriteModule (updateNameReference, visitExport, getExportedNameNode),
// which reads the siblings and parents of a name in whatever expression
// holds it.

// ProcessCommonJSModules rewrites a CommonJS module
// http://wiki.commonjs.org/wiki/Modules/1.1.1 into a form that can be safely
// concatenated. Does not add a function around the module body but instead
// adds suffixes to global variables to avoid conflicts. Calls to require are
// changed to reference the required module directly.
type ProcessCommonJSModules struct {
	NodeTraversal_AbstractPreOrderCallbackBase

	compiler AbstractCompiler
}

const processCommonJSModules_EXPORTS = "exports"
const processCommonJSModules_MODULE = "module"
const processCommonJSModules_REQUIRE = "require"
const processCommonJSModules_WEBPACK_REQUIRE = "__webpack_require__"

// Webpack transpiles import() statements to __webpack_require__.t(modulePath)
// Always imports the module namespace regardless of module type
const processCommonJSModules_WEBPACK_REQUIRE_NAMESPACE = "__webpack_require__.t"
const processCommonJSModules_EXPORT_PROPERTY_NAME = "default"

var ProcessCommonJSModules_UNKNOWN_REQUIRE_ENSURE = DiagnosticType_Warning(
	"JSC_COMMONJS_UNKNOWN_REQUIRE_ENSURE_ERROR", "Unrecognized require.ensure call: {0}")

var ProcessCommonJSModules_SUSPICIOUS_EXPORTS_ASSIGNMENT = DiagnosticType_Warning(
	"JSC_COMMONJS_SUSPICIOUS_EXPORTS_ASSIGNMENT",
	"Suspicious re-assignment of \"exports\" variable."+
		" Did you actually intend to export something?")

// NewProcessCommonJSModules creates a new ProcessCommonJSModules instance
// which can be used to rewrite CommonJS modules to a concatenable form.
//
// Java: ProcessCommonJSModules(AbstractCompiler)
func NewProcessCommonJSModules(compiler AbstractCompiler) *ProcessCommonJSModules {
	p := &ProcessCommonJSModules{}
	p.compiler = compiler
	return p
}

// Java: ProcessCommonJSModules.process(Node, Node)
func (p *ProcessCommonJSModules) Process(externs rhino.Node, root rhino.Node) {
	NodeTraversal_Traverse(p.compiler, root, p)
	NodeTraversal_Traverse(p.compiler, externs, p)
}

// Java: ProcessCommonJSModules.shouldTraverse(NodeTraversal, Node, Node)
func (p *ProcessCommonJSModules) ShouldTraverse(t *NodeTraversal, n rhino.Node, parent rhino.Node) bool {
	if n.IsRoot() {
		return true
	} else if n.IsScript() {
		if p.compiler.GetOptions().GetModuleResolutionMode() == ModuleLoader_ResolutionMode_WEBPACK {
			p.RemoveWebpackModuleShim(n)
		}

		finder := NewProcessCommonJSModules_FindImportsAndExports(p)
		moduleLoaderErrorHandler := p.compiler.GetModuleLoader().GetErrorHandler()
		p.compiler.GetModuleLoader().SetErrorHandler(finder)
		NodeTraversal_Traverse(p.compiler, n, finder)

		moduleType := p.compiler.GetInput(n.GetInputId()).GetJsModuleType()

		forceModuleDetection := moduleType == CompilerInput_ModuleType_IMPORTED_SCRIPT
		defaultExportIsConst := true

		isCommonJsModule := finder.IsCommonJsModule()
		var exports []*ProcessCommonJSModules_ExportInfo
		if isCommonJsModule || forceModuleDetection {
			if len(finder.umdPatterns) != 0 {
				needsRetraverse := false
				if finder.ReplaceUmdPatterns() {
					needsRetraverse = true
				}
				// Removing the IIFE rewrites vars. We need to re-traverse
				// to get the new references.
				if p.RemoveIIFEWrapper(n) {
					needsRetraverse = true
				}

				if needsRetraverse {
					finder = NewProcessCommonJSModules_FindImportsAndExports(p)
					p.compiler.GetModuleLoader().SetErrorHandler(finder)
					NodeTraversal_Traverse(p.compiler, n, finder)
				}
			}

			defaultExportIsConst = finder.InitializeModule()

			exports = append(exports, finder.GetModuleExports()...)
			exports = append(exports, finder.GetExports()...)
		}
		finder.ReportModuleErrors()
		p.compiler.GetModuleLoader().SetErrorHandler(moduleLoaderErrorHandler)

		NodeTraversal_Traverse(
			p.compiler,
			n,
			NewProcessCommonJSModules_RewriteModule(
				p, isCommonJsModule || forceModuleDetection, exports, defaultExportIsConst))
	}
	return false
}

// Java: ProcessCommonJSModules.getModuleName(CompilerInput)
func ProcessCommonJSModules_GetModuleName(input *CompilerInput) *string {
	modulePath := input.GetPath()
	if modulePath == nil {
		return nil
	}

	return jutil.StringPtr(ProcessCommonJSModules_GetModuleNameWithModulePath(modulePath))
}

// Java: ProcessCommonJSModules.getModuleName(ModulePath)
func ProcessCommonJSModules_GetModuleNameWithModulePath(input *ModuleLoader_ModulePath) string {
	return input.ToModuleName()
}

// Java: ProcessCommonJSModules.getBasePropertyImport(String, Node)
func (p *ProcessCommonJSModules) GetBasePropertyImportWithRequireCall(moduleName string, requireCall rhino.Node) string {
	jutil.CheckArgument(requireCall.IsCall())
	if p.compiler.GetOptions().GetModuleResolutionMode() == ModuleLoader_ResolutionMode_WEBPACK &&
		requireCall.GetFirstChild().MatchesQualifiedName(processCommonJSModules_WEBPACK_REQUIRE_NAMESPACE) {
		return moduleName
	}

	return p.GetBasePropertyImport(moduleName)
}

// Java: ProcessCommonJSModules.getBasePropertyImport(String)
func (p *ProcessCommonJSModules) GetBasePropertyImport(moduleName string) string {
	moduleType := p.compiler.GetModuleTypeByName(moduleName)
	if moduleType != nil && *moduleType != CompilerInput_ModuleType_COMMONJS {
		return moduleName
	}

	return moduleName + "." + processCommonJSModules_EXPORT_PROPERTY_NAME
}

// Java: ProcessCommonJSModules.isCommonJsImport(Node)
func (p *ProcessCommonJSModules) IsCommonJsImport(requireCall rhino.Node) bool {
	return ProcessCommonJSModules_IsCommonJsImport(requireCall, p.compiler.GetOptions().GetModuleResolutionMode())
}

// ProcessCommonJSModules_IsCommonJsImport recognizes if a node is a module
// import. We recognize two forms:
//
//   - require("something");
//   - __webpack_require__(4); // only when the module resolution is WEBPACK
//   - __webpack_require__.t(4); // only when the module resolution is WEBPACK
//
// Java: ProcessCommonJSModules.isCommonJsImport(Node, ModuleLoader.ResolutionMode)
func ProcessCommonJSModules_IsCommonJsImport(
	requireCall rhino.Node, resolutionMode ModuleLoader_ResolutionMode) bool {
	call, ok := ast.AsCall(requireCall)
	if !ok {
		return false
	}
	first, ok := call.FirstArg()
	if !ok {
		return false
	}
	callee := call.Callee().Node()
	second, hasSecond := first.Next()
	if !hasSecond {
		// require('path'), __webpack_require__(4), __webpack_require__.t(4)
		path := first.Node()
		if resolutionMode == ModuleLoader_ResolutionMode_WEBPACK &&
			(callee.MatchesQualifiedName(processCommonJSModules_WEBPACK_REQUIRE) ||
				callee.MatchesQualifiedName(processCommonJSModules_WEBPACK_REQUIRE_NAMESPACE)) &&
			(path.IsNumber() || path.IsStringLit()) {
			return true
		}
		return callee.MatchesQualifiedName(processCommonJSModules_REQUIRE) && path.IsStringLit()
	}
	// __webpack_require__.bind(null, 4)
	path := second.Node()
	return second.IsLast() &&
		resolutionMode == ModuleLoader_ResolutionMode_WEBPACK &&
		callee.MatchesQualifiedName(processCommonJSModules_WEBPACK_REQUIRE+".bind") &&
		first.Node().IsNull() &&
		(path.IsNumber() || path.IsStringLit())
}

// Java: ProcessCommonJSModules.getCommonJsImportPath(Node)
func (p *ProcessCommonJSModules) GetCommonJsImportPath(requireCall rhino.Node) string {
	return ProcessCommonJSModules_GetCommonJsImportPath(requireCall, p.compiler.GetOptions().GetModuleResolutionMode())
}

// Java: ProcessCommonJSModules.getCommonJsImportPath(Node, ModuleLoader.ResolutionMode)
//
// Go-only: requireCall is a call with at least one argument
// (isCommonJsImport).
func ProcessCommonJSModules_GetCommonJsImportPath(
	requireCall rhino.Node, resolutionMode ModuleLoader_ResolutionMode) string {
	first, ok := ast.MustCall(requireCall).FirstArg()
	jutil.CheckState(ok)
	if resolutionMode == ModuleLoader_ResolutionMode_WEBPACK {
		// The path is the second argument of __webpack_require__.bind(null, path).
		pathArgument := first.Node()
		if second, ok := first.Next(); ok {
			pathArgument = second.Node()
		}
		if pathArgument.IsNumber() {
			return strconv.Itoa(int(jutil.D2I(pathArgument.GetDouble())))
		} else {
			return pathArgument.GetString()
		}
	}

	return first.Node().GetString()
}

// Java: ProcessCommonJSModules.getImportedModuleName(NodeTraversal, Node)
func (p *ProcessCommonJSModules) GetImportedModuleName(t *NodeTraversal, requireCall rhino.Node) string {
	return p.GetImportedModuleNameWithImportPath(t, requireCall, p.GetCommonJsImportPath(requireCall))
}

// Java: ProcessCommonJSModules.getImportedModuleName(NodeTraversal, Node, String)
func (p *ProcessCommonJSModules) GetImportedModuleNameWithImportPath(t *NodeTraversal, n rhino.Node, importPath string) string {
	modulePath :=
		t.GetInput().
			GetPath().
			ResolveJsModule(importPath, processCommonJSModules_sourceFileName(n), n.GetLineno(), n.GetCharno())

	if modulePath == nil {
		return ModuleIdentifier_ForFile(importPath).ModuleName()
	}
	return modulePath.ToModuleName()
}

// processCommonJSModules_sourceFileName is n.getSourceFileName() as the
// String argument of ModulePath.resolveJsModule.
//
// Go-only: language-adapter — Node.getSourceFileName() is a nullable
// String; ModulePath.resolveJsModule takes a Go string.
func processCommonJSModules_sourceFileName(n rhino.Node) string {
	if s := n.GetSourceFileName(); s != nil {
		return *s
	}
	return ""
}

// ProcessCommonJSModules_IsCommonJsExport recognizes if a node is a module
// export. We recognize several forms:
//
//   - module.exports = something;
//   - module.exports.something = something;
//   - exports.something = something;
//
// In addition, we only recognize an export if the base export object is not
// defined or is defined in externs.
//
// Java: ProcessCommonJSModules.isCommonJsExport(NodeTraversal, Node, ModuleLoader.ResolutionMode)
func ProcessCommonJSModules_IsCommonJsExport(
	t *NodeTraversal, export rhino.Node, resolutionMode ModuleLoader_ResolutionMode) bool {
	if export.MatchesQualifiedName(processCommonJSModules_MODULE+"."+processCommonJSModules_EXPORTS) ||
		processCommonJSModules_isModuleExportsElement(export) {
		v := t.GetScope().GetVar(processCommonJSModules_MODULE)
		if v == nil || v.IsExtern() {
			return true
		}
	} else if export.IsName() && processCommonJSModules_EXPORTS == export.GetString() {
		v := t.GetScope().GetVar(export.GetString())
		if v == nil || v.IsGlobal() {
			return true
		}
	}
	return false
}

// processCommonJSModules_isModuleExportsElement reports whether n is
// `module['exports']`.
//
// Go-only: Java tests the element access inline, in isCommonJsExport and
// in FindImportsAndExports.visit.
func processCommonJSModules_isModuleExportsElement(n rhino.Node) bool {
	getElem, ok := ast.AsGetElem(n)
	if !ok || !getElem.Object().Node().MatchesQualifiedName(processCommonJSModules_MODULE) {
		return false
	}
	index, ok := getElem.Index().AsStringLit()
	return ok && index.Value() == processCommonJSModules_EXPORTS
}

// Java: ProcessCommonJSModules.isCommonJsExport(NodeTraversal, Node)
func (p *ProcessCommonJSModules) IsCommonJsExport(t *NodeTraversal, export rhino.Node) bool {
	return ProcessCommonJSModules_IsCommonJsExport(
		t, export, p.compiler.GetOptions().GetModuleResolutionMode())
}

// ProcessCommonJSModules_IsCommonJsDynamicImportCallback recognizes if a
// node is a dynamic module import. Currently only the webpack dynamic import
// is recognized:
//
//   - __webpack_require__.e(0).then(function() { return __webpack_require__(4);})
//   - Promise.all([__webpack_require__.e(0)]).then(function() { return
//     __webpack_require__(4);})
//
// Java: ProcessCommonJSModules.isCommonJsDynamicImportCallback(Node, ModuleLoader.ResolutionMode)
func ProcessCommonJSModules_IsCommonJsDynamicImportCallback(
	n rhino.Node, resolutionMode ModuleLoader_ResolutionMode) bool {
	if n.IsNil() || resolutionMode != ModuleLoader_ResolutionMode_WEBPACK {
		return false
	}
	if n.IsFunction() && ProcessCommonJSModules_IsWebpackRequireEnsureCallback(n) {
		return true
	}

	return false
}

var processCommonJSModules_PROMISE_ALL = rhino.QualifiedName_Of("Promise.all")

// ProcessCommonJSModules_IsWebpackRequireEnsureCallback recognizes
// __webpack_require__ calls that are the .then callback of a
// __webpack_require__.e call. Example:
//
//   - __webpack_require__.e(0).then(function() { return __webpack_require__(4); })
//   - Promise.all([__webpack_require__.e(0)]).then(function() { return
//     __webpack_require__(4);})
//
// Java: ProcessCommonJSModules.isWebpackRequireEnsureCallback(Node)
func ProcessCommonJSModules_IsWebpackRequireEnsureCallback(fnc rhino.Node) bool {
	jutil.CheckArgument(fnc.IsFunction())
	if fnc.GetParent().IsNil() {
		return false
	}

	// `x(...).then(fnc)`, where x is __webpack_require__.e or Promise.all
	callParent, ok := ast.AsCall(fnc.GetParent())
	if !ok {
		return false
	}
	then, ok := callParent.Callee().AsGetProp()
	if !ok {
		return false
	}
	ensureCall, ok := then.Object().AsCall()
	if !ok {
		return false
	}
	callParentTarget := ensureCall.Callee().Node()

	if callParentTarget.MatchesQualifiedName(processCommonJSModules_WEBPACK_REQUIRE+".e") &&
		then.Property() == "then" {
		return true
	} else if processCommonJSModules_PROMISE_ALL.Matches(callParentTarget) {
		// Promise.all([__webpack_require__.e(0), ...])
		firstArg, _ := ensureCall.FirstArg()
		array, ok := ast.AsArrayLit(firstArg.Node())
		if !ok {
			return false
		}
		allElementsAreDynamicImports := false
		for arrayItem := range array.Elements() {
			call, ok := ast.AsCall(arrayItem.Node())
			if !ok {
				return false
			}
			if _, ok := call.OnlyArg(); !ok ||
				!call.Callee().Node().MatchesQualifiedName(processCommonJSModules_WEBPACK_REQUIRE+".e") {
				return false
			}
			allElementsAreDynamicImports = true
		}
		return allElementsAreDynamicImports
	}
	return false
}

// ProcessCommonJSModules_UmdPattern is information on a Universal Module
// Definition A UMD is an IF statement and a reference to which branch
// contains the commonjs export
type ProcessCommonJSModules_UmdPattern struct {
	ifRoot       rhino.Node
	activeBranch rhino.Node
}

// Java: ProcessCommonJSModules.UmdPattern(Node, Node)
func NewProcessCommonJSModules_UmdPattern(ifRoot rhino.Node, activeBranch rhino.Node) *ProcessCommonJSModules_UmdPattern {
	u := &ProcessCommonJSModules_UmdPattern{}
	u.ifRoot = ifRoot
	u.activeBranch = activeBranch
	return u
}

// ProcessCommonJSModules_ExportInfo is the nested class
// ProcessCommonJSModules.ExportInfo (identity equality).
type ProcessCommonJSModules_ExportInfo struct {
	node               rhino.Node
	scope              *Scope
	isInSupportedScope bool
}

// Java: ProcessCommonJSModules.ExportInfo(Node, Scope)
func NewProcessCommonJSModules_ExportInfo(node rhino.Node, scope *Scope) *ProcessCommonJSModules_ExportInfo {
	e := &ProcessCommonJSModules_ExportInfo{}
	e.node = node
	e.scope = scope

	disqualifyingParent :=
		NodeUtil_GetEnclosingNode(
			node, func(n rhino.Node) bool { return n.IsIf() || n.IsHook() || n.IsFunction() || n.IsArrowFunction() })
	e.isInSupportedScope = disqualifyingParent.IsNil()
	return e
}

// Java: ProcessCommonJSModules.getBaseQualifiedNameNode(Node)
func (p *ProcessCommonJSModules) GetBaseQualifiedNameNode(n rhino.Node) rhino.Node {
	refParent := n
	for refParent.HasParent() && refParent.GetParent().IsQualifiedName() {
		refParent = refParent.GetParent()
	}

	return refParent
}

// RemoveIIFEWrapper: UMD modules are often wrapped in an IIFE for cases
// where they are used as scripts instead of modules. Remove the wrapper.
// Returns whether an IIFE wrapper was found and removed.
//
// Java: ProcessCommonJSModules.removeIIFEWrapper(Node)
func (p *ProcessCommonJSModules) RemoveIIFEWrapper(root rhino.Node) bool {
	jutil.CheckState(root.IsScript())
	exprStmt, ok := processCommonJSModules_onlyExprStatement(ast.MustScript(root))
	if !ok {
		return false
	}

	// Function expression can be forced with !, just skip !
	// TODO(ChadKillingsworth):
	//   Expression could also be forced with: + - ~ void
	//   ! ~ void can be repeated any number of times
	expr := exprStmt.Expr()
	if not, ok := expr.AsUnaryOf(rhino.Token_NOT); ok {
		expr = not.Operand()
	}

	iife, ok := expr.AsCall()
	if !ok {
		return false
	}
	call := iife.Node()

	// Find the IIFE call and function nodes
	var fnc rhino.Node
	if function, ok := iife.Callee().AsFunction(); ok {
		fnc = function.Node()
	} else if callMethod, ok := iife.Callee().AsGetProp(); ok && callMethod.Property() == "call" {
		function, ok := callMethod.Object().AsFunction()
		if !ok {
			return false
		}
		fnc = function.Node()

		// We only support explicitly binding "this" to the parent "this" or "exports"
		thisArg, ok := iife.FirstArg()
		if !ok || !(thisArg.Node().IsThis() || thisArg.Node().MatchesQualifiedName(processCommonJSModules_EXPORTS)) {
			return false
		}
	} else {
		return false
	}

	if NodeUtil_DoesFunctionReferenceOwnArgumentsObject(fnc) {
		return false
	}

	ci := p.compiler.GetInput(root.GetInputId())
	modulePath := ci.GetPath()
	if modulePath == nil {
		return false
	}

	iifeLabel := ProcessCommonJSModules_GetModuleNameWithModulePath(modulePath) + "_iifeWrapper"

	mutator :=
		NewFunctionToBlockMutator(p.compiler, p.compiler.GetUniqueNameIdSupplier())
	block := mutator.MutateWithoutRenaming(iifeLabel, fnc, call, "", false, false)
	root.RemoveChildren()
	root.AddChildrenToFront(block.RemoveChildren())
	p.ReportNestedScopesDeleted(fnc)
	p.compiler.ReportChangeToEnclosingScope(root)

	return true
}

// processCommonJSModules_onlyExprStatement returns the expression statement
// that is the only statement of script besides empty ones before it (an
// IIFE wrapper); ok is false if there is none.
//
// Go-only: Java writes the same walk in removeIIFEWrapper and
// removeWebpackModuleShim.
func processCommonJSModules_onlyExprStatement(script ast.Script) (ast.ExprStmt, bool) {
	// Sometimes scripts start with a semicolon for easy concatenation.
	// Skip any empty statements from those
	item, ok := script.FirstItem()
	for ok && item.Node().IsEmpty() {
		item, ok = item.Next()
	}

	// An IIFE wrapper must be the only non-empty statement in the script,
	// and it must be an expression statement.
	if !ok || !item.IsLast() {
		return ast.ExprStmt{}, false
	}
	return ast.AsExprStmt(item.Node())
}

// RemoveWebpackModuleShim: for AMD wrappers, webpack adds a shim for the
// "module" variable. We need that to be a free var so we remove the shim.
//
// Java: ProcessCommonJSModules.removeWebpackModuleShim(Node)
func (p *ProcessCommonJSModules) RemoveWebpackModuleShim(root rhino.Node) {
	jutil.CheckState(root.IsScript())
	exprStmt, ok := processCommonJSModules_onlyExprStatement(ast.MustScript(root))
	if !ok {
		return
	}

	iife, ok := exprStmt.Expr().AsCall()
	if !ok {
		return
	}
	call := iife.Node()

	// Find the IIFE call and function nodes

	function, ok := iife.Callee().AsFunction()
	if !ok {
		return
	}

	fnc := function.Node()

	var param rhino.Node
	paramNumber := 0
	for typedParam := range function.Params().Params() {
		paramNumber++
		if name, ok := typedParam.AsName(); ok && name.Ident() == processCommonJSModules_MODULE {
			param = name.Node()
			break
		}
	}
	if param.IsNil() {
		return
	}

	// Java quirk kept: the argument is read by child index, which counts the
	// callee; for a free call the index is one lower, so it reads the argument
	// before the one of the `module` parameter.
	isFreeCall := call.GetBooleanProp(rhino.Node_FREE_CALL)
	argIndex := paramNumber + 1
	if isFreeCall {
		argIndex = paramNumber
	}
	arg := call.GetChildAtIndex(argIndex)
	if arg.IsNil() {
		return
	}

	// `require('...')(module)`, webpack's shim
	argCall, ok := ast.AsCall(arg)
	if !ok {
		return
	}
	argCallTarget := argCall.Callee().Node()
	if moduleArg, ok := argCall.FirstArg(); ok &&
		p.IsCommonJsImport(argCallTarget) &&
		moduleArg.Node().MatchesName(processCommonJSModules_MODULE) {
		importPath := p.GetCommonJsImportPath(argCallTarget)

		modulePath :=
			p.compiler.
				GetInput(root.GetInputId()).
				GetPath().
				ResolveJsModule(
					importPath, processCommonJSModules_sourceFileName(arg), arg.GetLineno(), arg.GetCharno())
		if modulePath == nil {
			// The module loader will issue an error
			return
		}

		if strings.Contains(modulePath.String(), "/buildin/module.js") {
			arg.Detach()
			param.Detach()
			p.compiler.ReportChangeToChangeScope(fnc)
			p.compiler.ReportChangeToEnclosingScope(fnc)
		}
	}
}

var processCommonJSModules_DEFINE_AMD = rhino.QualifiedName_Of("define.amd")
var processCommonJSModules_WINDOW_DEFINE = rhino.QualifiedName_Of("window.define")
var processCommonJSModules_WINDOW_DEFINE_AMD = rhino.QualifiedName_Of("window.define.amd")
var processCommonJSModules_REQUIRE_ENSURE = rhino.QualifiedName_Of("require.ensure")
var processCommonJSModules_GOOG_PROVIDE = rhino.QualifiedName_Of("goog.provide")
var processCommonJSModules_GOOG_MODULE = rhino.QualifiedName_Of("goog.module")
var processCommonJSModules_MODULE_EXPORTS = rhino.QualifiedName_Of("module.exports")

// ProcessCommonJSModules_FindImportsAndExports traverses the script. Find
// all references to CommonJS require (import) and module.exports or export
// statements. Rewrites any require calls to reference the rewritten module
// name.
type ProcessCommonJSModules_FindImportsAndExports struct {
	outer *ProcessCommonJSModules // the enclosing instance

	hasGoogProvideOrModule bool
	script                 rhino.Node

	umdPatterns   []*ProcessCommonJSModules_UmdPattern
	moduleExports []*ProcessCommonJSModules_ExportInfo
	exports       []*ProcessCommonJSModules_ExportInfo
	errors        []*JSError
}

// Java: ProcessCommonJSModules.FindImportsAndExports()
func NewProcessCommonJSModules_FindImportsAndExports(outer *ProcessCommonJSModules) *ProcessCommonJSModules_FindImportsAndExports {
	f := &ProcessCommonJSModules_FindImportsAndExports{outer: outer}
	f.hasGoogProvideOrModule = false
	f.script = rhino.NilNode
	f.umdPatterns = []*ProcessCommonJSModules_UmdPattern{}
	f.moduleExports = []*ProcessCommonJSModules_ExportInfo{}
	f.exports = []*ProcessCommonJSModules_ExportInfo{}
	f.errors = []*JSError{}
	return f
}

// Java: ProcessCommonJSModules.FindImportsAndExports.isCommonJsModule()
func (f *ProcessCommonJSModules_FindImportsAndExports) IsCommonJsModule() bool {
	return (len(f.exports) != 0 || len(f.moduleExports) != 0) && !f.hasGoogProvideOrModule
}

// Java: ProcessCommonJSModules.FindImportsAndExports.report(CheckLevel, JSError)
func (f *ProcessCommonJSModules_FindImportsAndExports) Report(ignoredLevel *CheckLevel, error_ *JSError) {
	f.errors = append(f.errors, error_)
}

// Java: ProcessCommonJSModules.FindImportsAndExports.getModuleExports()
func (f *ProcessCommonJSModules_FindImportsAndExports) GetModuleExports() []*ProcessCommonJSModules_ExportInfo {
	return slices.Clone(f.moduleExports)
}

// Java: ProcessCommonJSModules.FindImportsAndExports.getExports()
func (f *ProcessCommonJSModules_FindImportsAndExports) GetExports() []*ProcessCommonJSModules_ExportInfo {
	return slices.Clone(f.exports)
}

// Java: ProcessCommonJSModules.FindImportsAndExports.shouldTraverse(NodeTraversal, Node, Node)
func (f *ProcessCommonJSModules_FindImportsAndExports) ShouldTraverse(nodeTraversal *NodeTraversal, n rhino.Node, parent rhino.Node) bool {
	if n.IsScript() {
		jutil.CheckState(f.script.IsNil())
		f.script = n
	}
	return true
}

// Java: ProcessCommonJSModules.FindImportsAndExports.visit(NodeTraversal, Node, Node)
func (f *ProcessCommonJSModules_FindImportsAndExports) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	// Check for goog.provide or goog.module statements
	if NodeUtil_IsShallowStatementTree(parent) {
		if n.IsExprResult() {
			maybeGetProp := n.GetFirstFirstChild()
			if !maybeGetProp.IsNil() &&
				(processCommonJSModules_GOOG_PROVIDE.Matches(maybeGetProp) || processCommonJSModules_GOOG_MODULE.Matches(maybeGetProp)) {
				f.hasGoogProvideOrModule = true
			}
		}
	}

	// Find require.ensure calls
	if call, ok := ast.AsCall(n); ok && processCommonJSModules_REQUIRE_ENSURE.Matches(call.Callee().Node()) {
		f.VisitRequireEnsureCall(t, call)
	}

	if n.MatchesQualifiedName(processCommonJSModules_MODULE+"."+processCommonJSModules_EXPORTS) ||
		processCommonJSModules_isModuleExportsElement(n) {
		if f.outer.IsCommonJsExport(t, n) {
			f.moduleExports = append(f.moduleExports, NewProcessCommonJSModules_ExportInfo(n, t.GetScope()))

			// If the module.exports statement is nested in an if statement,
			// and the test of the if checks for "module" or "define,
			// assume the if statement is an UMD pattern with a common js export in the then branch
			umdTestAncestor := f.GetOutermostUmdTest(parent)
			if umdTestAncestor != nil && !f.IsInIfTest(n) {
				existingPattern := ProcessCommonJSModules_FindUmdPattern(f.umdPatterns, umdTestAncestor.enclosingIf)
				if existingPattern == nil {
					f.umdPatterns = append(f.umdPatterns,
						NewProcessCommonJSModules_UmdPattern(umdTestAncestor.enclosingIf, umdTestAncestor.activeBranch))
				}
			}
		}
	} else if processCommonJSModules_DEFINE_AMD.Matches(n) || processCommonJSModules_WINDOW_DEFINE_AMD.Matches(n) {
		// If a define.amd statement is in the test of an if statement,
		// and the statement has no "else" block, we simply want to remove
		// the entire if statement.
		umdTestAncestor := f.GetOutermostUmdTest(parent)
		if umdTestAncestor != nil && f.IsInIfTest(n) {
			existingPattern := ProcessCommonJSModules_FindUmdPattern(f.umdPatterns, umdTestAncestor.enclosingIf)
			if existingPattern == nil && processCommonJSModules_isIfWithoutElse(umdTestAncestor.enclosingIf) {
				f.umdPatterns = append(f.umdPatterns, NewProcessCommonJSModules_UmdPattern(umdTestAncestor.enclosingIf, rhino.NilNode))
			}
		}
	}

	if n.IsName() && processCommonJSModules_EXPORTS == n.GetString() {
		v := t.GetScope().GetVar(processCommonJSModules_EXPORTS)
		if v == nil || v.IsGlobal() {
			qNameRoot := f.outer.GetBaseQualifiedNameNode(n)
			if !qNameRoot.IsNil() &&
				qNameRoot.MatchesQualifiedName(processCommonJSModules_EXPORTS) &&
				NodeUtil_IsLValue(qNameRoot) {
				// Match the special assignment
				// exports = module.exports
				if n.GetGrandparent().IsExprResult() &&
					!n.GetNext().IsNil() &&
					((n.GetNext().IsGetProp() &&
						n.GetNext().MatchesQualifiedName(processCommonJSModules_MODULE+"."+processCommonJSModules_EXPORTS)) ||
						(n.GetNext().IsAssign() &&
							n.GetNext().
								GetFirstChild().
								MatchesQualifiedName(processCommonJSModules_MODULE+"."+processCommonJSModules_EXPORTS))) {
					f.exports = append(f.exports, NewProcessCommonJSModules_ExportInfo(n, t.GetScope()))

					// Ignore inlining created identity vars
					// var exports = exports
				} else if !f.hasGoogProvideOrModule &&
					(v == nil ||
						(v.GetNameNode().IsNil() && v.GetNameNode().GetFirstChild() != n)) {
					f.errors = append(f.errors, JSError_MakeWithNode(qNameRoot, ProcessCommonJSModules_SUSPICIOUS_EXPORTS_ASSIGNMENT))
				}
			} else {
				f.exports = append(f.exports, NewProcessCommonJSModules_ExportInfo(n, t.GetScope()))

				// If the exports statement is nested in the then branch of an if statement,
				// and the test of the if checks for "module" or "define,
				// assume the if statement is an UMD pattern with a common js export in the then branch
				umdTestAncestor := f.GetOutermostUmdTest(parent)
				if umdTestAncestor != nil && !f.IsInIfTest(n) {
					existingPattern := ProcessCommonJSModules_FindUmdPattern(f.umdPatterns, umdTestAncestor.enclosingIf)
					if existingPattern == nil {
						f.umdPatterns = append(f.umdPatterns,
							NewProcessCommonJSModules_UmdPattern(umdTestAncestor.enclosingIf, umdTestAncestor.activeBranch))
					}
				}
			}
		}
	} else if n.IsThis() && n.GetParent().IsGetProp() && t.InGlobalScope() {
		f.exports = append(f.exports, NewProcessCommonJSModules_ExportInfo(n, t.GetScope()))
	}

	if f.outer.IsCommonJsImport(n) {
		f.VisitRequireCall(t, n, parent)
	}
}

// VisitRequireCall visits require calls.
//
// Java: ProcessCommonJSModules.FindImportsAndExports.visitRequireCall(NodeTraversal, Node, Node)
func (f *ProcessCommonJSModules_FindImportsAndExports) VisitRequireCall(t *NodeTraversal, require rhino.Node, parent rhino.Node) {
	// When require("name") is used as a standalone statement (the result isn't used)
	// it indicates that a module is being loaded for the side effects it produces.
	// In this case the require statement should just be removed as the dependency
	// sorting will insert the file for us.
	if !NodeUtil_IsExpressionResultUsed(require) &&
		parent.IsExprResult() &&
		NodeUtil_IsStatementBlock(parent.GetParent()) {

		// Attempt to resolve the module so that load warnings are issued
		t.GetInput().
			GetPath().
			ResolveJsModule(
				f.outer.GetCommonJsImportPath(require),
				processCommonJSModules_sourceFileName(require),
				require.GetLineno(),
				require.GetCharno())
		grandparent := parent.GetParent()
		parent.Detach()
		f.outer.compiler.ReportChangeToEnclosingScope(grandparent)
	}
}

// VisitRequireEnsureCall visits require.ensure calls. Replace the call with
// an IIFE. Require.ensure must always be of the form:
//
// require.ensure(['module1', ...], function(require) {})
//
// Java: ProcessCommonJSModules.FindImportsAndExports.visitRequireEnsureCall(NodeTraversal, Node)
func (f *ProcessCommonJSModules_FindImportsAndExports) VisitRequireEnsureCall(t *NodeTraversal, ensure ast.Call) {
	call := ensure.Node()
	dependenciesArg, hasDependencies := ensure.FirstArg()
	var callbackArg ast.Arg
	hasCallback := false
	if hasDependencies {
		callbackArg, hasCallback = dependenciesArg.Next()
	}
	if !hasCallback || !callbackArg.IsLast() {
		// Java quirk kept: the count is the call's child count, which includes the callee.
		f.outer.compiler.Report(
			JSError_MakeWithNode(
				call,
				ProcessCommonJSModules_UNKNOWN_REQUIRE_ENSURE,
				"Expected the function to have 2 arguments but instead found {0}",
				""+strconv.Itoa(call.GetChildCount())))
		return
	}

	dependencies, ok := ast.AsArrayLit(dependenciesArg.Node())
	if !ok {
		f.outer.compiler.Report(
			JSError_MakeWithNode(
				dependenciesArg.Node(),
				ProcessCommonJSModules_UNKNOWN_REQUIRE_ENSURE,
				"The first argument must be an array literal of string literals."))
		return
	}

	for dep := range dependencies.Elements() {
		if !dep.Node().IsStringLit() {
			f.outer.compiler.Report(
				JSError_MakeWithNode(
					dep.Node(),
					ProcessCommonJSModules_UNKNOWN_REQUIRE_ENSURE,
					"The first argument must be an array literal of string literals."))
			return
		}
	}
	callback, ok := ast.AsFunction(callbackArg.Node())
	if ok {
		requireParam, hasOneParam := callback.Params().OnlyParam()
		name, isName := requireParam.AsName()
		ok = hasOneParam && isName && "require" == name.Ident()
	}
	if !ok {
		f.outer.compiler.Report(
			JSError_MakeWithNode(
				callbackArg.Node(),
				ProcessCommonJSModules_UNKNOWN_REQUIRE_ENSURE,
				"The second argument must be a function"+
					" whose first argument is named \"require\"."))
		return
	}

	callback.Detach()

	// Remove the "require" argument from the parameter list.
	callback.Params().Node().RemoveChildren()
	call.RemoveChildren()
	call.PutBooleanProp(rhino.Node_FREE_CALL, true)
	call.AddChildToFront(callback.Node())

	t.ReportCodeChange()
}

// Java: ProcessCommonJSModules.FindImportsAndExports.reportModuleErrors()
func (f *ProcessCommonJSModules_FindImportsAndExports) ReportModuleErrors() {
	for _, error_ := range f.errors {
		f.outer.compiler.Report(error_)
	}
}

// InitializeModule: if the export is directly assigned more than once, or
// the assignments are not global, declare the module name variable.
//
// If all of the assignments are simply property assignments, initialize the
// module name variable as a namespace.
//
// # Returns whether the default export can be declared constant
//
// Java: ProcessCommonJSModules.FindImportsAndExports.initializeModule()
func (f *ProcessCommonJSModules_FindImportsAndExports) InitializeModule() bool {
	compiler := f.outer.compiler
	ci := compiler.GetInput(f.script.GetInputId())
	modulePath := ci.GetPath()
	if modulePath == nil {
		return true
	}

	moduleName := *ProcessCommonJSModules_GetModuleName(ci)

	var exportsToRemove []*ProcessCommonJSModules_ExportInfo
	for _, export := range f.exports {
		if NodeUtil_GetEnclosingScript(export.node).IsNil() {
			exportsToRemove = append(exportsToRemove, export)
			continue
		}
		qNameBase := f.outer.GetBaseQualifiedNameNode(export.node)
		exportAssign, isAssignTarget := processCommonJSModules_assignToTarget(export.node)
		if export.node == qNameBase &&
			isAssignTarget &&
			export.node.GetGrandparent().IsExprResult() {
			value := exportAssign.Value()

			// Find any identity assignments and just remove them
			// exports = module.exports;
			if getProp, ok := value.AsGetProp(); ok &&
				getProp.Node().MatchesQualifiedName(processCommonJSModules_MODULE+"."+processCommonJSModules_EXPORTS) {
				for i, moduleExport := range f.moduleExports {
					if value.Is(moduleExport.node) {
						f.moduleExports = slices.Delete(f.moduleExports, i, i+1)
						break
					}
				}

				changeRoot := export.node.GetGrandparent().GetParent()
				export.node.GetGrandparent().Detach()
				exportsToRemove = append(exportsToRemove, export)
				compiler.ReportChangeToEnclosingScope(changeRoot)

				// Find compound identity assignments and remove the exports = portion
				// exports = module.exports = foo;
			} else if assign, ok := value.AsAssign(); ok &&
				assign.Target().Node().MatchesQualifiedName(processCommonJSModules_MODULE+"."+processCommonJSModules_EXPORTS) {
				exportAssign.Node().ReplaceWith(assign.Detach().Node())
				exportsToRemove = append(exportsToRemove, export)
				compiler.ReportChangeToEnclosingScope(assign.Node())
			}
			// Find babel transpiled interop assignment
			// module.exports = exports['default'];
		} else if getElem, prop, ok := processCommonJSModules_babelDefaultInterop(export.node); ok {
			parent := getElem.Node()
			grandparent := parent.GetParent()
			prop.Detach()
			parent.ReplaceWith(
				rhino.IR_Getprop(export.node.Detach(), prop.Value()).Srcref(parent))

			compiler.ReportChangeToEnclosingScope(grandparent)
		}
	}

	f.exports = processCommonJSModules_removeAll(f.exports, exportsToRemove)
	exportsToRemove = exportsToRemove[:0]
	exportsToReplace := jutil.NewLinkedHashMap[*ProcessCommonJSModules_ExportInfo, *ProcessCommonJSModules_ExportInfo]()
	for _, export := range f.moduleExports {
		if NodeUtil_GetEnclosingScript(export.node).IsNil() {
			exportsToRemove = append(exportsToRemove, export)
		} else if getElem, ok := ast.AsGetElem(export.node); ok {
			// module['exports'] becomes module.exports.
			object, prop := getElem.Object(), getElem.Index() // read both before detaching either
			prop.Detach()
			newExport :=
				NewProcessCommonJSModules_ExportInfo(
					rhino.IR_Getprop(object.Detach().Node(), prop.Node().GetString()), export.scope)
			export.node.ReplaceWith(newExport.node)
			compiler.ReportChangeToEnclosingScope(newExport.node)
			exportsToReplace.Put(export, newExport)
		}
	}
	f.moduleExports = processCommonJSModules_removeAll(f.moduleExports, exportsToRemove)
	for oldExport := range exportsToReplace.Keys() {
		oldIndex := slices.Index(f.moduleExports, oldExport)
		f.moduleExports = slices.Delete(f.moduleExports, oldIndex, oldIndex+1)
		f.moduleExports = slices.Insert(f.moduleExports, oldIndex, exportsToReplace.Get(oldExport))
	}

	// If we assign to the variable more than once or all the assignments
	// are properties, initialize the variable as well.
	directAssignments := 0
	for _, export := range f.moduleExports {
		if NodeUtil_GetEnclosingScript(export.node).IsNil() {
			continue
		}

		base := f.outer.GetBaseQualifiedNameNode(export.node)
		if base == export.node && export.node.GetParent().IsAssign() {
			if rValue, ok := ast.AsObjectLit(NodeUtil_GetRValueOfLValue(export.node)); !ok {
				directAssignments++
			} else if !processCommonJSModules_isNamespaceLiteral(rValue) {
				directAssignments++
			}
		}
	}

	initModule := rhino.IR_VarWithValue(rhino.IR_Name(moduleName), rhino.IR_Objectlit())
	initModule.GetFirstChild().PutBooleanProp(rhino.Node_MODULE_EXPORT, true)
	initModule.GetFirstChild().MakeNonIndexable()
	builder := rhino.JSDocInfo_Builder_().ParseDocumentation()
	builder.RecordConstancy()
	initModule.SetJSDocInfo(builder.Build())
	if directAssignments == 0 || (len(f.exports) != 0 && len(f.moduleExports) != 0) {
		defaultProp := rhino.IR_StringKey(processCommonJSModules_EXPORT_PROPERTY_NAME)
		defaultProp.PutBooleanProp(rhino.Node_MODULE_EXPORT, true)
		defaultProp.AddChildToFront(rhino.IR_Objectlit())
		initModule.GetFirstFirstChild().AddChildToFront(defaultProp)
		if len(f.exports) == 0 || len(f.moduleExports) == 0 {
			builder = rhino.JSDocInfo_Builder_().ParseDocumentation()
			builder.RecordConstancy()
			defaultProp.SetJSDocInfo(builder.Build())
		}
	}
	f.script.AddChildToFront(initModule.SrcrefTree(f.script))
	compiler.ReportChangeToEnclosingScope(f.script)

	return directAssignments < 2 && (len(f.exports) == 0 || len(f.moduleExports) == 0)
}

// processCommonJSModules_assignToTarget returns the assignment whose target
// is n; ok is false if n is not the target of an assignment.
//
// Go-only: Java tests n.getParent().isAssign() and that n has no previous
// sibling, then reads the value as n.getNext().
func processCommonJSModules_assignToTarget(n rhino.Node) (ast.Assign, bool) {
	assign, ok := ast.AsAssign(n.GetParent())
	return assign, ok && assign.Target().Is(n)
}

// processCommonJSModules_babelDefaultInterop recognizes the `exports` of
// Babel's interop assignment `module.exports = exports['default'];` and
// returns its `exports['default']` and the 'default' string.
//
// Go-only: Java tests the siblings and parents of exports inline.
func processCommonJSModules_babelDefaultInterop(exports rhino.Node) (ast.GetElem, ast.StringLit, bool) {
	getElem, ok := ast.AsGetElem(exports.GetParent())
	if !ok || !getElem.Object().Is(exports) {
		return ast.GetElem{}, ast.StringLit{}, false
	}
	prop, ok := getElem.Index().AsStringLit()
	if !ok || prop.Value() != processCommonJSModules_EXPORT_PROPERTY_NAME {
		return ast.GetElem{}, ast.StringLit{}, false
	}
	assign, ok := ast.AsAssign(getElem.Node().GetParent())
	if !ok || !assign.Value().Is(getElem.Node()) ||
		!assign.Target().Node().MatchesQualifiedName(processCommonJSModules_MODULE+"."+processCommonJSModules_EXPORTS) {
		return ast.GetElem{}, ast.StringLit{}, false
	}
	return getElem, prop, true
}

// processCommonJSModules_isNamespaceLiteral reports whether every member of
// obj is an unquoted `key: value` or a method: an object literal whose
// members can become property assignments.
//
// Go-only: Java writes the member test twice, in initializeModule and in
// RewriteModule.expandObjectLitAssignment.
func processCommonJSModules_isNamespaceLiteral(obj ast.ObjectLit) bool {
	for member := range obj.Members() {
		if !processCommonJSModules_isNamespaceMember(member) {
			return false
		}
	}
	return true
}

// processCommonJSModules_isNamespaceMember reports whether member is an
// unquoted `key: value` or a method.
func processCommonJSModules_isNamespaceMember(member ast.ObjectMember) bool {
	if key, ok := member.AsStringKey(); ok {
		return !key.IsQuoted()
	}
	_, ok := member.AsMemberFunctionDef()
	return ok
}

// processCommonJSModules_removeAll is List.removeAll(Collection) on a list
// of identity-compared elements: it keeps the elements of list that are not
// in c, in order.
//
// Go-only: language-adapter — ArrayList.removeAll on a Go slice.
func processCommonJSModules_removeAll(
	list []*ProcessCommonJSModules_ExportInfo, c []*ProcessCommonJSModules_ExportInfo) []*ProcessCommonJSModules_ExportInfo {
	return slices.DeleteFunc(list, func(e *ProcessCommonJSModules_ExportInfo) bool { return slices.Contains(c, e) })
}

// ProcessCommonJSModules_FindImportsAndExports_UmdTestInfo is the private
// static nested class FindImportsAndExports.UmdTestInfo.
type ProcessCommonJSModules_FindImportsAndExports_UmdTestInfo struct {
	enclosingIf  rhino.Node
	activeBranch rhino.Node
}

// Java: ProcessCommonJSModules.FindImportsAndExports.UmdTestInfo(Node, Node)
func NewProcessCommonJSModules_FindImportsAndExports_UmdTestInfo(enclosingIf rhino.Node, activeBranch rhino.Node) *ProcessCommonJSModules_FindImportsAndExports_UmdTestInfo {
	u := &ProcessCommonJSModules_FindImportsAndExports_UmdTestInfo{}
	u.enclosingIf = enclosingIf
	u.activeBranch = activeBranch
	return u
}

// GetOutermostUmdTest finds the outermost if node ancestor for a node
// without leaving the function scope. To match, the test class of the "if"
// statement must reference "module" or "define" names.
//
// Java: ProcessCommonJSModules.FindImportsAndExports.getOutermostUmdTest(Node)
func (f *ProcessCommonJSModules_FindImportsAndExports) GetOutermostUmdTest(n rhino.Node) *ProcessCommonJSModules_FindImportsAndExports_UmdTestInfo {
	if n.IsNil() || NodeUtil_IsTopLevel(n) || n.IsFunction() {
		return nil
	}
	parent := n.GetParent()
	if parent.IsNil() {
		return nil
	}

	// When walking up ternary operations (hook), don't check if parent is the condition,
	// because one ternary operation can be then/else branch of another.
	if cond, ok := processCommonJSModules_ifOrHookCondition(parent); ok {
		umdTestInfo := f.GetOutermostUmdTest(parent)
		if umdTestInfo != nil {
			// If this is the then block of an else-if statement, set the active branch to the
			// then block rather than the "if" itself.
			if ifs, ok := ast.AsIf(parent); ok &&
				ifs.Then().Is(n) &&
				parent.GetParent().IsBlock() &&
				parent.GetNext().IsNil() &&
				umdTestInfo.activeBranch == parent.GetParent() {
				return NewProcessCommonJSModules_FindImportsAndExports_UmdTestInfo(umdTestInfo.enclosingIf, n)
			}

			return umdTestInfo
		}

		umdTests := []rhino.Node{}
		NodeUtil_VisitPreOrder(
			cond.Node(),
			NodeUtil_VisitorFunc(func(node rhino.Node) {
				switch node.GetToken() {
				case rhino.Token_NAME:
					if node.GetString() == processCommonJSModules_MODULE || node.GetString() == "define" {
						break
					}
					return
				case rhino.Token_GETPROP:
					if processCommonJSModules_WINDOW_DEFINE.Matches(node) {
						break
					}
					return
				case rhino.Token_STRINGLIT:
					if node.GetParent().IsIn() && node.GetString() == "amd" {
						break
					}
					return
				default:
					return
				}
				umdTests = append(umdTests, node)
			}))

		// Webpack replaces tests of `typeof module !== 'undefined'` with `true`
		if _, isTrue := cond.AsTrue(); len(umdTests) == 0 && isTrue {
			umdTests = append(umdTests, cond.Node())
		}

		if len(umdTests) != 0 {
			return NewProcessCommonJSModules_FindImportsAndExports_UmdTestInfo(parent, n)
		}

		return nil
	}

	return f.GetOutermostUmdTest(parent)
}

// IsInIfTest returns whether the node is within the test portion of an if
// statement
//
// Java: ProcessCommonJSModules.FindImportsAndExports.isInIfTest(Node)
func (f *ProcessCommonJSModules_FindImportsAndExports) IsInIfTest(n rhino.Node) bool {
	if n.IsNil() || NodeUtil_IsTopLevel(n) || n.IsFunction() {
		return false
	}
	parent := n.GetParent()
	if parent.IsNil() {
		return false
	}

	if cond, ok := processCommonJSModules_ifOrHookCondition(parent); ok && cond.Is(n) {
		return true
	}

	return f.IsInIfTest(parent)
}

// processCommonJSModules_ifOrHookCondition returns the condition of n if it
// is an `if` or a `?:`; ok is false otherwise.
//
// Go-only: Java tests `n.isIf() || n.isHook()` and reads the first child.
func processCommonJSModules_ifOrHookCondition(n rhino.Node) (ast.Expr, bool) {
	if ifs, ok := ast.AsIf(n); ok {
		return ifs.Cond(), true
	}
	if hook, ok := ast.AsHook(n); ok {
		return hook.Cond(), true
	}
	return ast.Expr{}, false
}

// processCommonJSModules_isIfWithoutElse reports whether n is an `if`
// without an `else` (Java: n.getChildAtIndex(2) == null, false for a `?:`,
// which always has a third child).
func processCommonJSModules_isIfWithoutElse(n rhino.Node) bool {
	ifs, ok := ast.AsIf(n)
	if !ok {
		return false
	}
	_, hasElse := ifs.Else()
	return !hasElse
}

// ReplaceUmdPatterns removes a Universal Module Definition and leave just
// the commonjs export statement
//
// Java: ProcessCommonJSModules.FindImportsAndExports.replaceUmdPatterns()
func (f *ProcessCommonJSModules_FindImportsAndExports) ReplaceUmdPatterns() bool {
	compiler := f.outer.compiler
	needsRetraverse := false
	var changeScope rhino.Node
	for _, umdPattern := range f.umdPatterns {
		if NodeUtil_GetEnclosingScript(umdPattern.ifRoot).IsNil() {
			f.outer.ReportNestedScopesDeleted(umdPattern.ifRoot)
			continue
		}

		parent := umdPattern.ifRoot.GetParent()
		newNode := umdPattern.activeBranch

		if newNode.IsNil() {
			umdPattern.ifRoot.Detach()
			f.outer.ReportNestedScopesDeleted(umdPattern.ifRoot)
			compiler.ReportChangeToEnclosingScope(parent)
			needsRetraverse = true
			continue
		}

		// Remove redundant block node. Not strictly necessary, but makes tests more legible.
		if block, ok := ast.AsBlock(umdPattern.activeBranch); ok {
			if only, ok := block.OnlyStatement(); ok {
				newNode = only.Detach().Node()
			} else {
				newNode.Detach()
			}
		} else {
			newNode.Detach()
		}
		needsRetraverse = true
		umdPattern.ifRoot.ReplaceWith(newNode)
		f.outer.ReportNestedScopesDeleted(umdPattern.ifRoot)
		changeScope = ChangeTracker_GetEnclosingChangeScopeRoot(newNode)
		if !changeScope.IsNil() {
			compiler.ReportChangeToEnclosingScope(newNode)
		}

		block := parent
		if block.IsExprResult() {
			block = block.GetParent()
		}

		// Detect UMD Factory Patterns and inline the functions
		if block.IsBlock() && block.GetParent().IsFunction() &&
			block.GetGrandparent().IsCall() &&
			parent.HasOneChild() {
			enclosingFnCall := block.GetGrandparent()
			fn := block.GetParent()

			enclosingScript := NodeUtil_GetEnclosingScript(enclosingFnCall)
			if enclosingScript.IsNil() {
				continue
			}
			ci := compiler.GetInput(
				NodeUtil_GetEnclosingScript(enclosingFnCall).GetInputId())
			modulePath := ci.GetPath()
			if modulePath == nil {
				continue
			}
			needsRetraverse = true
			factoryLabel :=
				modulePath.ToModuleName() +
					"_factory" +
					compiler.GetUniqueIdSupplier().GetUniqueId(ci)

			mutator :=
				NewFunctionToBlockMutator(compiler, compiler.GetUniqueNameIdSupplier())
			newStatements :=
				mutator.MutateWithoutRenaming(factoryLabel, fn, enclosingFnCall, "", false, false)

			// Check to see if the returned block is of the form:
			// {
			//   var jscomp$inline = function() {};
			//   jscomp$inline();
			// }
			//
			// or
			//
			// {
			//   var jscomp$inline = function() {};
			//   module.exports = jscomp$inline();
			// }
			//
			// If so, inline again
			if inlinedFn, exprStmt, ok := processCommonJSModules_inlinedFactory(newStatements); ok {
				expr := exprStmt.Expr()
				var call rhino.Node
				var assignedName string
				var exprAssign ast.Assign
				if assign, ok := expr.AsAssign(); ok {
					if valueCall, ok := assign.Value().AsCall(); ok {
						exprAssign = assign
						call = valueCall.Node()
						assignedName =
							modulePath.ToModuleName() +
								"_iife" +
								compiler.GetUniqueIdSupplier().GetUniqueId(ci)
					}
				} else if exprCall, ok := expr.AsCall(); ok {
					call = exprCall.Node()
				}

				if !call.IsNil() {
					newStatements =
						mutator.MutateWithoutRenaming(
							factoryLabel, inlinedFn.Node(), call, assignedName, false, false)
					if assignedName != "" {
						newName :=
							rhino.IR_Var(
								NodeUtil_NewNameWithOriginalName(
									compiler,
									assignedName,
									fn,
									*exprAssign.Target().Node().GetQualifiedName())).
								SrcrefTree(fn)
						// If the inlined function starts with `assignedName = value;`, declare it there.
						if first, ok := processCommonJSModules_firstAssignTo(newStatements, assignedName); ok {
							newName.
								GetFirstChild().
								AddChildToFront(first.Value().Detach().Node())
							newStatements.GetFirstChild().ReplaceWith(newName)
						} else {
							newStatements.AddChildToFront(newName)
						}
						exprAssign.Value().Node().ReplaceWith(newName.GetFirstChild().CloneNode())
						newStatements.AddChildToBack(exprStmt.Detach().Node())
					}
				}
			}

			callRoot := enclosingFnCall.GetParent()
			if callRoot.IsNot() {
				callRoot = callRoot.GetParent()
			}
			if callRoot.IsExprResult() {
				callRootParent := callRoot.GetParent()
				callRootParent.AddChildrenAfter(newStatements.RemoveChildren(), callRoot)
				callRoot.Detach()
				f.outer.ReportNestedScopesChanged(callRootParent)
				compiler.ReportChangeToEnclosingScope(callRootParent)
				f.outer.ReportNestedScopesDeleted(enclosingFnCall)
			} else {
				umdPattern.ifRoot.ReplaceWith(newNode)
				compiler.ReportChangeToEnclosingScope(newNode)
				f.outer.ReportNestedScopesDeleted(umdPattern.ifRoot)
			}
		}
	}
	return needsRetraverse
}

// processCommonJSModules_inlinedFactory recognizes the block a UMD factory
// was inlined into when it is of the form
//
//	{
//	  var jscomp$inline = function() {};
//	  jscomp$inline();  // or: module.exports = jscomp$inline();
//	}
//
// and returns the function and the second statement.
//
// Go-only: Java tests the shape with positional reads inline.
func processCommonJSModules_inlinedFactory(newStatements rhino.Node) (ast.Function, ast.ExprStmt, bool) {
	block, ok := ast.AsBlock(newStatements)
	if !ok {
		return ast.Function{}, ast.ExprStmt{}, false
	}
	first, ok := block.FirstStatement()
	if !ok {
		return ast.Function{}, ast.ExprStmt{}, false
	}
	second, ok := first.Next()
	if !ok || !second.IsLast() {
		return ast.Function{}, ast.ExprStmt{}, false
	}
	declaration, ok := first.AsNameDeclarationOf(rhino.Token_VAR)
	if !ok {
		return ast.Function{}, ast.ExprStmt{}, false
	}
	// The first declared name, with a function as its only child.
	firstDeclarator, _ := declaration.FirstDeclarator()
	name, ok := firstDeclarator.AsName()
	if !ok {
		return ast.Function{}, ast.ExprStmt{}, false
	}
	init, ok := name.Init()
	if !ok {
		return ast.Function{}, ast.ExprStmt{}, false
	}
	inlinedFn, ok := init.AsFunction()
	if !ok {
		return ast.Function{}, ast.ExprStmt{}, false
	}
	exprStmt, ok := second.AsExprStmt()
	return inlinedFn, exprStmt, ok
}

// processCommonJSModules_firstAssignTo returns the assignment of the first
// statement of block if that statement is `name = value;`.
//
// Go-only: Java tests the shape with positional reads inline.
func processCommonJSModules_firstAssignTo(block rhino.Node, name string) (ast.Assign, bool) {
	exprStmt, ok := ast.AsExprStmt(block.GetFirstChild())
	if !ok {
		return ast.Assign{}, false
	}
	assign, ok := exprStmt.Expr().AsAssign()
	if !ok {
		return ast.Assign{}, false
	}
	target, ok := assign.Target().AsName()
	return assign, ok && target.Ident() == name
}

// Java: ProcessCommonJSModules.reportNestedScopesDeleted(Node)
func (p *ProcessCommonJSModules) ReportNestedScopesDeleted(n rhino.Node) {
	NodeUtil_VisitPreOrder(
		n,
		NodeUtil_VisitorFunc(func(n rhino.Node) {
			if n.IsFunction() {
				p.compiler.ReportFunctionDeleted(n)
			}
		}))
}

// Java: ProcessCommonJSModules.reportNestedScopesChanged(Node)
func (p *ProcessCommonJSModules) ReportNestedScopesChanged(n rhino.Node) {
	NodeUtil_VisitPreOrder(
		n,
		NodeUtil_VisitorFunc(func(n rhino.Node) {
			if n.IsFunction() {
				p.compiler.ReportChangeToChangeScope(n)
			}
		}))
}

// Java: ProcessCommonJSModules.findUmdPattern(List<UmdPattern>, Node)
func ProcessCommonJSModules_FindUmdPattern(umdPatterns []*ProcessCommonJSModules_UmdPattern, n rhino.Node) *ProcessCommonJSModules_UmdPattern {
	for _, umd := range umdPatterns {
		if umd.ifRoot == n {
			return umd
		}
	}
	return nil
}

