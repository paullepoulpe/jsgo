/*
 * Copyright 2004 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/PeepholeFoldConstants.java (generated mechanically by j2go).

package jscomp

import (
	"math"
	"strconv"
	"sync/atomic"

	pkg_base "github.com/paullepoulpe/jsgo/transpiled/gen/base"
	pkg_colors "github.com/paullepoulpe/jsgo/transpiled/gen/colors"
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.jscomp.PeepholeFoldConstants
type PeepholeFoldConstants struct {
	AbstractPeepholeOptimizationBase
	late           bool `gson:"late"`
	shouldUseTypes bool `gson:"shouldUseTypes"`
}

func (this_ *PeepholeFoldConstants) BindSelf_(s any) {
	_ = s
	this_.AbstractPeepholeOptimizationBase.BindSelf_(s)
}
func (this_ *PeepholeFoldConstants) IsJavaNil_() bool { return this_ == nil }
func (*PeepholeFoldConstants) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.PeepholeFoldConstants", "PeepholeFoldConstants"
}
func (this_ *PeepholeFoldConstants) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.AbstractPeepholeOptimizationBase.IdentityHash_()
}
func (*PeepholeFoldConstants) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.PeepholeFoldConstants", "com.google.javascript.jscomp.PeepholeFoldConstants", func(v_ any) bool {
		return j2gort.InstanceOf[*PeepholeFoldConstants](v_)
	})
}
func (*PeepholeFoldConstants) Implements_jscomp_AbstractPeepholeOptimization() {}

var static_PeepholeFoldConstants_INVALID_GETELEM_INDEX_ERROR *DiagnosticType

func StaticGet_PeepholeFoldConstants_INVALID_GETELEM_INDEX_ERROR() *DiagnosticType {
	clinit_PeepholeFoldConstants()
	return static_PeepholeFoldConstants_INVALID_GETELEM_INDEX_ERROR
}
func StaticSet_PeepholeFoldConstants_INVALID_GETELEM_INDEX_ERROR(v *DiagnosticType) {
	clinit_PeepholeFoldConstants()
	static_PeepholeFoldConstants_INVALID_GETELEM_INDEX_ERROR = v
}

var static_PeepholeFoldConstants_FRACTIONAL_BITWISE_OPERAND *DiagnosticType

func StaticGet_PeepholeFoldConstants_FRACTIONAL_BITWISE_OPERAND() *DiagnosticType {
	clinit_PeepholeFoldConstants()
	return static_PeepholeFoldConstants_FRACTIONAL_BITWISE_OPERAND
}
func StaticSet_PeepholeFoldConstants_FRACTIONAL_BITWISE_OPERAND(v *DiagnosticType) {
	clinit_PeepholeFoldConstants()
	static_PeepholeFoldConstants_FRACTIONAL_BITWISE_OPERAND = v
}

var static_PeepholeFoldConstants_MAX_FOLD_NUMBER float64

func StaticGet_PeepholeFoldConstants_MAX_FOLD_NUMBER() float64 {
	clinit_PeepholeFoldConstants()
	return static_PeepholeFoldConstants_MAX_FOLD_NUMBER
}
func StaticSet_PeepholeFoldConstants_MAX_FOLD_NUMBER(v float64) {
	clinit_PeepholeFoldConstants()
	static_PeepholeFoldConstants_MAX_FOLD_NUMBER = v
}

var clinitDone_PeepholeFoldConstants uint32
var clinitState_PeepholeFoldConstants j2gort.ClinitState

func clinit_PeepholeFoldConstants() {
	if atomic.LoadUint32(&clinitDone_PeepholeFoldConstants) != 0 {
		return
	}
	clinitSlow_PeepholeFoldConstants()
}

func clinitSlow_PeepholeFoldConstants() {
	if !j2gort.EnterClinit(&clinitState_PeepholeFoldConstants, &clinitDone_PeepholeFoldConstants) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_PeepholeFoldConstants, &clinitDone_PeepholeFoldConstants)
	static_PeepholeFoldConstants_INVALID_GETELEM_INDEX_ERROR = DiagnosticType_Warning("JSC_INVALID_GETELEM_INDEX_ERROR", "Array index not integer: {0}")
	static_PeepholeFoldConstants_FRACTIONAL_BITWISE_OPERAND = DiagnosticType_Warning("JSC_FRACTIONAL_BITWISE_OPERAND", "Fractional bitwise operand: {0}")
	static_PeepholeFoldConstants_MAX_FOLD_NUMBER = math.Pow(float64(int32(2)), float64(int32(53)))
}

func (this_ *PeepholeFoldConstants) initFields_PeepholeFoldConstants() {
	_ = this_
}

// Java constructor: PeepholeFoldConstants(boolean late, boolean shouldUseTypes)
func NewPeepholeFoldConstants(late bool, shouldUseTypes bool) *PeepholeFoldConstants {
	clinit_PeepholeFoldConstants()
	this_ := &PeepholeFoldConstants{}
	this_.BindSelf_(this_)
	InitPeepholeFoldConstants(this_, late, shouldUseTypes)
	return this_
}

func InitPeepholeFoldConstants(this_ *PeepholeFoldConstants, late bool, shouldUseTypes bool) {
	_ = late
	_ = shouldUseTypes
	InitAbstractPeepholeOptimization(&this_.AbstractPeepholeOptimizationBase)
	this_.initFields_PeepholeFoldConstants()
	this_.late = late
	this_.shouldUseTypes = shouldUseTypes
}

// Java method: PeepholeFoldConstants.optimizeSubtree(com.google.javascript.rhino.Node subtree)
func (this_ *PeepholeFoldConstants) OptimizeSubtree(subtree pkg_rhino.Node) pkg_rhino.Node {
	_ = subtree
	switch subtree.TokenSwitchOrd_() {
	case int32(31) /* Token.OPTCHAIN_CALL */, int32(28) /* Token.CALL */ :
		return this_.priv_PeepholeFoldConstants_tryFoldUselessObjectDotDefinePropertiesCall_4218(subtree)
	case int32(23) /* Token.NEW */ :
		return this_.priv_PeepholeFoldConstants_tryFoldCtorCall_4217(subtree)
	case int32(25) /* Token.TYPEOF */ :
		return this_.priv_PeepholeFoldConstants_tryFoldTypeof_4187(subtree)
	case int32(127) /* Token.ITER_SPREAD */ :
		return this_.priv_PeepholeFoldConstants_tryFoldSpread_4223(subtree)
	case int32(46) /* Token.ARRAYLIT */, int32(47) /* Token.OBJECTLIT */ :
		return this_.priv_PeepholeFoldConstants_tryFlattenArrayOrObjectLit_4224(subtree)
	case int32(19) /* Token.NOT */, int32(21) /* Token.POS */, int32(22) /* Token.NEG */, int32(20) /* Token.BITNOT */ :
		this_.priv_PeepholeFoldConstants_tryReduceOperandsForOp_4184(subtree)
		return this_.priv_PeepholeFoldConstants_tryFoldUnaryOperator_4188(subtree)
	case int32(87) /* Token.VOID */ :
		return this_.priv_PeepholeFoldConstants_tryReduceVoid_4183(subtree)
	case int32(29) /* Token.OPTCHAIN_GETPROP */, int32(26) /* Token.GETPROP */ :
		return this_.priv_PeepholeFoldConstants_tryFoldGetProp_4221(subtree)
	case int32(131) /* Token.TEMPLATELIT */ :
		return this_.priv_PeepholeFoldConstants_tryFoldTemplateLiteralSubstitutions_4208(subtree)
	default:
		this_.priv_PeepholeFoldConstants_tryReduceOperandsForOp_4184(subtree)
		return this_.priv_PeepholeFoldConstants_tryFoldBinaryOperator_4182(subtree)
	}
	panic("unreachable switch expression")
}

// Java method: PeepholeFoldConstants.tryFoldBinaryOperator(com.google.javascript.rhino.Node subtree)
func (this_ *PeepholeFoldConstants) priv_PeepholeFoldConstants_tryFoldBinaryOperator_4182(subtree pkg_rhino.Node) pkg_rhino.Node {
	_ = subtree
	var left pkg_rhino.Node = subtree.GetFirstChild()
	_ = left
	if left == 0 {
		return subtree
	}
	var right pkg_rhino.Node = left.GetNext()
	_ = right
	if right == 0 {
		return subtree
	}
	switch subtree.TokenSwitchOrd_() {
	case int32(27) /* Token.GETELEM */, int32(30) /* Token.OPTCHAIN_GETELEM */ :
		return this_.priv_PeepholeFoldConstants_tryFoldGetElem_4220(subtree, left, right)
	case int32(45) /* Token.INSTANCEOF */ :
		return this_.priv_PeepholeFoldConstants_tryFoldInstanceof_4190(subtree, left, right)
	case int32(69) /* Token.AND */, int32(68) /* Token.OR */ :
		return this_.priv_PeepholeFoldConstants_tryFoldAndOr_4193(subtree, left, right)
	case int32(70) /* Token.COALESCE */ :
		return this_.priv_PeepholeFoldConstants_tryFoldCoalesce_4194(subtree, left, right)
	case int32(10) /* Token.LSH */, int32(11) /* Token.RSH */, int32(12) /* Token.URSH */ :
		return this_.priv_PeepholeFoldConstants_tryFoldShift_4215(subtree, left, right)
	case int32(51) /* Token.ASSIGN */ :
		return this_.priv_PeepholeFoldConstants_tryFoldAssign_4191(subtree, left, right)
	case int32(52) /* Token.ASSIGN_BITOR */, int32(53) /* Token.ASSIGN_BITXOR */, int32(54) /* Token.ASSIGN_BITAND */, int32(55) /* Token.ASSIGN_LSH */, int32(56) /* Token.ASSIGN_RSH */, int32(57) /* Token.ASSIGN_URSH */, int32(58) /* Token.ASSIGN_ADD */, int32(59) /* Token.ASSIGN_SUB */, int32(60) /* Token.ASSIGN_MUL */, int32(61) /* Token.ASSIGN_DIV */, int32(62) /* Token.ASSIGN_MOD */, int32(63) /* Token.ASSIGN_EXPONENT */ :
		return this_.priv_PeepholeFoldConstants_tryUnfoldAssignOp_4192(subtree, left, right)
	case int32(13) /* Token.ADD */ :
		return this_.priv_PeepholeFoldConstants_tryFoldAdd_4210(subtree, left, right)
	case int32(14) /* Token.SUB */, int32(16) /* Token.DIV */, int32(17) /* Token.MOD */, int32(15) /* Token.MUL */, int32(18) /* Token.EXPONENT */ :
		return this_.priv_PeepholeFoldConstants_tryFoldArithmeticOp_4197(subtree, left, right)
	case int32(3) /* Token.BITAND */, int32(1) /* Token.BITOR */, int32(2) /* Token.BITXOR */ :
		var result pkg_rhino.Node = this_.priv_PeepholeFoldConstants_tryFoldArithmeticOp_4197(subtree, left, right)
		_ = result
		if result != subtree {
			return result
		}
		return this_.priv_PeepholeFoldConstants_tryFoldLeftChildOp_4203(subtree, left, right)
	case int32(6) /* Token.LT */, int32(8) /* Token.GT */, int32(7) /* Token.LE */, int32(9) /* Token.GE */, int32(4) /* Token.EQ */, int32(5) /* Token.NE */, int32(40) /* Token.SHEQ */, int32(41) /* Token.SHNE */ :
		return this_.priv_PeepholeFoldConstants_tryFoldComparison_4216(subtree, left, right)
	default:
		return subtree
	}
	panic("unreachable switch expression")
}

// Java method: PeepholeFoldConstants.tryReduceVoid(com.google.javascript.rhino.Node n)
func (this_ *PeepholeFoldConstants) priv_PeepholeFoldConstants_tryReduceVoid_4183(n pkg_rhino.Node) pkg_rhino.Node {
	_ = n
	var child pkg_rhino.Node = n.GetFirstChild()
	_ = child
	if (!(child.IsNumber()) || (child.GetDouble() != float64(0.0))) && !(this_.MayHaveSideEffects(n)) {
		child.ReplaceWith(pkg_rhino.IR_Number(float64(int32(0))))
		this_.ReportChangeToEnclosingScope(n)
	}
	return n
}

// Java method: PeepholeFoldConstants.tryReduceOperandsForOp(com.google.javascript.rhino.Node n)
func (this_ *PeepholeFoldConstants) priv_PeepholeFoldConstants_tryReduceOperandsForOp_4184(n pkg_rhino.Node) {
	_ = n
	switch n.TokenSwitchOrd_() {
	case int32(13) /* Token.ADD */ :
		var left pkg_rhino.Node = n.GetFirstChild()
		_ = left
		var right pkg_rhino.Node = n.GetLastChild()
		_ = right
		if !(NodeUtil_MayBeString_2(left, this_.shouldUseTypes)) && !(NodeUtil_MayBeString_2(right, this_.shouldUseTypes)) {
			this_.priv_PeepholeFoldConstants_tryConvertOperandsToNumber_4185(n)
		}
	case int32(52) /* Token.ASSIGN_BITOR */, int32(53) /* Token.ASSIGN_BITXOR */, int32(54) /* Token.ASSIGN_BITAND */, int32(55) /* Token.ASSIGN_LSH */, int32(56) /* Token.ASSIGN_RSH */, int32(57) /* Token.ASSIGN_URSH */, int32(59) /* Token.ASSIGN_SUB */, int32(60) /* Token.ASSIGN_MUL */, int32(62) /* Token.ASSIGN_MOD */, int32(61) /* Token.ASSIGN_DIV */ :
		this_.priv_PeepholeFoldConstants_tryConvertToNumber_4186(n.GetLastChild())
	case int32(20) /* Token.BITNOT */, int32(1) /* Token.BITOR */, int32(2) /* Token.BITXOR */, int32(3) /* Token.BITAND */, int32(10) /* Token.LSH */, int32(11) /* Token.RSH */, int32(12) /* Token.URSH */, int32(14) /* Token.SUB */, int32(15) /* Token.MUL */, int32(17) /* Token.MOD */, int32(16) /* Token.DIV */, int32(21) /* Token.POS */, int32(22) /* Token.NEG */, int32(18) /* Token.EXPONENT */ :
		this_.priv_PeepholeFoldConstants_tryConvertOperandsToNumber_4185(n)
	default:
	}
}

