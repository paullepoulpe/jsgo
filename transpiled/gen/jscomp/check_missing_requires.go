/*
 * Copyright 2020 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/CheckMissingRequires.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
	"github.com/paullepoulpe/jsgo/transpiled/shims/gson"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.jscomp.CheckMissingRequires
type CheckMissingRequires struct {
	NodeTraversal_AbstractModuleCallbackBase
	templateParamNames    j2gort.OrderedSet[string]                         `gson:"templateParamNames"`
	moduleByNamespace     j2gort.Ext_com_google_common_collect_ImmutableMap `gson:"moduleByNamespace"`
	controlFlowScopeDepth int32                                             `gson:"controlFlowScopeDepth"`
}

func (this_ *CheckMissingRequires) BindSelf_(s any) {
	_ = s
	this_.NodeTraversal_AbstractModuleCallbackBase.BindSelf_(s)
}
func (this_ *CheckMissingRequires) IsJavaNil_() bool { return this_ == nil }
func (*CheckMissingRequires) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.CheckMissingRequires", "CheckMissingRequires"
}
func (this_ *CheckMissingRequires) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.NodeTraversal_AbstractModuleCallbackBase.IdentityHash_()
}
func (*CheckMissingRequires) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.CheckMissingRequires", "com.google.javascript.jscomp.CheckMissingRequires", func(v_ any) bool {
		return j2gort.InstanceOf[*CheckMissingRequires](v_)
	})
}
func (*CheckMissingRequires) Implements_jscomp_NodeTraversal_AbstractModuleCallback() {}
func (*CheckMissingRequires) Implements_jscomp_CompilerPass()                         {}
func (*CheckMissingRequires) Implements_jscomp_NodeTraversal_Callback()               {}

var static_CheckMissingRequires_MISSING_REQUIRE *DiagnosticType

func StaticGet_CheckMissingRequires_MISSING_REQUIRE() *DiagnosticType {
	clinit_CheckMissingRequires()
	return static_CheckMissingRequires_MISSING_REQUIRE
}
func StaticSet_CheckMissingRequires_MISSING_REQUIRE(v *DiagnosticType) {
	clinit_CheckMissingRequires()
	static_CheckMissingRequires_MISSING_REQUIRE = v
}

var static_CheckMissingRequires_MISSING_REQUIRE_TYPE *DiagnosticType

func StaticGet_CheckMissingRequires_MISSING_REQUIRE_TYPE() *DiagnosticType {
	clinit_CheckMissingRequires()
	return static_CheckMissingRequires_MISSING_REQUIRE_TYPE
}
func StaticSet_CheckMissingRequires_MISSING_REQUIRE_TYPE(v *DiagnosticType) {
	clinit_CheckMissingRequires()
	static_CheckMissingRequires_MISSING_REQUIRE_TYPE = v
}

var static_CheckMissingRequires_INCORRECT_NAMESPACE_ALIAS_REQUIRE *DiagnosticType

func StaticGet_CheckMissingRequires_INCORRECT_NAMESPACE_ALIAS_REQUIRE() *DiagnosticType {
	clinit_CheckMissingRequires()
	return static_CheckMissingRequires_INCORRECT_NAMESPACE_ALIAS_REQUIRE
}
func StaticSet_CheckMissingRequires_INCORRECT_NAMESPACE_ALIAS_REQUIRE(v *DiagnosticType) {
	clinit_CheckMissingRequires()
	static_CheckMissingRequires_INCORRECT_NAMESPACE_ALIAS_REQUIRE = v
}

var static_CheckMissingRequires_INCORRECT_NAMESPACE_ALIAS_REQUIRE_TYPE *DiagnosticType

func StaticGet_CheckMissingRequires_INCORRECT_NAMESPACE_ALIAS_REQUIRE_TYPE() *DiagnosticType {
	clinit_CheckMissingRequires()
	return static_CheckMissingRequires_INCORRECT_NAMESPACE_ALIAS_REQUIRE_TYPE
}
func StaticSet_CheckMissingRequires_INCORRECT_NAMESPACE_ALIAS_REQUIRE_TYPE(v *DiagnosticType) {
	clinit_CheckMissingRequires()
	static_CheckMissingRequires_INCORRECT_NAMESPACE_ALIAS_REQUIRE_TYPE = v
}

var static_CheckMissingRequires_INDIRECT_NAMESPACE_REF_REQUIRE *DiagnosticType

func StaticGet_CheckMissingRequires_INDIRECT_NAMESPACE_REF_REQUIRE() *DiagnosticType {
	clinit_CheckMissingRequires()
	return static_CheckMissingRequires_INDIRECT_NAMESPACE_REF_REQUIRE
}
func StaticSet_CheckMissingRequires_INDIRECT_NAMESPACE_REF_REQUIRE(v *DiagnosticType) {
	clinit_CheckMissingRequires()
	static_CheckMissingRequires_INDIRECT_NAMESPACE_REF_REQUIRE = v
}

var static_CheckMissingRequires_INDIRECT_NAMESPACE_REF_REQUIRE_TYPE *DiagnosticType

func StaticGet_CheckMissingRequires_INDIRECT_NAMESPACE_REF_REQUIRE_TYPE() *DiagnosticType {
	clinit_CheckMissingRequires()
	return static_CheckMissingRequires_INDIRECT_NAMESPACE_REF_REQUIRE_TYPE
}
func StaticSet_CheckMissingRequires_INDIRECT_NAMESPACE_REF_REQUIRE_TYPE(v *DiagnosticType) {
	clinit_CheckMissingRequires()
	static_CheckMissingRequires_INDIRECT_NAMESPACE_REF_REQUIRE_TYPE = v
}

var static_CheckMissingRequires_MISSING_REQUIRE_IN_PROVIDES_FILE *DiagnosticType

func StaticGet_CheckMissingRequires_MISSING_REQUIRE_IN_PROVIDES_FILE() *DiagnosticType {
	clinit_CheckMissingRequires()
	return static_CheckMissingRequires_MISSING_REQUIRE_IN_PROVIDES_FILE
}
func StaticSet_CheckMissingRequires_MISSING_REQUIRE_IN_PROVIDES_FILE(v *DiagnosticType) {
	clinit_CheckMissingRequires()
	static_CheckMissingRequires_MISSING_REQUIRE_IN_PROVIDES_FILE = v
}

var static_CheckMissingRequires_MISSING_REQUIRE_TYPE_IN_PROVIDES_FILE *DiagnosticType

func StaticGet_CheckMissingRequires_MISSING_REQUIRE_TYPE_IN_PROVIDES_FILE() *DiagnosticType {
	clinit_CheckMissingRequires()
	return static_CheckMissingRequires_MISSING_REQUIRE_TYPE_IN_PROVIDES_FILE
}
func StaticSet_CheckMissingRequires_MISSING_REQUIRE_TYPE_IN_PROVIDES_FILE(v *DiagnosticType) {
	clinit_CheckMissingRequires()
	static_CheckMissingRequires_MISSING_REQUIRE_TYPE_IN_PROVIDES_FILE = v
}

var static_CheckMissingRequires_NON_LEGACY_GOOG_MODULE_REFERENCE *DiagnosticType

func StaticGet_CheckMissingRequires_NON_LEGACY_GOOG_MODULE_REFERENCE() *DiagnosticType {
	clinit_CheckMissingRequires()
	return static_CheckMissingRequires_NON_LEGACY_GOOG_MODULE_REFERENCE
}
func StaticSet_CheckMissingRequires_NON_LEGACY_GOOG_MODULE_REFERENCE(v *DiagnosticType) {
	clinit_CheckMissingRequires()
	static_CheckMissingRequires_NON_LEGACY_GOOG_MODULE_REFERENCE = v
}

var static_CheckMissingRequires_MISSING_REQUIRE_IN_GOOG_SCOPE *DiagnosticType

func StaticGet_CheckMissingRequires_MISSING_REQUIRE_IN_GOOG_SCOPE() *DiagnosticType {
	clinit_CheckMissingRequires()
	return static_CheckMissingRequires_MISSING_REQUIRE_IN_GOOG_SCOPE
}
func StaticSet_CheckMissingRequires_MISSING_REQUIRE_IN_GOOG_SCOPE(v *DiagnosticType) {
	clinit_CheckMissingRequires()
	static_CheckMissingRequires_MISSING_REQUIRE_IN_GOOG_SCOPE = v
}

var static_CheckMissingRequires_MISSING_REQUIRE_TYPE_IN_GOOG_SCOPE *DiagnosticType

func StaticGet_CheckMissingRequires_MISSING_REQUIRE_TYPE_IN_GOOG_SCOPE() *DiagnosticType {
	clinit_CheckMissingRequires()
	return static_CheckMissingRequires_MISSING_REQUIRE_TYPE_IN_GOOG_SCOPE
}
func StaticSet_CheckMissingRequires_MISSING_REQUIRE_TYPE_IN_GOOG_SCOPE(v *DiagnosticType) {
	clinit_CheckMissingRequires()
	static_CheckMissingRequires_MISSING_REQUIRE_TYPE_IN_GOOG_SCOPE = v
}

var static_CheckMissingRequires_MISSING_REQUIRE_FOR_GOOG_MODULE_GET *DiagnosticType

func StaticGet_CheckMissingRequires_MISSING_REQUIRE_FOR_GOOG_MODULE_GET() *DiagnosticType {
	clinit_CheckMissingRequires()
	return static_CheckMissingRequires_MISSING_REQUIRE_FOR_GOOG_MODULE_GET
}
func StaticSet_CheckMissingRequires_MISSING_REQUIRE_FOR_GOOG_MODULE_GET(v *DiagnosticType) {
	clinit_CheckMissingRequires()
	static_CheckMissingRequires_MISSING_REQUIRE_FOR_GOOG_MODULE_GET = v
}

var static_CheckMissingRequires_GOOG_SCOPE pkg_rhino.QualifiedName

func StaticGet_CheckMissingRequires_GOOG_SCOPE() pkg_rhino.QualifiedName {
	clinit_CheckMissingRequires()
	return static_CheckMissingRequires_GOOG_SCOPE
}
func StaticSet_CheckMissingRequires_GOOG_SCOPE(v pkg_rhino.QualifiedName) {
	clinit_CheckMissingRequires()
	static_CheckMissingRequires_GOOG_SCOPE = v
}

var clinitDone_CheckMissingRequires uint32
var clinitState_CheckMissingRequires j2gort.ClinitState

func clinit_CheckMissingRequires() {
	if atomic.LoadUint32(&clinitDone_CheckMissingRequires) != 0 {
		return
	}
	clinitSlow_CheckMissingRequires()
}

func clinitSlow_CheckMissingRequires() {
	if !j2gort.EnterClinit(&clinitState_CheckMissingRequires, &clinitDone_CheckMissingRequires) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_CheckMissingRequires, &clinitDone_CheckMissingRequires)
	static_CheckMissingRequires_MISSING_REQUIRE = DiagnosticType_Warning("JSC_MISSING_REQUIRE", "''{0}'' references a fully qualified namespace, which is disallowed by the style guide.\nPlease add a goog.require, assign or destructure it into an alias, and use the alias instead.\nNamespace defined in {1}.")
	static_CheckMissingRequires_MISSING_REQUIRE_TYPE = DiagnosticType_Disabled("JSC_MISSING_REQUIRE_TYPE", "''{0}'' references a fully qualified namespace, which is disallowed by the style guide.\nPlease add a goog.requireType, assign or destructure it into an alias, and use the alias instead.\nNamespace defined in {1}.")
	static_CheckMissingRequires_INCORRECT_NAMESPACE_ALIAS_REQUIRE = DiagnosticType_Disabled("JSC_INCORRECT_NAMESPACE_ALIAS_REQUIRE", "''{0}'' is its own namespace.\nPlease add a separate goog.require and use that alias instead.")
	static_CheckMissingRequires_INCORRECT_NAMESPACE_ALIAS_REQUIRE_TYPE = DiagnosticType_Disabled("JSC_INCORRECT_NAMESPACE_ALIAS_REQUIRE_TYPE", "''{0}'' is its own namespace.\nPlease add a separate goog.requireType and use that alias instead.")
	static_CheckMissingRequires_INDIRECT_NAMESPACE_REF_REQUIRE = DiagnosticType_Disabled("JSC_INDIRECT_NAMESPACE_REF_REQUIRE", "''{0}'' should have its own goog.require.\nPlease add a separate goog.require and use that alias instead.")
	static_CheckMissingRequires_INDIRECT_NAMESPACE_REF_REQUIRE_TYPE = DiagnosticType_Disabled("JSC_INDIRECT_NAMESPACE_REF_REQUIRE_TYPE", "''{0}'' should have its own goog.requireType.\nPlease add a separate goog.requireType and use that alias instead.")
	static_CheckMissingRequires_MISSING_REQUIRE_IN_PROVIDES_FILE = DiagnosticType_Warning("JSC_MISSING_REQUIRE_IN_PROVIDES_FILE", "''{0}'' references a namespace which was not required by this file.\nPlease add a goog.require.\nNamespace defined in {1}.")
	static_CheckMissingRequires_MISSING_REQUIRE_TYPE_IN_PROVIDES_FILE = DiagnosticType_Disabled("JSC_MISSING_REQUIRE_TYPE_IN_PROVIDES_FILE", "''{0}'' references a namespace which was not required by this file.\nPlease add a goog.requireType.\nNamespace defined in {1}.")
	static_CheckMissingRequires_NON_LEGACY_GOOG_MODULE_REFERENCE = DiagnosticType_Error("JSC_NON_LEGACY_GOOG_MODULE_REFERENCE", "''{0}'' references the name of a module without goog.declareLegacyNamespace(), which is not actually defined. Use goog.module.get() instead.")
	static_CheckMissingRequires_MISSING_REQUIRE_IN_GOOG_SCOPE = DiagnosticType_Warning("JSC_MISSING_REQUIRE_IN_GOOG_SCOPE", "''{0}'' is its own namespace. Please add a goog.require and reference that alias.")
	static_CheckMissingRequires_MISSING_REQUIRE_TYPE_IN_GOOG_SCOPE = DiagnosticType_Warning("JSC_MISSING_REQUIRE_TYPE_IN_GOOG_SCOPE ", "''{0}'' is its own namespace. Please add a goog.requireType and reference that alias.")
	static_CheckMissingRequires_MISSING_REQUIRE_FOR_GOOG_MODULE_GET = DiagnosticType_Warning("JSC_MISSING_REQUIRE_FOR_GOOG_MODULE_GET", "''{0}'' references a namespace which was not required by this file.\nPlease add a goog.require.")
	static_CheckMissingRequires_GOOG_SCOPE = pkg_rhino.QualifiedName_Of("goog.scope")
}

func (this_ *CheckMissingRequires) initFields_CheckMissingRequires() {
	_ = this_
	this_.templateParamNames = j2gort.NewOrderedSet[string]()

}

// Java constructor: CheckMissingRequires(com.google.javascript.jscomp.AbstractCompiler compiler, com.google.javascript.jscomp.modules.ModuleMetadataMap moduleMetadataMap)
func NewCheckMissingRequires(compiler AbstractCompiler, moduleMetadataMap *ModuleMetadataMap) *CheckMissingRequires {
	clinit_CheckMissingRequires()
	this_ := &CheckMissingRequires{}
	this_.BindSelf_(this_)
	InitCheckMissingRequires(this_, compiler, moduleMetadataMap)
	return this_
}

func InitCheckMissingRequires(this_ *CheckMissingRequires, compiler AbstractCompiler, moduleMetadataMap *ModuleMetadataMap) {
	_ = compiler
	_ = moduleMetadataMap
	InitNodeTraversal_AbstractModuleCallback(&this_.NodeTraversal_AbstractModuleCallbackBase, compiler, moduleMetadataMap)
	this_.initFields_CheckMissingRequires()
	this_.moduleByNamespace = moduleMetadataMap.GetModulesByGoogNamespace()
}

// Java method: CheckMissingRequires.process(com.google.javascript.rhino.Node externs, com.google.javascript.rhino.Node root)
func (this_ *CheckMissingRequires) Process_2(externs pkg_rhino.Node, root pkg_rhino.Node) {
	_ = externs
	_ = root
	sc_ := NewSyntacticScopeCreator_1(this_.GetNodeTraversal_AbstractModuleCallbackBase().compiler)
	defer func() {
		this_.templateParamNames.Clear()
		sc_.ReleaseSlabs_()
	}()
	{
		clinit_NodeTraversal()
		var ntb_ NodeTraversal_Builder
		InitNodeTraversal_Builder(&ntb_)
		(&ntb_).SetCompiler(this_.GetNodeTraversal_AbstractModuleCallbackBase().compiler)
		(&ntb_).SetCallback_1_NodeTraversal_Callback(this_)
		(&ntb_).SetScopeCreator(sc_)
		(&ntb_).Traverse(root)
	}

}

// Java method: CheckMissingRequires.shouldTraverse(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.jscomp.modules.ModuleMetadataMap.@org.jspecify.annotations.Nullable ModuleMetadata currentModule, com.google.javascript.rhino.Node scopeRoot)
func (this_ *CheckMissingRequires) ShouldTraverse_4(t *NodeTraversal, n pkg_rhino.Node, currentModule ModuleMetadataMap_ModuleMetadata, scopeRoot pkg_rhino.Node) bool {
	_ = t
	_ = n
	_ = currentModule
	_ = scopeRoot
	if NodeUtil_IsValidCfgRoot(n) && !(this_.priv_CheckMissingRequires_isInGlobalControlFlowScope_1086(n)) {
		this_.controlFlowScopeDepth = this_.controlFlowScopeDepth + 1
	}
	if currentModule == nil || currentModule.IsJavaNil_() {
		return true
	}
	this_.priv_CheckMissingRequires_visitNode_1087(t, n, currentModule)
	return true
}

// Java method: CheckMissingRequires.isInGlobalControlFlowScope(com.google.javascript.rhino.Node n)
func (this_ *CheckMissingRequires) priv_CheckMissingRequires_isInGlobalControlFlowScope_1086(n pkg_rhino.Node) bool {
	_ = n
	switch n.TokenSwitchOrd_() {
	case int32(94) /* Token.SCRIPT */, int32(89) /* Token.ROOT */ :
		return true
	case int32(123) /* Token.MODULE_BODY */, int32(90) /* Token.BLOCK */ :
		return false
	case int32(73) /* Token.FUNCTION */ :
		if this_.controlFlowScopeDepth > int32(0) {
			return false
		}
		if CheckMissingRequires_IsGoogScopeBody(NodeUtil_GetFunctionBody(n)) {
			return true
		}
		return CheckMissingRequires_IsIIFE(n)
	default:
		panic(j2gort.AsIface[javalang.IAssertionError](j2gort.ConstructShimThrowable("java.lang.AssertionError", "Unexpected control flow root: "+j2gort.StrVal(n))))
	}
	panic("unreachable switch expression")
}

