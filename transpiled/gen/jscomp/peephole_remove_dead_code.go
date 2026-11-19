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
// Port of com/google/javascript/jscomp/PeepholeRemoveDeadCode.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	pkg_base "github.com/paullepoulpe/jsgo/transpiled/gen/base"
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_parser "github.com/paullepoulpe/jsgo/transpiled/gen/parser"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

var lambdaSingleton_peephole_remove_dead_code_1 j2gort.Ext_com_google_common_base_Predicate
var lambdaSingleton_peephole_remove_dead_code_2 j2gort.Ext_com_google_common_base_Predicate
var lambdaSingleton_peephole_remove_dead_code_3 j2gort.Ext_com_google_common_base_Predicate
var lambdaSingleton_peephole_remove_dead_code_4 j2gort.Ext_com_google_common_base_Predicate

func init() {
	lambdaSingleton_peephole_remove_dead_code_1 = j2gort.NewNodeLambda_Ext_com_google_common_base_Predicate(func(rawNode_7 uint32) bool {
		node := pkg_rhino.Node(rawNode_7)
		_ = node
		return (node.IsBreak() && !(node.HasChildren()))
	})
	lambdaSingleton_peephole_remove_dead_code_2 = j2gort.NewNodeLambda_Ext_com_google_common_base_Predicate(func(rawNode_8 uint32) bool {
		node := pkg_rhino.Node(rawNode_8)
		_ = node
		return ((!(pkg_rhino.IR_MayBeExpression(node)) && !(NodeUtil_IsLoopStructure(node))) && !(node.IsSwitch()))
	})
	lambdaSingleton_peephole_remove_dead_code_3 = j2gort.NewNodeLambda_Ext_com_google_common_base_Predicate(func(rawNode_9 uint32) bool {
		node := pkg_rhino.Node(rawNode_9)
		_ = node
		return (node.IsContinue() && !(node.HasChildren()))
	})
	lambdaSingleton_peephole_remove_dead_code_4 = j2gort.NewNodeLambda_Ext_com_google_common_base_Predicate(func(rawNode_10 uint32) bool {
		node := pkg_rhino.Node(rawNode_10)
		_ = node
		return (!(pkg_rhino.IR_MayBeExpression(node)) && !(NodeUtil_IsLoopStructure(node)))
	})
}

// Java class: com.google.javascript.jscomp.PeepholeRemoveDeadCode
type PeepholeRemoveDeadCode struct {
	AbstractPeepholeOptimizationBase
}

func (this_ *PeepholeRemoveDeadCode) BindSelf_(s any) {
	_ = s
	this_.AbstractPeepholeOptimizationBase.BindSelf_(s)
}
func (this_ *PeepholeRemoveDeadCode) IsJavaNil_() bool { return this_ == nil }
func (*PeepholeRemoveDeadCode) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.PeepholeRemoveDeadCode", "PeepholeRemoveDeadCode"
}
func (this_ *PeepholeRemoveDeadCode) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.AbstractPeepholeOptimizationBase.IdentityHash_()
}
func (*PeepholeRemoveDeadCode) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.PeepholeRemoveDeadCode", "com.google.javascript.jscomp.PeepholeRemoveDeadCode", func(v_ any) bool {
		return j2gort.InstanceOf[*PeepholeRemoveDeadCode](v_)
	})
}
func (*PeepholeRemoveDeadCode) Implements_jscomp_AbstractPeepholeOptimization() {}

var static_PeepholeRemoveDeadCode_MATCH_UNNAMED_BREAK j2gort.Ext_com_google_common_base_Predicate

func StaticGet_PeepholeRemoveDeadCode_MATCH_UNNAMED_BREAK() j2gort.Ext_com_google_common_base_Predicate {
	clinit_PeepholeRemoveDeadCode()
	return static_PeepholeRemoveDeadCode_MATCH_UNNAMED_BREAK
}
func StaticSet_PeepholeRemoveDeadCode_MATCH_UNNAMED_BREAK(v j2gort.Ext_com_google_common_base_Predicate) {
	clinit_PeepholeRemoveDeadCode()
	static_PeepholeRemoveDeadCode_MATCH_UNNAMED_BREAK = v
}

var clinitDone_PeepholeRemoveDeadCode uint32
var clinitState_PeepholeRemoveDeadCode j2gort.ClinitState

func clinit_PeepholeRemoveDeadCode() {
	if atomic.LoadUint32(&clinitDone_PeepholeRemoveDeadCode) != 0 {
		return
	}
	clinitSlow_PeepholeRemoveDeadCode()
}

func clinitSlow_PeepholeRemoveDeadCode() {
	if !j2gort.EnterClinit(&clinitState_PeepholeRemoveDeadCode, &clinitDone_PeepholeRemoveDeadCode) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_PeepholeRemoveDeadCode, &clinitDone_PeepholeRemoveDeadCode)
	static_PeepholeRemoveDeadCode_MATCH_UNNAMED_BREAK = NewPeepholeRemoveDeadCode_MatchUnnamedBreak()
}

func (this_ *PeepholeRemoveDeadCode) initFields_PeepholeRemoveDeadCode() {
	_ = this_
}

// Java constructor: PeepholeRemoveDeadCode()
func NewPeepholeRemoveDeadCode() *PeepholeRemoveDeadCode {
	clinit_PeepholeRemoveDeadCode()
	this_ := &PeepholeRemoveDeadCode{}
	this_.BindSelf_(this_)
	InitPeepholeRemoveDeadCode(this_)
	return this_
}

func InitPeepholeRemoveDeadCode(this_ *PeepholeRemoveDeadCode) {
	InitAbstractPeepholeOptimization(&this_.AbstractPeepholeOptimizationBase)
	this_.initFields_PeepholeRemoveDeadCode()
}

// Java method: PeepholeRemoveDeadCode.optimizeSubtree(com.google.javascript.rhino.Node subtree)
func (this_ *PeepholeRemoveDeadCode) OptimizeSubtree(subtree pkg_rhino.Node) pkg_rhino.Node {
	_ = subtree
	switch subtree.TokenSwitchOrd_() {
	case int32(51) /* Token.ASSIGN */ :
		return this_.priv_PeepholeRemoveDeadCode_tryFoldAssignment_4222(subtree)
	case int32(50) /* Token.COMMA */ :
		return this_.priv_PeepholeRemoveDeadCode_tryFoldComma_4236(subtree)
	case int32(94) /* Token.SCRIPT */, int32(90) /* Token.BLOCK */ :
		return this_.TryOptimizeBlock(subtree)
	case int32(93) /* Token.EXPR_RESULT */ :
		return this_.priv_PeepholeRemoveDeadCode_tryFoldExpr_4224(subtree)
	case int32(67) /* Token.HOOK */ :
		return this_.priv_PeepholeRemoveDeadCode_tryFoldHook_4246(subtree)
	case int32(75) /* Token.SWITCH */ :
		return this_.priv_PeepholeRemoveDeadCode_tryOptimizeSwitch_4230(subtree)
	case int32(74) /* Token.IF */ :
		return this_.priv_PeepholeRemoveDeadCode_tryFoldIf_4245(subtree)
	case int32(78) /* Token.WHILE */ :
		return subtree
	case int32(80) /* Token.FOR */ :
		var condition pkg_rhino.Node = NodeUtil_GetConditionExpression(subtree)
		_ = condition
		if condition != 0 {
			this_.priv_PeepholeRemoveDeadCode_tryFoldForCondition_4253(condition)
		}
		return this_.TryFoldFor(subtree)
	case int32(79) /* Token.DO */ :
		var foldedDo pkg_rhino.Node = this_.TryFoldDoAway(subtree)
		_ = foldedDo
		if foldedDo.IsDo() {
			return this_.TryFoldEmptyDo(foldedDo)
		}
		return foldedDo
	case int32(48) /* Token.TRY */ :
		return this_.priv_PeepholeRemoveDeadCode_tryFoldTry_4221(subtree)
	case int32(92) /* Token.LABEL */ :
		return this_.priv_PeepholeRemoveDeadCode_tryFoldLabel_4220(subtree)
	case int32(102) /* Token.ARRAY_PATTERN */ :
		return this_.TryOptimizeArrayPattern(subtree)
	case int32(103) /* Token.OBJECT_PATTERN */ :
		return this_.TryOptimizeObjectPattern(subtree)
	case int32(84) /* Token.VAR */, int32(97) /* Token.CONST */, int32(111) /* Token.LET */ :
		return this_.priv_PeepholeRemoveDeadCode_tryOptimizeNameDeclaration_4223(subtree)
	case int32(134) /* Token.DEFAULT_VALUE */ :
		return this_.priv_PeepholeRemoveDeadCode_tryRemoveDefaultValue_4219(subtree)
	case int32(29) /* Token.OPTCHAIN_GETPROP */, int32(31) /* Token.OPTCHAIN_CALL */, int32(30) /* Token.OPTCHAIN_GETELEM */ :
		return this_.priv_PeepholeRemoveDeadCode_tryRemoveOptionalChaining_4254(subtree)
	default:
		return subtree
	}
	panic("unreachable switch expression")
}

// Java method: PeepholeRemoveDeadCode.tryRemoveDefaultValue(com.google.javascript.rhino.Node defaultValue)
func (this_ *PeepholeRemoveDeadCode) priv_PeepholeRemoveDeadCode_tryRemoveDefaultValue_4219(defaultValue pkg_rhino.Node) pkg_rhino.Node {
	_ = defaultValue
	if !(defaultValue.IsDefaultValue()) {
		j2gort.CheckArgumentMsg(false, j2gort.BoxNode(defaultValue))
	}
	var lValue pkg_rhino.Node = defaultValue.GetFirstChild()
	_ = lValue
	var val pkg_rhino.Node = defaultValue.GetSecondChild()
	_ = val
	var removeVal bool = false
	_ = removeVal
	if val.IsName() && (j2gort.CheckNotNullString(val.GetString()) == "undefined") {
		removeVal = true
	}
	if val.IsVoid() {
		var voidArg pkg_rhino.Node = val.GetFirstChild()
		_ = voidArg
		removeVal = !(this_.MayHaveSideEffects(voidArg))
	}
	if removeVal {
		defaultValue.ReplaceWith(lValue.Detach())
		this_.ReportChangeToEnclosingScope(lValue)
		return lValue
	}
	return defaultValue
}

