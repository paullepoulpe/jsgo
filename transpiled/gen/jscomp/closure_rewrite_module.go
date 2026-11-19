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
// Port of com/google/javascript/jscomp/ClosureRewriteModule.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_parser "github.com/paullepoulpe/jsgo/transpiled/gen/parser"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
	"github.com/paullepoulpe/jsgo/transpiled/shims/gson"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.jscomp.ClosureRewriteModule
type ClosureRewriteModule struct {
	compiler                          AbstractCompiler                                      `gson:"compiler"`
	astFactory                        *AstFactory                                           `gson:"astFactory"`
	preprocessorSymbolTable           *PreprocessorSymbolTable                              `gson:"preprocessorSymbolTable"`
	syntheticExterns                  j2gort.OrderedMap[string, pkg_rhino.Node]             `gson:"syntheticExterns"`
	globalScope                       *Scope                                                `gson:"globalScope"`
	scriptStack                       j2gort.Deque[*ClosureRewriteModule_ScriptDescription] `gson:"scriptStack"`
	currentScript                     *ClosureRewriteModule_ScriptDescription               `gson:"currentScript"`
	rewriteState                      *ClosureRewriteModule_GlobalRewriteState              `gson:"rewriteState"`
	legacyScriptNamespacesAndPrefixes j2gort.OrderedSet[string]                             `gson:"legacyScriptNamespacesAndPrefixes"`
	unrecognizedRequires              []*ClosureRewriteModule_UnrecognizedRequire           `gson:"unrecognizedRequires"`
	googModuleGetCalls                j2gort.Ext_java_util_ArrayList                        `gson:"googModuleGetCalls"`
	googRequireDynamicCalls           j2gort.Ext_java_util_ArrayList                        `gson:"googRequireDynamicCalls"`
	globalTypedScope                  *TypedScope                                           `gson:"globalTypedScope"`
	idHash_                           atomic.Int32
	preserveSugar                     bool `gson:"preserveSugar"`
}

func (this_ *ClosureRewriteModule) BindSelf_(s any) {
	_ = s
}
func (this_ *ClosureRewriteModule) IsJavaNil_() bool { return this_ == nil }
func (*ClosureRewriteModule) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.ClosureRewriteModule", "ClosureRewriteModule"
}
func (this_ *ClosureRewriteModule) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*ClosureRewriteModule) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.ClosureRewriteModule", "com.google.javascript.jscomp.ClosureRewriteModule", func(v_ any) bool {
		return j2gort.InstanceOf[*ClosureRewriteModule](v_)
	})
}
func (*ClosureRewriteModule) Implements_jscomp_CompilerPass() {}

var static_ClosureRewriteModule_INVALID_MODULE_ID_ARG *DiagnosticType

func StaticGet_ClosureRewriteModule_INVALID_MODULE_ID_ARG() *DiagnosticType {
	clinit_ClosureRewriteModule()
	return static_ClosureRewriteModule_INVALID_MODULE_ID_ARG
}
func StaticSet_ClosureRewriteModule_INVALID_MODULE_ID_ARG(v *DiagnosticType) {
	clinit_ClosureRewriteModule()
	static_ClosureRewriteModule_INVALID_MODULE_ID_ARG = v
}

var static_ClosureRewriteModule_INVALID_PROVIDE_NAMESPACE *DiagnosticType

func StaticGet_ClosureRewriteModule_INVALID_PROVIDE_NAMESPACE() *DiagnosticType {
	clinit_ClosureRewriteModule()
	return static_ClosureRewriteModule_INVALID_PROVIDE_NAMESPACE
}
func StaticSet_ClosureRewriteModule_INVALID_PROVIDE_NAMESPACE(v *DiagnosticType) {
	clinit_ClosureRewriteModule()
	static_ClosureRewriteModule_INVALID_PROVIDE_NAMESPACE = v
}

var static_ClosureRewriteModule_INVALID_PROVIDE_CALL *DiagnosticType

func StaticGet_ClosureRewriteModule_INVALID_PROVIDE_CALL() *DiagnosticType {
	clinit_ClosureRewriteModule()
	return static_ClosureRewriteModule_INVALID_PROVIDE_CALL
}
func StaticSet_ClosureRewriteModule_INVALID_PROVIDE_CALL(v *DiagnosticType) {
	clinit_ClosureRewriteModule()
	static_ClosureRewriteModule_INVALID_PROVIDE_CALL = v
}

var static_ClosureRewriteModule_INVALID_GET_ALIAS *DiagnosticType

func StaticGet_ClosureRewriteModule_INVALID_GET_ALIAS() *DiagnosticType {
	clinit_ClosureRewriteModule()
	return static_ClosureRewriteModule_INVALID_GET_ALIAS
}
func StaticSet_ClosureRewriteModule_INVALID_GET_ALIAS(v *DiagnosticType) {
	clinit_ClosureRewriteModule()
	static_ClosureRewriteModule_INVALID_GET_ALIAS = v
}

var static_ClosureRewriteModule_INVALID_EXPORT_COMPUTED_PROPERTY *DiagnosticType

func StaticGet_ClosureRewriteModule_INVALID_EXPORT_COMPUTED_PROPERTY() *DiagnosticType {
	clinit_ClosureRewriteModule()
	return static_ClosureRewriteModule_INVALID_EXPORT_COMPUTED_PROPERTY
}
func StaticSet_ClosureRewriteModule_INVALID_EXPORT_COMPUTED_PROPERTY(v *DiagnosticType) {
	clinit_ClosureRewriteModule()
	static_ClosureRewriteModule_INVALID_EXPORT_COMPUTED_PROPERTY = v
}

var static_ClosureRewriteModule_USELESS_USE_STRICT_DIRECTIVE *DiagnosticType

func StaticGet_ClosureRewriteModule_USELESS_USE_STRICT_DIRECTIVE() *DiagnosticType {
	clinit_ClosureRewriteModule()
	return static_ClosureRewriteModule_USELESS_USE_STRICT_DIRECTIVE
}
func StaticSet_ClosureRewriteModule_USELESS_USE_STRICT_DIRECTIVE(v *DiagnosticType) {
	clinit_ClosureRewriteModule()
	static_ClosureRewriteModule_USELESS_USE_STRICT_DIRECTIVE = v
}

var static_ClosureRewriteModule_IMPORT_INLINING_SHADOWS_VAR *DiagnosticType

func StaticGet_ClosureRewriteModule_IMPORT_INLINING_SHADOWS_VAR() *DiagnosticType {
	clinit_ClosureRewriteModule()
	return static_ClosureRewriteModule_IMPORT_INLINING_SHADOWS_VAR
}
func StaticSet_ClosureRewriteModule_IMPORT_INLINING_SHADOWS_VAR(v *DiagnosticType) {
	clinit_ClosureRewriteModule()
	static_ClosureRewriteModule_IMPORT_INLINING_SHADOWS_VAR = v
}

var static_ClosureRewriteModule_ILLEGAL_DESTRUCTURING_DEFAULT_EXPORT *DiagnosticType

func StaticGet_ClosureRewriteModule_ILLEGAL_DESTRUCTURING_DEFAULT_EXPORT() *DiagnosticType {
	clinit_ClosureRewriteModule()
	return static_ClosureRewriteModule_ILLEGAL_DESTRUCTURING_DEFAULT_EXPORT
}
func StaticSet_ClosureRewriteModule_ILLEGAL_DESTRUCTURING_DEFAULT_EXPORT(v *DiagnosticType) {
	clinit_ClosureRewriteModule()
	static_ClosureRewriteModule_ILLEGAL_DESTRUCTURING_DEFAULT_EXPORT = v
}

var static_ClosureRewriteModule_ILLEGAL_DESTRUCTURING_NOT_EXPORTED *DiagnosticType

func StaticGet_ClosureRewriteModule_ILLEGAL_DESTRUCTURING_NOT_EXPORTED() *DiagnosticType {
	clinit_ClosureRewriteModule()
	return static_ClosureRewriteModule_ILLEGAL_DESTRUCTURING_NOT_EXPORTED
}
func StaticSet_ClosureRewriteModule_ILLEGAL_DESTRUCTURING_NOT_EXPORTED(v *DiagnosticType) {
	clinit_ClosureRewriteModule()
	static_ClosureRewriteModule_ILLEGAL_DESTRUCTURING_NOT_EXPORTED = v
}

var static_ClosureRewriteModule_LOAD_MODULE_FN_MISSING_RETURN *DiagnosticType

func StaticGet_ClosureRewriteModule_LOAD_MODULE_FN_MISSING_RETURN() *DiagnosticType {
	clinit_ClosureRewriteModule()
	return static_ClosureRewriteModule_LOAD_MODULE_FN_MISSING_RETURN
}
func StaticSet_ClosureRewriteModule_LOAD_MODULE_FN_MISSING_RETURN(v *DiagnosticType) {
	clinit_ClosureRewriteModule()
	static_ClosureRewriteModule_LOAD_MODULE_FN_MISSING_RETURN = v
}

var static_ClosureRewriteModule_ILLEGAL_MODULE_RENAMING_CONFLICT *DiagnosticType

func StaticGet_ClosureRewriteModule_ILLEGAL_MODULE_RENAMING_CONFLICT() *DiagnosticType {
	clinit_ClosureRewriteModule()
	return static_ClosureRewriteModule_ILLEGAL_MODULE_RENAMING_CONFLICT
}
func StaticSet_ClosureRewriteModule_ILLEGAL_MODULE_RENAMING_CONFLICT(v *DiagnosticType) {
	clinit_ClosureRewriteModule()
	static_ClosureRewriteModule_ILLEGAL_MODULE_RENAMING_CONFLICT = v
}

var static_ClosureRewriteModule_ILLEGAL_STMT_OF_GOOG_REQUIRE_DYNAMIC_IN_AWAIT *DiagnosticType

func StaticGet_ClosureRewriteModule_ILLEGAL_STMT_OF_GOOG_REQUIRE_DYNAMIC_IN_AWAIT() *DiagnosticType {
	clinit_ClosureRewriteModule()
	return static_ClosureRewriteModule_ILLEGAL_STMT_OF_GOOG_REQUIRE_DYNAMIC_IN_AWAIT
}
func StaticSet_ClosureRewriteModule_ILLEGAL_STMT_OF_GOOG_REQUIRE_DYNAMIC_IN_AWAIT(v *DiagnosticType) {
	clinit_ClosureRewriteModule()
	static_ClosureRewriteModule_ILLEGAL_STMT_OF_GOOG_REQUIRE_DYNAMIC_IN_AWAIT = v
}
func StaticGet_ClosureRewriteModule_MODULE_EXPORTS_PREFIX() string  { return "module$exports$" }
func StaticGet_ClosureRewriteModule_MODULE_CONTENTS_PREFIX() string { return "module$contents$" }

var static_ClosureRewriteModule_GOOG_FORWARDDECLARE pkg_rhino.Node

func StaticGet_ClosureRewriteModule_GOOG_FORWARDDECLARE() pkg_rhino.Node {
	clinit_ClosureRewriteModule()
	return static_ClosureRewriteModule_GOOG_FORWARDDECLARE
}
func StaticSet_ClosureRewriteModule_GOOG_FORWARDDECLARE(v pkg_rhino.Node) {
	clinit_ClosureRewriteModule()
	static_ClosureRewriteModule_GOOG_FORWARDDECLARE = v
}

var static_ClosureRewriteModule_GOOG_LOADMODULE pkg_rhino.Node

func StaticGet_ClosureRewriteModule_GOOG_LOADMODULE() pkg_rhino.Node {
	clinit_ClosureRewriteModule()
	return static_ClosureRewriteModule_GOOG_LOADMODULE
}
func StaticSet_ClosureRewriteModule_GOOG_LOADMODULE(v pkg_rhino.Node) {
	clinit_ClosureRewriteModule()
	static_ClosureRewriteModule_GOOG_LOADMODULE = v
}

var static_ClosureRewriteModule_GOOG_MODULE pkg_rhino.Node

func StaticGet_ClosureRewriteModule_GOOG_MODULE() pkg_rhino.Node {
	clinit_ClosureRewriteModule()
	return static_ClosureRewriteModule_GOOG_MODULE
}
func StaticSet_ClosureRewriteModule_GOOG_MODULE(v pkg_rhino.Node) {
	clinit_ClosureRewriteModule()
	static_ClosureRewriteModule_GOOG_MODULE = v
}

var static_ClosureRewriteModule_GOOG_MODULE_DECLARELEGACYNAMESPACE pkg_rhino.Node

func StaticGet_ClosureRewriteModule_GOOG_MODULE_DECLARELEGACYNAMESPACE() pkg_rhino.Node {
	clinit_ClosureRewriteModule()
	return static_ClosureRewriteModule_GOOG_MODULE_DECLARELEGACYNAMESPACE
}
func StaticSet_ClosureRewriteModule_GOOG_MODULE_DECLARELEGACYNAMESPACE(v pkg_rhino.Node) {
	clinit_ClosureRewriteModule()
	static_ClosureRewriteModule_GOOG_MODULE_DECLARELEGACYNAMESPACE = v
}

var static_ClosureRewriteModule_GOOG_MODULE_PREVENTMODULEEXPORTSEALING pkg_rhino.QualifiedName

func StaticGet_ClosureRewriteModule_GOOG_MODULE_PREVENTMODULEEXPORTSEALING() pkg_rhino.QualifiedName {
	clinit_ClosureRewriteModule()
	return static_ClosureRewriteModule_GOOG_MODULE_PREVENTMODULEEXPORTSEALING
}
func StaticSet_ClosureRewriteModule_GOOG_MODULE_PREVENTMODULEEXPORTSEALING(v pkg_rhino.QualifiedName) {
	clinit_ClosureRewriteModule()
	static_ClosureRewriteModule_GOOG_MODULE_PREVENTMODULEEXPORTSEALING = v
}

var static_ClosureRewriteModule_GOOG_MODULE_GET pkg_rhino.Node

func StaticGet_ClosureRewriteModule_GOOG_MODULE_GET() pkg_rhino.Node {
	clinit_ClosureRewriteModule()
	return static_ClosureRewriteModule_GOOG_MODULE_GET
}
func StaticSet_ClosureRewriteModule_GOOG_MODULE_GET(v pkg_rhino.Node) {
	clinit_ClosureRewriteModule()
	static_ClosureRewriteModule_GOOG_MODULE_GET = v
}

var static_ClosureRewriteModule_GOOG_PROVIDE pkg_rhino.Node

func StaticGet_ClosureRewriteModule_GOOG_PROVIDE() pkg_rhino.Node {
	clinit_ClosureRewriteModule()
	return static_ClosureRewriteModule_GOOG_PROVIDE
}
func StaticSet_ClosureRewriteModule_GOOG_PROVIDE(v pkg_rhino.Node) {
	clinit_ClosureRewriteModule()
	static_ClosureRewriteModule_GOOG_PROVIDE = v
}

var static_ClosureRewriteModule_GOOG_REQUIRE pkg_rhino.Node

func StaticGet_ClosureRewriteModule_GOOG_REQUIRE() pkg_rhino.Node {
	clinit_ClosureRewriteModule()
	return static_ClosureRewriteModule_GOOG_REQUIRE
}
func StaticSet_ClosureRewriteModule_GOOG_REQUIRE(v pkg_rhino.Node) {
	clinit_ClosureRewriteModule()
	static_ClosureRewriteModule_GOOG_REQUIRE = v
}

var static_ClosureRewriteModule_GOOG_REQUIRETYPE pkg_rhino.Node

func StaticGet_ClosureRewriteModule_GOOG_REQUIRETYPE() pkg_rhino.Node {
	clinit_ClosureRewriteModule()
	return static_ClosureRewriteModule_GOOG_REQUIRETYPE
}
func StaticSet_ClosureRewriteModule_GOOG_REQUIRETYPE(v pkg_rhino.Node) {
	clinit_ClosureRewriteModule()
	static_ClosureRewriteModule_GOOG_REQUIRETYPE = v
}

var static_ClosureRewriteModule_GOOG_REQUIREDYNAMIC pkg_rhino.Node

func StaticGet_ClosureRewriteModule_GOOG_REQUIREDYNAMIC() pkg_rhino.Node {
	clinit_ClosureRewriteModule()
	return static_ClosureRewriteModule_GOOG_REQUIREDYNAMIC
}
func StaticSet_ClosureRewriteModule_GOOG_REQUIREDYNAMIC(v pkg_rhino.Node) {
	clinit_ClosureRewriteModule()
	static_ClosureRewriteModule_GOOG_REQUIREDYNAMIC = v
}
func StaticGet_ClosureRewriteModule_GOOG_REQUIREDYNAMIC_NAME() string { return "goog.requireDynamic" }
func StaticGet_ClosureRewriteModule_IMPORT_HANDLER_NAME() string      { return "importHandler_" }

var clinitDone_ClosureRewriteModule uint32
var clinitState_ClosureRewriteModule j2gort.ClinitState

func clinit_ClosureRewriteModule() {
	if atomic.LoadUint32(&clinitDone_ClosureRewriteModule) != 0 {
		return
	}
	clinitSlow_ClosureRewriteModule()
}

func clinitSlow_ClosureRewriteModule() {
	if !j2gort.EnterClinit(&clinitState_ClosureRewriteModule, &clinitDone_ClosureRewriteModule) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_ClosureRewriteModule, &clinitDone_ClosureRewriteModule)
	static_ClosureRewriteModule_INVALID_MODULE_ID_ARG = DiagnosticType_Error("JSC_GOOG_MODULE_INVALID_MODULE_ID_ARG", "goog.module parameter must be a string literal")
	static_ClosureRewriteModule_INVALID_PROVIDE_NAMESPACE = DiagnosticType_Error("JSC_GOOG_MODULE_INVALID_PROVIDE_NAMESPACE", "goog.provide parameter must be a string literal.")
	static_ClosureRewriteModule_INVALID_PROVIDE_CALL = DiagnosticType_Error("JSC_GOOG_MODULE_INVALID_PROVIDE_CALL", "goog.provide can not be called in goog.module.")
	static_ClosureRewriteModule_INVALID_GET_ALIAS = DiagnosticType_Error("JSC_GOOG_MODULE_INVALID_GET_ALIAS", "goog.module.get should not be aliased.")
	static_ClosureRewriteModule_INVALID_EXPORT_COMPUTED_PROPERTY = DiagnosticType_Error("JSC_GOOG_MODULE_INVALID_EXPORT_COMPUTED_PROPERTY", "Computed properties are not yet supported in goog.module exports.")
	static_ClosureRewriteModule_USELESS_USE_STRICT_DIRECTIVE = DiagnosticType_Disabled("JSC_USELESS_USE_STRICT_DIRECTIVE", "'use strict' is unnecessary in goog.module files.")
	static_ClosureRewriteModule_IMPORT_INLINING_SHADOWS_VAR = DiagnosticType_Error("JSC_IMPORT_INLINING_SHADOWS_VAR", "Inlining of reference to import \"{1}\" shadows var \"{0}\".")
	static_ClosureRewriteModule_ILLEGAL_DESTRUCTURING_DEFAULT_EXPORT = DiagnosticType_Error("JSC_ILLEGAL_DESTRUCTURING_DEFAULT_EXPORT", "Destructuring import only allowed for importing module with named exports.\nSee https://github.com/google/closure-compiler/wiki/goog.module-style")
	static_ClosureRewriteModule_ILLEGAL_DESTRUCTURING_NOT_EXPORTED = DiagnosticType_Error("JSC_ILLEGAL_DESTRUCTURING_NOT_EXPORTED", "Destructuring import reference to name \"{0}\" was not exported in module {1}")
	static_ClosureRewriteModule_LOAD_MODULE_FN_MISSING_RETURN = DiagnosticType_Error("JSC_LOAD_MODULE_FN_MISSING_RETURN", "goog.loadModule function should end with 'return exports;'")
	static_ClosureRewriteModule_ILLEGAL_MODULE_RENAMING_CONFLICT = DiagnosticType_Error("JSC_ILLEGAL_MODULE_RENAMING_CONFLICT", "Internal compiler error: rewritten module global name {0} is already in use.\nOriginal definition: {1}")
	static_ClosureRewriteModule_ILLEGAL_STMT_OF_GOOG_REQUIRE_DYNAMIC_IN_AWAIT = DiagnosticType_Error("ILLEGAL_STMT_OF_GOOG_REQUIRE_DYNAMIC_IN_AWAIT", "Illegal use of dynamic import: LHS of await goog.requireDynamic() must be a destructing LHS or name, and it must be in a declaration statement.")
	static_ClosureRewriteModule_GOOG_FORWARDDECLARE = pkg_rhino.IR_Getprop_2(pkg_rhino.IR_Name("goog"), "forwardDeclare")
	static_ClosureRewriteModule_GOOG_LOADMODULE = pkg_rhino.IR_Getprop_2(pkg_rhino.IR_Name("goog"), "loadModule")
	static_ClosureRewriteModule_GOOG_MODULE = pkg_rhino.IR_Getprop_2(pkg_rhino.IR_Name("goog"), "module")
	static_ClosureRewriteModule_GOOG_MODULE_DECLARELEGACYNAMESPACE = pkg_rhino.IR_Getprop_2(StaticGet_ClosureRewriteModule_GOOG_MODULE(), "declareLegacyNamespace")
	static_ClosureRewriteModule_GOOG_MODULE_PREVENTMODULEEXPORTSEALING = pkg_rhino.QualifiedName_Of("goog.module.preventModuleExportSealing")
	static_ClosureRewriteModule_GOOG_MODULE_GET = pkg_rhino.IR_Getprop_2(StaticGet_ClosureRewriteModule_GOOG_MODULE().CloneTree_0(), "get")
	static_ClosureRewriteModule_GOOG_PROVIDE = pkg_rhino.IR_Getprop_2(pkg_rhino.IR_Name("goog"), "provide")
	static_ClosureRewriteModule_GOOG_REQUIRE = pkg_rhino.IR_Getprop_2(pkg_rhino.IR_Name("goog"), "require")
	static_ClosureRewriteModule_GOOG_REQUIRETYPE = pkg_rhino.IR_Getprop_2(pkg_rhino.IR_Name("goog"), "requireType")
	static_ClosureRewriteModule_GOOG_REQUIREDYNAMIC = pkg_rhino.IR_Getprop_2(pkg_rhino.IR_Name("goog"), "requireDynamic")
}