// Java method: PeepholeFoldConstants.tryConvertOperandsToNumber(com.google.javascript.rhino.Node n)
func (this_ *PeepholeFoldConstants) priv_PeepholeFoldConstants_tryConvertOperandsToNumber_4185(n pkg_rhino.Node) {
	_ = n
	var next pkg_rhino.Node = pkg_rhino.Node(0)
	_ = next
	{
		var c pkg_rhino.Node = n.GetFirstChild()
		_ = c
		for ; c != 0; c = next {
			next = c.GetNext()
			this_.priv_PeepholeFoldConstants_tryConvertToNumber_4186(c)
		}
	}
}

// Java method: PeepholeFoldConstants.tryConvertToNumber(com.google.javascript.rhino.Node n)
func (this_ *PeepholeFoldConstants) priv_PeepholeFoldConstants_tryConvertToNumber_4186(n pkg_rhino.Node) {
	_ = n
	switch n.TokenSwitchOrd_() {
	case int32(33) /* Token.NUMBER */ :
		return
	case int32(69) /* Token.AND */, int32(68) /* Token.OR */, int32(50) /* Token.COMMA */, int32(70) /* Token.COALESCE */ :
		this_.priv_PeepholeFoldConstants_tryConvertToNumber_4186(n.GetLastChild())
		return
	case int32(67) /* Token.HOOK */ :
		this_.priv_PeepholeFoldConstants_tryConvertToNumber_4186(n.GetSecondChild())
		this_.priv_PeepholeFoldConstants_tryConvertToNumber_4186(n.GetLastChild())
		return
	case int32(32) /* Token.NAME */ :
		if !(NodeUtil_IsUndefined(n)) {
			return
		}
	default:
	}
	var result *j2gort.BoxDouble = this_.GetSideEffectFreeNumberValue(n)
	_ = result
	if result == nil {
		return
	}
	var value float64 = j2gort.UnboxDouble(result)
	_ = value
	var replacement pkg_rhino.Node = NodeUtil_NumberNode(value, n)
	_ = replacement
	if replacement.IsEquivalentTo_1(n) {
		return
	}
	n.ReplaceWith(replacement)
	this_.ReportChangeToEnclosingScope(replacement)
}

// Java method: PeepholeFoldConstants.tryFoldTypeof(com.google.javascript.rhino.Node originalTypeofNode)
func (this_ *PeepholeFoldConstants) priv_PeepholeFoldConstants_tryFoldTypeof_4187(originalTypeofNode pkg_rhino.Node) pkg_rhino.Node {
	_ = originalTypeofNode
	if !(originalTypeofNode.IsTypeOf()) {
		j2gort.CheckArgument(false)
	}
	var argumentNode pkg_rhino.Node = originalTypeofNode.GetFirstChild()
	_ = argumentNode
	if (argumentNode == 0) || !(NodeUtil_IsLiteralValue(argumentNode, true)) {
		return originalTypeofNode
	}
	var typeNameString string = j2gort.NullString
	_ = typeNameString
	switch argumentNode.TokenSwitchOrd_() {
	case int32(73) /* Token.FUNCTION */ :
		typeNameString = "function"
	case int32(35) /* Token.STRINGLIT */ :
		typeNameString = "string"
	case int32(33) /* Token.NUMBER */ :
		typeNameString = "number"
	case int32(39) /* Token.TRUE */, int32(38) /* Token.FALSE */ :
		typeNameString = "boolean"
	case int32(36) /* Token.NULL */, int32(47) /* Token.OBJECTLIT */, int32(46) /* Token.ARRAYLIT */ :
		typeNameString = "object"
	case int32(87) /* Token.VOID */ :
		typeNameString = "undefined"
	case int32(32) /* Token.NAME */ :
		if j2gort.CheckNotNullString(argumentNode.GetString()) == "undefined" {
			typeNameString = "undefined"
		}
	default:
	}
	if typeNameString != j2gort.NullString {
		var newNode pkg_rhino.Node = pkg_rhino.IR_String(typeNameString)
		_ = newNode
		this_.ReportChangeToEnclosingScope(originalTypeofNode)
		originalTypeofNode.ReplaceWith(newNode)
		this_.MarkFunctionsDeleted(originalTypeofNode)
		return newNode
	}
	return originalTypeofNode
}

// Java method: PeepholeFoldConstants.tryFoldUnaryOperator(com.google.javascript.rhino.Node n)
func (this_ *PeepholeFoldConstants) priv_PeepholeFoldConstants_tryFoldUnaryOperator_4188(n pkg_rhino.Node) pkg_rhino.Node {
	_ = n
	if !(n.HasOneChild()) {
		j2gort.CheckStateMsg(false, j2gort.BoxNode(n))
	}
	var left pkg_rhino.Node = n.GetFirstChild()
	_ = left
	var parent pkg_rhino.Node = n.GetParent()
	_ = parent
	if left == 0 {
		return n
	}
	var leftVal pkg_base.Tri = this_.GetSideEffectFreeBooleanValue(left)
	_ = leftVal
	if leftVal == pkg_base.StaticGet_Tri_UNKNOWN() {
		return n
	}
	switch n.TokenSwitchOrd_() {
	case int32(19) /* Token.NOT */ :
		if this_.late && left.IsNumber() {
			var numValue float64 = left.GetDouble()
			_ = numValue
			if (numValue == float64(int32(0))) || (numValue == float64(int32(1))) {
				return n
			}
		}
		var replacementNode pkg_rhino.Node = NodeUtil_BooleanNode(!(leftVal.ToBoolean(true)))
		_ = replacementNode
		n.ReplaceWith(replacementNode)
		this_.ReportChangeToEnclosingScope(parent)
		return replacementNode
	case int32(21) /* Token.POS */ :
		if NodeUtil_IsNumericResult(left) {
			n.ReplaceWith(left.Detach())
			this_.ReportChangeToEnclosingScope(parent)
			return left
		}
		return n
	case int32(22) /* Token.NEG */ :
		var result pkg_rhino.Node = pkg_rhino.Node(0)
		_ = result
		if left.IsName() && (j2gort.CheckNotNullString(left.GetString()) == "NaN") {
			result = left.Detach()
		} else if left.IsNeg() {
			var leftLeft pkg_rhino.Node = left.GetOnlyChild()
			_ = leftLeft
			if leftLeft.IsBigInt() || leftLeft.IsNumber() {
				result = leftLeft.Detach()
			}
		}
		if result != 0 {
			n.ReplaceWith(result)
			this_.ReportChangeToEnclosingScope(parent)
			return result
		}
		return n
	case int32(20) /* Token.BITNOT */ :
		var doubleVal *j2gort.BoxDouble = this_.GetSideEffectFreeNumberValue(left)
		_ = doubleVal
		if doubleVal != nil {
			if pkg_base.JSCompDoubles_IsMathematicalInteger(j2gort.UnboxDouble(doubleVal)) {
				var intVal int32 = pkg_base.JSCompDoubles_EcmascriptToInt32(j2gort.UnboxDouble(doubleVal))
				_ = intVal
				var notIntValNode pkg_rhino.Node = NodeUtil_NumberNode(float64(^(intVal)), left)
				_ = notIntValNode
				n.ReplaceWith(notIntValNode)
				this_.ReportChangeToEnclosingScope(parent)
				return notIntValNode
			} else {
				this_.Report(StaticGet_PeepholeFoldConstants_FRACTIONAL_BITWISE_OPERAND(), left)
				return n
			}
		}
		var bigintVal *j2gort.Ext_java_math_BigInteger = this_.GetSideEffectFreeBigIntValue(n)
		_ = bigintVal
		if bigintVal != nil {
			var bigintNotNode pkg_rhino.Node = PeepholeFoldConstants_BigintNode(bigintVal, n)
			_ = bigintNotNode
			n.ReplaceWith(bigintNotNode)
			this_.ReportChangeToEnclosingScope(parent)
			return bigintNotNode
		}
		return n
	default:
		return n
	}
	panic("unreachable switch expression")
}

// Java method: PeepholeFoldConstants.isReasonableDoubleValue(java.lang.@org.jspecify.annotations.Nullable Double x)
func (this_ *PeepholeFoldConstants) priv_PeepholeFoldConstants_isReasonableDoubleValue_4189(x *j2gort.BoxDouble) bool {
	_ = x
	return (((x != nil) && !(x.IsInfinite())) && !(x.IsNaN()))
}

// Java method: PeepholeFoldConstants.tryFoldInstanceof(com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node left, com.google.javascript.rhino.Node right)
func (this_ *PeepholeFoldConstants) priv_PeepholeFoldConstants_tryFoldInstanceof_4190(n pkg_rhino.Node, left pkg_rhino.Node, right pkg_rhino.Node) pkg_rhino.Node {
	_ = n
	_ = left
	_ = right
	if !(n.IsInstanceOf()) {
		j2gort.CheckArgument(false)
	}
	if NodeUtil_IsLiteralValue(left, true) && !(this_.MayHaveSideEffects(right)) {
		var replacementNode pkg_rhino.Node = pkg_rhino.Node(0)
		_ = replacementNode
		if NodeUtil_IsImmutableValue(left) {
			replacementNode = pkg_rhino.IR_FalseNode()
		} else if left.IsTemplateLit() {
			if !(this_.MayHaveSideEffects(left)) {
				replacementNode = pkg_rhino.IR_FalseNode()
			}
		} else if right.IsName() && (right.GetString() == "Object") {
			replacementNode = pkg_rhino.IR_TrueNode()
		}
		if replacementNode != 0 {
			n.ReplaceWith(replacementNode)
			this_.ReportChangeToEnclosingScope(replacementNode)
			this_.MarkFunctionsDeleted(n)
			return replacementNode
		}
	}
	return n
}

// Java method: PeepholeFoldConstants.tryFoldAssign(com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node left, com.google.javascript.rhino.Node right)
func (this_ *PeepholeFoldConstants) priv_PeepholeFoldConstants_tryFoldAssign_4191(n pkg_rhino.Node, left pkg_rhino.Node, right pkg_rhino.Node) pkg_rhino.Node {
	_ = n
	_ = left
	_ = right
	if !(n.IsAssign()) {
		j2gort.CheckArgument(false)
	}
	if !(this_.late) {
		return n
	}
	var operator pkg_rhino.Token = right.GetToken()
	_ = operator
	var newType pkg_rhino.Token = pkg_rhino.Token(0)
	_ = newType
	switch (operator).Ordinal() {
	case int32(13) /* Token.ADD */ :
		newType = pkg_rhino.StaticGet_Token_ASSIGN_ADD()
	case int32(3) /* Token.BITAND */ :
		newType = pkg_rhino.StaticGet_Token_ASSIGN_BITAND()
	case int32(1) /* Token.BITOR */ :
		newType = pkg_rhino.StaticGet_Token_ASSIGN_BITOR()
	case int32(2) /* Token.BITXOR */ :
		newType = pkg_rhino.StaticGet_Token_ASSIGN_BITXOR()
	case int32(16) /* Token.DIV */ :
		newType = pkg_rhino.StaticGet_Token_ASSIGN_DIV()
	case int32(10) /* Token.LSH */ :
		newType = pkg_rhino.StaticGet_Token_ASSIGN_LSH()
	case int32(17) /* Token.MOD */ :
		newType = pkg_rhino.StaticGet_Token_ASSIGN_MOD()
	case int32(15) /* Token.MUL */ :
		newType = pkg_rhino.StaticGet_Token_ASSIGN_MUL()
	case int32(11) /* Token.RSH */ :
		newType = pkg_rhino.StaticGet_Token_ASSIGN_RSH()
	case int32(14) /* Token.SUB */ :
		newType = pkg_rhino.StaticGet_Token_ASSIGN_SUB()
	case int32(12) /* Token.URSH */ :
		newType = pkg_rhino.StaticGet_Token_ASSIGN_URSH()
	case int32(18) /* Token.EXPONENT */ :
		newType = pkg_rhino.StaticGet_Token_ASSIGN_EXPONENT()
	default:
		newType = pkg_rhino.Token(0)
	}
	if newType == 0 {
		return n
	}
	if !(right.HasChildren()) || (right.GetSecondChild() != right.GetLastChild()) {
		return n
	}
	if this_.MayHaveSideEffects(left) {
		return n
	}
	var newRight pkg_rhino.Node = pkg_rhino.Node(0)
	_ = newRight
	if this_.AreNodesEqualForInlining(left, right.GetFirstChild()) {
		newRight = right.GetLastChild()
	} else if (NodeUtil_IsCommutative(operator) && !(this_.MayHaveSideEffects(right.GetFirstChild()))) && this_.AreNodesEqualForInlining(left, right.GetLastChild()) {
		newRight = right.GetFirstChild()
	} else if (NodeUtil_IsAssociative(operator) && (right.GetFirstChild().GetToken() == operator)) && this_.AreNodesEqualForInlining(left, right.GetFirstFirstChild()) {
		var newLeft pkg_rhino.Node = right.GetFirstChild().GetLastChild().Detach()
		_ = newLeft
		right.RemoveFirstChild()
		right.AddChildToFront(newLeft)
		newRight = right
	} else {
		return n
	}
	var newNode pkg_rhino.Node = pkg_rhino.NewNode_3(newType, left.Detach(), newRight.Detach())
	_ = newNode
	n.ReplaceWith(newNode)
	this_.ReportChangeToEnclosingScope(newNode)
	return newNode
}

// Java method: PeepholeFoldConstants.tryUnfoldAssignOp(com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node left, com.google.javascript.rhino.Node right)
func (this_ *PeepholeFoldConstants) priv_PeepholeFoldConstants_tryUnfoldAssignOp_4192(n pkg_rhino.Node, left pkg_rhino.Node, right pkg_rhino.Node) pkg_rhino.Node {
	_ = n
	_ = left
	_ = right
	if this_.late {
		return n
	}
	if !(n.HasChildren()) || (n.GetSecondChild() != n.GetLastChild()) {
		return n
	}
	if this_.MayHaveSideEffects(left) {
		return n
	}
	var op pkg_rhino.Token = NodeUtil_GetOpFromAssignmentOp(n)
	_ = op
	var replacement pkg_rhino.Node = pkg_rhino.IR_Assign(left.Detach(), pkg_rhino.NewNode_3(op, left.CloneTree_0(), right.Detach()).Srcref(n))
	_ = replacement
	n.ReplaceWith(replacement)
	this_.ReportChangeToEnclosingScope(replacement)
	return replacement
}

