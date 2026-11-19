/*
 * Copyright 2010 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/PeepholeMinimizeConditions.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	pkg_base "github.com/paullepoulpe/jsgo/transpiled/gen/base"
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

var lambdaSingleton_peephole_minimize_conditions_1 j2gort.Ext_com_google_common_base_Predicate

func init() {
	lambdaSingleton_peephole_minimize_conditions_1 = j2gort.NewNodeLambda_Ext_com_google_common_base_Predicate(func(rawNode_5 uint32) bool {
		input := pkg_rhino.Node(rawNode_5)
		_ = input
		return (input.IsGetProp() && input.GetParent().IsAssign())
	})
}

// Java class: com.google.javascript.jscomp.PeepholeMinimizeConditions
type PeepholeMinimizeConditions struct {
	AbstractPeepholeOptimizationBase
	late                                 bool `gson:"late"`
	slab_MinimizedCondition_MeasuredNode j2gort.PooledSlab[MinimizedCondition_MeasuredNode]
	slab_MinimizedCondition              j2gort.PooledSlab[MinimizedCondition]
}

func (this_ *PeepholeMinimizeConditions) BindSelf_(s any) {
	_ = s
	this_.AbstractPeepholeOptimizationBase.BindSelf_(s)
}
func (this_ *PeepholeMinimizeConditions) IsJavaNil_() bool { return this_ == nil }
func (*PeepholeMinimizeConditions) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.PeepholeMinimizeConditions", "PeepholeMinimizeConditions"
}
func (this_ *PeepholeMinimizeConditions) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.AbstractPeepholeOptimizationBase.IdentityHash_()
}
func (*PeepholeMinimizeConditions) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.PeepholeMinimizeConditions", "com.google.javascript.jscomp.PeepholeMinimizeConditions", func(v_ any) bool {
		return j2gort.InstanceOf[*PeepholeMinimizeConditions](v_)
	})
}
func (*PeepholeMinimizeConditions) Implements_jscomp_AbstractPeepholeOptimization() {}

var static_PeepholeMinimizeConditions_AND_PRECEDENCE int32

func StaticGet_PeepholeMinimizeConditions_AND_PRECEDENCE() int32 {
	clinit_PeepholeMinimizeConditions()
	return static_PeepholeMinimizeConditions_AND_PRECEDENCE
}
func StaticSet_PeepholeMinimizeConditions_AND_PRECEDENCE(v int32) {
	clinit_PeepholeMinimizeConditions()
	static_PeepholeMinimizeConditions_AND_PRECEDENCE = v
}

var clinitDone_PeepholeMinimizeConditions uint32
var clinitState_PeepholeMinimizeConditions j2gort.ClinitState

func clinit_PeepholeMinimizeConditions() {
	if atomic.LoadUint32(&clinitDone_PeepholeMinimizeConditions) != 0 {
		return
	}
	clinitSlow_PeepholeMinimizeConditions()
}

func clinitSlow_PeepholeMinimizeConditions() {
	if !j2gort.EnterClinit(&clinitState_PeepholeMinimizeConditions, &clinitDone_PeepholeMinimizeConditions) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_PeepholeMinimizeConditions, &clinitDone_PeepholeMinimizeConditions)
	static_PeepholeMinimizeConditions_AND_PRECEDENCE = NodeUtil_Precedence(pkg_rhino.StaticGet_Token_AND())
}

func (this_ *PeepholeMinimizeConditions) initFields_PeepholeMinimizeConditions() {
	_ = this_
}

// Java constructor: PeepholeMinimizeConditions(boolean late)
func NewPeepholeMinimizeConditions(late bool) *PeepholeMinimizeConditions {
	clinit_PeepholeMinimizeConditions()
	this_ := &PeepholeMinimizeConditions{}
	this_.BindSelf_(this_)
	InitPeepholeMinimizeConditions(this_, late)
	return this_
}

func InitPeepholeMinimizeConditions(this_ *PeepholeMinimizeConditions, late bool) {
	_ = late
	InitAbstractPeepholeOptimization(&this_.AbstractPeepholeOptimizationBase)
	this_.initFields_PeepholeMinimizeConditions()
	this_.late = late
}

// Java method: PeepholeMinimizeConditions.optimizeSubtree(com.google.javascript.rhino.Node node)
func (this_ *PeepholeMinimizeConditions) OptimizeSubtree(node pkg_rhino.Node) pkg_rhino.Node {
	_ = node
	switch node.TokenSwitchOrd_() {
	case int32(43) /* Token.THROW */, int32(0) /* Token.RETURN */ :
		var result pkg_rhino.Node = this_.priv_PeepholeMinimizeConditions_tryRemoveRedundantExit_4203(node)
		_ = result
		if result != node {
			return result
		}
		return this_.priv_PeepholeMinimizeConditions_tryReplaceExitWithBreak_4202(node)
	case int32(19) /* Token.NOT */ :
		this_.priv_PeepholeMinimizeConditions_tryMinimizeCondition_4211(node.GetFirstChild())
		return this_.priv_PeepholeMinimizeConditions_tryMinimizeNot_4206(node)
	case int32(74) /* Token.IF */ :
		this_.priv_PeepholeMinimizeConditions_performConditionSubstitutions_4213(node.GetFirstChild())
		return this_.priv_PeepholeMinimizeConditions_tryMinimizeIf_4209(node)
	case int32(93) /* Token.EXPR_RESULT */ :
		this_.priv_PeepholeMinimizeConditions_performConditionSubstitutions_4213(node.GetFirstChild())
		return this_.priv_PeepholeMinimizeConditions_tryMinimizeExprResult_4207(node)
	case int32(67) /* Token.HOOK */ :
		this_.priv_PeepholeMinimizeConditions_performConditionSubstitutions_4213(node.GetFirstChild())
		return this_.priv_PeepholeMinimizeConditions_tryMinimizeHook_4208(node)
	case int32(78) /* Token.WHILE */, int32(79) /* Token.DO */ :
		this_.priv_PeepholeMinimizeConditions_tryMinimizeCondition_4211(NodeUtil_GetConditionExpression(node))
		return node
	case int32(80) /* Token.FOR */ :
		this_.priv_PeepholeMinimizeConditions_tryJoinForCondition_4200(node)
		this_.priv_PeepholeMinimizeConditions_tryMinimizeCondition_4211(NodeUtil_GetConditionExpression(node))
		return node
	case int32(90) /* Token.BLOCK */ :
		return this_.priv_PeepholeMinimizeConditions_tryReplaceIf_4201(node)
	default:
		return node
	}
	panic("unreachable switch expression")
}