func (this_ *ClosureRewriteModule) initFields_ClosureRewriteModule() {
	_ = this_
	this_.syntheticExterns = j2gort.NewOrderedMap[string, pkg_rhino.Node]()

	this_.scriptStack = j2gort.Deque[*ClosureRewriteModule_ScriptDescription]{}

	this_.rewriteState = NewClosureRewriteModule_GlobalRewriteState()
	this_.legacyScriptNamespacesAndPrefixes = j2gort.NewOrderedSet[string]()
	this_.unrecognizedRequires = []*ClosureRewriteModule_UnrecognizedRequire{}
	this_.googModuleGetCalls = j2gort.NewExt_java_util_ArrayList()
	this_.googRequireDynamicCalls = j2gort.NewExt_java_util_ArrayList()
}

// Java constructor: ClosureRewriteModule(com.google.javascript.jscomp.AbstractCompiler compiler, com.google.javascript.jscomp.PreprocessorSymbolTable preprocessorSymbolTable, com.google.javascript.jscomp.@org.jspecify.annotations.Nullable TypedScope globalTypedScope)
func NewClosureRewriteModule(compiler AbstractCompiler, preprocessorSymbolTable *PreprocessorSymbolTable, globalTypedScope *TypedScope) *ClosureRewriteModule {
	clinit_ClosureRewriteModule()
	this_ := &ClosureRewriteModule{}
	InitClosureRewriteModule(this_, compiler, preprocessorSymbolTable, globalTypedScope)
	return this_
}

func InitClosureRewriteModule(this_ *ClosureRewriteModule, compiler AbstractCompiler, preprocessorSymbolTable *PreprocessorSymbolTable, globalTypedScope *TypedScope) {
	_ = compiler
	_ = preprocessorSymbolTable
	_ = globalTypedScope
	this_.initFields_ClosureRewriteModule()
	if !((globalTypedScope == nil) || globalTypedScope.IsGlobal()) {
		j2gort.CheckArgument(false)
	}
	this_.compiler = compiler
	this_.astFactory = compiler.CreateAstFactory()
	this_.preprocessorSymbolTable = preprocessorSymbolTable
	this_.preserveSugar = compiler.GetOptions().ShouldPreserveGoogModule()
	this_.globalTypedScope = globalTypedScope
}

// Java static method: ClosureRewriteModule.getBinaryModuleNamespace(java.lang.String namespaceId)
func ClosureRewriteModule_GetBinaryModuleNamespace(namespaceId string) string {
	clinit_ClosureRewriteModule()
	_ = namespaceId
	return ("module$exports$" + j2gort.StrValString(j2gort.JStrReplaceChar(namespaceId, uint16(46), uint16(36))))
}

// Java method: ClosureRewriteModule.rewriteJsdoc(com.google.javascript.rhino.JSDocInfo info, com.google.javascript.jscomp.Scope scope)
func (this_ *ClosureRewriteModule) priv_ClosureRewriteModule_rewriteJsdoc_1181(info *pkg_rhino.JSDocInfo, scope *Scope) {
	_ = info
	_ = scope
	{
		var i int32 = int32(0)
		_ = i
		var n int32 = info.GetTypeNodeCount()
		_ = n
		for ; i < n; i = i + 1 {
			this_.priv_ClosureRewriteModule_rewriteJsdocTypeNode_1182(info.GetTypeNodeAt(i), scope)
		}
	}
}

// Java method: ClosureRewriteModule.rewriteJsdocTypeNode(com.google.javascript.rhino.Node typeRefNode, com.google.javascript.jscomp.Scope scope)
func (this_ *ClosureRewriteModule) priv_ClosureRewriteModule_rewriteJsdocTypeNode_1182(typeRefNode pkg_rhino.Node, scope *Scope) {
	_ = typeRefNode
	_ = scope
	if typeRefNode.IsStringLit() {
		this_.priv_ClosureRewriteModule_rewriteJsdocStringTypeRef_1183(typeRefNode, scope)
	}
	{
		var child pkg_rhino.Node = typeRefNode.GetFirstChild()
		_ = child
		for ; child != 0; child = child.GetNext() {
			this_.priv_ClosureRewriteModule_rewriteJsdocTypeNode_1182(child, scope)
		}
	}
}

// Java method: ClosureRewriteModule.rewriteJsdocStringTypeRef(com.google.javascript.rhino.Node typeRefNode, com.google.javascript.jscomp.Scope scope)
func (this_ *ClosureRewriteModule) priv_ClosureRewriteModule_rewriteJsdocStringTypeRef_1183(typeRefNode pkg_rhino.Node, scope *Scope) {
	_ = typeRefNode
	_ = scope
	var typeName string = typeRefNode.GetString()
	_ = typeName
	var dot int32 = j2gort.JStrIndexOfCp(typeName, int32(uint16(46)))
	_ = dot
	var rootOfType string
	_ = rootOfType
	if dot == int32(-1) {
		rootOfType = typeName
	} else {
		rootOfType = j2gort.JStrSub(typeName, int32(0), dot)
	}
	var rootVar *Var = scope.LookupVar(rootOfType)
	_ = rootVar
	if ((rootVar != nil) && !(j2gort.CastPtr[Scope](j2gort.CastPtr[Scope](rootVar.GetScope()).GetClosestHoistScope()).IsGlobal())) && !(j2gort.CastPtr[Scope](rootVar.GetScope()).IsModuleScope()) {
		return
	}
	var importRecord *ClosureRewriteModule_AliasName = this_.currentScript.namesToInlineByAlias.Get(rootOfType)
	_ = importRecord
	if importRecord != nil {
		if this_.preprocessorSymbolTable != nil {
			var moduleOnlyNode pkg_rhino.Node = typeRefNode.CloneNode_0()
			_ = moduleOnlyNode
			this_.priv_ClosureRewriteModule_safeSetString_1227(moduleOnlyNode, rootOfType)
			moduleOnlyNode.SetLength(j2gort.JStrLen(rootOfType))
			this_.priv_ClosureRewriteModule_maybeAddAliasToSymbolTable_1230(moduleOnlyNode, this_.currentScript.namespaceId)
		}
		var aliasedNamespace string = importRecord.newName
		_ = aliasedNamespace
		var remainder string
		_ = remainder
		if dot == int32(-1) {
			remainder = ""
		} else {
			remainder = j2gort.JStrSubFrom(typeName, dot)
		}
		this_.priv_ClosureRewriteModule_safeSetString_1227(typeRefNode, (j2gort.StrValString(aliasedNamespace) + j2gort.StrValString(remainder)))
	} else if this_.currentScript.isModule && this_.currentScript.topLevelNames.Contains(rootOfType) {
		this_.priv_ClosureRewriteModule_safeSetString_1227(typeRefNode, (j2gort.StrValString(this_.currentScript.contentsPrefix) + j2gort.StrValString(typeName)))
	} else if this_.currentScript.isModule && (j2gort.CheckNotNullString(rootOfType) == "exports") {
		var namespace string = this_.currentScript.GetBinaryNamespace()
		_ = namespace
		var remainder_2 string
		_ = remainder_2
		if dot == int32(-1) {
			remainder_2 = ""
		} else {
			remainder_2 = j2gort.JStrSubFrom(typeName, dot)
		}
		this_.priv_ClosureRewriteModule_safeSetString_1227(typeRefNode, (j2gort.StrValString(namespace) + j2gort.StrValString(remainder_2)))
	} else {
		this_.priv_ClosureRewriteModule_rewriteIfClosureNamespaceRef_1184(typeName, typeRefNode)
	}
}

// Java method: ClosureRewriteModule.rewriteIfClosureNamespaceRef(java.lang.String typeName, com.google.javascript.rhino.Node typeRefNode)
func (this_ *ClosureRewriteModule) priv_ClosureRewriteModule_rewriteIfClosureNamespaceRef_1184(typeName string, typeRefNode pkg_rhino.Node) {
	_ = typeName
	_ = typeRefNode
	var prefixTypeName string = typeName
	_ = prefixTypeName
	for {
		var binaryNamespaceIfModule string = this_.rewriteState.GetBinaryNamespace(prefixTypeName)
		_ = binaryNamespaceIfModule
		if this_.legacyScriptNamespacesAndPrefixes.Contains(prefixTypeName) && (binaryNamespaceIfModule == j2gort.NullString) {
			return
		}
		if binaryNamespaceIfModule != j2gort.NullString {
			var suffix string
			_ = suffix
			if j2gort.JStrLen(prefixTypeName) == j2gort.JStrLen(typeName) {
				suffix = ""
			} else {
				suffix = j2gort.JStrSubFrom(typeName, j2gort.JStrLen(prefixTypeName))
			}
			this_.priv_ClosureRewriteModule_safeSetString_1227(typeRefNode, (j2gort.StrValString(binaryNamespaceIfModule) + j2gort.StrValString(suffix)))
			return
		}
		var lastDot int32 = j2gort.JStrLastIndexOfCp(prefixTypeName, int32(uint16(46)))
		_ = lastDot
		if lastDot >= int32(0) {
			prefixTypeName = j2gort.JStrSub(typeName, int32(0), lastDot)
		} else {
			break
		}
	}
}

// Java method: ClosureRewriteModule.isWeakNamespace(java.lang.String namespaceId)
func (this_ *ClosureRewriteModule) priv_ClosureRewriteModule_isWeakNamespace_1185(namespaceId string) bool {
	_ = namespaceId
	var metadata ModuleMetadataMap_ModuleMetadata = j2gort.AsIface[ModuleMetadataMap_ModuleMetadata](j2gort.MapGetStr(this_.compiler.GetModuleMetadataMap().GetModulesByGoogNamespace(), namespaceId))
	_ = metadata
	return ((metadata != nil && !metadata.IsJavaNil_()) && metadata.RootNode().GetStaticSourceFile().IsWeak())
}

// Java method: ClosureRewriteModule.isUnrequiredWeakNamespace(com.google.javascript.rhino.Node call, java.lang.String namespaceId)
func (this_ *ClosureRewriteModule) priv_ClosureRewriteModule_isUnrequiredWeakNamespace_1186(call pkg_rhino.Node, namespaceId string) bool {
	_ = call
	_ = namespaceId
	if call.GetStaticSourceFile().IsWeak() || !(this_.priv_ClosureRewriteModule_isWeakNamespace_1185(namespaceId)) {
		return false
	}
	var callerMetadata ModuleMetadataMap_ModuleMetadata = j2gort.AsIface[ModuleMetadataMap_ModuleMetadata](j2gort.MapGetStr(this_.compiler.GetModuleMetadataMap().GetModulesByPath(), call.GetSourceFileName()))
	_ = callerMetadata
	return ((callerMetadata == nil || callerMetadata.IsJavaNil_()) || !(j2gort.SetContainsStr(callerMetadata.StronglyRequiredGoogNamespaces(), namespaceId)))
}

// Java method: ClosureRewriteModule.process(com.google.javascript.rhino.Node externs, com.google.javascript.rhino.Node root)
func (this_ *ClosureRewriteModule) Process_2(externs pkg_rhino.Node, root pkg_rhino.Node) {
	_ = externs
	_ = root
	sc_ := NewSyntacticScopeCreator_1(this_.compiler)
	defer func() {
		this_.globalScope = nil
		this_.currentScript = nil
		sc_.ReleaseSlabs_()
	}()

	var scriptDescriptions j2gort.Ext_java_util_Deque = j2gort.NewExt_java_util_ArrayDeque()
	_ = scriptDescriptions
	{
		cur_1 := j2gort.NewIterCursor(j2gort.CheckNotNull(j2gort.Ext_com_google_common_collect_ImmutableList_Of_2(j2gort.BoxNode(externs), j2gort.BoxNode(root))))
		for cur_1.HasNext() {
			parent := pkg_rhino.Node(cur_1.NextNode())
			_ = parent
			{
				var script pkg_rhino.Node = parent.GetFirstChild()
				_ = script
				for ; script != 0; script = script.GetNext() {
					if !(script.IsScript()) {
						j2gort.CheckStateMsg(false, j2gort.BoxNode(script))
					}
					NodeTraversal_Traverse(this_.compiler, script, NewClosureRewriteModule_UnwrapGoogLoadModule(this_))
					this_.priv_ClosureRewriteModule_pushScript_1221(NewClosureRewriteModule_ScriptDescription())
					{
						r_2 := this_.currentScript
						r_2.rootNode = script
					}
					scriptDescriptions.AddLast(j2gort.PtrToAny(this_.currentScript))
					NodeTraversal_Traverse(this_.compiler, script, NewClosureRewriteModule_ScriptPreprocessor(this_))
					{
						clinit_NodeTraversal()
						var ntb_ NodeTraversal_Builder
						InitNodeTraversal_Builder(&ntb_)
						(&ntb_).SetCompiler(this_.compiler)
						(&ntb_).SetCallback_1_NodeTraversal_Callback(NewClosureRewriteModule_ScriptRecorder(this_))
						(&ntb_).SetScopeCreator(sc_)
						(&ntb_).Traverse(script)
					}

					this_.priv_ClosureRewriteModule_popScript_1222()
				}
			}
		}
	}
	this_.priv_ClosureRewriteModule_reportUnrecognizedRequires_1226()
	if this_.compiler.HasHaltingErrors() {
		return
	}
	{
		clinit_NodeTraversal()
		var ntb_1_ NodeTraversal_Builder
		InitNodeTraversal_Builder(&ntb_1_)
		(&ntb_1_).SetCompiler(this_.compiler)
		(&ntb_1_).SetCallback_1_NodeTraversal_Callback(NewClosureRewriteModule_ScriptUpdater(this_, scriptDescriptions))
		(&ntb_1_).SetScopeCreator(sc_)
		(&ntb_1_).TraverseRoots(externs, root)
	}

	this_.priv_ClosureRewriteModule_declareSyntheticExterns_1187()
	this_.googModuleGetCalls.ForEach(func() j2gort.Ext_java_util_function_Consumer {
		bound_3 := this_
		return j2gort.NewLambda_Ext_java_util_function_Consumer(func(mrp0_ any) {
			bound_3.priv_ClosureRewriteModule_updateGoogModuleGetCall_1212(j2gort.CastNode[pkg_rhino.Node](mrp0_))
		})
	}())
	this_.googRequireDynamicCalls.ForEach(func() j2gort.Ext_java_util_function_Consumer {
		bound_4 := this_
		return j2gort.NewLambda_Ext_java_util_function_Consumer(func(mrp0_ any) {
			bound_4.priv_ClosureRewriteModule_updateGoogRequireDynamicCall_1208(j2gort.CastNode[pkg_rhino.Node](mrp0_))
		})
	}())
}

// Java method: ClosureRewriteModule.declareSyntheticExterns()
func (this_ *ClosureRewriteModule) priv_ClosureRewriteModule_declareSyntheticExterns_1187() {
	var vars j2gort.Ext_com_google_common_collect_ImmutableList = j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableList](j2gort.AsIface[j2gort.Ext_java_util_Collection](this_.syntheticExterns.LiveValues()).Stream().Filter(j2gort.NewNodeLambda_Ext_java_util_function_Predicate(func(rawNode_6 uint32) bool {
		lhs := pkg_rhino.Node(rawNode_6)
		_ = lhs
		return !(this_.priv_ClosureRewriteModule_isNameInGlobalScope_1188(lhs.GetString()))
	})).Map(j2gort.NewLambda_Ext_java_util_function_Function(func(lp_5 any) any {
		lhs := j2gort.CastNode[pkg_rhino.Node](lp_5)
		_ = lhs
		return j2gort.BoxNode(pkg_rhino.IR_Var_1(this_.astFactory.CreateNameWithUnknownType(lhs.GetString())).SrcrefTree(lhs))
	})).Collect_1(j2gort.Ext_com_google_common_collect_ImmutableList_ToImmutableList()))
	_ = vars
	if vars.IsEmpty() {
		return
	}
	var root pkg_rhino.Node = this_.compiler.GetSynthesizedExternsInput().GetAstRoot(this_.compiler)
	_ = root
	vars.ForEach(func() j2gort.Ext_java_util_function_Consumer {
		bound_7 := root
		return j2gort.NewLambda_Ext_java_util_function_Consumer(func(mrp0_ any) { bound_7.AddChildToBack(j2gort.CastNode[pkg_rhino.Node](mrp0_)) })
	}())
}

// Java method: ClosureRewriteModule.isNameInGlobalScope(java.lang.String name)
func (this_ *ClosureRewriteModule) priv_ClosureRewriteModule_isNameInGlobalScope_1188(name string) bool {
	_ = name
	return (this_.legacyScriptNamespacesAndPrefixes.Contains(name) || (this_.globalScope.LookupVar(name) != nil))
}

// Java method: ClosureRewriteModule.preprocessExportDeclaration(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n)
func (this_ *ClosureRewriteModule) priv_ClosureRewriteModule_preprocessExportDeclaration_1189(t *NodeTraversal, n pkg_rhino.Node) {
	_ = t
	_ = n
	if (!(ClosureRewriteModule_IsGoogModuleExportsRef(t, n)) || !(ClosureRewriteModule_IsAssignTarget(n))) || !(n.GetGrandparent().IsExprResult()) {
		return
	}
	if !(this_.currentScript.defaultExport == nil) {
		j2gort.CheckState(false)
	}
	var exportRhs pkg_rhino.Node = n.GetNext()
	_ = exportRhs
	if NodeUtil_IsNamedExportsLiteral(exportRhs) {
		var insertionPoint pkg_rhino.Node = n.GetGrandparent()
		_ = insertionPoint
		{
			var key pkg_rhino.Node = exportRhs.GetFirstChild()
			_ = key
			for ; key != 0; key = key.GetNext() {
				var exportName string = key.GetString()
				_ = exportName
				var jsdoc *pkg_rhino.JSDocInfo = key.GetJSDocInfo()
				_ = jsdoc
				var rhs pkg_rhino.Node = key.RemoveFirstChild()
				_ = rhs
				var lhs pkg_rhino.Node = this_.astFactory.CreateGetProp(this_.astFactory.CreateName_2_string_AstFactory_Type("exports", AstFactory_S_Type_1_Node(n)), exportName, AstFactory_S_Type_1_Node(rhs)).SrcrefTree(key)
				_ = lhs
				var newExport pkg_rhino.Node = pkg_rhino.IR_ExprResult(this_.astFactory.CreateAssign_2_Node_Node(lhs, rhs).Srcref(key).SetJSDocInfo(jsdoc)).Srcref(key)
				_ = newExport
				newExport.InsertAfter(insertionPoint)
				insertionPoint = newExport
			}
		}
		n.GetGrandparent().Detach()
	}
}

// Java method: ClosureRewriteModule.recordModuleBody(com.google.javascript.rhino.Node moduleRoot)
func (this_ *ClosureRewriteModule) priv_ClosureRewriteModule_recordModuleBody_1190(moduleRoot pkg_rhino.Node) {
	_ = moduleRoot
	this_.priv_ClosureRewriteModule_pushScript_1221(NewClosureRewriteModule_ScriptDescription())
	{
		r_8 := this_.currentScript
		r_8.rootNode = moduleRoot
	}
	{
		r_9 := this_.currentScript
		r_9.isModule = true
	}
}

// Java method: ClosureRewriteModule.recordGoogModule(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node call)
func (this_ *ClosureRewriteModule) priv_ClosureRewriteModule_recordGoogModule_1191(t *NodeTraversal, call pkg_rhino.Node) {
	_ = t
	_ = call
	var namespaceIdNode pkg_rhino.Node = call.GetLastChild()
	_ = namespaceIdNode
	if !(namespaceIdNode.IsStringLit()) {
		t.Report_3(namespaceIdNode, StaticGet_ClosureRewriteModule_INVALID_MODULE_ID_ARG(), []string{})
		return
	}
	var namespaceId string = namespaceIdNode.GetString()
	_ = namespaceId
	{
		r_10 := this_.currentScript
		r_10.namespaceId = namespaceId
	}
	{
		r_11 := this_.currentScript
		r_11.contentsPrefix = ClosureRewriteModule_ToModuleContentsPrefix(namespaceId)
	}
	var scriptNode pkg_rhino.Node = NodeUtil_GetEnclosingScript(this_.currentScript.rootNode)
	_ = scriptNode
	this_.rewriteState.scriptDescriptionsByGoogModuleNamespace.Put(namespaceId, this_.currentScript)
	j2gort.MultimapPutNode(this_.rewriteState.namespaceIdsByScriptNode, scriptNode, j2gort.BoxString(namespaceId))
}