// Java method: PeepholeRemoveDeadCode.tryFoldLabel(com.google.javascript.rhino.Node n)
func (this_ *PeepholeRemoveDeadCode) priv_PeepholeRemoveDeadCode_tryFoldLabel_4220(n pkg_rhino.Node) pkg_rhino.Node {
	_ = n
	var labelName string = n.GetFirstChild().GetString()
	_ = labelName
	var stmt pkg_rhino.Node = n.GetLastChild()
	_ = stmt
	if stmt.IsEmpty() {
		this_.ReportChangeToEnclosingScope(n)
		n.Detach()
		return pkg_rhino.Node(0)
	}
	if stmt.IsBlock() && !(stmt.HasChildren()) {
		this_.ReportChangeToEnclosingScope(n)
		if n.GetParent().IsLabel() {
			n.ReplaceWith(stmt.Detach())
		} else {
			n.Detach()
		}
		return pkg_rhino.Node(0)
	}
	var child pkg_rhino.Node = PeepholeRemoveDeadCode_GetOnlyInterestingChild(stmt)
	_ = child
	if child != 0 {
		stmt = child
	}
	if stmt.IsBreak() && (j2gort.CheckNotNullString(stmt.GetFirstChild().GetString()) == labelName) {
		if n.GetParent().IsLabel() {
			var replacement pkg_rhino.Node = pkg_rhino.IR_Block_0().Srcref(n)
			_ = replacement
			n.ReplaceWith(replacement)
			this_.ReportChangeToEnclosingScope(replacement)
			return replacement
		} else {
			var parent pkg_rhino.Node = n.GetParent()
			_ = parent
			n.Detach()
			this_.ReportChangeToEnclosingScope(parent)
			return pkg_rhino.Node(0)
		}
	}
	return n
}

// Java static method: PeepholeRemoveDeadCode.getOnlyInterestingChild(com.google.javascript.rhino.Node block)
func PeepholeRemoveDeadCode_GetOnlyInterestingChild(block pkg_rhino.Node) pkg_rhino.Node {
	clinit_PeepholeRemoveDeadCode()
	_ = block
	if !(block.IsBlock()) {
		return pkg_rhino.Node(0)
	}
	if block.HasOneChild() {
		return block.GetOnlyChild()
	}
	var ret pkg_rhino.Node = pkg_rhino.Node(0)
	_ = ret
	{
		var child pkg_rhino.Node = block.GetFirstChild()
		_ = child
		for ; child != 0; child = child.GetNext() {
			if child.IsSyntheticBlock() && !(child.HasChildren()) {
			} else if ret != 0 {
				return pkg_rhino.Node(0)
			} else {
				ret = child
			}
		}
	}
	return ret
}

// Java method: PeepholeRemoveDeadCode.tryFoldTry(com.google.javascript.rhino.Node n)
func (this_ *PeepholeRemoveDeadCode) priv_PeepholeRemoveDeadCode_tryFoldTry_4221(n pkg_rhino.Node) pkg_rhino.Node {
	_ = n
	if !(n.IsTry()) {
		j2gort.CheckStateMsg(false, j2gort.BoxNode(n))
	}
	var body pkg_rhino.Node = n.GetFirstChild()
	_ = body
	var catchBlock pkg_rhino.Node = body.GetNext()
	_ = catchBlock
	var finallyBlock pkg_rhino.Node = catchBlock.GetNext()
	_ = finallyBlock
	if !(catchBlock.HasChildren()) && ((finallyBlock == 0) || !(finallyBlock.HasChildren())) {
		if !(!(n.GetParent().IsLabel())) {
			j2gort.CheckState(false)
		}
		body.Detach()
		n.ReplaceWith(body)
		this_.ReportChangeToEnclosingScope(body)
		return body
	}
	if !(body.HasChildren()) {
		NodeUtil_RedeclareVarsInsideBranch(catchBlock)
		this_.ReportChangeToEnclosingScope(n)
		if finallyBlock != 0 {
			finallyBlock.Detach()
			if !(!(n.GetParent().IsLabel())) {
				j2gort.CheckState(false)
			}
			n.ReplaceWith(finallyBlock)
			return finallyBlock
		} else {
			if !(!(n.GetParent().IsLabel())) {
				j2gort.CheckState(false)
			}
			n.Detach()
			return pkg_rhino.Node(0)
		}
	}
	return n
}

// Java method: PeepholeRemoveDeadCode.tryFoldAssignment(com.google.javascript.rhino.Node subtree)
func (this_ *PeepholeRemoveDeadCode) priv_PeepholeRemoveDeadCode_tryFoldAssignment_4222(subtree pkg_rhino.Node) pkg_rhino.Node {
	_ = subtree
	if !(subtree.IsAssign()) {
		j2gort.CheckState(false)
	}
	var left pkg_rhino.Node = subtree.GetFirstChild()
	_ = left
	var right pkg_rhino.Node = subtree.GetLastChild()
	_ = right
	if (left.IsName() && right.IsName()) && func() bool {
		sr_1 := left.GetString()
		se_2 := right.GetString()
		return j2gort.CheckNotNullString(sr_1) == se_2
	}() {
		subtree.ReplaceWith(right.Detach())
		this_.ReportChangeToEnclosingScope(right)
		return right
	} else if left.IsDestructuringPattern() && !(left.HasChildren()) {
		subtree.ReplaceWith(right.Detach())
		this_.ReportChangeToEnclosingScope(right)
		return right
	}
	return subtree
}

// Java method: PeepholeRemoveDeadCode.tryOptimizeNameDeclaration(com.google.javascript.rhino.Node subtree)
func (this_ *PeepholeRemoveDeadCode) priv_PeepholeRemoveDeadCode_tryOptimizeNameDeclaration_4223(subtree pkg_rhino.Node) pkg_rhino.Node {
	_ = subtree
	if !(NodeUtil_IsNameDeclaration(subtree)) {
		j2gort.CheckState(false)
	}
	var left pkg_rhino.Node = subtree.GetFirstChild()
	_ = left
	if left.IsDestructuringLhs() && left.HasTwoChildren() {
		var pattern pkg_rhino.Node = left.GetFirstChild()
		_ = pattern
		if !(pattern.HasChildren()) {
			var value pkg_rhino.Node = left.GetSecondChild()
			_ = value
			subtree.ReplaceWith(pkg_rhino.IR_ExprResult(value.Detach()).Srcref(value))
			this_.ReportChangeToEnclosingScope(value)
		}
	}
	return subtree
}

// Java method: PeepholeRemoveDeadCode.tryFoldExpr(com.google.javascript.rhino.Node subtree)
func (this_ *PeepholeRemoveDeadCode) priv_PeepholeRemoveDeadCode_tryFoldExpr_4224(subtree pkg_rhino.Node) pkg_rhino.Node {
	_ = subtree
	var result pkg_rhino.Node = this_.priv_PeepholeRemoveDeadCode_trySimplifyUnusedResult_4225(subtree.GetFirstChild())
	_ = result
	if result == 0 {
		var parent pkg_rhino.Node = subtree.GetParent()
		_ = parent
		if parent.IsLabel() {
			var replacement pkg_rhino.Node = pkg_rhino.IR_Block_0().Srcref(subtree)
			_ = replacement
			subtree.ReplaceWith(replacement)
			subtree = replacement
		} else {
			subtree.Detach()
			subtree = pkg_rhino.Node(0)
		}
	}
	return subtree
}

// Java method: PeepholeRemoveDeadCode.trySimplifyUnusedResult(com.google.javascript.rhino.Node expression)
func (this_ *PeepholeRemoveDeadCode) priv_PeepholeRemoveDeadCode_trySimplifyUnusedResult_4225(expression pkg_rhino.Node) pkg_rhino.Node {
	_ = expression
	var sideEffectRoots j2gort.Deque[pkg_rhino.Node]
	_ = sideEffectRoots
	var atFixedPoint bool = this_.priv_PeepholeRemoveDeadCode_trySimplifyUnusedResultInternal_4226(expression, &sideEffectRoots)
	_ = atFixedPoint
	if atFixedPoint {
		return expression
	} else if sideEffectRoots.IsEmpty() {
		this_.DeleteNode(expression)
		return pkg_rhino.Node(0)
	} else if sideEffectRoots.PeekFirst() == expression {
		if !(sideEffectRoots.Size() == int32(1)) {
			j2gort.CheckStateStr(false, sideEffectRoots.String())
		}
		this_.ReportChangeToEnclosingScope(expression)
		return expression
	} else {
		var sideEffects pkg_rhino.Node = PeepholeRemoveDeadCode_AsDetachedExpression(sideEffectRoots.PollFirst())
		_ = sideEffects
		for !(sideEffectRoots.IsEmpty()) {
			var next pkg_rhino.Node = PeepholeRemoveDeadCode_AsDetachedExpression(sideEffectRoots.PollFirst())
			_ = next
			sideEffects = pkg_rhino.IR_Comma(sideEffects, next).Srcref(next)
		}
		sideEffects.InsertBefore(expression)
		this_.DeleteNode(expression)
		return sideEffects
	}
}

