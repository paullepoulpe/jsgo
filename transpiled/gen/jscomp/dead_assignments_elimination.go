/*
 * Copyright 2008 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/DeadAssignmentsElimination.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	pkg_graph "github.com/paullepoulpe/jsgo/transpiled/gen/graph"
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
	"github.com/paullepoulpe/jsgo/transpiled/shims/gson"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.jscomp.DeadAssignmentsElimination
type DeadAssignmentsElimination struct {
	NodeTraversal_AbstractCfgCallbackBase
	compiler      AbstractCompiler                                             `gson:"compiler"`
	liveness      *LiveVariablesAnalysis                                       `gson:"liveness"`
	functionStack j2gort.Deque[*DeadAssignmentsElimination_BailoutInformation] `gson:"functionStack"`
}

func (this_ *DeadAssignmentsElimination) BindSelf_(s any) {
	_ = s
	this_.NodeTraversal_AbstractCfgCallbackBase.BindSelf_(s)
}
func (this_ *DeadAssignmentsElimination) IsJavaNil_() bool { return this_ == nil }
func (*DeadAssignmentsElimination) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.DeadAssignmentsElimination", "DeadAssignmentsElimination"
}
func (this_ *DeadAssignmentsElimination) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.NodeTraversal_AbstractCfgCallbackBase.IdentityHash_()
}
func (*DeadAssignmentsElimination) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.DeadAssignmentsElimination", "com.google.javascript.jscomp.DeadAssignmentsElimination", func(v_ any) bool {
		return j2gort.InstanceOf[*DeadAssignmentsElimination](v_)
	})
}
func (*DeadAssignmentsElimination) Implements_jscomp_NodeTraversal_AbstractCfgCallback() {}
func (*DeadAssignmentsElimination) Implements_jscomp_CompilerPass()                      {}
func (*DeadAssignmentsElimination) Implements_jscomp_NodeTraversal_ScopedCallback()      {}
func (*DeadAssignmentsElimination) Implements_jscomp_NodeTraversal_Callback()            {}

func (this_ *DeadAssignmentsElimination) initFields_DeadAssignmentsElimination() {
	_ = this_
}

// Java constructor: DeadAssignmentsElimination(com.google.javascript.jscomp.AbstractCompiler compiler)
func NewDeadAssignmentsElimination(compiler AbstractCompiler) *DeadAssignmentsElimination {
	this_ := &DeadAssignmentsElimination{}
	this_.BindSelf_(this_)
	InitDeadAssignmentsElimination(this_, compiler)
	return this_
}

func InitDeadAssignmentsElimination(this_ *DeadAssignmentsElimination, compiler AbstractCompiler) {
	_ = compiler
	InitNodeTraversal_AbstractCfgCallback(&this_.NodeTraversal_AbstractCfgCallbackBase)
	this_.initFields_DeadAssignmentsElimination()
	this_.compiler = compiler
	this_.functionStack = j2gort.Deque[*DeadAssignmentsElimination_BailoutInformation]{}
}

// Java method: DeadAssignmentsElimination.process(com.google.javascript.rhino.Node externs, com.google.javascript.rhino.Node root)
func (this_ *DeadAssignmentsElimination) Process_2(externs pkg_rhino.Node, root pkg_rhino.Node) {
	_ = externs
	_ = root
	j2gort.CheckNotNullNode(externs)
	j2gort.CheckNotNullNode(root)
	if !(this_.compiler.GetLifeCycleStage().IsNormalized()) {
		j2gort.CheckState(false)
	}
	NodeTraversal_Traverse(this_.compiler, root, this_)
}

// Java method: DeadAssignmentsElimination.visit(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *DeadAssignmentsElimination) Visit_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = n
	_ = parent
	if this_.functionStack.IsEmpty() {
		return
	}
	if n.IsFunction() {
		{
			r_1 := this_.functionStack.PeekFirst()
			r_1.containsFunction = true
		}
	} else if this_.IsRemovableAssign(n) {
		{
			r_2 := this_.functionStack.PeekFirst()
			r_2.containsRemovableAssign = true
		}
	}
}

// Java method: DeadAssignmentsElimination.enterScopeWithCfg(com.google.javascript.jscomp.NodeTraversal t)
func (this_ *DeadAssignmentsElimination) EnterScopeWithCfg(t *NodeTraversal) {
	_ = t
	if t.InFunctionBlockScope() {
		this_.functionStack.AddFirst(NewDeadAssignmentsElimination_BailoutInformation())
	}
}

// Java method: DeadAssignmentsElimination.exitScopeWithCfg(com.google.javascript.jscomp.NodeTraversal t)
func (this_ *DeadAssignmentsElimination) ExitScopeWithCfg(t *NodeTraversal) {
	_ = t
	if t.InFunctionBlockScope() {
		this_.priv_DeadAssignmentsElimination_eliminateDeadAssignments_2331(t)
		this_.functionStack.RemoveFirst()
	}
}

// Java method: DeadAssignmentsElimination.eliminateDeadAssignments(com.google.javascript.jscomp.NodeTraversal t)
func (this_ *DeadAssignmentsElimination) priv_DeadAssignmentsElimination_eliminateDeadAssignments_2331(t *NodeTraversal) {
	_ = t
	if !(t.InFunctionBlockScope()) {
		j2gort.CheckArgument(false)
	}
	if !(!(this_.functionStack.IsEmpty())) {
		j2gort.CheckState(false)
	}
	if !(this_.compiler.HasScopeChanged(t.GetScopeRoot().GetParent())) {
		return
	}
	var currentFunction *DeadAssignmentsElimination_BailoutInformation = this_.functionStack.PeekFirst()
	_ = currentFunction
	if currentFunction.containsFunction {
		return
	}
	if !(currentFunction.containsRemovableAssign) {
		return
	}
	var blockScope *Scope = t.GetScope()
	_ = blockScope
	var functionScope *Scope = j2gort.CastPtr[Scope](blockScope.GetParent())
	_ = functionScope
	if int32(100) < (blockScope.GetVarCount() + functionScope.GetVarCount()) {
		return
	}
	var cfg ControlFlowGraph[pkg_rhino.Node] = this_.GetControlFlowGraph(this_.compiler)
	_ = cfg
	clinit_SyntacticScopeCreator()
	var val_3 SyntacticScopeCreator
	InitSyntacticScopeCreator_1(&val_3, this_.compiler)
	var scopeCreator *SyntacticScopeCreator = &val_3
	_ = scopeCreator
	defer func() {
		if this_.liveness != nil {
			this_.liveness.ReleaseSlabs_()
		}
		this_.liveness = nil
		scopeCreator.ReleaseSlabs_()
	}()

	var allVarsDeclaredInFunction *NodeUtil_AllVarsDeclaredInFunction = NodeUtil_GetAllVarsDeclaredInFunction(this_.compiler, scopeCreator, functionScope)
	_ = allVarsDeclaredInFunction
	this_.liveness = NewLiveVariablesAnalysis(cfg, functionScope, blockScope, this_.compiler, scopeCreator, allVarsDeclaredInFunction)
	this_.liveness.Analyze()
	var allVarsInFn j2gort.Ext_java_util_Map = this_.liveness.GetAllVariables()
	_ = allVarsInFn
	this_.priv_DeadAssignmentsElimination_tryRemoveDeadAssignments_2333(t, cfg, allVarsInFn)
}

// Java method: DeadAssignmentsElimination.isRemovableAssign(com.google.javascript.rhino.Node n)
func (this_ *DeadAssignmentsElimination) IsRemovableAssign(n pkg_rhino.Node) bool {
	_ = n
	return (((NodeUtil_IsAssignmentOp(n) && n.GetFirstChild().IsName()) || n.IsInc()) || n.IsDec())
}

// Java method: DeadAssignmentsElimination.tryRemoveDeadAssignments(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.jscomp.ControlFlowGraph<com.google.javascript.rhino.Node> cfg, java.util.Map<java.lang.String,com.google.javascript.jscomp.Var> allVarsInFn)
func (this_ *DeadAssignmentsElimination) priv_DeadAssignmentsElimination_tryRemoveDeadAssignments_2333(t *NodeTraversal, cfg ControlFlowGraph[pkg_rhino.Node], allVarsInFn j2gort.Ext_java_util_Map) {
	_ = t
	_ = cfg
	_ = allVarsInFn
	var nodes j2gort.Ext_java_lang_Iterable = j2gort.AsIface[j2gort.Ext_java_lang_Iterable](cfg.GetNodes())
	_ = nodes
	{
		cur_4 := j2gort.NewIterCursor(j2gort.CheckNotNull(nodes))
		for cur_4.HasNext() {
			rawElem_5 := cur_4.Next()
			cfgNode := j2gort.AsIface[pkg_graph.DiGraph_DiGraphNode[pkg_rhino.Node, ControlFlowGraph_Branch]](rawElem_5)
			_ = cfgNode
			var state *DataFlowAnalysis_LinearFlowState = j2gort.CastPtr[DataFlowAnalysis_LinearFlowState](cfgNode.GetAnnotation())
			_ = state
			var n pkg_rhino.Node = cfgNode.GetValue()
			_ = n
			if n == 0 {
				continue
			}
			switch n.TokenSwitchOrd_() {
			case int32(74) /* Token.IF */, int32(78) /* Token.WHILE */, int32(79) /* Token.DO */ :
				this_.priv_DeadAssignmentsElimination_tryRemoveAssignment_2334(t, NodeUtil_GetConditionExpression(n), state, allVarsInFn)
				continue
			case int32(80) /* Token.FOR */, int32(81) /* Token.FOR_IN */, int32(112) /* Token.FOR_OF */, int32(113) /* Token.FOR_AWAIT_OF */ :
				if n.IsVanillaFor() {
					this_.priv_DeadAssignmentsElimination_tryRemoveAssignment_2334(t, NodeUtil_GetConditionExpression(n), state, allVarsInFn)
				}
				continue
			case int32(75) /* Token.SWITCH */, int32(76) /* Token.CASE */, int32(0) /* Token.RETURN */ :
				if n.HasChildren() {
					this_.priv_DeadAssignmentsElimination_tryRemoveAssignment_2334(t, n.GetFirstChild(), state, allVarsInFn)
				}
				continue
			default:
			}
			this_.priv_DeadAssignmentsElimination_tryRemoveAssignment_2334(t, n, state, allVarsInFn)
		}
	}
}