// Java static method: CheckMissingRequires.isIIFE(com.google.javascript.rhino.Node n)
func CheckMissingRequires_IsIIFE(n pkg_rhino.Node) bool {
	clinit_CheckMissingRequires()
	_ = n
	return ((n.IsFunction() && n.GetParent().IsCall()) && n.IsFirstChildOf(n.GetParent()))
}

// Java method: CheckMissingRequires.visit(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.jscomp.modules.ModuleMetadataMap.@org.jspecify.annotations.Nullable ModuleMetadata currentModule, com.google.javascript.rhino.@org.jspecify.annotations.Nullable Node scopeRoot)
func (this_ *CheckMissingRequires) Visit_4(t *NodeTraversal, n pkg_rhino.Node, currentModule ModuleMetadataMap_ModuleMetadata, scopeRoot pkg_rhino.Node) {
	_ = t
	_ = n
	_ = currentModule
	_ = scopeRoot
	if NodeUtil_IsValidCfgRoot(n) && !(this_.priv_CheckMissingRequires_isInGlobalControlFlowScope_1086(n)) {
		this_.controlFlowScopeDepth = this_.controlFlowScopeDepth - 1
	}
	if (currentModule != nil && !currentModule.IsJavaNil_()) && (n == currentModule.RootNode()) {
		this_.templateParamNames.Clear()
	}
}

