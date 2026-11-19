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
// Port of com/google/javascript/jscomp/lint/CheckNullabilityModifiers.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
)

// Java class: com.google.javascript.jscomp.lint.CheckNullabilityModifiers
type CheckNullabilityModifiers struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	compiler              AbstractCompiler                  `gson:"compiler"`
	redundantCandidates   j2gort.OrderedSet[pkg_rhino.Node] `gson:"redundantCandidates"`
	missingCandidates     j2gort.OrderedSet[pkg_rhino.Node] `gson:"missingCandidates"`
	nullMissingCandidates j2gort.OrderedSet[pkg_rhino.Node] `gson:"nullMissingCandidates"`
	templateTypeNames     j2gort.OrderedSet[string]         `gson:"templateTypeNames"`
}

func (this_ *CheckNullabilityModifiers) BindSelf_(s any) {
	_ = s
}
func (this_ *CheckNullabilityModifiers) IsJavaNil_() bool { return this_ == nil }
func (*CheckNullabilityModifiers) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.lint.CheckNullabilityModifiers", "CheckNullabilityModifiers"
}
func (this_ *CheckNullabilityModifiers) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.NodeTraversal_AbstractPostOrderCallbackBase.IdentityHash_()
}
func (*CheckNullabilityModifiers) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.lint.CheckNullabilityModifiers", "com.google.javascript.jscomp.lint.CheckNullabilityModifiers", func(v_ any) bool {
		return j2gort.InstanceOf[*CheckNullabilityModifiers](v_)
	})
}
func (*CheckNullabilityModifiers) Implements_jscomp_NodeTraversal_AbstractPostOrderCallback() {}
func (*CheckNullabilityModifiers) Implements_jscomp_CompilerPass()                            {}
func (*CheckNullabilityModifiers) Implements_jscomp_NodeTraversal_Callback()                  {}

var static_CheckNullabilityModifiers_MISSING_NULLABILITY_MODIFIER_JSDOC *DiagnosticType

func StaticGet_CheckNullabilityModifiers_MISSING_NULLABILITY_MODIFIER_JSDOC() *DiagnosticType {
	clinit_CheckNullabilityModifiers()
	return static_CheckNullabilityModifiers_MISSING_NULLABILITY_MODIFIER_JSDOC
}
func StaticSet_CheckNullabilityModifiers_MISSING_NULLABILITY_MODIFIER_JSDOC(v *DiagnosticType) {
	clinit_CheckNullabilityModifiers()
	static_CheckNullabilityModifiers_MISSING_NULLABILITY_MODIFIER_JSDOC = v
}

var static_CheckNullabilityModifiers_NULL_MISSING_NULLABILITY_MODIFIER_JSDOC *DiagnosticType

func StaticGet_CheckNullabilityModifiers_NULL_MISSING_NULLABILITY_MODIFIER_JSDOC() *DiagnosticType {
	clinit_CheckNullabilityModifiers()
	return static_CheckNullabilityModifiers_NULL_MISSING_NULLABILITY_MODIFIER_JSDOC
}
func StaticSet_CheckNullabilityModifiers_NULL_MISSING_NULLABILITY_MODIFIER_JSDOC(v *DiagnosticType) {
	clinit_CheckNullabilityModifiers()
	static_CheckNullabilityModifiers_NULL_MISSING_NULLABILITY_MODIFIER_JSDOC = v
}

var static_CheckNullabilityModifiers_REDUNDANT_NULLABILITY_MODIFIER_JSDOC *DiagnosticType

func StaticGet_CheckNullabilityModifiers_REDUNDANT_NULLABILITY_MODIFIER_JSDOC() *DiagnosticType {
	clinit_CheckNullabilityModifiers()
	return static_CheckNullabilityModifiers_REDUNDANT_NULLABILITY_MODIFIER_JSDOC
}
func StaticSet_CheckNullabilityModifiers_REDUNDANT_NULLABILITY_MODIFIER_JSDOC(v *DiagnosticType) {
	clinit_CheckNullabilityModifiers()
	static_CheckNullabilityModifiers_REDUNDANT_NULLABILITY_MODIFIER_JSDOC = v
}