// Java method: DeadAssignmentsElimination.tryRemoveAssignment(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.jscomp.DataFlowAnalysis.LinearFlowState<com.google.javascript.jscomp.LiveVariablesAnalysis.LiveVariableLattice> state, java.util.Map<java.lang.String,com.google.javascript.jscomp.Var> allVarsInFn)
func (this_ *DeadAssignmentsElimination) priv_DeadAssignmentsElimination_tryRemoveAssignment_2334(t *NodeTraversal, n pkg_rhino.Node, state *DataFlowAnalysis_LinearFlowState, allVarsInFn j2gort.Ext_java_util_Map) {
	_ = t
	_ = n
	_ = state
	_ = allVarsInFn
	this_.priv_DeadAssignmentsElimination_tryRemoveAssignment_2335(t, n, n, state, allVarsInFn)
}

// Java method: DeadAssignmentsElimination.tryRemoveAssignment(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node exprRoot, com.google.javascript.jscomp.DataFlowAnalysis.LinearFlowState<com.google.javascript.jscomp.LiveVariablesAnalysis.LiveVariableLattice> state, java.util.Map<java.lang.String,com.google.javascript.jscomp.Var> allVarsInFn)
func (this_ *DeadAssignmentsElimination) priv_DeadAssignmentsElimination_tryRemoveAssignment_2335(t *NodeTraversal, n pkg_rhino.Node, exprRoot pkg_rhino.Node, state *DataFlowAnalysis_LinearFlowState, allVarsInFn j2gort.Ext_java_util_Map) {
	_ = t
	_ = n
	_ = exprRoot
	_ = state
	_ = allVarsInFn
	var parent pkg_rhino.Node = n.GetParent()
	_ = parent
	var isDeclarationNode bool = NodeUtil_IsNameDeclaration(parent)
	_ = isDeclarationNode
	if ((NodeUtil_IsAssignmentOp(n) || n.IsInc()) || n.IsDec()) || isDeclarationNode {
		if parent.IsConst() {
			return
		}
		var lhs pkg_rhino.Node
		_ = lhs
		if isDeclarationNode {
			lhs = n
		} else {
			lhs = n.GetFirstChild()
		}
		var rhs pkg_rhino.Node = NodeUtil_GetRValueOfLValue(lhs)
		_ = rhs
		if rhs != 0 {
			this_.priv_DeadAssignmentsElimination_tryRemoveAssignment_2335(t, rhs, exprRoot, state, allVarsInFn)
			rhs = NodeUtil_GetRValueOfLValue(lhs)
		}
		if isDeclarationNode && (lhs.GetNext() != 0) {
			this_.priv_DeadAssignmentsElimination_tryRemoveAssignment_2335(t, lhs.GetNext(), exprRoot, state, allVarsInFn)
		}
		if isDeclarationNode && ((rhs == 0) || NodeUtil_IsAnyFor(parent.GetParent())) {
			return
		}
		if !(lhs.IsName()) {
			return
		}
		var name string = lhs.GetString()
		_ = name
		var scope *Scope = t.GetScope()
		_ = scope
		if !(scope.IsFunctionBlockScope() || scope.IsBlockScope()) {
			j2gort.CheckState(false)
		}
		if !(j2gort.MapContainsKeyStr(allVarsInFn, name)) {
			return
		}
		var var_ *Var = j2gort.CastPtr[Var](j2gort.MapGetStr(allVarsInFn, name))
		_ = var_
		if this_.liveness.GetEscapedLocals().Contains(j2gort.PtrToAny(var_)) {
			return
		}
		if (((rhs != 0) && rhs.IsName()) && func() bool {
			sr_6 := rhs.GetString()
			se_7 := var_.GetName()
			return j2gort.CheckNotNullString(sr_6) == se_7
		}()) && n.IsAssign() {
			rhs.Detach()
			n.ReplaceWith(rhs)
			this_.compiler.ReportChangeToEnclosingScope(rhs)
			return
		}
		var index int32 = this_.liveness.GetVarIndex(var_.GetName())
		_ = index
		if j2gort.CastPtr[LiveVariablesAnalysis_LiveVariableLattice](state.GetOut()).IsLive(index) {
			return
		}
		if j2gort.CastPtr[LiveVariablesAnalysis_LiveVariableLattice](state.GetIn()).IsLive(index) && this_.priv_DeadAssignmentsElimination_isVariableStillLiveWithinExpression_2336(n, exprRoot, var_.GetName()) {
			return
		}
		if n.IsAssign() {
			rhs.Detach()
			n.ReplaceWith(rhs)
		} else if NodeUtil_IsAssignmentOp(n) {
			rhs.Detach()
			lhs.Detach()
			var op pkg_rhino.Node = pkg_rhino.NewNode_3(NodeUtil_GetOpFromAssignmentOp(n), lhs, rhs)
			_ = op
			n.ReplaceWith(op)
		} else if n.IsInc() || n.IsDec() {
			if parent.IsExprResult() {
				n.ReplaceWith(pkg_rhino.IR_VoidNode(pkg_rhino.IR_Number(float64(int32(0))).Srcref(n)))
			} else if n.IsComma() && (n != parent.GetLastChild()) {
				n.Detach()
			} else if parent.IsVanillaFor() && (NodeUtil_GetConditionExpression(parent) != n) {
				n.ReplaceWith(pkg_rhino.IR_Empty())
			} else {
				return
			}
		} else if isDeclarationNode {
			rhs.Detach()
			pkg_rhino.IR_ExprResult(rhs).InsertAfter(parent)
			rhs.GetParent().Srcref(rhs)
		} else {
			panic(j2gort.AsIface[javalang.IIllegalStateException](j2gort.ConstructShimThrowable("java.lang.IllegalStateException", "Unknown statement")))
		}
		this_.compiler.ReportChangeToEnclosingScope(parent)
		return
	} else {
		{
			var c pkg_rhino.Node = n.GetFirstChild()
			_ = c
			for c != 0 {
				var next pkg_rhino.Node = c.GetNext()
				_ = next
				if !(ControlFlowGraph_IsEnteringNewCfgNode(c)) {
					this_.priv_DeadAssignmentsElimination_tryRemoveAssignment_2335(t, c, exprRoot, state, allVarsInFn)
				}
				c = next
			}
		}
		return
	}
}