// Java method: ClosureRewriteModule.recordGoogDeclareLegacyNamespace()
func (this_ *ClosureRewriteModule) priv_ClosureRewriteModule_recordGoogDeclareLegacyNamespace_1192() {
	{
		r_12 := this_.currentScript
		r_12.declareLegacyNamespace = true
	}
	this_.priv_ClosureRewriteModule_updateLegacyScriptNamespacesAndPrefixes_1193(this_.currentScript.namespaceId)
}

// Java method: ClosureRewriteModule.updateLegacyScriptNamespacesAndPrefixes(java.lang.String namespace)
func (this_ *ClosureRewriteModule) priv_ClosureRewriteModule_updateLegacyScriptNamespacesAndPrefixes_1193(namespace string) {
	_ = namespace
	this_.legacyScriptNamespacesAndPrefixes.Add(namespace)
	{
		var dot int32 = j2gort.JStrLastIndexOfCp(namespace, int32(uint16(46)))
		_ = dot
		for ; dot != int32(-1); dot = j2gort.JStrLastIndexOfCp(namespace, int32(uint16(46))) {
			namespace = j2gort.JStrSub(namespace, int32(0), dot)
			this_.legacyScriptNamespacesAndPrefixes.Add(namespace)
		}
	}
}

// Java method: ClosureRewriteModule.recordGoogProvide(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node call)
func (this_ *ClosureRewriteModule) priv_ClosureRewriteModule_recordGoogProvide_1194(t *NodeTraversal, call pkg_rhino.Node) {
	_ = t
	_ = call
	var namespaceIdNode pkg_rhino.Node = call.GetLastChild()
	_ = namespaceIdNode
	if !(namespaceIdNode.IsStringLit()) {
		t.Report_3(namespaceIdNode, StaticGet_ClosureRewriteModule_INVALID_PROVIDE_NAMESPACE(), []string{})
		return
	}
	var namespaceId string = namespaceIdNode.GetString()
	_ = namespaceId
	if this_.currentScript.isModule {
		t.Report_3(namespaceIdNode, StaticGet_ClosureRewriteModule_INVALID_PROVIDE_CALL(), []string{})
	}
	var scriptNode pkg_rhino.Node = NodeUtil_GetEnclosingScript(call)
	_ = scriptNode
	this_.rewriteState.providedNamespaces.Add(namespaceId)
	j2gort.MultimapPutNode(this_.rewriteState.namespaceIdsByScriptNode, scriptNode, j2gort.BoxString(namespaceId))
	this_.priv_ClosureRewriteModule_updateLegacyScriptNamespacesAndPrefixes_1193(namespaceId)
}

// Java method: ClosureRewriteModule.recordGoogRequire(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node call)
func (this_ *ClosureRewriteModule) priv_ClosureRewriteModule_recordGoogRequire_1195(t *NodeTraversal, call pkg_rhino.Node) {
	_ = t
	_ = call
	ClosureRewriteModule_MaybeSplitMultiVar(call)
	var namespaceIdNode pkg_rhino.Node = call.GetLastChild()
	_ = namespaceIdNode
	if !(namespaceIdNode.IsStringLit()) {
		t.Report_3(namespaceIdNode, StaticGet_ClosurePrimitiveErrors_INVALID_REQUIRE_NAMESPACE(), []string{})
		return
	}
	var namespaceId string = namespaceIdNode.GetString()
	_ = namespaceId
	var targetIsAModule bool = this_.rewriteState.ContainsModule(namespaceId)
	_ = targetIsAModule
	var targetIsALegacyScript bool = this_.rewriteState.providedNamespaces.Contains(namespaceId)
	_ = targetIsALegacyScript
	if (this_.currentScript.isModule && !(targetIsAModule)) && !(targetIsALegacyScript) {
		this_.unrecognizedRequires = append(this_.unrecognizedRequires, NewClosureRewriteModule_UnrecognizedRequire(call, namespaceId))
	}
}

// Java method: ClosureRewriteModule.recordGoogRequireType(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node call)
func (this_ *ClosureRewriteModule) priv_ClosureRewriteModule_recordGoogRequireType_1196(t *NodeTraversal, call pkg_rhino.Node) {
	_ = t
	_ = call
	var namespaceIdNode pkg_rhino.Node = call.GetLastChild()
	_ = namespaceIdNode
	if !(namespaceIdNode.IsStringLit()) {
		t.Report_3(namespaceIdNode, StaticGet_ClosurePrimitiveErrors_INVALID_REQUIRE_TYPE_NAMESPACE(), []string{})
		return
	}
	this_.priv_ClosureRewriteModule_recordGoogRequire_1195(t, call)
}

// Java method: ClosureRewriteModule.recordGoogForwardDeclare(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node call)
func (this_ *ClosureRewriteModule) priv_ClosureRewriteModule_recordGoogForwardDeclare_1197(t *NodeTraversal, call pkg_rhino.Node) {
	_ = t
	_ = call
	var namespaceNode pkg_rhino.Node = call.GetLastChild()
	_ = namespaceNode
	if !(call.HasTwoChildren()) || !(namespaceNode.IsStringLit()) {
		t.Report_3(namespaceNode, StaticGet_ClosurePrimitiveErrors_INVALID_FORWARD_DECLARE_NAMESPACE(), []string{})
		return
	}
	this_.priv_ClosureRewriteModule_recordGoogRequire_1195(t, call)
}

// Java method: ClosureRewriteModule.recordGoogRequireDynamic(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node call)
func (this_ *ClosureRewriteModule) priv_ClosureRewriteModule_recordGoogRequireDynamic_1198(t *NodeTraversal, call pkg_rhino.Node) {
	_ = t
	_ = call
	var namespaceIdNode pkg_rhino.Node = call.GetLastChild()
	_ = namespaceIdNode
	if !(namespaceIdNode.IsStringLit()) {
		t.Report_3(namespaceIdNode, StaticGet_ClosurePrimitiveErrors_INVALID_REQUIRE_DYNAMIC(), []string{})
		return
	}
	var namespaceId string = namespaceIdNode.GetString()
	_ = namespaceId
	if !(this_.rewriteState.ContainsModule(namespaceId)) {
		this_.unrecognizedRequires = append(this_.unrecognizedRequires, NewClosureRewriteModule_UnrecognizedRequire(call, namespaceId))
	}
	this_.googRequireDynamicCalls.Add_1(j2gort.BoxNode(call))
}

// Java method: ClosureRewriteModule.recordGoogModuleGet(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node call)
func (this_ *ClosureRewriteModule) priv_ClosureRewriteModule_recordGoogModuleGet_1199(t *NodeTraversal, call pkg_rhino.Node) {
	_ = t
	_ = call
	var namespaceIdNode pkg_rhino.Node = call.GetLastChild()
	_ = namespaceIdNode
	if !(call.HasTwoChildren()) || !(namespaceIdNode.IsStringLit()) {
		t.Report_3(namespaceIdNode, StaticGet_ClosurePrimitiveErrors_INVALID_GET_NAMESPACE(), []string{})
		return
	}
	var namespaceId string = namespaceIdNode.GetString()
	_ = namespaceId
	if !(this_.rewriteState.ContainsModule(namespaceId)) {
		this_.unrecognizedRequires = append(this_.unrecognizedRequires, NewClosureRewriteModule_UnrecognizedRequire(call, namespaceId))
	}
	this_.googModuleGetCalls.Add_1(j2gort.BoxNode(call))
	var maybeAssign pkg_rhino.Node = call.GetParent()
	_ = maybeAssign
	var isFillingAnAlias bool = ((maybeAssign.IsAssign() && maybeAssign.GetFirstChild().IsName()) && maybeAssign.GetParent().IsExprResult())
	_ = isFillingAnAlias
	if !(isFillingAnAlias) || !(this_.currentScript.isModule) {
		return
	}
	var aliasName string = call.GetParent().GetFirstChild().GetString()
	_ = aliasName
	var aliasVar *Var = j2gort.CastPtr[Var](t.GetScope().GetVar(aliasName))
	_ = aliasVar
	if aliasVar == nil {
		return
	}
	var aliasVarNodeRhs pkg_rhino.Node = NodeUtil_GetRValueOfLValue(aliasVar.GetNode())
	_ = aliasVarNodeRhs
	if ((aliasVarNodeRhs == 0) || !(NodeUtil_IsCallTo_2_Node_Node(aliasVarNodeRhs, StaticGet_ClosureRewriteModule_GOOG_FORWARDDECLARE()))) || !(func() bool {
		sr_13 := namespaceId
		se_14 := aliasVarNodeRhs.GetLastChild().GetString()
		return j2gort.CheckNotNullString(sr_13) == se_14
	}()) {
		return
	}
	if this_.priv_ClosureRewriteModule_isUnrequiredWeakNamespace_1186(call, namespaceId) {
		return
	}
	this_.compiler.ReportChangeToEnclosingScope(maybeAssign)
	maybeAssign.GetParent().Detach()
	this_.googModuleGetCalls.Remove_1_int32((this_.googModuleGetCalls.Size() - int32(1)))
}

// Java method: ClosureRewriteModule.recordTopLevelClassOrFunctionName(com.google.javascript.rhino.Node classOrFunctionNode)
func (this_ *ClosureRewriteModule) priv_ClosureRewriteModule_recordTopLevelClassOrFunctionName_1200(classOrFunctionNode pkg_rhino.Node) {
	_ = classOrFunctionNode
	var nameNode pkg_rhino.Node = classOrFunctionNode.GetFirstChild()
	_ = nameNode
	if nameNode.IsName() && !(j2gort.Ext_com_google_common_base_Strings_IsNullOrEmpty(nameNode.GetString())) {
		var name string = nameNode.GetString()
		_ = name
		this_.currentScript.topLevelNames.Add(name)
	}
}

// Java method: ClosureRewriteModule.recordTopLevelVarNames(com.google.javascript.rhino.Node varNode)
func (this_ *ClosureRewriteModule) priv_ClosureRewriteModule_recordTopLevelVarNames_1201(varNode pkg_rhino.Node) {
	_ = varNode
	NodeUtil_VisitLhsNodesInNode(varNode, j2gort.NewLambda_Ext_java_util_function_Consumer(func(lp_15 any) {
		lhs := j2gort.CastNode[pkg_rhino.Node](lp_15)
		_ = lhs
		this_.currentScript.topLevelNames.Add(lhs.GetString())
	}))
}

// Java method: ClosureRewriteModule.maybeRecordExportDeclaration(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n)
func (this_ *ClosureRewriteModule) priv_ClosureRewriteModule_maybeRecordExportDeclaration_1202(t *NodeTraversal, n pkg_rhino.Node) {
	_ = t
	_ = n
	if (!(this_.currentScript.isModule) || !(ClosureRewriteModule_IsGoogModuleExportsRef(t, n))) || !(ClosureRewriteModule_IsAssignTarget(n)) {
		return
	}
	if this_.currentScript.defaultExport != nil {
		var previousExport *ClosureRewriteModule_ExportDefinition = this_.currentScript.defaultExport
		_ = previousExport
		var localName string = previousExport.GetLocalName()
		_ = localName
		if (localName != j2gort.NullString) && this_.currentScript.namesToInlineByAlias.ContainsKey(localName) {
			this_.currentScript.namesToInlineByAlias.Remove(localName)
		}
		{
			r_16 := this_.currentScript
			r_16.defaultExportLocalName = j2gort.NullString
		}
	}
	var exportRhs pkg_rhino.Node = n.GetNext()
	_ = exportRhs
	if !(!(NodeUtil_IsNamedExportsLiteral(exportRhs))) {
		j2gort.CheckStateStr(false, "Exports object should have been converted already")
	}
	{
		r_17 := this_.currentScript
		r_17.willCreateExportsObject = true
	}
	var defaultExport *ClosureRewriteModule_ExportDefinition = ClosureRewriteModule_ExportDefinition_NewDefaultExport(t, exportRhs)
	_ = defaultExport
	{
		r_18 := this_.currentScript
		r_18.defaultExport = defaultExport
	}
	if !(this_.currentScript.declareLegacyNamespace) && defaultExport.HasInlinableName(j2gort.AsIface[j2gort.Ext_java_util_Set](this_.currentScript.exportsToInline.LiveKeySet())) {
		var localName_2 string = defaultExport.GetLocalName()
		_ = localName_2
		{
			r_19 := this_.currentScript
			r_19.defaultExportLocalName = localName_2
		}
		this_.priv_ClosureRewriteModule_recordExportToInline_1224(defaultExport)
	}
}

// Java method: ClosureRewriteModule.updateModuleBodyEarly(com.google.javascript.rhino.Node moduleScopeRoot)
func (this_ *ClosureRewriteModule) priv_ClosureRewriteModule_updateModuleBodyEarly_1203(moduleScopeRoot pkg_rhino.Node) {
	_ = moduleScopeRoot
	this_.priv_ClosureRewriteModule_pushScript_1221(this_.currentScript.RemoveFirstChildScript())
	{
		r_20 := this_.currentScript
		r_20.rootNode = moduleScopeRoot
	}
}

// Java method: ClosureRewriteModule.updateGoogModule(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node call)
func (this_ *ClosureRewriteModule) priv_ClosureRewriteModule_updateGoogModule_1204(t *NodeTraversal, call pkg_rhino.Node) {
	_ = t
	_ = call
	if !(this_.currentScript.isModule) {
		this_.compiler.ReportChangeToEnclosingScope(call)
		var undefined pkg_rhino.Node = this_.astFactory.CreateVoid(this_.astFactory.CreateNumber(float64(int32(0)))).SrcrefTree(call)
		_ = undefined
		call.ReplaceWith(undefined)
		return
	}
	if this_.currentScript.declareLegacyNamespace {
		call.GetFirstChild().SetString("provide")
		this_.compiler.ReportChangeToEnclosingScope(call)
	}
	if !(this_.currentScript.willCreateExportsObject) {
		if !(!(this_.currentScript.hasCreatedExportObject)) {
			j2gort.CheckStateMsg(false, j2gort.PtrToAny(this_.currentScript))
		}
		this_.priv_ClosureRewriteModule_exportTheEmptyBinaryNamespaceAt_1223(NodeUtil_GetEnclosingStatement(call), StaticGet_ClosureRewriteModule_AddAt_AFTER(), t)
	}
	if !(this_.currentScript.declareLegacyNamespace) && !(this_.preserveSugar) {
		this_.compiler.ReportChangeToEnclosingScope(call)
		NodeUtil_GetEnclosingStatement(call).Detach()
	}
}

// Java static method: ClosureRewriteModule.updateGoogDeclareLegacyNamespace(com.google.javascript.rhino.Node call)
func ClosureRewriteModule_UpdateGoogDeclareLegacyNamespace(call pkg_rhino.Node) {
	clinit_ClosureRewriteModule()
	_ = call
	NodeUtil_GetEnclosingStatement(call).Detach()
}

// Java static method: ClosureRewriteModule.updateGoogPreventModuleExportsSealing(com.google.javascript.rhino.Node call)
func ClosureRewriteModule_UpdateGoogPreventModuleExportsSealing(call pkg_rhino.Node) {
	clinit_ClosureRewriteModule()
	_ = call
	NodeUtil_GetEnclosingStatement(call).Detach()
}

// Java method: ClosureRewriteModule.updateGoogRequire(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node call)
func (this_ *ClosureRewriteModule) priv_ClosureRewriteModule_updateGoogRequire_1205(t *NodeTraversal, call pkg_rhino.Node) {
	_ = t
	_ = call
	var namespaceIdNode pkg_rhino.Node = call.GetLastChild()
	_ = namespaceIdNode
	var statementNode pkg_rhino.Node = NodeUtil_GetEnclosingStatement(call)
	_ = statementNode
	var namespaceId string = namespaceIdNode.GetString()
	_ = namespaceId
	var targetIsNonLegacyGoogModule bool = (this_.rewriteState.ContainsModule(namespaceId) && !(this_.rewriteState.IsLegacyModule(namespaceId)))
	_ = targetIsNonLegacyGoogModule
	var importHasAlias bool = NodeUtil_IsNameDeclaration(statementNode)
	_ = importHasAlias
	var isDestructuring bool = statementNode.GetFirstChild().IsDestructuringLhs()
	_ = isDestructuring
	var currentScriptIsAModule bool = this_.currentScript.isModule
	_ = currentScriptIsAModule
	var requireDirectlyStoredInAlias bool = NodeUtil_IsNameDeclaration(call.GetGrandparent())
	_ = requireDirectlyStoredInAlias
	if (currentScriptIsAModule && requireDirectlyStoredInAlias) && this_.priv_ClosureRewriteModule_isTopLevel_1229(t, statementNode, StaticGet_ClosureRewriteModule_ScopeType_EXEC_CONTEXT()) {
		var lhs pkg_rhino.Node = call.GetParent()
		_ = lhs
		var exportedNamespace string = this_.rewriteState.priv_GlobalRewriteState_getExportedNamespaceOrScript_1244(namespaceId)
		_ = exportedNamespace
		if exportedNamespace == j2gort.NullString {
		} else if lhs.IsName() {
			var aliasName string = statementNode.GetFirstChild().GetString()
			_ = aliasName
			this_.priv_ClosureRewriteModule_recordNameToInline_1225(aliasName, exportedNamespace, namespaceId)
			this_.priv_ClosureRewriteModule_maybeAddAliasToSymbolTable_1230(statementNode.GetFirstChild(), this_.currentScript.namespaceId)
		} else if lhs.IsDestructuringLhs() && lhs.GetFirstChild().IsObjectPattern() {
			this_.priv_ClosureRewriteModule_maybeWarnForInvalidDestructuring_1206(t, lhs.GetParent(), namespaceId)
			{
				var importSpec pkg_rhino.Node = lhs.GetFirstFirstChild()
				_ = importSpec
				for ; importSpec != 0; importSpec = importSpec.GetNext() {
					if !(importSpec.HasChildren()) {
						j2gort.CheckStateMsg(false, j2gort.BoxNode(importSpec))
					}
					var importedProperty string = importSpec.GetString()
					_ = importedProperty
					var aliasNode pkg_rhino.Node = importSpec.GetFirstChild()
					_ = aliasNode
					var aliasName_2 string = aliasNode.GetString()
					_ = aliasName_2
					var fullName string = (j2gort.StrValString(exportedNamespace) + "." + j2gort.StrValString(importedProperty))
					_ = fullName
					this_.priv_ClosureRewriteModule_recordNameToInline_1225(aliasName_2, fullName, j2gort.NullString)
					this_.priv_ClosureRewriteModule_maybeAddAliasToSymbolTable_1230(aliasNode, this_.currentScript.namespaceId)
					this_.priv_ClosureRewriteModule_safeSetString_1227(aliasNode, (j2gort.StrValString(this_.currentScript.contentsPrefix) + j2gort.StrValString(aliasName_2)))
				}
			}
		} else {
			panic(j2gort.AsIface[javalang.IRuntimeException](j2gort.ConstructShimThrowable("java.lang.RuntimeException", "Illegal goog.module import: "+j2gort.StrVal(lhs))))
		}
	}
	if this_.currentScript.isModule || targetIsNonLegacyGoogModule {
		if isDestructuring {
			if !(this_.preserveSugar) {
				this_.compiler.ReportChangeToEnclosingScope(statementNode)
				statementNode.Detach()
			}
		} else if targetIsNonLegacyGoogModule {
			if !(this_.priv_ClosureRewriteModule_isTopLevel_1229(t, statementNode, StaticGet_ClosureRewriteModule_ScopeType_EXEC_CONTEXT())) {
				j2gort.CheckStateFmt(false, "Unexpected non-top-level require at %s", j2gort.BoxNode(call))
			}
			if importHasAlias || !(this_.rewriteState.IsLegacyModule(namespaceId)) {
				if !(this_.preserveSugar) {
					this_.compiler.ReportChangeToEnclosingScope(statementNode)
					statementNode.Detach()
				}
			}
		} else {
			call.Detach()
			statementNode.ReplaceWith(pkg_rhino.IR_ExprResult(call))
			this_.compiler.ReportChangeToEnclosingScope(call)
		}
	}
}