// Java method: PeepholeRemoveDeadCode.trySimplifyUnusedResultInternal(com.google.javascript.rhino.Node tree, java.util.ArrayDeque<com.google.javascript.rhino.Node> sideEffectRoots)
func (this_ *PeepholeRemoveDeadCode) priv_PeepholeRemoveDeadCode_trySimplifyUnusedResultInternal_4226(tree pkg_rhino.Node, sideEffectRoots *j2gort.Deque[pkg_rhino.Node]) bool {
	_ = tree
	_ = sideEffectRoots
	switch tree.TokenSwitchOrd_() {
	case int32(30) /* Token.OPTCHAIN_GETELEM */, int32(29) /* Token.OPTCHAIN_GETPROP */, int32(31) /* Token.OPTCHAIN_CALL */ :
		if this_.MayHaveSideEffects(tree) {
			sideEffectRoots.AddLast(j2gort.CheckNotNullNode(tree))
			return PeepholeRemoveDeadCode_HasFixedPointParent(tree)
		} else {
			return false
		}
	case int32(67) /* Token.HOOK */ :
		var trueNode pkg_rhino.Node = this_.priv_PeepholeRemoveDeadCode_trySimplifyUnusedResult_4225(tree.GetSecondChild())
		_ = trueNode
		var falseNode pkg_rhino.Node = this_.priv_PeepholeRemoveDeadCode_trySimplifyUnusedResult_4225(tree.GetLastChild())
		_ = falseNode
		if (trueNode == 0) && (falseNode != 0) {
			if !(tree.HasTwoChildren()) {
				j2gort.CheckStateMsg(false, j2gort.BoxNode(tree))
			}
			tree.SetToken(pkg_rhino.StaticGet_Token_OR())
			sideEffectRoots.AddLast(j2gort.CheckNotNullNode(tree))
			return false
		} else if (trueNode != 0) && (falseNode == 0) {
			if !(tree.HasTwoChildren()) {
				j2gort.CheckStateMsg(false, j2gort.BoxNode(tree))
			}
			tree.SetToken(pkg_rhino.StaticGet_Token_AND())
			sideEffectRoots.AddLast(j2gort.CheckNotNullNode(tree))
			return false
		} else if (trueNode == 0) && (falseNode == 0) {
			this_.priv_PeepholeRemoveDeadCode_trySimplifyUnusedResultInternal_4226(tree.GetOnlyChild(), sideEffectRoots)
			return false
		} else {
			sideEffectRoots.AddLast(j2gort.CheckNotNullNode(tree))
			return PeepholeRemoveDeadCode_HasFixedPointParent(tree)
		}
	case int32(69) /* Token.AND */, int32(68) /* Token.OR */, int32(70) /* Token.COALESCE */ :
		var conditionalResultNode pkg_rhino.Node = this_.priv_PeepholeRemoveDeadCode_trySimplifyUnusedResult_4225(tree.GetLastChild())
		_ = conditionalResultNode
		if conditionalResultNode == 0 {
			this_.priv_PeepholeRemoveDeadCode_trySimplifyUnusedResultInternal_4226(tree.GetOnlyChild(), sideEffectRoots)
			return false
		} else {
			sideEffectRoots.AddLast(j2gort.CheckNotNullNode(tree))
			return PeepholeRemoveDeadCode_HasFixedPointParent(tree)
		}
	case int32(73) /* Token.FUNCTION */ :
	case int32(108) /* Token.MEMBER_FIELD_DEF */ :
		if tree.IsStaticMember() && tree.HasChildren() {
			this_.priv_PeepholeRemoveDeadCode_trySimplifyUnusedResultInternal_4226(tree.GetFirstChild(), sideEffectRoots)
		}
	case int32(109) /* Token.COMPUTED_FIELD_DEF */ :
		this_.priv_PeepholeRemoveDeadCode_trySimplifyUnusedResultInternal_4226(tree.GetFirstChild(), sideEffectRoots)
		if tree.IsStaticMember() && tree.HasTwoChildren() {
			this_.priv_PeepholeRemoveDeadCode_trySimplifyUnusedResultInternal_4226(tree.GetLastChild(), sideEffectRoots)
		}
	default:
		if this_.NodeTypeMayHaveSideEffects(tree) {
			sideEffectRoots.AddLast(j2gort.CheckNotNullNode(tree))
			return PeepholeRemoveDeadCode_HasFixedPointParent(tree)
		} else if !(tree.HasChildren()) {
			return false
		}
		var atFixedPoint bool = PeepholeRemoveDeadCode_HasFixedPointParent(tree)
		_ = atFixedPoint
		{
			var child pkg_rhino.Node = tree.GetFirstChild()
			_ = child
			for ; child != 0; child = child.GetNext() {
				atFixedPoint = func() bool {
					ba_3 := atFixedPoint
					bb_4 := this_.priv_PeepholeRemoveDeadCode_trySimplifyUnusedResultInternal_4226(child, sideEffectRoots)
					return ba_3 && bb_4
				}()
			}
		}
		return atFixedPoint
	}
	return false
}

// Java static method: PeepholeRemoveDeadCode.asDetachedExpression(com.google.javascript.rhino.Node expr)
func PeepholeRemoveDeadCode_AsDetachedExpression(expr pkg_rhino.Node) pkg_rhino.Node {
	clinit_PeepholeRemoveDeadCode()
	_ = expr
	switch expr.TokenSwitchOrd_() {
	case int32(127) /* Token.ITER_SPREAD */, int32(128) /* Token.OBJECT_SPREAD */ :
		switch expr.GetParent().TokenSwitchOrd_() {
		case int32(46) /* Token.ARRAYLIT */, int32(23) /* Token.NEW */, int32(28) /* Token.CALL */, int32(31) /* Token.OPTCHAIN_CALL */ :
			expr = pkg_rhino.IR_Arraylit_1_Node([]pkg_rhino.Node{expr.Detach()}).Srcref(expr)
		case int32(47) /* Token.OBJECTLIT */ :
			expr = pkg_rhino.IR_Objectlit([]pkg_rhino.Node{expr.Detach()}).Srcref(expr)
		default:
			panic(j2gort.AsIface[javalang.IIllegalStateException](j2gort.ConstructShimThrowable("java.lang.IllegalStateException", j2gort.BoxString(expr.ToStringTree()))))
		}
	default:
	}
	if expr.HasParent() {
		expr.Detach()
	}
	if !(pkg_rhino.IR_MayBeExpression(expr)) {
		j2gort.CheckStateMsg(false, j2gort.BoxNode(expr))
	}
	return expr
}

// Java static method: PeepholeRemoveDeadCode.hasFixedPointParent(com.google.javascript.rhino.Node expr)
func PeepholeRemoveDeadCode_HasFixedPointParent(expr pkg_rhino.Node) bool {
	clinit_PeepholeRemoveDeadCode()
	_ = expr
	switch expr.GetParent().TokenSwitchOrd_() {
	case int32(69) /* Token.AND */, int32(50) /* Token.COMMA */, int32(67) /* Token.HOOK */, int32(68) /* Token.OR */, int32(70) /* Token.COALESCE */ :
		return true
	case int32(46) /* Token.ARRAYLIT */, int32(47) /* Token.OBJECTLIT */ :
		return expr.IsSpread()
	default:
		return NodeUtil_IsStatement(expr.GetParent())
	}
	panic("unreachable switch expression")
}

// Java method: PeepholeRemoveDeadCode.removeIfUnnamedBreak(com.google.javascript.rhino.Node maybeBreak)
func (this_ *PeepholeRemoveDeadCode) priv_PeepholeRemoveDeadCode_removeIfUnnamedBreak_4227(maybeBreak pkg_rhino.Node) {
	_ = maybeBreak
	if ((maybeBreak != 0) && maybeBreak.IsBreak()) && !(maybeBreak.HasChildren()) {
		this_.ReportChangeToEnclosingScope(maybeBreak)
		maybeBreak.Detach()
	}
}

// Java method: PeepholeRemoveDeadCode.tryRemoveSwitchWithSingleCase(com.google.javascript.rhino.Node switchNode, boolean shouldHoistCondition)
func (this_ *PeepholeRemoveDeadCode) priv_PeepholeRemoveDeadCode_tryRemoveSwitchWithSingleCase_4228(switchNode pkg_rhino.Node, shouldHoistCondition bool) pkg_rhino.Node {
	_ = switchNode
	_ = shouldHoistCondition
	var switchBody pkg_rhino.Node = switchNode.GetSecondChild()
	_ = switchBody
	var caseBlock pkg_rhino.Node = switchBody.GetOnlyChild().GetLastChild()
	_ = caseBlock
	this_.priv_PeepholeRemoveDeadCode_removeIfUnnamedBreak_4227(caseBlock.GetLastChild())
	if NodeUtil_Has(caseBlock, StaticGet_PeepholeRemoveDeadCode_MATCH_UNNAMED_BREAK(), StaticGet_NodeUtil_MATCH_NOT_FUNCTION()) {
		return switchNode
	}
	if shouldHoistCondition {
		pkg_rhino.IR_ExprResult(switchNode.RemoveFirstChild()).Srcref(switchNode).InsertBefore(switchNode)
	}
	switchNode.ReplaceWith(caseBlock.Detach())
	this_.ReportChangeToEnclosingScope(caseBlock)
	return caseBlock
}

// Java method: PeepholeRemoveDeadCode.tryRemoveSwitch(com.google.javascript.rhino.Node n)
func (this_ *PeepholeRemoveDeadCode) priv_PeepholeRemoveDeadCode_tryRemoveSwitch_4229(n pkg_rhino.Node) pkg_rhino.Node {
	_ = n
	var switchBody pkg_rhino.Node = n.GetSecondChild()
	_ = switchBody
	if !(switchBody.HasChildren()) {
		var condition pkg_rhino.Node = n.RemoveFirstChild()
		_ = condition
		var replacement pkg_rhino.Node = pkg_rhino.IR_ExprResult(condition).Srcref(n)
		_ = replacement
		n.ReplaceWith(replacement)
		this_.ReportChangeToEnclosingScope(replacement)
		return replacement
	} else if switchBody.HasOneChild() && switchBody.GetOnlyChild().IsDefaultCase() {
		if this_.MayHaveSideEffects(n.GetFirstChild()) {
			return this_.priv_PeepholeRemoveDeadCode_tryRemoveSwitchWithSingleCase_4228(n, true)
		} else {
			return this_.priv_PeepholeRemoveDeadCode_tryRemoveSwitchWithSingleCase_4228(n, false)
		}
	} else {
		return n
	}
}

