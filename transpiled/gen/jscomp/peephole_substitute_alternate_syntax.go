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
// Port of com/google/javascript/jscomp/PeepholeSubstituteAlternateSyntax.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
	"github.com/paullepoulpe/jsgo/transpiled/shims/gson"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.jscomp.PeepholeSubstituteAlternateSyntax
type PeepholeSubstituteAlternateSyntax struct {
	AbstractPeepholeOptimizationBase
	late bool `gson:"late"`
}

func (this_ *PeepholeSubstituteAlternateSyntax) BindSelf_(s any) {
	_ = s
	this_.AbstractPeepholeOptimizationBase.BindSelf_(s)
}
func (this_ *PeepholeSubstituteAlternateSyntax) IsJavaNil_() bool { return this_ == nil }
func (*PeepholeSubstituteAlternateSyntax) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.PeepholeSubstituteAlternateSyntax", "PeepholeSubstituteAlternateSyntax"
}
func (this_ *PeepholeSubstituteAlternateSyntax) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.AbstractPeepholeOptimizationBase.IdentityHash_()
}
func (*PeepholeSubstituteAlternateSyntax) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.PeepholeSubstituteAlternateSyntax", "com.google.javascript.jscomp.PeepholeSubstituteAlternateSyntax", func(v_ any) bool {
		return j2gort.InstanceOf[*PeepholeSubstituteAlternateSyntax](v_)
	})
}
func (*PeepholeSubstituteAlternateSyntax) Implements_jscomp_AbstractPeepholeOptimization() {}

var static_PeepholeSubstituteAlternateSyntax_STRING_SPLIT_OVERHEAD int32

func StaticGet_PeepholeSubstituteAlternateSyntax_STRING_SPLIT_OVERHEAD() int32 {
	clinit_PeepholeSubstituteAlternateSyntax()
	return static_PeepholeSubstituteAlternateSyntax_STRING_SPLIT_OVERHEAD
}
func StaticSet_PeepholeSubstituteAlternateSyntax_STRING_SPLIT_OVERHEAD(v int32) {
	clinit_PeepholeSubstituteAlternateSyntax()
	static_PeepholeSubstituteAlternateSyntax_STRING_SPLIT_OVERHEAD = v
}

var static_PeepholeSubstituteAlternateSyntax_BUILTIN_EXTERNS j2gort.Ext_com_google_common_collect_ImmutableSet

func StaticGet_PeepholeSubstituteAlternateSyntax_BUILTIN_EXTERNS() j2gort.Ext_com_google_common_collect_ImmutableSet {
	clinit_PeepholeSubstituteAlternateSyntax()
	return static_PeepholeSubstituteAlternateSyntax_BUILTIN_EXTERNS
}
func StaticSet_PeepholeSubstituteAlternateSyntax_BUILTIN_EXTERNS(v j2gort.Ext_com_google_common_collect_ImmutableSet) {
	clinit_PeepholeSubstituteAlternateSyntax()
	static_PeepholeSubstituteAlternateSyntax_BUILTIN_EXTERNS = v
}

var static_PeepholeSubstituteAlternateSyntax_STANDARD_OBJECT_CONSTRUCTORS j2gort.Ext_com_google_common_collect_ImmutableSet

func StaticGet_PeepholeSubstituteAlternateSyntax_STANDARD_OBJECT_CONSTRUCTORS() j2gort.Ext_com_google_common_collect_ImmutableSet {
	clinit_PeepholeSubstituteAlternateSyntax()
	return static_PeepholeSubstituteAlternateSyntax_STANDARD_OBJECT_CONSTRUCTORS
}
func StaticSet_PeepholeSubstituteAlternateSyntax_STANDARD_OBJECT_CONSTRUCTORS(v j2gort.Ext_com_google_common_collect_ImmutableSet) {
	clinit_PeepholeSubstituteAlternateSyntax()
	static_PeepholeSubstituteAlternateSyntax_STANDARD_OBJECT_CONSTRUCTORS = v
}

var clinitDone_PeepholeSubstituteAlternateSyntax uint32
var clinitState_PeepholeSubstituteAlternateSyntax j2gort.ClinitState

func clinit_PeepholeSubstituteAlternateSyntax() {
	if atomic.LoadUint32(&clinitDone_PeepholeSubstituteAlternateSyntax) != 0 {
		return
	}
	clinitSlow_PeepholeSubstituteAlternateSyntax()
}

func clinitSlow_PeepholeSubstituteAlternateSyntax() {
	if !j2gort.EnterClinit(&clinitState_PeepholeSubstituteAlternateSyntax, &clinitDone_PeepholeSubstituteAlternateSyntax) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_PeepholeSubstituteAlternateSyntax, &clinitDone_PeepholeSubstituteAlternateSyntax)
	static_PeepholeSubstituteAlternateSyntax_STRING_SPLIT_OVERHEAD = j2gort.JStrLen(".split('.')")
	static_PeepholeSubstituteAlternateSyntax_BUILTIN_EXTERNS = j2gort.Ext_com_google_common_collect_ImmutableSet_Of_5("Object", "Array", "Error", "RegExp", "Math")
	static_PeepholeSubstituteAlternateSyntax_STANDARD_OBJECT_CONSTRUCTORS = j2gort.Ext_com_google_common_collect_ImmutableSet_Of_3("Object", "Array", "Error")
}

func (this_ *PeepholeSubstituteAlternateSyntax) initFields_PeepholeSubstituteAlternateSyntax() {
	_ = this_
}