// Java method: ClosureRewriteModule.maybeWarnForInvalidDestructuring(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node importNode, java.lang.String importedNamespace)
func (this_ *ClosureRewriteModule) priv_ClosureRewriteModule_maybeWarnForInvalidDestructuring_1206(t *NodeTraversal, importNode pkg_rhino.Node, importedNamespace string) {
	_ = t
	_ = importNode
	_ = importedNamespace
	if !(importNode.GetFirstChild().IsDestructuringLhs()) {
		j2gort.CheckArgumentMsg(false, j2gort.BoxNode(importNode))
	}
	var importedModule *ClosureRewriteModule_ScriptDescription = this_.rewriteState.scriptDescriptionsByGoogModuleNamespace.Get(importedNamespace)
	_ = importedModule
	if importedModule == nil {
		return
	}
	if importedModule.defaultExport != nil {
		t.Report_3(importNode, StaticGet_ClosureRewriteModule_ILLEGAL_DESTRUCTURING_DEFAULT_EXPORT(), []string{})
		return
	}
	var objPattern pkg_rhino.Node = importNode.GetFirstFirstChild()
	_ = objPattern
	{
		var key pkg_rhino.Node = objPattern.GetFirstChild()
		_ = key
		for ; key != 0; key = key.GetNext() {
			var exportName string = key.GetString()
			_ = exportName
			if !(importedModule.namedExports.Contains(exportName)) {
				t.Report_3(importNode, StaticGet_ClosureRewriteModule_ILLEGAL_DESTRUCTURING_NOT_EXPORTED(), []string{exportName, importedNamespace})
			}
		}
	}
}

// Java method: ClosureRewriteModule.updateGoogForwardDeclare(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node call)
func (this_ *ClosureRewriteModule) priv_ClosureRewriteModule_updateGoogForwardDeclare_1207(t *NodeTraversal, call pkg_rhino.Node) {
	_ = t
	_ = call
	var namespaceId string = call.GetLastChild().GetString()
	_ = namespaceId
	if this_.priv_ClosureRewriteModule_isUnrequiredWeakNamespace_1186(call, namespaceId) {
		this_.compiler.ReportChangeToEnclosingScope(call)
		call.ReplaceWith(this_.astFactory.CreateNull().Srcref(call))
		return
	}
	this_.priv_ClosureRewriteModule_updateGoogRequire_1205(t, call)
}

// Java method: ClosureRewriteModule.updateGoogRequireDynamicCall(com.google.javascript.rhino.Node call)
func (this_ *ClosureRewriteModule) priv_ClosureRewriteModule_updateGoogRequireDynamicCall_1208(call pkg_rhino.Node) {
	_ = call
	var parent pkg_rhino.Node = call.GetParent()
	_ = parent
	if parent == 0 {
		return
	}
	if !(parent.IsAwait() || (parent.IsGetProp() && (j2gort.CheckNotNullString(parent.GetString()) == "then"))) {
		j2gort.CheckStateStr(false, "goog.requireDynamic() in only allowed in await/then expression")
	}
	if parent.IsAwait() {
		this_.priv_ClosureRewriteModule_updateGoogRequireDynamicCallInAwait_1210(call)
	} else {
		this_.priv_ClosureRewriteModule_updateGoogRequireDynamicCallInThen_1209(call)
	}
}

// Java method: ClosureRewriteModule.updateGoogRequireDynamicCallInThen(com.google.javascript.rhino.Node call)
func (this_ *ClosureRewriteModule) priv_ClosureRewriteModule_updateGoogRequireDynamicCallInThen_1209(call pkg_rhino.Node) {
	_ = call
	var namespaceIdNode pkg_rhino.Node = call.GetSecondChild()
	_ = namespaceIdNode
	var namespaceId string = namespaceIdNode.GetString()
	_ = namespaceId
	var exportedNamespace string = this_.rewriteState.priv_GlobalRewriteState_getExportedNamespaceOrScript_1244(namespaceId)
	_ = exportedNamespace
	if !(exportedNamespace != j2gort.NullString) {
		j2gort.CheckStateStr(false, "Exported namespace for goog.requireDynamic() canot be null")
	}
	var getProp pkg_rhino.Node = call.GetParent()
	_ = getProp
	var thenCallNode pkg_rhino.Node = getProp.GetParent()
	_ = thenCallNode
	if !((thenCallNode != 0) && thenCallNode.IsCall()) {
		j2gort.CheckStateStr(false, "must be a 'then' call expression")
	}
	var argNode pkg_rhino.Node = this_.astFactory.CreateString(this_.priv_ClosureRewriteModule_namespaceIdToXid_1211(namespaceId)).Srcref(namespaceIdNode)
	_ = argNode
	namespaceIdNode.ReplaceWith(argNode)
	var calleeNode pkg_rhino.Node = call.GetFirstChild()
	_ = calleeNode
	if !(calleeNode.MatchesQualifiedName_1_string("goog.requireDynamic")) {
		j2gort.CheckStateMsg(false, j2gort.BoxNode(calleeNode))
	}
	calleeNode.SetString("importHandler_")
	var exportedNamespaceNameNode pkg_rhino.Node = this_.astFactory.CreateName_2_StaticScope_string(pkg_rhino.ToIface_StaticScope(this_.globalTypedScope), exportedNamespace).SrcrefTree(call)
	_ = exportedNamespaceNameNode
	exportedNamespaceNameNode.SetOriginalName(namespaceId)
	var functionNode pkg_rhino.Node = thenCallNode.GetSecondChild()
	_ = functionNode
	if !((functionNode != 0) && functionNode.IsFunction()) {
		j2gort.CheckStateStr(false, "must be a function in `then`")
	}
	var paramListNode pkg_rhino.Node = functionNode.GetSecondChild()
	_ = paramListNode
	if !(paramListNode.GetChildCount() == int32(1)) {
		j2gort.CheckStateStr(false, "function must have only one parameter")
	}
	var functionBody pkg_rhino.Node = paramListNode.GetNext()
	_ = functionBody
	var objectPatternOrNameNode pkg_rhino.Node = paramListNode.GetOnlyChild()
	_ = objectPatternOrNameNode
	if !(objectPatternOrNameNode.IsObjectPattern() || objectPatternOrNameNode.IsName()) {
		j2gort.CheckStateStr(false, "parameter of callback function must be object pattern or name")
	}
	objectPatternOrNameNode.Detach()
	paramListNode.RemoveChildren()
	var languageIn CompilerOptions_LanguageMode = this_.compiler.GetOptions().GetLanguageIn()
	_ = languageIn
	if !(languageIn.ToFeatureSet().Contains_1_FeatureSet_Feature(pkg_parser.StaticGet_FeatureSet_Feature_CONST_DECLARATIONS())) {
		j2gort.CheckStateFmt(false, "'%s' does not contain '%s'", j2gort.BoxNullableVal(languageIn), j2gort.PtrToAny(pkg_parser.StaticGet_FeatureSet_Feature_CONST_DECLARATIONS()))
	}
	var enclosingScript pkg_rhino.Node = NodeUtil_GetEnclosingScript(call)
	_ = enclosingScript
	NodeUtil_AddFeatureToScript(enclosingScript, pkg_parser.StaticGet_FeatureSet_Feature_CONST_DECLARATIONS(), this_.compiler)
	var declarationNode pkg_rhino.Node = pkg_rhino.Node(0)
	_ = declarationNode
	if objectPatternOrNameNode.IsObjectPattern() {
		declarationNode = this_.astFactory.CreateSingleConstObjectPatternDeclaration(objectPatternOrNameNode, exportedNamespaceNameNode).SrcrefTreeIfMissing(call)
	} else {
		declarationNode = this_.astFactory.CreateSingleConstNameDeclaration(objectPatternOrNameNode.GetString(), exportedNamespaceNameNode).SrcrefTreeIfMissing(call)
	}
	if functionBody.IsBlock() {
		functionBody.AddChildToFront(declarationNode)
		this_.compiler.ReportChangeToEnclosingScope(declarationNode)
	} else {
		var returnStmt pkg_rhino.Node = this_.astFactory.CreateReturn(functionBody.Detach())
		_ = returnStmt
		var newBlock pkg_rhino.Node = this_.astFactory.CreateBlock([]pkg_rhino.Node{returnStmt}).SrcrefTree(call)
		_ = newBlock
		newBlock.InsertAfter(paramListNode)
		newBlock.AddChildToFront(declarationNode)
		this_.compiler.ReportChangeToEnclosingScope(newBlock)
	}
	this_.compiler.ReportChangeToEnclosingScope(call)
}

// Java method: ClosureRewriteModule.updateGoogRequireDynamicCallInAwait(com.google.javascript.rhino.Node call)
func (this_ *ClosureRewriteModule) priv_ClosureRewriteModule_updateGoogRequireDynamicCallInAwait_1210(call pkg_rhino.Node) {
	_ = call
	var namespaceIdNode pkg_rhino.Node = call.GetSecondChild()
	_ = namespaceIdNode
	var namespaceId string = namespaceIdNode.GetString()
	_ = namespaceId
	var exportedNamespace string = this_.rewriteState.priv_GlobalRewriteState_getExportedNamespaceOrScript_1244(namespaceId)
	_ = exportedNamespace
	if !(exportedNamespace != j2gort.NullString) {
		j2gort.CheckStateStr(false, "Exported namespace for goog.requireDynamic() cannot be null")
	}
	var existingAwait pkg_rhino.Node = call.GetParent()
	_ = existingAwait
	if !(existingAwait.IsAwait()) {
		j2gort.CheckStateStr(false, "Only goog.requireDynamic() in await expression is supported now")
	}
	var awaitParent pkg_rhino.Node = existingAwait.GetParent()
	_ = awaitParent
	if (awaitParent == 0) || (!(awaitParent.IsDestructuringLhs()) && !(awaitParent.IsName())) {
		this_.compiler.Report_1(JSError_Make_3(call, StaticGet_ClosureRewriteModule_ILLEGAL_STMT_OF_GOOG_REQUIRE_DYNAMIC_IN_AWAIT(), []string{}))
	}
	var declarationStatement pkg_rhino.Node = awaitParent.GetParent()
	_ = declarationStatement
	if !(NodeUtil_IsNameDeclaration(declarationStatement)) || !(NodeUtil_IsStatement(declarationStatement)) {
		this_.compiler.Report_1(JSError_Make_3(call, StaticGet_ClosureRewriteModule_ILLEGAL_STMT_OF_GOOG_REQUIRE_DYNAMIC_IN_AWAIT(), []string{}))
	}
	var exportedNamespaceNameNode pkg_rhino.Node = this_.astFactory.CreateQNameUsingJSTypeInfo(this_.globalTypedScope, exportedNamespace).SrcrefTree(call)
	_ = exportedNamespaceNameNode
	exportedNamespaceNameNode.SetJSType(this_.rewriteState.GetGoogModuleNamespaceType(namespaceId))
	exportedNamespaceNameNode.SetOriginalName(namespaceId)
	var argNode pkg_rhino.Node = this_.astFactory.CreateString(this_.priv_ClosureRewriteModule_namespaceIdToXid_1211(namespaceId)).Srcref(namespaceIdNode)
	_ = argNode
	namespaceIdNode.ReplaceWith(argNode)
	var calleeNode pkg_rhino.Node = call.GetFirstChild()
	_ = calleeNode
	if !(calleeNode.MatchesQualifiedName_1_string("goog.requireDynamic")) {
		j2gort.CheckStateMsg(false, j2gort.BoxNode(calleeNode))
	}
	calleeNode.SetString("importHandler_")
	existingAwait.ReplaceWith(exportedNamespaceNameNode)
	var awaitStatement pkg_rhino.Node = this_.astFactory.ExprResult(existingAwait).Srcref(existingAwait)
	_ = awaitStatement
	awaitStatement.InsertBefore(declarationStatement)
	this_.compiler.ReportChangeToEnclosingScope(call)
}

// Java method: ClosureRewriteModule.namespaceIdToXid(java.lang.String namespaceId)
func (this_ *ClosureRewriteModule) priv_ClosureRewriteModule_namespaceIdToXid_1211(namespaceId string) string {
	_ = namespaceId
	var hashFunction Xid_HashFunction = this_.compiler.GetOptions().GetChunkIdHashFunction()
	_ = hashFunction
	var xid *Xid
	_ = xid
	if hashFunction == nil || hashFunction.IsJavaNil_() {
		xid = NewXid()
	} else {
		xid = NewXid_1(hashFunction)
	}
	return xid.Get(namespaceId)
}

// Java method: ClosureRewriteModule.updateGoogModuleGetCall(com.google.javascript.rhino.Node call)
func (this_ *ClosureRewriteModule) priv_ClosureRewriteModule_updateGoogModuleGetCall_1212(call pkg_rhino.Node) {
	_ = call
	var namespaceIdNode pkg_rhino.Node = call.GetSecondChild()
	_ = namespaceIdNode
	var namespaceId string = namespaceIdNode.GetString()
	_ = namespaceId
	if this_.priv_ClosureRewriteModule_isUnrequiredWeakNamespace_1186(call, namespaceId) {
		this_.compiler.Report_1(JSError_Make_3(call, StaticGet_ClosurePrimitiveErrors_GOOG_MODULE_GET_OF_WEAK_MODULE(), []string{namespaceId}))
		this_.compiler.ReportChangeToEnclosingScope(call)
		var nullNode pkg_rhino.Node = this_.astFactory.CreateNull().Srcref(call)
		_ = nullNode
		call.ReplaceWith(nullNode)
		return
	}
	var exportedNamespace string = this_.rewriteState.priv_GlobalRewriteState_getExportedNamespaceOrScript_1244(namespaceId)
	_ = exportedNamespace
	if exportedNamespace != j2gort.NullString {
		this_.compiler.ReportChangeToEnclosingScope(call)
		var exportedNamespaceName pkg_rhino.Node = this_.astFactory.CreateQNameUsingJSTypeInfo(this_.globalTypedScope, exportedNamespace).SrcrefTree(call)
		_ = exportedNamespaceName
		exportedNamespaceName.SetJSType(this_.rewriteState.GetGoogModuleNamespaceType(namespaceId))
		exportedNamespaceName.SetOriginalName(namespaceId)
		call.ReplaceWith(exportedNamespaceName)
	}
}

// Java method: ClosureRewriteModule.recordExportsPropertyAssignment(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node getpropNode)
func (this_ *ClosureRewriteModule) priv_ClosureRewriteModule_recordExportsPropertyAssignment_1213(t *NodeTraversal, getpropNode pkg_rhino.Node) {
	_ = t
	_ = getpropNode
	if !(this_.currentScript.isModule) {
		return
	}
	var parent pkg_rhino.Node = getpropNode.GetParent()
	_ = parent
	if !(parent.IsAssign() || parent.IsExprResult()) {
		j2gort.CheckStateMsg(false, j2gort.BoxNode(parent))
	}
	var exportsNameNode pkg_rhino.Node = getpropNode.GetFirstChild()
	_ = exportsNameNode
	if !(j2gort.CheckNotNullString(exportsNameNode.GetString()) == "exports") {
		j2gort.CheckStateMsg(false, j2gort.BoxNode(exportsNameNode))
	}
	if t.InModuleScope() {
		var exportName string = getpropNode.GetString()
		_ = exportName
		this_.currentScript.namedExports.Add(exportName)
		var exportRhs pkg_rhino.Node = getpropNode.GetNext()
		_ = exportRhs
		var namedExport *ClosureRewriteModule_ExportDefinition = ClosureRewriteModule_ExportDefinition_NewNamedExport(t, exportName, exportRhs)
		_ = namedExport
		if (!(this_.currentScript.declareLegacyNamespace) && (this_.currentScript.defaultExport == nil)) && namedExport.HasInlinableName(j2gort.AsIface[j2gort.Ext_java_util_Set](this_.currentScript.exportsToInline.LiveKeySet())) {
			this_.priv_ClosureRewriteModule_recordExportToInline_1224(namedExport)
			parent.GetParent().Detach()
		}
	}
}

// Java method: ClosureRewriteModule.updateExportsPropertyAssignment(com.google.javascript.rhino.Node getpropNode, com.google.javascript.jscomp.NodeTraversal t)
func (this_ *ClosureRewriteModule) priv_ClosureRewriteModule_updateExportsPropertyAssignment_1214(getpropNode pkg_rhino.Node, t *NodeTraversal) {
	_ = getpropNode
	_ = t
	if !(this_.currentScript.isModule) {
		return
	}
	var parent pkg_rhino.Node = getpropNode.GetParent()
	_ = parent
	if !(parent.IsAssign() || parent.IsExprResult()) {
		j2gort.CheckStateMsg(false, j2gort.BoxNode(parent))
	}
	var exportsNameNode pkg_rhino.Node = getpropNode.GetFirstChild()
	_ = exportsNameNode
	if !(j2gort.CheckNotNullString(exportsNameNode.GetString()) == "exports") {
		j2gort.CheckState(false)
	}
	var exportedNamespace string = this_.currentScript.GetExportedNamespace()
	_ = exportedNamespace
	this_.priv_ClosureRewriteModule_safeSetMaybeQualifiedString_1228(exportsNameNode, exportedNamespace, false)
	var jsdocNode pkg_rhino.Node
	_ = jsdocNode
	if parent.IsAssign() {
		jsdocNode = parent
	} else {
		jsdocNode = getpropNode
	}
	ClosureRewriteModule_MarkConstAndCopyJsDoc(jsdocNode, jsdocNode)
	if !(this_.currentScript.hasCreatedExportObject) {
		this_.priv_ClosureRewriteModule_exportTheEmptyBinaryNamespaceAt_1223(NodeUtil_GetEnclosingStatement(parent), StaticGet_ClosureRewriteModule_AddAt_BEFORE(), t)
	}
}

// Java method: ClosureRewriteModule.maybeUpdateTopLevelName(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node nameNode)
func (this_ *ClosureRewriteModule) priv_ClosureRewriteModule_maybeUpdateTopLevelName_1215(t *NodeTraversal, nameNode pkg_rhino.Node) {
	_ = t
	_ = nameNode
	var name string = nameNode.GetString()
	_ = name
	if !(this_.currentScript.isModule) || !(this_.currentScript.topLevelNames.Contains(name)) {
		return
	}
	var var_ *Var = t.GetScope().LookupVar(name)
	_ = var_
	if (var_ == nil) || (j2gort.CastPtr[Scope](var_.GetScope()).GetRootNode_0() != this_.currentScript.rootNode) {
		return
	}
	if ((var_.GetNameNode() == nameNode) && nameNode.GetParent().IsStringKey()) && nameNode.GetGrandparent().IsObjectPattern() {
		var destructuringLhsNode pkg_rhino.Node = nameNode.GetGrandparent().GetParent()
		_ = destructuringLhsNode
		if NodeUtil_IsCallTo_2_Node_Node(destructuringLhsNode.GetLastChild(), StaticGet_ClosureRewriteModule_GOOG_REQUIRE()) || NodeUtil_IsCallTo_2_Node_Node(destructuringLhsNode.GetLastChild(), StaticGet_ClosureRewriteModule_GOOG_REQUIRETYPE()) {
			return
		}
	}
	var nameIsAnAlias bool = this_.currentScript.namesToInlineByAlias.ContainsKey(name)
	_ = nameIsAnAlias
	if nameIsAnAlias && (var_.GetNode() != nameNode) {
		this_.priv_ClosureRewriteModule_maybeAddAliasToSymbolTable_1230(nameNode, this_.currentScript.namespaceId)
		var inline *ClosureRewriteModule_AliasName = this_.currentScript.namesToInlineByAlias.Get(name)
		_ = inline
		var namespaceToInline string = inline.newName
		_ = namespaceToInline
		{
			var iifeCond_1 bool
			{

				sr_21 := namespaceToInline
				se_22 := this_.currentScript.GetBinaryNamespace()
				iifeCond_1 = j2gort.CheckNotNullString(sr_21) == se_22
			}
			if iifeCond_1 {
				{
					r_23 := this_.currentScript
					r_23.hasCreatedExportObject = true
				}
			}
		}
		var isModuleNamespace bool = (((inline.namespaceId != j2gort.NullString) && this_.rewriteState.scriptDescriptionsByGoogModuleNamespace.ContainsKey(inline.namespaceId)) && !(this_.rewriteState.scriptDescriptionsByGoogModuleNamespace.Get(inline.namespaceId).willCreateExportsObject))
		_ = isModuleNamespace
		this_.priv_ClosureRewriteModule_safeSetMaybeQualifiedString_1228(nameNode, namespaceToInline, isModuleNamespace)
		if j2gort.JStrIndexOfCp(namespaceToInline, int32(uint16(46))) != int32(-1) {
			var firstQualifiedName string = j2gort.JStrSub(namespaceToInline, int32(0), j2gort.JStrIndexOfCp(namespaceToInline, int32(uint16(46))))
			_ = firstQualifiedName
			var shadowedVar *Var = t.GetScope().LookupVar(firstQualifiedName)
			_ = shadowedVar
			if ((shadowedVar == nil) || shadowedVar.IsGlobal()) || j2gort.CastPtr[Scope](shadowedVar.GetScope()).IsModuleScope() {
				return
			}
			t.Report_3(shadowedVar.GetNode(), StaticGet_ClosureRewriteModule_IMPORT_INLINING_SHADOWS_VAR(), []string{shadowedVar.GetName(), namespaceToInline})
		}
		return
	}
	this_.priv_ClosureRewriteModule_safeSetString_1227(nameNode, (j2gort.StrValString(this_.currentScript.contentsPrefix) + j2gort.StrValString(name)))
}