// Java method: DeadAssignmentsElimination.isVariableStillLiveWithinExpression(com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node exprRoot, java.lang.String variable)
func (this_ *DeadAssignmentsElimination) priv_DeadAssignmentsElimination_isVariableStillLiveWithinExpression_2336(n pkg_rhino.Node, exprRoot pkg_rhino.Node, variable string) bool {
	_ = n
	_ = exprRoot
	_ = variable
	for n != exprRoot {
		var state DeadAssignmentsElimination_VariableLiveness = StaticGet_DeadAssignmentsElimination_VariableLiveness_MAYBE_LIVE()
		_ = state
		switch n.GetParent().TokenSwitchOrd_() {
		case int32(68) /* Token.OR */, int32(69) /* Token.AND */, int32(70) /* Token.COALESCE */, int32(29) /* Token.OPTCHAIN_GETPROP */, int32(30) /* Token.OPTCHAIN_GETELEM */, int32(31) /* Token.OPTCHAIN_CALL */ :
			{
				var sibling pkg_rhino.Node = n.GetNext()
				_ = sibling
				for ; sibling != 0; sibling = sibling.GetNext() {
					state = this_.priv_DeadAssignmentsElimination_isVariableReadBeforeKill_2337(sibling, variable)
					if state != StaticGet_DeadAssignmentsElimination_VariableLiveness_MAYBE_LIVE() {
						break
					}
				}
			}
			if (n == n.GetParent().GetFirstChild()) && (state == StaticGet_DeadAssignmentsElimination_VariableLiveness_KILL()) {
				state = StaticGet_DeadAssignmentsElimination_VariableLiveness_MAYBE_LIVE()
			}
		case int32(67) /* Token.HOOK */ :
			if (n.GetNext() != 0) && (n.GetNext().GetNext() != 0) {
				state = this_.priv_DeadAssignmentsElimination_checkHookBranchReadBeforeKill_2338(n.GetNext(), n.GetNext().GetNext(), variable)
			}
		default:
			{
				var sibling_2 pkg_rhino.Node = n.GetNext()
				_ = sibling_2
				for ; sibling_2 != 0; sibling_2 = sibling_2.GetNext() {
					state = this_.priv_DeadAssignmentsElimination_isVariableReadBeforeKill_2337(sibling_2, variable)
					if state != StaticGet_DeadAssignmentsElimination_VariableLiveness_MAYBE_LIVE() {
						break
					}
				}
			}
		}
		if state == StaticGet_DeadAssignmentsElimination_VariableLiveness_READ() {
			return true
		} else if state == StaticGet_DeadAssignmentsElimination_VariableLiveness_KILL() {
			return false
		}
		n = n.GetParent()
	}
	return false
}

