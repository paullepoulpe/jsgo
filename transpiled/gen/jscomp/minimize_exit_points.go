/*
 * Copyright 2009 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/MinimizeExitPoints.java (generated mechanically by j2go).

package jscomp

import (
	pkg_base "github.com/paullepoulpe/jsgo/transpiled/gen/base"
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
)

// Java class: com.google.javascript.jscomp.MinimizeExitPoints
type MinimizeExitPoints struct {
	AbstractPeepholeOptimizationBase
}

func (this_ *MinimizeExitPoints) BindSelf_(s any) {
	_ = s
	this_.AbstractPeepholeOptimizationBase.BindSelf_(s)
}
func (this_ *MinimizeExitPoints) IsJavaNil_() bool { return this_ == nil }
func (*MinimizeExitPoints) JavaClassInfo_() (string, string) { return "com.google.javascript.jscomp.MinimizeExitPoints", "MinimizeExitPoints" }
func (this_ *MinimizeExitPoints) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.AbstractPeepholeOptimizationBase.IdentityHash_()
}
func (*MinimizeExitPoints) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.MinimizeExitPoints", "com.google.javascript.jscomp.MinimizeExitPoints", func(v_ any) bool {
		return j2gort.InstanceOf[*MinimizeExitPoints](v_)
	})
}
func (*MinimizeExitPoints) Implements_jscomp_AbstractPeepholeOptimization() {}

func (this_ *MinimizeExitPoints) initFields_MinimizeExitPoints() {
	_ = this_
}

// Java constructor: MinimizeExitPoints()
func NewMinimizeExitPoints() *MinimizeExitPoints {
	this_ := &MinimizeExitPoints{}
	this_.BindSelf_(this_)
	InitMinimizeExitPoints(this_)
	return this_
}

func InitMinimizeExitPoints(this_ *MinimizeExitPoints) {
	InitAbstractPeepholeOptimization(&this_.AbstractPeepholeOptimizationBase)
	this_.initFields_MinimizeExitPoints()
}

// Java method: MinimizeExitPoints.optimizeSubtree(com.google.javascript.rhino.Node n)
func (this_ *MinimizeExitPoints) OptimizeSubtree(n pkg_rhino.Node) pkg_rhino.Node {
	_ = n
	if !(this_.IsASTNormalized()) { j2gort.CheckStateStr(false, "MinimizeExitPoints requires the AST to be normalized") }
	switch n.TokenSwitchOrd_() {
	case int32(92) /* Token.LABEL */:
		this_.TryMinimizeExits(n.GetLastChild(), pkg_rhino.StaticGet_Token_BREAK(), n.GetFirstChild().GetString())
	case int32(80) /* Token.FOR */, int32(81) /* Token.FOR_IN */, int32(112) /* Token.FOR_OF */, int32(113) /* Token.FOR_AWAIT_OF */, int32(78) /* Token.WHILE */:
		this_.TryMinimizeExits(NodeUtil_GetLoopCodeBlock(n), pkg_rhino.StaticGet_Token_CONTINUE(), j2gort.NullString)
	case int32(79) /* Token.DO */:
		this_.TryMinimizeExits(NodeUtil_GetLoopCodeBlock(n), pkg_rhino.StaticGet_Token_CONTINUE(), j2gort.NullString)
		var cond pkg_rhino.Node = NodeUtil_GetConditionExpression(n)
		_ = cond
		if ((this_.GetSideEffectFreeBooleanValue(cond) == pkg_base.StaticGet_Tri_FALSE())) {
			this_.TryMinimizeExits(n.GetFirstChild(), pkg_rhino.StaticGet_Token_BREAK(), j2gort.NullString)
		}
	case int32(90) /* Token.BLOCK */:
		if ((n.HasParent() && n.GetParent().IsFunction())) {
			this_.TryMinimizeExits(n, pkg_rhino.StaticGet_Token_RETURN(), j2gort.NullString)
		}
	case int32(75) /* Token.SWITCH */:
		this_.TryMinimizeSwitchExits(n, pkg_rhino.StaticGet_Token_BREAK(), j2gort.NullString)
	default:
	}
	return n
}