// Java method: CheckMissingRequires.visitNode(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.jscomp.modules.ModuleMetadataMap.ModuleMetadata currentModule)
func (this_ *CheckMissingRequires) priv_CheckMissingRequires_visitNode_1087(t *NodeTraversal, n pkg_rhino.Node, currentModule ModuleMetadataMap_ModuleMetadata) {
	_ = t
	_ = n
	_ = currentModule
	var info *pkg_rhino.JSDocInfo = n.GetJSDocInfo()
	_ = info
	if info != nil {
		this_.priv_CheckMissingRequires_visitJsDocInfo_1088(t, currentModule, info)
	}
	if !(n.GetParent().IsGetProp()) && n.IsQualifiedName() {
		var qualifiedName pkg_rhino.QualifiedName = n.GetQualifiedNameObject()
		_ = qualifiedName
		var root string = qualifiedName.GetRoot()
		_ = root
		if (j2gort.CheckNotNullString(root) == "this") || (j2gort.CheckNotNullString(root) == "super") {
			return
		}
		this_.priv_CheckMissingRequires_visitQualifiedName_1093(t, n, currentModule, qualifiedName, StaticGet_CheckMissingRequires_Strength_CODE())
	}
	if n.IsName() && !(j2gort.CheckNotNullString(n.GetString()) == "") {
		this_.priv_CheckMissingRequires_visitMaybeDeclaration_1091(t, n, currentModule)
	}
	if n.IsCall() && (this_.controlFlowScopeDepth == int32(0)) {
		this_.priv_CheckMissingRequires_visitMaybeGoogModuleGet_1090(t, n, currentModule)
	}
}

