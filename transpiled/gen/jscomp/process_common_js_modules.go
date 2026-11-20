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
// Port of com/google/javascript/jscomp/ProcessCommonJSModules.java (generated mechanically by j2go).

package jscomp

import (
	"strconv"
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

var lambdaSingleton_process_common_js_modules_1 j2gort.Ext_com_google_common_base_Predicate

func init() {
	lambdaSingleton_process_common_js_modules_1 = j2gort.NewNodeLambda_Ext_com_google_common_base_Predicate(func(rawNode_11 uint32) bool {
		n := pkg_rhino.Node(rawNode_11)
		_ = n
		return (((n.IsIf() || n.IsHook()) || n.IsFunction()) || n.IsArrowFunction())
	})
}

// Java class: com.google.javascript.jscomp.ProcessCommonJSModules
type ProcessCommonJSModules struct {
	NodeTraversal_AbstractPreOrderCallbackBase
	compiler AbstractCompiler `gson:"compiler"`
}

func (this_ *ProcessCommonJSModules) BindSelf_(s any) {
	_ = s
}
func (this_ *ProcessCommonJSModules) IsJavaNil_() bool { return this_ == nil }
func (*ProcessCommonJSModules) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.ProcessCommonJSModules", "ProcessCommonJSModules"
}
func (this_ *ProcessCommonJSModules) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.NodeTraversal_AbstractPreOrderCallbackBase.IdentityHash_()
}
func (*ProcessCommonJSModules) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.ProcessCommonJSModules", "com.google.javascript.jscomp.ProcessCommonJSModules", func(v_ any) bool {
		return j2gort.InstanceOf[*ProcessCommonJSModules](v_)
	})
}
func (*ProcessCommonJSModules) Implements_jscomp_NodeTraversal_AbstractPreOrderCallback() {}
func (*ProcessCommonJSModules) Implements_jscomp_CompilerPass()                           {}
func (*ProcessCommonJSModules) Implements_jscomp_NodeTraversal_Callback()                 {}

func StaticGet_ProcessCommonJSModules_EXPORTS() string         { return "exports" }
func StaticGet_ProcessCommonJSModules_MODULE() string          { return "module" }
func StaticGet_ProcessCommonJSModules_REQUIRE() string         { return "require" }
func StaticGet_ProcessCommonJSModules_WEBPACK_REQUIRE() string { return "__webpack_require__" }
func StaticGet_ProcessCommonJSModules_WEBPACK_REQUIRE_NAMESPACE() string {
	return "__webpack_require__.t"
}
func StaticGet_ProcessCommonJSModules_EXPORT_PROPERTY_NAME() string { return "default" }

var static_ProcessCommonJSModules_UNKNOWN_REQUIRE_ENSURE *DiagnosticType

func StaticGet_ProcessCommonJSModules_UNKNOWN_REQUIRE_ENSURE() *DiagnosticType {
	clinit_ProcessCommonJSModules()
	return static_ProcessCommonJSModules_UNKNOWN_REQUIRE_ENSURE
}
func StaticSet_ProcessCommonJSModules_UNKNOWN_REQUIRE_ENSURE(v *DiagnosticType) {
	clinit_ProcessCommonJSModules()
	static_ProcessCommonJSModules_UNKNOWN_REQUIRE_ENSURE = v
}

var static_ProcessCommonJSModules_SUSPICIOUS_EXPORTS_ASSIGNMENT *DiagnosticType

func StaticGet_ProcessCommonJSModules_SUSPICIOUS_EXPORTS_ASSIGNMENT() *DiagnosticType {
	clinit_ProcessCommonJSModules()
	return static_ProcessCommonJSModules_SUSPICIOUS_EXPORTS_ASSIGNMENT
}
func StaticSet_ProcessCommonJSModules_SUSPICIOUS_EXPORTS_ASSIGNMENT(v *DiagnosticType) {
	clinit_ProcessCommonJSModules()
	static_ProcessCommonJSModules_SUSPICIOUS_EXPORTS_ASSIGNMENT = v
}

var static_ProcessCommonJSModules_PROMISE_ALL pkg_rhino.QualifiedName

func StaticGet_ProcessCommonJSModules_PROMISE_ALL() pkg_rhino.QualifiedName {
	clinit_ProcessCommonJSModules()
	return static_ProcessCommonJSModules_PROMISE_ALL
}
func StaticSet_ProcessCommonJSModules_PROMISE_ALL(v pkg_rhino.QualifiedName) {
	clinit_ProcessCommonJSModules()
	static_ProcessCommonJSModules_PROMISE_ALL = v
}

var static_ProcessCommonJSModules_DEFINE_AMD pkg_rhino.QualifiedName

func StaticGet_ProcessCommonJSModules_DEFINE_AMD() pkg_rhino.QualifiedName {
	clinit_ProcessCommonJSModules()
	return static_ProcessCommonJSModules_DEFINE_AMD
}
func StaticSet_ProcessCommonJSModules_DEFINE_AMD(v pkg_rhino.QualifiedName) {
	clinit_ProcessCommonJSModules()
	static_ProcessCommonJSModules_DEFINE_AMD = v
}

var static_ProcessCommonJSModules_WINDOW_DEFINE pkg_rhino.QualifiedName

func StaticGet_ProcessCommonJSModules_WINDOW_DEFINE() pkg_rhino.QualifiedName {
	clinit_ProcessCommonJSModules()
	return static_ProcessCommonJSModules_WINDOW_DEFINE
}
func StaticSet_ProcessCommonJSModules_WINDOW_DEFINE(v pkg_rhino.QualifiedName) {
	clinit_ProcessCommonJSModules()
	static_ProcessCommonJSModules_WINDOW_DEFINE = v
}

var static_ProcessCommonJSModules_WINDOW_DEFINE_AMD pkg_rhino.QualifiedName

func StaticGet_ProcessCommonJSModules_WINDOW_DEFINE_AMD() pkg_rhino.QualifiedName {
	clinit_ProcessCommonJSModules()
	return static_ProcessCommonJSModules_WINDOW_DEFINE_AMD
}
func StaticSet_ProcessCommonJSModules_WINDOW_DEFINE_AMD(v pkg_rhino.QualifiedName) {
	clinit_ProcessCommonJSModules()
	static_ProcessCommonJSModules_WINDOW_DEFINE_AMD = v
}

var static_ProcessCommonJSModules_REQUIRE_ENSURE pkg_rhino.QualifiedName

func StaticGet_ProcessCommonJSModules_REQUIRE_ENSURE() pkg_rhino.QualifiedName {
	clinit_ProcessCommonJSModules()
	return static_ProcessCommonJSModules_REQUIRE_ENSURE
}
func StaticSet_ProcessCommonJSModules_REQUIRE_ENSURE(v pkg_rhino.QualifiedName) {
	clinit_ProcessCommonJSModules()
	static_ProcessCommonJSModules_REQUIRE_ENSURE = v
}

var static_ProcessCommonJSModules_GOOG_PROVIDE pkg_rhino.QualifiedName

func StaticGet_ProcessCommonJSModules_GOOG_PROVIDE() pkg_rhino.QualifiedName {
	clinit_ProcessCommonJSModules()
	return static_ProcessCommonJSModules_GOOG_PROVIDE
}
func StaticSet_ProcessCommonJSModules_GOOG_PROVIDE(v pkg_rhino.QualifiedName) {
	clinit_ProcessCommonJSModules()
	static_ProcessCommonJSModules_GOOG_PROVIDE = v
}

var static_ProcessCommonJSModules_GOOG_MODULE pkg_rhino.QualifiedName

func StaticGet_ProcessCommonJSModules_GOOG_MODULE() pkg_rhino.QualifiedName {
	clinit_ProcessCommonJSModules()
	return static_ProcessCommonJSModules_GOOG_MODULE
}
func StaticSet_ProcessCommonJSModules_GOOG_MODULE(v pkg_rhino.QualifiedName) {
	clinit_ProcessCommonJSModules()
	static_ProcessCommonJSModules_GOOG_MODULE = v
}

var static_ProcessCommonJSModules_MODULE_EXPORTS pkg_rhino.QualifiedName

func StaticGet_ProcessCommonJSModules_MODULE_EXPORTS() pkg_rhino.QualifiedName {
	clinit_ProcessCommonJSModules()
	return static_ProcessCommonJSModules_MODULE_EXPORTS
}
func StaticSet_ProcessCommonJSModules_MODULE_EXPORTS(v pkg_rhino.QualifiedName) {
	clinit_ProcessCommonJSModules()
	static_ProcessCommonJSModules_MODULE_EXPORTS = v
}

var clinitDone_ProcessCommonJSModules uint32
var clinitState_ProcessCommonJSModules j2gort.ClinitState

func clinit_ProcessCommonJSModules() {
	if atomic.LoadUint32(&clinitDone_ProcessCommonJSModules) != 0 {
		return
	}
	clinitSlow_ProcessCommonJSModules()
}

func clinitSlow_ProcessCommonJSModules() {
	if !j2gort.EnterClinit(&clinitState_ProcessCommonJSModules, &clinitDone_ProcessCommonJSModules) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_ProcessCommonJSModules, &clinitDone_ProcessCommonJSModules)
	static_ProcessCommonJSModules_UNKNOWN_REQUIRE_ENSURE = DiagnosticType_Warning("JSC_COMMONJS_UNKNOWN_REQUIRE_ENSURE_ERROR", "Unrecognized require.ensure call: {0}")
	static_ProcessCommonJSModules_SUSPICIOUS_EXPORTS_ASSIGNMENT = DiagnosticType_Warning("JSC_COMMONJS_SUSPICIOUS_EXPORTS_ASSIGNMENT", "Suspicious re-assignment of \"exports\" variable. Did you actually intend to export something?")
	static_ProcessCommonJSModules_PROMISE_ALL = pkg_rhino.QualifiedName_Of("Promise.all")
	static_ProcessCommonJSModules_DEFINE_AMD = pkg_rhino.QualifiedName_Of("define.amd")
	static_ProcessCommonJSModules_WINDOW_DEFINE = pkg_rhino.QualifiedName_Of("window.define")
	static_ProcessCommonJSModules_WINDOW_DEFINE_AMD = pkg_rhino.QualifiedName_Of("window.define.amd")
	static_ProcessCommonJSModules_REQUIRE_ENSURE = pkg_rhino.QualifiedName_Of("require.ensure")
	static_ProcessCommonJSModules_GOOG_PROVIDE = pkg_rhino.QualifiedName_Of("goog.provide")
	static_ProcessCommonJSModules_GOOG_MODULE = pkg_rhino.QualifiedName_Of("goog.module")
	static_ProcessCommonJSModules_MODULE_EXPORTS = pkg_rhino.QualifiedName_Of("module.exports")
}

func (this_ *ProcessCommonJSModules) initFields_ProcessCommonJSModules() {
	_ = this_
}

// Java constructor: ProcessCommonJSModules(com.google.javascript.jscomp.AbstractCompiler compiler)
func NewProcessCommonJSModules(compiler AbstractCompiler) *ProcessCommonJSModules {
	clinit_ProcessCommonJSModules()
	this_ := &ProcessCommonJSModules{}
	InitProcessCommonJSModules(this_, compiler)
	return this_
}

func InitProcessCommonJSModules(this_ *ProcessCommonJSModules, compiler AbstractCompiler) {
	_ = compiler
	InitNodeTraversal_AbstractPreOrderCallback(&this_.NodeTraversal_AbstractPreOrderCallbackBase)
	this_.initFields_ProcessCommonJSModules()
	this_.compiler = compiler
}

// Java method: ProcessCommonJSModules.process(com.google.javascript.rhino.Node externs, com.google.javascript.rhino.Node root)
func (this_ *ProcessCommonJSModules) Process_2(externs pkg_rhino.Node, root pkg_rhino.Node) {
	_ = externs
	_ = root
	NodeTraversal_Traverse(this_.compiler, root, this_)
	NodeTraversal_Traverse(this_.compiler, externs, this_)
}

// Java method: ProcessCommonJSModules.shouldTraverse(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *ProcessCommonJSModules) ShouldTraverse_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) bool {
	_ = t
	_ = n
	_ = parent
	if n.IsRoot() {
		return true
	} else if n.IsScript() {
		if this_.compiler.GetOptions().GetModuleResolutionMode() == StaticGet_ModuleLoader_ResolutionMode_WEBPACK() {
			this_.priv_ProcessCommonJSModules_removeWebpackModuleShim_4519(n)
		}
		var val_1 ProcessCommonJSModules_FindImportsAndExports
		val_1.outer_ = this_
		InitProcessCommonJSModules_FindImportsAndExports(&val_1)
		var finder *ProcessCommonJSModules_FindImportsAndExports = &val_1
		_ = finder
		var moduleLoaderErrorHandler ErrorHandler = this_.compiler.GetModuleLoader().GetErrorHandler()
		_ = moduleLoaderErrorHandler
		this_.compiler.GetModuleLoader().SetErrorHandler(ToIface_ErrorHandler(finder))
		NodeTraversal_Traverse(this_.compiler, n, ToIface_NodeTraversal_Callback(finder))
		var moduleType CompilerInput_ModuleType = this_.compiler.GetInput(n.GetInputId()).GetJsModuleType()
		_ = moduleType
		var forceModuleDetection bool = (moduleType == StaticGet_CompilerInput_ModuleType_IMPORTED_SCRIPT())
		_ = forceModuleDetection
		var defaultExportIsConst bool = true
		_ = defaultExportIsConst
		var isCommonJsModule bool = finder.IsCommonJsModule()
		_ = isCommonJsModule
		var exports *j2gort.Ext_com_google_common_collect_ImmutableList_Builder = (&j2gort.Ext_com_google_common_collect_ImmutableList_Builder{})
		_ = exports
		if isCommonJsModule || forceModuleDetection {
			if !(finder.umdPatterns.IsEmpty()) {
				var needsRetraverse bool = false
				_ = needsRetraverse
				if finder.ReplaceUmdPatterns() {
					needsRetraverse = true
				}
				if this_.priv_ProcessCommonJSModules_removeIIFEWrapper_4518(n) {
					needsRetraverse = true
				}
				if needsRetraverse {
					finder = NewProcessCommonJSModules_FindImportsAndExports(this_)
					this_.compiler.GetModuleLoader().SetErrorHandler(ToIface_ErrorHandler(finder))
					NodeTraversal_Traverse(this_.compiler, n, ToIface_NodeTraversal_Callback(finder))
				}
			}
			defaultExportIsConst = finder.InitializeModule()
			exports.AddAllVoid_(j2gort.AsIface[j2gort.Ext_java_lang_Iterable](finder.GetModuleExports()))
			exports.AddAllVoid_(j2gort.AsIface[j2gort.Ext_java_lang_Iterable](finder.GetExports()))
		}
		finder.ReportModuleErrors()
		this_.compiler.GetModuleLoader().SetErrorHandler(moduleLoaderErrorHandler)
		NodeTraversal_Traverse(this_.compiler, n, NewProcessCommonJSModules_RewriteModule(this_, (isCommonJsModule || forceModuleDetection), j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableCollection](exports.BuildImmList_()), defaultExportIsConst))
	}
	return false
}

// Java static method: ProcessCommonJSModules.getModuleName(com.google.javascript.jscomp.CompilerInput input)
func ProcessCommonJSModules_GetModuleName_1_CompilerInput(input *CompilerInput) string {
	clinit_ProcessCommonJSModules()
	_ = input
	var modulePath *ModuleLoader_ModulePath = input.GetPath()
	_ = modulePath
	if modulePath == nil {
		return j2gort.NullString
	}
	return ProcessCommonJSModules_GetModuleName_1_ModuleLoader_ModulePath(modulePath)
}

// Java static method: ProcessCommonJSModules.getModuleName(com.google.javascript.jscomp.deps.ModuleLoader.ModulePath input)
func ProcessCommonJSModules_GetModuleName_1_ModuleLoader_ModulePath(input *ModuleLoader_ModulePath) string {
	clinit_ProcessCommonJSModules()
	_ = input
	return input.ToModuleName()
}