var static_CheckNullabilityModifiers_PRIMITIVE_TYPE_NAMES j2gort.Ext_com_google_common_collect_ImmutableSet

func StaticGet_CheckNullabilityModifiers_PRIMITIVE_TYPE_NAMES() j2gort.Ext_com_google_common_collect_ImmutableSet {
	clinit_CheckNullabilityModifiers()
	return static_CheckNullabilityModifiers_PRIMITIVE_TYPE_NAMES
}
func StaticSet_CheckNullabilityModifiers_PRIMITIVE_TYPE_NAMES(v j2gort.Ext_com_google_common_collect_ImmutableSet) {
	clinit_CheckNullabilityModifiers()
	static_CheckNullabilityModifiers_PRIMITIVE_TYPE_NAMES = v
}

var clinitDone_CheckNullabilityModifiers uint32
var clinitState_CheckNullabilityModifiers j2gort.ClinitState

func clinit_CheckNullabilityModifiers() {
	if atomic.LoadUint32(&clinitDone_CheckNullabilityModifiers) != 0 {
		return
	}
	clinitSlow_CheckNullabilityModifiers()
}

func clinitSlow_CheckNullabilityModifiers() {
	if !j2gort.EnterClinit(&clinitState_CheckNullabilityModifiers, &clinitDone_CheckNullabilityModifiers) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_CheckNullabilityModifiers, &clinitDone_CheckNullabilityModifiers)
	static_CheckNullabilityModifiers_MISSING_NULLABILITY_MODIFIER_JSDOC = DiagnosticType_Disabled("JSC_MISSING_NULLABILITY_MODIFIER_JSDOC", "{0} is a reference type with no nullability modifier, which is disallowed by the style guide.\nPlease add a '!' to make it explicitly non-nullable, or a '?' to make it explicitly nullable.")
	static_CheckNullabilityModifiers_NULL_MISSING_NULLABILITY_MODIFIER_JSDOC = DiagnosticType_Disabled("JSC_NULL_MISSING_NULLABILITY_MODIFIER_JSDOC", "{0} is a reference type with no nullability modifier that is explicitly set to null.\nAdd a '?' to make it explicitly nullable.")
	static_CheckNullabilityModifiers_REDUNDANT_NULLABILITY_MODIFIER_JSDOC = DiagnosticType_Disabled("JSC_REDUNDANT_NULLABILITY_MODIFIER_JSDOC", "{0} is a non-reference type which is already non-nullable.\nPlease remove the redundant '!', which is disallowed by the style guide.")
	static_CheckNullabilityModifiers_PRIMITIVE_TYPE_NAMES = j2gort.Ext_com_google_common_collect_ImmutableSet_Of_7("boolean", "number", "bigint", "string", "symbol", "undefined", []any{"void", "null"})
}

func (this_ *CheckNullabilityModifiers) initFields_CheckNullabilityModifiers() {
	_ = this_
	this_.redundantCandidates = j2gort.NewOrderedSet[pkg_rhino.Node]()
	this_.missingCandidates = j2gort.NewOrderedSet[pkg_rhino.Node]()
	this_.nullMissingCandidates = j2gort.NewOrderedSet[pkg_rhino.Node]()
	this_.templateTypeNames = j2gort.NewOrderedSet[string]()
}

// Java constructor: CheckNullabilityModifiers(com.google.javascript.jscomp.AbstractCompiler compiler)
func NewCheckNullabilityModifiers(compiler AbstractCompiler) *CheckNullabilityModifiers {
	clinit_CheckNullabilityModifiers()
	this_ := &CheckNullabilityModifiers{}
	InitCheckNullabilityModifiers(this_, compiler)
	return this_
}

