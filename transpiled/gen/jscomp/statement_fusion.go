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
// Port of com/google/javascript/jscomp/StatementFusion.java (generated mechanically by j2go).

package jscomp

import (
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.jscomp.StatementFusion
type StatementFusion struct {
	AbstractPeepholeOptimizationBase
}

func (this_ *StatementFusion) BindSelf_(s any) {
	_ = s
	this_.AbstractPeepholeOptimizationBase.BindSelf_(s)
}
func (this_ *StatementFusion) IsJavaNil_() bool { return this_ == nil }
func (*StatementFusion) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.StatementFusion", "StatementFusion"
}
func (this_ *StatementFusion) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.AbstractPeepholeOptimizationBase.IdentityHash_()
}
func (*StatementFusion) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.StatementFusion", "com.google.javascript.jscomp.StatementFusion", func(v_ any) bool {
		return j2gort.InstanceOf[*StatementFusion](v_)
	})
}
func (*StatementFusion) Implements_jscomp_AbstractPeepholeOptimization() {}

func (this_ *StatementFusion) initFields_StatementFusion() {
	_ = this_
}

// Java constructor: StatementFusion()
func NewStatementFusion() *StatementFusion {
	this_ := &StatementFusion{}
	this_.BindSelf_(this_)
	InitStatementFusion(this_)
	return this_
}

func InitStatementFusion(this_ *StatementFusion) {
	InitAbstractPeepholeOptimization(&this_.AbstractPeepholeOptimizationBase)
	this_.initFields_StatementFusion()
}

// Java method: StatementFusion.optimizeSubtree(com.google.javascript.rhino.Node n)
func (this_ *StatementFusion) OptimizeSubtree(n pkg_rhino.Node) pkg_rhino.Node {
	_ = n
	if n.GetParent().IsFunction() || !(this_.priv_StatementFusion_canFuseIntoOneStatement_5199(n)) {
		return n
	}
	var start pkg_rhino.Node = n.GetFirstChild()
	_ = start
	var end pkg_rhino.Node = n.GetLastChild()
	_ = end
	var result pkg_rhino.Node = StatementFusion_FuseIntoOneStatement(start, end)
	_ = result
	StatementFusion_FuseExpressionIntoControlFlowStatement(result, n.GetLastChild())
	this_.ReportChangeToEnclosingScope(n)
	return n
}

// Java method: StatementFusion.canFuseIntoOneStatement(com.google.javascript.rhino.Node block)
func (this_ *StatementFusion) priv_StatementFusion_canFuseIntoOneStatement_5199(block pkg_rhino.Node) bool {
	_ = block
	if !(block.IsBlock()) {
		return false
	}
	if !(block.HasChildren()) || block.HasOneChild() {
		return false
	}
	var last pkg_rhino.Node = block.GetLastChild()
	_ = last
	{
		var c pkg_rhino.Node = block.GetFirstChild()
		_ = c
		for ; c != 0; c = c.GetNext() {
			if !(c.IsExprResult()) && (c != last) {
				return false
			}
		}
	}
	return this_.priv_StatementFusion_isFusableControlStatement_5200(last)
}

// Java method: StatementFusion.isFusableControlStatement(com.google.javascript.rhino.Node n)
func (this_ *StatementFusion) priv_StatementFusion_isFusableControlStatement_5200(n pkg_rhino.Node) bool {
	_ = n
	switch n.TokenSwitchOrd_() {
	case int32(74) /* Token.IF */, int32(43) /* Token.THROW */, int32(75) /* Token.SWITCH */, int32(93) /* Token.EXPR_RESULT */ :
		return true
	case int32(0) /* Token.RETURN */ :
		return n.HasChildren()
	case int32(80) /* Token.FOR */ :
		return !(NodeUtil_IsNameDeclaration(n.GetFirstChild()))
	case int32(81) /* Token.FOR_IN */ :
		return !(this_.MayHaveSideEffects(n.GetFirstChild()))
	case int32(92) /* Token.LABEL */ :
		return this_.priv_StatementFusion_isFusableControlStatement_5200(n.GetLastChild())
	case int32(90) /* Token.BLOCK */ :
		return (((this_.IsASTNormalized() || NodeUtil_CanMergeBlock(n)) && !(n.IsSyntheticBlock())) && this_.priv_StatementFusion_isFusableControlStatement_5200(n.GetFirstChild()))
	default:
	}
	return false
}