// Java method: DeadAssignmentsElimination.isVariableReadBeforeKill(com.google.javascript.rhino.Node n, java.lang.String variable)
func (this_ *DeadAssignmentsElimination) priv_DeadAssignmentsElimination_isVariableReadBeforeKill_2337(n pkg_rhino.Node, variable string) DeadAssignmentsElimination_VariableLiveness {
	_ = n
	_ = variable
	if ControlFlowGraph_IsEnteringNewCfgNode(n) {
		return StaticGet_DeadAssignmentsElimination_VariableLiveness_MAYBE_LIVE()
	}
	if n.IsName() && func() bool { sr_8 := variable; se_9 := n.GetString(); return j2gort.CheckNotNullString(sr_8) == se_9 }() {
		if NodeUtil_IsNameDeclOrSimpleAssignLhs(n, n.GetParent()) {
			if !(n.GetParent().IsAssign()) {
				j2gort.CheckStateMsg(false, j2gort.BoxNode(n.GetParent()))
			}
			var rhs pkg_rhino.Node = n.GetNext()
			_ = rhs
			var state DeadAssignmentsElimination_VariableLiveness = this_.priv_DeadAssignmentsElimination_isVariableReadBeforeKill_2337(rhs, variable)
			_ = state
			if state == StaticGet_DeadAssignmentsElimination_VariableLiveness_READ() {
				return state
			}
			return StaticGet_DeadAssignmentsElimination_VariableLiveness_KILL()
		} else {
			return StaticGet_DeadAssignmentsElimination_VariableLiveness_READ()
		}
	}
	switch n.TokenSwitchOrd_() {
	case int32(68) /* Token.OR */, int32(69) /* Token.AND */, int32(70) /* Token.COALESCE */, int32(29) /* Token.OPTCHAIN_GETPROP */, int32(30) /* Token.OPTCHAIN_GETELEM */ :
		var v1 DeadAssignmentsElimination_VariableLiveness = this_.priv_DeadAssignmentsElimination_isVariableReadBeforeKill_2337(n.GetFirstChild(), variable)
		_ = v1
		var v2 DeadAssignmentsElimination_VariableLiveness = this_.priv_DeadAssignmentsElimination_isVariableReadBeforeKill_2337(n.GetLastChild(), variable)
		_ = v2
		if v1 != StaticGet_DeadAssignmentsElimination_VariableLiveness_MAYBE_LIVE() {
			return v1
		} else if v2 == StaticGet_DeadAssignmentsElimination_VariableLiveness_READ() {
			return StaticGet_DeadAssignmentsElimination_VariableLiveness_READ()
		} else {
			return StaticGet_DeadAssignmentsElimination_VariableLiveness_MAYBE_LIVE()
		}
	case int32(31) /* Token.OPTCHAIN_CALL */ :
		var v1_2 DeadAssignmentsElimination_VariableLiveness = this_.priv_DeadAssignmentsElimination_isVariableReadBeforeKill_2337(n.GetFirstChild(), variable)
		_ = v1_2
		if v1_2 != StaticGet_DeadAssignmentsElimination_VariableLiveness_MAYBE_LIVE() {
			return v1_2
		}
		{
			var c pkg_rhino.Node = n.GetSecondChild()
			_ = c
			for ; c != 0; c = c.GetNext() {
				var state_2 DeadAssignmentsElimination_VariableLiveness = this_.priv_DeadAssignmentsElimination_isVariableReadBeforeKill_2337(c, variable)
				_ = state_2
				if state_2 != StaticGet_DeadAssignmentsElimination_VariableLiveness_MAYBE_LIVE() {
					if state_2 == StaticGet_DeadAssignmentsElimination_VariableLiveness_READ() {
						return StaticGet_DeadAssignmentsElimination_VariableLiveness_READ()
					} else {
						return StaticGet_DeadAssignmentsElimination_VariableLiveness_MAYBE_LIVE()
					}
				}
			}
		}
		return StaticGet_DeadAssignmentsElimination_VariableLiveness_MAYBE_LIVE()
	case int32(67) /* Token.HOOK */ :
		var first DeadAssignmentsElimination_VariableLiveness = this_.priv_DeadAssignmentsElimination_isVariableReadBeforeKill_2337(n.GetFirstChild(), variable)
		_ = first
		if first != StaticGet_DeadAssignmentsElimination_VariableLiveness_MAYBE_LIVE() {
			return first
		}
		return this_.priv_DeadAssignmentsElimination_checkHookBranchReadBeforeKill_2338(n.GetSecondChild(), n.GetLastChild(), variable)
	default:
		{
			var child pkg_rhino.Node = n.GetFirstChild()
			_ = child
			for ; child != 0; child = child.GetNext() {
				var state_3 DeadAssignmentsElimination_VariableLiveness = this_.priv_DeadAssignmentsElimination_isVariableReadBeforeKill_2337(child, variable)
				_ = state_3
				if state_3 != StaticGet_DeadAssignmentsElimination_VariableLiveness_MAYBE_LIVE() {
					return state_3
				}
			}
		}
	}
	return StaticGet_DeadAssignmentsElimination_VariableLiveness_MAYBE_LIVE()
}