// Java method: CheckMissingRequires.visitJsDocInfo(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.jscomp.modules.ModuleMetadataMap.ModuleMetadata currentModule, com.google.javascript.rhino.JSDocInfo info)
func (this_ *CheckMissingRequires) priv_CheckMissingRequires_visitJsDocInfo_1088(t *NodeTraversal, currentModule ModuleMetadataMap_ModuleMetadata, info *pkg_rhino.JSDocInfo) {
	_ = t
	_ = currentModule
	_ = info
	j2gort.OrderedSetAddAllStr(&this_.templateParamNames, j2gort.AsIface[j2gort.Ext_java_util_Collection](info.GetTemplateTypeNames()))
	j2gort.OrderedSetAddAllStr(&this_.templateParamNames, j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableSet](info.GetTypeTransformations().KeySet_0()))
	if info.HasType() {
		this_.priv_CheckMissingRequires_visitJsDocExpr_1089(t, currentModule, info.GetType(), StaticGet_CheckMissingRequires_Strength_WEAK_TYPE())
	}
	{
		var i int32 = int32(0)
		_ = i
		var paramCount int32 = info.GetParameterCount()
		_ = paramCount
		for ; i < paramCount; i = i + 1 {
			var paramType pkg_rhino.JSTypeExpression = info.GetParameterTypeAt(i)
			_ = paramType
			if paramType != (pkg_rhino.JSTypeExpression{}) {
				this_.priv_CheckMissingRequires_visitJsDocExpr_1089(t, currentModule, paramType, StaticGet_CheckMissingRequires_Strength_WEAK_TYPE())
			}
		}
	}
	if info.HasReturnType() {
		this_.priv_CheckMissingRequires_visitJsDocExpr_1089(t, currentModule, info.GetReturnType(), StaticGet_CheckMissingRequires_Strength_WEAK_TYPE())
	}
	if info.HasEnumParameterType() {
		this_.priv_CheckMissingRequires_visitJsDocExpr_1089(t, currentModule, info.GetEnumParameterType(), StaticGet_CheckMissingRequires_Strength_WEAK_TYPE())
	}
	if info.HasTypedefType() {
		this_.priv_CheckMissingRequires_visitJsDocExpr_1089(t, currentModule, info.GetTypedefType(), StaticGet_CheckMissingRequires_Strength_WEAK_TYPE())
	}
	if info.HasThisType() {
		this_.priv_CheckMissingRequires_visitJsDocExpr_1089(t, currentModule, info.GetThisType(), StaticGet_CheckMissingRequires_Strength_WEAK_TYPE())
	}
	if info.HasBaseType() {
		this_.priv_CheckMissingRequires_visitJsDocExpr_1089(t, currentModule, info.GetBaseType(), StaticGet_CheckMissingRequires_Strength_IMPLEMENTS_EXTENDS())
	}
	{
		var i_2 int32 = int32(0)
		_ = i_2
		var n int32 = info.GetExtendedInterfacesCount()
		_ = n
		for ; i_2 < n; i_2 = i_2 + 1 {
			this_.priv_CheckMissingRequires_visitJsDocExpr_1089(t, currentModule, info.GetExtendedInterfaceAt(i_2), StaticGet_CheckMissingRequires_Strength_IMPLEMENTS_EXTENDS())
		}
	}
	{
		var i_3 int32 = int32(0)
		_ = i_3
		var n_2 int32 = info.GetImplementedInterfaceCount()
		_ = n_2
		for ; i_3 < n_2; i_3 = i_3 + 1 {
			this_.priv_CheckMissingRequires_visitJsDocExpr_1089(t, currentModule, info.GetImplementedInterfaceAt(i_3), StaticGet_CheckMissingRequires_Strength_IMPLEMENTS_EXTENDS())
		}
	}
}

// Java method: CheckMissingRequires.visitJsDocExpr(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.jscomp.modules.ModuleMetadataMap.ModuleMetadata currentModule, com.google.javascript.rhino.JSTypeExpression expr, com.google.javascript.jscomp.CheckMissingRequires.Strength referenceStrength)
func (this_ *CheckMissingRequires) priv_CheckMissingRequires_visitJsDocExpr_1089(t *NodeTraversal, currentModule ModuleMetadataMap_ModuleMetadata, expr pkg_rhino.JSTypeExpression, referenceStrength CheckMissingRequires_Strength) {
	_ = t
	_ = currentModule
	_ = expr
	_ = referenceStrength
	{
		cur_1 := j2gort.NewIterCursor(j2gort.CheckNotNull(expr.GetAllTypeNodes()))
		for cur_1.HasNext() {
			typeNode := pkg_rhino.Node(cur_1.NextNode())
			_ = typeNode
			this_.priv_CheckMissingRequires_visitQualifiedName_1093(t, typeNode, currentModule, pkg_rhino.QualifiedName_Of(typeNode.GetString()), referenceStrength)
		}
	}
}

// Java static method: CheckMissingRequires.isGoogScopeBody(com.google.javascript.rhino.Node hoistScopeRoot)
func CheckMissingRequires_IsGoogScopeBody(hoistScopeRoot pkg_rhino.Node) bool {
	clinit_CheckMissingRequires()
	_ = hoistScopeRoot
	return (((hoistScopeRoot.IsBlock() && hoistScopeRoot.GetParent().IsFunction()) && hoistScopeRoot.GetGrandparent().IsCall()) && StaticGet_CheckMissingRequires_GOOG_SCOPE().Matches(hoistScopeRoot.GetGrandparent().GetFirstChild()))
}

// Java method: CheckMissingRequires.visitMaybeGoogModuleGet(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node googModuleGet, com.google.javascript.jscomp.modules.ModuleMetadataMap.ModuleMetadata currentFile)
func (this_ *CheckMissingRequires) priv_CheckMissingRequires_visitMaybeGoogModuleGet_1090(t *NodeTraversal, googModuleGet pkg_rhino.Node, currentFile ModuleMetadataMap_ModuleMetadata) {
	_ = t
	_ = googModuleGet
	_ = currentFile
	if !(googModuleGet.IsCall()) {
		j2gort.CheckState(false)
	}
	if !(NodeUtil_IsGoogModuleGetCall(googModuleGet)) {
		return
	}
	var importedNamespace string = googModuleGet.GetSecondChild().GetString()
	_ = importedNamespace
	var requiredFile ModuleMetadataMap_ModuleMetadata = j2gort.AsIface[ModuleMetadataMap_ModuleMetadata](j2gort.MapGetStr(this_.moduleByNamespace, importedNamespace))
	_ = requiredFile
	if requiredFile == nil || requiredFile.IsJavaNil_() {
		return
	}
	if !(CheckMissingRequires_HasAcceptableRequire(currentFile, pkg_rhino.QualifiedName_Of(importedNamespace), requiredFile, StaticGet_CheckMissingRequires_Strength_WEAK_TYPE())) {
		t.Report_3(googModuleGet, StaticGet_CheckMissingRequires_MISSING_REQUIRE_FOR_GOOG_MODULE_GET(), []string{importedNamespace})
	}
}

