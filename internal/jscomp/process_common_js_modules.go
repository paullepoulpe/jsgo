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

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

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
	if requireCall.IsCall() && requireCall.HasTwoChildren() {
		if resolutionMode == ModuleLoader_ResolutionMode_WEBPACK &&
			(requireCall.GetFirstChild().MatchesQualifiedName(processCommonJSModules_WEBPACK_REQUIRE) ||
				requireCall.GetFirstChild().MatchesQualifiedName(processCommonJSModules_WEBPACK_REQUIRE_NAMESPACE)) &&
			(requireCall.GetSecondChild().IsNumber() ||
				requireCall.GetSecondChild().IsStringLit()) {
			return true
		} else if requireCall.GetFirstChild().MatchesQualifiedName(processCommonJSModules_REQUIRE) &&
			requireCall.GetSecondChild().IsStringLit() {
			return true
		}
	} else if requireCall.IsCall() &&
		requireCall.HasXChildren(3) &&
		resolutionMode == ModuleLoader_ResolutionMode_WEBPACK &&
		requireCall.GetFirstChild().MatchesQualifiedName(processCommonJSModules_WEBPACK_REQUIRE+".bind") &&
		requireCall.GetSecondChild().IsNull() &&
		(requireCall.GetLastChild().IsNumber() || requireCall.GetLastChild().IsStringLit()) {
		return true
	}
	return false
}

// Java: ProcessCommonJSModules.getCommonJsImportPath(Node)
func (p *ProcessCommonJSModules) GetCommonJsImportPath(requireCall rhino.Node) string {
	return ProcessCommonJSModules_GetCommonJsImportPath(requireCall, p.compiler.GetOptions().GetModuleResolutionMode())
}