// Java method: PeepholeFoldConstants.tryFoldAndOr(com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node left, com.google.javascript.rhino.Node right)
func (this_ *PeepholeFoldConstants) priv_PeepholeFoldConstants_tryFoldAndOr_4193(n pkg_rhino.Node, left pkg_rhino.Node, right pkg_rhino.Node) pkg_rhino.Node {
	_ = n
	_ = left
	_ = right
	var parent pkg_rhino.Node = n.GetParent()
	_ = parent
	var result pkg_rhino.Node = pkg_rhino.Node(0)
	_ = result
	var dropped pkg_rhino.Node = pkg_rhino.Node(0)
	_ = dropped
	var type_ pkg_rhino.Token = n.GetToken()
	_ = type_
	var leftVal pkg_base.Tri = NodeUtil_GetBooleanValue(left)
	_ = leftVal
	if leftVal != pkg_base.StaticGet_Tri_UNKNOWN() {
		var lval bool = leftVal.ToBoolean(true)
		_ = lval
		if func() bool {
			if lval {
				return (type_ == pkg_rhino.StaticGet_Token_OR())
			}
			return (type_ == pkg_rhino.StaticGet_Token_AND())
		}() {
			result = left
			dropped = right
		} else if !(this_.MayHaveSideEffects(left)) {
			result = right
			dropped = left
		} else {
			n.DetachChildren()
			result = pkg_rhino.IR_Comma(left, right)
			dropped = pkg_rhino.Node(0)
		}
	} else if (parent.GetToken() == type_) && (n == parent.GetFirstChild()) {
		var rightValue pkg_base.Tri = NodeUtil_GetBooleanValue(right)
		_ = rightValue
		if !(this_.MayHaveSideEffects(right)) {
			if ((rightValue == pkg_base.StaticGet_Tri_FALSE()) && (type_ == pkg_rhino.StaticGet_Token_OR())) || ((rightValue == pkg_base.StaticGet_Tri_TRUE()) && (type_ == pkg_rhino.StaticGet_Token_AND())) {
				result = left
				dropped = right
			}
		}
	}
	if result != 0 {
		n.DetachChildren()
		n.ReplaceWith(result)
		this_.ReportChangeToEnclosingScope(result)
		if dropped != 0 {
			this_.MarkFunctionsDeleted(dropped)
		}
		return result
	} else {
		return n
	}
}

// Java method: PeepholeFoldConstants.tryFoldCoalesce(com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node left, com.google.javascript.rhino.Node right)
func (this_ *PeepholeFoldConstants) priv_PeepholeFoldConstants_tryFoldCoalesce_4194(n pkg_rhino.Node, left pkg_rhino.Node, right pkg_rhino.Node) pkg_rhino.Node {
	_ = n
	_ = left
	_ = right
	var result pkg_rhino.Node = pkg_rhino.Node(0)
	_ = result
	var leftVal NodeUtil_ValueType = NodeUtil_GetKnownValueType(left)
	_ = leftVal
	switch (leftVal).Ordinal() {
	case int32(1) /* ValueType.NULL */, int32(2) /* ValueType.VOID */ :
		if !(this_.MayHaveSideEffects(left)) {
			result = right
			this_.MarkFunctionsDeleted(left)
		} else {
			n.DetachChildren()
			result = pkg_rhino.IR_Comma(left, right)
		}
	case int32(3) /* ValueType.NUMBER */, int32(4) /* ValueType.BIGINT */, int32(5) /* ValueType.STRING */, int32(6) /* ValueType.BOOLEAN */, int32(7) /* ValueType.OBJECT */ :
		result = left
		this_.MarkFunctionsDeleted(right)
	case int32(0) /* ValueType.UNDETERMINED */ :
	}
	if result != 0 {
		n.DetachChildren()
		n.ReplaceWith(result)
		this_.ReportChangeToEnclosingScope(result)
		return result
	} else {
		return n
	}
}

// Java method: PeepholeFoldConstants.tryFoldAdjacentLiteralLeaves(com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node left, com.google.javascript.rhino.Node right)
func (this_ *PeepholeFoldConstants) priv_PeepholeFoldConstants_tryFoldAdjacentLiteralLeaves_4195(n pkg_rhino.Node, left pkg_rhino.Node, right pkg_rhino.Node) pkg_rhino.Node {
	_ = n
	_ = left
	_ = right
	var leftParent pkg_rhino.Node = n
	_ = leftParent
	for left.IsAdd() {
		leftParent = left
		left = left.GetSecondChild()
	}
	var rightParent pkg_rhino.Node = n
	_ = rightParent
	for right.IsAdd() {
		rightParent = right
		right = right.GetFirstChild()
	}
	if ((leftParent.IsAdd() && left.IsStringLit()) && rightParent.IsAdd()) && NodeUtil_IsLiteralValue(right, false) {
		var rightGrandparent pkg_rhino.Node = rightParent.GetParent()
		_ = rightGrandparent
		var rr pkg_rhino.Node = right.GetNext()
		_ = rr
		var foldIsTypeSafe bool = ((((rr != 0) && NodeUtil_IsStringResult(rr)) || (((((rr != 0) && right.IsStringLit()) && (rightGrandparent != 0)) && rightGrandparent.IsAdd()) && NodeUtil_IsStringResult(rightGrandparent.GetSecondChild()))) || (rr == 0))
		_ = foldIsTypeSafe
		if foldIsTypeSafe {
			var result string = (j2gort.StrValString(left.GetString()) + j2gort.StrValString(NodeUtil_GetStringValue(right)))
			_ = result
			if j2gort.CheckNotNullNode(rightParent.GetSecondChild()) == right {
				left.ReplaceWith(pkg_rhino.IR_String(result))
				this_.priv_PeepholeFoldConstants_replace_4213(rightParent, rightParent.GetFirstChild().CloneTree_1(true))
			} else {
				left.ReplaceWith(pkg_rhino.IR_String(result))
				this_.priv_PeepholeFoldConstants_replace_4213(rightParent, rightParent.GetSecondChild().CloneTree_1(true))
			}
		}
	}
	return n
}

// Java method: PeepholeFoldConstants.tryFoldAddConstantString(com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node left, com.google.javascript.rhino.Node right)
func (this_ *PeepholeFoldConstants) priv_PeepholeFoldConstants_tryFoldAddConstantString_4196(n pkg_rhino.Node, left pkg_rhino.Node, right pkg_rhino.Node) pkg_rhino.Node {
	_ = n
	_ = left
	_ = right
	if ((left.IsStringLit() || right.IsStringLit()) || left.IsArrayLit()) || right.IsArrayLit() {
		var leftString string = this_.GetSideEffectFreeStringValue(left)
		_ = leftString
		var rightString string = this_.GetSideEffectFreeStringValue(right)
		_ = rightString
		if (leftString != j2gort.NullString) && (rightString != j2gort.NullString) {
			var newStringNode pkg_rhino.Node = pkg_rhino.IR_String((j2gort.StrValString(leftString) + j2gort.StrValString(rightString)))
			_ = newStringNode
			n.ReplaceWith(newStringNode)
			this_.ReportChangeToEnclosingScope(newStringNode)
			return newStringNode
		}
	}
	return n
}

// Java method: PeepholeFoldConstants.tryFoldArithmeticOp(com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node left, com.google.javascript.rhino.Node right)
func (this_ *PeepholeFoldConstants) priv_PeepholeFoldConstants_tryFoldArithmeticOp_4197(n pkg_rhino.Node, left pkg_rhino.Node, right pkg_rhino.Node) pkg_rhino.Node {
	_ = n
	_ = left
	_ = right
	var result pkg_rhino.Node = this_.priv_PeepholeFoldConstants_performArithmeticOp_4198(n, left, right)
	_ = result
	if result != 0 {
		result.SrcrefTreeIfMissing(n)
		this_.ReportChangeToEnclosingScope(n)
		n.ReplaceWith(result)
		return result
	}
	return n
}

// Java method: PeepholeFoldConstants.performArithmeticOp(com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node left, com.google.javascript.rhino.Node right)
func (this_ *PeepholeFoldConstants) priv_PeepholeFoldConstants_performArithmeticOp_4198(n pkg_rhino.Node, left pkg_rhino.Node, right pkg_rhino.Node) pkg_rhino.Node {
	_ = n
	_ = left
	_ = right
	if n.IsAdd() && (NodeUtil_MayBeString_2(left, this_.shouldUseTypes) || NodeUtil_MayBeString_2(right, this_.shouldUseTypes)) {
		return pkg_rhino.Node(0)
	}
	if this_.priv_PeepholeFoldConstants_isBigInt_4201(left) && this_.priv_PeepholeFoldConstants_isBigInt_4201(right) {
		return this_.priv_PeepholeFoldConstants_performBigIntArithmeticOp_4199(n, left, right)
	}
	var result float64 = float64(0)
	_ = result
	var lValObj *j2gort.BoxDouble = this_.GetSideEffectFreeNumberValue(left)
	_ = lValObj
	var rValObj *j2gort.BoxDouble = this_.GetSideEffectFreeNumberValue(right)
	_ = rValObj
	if (((lValObj == nil) && (rValObj == nil)) || !(this_.priv_PeepholeFoldConstants_isNumeric_4200(left))) || !(this_.priv_PeepholeFoldConstants_isNumeric_4200(right)) {
		return pkg_rhino.Node(0)
	}
	switch n.TokenSwitchOrd_() {
	case int32(13) /* Token.ADD */ :
		if (lValObj != nil) && (rValObj != nil) {
			return this_.priv_PeepholeFoldConstants_maybeReplaceBinaryOpWithNumericResult_4202((j2gort.UnboxDouble(lValObj) + j2gort.UnboxDouble(rValObj)), j2gort.UnboxDouble(lValObj), j2gort.UnboxDouble(rValObj))
		}
		if (lValObj != nil) && (j2gort.UnboxDouble(lValObj) == float64(int32(0))) {
			return right.CloneTree_1(true)
		} else if (rValObj != nil) && (j2gort.UnboxDouble(rValObj) == float64(int32(0))) {
			return left.CloneTree_1(true)
		}
		return pkg_rhino.Node(0)
	case int32(14) /* Token.SUB */ :
		if (lValObj != nil) && (rValObj != nil) {
			return this_.priv_PeepholeFoldConstants_maybeReplaceBinaryOpWithNumericResult_4202((j2gort.UnboxDouble(lValObj) - j2gort.UnboxDouble(rValObj)), j2gort.UnboxDouble(lValObj), j2gort.UnboxDouble(rValObj))
		}
		if (lValObj != nil) && (j2gort.UnboxDouble(lValObj) == float64(int32(0))) {
			return pkg_rhino.IR_Neg(right.CloneTree_1(true))
		} else if (rValObj != nil) && (j2gort.UnboxDouble(rValObj) == float64(int32(0))) {
			return left.CloneTree_1(true)
		}
		return pkg_rhino.Node(0)
	case int32(15) /* Token.MUL */ :
		if (lValObj != nil) && (rValObj != nil) {
			return this_.priv_PeepholeFoldConstants_maybeReplaceBinaryOpWithNumericResult_4202((j2gort.UnboxDouble(lValObj) * j2gort.UnboxDouble(rValObj)), j2gort.UnboxDouble(lValObj), j2gort.UnboxDouble(rValObj))
		}
		if lValObj != nil {
			if j2gort.UnboxDouble(lValObj) == float64(int32(1)) {
				return right.CloneTree_1(true)
			}
		} else {
			if j2gort.UnboxDouble(rValObj) == float64(int32(1)) {
				return left.CloneTree_1(true)
			}
		}
		return pkg_rhino.Node(0)
	case int32(16) /* Token.DIV */ :
		if (lValObj != nil) && (rValObj != nil) {
			if j2gort.UnboxDouble(rValObj) == float64(int32(0)) {
				return pkg_rhino.Node(0)
			}
			return this_.priv_PeepholeFoldConstants_maybeReplaceBinaryOpWithNumericResult_4202((j2gort.UnboxDouble(lValObj) / j2gort.UnboxDouble(rValObj)), j2gort.UnboxDouble(lValObj), j2gort.UnboxDouble(rValObj))
		}
		if rValObj != nil {
			if j2gort.UnboxDouble(rValObj) == float64(int32(1)) {
				return left.CloneTree_1(true)
			}
		}
		return pkg_rhino.Node(0)
	case int32(18) /* Token.EXPONENT */ :
		if (lValObj != nil) && (rValObj != nil) {
			return this_.priv_PeepholeFoldConstants_maybeReplaceBinaryOpWithNumericResult_4202(math.Pow(j2gort.UnboxDouble(lValObj), j2gort.UnboxDouble(rValObj)), j2gort.UnboxDouble(lValObj), j2gort.UnboxDouble(rValObj))
		}
		return pkg_rhino.Node(0)
	default:
	}
	if (lValObj == nil) || (rValObj == nil) {
		return pkg_rhino.Node(0)
	}
	var lval float64 = j2gort.UnboxDouble(lValObj)
	_ = lval
	var rval float64 = j2gort.UnboxDouble(rValObj)
	_ = rval
	switch n.TokenSwitchOrd_() {
	case int32(3) /* Token.BITAND */ :
		result = float64((pkg_base.JSCompDoubles_EcmascriptToInt32(lval) & pkg_base.JSCompDoubles_EcmascriptToInt32(rval)))
	case int32(1) /* Token.BITOR */ :
		result = float64((pkg_base.JSCompDoubles_EcmascriptToInt32(lval) | pkg_base.JSCompDoubles_EcmascriptToInt32(rval)))
	case int32(2) /* Token.BITXOR */ :
		result = float64((pkg_base.JSCompDoubles_EcmascriptToInt32(lval) ^ pkg_base.JSCompDoubles_EcmascriptToInt32(rval)))
	case int32(17) /* Token.MOD */ :
		if rval == float64(int32(0)) {
			return pkg_rhino.Node(0)
		}
		result = math.Mod(lval, rval)
	default:
		panic(j2gort.AsIface[javalang.IIllegalStateException](j2gort.ConstructShimThrowable("java.lang.IllegalStateException", "Unexpected arithmetic operator: "+j2gort.StrVal(n.GetToken()))))
	}
	return this_.priv_PeepholeFoldConstants_maybeReplaceBinaryOpWithNumericResult_4202(result, lval, rval)
}