func InitCheckNullabilityModifiers(this_ *CheckNullabilityModifiers, compiler AbstractCompiler) {
	_ = compiler
	InitNodeTraversal_AbstractPostOrderCallback(&this_.NodeTraversal_AbstractPostOrderCallbackBase)
	this_.initFields_CheckNullabilityModifiers()
	this_.compiler = compiler
}

// Java method: CheckNullabilityModifiers.process(com.google.javascript.rhino.Node externs, com.google.javascript.rhino.Node root)
func (this_ *CheckNullabilityModifiers) Process_2(externs pkg_rhino.Node, root pkg_rhino.Node) {
	_ = externs
	_ = root
	NodeTraversal_TraverseRoots(this_.compiler, this_, externs, root)
}

// Java method: CheckNullabilityModifiers.visit(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *CheckNullabilityModifiers) Visit_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = n
	_ = parent
	var info *pkg_rhino.JSDocInfo = n.GetJSDocInfo()
	_ = info
	if info != nil {
		j2gort.OrderedSetAddAllStr(&this_.templateTypeNames, j2gort.AsIface[j2gort.Ext_java_util_Collection](info.GetTemplateTypeNames()))
		j2gort.OrderedSetAddAllStr(&this_.templateTypeNames, j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableSet](info.GetTypeTransformations().KeySet_0()))
		if info.HasType() {
			this_.priv_CheckNullabilityModifiers_handleHasType_6557(info.GetType(), n)
		}
		{
			var i int32 = int32(0)
			_ = i
			var paramCount int32 = info.GetParameterCount()
			_ = paramCount
			for ; i < paramCount; i = i + 1 {
				var paramType pkg_rhino.JSTypeExpression = info.GetParameterTypeAt(i)
				_ = paramType
				if !(paramType).IsJavaNil_() {
					this_.priv_CheckNullabilityModifiers_visitTypeExpression_6560(paramType, false)
				}
			}
		}
		if info.HasReturnType() {
			this_.priv_CheckNullabilityModifiers_visitTypeExpression_6560(info.GetReturnType(), false)
		}
		if info.HasEnumParameterType() {
			this_.priv_CheckNullabilityModifiers_visitTypeExpression_6560(info.GetEnumParameterType(), false)
		}
		if info.HasTypedefType() {
			this_.priv_CheckNullabilityModifiers_visitTypeExpression_6560(info.GetTypedefType(), false)
		}
		if info.HasThisType() {
			this_.priv_CheckNullabilityModifiers_visitTypeExpression_6560(info.GetThisType(), true)
		}
	}
	if n.IsScript() {
		this_.priv_CheckNullabilityModifiers_report_6558(t)
		this_.redundantCandidates.Clear()
		this_.missingCandidates.Clear()
		this_.nullMissingCandidates.Clear()
		this_.templateTypeNames.Clear()
		return
	}
}

// Java method: CheckNullabilityModifiers.handleHasType(com.google.javascript.rhino.JSTypeExpression expr, com.google.javascript.rhino.Node n)
func (this_ *CheckNullabilityModifiers) priv_CheckNullabilityModifiers_handleHasType_6557(expr pkg_rhino.JSTypeExpression, n pkg_rhino.Node) {
	_ = expr
	_ = n
	if NodeUtil_IsNameDeclOrSimpleAssignLhs(n.GetFirstChild(), n) {
		var rValue pkg_rhino.Node = NodeUtil_GetRValueOfLValue(n.GetFirstChild())
		_ = rValue
		if (rValue != 0) && rValue.IsNull() {
			this_.priv_CheckNullabilityModifiers_visitTypeExpression_6561(expr, false, rValue)
			return
		}
	}
	this_.priv_CheckNullabilityModifiers_visitTypeExpression_6560(expr, false)
}

