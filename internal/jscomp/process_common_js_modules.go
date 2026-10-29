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

	"github.com/paullepoulpe/jsgo/internal/ast"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
	"github.com/paullepoulpe/jsgo/internal/strpool"
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

// ProcessCommonJSModules rewrites a CommonJS module into a form that can be
// safely concatenated. Does not add a function around the module body but
// instead adds suffixes to global variables to avoid conflicts. Calls to
// require are changed to reference the required module directly.
type ProcessCommonJSModules struct {
	NodeTraversal_AbstractPreOrderCallbackBase

	compiler AbstractCompiler
	finder   commonJsFindImportsAndExports
	rewriter commonJsRewriteModule
}

const (
	processCommonJSModules_EXPORTS                   = "exports"
	processCommonJSModules_MODULE                    = "module"
	processCommonJSModules_REQUIRE                   = "require"
	processCommonJSModules_WEBPACK_REQUIRE           = "__webpack_require__"
	processCommonJSModules_WEBPACK_REQUIRE_NAMESPACE = "__webpack_require__.t"
	processCommonJSModules_EXPORT_PROPERTY_NAME      = "default"
)

var (
	processCommonJSModules_exportsID = uint32(strpool.Intern(processCommonJSModules_EXPORTS))
	processCommonJSModules_moduleID  = uint32(strpool.Intern(processCommonJSModules_MODULE))
)

var ProcessCommonJSModules_UNKNOWN_REQUIRE_ENSURE = DiagnosticType_Warning(
	"JSC_COMMONJS_UNKNOWN_REQUIRE_ENSURE_ERROR", "Unrecognized require.ensure call: {0}")

var ProcessCommonJSModules_SUSPICIOUS_EXPORTS_ASSIGNMENT = DiagnosticType_Warning(
	"JSC_COMMONJS_SUSPICIOUS_EXPORTS_ASSIGNMENT",
	"Suspicious re-assignment of \"exports\" variable."+
		" Did you actually intend to export something?")

// NewProcessCommonJSModules creates a new ProcessCommonJSModules pass.
//
// Java: ProcessCommonJSModules(AbstractCompiler)
func NewProcessCommonJSModules(compiler AbstractCompiler) *ProcessCommonJSModules {
	p := &ProcessCommonJSModules{compiler: compiler}
	p.finder.outer = p
	p.rewriter.outer = p
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
	}
	if !n.IsScript() {
		return false
	}

	if p.compiler.GetOptions().GetModuleResolutionMode() == ModuleLoader_ResolutionMode_WEBPACK {
		p.removeWebpackModuleShim(n)
	}

	finder := &p.finder
	finder.reset()
	moduleLoader := p.compiler.GetModuleLoader()
	moduleLoaderErrorHandler := moduleLoader.GetErrorHandler()
	moduleLoader.SetErrorHandler(finder)
	NodeTraversal_Traverse(p.compiler, n, finder)

	moduleType := p.compiler.GetInput(n.GetInputId()).GetJsModuleType()
	forceModuleDetection := moduleType == CompilerInput_ModuleType_IMPORTED_SCRIPT
	defaultExportIsConst := true

	isCommonJsModule := finder.isCommonJsModule()
	p.rewriter.reset(isCommonJsModule || forceModuleDetection, defaultExportIsConst)
	if isCommonJsModule || forceModuleDetection {
		if len(finder.umdPatterns) != 0 {
			needsRetraverse := finder.replaceUmdPatterns()
			if p.removeIIFEWrapper(n) {
				needsRetraverse = true
			}
			if needsRetraverse {
				finder.reset()
				moduleLoader.SetErrorHandler(finder)
				NodeTraversal_Traverse(p.compiler, n, finder)
			}
		}

		p.rewriter.defaultExportIsConst = finder.initializeModule()
		p.rewriter.exports = append(p.rewriter.exports, finder.moduleExports...)
		p.rewriter.exports = append(p.rewriter.exports, finder.exports...)
	}
	finder.reportModuleErrors()
	moduleLoader.SetErrorHandler(moduleLoaderErrorHandler)

	NodeTraversal_Traverse(p.compiler, n, &p.rewriter)
	return false
}

func commonJsModuleName(input *CompilerInput) string {
	return input.GetPath().ToModuleName()
}