// Java method: PeepholeFoldConstants.performBigIntArithmeticOp(com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node left, com.google.javascript.rhino.Node right)
func (this_ *PeepholeFoldConstants) priv_PeepholeFoldConstants_performBigIntArithmeticOp_4199(n pkg_rhino.Node, left pkg_rhino.Node, right pkg_rhino.Node) pkg_rhino.Node {
	_ = n
	_ = left
	_ = right
	var lVal *j2gort.Ext_java_math_BigInteger = this_.GetSideEffectFreeBigIntValue(left)
	_ = lVal
	var rVal *j2gort.Ext_java_math_BigInteger = this_.GetSideEffectFreeBigIntValue(right)
	_ = rVal
	if (lVal != nil) && (rVal != nil) {
		switch n.TokenSwitchOrd_() {
		case int32(13) /* Token.ADD */ :
			return PeepholeFoldConstants_BigintNode(lVal.Add(rVal), n)
		case int32(14) /* Token.SUB */ :
			return PeepholeFoldConstants_BigintNode(lVal.Subtract(rVal), n)
		case int32(15) /* Token.MUL */ :
			return PeepholeFoldConstants_BigintNode(lVal.Multiply(rVal), n)
		case int32(16) /* Token.DIV */ :
			if !(rVal.Equals(j2gort.PtrToAny(j2gort.Ext_java_math_BigInteger_F_ZERO))) {
				return PeepholeFoldConstants_BigintNode(lVal.Divide(rVal), n)
			} else {
				return pkg_rhino.Node(0)
			}
		case int32(18) /* Token.EXPONENT */ :
			var didRet_1 bool
			var retVal_2 pkg_rhino.Node
			func() {
				defer func() {
					if rec_3 := recover(); rec_3 != nil {
						jex_4 := javalang.RecoverException(rec_3)
						_ = jex_4
						if j2gort.InstanceOf[javalang.IArithmeticException](jex_4) {
							exception := j2gort.AsIface[javalang.IArithmeticException](jex_4)
							_ = exception
							retVal_2 = pkg_rhino.Node(0)
							didRet_1 = true
							return
						} else {
							panic(rec_3)
						}
					}
				}()
				retVal_2 = PeepholeFoldConstants_BigintNode(lVal.Pow(rVal.IntValueExact()), n)
				didRet_1 = true
				return
			}()
			if didRet_1 {
				return retVal_2
			}
		case int32(17) /* Token.MOD */ :
			if !(rVal.Equals(j2gort.PtrToAny(j2gort.Ext_java_math_BigInteger_F_ZERO))) {
				return PeepholeFoldConstants_BigintNode(lVal.Remainder(rVal), n)
			} else {
				return pkg_rhino.Node(0)
			}
		case int32(3) /* Token.BITAND */ :
			return PeepholeFoldConstants_BigintNode(lVal.And(rVal), n)
		case int32(1) /* Token.BITOR */ :
			return PeepholeFoldConstants_BigintNode(lVal.Or(rVal), n)
		case int32(2) /* Token.BITXOR */ :
			return PeepholeFoldConstants_BigintNode(lVal.Xor(rVal), n)
		default:
			return pkg_rhino.Node(0)
		}
	}
	switch n.TokenSwitchOrd_() {
	case int32(13) /* Token.ADD */ :
		if (lVal != nil) && lVal.Equals(j2gort.PtrToAny(j2gort.Ext_java_math_BigInteger_F_ZERO)) {
			return right.CloneTree_1(true)
		} else if (rVal != nil) && rVal.Equals(j2gort.PtrToAny(j2gort.Ext_java_math_BigInteger_F_ZERO)) {
			return left.CloneTree_1(true)
		}
		return pkg_rhino.Node(0)
	case int32(14) /* Token.SUB */ :
		if (lVal != nil) && lVal.Equals(j2gort.PtrToAny(j2gort.Ext_java_math_BigInteger_F_ZERO)) {
			return pkg_rhino.IR_Neg(right.CloneTree_1(true))
		} else if (rVal != nil) && rVal.Equals(j2gort.PtrToAny(j2gort.Ext_java_math_BigInteger_F_ZERO)) {
			return left.CloneTree_1(true)
		}
		return pkg_rhino.Node(0)
	case int32(15) /* Token.MUL */ :
		if (lVal != nil) && lVal.Equals(j2gort.PtrToAny(j2gort.Ext_java_math_BigInteger_F_ONE)) {
			return right.CloneTree_1(true)
		} else if (rVal != nil) && rVal.Equals(j2gort.PtrToAny(j2gort.Ext_java_math_BigInteger_F_ONE)) {
			return left.CloneTree_1(true)
		}
		return pkg_rhino.Node(0)
	case int32(16) /* Token.DIV */ :
		if (rVal != nil) && rVal.Equals(j2gort.PtrToAny(j2gort.Ext_java_math_BigInteger_F_ONE)) {
			return left.CloneTree_1(true)
		}
		return pkg_rhino.Node(0)
	default:
		return pkg_rhino.Node(0)
	}
	panic("unreachable switch expression")
}

// Java method: PeepholeFoldConstants.isNumeric(com.google.javascript.rhino.Node n)
func (this_ *PeepholeFoldConstants) priv_PeepholeFoldConstants_isNumeric_4200(n pkg_rhino.Node) bool {
	_ = n
	if NodeUtil_IsNumericResult(n) {
		return true
	}
	if this_.shouldUseTypes {
		return ((n.GetColor() != nil) && n.GetColor().Equals(j2gort.PtrToAny(pkg_colors.StaticGet_StandardColors_NUMBER())))
	}
	return false
}

// Java method: PeepholeFoldConstants.isBigInt(com.google.javascript.rhino.Node n)
func (this_ *PeepholeFoldConstants) priv_PeepholeFoldConstants_isBigInt_4201(n pkg_rhino.Node) bool {
	_ = n
	if NodeUtil_IsBigIntResult(n) {
		return true
	}
	if this_.shouldUseTypes {
		return ((n.GetColor() != nil) && n.GetColor().Equals(j2gort.PtrToAny(pkg_colors.StaticGet_StandardColors_BIGINT())))
	}
	return false
}

// Java method: PeepholeFoldConstants.maybeReplaceBinaryOpWithNumericResult(double result, double lval, double rval)
func (this_ *PeepholeFoldConstants) priv_PeepholeFoldConstants_maybeReplaceBinaryOpWithNumericResult_4202(result float64, lval float64, rval float64) pkg_rhino.Node {
	_ = result
	_ = lval
	_ = rval
	if ((((j2gort.JStrLen(javalang.DoubleToString(result)) <= ((j2gort.JStrLen(javalang.DoubleToString(lval)) + j2gort.JStrLen(javalang.DoubleToString(rval))) + int32(1))) && (math.Abs(result) <= StaticGet_PeepholeFoldConstants_MAX_FOLD_NUMBER())) || math.IsNaN(result)) || (result == math.Inf(1))) || (result == math.Inf(-1)) {
		return NodeUtil_NumberNode(result, pkg_rhino.Node(0))
	}
	return pkg_rhino.Node(0)
}

// Java method: PeepholeFoldConstants.tryFoldLeftChildOp(com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node left, com.google.javascript.rhino.Node right)
func (this_ *PeepholeFoldConstants) priv_PeepholeFoldConstants_tryFoldLeftChildOp_4203(n pkg_rhino.Node, left pkg_rhino.Node, right pkg_rhino.Node) pkg_rhino.Node {
	_ = n
	_ = left
	_ = right
	var opType pkg_rhino.Token = n.GetToken()
	_ = opType
	if !(NodeUtil_IsAssociative(opType) && NodeUtil_IsCommutative(opType)) {
		j2gort.CheckState(false)
	}
	var rightValObj *j2gort.BoxDouble = this_.GetSideEffectFreeNumberValue(right)
	_ = rightValObj
	var rightBigInt *j2gort.Ext_java_math_BigInteger = this_.GetSideEffectFreeBigIntValue(right)
	_ = rightBigInt
	if ((rightValObj != nil) || (rightBigInt != nil)) && (left.GetToken() == opType) {
		if !(left.HasTwoChildren()) {
			j2gort.CheckState(false)
		}
		var ll pkg_rhino.Node = left.GetFirstChild()
		_ = ll
		var lr pkg_rhino.Node = ll.GetNext()
		_ = lr
		var valueToCombine pkg_rhino.Node = ll
		_ = valueToCombine
		var replacement pkg_rhino.Node = this_.priv_PeepholeFoldConstants_performArithmeticOp_4198(n, valueToCombine, right)
		_ = replacement
		if replacement == 0 {
			valueToCombine = lr
			replacement = this_.priv_PeepholeFoldConstants_performArithmeticOp_4198(n, valueToCombine, right)
		}
		if replacement != 0 {
			valueToCombine.Detach()
			left.ReplaceWith(left.RemoveFirstChild())
			replacement.SrcrefTreeIfMissing(right)
			right.ReplaceWith(replacement)
			this_.ReportChangeToEnclosingScope(n)
		}
	}
	return n
}

// Java method: PeepholeFoldConstants.foldAddTemplateLiterals(com.google.javascript.rhino.Node oldNode, com.google.javascript.rhino.Node left, com.google.javascript.rhino.Node right)
func (this_ *PeepholeFoldConstants) priv_PeepholeFoldConstants_foldAddTemplateLiterals_4204(oldNode pkg_rhino.Node, left pkg_rhino.Node, right pkg_rhino.Node) pkg_rhino.Node {
	_ = oldNode
	_ = left
	_ = right
	var lastLeft pkg_rhino.Node = left.GetLastChild()
	_ = lastLeft
	var firstRight pkg_rhino.Node = right.GetFirstChild()
	_ = firstRight
	if this_.priv_PeepholeFoldConstants_endsWithUnescapedDollar_4207(lastLeft.GetRawString()) && javalang.StringStartsWith(j2gort.CheckNotNullString(firstRight.GetRawString()), "{") {
		return oldNode
	}
	var mergeNode pkg_rhino.Node = pkg_rhino.IR_TemplateLiteralString((j2gort.StrValString(lastLeft.GetCookedString()) + j2gort.StrValString(firstRight.GetCookedString())), (j2gort.StrValString(lastLeft.GetRawString()) + j2gort.StrValString(firstRight.GetRawString())))
	_ = mergeNode
	mergeNode.SrcrefTreeIfMissing(lastLeft)
	lastLeft.Detach()
	firstRight.Detach()
	left.AddChildToBack(mergeNode)
	left.AddChildrenToBack(right.RemoveChildren())
	left.Detach()
	return this_.priv_PeepholeFoldConstants_replace_4213(oldNode, left)
}

// Java method: PeepholeFoldConstants.isSimpleTemplateStringValue(java.lang.String s)
func (this_ *PeepholeFoldConstants) priv_PeepholeFoldConstants_isSimpleTemplateStringValue_4205(s string) bool {
	_ = s
	{
		var i int32 = int32(0)
		_ = i
		for ; i < j2gort.JStrLen(s); i = i + 1 {
			switch j2gort.JStrCharAt(s, i) {
			case uint16(96), uint16(36), uint16(123), uint16(92), uint16(10), uint16(13), uint16(8232), uint16(8233):
				return false
			default:
			}
		}
	}
	return true
}

// Java method: PeepholeFoldConstants.endsWithUnescapedBackslash(java.lang.String raw)
func (this_ *PeepholeFoldConstants) priv_PeepholeFoldConstants_endsWithUnescapedBackslash_4206(raw string) bool {
	_ = raw
	var backslashCount int32 = int32(0)
	_ = backslashCount
	{
		var i int32 = (j2gort.JStrLen(raw) - int32(1))
		_ = i
		for ; i >= int32(0); i = i - 1 {
			if int32(j2gort.JStrCharAt(raw, i)) == int32(uint16(92)) {
				backslashCount = backslashCount + 1
			} else {
				break
			}
		}
	}
	return ((backslashCount % int32(2)) != int32(0))
}

// Java method: PeepholeFoldConstants.endsWithUnescapedDollar(java.lang.String raw)
func (this_ *PeepholeFoldConstants) priv_PeepholeFoldConstants_endsWithUnescapedDollar_4207(raw string) bool {
	_ = raw
	if !(javalang.StringEndsWith(j2gort.CheckNotNullString(raw), "$")) {
		return false
	}
	var backslashCount int32 = int32(0)
	_ = backslashCount
	{
		var i int32 = (j2gort.JStrLen(raw) - int32(2))
		_ = i
		for ; i >= int32(0); i = i - 1 {
			if int32(j2gort.JStrCharAt(raw, i)) == int32(uint16(92)) {
				backslashCount = backslashCount + 1
			} else {
				break
			}
		}
	}
	return ((backslashCount % int32(2)) == int32(0))
}