// Java method: CheckMissingRequires.visitMaybeDeclaration(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.jscomp.modules.ModuleMetadataMap.ModuleMetadata currentFile)
func (this_ *CheckMissingRequires) priv_CheckMissingRequires_visitMaybeDeclaration_1091(t *NodeTraversal, n pkg_rhino.Node, currentFile ModuleMetadataMap_ModuleMetadata) {
	_ = t
	_ = n
	_ = currentFile
	if !(currentFile.IsModule()) {
		return
	}
	if this_.priv_CheckMissingRequires_isTypeScriptSource_1092(n) {
		return
	}
	if !(n.IsName()) {
		j2gort.CheckState(false)
	}
	var rootName string = n.GetString()
	_ = rootName
	var var_ *Var = j2gort.CastPtr[Var](t.GetScope().GetVar(rootName))
	_ = var_
	if var_ == nil {
		return
	}
	var declaration *Var = j2gort.CastPtr[Var](var_.GetDeclaration())
	_ = declaration
	if (declaration == nil) || (declaration.GetNode() != n) {
		return
	}
	if j2gort.CastPtr[Scope](var_.GetScope()).IsLocal() {
		var require *NodeUtil_GoogRequire = NodeUtil_GetGoogRequireInfo_1_Var(var_)
		_ = require
		if require == nil {
			return
		}
		if require.Property() == j2gort.NullString {
			return
		}
		var requiredFile ModuleMetadataMap_ModuleMetadata = j2gort.AsIface[ModuleMetadataMap_ModuleMetadata](j2gort.MapGetStr(this_.moduleByNamespace, require.Namespace()))
		_ = requiredFile
		if (requiredFile == nil || requiredFile.IsJavaNil_()) || !(requiredFile.HasLegacyGoogNamespaces()) {
			return
		}
		var reference string = (j2gort.StrValString(require.Namespace()) + "." + j2gort.StrValString(require.Property()))
		_ = reference
		if !(j2gort.SetContainsStr(currentFile.GoogNamespaces(), reference)) {
			var alternateFile ModuleMetadataMap_ModuleMetadata = j2gort.AsIface[ModuleMetadataMap_ModuleMetadata](j2gort.MapGetStr(this_.moduleByNamespace, reference))
			_ = alternateFile
			if (alternateFile != nil && !alternateFile.IsJavaNil_()) && !(alternateFile == requiredFile || ((alternateFile == nil || alternateFile.IsJavaNil_()) && (requiredFile == nil || requiredFile.IsJavaNil_()))) {
				var toReport *DiagnosticType
				_ = toReport
				if require.IsStrongRequire() {
					toReport = StaticGet_CheckMissingRequires_INCORRECT_NAMESPACE_ALIAS_REQUIRE()
				} else {
					toReport = StaticGet_CheckMissingRequires_INCORRECT_NAMESPACE_ALIAS_REQUIRE_TYPE()
				}
				t.Report_3(n, toReport, []string{reference})
			}
		}
	}
}

// Java method: CheckMissingRequires.isTypeScriptSource(com.google.javascript.rhino.Node n)
func (this_ *CheckMissingRequires) priv_CheckMissingRequires_isTypeScriptSource_1092(n pkg_rhino.Node) bool {
	_ = n
	return n.GetStaticSourceFile().IsTypeScriptSource()
}

// Java method: CheckMissingRequires.visitQualifiedName(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.jscomp.modules.ModuleMetadataMap.ModuleMetadata currentFile, com.google.javascript.rhino.QualifiedName qualifiedName, com.google.javascript.jscomp.CheckMissingRequires.Strength referenceStrength)
func (this_ *CheckMissingRequires) priv_CheckMissingRequires_visitQualifiedName_1093(t *NodeTraversal, n pkg_rhino.Node, currentFile ModuleMetadataMap_ModuleMetadata, qualifiedName pkg_rhino.QualifiedName, referenceStrength CheckMissingRequires_Strength) {
	_ = t
	_ = n
	_ = currentFile
	_ = qualifiedName
	_ = referenceStrength
	var rootName string = qualifiedName.GetRoot()
	_ = rootName
	if qualifiedName.IsSimple() {
		if this_.templateParamNames.Contains(rootName) {
			return
		}
	}
	if NodeUtil_IsDeclarationLValue(n) {
		return
	}
	var var_ *Var = j2gort.CastPtr[Var](t.GetScope().GetVar(rootName))
	_ = var_
	if (var_ != nil) && j2gort.CastPtr[Scope](var_.GetScope()).IsLocal() {
		this_.priv_CheckMissingRequires_checkMissingRequireThroughShortName_1094(t, n, var_, currentFile, qualifiedName, referenceStrength)
	} else {
		this_.priv_CheckMissingRequires_checkMissingRequireThroughFullyQualifiedName_1098(t, n, currentFile, qualifiedName, referenceStrength)
	}
}

// Java method: CheckMissingRequires.checkMissingRequireThroughShortName(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.jscomp.Var localVar, com.google.javascript.jscomp.modules.ModuleMetadataMap.ModuleMetadata currentFile, com.google.javascript.rhino.QualifiedName qualifiedName, com.google.javascript.jscomp.CheckMissingRequires.Strength referenceStrength)
func (this_ *CheckMissingRequires) priv_CheckMissingRequires_checkMissingRequireThroughShortName_1094(t *NodeTraversal, n pkg_rhino.Node, localVar *Var, currentFile ModuleMetadataMap_ModuleMetadata, qualifiedName pkg_rhino.QualifiedName, referenceStrength CheckMissingRequires_Strength) {
	_ = t
	_ = n
	_ = localVar
	_ = currentFile
	_ = qualifiedName
	_ = referenceStrength
	if this_.priv_CheckMissingRequires_isTypeScriptSource_1092(n) {
		return
	}
	var inModule bool = currentFile.IsModule()
	_ = inModule
	var inGoogScope bool = (!(inModule) && CheckMissingRequires_IsGoogScopeBody(localVar.GetScopeRoot()))
	_ = inGoogScope
	if !(inModule) && !(inGoogScope) {
		return
	}
	if qualifiedName.IsSimple() {
		return
	}
	if inGoogScope {
		this_.priv_CheckMissingRequires_checkMissingRequireThroughShortNameInGoogScope_1095(t, n, localVar, currentFile, qualifiedName, referenceStrength)
	} else {
		this_.priv_CheckMissingRequires_checkMissingRequireThroughShortNameInModule_1096(t, n, NodeUtil_GetGoogRequireInfo_1_Var(localVar), currentFile, qualifiedName, referenceStrength)
	}
}