// Java method: PeepholeRemoveDeadCode.tryOptimizeSwitch(com.google.javascript.rhino.Node n)
func (this_ *PeepholeRemoveDeadCode) priv_PeepholeRemoveDeadCode_tryOptimizeSwitch_4230(n pkg_rhino.Node) pkg_rhino.Node {
	_ = n
	if !(n.IsSwitch()) {
		j2gort.CheckStateMsg(false, j2gort.BoxNode(n))
	}
	var switchBody pkg_rhino.Node = n.GetSecondChild()
	_ = switchBody
	var defaultCase pkg_rhino.Node = this_.priv_PeepholeRemoveDeadCode_tryOptimizeDefaultCase_4231(n)
	_ = defaultCase
	if ((defaultCase == 0) || switchBody.GetLastChild().IsDefaultCase()) && this_.priv_PeepholeRemoveDeadCode_areAllCaseTagsLiterals_4235(switchBody) {
		var cond pkg_rhino.Node = n.GetFirstChild()
		_ = cond
		var prev pkg_rhino.Node = pkg_rhino.Node(0)
		_ = prev
		var next pkg_rhino.Node = pkg_rhino.Node(0)
		_ = next
		var cur pkg_rhino.Node = pkg_rhino.Node(0)
		_ = cur
		var foundMatchingCase bool = false
		_ = foundMatchingCase
		{
			cur = switchBody.GetFirstChild()
			for ; cur != 0; cur = next {
				next = cur.GetNext()
				var firstChild pkg_rhino.Node = cur.GetFirstChild()
				_ = firstChild
				foundMatchingCase = this_.priv_PeepholeRemoveDeadCode_isFirstSwitchMatch_4234(foundMatchingCase, cond, firstChild)
				if (!(foundMatchingCase) && !(this_.MayHaveSideEffects(firstChild))) && this_.priv_PeepholeRemoveDeadCode_isUselessCase_4233(cur, prev, defaultCase) {
					this_.priv_PeepholeRemoveDeadCode_removeCase_4232(n, cur)
				} else {
					prev = cur
				}
			}
		}
		if NodeUtil_IsLiteralValue(cond, false) {
			var caseLabel pkg_rhino.Node = pkg_rhino.Node(0)
			_ = caseLabel
			var caseMatches pkg_base.Tri = pkg_base.StaticGet_Tri_TRUE()
			_ = caseMatches
			{
				cur = switchBody.GetFirstChild()
				for ; cur != 0; cur = next {
					next = cur.GetNext()
					caseLabel = cur.GetFirstChild()
					caseMatches = PeepholeFoldConstants_EvaluateComparison(this_, pkg_rhino.StaticGet_Token_SHEQ(), cond, caseLabel)
					if caseMatches == pkg_base.StaticGet_Tri_TRUE() {
						break
					} else if caseMatches == pkg_base.StaticGet_Tri_UNKNOWN() {
						break
					} else {
						this_.priv_PeepholeRemoveDeadCode_removeCase_4232(n, cur)
					}
				}
			}
			if (cur != 0) && (caseMatches == pkg_base.StaticGet_Tri_TRUE()) {
				var matchingCase pkg_rhino.Node = cur
				_ = matchingCase
				var matchingCaseBlock pkg_rhino.Node = matchingCase.GetLastChild()
				_ = matchingCaseBlock
				for cur != 0 {
					var block pkg_rhino.Node = cur.GetLastChild()
					_ = block
					var lastStm pkg_rhino.Node = block.GetLastChild()
					_ = lastStm
					var isLastStmRemovableBreak bool = false
					_ = isLastStmRemovableBreak
					if (lastStm != 0) && PeepholeRemoveDeadCode_IsExit(lastStm) {
						this_.priv_PeepholeRemoveDeadCode_removeIfUnnamedBreak_4227(lastStm)
						isLastStmRemovableBreak = true
					}
					next = cur.GetNext()
					if cur != matchingCase {
						for block.HasChildren() {
							matchingCaseBlock.AddChildToBack(block.RemoveFirstChild())
						}
						this_.ReportChangeToEnclosingScope(cur)
						cur.Detach()
					}
					cur = next
					if isLastStmRemovableBreak {
						break
					}
				}
				for ; cur != 0; cur = next {
					next = cur.GetNext()
					this_.priv_PeepholeRemoveDeadCode_removeCase_4232(n, cur)
				}
				cur = cond.GetNext()
				if (cur != 0) && (cur.GetNext() == 0) {
					return this_.priv_PeepholeRemoveDeadCode_tryRemoveSwitchWithSingleCase_4228(n, false)
				}
			}
		}
	}
	return this_.priv_PeepholeRemoveDeadCode_tryRemoveSwitch_4229(n)
}

// Java method: PeepholeRemoveDeadCode.tryOptimizeDefaultCase(com.google.javascript.rhino.Node n)
func (this_ *PeepholeRemoveDeadCode) priv_PeepholeRemoveDeadCode_tryOptimizeDefaultCase_4231(n pkg_rhino.Node) pkg_rhino.Node {
	_ = n
	if !(n.IsSwitch()) {
		j2gort.CheckStateMsg(false, j2gort.BoxNode(n))
	}
	var switchBody pkg_rhino.Node = n.GetSecondChild()
	_ = switchBody
	var lastNonRemovable pkg_rhino.Node = pkg_rhino.Node(0)
	_ = lastNonRemovable
	{
		var c pkg_rhino.Node = switchBody.GetFirstChild()
		_ = c
		for ; c != 0; c = c.GetNext() {
			if c.IsDefaultCase() {
				var caseToRemove pkg_rhino.Node
				_ = caseToRemove
				if lastNonRemovable != 0 {
					caseToRemove = lastNonRemovable.GetNext()
				} else {
					caseToRemove = switchBody.GetFirstChild()
				}
				{
					var next pkg_rhino.Node = pkg_rhino.Node(0)
					_ = next
					for ; caseToRemove != c; caseToRemove = next {
						next = caseToRemove.GetNext()
						this_.priv_PeepholeRemoveDeadCode_removeCase_4232(n, caseToRemove)
					}
				}
				if this_.priv_PeepholeRemoveDeadCode_isUselessCase_4233(c, lastNonRemovable, c) {
					this_.priv_PeepholeRemoveDeadCode_removeCase_4232(n, c)
					return pkg_rhino.Node(0)
				}
				return c
			} else {
				if !(c.IsCase()) {
					j2gort.CheckState(false)
				}
				if c.GetLastChild().HasChildren() || this_.MayHaveSideEffects(c.GetFirstChild()) {
					lastNonRemovable = c
				}
			}
		}
	}
	return pkg_rhino.Node(0)
}

// Java method: PeepholeRemoveDeadCode.removeCase(com.google.javascript.rhino.Node switchNode, com.google.javascript.rhino.Node caseNode)
func (this_ *PeepholeRemoveDeadCode) priv_PeepholeRemoveDeadCode_removeCase_4232(switchNode pkg_rhino.Node, caseNode pkg_rhino.Node) {
	_ = switchNode
	_ = caseNode
	NodeUtil_RedeclareVarsInsideBranch(caseNode)
	caseNode.Detach()
	this_.ReportChangeToEnclosingScope(switchNode)
}

// Java method: PeepholeRemoveDeadCode.isUselessCase(com.google.javascript.rhino.Node caseNode, com.google.javascript.rhino.@org.jspecify.annotations.Nullable Node previousCase, com.google.javascript.rhino.@org.jspecify.annotations.Nullable Node defaultCase)
func (this_ *PeepholeRemoveDeadCode) priv_PeepholeRemoveDeadCode_isUselessCase_4233(caseNode pkg_rhino.Node, previousCase pkg_rhino.Node, defaultCase pkg_rhino.Node) bool {
	_ = caseNode
	_ = previousCase
	_ = defaultCase
	if !((previousCase == 0) || (previousCase.GetNext() == caseNode)) {
		j2gort.CheckState(false)
	}
	var switchBody pkg_rhino.Node = caseNode.GetParent()
	_ = switchBody
	if (switchBody.GetLastChild() != caseNode) && (previousCase != 0) {
		var previousBlock pkg_rhino.Node = previousCase.GetLastChild()
		_ = previousBlock
		if !(previousBlock.HasChildren()) || !(PeepholeRemoveDeadCode_IsExit(previousBlock.GetLastChild())) {
			return false
		}
	}
	var executingCase pkg_rhino.Node = caseNode
	_ = executingCase
	for executingCase != 0 {
		if !(executingCase.IsDefaultCase() || executingCase.IsCase()) {
			j2gort.CheckState(false)
		}
		if !((caseNode == executingCase) || !(executingCase.IsDefaultCase())) {
			j2gort.CheckState(false)
		}
		if !(executingCase.IsDefaultCase()) && this_.MayHaveSideEffects(executingCase.GetFirstChild()) {
			return false
		}
		var block pkg_rhino.Node = executingCase.GetLastChild()
		_ = block
		if !(block.IsBlock()) {
			j2gort.CheckState(false)
		}
		if block.HasChildren() {
			{
				var blockChild pkg_rhino.Node = block.GetFirstChild()
				_ = blockChild
				for ; blockChild != 0; blockChild = blockChild.GetNext() {
					switch blockChild.TokenSwitchOrd_() {
					case int32(82) /* Token.BREAK */ :
						return (!(blockChild.HasChildren()) && ((defaultCase == 0) || (defaultCase == executingCase)))
					case int32(84) /* Token.VAR */ :
						if blockChild.HasOneChild() && (blockChild.GetFirstFirstChild() == 0) {
							continue
						}
						return false
					default:
						return false
					}
				}
			}
		}
		executingCase = executingCase.GetNext()
	}
	return true
}