// Java method: ProcessCommonJSModules.getBasePropertyImport(java.lang.String moduleName, com.google.javascript.rhino.Node requireCall)
func (this_ *ProcessCommonJSModules) GetBasePropertyImport_2(moduleName string, requireCall pkg_rhino.Node) string {
	_ = moduleName
	_ = requireCall
	if !(requireCall.IsCall()) {
		j2gort.CheckArgument(false)
	}
	if (this_.compiler.GetOptions().GetModuleResolutionMode() == StaticGet_ModuleLoader_ResolutionMode_WEBPACK()) && requireCall.GetFirstChild().MatchesQualifiedName_1_string("__webpack_require__.t") {
		return moduleName
	}
	return this_.GetBasePropertyImport_1(moduleName)
}

// Java method: ProcessCommonJSModules.getBasePropertyImport(java.lang.String moduleName)
func (this_ *ProcessCommonJSModules) GetBasePropertyImport_1(moduleName string) string {
	_ = moduleName
	var moduleType CompilerInput_ModuleType = this_.compiler.GetModuleTypeByName(moduleName)
	_ = moduleType
	if (moduleType != 0) && (moduleType != StaticGet_CompilerInput_ModuleType_COMMONJS()) {
		return moduleName
	}
	return (j2gort.StrValString(moduleName) + "." + "default")
}

// Java method: ProcessCommonJSModules.isCommonJsImport(com.google.javascript.rhino.Node requireCall)
func (this_ *ProcessCommonJSModules) IsCommonJsImport(requireCall pkg_rhino.Node) bool {
	_ = requireCall
	return ProcessCommonJSModules_IsCommonJsImport(requireCall, this_.compiler.GetOptions().GetModuleResolutionMode())
}

// Java static method: ProcessCommonJSModules.isCommonJsImport(com.google.javascript.rhino.Node requireCall, com.google.javascript.jscomp.deps.ModuleLoader.ResolutionMode resolutionMode)
func ProcessCommonJSModules_IsCommonJsImport(requireCall pkg_rhino.Node, resolutionMode ModuleLoader_ResolutionMode) bool {
	clinit_ProcessCommonJSModules()
	_ = requireCall
	_ = resolutionMode
	if requireCall.IsCall() && requireCall.HasTwoChildren() {
		if ((resolutionMode == StaticGet_ModuleLoader_ResolutionMode_WEBPACK()) && (requireCall.GetFirstChild().MatchesQualifiedName_1_string("__webpack_require__") || requireCall.GetFirstChild().MatchesQualifiedName_1_string("__webpack_require__.t"))) && (requireCall.GetSecondChild().IsNumber() || requireCall.GetSecondChild().IsStringLit()) {
			return true
		} else if requireCall.GetFirstChild().MatchesQualifiedName_1_string("require") && requireCall.GetSecondChild().IsStringLit() {
			return true
		}
	} else if ((((requireCall.IsCall() && requireCall.HasXChildren(int32(3))) && (resolutionMode == StaticGet_ModuleLoader_ResolutionMode_WEBPACK())) && requireCall.GetFirstChild().MatchesQualifiedName_1_string(("__webpack_require__" + ".bind"))) && requireCall.GetSecondChild().IsNull()) && (requireCall.GetLastChild().IsNumber() || requireCall.GetLastChild().IsStringLit()) {
		return true
	}
	return false
}

// Java method: ProcessCommonJSModules.getCommonJsImportPath(com.google.javascript.rhino.Node requireCall)
func (this_ *ProcessCommonJSModules) GetCommonJsImportPath(requireCall pkg_rhino.Node) string {
	_ = requireCall
	return ProcessCommonJSModules_GetCommonJsImportPath(requireCall, this_.compiler.GetOptions().GetModuleResolutionMode())
}

// Java static method: ProcessCommonJSModules.getCommonJsImportPath(com.google.javascript.rhino.Node requireCall, com.google.javascript.jscomp.deps.ModuleLoader.ResolutionMode resolutionMode)
func ProcessCommonJSModules_GetCommonJsImportPath(requireCall pkg_rhino.Node, resolutionMode ModuleLoader_ResolutionMode) string {
	clinit_ProcessCommonJSModules()
	_ = requireCall
	_ = resolutionMode
	if resolutionMode == StaticGet_ModuleLoader_ResolutionMode_WEBPACK() {
		var pathArgument pkg_rhino.Node
		_ = pathArgument
		if requireCall.GetChildCount() >= int32(3) {
			pathArgument = requireCall.GetChildAtIndex(int32(2))
		} else {
			pathArgument = requireCall.GetSecondChild()
		}
		if pathArgument.IsNumber() {
			return strconv.FormatInt(int64(javalang.D2I(pathArgument.GetDouble())), 10)
		} else {
			return pathArgument.GetString()
		}
	}
	return requireCall.GetSecondChild().GetString()
}

// Java method: ProcessCommonJSModules.getImportedModuleName(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node requireCall)
func (this_ *ProcessCommonJSModules) priv_ProcessCommonJSModules_getImportedModuleName_4514(t *NodeTraversal, requireCall pkg_rhino.Node) string {
	_ = t
	_ = requireCall
	return this_.priv_ProcessCommonJSModules_getImportedModuleName_4515(t, requireCall, this_.GetCommonJsImportPath(requireCall))
}

// Java method: ProcessCommonJSModules.getImportedModuleName(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, java.lang.String importPath)
func (this_ *ProcessCommonJSModules) priv_ProcessCommonJSModules_getImportedModuleName_4515(t *NodeTraversal, n pkg_rhino.Node, importPath string) string {
	_ = t
	_ = n
	_ = importPath
	var modulePath *ModuleLoader_ModulePath = t.GetInput().GetPath().ResolveJsModule(importPath, n.GetSourceFileName(), n.GetLineno(), n.GetCharno())
	_ = modulePath
	if modulePath == nil {
		return ModuleIdentifier_ForFile(importPath).ModuleName()
	}
	return modulePath.ToModuleName()
}

// Java static method: ProcessCommonJSModules.isCommonJsExport(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node export, com.google.javascript.jscomp.deps.ModuleLoader.ResolutionMode resolutionMode)
func ProcessCommonJSModules_IsCommonJsExport(t *NodeTraversal, export pkg_rhino.Node, resolutionMode ModuleLoader_ResolutionMode) bool {
	clinit_ProcessCommonJSModules()
	_ = t
	_ = export
	_ = resolutionMode
	if export.MatchesQualifiedName_1_string(("module" + "." + "exports")) || (((export.IsGetElem() && export.GetFirstChild().MatchesQualifiedName_1_string("module")) && export.GetSecondChild().IsStringLit()) && func() bool {
		sr_2 := export.GetSecondChild().GetString()
		se_3 := "exports"
		return j2gort.CheckNotNullString(sr_2) == se_3
	}()) {
		var v *Var = j2gort.CastPtr[Var](t.GetScope().GetVar("module"))
		_ = v
		if (v == nil) || v.IsExtern() {
			return true
		}
	} else if export.IsName() && func() bool {
		sr_4 := "exports"
		se_5 := export.GetString()
		return j2gort.CheckNotNullString(sr_4) == se_5
	}() {
		var v_2 *Var = j2gort.CastPtr[Var](t.GetScope().GetVar(export.GetString()))
		_ = v_2
		if (v_2 == nil) || v_2.IsGlobal() {
			return true
		}
	}
	return false
}

// Java method: ProcessCommonJSModules.isCommonJsExport(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node export)
func (this_ *ProcessCommonJSModules) priv_ProcessCommonJSModules_isCommonJsExport_4516(t *NodeTraversal, export pkg_rhino.Node) bool {
	_ = t
	_ = export
	return ProcessCommonJSModules_IsCommonJsExport(t, export, this_.compiler.GetOptions().GetModuleResolutionMode())
}

// Java static method: ProcessCommonJSModules.isCommonJsDynamicImportCallback(com.google.javascript.rhino.Node n, com.google.javascript.jscomp.deps.ModuleLoader.ResolutionMode resolutionMode)
func ProcessCommonJSModules_IsCommonJsDynamicImportCallback(n pkg_rhino.Node, resolutionMode ModuleLoader_ResolutionMode) bool {
	clinit_ProcessCommonJSModules()
	_ = n
	_ = resolutionMode
	if (n == 0) || (resolutionMode != StaticGet_ModuleLoader_ResolutionMode_WEBPACK()) {
		return false
	}
	if n.IsFunction() && ProcessCommonJSModules_IsWebpackRequireEnsureCallback(n) {
		return true
	}
	return false
}

// Java static method: ProcessCommonJSModules.isWebpackRequireEnsureCallback(com.google.javascript.rhino.Node fnc)
func ProcessCommonJSModules_IsWebpackRequireEnsureCallback(fnc pkg_rhino.Node) bool {
	clinit_ProcessCommonJSModules()
	_ = fnc
	if !(fnc.IsFunction()) {
		j2gort.CheckArgument(false)
	}
	if fnc.GetParent() == 0 {
		return false
	}
	var callParent pkg_rhino.Node = fnc.GetParent()
	_ = callParent
	if !(callParent.IsCall()) {
		return false
	}
	if !((callParent.HasChildren() && callParent.GetFirstChild().IsGetProp()) && callParent.GetFirstFirstChild().IsCall()) {
		return false
	}
	var callParentTarget pkg_rhino.Node = callParent.GetFirstFirstChild().GetFirstChild()
	_ = callParentTarget
	if callParentTarget.MatchesQualifiedName_1_string(("__webpack_require__" + ".e")) && (j2gort.CheckNotNullString(callParent.GetFirstChild().GetString()) == "then") {
		return true
	} else if (StaticGet_ProcessCommonJSModules_PROMISE_ALL().Matches(callParentTarget) && (callParentTarget.GetNext() != 0)) && callParentTarget.GetNext().IsArrayLit() {
		var allElementsAreDynamicImports bool = false
		_ = allElementsAreDynamicImports
		{
			var arrayItem pkg_rhino.Node = callParentTarget.GetNext().GetFirstChild()
			_ = arrayItem
			for ; arrayItem != 0; arrayItem = arrayItem.GetNext() {
				if !((arrayItem.IsCall() && arrayItem.HasTwoChildren()) && arrayItem.GetFirstChild().MatchesQualifiedName_1_string(("__webpack_require__" + ".e"))) {
					return false
				}
				allElementsAreDynamicImports = true
			}
		}
		return allElementsAreDynamicImports
	}
	return false
}

// Java method: ProcessCommonJSModules.getBaseQualifiedNameNode(com.google.javascript.rhino.Node n)
func (this_ *ProcessCommonJSModules) priv_ProcessCommonJSModules_getBaseQualifiedNameNode_4517(n pkg_rhino.Node) pkg_rhino.Node {
	_ = n
	var refParent pkg_rhino.Node = n
	_ = refParent
	for refParent.HasParent() && refParent.GetParent().IsQualifiedName() {
		refParent = refParent.GetParent()
	}
	return refParent
}

// Java method: ProcessCommonJSModules.removeIIFEWrapper(com.google.javascript.rhino.Node root)
func (this_ *ProcessCommonJSModules) priv_ProcessCommonJSModules_removeIIFEWrapper_4518(root pkg_rhino.Node) bool {
	_ = root
	if !(root.IsScript()) {
		j2gort.CheckState(false)
	}
	var n pkg_rhino.Node = root.GetFirstChild()
	_ = n
	for (n != 0) && n.IsEmpty() {
		n = n.GetNext()
	}
	if ((n == 0) || !(n.IsExprResult())) || (n.GetNext() != 0) {
		return false
	}
	if ((n != 0) && n.HasChildren()) && n.GetFirstChild().IsNot() {
		n = n.GetFirstChild()
	}
	var call pkg_rhino.Node = n.GetFirstChild()
	_ = call
	if (call == 0) || !(call.IsCall()) {
		return false
	}
	var fnc pkg_rhino.Node = pkg_rhino.Node(0)
	_ = fnc
	if call.GetFirstChild().IsFunction() {
		fnc = n.GetFirstFirstChild()
	} else if (call.GetFirstChild().IsGetProp() && call.GetFirstFirstChild().IsFunction()) && (j2gort.CheckNotNullString(call.GetFirstChild().GetString()) == "call") {
		fnc = call.GetFirstFirstChild()
		if !((call.GetSecondChild() != 0) && (call.GetSecondChild().IsThis() || call.GetSecondChild().MatchesQualifiedName_1_string("exports"))) {
			return false
		}
	} else {
		return false
	}
	if NodeUtil_DoesFunctionReferenceOwnArgumentsObject(fnc) {
		return false
	}
	var ci *CompilerInput = this_.compiler.GetInput(root.GetInputId())
	_ = ci
	var modulePath *ModuleLoader_ModulePath = ci.GetPath()
	_ = modulePath
	if modulePath == nil {
		return false
	}
	var iifeLabel string = (j2gort.StrValString(ProcessCommonJSModules_GetModuleName_1_ModuleLoader_ModulePath(modulePath)) + "_iifeWrapper")
	_ = iifeLabel
	var val_6 FunctionToBlockMutator
	InitFunctionToBlockMutator_2(&val_6, this_.compiler, j2gort.AsIface[j2gort.Ext_java_util_function_Supplier](this_.compiler.GetUniqueNameIdSupplier()))
	var mutator *FunctionToBlockMutator = &val_6
	_ = mutator
	var block pkg_rhino.Node = mutator.MutateWithoutRenaming(iifeLabel, fnc, call, j2gort.NullString, false, false)
	_ = block
	root.RemoveChildren()
	root.AddChildrenToFront(block.RemoveChildren())
	this_.priv_ProcessCommonJSModules_reportNestedScopesDeleted_4520(fnc)
	this_.compiler.ReportChangeToEnclosingScope(root)
	return true
}

// Java method: ProcessCommonJSModules.removeWebpackModuleShim(com.google.javascript.rhino.Node root)
func (this_ *ProcessCommonJSModules) priv_ProcessCommonJSModules_removeWebpackModuleShim_4519(root pkg_rhino.Node) {
	_ = root
	if !(root.IsScript()) {
		j2gort.CheckState(false)
	}
	var n pkg_rhino.Node = root.GetFirstChild()
	_ = n
	for (n != 0) && n.IsEmpty() {
		n = n.GetNext()
	}
	if ((n == 0) || !(n.IsExprResult())) || (n.GetNext() != 0) {
		return
	}
	var call pkg_rhino.Node = n.GetFirstChild()
	_ = call
	if (call == 0) || !(call.IsCall()) {
		return
	}
	var callTarget pkg_rhino.Node = call.GetFirstChild()
	_ = callTarget
	if !(callTarget.IsFunction()) {
		return
	}
	var fnc pkg_rhino.Node = callTarget
	_ = fnc
	var params pkg_rhino.Node = NodeUtil_GetFunctionParameters(fnc)
	_ = params
	var moduleParam pkg_rhino.Node = pkg_rhino.Node(0)
	_ = moduleParam
	var param pkg_rhino.Node = params.GetFirstChild()
	_ = param
	var paramNumber int32 = int32(0)
	_ = paramNumber
	for param != 0 {
		paramNumber = paramNumber + 1
		if param.IsName() && func() bool {
			sr_7 := param.GetString()
			se_8 := "module"
			return j2gort.CheckNotNullString(sr_7) == se_8
		}() {
			moduleParam = param
			break
		}
		param = param.GetNext()
	}
	if moduleParam == 0 {
		return
	}
	var isFreeCall bool = call.GetBooleanProp(pkg_rhino.StaticGet_Node_FREE_CALL())
	_ = isFreeCall
	var arg pkg_rhino.Node = call.GetChildAtIndex(func() int32 {
		if isFreeCall {
			return paramNumber
		}
		return (paramNumber + int32(1))
	}())
	_ = arg
	if arg == 0 {
		return
	}
	var argCallTarget pkg_rhino.Node = arg.GetFirstChild()
	_ = argCallTarget
	if ((arg.IsCall() && argCallTarget.IsCall()) && this_.IsCommonJsImport(argCallTarget)) && argCallTarget.GetNext().MatchesName_1_string("module") {
		var importPath string = this_.GetCommonJsImportPath(argCallTarget)
		_ = importPath
		var modulePath *ModuleLoader_ModulePath = this_.compiler.GetInput(root.GetInputId()).GetPath().ResolveJsModule(importPath, arg.GetSourceFileName(), arg.GetLineno(), arg.GetCharno())
		_ = modulePath
		if modulePath == nil {
			return
		}
		if j2gort.JStrIndexOf(modulePath.String(), "/buildin/module.js") >= 0 {
			arg.Detach()
			param.Detach()
			this_.compiler.ReportChangeToChangeScope(fnc)
			this_.compiler.ReportChangeToEnclosingScope(fnc)
		}
	}
}