func (p *ProcessCommonJSModules) getBasePropertyImportWithRequireCall(moduleName string, requireCall rhino.Node) string {
	jutil.CheckArgument(requireCall.IsCall())
	if p.compiler.GetOptions().GetModuleResolutionMode() == ModuleLoader_ResolutionMode_WEBPACK &&
		requireCall.GetFirstChild().MatchesQualifiedName(processCommonJSModules_WEBPACK_REQUIRE_NAMESPACE) {
		return moduleName
	}
	return p.getBasePropertyImport(moduleName)
}

func (p *ProcessCommonJSModules) getBasePropertyImport(moduleName string) string {
	moduleType := p.compiler.GetModuleTypeByName(moduleName)
	if moduleType != nil && *moduleType != CompilerInput_ModuleType_COMMONJS {
		return moduleName
	}
	return moduleName + "." + processCommonJSModules_EXPORT_PROPERTY_NAME
}

func (p *ProcessCommonJSModules) isCommonJsImport(requireCall rhino.Node) bool {
	return ProcessCommonJSModules_IsCommonJsImport(requireCall, p.compiler.GetOptions().GetModuleResolutionMode())
}

// ProcessCommonJSModules_IsCommonJsImport recognizes if a node is a module import:
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

func (p *ProcessCommonJSModules) getCommonJsImportPath(requireCall rhino.Node) string {
	return ProcessCommonJSModules_GetCommonJsImportPath(requireCall, p.compiler.GetOptions().GetModuleResolutionMode())
}

// Java: ProcessCommonJSModules.getCommonJsImportPath(Node, ModuleLoader.ResolutionMode)
func ProcessCommonJSModules_GetCommonJsImportPath(
	requireCall rhino.Node, resolutionMode ModuleLoader_ResolutionMode) string {
	first, ok := ast.MustCall(requireCall).FirstArg()
	jutil.CheckState(ok)
	if resolutionMode == ModuleLoader_ResolutionMode_WEBPACK {
		pathArgument := first.Node()
		if second, ok := first.Next(); ok {
			pathArgument = second.Node()
		}
		if pathArgument.IsNumber() {
			return strconv.Itoa(int(jutil.D2I(pathArgument.GetDouble())))
		}
		return pathArgument.GetString()
	}
	return first.Node().GetString()
}

func (p *ProcessCommonJSModules) getImportedModuleName(t *NodeTraversal, requireCall rhino.Node) string {
	return p.getImportedModuleNameWithImportPath(t, requireCall, p.getCommonJsImportPath(requireCall))
}

func (p *ProcessCommonJSModules) getImportedModuleNameWithImportPath(t *NodeTraversal, n rhino.Node, importPath string) string {
	modulePath := t.GetInput().
		GetPath().
		ResolveJsModule(importPath, processCommonJSModules_sourceFileName(n), n.GetLineno(), n.GetCharno())
	if modulePath == nil {
		return ModuleIdentifier_ForFile(importPath).ModuleName()
	}
	return modulePath.ToModuleName()
}

// processCommonJSModules_sourceFileName returns n.GetSourceFileName() as a Go string.
func processCommonJSModules_sourceFileName(n rhino.Node) string {
	if s := n.GetSourceFileName(); s != nil {
		return *s
	}
	return ""
}

// ProcessCommonJSModules_IsCommonJsExport recognizes if a node is a module export:
//   - module.exports = something;
//   - module.exports.something = something;
//   - exports.something = something;
//
// Java: ProcessCommonJSModules.isCommonJsExport(NodeTraversal, Node, ModuleLoader.ResolutionMode)
func ProcessCommonJSModules_IsCommonJsExport(
	t *NodeTraversal, export rhino.Node, resolutionMode ModuleLoader_ResolutionMode) bool {
	if processCommonJSModules_MODULE_EXPORTS.Matches(export) ||
		processCommonJSModules_isModuleExportsElement(export) {
		v := t.GetScope().GetVarByID(processCommonJSModules_moduleID)
		return v == nil || v.IsExtern()
	}
	if export.IsName() && export.GetStringID() == processCommonJSModules_exportsID {
		v := t.GetScope().GetVarByID(processCommonJSModules_exportsID)
		return v == nil || v.IsGlobal()
	}
	return false
}

func processCommonJSModules_isModuleExportsElement(n rhino.Node) bool {
	getElem, ok := ast.AsGetElem(n)
	if !ok || !getElem.Object().Node().MatchesQualifiedName(processCommonJSModules_MODULE) {
		return false
	}
	index, ok := getElem.Index().AsStringLit()
	return ok && index.Value() == processCommonJSModules_EXPORTS
}