// Java method: PeepholeMinimizeConditions.tryJoinForCondition(com.google.javascript.rhino.Node n)
func (this_ *PeepholeMinimizeConditions) priv_PeepholeMinimizeConditions_tryJoinForCondition_4200(n pkg_rhino.Node) {
	_ = n
	if !(this_.late) {
		return
	}
	var block pkg_rhino.Node = n.GetLastChild()
	_ = block
	var maybeIf pkg_rhino.Node = block.GetFirstChild()
	_ = maybeIf
	if (maybeIf != 0) && maybeIf.IsIf() {
		var thenBlock pkg_rhino.Node = maybeIf.GetSecondChild()
		_ = thenBlock
		var maybeBreak pkg_rhino.Node = thenBlock.GetFirstChild()
		_ = maybeBreak
		if ((maybeBreak != 0) && maybeBreak.IsBreak()) && !(maybeBreak.HasChildren()) {
			if maybeIf.HasXChildren(int32(3)) {
				maybeIf.ReplaceWith(maybeIf.GetLastChild().Detach())
			} else {
				NodeUtil_RedeclareVarsInsideBranch(thenBlock)
				block.RemoveFirstChild()
			}
			var ifCondition pkg_rhino.Node = maybeIf.RemoveFirstChild()
			_ = ifCondition
			var fixedIfCondition pkg_rhino.Node = pkg_rhino.IR_Not(ifCondition).Srcref(ifCondition)
			_ = fixedIfCondition
			var forCondition pkg_rhino.Node = NodeUtil_GetConditionExpression(n)
			_ = forCondition
			if forCondition.IsEmpty() {
				forCondition.ReplaceWith(fixedIfCondition)
				this_.ReportChangeToEnclosingScope(fixedIfCondition)
			} else {
				var replacement pkg_rhino.Node = pkg_rhino.NewNode_1(pkg_rhino.StaticGet_Token_AND())
				_ = replacement
				forCondition.ReplaceWith(replacement)
				replacement.AddChildToBack(forCondition)
				replacement.AddChildToBack(fixedIfCondition)
				this_.ReportChangeToEnclosingScope(replacement)
			}
		}
	}
}

// Java method: PeepholeMinimizeConditions.tryReplaceIf(com.google.javascript.rhino.Node n)
func (this_ *PeepholeMinimizeConditions) priv_PeepholeMinimizeConditions_tryReplaceIf_4201(n pkg_rhino.Node) pkg_rhino.Node {
	_ = n
	var next pkg_rhino.Node = pkg_rhino.Node(0)
	_ = next
	{
		var child pkg_rhino.Node = n.GetFirstChild()
		_ = child
		for ; child != 0; child = next {
			next = child.GetNext()
			if child.IsIf() {
				var cond pkg_rhino.Node = child.GetFirstChild()
				_ = cond
				var thenBranch pkg_rhino.Node = cond.GetNext()
				_ = thenBranch
				var elseBranch pkg_rhino.Node = thenBranch.GetNext()
				_ = elseBranch
				var nextNode pkg_rhino.Node = child.GetNext()
				_ = nextNode
				if (((nextNode != 0) && (elseBranch == 0)) && PeepholeMinimizeConditions_IsReturnBlock(thenBranch)) && nextNode.IsIf() {
					var nextCond pkg_rhino.Node = nextNode.GetFirstChild()
					_ = nextCond
					var nextThen pkg_rhino.Node = nextCond.GetNext()
					_ = nextThen
					var nextElse pkg_rhino.Node = nextThen.GetNext()
					_ = nextElse
					if this_.AreNodesEqualForInlining(thenBranch, nextThen) {
						child.Detach()
						child.DetachChildren()
						var newCond pkg_rhino.Node = pkg_rhino.NewNode_2(pkg_rhino.StaticGet_Token_OR(), cond)
						_ = newCond
						nextCond.ReplaceWith(newCond)
						newCond.AddChildToBack(nextCond)
						this_.ReportChangeToEnclosingScope(newCond)
					} else if (nextElse != 0) && this_.AreNodesEqualForInlining(thenBranch, nextElse) {
						child.Detach()
						child.DetachChildren()
						var newCond_2 pkg_rhino.Node = pkg_rhino.NewNode_2(pkg_rhino.StaticGet_Token_AND(), pkg_rhino.IR_Not(cond).Srcref(cond))
						_ = newCond_2
						nextCond.ReplaceWith(newCond_2)
						newCond_2.AddChildToBack(nextCond)
						this_.ReportChangeToEnclosingScope(newCond_2)
					}
				} else if (((nextNode != 0) && (elseBranch == 0)) && PeepholeMinimizeConditions_IsReturnBlock(thenBranch)) && PeepholeMinimizeConditions_IsReturnExpression(nextNode) {
					var thenExpr pkg_rhino.Node = pkg_rhino.Node(0)
					_ = thenExpr
					if PeepholeMinimizeConditions_IsReturnExpressBlock(thenBranch) {
						thenExpr = PeepholeMinimizeConditions_GetBlockReturnExpression(thenBranch)
						thenExpr.Detach()
					} else {
						thenExpr = NodeUtil_NewUndefinedNode(child)
					}
					var elseExpr pkg_rhino.Node = nextNode.GetFirstChild()
					_ = elseExpr
					cond.Detach()
					elseExpr.Detach()
					var returnNode pkg_rhino.Node = pkg_rhino.IR_ReturnNode_1(pkg_rhino.IR_Hook(cond, thenExpr, elseExpr).Srcref(child))
					_ = returnNode
					child.ReplaceWith(returnNode)
					nextNode.Detach()
					this_.ReportChangeToEnclosingScope(n)
					break
				} else if (elseBranch != 0) && PeepholeMinimizeConditions_StatementMustExitParent(thenBranch) {
					elseBranch.Detach()
					elseBranch.InsertAfter(child)
					this_.ReportChangeToEnclosingScope(n)
				}
			}
		}
	}
	return n
}

// Java static method: PeepholeMinimizeConditions.statementMustExitParent(com.google.javascript.rhino.Node n)
func PeepholeMinimizeConditions_StatementMustExitParent(n pkg_rhino.Node) bool {
	clinit_PeepholeMinimizeConditions()
	_ = n
	switch n.TokenSwitchOrd_() {
	case int32(43) /* Token.THROW */, int32(0) /* Token.RETURN */ :
		return true
	case int32(90) /* Token.BLOCK */ :
		if n.HasChildren() {
			var child pkg_rhino.Node = n.GetLastChild()
			_ = child
			return PeepholeMinimizeConditions_StatementMustExitParent(child)
		}
		return false
	case int32(73) /* Token.FUNCTION */ :
		fallthrough
	default:
		return false
	}
	panic("unreachable")
}