// Java constructor: PeepholeSubstituteAlternateSyntax(boolean late)
func NewPeepholeSubstituteAlternateSyntax(late bool) *PeepholeSubstituteAlternateSyntax {
	clinit_PeepholeSubstituteAlternateSyntax()
	this_ := &PeepholeSubstituteAlternateSyntax{}
	this_.BindSelf_(this_)
	InitPeepholeSubstituteAlternateSyntax(this_, late)
	return this_
}

func InitPeepholeSubstituteAlternateSyntax(this_ *PeepholeSubstituteAlternateSyntax, late bool) {
	_ = late
	InitAbstractPeepholeOptimization(&this_.AbstractPeepholeOptimizationBase)
	this_.initFields_PeepholeSubstituteAlternateSyntax()
	this_.late = late
}

// Java method: PeepholeSubstituteAlternateSyntax.optimizeSubtree(com.google.javascript.rhino.Node node)
func (this_ *PeepholeSubstituteAlternateSyntax) OptimizeSubtree(node pkg_rhino.Node) pkg_rhino.Node {
	_ = node
	switch node.TokenSwitchOrd_() {
	case int32(59) /* Token.ASSIGN_SUB */ :
		return this_.priv_PeepholeSubstituteAlternateSyntax_reduceSubstractionAssignment_4287(node)
	case int32(39) /* Token.TRUE */, int32(38) /* Token.FALSE */ :
		return this_.priv_PeepholeSubstituteAlternateSyntax_reduceTrueFalse_4288(node)
	case int32(23) /* Token.NEW */ :
		var result pkg_rhino.Node = this_.priv_PeepholeSubstituteAlternateSyntax_tryFoldStandardConstructors_4284(node)
		_ = result
		if !(result.IsCall()) {
			return result
		}
		return this_.priv_PeepholeSubstituteAlternateSyntax_tryReduceCall_4282(result)
	case int32(28) /* Token.CALL */ :
		return this_.priv_PeepholeSubstituteAlternateSyntax_tryReduceCall_4282(node)
	case int32(0) /* Token.RETURN */ :
		return this_.priv_PeepholeSubstituteAlternateSyntax_tryReduceReturn_4283(node)
	case int32(93) /* Token.EXPR_RESULT */ :
		return func() pkg_rhino.Node {
			if node.GetFirstChild().IsComma() {
				return this_.priv_PeepholeSubstituteAlternateSyntax_trySplitComma_4280(node)
			}
			return node
		}()
	case int32(32) /* Token.NAME */ :
		return this_.priv_PeepholeSubstituteAlternateSyntax_tryReplaceUndefined_4281(node)
	case int32(46) /* Token.ARRAYLIT */ :
		return this_.priv_PeepholeSubstituteAlternateSyntax_tryMinimizeArrayLiteral_4289(node)
	case int32(26) /* Token.GETPROP */ :
		return this_.priv_PeepholeSubstituteAlternateSyntax_tryMinimizeWindowRefs_4274(node)
	case int32(131) /* Token.TEMPLATELIT */ :
		return this_.priv_PeepholeSubstituteAlternateSyntax_tryTurnTemplateStringsToStrings_4291(node)
	case int32(69) /* Token.AND */, int32(68) /* Token.OR */, int32(1) /* Token.BITOR */, int32(2) /* Token.BITXOR */, int32(3) /* Token.BITAND */, int32(70) /* Token.COALESCE */ :
		return this_.priv_PeepholeSubstituteAlternateSyntax_tryRotateAssociativeOperator_4276(node)
	case int32(15) /* Token.MUL */ :
		return func() pkg_rhino.Node {
			if this_.MayHaveSideEffects(node) {
				return node
			}
			return this_.priv_PeepholeSubstituteAlternateSyntax_tryRotateCommutativeOperator_4275(node)
		}()
	default:
		return node
	}
	panic("unreachable switch expression")
}

// Java method: PeepholeSubstituteAlternateSyntax.tryMinimizeWindowRefs(com.google.javascript.rhino.Node node)
func (this_ *PeepholeSubstituteAlternateSyntax) priv_PeepholeSubstituteAlternateSyntax_tryMinimizeWindowRefs_4274(node pkg_rhino.Node) pkg_rhino.Node {
	_ = node
	if !(this_.IsASTNormalized()) {
		return node
	}
	if !(node.IsGetProp()) {
		j2gort.CheckArgument(false)
	}
	if node.GetFirstChild().IsName() {
		var nameNode pkg_rhino.Node = node.GetFirstChild()
		_ = nameNode
		{
			var iifeCond_1 bool
			{

				sr_1 := "window"
				se_2 := nameNode.GetString()
				iifeCond_1 = j2gort.CheckNotNullString(sr_1) == se_2
			}
			if iifeCond_1 && j2gort.SetContainsStr(StaticGet_PeepholeSubstituteAlternateSyntax_BUILTIN_EXTERNS(), node.GetString()) {
				var newNameNode pkg_rhino.Node = pkg_rhino.IR_Name(node.GetString())
				_ = newNameNode
				var parentNode pkg_rhino.Node = node.GetParent()
				_ = parentNode
				newNameNode.Srcref(node)
				node.ReplaceWith(newNameNode)
				if parentNode.IsCall() || parentNode.IsOptChainCall() {
					parentNode.PutBooleanProp(pkg_rhino.StaticGet_Node_FREE_CALL(), true)
				}
				this_.ReportChangeToEnclosingScope(parentNode)
				return newNameNode
			}
		}
	}
	return node
}