// Java method: PeepholeFoldConstants.tryFoldTemplateLiteralSubstitutions(com.google.javascript.rhino.Node n)
func (this_ *PeepholeFoldConstants) priv_PeepholeFoldConstants_tryFoldTemplateLiteralSubstitutions_4208(n pkg_rhino.Node) pkg_rhino.Node {
	_ = n
	if (n.GetParent() != 0) && n.GetParent().IsTaggedTemplateLit() {
		return n
	}
	var changed bool = false
	_ = changed
	var cur pkg_rhino.Node = n.GetFirstChild()
	_ = cur
	if cur == 0 {
		return n
	}
	for cur != 0 {
		var next pkg_rhino.Node = cur.GetNext()
		_ = next
		if next == 0 {
			break
		}
		if cur.IsTemplateLitString() && next.IsTemplateLitSub() {
			var expr pkg_rhino.Node = next.GetFirstChild()
			_ = expr
			if (expr != 0) && (((expr.IsStringLit() || expr.IsNumber()) || expr.IsTrue()) || expr.IsFalse()) {
				var strValue string = this_.GetSideEffectFreeStringValue(expr)
				_ = strValue
				if (strValue != j2gort.NullString) && this_.priv_PeepholeFoldConstants_isSimpleTemplateStringValue_4205(strValue) {
					var nextNext pkg_rhino.Node = next.GetNext()
					_ = nextNext
					if (nextNext != 0) && nextNext.IsTemplateLitString() {
						var curCooked string = cur.GetCookedString()
						_ = curCooked
						var nextNextCooked string = nextNext.GetCookedString()
						_ = nextNextCooked
						if ((j2gort.CheckNotNullString(strValue) == "") && this_.priv_PeepholeFoldConstants_endsWithUnescapedDollar_4207(cur.GetRawString())) && javalang.StringStartsWith(j2gort.CheckNotNullString(nextNext.GetRawString()), "{") {
							cur = next
							continue
						}
						var combinedCooked string = (j2gort.StrValString(curCooked) + j2gort.StrValString(strValue) + j2gort.StrValString(nextNextCooked))
						_ = combinedCooked
						var combinedRaw string = (j2gort.StrValString(cur.GetRawString()) + j2gort.StrValString(strValue) + j2gort.StrValString(nextNext.GetRawString()))
						_ = combinedRaw
						if this_.priv_PeepholeFoldConstants_isRestrictedOctalEscape_4209(cur.GetRawString(), strValue) || this_.priv_PeepholeFoldConstants_isRestrictedOctalEscape_4209((j2gort.StrValString(cur.GetRawString())+j2gort.StrValString(strValue)), nextNext.GetRawString()) {
							cur = next
							continue
						}
						var newStringNode pkg_rhino.Node = pkg_rhino.IR_TemplateLiteralString(combinedCooked, combinedRaw)
						_ = newStringNode
						newStringNode.SrcrefTree(cur)
						cur.ReplaceWith(newStringNode)
						next.Detach()
						nextNext.Detach()
						cur = newStringNode
						changed = true
						continue
					}
				}
			}
		}
		cur = cur.GetNext()
	}
	if changed {
		this_.ReportChangeToEnclosingScope(n)
	}
	return n
}

// Java method: PeepholeFoldConstants.isRestrictedOctalEscape(java.lang.String prefixRaw, java.lang.String suffixRaw)
func (this_ *PeepholeFoldConstants) priv_PeepholeFoldConstants_isRestrictedOctalEscape_4209(prefixRaw string, suffixRaw string) bool {
	_ = prefixRaw
	_ = suffixRaw
	if j2gort.CheckNotNullString(suffixRaw) == "" {
		return false
	}
	var firstChar uint16 = j2gort.JStrCharAt(suffixRaw, int32(0))
	_ = firstChar
	if (int32(firstChar) < int32(uint16(48))) || (int32(firstChar) > int32(uint16(57))) {
		return false
	}
	if javalang.StringEndsWith(j2gort.CheckNotNullString(prefixRaw), "0") {
		var beforeZero string = j2gort.JStrSub(prefixRaw, int32(0), (j2gort.JStrLen(prefixRaw) - int32(1)))
		_ = beforeZero
		return this_.priv_PeepholeFoldConstants_endsWithUnescapedBackslash_4206(beforeZero)
	}
	return false
}

// Java method: PeepholeFoldConstants.tryFoldAdd(com.google.javascript.rhino.Node node, com.google.javascript.rhino.Node left, com.google.javascript.rhino.Node right)
func (this_ *PeepholeFoldConstants) priv_PeepholeFoldConstants_tryFoldAdd_4210(node pkg_rhino.Node, left pkg_rhino.Node, right pkg_rhino.Node) pkg_rhino.Node {
	_ = node
	_ = left
	_ = right
	if !(node.IsAdd()) {
		j2gort.CheckArgument(false)
	}
	if NodeUtil_MayBeString_2(node, this_.shouldUseTypes) {
		if left.IsTemplateLit() && right.IsTemplateLit() {
			return this_.priv_PeepholeFoldConstants_foldAddTemplateLiterals_4204(node, left, right)
		} else if left.IsTemplateLit() {
			var rightStr string = this_.GetSideEffectFreeStringValue(right)
			_ = rightStr
			if (rightStr != j2gort.NullString) && this_.priv_PeepholeFoldConstants_isSimpleTemplateStringValue_4205(rightStr) {
				return this_.priv_PeepholeFoldConstants_foldTemplateLiteralAndConstant_4211(node, left, rightStr)
			}
		} else if right.IsTemplateLit() {
			var leftStr string = this_.GetSideEffectFreeStringValue(left)
			_ = leftStr
			if (leftStr != j2gort.NullString) && this_.priv_PeepholeFoldConstants_isSimpleTemplateStringValue_4205(leftStr) {
				return this_.priv_PeepholeFoldConstants_foldConstantAndTemplateLiteral_4212(node, leftStr, right)
			}
		}
		if NodeUtil_IsLiteralValue(left, false) && NodeUtil_IsLiteralValue(right, false) {
			return this_.priv_PeepholeFoldConstants_tryFoldAddConstantString_4196(node, left, right)
		} else if (left.IsStringLit() && (j2gort.CheckNotNullString(left.GetString()) == "")) && this_.priv_PeepholeFoldConstants_isStringTyped_4214(right) {
			return this_.priv_PeepholeFoldConstants_replace_4213(node, right.CloneTree_1(true))
		} else if (right.IsStringLit() && (j2gort.CheckNotNullString(right.GetString()) == "")) && this_.priv_PeepholeFoldConstants_isStringTyped_4214(left) {
			return this_.priv_PeepholeFoldConstants_replace_4213(node, left.CloneTree_1(true))
		} else {
			return this_.priv_PeepholeFoldConstants_tryFoldAdjacentLiteralLeaves_4195(node, left, right)
		}
	} else {
		return this_.priv_PeepholeFoldConstants_tryFoldArithmeticOp_4197(node, left, right)
	}
}

// Java method: PeepholeFoldConstants.foldTemplateLiteralAndConstant(com.google.javascript.rhino.Node oldNode, com.google.javascript.rhino.Node left, java.lang.String rightStr)
func (this_ *PeepholeFoldConstants) priv_PeepholeFoldConstants_foldTemplateLiteralAndConstant_4211(oldNode pkg_rhino.Node, left pkg_rhino.Node, rightStr string) pkg_rhino.Node {
	_ = oldNode
	_ = left
	_ = rightStr
	var lastLeft pkg_rhino.Node = left.GetLastChild()
	_ = lastLeft
	if !(lastLeft.IsTemplateLitString()) {
		j2gort.CheckState(false)
	}
	if this_.priv_PeepholeFoldConstants_isRestrictedOctalEscape_4209(lastLeft.GetRawString(), rightStr) {
		return oldNode
	}
	var mergeNode pkg_rhino.Node = pkg_rhino.IR_TemplateLiteralString((j2gort.StrValString(lastLeft.GetCookedString()) + j2gort.StrValString(rightStr)), (j2gort.StrValString(lastLeft.GetRawString()) + j2gort.StrValString(rightStr)))
	_ = mergeNode
	mergeNode.SrcrefTreeIfMissing(lastLeft)
	lastLeft.ReplaceWith(mergeNode)
	left.Detach()
	return this_.priv_PeepholeFoldConstants_replace_4213(oldNode, left)
}

// Java method: PeepholeFoldConstants.foldConstantAndTemplateLiteral(com.google.javascript.rhino.Node oldNode, java.lang.String leftStr, com.google.javascript.rhino.Node right)
func (this_ *PeepholeFoldConstants) priv_PeepholeFoldConstants_foldConstantAndTemplateLiteral_4212(oldNode pkg_rhino.Node, leftStr string, right pkg_rhino.Node) pkg_rhino.Node {
	_ = oldNode
	_ = leftStr
	_ = right
	var firstRight pkg_rhino.Node = right.GetFirstChild()
	_ = firstRight
	if !(firstRight.IsTemplateLitString()) {
		j2gort.CheckState(false)
	}
	if this_.priv_PeepholeFoldConstants_isRestrictedOctalEscape_4209(leftStr, firstRight.GetRawString()) {
		return oldNode
	}
	var mergeNode pkg_rhino.Node = pkg_rhino.IR_TemplateLiteralString((j2gort.StrValString(leftStr) + j2gort.StrValString(firstRight.GetCookedString())), (j2gort.StrValString(leftStr) + j2gort.StrValString(firstRight.GetRawString())))
	_ = mergeNode
	mergeNode.SrcrefTreeIfMissing(firstRight)
	firstRight.ReplaceWith(mergeNode)
	right.Detach()
	return this_.priv_PeepholeFoldConstants_replace_4213(oldNode, right)
}

// Java method: PeepholeFoldConstants.replace(com.google.javascript.rhino.Node oldNode, com.google.javascript.rhino.Node newNode)
func (this_ *PeepholeFoldConstants) priv_PeepholeFoldConstants_replace_4213(oldNode pkg_rhino.Node, newNode pkg_rhino.Node) pkg_rhino.Node {
	_ = oldNode
	_ = newNode
	oldNode.ReplaceWith(newNode)
	this_.ReportChangeToEnclosingScope(newNode)
	return newNode
}

// Java method: PeepholeFoldConstants.isStringTyped(com.google.javascript.rhino.Node n)
func (this_ *PeepholeFoldConstants) priv_PeepholeFoldConstants_isStringTyped_4214(n pkg_rhino.Node) bool {
	_ = n
	if NodeUtil_IsStringResult(n) {
		return true
	}
	if this_.shouldUseTypes {
		return ((n.GetColor() != nil) && n.GetColor().Equals(j2gort.PtrToAny(pkg_colors.StaticGet_StandardColors_STRING())))
	}
	return false
}

// Java method: PeepholeFoldConstants.tryFoldShift(com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node left, com.google.javascript.rhino.Node right)
func (this_ *PeepholeFoldConstants) priv_PeepholeFoldConstants_tryFoldShift_4215(n pkg_rhino.Node, left pkg_rhino.Node, right pkg_rhino.Node) pkg_rhino.Node {
	_ = n
	_ = left
	_ = right
	var leftVal *j2gort.BoxDouble = this_.GetSideEffectFreeNumberValue(left)
	_ = leftVal
	var rightVal *j2gort.BoxDouble = this_.GetSideEffectFreeNumberValue(right)
	_ = rightVal
	if !(this_.priv_PeepholeFoldConstants_isReasonableDoubleValue_4189(leftVal)) || !(this_.priv_PeepholeFoldConstants_isReasonableDoubleValue_4189(rightVal)) {
		return n
	}
	if !(pkg_base.JSCompDoubles_IsMathematicalInteger(j2gort.UnboxDouble(leftVal))) {
		this_.Report(StaticGet_PeepholeFoldConstants_FRACTIONAL_BITWISE_OPERAND(), left)
		return n
	}
	if !(pkg_base.JSCompDoubles_IsMathematicalInteger(j2gort.UnboxDouble(rightVal))) {
		this_.Report(StaticGet_PeepholeFoldConstants_FRACTIONAL_BITWISE_OPERAND(), right)
		return n
	}
	if !((float64(int32(0)) <= j2gort.UnboxDouble(rightVal)) && (j2gort.UnboxDouble(rightVal) < float64(int32(32)))) {
		return n
	}
	var rvalInt int32 = rightVal.IntValue()
	_ = rvalInt
	var bits int32 = pkg_base.JSCompDoubles_EcmascriptToInt32(j2gort.UnboxDouble(leftVal))
	_ = bits
	var result float64 = float64(0)
	_ = result
	switch n.TokenSwitchOrd_() {
	case int32(10) /* Token.LSH */ :
		result = float64(javalang.Shl32(bits, int(rvalInt)))
	case int32(11) /* Token.RSH */ :
		result = float64(javalang.Shr32(bits, int(rvalInt)))
	case int32(12) /* Token.URSH */ :
		result = float64((int64(4294967295) & int64((javalang.Ushr32(bits, int(rvalInt))))))
	default:
		panic(j2gort.AsIface[javalang.IAssertionError](j2gort.ConstructShimThrowable("java.lang.AssertionError", "Unknown shift operator: "+j2gort.StrVal(n.GetToken()))))
	}
	var newNumber pkg_rhino.Node = NodeUtil_NumberNode(result, n)
	_ = newNumber
	this_.ReportChangeToEnclosingScope(n)
	n.ReplaceWith(newNumber)
	return newNumber
}

// Java method: PeepholeFoldConstants.tryFoldComparison(com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node left, com.google.javascript.rhino.Node right)
func (this_ *PeepholeFoldConstants) priv_PeepholeFoldConstants_tryFoldComparison_4216(n pkg_rhino.Node, left pkg_rhino.Node, right pkg_rhino.Node) pkg_rhino.Node {
	_ = n
	_ = left
	_ = right
	var result pkg_base.Tri = PeepholeFoldConstants_EvaluateComparison(this_, n.GetToken(), left, right)
	_ = result
	if result == pkg_base.StaticGet_Tri_UNKNOWN() {
		return n
	}
	var newNode pkg_rhino.Node = NodeUtil_BooleanNode(result.ToBoolean(true))
	_ = newNode
	this_.ReportChangeToEnclosingScope(n)
	n.ReplaceWith(newNode)
	this_.MarkFunctionsDeleted(n)
	return newNode
}