// Java method: DeadAssignmentsElimination.checkHookBranchReadBeforeKill(com.google.javascript.rhino.Node trueCase, com.google.javascript.rhino.Node falseCase, java.lang.String variable)
func (this_ *DeadAssignmentsElimination) priv_DeadAssignmentsElimination_checkHookBranchReadBeforeKill_2338(trueCase pkg_rhino.Node, falseCase pkg_rhino.Node, variable string) DeadAssignmentsElimination_VariableLiveness {
	_ = trueCase
	_ = falseCase
	_ = variable
	var v1 DeadAssignmentsElimination_VariableLiveness = this_.priv_DeadAssignmentsElimination_isVariableReadBeforeKill_2337(trueCase, variable)
	_ = v1
	var v2 DeadAssignmentsElimination_VariableLiveness = this_.priv_DeadAssignmentsElimination_isVariableReadBeforeKill_2337(falseCase, variable)
	_ = v2
	if (v1 == StaticGet_DeadAssignmentsElimination_VariableLiveness_READ()) || (v2 == StaticGet_DeadAssignmentsElimination_VariableLiveness_READ()) {
		return StaticGet_DeadAssignmentsElimination_VariableLiveness_READ()
	} else if (v1 == StaticGet_DeadAssignmentsElimination_VariableLiveness_KILL()) && (v2 == StaticGet_DeadAssignmentsElimination_VariableLiveness_KILL()) {
		return StaticGet_DeadAssignmentsElimination_VariableLiveness_KILL()
	} else {
		return StaticGet_DeadAssignmentsElimination_VariableLiveness_MAYBE_LIVE()
	}
}