// Java method: ClosureRewriteModule.maybeUpdateExportObjectLiteral(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n)
func (this_ *ClosureRewriteModule) priv_ClosureRewriteModule_maybeUpdateExportObjectLiteral_1216(t *NodeTraversal, n pkg_rhino.Node) {
	_ = t
	_ = n
	if !(this_.currentScript.isModule) {
		return
	}
	var parent pkg_rhino.Node = n.GetParent()
	_ = parent
	var rhs pkg_rhino.Node = parent.GetLastChild()
	_ = rhs
	if rhs.IsObjectLit() {
		{
			var c pkg_rhino.Node = rhs.GetFirstChild()
			_ = c
			for ; c != 0; c = c.GetNext() {
				if c.IsComputedProp() {
					t.Report_3(c, StaticGet_ClosureRewriteModule_INVALID_EXPORT_COMPUTED_PROPERTY(), []string{})
				} else if c.IsStringKey() {
					var value pkg_rhino.Node = c.GetFirstChild()
					_ = value
					this_.priv_ClosureRewriteModule_maybeUpdateExportDeclToNode_1217(t, c, value)
				}
			}
		}
	}
}

// Java method: ClosureRewriteModule.maybeUpdateExportDeclToNode(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node target, com.google.javascript.rhino.Node value)
func (this_ *ClosureRewriteModule) priv_ClosureRewriteModule_maybeUpdateExportDeclToNode_1217(t *NodeTraversal, target pkg_rhino.Node, value pkg_rhino.Node) {
	_ = t
	_ = target
	_ = value
	if !(this_.currentScript.isModule) {
		return
	}
	if value.IsName() {
		var currentScope *Scope = t.GetScope()
		_ = currentScope
		var v *Var = t.GetScope().LookupVar(value.GetString())
		_ = v
		if v != nil {
			var varScope AbstractScope = ToIface_AbstractScope(j2gort.CastPtr[Scope](v.GetScope()))
			_ = varScope
			if varScope.GetDepth() == currentScope.GetDepth() {
				var info *pkg_rhino.JSDocInfo = v.GetJSDocInfo()
				_ = info
				if (info != nil) && info.HasTypedefType() {
					var builder *pkg_rhino.JSDocInfo_Builder = pkg_rhino.JSDocInfo_Builder_CopyFrom(info)
					_ = builder
					target.SetJSDocInfo(builder.Build_0())
					return
				}
			}
		}
	}
	ClosureRewriteModule_MarkConstAndCopyJsDoc(target, target)
}

// Java method: ClosureRewriteModule.maybeUpdateExportDeclaration(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n)
func (this_ *ClosureRewriteModule) priv_ClosureRewriteModule_maybeUpdateExportDeclaration_1218(t *NodeTraversal, n pkg_rhino.Node) {
	_ = t
	_ = n
	if (!(this_.currentScript.isModule) || !(ClosureRewriteModule_IsGoogModuleExportsRef(t, n))) || !(ClosureRewriteModule_IsAssignTarget(n)) {
		return
	}
	var assignNode pkg_rhino.Node = n.GetParent()
	_ = assignNode
	var rhs pkg_rhino.Node = assignNode.GetLastChild()
	_ = rhs
	if rhs != this_.currentScript.defaultExport.rhs {
		assignNode.ReplaceWith(rhs.Detach())
		return
	}
	if !(this_.currentScript.declareLegacyNamespace) && (this_.currentScript.defaultExportLocalName != j2gort.NullString) {
		assignNode.GetParent().Detach()
		var binaryNamespaceName pkg_rhino.Node = this_.astFactory.CreateName_2_string_AstFactory_Type(this_.currentScript.GetBinaryNamespace(), AstFactory_S_Type_1_Node(n))
		_ = binaryNamespaceName
		this_.priv_ClosureRewriteModule_declareGlobalVariable_1231(binaryNamespaceName, t)
		return
	}
	var jsdocNode pkg_rhino.Node = pkg_rhino.Node(0)
	_ = jsdocNode
	if this_.currentScript.declareLegacyNamespace {
		var legacyQname pkg_rhino.Node = this_.astFactory.CreateQNameUsingJSTypeInfo(this_.globalTypedScope, this_.currentScript.namespaceId).SrcrefTree(n)
		_ = legacyQname
		legacyQname.SetJSType(n.GetJSType())
		n.ReplaceWith(legacyQname)
		jsdocNode = assignNode
	} else {
		rhs.Detach()
		var exprResultNode pkg_rhino.Node = assignNode.GetParent()
		_ = exprResultNode
		var binaryNamespaceName_2 pkg_rhino.Node = this_.astFactory.CreateName_2_string_AstFactory_Type(this_.currentScript.GetBinaryNamespace(), AstFactory_S_Type_1_Node(n))
		_ = binaryNamespaceName_2
		binaryNamespaceName_2.SetOriginalName("exports")
		this_.priv_ClosureRewriteModule_declareGlobalVariable_1231(binaryNamespaceName_2, t)
		var exportsObjectCreationNode pkg_rhino.Node = pkg_rhino.IR_Var_2(binaryNamespaceName_2, rhs)
		_ = exportsObjectCreationNode
		exportsObjectCreationNode.SrcrefTreeIfMissing(exprResultNode)
		exportsObjectCreationNode.PutBooleanProp(pkg_rhino.StaticGet_Node_IS_NAMESPACE(), true)
		exprResultNode.ReplaceWith(exportsObjectCreationNode)
		jsdocNode = exportsObjectCreationNode
		{
			r_24 := this_.currentScript
			r_24.hasCreatedExportObject = true
		}
	}
	ClosureRewriteModule_MarkConstAndCopyJsDoc(assignNode, jsdocNode)
	this_.compiler.ReportChangeToEnclosingScope(jsdocNode)
	this_.priv_ClosureRewriteModule_maybeUpdateExportObjectLiteral_1216(t, rhs)
}

// Java method: ClosureRewriteModule.maybeUpdateExportNameRef(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n)
func (this_ *ClosureRewriteModule) priv_ClosureRewriteModule_maybeUpdateExportNameRef_1219(t *NodeTraversal, n pkg_rhino.Node) {
	_ = t
	_ = n
	if (!(this_.currentScript.isModule) || !(ClosureRewriteModule_IsGoogModuleExportsRef(t, n))) || (n.GetParent() == 0) {
		return
	}
	if n.GetParent().IsParamList() {
		return
	}
	if this_.currentScript.declareLegacyNamespace {
		var legacyQname pkg_rhino.Node = this_.astFactory.CreateQNameUsingJSTypeInfo(this_.globalTypedScope, this_.currentScript.namespaceId).SrcrefTree(n)
		_ = legacyQname
		legacyQname.SetJSType(n.GetJSType())
		n.ReplaceWith(legacyQname)
		this_.compiler.ReportChangeToEnclosingScope(legacyQname)
		return
	}
	this_.priv_ClosureRewriteModule_safeSetString_1227(n, this_.currentScript.GetBinaryNamespace())
	if !(this_.currentScript.willCreateExportsObject || this_.currentScript.hasCreatedExportObject) {
		j2gort.CheckState(false)
	}
}

// Java method: ClosureRewriteModule.updateModuleBody(com.google.javascript.rhino.Node moduleBody)
func (this_ *ClosureRewriteModule) UpdateModuleBody(moduleBody pkg_rhino.Node) {
	_ = moduleBody
	if !(moduleBody.IsModuleBody() && moduleBody.GetParent().GetBooleanProp(pkg_rhino.StaticGet_Node_GOOG_MODULE())) {
		j2gort.CheckArgumentMsg(false, j2gort.BoxNode(moduleBody))
	}
	moduleBody.SetToken(pkg_rhino.StaticGet_Token_BLOCK())
	NodeUtil_TryMergeBlock(moduleBody, true)
	{
		cur_25 := this_.currentScript.exportsToInline.IterValues()
		for cur_25.HasNext() {

			export := cur_25.
				Next()
			_ = export
			var nameNode pkg_rhino.Node = export.nameDecl.GetNameNode()
			_ = nameNode
			this_.priv_ClosureRewriteModule_safeSetMaybeQualifiedString_1228(nameNode, (j2gort.StrValString(this_.currentScript.GetBinaryNamespace()) + j2gort.StrValString(export.GetExportPostfix())), false)
		}
	}
	if !(this_.currentScript.isModule) {
		j2gort.CheckStateMsg(false, j2gort.PtrToAny(this_.currentScript))
	}
	if !(this_.currentScript.declareLegacyNamespace || this_.currentScript.hasCreatedExportObject) {
		j2gort.CheckStateMsg(false, j2gort.PtrToAny(this_.currentScript))
	}
	this_.priv_ClosureRewriteModule_popScript_1222()
}

// Java method: ClosureRewriteModule.pushScript(com.google.javascript.jscomp.ClosureRewriteModule.ScriptDescription newCurrentScript)
func (this_ *ClosureRewriteModule) priv_ClosureRewriteModule_pushScript_1221(newCurrentScript *ClosureRewriteModule_ScriptDescription) {
	_ = newCurrentScript
	this_.currentScript = newCurrentScript
	if !(this_.scriptStack.IsEmpty()) {
		var parentScript *ClosureRewriteModule_ScriptDescription = this_.scriptStack.Peek()
		_ = parentScript
		parentScript.AddChildScript(this_.currentScript)
	}
	this_.scriptStack.AddFirst(this_.currentScript)
}

// Java method: ClosureRewriteModule.popScript()
func (this_ *ClosureRewriteModule) priv_ClosureRewriteModule_popScript_1222() {
	this_.scriptStack.RemoveFirst()
	this_.currentScript = this_.scriptStack.PeekFirst()
}

// Java method: ClosureRewriteModule.exportTheEmptyBinaryNamespaceAt(com.google.javascript.rhino.Node atNode, com.google.javascript.jscomp.ClosureRewriteModule.AddAt addAt, com.google.javascript.jscomp.NodeTraversal t)
func (this_ *ClosureRewriteModule) priv_ClosureRewriteModule_exportTheEmptyBinaryNamespaceAt_1223(atNode pkg_rhino.Node, addAt ClosureRewriteModule_AddAt, t *NodeTraversal) {
	_ = atNode
	_ = addAt
	_ = t
	if this_.currentScript.declareLegacyNamespace {
		return
	}
	var binaryNamespaceString string = this_.currentScript.GetBinaryNamespace()
	_ = binaryNamespaceString
	var moduleType AstFactory_Type = AstFactory_S_Type_1_Node(this_.currentScript.rootNode)
	_ = moduleType
	var binaryNamespaceName pkg_rhino.Node = this_.astFactory.CreateName_2_string_AstFactory_Type(binaryNamespaceString, moduleType)
	_ = binaryNamespaceName
	binaryNamespaceName.SetOriginalName(this_.currentScript.namespaceId)
	this_.priv_ClosureRewriteModule_declareGlobalVariable_1231(binaryNamespaceName, t)
	var binaryNamespaceExportNode pkg_rhino.Node = pkg_rhino.IR_Var_2(binaryNamespaceName, this_.astFactory.CreateObjectLit_1([]pkg_rhino.Node{}))
	_ = binaryNamespaceExportNode
	if addAt == StaticGet_ClosureRewriteModule_AddAt_BEFORE() {
		binaryNamespaceExportNode.InsertBefore(atNode)
	} else if addAt == StaticGet_ClosureRewriteModule_AddAt_AFTER() {
		binaryNamespaceExportNode.InsertAfter(atNode)
	}
	binaryNamespaceExportNode.PutBooleanProp(pkg_rhino.StaticGet_Node_IS_NAMESPACE(), true)
	binaryNamespaceExportNode.SrcrefTree(atNode)
	ClosureRewriteModule_MarkConst(binaryNamespaceExportNode)
	this_.compiler.ReportChangeToEnclosingScope(binaryNamespaceExportNode)
	{
		r_27 := this_.currentScript
		r_27.hasCreatedExportObject = true
	}
}

// Java static method: ClosureRewriteModule.checkAndSetStrictModeDirective(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n)
func ClosureRewriteModule_CheckAndSetStrictModeDirective(t *NodeTraversal, n pkg_rhino.Node) {
	clinit_ClosureRewriteModule()
	_ = t
	_ = n
	if !(n.IsScript()) {
		j2gort.CheckStateMsg(false, j2gort.BoxNode(n))
	}
	if n.IsUseStrict() {
		t.Report_3(n, StaticGet_ClosureRewriteModule_USELESS_USE_STRICT_DIRECTIVE(), []string{})
	} else {
		n.SetUseStrict(true)
	}
}

// Java static method: ClosureRewriteModule.markConst(com.google.javascript.rhino.Node n)
func ClosureRewriteModule_MarkConst(n pkg_rhino.Node) {
	clinit_ClosureRewriteModule()
	_ = n
	var builder *pkg_rhino.JSDocInfo_Builder = pkg_rhino.JSDocInfo_Builder_MaybeCopyFrom(n.GetJSDocInfo())
	_ = builder
	builder.RecordConstancy()
	n.SetJSDocInfo(builder.Build_0())
}

// Java static method: ClosureRewriteModule.maybeSplitMultiVar(com.google.javascript.rhino.Node rhsNode)
func ClosureRewriteModule_MaybeSplitMultiVar(rhsNode pkg_rhino.Node) {
	clinit_ClosureRewriteModule()
	_ = rhsNode
	var statementNode pkg_rhino.Node = rhsNode.GetGrandparent()
	_ = statementNode
	if !(statementNode.IsVar()) || !(statementNode.HasMoreThanOneChild()) {
		return
	}
	var nameNode pkg_rhino.Node = rhsNode.GetParent()
	_ = nameNode
	nameNode.Detach()
	rhsNode.Detach()
	pkg_rhino.IR_Var_2(nameNode, rhsNode).InsertBefore(statementNode)
}

// Java static method: ClosureRewriteModule.markConstAndCopyJsDoc(com.google.javascript.rhino.Node from, com.google.javascript.rhino.Node target)
func ClosureRewriteModule_MarkConstAndCopyJsDoc(from pkg_rhino.Node, target pkg_rhino.Node) {
	clinit_ClosureRewriteModule()
	_ = from
	_ = target
	var info *pkg_rhino.JSDocInfo = from.GetJSDocInfo()
	_ = info
	var builder *pkg_rhino.JSDocInfo_Builder = pkg_rhino.JSDocInfo_Builder_MaybeCopyFrom(info)
	_ = builder
	builder.RecordConstancy()
	target.SetJSDocInfo(builder.Build_0())
}

// Java method: ClosureRewriteModule.recordExportToInline(com.google.javascript.jscomp.ClosureRewriteModule.ExportDefinition exportDefinition)
func (this_ *ClosureRewriteModule) priv_ClosureRewriteModule_recordExportToInline_1224(exportDefinition *ClosureRewriteModule_ExportDefinition) {
	_ = exportDefinition
	j2gort.CheckStateFmt(exportDefinition.HasInlinableName(j2gort.AsIface[j2gort.Ext_java_util_Set](this_.currentScript.exportsToInline.LiveKeySet())), "exportDefinition: %s\n\nexportsToInline keys: %s", j2gort.PtrToAny(exportDefinition), j2gort.AsIface[j2gort.Ext_java_util_Set](this_.currentScript.exportsToInline.LiveKeySet()))
	if !(this_.currentScript.exportsToInline.Put(exportDefinition.nameDecl, exportDefinition) == nil) {
		j2gort.CheckStateFmt(false, "Already found a mapping for inlining export: %s", j2gort.PtrToAny(exportDefinition.nameDecl))
	}
	var localName string = exportDefinition.GetLocalName()
	_ = localName
	var fullExportedName string = (j2gort.StrValString(this_.currentScript.GetBinaryNamespace()) + j2gort.StrValString(exportDefinition.GetExportPostfix()))
	_ = fullExportedName
	this_.priv_ClosureRewriteModule_recordNameToInline_1225(localName, fullExportedName, j2gort.NullString)
}

// Java method: ClosureRewriteModule.recordNameToInline(java.lang.String aliasName, java.lang.String newName, java.lang.@org.jspecify.annotations.Nullable String namespaceId)
func (this_ *ClosureRewriteModule) priv_ClosureRewriteModule_recordNameToInline_1225(aliasName string, newName string, namespaceId string) {
	_ = aliasName
	_ = newName
	_ = namespaceId
	j2gort.CheckNotNullString(aliasName)
	j2gort.CheckNotNullString(newName)
	this_.currentScript.namesToInlineByAlias.Put(aliasName, NewClosureRewriteModule_AliasName(newName, namespaceId))
}

// Java method: ClosureRewriteModule.reportUnrecognizedRequires()
func (this_ *ClosureRewriteModule) priv_ClosureRewriteModule_reportUnrecognizedRequires_1226() {
	for _, unrecognizedRequire := range this_.unrecognizedRequires {
		_ = unrecognizedRequire
		var namespaceId string = unrecognizedRequire.namespaceId
		_ = namespaceId
		var requireNode pkg_rhino.Node = unrecognizedRequire.requireNode
		_ = requireNode
		var targetGoogModuleExists bool = this_.rewriteState.ContainsModule(namespaceId)
		_ = targetGoogModuleExists
		var targetLegacyScriptExists bool = this_.rewriteState.providedNamespaces.Contains(namespaceId)
		_ = targetLegacyScriptExists
		if targetGoogModuleExists || targetLegacyScriptExists {
			continue
		}
		if this_.preserveSugar {
			continue
		}
		if NodeUtil_GetEnclosingScript(requireNode) == 0 {
			continue
		}
		this_.compiler.ReportChangeToEnclosingScope(requireNode)
		var enclosingStatement pkg_rhino.Node = NodeUtil_GetEnclosingStatement(requireNode)
		_ = enclosingStatement
		if !(NodeUtil_IsNameDeclaration(enclosingStatement)) {
			requireNode.ReplaceWith(this_.astFactory.CreateNull().Srcref(requireNode))
			continue
		}
		enclosingStatement.Detach()
		NodeUtil_VisitLhsNodesInNode(enclosingStatement, j2gort.NewLambda_Ext_java_util_function_Consumer(func(lp_30 any) {
			lhs := j2gort.CastNode[pkg_rhino.Node](lp_30)
			_ = lhs
			this_.syntheticExterns.PutIfAbsent(lhs.GetString(), lhs)
		}))
	}

}

// Java method: ClosureRewriteModule.safeSetString(com.google.javascript.rhino.Node n, java.lang.String newString)
func (this_ *ClosureRewriteModule) priv_ClosureRewriteModule_safeSetString_1227(n pkg_rhino.Node, newString string) {
	_ = n
	_ = newString
	if j2gort.CheckNotNullString(n.GetString()) == newString {
		return
	}
	var originalName string = n.GetString()
	_ = originalName
	n.SetString(newString)
	if n.GetOriginalName() == j2gort.NullString {
		n.SetOriginalName(originalName)
	}
	var changeScope pkg_rhino.Node = ChangeTracker_GetEnclosingChangeScopeRoot(n)
	_ = changeScope
	if changeScope != 0 {
		this_.compiler.ReportChangeToChangeScope(changeScope)
	}
}

// Java method: ClosureRewriteModule.safeSetMaybeQualifiedString(com.google.javascript.rhino.Node nameNode, java.lang.String newString, boolean isModuleNamespace)
func (this_ *ClosureRewriteModule) priv_ClosureRewriteModule_safeSetMaybeQualifiedString_1228(nameNode pkg_rhino.Node, newString string, isModuleNamespace bool) {
	_ = nameNode
	_ = newString
	_ = isModuleNamespace
	if !(j2gort.JStrIndexOf(newString, ".") >= 0) {
		this_.priv_ClosureRewriteModule_safeSetString_1227(nameNode, newString)
		var parent pkg_rhino.Node = nameNode.GetParent()
		_ = parent
		if ((isModuleNamespace && parent.IsGetProp()) && nameNode.GetGrandparent().IsCall()) && parent.IsFirstChildOf(nameNode.GetGrandparent()) {
			nameNode.GetGrandparent().PutBooleanProp(pkg_rhino.StaticGet_Node_FREE_CALL(), true)
		}
		return
	}
	var nameParent pkg_rhino.Node = nameNode.GetParent()
	_ = nameParent
	var newQualifiedName pkg_rhino.Node = this_.astFactory.CreateQNameUsingJSTypeInfo(this_.globalTypedScope, newString).SrcrefTree(nameNode)
	_ = newQualifiedName
	newQualifiedName.SetJSType(nameNode.GetJSType())
	var replaced bool = ClosureRewriteModule_SafeSetStringIfDeclaration(nameParent, nameNode, newQualifiedName)
	_ = replaced
	if replaced {
		return
	}
	nameNode.ReplaceWith(newQualifiedName)
	if newQualifiedName.HasChildren() {
		newQualifiedName.GetFirstChild().MakeNonIndexableRecursive()
	}
	this_.compiler.ReportChangeToEnclosingScope(newQualifiedName)
}