// Java method: ProcessCommonJSModules.reportNestedScopesDeleted(com.google.javascript.rhino.Node n)
func (this_ *ProcessCommonJSModules) priv_ProcessCommonJSModules_reportNestedScopesDeleted_4520(n pkg_rhino.Node) {
	_ = n
	NodeUtil_VisitPreOrder_2(n, NewProcessCommonJSModules__Anon_74(this_))
}

// Java method: ProcessCommonJSModules.reportNestedScopesChanged(com.google.javascript.rhino.Node n)
func (this_ *ProcessCommonJSModules) priv_ProcessCommonJSModules_reportNestedScopesChanged_4521(n pkg_rhino.Node) {
	_ = n
	NodeUtil_VisitPreOrder_2(n, NewProcessCommonJSModules__Anon_75(this_))
}

// Java static method: ProcessCommonJSModules.findUmdPattern(java.util.List<com.google.javascript.jscomp.ProcessCommonJSModules.UmdPattern> umdPatterns, com.google.javascript.rhino.Node n)
func ProcessCommonJSModules_FindUmdPattern(umdPatterns j2gort.Ext_java_util_List, n pkg_rhino.Node) *ProcessCommonJSModules_UmdPattern {
	clinit_ProcessCommonJSModules()
	_ = umdPatterns
	_ = n
	{
		cur_9 := j2gort.NewIterCursor(j2gort.CheckNotNull(umdPatterns))
		for cur_9.HasNext() {
			rawElem_10 := cur_9.Next()
			umd := j2gort.CastPtr[ProcessCommonJSModules_UmdPattern](rawElem_10)
			_ = umd
			if umd.ifRoot == n {
				return umd
			}
		}
	}
	return nil
}

// Java class: com.google.javascript.jscomp.ProcessCommonJSModules.UmdPattern
type ProcessCommonJSModules_UmdPattern struct {
	ifRoot       pkg_rhino.Node `gson:"ifRoot"`
	activeBranch pkg_rhino.Node `gson:"activeBranch"`
	idHash_      atomic.Int32
}

func (this_ *ProcessCommonJSModules_UmdPattern) BindSelf_(s any) {
	_ = s
}
func (this_ *ProcessCommonJSModules_UmdPattern) IsJavaNil_() bool { return this_ == nil }
func (*ProcessCommonJSModules_UmdPattern) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.ProcessCommonJSModules$UmdPattern", "UmdPattern"
}
func (this_ *ProcessCommonJSModules_UmdPattern) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*ProcessCommonJSModules_UmdPattern) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.ProcessCommonJSModules.UmdPattern", "com.google.javascript.jscomp.ProcessCommonJSModules$UmdPattern", func(v_ any) bool {
		return j2gort.InstanceOf[*ProcessCommonJSModules_UmdPattern](v_)
	})
}

func (this_ *ProcessCommonJSModules_UmdPattern) initFields_ProcessCommonJSModules_UmdPattern() {
	_ = this_
}

// Java constructor: UmdPattern(com.google.javascript.rhino.Node ifRoot, com.google.javascript.rhino.@org.jspecify.annotations.Nullable Node activeBranch)
func NewProcessCommonJSModules_UmdPattern(ifRoot pkg_rhino.Node, activeBranch pkg_rhino.Node) *ProcessCommonJSModules_UmdPattern {
	this_ := &ProcessCommonJSModules_UmdPattern{}
	InitProcessCommonJSModules_UmdPattern(this_, ifRoot, activeBranch)
	return this_
}

func InitProcessCommonJSModules_UmdPattern(this_ *ProcessCommonJSModules_UmdPattern, ifRoot pkg_rhino.Node, activeBranch pkg_rhino.Node) {
	_ = ifRoot
	_ = activeBranch
	this_.initFields_ProcessCommonJSModules_UmdPattern()
	this_.ifRoot = ifRoot
	this_.activeBranch = activeBranch
}

// Java class: com.google.javascript.jscomp.ProcessCommonJSModules.ExportInfo
type ProcessCommonJSModules_ExportInfo struct {
	scope              *Scope         `gson:"scope"`
	node               pkg_rhino.Node `gson:"node"`
	idHash_            atomic.Int32
	isInSupportedScope bool `gson:"isInSupportedScope"`
}

func (this_ *ProcessCommonJSModules_ExportInfo) BindSelf_(s any) {
	_ = s
}
func (this_ *ProcessCommonJSModules_ExportInfo) IsJavaNil_() bool { return this_ == nil }
func (*ProcessCommonJSModules_ExportInfo) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.ProcessCommonJSModules$ExportInfo", "ExportInfo"
}
func (this_ *ProcessCommonJSModules_ExportInfo) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*ProcessCommonJSModules_ExportInfo) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.ProcessCommonJSModules.ExportInfo", "com.google.javascript.jscomp.ProcessCommonJSModules$ExportInfo", func(v_ any) bool {
		return j2gort.InstanceOf[*ProcessCommonJSModules_ExportInfo](v_)
	})
}

func (this_ *ProcessCommonJSModules_ExportInfo) initFields_ProcessCommonJSModules_ExportInfo() {
	_ = this_
}

// Java constructor: ExportInfo(com.google.javascript.rhino.Node node, com.google.javascript.jscomp.Scope scope)
func NewProcessCommonJSModules_ExportInfo(node pkg_rhino.Node, scope *Scope) *ProcessCommonJSModules_ExportInfo {
	this_ := &ProcessCommonJSModules_ExportInfo{}
	InitProcessCommonJSModules_ExportInfo(this_, node, scope)
	return this_
}

func InitProcessCommonJSModules_ExportInfo(this_ *ProcessCommonJSModules_ExportInfo, node pkg_rhino.Node, scope *Scope) {
	_ = node
	_ = scope
	this_.initFields_ProcessCommonJSModules_ExportInfo()
	this_.node = node
	this_.scope = scope
	var disqualifyingParent pkg_rhino.Node = NodeUtil_GetEnclosingNode(node, lambdaSingleton_process_common_js_modules_1)
	_ = disqualifyingParent
	this_.isInSupportedScope = (disqualifyingParent == 0)
}

// Java class: com.google.javascript.jscomp.ProcessCommonJSModules.FindImportsAndExports
type ProcessCommonJSModules_FindImportsAndExports struct {
	outer_                 *ProcessCommonJSModules
	umdPatterns            j2gort.Ext_java_util_List `gson:"umdPatterns"`
	moduleExports          j2gort.Ext_java_util_List `gson:"moduleExports"`
	exports                j2gort.Ext_java_util_List `gson:"exports"`
	errors                 []*JSError                `gson:"errors"`
	script                 pkg_rhino.Node            `gson:"script"`
	idHash_                atomic.Int32
	hasGoogProvideOrModule bool `gson:"hasGoogProvideOrModule"`
}

func (this_ *ProcessCommonJSModules_FindImportsAndExports) BindSelf_(s any) {
	_ = s
}
func (this_ *ProcessCommonJSModules_FindImportsAndExports) IsJavaNil_() bool { return this_ == nil }
func (*ProcessCommonJSModules_FindImportsAndExports) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.ProcessCommonJSModules$FindImportsAndExports", "FindImportsAndExports"
}
func (this_ *ProcessCommonJSModules_FindImportsAndExports) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*ProcessCommonJSModules_FindImportsAndExports) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.ProcessCommonJSModules.FindImportsAndExports", "com.google.javascript.jscomp.ProcessCommonJSModules$FindImportsAndExports", func(v_ any) bool {
		return j2gort.InstanceOf[*ProcessCommonJSModules_FindImportsAndExports](v_)
	})
}
func (this_ *ProcessCommonJSModules_FindImportsAndExports) SetOuter_(o *ProcessCommonJSModules) {
	this_.outer_ = o
}
func (this_ *ProcessCommonJSModules_FindImportsAndExports) GetOuter_() *ProcessCommonJSModules {
	return this_.outer_
}
func (*ProcessCommonJSModules_FindImportsAndExports) Implements_jscomp_NodeTraversal_Callback() {}
func (*ProcessCommonJSModules_FindImportsAndExports) Implements_jscomp_ErrorHandler()           {}

func (this_ *ProcessCommonJSModules_FindImportsAndExports) initFields_ProcessCommonJSModules_FindImportsAndExports() {
	_ = this_

	this_.umdPatterns = j2gort.NewExt_java_util_ArrayList()
	this_.moduleExports = j2gort.NewExt_java_util_ArrayList()
	this_.exports = j2gort.NewExt_java_util_ArrayList()
	this_.errors = []*JSError{}
}

// Java constructor: FindImportsAndExports()
func NewProcessCommonJSModules_FindImportsAndExports(outer_ *ProcessCommonJSModules) *ProcessCommonJSModules_FindImportsAndExports {
	this_ := &ProcessCommonJSModules_FindImportsAndExports{}
	this_.outer_ = outer_
	InitProcessCommonJSModules_FindImportsAndExports(this_)
	return this_
}

func InitProcessCommonJSModules_FindImportsAndExports(this_ *ProcessCommonJSModules_FindImportsAndExports) {
	this_.initFields_ProcessCommonJSModules_FindImportsAndExports()
}

// Java method: FindImportsAndExports.isCommonJsModule()
func (this_ *ProcessCommonJSModules_FindImportsAndExports) IsCommonJsModule() bool {
	return ((!(this_.exports.IsEmpty()) || !(this_.moduleExports.IsEmpty())) && !(this_.hasGoogProvideOrModule))
}

// Java method: FindImportsAndExports.report(com.google.javascript.jscomp.CheckLevel ignoredLevel, com.google.javascript.jscomp.JSError error)
func (this_ *ProcessCommonJSModules_FindImportsAndExports) Report_2(ignoredLevel *CheckLevel, error_ *JSError) {
	_ = ignoredLevel
	_ = error_
	this_.errors = append(this_.errors, error_)
}

// Java method: FindImportsAndExports.getModuleExports()
func (this_ *ProcessCommonJSModules_FindImportsAndExports) GetModuleExports() j2gort.Ext_java_util_List {
	return j2gort.AsIface[j2gort.Ext_java_util_List](j2gort.Ext_com_google_common_collect_ImmutableList_CopyOf_1_Ext_java_util_Collection(j2gort.AsIface[j2gort.Ext_java_util_Collection](this_.moduleExports)))
}

// Java method: FindImportsAndExports.getExports()
func (this_ *ProcessCommonJSModules_FindImportsAndExports) GetExports() j2gort.Ext_java_util_List {
	return j2gort.AsIface[j2gort.Ext_java_util_List](j2gort.Ext_com_google_common_collect_ImmutableList_CopyOf_1_Ext_java_util_Collection(j2gort.AsIface[j2gort.Ext_java_util_Collection](this_.exports)))
}

// Java method: FindImportsAndExports.shouldTraverse(com.google.javascript.jscomp.NodeTraversal nodeTraversal, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *ProcessCommonJSModules_FindImportsAndExports) ShouldTraverse_3(nodeTraversal *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) bool {
	_ = nodeTraversal
	_ = n
	_ = parent
	if n.IsScript() {
		if !(this_.script == 0) {
			j2gort.CheckState(false)
		}
		this_.script = n
	}
	return true
}

// Java method: FindImportsAndExports.visit(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *ProcessCommonJSModules_FindImportsAndExports) Visit_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = n
	_ = parent
	if NodeUtil_IsShallowStatementTree(parent) {
		if n.IsExprResult() {
			var maybeGetProp pkg_rhino.Node = n.GetFirstFirstChild()
			_ = maybeGetProp
			if (maybeGetProp != 0) && (StaticGet_ProcessCommonJSModules_GOOG_PROVIDE().Matches(maybeGetProp) || StaticGet_ProcessCommonJSModules_GOOG_MODULE().Matches(maybeGetProp)) {
				this_.hasGoogProvideOrModule = true
			}
		}
	}
	if n.IsCall() && StaticGet_ProcessCommonJSModules_REQUIRE_ENSURE().Matches(n.GetFirstChild()) {
		this_.priv_FindImportsAndExports_visitRequireEnsureCall_4526(t, n)
	}
	if n.MatchesQualifiedName_1_string(("module" + "." + "exports")) || (((n.IsGetElem() && n.GetFirstChild().MatchesQualifiedName_1_string("module")) && n.GetSecondChild().IsStringLit()) && func() bool {
		sr_12 := n.GetSecondChild().GetString()
		se_13 := "exports"
		return j2gort.CheckNotNullString(sr_12) == se_13
	}()) {
		if this_.outer_.priv_ProcessCommonJSModules_isCommonJsExport_4516(t, n) {
			this_.moduleExports.Add_1(NewProcessCommonJSModules_ExportInfo(n, t.GetScope()))
			var umdTestAncestor *ProcessCommonJSModules_FindImportsAndExports_UmdTestInfo = this_.priv_FindImportsAndExports_getOutermostUmdTest_4529(parent)
			_ = umdTestAncestor
			if (umdTestAncestor != nil) && !(this_.priv_FindImportsAndExports_isInIfTest_4530(n)) {
				var existingPattern *ProcessCommonJSModules_UmdPattern = ProcessCommonJSModules_FindUmdPattern(this_.umdPatterns, umdTestAncestor.enclosingIf)
				_ = existingPattern
				if existingPattern == nil {
					this_.umdPatterns.Add_1(NewProcessCommonJSModules_UmdPattern(umdTestAncestor.enclosingIf, umdTestAncestor.activeBranch))
				}
			}
		}
	} else if StaticGet_ProcessCommonJSModules_DEFINE_AMD().Matches(n) || StaticGet_ProcessCommonJSModules_WINDOW_DEFINE_AMD().Matches(n) {
		var umdTestAncestor_2 *ProcessCommonJSModules_FindImportsAndExports_UmdTestInfo = this_.priv_FindImportsAndExports_getOutermostUmdTest_4529(parent)
		_ = umdTestAncestor_2
		if (umdTestAncestor_2 != nil) && this_.priv_FindImportsAndExports_isInIfTest_4530(n) {
			var existingPattern_2 *ProcessCommonJSModules_UmdPattern = ProcessCommonJSModules_FindUmdPattern(this_.umdPatterns, umdTestAncestor_2.enclosingIf)
			_ = existingPattern_2
			if (existingPattern_2 == nil) && (umdTestAncestor_2.enclosingIf.GetChildAtIndex(int32(2)) == 0) {
				this_.umdPatterns.Add_1(NewProcessCommonJSModules_UmdPattern(umdTestAncestor_2.enclosingIf, pkg_rhino.Node(0)))
			}
		}
	}
	if n.IsName() && func() bool {
		sr_14 := "exports"
		se_15 := n.GetString()
		return j2gort.CheckNotNullString(sr_14) == se_15
	}() {
		var v *Var = j2gort.CastPtr[Var](t.GetScope().GetVar("exports"))
		_ = v
		if (v == nil) || v.IsGlobal() {
			var qNameRoot pkg_rhino.Node = this_.outer_.priv_ProcessCommonJSModules_getBaseQualifiedNameNode_4517(n)
			_ = qNameRoot
			if ((qNameRoot != 0) && qNameRoot.MatchesQualifiedName_1_string("exports")) && NodeUtil_IsLValue(qNameRoot) {
				if (n.GetGrandparent().IsExprResult() && (n.GetNext() != 0)) && ((n.GetNext().IsGetProp() && n.GetNext().MatchesQualifiedName_1_string(("module" + "." + "exports"))) || (n.GetNext().IsAssign() && n.GetNext().GetFirstChild().MatchesQualifiedName_1_string(("module" + "." + "exports")))) {
					this_.exports.Add_1(NewProcessCommonJSModules_ExportInfo(n, t.GetScope()))
				} else if !(this_.hasGoogProvideOrModule) && ((v == nil) || ((v.GetNameNode() == 0) && (v.GetNameNode().GetFirstChild() != n))) {
					this_.errors = append(this_.errors, JSError_Make_3(qNameRoot, StaticGet_ProcessCommonJSModules_SUSPICIOUS_EXPORTS_ASSIGNMENT(), []string{}))
				}
			} else {
				this_.exports.Add_1(NewProcessCommonJSModules_ExportInfo(n, t.GetScope()))
				var umdTestAncestor_3 *ProcessCommonJSModules_FindImportsAndExports_UmdTestInfo = this_.priv_FindImportsAndExports_getOutermostUmdTest_4529(parent)
				_ = umdTestAncestor_3
				if (umdTestAncestor_3 != nil) && !(this_.priv_FindImportsAndExports_isInIfTest_4530(n)) {
					var existingPattern_3 *ProcessCommonJSModules_UmdPattern = ProcessCommonJSModules_FindUmdPattern(this_.umdPatterns, umdTestAncestor_3.enclosingIf)
					_ = existingPattern_3
					if existingPattern_3 == nil {
						this_.umdPatterns.Add_1(NewProcessCommonJSModules_UmdPattern(umdTestAncestor_3.enclosingIf, umdTestAncestor_3.activeBranch))
					}
				}
			}
		}
	} else if (n.IsThis() && n.GetParent().IsGetProp()) && t.InGlobalScope() {
		this_.exports.Add_1(NewProcessCommonJSModules_ExportInfo(n, t.GetScope()))
	}
	if this_.outer_.IsCommonJsImport(n) {
		this_.priv_FindImportsAndExports_visitRequireCall_4525(t, n, parent)
	}
}

