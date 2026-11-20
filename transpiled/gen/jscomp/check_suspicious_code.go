/*
 * Copyright 2012 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/CheckSuspiciousCode.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	pkg_base "github.com/paullepoulpe/jsgo/transpiled/gen/base"
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
)

// Java class: com.google.javascript.jscomp.CheckSuspiciousCode
type CheckSuspiciousCode struct {
	NodeTraversal_AbstractPostOrderCallbackBase
}

func (this_ *CheckSuspiciousCode) BindSelf_(s any) {
	_ = s
}
func (this_ *CheckSuspiciousCode) IsJavaNil_() bool { return this_ == nil }
func (*CheckSuspiciousCode) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.CheckSuspiciousCode", "CheckSuspiciousCode"
}
func (this_ *CheckSuspiciousCode) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.NodeTraversal_AbstractPostOrderCallbackBase.IdentityHash_()
}
func (*CheckSuspiciousCode) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.CheckSuspiciousCode", "com.google.javascript.jscomp.CheckSuspiciousCode", func(v_ any) bool {
		return j2gort.InstanceOf[*CheckSuspiciousCode](v_)
	})
}
func (*CheckSuspiciousCode) Implements_jscomp_NodeTraversal_AbstractPostOrderCallback() {}
func (*CheckSuspiciousCode) Implements_jscomp_NodeTraversal_Callback()                  {}

var static_CheckSuspiciousCode_SUSPICIOUS_SEMICOLON *DiagnosticType

func StaticGet_CheckSuspiciousCode_SUSPICIOUS_SEMICOLON() *DiagnosticType {
	clinit_CheckSuspiciousCode()
	return static_CheckSuspiciousCode_SUSPICIOUS_SEMICOLON
}
func StaticSet_CheckSuspiciousCode_SUSPICIOUS_SEMICOLON(v *DiagnosticType) {
	clinit_CheckSuspiciousCode()
	static_CheckSuspiciousCode_SUSPICIOUS_SEMICOLON = v
}

var static_CheckSuspiciousCode_SUSPICIOUS_COMPARISON_WITH_NAN *DiagnosticType

func StaticGet_CheckSuspiciousCode_SUSPICIOUS_COMPARISON_WITH_NAN() *DiagnosticType {
	clinit_CheckSuspiciousCode()
	return static_CheckSuspiciousCode_SUSPICIOUS_COMPARISON_WITH_NAN
}
func StaticSet_CheckSuspiciousCode_SUSPICIOUS_COMPARISON_WITH_NAN(v *DiagnosticType) {
	clinit_CheckSuspiciousCode()
	static_CheckSuspiciousCode_SUSPICIOUS_COMPARISON_WITH_NAN = v
}

var static_CheckSuspiciousCode_SUSPICIOUS_IN_OPERATOR *DiagnosticType

func StaticGet_CheckSuspiciousCode_SUSPICIOUS_IN_OPERATOR() *DiagnosticType {
	clinit_CheckSuspiciousCode()
	return static_CheckSuspiciousCode_SUSPICIOUS_IN_OPERATOR
}
func StaticSet_CheckSuspiciousCode_SUSPICIOUS_IN_OPERATOR(v *DiagnosticType) {
	clinit_CheckSuspiciousCode()
	static_CheckSuspiciousCode_SUSPICIOUS_IN_OPERATOR = v
}

var static_CheckSuspiciousCode_SUSPICIOUS_INSTANCEOF_LEFT_OPERAND *DiagnosticType

func StaticGet_CheckSuspiciousCode_SUSPICIOUS_INSTANCEOF_LEFT_OPERAND() *DiagnosticType {
	clinit_CheckSuspiciousCode()
	return static_CheckSuspiciousCode_SUSPICIOUS_INSTANCEOF_LEFT_OPERAND
}
func StaticSet_CheckSuspiciousCode_SUSPICIOUS_INSTANCEOF_LEFT_OPERAND(v *DiagnosticType) {
	clinit_CheckSuspiciousCode()
	static_CheckSuspiciousCode_SUSPICIOUS_INSTANCEOF_LEFT_OPERAND = v
}

var static_CheckSuspiciousCode_SUSPICIOUS_LEFT_OPERAND_OF_LOGICAL_OPERATOR *DiagnosticType

func StaticGet_CheckSuspiciousCode_SUSPICIOUS_LEFT_OPERAND_OF_LOGICAL_OPERATOR() *DiagnosticType {
	clinit_CheckSuspiciousCode()
	return static_CheckSuspiciousCode_SUSPICIOUS_LEFT_OPERAND_OF_LOGICAL_OPERATOR
}
func StaticSet_CheckSuspiciousCode_SUSPICIOUS_LEFT_OPERAND_OF_LOGICAL_OPERATOR(v *DiagnosticType) {
	clinit_CheckSuspiciousCode()
	static_CheckSuspiciousCode_SUSPICIOUS_LEFT_OPERAND_OF_LOGICAL_OPERATOR = v
}

var static_CheckSuspiciousCode_SUSPICIOUS_NEGATED_LEFT_OPERAND_OF_IN_OPERATOR *DiagnosticType

func StaticGet_CheckSuspiciousCode_SUSPICIOUS_NEGATED_LEFT_OPERAND_OF_IN_OPERATOR() *DiagnosticType {
	clinit_CheckSuspiciousCode()
	return static_CheckSuspiciousCode_SUSPICIOUS_NEGATED_LEFT_OPERAND_OF_IN_OPERATOR
}
func StaticSet_CheckSuspiciousCode_SUSPICIOUS_NEGATED_LEFT_OPERAND_OF_IN_OPERATOR(v *DiagnosticType) {
	clinit_CheckSuspiciousCode()
	static_CheckSuspiciousCode_SUSPICIOUS_NEGATED_LEFT_OPERAND_OF_IN_OPERATOR = v
}

var static_CheckSuspiciousCode_SUSPICIOUS_BREAKING_OUT_OF_OPTIONAL_CHAIN *DiagnosticType

func StaticGet_CheckSuspiciousCode_SUSPICIOUS_BREAKING_OUT_OF_OPTIONAL_CHAIN() *DiagnosticType {
	clinit_CheckSuspiciousCode()
	return static_CheckSuspiciousCode_SUSPICIOUS_BREAKING_OUT_OF_OPTIONAL_CHAIN
}
func StaticSet_CheckSuspiciousCode_SUSPICIOUS_BREAKING_OUT_OF_OPTIONAL_CHAIN(v *DiagnosticType) {
	clinit_CheckSuspiciousCode()
	static_CheckSuspiciousCode_SUSPICIOUS_BREAKING_OUT_OF_OPTIONAL_CHAIN = v
}

var clinitDone_CheckSuspiciousCode uint32
var clinitState_CheckSuspiciousCode j2gort.ClinitState

func clinit_CheckSuspiciousCode() {
	if atomic.LoadUint32(&clinitDone_CheckSuspiciousCode) != 0 {
		return
	}
	clinitSlow_CheckSuspiciousCode()
}

func clinitSlow_CheckSuspiciousCode() {
	if !j2gort.EnterClinit(&clinitState_CheckSuspiciousCode, &clinitDone_CheckSuspiciousCode) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_CheckSuspiciousCode, &clinitDone_CheckSuspiciousCode)
	static_CheckSuspiciousCode_SUSPICIOUS_SEMICOLON = DiagnosticType_Warning("JSC_SUSPICIOUS_SEMICOLON", "If this if/for/while really shouldn''t have a body, use '{}'")
	static_CheckSuspiciousCode_SUSPICIOUS_COMPARISON_WITH_NAN = DiagnosticType_Warning("JSC_SUSPICIOUS_NAN", "Comparison against NaN is always false. Did you mean isNaN()?")
	static_CheckSuspiciousCode_SUSPICIOUS_IN_OPERATOR = DiagnosticType_Warning("JSC_SUSPICIOUS_IN", "Use of the \"in\" keyword on non-object types throws an exception.")
	static_CheckSuspiciousCode_SUSPICIOUS_INSTANCEOF_LEFT_OPERAND = DiagnosticType_Warning("JSC_SUSPICIOUS_INSTANCEOF_LEFT", "\"instanceof\" with left non-object operand is always false.")
	static_CheckSuspiciousCode_SUSPICIOUS_LEFT_OPERAND_OF_LOGICAL_OPERATOR = DiagnosticType_Warning("JSC_SUSPICIOUS_LEFT_OPERAND_OF_LOGICAL_OPERATOR", "Left operand of {0} operator is always {1}.")
	static_CheckSuspiciousCode_SUSPICIOUS_NEGATED_LEFT_OPERAND_OF_IN_OPERATOR = DiagnosticType_Warning("JSC_SUSPICIOUS_NEGATED_LEFT_OPERAND_OF_IN_OPERATOR", "Suspicious negated left operand of 'in' operator.")
	static_CheckSuspiciousCode_SUSPICIOUS_BREAKING_OUT_OF_OPTIONAL_CHAIN = DiagnosticType_Warning("SUSPICIOUS_BREAKING_OUT_OF_OPTIONAL_CHAIN", "Suspicious breaking out of optional chain. May result in TypeError if optional chain is undefined.")
}

func (this_ *CheckSuspiciousCode) initFields_CheckSuspiciousCode() {
	_ = this_
}

// Java constructor: CheckSuspiciousCode()
func NewCheckSuspiciousCode() *CheckSuspiciousCode {
	clinit_CheckSuspiciousCode()
	this_ := &CheckSuspiciousCode{}
	InitCheckSuspiciousCode(this_)
	return this_
}

func InitCheckSuspiciousCode(this_ *CheckSuspiciousCode) {
	InitNodeTraversal_AbstractPostOrderCallback(&this_.NodeTraversal_AbstractPostOrderCallbackBase)
	this_.initFields_CheckSuspiciousCode()
}

// Java method: CheckSuspiciousCode.visit(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *CheckSuspiciousCode) Visit_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = n
	_ = parent
	this_.priv_CheckSuspiciousCode_checkMissingSemicolon_1121(t, n)
	this_.priv_CheckSuspiciousCode_checkNaN_1122(t, n)
	this_.priv_CheckSuspiciousCode_checkInvalidIn_1123(t, n)
	this_.priv_CheckSuspiciousCode_checkNonObjectInstanceOf_1124(t, n)
	this_.priv_CheckSuspiciousCode_checkNegatedLeftOperandOfInOperator_1125(t, n)
	this_.priv_CheckSuspiciousCode_checkLeftOperandOfLogicalOperator_1126(t, n)
	CheckSuspiciousCode_CheckSuspiciousBreakingOutOfOptionalChain(t, n)
}

// Java method: CheckSuspiciousCode.checkMissingSemicolon(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n)
func (this_ *CheckSuspiciousCode) priv_CheckSuspiciousCode_checkMissingSemicolon_1121(t *NodeTraversal, n pkg_rhino.Node) {
	_ = t
	_ = n
	switch n.TokenSwitchOrd_() {
	case int32(74) /* Token.IF */ :
		var trueCase pkg_rhino.Node = n.GetSecondChild()
		_ = trueCase
		CheckSuspiciousCode_ReportIfWasEmpty(t, trueCase)
		var elseCase pkg_rhino.Node = trueCase.GetNext()
		_ = elseCase
		if elseCase != 0 {
			CheckSuspiciousCode_ReportIfWasEmpty(t, elseCase)
		}
	case int32(78) /* Token.WHILE */, int32(80) /* Token.FOR */, int32(81) /* Token.FOR_IN */, int32(112) /* Token.FOR_OF */, int32(113) /* Token.FOR_AWAIT_OF */ :
		CheckSuspiciousCode_ReportIfWasEmpty(t, NodeUtil_GetLoopCodeBlock(n))
	default:
	}
}