// Java method: PeepholeMinimizeConditions.tryReplaceExitWithBreak(com.google.javascript.rhino.Node n)
func (this_ *PeepholeMinimizeConditions) priv_PeepholeMinimizeConditions_tryReplaceExitWithBreak_4202(n pkg_rhino.Node) pkg_rhino.Node {
	_ = n
	var result pkg_rhino.Node = n.GetFirstChild()
	_ = result
	var breakTarget pkg_rhino.Node = n
	_ = breakTarget
	for ; !(ControlFlowAnalysis_IsBreakTarget(breakTarget, j2gort.NullString)); breakTarget = breakTarget.GetParent() {
		if breakTarget.IsFunction() || breakTarget.IsScript() {
			return n
		}
	}
	var follow pkg_rhino.Node = ControlFlowAnalysis_ComputeFollowNode_1(breakTarget)
	_ = follow
	var prefinallyFollows pkg_rhino.Node = follow
	_ = prefinallyFollows
	follow = PeepholeMinimizeConditions_SkipFinallyNodes(follow)
	var innerFollow pkg_rhino.Node = ControlFlowAnalysis_ComputeFollowNode_1(n)
	_ = innerFollow
	if (prefinallyFollows != follow) || (innerFollow != PeepholeMinimizeConditions_SkipFinallyNodes(innerFollow)) {
		if !(this_.IsPure(result)) {
			return n
		}
	}
	if (follow == 0) && (n.IsThrow() || (result != 0)) {
		return n
	}
	if (follow == 0) || this_.AreMatchingExits(n, follow) {
		var replacement pkg_rhino.Node = pkg_rhino.IR_BreakNode_0()
		_ = replacement
		n.ReplaceWith(replacement)
		this_.ReportChangeToEnclosingScope(replacement)
		return replacement
	}
	return n
}

// Java method: PeepholeMinimizeConditions.tryRemoveRedundantExit(com.google.javascript.rhino.Node n)
func (this_ *PeepholeMinimizeConditions) priv_PeepholeMinimizeConditions_tryRemoveRedundantExit_4203(n pkg_rhino.Node) pkg_rhino.Node {
	_ = n
	var exitExpr pkg_rhino.Node = n.GetFirstChild()
	_ = exitExpr
	var follow pkg_rhino.Node = ControlFlowAnalysis_ComputeFollowNode_1(n)
	_ = follow
	var prefinallyFollows pkg_rhino.Node = follow
	_ = prefinallyFollows
	follow = PeepholeMinimizeConditions_SkipFinallyNodes(follow)
	if prefinallyFollows != follow {
		if !(this_.IsPure(exitExpr)) {
			return n
		}
	}
	if (follow == 0) && (n.IsThrow() || (exitExpr != 0)) {
		return n
	}
	if (follow == 0) || this_.AreMatchingExits(n, follow) {
		this_.ReportChangeToEnclosingScope(n)
		n.Detach()
		return pkg_rhino.Node(0)
	}
	return n
}

// Java method: PeepholeMinimizeConditions.isPure(com.google.javascript.rhino.Node n)
func (this_ *PeepholeMinimizeConditions) IsPure(n pkg_rhino.Node) bool {
	_ = n
	return ((n == 0) || (!(NodeUtil_CanBeSideEffected_1(n)) && !(this_.MayHaveSideEffects(n))))
}

// Java static method: PeepholeMinimizeConditions.skipFinallyNodes(com.google.javascript.rhino.Node n)
func PeepholeMinimizeConditions_SkipFinallyNodes(n pkg_rhino.Node) pkg_rhino.Node {
	clinit_PeepholeMinimizeConditions()
	_ = n
	for (n != 0) && NodeUtil_IsTryFinallyNode(n.GetParent(), n) {
		n = ControlFlowAnalysis_ComputeFollowNode_1(n)
	}
	return n
}

// Java method: PeepholeMinimizeConditions.areMatchingExits(com.google.javascript.rhino.Node nodeThis, com.google.javascript.rhino.Node nodeThat)
func (this_ *PeepholeMinimizeConditions) AreMatchingExits(nodeThis pkg_rhino.Node, nodeThat pkg_rhino.Node) bool {
	_ = nodeThis
	_ = nodeThat
	if (!(this_.IsASTNormalized()) && (nodeThis.IsThrow() || nodeThis.IsReturn())) && nodeThis.HasChildren() {
		return false
	}
	return (nodeThis.IsEquivalentTo_1(nodeThat) && (!(PeepholeMinimizeConditions_IsExceptionPossible(nodeThis)) || (PeepholeMinimizeConditions_GetExceptionHandler(nodeThis) == PeepholeMinimizeConditions_GetExceptionHandler(nodeThat))))
}

// Java static method: PeepholeMinimizeConditions.isExceptionPossible(com.google.javascript.rhino.Node n)
func PeepholeMinimizeConditions_IsExceptionPossible(n pkg_rhino.Node) bool {
	clinit_PeepholeMinimizeConditions()
	_ = n
	if !(n.IsReturn() || n.IsThrow()) {
		j2gort.CheckStateMsg(false, j2gort.BoxNode(n))
	}
	return (n.IsThrow() || (n.HasChildren() && !(NodeUtil_IsLiteralValue(n.GetLastChild(), true))))
}

// Java static method: PeepholeMinimizeConditions.getExceptionHandler(com.google.javascript.rhino.Node n)
func PeepholeMinimizeConditions_GetExceptionHandler(n pkg_rhino.Node) pkg_rhino.Node {
	clinit_PeepholeMinimizeConditions()
	_ = n
	return ControlFlowAnalysis_GetExceptionHandler(n)
}

// Java method: PeepholeMinimizeConditions.tryMinimizeNot(com.google.javascript.rhino.Node n)
func (this_ *PeepholeMinimizeConditions) priv_PeepholeMinimizeConditions_tryMinimizeNot_4206(n pkg_rhino.Node) pkg_rhino.Node {
	_ = n
	if !(n.IsNot()) {
		j2gort.CheckArgument(false)
	}
	var parent pkg_rhino.Node = n.GetParent()
	_ = parent
	var notChild pkg_rhino.Node = n.GetFirstChild()
	_ = notChild
	var complementOperator pkg_rhino.Token = pkg_rhino.Token(0)
	_ = complementOperator
	switch notChild.TokenSwitchOrd_() {
	case int32(4) /* Token.EQ */ :
		complementOperator = pkg_rhino.StaticGet_Token_NE()
	case int32(5) /* Token.NE */ :
		complementOperator = pkg_rhino.StaticGet_Token_EQ()
	case int32(40) /* Token.SHEQ */ :
		complementOperator = pkg_rhino.StaticGet_Token_SHNE()
	case int32(41) /* Token.SHNE */ :
		complementOperator = pkg_rhino.StaticGet_Token_SHEQ()
	case int32(19) /* Token.NOT */ :
		if NodeUtil_IsBooleanResult(notChild.GetFirstChild()) {
			var expr pkg_rhino.Node = notChild.RemoveFirstChild()
			_ = expr
			n.ReplaceWith(expr)
			this_.ReportChangeToEnclosingScope(parent)
			return expr
		}
		return n
	default:
		return n
	}
	var newOperator pkg_rhino.Node = n.RemoveFirstChild()
	_ = newOperator
	newOperator.SetToken(complementOperator)
	n.ReplaceWith(newOperator)
	this_.ReportChangeToEnclosingScope(parent)
	return newOperator
}