// Java static method: ClosureRewriteModule.safeSetStringIfDeclaration(com.google.javascript.rhino.Node nameParent, com.google.javascript.rhino.Node nameNode, com.google.javascript.rhino.Node newQualifiedName)
func ClosureRewriteModule_SafeSetStringIfDeclaration(nameParent pkg_rhino.Node, nameNode pkg_rhino.Node, newQualifiedName pkg_rhino.Node) bool {
	clinit_ClosureRewriteModule()
	_ = nameParent
	_ = nameNode
	_ = newQualifiedName
	var jsdoc *pkg_rhino.JSDocInfo = nameParent.GetJSDocInfo()
	_ = jsdoc
	switch nameParent.TokenSwitchOrd_() {
	case int32(73) /* Token.FUNCTION */, int32(105) /* Token.CLASS */ :
		if !(NodeUtil_IsStatement(nameParent)) || (nameParent.GetFirstChild() != nameNode) {
			return false
		}
		var placeholder pkg_rhino.Node = pkg_rhino.IR_Empty()
		_ = placeholder
		nameParent.ReplaceWith(placeholder)
		var newDeclaration pkg_rhino.Node = NodeUtil_GetDeclarationFromName(newQualifiedName, nameParent, pkg_rhino.StaticGet_Token_VAR(), jsdoc)
		_ = newDeclaration
		if NodeUtil_IsExprAssign(newDeclaration) {
			var assign pkg_rhino.Node = newDeclaration.GetOnlyChild()
			_ = assign
			assign.SetJSType(nameNode.GetJSType())
			ClosureRewriteModule_UpdateSourceInfoForExportedTopLevelVariable(assign, nameNode)
		}
		nameParent.SetJSDocInfo(nil)
		newDeclaration.SrcrefTreeIfMissing(nameParent)
		placeholder.ReplaceWith(newDeclaration)
		NodeUtil_RemoveName(nameParent)
		return true
	case int32(84) /* Token.VAR */, int32(111) /* Token.LET */, int32(97) /* Token.CONST */ :
		var rhs pkg_rhino.Node
		_ = rhs
		if nameNode.HasChildren() {
			rhs = nameNode.GetLastChild().Detach()
		} else {
			rhs = pkg_rhino.Node(0)
		}
		if jsdoc == nil {
			jsdoc = nameNode.GetJSDocInfo()
		}
		var newStatement pkg_rhino.Node = NodeUtil_GetDeclarationFromName(newQualifiedName, rhs, pkg_rhino.StaticGet_Token_VAR(), jsdoc)
		_ = newStatement
		if NodeUtil_IsExprAssign(newStatement) {
			var assign_2 pkg_rhino.Node = newStatement.GetOnlyChild()
			_ = assign_2
			assign_2.SetJSType(nameNode.GetJSType())
			ClosureRewriteModule_UpdateSourceInfoForExportedTopLevelVariable(assign_2, nameNode)
			if nameParent.IsConst() {
				var jsdocBuilder *pkg_rhino.JSDocInfo_Builder = pkg_rhino.JSDocInfo_Builder_MaybeCopyFrom(jsdoc)
				_ = jsdocBuilder
				jsdocBuilder.RecordConstancy()
				jsdoc = jsdocBuilder.Build_0()
				assign_2.SetJSDocInfo(jsdoc)
			}
		}
		newStatement.SrcrefTreeIfMissing(nameParent)
		NodeUtil_ReplaceDeclarationChild(nameNode, newStatement)
		return true
	case int32(103) /* Token.OBJECT_PATTERN */, int32(102) /* Token.ARRAY_PATTERN */, int32(49) /* Token.PARAM_LIST */ :
		panic(j2gort.AsIface[javalang.IRuntimeException](j2gort.ConstructShimThrowable("java.lang.RuntimeException", "Not supported")))
	default:
	}
	return false
}

// Java static method: ClosureRewriteModule.updateSourceInfoForExportedTopLevelVariable(com.google.javascript.rhino.Node assign, com.google.javascript.rhino.Node sourceName)
func ClosureRewriteModule_UpdateSourceInfoForExportedTopLevelVariable(assign pkg_rhino.Node, sourceName pkg_rhino.Node) {
	clinit_ClosureRewriteModule()
	_ = assign
	_ = sourceName
	if !(assign.IsAssign()) {
		j2gort.CheckState(false)
	}
	if !(sourceName.IsName()) {
		j2gort.CheckState(false)
	}
	var getProp pkg_rhino.Node = assign.GetFirstChild()
	_ = getProp
	if !(getProp.IsGetProp()) {
		return
	}
	var name string = sourceName.GetOriginalName()
	_ = name
	if name == j2gort.NullString {
		name = sourceName.GetString()
	}
	sourceName = sourceName.CloneNode_0()
	sourceName.SetLength(j2gort.JStrLen(name))
	getProp.SrcrefTree(sourceName)
}

// Java method: ClosureRewriteModule.isTopLevel(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.jscomp.ClosureRewriteModule.ScopeType scopeType)
func (this_ *ClosureRewriteModule) priv_ClosureRewriteModule_isTopLevel_1229(t *NodeTraversal, n pkg_rhino.Node, scopeType ClosureRewriteModule_ScopeType) bool {
	_ = t
	_ = n
	_ = scopeType
	if scopeType == StaticGet_ClosureRewriteModule_ScopeType_EXEC_CONTEXT() {
		return (t.InGlobalScope() || (t.GetClosestHoistScopeRoot() == this_.currentScript.rootNode))
	} else {
		return (n.GetParent() == this_.currentScript.rootNode)
	}
}

// Java static method: ClosureRewriteModule.toModuleContentsPrefix(java.lang.String namespaceId)
func ClosureRewriteModule_ToModuleContentsPrefix(namespaceId string) string {
	clinit_ClosureRewriteModule()
	_ = namespaceId
	return ("module$contents$" + j2gort.StrValString(j2gort.JStrReplaceChar(namespaceId, uint16(46), uint16(36))) + "_")
}

// Java static method: ClosureRewriteModule.isModuleExport(java.lang.String name)
func ClosureRewriteModule_IsModuleExport(name string) bool {
	clinit_ClosureRewriteModule()
	_ = name
	return javalang.StringStartsWith(j2gort.CheckNotNullString(name), "module$exports$")
}

// Java static method: ClosureRewriteModule.isModuleContent(java.lang.String name)
func ClosureRewriteModule_IsModuleContent(name string) bool {
	clinit_ClosureRewriteModule()
	_ = name
	return javalang.StringStartsWith(j2gort.CheckNotNullString(name), "module$contents$")
}

// Java static method: ClosureRewriteModule.isGoogModuleExportsRef(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node target)
func ClosureRewriteModule_IsGoogModuleExportsRef(t *NodeTraversal, target pkg_rhino.Node) bool {
	clinit_ClosureRewriteModule()
	_ = t
	_ = target
	if !(target.IsName()) || !(target.MatchesName_1_string("exports")) {
		return false
	}
	var exportsVar *Var = t.GetScope().LookupVar("exports")
	_ = exportsVar
	return ((exportsVar != nil) && exportsVar.IsGoogModuleExports())
}

// Java static method: ClosureRewriteModule.isExportPropertyAssignment(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n)
func ClosureRewriteModule_IsExportPropertyAssignment(t *NodeTraversal, n pkg_rhino.Node) bool {
	clinit_ClosureRewriteModule()
	_ = t
	_ = n
	var target pkg_rhino.Node = n.GetFirstChild()
	_ = target
	return ((ClosureRewriteModule_IsAssignTarget(n) || ClosureRewriteModule_IsTypedefTarget(n)) && ClosureRewriteModule_IsGoogModuleExportsRef(t, target))
}

// Java static method: ClosureRewriteModule.isAssignTarget(com.google.javascript.rhino.Node n)
func ClosureRewriteModule_IsAssignTarget(n pkg_rhino.Node) bool {
	clinit_ClosureRewriteModule()
	_ = n
	var parent pkg_rhino.Node = n.GetParent()
	_ = parent
	return (parent.IsAssign() && (parent.GetFirstChild() == n))
}

// Java static method: ClosureRewriteModule.isTypedefTarget(com.google.javascript.rhino.Node n)
func ClosureRewriteModule_IsTypedefTarget(n pkg_rhino.Node) bool {
	clinit_ClosureRewriteModule()
	_ = n
	var parent pkg_rhino.Node = n.GetParent()
	_ = parent
	return (parent.IsExprResult() && (parent.GetFirstChild() == n))
}

// Java method: ClosureRewriteModule.maybeAddAliasToSymbolTable(com.google.javascript.rhino.Node n, java.lang.String module)
func (this_ *ClosureRewriteModule) priv_ClosureRewriteModule_maybeAddAliasToSymbolTable_1230(n pkg_rhino.Node, module string) {
	_ = n
	_ = module
	if this_.preprocessorSymbolTable != nil {
		n.PutBooleanProp(pkg_rhino.StaticGet_Node_MODULE_ALIAS(), true)
		var nodeName string
		_ = nodeName
		if n.IsStringLit() {
			nodeName = n.GetString()
		} else {
			nodeName = this_.preprocessorSymbolTable.GetQualifiedName(n)
		}
		var name string = ("alias_" + j2gort.StrValString(module) + "_" + j2gort.StrValString(nodeName))
		_ = name
		this_.preprocessorSymbolTable.AddReference_2(n, name)
	}
}

// Java method: ClosureRewriteModule.declareGlobalVariable(com.google.javascript.rhino.Node n, com.google.javascript.jscomp.NodeTraversal t)
func (this_ *ClosureRewriteModule) priv_ClosureRewriteModule_declareGlobalVariable_1231(n pkg_rhino.Node, t *NodeTraversal) {
	_ = n
	_ = t
	if !(n.IsName()) {
		j2gort.CheckState(false)
	}
	if this_.globalTypedScope == nil {
		return
	}
	var name string = n.GetString()
	_ = name
	if this_.globalTypedScope.HasOwnSlot(name) {
		var original pkg_rhino.Node = j2gort.CastPtr[TypedVar](this_.globalTypedScope.GetOwnSlot(name)).GetNode()
		_ = original
		t.Report_3(t.GetCurrentScript(), StaticGet_ClosureRewriteModule_ILLEGAL_MODULE_RENAMING_CONFLICT(), []string{name, func() string {
			if original != 0 {
				return original.String()
			}
			return "<unknown>"
		}()})
	} else {
		var type_ pkg_rhino.JSType = j2gort.CheckNotNullIface(n.GetJSType())
		_ = type_
		this_.globalTypedScope.Declare(name, n, type_, t.GetInput(), false)
	}
}

// Java enum: com.google.javascript.jscomp.ClosureRewriteModule.AddAt
type ClosureRewriteModule_AddAt uint8

const (
	enum_ClosureRewriteModule_AddAt_BEFORE ClosureRewriteModule_AddAt = 1
	enum_ClosureRewriteModule_AddAt_AFTER  ClosureRewriteModule_AddAt = 2
)

var enumBaseTable_ClosureRewriteModule_AddAt = [...]j2gort.EnumBase{
	{Name_: "BEFORE", Ordinal_: 0},
	{Name_: "AFTER", Ordinal_: 1},
}
var enumValues_ClosureRewriteModule_AddAt = []ClosureRewriteModule_AddAt{1, 2}

func StaticGet_ClosureRewriteModule_AddAt_BEFORE() ClosureRewriteModule_AddAt {
	return enum_ClosureRewriteModule_AddAt_BEFORE
}
func StaticGet_ClosureRewriteModule_AddAt_AFTER() ClosureRewriteModule_AddAt {
	return enum_ClosureRewriteModule_AddAt_AFTER
}

func ClosureRewriteModule_AddAt_Values() []ClosureRewriteModule_AddAt {
	return j2gort.CloneSlice(enumValues_ClosureRewriteModule_AddAt)
}

func ClosureRewriteModule_AddAt_ValueOf(name string) ClosureRewriteModule_AddAt {
	if name == j2gort.NullString {
		j2gort.ThrowNPE()
	}
	for i_ := range enumBaseTable_ClosureRewriteModule_AddAt {
		if enumBaseTable_ClosureRewriteModule_AddAt[i_].Name_ == name {
			return ClosureRewriteModule_AddAt(i_ + 1)
		}
	}
	panic(javalang.NewIllegalArgumentException("No enum constant com.google.javascript.jscomp.ClosureRewriteModule.AddAt." + name))
}

func (this_ ClosureRewriteModule_AddAt) IsJavaNil_() bool                              { return this_ == 0 }
func (ClosureRewriteModule_AddAt) HasIdentityEquals_()                                 {}
func (ClosureRewriteModule_AddAt) Implements_j2gort_Ext_java_lang_Enum()               {}
func (ClosureRewriteModule_AddAt) Implements_j2gort_Ext_java_lang_constant_Constable() {}
func (this_ ClosureRewriteModule_AddAt) GetExt_java_lang_EnumBase() *j2gort.Ext_java_lang_EnumBase {
	return &j2gort.Ext_java_lang_EnumBase{Core_: this_.GetEnumBase()}
}
func (ClosureRewriteModule_AddAt) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.ClosureRewriteModule$AddAt", "AddAt"
}
func (this_ ClosureRewriteModule_AddAt) GetEnumBase() *j2gort.EnumBase {
	if this_ == 0 {
		return nil
	}
	return &enumBaseTable_ClosureRewriteModule_AddAt[this_-1]
}
func (this_ ClosureRewriteModule_AddAt) Ordinal() int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return int32(this_ - 1)
}
func (this_ ClosureRewriteModule_AddAt) EnumOrdinal_() int32 { return int32(this_) - 1 }
func (this_ ClosureRewriteModule_AddAt) Name() string {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_ClosureRewriteModule_AddAt[this_-1].Name_
}
func (this_ ClosureRewriteModule_AddAt) String() string {
	if this_ == 0 {
		return "null"
	}
	return enumBaseTable_ClosureRewriteModule_AddAt[this_-1].Name_
}
func (this_ ClosureRewriteModule_AddAt) Equals(other_ any) bool {
	o_, ok_ := other_.(ClosureRewriteModule_AddAt)
	return this_ != 0 && ok_ && this_ == o_
}
func (this_ ClosureRewriteModule_AddAt) HashCode() int32 {
	if this_ == 0 {
		return 0
	}
	return enumBaseTable_ClosureRewriteModule_AddAt[this_-1].HashCode()
}
func (this_ ClosureRewriteModule_AddAt) CompareTo(other_ any) int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_ClosureRewriteModule_AddAt[this_-1].CompareTo(other_)
}
func (this_ ClosureRewriteModule_AddAt) WriteGson(out_ *gson.JsonWriter) {
	if this_ == 0 {
		out_.NullValue()
		return
	}
	out_.ValueString(enumBaseTable_ClosureRewriteModule_AddAt[this_-1].Name_)
}
func init() {
	vals_ := []any{ClosureRewriteModule_AddAt(1), ClosureRewriteModule_AddAt(2)}
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.ClosureRewriteModule.AddAt", vals_)
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.ClosureRewriteModule$AddAt", vals_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.ClosureRewriteModule.AddAt", "com.google.javascript.jscomp.ClosureRewriteModule$AddAt", func(v_ any) bool {
		return j2gort.InstanceOf[ClosureRewriteModule_AddAt](v_)
	})
}

// Java enum: com.google.javascript.jscomp.ClosureRewriteModule.ScopeType
type ClosureRewriteModule_ScopeType uint8

const (
	enum_ClosureRewriteModule_ScopeType_EXEC_CONTEXT ClosureRewriteModule_ScopeType = 1
	enum_ClosureRewriteModule_ScopeType_BLOCK        ClosureRewriteModule_ScopeType = 2
)

var enumBaseTable_ClosureRewriteModule_ScopeType = [...]j2gort.EnumBase{
	{Name_: "EXEC_CONTEXT", Ordinal_: 0},
	{Name_: "BLOCK", Ordinal_: 1},
}
var enumValues_ClosureRewriteModule_ScopeType = []ClosureRewriteModule_ScopeType{1, 2}

func StaticGet_ClosureRewriteModule_ScopeType_EXEC_CONTEXT() ClosureRewriteModule_ScopeType {
	return enum_ClosureRewriteModule_ScopeType_EXEC_CONTEXT
}
func StaticGet_ClosureRewriteModule_ScopeType_BLOCK() ClosureRewriteModule_ScopeType {
	return enum_ClosureRewriteModule_ScopeType_BLOCK
}

func ClosureRewriteModule_ScopeType_Values() []ClosureRewriteModule_ScopeType {
	return j2gort.CloneSlice(enumValues_ClosureRewriteModule_ScopeType)
}

func ClosureRewriteModule_ScopeType_ValueOf(name string) ClosureRewriteModule_ScopeType {
	if name == j2gort.NullString {
		j2gort.ThrowNPE()
	}
	for i_ := range enumBaseTable_ClosureRewriteModule_ScopeType {
		if enumBaseTable_ClosureRewriteModule_ScopeType[i_].Name_ == name {
			return ClosureRewriteModule_ScopeType(i_ + 1)
		}
	}
	panic(javalang.NewIllegalArgumentException("No enum constant com.google.javascript.jscomp.ClosureRewriteModule.ScopeType." + name))
}

func (this_ ClosureRewriteModule_ScopeType) IsJavaNil_() bool                              { return this_ == 0 }
func (ClosureRewriteModule_ScopeType) HasIdentityEquals_()                                 {}
func (ClosureRewriteModule_ScopeType) Implements_j2gort_Ext_java_lang_Enum()               {}
func (ClosureRewriteModule_ScopeType) Implements_j2gort_Ext_java_lang_constant_Constable() {}
func (this_ ClosureRewriteModule_ScopeType) GetExt_java_lang_EnumBase() *j2gort.Ext_java_lang_EnumBase {
	return &j2gort.Ext_java_lang_EnumBase{Core_: this_.GetEnumBase()}
}
func (ClosureRewriteModule_ScopeType) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.ClosureRewriteModule$ScopeType", "ScopeType"
}
func (this_ ClosureRewriteModule_ScopeType) GetEnumBase() *j2gort.EnumBase {
	if this_ == 0 {
		return nil
	}
	return &enumBaseTable_ClosureRewriteModule_ScopeType[this_-1]
}
func (this_ ClosureRewriteModule_ScopeType) Ordinal() int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return int32(this_ - 1)
}
func (this_ ClosureRewriteModule_ScopeType) EnumOrdinal_() int32 { return int32(this_) - 1 }
func (this_ ClosureRewriteModule_ScopeType) Name() string {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_ClosureRewriteModule_ScopeType[this_-1].Name_
}
func (this_ ClosureRewriteModule_ScopeType) String() string {
	if this_ == 0 {
		return "null"
	}
	return enumBaseTable_ClosureRewriteModule_ScopeType[this_-1].Name_
}
func (this_ ClosureRewriteModule_ScopeType) Equals(other_ any) bool {
	o_, ok_ := other_.(ClosureRewriteModule_ScopeType)
	return this_ != 0 && ok_ && this_ == o_
}
func (this_ ClosureRewriteModule_ScopeType) HashCode() int32 {
	if this_ == 0 {
		return 0
	}
	return enumBaseTable_ClosureRewriteModule_ScopeType[this_-1].HashCode()
}
func (this_ ClosureRewriteModule_ScopeType) CompareTo(other_ any) int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_ClosureRewriteModule_ScopeType[this_-1].CompareTo(other_)
}
func (this_ ClosureRewriteModule_ScopeType) WriteGson(out_ *gson.JsonWriter) {
	if this_ == 0 {
		out_.NullValue()
		return
	}
	out_.ValueString(enumBaseTable_ClosureRewriteModule_ScopeType[this_-1].Name_)
}
func init() {
	vals_ := []any{ClosureRewriteModule_ScopeType(1), ClosureRewriteModule_ScopeType(2)}
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.ClosureRewriteModule.ScopeType", vals_)
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.ClosureRewriteModule$ScopeType", vals_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.ClosureRewriteModule.ScopeType", "com.google.javascript.jscomp.ClosureRewriteModule$ScopeType", func(v_ any) bool {
		return j2gort.InstanceOf[ClosureRewriteModule_ScopeType](v_)
	})
}

// Java class: com.google.javascript.jscomp.ClosureRewriteModule.UnrecognizedRequire
type ClosureRewriteModule_UnrecognizedRequire struct {
	namespaceId string         `gson:"namespaceId"`
	requireNode pkg_rhino.Node `gson:"requireNode"`
	idHash_     atomic.Int32
}

func (this_ *ClosureRewriteModule_UnrecognizedRequire) BindSelf_(s any) {
	_ = s
}
func (this_ *ClosureRewriteModule_UnrecognizedRequire) IsJavaNil_() bool { return this_ == nil }
func (*ClosureRewriteModule_UnrecognizedRequire) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.ClosureRewriteModule$UnrecognizedRequire", "UnrecognizedRequire"
}
func (this_ *ClosureRewriteModule_UnrecognizedRequire) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*ClosureRewriteModule_UnrecognizedRequire) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.ClosureRewriteModule.UnrecognizedRequire", "com.google.javascript.jscomp.ClosureRewriteModule$UnrecognizedRequire", func(v_ any) bool {
		return j2gort.InstanceOf[*ClosureRewriteModule_UnrecognizedRequire](v_)
	})
}