func (p *ProcessCommonJSModules) isCommonJsExport(t *NodeTraversal, export rhino.Node) bool {
	return ProcessCommonJSModules_IsCommonJsExport(
		t, export, p.compiler.GetOptions().GetModuleResolutionMode())
}

type commonJsExportInfo struct {
	node               rhino.Node
	scope              *Scope
	isInSupportedScope bool
}

func newCommonJsExportInfo(node rhino.Node, scope *Scope) commonJsExportInfo {
	disqualifyingParent := NodeUtil_GetEnclosingNode(
		node, func(n rhino.Node) bool { return n.IsIf() || n.IsHook() || n.IsFunction() || n.IsArrowFunction() })
	return commonJsExportInfo{
		node:               node,
		scope:              scope,
		isInSupportedScope: disqualifyingParent.IsNil(),
	}
}

func baseQualifiedNameNode(n rhino.Node) rhino.Node {
	refParent := n
	for refParent.HasParent() && refParent.GetParent().IsQualifiedName() {
		refParent = refParent.GetParent()
	}
	return refParent
}

var (
	processCommonJSModules_DEFINE_AMD        = rhino.QualifiedName_Of("define.amd")
	processCommonJSModules_WINDOW_DEFINE     = rhino.QualifiedName_Of("window.define")
	processCommonJSModules_WINDOW_DEFINE_AMD = rhino.QualifiedName_Of("window.define.amd")
	processCommonJSModules_REQUIRE_ENSURE    = rhino.QualifiedName_Of("require.ensure")
	processCommonJSModules_GOOG_PROVIDE      = rhino.QualifiedName_Of("goog.provide")
	processCommonJSModules_GOOG_MODULE       = rhino.QualifiedName_Of("goog.module")
	processCommonJSModules_MODULE_EXPORTS    = rhino.QualifiedName_Of("module.exports")
)

type commonJsFindImportsAndExports struct {
	outer *ProcessCommonJSModules

	hasGoogProvideOrModule bool
	script                 rhino.Node

	umdPatterns   []commonJsUmdPattern
	moduleExports []commonJsExportInfo
	exports       []commonJsExportInfo
	errors        []*JSError
}

func (f *commonJsFindImportsAndExports) reset() {
	f.hasGoogProvideOrModule = false
	f.script = rhino.NilNode
	f.umdPatterns = f.umdPatterns[:0]
	f.moduleExports = f.moduleExports[:0]
	f.exports = f.exports[:0]
	clear(f.errors)
	f.errors = f.errors[:0]
}

func (f *commonJsFindImportsAndExports) isCommonJsModule() bool {
	return (len(f.exports) != 0 || len(f.moduleExports) != 0) && !f.hasGoogProvideOrModule
}

// Report implements ErrorHandler for ModuleLoader.SetErrorHandler.
func (f *commonJsFindImportsAndExports) Report(ignoredLevel *CheckLevel, error_ *JSError) {
	f.errors = append(f.errors, error_)
}

var commonJsFindImportsAndExportsVisitTokens = func() ast.TokenSet {
	var s ast.TokenSet
	s.Add(rhino.Token_EXPR_RESULT)
	s.Add(rhino.Token_CALL)
	s.Add(rhino.Token_GETPROP)
	s.Add(rhino.Token_GETELEM)
	s.Add(rhino.Token_NAME)
	s.Add(rhino.Token_THIS)
	return s
}()

var (
	_ nodeTraversal_visitTokens                         = (*commonJsFindImportsAndExports)(nil)
	_ nodeTraversal_onlyScriptAndFunctionShouldTraverse = (*commonJsFindImportsAndExports)(nil)
)

func (*commonJsFindImportsAndExports) visitTokens() *ast.TokenSet {
	return &commonJsFindImportsAndExportsVisitTokens
}

func (*commonJsFindImportsAndExports) shouldTraverseOnlyScriptAndFunction() bool {
	return true
}

func (f *commonJsFindImportsAndExports) ShouldTraverse(nodeTraversal *NodeTraversal, n rhino.Node, parent rhino.Node) bool {
	if n.IsScript() {
		jutil.CheckState(f.script.IsNil())
		f.script = n
	}
	return true
}