// Java static method: PeepholeFoldConstants.tryAbstractRelationalComparison(com.google.javascript.jscomp.AbstractPeepholeOptimization peepholeOptimization, com.google.javascript.rhino.Node left, com.google.javascript.rhino.Node right, boolean willNegate)
func PeepholeFoldConstants_TryAbstractRelationalComparison(peepholeOptimization AbstractPeepholeOptimization, left pkg_rhino.Node, right pkg_rhino.Node, willNegate bool) pkg_base.Tri {
	clinit_PeepholeFoldConstants()
	_ = peepholeOptimization
	_ = left
	_ = right
	_ = willNegate
	var leftValueType NodeUtil_ValueType = NodeUtil_GetKnownValueType(left)
	_ = leftValueType
	var rightValueType NodeUtil_ValueType = NodeUtil_GetKnownValueType(right)
	_ = rightValueType
	if (leftValueType == StaticGet_NodeUtil_ValueType_STRING()) && (rightValueType == StaticGet_NodeUtil_ValueType_STRING()) {
		var lvStr string = peepholeOptimization.GetAbstractPeepholeOptimizationBase().GetSideEffectFreeStringValue(left)
		_ = lvStr
		var rvStr string = peepholeOptimization.GetAbstractPeepholeOptimizationBase().GetSideEffectFreeStringValue(right)
		_ = rvStr
		if (lvStr != j2gort.NullString) && (rvStr != j2gort.NullString) {
			if (j2gort.JStrIndexOfCp(lvStr, int32(uint16(11))) != int32(-1)) || (j2gort.JStrIndexOfCp(rvStr, int32(uint16(11))) != int32(-1)) {
				return pkg_base.StaticGet_Tri_UNKNOWN()
			} else {
				return pkg_base.Tri_ForBoolean((j2gort.JStrCompareTo(lvStr, rvStr) < int32(0)))
			}
		} else if (((left.IsTypeOf() && right.IsTypeOf()) && left.GetFirstChild().IsName()) && right.GetFirstChild().IsName()) && func() bool {
			sr_5 := left.GetFirstChild().GetString()
			se_6 := right.GetFirstChild().GetString()
			return j2gort.CheckNotNullString(sr_5) == se_6
		}() {
			return pkg_base.StaticGet_Tri_FALSE()
		}
	}
	if (leftValueType == StaticGet_NodeUtil_ValueType_BIGINT()) && (rightValueType == StaticGet_NodeUtil_ValueType_STRING()) {
		var rvStr_2 string = peepholeOptimization.GetAbstractPeepholeOptimizationBase().GetSideEffectFreeStringValue(right)
		_ = rvStr_2
		if rvStr_2 != j2gort.NullString {
			var rvBig *j2gort.Ext_java_math_BigInteger = NodeUtil_GetStringBigIntValue(rvStr_2)
			_ = rvBig
			if rvBig == nil {
				return pkg_base.Tri_ForBoolean(willNegate)
			}
			var lvBig *j2gort.Ext_java_math_BigInteger = peepholeOptimization.GetAbstractPeepholeOptimizationBase().GetSideEffectFreeBigIntValue(left)
			_ = lvBig
			if lvBig != nil {
				return pkg_base.Tri_ForBoolean((lvBig.CompareTo(j2gort.PtrToAny(rvBig)) < int32(0)))
			}
		}
		return pkg_base.StaticGet_Tri_UNKNOWN()
	}
	if (leftValueType == StaticGet_NodeUtil_ValueType_STRING()) && (rightValueType == StaticGet_NodeUtil_ValueType_BIGINT()) {
		var lvStr_2 string = peepholeOptimization.GetAbstractPeepholeOptimizationBase().GetSideEffectFreeStringValue(left)
		_ = lvStr_2
		if lvStr_2 != j2gort.NullString {
			var lvBig_2 *j2gort.Ext_java_math_BigInteger = NodeUtil_GetStringBigIntValue(lvStr_2)
			_ = lvBig_2
			if lvBig_2 == nil {
				return pkg_base.Tri_ForBoolean(willNegate)
			}
			var rvBig_2 *j2gort.Ext_java_math_BigInteger = peepholeOptimization.GetAbstractPeepholeOptimizationBase().GetSideEffectFreeBigIntValue(right)
			_ = rvBig_2
			if rvBig_2 != nil {
				return pkg_base.Tri_ForBoolean((lvBig_2.CompareTo(j2gort.PtrToAny(rvBig_2)) < int32(0)))
			}
		}
		return pkg_base.StaticGet_Tri_UNKNOWN()
	}
	base_peepholeOptimization_ := peepholeOptimization.GetAbstractPeepholeOptimizationBase()
	var lvBig_3 *j2gort.Ext_java_math_BigInteger = base_peepholeOptimization_.GetSideEffectFreeBigIntValue(left)
	_ = lvBig_3
	var rvBig_3 *j2gort.Ext_java_math_BigInteger = base_peepholeOptimization_.GetSideEffectFreeBigIntValue(right)
	_ = rvBig_3
	if (lvBig_3 != nil) && (rvBig_3 != nil) {
		return pkg_base.Tri_ForBoolean((lvBig_3.CompareTo(j2gort.PtrToAny(rvBig_3)) < int32(0)))
	}
	var lvNum *j2gort.BoxDouble = base_peepholeOptimization_.GetSideEffectFreeNumberValue(left)
	_ = lvNum
	var rvNum *j2gort.BoxDouble = base_peepholeOptimization_.GetSideEffectFreeNumberValue(right)
	_ = rvNum
	if (lvNum != nil) && (rvNum != nil) {
		if math.IsNaN(j2gort.UnboxDouble(lvNum)) || math.IsNaN(j2gort.UnboxDouble(rvNum)) {
			return pkg_base.Tri_ForBoolean(willNegate)
		} else {
			return pkg_base.Tri_ForBoolean((lvNum.DoubleValue() < rvNum.DoubleValue()))
		}
	}
	if (lvBig_3 != nil) && (rvNum != nil) {
		return PeepholeFoldConstants_BigintLessThanDouble(lvBig_3, j2gort.UnboxDouble(rvNum), pkg_base.StaticGet_Tri_FALSE(), willNegate)
	}
	if (lvNum != nil) && (rvBig_3 != nil) {
		return PeepholeFoldConstants_BigintLessThanDouble(rvBig_3, j2gort.UnboxDouble(lvNum), pkg_base.StaticGet_Tri_TRUE(), willNegate)
	}
	if (!(willNegate) && left.IsName()) && right.IsName() {
		{
			var iifeCond_1 bool
			{

				sr_7 := left.GetString()
				se_8 := right.GetString()
				iifeCond_1 = j2gort.CheckNotNullString(sr_7) == se_8
			}
			if iifeCond_1 {
				return pkg_base.StaticGet_Tri_FALSE()
			}
		}
	}
	return pkg_base.StaticGet_Tri_UNKNOWN()
}

// Java static method: PeepholeFoldConstants.bigintLessThanDouble(java.math.BigInteger bigint, double number, com.google.javascript.jscomp.base.Tri invert, boolean willNegate)
func PeepholeFoldConstants_BigintLessThanDouble(bigint *j2gort.Ext_java_math_BigInteger, number float64, invert pkg_base.Tri, willNegate bool) pkg_base.Tri {
	clinit_PeepholeFoldConstants()
	_ = bigint
	_ = number
	_ = invert
	_ = willNegate
	if math.IsNaN(number) {
		return pkg_base.Tri_ForBoolean(willNegate)
	} else if number == math.Inf(1) {
		return pkg_base.StaticGet_Tri_TRUE().Xor(invert)
	} else if number == math.Inf(-1) {
		return pkg_base.StaticGet_Tri_FALSE().Xor(invert)
	} else if !(pkg_base.JSCompDoubles_IsAtLeastIntegerPrecision(number)) {
		return pkg_base.StaticGet_Tri_UNKNOWN()
	}
	var numberAsBigInt *j2gort.Ext_java_math_BigInteger = j2gort.Ext_java_math_BigInteger_ValueOf(javalang.D2L(number))
	_ = numberAsBigInt
	var negativeMeansBigintSmaller int32 = bigint.CompareTo(j2gort.PtrToAny(numberAsBigInt))
	_ = negativeMeansBigintSmaller
	if negativeMeansBigintSmaller < int32(0) {
		return pkg_base.StaticGet_Tri_TRUE().Xor(invert)
	} else if negativeMeansBigintSmaller > int32(0) {
		return pkg_base.StaticGet_Tri_FALSE().Xor(invert)
	} else if pkg_base.JSCompDoubles_IsExactInt64(number) {
		return pkg_base.StaticGet_Tri_FALSE()
	} else {
		return pkg_base.Tri_ForBoolean(pkg_base.JSCompDoubles_IsPositive(number)).Xor(invert)
	}
}

// Java static method: PeepholeFoldConstants.tryAbstractEqualityComparison(com.google.javascript.jscomp.AbstractPeepholeOptimization peepholeOptimization, com.google.javascript.rhino.Node left, com.google.javascript.rhino.Node right)
func PeepholeFoldConstants_TryAbstractEqualityComparison(peepholeOptimization AbstractPeepholeOptimization, left pkg_rhino.Node, right pkg_rhino.Node) pkg_base.Tri {
	clinit_PeepholeFoldConstants()
	_ = peepholeOptimization
	_ = left
	_ = right
	var leftValueType NodeUtil_ValueType = NodeUtil_GetKnownValueType(left)
	_ = leftValueType
	var rightValueType NodeUtil_ValueType = NodeUtil_GetKnownValueType(right)
	_ = rightValueType
	if (leftValueType != StaticGet_NodeUtil_ValueType_UNDETERMINED()) && (rightValueType != StaticGet_NodeUtil_ValueType_UNDETERMINED()) {
		if leftValueType == rightValueType {
			return PeepholeFoldConstants_TryStrictEqualityComparison(peepholeOptimization, left, right)
		}
		if ((leftValueType == StaticGet_NodeUtil_ValueType_NULL()) && (rightValueType == StaticGet_NodeUtil_ValueType_VOID())) || ((leftValueType == StaticGet_NodeUtil_ValueType_VOID()) && (rightValueType == StaticGet_NodeUtil_ValueType_NULL())) {
			return pkg_base.StaticGet_Tri_TRUE()
		}
		if ((leftValueType == StaticGet_NodeUtil_ValueType_NUMBER()) && (rightValueType == StaticGet_NodeUtil_ValueType_STRING())) || (rightValueType == StaticGet_NodeUtil_ValueType_BOOLEAN()) {
			var rv *j2gort.BoxDouble = peepholeOptimization.GetAbstractPeepholeOptimizationBase().GetSideEffectFreeNumberValue(right)
			_ = rv
			if rv == nil {
				return pkg_base.StaticGet_Tri_UNKNOWN()
			} else {
				return PeepholeFoldConstants_TryAbstractEqualityComparison(peepholeOptimization, left, NodeUtil_NumberNode(j2gort.UnboxDouble(rv), right))
			}
		}
		if ((leftValueType == StaticGet_NodeUtil_ValueType_STRING()) && (rightValueType == StaticGet_NodeUtil_ValueType_NUMBER())) || (leftValueType == StaticGet_NodeUtil_ValueType_BOOLEAN()) {
			var lv *j2gort.BoxDouble = peepholeOptimization.GetAbstractPeepholeOptimizationBase().GetSideEffectFreeNumberValue(left)
			_ = lv
			if lv == nil {
				return pkg_base.StaticGet_Tri_UNKNOWN()
			} else {
				return PeepholeFoldConstants_TryAbstractEqualityComparison(peepholeOptimization, NodeUtil_NumberNode(j2gort.UnboxDouble(lv), left), right)
			}
		}
		if (leftValueType == StaticGet_NodeUtil_ValueType_BIGINT()) && (rightValueType == StaticGet_NodeUtil_ValueType_STRING()) {
			var rvStr string = peepholeOptimization.GetAbstractPeepholeOptimizationBase().GetSideEffectFreeStringValue(right)
			_ = rvStr
			if rvStr != j2gort.NullString {
				var rvBig *j2gort.Ext_java_math_BigInteger = NodeUtil_GetStringBigIntValue(rvStr)
				_ = rvBig
				if rvBig == nil {
					return pkg_base.StaticGet_Tri_FALSE()
				}
				var lvBig *j2gort.Ext_java_math_BigInteger = peepholeOptimization.GetAbstractPeepholeOptimizationBase().GetSideEffectFreeBigIntValue(left)
				_ = lvBig
				if lvBig != nil {
					return pkg_base.Tri_ForBoolean(lvBig.Equals(j2gort.PtrToAny(rvBig)))
				} else {
					return pkg_base.StaticGet_Tri_UNKNOWN()
				}
			}
			return pkg_base.StaticGet_Tri_UNKNOWN()
		}
		if (leftValueType == StaticGet_NodeUtil_ValueType_STRING()) && (rightValueType == StaticGet_NodeUtil_ValueType_BIGINT()) {
			var lvStr string = peepholeOptimization.GetAbstractPeepholeOptimizationBase().GetSideEffectFreeStringValue(left)
			_ = lvStr
			if lvStr != j2gort.NullString {
				var lvBig_2 *j2gort.Ext_java_math_BigInteger = NodeUtil_GetStringBigIntValue(lvStr)
				_ = lvBig_2
				if lvBig_2 == nil {
					return pkg_base.StaticGet_Tri_FALSE()
				}
				var rvBig_2 *j2gort.Ext_java_math_BigInteger = peepholeOptimization.GetAbstractPeepholeOptimizationBase().GetSideEffectFreeBigIntValue(right)
				_ = rvBig_2
				if rvBig_2 != nil {
					return pkg_base.Tri_ForBoolean(lvBig_2.Equals(j2gort.PtrToAny(rvBig_2)))
				} else {
					return pkg_base.StaticGet_Tri_UNKNOWN()
				}
			}
			return pkg_base.StaticGet_Tri_UNKNOWN()
		}
		if (leftValueType == StaticGet_NodeUtil_ValueType_BIGINT()) || (rightValueType == StaticGet_NodeUtil_ValueType_BIGINT()) {
			var lv_2 *j2gort.Ext_java_math_BigInteger = peepholeOptimization.GetAbstractPeepholeOptimizationBase().GetSideEffectFreeBigIntValue(left)
			_ = lv_2
			var rv_2 *j2gort.Ext_java_math_BigInteger = peepholeOptimization.GetAbstractPeepholeOptimizationBase().GetSideEffectFreeBigIntValue(right)
			_ = rv_2
			if (lv_2 != nil) && (rv_2 != nil) {
				return pkg_base.Tri_ForBoolean(lv_2.Equals(j2gort.PtrToAny(rv_2)))
			}
			if (((leftValueType == StaticGet_NodeUtil_ValueType_NULL()) || (leftValueType == StaticGet_NodeUtil_ValueType_VOID())) || (rightValueType == StaticGet_NodeUtil_ValueType_NULL())) || (rightValueType == StaticGet_NodeUtil_ValueType_VOID()) {
				return pkg_base.StaticGet_Tri_FALSE()
			}
			if ((leftValueType == StaticGet_NodeUtil_ValueType_STRING()) && (peepholeOptimization.GetAbstractPeepholeOptimizationBase().GetSideEffectFreeStringValue(left) != j2gort.NullString)) && (lv_2 == nil) {
				return pkg_base.StaticGet_Tri_FALSE()
			}
			if ((rightValueType == StaticGet_NodeUtil_ValueType_STRING()) && (peepholeOptimization.GetAbstractPeepholeOptimizationBase().GetSideEffectFreeStringValue(right) != j2gort.NullString)) && (rv_2 == nil) {
				return pkg_base.StaticGet_Tri_FALSE()
			}
			return pkg_base.StaticGet_Tri_UNKNOWN()
		}
		if ((leftValueType == StaticGet_NodeUtil_ValueType_STRING()) || (leftValueType == StaticGet_NodeUtil_ValueType_NUMBER())) && (rightValueType == StaticGet_NodeUtil_ValueType_OBJECT()) {
			return pkg_base.StaticGet_Tri_UNKNOWN()
		}
		if (leftValueType == StaticGet_NodeUtil_ValueType_OBJECT()) && ((rightValueType == StaticGet_NodeUtil_ValueType_STRING()) || (rightValueType == StaticGet_NodeUtil_ValueType_NUMBER())) {
			return pkg_base.StaticGet_Tri_UNKNOWN()
		}
		return pkg_base.StaticGet_Tri_FALSE()
	}
	return pkg_base.StaticGet_Tri_UNKNOWN()
}