// Java static method: CheckSuspiciousCode.reportIfWasEmpty(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node block)
func CheckSuspiciousCode_ReportIfWasEmpty(t *NodeTraversal, block pkg_rhino.Node) {
	clinit_CheckSuspiciousCode()
	_ = t
	_ = block
	if !(block.IsBlock()) {
		j2gort.CheckState(false)
	}
	if !(block.HasChildren()) && block.IsAddedBlock() {
		t.GetCompiler().Report_1(JSError_Make_3(block, StaticGet_CheckSuspiciousCode_SUSPICIOUS_SEMICOLON(), []string{}))
	}
}

// Java method: CheckSuspiciousCode.checkNaN(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n)
func (this_ *CheckSuspiciousCode) priv_CheckSuspiciousCode_checkNaN_1122(t *NodeTraversal, n pkg_rhino.Node) {
	_ = t
	_ = n
	switch n.TokenSwitchOrd_() {
	case int32(4) /* Token.EQ */, int32(9) /* Token.GE */, int32(8) /* Token.GT */, int32(7) /* Token.LE */, int32(6) /* Token.LT */, int32(5) /* Token.NE */, int32(40) /* Token.SHEQ */, int32(41) /* Token.SHNE */ :
		CheckSuspiciousCode_ReportIfNaN(t, n.GetFirstChild())
		CheckSuspiciousCode_ReportIfNaN(t, n.GetLastChild())
	default:
	}
}