// Java method: PeepholeRemoveDeadCode.isFirstSwitchMatch(boolean foundMatchingCase, com.google.javascript.rhino.Node condition, com.google.javascript.rhino.Node tag)
func (this_ *PeepholeRemoveDeadCode) priv_PeepholeRemoveDeadCode_isFirstSwitchMatch_4234(foundMatchingCase bool, condition pkg_rhino.Node, tag pkg_rhino.Node) bool {
	_ = foundMatchingCase
	_ = condition
	_ = tag
	if foundMatchingCase {
		return false
	}
	return (PeepholeFoldConstants_EvaluateComparison(this_, pkg_rhino.StaticGet_Token_SHEQ(), condition, tag) == pkg_base.StaticGet_Tri_TRUE())
}

// Java method: PeepholeRemoveDeadCode.areAllCaseTagsLiterals(com.google.javascript.rhino.Node switchBody)
func (this_ *PeepholeRemoveDeadCode) priv_PeepholeRemoveDeadCode_areAllCaseTagsLiterals_4235(switchBody pkg_rhino.Node) bool {
	_ = switchBody
	{
		var caseNode pkg_rhino.Node = switchBody.GetFirstChild()
		_ = caseNode
		for ; caseNode != 0; caseNode = caseNode.GetNext() {
			if caseNode.IsDefaultCase() {
				continue
			}
			var caseTag pkg_rhino.Node = caseNode.GetFirstChild()
			_ = caseTag
			if !(NodeUtil_IsLiteralValue(caseTag, false)) {
				return false
			}
		}
	}
	return true
}

// Java static method: PeepholeRemoveDeadCode.isExit(com.google.javascript.rhino.Node n)
func PeepholeRemoveDeadCode_IsExit(n pkg_rhino.Node) bool {
	clinit_PeepholeRemoveDeadCode()
	_ = n
	switch n.TokenSwitchOrd_() {
	case int32(82) /* Token.BREAK */, int32(83) /* Token.CONTINUE */, int32(0) /* Token.RETURN */, int32(43) /* Token.THROW */ :
		return true
	case int32(75) /* Token.SWITCH */ :
		return PeepholeRemoveDeadCode_IsSwitchExit(n)
	case int32(48) /* Token.TRY */ :
		return PeepholeRemoveDeadCode_IsTryExit(n)
	default:
		return false
	}
	panic("unreachable switch expression")
}

// Java static method: PeepholeRemoveDeadCode.isUnconditionalBlockExit(com.google.javascript.rhino.Node n)
func PeepholeRemoveDeadCode_IsUnconditionalBlockExit(n pkg_rhino.Node) bool {
	clinit_PeepholeRemoveDeadCode()
	_ = n
	if !(n.IsBlock()) {
		j2gort.CheckStateMsg(false, j2gort.BoxNode(n))
	}
	if !(!(n.GetParent().IsLabel())) {
		j2gort.CheckStateMsg(false, j2gort.BoxNode(n))
	}
	var lastStm pkg_rhino.Node = n.GetLastChild()
	_ = lastStm
	if lastStm == 0 {
		return false
	}
	switch lastStm.TokenSwitchOrd_() {
	case int32(82) /* Token.BREAK */ :
		if !(lastStm.HasChildren()) {
			return false
		}
	case int32(0) /* Token.RETURN */, int32(43) /* Token.THROW */ :
	default:
		return false
	}
	{
		var child pkg_rhino.Node = n.GetFirstChild()
		_ = child
		for ; child != lastStm; child = child.GetNext() {
			switch child.TokenSwitchOrd_() {
			case int32(82) /* Token.BREAK */ :
				if !(child.HasChildren()) {
					return false
				}
			case int32(0) /* Token.RETURN */, int32(43) /* Token.THROW */, int32(73) /* Token.FUNCTION */, int32(84) /* Token.VAR */, int32(111) /* Token.LET */, int32(97) /* Token.CONST */, int32(93) /* Token.EXPR_RESULT */ :
			default:
				return false
			}
		}
	}
	return true
}

// Java static method: PeepholeRemoveDeadCode.isSwitchExit(com.google.javascript.rhino.Node n)
func PeepholeRemoveDeadCode_IsSwitchExit(n pkg_rhino.Node) bool {
	clinit_PeepholeRemoveDeadCode()
	_ = n
	if !(n.IsSwitch()) {
		j2gort.CheckStateMsg(false, j2gort.BoxNode(n))
	}
	var hasDefaultCase bool = false
	_ = hasDefaultCase
	var switchBody pkg_rhino.Node = n.GetSecondChild()
	_ = switchBody
	{
		var switchCase pkg_rhino.Node = switchBody.GetFirstChild()
		_ = switchCase
		for ; switchCase != 0; switchCase = switchCase.GetNext() {
			if switchCase.IsDefaultCase() {
				hasDefaultCase = true
			}
			var block pkg_rhino.Node = switchCase.GetLastChild()
			_ = block
			if (block.HasChildren() || (switchCase.GetNext() == 0)) && !(PeepholeRemoveDeadCode_IsUnconditionalBlockExit(block)) {
				return false
			}
		}
	}
	return hasDefaultCase
}

// Java static method: PeepholeRemoveDeadCode.isTryExit(com.google.javascript.rhino.Node n)
func PeepholeRemoveDeadCode_IsTryExit(n pkg_rhino.Node) bool {
	clinit_PeepholeRemoveDeadCode()
	_ = n
	if !(n.IsTry()) {
		j2gort.CheckStateMsg(false, j2gort.BoxNode(n))
	}
	if n.HasXChildren(int32(3)) {
		if PeepholeRemoveDeadCode_IsUnconditionalBlockExit(n.GetLastChild()) {
			return true
		}
	}
	if !(PeepholeRemoveDeadCode_IsUnconditionalBlockExit(n.GetFirstChild())) {
		return false
	}
	var catches pkg_rhino.Node = n.GetSecondChild()
	_ = catches
	return (!(catches.HasChildren()) || PeepholeRemoveDeadCode_IsUnconditionalBlockExit(catches.GetOnlyChild().GetLastChild()))
}

// Java method: PeepholeRemoveDeadCode.tryFoldComma(com.google.javascript.rhino.Node n)
func (this_ *PeepholeRemoveDeadCode) priv_PeepholeRemoveDeadCode_tryFoldComma_4236(n pkg_rhino.Node) pkg_rhino.Node {
	_ = n
	var parent pkg_rhino.Node = n.GetParent()
	_ = parent
	var left pkg_rhino.Node = n.GetFirstChild()
	_ = left
	var right pkg_rhino.Node = left.GetNext()
	_ = right
	left = this_.priv_PeepholeRemoveDeadCode_trySimplifyUnusedResult_4225(left)
	if (left == 0) || !(this_.MayHaveSideEffects(left)) {
		right.Detach()
		n.ReplaceWith(right)
		this_.ReportChangeToEnclosingScope(parent)
		return right
	}
	return n
}

// Java method: PeepholeRemoveDeadCode.tryOptimizeBlock(com.google.javascript.rhino.Node n)
func (this_ *PeepholeRemoveDeadCode) TryOptimizeBlock(n pkg_rhino.Node) pkg_rhino.Node {
	_ = n
	{
		var c pkg_rhino.Node = n.GetFirstChild()
		_ = c
		for c != 0 {
			var next pkg_rhino.Node = c.GetNext()
			_ = next
			if !(PeepholeRemoveDeadCode_IsUnremovableNode(c)) && !(this_.MayHaveSideEffects(c)) {
				c.Detach()
				this_.ReportChangeToEnclosingScope(n)
				this_.MarkFunctionsDeleted(c)
			} else if PeepholeRemoveDeadCode_IsExit(c) {
				this_.priv_PeepholeRemoveDeadCode_removeFollowingNodes_4238(c)
				break
			} else {
				this_.priv_PeepholeRemoveDeadCode_tryOptimizeConditionalAfterAssign_4240(c)
			}
			c = next
		}
	}
	if (n.IsSyntheticBlock() || n.IsScript()) || (n.GetParent() == 0) {
		return n
	}
	var parent pkg_rhino.Node = n.GetParent()
	_ = parent
	if NodeUtil_TryMergeBlock(n, this_.IsASTNormalized()) {
		this_.ReportChangeToEnclosingScope(parent)
		return pkg_rhino.Node(0)
	} else if parent.IsClassMembers() && !(n.HasChildren()) {
		n.Detach()
		this_.ReportChangeToEnclosingScope(parent)
		return pkg_rhino.Node(0)
	}
	return n
}

// Java method: PeepholeRemoveDeadCode.removeFollowingNodes(com.google.javascript.rhino.Node n)
func (this_ *PeepholeRemoveDeadCode) priv_PeepholeRemoveDeadCode_removeFollowingNodes_4238(n pkg_rhino.Node) {
	_ = n
	var parent pkg_rhino.Node = n.GetParent()
	_ = parent
	if parent.GetLastChild() != n {
		var changed bool = false
		_ = changed
		for n.GetNext() != 0 {
			var dead pkg_rhino.Node = n.GetNext()
			_ = dead
			if NodeUtil_IsFunctionDeclaration(dead) {
				n = dead
				continue
			}
			changed = true
			NodeUtil_RedeclareVarsInsideBranch(dead)
			this_.priv_PeepholeRemoveDeadCode_redeclareIfBlockScopedVar_4239(dead, parent)
			dead.Detach()
			this_.MarkFunctionsDeleted(dead)
		}
		if changed {
			this_.ReportChangeToEnclosingScope(parent)
		}
	}
}