// Java method: PeepholeSubstituteAlternateSyntax.tryRotateCommutativeOperator(com.google.javascript.rhino.Node n)
func (this_ *PeepholeSubstituteAlternateSyntax) priv_PeepholeSubstituteAlternateSyntax_tryRotateCommutativeOperator_4275(n pkg_rhino.Node) pkg_rhino.Node {
	_ = n
	if !(this_.late) {
		return n
	}
	var rhs pkg_rhino.Node = n.GetLastChild()
	_ = rhs
	var lhs pkg_rhino.Node = n.GetFirstChild()
	_ = lhs
	for (lhs.RawToken_() == n.RawToken_()) && NodeUtil_IsAssociative(n.GetToken()) {
		lhs = lhs.GetFirstChild()
	}
	var precedence int32 = NodeUtil_Precedence(n.GetToken())
	_ = precedence
	var lhsPrecedence int32 = NodeUtil_Precedence(lhs.GetToken())
	_ = lhsPrecedence
	var rhsPrecedence int32 = NodeUtil_Precedence(rhs.GetToken())
	_ = rhsPrecedence
	if (rhsPrecedence == precedence) && (lhsPrecedence != precedence) {
		rhs.Detach()
		lhs.ReplaceWith(rhs)
		n.AddChildToBack(lhs)
		this_.ReportChangeToEnclosingScope(n)
	}
	return n
}

// Java method: PeepholeSubstituteAlternateSyntax.tryRotateAssociativeOperator(com.google.javascript.rhino.Node n)
func (this_ *PeepholeSubstituteAlternateSyntax) priv_PeepholeSubstituteAlternateSyntax_tryRotateAssociativeOperator_4276(n pkg_rhino.Node) pkg_rhino.Node {
	_ = n
	if !(this_.late) {
		return n
	}
	if !(NodeUtil_IsAssociative(n.GetToken())) {
		j2gort.CheckArgument(false)
	}
	var rhs pkg_rhino.Node = n.GetLastChild()
	_ = rhs
	if n.RawToken_() == rhs.RawToken_() {
		var first pkg_rhino.Node = n.RemoveFirstChild()
		_ = first
		var second pkg_rhino.Node = rhs.RemoveFirstChild()
		_ = second
		var third pkg_rhino.Node = rhs.GetLastChild().Detach()
		_ = third
		var newLhs pkg_rhino.Node = pkg_rhino.NewNode_3(n.GetToken(), first, second).SrcrefIfMissing(n)
		_ = newLhs
		var newRoot pkg_rhino.Node = pkg_rhino.NewNode_3(rhs.GetToken(), newLhs, third).SrcrefIfMissing(rhs)
		_ = newRoot
		n.ReplaceWith(newRoot)
		this_.ReportChangeToEnclosingScope(newRoot)
		return newRoot
	} else if NodeUtil_IsCommutative(n.GetToken()) && !(this_.MayHaveSideEffects(n)) {
		return this_.priv_PeepholeSubstituteAlternateSyntax_tryRotateCommutativeOperator_4275(n)
	}
	return n
}

// Java method: PeepholeSubstituteAlternateSyntax.tryFoldSimpleFunctionCall(com.google.javascript.rhino.Node n)
func (this_ *PeepholeSubstituteAlternateSyntax) priv_PeepholeSubstituteAlternateSyntax_tryFoldSimpleFunctionCall_4277(n pkg_rhino.Node) pkg_rhino.Node {
	_ = n
	if !(n.IsCall()) {
		j2gort.CheckStateMsg(false, j2gort.BoxNode(n))
	}
	var callTarget pkg_rhino.Node = n.GetFirstChild()
	_ = callTarget
	if (callTarget == 0) || !(callTarget.IsName()) {
		return n
	}
	var targetName string = callTarget.GetString()
	_ = targetName
	switch j2gort.CheckNotNullString(targetName) {
	case "Boolean":
		var paramCount int32 = (n.GetChildCount() - int32(1))
		_ = paramCount
		if paramCount == int32(1) {
			var value pkg_rhino.Node = n.GetLastChild().Detach()
			_ = value
			var replacement pkg_rhino.Node = pkg_rhino.Node(0)
			_ = replacement
			if NodeUtil_IsBooleanResult(value) {
				replacement = value
			} else {
				replacement = pkg_rhino.IR_Not(pkg_rhino.IR_Not(value).Srcref(n))
			}
			n.ReplaceWith(replacement)
			this_.ReportChangeToEnclosingScope(replacement)
		}
	case "String":
		var value_2 pkg_rhino.Node = callTarget.GetNext()
		_ = value_2
		if ((value_2 != 0) && (value_2.GetNext() == 0)) && NodeUtil_IsImmutableValue(value_2) {
			var addition pkg_rhino.Node = pkg_rhino.IR_Add(pkg_rhino.IR_String("").Srcref(callTarget), value_2.Detach())
			_ = addition
			n.ReplaceWith(addition)
			this_.ReportChangeToEnclosingScope(addition)
			return addition
		}
	default:
	}
	return n
}