// Java method: CheckMissingRequires.checkMissingRequireThroughShortNameInGoogScope(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.jscomp.Var localVar, com.google.javascript.jscomp.modules.ModuleMetadataMap.ModuleMetadata currentFile, com.google.javascript.rhino.QualifiedName qualifiedName, com.google.javascript.jscomp.CheckMissingRequires.Strength referenceStrength)
func (this_ *CheckMissingRequires) priv_CheckMissingRequires_checkMissingRequireThroughShortNameInGoogScope_1095(t *NodeTraversal, n pkg_rhino.Node, localVar *Var, currentFile ModuleMetadataMap_ModuleMetadata, qualifiedName pkg_rhino.QualifiedName, referenceStrength CheckMissingRequires_Strength) {
	_ = t
	_ = n
	_ = localVar
	_ = currentFile
	_ = qualifiedName
	_ = referenceStrength
	var initialValue pkg_rhino.Node = localVar.GetInitialValue()
	_ = initialValue
	if initialValue == 0 {
		return
	}
	if initialValue.IsQualifiedName() {
		var aliasedName pkg_rhino.QualifiedName = initialValue.GetQualifiedNameObject()
		_ = aliasedName
		var normalizedQualifiedName pkg_rhino.QualifiedName = CheckMissingRequires_SwapRootOfQualifiedName(qualifiedName, aliasedName)
		_ = normalizedQualifiedName
		var originalImport ModuleMetadataMap_ModuleMetadata = nil
		_ = originalImport
		{
			var subName pkg_rhino.QualifiedName = aliasedName
			_ = subName
			for ; subName != nil && !subName.IsJavaNil_(); subName = subName.GetOwner() {
				var imported ModuleMetadataMap_ModuleMetadata = j2gort.AsIface[ModuleMetadataMap_ModuleMetadata](j2gort.MapGetStr(this_.moduleByNamespace, subName.Join_0()))
				_ = imported
				if imported != nil && !imported.IsJavaNil_() {
					originalImport = imported
					break
				}
			}
		}
		this_.priv_CheckMissingRequires_checkNormalizedShortNameImport_1097(t, n, currentFile, originalImport, qualifiedName, normalizedQualifiedName, referenceStrength)
		return
	}
	for initialValue.IsGetProp() {
		initialValue = initialValue.GetFirstChild()
	}
	if !(NodeUtil_IsGoogModuleGetCall(initialValue)) {
		return
	}
	var importedNamespace string = initialValue.GetSecondChild().GetString()
	_ = importedNamespace
	var originalImport_2 ModuleMetadataMap_ModuleMetadata = j2gort.AsIface[ModuleMetadataMap_ModuleMetadata](j2gort.MapGetStr(this_.moduleByNamespace, importedNamespace))
	_ = originalImport_2
	if (originalImport_2 != nil && !originalImport_2.IsJavaNil_()) && !(originalImport_2.HasLegacyGoogNamespaces()) {
		return
	}
	var aliasedName_2 pkg_rhino.QualifiedName = pkg_rhino.QualifiedName_Of(importedNamespace)
	_ = aliasedName_2
	{
		var prop pkg_rhino.Node = initialValue.GetParent()
		_ = prop
		for ; prop.IsGetProp(); prop = prop.GetParent() {
			aliasedName_2 = aliasedName_2.Getprop(prop.GetString())
		}
	}
	var normalizedQualifiedName_2 pkg_rhino.QualifiedName = CheckMissingRequires_SwapRootOfQualifiedName(qualifiedName, aliasedName_2)
	_ = normalizedQualifiedName_2
	this_.priv_CheckMissingRequires_checkNormalizedShortNameImport_1097(t, n, currentFile, originalImport_2, qualifiedName, normalizedQualifiedName_2, referenceStrength)
}

// Java method: CheckMissingRequires.checkMissingRequireThroughShortNameInModule(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.jscomp.NodeUtil.GoogRequire require, com.google.javascript.jscomp.modules.ModuleMetadataMap.ModuleMetadata currentFile, com.google.javascript.rhino.QualifiedName qualifiedName, com.google.javascript.jscomp.CheckMissingRequires.Strength referenceStrength)
func (this_ *CheckMissingRequires) priv_CheckMissingRequires_checkMissingRequireThroughShortNameInModule_1096(t *NodeTraversal, n pkg_rhino.Node, require *NodeUtil_GoogRequire, currentFile ModuleMetadataMap_ModuleMetadata, qualifiedName pkg_rhino.QualifiedName, referenceStrength CheckMissingRequires_Strength) {
	_ = t
	_ = n
	_ = require
	_ = currentFile
	_ = qualifiedName
	_ = referenceStrength
	if require == nil {
		return
	}
	var originalRequiredFile ModuleMetadataMap_ModuleMetadata = j2gort.AsIface[ModuleMetadataMap_ModuleMetadata](j2gort.MapGetStr(this_.moduleByNamespace, require.Namespace()))
	_ = originalRequiredFile
	if (originalRequiredFile == nil || originalRequiredFile.IsJavaNil_()) || !(originalRequiredFile.HasLegacyGoogNamespaces()) {
		return
	}
	var normalizeQualifiedName pkg_rhino.QualifiedName = this_.priv_CheckMissingRequires_normalizeQualifiedNamePlusImport_1099(qualifiedName, require)
	_ = normalizeQualifiedName
	this_.priv_CheckMissingRequires_checkNormalizedShortNameImport_1097(t, n, currentFile, originalRequiredFile, qualifiedName, normalizeQualifiedName, referenceStrength)
}

// Java method: CheckMissingRequires.checkNormalizedShortNameImport(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.jscomp.modules.ModuleMetadataMap.ModuleMetadata currentFile, com.google.javascript.jscomp.modules.ModuleMetadataMap.ModuleMetadata originalImport, com.google.javascript.rhino.QualifiedName originalQualifiedName, com.google.javascript.rhino.QualifiedName normalizedQualifiedName, com.google.javascript.jscomp.CheckMissingRequires.Strength referenceStrength)
func (this_ *CheckMissingRequires) priv_CheckMissingRequires_checkNormalizedShortNameImport_1097(t *NodeTraversal, n pkg_rhino.Node, currentFile ModuleMetadataMap_ModuleMetadata, originalImport ModuleMetadataMap_ModuleMetadata, originalQualifiedName pkg_rhino.QualifiedName, normalizedQualifiedName pkg_rhino.QualifiedName, referenceStrength CheckMissingRequires_Strength) {
	_ = t
	_ = n
	_ = currentFile
	_ = originalImport
	_ = originalQualifiedName
	_ = normalizedQualifiedName
	_ = referenceStrength
	{
		var subName pkg_rhino.QualifiedName = normalizedQualifiedName
		_ = subName
		for ; subName != nil && !subName.IsJavaNil_(); subName = subName.GetOwner() {
			var namespace string = subName.Join_0()
			_ = namespace
			if this_.IsAllowedNamespace(currentFile, namespace) {
				return
			}
			var requiredFile ModuleMetadataMap_ModuleMetadata = j2gort.AsIface[ModuleMetadataMap_ModuleMetadata](j2gort.MapGetStr(this_.moduleByNamespace, namespace))
			_ = requiredFile
			if requiredFile == nil || requiredFile.IsJavaNil_() {
				continue
			}
			if (!(CheckMissingRequires_HasAcceptableRequire(currentFile, subName, requiredFile, referenceStrength)) || (originalImport == nil || originalImport.IsJavaNil_())) || !(originalImport == requiredFile || ((originalImport == nil || originalImport.IsJavaNil_()) && (requiredFile == nil || requiredFile.IsJavaNil_()))) {
				var toReport *DiagnosticType = nil
				_ = toReport
				if currentFile.IsModule() {
					if referenceStrength.FieldGet_isStrong() {
						toReport = StaticGet_CheckMissingRequires_INDIRECT_NAMESPACE_REF_REQUIRE()
					} else {
						toReport = StaticGet_CheckMissingRequires_INDIRECT_NAMESPACE_REF_REQUIRE_TYPE()
					}
				} else {
					if referenceStrength.FieldGet_isStrong() {
						toReport = StaticGet_CheckMissingRequires_MISSING_REQUIRE_IN_GOOG_SCOPE()
					} else {
						toReport = StaticGet_CheckMissingRequires_MISSING_REQUIRE_TYPE_IN_GOOG_SCOPE()
					}
				}
				var errorNode pkg_rhino.Node = n
				_ = errorNode
				var diff int32 = (normalizedQualifiedName.GetComponentCount() - subName.GetComponentCount())
				_ = diff
				if n.IsQualifiedName() {
					{
						var i int32 = int32(0)
						_ = i
						for ; (i < diff) && errorNode.IsGetProp(); i = i + 1 {
							errorNode = errorNode.GetFirstChild()
						}
					}
					t.Report_4(NodeUtil_GetRootOfQualifiedName_1_Node(errorNode), errorNode, toReport, []string{namespace})
				} else {
					var errorQName pkg_rhino.QualifiedName = originalQualifiedName
					_ = errorQName
					{
						var i_2 int32 = int32(0)
						_ = i_2
						for ; i_2 < diff; i_2 = i_2 + 1 {
							errorQName = errorQName.GetOwner()
						}
					}
					var correctName string = errorQName.Join_0()
					_ = correctName
					var newErrorNode pkg_rhino.Node = n.CloneNode_0()
					_ = newErrorNode
					newErrorNode.SetString(correctName)
					newErrorNode.SetLength(j2gort.JStrLen(correctName))
					t.Report_3(newErrorNode, toReport, []string{namespace})
				}
			}
			return
		}
	}
}