// Java class: com.google.javascript.jscomp.DeadAssignmentsElimination.BailoutInformation
type DeadAssignmentsElimination_BailoutInformation struct {
	idHash_                 atomic.Int32
	containsFunction        bool `gson:"containsFunction"`
	containsRemovableAssign bool `gson:"containsRemovableAssign"`
}

func (this_ *DeadAssignmentsElimination_BailoutInformation) BindSelf_(s any) {
	_ = s
}
func (this_ *DeadAssignmentsElimination_BailoutInformation) IsJavaNil_() bool { return this_ == nil }
func (*DeadAssignmentsElimination_BailoutInformation) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.DeadAssignmentsElimination$BailoutInformation", "BailoutInformation"
}
func (this_ *DeadAssignmentsElimination_BailoutInformation) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*DeadAssignmentsElimination_BailoutInformation) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.DeadAssignmentsElimination.BailoutInformation", "com.google.javascript.jscomp.DeadAssignmentsElimination$BailoutInformation", func(v_ any) bool {
		return j2gort.InstanceOf[*DeadAssignmentsElimination_BailoutInformation](v_)
	})
}

func (this_ *DeadAssignmentsElimination_BailoutInformation) initFields_DeadAssignmentsElimination_BailoutInformation() {
	_ = this_
}