func (f *commonJsFindImportsAndExports) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	if NodeUtil_IsShallowStatementTree(parent) && n.IsExprResult() {
		maybeGetProp := n.GetFirstFirstChild()
		if !maybeGetProp.IsNil() &&
			(processCommonJSModules_GOOG_PROVIDE.Matches(maybeGetProp) || processCommonJSModules_GOOG_MODULE.Matches(maybeGetProp)) {
			f.hasGoogProvideOrModule = true
		}
	}

	if call, ok := ast.AsCall(n); ok && processCommonJSModules_REQUIRE_ENSURE.Matches(call.Callee().Node()) {
		f.visitRequireEnsureCall(t, call)
	}

	if processCommonJSModules_MODULE_EXPORTS.Matches(n) ||
		processCommonJSModules_isModuleExportsElement(n) {
		if f.outer.isCommonJsExport(t, n) {
			f.moduleExports = append(f.moduleExports, newCommonJsExportInfo(n, t.GetScope()))

			if umdTestAncestor, ok := getOutermostUmdTest(parent); ok && !isInIfTest(n) {
				if !hasUmdPattern(f.umdPatterns, umdTestAncestor.enclosingIf) {
					f.umdPatterns = append(f.umdPatterns, commonJsUmdPattern{
						ifRoot:       umdTestAncestor.enclosingIf,
						activeBranch: umdTestAncestor.activeBranch,
					})
				}
			}
		}
	} else if processCommonJSModules_DEFINE_AMD.Matches(n) || processCommonJSModules_WINDOW_DEFINE_AMD.Matches(n) {
		if umdTestAncestor, ok := getOutermostUmdTest(parent); ok && isInIfTest(n) {
			if !hasUmdPattern(f.umdPatterns, umdTestAncestor.enclosingIf) &&
				processCommonJSModules_isIfWithoutElse(umdTestAncestor.enclosingIf) {
				f.umdPatterns = append(f.umdPatterns, commonJsUmdPattern{
					ifRoot:       umdTestAncestor.enclosingIf,
					activeBranch: rhino.NilNode,
				})
			}
		}
	}

	if n.IsName() && n.GetStringID() == processCommonJSModules_exportsID {
		v := t.GetScope().GetVarByID(processCommonJSModules_exportsID)
		if v == nil || v.IsGlobal() {
			qNameRoot := baseQualifiedNameNode(n)
			if !qNameRoot.IsNil() &&
				qNameRoot.MatchesQualifiedName(processCommonJSModules_EXPORTS) &&
				NodeUtil_IsLValue(qNameRoot) {
				// Match `exports = module.exports` or `exports = module.exports = ...`
				next := n.GetNext()
				if n.GetGrandparent().IsExprResult() &&
					!next.IsNil() &&
					((next.IsGetProp() &&
						processCommonJSModules_MODULE_EXPORTS.Matches(next)) ||
						(next.IsAssign() &&
							processCommonJSModules_MODULE_EXPORTS.Matches(next.GetFirstChild()))) {
					f.exports = append(f.exports, newCommonJsExportInfo(n, t.GetScope()))
				} else if !f.hasGoogProvideOrModule &&
					(v == nil || (v.GetNameNode().IsNil() && v.GetNameNode().GetFirstChild() != n)) {
					f.errors = append(f.errors, JSError_MakeWithNode(qNameRoot, ProcessCommonJSModules_SUSPICIOUS_EXPORTS_ASSIGNMENT))
				}
			} else {
				f.exports = append(f.exports, newCommonJsExportInfo(n, t.GetScope()))

				if umdTestAncestor, ok := getOutermostUmdTest(parent); ok && !isInIfTest(n) {
					if !hasUmdPattern(f.umdPatterns, umdTestAncestor.enclosingIf) {
						f.umdPatterns = append(f.umdPatterns, commonJsUmdPattern{
							ifRoot:       umdTestAncestor.enclosingIf,
							activeBranch: umdTestAncestor.activeBranch,
						})
					}
				}
			}
		}
	} else if n.IsThis() && n.GetParent().IsGetProp() && t.InGlobalScope() {
		f.exports = append(f.exports, newCommonJsExportInfo(n, t.GetScope()))
	}

	if f.outer.isCommonJsImport(n) {
		f.visitRequireCall(t, n, parent)
	}
}