// Java method: PeepholeMinimizeConditions.tryMinimizeExprResult(com.google.javascript.rhino.Node n)
func (this_ *PeepholeMinimizeConditions) priv_PeepholeMinimizeConditions_tryMinimizeExprResult_4207(n pkg_rhino.Node) pkg_rhino.Node {
	_ = n
	var originalExpr pkg_rhino.Node = n.GetFirstChild()
	_ = originalExpr
	if !(MinimizedCondition_IsMinimizable(originalExpr)) {
		return n
	}
	var minCond *MinimizedCondition = minimizedConditionFromConditionNodeWithSlab_(&this_.slab_MinimizedCondition_MeasuredNode, &this_.slab_MinimizedCondition, originalExpr)
	_ = minCond
	var mNode *MinimizedCondition_MeasuredNode = minCond.getMinimizedWithSlab_(&this_.slab_MinimizedCondition_MeasuredNode, StaticGet_MinimizedCondition_MinimizationStyle_ALLOW_LEADING_NOT())
	_ = mNode
	if mNode.IsNot() {
		this_.priv_PeepholeMinimizeConditions_replaceNode_4212(originalExpr, mNode.withoutNotWithSlab_(&this_.slab_MinimizedCondition_MeasuredNode))
	} else {
		this_.priv_PeepholeMinimizeConditions_replaceNode_4212(originalExpr, mNode)
	}
	return n
}

// Java method: PeepholeMinimizeConditions.tryMinimizeHook(com.google.javascript.rhino.Node n)
func (this_ *PeepholeMinimizeConditions) priv_PeepholeMinimizeConditions_tryMinimizeHook_4208(n pkg_rhino.Node) pkg_rhino.Node {
	_ = n
	var originalCond pkg_rhino.Node = n.GetFirstChild()
	_ = originalCond
	var thenBranch pkg_rhino.Node = originalCond.GetNext()
	_ = thenBranch
	var elseBranch pkg_rhino.Node = thenBranch.GetNext()
	_ = elseBranch
	if NodeUtil_IsBooleanResult(thenBranch) && NodeUtil_IsBooleanResult(elseBranch) {
		var thenVal pkg_base.Tri = this_.GetSideEffectFreeBooleanValue(thenBranch)
		_ = thenVal
		var elseVal pkg_base.Tri = this_.GetSideEffectFreeBooleanValue(elseBranch)
		_ = elseVal
		var replacement pkg_rhino.Node = pkg_rhino.Node(0)
		_ = replacement
		if (thenVal == pkg_base.StaticGet_Tri_TRUE()) && (elseVal == pkg_base.StaticGet_Tri_FALSE()) {
			var cond pkg_rhino.Node = originalCond.Detach()
			_ = cond
			if NodeUtil_IsBooleanResult(cond) {
				replacement = cond
			} else {
				replacement = pkg_rhino.IR_Not(pkg_rhino.IR_Not(cond).Srcref(n)).Srcref(n)
			}
		} else if (thenVal == pkg_base.StaticGet_Tri_FALSE()) && (elseVal == pkg_base.StaticGet_Tri_TRUE()) {
			var cond_2 pkg_rhino.Node = originalCond.Detach()
			_ = cond_2
			replacement = pkg_rhino.IR_Not(cond_2).Srcref(n)
		}
		if replacement != 0 {
			n.ReplaceWith(replacement)
			this_.ReportChangeToEnclosingScope(replacement)
			this_.MarkFunctionsDeleted(n)
			if replacement.IsNot() {
				return this_.priv_PeepholeMinimizeConditions_tryMinimizeNot_4206(replacement)
			} else {
				return replacement
			}
		}
	}
	if MinimizedCondition_IsMinimizable(originalCond) {
		var minCond *MinimizedCondition = minimizedConditionFromConditionNodeWithSlab_(&this_.slab_MinimizedCondition_MeasuredNode, &this_.slab_MinimizedCondition, originalCond)
		_ = minCond
		var mNode *MinimizedCondition_MeasuredNode = minCond.getMinimizedWithSlab_(&this_.slab_MinimizedCondition_MeasuredNode, StaticGet_MinimizedCondition_MinimizationStyle_ALLOW_LEADING_NOT())
		_ = mNode
		if mNode.IsNot() {
			this_.priv_PeepholeMinimizeConditions_replaceNode_4212(originalCond, mNode.withoutNotWithSlab_(&this_.slab_MinimizedCondition_MeasuredNode))
			thenBranch.Detach()
			n.AddChildToBack(thenBranch)
			this_.ReportChangeToEnclosingScope(n)
		} else {
			this_.priv_PeepholeMinimizeConditions_replaceNode_4212(originalCond, mNode)
		}
	}
	return n
}