// Java method: PeepholeRemoveDeadCode.redeclareIfBlockScopedVar(com.google.javascript.rhino.Node decl, com.google.javascript.rhino.Node parent)
func (this_ *PeepholeRemoveDeadCode) priv_PeepholeRemoveDeadCode_redeclareIfBlockScopedVar_4239(decl pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = decl
	_ = parent
	if !(!(NodeUtil_IsFunctionDeclaration(decl))) {
		j2gort.CheckStateFmt(false, "Unexpected function declaration %s", j2gort.BoxNode(decl))
	}
	if !(NodeUtil_IsBlockScopedDeclaration(decl)) {
		return
	}
	var blockScopedVars []pkg_rhino.Node = []pkg_rhino.Node{}
	_ = blockScopedVars
	if decl.IsClass() {
		blockScopedVars = append(blockScopedVars, decl.GetFirstChild())
	} else {
		NodeUtil_VisitLhsNodesInNode(decl, j2gort.NewLambda_Ext_java_util_function_Consumer(func(lp_5 any) {
			n := j2gort.CastNode[pkg_rhino.Node](lp_5)
			_ = n
			if n.IsName() {
				blockScopedVars = append(blockScopedVars, n)
			}
		}))
	}
	var addBefore pkg_rhino.Node = NodeUtil_GetInsertionPointAfterAllInnerFunctionDeclarations(parent)
	_ = addBefore
	for _, nameNode := range blockScopedVars {
		_ = nameNode
		var var_ pkg_rhino.Node = pkg_rhino.IR_Let_1(pkg_rhino.IR_Name(nameNode.GetString()).Srcref(nameNode)).Srcref(nameNode)
		_ = var_
		NodeUtil_CopyNameAnnotations(nameNode, var_.GetFirstChild())
		if addBefore != 0 {
			var_.InsertBefore(addBefore)
		} else {
			parent.AddChildToBack(var_)
		}
	}

	this_.AddFeatureToEnclosingScript(pkg_parser.StaticGet_FeatureSet_Feature_LET_DECLARATIONS())
}

// Java static method: PeepholeRemoveDeadCode.isUnremovableNode(com.google.javascript.rhino.Node n)
func PeepholeRemoveDeadCode_IsUnremovableNode(n pkg_rhino.Node) bool {
	clinit_PeepholeRemoveDeadCode()
	_ = n
	return ((n.IsBlock() && n.IsSyntheticBlock()) || n.IsScript())
}

// Java method: PeepholeRemoveDeadCode.tryOptimizeConditionalAfterAssign(com.google.javascript.rhino.Node n)
func (this_ *PeepholeRemoveDeadCode) priv_PeepholeRemoveDeadCode_tryOptimizeConditionalAfterAssign_4240(n pkg_rhino.Node) {
	_ = n
	var next pkg_rhino.Node = n.GetNext()
	_ = next
	if !(PeepholeRemoveDeadCode_IsSimpleAssignment(n)) {
		return
	}
	var conditionalRoot pkg_rhino.Node = this_.priv_PeepholeRemoveDeadCode_getConditionalRoot_4243(next)
	_ = conditionalRoot
	if conditionalRoot == 0 {
		return
	}
	var lhsAssign pkg_rhino.Node = this_.priv_PeepholeRemoveDeadCode_getSimpleAssignmentName_4241(n)
	_ = lhsAssign
	var condition pkg_rhino.Node = this_.priv_PeepholeRemoveDeadCode_getConditionalStatementCondition_4244(next)
	_ = condition
	if !(lhsAssign.MatchesName_1_Node(condition)) {
		return
	}
	var rhsAssign pkg_rhino.Node = this_.priv_PeepholeRemoveDeadCode_getSimpleAssignmentValue_4242(n)
	_ = rhsAssign
	switch conditionalRoot.TokenSwitchOrd_() {
	case int32(69) /* Token.AND */, int32(68) /* Token.OR */, int32(74) /* Token.IF */, int32(67) /* Token.HOOK */ :
		var value pkg_base.Tri = NodeUtil_GetBooleanValue(rhsAssign)
		_ = value
		if value != pkg_base.StaticGet_Tri_UNKNOWN() {
			var replacementConditionNode pkg_rhino.Node = NodeUtil_BooleanNode(value.ToBoolean(true))
			_ = replacementConditionNode
			condition.ReplaceWith(replacementConditionNode)
			this_.ReportChangeToEnclosingScope(replacementConditionNode)
		}
		return
	case int32(70) /* Token.COALESCE */ :
		var valueType NodeUtil_ValueType = NodeUtil_GetKnownValueType(rhsAssign)
		_ = valueType
		switch (valueType).Ordinal() {
		case int32(1) /* ValueType.NULL */, int32(2) /* ValueType.VOID */ :
			condition.ReplaceWith(NodeUtil_NewUndefinedNode(condition))
			this_.ReportChangeToEnclosingScope(conditionalRoot)
		case int32(3) /* ValueType.NUMBER */, int32(4) /* ValueType.BIGINT */, int32(5) /* ValueType.STRING */, int32(6) /* ValueType.BOOLEAN */, int32(7) /* ValueType.OBJECT */ :
			condition.ReplaceWith(pkg_rhino.IR_Number(float64(int32(0))).Srcref(condition))
			this_.ReportChangeToEnclosingScope(conditionalRoot)
		case int32(0) /* ValueType.UNDETERMINED */ :
		}
		return
	default:
		panic(j2gort.AsIface[javalang.IAssertionError](j2gort.ConstructShimThrowable("java.lang.AssertionError", "Unhandled condition "+j2gort.StrVal(conditionalRoot))))
	}
}

// Java static method: PeepholeRemoveDeadCode.isSimpleAssignment(com.google.javascript.rhino.Node n)
func PeepholeRemoveDeadCode_IsSimpleAssignment(n pkg_rhino.Node) bool {
	clinit_PeepholeRemoveDeadCode()
	_ = n
	if NodeUtil_IsExprAssign(n) && n.GetFirstFirstChild().IsName() {
		return true
	} else if (NodeUtil_IsNameDeclaration(n) && n.HasOneChild()) && (n.GetFirstFirstChild() != 0) {
		return true
	}
	return false
}

// Java method: PeepholeRemoveDeadCode.getSimpleAssignmentName(com.google.javascript.rhino.Node n)
func (this_ *PeepholeRemoveDeadCode) priv_PeepholeRemoveDeadCode_getSimpleAssignmentName_4241(n pkg_rhino.Node) pkg_rhino.Node {
	_ = n
	if !(PeepholeRemoveDeadCode_IsSimpleAssignment(n)) {
		j2gort.CheckState(false)
	}
	if NodeUtil_IsExprAssign(n) {
		return n.GetFirstFirstChild()
	} else {
		return n.GetFirstChild()
	}
}

// Java method: PeepholeRemoveDeadCode.getSimpleAssignmentValue(com.google.javascript.rhino.Node n)
func (this_ *PeepholeRemoveDeadCode) priv_PeepholeRemoveDeadCode_getSimpleAssignmentValue_4242(n pkg_rhino.Node) pkg_rhino.Node {
	_ = n
	if !(PeepholeRemoveDeadCode_IsSimpleAssignment(n)) {
		j2gort.CheckState(false)
	}
	return n.GetFirstChild().GetLastChild()
}

// Java method: PeepholeRemoveDeadCode.getConditionalRoot(com.google.javascript.rhino.Node n)
func (this_ *PeepholeRemoveDeadCode) priv_PeepholeRemoveDeadCode_getConditionalRoot_4243(n pkg_rhino.Node) pkg_rhino.Node {
	_ = n
	if n == 0 {
		return pkg_rhino.Node(0)
	}
	if n.IsIf() {
		return n
	} else if PeepholeRemoveDeadCode_IsExprConditional(n) {
		return n.GetFirstChild()
	}
	return pkg_rhino.Node(0)
}

// Java static method: PeepholeRemoveDeadCode.isExprConditional(com.google.javascript.rhino.Node n)
func PeepholeRemoveDeadCode_IsExprConditional(n pkg_rhino.Node) bool {
	clinit_PeepholeRemoveDeadCode()
	_ = n
	if n.IsExprResult() {
		switch n.GetFirstChild().TokenSwitchOrd_() {
		case int32(67) /* Token.HOOK */, int32(69) /* Token.AND */, int32(68) /* Token.OR */, int32(70) /* Token.COALESCE */ :
			return true
		default:
		}
	}
	return false
}

// Java method: PeepholeRemoveDeadCode.getConditionalStatementCondition(com.google.javascript.rhino.Node n)
func (this_ *PeepholeRemoveDeadCode) priv_PeepholeRemoveDeadCode_getConditionalStatementCondition_4244(n pkg_rhino.Node) pkg_rhino.Node {
	_ = n
	if n.IsIf() {
		return NodeUtil_GetConditionExpression(n)
	} else {
		if !(PeepholeRemoveDeadCode_IsExprConditional(n)) {
			j2gort.CheckState(false)
		}
		return n.GetFirstFirstChild()
	}
}