// Java method: MinimizeExitPoints.tryMinimizeExits(com.google.javascript.rhino.Node n, com.google.javascript.rhino.Token exitType, java.lang.@org.jspecify.annotations.Nullable String labelName)
func (this_ *MinimizeExitPoints) TryMinimizeExits(n pkg_rhino.Node, exitType pkg_rhino.Token, labelName string) {
	_ = n
	_ = exitType
	_ = labelName
	if (MinimizeExitPoints_MatchingExitNode(n, exitType, labelName)) {
		this_.ReportChangeToEnclosingScope(n)
		NodeUtil_RemoveChild(n.GetParent(), n)
		return
	}
	if (n.IsIf()) {
		var ifBlock pkg_rhino.Node = n.GetSecondChild()
		_ = ifBlock
		this_.TryMinimizeExits(ifBlock, exitType, labelName)
		var elseBlock pkg_rhino.Node = ifBlock.GetNext()
		_ = elseBlock
		if ((elseBlock != 0)) {
			this_.TryMinimizeExits(elseBlock, exitType, labelName)
		}
		return
	}
	if (n.IsTry()) {
		var tryBlock pkg_rhino.Node = n.GetFirstChild()
		_ = tryBlock
		this_.TryMinimizeExits(tryBlock, exitType, labelName)
		var allCatchNodes pkg_rhino.Node = NodeUtil_GetCatchBlock(n)
		_ = allCatchNodes
		if (NodeUtil_HasCatchHandler(allCatchNodes)) {
			if !(allCatchNodes.HasOneChild()) { j2gort.CheckState(false) }
			var catchNode pkg_rhino.Node = allCatchNodes.GetFirstChild()
			_ = catchNode
			var catchCodeBlock pkg_rhino.Node = catchNode.GetLastChild()
			_ = catchCodeBlock
			this_.TryMinimizeExits(catchCodeBlock, exitType, labelName)
		}
	}
	if (n.IsLabel()) {
		var labelBlock pkg_rhino.Node = n.GetLastChild()
		_ = labelBlock
		this_.TryMinimizeExits(labelBlock, exitType, labelName)
	}
	if ((n.IsSwitch() && (((exitType != pkg_rhino.StaticGet_Token_BREAK()) || (labelName != j2gort.NullString))))) {
		this_.TryMinimizeSwitchExits(n, exitType, labelName)
		return
	}
	if ((!(n.IsBlock()) || !(n.HasChildren()))) {
		return
	}
	{
		var c pkg_rhino.Node = n.GetFirstChild()
		_ = c
		for ; (c != 0); c = c.GetNext() {
			if (c.IsIf()) {
				var ifTree pkg_rhino.Node = c
				_ = ifTree
				var trueBlock pkg_rhino.Node = ifTree.GetSecondChild()
				_ = trueBlock
				var falseBlock pkg_rhino.Node = trueBlock.GetNext()
				_ = falseBlock
				this_.priv_MinimizeExitPoints_tryMinimizeIfBlockExits_3871(trueBlock, falseBlock, ifTree, exitType, labelName)
				trueBlock = ifTree.GetSecondChild()
				falseBlock = trueBlock.GetNext()
				if ((falseBlock != 0)) {
					this_.priv_MinimizeExitPoints_tryMinimizeIfBlockExits_3871(falseBlock, trueBlock, ifTree, exitType, labelName)
				}
			}
			if ((c == n.GetLastChild())) {
				break
			}
		}
	}
	{
		var c_2 pkg_rhino.Node = n.GetLastChild()
		_ = c_2
		for ; (c_2 != 0); c_2 = n.GetLastChild() {
			this_.TryMinimizeExits(c_2, exitType, labelName)
			if ((c_2 == n.GetLastChild())) {
				break
			}
		}
	}
}

// Java method: MinimizeExitPoints.tryMinimizeSwitchExits(com.google.javascript.rhino.Node n, com.google.javascript.rhino.Token exitType, java.lang.@org.jspecify.annotations.Nullable String labelName)
func (this_ *MinimizeExitPoints) TryMinimizeSwitchExits(n pkg_rhino.Node, exitType pkg_rhino.Token, labelName string) {
	_ = n
	_ = exitType
	_ = labelName
	if !(n.IsSwitch()) { j2gort.CheckState(false) }
	var switchBody pkg_rhino.Node = n.GetSecondChild()
	_ = switchBody
	{
		var c pkg_rhino.Node = switchBody.GetFirstChild()
		_ = c
		for ; (c != 0); c = c.GetNext() {
			if ((c != switchBody.GetLastChild())) {
				this_.TryMinimizeSwitchCaseExits(c, exitType, labelName)
			} else {
				this_.TryMinimizeExits(c.GetLastChild(), exitType, labelName)
			}
		}
	}
}

// Java method: MinimizeExitPoints.tryMinimizeSwitchCaseExits(com.google.javascript.rhino.Node n, com.google.javascript.rhino.Token exitType, java.lang.@org.jspecify.annotations.Nullable String labelName)
func (this_ *MinimizeExitPoints) TryMinimizeSwitchCaseExits(n pkg_rhino.Node, exitType pkg_rhino.Token, labelName string) {
	_ = n
	_ = exitType
	_ = labelName
	if !(NodeUtil_IsSwitchCase(n)) { j2gort.CheckState(false) }
	if !((n != n.GetParent().GetLastChild())) { j2gort.CheckState(false) }
	var block pkg_rhino.Node = n.GetLastChild()
	_ = block
	var maybeBreak pkg_rhino.Node = block.GetLastChild()
	_ = maybeBreak
	if ((((maybeBreak == 0) || !(maybeBreak.IsBreak())) || maybeBreak.HasChildren())) {
		return
	}
	var childBeforeBreak pkg_rhino.Node = maybeBreak.GetPrevious()
	_ = childBeforeBreak
	for ((childBeforeBreak != 0)) {
		var c pkg_rhino.Node = childBeforeBreak
		_ = c
		this_.TryMinimizeExits(c, exitType, labelName)
		childBeforeBreak = maybeBreak.GetPrevious()
		if ((c == childBeforeBreak)) {
			break
		}
	}
}