// Java method: PeepholeMinimizeConditions.tryMinimizeIf(com.google.javascript.rhino.Node n)
func (this_ *PeepholeMinimizeConditions) priv_PeepholeMinimizeConditions_tryMinimizeIf_4209(n pkg_rhino.Node) pkg_rhino.Node {
	_ = n
	var parent pkg_rhino.Node = n.GetParent()
	_ = parent
	var originalCond pkg_rhino.Node = n.GetFirstChild()
	_ = originalCond
	if NodeUtil_IsLiteralValue(originalCond, true) {
		return n
	}
	var thenBranch pkg_rhino.Node = originalCond.GetNext()
	_ = thenBranch
	var elseBranch pkg_rhino.Node = thenBranch.GetNext()
	_ = elseBranch
	var unnegatedCond *MinimizedCondition_MeasuredNode = nil
	_ = unnegatedCond
	var shortCond *MinimizedCondition_MeasuredNode = nil
	_ = shortCond
	if MinimizedCondition_IsMinimizable(originalCond) {
		var minCond *MinimizedCondition = minimizedConditionFromConditionNodeWithSlab_(&this_.slab_MinimizedCondition_MeasuredNode, &this_.slab_MinimizedCondition, originalCond)
		_ = minCond
		unnegatedCond = minCond.getMinimizedWithSlab_(&this_.slab_MinimizedCondition_MeasuredNode, StaticGet_MinimizedCondition_MinimizationStyle_PREFER_UNNEGATED())
		shortCond = minCond.getMinimizedWithSlab_(&this_.slab_MinimizedCondition_MeasuredNode, StaticGet_MinimizedCondition_MinimizationStyle_ALLOW_LEADING_NOT())
	}
	if elseBranch == 0 {
		if PeepholeMinimizeConditions_IsFoldableExpressBlock(thenBranch) {
			var expr pkg_rhino.Node = PeepholeMinimizeConditions_GetBlockExpression(thenBranch)
			_ = expr
			if !(this_.late) && PeepholeMinimizeConditions_IsPropertyAssignmentInExpression(expr) {
				this_.priv_PeepholeMinimizeConditions_replaceNode_4212(originalCond, unnegatedCond)
				return n
			}
			if (shortCond != nil) && shortCond.IsNot() {
				var replacementCond pkg_rhino.Node = this_.priv_PeepholeMinimizeConditions_replaceNode_4212(originalCond, shortCond.withoutNotWithSlab_(&this_.slab_MinimizedCondition_MeasuredNode)).Detach()
				_ = replacementCond
				var or pkg_rhino.Node = pkg_rhino.IR_Or(replacementCond, expr.RemoveFirstChild()).Srcref(n)
				_ = or
				var newExpr pkg_rhino.Node = NodeUtil_NewExpr(or)
				_ = newExpr
				n.ReplaceWith(newExpr)
				this_.ReportChangeToEnclosingScope(parent)
				return newExpr
			}
			var shortLowerThanAnd bool
			_ = shortLowerThanAnd
			if shortCond != nil {
				shortLowerThanAnd = shortCond.IsLowerPrecedenceThan(StaticGet_PeepholeMinimizeConditions_AND_PRECEDENCE())
			} else {
				shortLowerThanAnd = PeepholeMinimizeConditions_IsLowerPrecedence(originalCond, StaticGet_PeepholeMinimizeConditions_AND_PRECEDENCE())
			}
			if shortLowerThanAnd && PeepholeMinimizeConditions_IsLowerPrecedence(expr.GetFirstChild(), StaticGet_PeepholeMinimizeConditions_AND_PRECEDENCE()) {
				this_.priv_PeepholeMinimizeConditions_replaceNode_4212(originalCond, shortCond)
				return n
			}
			var replacementCond_2 pkg_rhino.Node = this_.priv_PeepholeMinimizeConditions_replaceNode_4212(originalCond, shortCond).Detach()
			_ = replacementCond_2
			var and pkg_rhino.Node = pkg_rhino.IR_And(replacementCond_2, expr.RemoveFirstChild()).Srcref(n)
			_ = and
			var newExpr_2 pkg_rhino.Node = NodeUtil_NewExpr(and)
			_ = newExpr_2
			n.ReplaceWith(newExpr_2)
			this_.ReportChangeToEnclosingScope(parent)
			return newExpr_2
		} else {
			if NodeUtil_IsStatementBlock(thenBranch) && thenBranch.HasOneChild() {
				var innerIf pkg_rhino.Node = thenBranch.GetFirstChild()
				_ = innerIf
				if innerIf.IsIf() {
					var innerCond pkg_rhino.Node = innerIf.GetFirstChild()
					_ = innerCond
					var innerThenBranch pkg_rhino.Node = innerCond.GetNext()
					_ = innerThenBranch
					var innerElseBranch pkg_rhino.Node = innerThenBranch.GetNext()
					_ = innerElseBranch
					var unnegatedLowerThanAnd bool
					_ = unnegatedLowerThanAnd
					if unnegatedCond != nil {
						unnegatedLowerThanAnd = unnegatedCond.IsLowerPrecedenceThan(StaticGet_PeepholeMinimizeConditions_AND_PRECEDENCE())
					} else {
						unnegatedLowerThanAnd = PeepholeMinimizeConditions_IsLowerPrecedence(originalCond, StaticGet_PeepholeMinimizeConditions_AND_PRECEDENCE())
					}
					if (innerElseBranch == 0) && !(unnegatedLowerThanAnd && PeepholeMinimizeConditions_IsLowerPrecedence(innerCond, StaticGet_PeepholeMinimizeConditions_AND_PRECEDENCE())) {
						var replacementCond_3 pkg_rhino.Node = this_.priv_PeepholeMinimizeConditions_replaceNode_4212(originalCond, unnegatedCond).Detach()
						_ = replacementCond_3
						n.DetachChildren()
						n.AddChildToBack(pkg_rhino.IR_And(replacementCond_3, innerCond.Detach()).Srcref(originalCond))
						n.AddChildToBack(innerThenBranch.Detach())
						this_.ReportChangeToEnclosingScope(n)
						return n
					}
				}
			}
		}
		this_.priv_PeepholeMinimizeConditions_replaceNode_4212(originalCond, unnegatedCond)
		return n
	}
	this_.priv_PeepholeMinimizeConditions_tryRemoveRepeatedStatements_4210(n)
	if ((shortCond != nil) && shortCond.IsNot()) && !(PeepholeMinimizeConditions_ConsumesDanglingElse(elseBranch)) {
		this_.priv_PeepholeMinimizeConditions_replaceNode_4212(originalCond, shortCond.withoutNotWithSlab_(&this_.slab_MinimizedCondition_MeasuredNode))
		thenBranch.Detach()
		n.AddChildToBack(thenBranch)
		this_.ReportChangeToEnclosingScope(n)
		return n
	}
	if PeepholeMinimizeConditions_IsReturnExpressBlock(thenBranch) && PeepholeMinimizeConditions_IsReturnExpressBlock(elseBranch) {
		var thenExpr pkg_rhino.Node = PeepholeMinimizeConditions_GetBlockReturnExpression(thenBranch)
		_ = thenExpr
		var elseExpr pkg_rhino.Node = PeepholeMinimizeConditions_GetBlockReturnExpression(elseBranch)
		_ = elseExpr
		var replacementCond_4 pkg_rhino.Node = this_.priv_PeepholeMinimizeConditions_replaceNode_4212(originalCond, shortCond).Detach()
		_ = replacementCond_4
		thenExpr.Detach()
		elseExpr.Detach()
		var returnNode pkg_rhino.Node = pkg_rhino.IR_ReturnNode_1(pkg_rhino.IR_Hook(replacementCond_4, thenExpr, elseExpr).Srcref(n))
		_ = returnNode
		n.ReplaceWith(returnNode)
		this_.ReportChangeToEnclosingScope(returnNode)
		return returnNode
	}
	var thenBranchIsExpressionBlock bool = PeepholeMinimizeConditions_IsFoldableExpressBlock(thenBranch)
	_ = thenBranchIsExpressionBlock
	var elseBranchIsExpressionBlock bool = PeepholeMinimizeConditions_IsFoldableExpressBlock(elseBranch)
	_ = elseBranchIsExpressionBlock
	if thenBranchIsExpressionBlock && elseBranchIsExpressionBlock {
		var thenOp pkg_rhino.Node = PeepholeMinimizeConditions_GetBlockExpression(thenBranch).GetFirstChild()
		_ = thenOp
		var elseOp pkg_rhino.Node = PeepholeMinimizeConditions_GetBlockExpression(elseBranch).GetFirstChild()
		_ = elseOp
		if thenOp.RawToken_() == elseOp.RawToken_() {
			if NodeUtil_IsAssignmentOp(thenOp) {
				var lhs pkg_rhino.Node = thenOp.GetFirstChild()
				_ = lhs
				if (this_.AreNodesEqualForInlining(lhs, elseOp.GetFirstChild()) && !(this_.MayEffectMutableState(lhs))) && (!(this_.MayHaveSideEffects(originalCond)) || (thenOp.IsAssign() && thenOp.GetFirstChild().IsName())) {
					var replacementCond_5 pkg_rhino.Node = this_.priv_PeepholeMinimizeConditions_replaceNode_4212(originalCond, shortCond).Detach()
					_ = replacementCond_5
					var assignName pkg_rhino.Node = thenOp.RemoveFirstChild()
					_ = assignName
					var thenExpr_2 pkg_rhino.Node = thenOp.RemoveFirstChild()
					_ = thenExpr_2
					var elseExpr_2 pkg_rhino.Node = elseOp.GetLastChild()
					_ = elseExpr_2
					elseExpr_2.Detach()
					var hookNode pkg_rhino.Node = pkg_rhino.IR_Hook(replacementCond_5, thenExpr_2, elseExpr_2).Srcref(n)
					_ = hookNode
					var assign pkg_rhino.Node = pkg_rhino.NewNode_3(thenOp.GetToken(), assignName, hookNode).Srcref(thenOp)
					_ = assign
					var expr_2 pkg_rhino.Node = NodeUtil_NewExpr(assign)
					_ = expr_2
					n.ReplaceWith(expr_2)
					this_.ReportChangeToEnclosingScope(parent)
					return expr_2
				}
			}
		}
		var replacementCond_6 pkg_rhino.Node = this_.priv_PeepholeMinimizeConditions_replaceNode_4212(originalCond, shortCond).Detach()
		_ = replacementCond_6
		thenOp.Detach()
		elseOp.Detach()
		var expr_3 pkg_rhino.Node = pkg_rhino.IR_ExprResult(pkg_rhino.IR_Hook(replacementCond_6, thenOp, elseOp).Srcref(n))
		_ = expr_3
		n.ReplaceWith(expr_3)
		this_.ReportChangeToEnclosingScope(parent)
		return expr_3
	}
	var thenBranchIsVar bool = PeepholeMinimizeConditions_IsVarBlock(thenBranch)
	_ = thenBranchIsVar
	var elseBranchIsVar bool = PeepholeMinimizeConditions_IsVarBlock(elseBranch)
	_ = elseBranchIsVar
	if (thenBranchIsVar && elseBranchIsExpressionBlock) && PeepholeMinimizeConditions_GetBlockExpression(elseBranch).GetFirstChild().IsAssign() {
		var var_ pkg_rhino.Node = PeepholeMinimizeConditions_GetBlockVar(thenBranch)
		_ = var_
		var elseAssign pkg_rhino.Node = PeepholeMinimizeConditions_GetBlockExpression(elseBranch).GetFirstChild()
		_ = elseAssign
		var name1 pkg_rhino.Node = var_.GetFirstChild()
		_ = name1
		var maybeName2 pkg_rhino.Node = elseAssign.GetFirstChild()
		_ = maybeName2
		if (name1.HasChildren() && maybeName2.IsName()) && func() bool {
			sr_1 := name1.GetString()
			se_2 := maybeName2.GetString()
			return j2gort.CheckNotNullString(sr_1) == se_2
		}() {
			if !(name1.HasOneChild()) {
				j2gort.CheckState(false)
			}
			var thenExpr_3 pkg_rhino.Node = name1.RemoveFirstChild()
			_ = thenExpr_3
			var elseExpr_3 pkg_rhino.Node = elseAssign.GetLastChild().Detach()
			_ = elseExpr_3
			var replacementCond_7 pkg_rhino.Node = this_.priv_PeepholeMinimizeConditions_replaceNode_4212(originalCond, shortCond).Detach()
			_ = replacementCond_7
			var hookNode_2 pkg_rhino.Node = pkg_rhino.IR_Hook(replacementCond_7, thenExpr_3, elseExpr_3).Srcref(n)
			_ = hookNode_2
			var_.Detach()
			name1.AddChildToBack(hookNode_2)
			n.ReplaceWith(var_)
			this_.ReportChangeToEnclosingScope(parent)
			return var_
		}
	} else if (elseBranchIsVar && thenBranchIsExpressionBlock) && PeepholeMinimizeConditions_GetBlockExpression(thenBranch).GetFirstChild().IsAssign() {
		var var__2 pkg_rhino.Node = PeepholeMinimizeConditions_GetBlockVar(elseBranch)
		_ = var__2
		var thenAssign pkg_rhino.Node = PeepholeMinimizeConditions_GetBlockExpression(thenBranch).GetFirstChild()
		_ = thenAssign
		var maybeName1 pkg_rhino.Node = thenAssign.GetFirstChild()
		_ = maybeName1
		var name2 pkg_rhino.Node = var__2.GetFirstChild()
		_ = name2
		if (name2.HasChildren() && maybeName1.IsName()) && func() bool {
			sr_3 := maybeName1.GetString()
			se_4 := name2.GetString()
			return j2gort.CheckNotNullString(sr_3) == se_4
		}() {
			var thenExpr_4 pkg_rhino.Node = thenAssign.GetLastChild().Detach()
			_ = thenExpr_4
			if !(name2.HasOneChild()) {
				j2gort.CheckState(false)
			}
			var elseExpr_4 pkg_rhino.Node = name2.RemoveFirstChild()
			_ = elseExpr_4
			var replacementCond_8 pkg_rhino.Node = this_.priv_PeepholeMinimizeConditions_replaceNode_4212(originalCond, shortCond).Detach()
			_ = replacementCond_8
			var hookNode_3 pkg_rhino.Node = pkg_rhino.IR_Hook(replacementCond_8, thenExpr_4, elseExpr_4).Srcref(n)
			_ = hookNode_3
			var__2.Detach()
			name2.AddChildToBack(hookNode_3)
			n.ReplaceWith(var__2)
			this_.ReportChangeToEnclosingScope(parent)
			return var__2
		}
	}
	this_.priv_PeepholeMinimizeConditions_replaceNode_4212(originalCond, unnegatedCond)
	return n
}