// Java method: PeepholeRemoveDeadCode.tryFoldIf(com.google.javascript.rhino.Node n)
func (this_ *PeepholeRemoveDeadCode) priv_PeepholeRemoveDeadCode_tryFoldIf_4245(n pkg_rhino.Node) pkg_rhino.Node {
	_ = n
	if !(n.IsIf()) {
		j2gort.CheckStateMsg(false, j2gort.BoxNode(n))
	}
	var parent pkg_rhino.Node = n.GetParent()
	_ = parent
	j2gort.CheckNotNullNode(parent)
	var type_ pkg_rhino.Token = n.GetToken()
	_ = type_
	var cond pkg_rhino.Node = n.GetFirstChild()
	_ = cond
	var thenBody pkg_rhino.Node = cond.GetNext()
	_ = thenBody
	var elseBody pkg_rhino.Node = thenBody.GetNext()
	_ = elseBody
	if (elseBody != 0) && !(this_.MayHaveSideEffects(elseBody)) {
		elseBody.Detach()
		this_.ReportChangeToEnclosingScope(n)
		elseBody = pkg_rhino.Node(0)
	}
	if !(this_.MayHaveSideEffects(thenBody)) && (elseBody != 0) {
		elseBody.Detach()
		thenBody.ReplaceWith(elseBody)
		var notCond pkg_rhino.Node = pkg_rhino.NewNode_1(pkg_rhino.StaticGet_Token_NOT())
		_ = notCond
		cond.ReplaceWith(notCond)
		this_.ReportChangeToEnclosingScope(n)
		notCond.AddChildToFront(cond)
		cond = notCond
		thenBody = cond.GetNext()
		elseBody = pkg_rhino.Node(0)
	}
	if !(this_.MayHaveSideEffects(thenBody)) && (elseBody == 0) {
		if this_.MayHaveSideEffects(cond) {
			cond.Detach()
			var replacement pkg_rhino.Node = NodeUtil_NewExpr(cond)
			_ = replacement
			n.ReplaceWith(replacement)
			this_.ReportChangeToEnclosingScope(parent)
			return replacement
		} else {
			NodeUtil_RemoveChild(parent, n)
			this_.ReportChangeToEnclosingScope(parent)
			return pkg_rhino.Node(0)
		}
	}
	var condValue pkg_base.Tri = NodeUtil_GetBooleanValue(cond)
	_ = condValue
	if condValue == pkg_base.StaticGet_Tri_UNKNOWN() {
		return n
	}
	if this_.MayHaveSideEffects(cond) {
		var newConditionValue bool = (condValue == pkg_base.StaticGet_Tri_TRUE())
		_ = newConditionValue
		if !(newConditionValue) && (elseBody == 0) {
			elseBody = pkg_rhino.IR_Block_0().Srcref(n)
			n.AddChildToBack(elseBody)
		}
		var newCond pkg_rhino.Node = NodeUtil_BooleanNode(newConditionValue)
		_ = newCond
		cond.ReplaceWith(newCond)
		var branchToKeep pkg_rhino.Node
		_ = branchToKeep
		if newConditionValue {
			branchToKeep = thenBody
		} else {
			branchToKeep = elseBody
		}
		branchToKeep.AddChildToFront(pkg_rhino.IR_ExprResult(cond).Srcref(cond))
		this_.ReportChangeToEnclosingScope(branchToKeep)
		cond = newCond
	}
	var condTrue bool = condValue.ToBoolean(true)
	_ = condTrue
	if n.HasTwoChildren() {
		if !(type_ == pkg_rhino.StaticGet_Token_IF()) {
			j2gort.CheckState(false)
		}
		if condTrue {
			var thenStmt pkg_rhino.Node = n.GetSecondChild()
			_ = thenStmt
			thenStmt.Detach()
			n.ReplaceWith(thenStmt)
			this_.ReportChangeToEnclosingScope(thenStmt)
			return thenStmt
		} else {
			NodeUtil_RedeclareVarsInsideBranch(n)
			NodeUtil_RemoveChild(parent, n)
			this_.ReportChangeToEnclosingScope(parent)
			this_.MarkFunctionsDeleted(n)
			return pkg_rhino.Node(0)
		}
	} else {
		var trueBranch pkg_rhino.Node = n.GetSecondChild()
		_ = trueBranch
		var falseBranch pkg_rhino.Node = trueBranch.GetNext()
		_ = falseBranch
		var branchToKeep_2 pkg_rhino.Node
		_ = branchToKeep_2
		if condTrue {
			branchToKeep_2 = trueBranch
		} else {
			branchToKeep_2 = falseBranch
		}
		var branchToRemove pkg_rhino.Node
		_ = branchToRemove
		if condTrue {
			branchToRemove = falseBranch
		} else {
			branchToRemove = trueBranch
		}
		NodeUtil_RedeclareVarsInsideBranch(branchToRemove)
		branchToKeep_2.Detach()
		n.ReplaceWith(branchToKeep_2)
		this_.ReportChangeToEnclosingScope(branchToKeep_2)
		this_.MarkFunctionsDeleted(n)
		return branchToKeep_2
	}
}

// Java method: PeepholeRemoveDeadCode.tryFoldHook(com.google.javascript.rhino.Node n)
func (this_ *PeepholeRemoveDeadCode) priv_PeepholeRemoveDeadCode_tryFoldHook_4246(n pkg_rhino.Node) pkg_rhino.Node {
	_ = n
	if !(n.IsHook()) {
		j2gort.CheckStateMsg(false, j2gort.BoxNode(n))
	}
	var parent pkg_rhino.Node = n.GetParent()
	_ = parent
	j2gort.CheckNotNullNode(parent)
	var cond pkg_rhino.Node = n.GetFirstChild()
	_ = cond
	var thenBody pkg_rhino.Node = cond.GetNext()
	_ = thenBody
	var elseBody pkg_rhino.Node = thenBody.GetNext()
	_ = elseBody
	var condValue pkg_base.Tri = NodeUtil_GetBooleanValue(cond)
	_ = condValue
	if condValue == pkg_base.StaticGet_Tri_UNKNOWN() {
		if !(this_.AreNodesEqualForInlining(thenBody, elseBody)) {
			return n
		}
	}
	var branchToKeep pkg_rhino.Node = pkg_rhino.Node(0)
	_ = branchToKeep
	var branchToRemove pkg_rhino.Node = pkg_rhino.Node(0)
	_ = branchToRemove
	if condValue.ToBoolean(true) {
		branchToKeep = thenBody
		branchToRemove = elseBody
	} else {
		branchToKeep = elseBody
		branchToRemove = thenBody
	}
	var replacement pkg_rhino.Node = pkg_rhino.Node(0)
	_ = replacement
	var condHasSideEffects bool = this_.MayHaveSideEffects(cond)
	_ = condHasSideEffects
	n.DetachChildren()
	if condHasSideEffects {
		replacement = pkg_rhino.IR_Comma(cond, branchToKeep).Srcref(n)
	} else {
		replacement = branchToKeep
		this_.MarkFunctionsDeleted(cond)
	}
	n.ReplaceWith(replacement)
	this_.ReportChangeToEnclosingScope(replacement)
	this_.MarkFunctionsDeleted(branchToRemove)
	return replacement
}

// Java method: PeepholeRemoveDeadCode.tryFoldFor(com.google.javascript.rhino.Node n)
func (this_ *PeepholeRemoveDeadCode) TryFoldFor(n pkg_rhino.Node) pkg_rhino.Node {
	_ = n
	if !(n.IsVanillaFor()) {
		j2gort.CheckArgument(false)
	}
	var init pkg_rhino.Node = n.GetFirstChild()
	_ = init
	var cond pkg_rhino.Node = init.GetNext()
	_ = cond
	var increment pkg_rhino.Node = cond.GetNext()
	_ = increment
	if !(init.IsEmpty()) && !(NodeUtil_IsNameDeclaration(init)) {
		init = this_.priv_PeepholeRemoveDeadCode_trySimplifyUnusedResult_4225(init)
		if init == 0 {
			init = pkg_rhino.IR_Empty().Srcref(n)
			n.AddChildToFront(init)
		}
	}
	if !(increment.IsEmpty()) {
		increment = this_.priv_PeepholeRemoveDeadCode_trySimplifyUnusedResult_4225(increment)
		if increment == 0 {
			increment = pkg_rhino.IR_Empty().Srcref(n)
			increment.InsertAfter(cond)
		}
	}
	if !(n.GetFirstChild().IsEmpty()) {
		return n
	}
	if NodeUtil_GetBooleanValue(cond) != pkg_base.StaticGet_Tri_FALSE() {
		return n
	}
	var parent pkg_rhino.Node = n.GetParent()
	_ = parent
	NodeUtil_RedeclareVarsInsideBranch(n)
	if !(this_.MayHaveSideEffects(cond)) {
		for parent.IsLabel() {
			n = parent
			parent = parent.GetParent()
		}
		n.Detach()
	} else {
		var statement pkg_rhino.Node = pkg_rhino.IR_ExprResult(cond.Detach()).SrcrefIfMissing(cond)
		_ = statement
		if parent.IsLabel() {
			var block pkg_rhino.Node = pkg_rhino.IR_Block_0()
			_ = block
			block.SrcrefIfMissing(statement)
			block.AddChildToFront(statement)
			statement = block
		}
		n.ReplaceWith(statement)
	}
	this_.ReportChangeToEnclosingScope(parent)
	return pkg_rhino.Node(0)
}

// Java method: PeepholeRemoveDeadCode.tryFoldDoAway(com.google.javascript.rhino.Node n)
func (this_ *PeepholeRemoveDeadCode) TryFoldDoAway(n pkg_rhino.Node) pkg_rhino.Node {
	_ = n
	if !(n.IsDo()) {
		j2gort.CheckArgument(false)
	}
	var cond pkg_rhino.Node = NodeUtil_GetConditionExpression(n)
	_ = cond
	if NodeUtil_GetBooleanValue(cond) != pkg_base.StaticGet_Tri_FALSE() {
		return n
	}
	var block pkg_rhino.Node = NodeUtil_GetLoopCodeBlock(n)
	_ = block
	if n.GetParent().IsLabel() || PeepholeRemoveDeadCode_HasUnnamedBreakOrContinue(block) {
		return n
	}
	var parent pkg_rhino.Node = n.GetParent()
	_ = parent
	n.ReplaceWith(block.Detach())
	if this_.MayHaveSideEffects(cond) {
		var condStatement pkg_rhino.Node = pkg_rhino.IR_ExprResult(cond.Detach()).Srcref(cond)
		_ = condStatement
		condStatement.InsertAfter(block)
	}
	this_.ReportChangeToEnclosingScope(parent)
	return block
}