// Java static method: PeepholeFoldConstants.tryStrictEqualityComparison(com.google.javascript.jscomp.AbstractPeepholeOptimization peepholeOptimization, com.google.javascript.rhino.Node left, com.google.javascript.rhino.Node right)
func PeepholeFoldConstants_TryStrictEqualityComparison(peepholeOptimization AbstractPeepholeOptimization, left pkg_rhino.Node, right pkg_rhino.Node) pkg_base.Tri {
	clinit_PeepholeFoldConstants()
	_ = peepholeOptimization
	_ = left
	_ = right
	var leftValueType NodeUtil_ValueType = NodeUtil_GetKnownValueType(left)
	_ = leftValueType
	var rightValueType NodeUtil_ValueType = NodeUtil_GetKnownValueType(right)
	_ = rightValueType
	if (leftValueType != StaticGet_NodeUtil_ValueType_UNDETERMINED()) && (rightValueType != StaticGet_NodeUtil_ValueType_UNDETERMINED()) {
		if leftValueType != rightValueType {
			return pkg_base.StaticGet_Tri_FALSE()
		}
		switch (leftValueType).Ordinal() {
		case int32(2) /* ValueType.VOID */, int32(1) /* ValueType.NULL */ :
			return pkg_base.StaticGet_Tri_TRUE()
		case int32(3) /* ValueType.NUMBER */ :
			if NodeUtil_IsNaN(left) {
				return pkg_base.StaticGet_Tri_FALSE()
			}
			if NodeUtil_IsNaN(right) {
				return pkg_base.StaticGet_Tri_FALSE()
			}
			var lv *j2gort.BoxDouble = peepholeOptimization.GetAbstractPeepholeOptimizationBase().GetSideEffectFreeNumberValue(left)
			_ = lv
			var rv *j2gort.BoxDouble = peepholeOptimization.GetAbstractPeepholeOptimizationBase().GetSideEffectFreeNumberValue(right)
			_ = rv
			if (lv != nil) && (rv != nil) {
				return pkg_base.Tri_ForBoolean((lv.DoubleValue() == rv.DoubleValue()))
			}
		case int32(5) /* ValueType.STRING */ :
			var lv_2 string = peepholeOptimization.GetAbstractPeepholeOptimizationBase().GetSideEffectFreeStringValue(left)
			_ = lv_2
			var rv_2 string = peepholeOptimization.GetAbstractPeepholeOptimizationBase().GetSideEffectFreeStringValue(right)
			_ = rv_2
			if (lv_2 != j2gort.NullString) && (rv_2 != j2gort.NullString) {
				if (j2gort.JStrIndexOfCp(lv_2, int32(uint16(11))) != int32(-1)) || (j2gort.JStrIndexOfCp(rv_2, int32(uint16(11))) != int32(-1)) {
					return pkg_base.StaticGet_Tri_UNKNOWN()
				} else {
					if j2gort.CheckNotNullString(lv_2) == rv_2 {
						return pkg_base.StaticGet_Tri_TRUE()
					} else {
						return pkg_base.StaticGet_Tri_FALSE()
					}
				}
			} else if (((left.IsTypeOf() && right.IsTypeOf()) && left.GetFirstChild().IsName()) && right.GetFirstChild().IsName()) && func() bool {
				sr_9 := left.GetFirstChild().GetString()
				se_10 := right.GetFirstChild().GetString()
				return j2gort.CheckNotNullString(sr_9) == se_10
			}() {
				return pkg_base.StaticGet_Tri_TRUE()
			}
		case int32(6) /* ValueType.BOOLEAN */ :
			var lv_3 pkg_base.Tri = peepholeOptimization.GetAbstractPeepholeOptimizationBase().GetSideEffectFreeBooleanValue(left)
			_ = lv_3
			var rv_3 pkg_base.Tri = peepholeOptimization.GetAbstractPeepholeOptimizationBase().GetSideEffectFreeBooleanValue(right)
			_ = rv_3
			return lv_3.And(rv_3).Or(lv_3.Not().And(rv_3.Not()))
		case int32(4) /* ValueType.BIGINT */ :
			var lv_4 *j2gort.Ext_java_math_BigInteger = peepholeOptimization.GetAbstractPeepholeOptimizationBase().GetSideEffectFreeBigIntValue(left)
			_ = lv_4
			var rv_4 *j2gort.Ext_java_math_BigInteger = peepholeOptimization.GetAbstractPeepholeOptimizationBase().GetSideEffectFreeBigIntValue(right)
			_ = rv_4
			if (lv_4 != nil) && (rv_4 != nil) {
				return pkg_base.Tri_ForBoolean(lv_4.Equals(j2gort.PtrToAny(rv_4)))
			}
		default:
			return pkg_base.StaticGet_Tri_UNKNOWN()
		}
	}
	if NodeUtil_IsNaN(left) || NodeUtil_IsNaN(right) {
		return pkg_base.StaticGet_Tri_FALSE()
	}
	return pkg_base.StaticGet_Tri_UNKNOWN()
}

// Java static method: PeepholeFoldConstants.evaluateComparison(com.google.javascript.jscomp.AbstractPeepholeOptimization peepholeOptimization, com.google.javascript.rhino.Token op, com.google.javascript.rhino.Node left, com.google.javascript.rhino.Node right)
func PeepholeFoldConstants_EvaluateComparison(peepholeOptimization AbstractPeepholeOptimization, op pkg_rhino.Token, left pkg_rhino.Node, right pkg_rhino.Node) pkg_base.Tri {
	clinit_PeepholeFoldConstants()
	_ = peepholeOptimization
	_ = op
	_ = left
	_ = right
	base_peepholeOptimization_ := peepholeOptimization.GetAbstractPeepholeOptimizationBase()
	if base_peepholeOptimization_.MayHaveSideEffects(left) || base_peepholeOptimization_.MayHaveSideEffects(right) {
		return pkg_base.StaticGet_Tri_UNKNOWN()
	}
	switch (op).Ordinal() {
	case int32(4) /* Token.EQ */ :
		return PeepholeFoldConstants_TryAbstractEqualityComparison(peepholeOptimization, left, right)
	case int32(5) /* Token.NE */ :
		return PeepholeFoldConstants_TryAbstractEqualityComparison(peepholeOptimization, left, right).Not()
	case int32(40) /* Token.SHEQ */ :
		return PeepholeFoldConstants_TryStrictEqualityComparison(peepholeOptimization, left, right)
	case int32(41) /* Token.SHNE */ :
		return PeepholeFoldConstants_TryStrictEqualityComparison(peepholeOptimization, left, right).Not()
	case int32(6) /* Token.LT */ :
		return PeepholeFoldConstants_TryAbstractRelationalComparison(peepholeOptimization, left, right, false)
	case int32(8) /* Token.GT */ :
		return PeepholeFoldConstants_TryAbstractRelationalComparison(peepholeOptimization, right, left, false)
	case int32(7) /* Token.LE */ :
		return PeepholeFoldConstants_TryAbstractRelationalComparison(peepholeOptimization, right, left, true).Not()
	case int32(9) /* Token.GE */ :
		return PeepholeFoldConstants_TryAbstractRelationalComparison(peepholeOptimization, left, right, true).Not()
	default:
	}
	panic(j2gort.AsIface[javalang.IIllegalStateException](j2gort.ConstructShimThrowable("java.lang.IllegalStateException", "Unexpected operator for comparison")))
}

// Java method: PeepholeFoldConstants.tryFoldCtorCall(com.google.javascript.rhino.Node n)
func (this_ *PeepholeFoldConstants) priv_PeepholeFoldConstants_tryFoldCtorCall_4217(n pkg_rhino.Node) pkg_rhino.Node {
	_ = n
	if !(n.IsNew()) {
		j2gort.CheckArgument(false)
	}
	if PeepholeFoldConstants_InForcedStringContext(n) {
		return this_.priv_PeepholeFoldConstants_tryFoldInForcedStringContext_4219(n)
	}
	return n
}

// Java method: PeepholeFoldConstants.tryFoldUselessObjectDotDefinePropertiesCall(com.google.javascript.rhino.Node n)
func (this_ *PeepholeFoldConstants) priv_PeepholeFoldConstants_tryFoldUselessObjectDotDefinePropertiesCall_4218(n pkg_rhino.Node) pkg_rhino.Node {
	_ = n
	if !(n.IsCall() || n.IsOptChainCall()) {
		j2gort.CheckArgument(false)
	}
	if NodeUtil_IsObjectDefinePropertiesDefinition(n) {
		var srcObj pkg_rhino.Node = n.GetLastChild()
		_ = srcObj
		if srcObj.IsObjectLit() && !(srcObj.HasChildren()) {
			var parent pkg_rhino.Node = n.GetParent()
			_ = parent
			var destObj pkg_rhino.Node = n.GetSecondChild().Detach()
			_ = destObj
			n.ReplaceWith(destObj)
			this_.ReportChangeToEnclosingScope(parent)
		}
	}
	return n
}

// Java static method: PeepholeFoldConstants.inForcedStringContext(com.google.javascript.rhino.Node n)
func PeepholeFoldConstants_InForcedStringContext(n pkg_rhino.Node) bool {
	clinit_PeepholeFoldConstants()
	_ = n
	if n.GetParent().IsGetElem() && (n.GetParent().GetLastChild() == n) {
		return true
	}
	return n.GetParent().IsAdd()
}

// Java method: PeepholeFoldConstants.tryFoldInForcedStringContext(com.google.javascript.rhino.Node n)
func (this_ *PeepholeFoldConstants) priv_PeepholeFoldConstants_tryFoldInForcedStringContext_4219(n pkg_rhino.Node) pkg_rhino.Node {
	_ = n
	if !(n.IsNew()) {
		j2gort.CheckArgument(false)
	}
	var objectType pkg_rhino.Node = n.GetFirstChild()
	_ = objectType
	if !(objectType.IsName()) {
		return n
	}
	if j2gort.CheckNotNullString(objectType.GetString()) == "String" {
		var value pkg_rhino.Node = objectType.GetNext()
		_ = value
		var stringValue string = j2gort.NullString
		_ = stringValue
		if value == 0 {
			stringValue = ""
		} else {
			stringValue = this_.GetSideEffectFreeStringValue(value)
		}
		if stringValue == j2gort.NullString {
			return n
		}
		var parent pkg_rhino.Node = n.GetParent()
		_ = parent
		var newString pkg_rhino.Node = pkg_rhino.IR_String(stringValue)
		_ = newString
		n.ReplaceWith(newString)
		newString.SrcrefIfMissing(parent)
		this_.ReportChangeToEnclosingScope(parent)
		return newString
	}
	return n
}

// Java method: PeepholeFoldConstants.tryFoldGetElem(com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node left, com.google.javascript.rhino.Node right)
func (this_ *PeepholeFoldConstants) priv_PeepholeFoldConstants_tryFoldGetElem_4220(n pkg_rhino.Node, left pkg_rhino.Node, right pkg_rhino.Node) pkg_rhino.Node {
	_ = n
	_ = left
	_ = right
	if !(n.IsGetElem() || n.IsOptChainGetElem()) {
		j2gort.CheckArgument(false)
	}
	if left.IsObjectLit() {
		if right.IsStringLit() {
			return this_.priv_PeepholeFoldConstants_tryFoldObjectPropAccess_4226(n, left, right.GetString())
		}
	} else if left.IsArrayLit() {
		return this_.priv_PeepholeFoldConstants_tryFoldArrayAccess_4222(n, left, right)
	} else if left.IsStringLit() {
		return this_.priv_PeepholeFoldConstants_tryFoldStringArrayAccess_4225(n, left, right)
	}
	return n
}

// Java method: PeepholeFoldConstants.tryFoldGetProp(com.google.javascript.rhino.Node n)
func (this_ *PeepholeFoldConstants) priv_PeepholeFoldConstants_tryFoldGetProp_4221(n pkg_rhino.Node) pkg_rhino.Node {
	_ = n
	if !(n.IsGetProp() || n.IsOptChainGetProp()) {
		j2gort.CheckArgument(false)
	}
	var left pkg_rhino.Node = n.GetFirstChild()
	_ = left
	var name string = n.GetString()
	_ = name
	if left.IsObjectLit() {
		return this_.priv_PeepholeFoldConstants_tryFoldObjectPropAccess_4226(n, left, name)
	}
	if j2gort.CheckNotNullString(name) == "length" {
		if NodeUtil_IsLValue(n) {
			return n
		}
		var knownLength int32 = int32(-1)
		_ = knownLength
		switch left.TokenSwitchOrd_() {
		case int32(46) /* Token.ARRAYLIT */ :
			if this_.MayHaveSideEffects(left) {
				return n
			}
			{
				var child pkg_rhino.Node = left.GetFirstChild()
				_ = child
				for ; child != 0; child = child.GetNext() {
					if child.IsSpread() {
						return n
					}
				}
			}
			knownLength = left.GetChildCount()
		case int32(35) /* Token.STRINGLIT */ :
			knownLength = j2gort.JStrLen(left.GetString())
		default:
			return n
		}
		if !(knownLength != int32(-1)) {
			j2gort.CheckState(false)
		}
		var lengthNode pkg_rhino.Node = pkg_rhino.IR_Number(float64(knownLength))
		_ = lengthNode
		this_.ReportChangeToEnclosingScope(n)
		n.ReplaceWith(lengthNode)
		return lengthNode
	}
	return n
}