// Java method: CheckMissingRequires.checkMissingRequireThroughFullyQualifiedName(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.jscomp.modules.ModuleMetadataMap.ModuleMetadata currentFile, com.google.javascript.rhino.QualifiedName qualifiedName, com.google.javascript.jscomp.CheckMissingRequires.Strength referenceStrength)
func (this_ *CheckMissingRequires) priv_CheckMissingRequires_checkMissingRequireThroughFullyQualifiedName_1098(t *NodeTraversal, n pkg_rhino.Node, currentFile ModuleMetadataMap_ModuleMetadata, qualifiedName pkg_rhino.QualifiedName, referenceStrength CheckMissingRequires_Strength) {
	_ = t
	_ = n
	_ = currentFile
	_ = qualifiedName
	_ = referenceStrength
	{
		var subName pkg_rhino.QualifiedName = qualifiedName
		_ = subName
		for ; subName != nil && !subName.IsJavaNil_(); subName = subName.GetOwner() {
			var namespace string = subName.Join_0()
			_ = namespace
			if this_.IsAllowedNamespace(currentFile, namespace) {
				return
			}
			var requiredFile ModuleMetadataMap_ModuleMetadata = j2gort.AsIface[ModuleMetadataMap_ModuleMetadata](j2gort.MapGetStr(this_.moduleByNamespace, namespace))
			_ = requiredFile
			if requiredFile == nil || requiredFile.IsJavaNil_() {
				continue
			}
			var fileName string = requiredFile.RootNode().GetSourceFileName()
			_ = fileName
			var toReport *DiagnosticType = nil
			_ = toReport
			if currentFile.IsModule() {
				if referenceStrength.FieldGet_isStrong() {
					toReport = StaticGet_CheckMissingRequires_MISSING_REQUIRE()
				} else {
					toReport = StaticGet_CheckMissingRequires_MISSING_REQUIRE_TYPE()
				}
			} else if !(CheckMissingRequires_HasAcceptableRequire(currentFile, subName, requiredFile, referenceStrength)) {
				if referenceStrength.FieldGet_isStrong() {
					toReport = StaticGet_CheckMissingRequires_MISSING_REQUIRE_IN_PROVIDES_FILE()
				} else {
					toReport = StaticGet_CheckMissingRequires_MISSING_REQUIRE_TYPE_IN_PROVIDES_FILE()
				}
			} else if !(referenceStrength.FieldGet_isTypeOnly()) && requiredFile.IsNonLegacyGoogModule() {
				toReport = StaticGet_CheckMissingRequires_NON_LEGACY_GOOG_MODULE_REFERENCE()
			} else {
				return
			}
			t.Report_3(n, toReport, []string{namespace, fileName})
			return
		}
	}
}

// Java method: CheckMissingRequires.normalizeQualifiedNamePlusImport(com.google.javascript.rhino.QualifiedName qualifiedName, com.google.javascript.jscomp.NodeUtil.GoogRequire imported)
func (this_ *CheckMissingRequires) priv_CheckMissingRequires_normalizeQualifiedNamePlusImport_1099(qualifiedName pkg_rhino.QualifiedName, imported *NodeUtil_GoogRequire) pkg_rhino.QualifiedName {
	_ = qualifiedName
	_ = imported
	var newQName pkg_rhino.QualifiedName = pkg_rhino.QualifiedName_Of(imported.Namespace())
	_ = newQName
	if imported.Property() != j2gort.NullString {
		newQName = newQName.Getprop(imported.Property())
	}
	return CheckMissingRequires_SwapRootOfQualifiedName(qualifiedName, newQName)
}

// Java static method: CheckMissingRequires.swapRootOfQualifiedName(com.google.javascript.rhino.QualifiedName original, com.google.javascript.rhino.QualifiedName newRoot)
func CheckMissingRequires_SwapRootOfQualifiedName(original pkg_rhino.QualifiedName, newRoot pkg_rhino.QualifiedName) pkg_rhino.QualifiedName {
	clinit_CheckMissingRequires()
	_ = original
	_ = newRoot
	var newQName pkg_rhino.QualifiedName = newRoot
	_ = newQName
	var components j2gort.Ext_java_util_Iterator = j2gort.AsIface[j2gort.Ext_java_util_Iterator](original.Components().Iterator())
	_ = components
	components.Next()
	for components.HasNext() {
		var component string = j2gort.CastString(components.Next())
		_ = component
		newQName = newQName.Getprop(component)
	}
	return newQName
}

// Java method: CheckMissingRequires.isAllowedNamespace(com.google.javascript.jscomp.modules.ModuleMetadataMap.ModuleMetadata currentFile, java.lang.String namespace)
func (this_ *CheckMissingRequires) IsAllowedNamespace(currentFile ModuleMetadataMap_ModuleMetadata, namespace string) bool {
	_ = currentFile
	_ = namespace
	if j2gort.CheckNotNullString(namespace) == "goog.module" {
		return true
	}
	if j2gort.SetContainsStr(currentFile.GoogNamespaces(), namespace) {
		return true
	}
	return false
}