// Java method: PeepholeRemoveDeadCode.tryFoldEmptyDo(com.google.javascript.rhino.Node n)
func (this_ *PeepholeRemoveDeadCode) TryFoldEmptyDo(n pkg_rhino.Node) pkg_rhino.Node {
	_ = n
	if !(n.IsDo()) {
		j2gort.CheckArgument(false)
	}
	var body pkg_rhino.Node = NodeUtil_GetLoopCodeBlock(n)
	_ = body
	if body.IsBlock() && !(body.HasChildren()) {
		var cond pkg_rhino.Node = NodeUtil_GetConditionExpression(n)
		_ = cond
		var forNode pkg_rhino.Node = pkg_rhino.IR_ForNode(pkg_rhino.IR_Empty().Srcref(n), cond.Detach(), pkg_rhino.IR_Empty().Srcref(n), body.Detach())
		_ = forNode
		n.ReplaceWith(forNode)
		this_.ReportChangeToEnclosingScope(forNode)
		return forNode
	}
	return n
}

// Java method: PeepholeRemoveDeadCode.tryOptimizeObjectPattern(com.google.javascript.rhino.Node pattern)
func (this_ *PeepholeRemoveDeadCode) TryOptimizeObjectPattern(pattern pkg_rhino.Node) pkg_rhino.Node {
	_ = pattern
	if !(pattern.IsObjectPattern()) {
		j2gort.CheckArgumentMsg(false, j2gort.BoxNode(pattern))
	}
	if pattern.HasChildren() && pattern.GetLastChild().IsRest() {
		return pattern
	}
	{
		var child pkg_rhino.Node = pattern.GetFirstChild()
		_ = child
		for child != 0 {
			var key pkg_rhino.Node = child
			_ = key
			child = key.GetNext()
			if !(key.IsStringKey()) {
				continue
			}
			if this_.priv_PeepholeRemoveDeadCode_isRemovableDestructuringTarget_4252(key.GetOnlyChild()) {
				key.Detach()
				this_.ReportChangeToEnclosingScope(pattern)
			}
		}
	}
	return pattern
}

// Java method: PeepholeRemoveDeadCode.tryOptimizeArrayPattern(com.google.javascript.rhino.Node pattern)
func (this_ *PeepholeRemoveDeadCode) TryOptimizeArrayPattern(pattern pkg_rhino.Node) pkg_rhino.Node {
	_ = pattern
	if !(pattern.IsArrayPattern()) {
		j2gort.CheckArgumentMsg(false, j2gort.BoxNode(pattern))
	}
	{
		var lastChild pkg_rhino.Node = pattern.GetLastChild()
		_ = lastChild
		for lastChild != 0 {
			if lastChild.IsEmpty() || this_.priv_PeepholeRemoveDeadCode_isRemovableDestructuringTarget_4252(lastChild) {
				var prev pkg_rhino.Node = lastChild.GetPrevious()
				_ = prev
				lastChild.Detach()
				lastChild = prev
				this_.ReportChangeToEnclosingScope(pattern)
			} else {
				break
			}
		}
	}
	return pattern
}

// Java method: PeepholeRemoveDeadCode.isRemovableDestructuringTarget(com.google.javascript.rhino.Node destructruringElement)
func (this_ *PeepholeRemoveDeadCode) priv_PeepholeRemoveDeadCode_isRemovableDestructuringTarget_4252(destructruringElement pkg_rhino.Node) bool {
	_ = destructruringElement
	var target pkg_rhino.Node = destructruringElement
	_ = target
	var defaultValue pkg_rhino.Node = pkg_rhino.Node(0)
	_ = defaultValue
	if destructruringElement.IsDefaultValue() {
		target = destructruringElement.GetFirstChild()
		defaultValue = destructruringElement.GetSecondChild()
	}
	if !(target.IsDestructuringPattern()) || target.HasChildren() {
		return false
	}
	return ((defaultValue == 0) || !(this_.MayHaveSideEffects(defaultValue)))
}

// Java static method: PeepholeRemoveDeadCode.hasUnnamedBreakOrContinue(com.google.javascript.rhino.Node n)
func PeepholeRemoveDeadCode_HasUnnamedBreakOrContinue(n pkg_rhino.Node) bool {
	clinit_PeepholeRemoveDeadCode()
	_ = n
	return (NodeUtil_Has(n, lambdaSingleton_peephole_remove_dead_code_1, lambdaSingleton_peephole_remove_dead_code_2) || NodeUtil_Has(n, lambdaSingleton_peephole_remove_dead_code_3, lambdaSingleton_peephole_remove_dead_code_4))
}

// Java method: PeepholeRemoveDeadCode.tryFoldForCondition(com.google.javascript.rhino.Node forCondition)
func (this_ *PeepholeRemoveDeadCode) priv_PeepholeRemoveDeadCode_tryFoldForCondition_4253(forCondition pkg_rhino.Node) {
	_ = forCondition
	if this_.GetSideEffectFreeBooleanValue(forCondition) == pkg_base.StaticGet_Tri_TRUE() {
		this_.ReportChangeToEnclosingScope(forCondition)
		forCondition.ReplaceWith(pkg_rhino.IR_Empty())
	}
}

// Java method: PeepholeRemoveDeadCode.tryRemoveOptionalChaining(com.google.javascript.rhino.Node optionalChain)
func (this_ *PeepholeRemoveDeadCode) priv_PeepholeRemoveDeadCode_tryRemoveOptionalChaining_4254(optionalChain pkg_rhino.Node) pkg_rhino.Node {
	_ = optionalChain
	var callee pkg_rhino.Node = optionalChain.GetFirstChild()
	_ = callee
	if !(NodeUtil_IsNullOrUndefined(callee)) {
		return optionalChain
	}
	var result pkg_rhino.Node = pkg_rhino.Node(0)
	_ = result
	if this_.MayHaveSideEffects(callee) {
		optionalChain.ReplaceWith(callee.Detach())
		result = callee
	} else {
		result = NodeUtil_NewUndefinedNode(callee)
		optionalChain.ReplaceWith(result)
	}
	this_.MarkFunctionsDeleted(optionalChain)
	this_.ReportChangeToEnclosingScope(result)
	return result
}

// Java class: com.google.javascript.jscomp.PeepholeRemoveDeadCode.MatchUnnamedBreak
type PeepholeRemoveDeadCode_MatchUnnamedBreak struct {
	idHash_ atomic.Int32
}

func (this_ *PeepholeRemoveDeadCode_MatchUnnamedBreak) BindSelf_(s any) {
	_ = s
}
func (this_ *PeepholeRemoveDeadCode_MatchUnnamedBreak) IsJavaNil_() bool { return this_ == nil }
func (*PeepholeRemoveDeadCode_MatchUnnamedBreak) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.PeepholeRemoveDeadCode$MatchUnnamedBreak", "MatchUnnamedBreak"
}
func (this_ *PeepholeRemoveDeadCode_MatchUnnamedBreak) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*PeepholeRemoveDeadCode_MatchUnnamedBreak) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.PeepholeRemoveDeadCode.MatchUnnamedBreak", "com.google.javascript.jscomp.PeepholeRemoveDeadCode$MatchUnnamedBreak", func(v_ any) bool {
		return j2gort.InstanceOf[*PeepholeRemoveDeadCode_MatchUnnamedBreak](v_)
	})
}
func (*PeepholeRemoveDeadCode_MatchUnnamedBreak) Implements_j2gort_Ext_com_google_common_base_Predicate() {
}
func (*PeepholeRemoveDeadCode_MatchUnnamedBreak) Implements_j2gort_Ext_java_util_function_Predicate() {
}

func (this_ *PeepholeRemoveDeadCode_MatchUnnamedBreak) initFields_PeepholeRemoveDeadCode_MatchUnnamedBreak() {
	_ = this_
}

// Java constructor: MatchUnnamedBreak()
func NewPeepholeRemoveDeadCode_MatchUnnamedBreak() *PeepholeRemoveDeadCode_MatchUnnamedBreak {
	this_ := &PeepholeRemoveDeadCode_MatchUnnamedBreak{}
	InitPeepholeRemoveDeadCode_MatchUnnamedBreak(this_)
	return this_
}

func InitPeepholeRemoveDeadCode_MatchUnnamedBreak(this_ *PeepholeRemoveDeadCode_MatchUnnamedBreak) {
	this_.initFields_PeepholeRemoveDeadCode_MatchUnnamedBreak()
}

// Java method: MatchUnnamedBreak.apply(com.google.javascript.rhino.Node n)
func (this_ *PeepholeRemoveDeadCode_MatchUnnamedBreak) TestNode_(rawNode_11 uint32) bool {
	if this_ == nil {
		j2gort.ThrowNPE()
	}
	n := pkg_rhino.Node(rawNode_11)
	_ = n
	return (n.IsBreak() && !(n.HasChildren()))
}

func (this_ *PeepholeRemoveDeadCode_MatchUnnamedBreak) Apply(arg_ any) bool {
	if this_ == nil {
		j2gort.ThrowNPE()
	}
	return this_.TestNode_(j2gort.CastNode[uint32](arg_))
}

func (this_ *PeepholeRemoveDeadCode_MatchUnnamedBreak) Equals(p0_ any) bool {
	return j2gort.RefEquals(this_, p0_)
}

func (this_ *PeepholeRemoveDeadCode_MatchUnnamedBreak) Test(p0_ any) bool {
	return this_.Apply(p0_)
}