// Java method: PeepholeSubstituteAlternateSyntax.tryFoldImmediateCallToBoundFunction(com.google.javascript.rhino.Node n)
func (this_ *PeepholeSubstituteAlternateSyntax) priv_PeepholeSubstituteAlternateSyntax_tryFoldImmediateCallToBoundFunction_4278(n pkg_rhino.Node) pkg_rhino.Node {
	_ = n
	if !(n.IsCall()) {
		j2gort.CheckState(false)
	}
	var callTarget pkg_rhino.Node = n.GetFirstChild()
	_ = callTarget
	var bind *CodingConvention_Bind = this_.GetCodingConvention().DescribeFunctionBind(callTarget, false)
	_ = bind
	if bind != nil {
		bind.target.Detach()
		callTarget.ReplaceWith(bind.target)
		callTarget = bind.target
		PeepholeSubstituteAlternateSyntax_AddParameterAfter(bind.parameters, callTarget)
		if (bind.thisValue != 0) && !(NodeUtil_IsUndefined(bind.thisValue)) {
			var newCallTarget pkg_rhino.Node = pkg_rhino.IR_Getprop_2(callTarget.CloneTree_0(), "call")
			_ = newCallTarget
			this_.MarkNewScopesChanged(newCallTarget)
			callTarget.ReplaceWith(newCallTarget)
			this_.MarkFunctionsDeleted(callTarget)
			bind.thisValue.CloneTree_0().InsertAfter(newCallTarget)
			n.PutBooleanProp(pkg_rhino.StaticGet_Node_FREE_CALL(), false)
		} else {
			n.PutBooleanProp(pkg_rhino.StaticGet_Node_FREE_CALL(), true)
		}
		this_.ReportChangeToEnclosingScope(n)
	}
	return n
}

// Java static method: PeepholeSubstituteAlternateSyntax.addParameterAfter(com.google.javascript.rhino.Node parameterList, com.google.javascript.rhino.Node after)
func PeepholeSubstituteAlternateSyntax_AddParameterAfter(parameterList pkg_rhino.Node, after pkg_rhino.Node) {
	clinit_PeepholeSubstituteAlternateSyntax()
	_ = parameterList
	_ = after
	if parameterList != 0 {
		PeepholeSubstituteAlternateSyntax_AddParameterAfter(parameterList.GetNext(), after)
		parameterList.CloneTree_0().InsertAfter(after)
	}
}

// Java method: PeepholeSubstituteAlternateSyntax.splitComma(com.google.javascript.rhino.Node n, boolean insert, com.google.javascript.rhino.Node insertAfter)
func (this_ *PeepholeSubstituteAlternateSyntax) priv_PeepholeSubstituteAlternateSyntax_splitComma_4279(n pkg_rhino.Node, insert bool, insertAfter pkg_rhino.Node) pkg_rhino.Node {
	_ = n
	_ = insert
	_ = insertAfter
	for n.IsComma() {
		var left pkg_rhino.Node = n.GetFirstChild()
		_ = left
		var right pkg_rhino.Node = n.GetLastChild()
		_ = right
		n.DetachChildren()
		if right.IsComma() {
			this_.priv_PeepholeSubstituteAlternateSyntax_splitComma_4279(right, true, insertAfter)
		} else {
			var newStatement pkg_rhino.Node = pkg_rhino.IR_ExprResult(right)
			_ = newStatement
			newStatement.SrcrefIfMissing(right)
			newStatement.InsertAfter(insertAfter)
		}
		n = left
	}
	if insert {
		var newStatement_2 pkg_rhino.Node = pkg_rhino.IR_ExprResult(n)
		_ = newStatement_2
		newStatement_2.SrcrefIfMissing(n)
		newStatement_2.InsertAfter(insertAfter)
		return newStatement_2
	}
	return n
}

// Java method: PeepholeSubstituteAlternateSyntax.trySplitComma(com.google.javascript.rhino.Node n)
func (this_ *PeepholeSubstituteAlternateSyntax) priv_PeepholeSubstituteAlternateSyntax_trySplitComma_4280(n pkg_rhino.Node) pkg_rhino.Node {
	_ = n
	if this_.late {
		return n
	}
	if !(n.IsExprResult()) {
		j2gort.CheckState(false)
	}
	if n.GetParent().IsLabel() {
		return n
	}
	var comma pkg_rhino.Node = n.GetFirstChild()
	_ = comma
	if !(comma.IsComma()) {
		j2gort.CheckState(false)
	}
	var leftmost pkg_rhino.Node = this_.priv_PeepholeSubstituteAlternateSyntax_splitComma_4279(comma, false, n)
	_ = leftmost
	n.RemoveChildren()
	n.AddChildToFront(leftmost)
	n.Srcref(leftmost)
	this_.ReportChangeToEnclosingScope(leftmost)
	return leftmost
}

// Java method: PeepholeSubstituteAlternateSyntax.tryReplaceUndefined(com.google.javascript.rhino.Node n)
func (this_ *PeepholeSubstituteAlternateSyntax) priv_PeepholeSubstituteAlternateSyntax_tryReplaceUndefined_4281(n pkg_rhino.Node) pkg_rhino.Node {
	_ = n
	if (this_.IsASTNormalized() && NodeUtil_IsUndefined(n)) && !(NodeUtil_IsLValue(n)) {
		var replacement pkg_rhino.Node = NodeUtil_NewUndefinedNode(n)
		_ = replacement
		n.ReplaceWith(replacement)
		this_.ReportChangeToEnclosingScope(replacement)
		return replacement
	}
	return n
}

// Java method: PeepholeSubstituteAlternateSyntax.tryReduceCall(com.google.javascript.rhino.Node node)
func (this_ *PeepholeSubstituteAlternateSyntax) priv_PeepholeSubstituteAlternateSyntax_tryReduceCall_4282(node pkg_rhino.Node) pkg_rhino.Node {
	_ = node
	var result pkg_rhino.Node = this_.priv_PeepholeSubstituteAlternateSyntax_tryFoldLiteralConstructor_4286(node)
	_ = result
	if result == node {
		result = this_.priv_PeepholeSubstituteAlternateSyntax_tryFoldSimpleFunctionCall_4277(node)
		if result == node {
			result = this_.priv_PeepholeSubstituteAlternateSyntax_tryFoldImmediateCallToBoundFunction_4278(node)
		}
	}
	return result
}

