// Port of com/google/javascript/jscomp/DeadAssignmentsElimination.java @ closure-compiler v20261015

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

package jscomp

import (
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// DeadAssignmentsElimination removes local variable assignments that are
// useless based on information from LiveVariablesAnalysis. If there is an
// assignment to variable x and x is dead after this assignment, we know that
// the current content of x will not be read and this assignment is useless.
type DeadAssignmentsElimination struct {
	NodeTraversal_AbstractCfgCallbackBase

	compiler     AbstractCompiler
	scopeCreator *SyntacticScopeCreator
	liveness     *LiveVariablesAnalysis
	// Go-only: language-adapter — the ArrayDeque used as a stack
	// (addFirst/peekFirst/removeFirst); the first element is the last one
	// of the slice.
	functionStack    []DeadAssignmentsElimination_BailoutInformation
	functionStackBuf [8]DeadAssignmentsElimination_BailoutInformation

	livenessStorage LiveVariablesAnalysis
	varsFinder      lvaAllVarsFinder
}

// DeadAssignmentsElimination_BailoutInformation is the nested class
// BailoutInformation.
type DeadAssignmentsElimination_BailoutInformation struct {
	containsFunction        bool
	containsRemovableAssign bool
}

// Java: DeadAssignmentsElimination(AbstractCompiler)
func NewDeadAssignmentsElimination(compiler AbstractCompiler) *DeadAssignmentsElimination {
	d := &DeadAssignmentsElimination{}
	d.InitNodeTraversal_AbstractCfgCallbackBase(d)
	d.compiler = compiler
	d.functionStack = d.functionStackBuf[:0]
	return d
}

// Java: DeadAssignmentsElimination.process(Node, Node)
func (d *DeadAssignmentsElimination) Process(externs *rhino.Node, root *rhino.Node) {
	jutil.CheckNotNull(externs)
	jutil.CheckNotNull(root)
	jutil.CheckState(d.compiler.GetLifeCycleStage().IsNormalized())
	sc := NewSyntacticScopeCreator(d.compiler)
	defer sc.ReleaseArena()
	b := NodeTraversal_Builder{
		compiler:     d.compiler,
		callback:     d,
		scopeCreator: sc,
	}
	b.Traverse(root)
}

// functionStackPeekFirst is ArrayDeque.peekFirst() (null when empty).
//
// Go-only: language-adapter — see functionStack.
func (d *DeadAssignmentsElimination) functionStackPeekFirst() *DeadAssignmentsElimination_BailoutInformation {
	if len(d.functionStack) == 0 {
		return nil
	}
	return &d.functionStack[len(d.functionStack)-1]
}

// Java: DeadAssignmentsElimination.visit(NodeTraversal, Node, Node)
func (d *DeadAssignmentsElimination) Visit(t *NodeTraversal, n *rhino.Node, parent *rhino.Node) {
	if len(d.functionStack) == 0 {
		return
	}
	top := &d.functionStack[len(d.functionStack)-1]
	if n.IsFunction() {
		top.containsFunction = true
	} else if !top.containsRemovableAssign && d.IsRemovableAssign(n) {
		top.containsRemovableAssign = true
	}
}

// Java: DeadAssignmentsElimination.enterScopeWithCfg(NodeTraversal)
func (d *DeadAssignmentsElimination) EnterScopeWithCfg(t *NodeTraversal) {
	if t.InFunctionBlockScope() {
		d.functionStack = append(d.functionStack, DeadAssignmentsElimination_BailoutInformation{})
	}
}

// Java: DeadAssignmentsElimination.exitScopeWithCfg(NodeTraversal)
func (d *DeadAssignmentsElimination) ExitScopeWithCfg(t *NodeTraversal) {
	if t.InFunctionBlockScope() {
		d.EliminateDeadAssignments(t)
		if len(d.functionStack) == 0 {
			panic(jutil.NewNoSuchElementException()) // ArrayDeque.removeFirst()
		}
		d.functionStack = d.functionStack[:len(d.functionStack)-1]
	}
}

// Java: DeadAssignmentsElimination.eliminateDeadAssignments(NodeTraversal)
func (d *DeadAssignmentsElimination) EliminateDeadAssignments(t *NodeTraversal) {
	jutil.CheckArgument(t.InFunctionBlockScope())
	jutil.CheckState(len(d.functionStack) != 0)

	// Skip unchanged functions (note that the scope root is the function block, not the function).
	if !d.compiler.HasScopeChanged(t.GetScopeRoot().GetParent()) {
		return
	}

	currentFunction := d.functionStackPeekFirst()
	// We are not going to do any dead assignment elimination in when there is
	// at least one inner function because in most browsers, when there is a
	// closure, ALL the variables are saved (escaped).
	if currentFunction.containsFunction {
		return
	}

	// We don't do any dead assignment elimination if there are no assigns
	// to eliminate. :)
	if !currentFunction.containsRemovableAssign {
		return
	}

	blockScope := t.GetScope()
	functionScope := blockScope.GetParent()
	if LiveVariablesAnalysis_MAX_VARIABLES_TO_ANALYZE <
		blockScope.GetVarCount()+functionScope.GetVarCount() {
		return
	}

	// Computes liveness information first.
	cfg := d.GetControlFlowGraph(d.compiler)
	if d.scopeCreator == nil {
		d.scopeCreator = NewSyntacticScopeCreator(d.compiler)
	}
	scopeCreator := d.scopeCreator
	allVarsDeclaredInFunction := d.varsFinder.collect(d.compiler, scopeCreator, functionScope)
	d.livenessStorage.Reset(
		cfg, functionScope, blockScope, d.compiler, scopeCreator, allVarsDeclaredInFunction)
	d.liveness = &d.livenessStorage
	d.liveness.Analyze()
	allVarsInFn := d.liveness.GetAllVariables()
	d.TryRemoveDeadAssignments(t, cfg, allVarsInFn)
}

// IsRemovableAssign matches all assignment operators and increment/decrement
// operators. Does *not* match VAR initialization, since
// RemoveUnusedVariables will already remove variables that are initialized
// but unused.
//
// Java: DeadAssignmentsElimination.isRemovableAssign(Node)
func (d *DeadAssignmentsElimination) IsRemovableAssign(n *rhino.Node) bool {
	return (NodeUtil_IsAssignmentOp(n) && n.GetFirstChild().IsName()) || n.IsInc() || n.IsDec()
}

// TryRemoveDeadAssignments tries to remove useless assignments from a
// control flow graph that has been annotated with liveness information.
//
// t is the node traversal; cfg the control flow graph of the program
// annotated with liveness information.
//
// Java: DeadAssignmentsElimination.tryRemoveDeadAssignments(NodeTraversal, ControlFlowGraph<Node>, Map<String, Var>)
func (d *DeadAssignmentsElimination) TryRemoveDeadAssignments(t *NodeTraversal,
	cfg ControlFlowGraph[*rhino.Node],
	allVarsInFn *jutil.LinkedHashMap[string, *Var]) {

	for cfgNode := range cfg.GetControlFlowGraphBase().nodes.Values() {
		state := dataFlowAnalysis_linearFlowStateOf[*LiveVariablesAnalysis_LiveVariableLattice](cfgNode)
		n := cfgNode.GetValue()
		if n == nil {
			continue
		}
		switch n.GetToken() {
		case rhino.Token_IF, rhino.Token_WHILE, rhino.Token_DO:
			d.TryRemoveAssignment(t, NodeUtil_GetConditionExpression(n), state, allVarsInFn)
			continue
		case rhino.Token_FOR, rhino.Token_FOR_IN, rhino.Token_FOR_OF, rhino.Token_FOR_AWAIT_OF:
			if n.IsVanillaFor() {
				d.TryRemoveAssignment(t, NodeUtil_GetConditionExpression(n), state, allVarsInFn)
			}
			continue
		case rhino.Token_SWITCH, rhino.Token_CASE, rhino.Token_RETURN:
			if n.HasChildren() {
				d.TryRemoveAssignment(t, n.GetFirstChild(), state, allVarsInFn)
			}
			continue
			// TODO(user): case VAR: Remove var a=1;a=2;.....
		default:
		}

		d.TryRemoveAssignment(t, n, state, allVarsInFn)
	}
}

// Java: DeadAssignmentsElimination.tryRemoveAssignment(NodeTraversal, Node, LinearFlowState<LiveVariableLattice>, Map<String, Var>)
func (d *DeadAssignmentsElimination) TryRemoveAssignment(
	t *NodeTraversal,
	n *rhino.Node,
	state *DataFlowAnalysis_LinearFlowState[*LiveVariablesAnalysis_LiveVariableLattice],
	allVarsInFn *jutil.LinkedHashMap[string, *Var]) {
	d.TryRemoveAssignmentWithExprRoot(t, n, n, state, allVarsInFn)
}

// TryRemoveAssignmentWithExprRoot determines if any local variables are dead
// after the instruction n and are assigned within the subtree of n. Removes
// those assignments if there are any.
//
// n is the target instruction; exprRoot the CFG node where the liveness
// information in state is still correct; state the liveness information at
// n.
//
// Java: DeadAssignmentsElimination.tryRemoveAssignment(NodeTraversal, Node, Node, LinearFlowState<LiveVariableLattice>, Map<String, Var>)
func (d *DeadAssignmentsElimination) TryRemoveAssignmentWithExprRoot(
	t *NodeTraversal,
	n *rhino.Node,
	exprRoot *rhino.Node,
	state *DataFlowAnalysis_LinearFlowState[*LiveVariablesAnalysis_LiveVariableLattice],
	allVarsInFn *jutil.LinkedHashMap[string, *Var]) {

	parent := n.GetParent()
	isDeclarationNode := NodeUtil_IsNameDeclaration(parent)

	if NodeUtil_IsAssignmentOp(n) || n.IsInc() || n.IsDec() || isDeclarationNode {

		if parent.IsConst() {
			// Removing the RHS of a const produces as invalid AST.
			return
		}

		var lhs *rhino.Node
		if isDeclarationNode {
			lhs = n
		} else {
			lhs = n.GetFirstChild()
		}
		rhs := NodeUtil_GetRValueOfLValue(lhs)

		// Recurse first. Example: dead_x = dead_y = 1; We try to clean up dead_y
		// first.
		if rhs != nil {
			d.TryRemoveAssignmentWithExprRoot(t, rhs, exprRoot, state, allVarsInFn)
			rhs = NodeUtil_GetRValueOfLValue(lhs)
		}

		// Multiple declarations should be processed from right-to-left to ensure side-effects
		// are run in the correct order.
		if isDeclarationNode && lhs.GetNext() != nil {
			d.TryRemoveAssignmentWithExprRoot(t, lhs.GetNext(), exprRoot, state, allVarsInFn)
		}

		// Ignore declarations that don't initialize a value. Dead code removal will kill those nodes.
		// Also ignore the var declaration if it's in a for-loop instantiation since there's not a
		// safe place to move the side-effects.
		if isDeclarationNode && (rhs == nil || NodeUtil_IsAnyFor(parent.GetParent())) {
			return
		}

		if !lhs.IsName() {
			return // Not a local variable assignment.
		}
		name := lhs.GetString()
		scope := t.GetScope()
		jutil.CheckState(scope.IsFunctionBlockScope() || scope.IsBlockScope())
		var_ := allVarsInFn.Get(name)
		if var_ == nil {
			return
		}

		if d.liveness.GetEscapedLocals().Contains(var_) {
			return // Local variable that might be escaped due to closures.
		}

		// If we have an identity assignment such as a=a, always remove it
		// regardless of what the liveness results because it
		// does not change the result afterward.
		if rhs != nil && rhs.IsName() && rhs.GetString() == var_.GetName() && n.IsAssign() {
			rhs.Detach()
			n.ReplaceWith(rhs)
			d.compiler.ReportChangeToEnclosingScope(rhs)
			return
		}

		index := d.liveness.GetVarIndex(var_.GetName())
		if state.GetOut().IsLive(index) {
			return // Variable not dead.
		}

		if state.GetIn().IsLive(index) &&
			d.IsVariableStillLiveWithinExpression(n, exprRoot, var_.GetName()) {
			// The variable is killed here but it is also live before it.
			// This is possible if we have say:
			//    if (X = a && a = C) {..} ; .......; a = S;
			// In this case we are safe to remove "a = C" because it is dead.
			// However if we have:
			//    if (a = C && X = a) {..} ; .......; a = S;
			// removing "a = C" is NOT correct, although the live set at the node
			// is exactly the same.
			// TODO(user): We need more fine grain CFA or we need to keep track
			// of GEN sets when we recurse here.
			return
		}

		if n.IsAssign() {
			rhs.Detach()
			n.ReplaceWith(rhs)
		} else if NodeUtil_IsAssignmentOp(n) {
			rhs.Detach()
			lhs.Detach()
			op := rhino.NewNodeWith2Children(NodeUtil_GetOpFromAssignmentOp(n), lhs, rhs)
			n.ReplaceWith(op)
		} else if n.IsInc() || n.IsDec() {
			if parent.IsExprResult() {
				n.ReplaceWith(rhino.IR_VoidNode(rhino.IR_Number(0).Srcref(n)))
			} else if n.IsComma() && n != parent.GetLastChild() {
				n.Detach()
			} else if parent.IsVanillaFor() && NodeUtil_GetConditionExpression(parent) != n {
				n.ReplaceWith(rhino.IR_Empty())
			} else {
				// Cannot replace x = a++ with x = a because that's not valid
				// when a is not a number.
				return
			}
		} else if isDeclarationNode {
			rhs.Detach()
			rhino.IR_ExprResult(rhs).InsertAfter(parent)
			rhs.GetParent().Srcref(rhs)
		} else {
			// Not reachable.
			panic(jutil.NewIllegalStateExceptionMsg("Unknown statement"))
		}

		d.compiler.ReportChangeToEnclosingScope(parent)
		return
	} else {
		for c := n.GetFirstChild(); c != nil; {
			next := c.GetNext()
			if !ControlFlowGraph_IsEnteringNewCfgNode(c) {
				d.TryRemoveAssignmentWithExprRoot(t, c, exprRoot, state, allVarsInFn)
			}
			c = next
		}
		return
	}
}

// IsVariableStillLiveWithinExpression: given a variable, node n in the tree
// and a sub-tree denoted by exprRoot as the root, this function returns true
// if there exists a read of that variable before a write to that variable
// that is on the right side of n.
//
// For example, suppose the node is x = 1:
//
//	y = 1, x = 1; // false, there is no reads at all.
//	y = 1, x = 1, print(x) // true, there is a read right of n.
//	y = 1, x = 1, x = 2, print(x) // false, there is a read right of n but
//	                              // it is after a write.
//
// n is the current node we should look at; exprRoot the node.
//
// Java: DeadAssignmentsElimination.isVariableStillLiveWithinExpression(Node, Node, String)
func (d *DeadAssignmentsElimination) IsVariableStillLiveWithinExpression(
	n *rhino.Node, exprRoot *rhino.Node, variable string) bool {
	for n != exprRoot {
		state := DeadAssignmentsElimination_VariableLiveness_MAYBE_LIVE
		switch n.GetParent().GetToken() {
		case rhino.Token_OR, rhino.Token_AND, rhino.Token_COALESCE, rhino.Token_OPTCHAIN_GETPROP, rhino.Token_OPTCHAIN_GETELEM, rhino.Token_OPTCHAIN_CALL:
			// For short-circuiting and optional chaining operations, if the current node is the
			// first child, subsequent siblings (the RHS or arguments) are only conditionally
			// evaluated. Therefore, be conservative: a READ in subsequent siblings means the variable
			// might be read, but a KILL cannot be guaranteed to occur.
			for sibling := n.GetNext(); sibling != nil; sibling = sibling.GetNext() {
				state = d.IsVariableReadBeforeKill(sibling, variable)
				if state != DeadAssignmentsElimination_VariableLiveness_MAYBE_LIVE {
					break
				}
			}
			if n == n.GetParent().GetFirstChild() && state == DeadAssignmentsElimination_VariableLiveness_KILL {
				state = DeadAssignmentsElimination_VariableLiveness_MAYBE_LIVE
			}
		case rhino.Token_HOOK:
			// If current node is the condition, check each following
			// branch, otherwise it is a conditional branch and the
			// other branch can be ignored.
			if n.GetNext() != nil && n.GetNext().GetNext() != nil {
				state = d.CheckHookBranchReadBeforeKill(n.GetNext(), n.GetNext().GetNext(), variable)
			}
		default:
			for sibling := n.GetNext(); sibling != nil; sibling = sibling.GetNext() {
				state = d.IsVariableReadBeforeKill(sibling, variable)
				if state != DeadAssignmentsElimination_VariableLiveness_MAYBE_LIVE {
					break
				}
			}
		}

		// If we see a READ or KILL there is no need to continue.
		if state == DeadAssignmentsElimination_VariableLiveness_READ {
			return true
		} else if state == DeadAssignmentsElimination_VariableLiveness_KILL {
			return false
		}
		n = n.GetParent()
	}
	return false
}

// DeadAssignmentsElimination_VariableLiveness is the current liveness of the
// variable.
type DeadAssignmentsElimination_VariableLiveness int32

const (
	DeadAssignmentsElimination_VariableLiveness_MAYBE_LIVE DeadAssignmentsElimination_VariableLiveness = iota // May be still live in the current expression tree.
	DeadAssignmentsElimination_VariableLiveness_READ                                                          // Known there is a read left of it.
	DeadAssignmentsElimination_VariableLiveness_KILL                                                          // Known there is a write before any read.
)

var deadAssignmentsElimination_VariableLiveness_names = [...]string{"MAYBE_LIVE", "READ", "KILL"}

// Java: DeadAssignmentsElimination.VariableLiveness.values()
func DeadAssignmentsElimination_VariableLiveness_values() []DeadAssignmentsElimination_VariableLiveness {
	return []DeadAssignmentsElimination_VariableLiveness{
		DeadAssignmentsElimination_VariableLiveness_MAYBE_LIVE,
		DeadAssignmentsElimination_VariableLiveness_READ,
		DeadAssignmentsElimination_VariableLiveness_KILL,
	}
}

// Java: DeadAssignmentsElimination.VariableLiveness.name()
func (v DeadAssignmentsElimination_VariableLiveness) String() string {
	return deadAssignmentsElimination_VariableLiveness_names[v]
}

// IsVariableReadBeforeKill: give an expression and a variable. It returns
// READ, if the first reference of that variable is a read. It returns KILL,
// if the first reference of that variable is an assignment. It returns
// MAY_LIVE otherwise.
//
// Java: DeadAssignmentsElimination.isVariableReadBeforeKill(Node, String)
func (d *DeadAssignmentsElimination) IsVariableReadBeforeKill(
	n *rhino.Node, variable string) DeadAssignmentsElimination_VariableLiveness {
	if ControlFlowGraph_IsEnteringNewCfgNode(n) { // Not a FUNCTION
		return DeadAssignmentsElimination_VariableLiveness_MAYBE_LIVE
	}

	if n.IsName() && variable == n.GetString() {
		if NodeUtil_IsNameDeclOrSimpleAssignLhs(n, n.GetParent()) {
			jutil.CheckStateMsg1(n.GetParent().IsAssign(), "%s", n.GetParent())
			// The expression to which the assignment is made is evaluated before
			// the RHS is evaluated (normal left to right evaluation) but the KILL
			// occurs after the RHS is evaluated.
			rhs := n.GetNext()
			state := d.IsVariableReadBeforeKill(rhs, variable)
			if state == DeadAssignmentsElimination_VariableLiveness_READ {
				return state
			}
			return DeadAssignmentsElimination_VariableLiveness_KILL
		} else {
			return DeadAssignmentsElimination_VariableLiveness_READ
		}
	}

	switch n.GetToken() {
	case rhino.Token_OR, rhino.Token_AND, rhino.Token_COALESCE, rhino.Token_OPTCHAIN_GETPROP, rhino.Token_OPTCHAIN_GETELEM:
		{
			// Conditionals
			v1 := d.IsVariableReadBeforeKill(n.GetFirstChild(), variable)
			v2 := d.IsVariableReadBeforeKill(n.GetLastChild(), variable)
			// With a AND/OR/COALESCE the first branch always runs, but the second is
			// may not.
			if v1 != DeadAssignmentsElimination_VariableLiveness_MAYBE_LIVE {
				return v1
			} else if v2 == DeadAssignmentsElimination_VariableLiveness_READ {
				return DeadAssignmentsElimination_VariableLiveness_READ
			} else {
				return DeadAssignmentsElimination_VariableLiveness_MAYBE_LIVE
			}
		}
	case rhino.Token_OPTCHAIN_CALL:
		{
			v1 := d.IsVariableReadBeforeKill(n.GetFirstChild(), variable)
			if v1 != DeadAssignmentsElimination_VariableLiveness_MAYBE_LIVE {
				return v1
			}
			for c := n.GetSecondChild(); c != nil; c = c.GetNext() {
				state := d.IsVariableReadBeforeKill(c, variable)
				if state != DeadAssignmentsElimination_VariableLiveness_MAYBE_LIVE {
					if state == DeadAssignmentsElimination_VariableLiveness_READ {
						return DeadAssignmentsElimination_VariableLiveness_READ
					}
					return DeadAssignmentsElimination_VariableLiveness_MAYBE_LIVE
				}
			}
			return DeadAssignmentsElimination_VariableLiveness_MAYBE_LIVE
		}
	case rhino.Token_HOOK:
		{
			first := d.IsVariableReadBeforeKill(n.GetFirstChild(), variable)
			if first != DeadAssignmentsElimination_VariableLiveness_MAYBE_LIVE {
				return first
			}
			return d.CheckHookBranchReadBeforeKill(n.GetSecondChild(), n.GetLastChild(), variable)
		}
	default:
		// Expressions are evaluated left-right, depth first.
		for child := n.GetFirstChild(); child != nil; child = child.GetNext() {
			state := d.IsVariableReadBeforeKill(child, variable)
			if state != DeadAssignmentsElimination_VariableLiveness_MAYBE_LIVE {
				return state
			}
		}
	}

	return DeadAssignmentsElimination_VariableLiveness_MAYBE_LIVE
}

// Java: DeadAssignmentsElimination.checkHookBranchReadBeforeKill(Node, Node, String)
func (d *DeadAssignmentsElimination) CheckHookBranchReadBeforeKill(
	trueCase *rhino.Node, falseCase *rhino.Node, variable string) DeadAssignmentsElimination_VariableLiveness {
	v1 := d.IsVariableReadBeforeKill(
		trueCase, variable)
	v2 := d.IsVariableReadBeforeKill(
		falseCase, variable)
	// With a hook it is unknown which branch will run, so
	// we must be conservative.  A read by either is a READ, and
	// a KILL is only considered if both KILL.
	if v1 == DeadAssignmentsElimination_VariableLiveness_READ || v2 == DeadAssignmentsElimination_VariableLiveness_READ {
		return DeadAssignmentsElimination_VariableLiveness_READ
	} else if v1 == DeadAssignmentsElimination_VariableLiveness_KILL && v2 == DeadAssignmentsElimination_VariableLiveness_KILL {
		return DeadAssignmentsElimination_VariableLiveness_KILL
	} else {
		return DeadAssignmentsElimination_VariableLiveness_MAYBE_LIVE
	}
}