// Java constructor: BailoutInformation()
func NewDeadAssignmentsElimination_BailoutInformation() *DeadAssignmentsElimination_BailoutInformation {
	this_ := &DeadAssignmentsElimination_BailoutInformation{}
	InitDeadAssignmentsElimination_BailoutInformation(this_)
	return this_
}

func InitDeadAssignmentsElimination_BailoutInformation(this_ *DeadAssignmentsElimination_BailoutInformation) {
	this_.initFields_DeadAssignmentsElimination_BailoutInformation()
}

// Java enum: com.google.javascript.jscomp.DeadAssignmentsElimination.VariableLiveness
type DeadAssignmentsElimination_VariableLiveness uint8

const (
	enum_DeadAssignmentsElimination_VariableLiveness_MAYBE_LIVE DeadAssignmentsElimination_VariableLiveness = 1
	enum_DeadAssignmentsElimination_VariableLiveness_READ       DeadAssignmentsElimination_VariableLiveness = 2
	enum_DeadAssignmentsElimination_VariableLiveness_KILL       DeadAssignmentsElimination_VariableLiveness = 3
)

var enumBaseTable_DeadAssignmentsElimination_VariableLiveness = [...]j2gort.EnumBase{
	{Name_: "MAYBE_LIVE", Ordinal_: 0},
	{Name_: "READ", Ordinal_: 1},
	{Name_: "KILL", Ordinal_: 2},
}
var enumValues_DeadAssignmentsElimination_VariableLiveness = []DeadAssignmentsElimination_VariableLiveness{1, 2, 3}

func StaticGet_DeadAssignmentsElimination_VariableLiveness_MAYBE_LIVE() DeadAssignmentsElimination_VariableLiveness {
	return enum_DeadAssignmentsElimination_VariableLiveness_MAYBE_LIVE
}
func StaticGet_DeadAssignmentsElimination_VariableLiveness_READ() DeadAssignmentsElimination_VariableLiveness {
	return enum_DeadAssignmentsElimination_VariableLiveness_READ
}
func StaticGet_DeadAssignmentsElimination_VariableLiveness_KILL() DeadAssignmentsElimination_VariableLiveness {
	return enum_DeadAssignmentsElimination_VariableLiveness_KILL
}

func DeadAssignmentsElimination_VariableLiveness_Values() []DeadAssignmentsElimination_VariableLiveness {
	return j2gort.CloneSlice(enumValues_DeadAssignmentsElimination_VariableLiveness)
}