// Java method: FindImportsAndExports.visitRequireCall(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node require, com.google.javascript.rhino.Node parent)
func (this_ *ProcessCommonJSModules_FindImportsAndExports) priv_FindImportsAndExports_visitRequireCall_4525(t *NodeTraversal, require pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = require
	_ = parent
	if (!(NodeUtil_IsExpressionResultUsed(require)) && parent.IsExprResult()) && NodeUtil_IsStatementBlock(parent.GetParent()) {
		t.GetInput().GetPath().ResolveJsModule(this_.outer_.GetCommonJsImportPath(require), require.GetSourceFileName(), require.GetLineno(), require.GetCharno())
		var grandparent pkg_rhino.Node = parent.GetParent()
		_ = grandparent
		parent.Detach()
		this_.outer_.compiler.ReportChangeToEnclosingScope(grandparent)
	}
}

// Java method: FindImportsAndExports.visitRequireEnsureCall(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node call)
func (this_ *ProcessCommonJSModules_FindImportsAndExports) priv_FindImportsAndExports_visitRequireEnsureCall_4526(t *NodeTraversal, call pkg_rhino.Node) {
	_ = t
	_ = call
	if !(call.HasXChildren(int32(3))) {
		this_.outer_.compiler.Report_1(JSError_Make_3(call, StaticGet_ProcessCommonJSModules_UNKNOWN_REQUIRE_ENSURE(), []string{"Expected the function to have 2 arguments but instead found {0}", ("" + j2gort.StrValInt32(call.GetChildCount()))}))
		return
	}
	var dependencies pkg_rhino.Node = call.GetSecondChild()
	_ = dependencies
	if !(dependencies.IsArrayLit()) {
		this_.outer_.compiler.Report_1(JSError_Make_3(dependencies, StaticGet_ProcessCommonJSModules_UNKNOWN_REQUIRE_ENSURE(), []string{"The first argument must be an array literal of string literals."}))
		return
	}
	{
		var dep pkg_rhino.Node = dependencies.GetFirstChild()
		_ = dep
		for ; dep != 0; dep = dep.GetNext() {
			if !(dep.IsStringLit()) {
				this_.outer_.compiler.Report_1(JSError_Make_3(dep, StaticGet_ProcessCommonJSModules_UNKNOWN_REQUIRE_ENSURE(), []string{"The first argument must be an array literal of string literals."}))
				return
			}
		}
	}
	var callback pkg_rhino.Node = dependencies.GetNext()
	_ = callback
	if !(((callback.IsFunction() && callback.GetSecondChild().HasOneChild()) && callback.GetSecondChild().GetFirstChild().IsName()) && func() bool {
		sr_16 := "require"
		se_17 := callback.GetSecondChild().GetFirstChild().GetString()
		return j2gort.CheckNotNullString(sr_16) == se_17
	}()) {
		this_.outer_.compiler.Report_1(JSError_Make_3(callback, StaticGet_ProcessCommonJSModules_UNKNOWN_REQUIRE_ENSURE(), []string{"The second argument must be a function whose first argument is named \"require\"."}))
		return
	}
	callback.Detach()
	callback.GetSecondChild().RemoveChildren()
	call.RemoveChildren()
	call.PutBooleanProp(pkg_rhino.StaticGet_Node_FREE_CALL(), true)
	call.AddChildToFront(callback)
	t.ReportCodeChange_0()
}

// Java method: FindImportsAndExports.reportModuleErrors()
func (this_ *ProcessCommonJSModules_FindImportsAndExports) ReportModuleErrors() {
	for _, error_ := range this_.errors {
		_ = error_
		this_.outer_.compiler.Report_1(error_)
	}

}

// Java method: FindImportsAndExports.initializeModule()
func (this_ *ProcessCommonJSModules_FindImportsAndExports) InitializeModule() bool {
	var ci *CompilerInput = this_.outer_.compiler.GetInput(this_.script.GetInputId())
	_ = ci
	var modulePath *ModuleLoader_ModulePath = ci.GetPath()
	_ = modulePath
	if modulePath == nil {
		return true
	}
	var moduleName string = ProcessCommonJSModules_GetModuleName_1_CompilerInput(ci)
	_ = moduleName
	var exportsToRemove j2gort.Ext_java_util_List = j2gort.NewExt_java_util_ArrayList()
	_ = exportsToRemove
	{
		cur_20 := j2gort.NewIterCursor(j2gort.CheckNotNull(this_.exports))
		for cur_20.HasNext() {
			rawElem_21 := cur_20.Next()
			export := j2gort.CastPtr[ProcessCommonJSModules_ExportInfo](rawElem_21)
			_ = export
			if NodeUtil_GetEnclosingScript(export.node) == 0 {
				exportsToRemove.Add_1(j2gort.PtrToAny(export))
				continue
			}
			var qNameBase pkg_rhino.Node = this_.outer_.priv_ProcessCommonJSModules_getBaseQualifiedNameNode_4517(export.node)
			_ = qNameBase
			if ((((export.node == qNameBase) && export.node.GetParent().IsAssign()) && export.node.GetGrandparent().IsExprResult()) && (export.node.GetPrevious() == 0)) && (export.node.GetNext() != 0) {
				if export.node.GetNext().IsGetProp() && export.node.GetNext().MatchesQualifiedName_1_string(("module" + "." + "exports")) {
					{
						cur_22 := j2gort.NewIterCursor(j2gort.CheckNotNull(this_.moduleExports))
						for cur_22.HasNext() {
							rawElem_23 := cur_22.Next()
							moduleExport := j2gort.CastPtr[ProcessCommonJSModules_ExportInfo](rawElem_23)
							_ = moduleExport
							if moduleExport.node == export.node.GetNext() {
								this_.moduleExports.Remove_1_any(j2gort.PtrToAny(moduleExport))
								break
							}
						}
					}
					var changeRoot pkg_rhino.Node = export.node.GetGrandparent().GetParent()
					_ = changeRoot
					export.node.GetGrandparent().Detach()
					exportsToRemove.Add_1(j2gort.PtrToAny(export))
					this_.outer_.compiler.ReportChangeToEnclosingScope(changeRoot)
				} else if export.node.GetNext().IsAssign() && export.node.GetNext().GetFirstChild().MatchesQualifiedName_1_string(("module" + "." + "exports")) {
					var assign pkg_rhino.Node = export.node.GetNext()
					_ = assign
					export.node.GetParent().ReplaceWith(assign.Detach())
					exportsToRemove.Add_1(j2gort.PtrToAny(export))
					this_.outer_.compiler.ReportChangeToEnclosingScope(assign)
				}
			} else if ((((export.node.GetParent().IsGetElem() && export.node.GetNext().IsStringLit()) && func() bool {
				sr_24 := export.node.GetNext().GetString()
				se_25 := "default"
				return j2gort.CheckNotNullString(sr_24) == se_25
			}()) && export.node.GetGrandparent().IsAssign()) && (export.node.GetParent().GetPrevious() != 0)) && export.node.GetParent().GetPrevious().MatchesQualifiedName_1_string(("module" + "." + "exports")) {
				var parent pkg_rhino.Node = export.node.GetParent()
				_ = parent
				var grandparent pkg_rhino.Node = parent.GetParent()
				_ = grandparent
				var prop pkg_rhino.Node = export.node.GetNext()
				_ = prop
				parent.ReplaceWith(pkg_rhino.IR_Getprop_2(export.node.Detach(), prop.Detach().GetString()).Srcref(parent))
				this_.outer_.compiler.ReportChangeToEnclosingScope(grandparent)
			}
		}
	}
	this_.exports.RemoveAll(exportsToRemove)
	exportsToRemove.Clear()
	var exportsToReplace j2gort.OrderedMap[*ProcessCommonJSModules_ExportInfo, *ProcessCommonJSModules_ExportInfo] = j2gort.NewOrderedMap[*ProcessCommonJSModules_ExportInfo, *ProcessCommonJSModules_ExportInfo]()
	_ = exportsToReplace
	{
		cur_26 := j2gort.NewIterCursor(j2gort.CheckNotNull(this_.moduleExports))
		for cur_26.HasNext() {
			rawElem_27 := cur_26.Next()
			export_2 := j2gort.CastPtr[ProcessCommonJSModules_ExportInfo](rawElem_27)
			_ = export_2
			if NodeUtil_GetEnclosingScript(export_2.node) == 0 {
				exportsToRemove.Add_1(j2gort.PtrToAny(export_2))
			} else if export_2.node.IsGetElem() {
				var prop_2 pkg_rhino.Node = export_2.node.GetSecondChild().Detach()
				_ = prop_2
				var newExport *ProcessCommonJSModules_ExportInfo = NewProcessCommonJSModules_ExportInfo(pkg_rhino.IR_Getprop_2(export_2.node.RemoveFirstChild(), prop_2.GetString()), export_2.scope)
				_ = newExport
				export_2.node.ReplaceWith(newExport.node)
				this_.outer_.compiler.ReportChangeToEnclosingScope(newExport.node)
				exportsToReplace.Put(export_2, newExport)
			}
		}
	}
	this_.moduleExports.RemoveAll(exportsToRemove)
	{
		cur_28 := exportsToReplace.IterKeys()
		for cur_28.HasNext() {

			oldExport := cur_28.
				Next()
			_ = oldExport
			var oldIndex int32 = this_.moduleExports.IndexOf_1(j2gort.PtrToAny(oldExport))
			_ = oldIndex
			this_.moduleExports.Remove_1_int32(oldIndex)
			this_.moduleExports.Add_2(oldIndex, j2gort.PtrToAny(exportsToReplace.Get(oldExport)))
		}
	}
	var directAssignments int32 = int32(0)
	_ = directAssignments
	{
		cur_30 := j2gort.NewIterCursor(j2gort.CheckNotNull(this_.moduleExports))
		for cur_30.HasNext() {
			rawElem_31 := cur_30.Next()
			export_3 := j2gort.CastPtr[ProcessCommonJSModules_ExportInfo](rawElem_31)
			_ = export_3
			if NodeUtil_GetEnclosingScript(export_3.node) == 0 {
				continue
			}
			var base pkg_rhino.Node = this_.outer_.priv_ProcessCommonJSModules_getBaseQualifiedNameNode_4517(export_3.node)
			_ = base
			if (base == export_3.node) && export_3.node.GetParent().IsAssign() {
				var rValue pkg_rhino.Node = NodeUtil_GetRValueOfLValue(export_3.node)
				_ = rValue
				if (rValue == 0) || !(rValue.IsObjectLit()) {
					directAssignments = directAssignments + 1
				} else if rValue.IsObjectLit() {
					{
						var key pkg_rhino.Node = rValue.GetFirstChild()
						_ = key
						for ; key != 0; key = key.GetNext() {
							if (!(key.IsStringKey()) || key.IsQuotedStringKey()) && !(key.IsMemberFunctionDef()) {
								directAssignments = directAssignments + 1
								break
							}
						}
					}
				}
			}
		}
	}
	var initModule pkg_rhino.Node = pkg_rhino.IR_Var_2(pkg_rhino.IR_Name(moduleName), pkg_rhino.IR_Objectlit([]pkg_rhino.Node{}))
	_ = initModule
	initModule.GetFirstChild().PutBooleanProp(pkg_rhino.StaticGet_Node_MODULE_EXPORT(), true)
	initModule.GetFirstChild().MakeNonIndexable()
	var builder *pkg_rhino.JSDocInfo_Builder = pkg_rhino.JSDocInfo_S_Builder().ParseDocumentation()
	_ = builder
	builder.RecordConstancy()
	initModule.SetJSDocInfo(builder.Build_0())
	if (directAssignments == int32(0)) || (!(this_.exports.IsEmpty()) && !(this_.moduleExports.IsEmpty())) {
		var defaultProp pkg_rhino.Node = pkg_rhino.IR_StringKey_1("default")
		_ = defaultProp
		defaultProp.PutBooleanProp(pkg_rhino.StaticGet_Node_MODULE_EXPORT(), true)
		defaultProp.AddChildToFront(pkg_rhino.IR_Objectlit([]pkg_rhino.Node{}))
		initModule.GetFirstFirstChild().AddChildToFront(defaultProp)
		if this_.exports.IsEmpty() || this_.moduleExports.IsEmpty() {
			builder = pkg_rhino.JSDocInfo_S_Builder().ParseDocumentation()
			builder.RecordConstancy()
			defaultProp.SetJSDocInfo(builder.Build_0())
		}
	}
	this_.script.AddChildToFront(initModule.SrcrefTree(this_.script))
	this_.outer_.compiler.ReportChangeToEnclosingScope(this_.script)
	return ((directAssignments < int32(2)) && (this_.exports.IsEmpty() || this_.moduleExports.IsEmpty()))
}

// Java method: FindImportsAndExports.getOutermostUmdTest(com.google.javascript.rhino.Node n)
func (this_ *ProcessCommonJSModules_FindImportsAndExports) priv_FindImportsAndExports_getOutermostUmdTest_4529(n pkg_rhino.Node) *ProcessCommonJSModules_FindImportsAndExports_UmdTestInfo {
	_ = n
	if ((n == 0) || NodeUtil_IsTopLevel(n)) || n.IsFunction() {
		return nil
	}
	var parent pkg_rhino.Node = n.GetParent()
	_ = parent
	if parent == 0 {
		return nil
	}
	if parent.IsIf() || parent.IsHook() {
		var umdTestInfo *ProcessCommonJSModules_FindImportsAndExports_UmdTestInfo = this_.priv_FindImportsAndExports_getOutermostUmdTest_4529(parent)
		_ = umdTestInfo
		if umdTestInfo != nil {
			if ((((parent.IsIf() && (parent.GetSecondChild() == n)) && parent.HasParent()) && parent.GetParent().IsBlock()) && (parent.GetNext() == 0)) && (umdTestInfo.activeBranch == parent.GetParent()) {
				return NewProcessCommonJSModules_FindImportsAndExports_UmdTestInfo(umdTestInfo.enclosingIf, n)
			}
			return umdTestInfo
		}
		var umdTests j2gort.Ext_java_util_List = j2gort.NewExt_java_util_ArrayList()
		_ = umdTests
		NodeUtil_VisitPreOrder_2(parent.GetFirstChild(), NewProcessCommonJSModules__Anon_73(this_, umdTests))
		if umdTests.IsEmpty() && parent.GetFirstChild().IsTrue() {
			umdTests.Add_1(j2gort.BoxNode(parent.GetFirstChild()))
		}
		if !(umdTests.IsEmpty()) {
			return NewProcessCommonJSModules_FindImportsAndExports_UmdTestInfo(parent, n)
		}
		return nil
	}
	return this_.priv_FindImportsAndExports_getOutermostUmdTest_4529(parent)
}

// Java method: FindImportsAndExports.isInIfTest(com.google.javascript.rhino.Node n)
func (this_ *ProcessCommonJSModules_FindImportsAndExports) priv_FindImportsAndExports_isInIfTest_4530(n pkg_rhino.Node) bool {
	_ = n
	if ((n == 0) || NodeUtil_IsTopLevel(n)) || n.IsFunction() {
		return false
	}
	var parent pkg_rhino.Node = n.GetParent()
	_ = parent
	if parent == 0 {
		return false
	}
	if (parent.IsIf() || parent.IsHook()) && (parent.GetFirstChild() == n) {
		return true
	}
	return this_.priv_FindImportsAndExports_isInIfTest_4530(parent)
}