// Java method: CheckNullabilityModifiers.report(com.google.javascript.jscomp.NodeTraversal t)
func (this_ *CheckNullabilityModifiers) priv_CheckNullabilityModifiers_report_6558(t *NodeTraversal) {
	_ = t
	{
		cur_1 := this_.missingCandidates.Iter()
		for cur_1.HasNext() {
			n := cur_1.
				Next()
			_ = n
			if this_.priv_CheckNullabilityModifiers_shouldReport_6559(n) {
				t.Report_3(n, StaticGet_CheckNullabilityModifiers_MISSING_NULLABILITY_MODIFIER_JSDOC(), []string{CheckNullabilityModifiers_GetReportedTypeName(n)})
			}
		}
	}
	{
		cur_2 := this_.nullMissingCandidates.Iter()
		for cur_2.HasNext() {
			n_2 := cur_2.
				Next()
			_ = n_2
			if this_.priv_CheckNullabilityModifiers_shouldReport_6559(n_2) {
				t.Report_3(n_2, StaticGet_CheckNullabilityModifiers_NULL_MISSING_NULLABILITY_MODIFIER_JSDOC(), []string{CheckNullabilityModifiers_GetReportedTypeName(n_2)})
			}
		}
	}
	{
		cur_3 := this_.redundantCandidates.Iter()
		for cur_3.HasNext() {
			n_3 := cur_3.
				Next()
			_ = n_3
			if this_.priv_CheckNullabilityModifiers_shouldReport_6559(n_3) {
				t.Report_3(n_3.GetParent(), StaticGet_CheckNullabilityModifiers_REDUNDANT_NULLABILITY_MODIFIER_JSDOC(), []string{CheckNullabilityModifiers_GetReportedTypeName(n_3)})
			}
		}
	}
}

// Java method: CheckNullabilityModifiers.shouldReport(com.google.javascript.rhino.Node n)
func (this_ *CheckNullabilityModifiers) priv_CheckNullabilityModifiers_shouldReport_6559(n pkg_rhino.Node) bool {
	_ = n
	return (!(n.IsStringLit()) || !(this_.templateTypeNames.Contains(n.GetString())))
}

// Java method: CheckNullabilityModifiers.visitTypeExpression(com.google.javascript.rhino.JSTypeExpression expr, boolean hasArtificialTopLevelBang)
func (this_ *CheckNullabilityModifiers) priv_CheckNullabilityModifiers_visitTypeExpression_6560(expr pkg_rhino.JSTypeExpression, hasArtificialTopLevelBang bool) {
	_ = expr
	_ = hasArtificialTopLevelBang
	this_.priv_CheckNullabilityModifiers_visitTypeExpression_6561(expr, hasArtificialTopLevelBang, pkg_rhino.Node(0))
}