// Java method: PeepholeSubstituteAlternateSyntax.tryReduceReturn(com.google.javascript.rhino.Node n)
func (this_ *PeepholeSubstituteAlternateSyntax) priv_PeepholeSubstituteAlternateSyntax_tryReduceReturn_4283(n pkg_rhino.Node) pkg_rhino.Node {
	_ = n
	var result pkg_rhino.Node = n.GetFirstChild()
	_ = result
	if result != 0 {
		switch result.TokenSwitchOrd_() {
		case int32(87) /* Token.VOID */ :
			var operand pkg_rhino.Node = result.GetFirstChild()
			_ = operand
			if !(this_.MayHaveSideEffects(operand)) {
				n.RemoveFirstChild()
				this_.ReportChangeToEnclosingScope(n)
			}
		case int32(32) /* Token.NAME */ :
			var name string = result.GetString()
			_ = name
			if j2gort.CheckNotNullString(name) == "undefined" {
				n.RemoveFirstChild()
				this_.ReportChangeToEnclosingScope(n)
			}
		default:
		}
	}
	return n
}

// Java method: PeepholeSubstituteAlternateSyntax.tryFoldStandardConstructors(com.google.javascript.rhino.Node n)
func (this_ *PeepholeSubstituteAlternateSyntax) priv_PeepholeSubstituteAlternateSyntax_tryFoldStandardConstructors_4284(n pkg_rhino.Node) pkg_rhino.Node {
	_ = n
	if !(n.IsNew()) {
		j2gort.CheckState(false)
	}
	if this_.priv_PeepholeSubstituteAlternateSyntax_canFoldStandardConstructors_4285(n) {
		n.SetToken(pkg_rhino.StaticGet_Token_CALL())
		n.PutBooleanProp(pkg_rhino.StaticGet_Node_FREE_CALL(), true)
		this_.ReportChangeToEnclosingScope(n)
	}
	return n
}

// Java method: PeepholeSubstituteAlternateSyntax.canFoldStandardConstructors(com.google.javascript.rhino.Node n)
func (this_ *PeepholeSubstituteAlternateSyntax) priv_PeepholeSubstituteAlternateSyntax_canFoldStandardConstructors_4285(n pkg_rhino.Node) bool {
	_ = n
	if this_.IsASTNormalized() && n.GetFirstChild().IsName() {
		var className string = n.GetFirstChild().GetString()
		_ = className
		if j2gort.SetContainsStr(StaticGet_PeepholeSubstituteAlternateSyntax_STANDARD_OBJECT_CONSTRUCTORS(), className) {
			return true
		}
		if "RegExp" == className {
			if (n.GetSecondChild() == 0) || n.GetSecondChild().IsStringLit() {
				return true
			}
		}
	}
	return false
}

// Java method: PeepholeSubstituteAlternateSyntax.tryFoldLiteralConstructor(com.google.javascript.rhino.Node n)
func (this_ *PeepholeSubstituteAlternateSyntax) priv_PeepholeSubstituteAlternateSyntax_tryFoldLiteralConstructor_4286(n pkg_rhino.Node) pkg_rhino.Node {
	_ = n
	if !(n.IsCall() || n.IsNew()) {
		j2gort.CheckArgument(false)
	}
	var constructorNameNode pkg_rhino.Node = n.GetFirstChild()
	_ = constructorNameNode
	var newLiteralNode pkg_rhino.Node = pkg_rhino.Node(0)
	_ = newLiteralNode
	if this_.IsASTNormalized() && constructorNameNode.IsName() {
		var className string = constructorNameNode.GetString()
		_ = className
		var constructorHasArgs bool = (constructorNameNode.GetNext() != 0)
		_ = constructorHasArgs
		if ("Object" == className) && !(constructorHasArgs) {
			newLiteralNode = pkg_rhino.IR_Objectlit([]pkg_rhino.Node{})
		} else if "Array" == className {
			var arg0 pkg_rhino.Node = constructorNameNode.GetNext()
			_ = arg0
			var action PeepholeSubstituteAlternateSyntax_FoldArrayAction = PeepholeSubstituteAlternateSyntax_IsSafeToFoldArrayConstructor(arg0)
			_ = action
			if (action == StaticGet_PeepholeSubstituteAlternateSyntax_FoldArrayAction_SAFE_TO_FOLD_WITH_ARGS()) || (action == StaticGet_PeepholeSubstituteAlternateSyntax_FoldArrayAction_SAFE_TO_FOLD_WITHOUT_ARGS()) {
				newLiteralNode = pkg_rhino.IR_Arraylit_1_Node([]pkg_rhino.Node{})
				n.RemoveFirstChild()
				var elements pkg_rhino.Node = n.RemoveChildren()
				_ = elements
				if action == StaticGet_PeepholeSubstituteAlternateSyntax_FoldArrayAction_SAFE_TO_FOLD_WITH_ARGS() {
					newLiteralNode.AddChildrenToFront(elements)
				}
			}
		}
		if newLiteralNode != 0 {
			n.ReplaceWith(newLiteralNode)
			this_.ReportChangeToEnclosingScope(newLiteralNode)
			return newLiteralNode
		}
	}
	return n
}