// Java method: FindImportsAndExports.replaceUmdPatterns()
func (this_ *ProcessCommonJSModules_FindImportsAndExports) ReplaceUmdPatterns() bool {
	var needsRetraverse bool = false
	_ = needsRetraverse
	var changeScope pkg_rhino.Node = pkg_rhino.Node(0)
	_ = changeScope
	{
		cur_32 := j2gort.NewIterCursor(j2gort.CheckNotNull(this_.umdPatterns))
		for cur_32.HasNext() {
			rawElem_33 := cur_32.Next()
			umdPattern := j2gort.CastPtr[ProcessCommonJSModules_UmdPattern](rawElem_33)
			_ = umdPattern
			if NodeUtil_GetEnclosingScript(umdPattern.ifRoot) == 0 {
				this_.outer_.priv_ProcessCommonJSModules_reportNestedScopesDeleted_4520(umdPattern.ifRoot)
				continue
			}
			var parent pkg_rhino.Node = umdPattern.ifRoot.GetParent()
			_ = parent
			var newNode pkg_rhino.Node = umdPattern.activeBranch
			_ = newNode
			if newNode == 0 {
				umdPattern.ifRoot.Detach()
				this_.outer_.priv_ProcessCommonJSModules_reportNestedScopesDeleted_4520(umdPattern.ifRoot)
				this_.outer_.compiler.ReportChangeToEnclosingScope(parent)
				needsRetraverse = true
				continue
			}
			if umdPattern.activeBranch.IsBlock() && umdPattern.activeBranch.HasOneChild() {
				newNode = umdPattern.activeBranch.RemoveFirstChild()
			} else {
				newNode.Detach()
			}
			needsRetraverse = true
			umdPattern.ifRoot.ReplaceWith(newNode)
			this_.outer_.priv_ProcessCommonJSModules_reportNestedScopesDeleted_4520(umdPattern.ifRoot)
			changeScope = ChangeTracker_GetEnclosingChangeScopeRoot(newNode)
			if changeScope != 0 {
				this_.outer_.compiler.ReportChangeToEnclosingScope(newNode)
			}
			var block pkg_rhino.Node = parent
			_ = block
			if block.IsExprResult() {
				block = block.GetParent()
			}
			if ((block.IsBlock() && block.GetParent().IsFunction()) && block.GetGrandparent().IsCall()) && parent.HasOneChild() {
				var enclosingFnCall pkg_rhino.Node = block.GetGrandparent()
				_ = enclosingFnCall
				var fn pkg_rhino.Node = block.GetParent()
				_ = fn
				var enclosingScript pkg_rhino.Node = NodeUtil_GetEnclosingScript(enclosingFnCall)
				_ = enclosingScript
				if enclosingScript == 0 {
					continue
				}
				var ci *CompilerInput = this_.outer_.compiler.GetInput(NodeUtil_GetEnclosingScript(enclosingFnCall).GetInputId())
				_ = ci
				var modulePath *ModuleLoader_ModulePath = ci.GetPath()
				_ = modulePath
				if modulePath == nil {
					continue
				}
				needsRetraverse = true
				var factoryLabel string = (j2gort.StrValString(modulePath.ToModuleName()) + "_factory" + j2gort.StrValString(this_.outer_.compiler.GetUniqueIdSupplier().GetUniqueId(ci)))
				_ = factoryLabel
				var mutator *FunctionToBlockMutator = NewFunctionToBlockMutator_2(this_.outer_.compiler, j2gort.AsIface[j2gort.Ext_java_util_function_Supplier](this_.outer_.compiler.GetUniqueNameIdSupplier()))
				_ = mutator
				var newStatements pkg_rhino.Node = mutator.MutateWithoutRenaming(factoryLabel, fn, enclosingFnCall, j2gort.NullString, false, false)
				_ = newStatements
				if ((((newStatements.IsBlock() && newStatements.HasTwoChildren()) && newStatements.GetFirstChild().IsVar()) && newStatements.GetFirstFirstChild().HasOneChild()) && newStatements.GetFirstFirstChild().GetFirstChild().IsFunction()) && newStatements.GetSecondChild().IsExprResult() {
					var inlinedFn pkg_rhino.Node = newStatements.GetFirstFirstChild().GetFirstChild()
					_ = inlinedFn
					var expr pkg_rhino.Node = newStatements.GetSecondChild().GetFirstChild()
					_ = expr
					var call pkg_rhino.Node = pkg_rhino.Node(0)
					_ = call
					var assignedName string = j2gort.NullString
					_ = assignedName
					if expr.IsAssign() && expr.GetSecondChild().IsCall() {
						call = expr.GetSecondChild()
						assignedName = (j2gort.StrValString(modulePath.ToModuleName()) + "_iife" + j2gort.StrValString(this_.outer_.compiler.GetUniqueIdSupplier().GetUniqueId(ci)))
					} else if expr.IsCall() {
						call = expr
					}
					if call != 0 {
						newStatements = mutator.MutateWithoutRenaming(factoryLabel, inlinedFn, call, assignedName, false, false)
						if assignedName != j2gort.NullString {
							var newName pkg_rhino.Node = pkg_rhino.IR_Var_1(NodeUtil_NewName_4_AbstractCompiler_string_Node_string(this_.outer_.compiler, assignedName, fn, expr.GetFirstChild().GetQualifiedName())).SrcrefTree(fn)
							_ = newName
							if (((newStatements.HasChildren() && newStatements.GetFirstChild().IsExprResult()) && newStatements.GetFirstFirstChild().IsAssign()) && newStatements.GetFirstFirstChild().GetFirstChild().IsName()) && (j2gort.CheckNotNullString(newStatements.GetFirstFirstChild().GetFirstChild().GetString()) == assignedName) {
								newName.GetFirstChild().AddChildToFront(newStatements.GetFirstFirstChild().GetSecondChild().Detach())
								newStatements.GetFirstChild().ReplaceWith(newName)
							} else {
								newStatements.AddChildToFront(newName)
							}
							expr.GetSecondChild().ReplaceWith(newName.GetFirstChild().CloneNode_0())
							newStatements.AddChildToBack(expr.GetParent().Detach())
						}
					}
				}
				var callRoot pkg_rhino.Node = enclosingFnCall.GetParent()
				_ = callRoot
				if callRoot.IsNot() {
					callRoot = callRoot.GetParent()
				}
				if callRoot.IsExprResult() {
					var callRootParent pkg_rhino.Node = callRoot.GetParent()
					_ = callRootParent
					callRootParent.AddChildrenAfter(newStatements.RemoveChildren(), callRoot)
					callRoot.Detach()
					this_.outer_.priv_ProcessCommonJSModules_reportNestedScopesChanged_4521(callRootParent)
					this_.outer_.compiler.ReportChangeToEnclosingScope(callRootParent)
					this_.outer_.priv_ProcessCommonJSModules_reportNestedScopesDeleted_4520(enclosingFnCall)
				} else {
					umdPattern.ifRoot.ReplaceWith(newNode)
					this_.outer_.compiler.ReportChangeToEnclosingScope(newNode)
					this_.outer_.priv_ProcessCommonJSModules_reportNestedScopesDeleted_4520(umdPattern.ifRoot)
				}
			}
		}
	}
	return needsRetraverse
}

// Java class: com.google.javascript.jscomp.ProcessCommonJSModules.FindImportsAndExports.UmdTestInfo
type ProcessCommonJSModules_FindImportsAndExports_UmdTestInfo struct {
	enclosingIf  pkg_rhino.Node `gson:"enclosingIf"`
	activeBranch pkg_rhino.Node `gson:"activeBranch"`
	idHash_      atomic.Int32
}

func (this_ *ProcessCommonJSModules_FindImportsAndExports_UmdTestInfo) BindSelf_(s any) {
	_ = s
}
func (this_ *ProcessCommonJSModules_FindImportsAndExports_UmdTestInfo) IsJavaNil_() bool {
	return this_ == nil
}
func (*ProcessCommonJSModules_FindImportsAndExports_UmdTestInfo) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.ProcessCommonJSModules$FindImportsAndExports$UmdTestInfo", "UmdTestInfo"
}
func (this_ *ProcessCommonJSModules_FindImportsAndExports_UmdTestInfo) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*ProcessCommonJSModules_FindImportsAndExports_UmdTestInfo) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.ProcessCommonJSModules.FindImportsAndExports.UmdTestInfo", "com.google.javascript.jscomp.ProcessCommonJSModules$FindImportsAndExports$UmdTestInfo", func(v_ any) bool {
		return j2gort.InstanceOf[*ProcessCommonJSModules_FindImportsAndExports_UmdTestInfo](v_)
	})
}

func (this_ *ProcessCommonJSModules_FindImportsAndExports_UmdTestInfo) initFields_ProcessCommonJSModules_FindImportsAndExports_UmdTestInfo() {
	_ = this_
}

// Java constructor: UmdTestInfo(com.google.javascript.rhino.Node enclosingIf, com.google.javascript.rhino.Node activeBranch)
func NewProcessCommonJSModules_FindImportsAndExports_UmdTestInfo(enclosingIf pkg_rhino.Node, activeBranch pkg_rhino.Node) *ProcessCommonJSModules_FindImportsAndExports_UmdTestInfo {
	this_ := &ProcessCommonJSModules_FindImportsAndExports_UmdTestInfo{}
	InitProcessCommonJSModules_FindImportsAndExports_UmdTestInfo(this_, enclosingIf, activeBranch)
	return this_
}

func InitProcessCommonJSModules_FindImportsAndExports_UmdTestInfo(this_ *ProcessCommonJSModules_FindImportsAndExports_UmdTestInfo, enclosingIf pkg_rhino.Node, activeBranch pkg_rhino.Node) {
	_ = enclosingIf
	_ = activeBranch
	this_.initFields_ProcessCommonJSModules_FindImportsAndExports_UmdTestInfo()
	this_.enclosingIf = enclosingIf
	this_.activeBranch = activeBranch
}

// Java class: com.google.javascript.jscomp.ProcessCommonJSModules.RewriteModule
type ProcessCommonJSModules_RewriteModule struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	outer_                    *ProcessCommonJSModules
	exports                   j2gort.Ext_com_google_common_collect_ImmutableCollection `gson:"exports"`
	imports                   []pkg_rhino.Node                                         `gson:"imports"`
	rewrittenClassExpressions []pkg_rhino.Node                                         `gson:"rewrittenClassExpressions"`
	functionsToHoist          []pkg_rhino.Node                                         `gson:"functionsToHoist"`
	allowFullRewrite          bool                                                     `gson:"allowFullRewrite"`
	defaultExportIsConst      bool                                                     `gson:"defaultExportIsConst"`
}

func (this_ *ProcessCommonJSModules_RewriteModule) BindSelf_(s any) {
	_ = s
}
func (this_ *ProcessCommonJSModules_RewriteModule) IsJavaNil_() bool { return this_ == nil }
func (*ProcessCommonJSModules_RewriteModule) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.ProcessCommonJSModules$RewriteModule", "RewriteModule"
}
func (this_ *ProcessCommonJSModules_RewriteModule) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.NodeTraversal_AbstractPostOrderCallbackBase.IdentityHash_()
}
func (*ProcessCommonJSModules_RewriteModule) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.ProcessCommonJSModules.RewriteModule", "com.google.javascript.jscomp.ProcessCommonJSModules$RewriteModule", func(v_ any) bool {
		return j2gort.InstanceOf[*ProcessCommonJSModules_RewriteModule](v_)
	})
}
func (this_ *ProcessCommonJSModules_RewriteModule) SetOuter_(o *ProcessCommonJSModules) {
	this_.outer_ = o
}
func (this_ *ProcessCommonJSModules_RewriteModule) GetOuter_() *ProcessCommonJSModules {
	return this_.outer_
}
func (*ProcessCommonJSModules_RewriteModule) Implements_jscomp_NodeTraversal_AbstractPostOrderCallback() {
}
func (*ProcessCommonJSModules_RewriteModule) Implements_jscomp_NodeTraversal_Callback() {}

func (this_ *ProcessCommonJSModules_RewriteModule) initFields_ProcessCommonJSModules_RewriteModule() {
	_ = this_
	this_.imports = []pkg_rhino.Node{}
	this_.rewrittenClassExpressions = []pkg_rhino.Node{}
	this_.functionsToHoist = []pkg_rhino.Node{}
}

// Java constructor: RewriteModule(boolean allowFullRewrite, com.google.common.collect.ImmutableCollection<com.google.javascript.jscomp.ProcessCommonJSModules.ExportInfo> exports, boolean defaultExportIsConst)
func NewProcessCommonJSModules_RewriteModule(outer_ *ProcessCommonJSModules, allowFullRewrite bool, exports j2gort.Ext_com_google_common_collect_ImmutableCollection, defaultExportIsConst bool) *ProcessCommonJSModules_RewriteModule {
	this_ := &ProcessCommonJSModules_RewriteModule{}
	this_.outer_ = outer_
	InitProcessCommonJSModules_RewriteModule(this_, allowFullRewrite, exports, defaultExportIsConst)
	return this_
}

func InitProcessCommonJSModules_RewriteModule(this_ *ProcessCommonJSModules_RewriteModule, allowFullRewrite bool, exports j2gort.Ext_com_google_common_collect_ImmutableCollection, defaultExportIsConst bool) {
	_ = allowFullRewrite
	_ = exports
	_ = defaultExportIsConst
	InitNodeTraversal_AbstractPostOrderCallback(&this_.NodeTraversal_AbstractPostOrderCallbackBase)
	this_.initFields_ProcessCommonJSModules_RewriteModule()
	this_.allowFullRewrite = allowFullRewrite
	this_.exports = exports
	this_.defaultExportIsConst = defaultExportIsConst
}