func (this_ *ClosureRewriteModule_UnrecognizedRequire) initFields_ClosureRewriteModule_UnrecognizedRequire() {
	_ = this_
	this_.namespaceId = j2gort.NullString
}

// Java constructor: UnrecognizedRequire(com.google.javascript.rhino.Node requireNode, java.lang.String namespaceId)
func NewClosureRewriteModule_UnrecognizedRequire(requireNode pkg_rhino.Node, namespaceId string) *ClosureRewriteModule_UnrecognizedRequire {
	this_ := &ClosureRewriteModule_UnrecognizedRequire{}
	InitClosureRewriteModule_UnrecognizedRequire(this_, requireNode, namespaceId)
	return this_
}

func InitClosureRewriteModule_UnrecognizedRequire(this_ *ClosureRewriteModule_UnrecognizedRequire, requireNode pkg_rhino.Node, namespaceId string) {
	_ = requireNode
	_ = namespaceId
	this_.initFields_ClosureRewriteModule_UnrecognizedRequire()
	this_.requireNode = requireNode
	this_.namespaceId = namespaceId
}

// Java class: com.google.javascript.jscomp.ClosureRewriteModule.ExportDefinition
type ClosureRewriteModule_ExportDefinition struct {
	exportName string         `gson:"exportName"`
	nameDecl   *Var           `gson:"nameDecl"`
	rhs        pkg_rhino.Node `gson:"rhs"`
	idHash_    atomic.Int32
}

func (this_ *ClosureRewriteModule_ExportDefinition) BindSelf_(s any) {
	_ = s
}
func (this_ *ClosureRewriteModule_ExportDefinition) IsJavaNil_() bool { return this_ == nil }
func (*ClosureRewriteModule_ExportDefinition) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.ClosureRewriteModule$ExportDefinition", "ExportDefinition"
}
func (this_ *ClosureRewriteModule_ExportDefinition) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*ClosureRewriteModule_ExportDefinition) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.ClosureRewriteModule.ExportDefinition", "com.google.javascript.jscomp.ClosureRewriteModule$ExportDefinition", func(v_ any) bool {
		return j2gort.InstanceOf[*ClosureRewriteModule_ExportDefinition](v_)
	})
}

var static_ClosureRewriteModule_ExportDefinition_INLINABLE_NAME_PARENTS j2gort.Ext_com_google_common_collect_ImmutableSet

func StaticGet_ClosureRewriteModule_ExportDefinition_INLINABLE_NAME_PARENTS() j2gort.Ext_com_google_common_collect_ImmutableSet {
	clinit_ClosureRewriteModule_ExportDefinition()
	return static_ClosureRewriteModule_ExportDefinition_INLINABLE_NAME_PARENTS
}
func StaticSet_ClosureRewriteModule_ExportDefinition_INLINABLE_NAME_PARENTS(v j2gort.Ext_com_google_common_collect_ImmutableSet) {
	clinit_ClosureRewriteModule_ExportDefinition()
	static_ClosureRewriteModule_ExportDefinition_INLINABLE_NAME_PARENTS = v
}

var clinitDone_ClosureRewriteModule_ExportDefinition uint32
var clinitState_ClosureRewriteModule_ExportDefinition j2gort.ClinitState

func clinit_ClosureRewriteModule_ExportDefinition() {
	if atomic.LoadUint32(&clinitDone_ClosureRewriteModule_ExportDefinition) != 0 {
		return
	}
	clinitSlow_ClosureRewriteModule_ExportDefinition()
}

func clinitSlow_ClosureRewriteModule_ExportDefinition() {
	if !j2gort.EnterClinit(&clinitState_ClosureRewriteModule_ExportDefinition, &clinitDone_ClosureRewriteModule_ExportDefinition) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_ClosureRewriteModule_ExportDefinition, &clinitDone_ClosureRewriteModule_ExportDefinition)
	static_ClosureRewriteModule_ExportDefinition_INLINABLE_NAME_PARENTS = j2gort.Ext_com_google_common_collect_Sets_ImmutableEnumSet_2(j2gort.ToIface_Ext_java_lang_Enum(pkg_rhino.StaticGet_Token_VAR()), []j2gort.Ext_java_lang_Enum{j2gort.ToIface_Ext_java_lang_Enum(pkg_rhino.StaticGet_Token_CONST()), j2gort.ToIface_Ext_java_lang_Enum(pkg_rhino.StaticGet_Token_LET()), j2gort.ToIface_Ext_java_lang_Enum(pkg_rhino.StaticGet_Token_FUNCTION()), j2gort.ToIface_Ext_java_lang_Enum(pkg_rhino.StaticGet_Token_CLASS())})
}

func (this_ *ClosureRewriteModule_ExportDefinition) initFields_ClosureRewriteModule_ExportDefinition() {
	_ = this_
	this_.exportName = j2gort.NullString
}

// Java constructor: ExportDefinition()
func NewClosureRewriteModule_ExportDefinition() *ClosureRewriteModule_ExportDefinition {
	this_ := &ClosureRewriteModule_ExportDefinition{}
	InitClosureRewriteModule_ExportDefinition(this_)
	return this_
}

func InitClosureRewriteModule_ExportDefinition(this_ *ClosureRewriteModule_ExportDefinition) {
	this_.initFields_ClosureRewriteModule_ExportDefinition()
}

// Java method: ExportDefinition.toString()
func (this_ *ClosureRewriteModule_ExportDefinition) String() string {
	return j2gort.Ext_com_google_common_base_MoreObjects_S_ToStringHelper_1_any(this_).Add_2_string_any("exportName", j2gort.BoxString(this_.exportName)).Add_2_string_any("rhs", j2gort.BoxNode(this_.rhs)).Add_2_string_any("nameDecl", j2gort.PtrToAny(this_.nameDecl)).OmitNullValues().String()
}

// Java static method: ExportDefinition.newDefaultExport(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node rhs)
func ClosureRewriteModule_ExportDefinition_NewDefaultExport(t *NodeTraversal, rhs pkg_rhino.Node) *ClosureRewriteModule_ExportDefinition {
	_ = t
	_ = rhs
	return ClosureRewriteModule_ExportDefinition_NewNamedExport(t, j2gort.NullString, rhs)
}

// Java static method: ExportDefinition.newNamedExport(com.google.javascript.jscomp.NodeTraversal t, java.lang.@org.jspecify.annotations.Nullable String name, com.google.javascript.rhino.Node rhs)
func ClosureRewriteModule_ExportDefinition_NewNamedExport(t *NodeTraversal, name string, rhs pkg_rhino.Node) *ClosureRewriteModule_ExportDefinition {
	_ = t
	_ = name
	_ = rhs
	var val_31 ClosureRewriteModule_ExportDefinition
	InitClosureRewriteModule_ExportDefinition(&val_31)
	var newExport *ClosureRewriteModule_ExportDefinition = &val_31
	_ = newExport
	{
		r_32 := newExport
		r_32.exportName = name
	}
	{
		r_33 := newExport
		r_33.rhs = rhs
	}
	if (rhs != 0) && (rhs.IsName() || rhs.IsStringKey()) {
		{
			r_34 := newExport
			r_34.nameDecl = j2gort.CastPtr[Var](t.GetScope().GetVar(rhs.GetString()))
		}
	}
	return newExport
}

// Java method: ExportDefinition.getExportPostfix()
func (this_ *ClosureRewriteModule_ExportDefinition) GetExportPostfix() string {
	if this_.exportName == j2gort.NullString {
		return ""
	}
	return ("." + j2gort.StrValString(this_.exportName))
}

// Java method: ExportDefinition.hasInlinableName(java.util.Set<com.google.javascript.jscomp.Var> exportedNames)
func (this_ *ClosureRewriteModule_ExportDefinition) HasInlinableName(exportedNames j2gort.Ext_java_util_Set) bool {
	_ = exportedNames
	if (((this_.nameDecl == nil) || exportedNames.Contains(j2gort.PtrToAny(this_.nameDecl))) || !(j2gort.ImmSetContainsTypedU8_(StaticGet_ClosureRewriteModule_ExportDefinition_INLINABLE_NAME_PARENTS(), uint8(this_.nameDecl.GetParentNode().GetToken())))) || NodeUtil_IsFunctionDeclaration(this_.nameDecl.GetParentNode()) {
		return false
	}
	var initialValue pkg_rhino.Node = this_.nameDecl.GetInitialValue()
	_ = initialValue
	if (initialValue == 0) || !(initialValue.IsCall()) {
		return true
	}
	var method pkg_rhino.Node = initialValue.GetFirstChild()
	_ = method
	if !(method.IsGetProp()) {
		return true
	}
	var maybeGoog pkg_rhino.Node = method.GetFirstChild()
	_ = maybeGoog
	if !(maybeGoog.IsName()) || !(j2gort.CheckNotNullString(maybeGoog.GetString()) == "goog") {
		return true
	}
	var name string = method.GetString()
	_ = name
	return ((!(j2gort.CheckNotNullString(name) == "require") && !(j2gort.CheckNotNullString(name) == "forwardDeclare")) && !(j2gort.CheckNotNullString(name) == "getMsg"))
}

// Java method: ExportDefinition.getLocalName()
func (this_ *ClosureRewriteModule_ExportDefinition) GetLocalName() string {
	if this_.nameDecl != nil {
		return this_.nameDecl.GetName()
	} else {
		return j2gort.NullString
	}
}

// Java class: com.google.javascript.jscomp.ClosureRewriteModule.AliasName
type ClosureRewriteModule_AliasName struct {
	newName     string `gson:"newName"`
	namespaceId string `gson:"namespaceId"`
	idHash_     atomic.Int32
}

func (this_ *ClosureRewriteModule_AliasName) BindSelf_(s any) {
	_ = s
}
func (this_ *ClosureRewriteModule_AliasName) IsJavaNil_() bool { return this_ == nil }
func (*ClosureRewriteModule_AliasName) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.ClosureRewriteModule$AliasName", "AliasName"
}
func (this_ *ClosureRewriteModule_AliasName) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*ClosureRewriteModule_AliasName) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.ClosureRewriteModule.AliasName", "com.google.javascript.jscomp.ClosureRewriteModule$AliasName", func(v_ any) bool {
		return j2gort.InstanceOf[*ClosureRewriteModule_AliasName](v_)
	})
}

func (this_ *ClosureRewriteModule_AliasName) initFields_ClosureRewriteModule_AliasName() {
	_ = this_
	this_.newName = j2gort.NullString
	this_.namespaceId = j2gort.NullString
}

// Java constructor: AliasName(java.lang.String newName, java.lang.@org.jspecify.annotations.Nullable String namespaceId)
func NewClosureRewriteModule_AliasName(newName string, namespaceId string) *ClosureRewriteModule_AliasName {
	this_ := &ClosureRewriteModule_AliasName{}
	InitClosureRewriteModule_AliasName(this_, newName, namespaceId)
	return this_
}

func InitClosureRewriteModule_AliasName(this_ *ClosureRewriteModule_AliasName, newName string, namespaceId string) {
	_ = newName
	_ = namespaceId
	this_.initFields_ClosureRewriteModule_AliasName()
	this_.newName = newName
	this_.namespaceId = namespaceId
}

// Java class: com.google.javascript.jscomp.ClosureRewriteModule.ScriptDescription
type ClosureRewriteModule_ScriptDescription struct {
	namespaceId             string                                                                `gson:"namespaceId"`
	contentsPrefix          string                                                                `gson:"contentsPrefix"`
	topLevelNames           j2gort.OrderedSet[string]                                             `gson:"topLevelNames"`
	childScripts            j2gort.Deque[*ClosureRewriteModule_ScriptDescription]                 `gson:"childScripts"`
	namesToInlineByAlias    j2gort.OrderedMap[string, *ClosureRewriteModule_AliasName]            `gson:"namesToInlineByAlias"`
	defaultExport           *ClosureRewriteModule_ExportDefinition                                `gson:"defaultExport"`
	defaultExportLocalName  string                                                                `gson:"defaultExportLocalName"`
	namedExports            j2gort.OrderedSet[string]                                             `gson:"namedExports"`
	exportsToInline         j2gort.CustomOrderedMap[*Var, *ClosureRewriteModule_ExportDefinition] `gson:"exportsToInline"`
	rootNode                pkg_rhino.Node                                                        `gson:"rootNode"`
	idHash_                 atomic.Int32
	isModule                bool `gson:"isModule"`
	declareLegacyNamespace  bool `gson:"declareLegacyNamespace"`
	willCreateExportsObject bool `gson:"willCreateExportsObject"`
	hasCreatedExportObject  bool `gson:"hasCreatedExportObject"`
}

func (this_ *ClosureRewriteModule_ScriptDescription) BindSelf_(s any) {
	_ = s
}
func (this_ *ClosureRewriteModule_ScriptDescription) IsJavaNil_() bool { return this_ == nil }
func (*ClosureRewriteModule_ScriptDescription) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.ClosureRewriteModule$ScriptDescription", "ScriptDescription"
}
func (this_ *ClosureRewriteModule_ScriptDescription) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*ClosureRewriteModule_ScriptDescription) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.ClosureRewriteModule.ScriptDescription", "com.google.javascript.jscomp.ClosureRewriteModule$ScriptDescription", func(v_ any) bool {
		return j2gort.InstanceOf[*ClosureRewriteModule_ScriptDescription](v_)
	})
}

func (this_ *ClosureRewriteModule_ScriptDescription) initFields_ClosureRewriteModule_ScriptDescription() {
	_ = this_
	this_.namespaceId = j2gort.NullString
	this_.contentsPrefix = j2gort.NullString
	this_.topLevelNames = j2gort.NewOrderedSet[string]()
	this_.childScripts = j2gort.Deque[*ClosureRewriteModule_ScriptDescription]{}
	this_.namesToInlineByAlias = j2gort.NewOrderedMap[string, *ClosureRewriteModule_AliasName]()
	this_.defaultExportLocalName = j2gort.NullString
	this_.namedExports = j2gort.NewOrderedSet[string]()
	this_.exportsToInline = j2gort.NewCustomOrderedMap[*Var, *ClosureRewriteModule_ExportDefinition](

	// Java constructor: ScriptDescription()
	)
}

func NewClosureRewriteModule_ScriptDescription() *ClosureRewriteModule_ScriptDescription {
	this_ := &ClosureRewriteModule_ScriptDescription{}
	InitClosureRewriteModule_ScriptDescription(this_)
	return this_
}

func InitClosureRewriteModule_ScriptDescription(this_ *ClosureRewriteModule_ScriptDescription) {
	this_.initFields_ClosureRewriteModule_ScriptDescription()
}

// Java method: ScriptDescription.addChildScript(com.google.javascript.jscomp.ClosureRewriteModule.ScriptDescription childScript)
func (this_ *ClosureRewriteModule_ScriptDescription) AddChildScript(childScript *ClosureRewriteModule_ScriptDescription) {
	_ = childScript
	this_.childScripts.AddLast(childScript)
}

// Java method: ScriptDescription.removeFirstChildScript()
func (this_ *ClosureRewriteModule_ScriptDescription) RemoveFirstChildScript() *ClosureRewriteModule_ScriptDescription {
	return this_.childScripts.RemoveFirst()
}

// Java method: ScriptDescription.getBinaryNamespace()
func (this_ *ClosureRewriteModule_ScriptDescription) GetBinaryNamespace() string {
	if !(this_.isModule) || this_.declareLegacyNamespace {
		return j2gort.NullString
	}
	return ClosureRewriteModule_GetBinaryModuleNamespace(this_.namespaceId)
}

// Java method: ScriptDescription.getExportedNamespace()
func (this_ *ClosureRewriteModule_ScriptDescription) GetExportedNamespace() string {
	if this_.declareLegacyNamespace {
		return this_.namespaceId
	}
	return this_.GetBinaryNamespace()
}

// Java class: com.google.javascript.jscomp.ClosureRewriteModule.ScriptPreprocessor
type ClosureRewriteModule_ScriptPreprocessor struct {
	NodeTraversal_AbstractPreOrderCallbackBase
	outer_ *ClosureRewriteModule
}

func (this_ *ClosureRewriteModule_ScriptPreprocessor) BindSelf_(s any) {
	_ = s
}
func (this_ *ClosureRewriteModule_ScriptPreprocessor) IsJavaNil_() bool { return this_ == nil }
func (*ClosureRewriteModule_ScriptPreprocessor) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.ClosureRewriteModule$ScriptPreprocessor", "ScriptPreprocessor"
}
func (this_ *ClosureRewriteModule_ScriptPreprocessor) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.NodeTraversal_AbstractPreOrderCallbackBase.IdentityHash_()
}
func (*ClosureRewriteModule_ScriptPreprocessor) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.ClosureRewriteModule.ScriptPreprocessor", "com.google.javascript.jscomp.ClosureRewriteModule$ScriptPreprocessor", func(v_ any) bool {
		return j2gort.InstanceOf[*ClosureRewriteModule_ScriptPreprocessor](v_)
	})
}
func (this_ *ClosureRewriteModule_ScriptPreprocessor) SetOuter_(o *ClosureRewriteModule) {
	this_.outer_ = o
}
func (this_ *ClosureRewriteModule_ScriptPreprocessor) GetOuter_() *ClosureRewriteModule {
	return this_.outer_
}
func (*ClosureRewriteModule_ScriptPreprocessor) Implements_jscomp_NodeTraversal_AbstractPreOrderCallback() {
}
func (*ClosureRewriteModule_ScriptPreprocessor) Implements_jscomp_NodeTraversal_Callback() {}

func (this_ *ClosureRewriteModule_ScriptPreprocessor) initFields_ClosureRewriteModule_ScriptPreprocessor() {
	_ = this_
}

// Java constructor: ScriptPreprocessor()
func NewClosureRewriteModule_ScriptPreprocessor(outer_ *ClosureRewriteModule) *ClosureRewriteModule_ScriptPreprocessor {
	this_ := &ClosureRewriteModule_ScriptPreprocessor{}
	this_.outer_ = outer_
	InitClosureRewriteModule_ScriptPreprocessor(this_)
	return this_
}

func InitClosureRewriteModule_ScriptPreprocessor(this_ *ClosureRewriteModule_ScriptPreprocessor) {
	InitNodeTraversal_AbstractPreOrderCallback(&this_.NodeTraversal_AbstractPreOrderCallbackBase)
	this_.initFields_ClosureRewriteModule_ScriptPreprocessor()
}

// Java method: ScriptPreprocessor.shouldTraverse(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *ClosureRewriteModule_ScriptPreprocessor) ShouldTraverse_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) bool {
	_ = t
	_ = n
	_ = parent
	switch n.TokenSwitchOrd_() {
	case int32(89) /* Token.ROOT */, int32(123) /* Token.MODULE_BODY */ :
		return true
	case int32(94) /* Token.SCRIPT */ :
		if NodeUtil_IsGoogModuleFile(n) {
			ClosureRewriteModule_CheckAndSetStrictModeDirective(t, n)
		}
		return true
	case int32(32) /* Token.NAME */ :
		this_.outer_.priv_ClosureRewriteModule_preprocessExportDeclaration_1189(t, n)
		return true
	default:
		return !(parent.IsScript())
	}
	panic("unreachable switch expression")
}

// Java class: com.google.javascript.jscomp.ClosureRewriteModule.ScriptRecorder
type ClosureRewriteModule_ScriptRecorder struct {
	outer_  *ClosureRewriteModule
	idHash_ atomic.Int32
}

func (this_ *ClosureRewriteModule_ScriptRecorder) BindSelf_(s any) {
	_ = s
}
func (this_ *ClosureRewriteModule_ScriptRecorder) IsJavaNil_() bool { return this_ == nil }
func (*ClosureRewriteModule_ScriptRecorder) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.ClosureRewriteModule$ScriptRecorder", "ScriptRecorder"
}
func (this_ *ClosureRewriteModule_ScriptRecorder) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*ClosureRewriteModule_ScriptRecorder) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.ClosureRewriteModule.ScriptRecorder", "com.google.javascript.jscomp.ClosureRewriteModule$ScriptRecorder", func(v_ any) bool {
		return j2gort.InstanceOf[*ClosureRewriteModule_ScriptRecorder](v_)
	})
}
func (this_ *ClosureRewriteModule_ScriptRecorder) SetOuter_(o *ClosureRewriteModule) {
	this_.outer_ = o
}
func (this_ *ClosureRewriteModule_ScriptRecorder) GetOuter_() *ClosureRewriteModule {
	return this_.outer_
}
func (*ClosureRewriteModule_ScriptRecorder) Implements_jscomp_NodeTraversal_Callback() {}

func (this_ *ClosureRewriteModule_ScriptRecorder) initFields_ClosureRewriteModule_ScriptRecorder() {
	_ = this_
}

// Java constructor: ScriptRecorder()
func NewClosureRewriteModule_ScriptRecorder(outer_ *ClosureRewriteModule) *ClosureRewriteModule_ScriptRecorder {
	this_ := &ClosureRewriteModule_ScriptRecorder{}
	this_.outer_ = outer_
	InitClosureRewriteModule_ScriptRecorder(this_)
	return this_
}