// Java method: PeepholeMinimizeConditions.tryRemoveRepeatedStatements(com.google.javascript.rhino.Node n)
func (this_ *PeepholeMinimizeConditions) priv_PeepholeMinimizeConditions_tryRemoveRepeatedStatements_4210(n pkg_rhino.Node) {
	_ = n
	if !(this_.IsASTNormalized()) {
		return
	}
	if !(n.IsIf()) {
		j2gort.CheckStateMsg(false, j2gort.BoxNode(n))
	}
	var parent pkg_rhino.Node = n.GetParent()
	_ = parent
	if !(NodeUtil_IsStatementBlock(parent)) {
		return
	}
	var cond pkg_rhino.Node = n.GetFirstChild()
	_ = cond
	var trueBranch pkg_rhino.Node = cond.GetNext()
	_ = trueBranch
	var falseBranch pkg_rhino.Node = trueBranch.GetNext()
	_ = falseBranch
	j2gort.CheckNotNullNode(trueBranch)
	j2gort.CheckNotNullNode(falseBranch)
	for {
		var lastTrue pkg_rhino.Node = trueBranch.GetLastChild()
		_ = lastTrue
		var lastFalse pkg_rhino.Node = falseBranch.GetLastChild()
		_ = lastFalse
		if ((lastTrue == 0) || (lastFalse == 0)) || !(this_.AreNodesEqualForInlining(lastTrue, lastFalse)) {
			break
		}
		lastTrue.Detach()
		lastFalse.Detach()
		lastTrue.InsertAfter(n)
		this_.ReportChangeToEnclosingScope(parent)
	}
}