// Java static method: PeepholeFoldConstants.isFoldableElementIndexKey(com.google.javascript.rhino.Node right, int intIndex)
func PeepholeFoldConstants_IsFoldableElementIndexKey(right pkg_rhino.Node, intIndex int32) bool {
	clinit_PeepholeFoldConstants()
	_ = right
	_ = intIndex
	if right.IsStringLit() {
		return ((intIndex >= int32(0)) && func() bool {
			sr_11 := right.GetString()
			se_12 := strconv.FormatInt(int64(intIndex), 10)
			return j2gort.CheckNotNullString(sr_11) == se_12
		}())
	}
	return NodeUtil_IsNumericResult(right)
}

// Java method: PeepholeFoldConstants.tryFoldArrayAccess(com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node left, com.google.javascript.rhino.Node right)
func (this_ *PeepholeFoldConstants) priv_PeepholeFoldConstants_tryFoldArrayAccess_4222(n pkg_rhino.Node, left pkg_rhino.Node, right pkg_rhino.Node) pkg_rhino.Node {
	_ = n
	_ = left
	_ = right
	if NodeUtil_IsLValue(n) {
		return n
	}
	var index *j2gort.BoxDouble = this_.GetSideEffectFreeNumberValue(right)
	_ = index
	if !(this_.priv_PeepholeFoldConstants_isReasonableDoubleValue_4189(index)) {
		return n
	}
	var intIndex int32 = index.IntValue()
	_ = intIndex
	if !(PeepholeFoldConstants_IsFoldableElementIndexKey(right, intIndex)) {
		return n
	}
	if !(pkg_base.JSCompDoubles_IsExactInt32(j2gort.UnboxDouble(index))) {
		this_.Report(StaticGet_PeepholeFoldConstants_INVALID_GETELEM_INDEX_ERROR(), right)
		return n
	}
	var current pkg_rhino.Node = left.GetFirstChild()
	_ = current
	var elem pkg_rhino.Node = pkg_rhino.Node(0)
	_ = elem
	{
		var i int32 = int32(0)
		_ = i
		for ; current != 0; i = i + 1 {
			if current.IsSpread() {
				return n
			}
			if i != intIndex {
				if this_.MayHaveSideEffects(current) {
					return n
				}
			} else {
				elem = current
			}
			current = current.GetNext()
		}
	}
	if elem == 0 {
		elem = NodeUtil_NewUndefinedNode(left)
	} else if elem.IsEmpty() {
		elem = NodeUtil_NewUndefinedNode(elem)
	} else {
		elem.Detach()
	}
	n.ReplaceWith(elem)
	this_.ReportChangeToEnclosingScope(elem)
	return elem
}

// Java method: PeepholeFoldConstants.tryFoldSpread(com.google.javascript.rhino.Node spread)
func (this_ *PeepholeFoldConstants) priv_PeepholeFoldConstants_tryFoldSpread_4223(spread pkg_rhino.Node) pkg_rhino.Node {
	_ = spread
	if !(spread.IsSpread()) {
		j2gort.CheckState(false)
	}
	var parent pkg_rhino.Node = spread.GetParent()
	_ = parent
	var child pkg_rhino.Node = spread.GetOnlyChild()
	_ = child
	if child.IsArrayLit() && PeepholeFoldConstants_IsSafeToFlattenArrayLitSpread(child) {
		parent.AddChildrenAfter(child.RemoveChildren(), spread)
		spread.Detach()
		this_.ReportChangeToEnclosingScope(parent)
	}
	return parent
}

// Java method: PeepholeFoldConstants.tryFlattenArrayOrObjectLit(com.google.javascript.rhino.Node parentLit)
func (this_ *PeepholeFoldConstants) priv_PeepholeFoldConstants_tryFlattenArrayOrObjectLit_4224(parentLit pkg_rhino.Node) pkg_rhino.Node {
	_ = parentLit
	{
		var child pkg_rhino.Node = parentLit.GetFirstChild()
		_ = child
		for child != 0 {
			var spread pkg_rhino.Node = child
			_ = spread
			child = child.GetNext()
			if !(spread.IsSpread()) {
				continue
			}
			var innerLit pkg_rhino.Node = spread.GetOnlyChild()
			_ = innerLit
			{
				var iifeCond_1 bool
				{

					re_13 := parentLit.GetToken()
					ae_14 := innerLit.GetToken()
					if re_13 == 0 {
						j2gort.ThrowNPE()
					}
					iifeCond_1 = re_13 == ae_14
				}
				if !(iifeCond_1) {
					continue
				}
			}
			if parentLit.IsObjectLit() && !(PeepholeFoldConstants_IsSafeToFlattenObjectLitSpread(innerLit)) {
				continue
			}
			if parentLit.IsArrayLit() && !(PeepholeFoldConstants_IsSafeToFlattenArrayLitSpread(innerLit)) {
				continue
			}
			parentLit.AddChildrenAfter(innerLit.RemoveChildren(), spread)
			spread.Detach()
			this_.ReportChangeToEnclosingScope(parentLit)
		}
	}
	return parentLit
}

// Java method: PeepholeFoldConstants.tryFoldStringArrayAccess(com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node left, com.google.javascript.rhino.Node right)
func (this_ *PeepholeFoldConstants) priv_PeepholeFoldConstants_tryFoldStringArrayAccess_4225(n pkg_rhino.Node, left pkg_rhino.Node, right pkg_rhino.Node) pkg_rhino.Node {
	_ = n
	_ = left
	_ = right
	if NodeUtil_IsLValue(n) {
		return n
	}
	var index *j2gort.BoxDouble = this_.GetSideEffectFreeNumberValue(right)
	_ = index
	if !(this_.priv_PeepholeFoldConstants_isReasonableDoubleValue_4189(index)) {
		return n
	}
	var intIndex int32 = index.IntValue()
	_ = intIndex
	if !(PeepholeFoldConstants_IsFoldableElementIndexKey(right, intIndex)) {
		return n
	}
	if float64(intIndex) != j2gort.UnboxDouble(index) {
		this_.Report(StaticGet_PeepholeFoldConstants_INVALID_GETELEM_INDEX_ERROR(), right)
		return n
	}
	if !(left.IsStringLit()) {
		j2gort.CheckState(false)
	}
	var value string = left.GetString()
	_ = value
	if (intIndex < int32(0)) || (intIndex >= j2gort.JStrLen(value)) {
		var undefined pkg_rhino.Node = NodeUtil_NewUndefinedNode(left)
		_ = undefined
		n.ReplaceWith(undefined)
		this_.ReportChangeToEnclosingScope(undefined)
		return undefined
	}
	var c uint16 = uint16(int32(0))
	_ = c
	{
		var i int32 = int32(0)
		_ = i
		for ; i <= intIndex; i = i + 1 {
			c = j2gort.JStrCharAt(value, i)
			if (int32(c) < int32(32)) || (int32(c) > int32(127)) {
				return n
			}
		}
	}
	var elem pkg_rhino.Node = pkg_rhino.IR_String(javalang.StringFromChar(c))
	_ = elem
	n.ReplaceWith(elem)
	this_.ReportChangeToEnclosingScope(elem)
	return elem
}

// Java method: PeepholeFoldConstants.tryFoldObjectPropAccess(com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node left, java.lang.String name)
func (this_ *PeepholeFoldConstants) priv_PeepholeFoldConstants_tryFoldObjectPropAccess_4226(n pkg_rhino.Node, left pkg_rhino.Node, name string) pkg_rhino.Node {
	_ = n
	_ = left
	_ = name
	if !(NodeUtil_IsNormalOrOptChainGet(n)) {
		j2gort.CheckArgument(false)
	}
	if !(left.IsObjectLit()) {
		return n
	}
	if NodeUtil_IsLValue(n) {
		if !(!(NodeUtil_IsOptChainNode(n))) {
			j2gort.CheckState(false)
		}
		return n
	}
	var key pkg_rhino.Node = pkg_rhino.Node(0)
	_ = key
	var value pkg_rhino.Node = pkg_rhino.Node(0)
	_ = value
	{
		var c pkg_rhino.Node = left.GetFirstChild()
		_ = c
		for ; c != 0; c = c.GetNext() {
			switch c.TokenSwitchOrd_() {
			case int32(96) /* Token.SETTER_DEF */ :
				if j2gort.CheckNotNullString(c.GetString()) == name {
					if (key != 0) && !(PeepholeFoldConstants_IsGetterProp(key)) {
						key = c
						value = NodeUtil_NewUndefinedNode(c)
					}
				}
			case int32(128) /* Token.OBJECT_SPREAD */ :
				key = pkg_rhino.Node(0)
				value = pkg_rhino.Node(0)
			case int32(129) /* Token.COMPUTED_PROP */ :
				var prop pkg_rhino.Node = c.GetFirstChild()
				_ = prop
				if !(prop.IsStringLit()) {
					return n
				}
				if j2gort.CheckNotNullString(prop.GetString()) == name {
					if c.GetBooleanProp(pkg_rhino.StaticGet_Node_COMPUTED_PROP_SETTER()) {
						if (key != 0) && !(PeepholeFoldConstants_IsGetterProp(key)) {
							key = c
							value = NodeUtil_NewUndefinedNode(c)
						}
					} else {
						key = c
						value = c.GetSecondChild()
					}
				}
			case int32(95) /* Token.GETTER_DEF */, int32(100) /* Token.STRING_KEY */, int32(107) /* Token.MEMBER_FUNCTION_DEF */ :
				if j2gort.CheckNotNullString(c.GetString()) == name {
					key = c
					value = key.GetFirstChild()
				}
			default:
				panic(j2gort.AsIface[javalang.IIllegalStateException](j2gort.ConstructShimThrowable("java.lang.IllegalStateException")))
			}
		}
	}
	if value == 0 {
		return n
	}
	if NodeUtil_ReferencesSuper(value) {
		return n
	}
	if value.GetParent() != 0 {
		var tempValue pkg_rhino.Node = pkg_rhino.IR_NullNode()
		_ = tempValue
		value.ReplaceWith(tempValue)
		var hasSideEffectBesidesValue bool = this_.MayHaveSideEffects(left)
		_ = hasSideEffectBesidesValue
		tempValue.ReplaceWith(value)
		if hasSideEffectBesidesValue {
			return n
		}
	} else {
		if this_.MayHaveSideEffects(left) {
			return n
		}
	}
	var keyIsGetter bool = PeepholeFoldConstants_IsGetterProp(key)
	_ = keyIsGetter
	var nIsInvoked bool = NodeUtil_IsInvocationTarget(n)
	_ = nIsInvoked
	if keyIsGetter || nIsInvoked {
		if !(value.IsFunction()) || NodeUtil_ReferencesOwnReceiver(value) {
			return n
		}
	}
	if value.GetParent() != 0 {
		value.Detach()
	}
	var parent pkg_rhino.Node = n.GetParent()
	_ = parent
	if (NodeUtil_IsOptChainNode(parent) && n.IsFirstChildOf(parent)) && !(parent.IsOptionalChainStart()) {
		var endOfCurrentChain pkg_rhino.Node = NodeUtil_GetEndOfOptChainSegment(parent)
		_ = endOfCurrentChain
		NodeUtil_ConvertToNonOptionalChainSegment(endOfCurrentChain)
	}
	if keyIsGetter {
		value = pkg_rhino.IR_Call(value, []pkg_rhino.Node{})
		value.PutBooleanProp(pkg_rhino.StaticGet_Node_FREE_CALL(), true)
	} else if nIsInvoked {
		n.GetParent().PutBooleanProp(pkg_rhino.StaticGet_Node_FREE_CALL(), true)
	}
	n.ReplaceWith(value)
	this_.ReportChangeToEnclosingScope(value)
	this_.MarkFunctionsDeleted(n)
	return n
}

// Java static method: PeepholeFoldConstants.bigintNode(java.math.BigInteger val, com.google.javascript.rhino.Node srcref)
func PeepholeFoldConstants_BigintNode(val *j2gort.Ext_java_math_BigInteger, srcref pkg_rhino.Node) pkg_rhino.Node {
	clinit_PeepholeFoldConstants()
	_ = val
	_ = srcref
	var tree pkg_rhino.Node
	_ = tree
	if val.Signum() < int32(0) {
		tree = pkg_rhino.IR_Neg(pkg_rhino.IR_Bigint(val.Negate()))
	} else {
		tree = pkg_rhino.IR_Bigint(val)
	}
	return tree.SrcrefTree(srcref)
}

// Java static method: PeepholeFoldConstants.isGetterProp(com.google.javascript.rhino.Node key)
func PeepholeFoldConstants_IsGetterProp(key pkg_rhino.Node) bool {
	clinit_PeepholeFoldConstants()
	_ = key
	return ((key != 0) && (key.IsGetterDef() || (key.IsComputedProp() && key.GetBooleanProp(pkg_rhino.StaticGet_Node_COMPUTED_PROP_GETTER()))))
}

// Java static method: PeepholeFoldConstants.isSafeToFlattenObjectLitSpread(com.google.javascript.rhino.Node innerLit)
func PeepholeFoldConstants_IsSafeToFlattenObjectLitSpread(innerLit pkg_rhino.Node) bool {
	clinit_PeepholeFoldConstants()
	_ = innerLit
	{
		var prop pkg_rhino.Node = innerLit.GetFirstChild()
		_ = prop
		for ; prop != 0; prop = prop.GetNext() {
			if NodeUtil_IsGetOrSetKey(prop) {
				return false
			}
		}
	}
	return true
}

// Java static method: PeepholeFoldConstants.isSafeToFlattenArrayLitSpread(com.google.javascript.rhino.Node innerLit)
func PeepholeFoldConstants_IsSafeToFlattenArrayLitSpread(innerLit pkg_rhino.Node) bool {
	clinit_PeepholeFoldConstants()
	_ = innerLit
	{
		var n pkg_rhino.Node = innerLit.GetFirstChild()
		_ = n
		for ; n != 0; n = n.GetNext() {
			if n.IsEmpty() {
				return false
			}
		}
	}
	return true
}