func (f *commonJsFindImportsAndExports) visitRequireCall(t *NodeTraversal, require rhino.Node, parent rhino.Node) {
	if !NodeUtil_IsExpressionResultUsed(require) &&
		parent.IsExprResult() &&
		NodeUtil_IsStatementBlock(parent.GetParent()) {
		t.GetInput().
			GetPath().
			ResolveJsModule(
				f.outer.getCommonJsImportPath(require),
				processCommonJSModules_sourceFileName(require),
				require.GetLineno(),
				require.GetCharno())
		grandparent := parent.GetParent()
		parent.Detach()
		f.outer.compiler.ReportChangeToEnclosingScope(grandparent)
	}
}

func (f *commonJsFindImportsAndExports) visitRequireEnsureCall(t *NodeTraversal, ensure ast.Call) {
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
				strconv.Itoa(call.GetChildCount())))
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
	callback.Params().Node().RemoveChildren()
	call.RemoveChildren()
	call.PutBooleanProp(rhino.Node_FREE_CALL, true)
	call.AddChildToFront(callback.Node())

	t.ReportCodeChange()
}

func (f *commonJsFindImportsAndExports) reportModuleErrors() {
	for _, error_ := range f.errors {
		f.outer.compiler.Report(error_)
	}
}

func (f *commonJsFindImportsAndExports) initializeModule() bool {
	compiler := f.outer.compiler
	ci := compiler.GetInput(f.script.GetInputId())
	modulePath := ci.GetPath()
	if modulePath == nil {
		return true
	}

	moduleName := modulePath.ToModuleName()

	w := 0
	for _, export := range f.exports {
		if NodeUtil_GetEnclosingScript(export.node).IsNil() {
			continue
		}
		qNameBase := baseQualifiedNameNode(export.node)
		exportAssign, isAssignTarget := processCommonJSModules_assignToTarget(export.node)
		if export.node == qNameBase &&
			isAssignTarget &&
			export.node.GetGrandparent().IsExprResult() {
			value := exportAssign.Value()

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
				compiler.ReportChangeToEnclosingScope(changeRoot)
				continue
			} else if assign, ok := value.AsAssign(); ok &&
				assign.Target().Node().MatchesQualifiedName(processCommonJSModules_MODULE+"."+processCommonJSModules_EXPORTS) {
				exportAssign.Node().ReplaceWith(assign.Detach().Node())
				compiler.ReportChangeToEnclosingScope(assign.Node())
				continue
			}
		} else if getElem, prop, ok := processCommonJSModules_babelDefaultInterop(export.node); ok {
			parent := getElem.Node()
			grandparent := parent.GetParent()
			prop.Detach()
			parent.ReplaceWith(rhino.IR_Getprop(export.node.Detach(), prop.Value()).Srcref(parent))
			compiler.ReportChangeToEnclosingScope(grandparent)
		}
		f.exports[w] = export
		w++
	}
	f.exports = f.exports[:w]

	w = 0
	for _, export := range f.moduleExports {
		if NodeUtil_GetEnclosingScript(export.node).IsNil() {
			continue
		}
		if getElem, ok := ast.AsGetElem(export.node); ok {
			// module['exports'] becomes module.exports.
			object, prop := getElem.Object(), getElem.Index()
			prop.Detach()
			newExport := newCommonJsExportInfo(
				rhino.IR_Getprop(object.Detach().Node(), prop.Node().GetString()), export.scope)
			export.node.ReplaceWith(newExport.node)
			compiler.ReportChangeToEnclosingScope(newExport.node)
			f.moduleExports[w] = newExport
			w++
			continue
		}
		f.moduleExports[w] = export
		w++
	}
	f.moduleExports = f.moduleExports[:w]

	directAssignments := 0
	for _, export := range f.moduleExports {
		if NodeUtil_GetEnclosingScript(export.node).IsNil() {
			continue
		}
		base := baseQualifiedNameNode(export.node)
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

func processCommonJSModules_assignToTarget(n rhino.Node) (ast.Assign, bool) {
	assign, ok := ast.AsAssign(n.GetParent())
	return assign, ok && assign.Target().Is(n)
}

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

func processCommonJSModules_isNamespaceLiteral(obj ast.ObjectLit) bool {
	for member := range obj.Members() {
		if !processCommonJSModules_isNamespaceMember(member) {
			return false
		}
	}
	return true
}

func processCommonJSModules_isNamespaceMember(member ast.ObjectMember) bool {
	if key, ok := member.AsStringKey(); ok {
		return !key.IsQuoted()
	}
	_, ok := member.AsMemberFunctionDef()
	return ok
}