// Java static method: PeepholeSubstituteAlternateSyntax.isSafeToFoldArrayConstructor(com.google.javascript.rhino.Node arg)
func PeepholeSubstituteAlternateSyntax_IsSafeToFoldArrayConstructor(arg pkg_rhino.Node) PeepholeSubstituteAlternateSyntax_FoldArrayAction {
	clinit_PeepholeSubstituteAlternateSyntax()
	_ = arg
	var action PeepholeSubstituteAlternateSyntax_FoldArrayAction = StaticGet_PeepholeSubstituteAlternateSyntax_FoldArrayAction_NOT_SAFE_TO_FOLD()
	_ = action
	if arg == 0 {
		action = StaticGet_PeepholeSubstituteAlternateSyntax_FoldArrayAction_SAFE_TO_FOLD_WITHOUT_ARGS()
	} else if arg.GetNext() != 0 {
		action = StaticGet_PeepholeSubstituteAlternateSyntax_FoldArrayAction_SAFE_TO_FOLD_WITH_ARGS()
	} else {
		switch arg.TokenSwitchOrd_() {
		case int32(35) /* Token.STRINGLIT */ :
			action = StaticGet_PeepholeSubstituteAlternateSyntax_FoldArrayAction_SAFE_TO_FOLD_WITH_ARGS()
		case int32(33) /* Token.NUMBER */ :
			if arg.GetDouble() == float64(int32(0)) {
				action = StaticGet_PeepholeSubstituteAlternateSyntax_FoldArrayAction_SAFE_TO_FOLD_WITHOUT_ARGS()
			}
		case int32(46) /* Token.ARRAYLIT */ :
			action = StaticGet_PeepholeSubstituteAlternateSyntax_FoldArrayAction_SAFE_TO_FOLD_WITH_ARGS()
		default:
		}
	}
	return action
}

// Java method: PeepholeSubstituteAlternateSyntax.reduceSubstractionAssignment(com.google.javascript.rhino.Node n)
func (this_ *PeepholeSubstituteAlternateSyntax) priv_PeepholeSubstituteAlternateSyntax_reduceSubstractionAssignment_4287(n pkg_rhino.Node) pkg_rhino.Node {
	_ = n
	var right pkg_rhino.Node = n.GetLastChild()
	_ = right
	var isNegative bool = false
	_ = isNegative
	if right.IsNeg() {
		isNegative = true
		right = right.GetOnlyChild()
	}
	if right.IsNumber() && (right.GetDouble() == float64(int32(1))) {
		var left pkg_rhino.Node = n.RemoveFirstChild()
		_ = left
		var newNode pkg_rhino.Node
		_ = newNode
		if isNegative {
			newNode = pkg_rhino.IR_Inc(left, false)
		} else {
			newNode = pkg_rhino.IR_Dec(left, false)
		}
		n.ReplaceWith(newNode)
		this_.ReportChangeToEnclosingScope(newNode)
		return newNode
	}
	return n
}

// Java method: PeepholeSubstituteAlternateSyntax.reduceTrueFalse(com.google.javascript.rhino.Node n)
func (this_ *PeepholeSubstituteAlternateSyntax) priv_PeepholeSubstituteAlternateSyntax_reduceTrueFalse_4288(n pkg_rhino.Node) pkg_rhino.Node {
	_ = n
	if this_.late {
		switch n.GetParent().TokenSwitchOrd_() {
		case int32(4) /* Token.EQ */, int32(8) /* Token.GT */, int32(9) /* Token.GE */, int32(7) /* Token.LE */, int32(6) /* Token.LT */, int32(5) /* Token.NE */ :
			var number pkg_rhino.Node = pkg_rhino.IR_Number(float64(func() int32 {
				if n.IsTrue() {
					return int32(1)
				}
				return int32(0)
			}()))
			_ = number
			n.ReplaceWith(number)
			this_.ReportChangeToEnclosingScope(number)
			return number
		default:
		}
		var not pkg_rhino.Node = pkg_rhino.IR_Not(pkg_rhino.IR_Number(float64(func() int32 {
			if n.IsTrue() {
				return int32(0)
			}
			return int32(1)
		}())))
		_ = not
		not.SrcrefTreeIfMissing(n)
		n.ReplaceWith(not)
		this_.ReportChangeToEnclosingScope(not)
		return not
	}
	return n
}

// Java method: PeepholeSubstituteAlternateSyntax.tryMinimizeArrayLiteral(com.google.javascript.rhino.Node n)
func (this_ *PeepholeSubstituteAlternateSyntax) priv_PeepholeSubstituteAlternateSyntax_tryMinimizeArrayLiteral_4289(n pkg_rhino.Node) pkg_rhino.Node {
	_ = n
	var allStrings bool = true
	_ = allStrings
	{
		var cur pkg_rhino.Node = n.GetFirstChild()
		_ = cur
		for ; cur != 0; cur = cur.GetNext() {
			if !(cur.IsStringLit()) {
				allStrings = false
			}
		}
	}
	if allStrings {
		return this_.priv_PeepholeSubstituteAlternateSyntax_tryMinimizeStringArrayLiteral_4290(n)
	} else {
		return n
	}
}