// Java: ProcessCommonJSModules.getCommonJsImportPath(Node, ModuleLoader.ResolutionMode)
func ProcessCommonJSModules_GetCommonJsImportPath(
	requireCall rhino.Node, resolutionMode ModuleLoader_ResolutionMode) string {
	if resolutionMode == ModuleLoader_ResolutionMode_WEBPACK {
		var pathArgument rhino.Node
		if requireCall.GetChildCount() >= 3 {
			pathArgument = requireCall.GetChildAtIndex(2)
		} else {
			pathArgument = requireCall.GetSecondChild()
		}
		if pathArgument.IsNumber() {
			return strconv.Itoa(int(jutil.D2I(pathArgument.GetDouble())))
		} else {
			return pathArgument.GetString()
		}
	}

	return requireCall.GetSecondChild().GetString()
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
		(export.IsGetElem() &&
			export.GetFirstChild().MatchesQualifiedName(processCommonJSModules_MODULE) &&
			export.GetSecondChild().IsStringLit() &&
			export.GetSecondChild().GetString() == processCommonJSModules_EXPORTS) {
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

	callParent := fnc.GetParent()
	if !callParent.IsCall() {
		return false
	}

	if !(callParent.HasChildren() &&
		callParent.GetFirstChild().IsGetProp() &&
		callParent.GetFirstFirstChild().IsCall()) {
		return false
	}
	callParentTarget := callParent.GetFirstFirstChild().GetFirstChild()

	if callParentTarget.MatchesQualifiedName(processCommonJSModules_WEBPACK_REQUIRE+".e") &&
		callParent.GetFirstChild().GetString() == "then" {
		return true
	} else if processCommonJSModules_PROMISE_ALL.Matches(callParentTarget) &&
		!callParentTarget.GetNext().IsNil() &&
		callParentTarget.GetNext().IsArrayLit() {
		allElementsAreDynamicImports := false
		for arrayItem := callParentTarget.GetNext().GetFirstChild(); !arrayItem.IsNil(); arrayItem = arrayItem.GetNext() {
			if !(arrayItem.IsCall() &&
				arrayItem.HasTwoChildren() &&
				arrayItem.GetFirstChild().MatchesQualifiedName(processCommonJSModules_WEBPACK_REQUIRE+".e")) {
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
	n := root.GetFirstChild()

	// Sometimes scripts start with a semicolon for easy concatenation.
	// Skip any empty statements from those
	for !n.IsNil() && n.IsEmpty() {
		n = n.GetNext()
	}

	// An IIFE wrapper must be the only non-empty statement in the script,
	// and it must be an expression statement.
	if n.IsNil() || !n.IsExprResult() || !n.GetNext().IsNil() {
		return false
	}

	// Function expression can be forced with !, just skip !
	// TODO(ChadKillingsworth):
	//   Expression could also be forced with: + - ~ void
	//   ! ~ void can be repeated any number of times
	if !n.IsNil() && n.HasChildren() && n.GetFirstChild().IsNot() {
		n = n.GetFirstChild()
	}

	call := n.GetFirstChild()
	if call.IsNil() || !call.IsCall() {
		return false
	}

	// Find the IIFE call and function nodes
	var fnc rhino.Node
	if call.GetFirstChild().IsFunction() {
		fnc = n.GetFirstFirstChild()
	} else if call.GetFirstChild().IsGetProp() &&
		call.GetFirstFirstChild().IsFunction() &&
		call.GetFirstChild().GetString() == "call" {
		fnc = call.GetFirstFirstChild()

		// We only support explicitly binding "this" to the parent "this" or "exports"
		if !(!call.GetSecondChild().IsNil() &&
			(call.GetSecondChild().IsThis() ||
				call.GetSecondChild().MatchesQualifiedName(processCommonJSModules_EXPORTS))) {
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
	block := mutator.MutateWithoutRenaming(iifeLabel, fnc, call, nil, false, false)
	root.RemoveChildren()
	root.AddChildrenToFront(block.RemoveChildren())
	p.ReportNestedScopesDeleted(fnc)
	p.compiler.ReportChangeToEnclosingScope(root)

	return true
}

// RemoveWebpackModuleShim: for AMD wrappers, webpack adds a shim for the
// "module" variable. We need that to be a free var so we remove the shim.
//
// Java: ProcessCommonJSModules.removeWebpackModuleShim(Node)
func (p *ProcessCommonJSModules) RemoveWebpackModuleShim(root rhino.Node) {
	jutil.CheckState(root.IsScript())
	n := root.GetFirstChild()

	// Sometimes scripts start with a semicolon for easy concatenation.
	// Skip any empty statements from those
	for !n.IsNil() && n.IsEmpty() {
		n = n.GetNext()
	}

	// An IIFE wrapper must be the only non-empty statement in the script,
	// and it must be an expression statement.
	if n.IsNil() || !n.IsExprResult() || !n.GetNext().IsNil() {
		return
	}

	call := n.GetFirstChild()
	if call.IsNil() || !call.IsCall() {
		return
	}

	// Find the IIFE call and function nodes

	callTarget := call.GetFirstChild()
	if !callTarget.IsFunction() {
		return
	}

	fnc := callTarget

	params := NodeUtil_GetFunctionParameters(fnc)
	var moduleParam rhino.Node
	param := params.GetFirstChild()
	paramNumber := 0
	for !param.IsNil() {
		paramNumber++
		if param.IsName() && param.GetString() == processCommonJSModules_MODULE {
			moduleParam = param
			break
		}
		param = param.GetNext()
	}
	if moduleParam.IsNil() {
		return
	}

	isFreeCall := call.GetBooleanProp(rhino.Node_FREE_CALL)
	argIndex := paramNumber + 1
	if isFreeCall {
		argIndex = paramNumber
	}
	arg := call.GetChildAtIndex(argIndex)
	if arg.IsNil() {
		return
	}

	argCallTarget := arg.GetFirstChild()
	if arg.IsCall() &&
		argCallTarget.IsCall() &&
		p.IsCommonJsImport(argCallTarget) &&
		argCallTarget.GetNext().MatchesName(processCommonJSModules_MODULE) {
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
	if n.IsCall() && processCommonJSModules_REQUIRE_ENSURE.Matches(n.GetFirstChild()) {
		f.VisitRequireEnsureCall(t, n)
	}

	if n.MatchesQualifiedName(processCommonJSModules_MODULE+"."+processCommonJSModules_EXPORTS) ||
		(n.IsGetElem() &&
			n.GetFirstChild().MatchesQualifiedName(processCommonJSModules_MODULE) &&
			n.GetSecondChild().IsStringLit() &&
			n.GetSecondChild().GetString() == processCommonJSModules_EXPORTS) {
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
			if existingPattern == nil && umdTestAncestor.enclosingIf.GetChildAtIndex(2).IsNil() {
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
func (f *ProcessCommonJSModules_FindImportsAndExports) VisitRequireEnsureCall(t *NodeTraversal, call rhino.Node) {
	if !call.HasXChildren(3) {
		f.outer.compiler.Report(
			JSError_MakeWithNode(
				call,
				ProcessCommonJSModules_UNKNOWN_REQUIRE_ENSURE,
				"Expected the function to have 2 arguments but instead found {0}",
				""+strconv.Itoa(call.GetChildCount())))
		return
	}

	dependencies := call.GetSecondChild()
	if !dependencies.IsArrayLit() {
		f.outer.compiler.Report(
			JSError_MakeWithNode(
				dependencies,
				ProcessCommonJSModules_UNKNOWN_REQUIRE_ENSURE,
				"The first argument must be an array literal of string literals."))
		return
	}

	for dep := dependencies.GetFirstChild(); !dep.IsNil(); dep = dep.GetNext() {
		if !dep.IsStringLit() {
			f.outer.compiler.Report(
				JSError_MakeWithNode(
					dep,
					ProcessCommonJSModules_UNKNOWN_REQUIRE_ENSURE,
					"The first argument must be an array literal of string literals."))
			return
		}
	}
	callback := dependencies.GetNext()
	if !(callback.IsFunction() &&
		callback.GetSecondChild().HasOneChild() &&
		callback.GetSecondChild().GetFirstChild().IsName() &&
		"require" == callback.GetSecondChild().GetFirstChild().GetString()) {
		f.outer.compiler.Report(
			JSError_MakeWithNode(
				callback,
				ProcessCommonJSModules_UNKNOWN_REQUIRE_ENSURE,
				"The second argument must be a function"+
					" whose first argument is named \"require\"."))
		return
	}

	callback.Detach()

	// Remove the "require" argument from the parameter list.
	callback.GetSecondChild().RemoveChildren()
	call.RemoveChildren()
	call.PutBooleanProp(rhino.Node_FREE_CALL, true)
	call.AddChildToFront(callback)

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
		if export.node == qNameBase &&
			export.node.GetParent().IsAssign() &&
			export.node.GetGrandparent().IsExprResult() &&
			export.node.GetPrevious().IsNil() &&
			!export.node.GetNext().IsNil() {

			// Find any identity assignments and just remove them
			// exports = module.exports;
			if export.node.GetNext().IsGetProp() &&
				export.node.GetNext().MatchesQualifiedName(processCommonJSModules_MODULE+"."+processCommonJSModules_EXPORTS) {
				for i, moduleExport := range f.moduleExports {
					if moduleExport.node == export.node.GetNext() {
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
			} else if export.node.GetNext().IsAssign() &&
				export.
					node.
					GetNext().
					GetFirstChild().
					MatchesQualifiedName(processCommonJSModules_MODULE+"."+processCommonJSModules_EXPORTS) {
				assign := export.node.GetNext()
				export.node.GetParent().ReplaceWith(assign.Detach())
				exportsToRemove = append(exportsToRemove, export)
				compiler.ReportChangeToEnclosingScope(assign)
			}
			// Find babel transpiled interop assignment
			// module.exports = exports['default'];
		} else if export.node.GetParent().IsGetElem() &&
			export.node.GetNext().IsStringLit() &&
			export.node.GetNext().GetString() == processCommonJSModules_EXPORT_PROPERTY_NAME &&
			export.node.GetGrandparent().IsAssign() &&
			!export.node.GetParent().GetPrevious().IsNil() &&
			export.node.GetParent().GetPrevious().MatchesQualifiedName(processCommonJSModules_MODULE+"."+processCommonJSModules_EXPORTS) {
			parent := export.node.GetParent()
			grandparent := parent.GetParent()
			prop := export.node.GetNext()
			parent.ReplaceWith(
				rhino.IR_Getprop(export.node.Detach(), prop.Detach().GetString()).Srcref(parent))

			compiler.ReportChangeToEnclosingScope(grandparent)
		}
	}

	f.exports = processCommonJSModules_removeAll(f.exports, exportsToRemove)
	exportsToRemove = exportsToRemove[:0]
	exportsToReplace := jutil.NewLinkedHashMap[*ProcessCommonJSModules_ExportInfo, *ProcessCommonJSModules_ExportInfo]()
	for _, export := range f.moduleExports {
		if NodeUtil_GetEnclosingScript(export.node).IsNil() {
			exportsToRemove = append(exportsToRemove, export)
		} else if export.node.IsGetElem() {
			prop := export.node.GetSecondChild().Detach()
			newExport :=
				NewProcessCommonJSModules_ExportInfo(
					rhino.IR_Getprop(export.node.RemoveFirstChild(), prop.GetString()), export.scope)
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
			rValue := NodeUtil_GetRValueOfLValue(export.node)
			if rValue.IsNil() || !rValue.IsObjectLit() {
				directAssignments++
			} else if rValue.IsObjectLit() {
				for key := rValue.GetFirstChild(); !key.IsNil(); key = key.GetNext() {
					if (!key.IsStringKey() || key.IsQuotedStringKey()) && !key.IsMemberFunctionDef() {
						directAssignments++
						break
					}
				}
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
	if parent.IsIf() || parent.IsHook() {
		umdTestInfo := f.GetOutermostUmdTest(parent)
		if umdTestInfo != nil {
			// If this is the then block of an else-if statement, set the active branch to the
			// then block rather than the "if" itself.
			if parent.IsIf() &&
				parent.GetSecondChild() == n &&
				parent.HasParent() &&
				parent.GetParent().IsBlock() &&
				parent.GetNext().IsNil() &&
				umdTestInfo.activeBranch == parent.GetParent() {
				return NewProcessCommonJSModules_FindImportsAndExports_UmdTestInfo(umdTestInfo.enclosingIf, n)
			}

			return umdTestInfo
		}

		umdTests := []rhino.Node{}
		NodeUtil_VisitPreOrder(
			parent.GetFirstChild(),
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
		if len(umdTests) == 0 && parent.GetFirstChild().IsTrue() {
			umdTests = append(umdTests, parent.GetFirstChild())
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

	if (parent.IsIf() || parent.IsHook()) && parent.GetFirstChild() == n {
		return true
	}

	return f.IsInIfTest(parent)
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
		if umdPattern.activeBranch.IsBlock() && umdPattern.activeBranch.HasOneChild() {
			newNode = umdPattern.activeBranch.RemoveFirstChild()
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
				mutator.MutateWithoutRenaming(factoryLabel, fn, enclosingFnCall, nil, false, false)

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
			if newStatements.IsBlock() &&
				newStatements.HasTwoChildren() &&
				newStatements.GetFirstChild().IsVar() &&
				newStatements.GetFirstFirstChild().HasOneChild() &&
				newStatements.GetFirstFirstChild().GetFirstChild().IsFunction() &&
				newStatements.GetSecondChild().IsExprResult() {
				inlinedFn := newStatements.GetFirstFirstChild().GetFirstChild()
				expr := newStatements.GetSecondChild().GetFirstChild()
				var call rhino.Node
				var assignedName *string
				if expr.IsAssign() && expr.GetSecondChild().IsCall() {
					call = expr.GetSecondChild()
					assignedName =
						jutil.StringPtr(modulePath.ToModuleName() +
							"_iife" +
							compiler.GetUniqueIdSupplier().GetUniqueId(ci))
				} else if expr.IsCall() {
					call = expr
				}

				if !call.IsNil() {
					newStatements =
						mutator.MutateWithoutRenaming(
							factoryLabel, inlinedFn, call, assignedName, false, false)
					if assignedName != nil {
						newName :=
							rhino.IR_Var(
								NodeUtil_NewNameWithOriginalName(
									compiler,
									*assignedName,
									fn,
									*expr.GetFirstChild().GetQualifiedName())).
								SrcrefTree(fn)
						if newStatements.HasChildren() &&
							newStatements.GetFirstChild().IsExprResult() &&
							newStatements.GetFirstFirstChild().IsAssign() &&
							newStatements.GetFirstFirstChild().GetFirstChild().IsName() &&
							newStatements.
								GetFirstFirstChild().
								GetFirstChild().
								GetString() == *assignedName {
							newName.
								GetFirstChild().
								AddChildToFront(
									newStatements.GetFirstFirstChild().GetSecondChild().Detach())
							newStatements.GetFirstChild().ReplaceWith(newName)
						} else {
							newStatements.AddChildToFront(newName)
						}
						expr.GetSecondChild().ReplaceWith(newName.GetFirstChild().CloneNode())
						newStatements.AddChildToBack(expr.GetParent().Detach())
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

// ProcessCommonJSModules_RewriteModule traverses a file and rewrite all
// references to imported names directly to the targeted module name.
//
// If a file is a CommonJS module, rewrite export statements. Typically
// exports create an alias - the rewriting tries to avoid such aliases.
type ProcessCommonJSModules_RewriteModule struct {
	NodeTraversal_AbstractPostOrderCallbackBase

	outer *ProcessCommonJSModules // the enclosing instance

	allowFullRewrite          bool
	exports                   []*ProcessCommonJSModules_ExportInfo
	imports                   []rhino.Node
	rewrittenClassExpressions []rhino.Node
	functionsToHoist          []rhino.Node
	defaultExportIsConst      bool
}

// Java: ProcessCommonJSModules.RewriteModule(boolean, ImmutableCollection<ExportInfo>, boolean)
func NewProcessCommonJSModules_RewriteModule(
	outer *ProcessCommonJSModules,
	allowFullRewrite bool,
	exports []*ProcessCommonJSModules_ExportInfo,
	defaultExportIsConst bool) *ProcessCommonJSModules_RewriteModule {
	r := &ProcessCommonJSModules_RewriteModule{outer: outer}
	r.imports = []rhino.Node{}
	r.rewrittenClassExpressions = []rhino.Node{}
	r.functionsToHoist = []rhino.Node{}
	r.allowFullRewrite = allowFullRewrite
	r.exports = exports
	r.defaultExportIsConst = defaultExportIsConst
	return r
}

// Java: ProcessCommonJSModules.RewriteModule.visit(NodeTraversal, Node, Node)
func (r *ProcessCommonJSModules_RewriteModule) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	compiler := r.outer.compiler
	switch n.GetToken() {
	case rhino.Token_SCRIPT:
		// Class names can't be changed during the middle of a traversal. Unlike functions,
		// the name can be the EMPTY token rather than just a zero length string.
		for _, clazz := range r.rewrittenClassExpressions {
			clazz.GetFirstChild().ReplaceWith(rhino.IR_Empty().Srcref(clazz.GetFirstChild()))
			t.ReportCodeChange()
		}

		ci := compiler.GetInput(n.GetInputId())
		moduleName := *ProcessCommonJSModules_GetModuleName(ci)

		// If a function is the direct module export, move it to the top.
		for i := 1; i < len(r.functionsToHoist); i++ {
			if r.functionsToHoist[i].
				GetFirstFirstChild().
				MatchesQualifiedName(r.outer.GetBasePropertyImport(moduleName)) {
				fncVar := r.functionsToHoist[i]
				r.functionsToHoist = slices.Delete(r.functionsToHoist, i, i+1)
				r.functionsToHoist = slices.Insert(r.functionsToHoist, 0, fncVar)
				break
			}
		}

		// Hoist functions in reverse order so that they maintain the same relative
		// order after hoisting.
		for i := len(r.functionsToHoist) - 1; i >= 0; i-- {
			functionExpr := r.functionsToHoist[i]
			scopeRoot := t.GetClosestHoistScopeRoot()
			insertionPoint := scopeRoot.GetFirstChild()
			if insertionPoint.IsNil() ||
				!(insertionPoint.IsVar() &&
					insertionPoint.GetFirstChild().GetString() == moduleName) {
				insertionPoint = rhino.NilNode
			}

			if insertionPoint.IsNil() {
				if scopeRoot.GetFirstChild() != functionExpr {
					scopeRoot.AddChildToFront(functionExpr.Detach())
				}
			} else if insertionPoint != functionExpr && insertionPoint.GetNext() != functionExpr {
				functionExpr.Detach().InsertAfter(insertionPoint)
			}
		}

		for _, export := range r.exports {
			r.VisitExport(t, export)
		}

		for _, require := range r.imports {
			r.VisitRequireCall(t, require)
		}
	case rhino.Token_CALL:
		if r.outer.IsCommonJsImport(n) {
			r.imports = append(r.imports, n)
		}
	case rhino.Token_VAR, rhino.Token_LET, rhino.Token_CONST:
		// Multiple declarations need split apart so that they can be refactored into
		// property assignments or removed altogether. Don't split declarations for
		// ES module export calls as it breaks AST and ES modules shouldn't be affected at all
		// by this pass.
		if n.HasMoreThanOneChild() && !NodeUtil_IsAnyFor(parent) && !parent.IsExport() {
			vars := r.SplitMultipleDeclarations(n)
			t.ReportCodeChange()
			for _, var_ := range vars {
				r.Visit(t, var_.GetFirstChild(), var_)
			}
		}

		// UMD Inlining can shadow global variables - these are just removed.
		//
		// var exports = exports;
		if n.GetFirstChild().HasChildren() &&
			n.GetFirstFirstChild().IsName() &&
			n.GetFirstChild().GetString() == n.GetFirstFirstChild().GetString() {
			n.Detach()
			t.ReportCodeChange()
			return
		}
	case rhino.Token_NAME:
		// If this is a name declaration with multiple names, it will be split apart when
		// the parent is visited and then revisit the children.
		if NodeUtil_IsNameDeclaration(n.GetParent()) &&
			n.GetParent().HasMoreThanOneChild() &&
			!NodeUtil_IsAnyFor(n.GetGrandparent()) {
			break
		}

		qName := n.GetQualifiedName()
		if qName == nil {
			break
		}
		nameDeclaration := t.GetScope().GetVar(*qName)
		if nameDeclaration != nil {
			if NodeUtil_IsLhsByDestructuring(n) {
				r.MaybeUpdateName(t, n, nameDeclaration)
			} else if !nameDeclaration.GetNode().IsNil() &&
				jutil.ObjectsEquals(nameDeclaration.GetNode().GetInputId(), n.GetInputId()) {
				// Avoid renaming a shadowed global
				//
				// var angular = angular;  // value is global ref
				enclosingDeclaration :=
					NodeUtil_GetEnclosingNode(
						n, func(node rhino.Node) bool { return node == nameDeclaration.GetNameNode() })

				if enclosingDeclaration.IsNil() ||
					enclosingDeclaration == n ||
					nameDeclaration.GetScope().(*Scope) != t.GetScope() {
					r.MaybeUpdateName(t, n, nameDeclaration)
				}
			}
			// Replace loose "module" references with an object literal
		} else if r.allowFullRewrite &&
			n.GetString() == processCommonJSModules_MODULE &&
			!(n.GetParent().IsGetProp() ||
				n.GetParent().IsGetElem() ||
				n.GetParent().IsTypeOf()) {
			objectLit := rhino.IR_Objectlit().Srcref(n)
			n.ReplaceWith(objectLit)
			t.ReportCodeChangeWithNode(objectLit)
		}
	case rhino.Token_GETPROP:
		if n.MatchesQualifiedName(processCommonJSModules_MODULE+".id") ||
			n.MatchesQualifiedName(processCommonJSModules_MODULE+".filename") {
			v := t.GetScope().GetVar(processCommonJSModules_MODULE)
			if v == nil || v.IsExtern() {
				n.ReplaceWith(rhino.IR_String(t.GetInput().GetPath().String()).Srcref(n))
			}
		} else if r.allowFullRewrite &&
			n.GetFirstChild().IsName() &&
			n.GetFirstChild().GetString() == processCommonJSModules_MODULE &&
			!n.MatchesQualifiedName(processCommonJSModules_MODULE+"."+processCommonJSModules_EXPORTS) {
			v := t.GetScope().GetVar(processCommonJSModules_MODULE)
			if v == nil || v.IsExtern() {
				n.GetFirstChild().ReplaceWith(rhino.IR_Objectlit().Srcref(n.GetFirstChild()))
			}
		}
	case rhino.Token_TYPEOF:
		if r.allowFullRewrite &&
			n.GetFirstChild().IsName() &&
			(n.GetFirstChild().GetString() == processCommonJSModules_MODULE ||
				n.GetFirstChild().GetString() == processCommonJSModules_EXPORTS) {
			v := t.GetScope().GetVar(n.GetFirstChild().GetString())
			if v == nil || v.IsExtern() {
				n.ReplaceWith(rhino.IR_String("object"))
			}
		}
	default:
	}

	r.FixTypeAnnotationsForNode(t, n)
}

// Java: ProcessCommonJSModules.RewriteModule.fixTypeAnnotationsForNode(NodeTraversal, Node)
func (r *ProcessCommonJSModules_RewriteModule) FixTypeAnnotationsForNode(t *NodeTraversal, n rhino.Node) {
	info := n.GetJSDocInfo()
	if info != nil {
		for it := info.TypeExprIter(); it.Next(); {
			r.FixTypeNode(t, it.Node())
		}
	}
}

// VisitRequireCall visits require calls. Rewrite require statements to be a
// direct reference to name of require module. By this point all references
// to the import alias should have already been renamed.
//
// Java: ProcessCommonJSModules.RewriteModule.visitRequireCall(NodeTraversal, Node)
func (r *ProcessCommonJSModules_RewriteModule) VisitRequireCall(t *NodeTraversal, require rhino.Node) {
	moduleName := r.outer.GetImportedModuleName(t, require)
	moduleRef :=
		NodeUtil_NewQName(r.outer.compiler, r.outer.GetBasePropertyImportWithRequireCall(moduleName, require)).
			SrcrefTree(require.GetSecondChild())
	require.ReplaceWith(moduleRef)

	t.ReportCodeChange()
}

// VisitExport visits export statements. Export statements can be either a
// direct assignment: module.exports = foo or a property assignment:
// module.exports.foo = foo; exports.foo = foo;
//
// Java: ProcessCommonJSModules.RewriteModule.visitExport(NodeTraversal, ExportInfo)
func (r *ProcessCommonJSModules_RewriteModule) VisitExport(t *NodeTraversal, export *ProcessCommonJSModules_ExportInfo) {
	compiler := r.outer.compiler
	root := r.outer.GetBaseQualifiedNameNode(export.node)
	rValue := NodeUtil_GetRValueOfLValue(root)

	// For object literal assignments to module.exports, convert them to
	// individual property assignments.
	//
	//     module.exports = { foo: bar};
	//
	// becomes
	//
	//     module.exports = {};
	//     module.exports.foo = bar;
	if processCommonJSModules_MODULE_EXPORTS.Matches(root) {
		if !rValue.IsNil() &&
			rValue.IsObjectLit() &&
			root.GetParent().IsAssign() &&
			root.GetGrandparent().IsExprResult() {
			if r.ExpandObjectLitAssignment(t, root, export.scope) {
				return
			}
		}
	}

	moduleName := *ProcessCommonJSModules_GetModuleName(t.GetInput())
	moduleInitialization := t.GetScope().GetVar(moduleName)

	// If this is an assignment to module.exports or exports, renaming
	// has already handled this case. Remove the export.
	var rValueVar *Var
	if !rValue.IsNil() && rValue.IsQualifiedName() {
		rValueVar = export.scope.GetVar(*rValue.GetQualifiedName())
	}

	if root.GetParent().IsAssign() &&
		root.GetGrandparent().IsExprResult() &&
		(!root.GetNext().IsNil() && (root.GetNext().IsName() || root.GetNext().IsGetProp())) &&
		root.GetGrandparent().IsExprResult() &&
		rValueVar != nil &&
		(NodeUtil_GetEnclosingScript(rValueVar.GetNameNode()).IsNil() ||
			(rValueVar.GetNameNode().HasParent() && !rValueVar.IsParam())) &&
		export.isInSupportedScope &&
		(rValueVar.GetNameNode().GetParent().IsNil() ||
			!NodeUtil_IsLhsByDestructuring(rValueVar.GetNameNode())) {
		root.GetGrandparent().Detach()
		t.ReportCodeChange()
		return
	}

	moduleName = moduleName + "." + processCommonJSModules_EXPORT_PROPERTY_NAME

	updatedExport :=
		NodeUtil_NewQNameWithOriginalName(compiler, moduleName, export.node, *export.node.GetQualifiedName())
	updatedExport.PutBooleanProp(rhino.Node_MODULE_EXPORT, true)
	exportIsConst :=
		r.defaultExportIsConst &&
			updatedExport.MatchesQualifiedName(
				r.outer.GetBasePropertyImport(*ProcessCommonJSModules_GetModuleName(t.GetInput()))) &&
			root == export.node &&
			NodeUtil_IsLValue(export.node)

	var changeScope rhino.Node

	if processCommonJSModules_MODULE_EXPORTS.Matches(root) &&
		!rValue.IsNil() &&
		export.scope.GetVar("module.exports") == nil &&
		root.GetParent().IsAssign() &&
		root.GetGrandparent().IsExprResult() &&
		moduleInitialization == nil {
		// Rewrite "module.exports = foo;" to "var moduleName = {default: foo};"
		parent := root.GetParent()
		exportName := rhino.IR_ExprResult(rhino.IR_Assign(updatedExport, rValue.Detach()))
		if exportIsConst {
			info := rhino.JSDocInfo_Builder_()
			info.RecordConstancy()
			exportName.GetFirstChild().SetJSDocInfo(info.Build())
		}
		parent.GetParent().ReplaceWith(exportName.SrcrefTree(root.GetParent()))
		changeScope = ChangeTracker_GetEnclosingChangeScopeRoot(parent)
	} else if !root.GetNext().IsNil() &&
		root.GetNext().IsName() &&
		rValueVar != nil &&
		rValueVar.IsGlobal() &&
		export.isInSupportedScope &&
		(rValueVar.GetNameNode().GetParent().IsNil() ||
			!NodeUtil_IsLhsByDestructuring(rValueVar.GetNameNode())) {
		// This is a where a module export assignment is used in a complex expression.
		// Before: `SOME_VALUE !== undefined && module.exports = SOME_VALUE`
		// After: `SOME_VALUE !== undefined && module$name`
		parent := root.GetParent()
		root.Detach()
		parent.ReplaceWith(root)
		if root == export.node {
			root = updatedExport
		}
		export.node.ReplaceWith(updatedExport)
		changeScope = ChangeTracker_GetEnclosingChangeScopeRoot(root)
	} else {
		// Other references to "module.exports" are just replaced with the module name.
		export.node.ReplaceWith(updatedExport)
		if updatedExport.GetParent().IsAssign() && exportIsConst {
			infoBuilder :=
				rhino.JSDocInfo_Builder_MaybeCopyFrom(updatedExport.GetParent().GetJSDocInfo())
			infoBuilder.RecordConstancy()
			updatedExport.GetParent().SetJSDocInfo(infoBuilder.Build())
		}

		changeScope = ChangeTracker_GetEnclosingChangeScopeRoot(updatedExport)
	}

	if !changeScope.IsNil() {
		compiler.ReportChangeToChangeScope(changeScope)
	}
}

// ExpandObjectLitAssignment: since CommonJS modules may have only a single
// export, it's common to see the export be an object pattern. We want to
// expand this to individual property assignments. If any individual
// property assignment has been renamed, it will be removed.
//
// We need to keep assignments which aren't names
//
// module.exports = { foo: bar, baz: function() {} }
//
// becomes
//
// module.exports.foo = bar; // removed later module.exports.baz = function() {};
//
// Java: ProcessCommonJSModules.RewriteModule.expandObjectLitAssignment(NodeTraversal, Node, Scope)
func (r *ProcessCommonJSModules_RewriteModule) ExpandObjectLitAssignment(t *NodeTraversal, export rhino.Node, scope *Scope) bool {
	jutil.CheckState(export.GetParent().IsAssign())
	insertionRef := export.GetGrandparent()
	jutil.CheckState(insertionRef.IsExprResult())
	insertionParent := insertionRef.GetParent()
	rhino.CheckNotNull(insertionParent)

	rValue := NodeUtil_GetRValueOfLValue(export)
	key := rValue.GetFirstChild()

	removedNodes := false
	for !key.IsNil() {
		if (!key.IsStringKey() || key.IsQuotedStringKey()) && !key.IsMemberFunctionDef() {
			key = key.GetNext()
			continue
		}

		lhs := rhino.IR_Getprop(export.CloneTree(), key.GetString())
		var value rhino.Node
		if key.IsStringKey() {
			value = key.RemoveFirstChild()
		} else if key.IsMemberFunctionDef() {
			value = key.RemoveFirstChild()
		}

		expr := rhino.IR_ExprResult(rhino.IR_Assign(lhs, value)).SrcrefTreeIfMissing(key)
		expr.InsertAfter(insertionRef)
		newExport := NewProcessCommonJSModules_ExportInfo(lhs.GetFirstChild(), scope)
		r.VisitExport(t, newExport)

		// Export statements can be removed in visitExport
		if !expr.IsNil() && expr.HasParent() {
			insertionRef = expr
		}

		currentKey := key
		key = key.GetNext()
		currentKey.Detach()
		removedNodes = true
	}

	if !rValue.HasChildren() {
		export.GetGrandparent().Detach()
		return true
	}

	if removedNodes {
		t.ReportCodeChangeWithNode(rValue)
	}
	return false
}

// MaybeUpdateName: given a name reference, check to see if it needs
// renamed.
//
// We handle 3 main cases: 1. References to an import alias. These are
// replaced with a direct reference to the imported module. 2. Names which
// are exported. These are rewritten to be the export assignment directly. 3.
// Global names: If a name is global to the script, add a suffix so it
// doesn't collide with any other global.
//
// Rewriting case 1 is safe to perform on all files. Cases 2 and 3 can only
// be done if this file is a commonjs module.
//
// Java: ProcessCommonJSModules.RewriteModule.maybeUpdateName(NodeTraversal, Node, Var)
func (r *ProcessCommonJSModules_RewriteModule) MaybeUpdateName(t *NodeTraversal, n rhino.Node, var_ *Var) {
	jutil.CheckNotNull(var_)
	jutil.CheckState(n.IsName() || n.IsGetProp())
	jutil.CheckState(n.HasParent())
	importedModuleName := r.GetModuleImportName(t, var_.GetNode())
	name := *n.GetQualifiedName()

	// Check if the name refers to a alias for a require('foo') import.
	if importedModuleName != nil && n != var_.GetNode() {
		// Reference the imported name directly, rather than the alias
		r.UpdateNameReference(t, n, name, *importedModuleName, false, false)

	} else if r.allowFullRewrite {
		exportedName := r.GetExportedName(t, n, var_)

		// We need to exclude the alias created by the require import. We assume dead
		// code elimination will remove these later.
		if (n != var_.GetNode() || n.GetParent().IsClass()) &&
			exportedName == nil &&
			(var_ == nil ||
				var_.GetNameNode().GetParent().IsNil() ||
				!NodeUtil_IsLhsByDestructuring(var_.GetNameNode())) {
			// The name is actually the export reference itself.
			// This will be handled later by visitExports.
			if n.GetParent().IsClass() && n.GetParent().GetFirstChild() == n {
				r.rewrittenClassExpressions = append(r.rewrittenClassExpressions, n.GetParent())
			}

			return
		}

		// Check if the name is used as an export
		if importedModuleName == nil &&
			exportedName != nil &&
			*exportedName != name &&
			!var_.IsParam() &&
			(var_.GetNameNode().GetParent().IsNil() ||
				!NodeUtil_IsLhsByDestructuring(var_.GetNameNode())) {
			exportPropIsConst :=
				r.defaultExportIsConst &&
					r.outer.GetBasePropertyImport(*ProcessCommonJSModules_GetModuleName(t.GetInput())) == *exportedName &&
					r.outer.GetBaseQualifiedNameNode(n) == n &&
					NodeUtil_IsLValue(n)
			r.UpdateNameReference(t, n, name, *exportedName, true, exportPropIsConst)

			// If it's a global name, rename it to prevent conflicts with other scripts
		} else if var_.IsGlobal() {
			currentModuleName := *ProcessCommonJSModules_GetModuleName(t.GetInput())

			if currentModuleName == name {
				return
			}

			// refs to 'exports' are handled separately.
			if processCommonJSModules_EXPORTS == name {
				return
			}

			// closure_test_suite looks for test*() functions
			if r.outer.compiler.GetOptions().ShouldExportTestFunctions() &&
				strings.HasPrefix(currentModuleName, "test") {
				return
			}

			newName := name + "$$" + currentModuleName
			r.UpdateNameReference(t, n, name, newName, false, false)
		}
	}
}

// UpdateNameReference: nameRef is the qualified name node; originalName of
// nameRef; newName for nameRef; requireFunctionExpressions whether named
// class or functions should be rewritten to variable assignments.
//
// Java: ProcessCommonJSModules.RewriteModule.updateNameReference(NodeTraversal, Node, String, String, boolean, boolean)
func (r *ProcessCommonJSModules_RewriteModule) UpdateNameReference(
	t *NodeTraversal,
	nameRef rhino.Node,
	originalName string,
	newName string,
	requireFunctionExpressions bool,
	qualifiedNameIsConst bool) {
	compiler := r.outer.compiler
	parent := nameRef.GetParent()
	rhino.CheckNotNull(parent)
	// checkNotNull(newName): a Go string cannot be null.
	newNameIsQualified := strings.IndexByte(newName, '.') >= 0
	newNameIsModuleExport :=
		newName == r.outer.GetBasePropertyImport(*ProcessCommonJSModules_GetModuleName(t.GetInput()))

	newNameDeclaration := t.GetScope().GetVar(newName)

	switch parent.GetToken() {
	case rhino.Token_CLASS:
		if parent.GetIndexOfChild(nameRef) == 0 &&
			(newNameIsQualified || requireFunctionExpressions) {
			// Refactor a named class to a class expression
			// We can't remove the class name during a traversal, so save it for later
			r.rewrittenClassExpressions = append(r.rewrittenClassExpressions, parent)

			newNameRef := NodeUtil_NewQNameWithOriginalName(compiler, newName, nameRef, originalName)
			if newNameIsModuleExport {
				newNameRef.PutBooleanProp(rhino.Node_MODULE_EXPORT, true)
			}

			var expr rhino.Node
			if !newNameIsQualified && newNameDeclaration == nil {
				expr = rhino.IR_LetWithValue(newNameRef, rhino.IR_NullNode()).SrcrefTreeIfMissing(nameRef)
			} else {
				expr =
					rhino.IR_ExprResult(rhino.IR_Assign(newNameRef, rhino.IR_NullNode())).SrcrefTreeIfMissing(nameRef)
				info := rhino.JSDocInfo_Builder_MaybeCopyFrom(parent.GetJSDocInfo())
				parent.SetJSDocInfo(nil)
				if qualifiedNameIsConst {
					info.RecordConstancy()
				}
				expr.GetFirstChild().SetJSDocInfo(info.Build())
				r.FixTypeAnnotationsForNode(t, expr.GetFirstChild())
			}
			parent.ReplaceWith(expr)
			if expr.IsLet() {
				expr.GetFirstFirstChild().ReplaceWith(parent)
			} else {
				expr.GetFirstChild().GetSecondChild().ReplaceWith(parent)
			}
		} else if parent.GetIndexOfChild(nameRef) == 1 {
			newNameRef := NodeUtil_NewQNameWithOriginalName(compiler, newName, nameRef, originalName)
			if newNameIsModuleExport {
				newNameRef.PutBooleanProp(rhino.Node_MODULE_EXPORT, true)
			}
			nameRef.ReplaceWith(newNameRef)
		} else {
			nameRef.SetString(newName)
			nameRef.SetOriginalName(jutil.StringPtr(originalName))
		}
	case rhino.Token_FUNCTION:
		if newNameIsQualified || requireFunctionExpressions {
			// Refactor a named function to a function expression
			if NodeUtil_IsFunctionExpression(parent) {
				// Don't refactor if the parent is a named function expression.
				// e.g. var foo = function foo() {};
				return
			}
			newNameRef := NodeUtil_NewQNameWithOriginalName(compiler, newName, nameRef, originalName)
			if newNameIsModuleExport {
				newNameRef.PutBooleanProp(rhino.Node_MODULE_EXPORT, true)
			}
			nameRef.SetString("")

			var expr rhino.Node
			if !newNameIsQualified && newNameDeclaration == nil {
				expr = rhino.IR_VarWithValue(newNameRef, rhino.IR_NullNode()).SrcrefTreeIfMissing(nameRef)
			} else {
				expr =
					rhino.IR_ExprResult(rhino.IR_Assign(newNameRef, rhino.IR_NullNode())).SrcrefTreeIfMissing(nameRef)
			}
			parent.ReplaceWith(expr)
			if expr.IsVar() {
				expr.GetFirstFirstChild().ReplaceWith(parent)
			} else {
				expr.GetFirstChild().GetSecondChild().ReplaceWith(parent)
				info := rhino.JSDocInfo_Builder_MaybeCopyFrom(parent.GetJSDocInfo())
				parent.SetJSDocInfo(nil)
				if qualifiedNameIsConst {
					info.RecordConstancy()
				}
				expr.GetFirstChild().SetJSDocInfo(info.Build())
				r.FixTypeAnnotationsForNode(t, expr.GetFirstChild())
			}
			r.functionsToHoist = append(r.functionsToHoist, expr)
		} else {
			nameRef.SetString(newName)
			nameRef.SetOriginalName(jutil.StringPtr(originalName))
		}
	case rhino.Token_VAR, rhino.Token_LET, rhino.Token_CONST:
		// Multiple declaration - needs split apart.
		if parent.HasMoreThanOneChild() && !NodeUtil_IsAnyFor(parent.GetParent()) {
			r.SplitMultipleDeclarations(parent)
			parent = nameRef.GetParent()
			newNameDeclaration = t.GetScope().GetVar(newName)
		}

		if newNameIsQualified {
			// Var declarations without initialization can simply
			// be removed if they are being converted to a property.
			if !nameRef.HasChildren() && parent.GetJSDocInfo() == nil {
				parent.Detach()
				break
			}

			// Refactor a var declaration to a getprop assignment
			getProp := NodeUtil_NewQNameWithOriginalName(compiler, newName, nameRef, originalName)
			if newNameIsModuleExport {
				getProp.PutBooleanProp(rhino.Node_MODULE_EXPORT, true)
			}
			info := parent.GetJSDocInfo()
			parent.SetJSDocInfo(nil)
			if nameRef.HasChildren() {
				assign := rhino.IR_Assign(getProp, nameRef.RemoveFirstChild())
				assign.SetJSDocInfo(info)
				expr := rhino.IR_ExprResult(assign).SrcrefTreeIfMissing(nameRef)
				parent.ReplaceWith(expr)
				infoBuilder := rhino.JSDocInfo_Builder_MaybeCopyFrom(info)
				parent.SetJSDocInfo(nil)
				if qualifiedNameIsConst {
					infoBuilder.RecordConstancy()
				}
				assign.SetJSDocInfo(infoBuilder.Build())
				r.FixTypeAnnotationsForNode(t, assign)
			} else {
				getProp.SetJSDocInfo(info)
				parent.ReplaceWith(rhino.IR_ExprResult(getProp).Srcref(getProp))
			}
		} else if newNameDeclaration != nil && newNameDeclaration.GetNameNode() != nameRef {
			// Variable is already defined. Convert this to an assignment.
			// If the variable declaration has no initialization, we simply
			// remove the node. This can occur when the variable which is exported
			// is declared in an outer scope but assigned in an inner one.
			if !nameRef.HasChildren() {
				parent.Detach()
				break
			}

			name := NodeUtil_NewNameWithOriginalName(compiler, newName, nameRef, originalName)
			assign := rhino.IR_Assign(name, nameRef.RemoveFirstChild())
			info := parent.GetJSDocInfo()
			if info != nil {
				parent.SetJSDocInfo(nil)
				assign.SetJSDocInfo(info)
			}

			parent.ReplaceWith(rhino.IR_ExprResult(assign).SrcrefTree(nameRef))
		} else {
			nameRef.SetString(newName)
			nameRef.SetOriginalName(jutil.StringPtr(originalName))
		}
	default:
		// Whenever possible, reuse the existing reference
		if !newNameIsQualified && nameRef.IsName() {
			nameRef.SetString(newName)
			nameRef.SetOriginalName(jutil.StringPtr(originalName))
		} else {
			var name rhino.Node
			if newNameIsQualified {
				name = NodeUtil_NewQNameWithOriginalName(compiler, newName, nameRef, originalName)
			} else {
				name = NodeUtil_NewNameWithOriginalName(compiler, newName, nameRef, originalName)
			}

			if newNameIsModuleExport {
				name.PutBooleanProp(rhino.Node_MODULE_EXPORT, true)
			}
			info := nameRef.GetJSDocInfo()
			if info != nil {
				nameRef.SetJSDocInfo(nil)
				name.SetJSDocInfo(info)
			}
			name.SrcrefTree(nameRef)
			nameRef.ReplaceWith(name)
			if nameRef.HasChildren() {
				name.AddChildrenToFront(nameRef.RemoveChildren())
			}
		}
	}

	t.ReportCodeChange()
}

// GetExportedName determines whether the given name Node n is referenced in
// an export. Returns: if the name is not used in an export, its own name; if
// the name node is actually the export target itself, nil.
//
// Java: ProcessCommonJSModules.RewriteModule.getExportedName(NodeTraversal, Node, Var)
func (r *ProcessCommonJSModules_RewriteModule) GetExportedName(t *NodeTraversal, n rhino.Node, var_ *Var) *string {
	if var_ == nil || !jutil.ObjectsEquals(var_.GetNode().GetInputId(), n.GetInputId()) {
		return n.GetQualifiedName()
	}

	baseExportName := r.outer.GetBasePropertyImport(*ProcessCommonJSModules_GetModuleName(t.GetInput()))

	for _, export := range r.exports {
		exportBase := r.outer.GetBaseQualifiedNameNode(export.node)
		exportRValue := NodeUtil_GetRValueOfLValue(exportBase)

		if exportRValue.IsNil() {
			continue
		}

		exportedName := r.GetExportedNameNode(export)
		// We don't want to handle the export itself
		if export.isInSupportedScope &&
			(exportRValue == n ||
				((NodeUtil_IsClassExpression(exportRValue) ||
					NodeUtil_IsFunctionExpression(exportRValue)) &&
					exportedName == n)) {
			return nil
		}

		exportBaseQName := exportBase.GetQualifiedName()

		if exportRValue.IsObjectLit() {
			if exportBaseQName == nil || "module.exports" != *exportBaseQName {
				return n.GetQualifiedName()
			}

			key := exportRValue.GetFirstChild()
			keyIsExport := false
			for !key.IsNil() {
				if key.IsStringKey() &&
					!key.IsQuotedStringKey() &&
					NodeUtil_IsValidPropertyName(
						r.outer.compiler.GetOptions().GetLanguageIn().ToFeatureSet(), key.GetString()) {
					if key.GetFirstChild().IsQualifiedName() {
						if key.GetFirstChild() == n {
							return nil
						}

						valVar := t.GetScope().GetVar(*key.GetFirstChild().GetQualifiedName())
						if valVar != nil && valVar.GetNameNode() == var_.GetNameNode() {
							keyIsExport = true
							break
						}
					}
				}

				key = key.GetNext()
			}
			if !key.IsNil() && keyIsExport {
				if export.isInSupportedScope {
					return jutil.StringPtr(baseExportName + "." + key.GetString())
				} else {
					return n.GetQualifiedName()
				}
			}
		} else {
			if !export.isInSupportedScope {
				return n.GetQualifiedName()
			}
			if var_.GetNameNode() == exportedName {
				var exportPrefix string
				if strings.HasPrefix(*exportBaseQName, processCommonJSModules_MODULE) {
					exportPrefix = processCommonJSModules_MODULE + "." + processCommonJSModules_EXPORTS
				} else {
					exportPrefix = processCommonJSModules_EXPORTS
				}

				if jutil.StringLength(*exportBaseQName) == jutil.StringLength(exportPrefix) {
					return jutil.StringPtr(baseExportName)
				}

				return jutil.StringPtr(baseExportName + jutil.StringSubstringFrom(*exportBaseQName, jutil.StringLength(exportPrefix)))
			}
		}
	}
	return n.GetQualifiedName()
}

// Java: ProcessCommonJSModules.RewriteModule.getExportedNameNode(ExportInfo)
func (r *ProcessCommonJSModules_RewriteModule) GetExportedNameNode(info *ProcessCommonJSModules_ExportInfo) rhino.Node {
	qNameBase := r.outer.GetBaseQualifiedNameNode(info.node)
	rValue := NodeUtil_GetRValueOfLValue(qNameBase)

	if rValue.IsNil() {
		return rhino.NilNode
	}

	if NodeUtil_IsFunctionExpression(rValue) {
		return rValue.GetFirstChild()
	}

	if NodeUtil_IsClassExpression(rValue) && rValue.GetFirstChild() == qNameBase {
		return rValue.GetFirstChild()
	}

	if !rValue.IsName() {
		return rhino.NilNode
	}

	var_ := info.scope.GetVar(rValue.GetString())
	if var_ == nil {
		return rhino.NilNode
	}

	return var_.GetNameNode()
}

// GetModuleImportName determines if the given Node n is an alias created by
// a module import. Returns nil if it's not an alias or the imported module
// name.
//
// Java: ProcessCommonJSModules.RewriteModule.getModuleImportName(NodeTraversal, Node)
func (r *ProcessCommonJSModules_RewriteModule) GetModuleImportName(t *NodeTraversal, n rhino.Node) *string {
	var rValue rhino.Node
	propSuffix := ""
	parent := n.GetParent()
	if !parent.IsNil() && parent.IsStringKey() {
		grandparent := parent.GetParent()
		if grandparent.IsObjectPattern() && grandparent.GetParent().IsDestructuringLhs() {
			rValue = grandparent.GetNext()
			propSuffix = "." + parent.GetString()
		}
	}
	if propSuffix == "" && !parent.IsNil() {
		rValue = NodeUtil_GetRValueOfLValue(n)
	}

	if rValue.IsNil() {
		return nil
	}

	if rValue.IsCall() && r.outer.IsCommonJsImport(rValue) {
		return jutil.StringPtr(r.outer.GetBasePropertyImportWithRequireCall(r.outer.GetImportedModuleName(t, rValue), rValue) + propSuffix)
	} else if rValue.IsGetProp() && r.outer.IsCommonJsImport(rValue.GetFirstChild()) {
		// var foo = require('bar').foo;
		importName :=
			r.outer.GetBasePropertyImportWithRequireCall(
				r.outer.GetImportedModuleName(t, rValue.GetFirstChild()), rValue.GetFirstChild())

		suffix := rValue.GetString()

		return jutil.StringPtr(importName + "." + suffix + propSuffix)
	}

	return nil
}

// FixTypeNode updates any type references in JSDoc annotations to account
// for all the rewriting we've done.
//
// Java: ProcessCommonJSModules.RewriteModule.fixTypeNode(NodeTraversal, Node)
func (r *ProcessCommonJSModules_RewriteModule) FixTypeNode(t *NodeTraversal, typeNode rhino.Node) {
	if typeNode.IsStringLit() {
		name := typeNode.GetString()
		// Type nodes can be module paths.
		if ModuleLoader_IsPathIdentifier(name) {
			lastSlash := jutil.StringLastIndexOf(name, "/")
			endIndex := jutil.StringIndexOfCharFrom(name, '.', lastSlash)
			var localTypeName *string
			if endIndex == -1 {
				endIndex = jutil.StringLength(name)
			} else {
				localTypeName = jutil.StringPtr(jutil.StringSubstringFrom(name, endIndex))
			}

			moduleName := jutil.StringSubstring(name, 0, endIndex)
			globalModuleName := r.outer.GetImportedModuleNameWithImportPath(t, typeNode, moduleName)
			baseImportProperty := r.outer.GetBasePropertyImport(globalModuleName)
			if localTypeName == nil {
				typeNode.SetString(baseImportProperty)
			} else {
				typeNode.SetString(baseImportProperty + *localTypeName)
			}

		} else {
			// A type node can be a getprop. Any portion of the getprop
			// can be either an import alias or export alias. Check each
			// segment.
			wasRewritten := false
			endIndex := -1
			for endIndex < jutil.StringLength(name) {
				endIndex = jutil.StringIndexOfCharFrom(name, '.', endIndex+1)
				if endIndex == -1 {
					endIndex = jutil.StringLength(name)
				}
				baseName := jutil.StringSubstring(name, 0, endIndex)
				suffix := ""
				if endIndex < jutil.StringLength(name) {
					suffix = jutil.StringSubstringFrom(name, endIndex)
				}
				typeDeclaration := t.GetScope().GetVar(baseName)

				// Make sure we can find a variable declaration (and it's in this file)
				if typeDeclaration != nil && !typeDeclaration.GetNode().IsNil() &&
					jutil.ObjectsEquals(typeDeclaration.GetNode().GetInputId(), typeNode.GetInputId()) {
					importedModuleName := r.GetModuleImportName(t, typeDeclaration.GetNode())

					// If the name is an import alias, rewrite it to be a reference to the
					// module name directly
					if importedModuleName != nil {
						typeNode.SetString(*importedModuleName + suffix)
						typeNode.SetOriginalName(jutil.StringPtr(name))
						wasRewritten = true
						break
					} else if r.allowFullRewrite {
						// Names referenced in export statements can only be rewritten in
						// commonjs modules.
						exportedName := r.GetExportedName(t, typeNode, typeDeclaration)
						if exportedName != nil && *exportedName != name {
							typeNode.SetString(*exportedName + suffix)
							typeNode.SetOriginalName(jutil.StringPtr(name))
							wasRewritten = true
							break
						}
					}
				}
			}

			// If the name was neither an import alias or referenced in an export,
			// We still may need to rename it if it's global
			if !wasRewritten && r.allowFullRewrite {
				endIndex = jutil.StringIndexOfChar(name, '.')
				if endIndex == -1 {
					endIndex = jutil.StringLength(name)
				}
				baseName := jutil.StringSubstring(name, 0, endIndex)
				typeDeclaration := t.GetScope().GetVar(baseName)
				if typeDeclaration != nil && typeDeclaration.IsGlobal() {
					moduleName := *ProcessCommonJSModules_GetModuleName(t.GetInput())
					newName := baseName + "$$" + moduleName
					if endIndex < jutil.StringLength(name) {
						newName += jutil.StringSubstringFrom(name, endIndex)
					}

					typeNode.SetString(newName)
					typeNode.SetOriginalName(jutil.StringPtr(name))
				}
			}
		}
	}

	for child := typeNode.GetFirstChild(); !child.IsNil(); child = child.GetNext() {
		r.FixTypeNode(t, child)
	}
}

// Java: ProcessCommonJSModules.RewriteModule.splitMultipleDeclarations(Node)
func (r *ProcessCommonJSModules_RewriteModule) SplitMultipleDeclarations(var_ rhino.Node) []rhino.Node {
	jutil.CheckState(NodeUtil_IsNameDeclaration(var_))
	vars := []rhino.Node{}
	info := var_.GetJSDocInfo()
	for !var_.GetSecondChild().IsNil() {
		newVar := rhino.NewNodeWithChild(var_.GetToken(), var_.RemoveFirstChild())

		if info != nil {
			newVar.SetJSDocInfo(info.Clone())
		}

		newVar.Srcref(var_)
		newVar.InsertBefore(var_)
		vars = append(vars, newVar)
	}
	vars = append(vars, var_)
	return vars
}