// Java static method: CheckMissingRequires.hasAcceptableRequire(com.google.javascript.jscomp.modules.ModuleMetadataMap.ModuleMetadata rdep, com.google.javascript.rhino.QualifiedName namespace, com.google.javascript.jscomp.modules.ModuleMetadataMap.ModuleMetadata dep, com.google.javascript.jscomp.CheckMissingRequires.Strength referenceStrength)
func CheckMissingRequires_HasAcceptableRequire(rdep ModuleMetadataMap_ModuleMetadata, namespace pkg_rhino.QualifiedName, dep ModuleMetadataMap_ModuleMetadata, referenceStrength CheckMissingRequires_Strength) bool {
	clinit_CheckMissingRequires()
	_ = rdep
	_ = namespace
	_ = dep
	_ = referenceStrength
	var acceptableRequires j2gort.Ext_java_util_Set = j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableSet](rdep.StronglyRequiredGoogNamespaces().ElementSet())
	_ = acceptableRequires
	if !(referenceStrength.FieldGet_isStrong()) {
		acceptableRequires = j2gort.AsIface[j2gort.Ext_java_util_Set](j2gort.Ext_com_google_common_collect_Sets_Union(acceptableRequires, j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableSet](rdep.WeaklyRequiredGoogNamespaces().ElementSet())))
	}
	acceptableRequires = j2gort.AsIface[j2gort.Ext_java_util_Set](j2gort.Ext_com_google_common_collect_Sets_Intersection(acceptableRequires, j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableSet](dep.GoogNamespaces().ElementSet())))
	{
		var parent pkg_rhino.QualifiedName = namespace
		_ = parent
		for ; parent != nil && !parent.IsJavaNil_(); parent = parent.GetOwner() {
			if j2gort.SetContainsStr(acceptableRequires, parent.Join_0()) {
				return true
			}
		}
	}
	return false
}

// Java enum: com.google.javascript.jscomp.CheckMissingRequires.Strength
type CheckMissingRequires_Strength uint8

const (
	enum_CheckMissingRequires_Strength_CODE               CheckMissingRequires_Strength = 1
	enum_CheckMissingRequires_Strength_IMPLEMENTS_EXTENDS CheckMissingRequires_Strength = 2
	enum_CheckMissingRequires_Strength_WEAK_TYPE          CheckMissingRequires_Strength = 3
)

var enumBaseTable_CheckMissingRequires_Strength = [...]j2gort.EnumBase{
	{Name_: "CODE", Ordinal_: 0},
	{Name_: "IMPLEMENTS_EXTENDS", Ordinal_: 1},
	{Name_: "WEAK_TYPE", Ordinal_: 2},
}
var enumValues_CheckMissingRequires_Strength = []CheckMissingRequires_Strength{1, 2, 3}

var enumField_CheckMissingRequires_Strength_isStrong = [...]bool{
	true,
	true,
	false,
}

func (this_ CheckMissingRequires_Strength) FieldGet_isStrong() bool {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumField_CheckMissingRequires_Strength_isStrong[this_-1]
}

var enumField_CheckMissingRequires_Strength_isTypeOnly = [...]bool{
	false,
	true,
	true,
}

func (this_ CheckMissingRequires_Strength) FieldGet_isTypeOnly() bool {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumField_CheckMissingRequires_Strength_isTypeOnly[this_-1]
}

func StaticGet_CheckMissingRequires_Strength_CODE() CheckMissingRequires_Strength {
	return enum_CheckMissingRequires_Strength_CODE
}
func StaticGet_CheckMissingRequires_Strength_IMPLEMENTS_EXTENDS() CheckMissingRequires_Strength {
	return enum_CheckMissingRequires_Strength_IMPLEMENTS_EXTENDS
}
func StaticGet_CheckMissingRequires_Strength_WEAK_TYPE() CheckMissingRequires_Strength {
	return enum_CheckMissingRequires_Strength_WEAK_TYPE
}

func CheckMissingRequires_Strength_Values() []CheckMissingRequires_Strength {
	return j2gort.CloneSlice(enumValues_CheckMissingRequires_Strength)
}

func CheckMissingRequires_Strength_ValueOf(name string) CheckMissingRequires_Strength {
	if name == j2gort.NullString {
		j2gort.ThrowNPE()
	}
	for i_ := range enumBaseTable_CheckMissingRequires_Strength {
		if enumBaseTable_CheckMissingRequires_Strength[i_].Name_ == name {
			return CheckMissingRequires_Strength(i_ + 1)
		}
	}
	panic(javalang.NewIllegalArgumentException("No enum constant com.google.javascript.jscomp.CheckMissingRequires.Strength." + name))
}

func (this_ CheckMissingRequires_Strength) IsJavaNil_() bool                              { return this_ == 0 }
func (CheckMissingRequires_Strength) HasIdentityEquals_()                                 {}
func (CheckMissingRequires_Strength) Implements_j2gort_Ext_java_lang_Enum()               {}
func (CheckMissingRequires_Strength) Implements_j2gort_Ext_java_lang_constant_Constable() {}
func (this_ CheckMissingRequires_Strength) GetExt_java_lang_EnumBase() *j2gort.Ext_java_lang_EnumBase {
	return &j2gort.Ext_java_lang_EnumBase{Core_: this_.GetEnumBase()}
}
func (CheckMissingRequires_Strength) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.CheckMissingRequires$Strength", "Strength"
}
func (this_ CheckMissingRequires_Strength) GetEnumBase() *j2gort.EnumBase {
	if this_ == 0 {
		return nil
	}
	return &enumBaseTable_CheckMissingRequires_Strength[this_-1]
}
func (this_ CheckMissingRequires_Strength) Ordinal() int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return int32(this_ - 1)
}
func (this_ CheckMissingRequires_Strength) EnumOrdinal_() int32 { return int32(this_) - 1 }
func (this_ CheckMissingRequires_Strength) Name() string {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_CheckMissingRequires_Strength[this_-1].Name_
}
func (this_ CheckMissingRequires_Strength) String() string {
	if this_ == 0 {
		return "null"
	}
	return enumBaseTable_CheckMissingRequires_Strength[this_-1].Name_
}
func (this_ CheckMissingRequires_Strength) Equals(other_ any) bool {
	o_, ok_ := other_.(CheckMissingRequires_Strength)
	return this_ != 0 && ok_ && this_ == o_
}
func (this_ CheckMissingRequires_Strength) HashCode() int32 {
	if this_ == 0 {
		return 0
	}
	return enumBaseTable_CheckMissingRequires_Strength[this_-1].HashCode()
}
func (this_ CheckMissingRequires_Strength) CompareTo(other_ any) int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_CheckMissingRequires_Strength[this_-1].CompareTo(other_)
}
func (this_ CheckMissingRequires_Strength) WriteGson(out_ *gson.JsonWriter) {
	if this_ == 0 {
		out_.NullValue()
		return
	}
	out_.ValueString(enumBaseTable_CheckMissingRequires_Strength[this_-1].Name_)
}
func init() {
	vals_ := []any{CheckMissingRequires_Strength(1), CheckMissingRequires_Strength(2), CheckMissingRequires_Strength(3)}
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.CheckMissingRequires.Strength", vals_)
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.CheckMissingRequires$Strength", vals_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.CheckMissingRequires.Strength", "com.google.javascript.jscomp.CheckMissingRequires$Strength", func(v_ any) bool {
		return j2gort.InstanceOf[CheckMissingRequires_Strength](v_)
	})
}