// Java method: PeepholeSubstituteAlternateSyntax.tryMinimizeStringArrayLiteral(com.google.javascript.rhino.Node n)
func (this_ *PeepholeSubstituteAlternateSyntax) priv_PeepholeSubstituteAlternateSyntax_tryMinimizeStringArrayLiteral_4290(n pkg_rhino.Node) pkg_rhino.Node {
	_ = n
	if !(this_.late) {
		return n
	}
	var numElements int32 = n.GetChildCount()
	_ = numElements
	var saving int32 = ((numElements * int32(2)) - StaticGet_PeepholeSubstituteAlternateSyntax_STRING_SPLIT_OVERHEAD())
	_ = saving
	if saving <= int32(0) {
		return n
	}
	var strings []string = j2gort.NewStringArray(numElements)
	_ = strings
	var idx int32 = int32(0)
	_ = idx
	{
		var cur pkg_rhino.Node = n.GetFirstChild()
		_ = cur
		for ; cur != 0; cur = cur.GetNext() {
			{
				a_3 := strings
				i_4 := func() int32 { old_5 := idx; nw_6 := old_5 + 1; idx = nw_6; return old_5 }()
				a_3[i_4] = cur.GetString()
			}
		}
	}
	var delimiter string = PeepholeSubstituteAlternateSyntax_PickDelimiter(strings)
	_ = delimiter
	if delimiter != j2gort.NullString {
		var template string = j2gort.Ext_com_google_common_base_Joiner_On_1_string(delimiter).Join_1_slice_any(func() []any {
			arr_7 := strings
			if arr_7 == nil {
				return nil
			}
			out_8 := make([]any, len(arr_7))
			for i_, el_9 := range arr_7 {
				out_8[i_] = j2gort.BoxString(el_9)
			}
			return out_8
		}())
		_ = template
		var call pkg_rhino.Node = pkg_rhino.IR_Call(pkg_rhino.IR_Getprop_2(pkg_rhino.IR_String(template), "split"), []pkg_rhino.Node{pkg_rhino.IR_String(("" + j2gort.StrValString(delimiter)))})
		_ = call
		call.SrcrefTreeIfMissing(n)
		n.ReplaceWith(call)
		this_.ReportChangeToEnclosingScope(call)
		return call
	}
	return n
}

// Java method: PeepholeSubstituteAlternateSyntax.tryTurnTemplateStringsToStrings(com.google.javascript.rhino.Node n)
func (this_ *PeepholeSubstituteAlternateSyntax) priv_PeepholeSubstituteAlternateSyntax_tryTurnTemplateStringsToStrings_4291(n pkg_rhino.Node) pkg_rhino.Node {
	_ = n
	if !(n.IsTemplateLit()) {
		j2gort.CheckStateMsg(false, j2gort.BoxNode(n))
	}
	if n.GetParent().IsTaggedTemplateLit() {
		return n
	}
	var string_ string = this_.GetSideEffectFreeStringValue(n)
	_ = string_
	if string_ == j2gort.NullString {
		return n
	}
	var stringNode pkg_rhino.Node = pkg_rhino.IR_String(string_).Srcref(n)
	_ = stringNode
	n.ReplaceWith(stringNode)
	this_.ReportChangeToEnclosingScope(stringNode)
	return stringNode
}

// Java static method: PeepholeSubstituteAlternateSyntax.pickDelimiter(java.lang.String[] strings)
func PeepholeSubstituteAlternateSyntax_PickDelimiter(strings []string) string {
	clinit_PeepholeSubstituteAlternateSyntax()
	_ = strings
	var allLength1 bool = true
	_ = allLength1
	for _, el_10 := range j2gort.CheckNotNullSlice(strings) {
		s := el_10
		_ = s
		if j2gort.JStrLen(s) != int32(1) {
			allLength1 = false
			break
		}
	}
	if allLength1 {
		return ""
	}
	var delimiters []string = []string{" ", ";", ",", "{", "}", j2gort.NullString}
	_ = delimiters
	var i int32 = int32(0)
	_ = i
l_NEXT_DELIMITER:
	for ; delimiters[i] != j2gort.NullString; i = i + 1 {
		for _, el_11 := range j2gort.CheckNotNullSlice(strings) {
			cur := el_11
			_ = cur
			if j2gort.JStrIndexOf(cur, j2gort.CheckNotNullString(delimiters[i])) >= 0 {
				continue l_NEXT_DELIMITER
			}
		}
		break
	}
	return delimiters[i]
}

// Java enum: com.google.javascript.jscomp.PeepholeSubstituteAlternateSyntax.FoldArrayAction
type PeepholeSubstituteAlternateSyntax_FoldArrayAction uint8

const (
	enum_PeepholeSubstituteAlternateSyntax_FoldArrayAction_NOT_SAFE_TO_FOLD          PeepholeSubstituteAlternateSyntax_FoldArrayAction = 1
	enum_PeepholeSubstituteAlternateSyntax_FoldArrayAction_SAFE_TO_FOLD_WITH_ARGS    PeepholeSubstituteAlternateSyntax_FoldArrayAction = 2
	enum_PeepholeSubstituteAlternateSyntax_FoldArrayAction_SAFE_TO_FOLD_WITHOUT_ARGS PeepholeSubstituteAlternateSyntax_FoldArrayAction = 3
)

var enumBaseTable_PeepholeSubstituteAlternateSyntax_FoldArrayAction = [...]j2gort.EnumBase{
	{Name_: "NOT_SAFE_TO_FOLD", Ordinal_: 0},
	{Name_: "SAFE_TO_FOLD_WITH_ARGS", Ordinal_: 1},
	{Name_: "SAFE_TO_FOLD_WITHOUT_ARGS", Ordinal_: 2},
}
var enumValues_PeepholeSubstituteAlternateSyntax_FoldArrayAction = []PeepholeSubstituteAlternateSyntax_FoldArrayAction{1, 2, 3}