func DeadAssignmentsElimination_VariableLiveness_ValueOf(name string) DeadAssignmentsElimination_VariableLiveness {
	if name == j2gort.NullString {
		j2gort.ThrowNPE()
	}
	for i_ := range enumBaseTable_DeadAssignmentsElimination_VariableLiveness {
		if enumBaseTable_DeadAssignmentsElimination_VariableLiveness[i_].Name_ == name {
			return DeadAssignmentsElimination_VariableLiveness(i_ + 1)
		}
	}
	panic(javalang.NewIllegalArgumentException("No enum constant com.google.javascript.jscomp.DeadAssignmentsElimination.VariableLiveness." + name))
}

func (this_ DeadAssignmentsElimination_VariableLiveness) IsJavaNil_() bool                { return this_ == 0 }
func (DeadAssignmentsElimination_VariableLiveness) HasIdentityEquals_()                   {}
func (DeadAssignmentsElimination_VariableLiveness) Implements_j2gort_Ext_java_lang_Enum() {}
func (DeadAssignmentsElimination_VariableLiveness) Implements_j2gort_Ext_java_lang_constant_Constable() {
}
func (this_ DeadAssignmentsElimination_VariableLiveness) GetExt_java_lang_EnumBase() *j2gort.Ext_java_lang_EnumBase {
	return &j2gort.Ext_java_lang_EnumBase{Core_: this_.GetEnumBase()}
}
func (DeadAssignmentsElimination_VariableLiveness) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.DeadAssignmentsElimination$VariableLiveness", "VariableLiveness"
}
func (this_ DeadAssignmentsElimination_VariableLiveness) GetEnumBase() *j2gort.EnumBase {
	if this_ == 0 {
		return nil
	}
	return &enumBaseTable_DeadAssignmentsElimination_VariableLiveness[this_-1]
}
func (this_ DeadAssignmentsElimination_VariableLiveness) Ordinal() int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return int32(this_ - 1)
}
func (this_ DeadAssignmentsElimination_VariableLiveness) EnumOrdinal_() int32 {
	return int32(this_) - 1
}
func (this_ DeadAssignmentsElimination_VariableLiveness) Name() string {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_DeadAssignmentsElimination_VariableLiveness[this_-1].Name_
}
func (this_ DeadAssignmentsElimination_VariableLiveness) String() string {
	if this_ == 0 {
		return "null"
	}
	return enumBaseTable_DeadAssignmentsElimination_VariableLiveness[this_-1].Name_
}
func (this_ DeadAssignmentsElimination_VariableLiveness) Equals(other_ any) bool {
	o_, ok_ := other_.(DeadAssignmentsElimination_VariableLiveness)
	return this_ != 0 && ok_ && this_ == o_
}
func (this_ DeadAssignmentsElimination_VariableLiveness) HashCode() int32 {
	if this_ == 0 {
		return 0
	}
	return enumBaseTable_DeadAssignmentsElimination_VariableLiveness[this_-1].HashCode()
}
func (this_ DeadAssignmentsElimination_VariableLiveness) CompareTo(other_ any) int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_DeadAssignmentsElimination_VariableLiveness[this_-1].CompareTo(other_)
}
func (this_ DeadAssignmentsElimination_VariableLiveness) WriteGson(out_ *gson.JsonWriter) {
	if this_ == 0 {
		out_.NullValue()
		return
	}
	out_.ValueString(enumBaseTable_DeadAssignmentsElimination_VariableLiveness[this_-1].Name_)
}
func init() {
	vals_ := []any{DeadAssignmentsElimination_VariableLiveness(1), DeadAssignmentsElimination_VariableLiveness(2), DeadAssignmentsElimination_VariableLiveness(3)}
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.DeadAssignmentsElimination.VariableLiveness", vals_)
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.DeadAssignmentsElimination$VariableLiveness", vals_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.DeadAssignmentsElimination.VariableLiveness", "com.google.javascript.jscomp.DeadAssignmentsElimination$VariableLiveness", func(v_ any) bool {
		return j2gort.InstanceOf[DeadAssignmentsElimination_VariableLiveness](v_)
	})
}
func (*DeadAssignmentsElimination) AlwaysShouldTraverse_()      {}
func (*DeadAssignmentsElimination) CanReleaseSyntacticScopes_() {}