// Java method: CheckNullabilityModifiers.visitTypeExpression(com.google.javascript.rhino.JSTypeExpression expr, boolean hasArtificialTopLevelBang, com.google.javascript.rhino.@org.jspecify.annotations.Nullable Node rValue)
func (this_ *CheckNullabilityModifiers) priv_CheckNullabilityModifiers_visitTypeExpression_6561(expr pkg_rhino.JSTypeExpression, hasArtificialTopLevelBang bool, rValue pkg_rhino.Node) {
	_ = expr
	_ = hasArtificialTopLevelBang
	_ = rValue
	var root pkg_rhino.Node = expr.GetRoot()
	_ = root
	NodeUtil_VisitPreOrder_2(root, NewLambda_NodeUtil_Visitor(func(node pkg_rhino.Node) {
		_ = node
		var parent pkg_rhino.Node = node.GetParent()
		_ = parent
		var isPrimitiveOrLiteral bool = ((CheckNullabilityModifiers_IsPrimitiveType(node) || CheckNullabilityModifiers_IsFunctionLiteral(node)) || CheckNullabilityModifiers_IsRecordLiteral(node))
		_ = isPrimitiveOrLiteral
		var isReference bool = CheckNullabilityModifiers_IsReferenceType(node)
		_ = isReference
		var hasBang bool = ((parent != 0) && (parent.RawToken_() == uint8(159)))
		_ = hasBang
		var hasQmark bool = ((parent != 0) && (parent.RawToken_() == uint8(158)))
		_ = hasQmark
		var hasNonArtificialBang bool = (hasBang && !(hasArtificialTopLevelBang && (parent == root)))
		_ = hasNonArtificialBang
		var isNewOrThis bool = ((parent != 0) && (parent.IsNew() || parent.IsThis()))
		_ = isNewOrThis
		var isTypeOfType bool = ((parent != 0) && parent.IsTypeOf())
		_ = isTypeOfType
		var isAppliedToRValue bool = ((node == root) || ((parent != 0) && (parent.RawToken_() == uint8(155))))
		_ = isAppliedToRValue
		if (((isReference && !(hasBang)) && !(hasQmark)) && !(isNewOrThis)) && !(isTypeOfType) {
			if (isAppliedToRValue && (rValue != 0)) && rValue.IsNull() {
				this_.nullMissingCandidates.Add(node)
			} else {
				this_.missingCandidates.Add(node)
			}
		} else if isPrimitiveOrLiteral && hasNonArtificialBang {
			this_.redundantCandidates.Add(node)
		}
	}))
}

// Java static method: CheckNullabilityModifiers.isPrimitiveType(com.google.javascript.rhino.Node node)
func CheckNullabilityModifiers_IsPrimitiveType(node pkg_rhino.Node) bool {
	clinit_CheckNullabilityModifiers()
	_ = node
	return (node.IsStringLit() && j2gort.SetContainsStr(StaticGet_CheckNullabilityModifiers_PRIMITIVE_TYPE_NAMES(), node.GetString()))
}

// Java static method: CheckNullabilityModifiers.isReferenceType(com.google.javascript.rhino.Node node)
func CheckNullabilityModifiers_IsReferenceType(node pkg_rhino.Node) bool {
	clinit_CheckNullabilityModifiers()
	_ = node
	return (node.IsStringLit() && !(j2gort.SetContainsStr(StaticGet_CheckNullabilityModifiers_PRIMITIVE_TYPE_NAMES(), node.GetString())))
}

// Java static method: CheckNullabilityModifiers.isFunctionLiteral(com.google.javascript.rhino.Node node)
func CheckNullabilityModifiers_IsFunctionLiteral(node pkg_rhino.Node) bool {
	clinit_CheckNullabilityModifiers()
	_ = node
	return node.IsFunction()
}

// Java static method: CheckNullabilityModifiers.isRecordLiteral(com.google.javascript.rhino.Node node)
func CheckNullabilityModifiers_IsRecordLiteral(node pkg_rhino.Node) bool {
	clinit_CheckNullabilityModifiers()
	_ = node
	return (node.RawToken_() == uint8(162))
}

// Java static method: CheckNullabilityModifiers.getReportedTypeName(com.google.javascript.rhino.Node node)
func CheckNullabilityModifiers_GetReportedTypeName(node pkg_rhino.Node) string {
	clinit_CheckNullabilityModifiers()
	_ = node
	if CheckNullabilityModifiers_IsFunctionLiteral(node) {
		return "Function"
	}
	if CheckNullabilityModifiers_IsRecordLiteral(node) {
		return "Record literal"
	}
	if !(CheckNullabilityModifiers_IsPrimitiveType(node) || CheckNullabilityModifiers_IsReferenceType(node)) {
		j2gort.CheckState(false)
	}
	return node.GetString()
}
func (*CheckNullabilityModifiers) AlwaysShouldTraverse_()      {}
func (*CheckNullabilityModifiers) CanReleaseSyntacticScopes_() {}