// Java static method: CheckSuspiciousCode.reportIfNaN(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n)
func CheckSuspiciousCode_ReportIfNaN(t *NodeTraversal, n pkg_rhino.Node) {
	clinit_CheckSuspiciousCode()
	_ = t
	_ = n
	if NodeUtil_IsNaN(n) {
		t.GetCompiler().Report_1(JSError_Make_3(n.GetParent(), StaticGet_CheckSuspiciousCode_SUSPICIOUS_COMPARISON_WITH_NAN(), []string{}))
	}
}

// Java method: CheckSuspiciousCode.checkInvalidIn(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n)
func (this_ *CheckSuspiciousCode) priv_CheckSuspiciousCode_checkInvalidIn_1123(t *NodeTraversal, n pkg_rhino.Node) {
	_ = t
	_ = n
	if n.IsIn() {
		CheckSuspiciousCode_ReportIfNonObject(t, n.GetLastChild(), StaticGet_CheckSuspiciousCode_SUSPICIOUS_IN_OPERATOR())
	}
}

// Java method: CheckSuspiciousCode.checkNonObjectInstanceOf(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n)
func (this_ *CheckSuspiciousCode) priv_CheckSuspiciousCode_checkNonObjectInstanceOf_1124(t *NodeTraversal, n pkg_rhino.Node) {
	_ = t
	_ = n
	if n.IsInstanceOf() {
		CheckSuspiciousCode_ReportIfNonObject(t, n.GetFirstChild(), StaticGet_CheckSuspiciousCode_SUSPICIOUS_INSTANCEOF_LEFT_OPERAND())
	}
}