// Java method: RewriteModule.visit(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *ProcessCommonJSModules_RewriteModule) Visit_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = n
	_ = parent
	switch n.TokenSwitchOrd_() {
	case int32(94) /* Token.SCRIPT */ :
		{
			for _, clazz := range this_.rewrittenClassExpressions {
				_ = clazz
				clazz.GetFirstChild().ReplaceWith(pkg_rhino.IR_Empty().Srcref(clazz.GetFirstChild()))
				t.ReportCodeChange_0()
			}
		}
		var ci *CompilerInput = this_.outer_.compiler.GetInput(n.GetInputId())
		_ = ci
		var moduleName string = ProcessCommonJSModules_GetModuleName_1_CompilerInput(ci)
		_ = moduleName
		{
			var i int32 = int32(1)
			_ = i
			for ; i < int32(len(this_.functionsToHoist)); i = i + 1 {
				if this_.functionsToHoist[i].GetFirstFirstChild().MatchesQualifiedName_1_string(this_.outer_.GetBasePropertyImport_1(moduleName)) {
					var fncVar pkg_rhino.Node = this_.functionsToHoist[i]
					_ = fncVar
					this_.functionsToHoist = j2gort.SliceRemoveAt(this_.functionsToHoist, i)
					this_.functionsToHoist = j2gort.SliceInsertAt(this_.functionsToHoist, int32(0), fncVar)
					break
				}
			}
		}
		{
			var i_2 int32 = (int32(len(this_.functionsToHoist)) - int32(1))
			_ = i_2
			for ; i_2 >= int32(0); i_2 = i_2 - 1 {
				var functionExpr pkg_rhino.Node = this_.functionsToHoist[i_2]
				_ = functionExpr
				var scopeRoot pkg_rhino.Node = t.GetClosestHoistScopeRoot()
				_ = scopeRoot
				var insertionPoint pkg_rhino.Node = scopeRoot.GetFirstChild()
				_ = insertionPoint
				if (insertionPoint == 0) || !(insertionPoint.IsVar() && (j2gort.CheckNotNullString(insertionPoint.GetFirstChild().GetString()) == moduleName)) {
					insertionPoint = pkg_rhino.Node(0)
				}
				if insertionPoint == 0 {
					if scopeRoot.GetFirstChild() != functionExpr {
						scopeRoot.AddChildToFront(functionExpr.Detach())
					}
				} else if (insertionPoint != functionExpr) && (insertionPoint.GetNext() != functionExpr) {
					functionExpr.Detach().InsertAfter(insertionPoint)
				}
			}
		}
		{
			cur_35 := j2gort.NewIterCursor(j2gort.CheckNotNull(this_.exports))
			for cur_35.HasNext() {
				rawElem_36 := cur_35.Next()
				export := j2gort.CastPtr[ProcessCommonJSModules_ExportInfo](rawElem_36)
				_ = export
				this_.priv_RewriteModule_visitExport_4534(t, export)
			}
		}
		{
			for _, require := range this_.imports {
				_ = require
				this_.priv_RewriteModule_visitRequireCall_4533(t, require)
			}
		}
	case int32(28) /* Token.CALL */ :
		if this_.outer_.IsCommonJsImport(n) {
			this_.imports = append(this_.imports, n)
		}
	case int32(84) /* Token.VAR */, int32(111) /* Token.LET */, int32(97) /* Token.CONST */ :
		if (n.HasMoreThanOneChild() && !(NodeUtil_IsAnyFor(parent))) && !(parent.IsExport()) {
			var vars j2gort.Ext_java_util_List = this_.priv_RewriteModule_splitMultipleDeclarations_4542(n)
			_ = vars
			t.ReportCodeChange_0()
			{
				cur_38 := j2gort.NewIterCursor(j2gort.CheckNotNull(vars))
				for cur_38.HasNext() {
					var_ := pkg_rhino.Node(cur_38.NextNode())
					_ = var_
					this_.Visit_3(t, var_.GetFirstChild(), var_)
				}
			}
		}
		if (n.GetFirstChild().HasChildren() && n.GetFirstFirstChild().IsName()) && func() bool {
			sr_39 := n.GetFirstChild().GetString()
			se_40 := n.GetFirstFirstChild().GetString()
			return j2gort.CheckNotNullString(sr_39) == se_40
		}() {
			n.Detach()
			t.ReportCodeChange_0()
			return
		}
	case int32(32) /* Token.NAME */ :
		if (NodeUtil_IsNameDeclaration(n.GetParent()) && n.GetParent().HasMoreThanOneChild()) && !(NodeUtil_IsAnyFor(n.GetGrandparent())) {
			break
		}
		var qName string = n.GetQualifiedName()
		_ = qName
		if qName == j2gort.NullString {
			break
		}
		var nameDeclaration *Var = j2gort.CastPtr[Var](t.GetScope().GetVar(qName))
		_ = nameDeclaration
		if nameDeclaration != nil {
			if NodeUtil_IsLhsByDestructuring(n) {
				this_.priv_RewriteModule_maybeUpdateName_4536(t, n, nameDeclaration)
			} else if (nameDeclaration.GetNode() != 0) && j2gort.ObjEqCustomPtr(nameDeclaration.GetNode().GetInputId(), n.GetInputId()) {
				var enclosingDeclaration pkg_rhino.Node = NodeUtil_GetEnclosingNode(n, j2gort.NewNodeLambda_Ext_com_google_common_base_Predicate(func(rawNode_41 uint32) bool {
					node := pkg_rhino.Node(rawNode_41)
					_ = node
					return (node == nameDeclaration.GetNameNode())
				}))
				_ = enclosingDeclaration
				if ((enclosingDeclaration == 0) || (enclosingDeclaration == n)) || (j2gort.CastPtr[Scope](nameDeclaration.GetScope()) != t.GetScope()) {
					this_.priv_RewriteModule_maybeUpdateName_4536(t, n, nameDeclaration)
				}
			}
		} else if (this_.allowFullRewrite && func() bool {
			sr_42 := n.GetString()
			se_43 := "module"
			return j2gort.CheckNotNullString(sr_42) == se_43
		}()) && !((n.GetParent().IsGetProp() || n.GetParent().IsGetElem()) || n.GetParent().IsTypeOf()) {
			var objectLit pkg_rhino.Node = pkg_rhino.IR_Objectlit([]pkg_rhino.Node{}).Srcref(n)
			_ = objectLit
			n.ReplaceWith(objectLit)
			t.ReportCodeChange_1(objectLit)
		}
	case int32(26) /* Token.GETPROP */ :
		if n.MatchesQualifiedName_1_string(("module" + ".id")) || n.MatchesQualifiedName_1_string(("module" + ".filename")) {
			var v *Var = j2gort.CastPtr[Var](t.GetScope().GetVar("module"))
			_ = v
			if (v == nil) || v.IsExtern() {
				n.ReplaceWith(pkg_rhino.IR_String(t.GetInput().GetPath().String()).Srcref(n))
			}
		} else if ((this_.allowFullRewrite && n.GetFirstChild().IsName()) && func() bool {
			sr_44 := n.GetFirstChild().GetString()
			se_45 := "module"
			return j2gort.CheckNotNullString(sr_44) == se_45
		}()) && !(n.MatchesQualifiedName_1_string(("module" + "." + "exports"))) {
			var v_2 *Var = j2gort.CastPtr[Var](t.GetScope().GetVar("module"))
			_ = v_2
			if (v_2 == nil) || v_2.IsExtern() {
				n.GetFirstChild().ReplaceWith(pkg_rhino.IR_Objectlit([]pkg_rhino.Node{}).Srcref(n.GetFirstChild()))
			}
		}
	case int32(25) /* Token.TYPEOF */ :
		if (this_.allowFullRewrite && n.GetFirstChild().IsName()) && (func() bool {
			sr_46 := n.GetFirstChild().GetString()
			se_47 := "module"
			return j2gort.CheckNotNullString(sr_46) == se_47
		}() || func() bool {
			sr_48 := n.GetFirstChild().GetString()
			se_49 := "exports"
			return j2gort.CheckNotNullString(sr_48) == se_49
		}()) {
			var v_3 *Var = j2gort.CastPtr[Var](t.GetScope().GetVar(n.GetFirstChild().GetString()))
			_ = v_3
			if (v_3 == nil) || v_3.IsExtern() {
				n.ReplaceWith(pkg_rhino.IR_String("object"))
			}
		}
	default:
	}
	this_.priv_RewriteModule_fixTypeAnnotationsForNode_4532(t, n)
}

// Java method: RewriteModule.fixTypeAnnotationsForNode(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n)
func (this_ *ProcessCommonJSModules_RewriteModule) priv_RewriteModule_fixTypeAnnotationsForNode_4532(t *NodeTraversal, n pkg_rhino.Node) {
	_ = t
	_ = n
	var info *pkg_rhino.JSDocInfo = n.GetJSDocInfo()
	_ = info
	if info != nil {
		{
			var i int32 = int32(0)
			_ = i
			var count int32 = info.GetTypeNodeCount()
			_ = count
			for ; i < count; i = i + 1 {
				this_.priv_RewriteModule_fixTypeNode_4541(t, info.GetTypeNodeAt(i))
			}
		}
	}
}

// Java method: RewriteModule.visitRequireCall(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node require)
func (this_ *ProcessCommonJSModules_RewriteModule) priv_RewriteModule_visitRequireCall_4533(t *NodeTraversal, require pkg_rhino.Node) {
	_ = t
	_ = require
	var moduleName string = this_.outer_.priv_ProcessCommonJSModules_getImportedModuleName_4514(t, require)
	_ = moduleName
	var moduleRef pkg_rhino.Node = NodeUtil_NewQName_2_AbstractCompiler_string(this_.outer_.compiler, this_.outer_.GetBasePropertyImport_2(moduleName, require)).SrcrefTree(require.GetSecondChild())
	_ = moduleRef
	require.ReplaceWith(moduleRef)
	t.ReportCodeChange_0()
}

// Java method: RewriteModule.visitExport(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.jscomp.ProcessCommonJSModules.ExportInfo export)
func (this_ *ProcessCommonJSModules_RewriteModule) priv_RewriteModule_visitExport_4534(t *NodeTraversal, export *ProcessCommonJSModules_ExportInfo) {
	_ = t
	_ = export
	var root pkg_rhino.Node = this_.outer_.priv_ProcessCommonJSModules_getBaseQualifiedNameNode_4517(export.node)
	_ = root
	var rValue pkg_rhino.Node = NodeUtil_GetRValueOfLValue(root)
	_ = rValue
	if StaticGet_ProcessCommonJSModules_MODULE_EXPORTS().Matches(root) {
		if (((rValue != 0) && rValue.IsObjectLit()) && root.GetParent().IsAssign()) && root.GetGrandparent().IsExprResult() {
			if this_.priv_RewriteModule_expandObjectLitAssignment_4535(t, root, export.scope) {
				return
			}
		}
	}
	var moduleName string = ProcessCommonJSModules_GetModuleName_1_CompilerInput(t.GetInput())
	_ = moduleName
	var moduleInitialization *Var = j2gort.CastPtr[Var](t.GetScope().GetVar(moduleName))
	_ = moduleInitialization
	var rValueVar *Var = nil
	_ = rValueVar
	if (rValue != 0) && rValue.IsQualifiedName() {
		rValueVar = j2gort.CastPtr[Var](export.scope.GetVar(rValue.GetQualifiedName()))
	}
	if ((((((root.GetParent().IsAssign() && root.GetGrandparent().IsExprResult()) && ((root.GetNext() != 0) && (root.GetNext().IsName() || root.GetNext().IsGetProp()))) && root.GetGrandparent().IsExprResult()) && (rValueVar != nil)) && ((NodeUtil_GetEnclosingScript(rValueVar.GetNameNode()) == 0) || (rValueVar.GetNameNode().HasParent() && !(rValueVar.IsParam())))) && export.isInSupportedScope) && ((rValueVar.GetNameNode().GetParent() == 0) || !(NodeUtil_IsLhsByDestructuring(rValueVar.GetNameNode()))) {
		root.GetGrandparent().Detach()
		t.ReportCodeChange_0()
		return
	}
	moduleName = (j2gort.StrValString(moduleName) + "." + "default")
	var updatedExport pkg_rhino.Node = NodeUtil_NewQName_4_AbstractCompiler_string_Node_string(this_.outer_.compiler, moduleName, export.node, export.node.GetQualifiedName())
	_ = updatedExport
	updatedExport.PutBooleanProp(pkg_rhino.StaticGet_Node_MODULE_EXPORT(), true)
	var exportIsConst bool = (((this_.defaultExportIsConst && updatedExport.MatchesQualifiedName_1_string(this_.outer_.GetBasePropertyImport_1(ProcessCommonJSModules_GetModuleName_1_CompilerInput(t.GetInput())))) && (root == export.node)) && NodeUtil_IsLValue(export.node))
	_ = exportIsConst
	var changeScope pkg_rhino.Node = pkg_rhino.Node(0)
	_ = changeScope
	if ((((StaticGet_ProcessCommonJSModules_MODULE_EXPORTS().Matches(root) && (rValue != 0)) && (j2gort.CastPtr[Var](export.scope.GetVar("module.exports")) == nil)) && root.GetParent().IsAssign()) && root.GetGrandparent().IsExprResult()) && (moduleInitialization == nil) {
		var parent pkg_rhino.Node = root.GetParent()
		_ = parent
		var exportName pkg_rhino.Node = pkg_rhino.IR_ExprResult(pkg_rhino.IR_Assign(updatedExport, rValue.Detach()))
		_ = exportName
		if exportIsConst {
			var info *pkg_rhino.JSDocInfo_Builder = pkg_rhino.JSDocInfo_S_Builder()
			_ = info
			info.RecordConstancy()
			exportName.GetFirstChild().SetJSDocInfo(info.Build_0())
		}
		parent.GetParent().ReplaceWith(exportName.SrcrefTree(root.GetParent()))
		changeScope = ChangeTracker_GetEnclosingChangeScopeRoot(parent)
	} else if (((((root.GetNext() != 0) && root.GetNext().IsName()) && (rValueVar != nil)) && rValueVar.IsGlobal()) && export.isInSupportedScope) && ((rValueVar.GetNameNode().GetParent() == 0) || !(NodeUtil_IsLhsByDestructuring(rValueVar.GetNameNode()))) {
		var parent_2 pkg_rhino.Node = root.GetParent()
		_ = parent_2
		root.Detach()
		parent_2.ReplaceWith(root)
		if root == export.node {
			root = updatedExport
		}
		export.node.ReplaceWith(updatedExport)
		changeScope = ChangeTracker_GetEnclosingChangeScopeRoot(root)
	} else {
		export.node.ReplaceWith(updatedExport)
		if updatedExport.GetParent().IsAssign() && exportIsConst {
			var infoBuilder *pkg_rhino.JSDocInfo_Builder = pkg_rhino.JSDocInfo_Builder_MaybeCopyFrom(updatedExport.GetParent().GetJSDocInfo())
			_ = infoBuilder
			infoBuilder.RecordConstancy()
			updatedExport.GetParent().SetJSDocInfo(infoBuilder.Build_0())
		}
		changeScope = ChangeTracker_GetEnclosingChangeScopeRoot(updatedExport)
	}
	if changeScope != 0 {
		this_.outer_.compiler.ReportChangeToChangeScope(changeScope)
	}
}

// Java method: RewriteModule.expandObjectLitAssignment(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node export, com.google.javascript.jscomp.Scope scope)
func (this_ *ProcessCommonJSModules_RewriteModule) priv_RewriteModule_expandObjectLitAssignment_4535(t *NodeTraversal, export pkg_rhino.Node, scope *Scope) bool {
	_ = t
	_ = export
	_ = scope
	if !(export.GetParent().IsAssign()) {
		j2gort.CheckState(false)
	}
	var insertionRef pkg_rhino.Node = export.GetGrandparent()
	_ = insertionRef
	if !(insertionRef.IsExprResult()) {
		j2gort.CheckState(false)
	}
	var insertionParent pkg_rhino.Node = insertionRef.GetParent()
	_ = insertionParent
	j2gort.CheckNotNullNode(insertionParent)
	var rValue pkg_rhino.Node = NodeUtil_GetRValueOfLValue(export)
	_ = rValue
	var key pkg_rhino.Node = rValue.GetFirstChild()
	_ = key
	var removedNodes bool = false
	_ = removedNodes
	for key != 0 {
		if (!(key.IsStringKey()) || key.IsQuotedStringKey()) && !(key.IsMemberFunctionDef()) {
			key = key.GetNext()
			continue
		}
		var lhs pkg_rhino.Node = pkg_rhino.IR_Getprop_2(export.CloneTree_0(), key.GetString())
		_ = lhs
		var value pkg_rhino.Node = pkg_rhino.Node(0)
		_ = value
		if key.IsStringKey() {
			value = key.RemoveFirstChild()
		} else if key.IsMemberFunctionDef() {
			value = key.RemoveFirstChild()
		}
		var expr pkg_rhino.Node = pkg_rhino.IR_ExprResult(pkg_rhino.IR_Assign(lhs, value)).SrcrefTreeIfMissing(key)
		_ = expr
		expr.InsertAfter(insertionRef)
		var newExport *ProcessCommonJSModules_ExportInfo = NewProcessCommonJSModules_ExportInfo(lhs.GetFirstChild(), scope)
		_ = newExport
		this_.priv_RewriteModule_visitExport_4534(t, newExport)
		if (expr != 0) && expr.HasParent() {
			insertionRef = expr
		}
		var currentKey pkg_rhino.Node = key
		_ = currentKey
		key = key.GetNext()
		currentKey.Detach()
		removedNodes = true
	}
	if !(rValue.HasChildren()) {
		export.GetGrandparent().Detach()
		return true
	}
	if removedNodes {
		t.ReportCodeChange_1(rValue)
	}
	return false
}