// Java static method: StatementFusion.fuseIntoOneStatement(com.google.javascript.rhino.Node first, com.google.javascript.rhino.Node last)
func StatementFusion_FuseIntoOneStatement(first pkg_rhino.Node, last pkg_rhino.Node) pkg_rhino.Node {
	_ = first
	_ = last
	if first.GetNext() == last {
		return first
	}
	var commaTree pkg_rhino.Node = first.RemoveFirstChild()
	_ = commaTree
	var next pkg_rhino.Node = pkg_rhino.Node(0)
	_ = next
	{
		var cur pkg_rhino.Node = first.GetNext()
		_ = cur
		for ; cur != last; cur = next {
			commaTree = AstManipulations_FuseExpressions(commaTree, cur.RemoveFirstChild())
			next = cur.GetNext()
			cur.Detach()
		}
	}
	first.AddChildToBack(commaTree)
	return first
}

// Java static method: StatementFusion.fuseExpressionIntoControlFlowStatement(com.google.javascript.rhino.Node before, com.google.javascript.rhino.Node control)
func StatementFusion_FuseExpressionIntoControlFlowStatement(before pkg_rhino.Node, control pkg_rhino.Node) {
	_ = before
	_ = control
	if !(before.IsExprResult()) {
		j2gort.CheckArgumentStr(false, "before must be expression result")
	}
	switch control.TokenSwitchOrd_() {
	case int32(74) /* Token.IF */, int32(0) /* Token.RETURN */, int32(43) /* Token.THROW */, int32(75) /* Token.SWITCH */, int32(93) /* Token.EXPR_RESULT */, int32(80) /* Token.FOR */ :
		before.Detach()
		StatementFusion_FuseExpressionIntoFirstChild(before.RemoveFirstChild(), control)
	case int32(81) /* Token.FOR_IN */ :
		before.Detach()
		StatementFusion_FuseExpressionIntoSecondChild(before.RemoveFirstChild(), control)
	case int32(92) /* Token.LABEL */ :
		StatementFusion_FuseExpressionIntoControlFlowStatement(before, control.GetLastChild())
	case int32(90) /* Token.BLOCK */ :
		StatementFusion_FuseExpressionIntoControlFlowStatement(before, control.GetFirstChild())
	default:
		panic(j2gort.AsIface[javalang.IIllegalStateException](j2gort.ConstructShimThrowable("java.lang.IllegalStateException", "Statement fusion missing.")))
	}
}

// Java static method: StatementFusion.fuseExpressionIntoFirstChild(com.google.javascript.rhino.Node exp, com.google.javascript.rhino.Node stmt)
func StatementFusion_FuseExpressionIntoFirstChild(exp pkg_rhino.Node, stmt pkg_rhino.Node) {
	_ = exp
	_ = stmt
	var val pkg_rhino.Node = stmt.RemoveFirstChild()
	_ = val
	var comma pkg_rhino.Node = AstManipulations_FuseExpressions(exp, val)
	_ = comma
	stmt.AddChildToFront(comma)
}

// Java static method: StatementFusion.fuseExpressionIntoSecondChild(com.google.javascript.rhino.Node exp, com.google.javascript.rhino.Node stmt)
func StatementFusion_FuseExpressionIntoSecondChild(exp pkg_rhino.Node, stmt pkg_rhino.Node) {
	_ = exp
	_ = stmt
	var val pkg_rhino.Node = stmt.GetSecondChild().Detach()
	_ = val
	var comma pkg_rhino.Node = AstManipulations_FuseExpressions(exp, val)
	_ = comma
	comma.InsertAfter(stmt.GetFirstChild())
}