// Java static method: PeepholeMinimizeConditions.isFoldableExpressBlock(com.google.javascript.rhino.Node n)
func PeepholeMinimizeConditions_IsFoldableExpressBlock(n pkg_rhino.Node) bool {
	clinit_PeepholeMinimizeConditions()
	_ = n
	if n.IsBlock() {
		if n.HasOneChild() {
			var maybeExpr pkg_rhino.Node = n.GetFirstChild()
			_ = maybeExpr
			if maybeExpr.IsExprResult() {
				if maybeExpr.GetFirstChild().IsCall() || maybeExpr.GetFirstChild().IsOptChainCall() {
					var calledFn pkg_rhino.Node = maybeExpr.GetFirstFirstChild()
					_ = calledFn
					if calledFn.IsGetElem() || calledFn.IsOptChainGetElem() {
						return false
					} else if (calledFn.IsGetProp() || calledFn.IsOptChainGetProp()) && javalang.StringStartsWith(j2gort.CheckNotNullString(calledFn.GetString()), "on") {
						return false
					}
				}
				return true
			}
			return false
		}
	}
	return false
}

// Java static method: PeepholeMinimizeConditions.getBlockExpression(com.google.javascript.rhino.Node n)
func PeepholeMinimizeConditions_GetBlockExpression(n pkg_rhino.Node) pkg_rhino.Node {
	clinit_PeepholeMinimizeConditions()
	_ = n
	if !(PeepholeMinimizeConditions_IsFoldableExpressBlock(n)) {
		j2gort.CheckState(false)
	}
	return n.GetFirstChild()
}

// Java static method: PeepholeMinimizeConditions.isReturnBlock(com.google.javascript.rhino.Node n)
func PeepholeMinimizeConditions_IsReturnBlock(n pkg_rhino.Node) bool {
	clinit_PeepholeMinimizeConditions()
	_ = n
	if n.IsBlock() {
		if n.HasOneChild() {
			var first pkg_rhino.Node = n.GetFirstChild()
			_ = first
			return first.IsReturn()
		}
	}
	return false
}

// Java static method: PeepholeMinimizeConditions.isReturnExpressBlock(com.google.javascript.rhino.Node n)
func PeepholeMinimizeConditions_IsReturnExpressBlock(n pkg_rhino.Node) bool {
	clinit_PeepholeMinimizeConditions()
	_ = n
	if n.IsBlock() {
		if n.HasOneChild() {
			var first pkg_rhino.Node = n.GetFirstChild()
			_ = first
			if first.IsReturn() {
				return first.HasOneChild()
			}
		}
	}
	return false
}

// Java static method: PeepholeMinimizeConditions.isReturnExpression(com.google.javascript.rhino.Node n)
func PeepholeMinimizeConditions_IsReturnExpression(n pkg_rhino.Node) bool {
	clinit_PeepholeMinimizeConditions()
	_ = n
	if n.IsReturn() {
		return n.HasOneChild()
	}
	return false
}

// Java static method: PeepholeMinimizeConditions.getBlockReturnExpression(com.google.javascript.rhino.Node n)
func PeepholeMinimizeConditions_GetBlockReturnExpression(n pkg_rhino.Node) pkg_rhino.Node {
	clinit_PeepholeMinimizeConditions()
	_ = n
	if !(PeepholeMinimizeConditions_IsReturnExpressBlock(n)) {
		j2gort.CheckState(false)
	}
	return n.GetFirstFirstChild()
}

// Java static method: PeepholeMinimizeConditions.isVarBlock(com.google.javascript.rhino.Node n)
func PeepholeMinimizeConditions_IsVarBlock(n pkg_rhino.Node) bool {
	clinit_PeepholeMinimizeConditions()
	_ = n
	if n.IsBlock() {
		if n.HasOneChild() {
			var first pkg_rhino.Node = n.GetFirstChild()
			_ = first
			if first.IsVar() {
				return first.HasOneChild()
			}
		}
	}
	return false
}

// Java static method: PeepholeMinimizeConditions.getBlockVar(com.google.javascript.rhino.Node n)
func PeepholeMinimizeConditions_GetBlockVar(n pkg_rhino.Node) pkg_rhino.Node {
	clinit_PeepholeMinimizeConditions()
	_ = n
	if !(PeepholeMinimizeConditions_IsVarBlock(n)) {
		j2gort.CheckState(false)
	}
	return n.GetFirstChild()
}

// Java static method: PeepholeMinimizeConditions.consumesDanglingElse(com.google.javascript.rhino.Node n)
func PeepholeMinimizeConditions_ConsumesDanglingElse(n pkg_rhino.Node) bool {
	clinit_PeepholeMinimizeConditions()
	_ = n
	for {
		switch n.TokenSwitchOrd_() {
		case int32(74) /* Token.IF */ :
			if n.GetChildCount() < int32(3) {
				return true
			}
			n = n.GetLastChild()
			continue
		case int32(90) /* Token.BLOCK */ :
			if !(n.HasOneChild()) {
				return false
			}
			n = n.GetLastChild()
			continue
		case int32(85) /* Token.WITH */, int32(78) /* Token.WHILE */, int32(80) /* Token.FOR */, int32(81) /* Token.FOR_IN */ :
			n = n.GetLastChild()
			continue
		default:
			return false
		}
	}
	panic("unreachable")
}

// Java static method: PeepholeMinimizeConditions.isLowerPrecedence(com.google.javascript.rhino.Node n, int precedence)
func PeepholeMinimizeConditions_IsLowerPrecedence(n pkg_rhino.Node, precedence int32) bool {
	clinit_PeepholeMinimizeConditions()
	_ = n
	_ = precedence
	return (NodeUtil_Precedence(n.GetToken()) < precedence)
}

// Java static method: PeepholeMinimizeConditions.isPropertyAssignmentInExpression(com.google.javascript.rhino.Node n)
func PeepholeMinimizeConditions_IsPropertyAssignmentInExpression(n pkg_rhino.Node) bool {
	clinit_PeepholeMinimizeConditions()
	_ = n
	var isPropertyAssignmentInExpressionPredicate j2gort.Ext_com_google_common_base_Predicate = lambdaSingleton_peephole_minimize_conditions_1
	_ = isPropertyAssignmentInExpressionPredicate
	return NodeUtil_Has(n, isPropertyAssignmentInExpressionPredicate, StaticGet_NodeUtil_MATCH_NOT_FUNCTION())
}

// Java method: PeepholeMinimizeConditions.tryMinimizeCondition(com.google.javascript.rhino.Node n)
func (this_ *PeepholeMinimizeConditions) priv_PeepholeMinimizeConditions_tryMinimizeCondition_4211(n pkg_rhino.Node) pkg_rhino.Node {
	_ = n
	n = this_.priv_PeepholeMinimizeConditions_performConditionSubstitutions_4213(n)
	if !(MinimizedCondition_IsMinimizable(n)) {
		return n
	}
	var minCond *MinimizedCondition = minimizedConditionFromConditionNodeWithSlab_(&this_.slab_MinimizedCondition_MeasuredNode, &this_.slab_MinimizedCondition, n)
	_ = minCond
	return this_.priv_PeepholeMinimizeConditions_replaceNode_4212(n, minCond.getMinimizedWithSlab_(&this_.slab_MinimizedCondition_MeasuredNode, StaticGet_MinimizedCondition_MinimizationStyle_PREFER_UNNEGATED()))
}