// Java method: RewriteModule.maybeUpdateName(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.jscomp.Var var)
func (this_ *ProcessCommonJSModules_RewriteModule) priv_RewriteModule_maybeUpdateName_4536(t *NodeTraversal, n pkg_rhino.Node, var_ *Var) {
	_ = t
	_ = n
	_ = var_
	j2gort.CheckNotNullPtr(var_)
	if !(n.IsName() || n.IsGetProp()) {
		j2gort.CheckState(false)
	}
	if !(n.HasParent()) {
		j2gort.CheckState(false)
	}
	var importedModuleName string = this_.priv_RewriteModule_getModuleImportName_4540(t, var_.GetNode())
	_ = importedModuleName
	var name string = n.GetQualifiedName()
	_ = name
	if (importedModuleName != j2gort.NullString) && (n != var_.GetNode()) {
		this_.priv_RewriteModule_updateNameReference_4537(t, n, name, importedModuleName, false, false)
	} else if this_.allowFullRewrite {
		var exportedName string = this_.priv_RewriteModule_getExportedName_4538(t, n, var_)
		_ = exportedName
		if (((n != var_.GetNode()) || n.GetParent().IsClass()) && (exportedName == j2gort.NullString)) && (((var_ == nil) || (var_.GetNameNode().GetParent() == 0)) || !(NodeUtil_IsLhsByDestructuring(var_.GetNameNode()))) {
			if n.GetParent().IsClass() && (n.GetParent().GetFirstChild() == n) {
				this_.rewrittenClassExpressions = append(this_.rewrittenClassExpressions, n.GetParent())
			}
			return
		}
		if ((((importedModuleName == j2gort.NullString) && (exportedName != j2gort.NullString)) && !(j2gort.CheckNotNullString(exportedName) == name)) && !(var_.IsParam())) && ((var_.GetNameNode().GetParent() == 0) || !(NodeUtil_IsLhsByDestructuring(var_.GetNameNode()))) {
			var exportPropIsConst bool = (((this_.defaultExportIsConst && (j2gort.CheckNotNullString(this_.outer_.GetBasePropertyImport_1(ProcessCommonJSModules_GetModuleName_1_CompilerInput(t.GetInput()))) == exportedName)) && (this_.outer_.priv_ProcessCommonJSModules_getBaseQualifiedNameNode_4517(n) == n)) && NodeUtil_IsLValue(n))
			_ = exportPropIsConst
			this_.priv_RewriteModule_updateNameReference_4537(t, n, name, exportedName, true, exportPropIsConst)
		} else if var_.IsGlobal() {
			var currentModuleName string = ProcessCommonJSModules_GetModuleName_1_CompilerInput(t.GetInput())
			_ = currentModuleName
			if j2gort.CheckNotNullString(currentModuleName) == name {
				return
			}
			if "exports" == name {
				return
			}
			if this_.outer_.compiler.GetOptions().ShouldExportTestFunctions() && javalang.StringStartsWith(j2gort.CheckNotNullString(currentModuleName), "test") {
				return
			}
			var newName string = (j2gort.StrValString(name) + "$$" + j2gort.StrValString(currentModuleName))
			_ = newName
			this_.priv_RewriteModule_updateNameReference_4537(t, n, name, newName, false, false)
		}
	}
}

// Java method: RewriteModule.updateNameReference(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node nameRef, java.lang.String originalName, java.lang.String newName, boolean requireFunctionExpressions, boolean qualifiedNameIsConst)
func (this_ *ProcessCommonJSModules_RewriteModule) priv_RewriteModule_updateNameReference_4537(t *NodeTraversal, nameRef pkg_rhino.Node, originalName string, newName string, requireFunctionExpressions bool, qualifiedNameIsConst bool) {
	_ = t
	_ = nameRef
	_ = originalName
	_ = newName
	_ = requireFunctionExpressions
	_ = qualifiedNameIsConst
	var parent pkg_rhino.Node = nameRef.GetParent()
	_ = parent
	j2gort.CheckNotNullNode(parent)
	j2gort.CheckNotNullString(newName)
	var newNameIsQualified bool = (j2gort.JStrIndexOfCp(newName, int32(uint16(46))) >= int32(0))
	_ = newNameIsQualified
	var newNameIsModuleExport bool
	{
		sr_50 := newName
		se_51 := this_.outer_.GetBasePropertyImport_1(ProcessCommonJSModules_GetModuleName_1_CompilerInput(t.GetInput()))
		newNameIsModuleExport = j2gort.CheckNotNullString(sr_50) == se_51
	}

	_ = newNameIsModuleExport
	var newNameDeclaration *Var = j2gort.CastPtr[Var](t.GetScope().GetVar(newName))
	_ = newNameDeclaration
	switch parent.TokenSwitchOrd_() {
	case int32(105) /* Token.CLASS */ :
		if (parent.GetIndexOfChild(nameRef) == int32(0)) && (newNameIsQualified || requireFunctionExpressions) {
			this_.rewrittenClassExpressions = append(this_.rewrittenClassExpressions, parent)
			var newNameRef pkg_rhino.Node = NodeUtil_NewQName_4_AbstractCompiler_string_Node_string(this_.outer_.compiler, newName, nameRef, originalName)
			_ = newNameRef
			if newNameIsModuleExport {
				newNameRef.PutBooleanProp(pkg_rhino.StaticGet_Node_MODULE_EXPORT(), true)
			}
			var expr pkg_rhino.Node = pkg_rhino.Node(0)
			_ = expr
			if !(newNameIsQualified) && (newNameDeclaration == nil) {
				expr = pkg_rhino.IR_Let_2(newNameRef, pkg_rhino.IR_NullNode()).SrcrefTreeIfMissing(nameRef)
			} else {
				expr = pkg_rhino.IR_ExprResult(pkg_rhino.IR_Assign(newNameRef, pkg_rhino.IR_NullNode())).SrcrefTreeIfMissing(nameRef)
				var info *pkg_rhino.JSDocInfo_Builder = pkg_rhino.JSDocInfo_Builder_MaybeCopyFrom(parent.GetJSDocInfo())
				_ = info
				parent.SetJSDocInfo(nil)
				if qualifiedNameIsConst {
					info.RecordConstancy()
				}
				expr.GetFirstChild().SetJSDocInfo(info.Build_0())
				this_.priv_RewriteModule_fixTypeAnnotationsForNode_4532(t, expr.GetFirstChild())
			}
			parent.ReplaceWith(expr)
			if expr.IsLet() {
				expr.GetFirstFirstChild().ReplaceWith(parent)
			} else {
				expr.GetFirstChild().GetSecondChild().ReplaceWith(parent)
			}
		} else if parent.GetIndexOfChild(nameRef) == int32(1) {
			var newNameRef_2 pkg_rhino.Node = NodeUtil_NewQName_4_AbstractCompiler_string_Node_string(this_.outer_.compiler, newName, nameRef, originalName)
			_ = newNameRef_2
			if newNameIsModuleExport {
				newNameRef_2.PutBooleanProp(pkg_rhino.StaticGet_Node_MODULE_EXPORT(), true)
			}
			nameRef.ReplaceWith(newNameRef_2)
		} else {
			nameRef.SetString(newName)
			nameRef.SetOriginalName(originalName)
		}
	case int32(73) /* Token.FUNCTION */ :
		if newNameIsQualified || requireFunctionExpressions {
			if NodeUtil_IsFunctionExpression(parent) {
				return
			}
			var newNameRef_3 pkg_rhino.Node = NodeUtil_NewQName_4_AbstractCompiler_string_Node_string(this_.outer_.compiler, newName, nameRef, originalName)
			_ = newNameRef_3
			if newNameIsModuleExport {
				newNameRef_3.PutBooleanProp(pkg_rhino.StaticGet_Node_MODULE_EXPORT(), true)
			}
			nameRef.SetString("")
			var expr_2 pkg_rhino.Node = pkg_rhino.Node(0)
			_ = expr_2
			if !(newNameIsQualified) && (newNameDeclaration == nil) {
				expr_2 = pkg_rhino.IR_Var_2(newNameRef_3, pkg_rhino.IR_NullNode()).SrcrefTreeIfMissing(nameRef)
			} else {
				expr_2 = pkg_rhino.IR_ExprResult(pkg_rhino.IR_Assign(newNameRef_3, pkg_rhino.IR_NullNode())).SrcrefTreeIfMissing(nameRef)
			}
			parent.ReplaceWith(expr_2)
			if expr_2.IsVar() {
				expr_2.GetFirstFirstChild().ReplaceWith(parent)
			} else {
				expr_2.GetFirstChild().GetSecondChild().ReplaceWith(parent)
				var info_2 *pkg_rhino.JSDocInfo_Builder = pkg_rhino.JSDocInfo_Builder_MaybeCopyFrom(parent.GetJSDocInfo())
				_ = info_2
				parent.SetJSDocInfo(nil)
				if qualifiedNameIsConst {
					info_2.RecordConstancy()
				}
				expr_2.GetFirstChild().SetJSDocInfo(info_2.Build_0())
				this_.priv_RewriteModule_fixTypeAnnotationsForNode_4532(t, expr_2.GetFirstChild())
			}
			this_.functionsToHoist = append(this_.functionsToHoist, expr_2)
		} else {
			nameRef.SetString(newName)
			nameRef.SetOriginalName(originalName)
		}
	case int32(84) /* Token.VAR */, int32(111) /* Token.LET */, int32(97) /* Token.CONST */ :
		if parent.HasMoreThanOneChild() && !(NodeUtil_IsAnyFor(parent.GetParent())) {
			this_.priv_RewriteModule_splitMultipleDeclarations_4542(parent)
			parent = nameRef.GetParent()
			newNameDeclaration = j2gort.CastPtr[Var](t.GetScope().GetVar(newName))
		}
		if newNameIsQualified {
			if !(nameRef.HasChildren()) && (parent.GetJSDocInfo() == nil) {
				parent.Detach()
				break
			}
			var getProp pkg_rhino.Node = NodeUtil_NewQName_4_AbstractCompiler_string_Node_string(this_.outer_.compiler, newName, nameRef, originalName)
			_ = getProp
			if newNameIsModuleExport {
				getProp.PutBooleanProp(pkg_rhino.StaticGet_Node_MODULE_EXPORT(), true)
			}
			var info_3 *pkg_rhino.JSDocInfo = parent.GetJSDocInfo()
			_ = info_3
			parent.SetJSDocInfo(nil)
			if nameRef.HasChildren() {
				var assign pkg_rhino.Node = pkg_rhino.IR_Assign(getProp, nameRef.RemoveFirstChild())
				_ = assign
				assign.SetJSDocInfo(info_3)
				var expr_3 pkg_rhino.Node = pkg_rhino.IR_ExprResult(assign).SrcrefTreeIfMissing(nameRef)
				_ = expr_3
				parent.ReplaceWith(expr_3)
				var infoBuilder *pkg_rhino.JSDocInfo_Builder = pkg_rhino.JSDocInfo_Builder_MaybeCopyFrom(info_3)
				_ = infoBuilder
				parent.SetJSDocInfo(nil)
				if qualifiedNameIsConst {
					infoBuilder.RecordConstancy()
				}
				assign.SetJSDocInfo(infoBuilder.Build_0())
				this_.priv_RewriteModule_fixTypeAnnotationsForNode_4532(t, assign)
			} else {
				getProp.SetJSDocInfo(info_3)
				parent.ReplaceWith(pkg_rhino.IR_ExprResult(getProp).Srcref(getProp))
			}
		} else if (newNameDeclaration != nil) && (newNameDeclaration.GetNameNode() != nameRef) {
			if !(nameRef.HasChildren()) {
				parent.Detach()
				break
			}
			var name pkg_rhino.Node = NodeUtil_NewName_4_AbstractCompiler_string_Node_string(this_.outer_.compiler, newName, nameRef, originalName)
			_ = name
			var assign_2 pkg_rhino.Node = pkg_rhino.IR_Assign(name, nameRef.RemoveFirstChild())
			_ = assign_2
			var info_4 *pkg_rhino.JSDocInfo = parent.GetJSDocInfo()
			_ = info_4
			if info_4 != nil {
				parent.SetJSDocInfo(nil)
				assign_2.SetJSDocInfo(info_4)
			}
			parent.ReplaceWith(pkg_rhino.IR_ExprResult(assign_2).SrcrefTree(nameRef))
		} else {
			nameRef.SetString(newName)
			nameRef.SetOriginalName(originalName)
		}
	default:
		if !(newNameIsQualified) && nameRef.IsName() {
			nameRef.SetString(newName)
			nameRef.SetOriginalName(originalName)
		} else {
			var name_2 pkg_rhino.Node
			_ = name_2
			if newNameIsQualified {
				name_2 = NodeUtil_NewQName_4_AbstractCompiler_string_Node_string(this_.outer_.compiler, newName, nameRef, originalName)
			} else {
				name_2 = NodeUtil_NewName_4_AbstractCompiler_string_Node_string(this_.outer_.compiler, newName, nameRef, originalName)
			}
			if newNameIsModuleExport {
				name_2.PutBooleanProp(pkg_rhino.StaticGet_Node_MODULE_EXPORT(), true)
			}
			var info_5 *pkg_rhino.JSDocInfo = nameRef.GetJSDocInfo()
			_ = info_5
			if info_5 != nil {
				nameRef.SetJSDocInfo(nil)
				name_2.SetJSDocInfo(info_5)
			}
			name_2.SrcrefTree(nameRef)
			nameRef.ReplaceWith(name_2)
			if nameRef.HasChildren() {
				name_2.AddChildrenToFront(nameRef.RemoveChildren())
			}
		}
	}
	t.ReportCodeChange_0()
}

// Java method: RewriteModule.getExportedName(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.jscomp.Var var)
func (this_ *ProcessCommonJSModules_RewriteModule) priv_RewriteModule_getExportedName_4538(t *NodeTraversal, n pkg_rhino.Node, var_ *Var) string {
	_ = t
	_ = n
	_ = var_
	if (var_ == nil) || !(j2gort.ObjEqCustomPtr(var_.GetNode().GetInputId(), n.GetInputId())) {
		return n.GetQualifiedName()
	}
	var baseExportName string = this_.outer_.GetBasePropertyImport_1(ProcessCommonJSModules_GetModuleName_1_CompilerInput(t.GetInput()))
	_ = baseExportName
	{
		cur_52 := j2gort.NewIterCursor(j2gort.CheckNotNull(this_.exports))
		for cur_52.HasNext() {
			rawElem_53 := cur_52.Next()
			export := j2gort.CastPtr[ProcessCommonJSModules_ExportInfo](rawElem_53)
			_ = export
			var exportBase pkg_rhino.Node = this_.outer_.priv_ProcessCommonJSModules_getBaseQualifiedNameNode_4517(export.node)
			_ = exportBase
			var exportRValue pkg_rhino.Node = NodeUtil_GetRValueOfLValue(exportBase)
			_ = exportRValue
			if exportRValue == 0 {
				continue
			}
			var exportedName pkg_rhino.Node = this_.priv_RewriteModule_getExportedNameNode_4539(export)
			_ = exportedName
			if export.isInSupportedScope && ((exportRValue == n) || ((NodeUtil_IsClassExpression(exportRValue) || NodeUtil_IsFunctionExpression(exportRValue)) && (exportedName == n))) {
				return j2gort.NullString
			}
			var exportBaseQName string = exportBase.GetQualifiedName()
			_ = exportBaseQName
			if exportRValue.IsObjectLit() {
				if !("module.exports" == exportBaseQName) {
					return n.GetQualifiedName()
				}
				var key pkg_rhino.Node = exportRValue.GetFirstChild()
				_ = key
				var keyIsExport bool = false
				_ = keyIsExport
				for key != 0 {
					if (key.IsStringKey() && !(key.IsQuotedStringKey())) && NodeUtil_IsValidPropertyName(this_.outer_.compiler.GetOptions().GetLanguageIn().ToFeatureSet(), key.GetString()) {
						if key.GetFirstChild().IsQualifiedName() {
							if key.GetFirstChild() == n {
								return j2gort.NullString
							}
							var valVar *Var = j2gort.CastPtr[Var](t.GetScope().GetVar(key.GetFirstChild().GetQualifiedName()))
							_ = valVar
							if (valVar != nil) && (valVar.GetNameNode() == var_.GetNameNode()) {
								keyIsExport = true
								break
							}
						}
					}
					key = key.GetNext()
				}
				if (key != 0) && keyIsExport {
					if export.isInSupportedScope {
						return (j2gort.StrValString(baseExportName) + "." + j2gort.StrValString(key.GetString()))
					} else {
						return n.GetQualifiedName()
					}
				}
			} else {
				if !(export.isInSupportedScope) {
					return n.GetQualifiedName()
				}
				if var_.GetNameNode() == exportedName {
					var exportPrefix string = j2gort.NullString
					_ = exportPrefix
					if javalang.StringStartsWith(j2gort.CheckNotNullString(exportBaseQName), "module") {
						exportPrefix = ("module" + "." + "exports")
					} else {
						exportPrefix = "exports"
					}
					if j2gort.JStrLen(exportBaseQName) == j2gort.JStrLen(exportPrefix) {
						return baseExportName
					}
					return (j2gort.StrValString(baseExportName) + j2gort.StrValString(j2gort.JStrSubFrom(exportBaseQName, j2gort.JStrLen(exportPrefix))))
				}
			}
		}
	}
	return n.GetQualifiedName()
}