func InitClosureRewriteModule_ScriptRecorder(this_ *ClosureRewriteModule_ScriptRecorder) {
	this_.initFields_ClosureRewriteModule_ScriptRecorder()
}

// Java method: ScriptRecorder.shouldTraverse(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *ClosureRewriteModule_ScriptRecorder) ShouldTraverse_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) bool {
	_ = t
	_ = n
	_ = parent
	switch n.TokenSwitchOrd_() {
	case int32(123) /* Token.MODULE_BODY */ :
		this_.outer_.priv_ClosureRewriteModule_recordModuleBody_1190(n)
	case int32(28) /* Token.CALL */ :
		var method pkg_rhino.Node = n.GetFirstChild()
		_ = method
		if !(method.IsGetProp()) {
			break
		}
		if method.MatchesQualifiedName_1_Node(StaticGet_ClosureRewriteModule_GOOG_MODULE()) {
			this_.outer_.priv_ClosureRewriteModule_recordGoogModule_1191(t, n)
		} else if method.MatchesQualifiedName_1_Node(StaticGet_ClosureRewriteModule_GOOG_MODULE_DECLARELEGACYNAMESPACE()) {
			this_.outer_.priv_ClosureRewriteModule_recordGoogDeclareLegacyNamespace_1192()
		} else if method.MatchesQualifiedName_1_Node(StaticGet_ClosureRewriteModule_GOOG_PROVIDE()) {
			this_.outer_.priv_ClosureRewriteModule_recordGoogProvide_1194(t, n)
		} else if method.MatchesQualifiedName_1_Node(StaticGet_ClosureRewriteModule_GOOG_REQUIRE()) {
			this_.outer_.priv_ClosureRewriteModule_recordGoogRequire_1195(t, n)
		} else if method.MatchesQualifiedName_1_Node(StaticGet_ClosureRewriteModule_GOOG_REQUIRETYPE()) {
			this_.outer_.priv_ClosureRewriteModule_recordGoogRequireType_1196(t, n)
		} else if method.MatchesQualifiedName_1_Node(StaticGet_ClosureRewriteModule_GOOG_REQUIREDYNAMIC()) {
			this_.outer_.priv_ClosureRewriteModule_recordGoogRequireDynamic_1198(t, n)
		} else if method.MatchesQualifiedName_1_Node(StaticGet_ClosureRewriteModule_GOOG_FORWARDDECLARE()) && !(parent.IsExprResult()) {
			this_.outer_.priv_ClosureRewriteModule_recordGoogForwardDeclare_1197(t, n)
		} else if method.MatchesQualifiedName_1_Node(StaticGet_ClosureRewriteModule_GOOG_MODULE_GET()) {
			this_.outer_.priv_ClosureRewriteModule_recordGoogModuleGet_1199(t, n)
		}
	case int32(105) /* Token.CLASS */, int32(73) /* Token.FUNCTION */ :
		if this_.outer_.priv_ClosureRewriteModule_isTopLevel_1229(t, n, StaticGet_ClosureRewriteModule_ScopeType_BLOCK()) {
			this_.outer_.priv_ClosureRewriteModule_recordTopLevelClassOrFunctionName_1200(n)
		}
	case int32(97) /* Token.CONST */, int32(111) /* Token.LET */, int32(84) /* Token.VAR */ :
		if this_.outer_.priv_ClosureRewriteModule_isTopLevel_1229(t, n, func() ClosureRewriteModule_ScopeType {
			if n.IsVar() {
				return StaticGet_ClosureRewriteModule_ScopeType_EXEC_CONTEXT()
			}
			return StaticGet_ClosureRewriteModule_ScopeType_BLOCK()
		}()) {
			this_.outer_.priv_ClosureRewriteModule_recordTopLevelVarNames_1201(n)
		}
	case int32(26) /* Token.GETPROP */ :
		if ClosureRewriteModule_IsExportPropertyAssignment(t, n) {
			this_.outer_.priv_ClosureRewriteModule_recordExportsPropertyAssignment_1213(t, n)
		}
	case int32(32) /* Token.NAME */ :
		this_.outer_.priv_ClosureRewriteModule_maybeRecordExportDeclaration_1202(t, n)
	default:
	}
	return true
}

// Java method: ScriptRecorder.visit(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *ClosureRewriteModule_ScriptRecorder) Visit_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = n
	_ = parent
	if n.IsModuleBody() {
		this_.outer_.priv_ClosureRewriteModule_popScript_1222()
	}
}

// Java class: com.google.javascript.jscomp.ClosureRewriteModule.ScriptUpdater
type ClosureRewriteModule_ScriptUpdater struct {
	outer_             *ClosureRewriteModule
	scriptDescriptions j2gort.Ext_java_util_Deque `gson:"scriptDescriptions"`
	idHash_            atomic.Int32
}

func (this_ *ClosureRewriteModule_ScriptUpdater) BindSelf_(s any) {
	_ = s
}
func (this_ *ClosureRewriteModule_ScriptUpdater) IsJavaNil_() bool { return this_ == nil }
func (*ClosureRewriteModule_ScriptUpdater) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.ClosureRewriteModule$ScriptUpdater", "ScriptUpdater"
}
func (this_ *ClosureRewriteModule_ScriptUpdater) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*ClosureRewriteModule_ScriptUpdater) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.ClosureRewriteModule.ScriptUpdater", "com.google.javascript.jscomp.ClosureRewriteModule$ScriptUpdater", func(v_ any) bool {
		return j2gort.InstanceOf[*ClosureRewriteModule_ScriptUpdater](v_)
	})
}
func (this_ *ClosureRewriteModule_ScriptUpdater) SetOuter_(o *ClosureRewriteModule) { this_.outer_ = o }
func (this_ *ClosureRewriteModule_ScriptUpdater) GetOuter_() *ClosureRewriteModule {
	return this_.outer_
}
func (*ClosureRewriteModule_ScriptUpdater) Implements_jscomp_NodeTraversal_Callback() {}

func (this_ *ClosureRewriteModule_ScriptUpdater) initFields_ClosureRewriteModule_ScriptUpdater() {
	_ = this_
}

// Java constructor: ScriptUpdater(java.util.Deque<com.google.javascript.jscomp.ClosureRewriteModule.ScriptDescription> scriptDescriptions)
func NewClosureRewriteModule_ScriptUpdater(outer_ *ClosureRewriteModule, scriptDescriptions j2gort.Ext_java_util_Deque) *ClosureRewriteModule_ScriptUpdater {
	this_ := &ClosureRewriteModule_ScriptUpdater{}
	this_.outer_ = outer_
	InitClosureRewriteModule_ScriptUpdater(this_, scriptDescriptions)
	return this_
}

func InitClosureRewriteModule_ScriptUpdater(this_ *ClosureRewriteModule_ScriptUpdater, scriptDescriptions j2gort.Ext_java_util_Deque) {
	_ = scriptDescriptions
	this_.initFields_ClosureRewriteModule_ScriptUpdater()
	this_.scriptDescriptions = scriptDescriptions
}

// Java method: ScriptUpdater.shouldTraverse(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *ClosureRewriteModule_ScriptUpdater) ShouldTraverse_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) bool {
	_ = t
	_ = n
	_ = parent
	switch n.TokenSwitchOrd_() {
	case int32(94) /* Token.SCRIPT */ :
		var currentDescription *ClosureRewriteModule_ScriptDescription = j2gort.CastPtr[ClosureRewriteModule_ScriptDescription](this_.scriptDescriptions.RemoveFirst())
		_ = currentDescription
		if !(currentDescription.rootNode == n) {
			j2gort.CheckState(false)
		}
		if n.IsFromExterns() && !(NodeUtil_IsFromTypeSummary(n)) {
			return false
		}
		if !(this_.outer_.scriptStack.IsEmpty()) {
			j2gort.CheckState(false)
		}
		this_.outer_.priv_ClosureRewriteModule_pushScript_1221(currentDescription)
		t.GetScope()
		if this_.outer_.globalScope == nil {
			this_.outer_.globalScope = j2gort.CastPtr[Scope](t.GetScope().GetGlobalScope())
		}
	case int32(123) /* Token.MODULE_BODY */ :
		if parent.GetBooleanProp(pkg_rhino.StaticGet_Node_GOOG_MODULE()) {
			this_.outer_.priv_ClosureRewriteModule_updateModuleBodyEarly_1203(n)
		} else {
			return false
		}
	case int32(28) /* Token.CALL */ :
		var method pkg_rhino.Node = n.GetFirstChild()
		_ = method
		if !(method.IsGetProp()) {
			break
		}
		if method.MatchesQualifiedName_1_Node(StaticGet_ClosureRewriteModule_GOOG_MODULE()) {
			this_.outer_.priv_ClosureRewriteModule_updateGoogModule_1204(t, n)
		} else if method.MatchesQualifiedName_1_Node(StaticGet_ClosureRewriteModule_GOOG_MODULE_DECLARELEGACYNAMESPACE()) {
			ClosureRewriteModule_UpdateGoogDeclareLegacyNamespace(n)
		} else if method.MatchesQualifiedName_1_Node(StaticGet_ClosureRewriteModule_GOOG_REQUIRE()) || method.MatchesQualifiedName_1_Node(StaticGet_ClosureRewriteModule_GOOG_REQUIRETYPE()) {
			this_.outer_.priv_ClosureRewriteModule_updateGoogRequire_1205(t, n)
		} else if method.MatchesQualifiedName_1_Node(StaticGet_ClosureRewriteModule_GOOG_FORWARDDECLARE()) && !(parent.IsExprResult()) {
			this_.outer_.priv_ClosureRewriteModule_updateGoogForwardDeclare_1207(t, n)
		} else if StaticGet_ClosureRewriteModule_GOOG_MODULE_PREVENTMODULEEXPORTSEALING().Matches(method) {
			ClosureRewriteModule_UpdateGoogPreventModuleExportsSealing(n)
		}
	case int32(26) /* Token.GETPROP */ :
		if ClosureRewriteModule_IsExportPropertyAssignment(t, n) {
			this_.outer_.priv_ClosureRewriteModule_updateExportsPropertyAssignment_1214(n, t)
		}
	default:
	}
	if n.GetJSDocInfo() != nil {
		this_.outer_.priv_ClosureRewriteModule_rewriteJsdoc_1181(n.GetJSDocInfo(), t.GetScope())
	}
	return true
}

// Java method: ScriptUpdater.visit(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *ClosureRewriteModule_ScriptUpdater) Visit_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = n
	_ = parent
	switch n.TokenSwitchOrd_() {
	case int32(123) /* Token.MODULE_BODY */ :
		this_.outer_.UpdateModuleBody(n)
	case int32(32) /* Token.NAME */ :
		this_.outer_.priv_ClosureRewriteModule_maybeUpdateTopLevelName_1215(t, n)
		this_.outer_.priv_ClosureRewriteModule_maybeUpdateExportDeclaration_1218(t, n)
		t.GetScope()
		this_.outer_.priv_ClosureRewriteModule_maybeUpdateExportNameRef_1219(t, n)
	case int32(94) /* Token.SCRIPT */ :
		if !(this_.outer_.currentScript.rootNode == n) {
			j2gort.CheckState(false)
		}
		this_.outer_.priv_ClosureRewriteModule_popScript_1222()
	default:
	}
}

// Java class: com.google.javascript.jscomp.ClosureRewriteModule.GlobalRewriteState
type ClosureRewriteModule_GlobalRewriteState struct {
	scriptDescriptionsByGoogModuleNamespace j2gort.OrderedMap[string, *ClosureRewriteModule_ScriptDescription] `gson:"scriptDescriptionsByGoogModuleNamespace"`
	namespaceIdsByScriptNode                j2gort.Ext_com_google_common_collect_Multimap                      `gson:"namespaceIdsByScriptNode"`
	providedNamespaces                      j2gort.OrderedSet[string]                                          `gson:"providedNamespaces"`
	idHash_                                 atomic.Int32
}

func (this_ *ClosureRewriteModule_GlobalRewriteState) BindSelf_(s any) {
	_ = s
}
func (this_ *ClosureRewriteModule_GlobalRewriteState) IsJavaNil_() bool { return this_ == nil }
func (*ClosureRewriteModule_GlobalRewriteState) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.ClosureRewriteModule$GlobalRewriteState", "GlobalRewriteState"
}
func (this_ *ClosureRewriteModule_GlobalRewriteState) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*ClosureRewriteModule_GlobalRewriteState) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.ClosureRewriteModule.GlobalRewriteState", "com.google.javascript.jscomp.ClosureRewriteModule$GlobalRewriteState", func(v_ any) bool {
		return j2gort.InstanceOf[*ClosureRewriteModule_GlobalRewriteState](v_)
	})
}

func (this_ *ClosureRewriteModule_GlobalRewriteState) initFields_ClosureRewriteModule_GlobalRewriteState() {
	_ = this_
	this_.scriptDescriptionsByGoogModuleNamespace = j2gort.NewOrderedMap[string, *ClosureRewriteModule_ScriptDescription]()
	this_.namespaceIdsByScriptNode = j2gort.AsIface[j2gort.Ext_com_google_common_collect_Multimap](j2gort.Ext_com_google_common_collect_HashMultimap_Create_0())
	this_.providedNamespaces = j2gort.NewOrderedSet[string]()
}

// Java constructor: GlobalRewriteState()
func NewClosureRewriteModule_GlobalRewriteState() *ClosureRewriteModule_GlobalRewriteState {
	this_ := &ClosureRewriteModule_GlobalRewriteState{}
	InitClosureRewriteModule_GlobalRewriteState(this_)
	return this_
}

func InitClosureRewriteModule_GlobalRewriteState(this_ *ClosureRewriteModule_GlobalRewriteState) {
	this_.initFields_ClosureRewriteModule_GlobalRewriteState()
}

// Java method: GlobalRewriteState.containsModule(java.lang.String namespaceId)
func (this_ *ClosureRewriteModule_GlobalRewriteState) ContainsModule(namespaceId string) bool {
	_ = namespaceId
	return this_.scriptDescriptionsByGoogModuleNamespace.ContainsKey(namespaceId)
}

// Java method: GlobalRewriteState.isLegacyModule(java.lang.String namespaceId)
func (this_ *ClosureRewriteModule_GlobalRewriteState) IsLegacyModule(namespaceId string) bool {
	_ = namespaceId
	if !(this_.ContainsModule(namespaceId)) {
		j2gort.CheckArgument(false)
	}
	return this_.scriptDescriptionsByGoogModuleNamespace.Get(namespaceId).declareLegacyNamespace
}

// Java method: GlobalRewriteState.getBinaryNamespace(java.lang.String namespaceId)
func (this_ *ClosureRewriteModule_GlobalRewriteState) GetBinaryNamespace(namespaceId string) string {
	_ = namespaceId
	var script *ClosureRewriteModule_ScriptDescription = this_.scriptDescriptionsByGoogModuleNamespace.Get(namespaceId)
	_ = script
	if script == nil {
		return j2gort.NullString
	} else {
		return script.GetBinaryNamespace()
	}
}

// Java method: GlobalRewriteState.getGoogModuleNamespaceType(java.lang.String namespaceId)
func (this_ *ClosureRewriteModule_GlobalRewriteState) GetGoogModuleNamespaceType(namespaceId string) pkg_rhino.JSType {
	_ = namespaceId
	var googModule *ClosureRewriteModule_ScriptDescription = this_.scriptDescriptionsByGoogModuleNamespace.Get(namespaceId)
	_ = googModule
	if googModule == nil {
		return nil
	} else {
		return googModule.rootNode.GetJSType()
	}
}

// Java method: GlobalRewriteState.getExportedNamespaceOrScript(java.lang.String namespaceId)
func (this_ *ClosureRewriteModule_GlobalRewriteState) priv_GlobalRewriteState_getExportedNamespaceOrScript_1244(namespaceId string) string {
	_ = namespaceId
	if this_.providedNamespaces.Contains(namespaceId) {
		return namespaceId
	}
	var script *ClosureRewriteModule_ScriptDescription = this_.scriptDescriptionsByGoogModuleNamespace.Get(namespaceId)
	_ = script
	if script == nil {
		return j2gort.NullString
	} else {
		return script.GetExportedNamespace()
	}
}

// Java class: com.google.javascript.jscomp.ClosureRewriteModule.UnwrapGoogLoadModule
type ClosureRewriteModule_UnwrapGoogLoadModule struct {
	NodeTraversal_AbstractPreOrderCallbackBase
	outer_ *ClosureRewriteModule
}

func (this_ *ClosureRewriteModule_UnwrapGoogLoadModule) BindSelf_(s any) {
	_ = s
}
func (this_ *ClosureRewriteModule_UnwrapGoogLoadModule) IsJavaNil_() bool { return this_ == nil }
func (*ClosureRewriteModule_UnwrapGoogLoadModule) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.ClosureRewriteModule$UnwrapGoogLoadModule", "UnwrapGoogLoadModule"
}
func (this_ *ClosureRewriteModule_UnwrapGoogLoadModule) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.NodeTraversal_AbstractPreOrderCallbackBase.IdentityHash_()
}
func (*ClosureRewriteModule_UnwrapGoogLoadModule) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.ClosureRewriteModule.UnwrapGoogLoadModule", "com.google.javascript.jscomp.ClosureRewriteModule$UnwrapGoogLoadModule", func(v_ any) bool {
		return j2gort.InstanceOf[*ClosureRewriteModule_UnwrapGoogLoadModule](v_)
	})
}
func (this_ *ClosureRewriteModule_UnwrapGoogLoadModule) SetOuter_(o *ClosureRewriteModule) {
	this_.outer_ = o
}
func (this_ *ClosureRewriteModule_UnwrapGoogLoadModule) GetOuter_() *ClosureRewriteModule {
	return this_.outer_
}
func (*ClosureRewriteModule_UnwrapGoogLoadModule) Implements_jscomp_NodeTraversal_AbstractPreOrderCallback() {
}
func (*ClosureRewriteModule_UnwrapGoogLoadModule) Implements_jscomp_NodeTraversal_Callback() {}

func (this_ *ClosureRewriteModule_UnwrapGoogLoadModule) initFields_ClosureRewriteModule_UnwrapGoogLoadModule() {
	_ = this_
}

// Java constructor: UnwrapGoogLoadModule()
func NewClosureRewriteModule_UnwrapGoogLoadModule(outer_ *ClosureRewriteModule) *ClosureRewriteModule_UnwrapGoogLoadModule {
	this_ := &ClosureRewriteModule_UnwrapGoogLoadModule{}
	this_.outer_ = outer_
	InitClosureRewriteModule_UnwrapGoogLoadModule(this_)
	return this_
}

func InitClosureRewriteModule_UnwrapGoogLoadModule(this_ *ClosureRewriteModule_UnwrapGoogLoadModule) {
	InitNodeTraversal_AbstractPreOrderCallback(&this_.NodeTraversal_AbstractPreOrderCallbackBase)
	this_.initFields_ClosureRewriteModule_UnwrapGoogLoadModule()
}

// Java method: UnwrapGoogLoadModule.shouldTraverse(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *ClosureRewriteModule_UnwrapGoogLoadModule) ShouldTraverse_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) bool {
	_ = t
	_ = n
	_ = parent
	switch n.TokenSwitchOrd_() {
	case int32(89) /* Token.ROOT */, int32(94) /* Token.SCRIPT */ :
		return true
	case int32(93) /* Token.EXPR_RESULT */ :
		var call pkg_rhino.Node = n.GetFirstChild()
		_ = call
		if NodeUtil_IsCallTo_2_Node_Node(call, StaticGet_ClosureRewriteModule_GOOG_LOADMODULE()) && call.GetLastChild().IsFunction() {
			parent.PutBooleanProp(pkg_rhino.StaticGet_Node_GOOG_MODULE(), true)
			var functionNode pkg_rhino.Node = call.GetLastChild()
			_ = functionNode
			this_.outer_.compiler.ReportFunctionDeleted(functionNode)
			var moduleBody pkg_rhino.Node = functionNode.GetLastChild().Detach()
			_ = moduleBody
			moduleBody.SetToken(pkg_rhino.StaticGet_Token_MODULE_BODY())
			var exportsParameter pkg_rhino.Node = NodeUtil_GetFunctionParameters(functionNode).GetOnlyChild()
			_ = exportsParameter
			moduleBody.SetJSType(exportsParameter.GetJSType())
			n.ReplaceWith(moduleBody)
			var returnNode pkg_rhino.Node = moduleBody.GetLastChild()
			_ = returnNode
			if !(returnNode.IsReturn()) {
				this_.outer_.compiler.Report_1(JSError_Make_3(moduleBody, StaticGet_ClosureRewriteModule_LOAD_MODULE_FN_MISSING_RETURN(), []string{}))
			} else {
				returnNode.Detach()
			}
			t.ReportCodeChange_0()
		}
		return false
	default:
		return false
	}
	panic("unreachable switch expression")
}
func (*ClosureRewriteModule_ScriptPreprocessor) NoOpVisit_()                   {}
func (*ClosureRewriteModule_UnwrapGoogLoadModule) NoOpVisit_()                 {}
func (*ClosureRewriteModule_UnwrapGoogLoadModule) CanReleaseSyntacticScopes_() {}