// Java method: PeepholeMinimizeConditions.replaceNode(com.google.javascript.rhino.Node original, com.google.javascript.jscomp.MinimizedCondition.@org.jspecify.annotations.Nullable MeasuredNode measuredNodeReplacement)
func (this_ *PeepholeMinimizeConditions) priv_PeepholeMinimizeConditions_replaceNode_4212(original pkg_rhino.Node, measuredNodeReplacement *MinimizedCondition_MeasuredNode) pkg_rhino.Node {
	_ = original
	_ = measuredNodeReplacement
	if (measuredNodeReplacement != nil) && measuredNodeReplacement.WillChange(original) {
		var replacement pkg_rhino.Node = measuredNodeReplacement.ApplyTo(original)
		_ = replacement
		this_.ReportChangeToEnclosingScope(replacement)
		return replacement
	}
	return original
}

// Java method: PeepholeMinimizeConditions.performConditionSubstitutions(com.google.javascript.rhino.Node n)
func (this_ *PeepholeMinimizeConditions) priv_PeepholeMinimizeConditions_performConditionSubstitutions_4213(n pkg_rhino.Node) pkg_rhino.Node {
	_ = n
	var parent pkg_rhino.Node = n.GetParent()
	_ = parent
	switch n.TokenSwitchOrd_() {
	case int32(68) /* Token.OR */, int32(69) /* Token.AND */ :
		var left pkg_rhino.Node = n.GetFirstChild()
		_ = left
		var right pkg_rhino.Node = n.GetLastChild()
		_ = right
		left = this_.priv_PeepholeMinimizeConditions_performConditionSubstitutions_4213(left)
		right = this_.priv_PeepholeMinimizeConditions_performConditionSubstitutions_4213(right)
		var rightVal pkg_base.Tri = this_.GetSideEffectFreeBooleanValue(right)
		_ = rightVal
		if this_.GetSideEffectFreeBooleanValue(right) != pkg_base.StaticGet_Tri_UNKNOWN() {
			var type_ pkg_rhino.Token = n.GetToken()
			_ = type_
			var replacement pkg_rhino.Node = pkg_rhino.Node(0)
			_ = replacement
			var rval bool = rightVal.ToBoolean(true)
			_ = rval
			if ((type_ == pkg_rhino.StaticGet_Token_OR()) && !(rval)) || ((type_ == pkg_rhino.StaticGet_Token_AND()) && rval) {
				replacement = left
			} else if !(this_.MayHaveSideEffects(left)) {
				replacement = right
			} else {
				n.DetachChildren()
				replacement = pkg_rhino.IR_Comma(left, right)
			}
			if replacement != 0 {
				n.DetachChildren()
				n.ReplaceWith(replacement)
				this_.ReportChangeToEnclosingScope(parent)
				return replacement
			}
		}
		return n
	case int32(67) /* Token.HOOK */ :
		var condition pkg_rhino.Node = n.GetFirstChild()
		_ = condition
		var trueNode pkg_rhino.Node = n.GetSecondChild()
		_ = trueNode
		var falseNode pkg_rhino.Node = n.GetLastChild()
		_ = falseNode
		trueNode = this_.priv_PeepholeMinimizeConditions_performConditionSubstitutions_4213(trueNode)
		falseNode = this_.priv_PeepholeMinimizeConditions_performConditionSubstitutions_4213(falseNode)
		var replacement_2 pkg_rhino.Node = pkg_rhino.Node(0)
		_ = replacement_2
		var trueNodeVal pkg_base.Tri = this_.GetSideEffectFreeBooleanValue(trueNode)
		_ = trueNodeVal
		var falseNodeVal pkg_base.Tri = this_.GetSideEffectFreeBooleanValue(falseNode)
		_ = falseNodeVal
		if (trueNodeVal == pkg_base.StaticGet_Tri_TRUE()) && (falseNodeVal == pkg_base.StaticGet_Tri_FALSE()) {
			condition.Detach()
			replacement_2 = condition
		} else if (trueNodeVal == pkg_base.StaticGet_Tri_FALSE()) && (falseNodeVal == pkg_base.StaticGet_Tri_TRUE()) {
			condition.Detach()
			replacement_2 = pkg_rhino.IR_Not(condition)
		} else if trueNodeVal == pkg_base.StaticGet_Tri_TRUE() {
			n.DetachChildren()
			replacement_2 = pkg_rhino.IR_Or(condition, falseNode)
		} else if falseNodeVal == pkg_base.StaticGet_Tri_FALSE() {
			n.DetachChildren()
			replacement_2 = pkg_rhino.IR_And(condition, trueNode)
		} else if (!(this_.MayHaveSideEffects(condition)) && !(this_.MayHaveSideEffects(trueNode))) && condition.IsEquivalentTo_1(trueNode) {
			n.DetachChildren()
			replacement_2 = pkg_rhino.IR_Or(trueNode, falseNode)
		}
		if replacement_2 != 0 {
			n.ReplaceWith(replacement_2)
			this_.ReportChangeToEnclosingScope(replacement_2)
			n = replacement_2
		}
		return n
	default:
		var nVal pkg_base.Tri = this_.GetSideEffectFreeBooleanValue(n)
		_ = nVal
		if nVal != pkg_base.StaticGet_Tri_UNKNOWN() {
			var result bool = nVal.ToBoolean(true)
			_ = result
			var equivalentResult int32
			_ = equivalentResult
			if result {
				equivalentResult = int32(1)
			} else {
				equivalentResult = int32(0)
			}
			return this_.priv_PeepholeMinimizeConditions_maybeReplaceChildWithNumber_4214(n, equivalentResult)
		}
		return n
	}
	panic("unreachable switch expression")
}

// Java method: PeepholeMinimizeConditions.maybeReplaceChildWithNumber(com.google.javascript.rhino.Node n, int num)
func (this_ *PeepholeMinimizeConditions) priv_PeepholeMinimizeConditions_maybeReplaceChildWithNumber_4214(n pkg_rhino.Node, num int32) pkg_rhino.Node {
	_ = n
	_ = num
	var newNode pkg_rhino.Node = pkg_rhino.IR_Number(float64(num))
	_ = newNode
	if !(newNode.IsEquivalentTo_1(n)) {
		n.ReplaceWith(newNode)
		this_.ReportChangeToEnclosingScope(newNode)
		this_.MarkFunctionsDeleted(n)
		return newNode
	}
	return n
}
func (this_ *PeepholeMinimizeConditions) ReleaseSlabs_() {
	if this_ == nil {
		return
	}
	this_.slab_MinimizedCondition_MeasuredNode.Release()
	this_.slab_MinimizedCondition.Release()
}
func (this_ *PeepholeMinimizeConditions) EndTraversal() {
	this_.AbstractPeepholeOptimizationBase.EndTraversal()
	this_.ReleaseSlabs_()
}