// Java static method: CheckSuspiciousCode.reportIfNonObject(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.jscomp.DiagnosticType diagnosticType)
func CheckSuspiciousCode_ReportIfNonObject(t *NodeTraversal, n pkg_rhino.Node, diagnosticType *DiagnosticType) bool {
	clinit_CheckSuspiciousCode()
	_ = t
	_ = n
	_ = diagnosticType
	if n.IsAdd() || !(NodeUtil_MayBeObject_1_Node(n)) {
		t.Report_3(n.GetParent(), diagnosticType, []string{})
		return true
	}
	return false
}

// Java method: CheckSuspiciousCode.checkNegatedLeftOperandOfInOperator(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n)
func (this_ *CheckSuspiciousCode) priv_CheckSuspiciousCode_checkNegatedLeftOperandOfInOperator_1125(t *NodeTraversal, n pkg_rhino.Node) {
	_ = t
	_ = n
	if n.IsIn() && n.GetFirstChild().IsNot() {
		t.Report_3(n.GetFirstChild(), StaticGet_CheckSuspiciousCode_SUSPICIOUS_NEGATED_LEFT_OPERAND_OF_IN_OPERATOR(), []string{})
	}
}

// Java method: CheckSuspiciousCode.checkLeftOperandOfLogicalOperator(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n)
func (this_ *CheckSuspiciousCode) priv_CheckSuspiciousCode_checkLeftOperandOfLogicalOperator_1126(t *NodeTraversal, n pkg_rhino.Node) {
	_ = t
	_ = n
	if n.IsOr() || n.IsAnd() {
		var operator string
		_ = operator
		if n.IsOr() {
			operator = "||"
		} else {
			operator = "&&"
		}
		var v pkg_base.Tri = this_.priv_CheckSuspiciousCode_getBooleanValueWithTypes_1127(n.GetFirstChild())
		_ = v
		if v != pkg_base.StaticGet_Tri_UNKNOWN() {
			var result string
			_ = result
			if v == pkg_base.StaticGet_Tri_TRUE() {
				result = "truthy"
			} else {
				result = "falsy"
			}
			t.Report_3(n, StaticGet_CheckSuspiciousCode_SUSPICIOUS_LEFT_OPERAND_OF_LOGICAL_OPERATOR(), []string{operator, result})
		}
	}
}