// Java method: RewriteModule.getExportedNameNode(com.google.javascript.jscomp.ProcessCommonJSModules.ExportInfo info)
func (this_ *ProcessCommonJSModules_RewriteModule) priv_RewriteModule_getExportedNameNode_4539(info *ProcessCommonJSModules_ExportInfo) pkg_rhino.Node {
	_ = info
	var qNameBase pkg_rhino.Node = this_.outer_.priv_ProcessCommonJSModules_getBaseQualifiedNameNode_4517(info.node)
	_ = qNameBase
	var rValue pkg_rhino.Node = NodeUtil_GetRValueOfLValue(qNameBase)
	_ = rValue
	if rValue == 0 {
		return pkg_rhino.Node(0)
	}
	if NodeUtil_IsFunctionExpression(rValue) {
		return rValue.GetFirstChild()
	}
	if NodeUtil_IsClassExpression(rValue) && (rValue.GetFirstChild() == qNameBase) {
		return rValue.GetFirstChild()
	}
	if !(rValue.IsName()) {
		return pkg_rhino.Node(0)
	}
	var var_ *Var = j2gort.CastPtr[Var](info.scope.GetVar(rValue.GetString()))
	_ = var_
	if var_ == nil {
		return pkg_rhino.Node(0)
	}
	return var_.GetNameNode()
}

// Java method: RewriteModule.getModuleImportName(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n)
func (this_ *ProcessCommonJSModules_RewriteModule) priv_RewriteModule_getModuleImportName_4540(t *NodeTraversal, n pkg_rhino.Node) string {
	_ = t
	_ = n
	var rValue pkg_rhino.Node = pkg_rhino.Node(0)
	_ = rValue
	var propSuffix string = ""
	_ = propSuffix
	var parent pkg_rhino.Node = n.GetParent()
	_ = parent
	if (parent != 0) && parent.IsStringKey() {
		var grandparent pkg_rhino.Node = parent.GetParent()
		_ = grandparent
		if grandparent.IsObjectPattern() && grandparent.GetParent().IsDestructuringLhs() {
			rValue = grandparent.GetNext()
			propSuffix = ("." + j2gort.StrValString(parent.GetString()))
		}
	}
	if (j2gort.CheckNotNullString(propSuffix) == "") && (parent != 0) {
		rValue = NodeUtil_GetRValueOfLValue(n)
	}
	if rValue == 0 {
		return j2gort.NullString
	}
	if rValue.IsCall() && this_.outer_.IsCommonJsImport(rValue) {
		return (j2gort.StrValString(this_.outer_.GetBasePropertyImport_2(this_.outer_.priv_ProcessCommonJSModules_getImportedModuleName_4514(t, rValue), rValue)) + j2gort.StrValString(propSuffix))
	} else if rValue.IsGetProp() && this_.outer_.IsCommonJsImport(rValue.GetFirstChild()) {
		var importName string = this_.outer_.GetBasePropertyImport_2(this_.outer_.priv_ProcessCommonJSModules_getImportedModuleName_4514(t, rValue.GetFirstChild()), rValue.GetFirstChild())
		_ = importName
		var suffix string = rValue.GetString()
		_ = suffix
		return (j2gort.StrValString(importName) + "." + j2gort.StrValString(suffix) + j2gort.StrValString(propSuffix))
	}
	return j2gort.NullString
}

// Java method: RewriteModule.fixTypeNode(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node typeNode)
func (this_ *ProcessCommonJSModules_RewriteModule) priv_RewriteModule_fixTypeNode_4541(t *NodeTraversal, typeNode pkg_rhino.Node) {
	_ = t
	_ = typeNode
	if typeNode.IsStringLit() {
		var name string = typeNode.GetString()
		_ = name
		if ModuleLoader_IsPathIdentifier(name) {
			var lastSlash int32 = j2gort.JStrLastIndexOfCp(name, int32(uint16(47)))
			_ = lastSlash
			var endIndex int32 = j2gort.JStrIndexOfCpFrom(name, int32(uint16(46)), lastSlash)
			_ = endIndex
			var localTypeName string = j2gort.NullString
			_ = localTypeName
			if endIndex == int32(-1) {
				endIndex = j2gort.JStrLen(name)
			} else {
				localTypeName = j2gort.JStrSubFrom(name, endIndex)
			}
			var moduleName string = j2gort.JStrSub(name, int32(0), endIndex)
			_ = moduleName
			var globalModuleName string = this_.outer_.priv_ProcessCommonJSModules_getImportedModuleName_4515(t, typeNode, moduleName)
			_ = globalModuleName
			var baseImportProperty string = this_.outer_.GetBasePropertyImport_1(globalModuleName)
			_ = baseImportProperty
			typeNode.SetString(func() string {
				if localTypeName == j2gort.NullString {
					return baseImportProperty
				}
				return (j2gort.StrValString(baseImportProperty) + j2gort.StrValString(localTypeName))
			}())
		} else {
			var wasRewritten bool = false
			_ = wasRewritten
			var endIndex_2 int32 = int32(-1)
			_ = endIndex_2
			for endIndex_2 < j2gort.JStrLen(name) {
				endIndex_2 = j2gort.JStrIndexOfCpFrom(name, int32(uint16(46)), (endIndex_2 + int32(1)))
				if endIndex_2 == int32(-1) {
					endIndex_2 = j2gort.JStrLen(name)
				}
				var baseName string = j2gort.JStrSub(name, int32(0), endIndex_2)
				_ = baseName
				var suffix string
				_ = suffix
				if endIndex_2 < j2gort.JStrLen(name) {
					suffix = j2gort.JStrSubFrom(name, endIndex_2)
				} else {
					suffix = ""
				}
				var typeDeclaration *Var = j2gort.CastPtr[Var](t.GetScope().GetVar(baseName))
				_ = typeDeclaration
				if ((typeDeclaration != nil) && (typeDeclaration.GetNode() != 0)) && j2gort.ObjEqCustomPtr(typeDeclaration.GetNode().GetInputId(), typeNode.GetInputId()) {
					var importedModuleName string = this_.priv_RewriteModule_getModuleImportName_4540(t, typeDeclaration.GetNode())
					_ = importedModuleName
					if importedModuleName != j2gort.NullString {
						typeNode.SetString((j2gort.StrValString(importedModuleName) + j2gort.StrValString(suffix)))
						typeNode.SetOriginalName(name)
						wasRewritten = true
						break
					} else if this_.allowFullRewrite {
						var exportedName string = this_.priv_RewriteModule_getExportedName_4538(t, typeNode, typeDeclaration)
						_ = exportedName
						if (exportedName != j2gort.NullString) && !(j2gort.CheckNotNullString(exportedName) == name) {
							typeNode.SetString((j2gort.StrValString(exportedName) + j2gort.StrValString(suffix)))
							typeNode.SetOriginalName(name)
							wasRewritten = true
							break
						}
					}
				}
			}
			if !(wasRewritten) && this_.allowFullRewrite {
				endIndex_2 = j2gort.JStrIndexOfCp(name, int32(uint16(46)))
				if endIndex_2 == int32(-1) {
					endIndex_2 = j2gort.JStrLen(name)
				}
				var baseName_2 string = j2gort.JStrSub(name, int32(0), endIndex_2)
				_ = baseName_2
				var typeDeclaration_2 *Var = j2gort.CastPtr[Var](t.GetScope().GetVar(baseName_2))
				_ = typeDeclaration_2
				if (typeDeclaration_2 != nil) && typeDeclaration_2.IsGlobal() {
					var moduleName_2 string = ProcessCommonJSModules_GetModuleName_1_CompilerInput(t.GetInput())
					_ = moduleName_2
					var newName string = (j2gort.StrValString(baseName_2) + "$$" + j2gort.StrValString(moduleName_2))
					_ = newName
					if endIndex_2 < j2gort.JStrLen(name) {
						newName = (j2gort.StrValString(newName) + j2gort.StrValString(j2gort.JStrSubFrom(name, endIndex_2)))
					}
					typeNode.SetString(newName)
					typeNode.SetOriginalName(name)
				}
			}
		}
	}
	{
		var child pkg_rhino.Node = typeNode.GetFirstChild()
		_ = child
		for ; child != 0; child = child.GetNext() {
			this_.priv_RewriteModule_fixTypeNode_4541(t, child)
		}
	}
}

// Java method: RewriteModule.splitMultipleDeclarations(com.google.javascript.rhino.Node var)
func (this_ *ProcessCommonJSModules_RewriteModule) priv_RewriteModule_splitMultipleDeclarations_4542(var_ pkg_rhino.Node) j2gort.Ext_java_util_List {
	_ = var_
	if !(NodeUtil_IsNameDeclaration(var_)) {
		j2gort.CheckState(false)
	}
	var vars []pkg_rhino.Node = []pkg_rhino.Node{}
	_ = vars
	var info *pkg_rhino.JSDocInfo = var_.GetJSDocInfo()
	_ = info
	for var_.GetSecondChild() != 0 {
		var newVar pkg_rhino.Node = pkg_rhino.NewNode_2(var_.GetToken(), var_.RemoveFirstChild())
		_ = newVar
		if info != nil {
			newVar.SetJSDocInfo(info.Clone_0())
		}
		newVar.Srcref(var_)
		newVar.InsertBefore(var_)
		vars = append(vars, newVar)
	}
	vars = append(vars, var_)
	return j2gort.WrapTypedArrayList(vars)
}

// Java class:
type ProcessCommonJSModules__Anon_74 struct {
	outer_  *ProcessCommonJSModules
	idHash_ atomic.Int32
}

func (this_ *ProcessCommonJSModules__Anon_74) BindSelf_(s any) {
	_ = s
}
func (this_ *ProcessCommonJSModules__Anon_74) IsJavaNil_() bool { return this_ == nil }
func (*ProcessCommonJSModules__Anon_74) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.ProcessCommonJSModules$1", ""
}
func (this_ *ProcessCommonJSModules__Anon_74) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*ProcessCommonJSModules__Anon_74) HasIdentityEquals_()                       {}
func (this_ *ProcessCommonJSModules__Anon_74) SetOuter_(o *ProcessCommonJSModules) { this_.outer_ = o }
func (this_ *ProcessCommonJSModules__Anon_74) GetOuter_() *ProcessCommonJSModules {
	return this_.outer_
}
func (*ProcessCommonJSModules__Anon_74) Implements_jscomp_NodeUtil_Visitor() {}

func (this_ *ProcessCommonJSModules__Anon_74) initFields_ProcessCommonJSModules__Anon_74() {
	_ = this_
}

func NewProcessCommonJSModules__Anon_74(outer_ *ProcessCommonJSModules) *ProcessCommonJSModules__Anon_74 {
	this_ := &ProcessCommonJSModules__Anon_74{}
	this_.outer_ = outer_
	this_.initFields_ProcessCommonJSModules__Anon_74()
	return this_
}

// Java method: .visit(com.google.javascript.rhino.Node n)
func (this_ *ProcessCommonJSModules__Anon_74) Visit(n pkg_rhino.Node) {
	_ = n
	if n.IsFunction() {
		this_.outer_.compiler.ReportFunctionDeleted(n)
	}
}

// Java class:
type ProcessCommonJSModules__Anon_75 struct {
	outer_  *ProcessCommonJSModules
	idHash_ atomic.Int32
}

func (this_ *ProcessCommonJSModules__Anon_75) BindSelf_(s any) {
	_ = s
}
func (this_ *ProcessCommonJSModules__Anon_75) IsJavaNil_() bool { return this_ == nil }
func (*ProcessCommonJSModules__Anon_75) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.ProcessCommonJSModules$2", ""
}
func (this_ *ProcessCommonJSModules__Anon_75) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*ProcessCommonJSModules__Anon_75) HasIdentityEquals_()                       {}
func (this_ *ProcessCommonJSModules__Anon_75) SetOuter_(o *ProcessCommonJSModules) { this_.outer_ = o }
func (this_ *ProcessCommonJSModules__Anon_75) GetOuter_() *ProcessCommonJSModules {
	return this_.outer_
}
func (*ProcessCommonJSModules__Anon_75) Implements_jscomp_NodeUtil_Visitor() {}

func (this_ *ProcessCommonJSModules__Anon_75) initFields_ProcessCommonJSModules__Anon_75() {
	_ = this_
}

func NewProcessCommonJSModules__Anon_75(outer_ *ProcessCommonJSModules) *ProcessCommonJSModules__Anon_75 {
	this_ := &ProcessCommonJSModules__Anon_75{}
	this_.outer_ = outer_
	this_.initFields_ProcessCommonJSModules__Anon_75()
	return this_
}

// Java method: .visit(com.google.javascript.rhino.Node n)
func (this_ *ProcessCommonJSModules__Anon_75) Visit(n pkg_rhino.Node) {
	_ = n
	if n.IsFunction() {
		this_.outer_.compiler.ReportChangeToChangeScope(n)
	}
}

// Java class:
type ProcessCommonJSModules__Anon_73 struct {
	outer_         *ProcessCommonJSModules_FindImportsAndExports
	cap_0_umdTests j2gort.Ext_java_util_List
	idHash_        atomic.Int32
}

func (this_ *ProcessCommonJSModules__Anon_73) BindSelf_(s any) {
	_ = s
}
func (this_ *ProcessCommonJSModules__Anon_73) IsJavaNil_() bool { return this_ == nil }
func (*ProcessCommonJSModules__Anon_73) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.ProcessCommonJSModules$FindImportsAndExports$1", ""
}
func (this_ *ProcessCommonJSModules__Anon_73) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*ProcessCommonJSModules__Anon_73) HasIdentityEquals_() {}
func (this_ *ProcessCommonJSModules__Anon_73) SetOuter_(o *ProcessCommonJSModules_FindImportsAndExports) {
	this_.outer_ = o
}
func (this_ *ProcessCommonJSModules__Anon_73) GetOuter_() *ProcessCommonJSModules_FindImportsAndExports {
	return this_.outer_
}
func (*ProcessCommonJSModules__Anon_73) Implements_jscomp_NodeUtil_Visitor() {}

func (this_ *ProcessCommonJSModules__Anon_73) initFields_ProcessCommonJSModules__Anon_73() {
	_ = this_
}

func NewProcessCommonJSModules__Anon_73(outer_ *ProcessCommonJSModules_FindImportsAndExports, cap_p_umdTests j2gort.Ext_java_util_List) *ProcessCommonJSModules__Anon_73 {
	this_ := &ProcessCommonJSModules__Anon_73{}
	this_.outer_ = outer_
	this_.cap_0_umdTests = cap_p_umdTests
	this_.initFields_ProcessCommonJSModules__Anon_73()
	return this_
}

// Java method: .visit(com.google.javascript.rhino.Node node)
func (this_ *ProcessCommonJSModules__Anon_73) Visit(node pkg_rhino.Node) {
	_ = node
	switch node.TokenSwitchOrd_() {
	case int32(32) /* Token.NAME */ :
		{
			var iifeCond_1 bool
			{

				sr_54 := node.GetString()
				se_55 := "module"
				iifeCond_1 = j2gort.CheckNotNullString(sr_54) == se_55
			}
			if iifeCond_1 || (j2gort.CheckNotNullString(node.GetString()) == "define") {
				break
			}
		}
		return
	case int32(26) /* Token.GETPROP */ :
		if StaticGet_ProcessCommonJSModules_WINDOW_DEFINE().Matches(node) {
			break
		}
		return
	case int32(35) /* Token.STRINGLIT */ :
		if node.GetParent().IsIn() && (j2gort.CheckNotNullString(node.GetString()) == "amd") {
			break
		}
		return
	default:
		return
	}
	this_.cap_0_umdTests.Add_1(j2gort.BoxNode(node))
}
func (*ProcessCommonJSModules) NoOpVisit_()                          {}
func (*ProcessCommonJSModules) CanReleaseSyntacticScopes_()          {}
func (*ProcessCommonJSModules_RewriteModule) AlwaysShouldTraverse_() {}