// Java method: MinimizeExitPoints.tryMinimizeIfBlockExits(com.google.javascript.rhino.Node srcBlock, com.google.javascript.rhino.Node destBlock, com.google.javascript.rhino.Node ifNode, com.google.javascript.rhino.Token exitType, java.lang.@org.jspecify.annotations.Nullable String labelName)
func (this_ *MinimizeExitPoints) priv_MinimizeExitPoints_tryMinimizeIfBlockExits_3871(srcBlock pkg_rhino.Node, destBlock pkg_rhino.Node, ifNode pkg_rhino.Node, exitType pkg_rhino.Token, labelName string) {
	_ = srcBlock
	_ = destBlock
	_ = ifNode
	_ = exitType
	_ = labelName
	var exitNodeParent pkg_rhino.Node = pkg_rhino.Node(0)
	_ = exitNodeParent
	var exitNode pkg_rhino.Node = pkg_rhino.Node(0)
	_ = exitNode
	if (srcBlock.IsBlock()) {
		if (!(srcBlock.HasChildren())) {
			return
		}
		exitNodeParent = srcBlock
		exitNode = exitNodeParent.GetLastChild()
	} else {
		exitNodeParent = ifNode
		exitNode = srcBlock
	}
	if (!(MinimizeExitPoints_MatchingExitNode(exitNode, exitType, labelName))) {
		return
	}
	if (!(MinimizeExitPoints_TryConvertAllBlockScopedFollowing(ifNode))) {
		return
	}
	if ((ifNode.GetNext() != 0)) {
		var newDestBlock pkg_rhino.Node = pkg_rhino.IR_Block_0().Srcref(ifNode)
		_ = newDestBlock
		if ((destBlock == 0)) {
			ifNode.AddChildToBack(newDestBlock)
		} else if (destBlock.IsEmpty()) {
			destBlock.ReplaceWith(newDestBlock)
		} else if (destBlock.IsBlock()) {
			newDestBlock = destBlock
		} else {
			destBlock.ReplaceWith(newDestBlock)
			newDestBlock.AddChildToBack(destBlock)
		}
		MinimizeExitPoints_MoveAllFollowing(ifNode, newDestBlock)
		this_.ReportChangeToEnclosingScope(ifNode)
	}
}

// Java static method: MinimizeExitPoints.matchingExitNode(com.google.javascript.rhino.Node n, com.google.javascript.rhino.Token type, java.lang.@org.jspecify.annotations.Nullable String labelName)
func MinimizeExitPoints_MatchingExitNode(n pkg_rhino.Node, type_ pkg_rhino.Token, labelName string) bool {
	_ = n
	_ = type_
	_ = labelName
	if ((n.GetToken() == type_)) {
		if ((type_ == pkg_rhino.StaticGet_Token_RETURN())) {
			return !(n.HasChildren())
		} else {
			if ((labelName == j2gort.NullString)) {
				return !(n.HasChildren())
			} else {
				return (n.HasChildren() && func() bool { sr_1 := labelName; se_2 := n.GetFirstChild().GetString(); return j2gort.CheckNotNullString(sr_1) == se_2 }())
			}
		}
	}
	return false
}

// Java static method: MinimizeExitPoints.moveAllFollowing(com.google.javascript.rhino.Node start, com.google.javascript.rhino.Node destParent)
func MinimizeExitPoints_MoveAllFollowing(start pkg_rhino.Node, destParent pkg_rhino.Node) {
	_ = start
	_ = destParent
	{
		var n pkg_rhino.Node = start.GetNext()
		_ = n
		for ; (n != 0); n = start.GetNext() {
			var isFunctionDeclaration bool = NodeUtil_IsFunctionDeclaration(n)
			_ = isFunctionDeclaration
			n.Detach()
			if (isFunctionDeclaration) {
				destParent.AddChildToFront(n)
			} else {
				destParent.AddChildToBack(n)
			}
		}
	}
}

// Java static method: MinimizeExitPoints.tryConvertAllBlockScopedFollowing(com.google.javascript.rhino.Node start)
func MinimizeExitPoints_TryConvertAllBlockScopedFollowing(start pkg_rhino.Node) bool {
	_ = start
	if (NodeUtil_IsWithinLoop(start)) {
		return !(MinimizeExitPoints_HasBlockScopedVarsFollowing(start))
	}
	{
		var n pkg_rhino.Node = start.GetNext()
		_ = n
		for ; (n != 0); n = n.GetNext() {
			if ((n.IsLet() || n.IsConst())) {
				n.SetToken(pkg_rhino.StaticGet_Token_VAR())
			}
		}
	}
	return true
}

// Java static method: MinimizeExitPoints.hasBlockScopedVarsFollowing(com.google.javascript.rhino.Node start)
func MinimizeExitPoints_HasBlockScopedVarsFollowing(start pkg_rhino.Node) bool {
	_ = start
	{
		var n pkg_rhino.Node = start.GetNext()
		_ = n
		for ; (n != 0); n = n.GetNext() {
			if ((n.IsLet() || n.IsConst())) {
				return true
			}
		}
	}
	return false
}