// Java static method: CheckSuspiciousCode.checkSuspiciousBreakingOutOfOptionalChain(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n)
func CheckSuspiciousCode_CheckSuspiciousBreakingOutOfOptionalChain(t *NodeTraversal, n pkg_rhino.Node) {
	clinit_CheckSuspiciousCode()
	_ = t
	_ = n
	if NodeUtil_IsOptChainNode(n) {
		var parent pkg_rhino.Node = n.GetParent()
		_ = parent
		if n.IsFirstChildOf(parent) && ((parent.IsGetProp() || parent.IsGetElem()) || parent.IsCall()) {
			t.Report_3(n, StaticGet_CheckSuspiciousCode_SUSPICIOUS_BREAKING_OUT_OF_OPTIONAL_CHAIN(), []string{})
		}
	}
}

// Java method: CheckSuspiciousCode.getBooleanValueWithTypes(com.google.javascript.rhino.Node n)
func (this_ *CheckSuspiciousCode) priv_CheckSuspiciousCode_getBooleanValueWithTypes_1127(n pkg_rhino.Node) pkg_base.Tri {
	_ = n
	switch n.TokenSwitchOrd_() {
	case int32(51) /* Token.ASSIGN */, int32(50) /* Token.COMMA */ :
		return this_.priv_CheckSuspiciousCode_getBooleanValueWithTypes_1127(n.GetLastChild())
	case int32(19) /* Token.NOT */ :
		return this_.priv_CheckSuspiciousCode_getBooleanValueWithTypes_1127(n.GetLastChild()).Not()
	case int32(69) /* Token.AND */ :
		return pkg_base.StaticGet_Tri_UNKNOWN().And(this_.priv_CheckSuspiciousCode_getBooleanValueWithTypes_1127(n.GetLastChild()))
	case int32(68) /* Token.OR */ :
		return pkg_base.StaticGet_Tri_UNKNOWN().Or(this_.priv_CheckSuspiciousCode_getBooleanValueWithTypes_1127(n.GetLastChild()))
	case int32(67) /* Token.HOOK */ :
		var trueValue pkg_base.Tri = this_.priv_CheckSuspiciousCode_getBooleanValueWithTypes_1127(n.GetSecondChild())
		_ = trueValue
		var falseValue pkg_base.Tri = this_.priv_CheckSuspiciousCode_getBooleanValueWithTypes_1127(n.GetLastChild())
		_ = falseValue
		if j2gort.CheckNotNullEnum(trueValue) == falseValue {
			return trueValue
		} else {
			return pkg_base.StaticGet_Tri_UNKNOWN()
		}
	case int32(73) /* Token.FUNCTION */, int32(105) /* Token.CLASS */, int32(23) /* Token.NEW */, int32(46) /* Token.ARRAYLIT */, int32(47) /* Token.OBJECTLIT */ :
		return pkg_base.StaticGet_Tri_TRUE()
	case int32(87) /* Token.VOID */ :
		return pkg_base.StaticGet_Tri_FALSE()
	case int32(26) /* Token.GETPROP */, int32(27) /* Token.GETELEM */, int32(30) /* Token.OPTCHAIN_GETELEM */, int32(29) /* Token.OPTCHAIN_GETPROP */ :
		return pkg_base.StaticGet_Tri_UNKNOWN()
	default:
	}
	var literalValue pkg_base.Tri = NodeUtil_GetBooleanValue(n)
	_ = literalValue
	if (literalValue != pkg_base.StaticGet_Tri_UNKNOWN()) || n.IsName() {
		return literalValue
	}
	var type_ pkg_rhino.JSType = n.GetJSType()
	_ = type_
	if type_ != nil && !type_.IsJavaNil_() {
		return pkg_base.StaticGet_Tri_UNKNOWN().And(type_.GetPossibleToBooleanOutcomes().ToTri())
	}
	return pkg_base.StaticGet_Tri_UNKNOWN()
}
func (*CheckSuspiciousCode) AlwaysShouldTraverse_()      {}
func (*CheckSuspiciousCode) CanReleaseSyntacticScopes_() {}