func StaticGet_PeepholeSubstituteAlternateSyntax_FoldArrayAction_NOT_SAFE_TO_FOLD() PeepholeSubstituteAlternateSyntax_FoldArrayAction {
	return enum_PeepholeSubstituteAlternateSyntax_FoldArrayAction_NOT_SAFE_TO_FOLD
}
func StaticGet_PeepholeSubstituteAlternateSyntax_FoldArrayAction_SAFE_TO_FOLD_WITH_ARGS() PeepholeSubstituteAlternateSyntax_FoldArrayAction {
	return enum_PeepholeSubstituteAlternateSyntax_FoldArrayAction_SAFE_TO_FOLD_WITH_ARGS
}
func StaticGet_PeepholeSubstituteAlternateSyntax_FoldArrayAction_SAFE_TO_FOLD_WITHOUT_ARGS() PeepholeSubstituteAlternateSyntax_FoldArrayAction {
	return enum_PeepholeSubstituteAlternateSyntax_FoldArrayAction_SAFE_TO_FOLD_WITHOUT_ARGS
}

func PeepholeSubstituteAlternateSyntax_FoldArrayAction_Values() []PeepholeSubstituteAlternateSyntax_FoldArrayAction {
	return j2gort.CloneSlice(enumValues_PeepholeSubstituteAlternateSyntax_FoldArrayAction)
}

func PeepholeSubstituteAlternateSyntax_FoldArrayAction_ValueOf(name string) PeepholeSubstituteAlternateSyntax_FoldArrayAction {
	if name == j2gort.NullString {
		j2gort.ThrowNPE()
	}
	for i_ := range enumBaseTable_PeepholeSubstituteAlternateSyntax_FoldArrayAction {
		if enumBaseTable_PeepholeSubstituteAlternateSyntax_FoldArrayAction[i_].Name_ == name {
			return PeepholeSubstituteAlternateSyntax_FoldArrayAction(i_ + 1)
		}
	}
	panic(javalang.NewIllegalArgumentException("No enum constant com.google.javascript.jscomp.PeepholeSubstituteAlternateSyntax.FoldArrayAction." + name))
}

func (this_ PeepholeSubstituteAlternateSyntax_FoldArrayAction) IsJavaNil_() bool                { return this_ == 0 }
func (PeepholeSubstituteAlternateSyntax_FoldArrayAction) HasIdentityEquals_()                   {}
func (PeepholeSubstituteAlternateSyntax_FoldArrayAction) Implements_j2gort_Ext_java_lang_Enum() {}
func (PeepholeSubstituteAlternateSyntax_FoldArrayAction) Implements_j2gort_Ext_java_lang_constant_Constable() {
}
func (this_ PeepholeSubstituteAlternateSyntax_FoldArrayAction) GetExt_java_lang_EnumBase() *j2gort.Ext_java_lang_EnumBase {
	return &j2gort.Ext_java_lang_EnumBase{Core_: this_.GetEnumBase()}
}
func (PeepholeSubstituteAlternateSyntax_FoldArrayAction) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.PeepholeSubstituteAlternateSyntax$FoldArrayAction", "FoldArrayAction"
}
func (this_ PeepholeSubstituteAlternateSyntax_FoldArrayAction) GetEnumBase() *j2gort.EnumBase {
	if this_ == 0 {
		return nil
	}
	return &enumBaseTable_PeepholeSubstituteAlternateSyntax_FoldArrayAction[this_-1]
}
func (this_ PeepholeSubstituteAlternateSyntax_FoldArrayAction) Ordinal() int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return int32(this_ - 1)
}
func (this_ PeepholeSubstituteAlternateSyntax_FoldArrayAction) EnumOrdinal_() int32 {
	return int32(this_) - 1
}
func (this_ PeepholeSubstituteAlternateSyntax_FoldArrayAction) Name() string {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_PeepholeSubstituteAlternateSyntax_FoldArrayAction[this_-1].Name_
}
func (this_ PeepholeSubstituteAlternateSyntax_FoldArrayAction) String() string {
	if this_ == 0 {
		return "null"
	}
	return enumBaseTable_PeepholeSubstituteAlternateSyntax_FoldArrayAction[this_-1].Name_
}
func (this_ PeepholeSubstituteAlternateSyntax_FoldArrayAction) Equals(other_ any) bool {
	o_, ok_ := other_.(PeepholeSubstituteAlternateSyntax_FoldArrayAction)
	return this_ != 0 && ok_ && this_ == o_
}
func (this_ PeepholeSubstituteAlternateSyntax_FoldArrayAction) HashCode() int32 {
	if this_ == 0 {
		return 0
	}
	return enumBaseTable_PeepholeSubstituteAlternateSyntax_FoldArrayAction[this_-1].HashCode()
}
func (this_ PeepholeSubstituteAlternateSyntax_FoldArrayAction) CompareTo(other_ any) int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_PeepholeSubstituteAlternateSyntax_FoldArrayAction[this_-1].CompareTo(other_)
}
func (this_ PeepholeSubstituteAlternateSyntax_FoldArrayAction) WriteGson(out_ *gson.JsonWriter) {
	if this_ == 0 {
		out_.NullValue()
		return
	}
	out_.ValueString(enumBaseTable_PeepholeSubstituteAlternateSyntax_FoldArrayAction[this_-1].Name_)
}
func init() {
	vals_ := []any{PeepholeSubstituteAlternateSyntax_FoldArrayAction(1), PeepholeSubstituteAlternateSyntax_FoldArrayAction(2), PeepholeSubstituteAlternateSyntax_FoldArrayAction(3)}
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.PeepholeSubstituteAlternateSyntax.FoldArrayAction", vals_)
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.PeepholeSubstituteAlternateSyntax$FoldArrayAction", vals_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.PeepholeSubstituteAlternateSyntax.FoldArrayAction", "com.google.javascript.jscomp.PeepholeSubstituteAlternateSyntax$FoldArrayAction", func(v_ any) bool {
		return j2gort.InstanceOf[PeepholeSubstituteAlternateSyntax_FoldArrayAction](v_)
	})
}
