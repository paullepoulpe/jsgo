// Port of com/google/javascript/jscomp/DeadAssignmentsElimination.java @ closure-compiler v20261021

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
	"sync"
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/internal/ast"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// DeadAssignmentsElimination removes local variable assignments that are
// useless based on information from LiveVariablesAnalysis. If there is an
// assignment to variable x and x is dead after this assignment, we know that
// the current content of x will not be read and this assignment is useless.
type DeadAssignmentsElimination struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	cfgScopeStack

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

	deferChanges bool
	changeCount  int
}

// DeadAssignmentsElimination_BailoutInformation is the nested class
// BailoutInformation.
type DeadAssignmentsElimination_BailoutInformation struct {
	containsFunction        bool
	containsRemovableAssign bool
}

// Java: DeadAssignmentsElimination(AbstractCompiler)
func NewDeadAssignmentsElimination(compiler AbstractCompiler) *DeadAssignmentsElimination {
	d := &DeadAssignmentsElimination{compiler: compiler}
	d.cfgs = d.cfgsBuf[:0]
	d.functionStack = d.functionStackBuf[:0]
	return d
}

func isRemovableAssignNode(n rhino.Node) bool {
	if assign, ok := ast.AsAssignment(n); ok {
		_, isName := assign.Target().AsName()
		return isName
	}
	_, isIncDec := ast.AsIncDec(n)
	return isIncDec
}

func daeScanNode(n rhino.Node, compiler AbstractCompiler, candidates *[]rhino.Node) (hasFn bool, hasAssign bool) {
	if fn, ok := ast.AsFunction(n); ok {
		if fnBody, isBlock := fn.Body().AsBlock(); isBlock {
			daeScanNode(fn.NameNode().Node(), compiler, candidates)
			daeScanNode(fn.Params().Node(), compiler, candidates)
			bodyHasFn := false
			bodyHasAssign := false
			for c := fnBody.Node().GetFirstChild(); !c.IsNil(); c = c.GetNext() {
				cFn, cAssign := daeScanNode(c, compiler, candidates)
				if cFn {
					bodyHasFn = true
				}
				if cAssign {
					bodyHasAssign = true
				}
			}
			if !bodyHasFn && bodyHasAssign && compiler.HasScopeChanged(n) {
				*candidates = append(*candidates, n)
			}
			return true, false
		}
		hasFn = true
	}
	for c := n.GetFirstChild(); !c.IsNil(); c = c.GetNext() {
		cFn, cAssign := daeScanNode(c, compiler, candidates)
		if cFn {
			hasFn = true
		}
		if cAssign {
			hasAssign = true
		}
	}
	if !hasAssign && isRemovableAssignNode(n) {
		hasAssign = true
	}
	return hasFn, hasAssign
}

func (d *DeadAssignmentsElimination) eliminateFunction(fn rhino.Node, dummyGlobal *Scope) {
	fnBody := ast.MustFunction(fn).Body().Node()
	scopeCreator := d.scopeCreator
	functionScope := scopeCreator.CreateScopeOfScope(fn, dummyGlobal)
	blockScope := scopeCreator.CreateScopeOfScope(fnBody, functionScope)
	if LiveVariablesAnalysis_MAX_VARIABLES_TO_ANALYZE <
		blockScope.GetVarCount()+functionScope.GetVarCount() {
		return
	}

	allVarsDeclaredInFunction := d.varsFinder.collect(d.compiler, scopeCreator, functionScope)
	if allVarsDeclaredInFunction.GetAllVariables().IsEmpty() {
		return
	}
	d.pushCfgRoot(fn)
	defer d.popCfgRoot()
	cfg := d.GetControlFlowGraph(d.compiler)

	d.livenessStorage.Reset(
		cfg, functionScope, blockScope, d.compiler, scopeCreator, allVarsDeclaredInFunction)
	d.liveness = &d.livenessStorage
	if d.liveness.allInFnVarsEscaped() {
		return
	}
	d.liveness.Analyze()
	d.TryRemoveDeadAssignments(nil, cfg, d.liveness.GetAllVariables())
}

func (d *DeadAssignmentsElimination) reportChange(n rhino.Node) {
	if d.deferChanges {
		d.changeCount++
		return
	}
	d.compiler.ReportChangeToEnclosingScope(n)
}

// Java: DeadAssignmentsElimination.process(Node, Node)
func (d *DeadAssignmentsElimination) Process(externs rhino.Node, root rhino.Node) {
	rhino.CheckNotNull(externs)
	rhino.CheckNotNull(root)
	jutil.CheckState(d.compiler.GetLifeCycleStage().IsNormalized())

	var candidates []rhino.Node
	daeScanNode(root, d.compiler, &candidates)
	if len(candidates) == 0 {
		return
	}

	dummyGlobal := Scope_CreateGlobalScope(root)
	numThreads := 1
	if opts := d.compiler.GetOptions(); opts != nil && opts.GetNumParallelThreads() > 1 {
		numThreads = opts.GetNumParallelThreads()
	}

	if numThreads <= 1 || len(candidates) <= 1 {
		sc := NewSyntacticScopeCreator(d.compiler)
		d.scopeCreator = sc
		defer func() {
			d.scopeCreator = nil
		}()
		for _, fn := range candidates {
			d.eliminateFunction(fn, dummyGlobal)
		}
		return
	}

	workers := min(numThreads, len(candidates))
	changeCounts := make([]int, len(candidates))
	failures := make([]jutil.Throwable, len(candidates))
	var next atomic.Int64
	var wg sync.WaitGroup
	for range workers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			wd := NewDeadAssignmentsElimination(d.compiler)
			wd.deferChanges = true
			sc := NewSyntacticScopeCreator(d.compiler)
			wd.scopeCreator = sc
			defer func() {
				wd.scopeCreator = nil
			}()
			curIdx := -1
			defer func() {
				if r := recover(); r != nil && curIdx >= 0 {
					failures[curIdx] = jutil.AsThrowable(r)
				}
			}()
			for {
				idx := int(next.Add(1) - 1)
				if idx >= len(candidates) {
					return
				}
				curIdx = idx
				wd.changeCount = 0
				wd.eliminateFunction(candidates[idx], dummyGlobal)
				changeCounts[idx] = wd.changeCount
				curIdx = -1
			}
		}()
	}
	wg.Wait()
	for _, f := range failures {
		if f != nil {
			panic(f)
		}
	}
	for i, count := range changeCounts {
		for range count {
			d.compiler.ReportChangeToChangeScope(candidates[i])
		}
	}
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
func (d *DeadAssignmentsElimination) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
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

func (d *DeadAssignmentsElimination) EnterScope(t *NodeTraversal) {
	d.pushScopeRoot(t.GetScopeRoot())
	d.EnterScopeWithCfg(t)
}

func (d *DeadAssignmentsElimination) ExitScope(t *NodeTraversal) {
	d.ExitScopeWithCfg(t)
	d.popScopeRoot(t.GetScopeRoot())
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
	if d.scopeCreator == nil {
		d.scopeCreator = NewSyntacticScopeCreator(d.compiler)
	}
	scopeCreator := d.scopeCreator
	allVarsDeclaredInFunction := d.varsFinder.collect(d.compiler, scopeCreator, functionScope)
	if allVarsDeclaredInFunction.GetAllVariables().IsEmpty() {
		return
	}
	cfg := d.GetControlFlowGraph(d.compiler)
	d.livenessStorage.Reset(
		cfg, functionScope, blockScope, d.compiler, scopeCreator, allVarsDeclaredInFunction)
	d.liveness = &d.livenessStorage
	if d.liveness.allInFnVarsEscaped() {
		return
	}
	d.liveness.Analyze()
	d.TryRemoveDeadAssignments(t, cfg, d.liveness.GetAllVariables())
}

// IsRemovableAssign matches all assignment operators and increment/decrement
// operators. Does *not* match VAR initialization, since
// RemoveUnusedVariables will already remove variables that are initialized
// but unused.
//
// Java: DeadAssignmentsElimination.isRemovableAssign(Node)
func (d *DeadAssignmentsElimination) IsRemovableAssign(n rhino.Node) bool {
	return isRemovableAssignNode(n)
}

// TryRemoveDeadAssignments tries to remove useless assignments from a
// control flow graph that has been annotated with liveness information.
//
// t is the node traversal; cfg the control flow graph of the program
// annotated with liveness information.
//
// Java: DeadAssignmentsElimination.tryRemoveDeadAssignments(NodeTraversal, ControlFlowGraph<Node>, Map<String, Var>)
func (d *DeadAssignmentsElimination) TryRemoveDeadAssignments(t *NodeTraversal,
	cfg *ControlFlowGraph[rhino.Node],
	allVarsInFn *jutil.LinkedHashMap[string, *Var]) {

	for cfgNode := range cfg.Nodes.Values() {
		state := dataFlowAnalysis_linearFlowStateOf[*LiveVariablesAnalysis_LiveVariableLattice](cfgNode)
		n := cfgNode.Value
		if n.IsNil() {
			continue
		}
		switch n.GetToken() {
		case rhino.Token_IF, rhino.Token_WHILE, rhino.Token_DO:
			d.TryRemoveAssignment(t, NodeUtil_GetConditionExpression(n), state, allVarsInFn)
			continue
		case rhino.Token_FOR:
			d.TryRemoveAssignment(t, NodeUtil_GetConditionExpression(n), state, allVarsInFn)
			continue
		case rhino.Token_FOR_IN, rhino.Token_FOR_OF, rhino.Token_FOR_AWAIT_OF:
			continue
		case rhino.Token_SWITCH:
			d.TryRemoveAssignment(t, ast.MustSwitch(n).Discriminant().Node(), state, allVarsInFn)
			continue
		case rhino.Token_CASE:
			d.TryRemoveAssignment(t, ast.MustCase(n).Value().Node(), state, allVarsInFn)
			continue
		case rhino.Token_RETURN:
			if val, ok := ast.MustReturn(n).Value(); ok {
				d.TryRemoveAssignment(t, val.Node(), state, allVarsInFn)
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
	n rhino.Node,
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
	n rhino.Node,
	exprRoot rhino.Node,
	state *DataFlowAnalysis_LinearFlowState[*LiveVariablesAnalysis_LiveVariableLattice],
	allVarsInFn *jutil.LinkedHashMap[string, *Var]) {

	parent := n.GetParent()
	parentDecl, isDeclarationNode := ast.AsNameDeclaration(parent)
	assign, isAssignOp := ast.AsAssignment(n)
	incDec, isIncDec := ast.AsIncDec(n)

	if isAssignOp || isIncDec || isDeclarationNode {

		if isDeclarationNode && parentDecl.Token() == rhino.Token_CONST {
			// Removing the RHS of a const produces as invalid AST.
			return
		}

		var lhs rhino.Node
		if isDeclarationNode {
			lhs = n
		} else if isAssignOp {
			lhs = assign.Target().Node()
		} else {
			lhs = incDec.Target().Node()
		}
		rhs := NodeUtil_GetRValueOfLValue(lhs)

		// Recurse first. Example: dead_x = dead_y = 1; We try to clean up dead_y
		// first.
		if !rhs.IsNil() {
			d.TryRemoveAssignmentWithExprRoot(t, rhs, exprRoot, state, allVarsInFn)
			rhs = NodeUtil_GetRValueOfLValue(lhs)
		}

		// Multiple declarations should be processed from right-to-left to ensure side-effects
		// are run in the correct order.
		if isDeclarationNode && !lhs.GetNext().IsNil() {
			d.TryRemoveAssignmentWithExprRoot(t, lhs.GetNext(), exprRoot, state, allVarsInFn)
		}

		// Ignore declarations that don't initialize a value. Dead code removal will kill those nodes.
		// Also ignore the var declaration if it's in a for-loop instantiation since there's not a
		// safe place to move the side-effects.
		if isDeclarationNode && (rhs.IsNil() || NodeUtil_IsAnyFor(parent.GetParent())) {
			return
		}

		if !lhs.IsName() {
			return // Not a local variable assignment.
		}
		if t != nil {
			scope := t.GetScope()
			jutil.CheckState(scope.IsFunctionBlockScope() || scope.IsBlockScope())
		}
		lhsID := lhs.GetStringID()
		var index int
		if allVarsInFn == d.liveness.allVarsInFn {
			var ok bool
			index, ok = d.liveness.lookupVarIndexByID(lhsID)
			if !ok || d.liveness.isEscapedIndex(index) {
				return
			}
		} else {
			var_ := allVarsInFn.Get(lhs.GetString())
			if var_ == nil {
				return
			}
			if d.liveness.GetEscapedLocals().Contains(var_) {
				return // Local variable that might be escaped due to closures.
			}
			lhsID = var_.nameID
			index = d.liveness.getVarIndexByID(lhsID)
		}

		// If we have an identity assignment such as a=a, always remove it
		// regardless of what the liveness results because it
		// does not change the result afterward.
		if isAssignOp && assign.Token() == rhino.Token_ASSIGN && !rhs.IsNil() && rhs.IsName() && rhs.GetStringID() == lhsID {
			rhs.Detach()
			n.ReplaceWith(rhs)
			d.reportChange(rhs)
			return
		}

		if state.GetOut().IsLive(index) {
			return // Variable not dead.
		}

		if state.GetIn().IsLive(index) &&
			d.isVariableStillLiveWithinExpressionByID(n, exprRoot, lhsID) {
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

		if isAssignOp && assign.Token() == rhino.Token_ASSIGN {
			rhs.Detach()
			n.ReplaceWith(rhs)
		} else if isAssignOp {
			rhs.Detach()
			lhs.Detach()
			op := ast.NewBinary(NodeUtil_GetOpFromAssignmentOp(n), ast.MustExpr(lhs), ast.MustExpr(rhs)).Node()
			n.ReplaceWith(op)
		} else if isIncDec {
			if _, isExprStmt := ast.AsExprStmt(parent); isExprStmt {
				n.ReplaceWith(rhino.IR_VoidNode(rhino.IR_Number(0).Srcref(n)))
			} else if n.IsComma() && n != parent.GetLastChild() {
				// Note: upstream Java checks n.isComma() (always false inside isInc/isDec)
				// rather than parent.isComma(); preserved for parity.
				n.Detach()
			} else if forStmt, isFor := ast.AsFor(parent); isFor && NodeUtil_GetConditionExpression(forStmt.Node()) != n {
				n.ReplaceWith(ast.NewEmpty().Node())
			} else {
				// Cannot replace x = a++ with x = a because that's not valid
				// when a is not a number.
				return
			}
		} else if isDeclarationNode {
			rhs.Detach()
			ast.NewExprStmt(ast.MustExpr(rhs)).Node().Srcref(rhs).InsertAfter(parent)
		} else {
			// Not reachable.
			panic(jutil.NewIllegalStateExceptionMsg("Unknown statement"))
		}

		d.reportChange(parent)
		return
	} else {
		for c := n.GetFirstChild(); !c.IsNil(); {
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
	n rhino.Node, exprRoot rhino.Node, variable string) bool {
	id := rhino.RhinoStringPool_LookupID(variable)
	if id == 0 {
		return false
	}
	return d.isVariableStillLiveWithinExpressionByID(n, exprRoot, id)
}

func (d *DeadAssignmentsElimination) isVariableStillLiveWithinExpressionByID(
	n rhino.Node, exprRoot rhino.Node, varID uint32) bool {
	for n != exprRoot {
		parent := n.GetParent()
		state := DeadAssignmentsElimination_VariableLiveness_MAYBE_LIVE
		switch parent.GetToken() {
		case rhino.Token_OR, rhino.Token_AND, rhino.Token_COALESCE, rhino.Token_OPTCHAIN_GETPROP, rhino.Token_OPTCHAIN_GETELEM, rhino.Token_OPTCHAIN_CALL:
			// For short-circuiting and optional chaining operations, if the current node is the
			// first child, subsequent siblings (the RHS or arguments) are only conditionally
			// evaluated. Therefore, be conservative: a READ in subsequent siblings means the variable
			// might be read, but a KILL cannot be guaranteed to occur.
			for sibling := n.GetNext(); !sibling.IsNil(); sibling = sibling.GetNext() {
				state = d.isVariableReadBeforeKillByID(sibling, varID)
				if state != DeadAssignmentsElimination_VariableLiveness_MAYBE_LIVE {
					break
				}
			}
			if n == parent.GetFirstChild() && state == DeadAssignmentsElimination_VariableLiveness_KILL {
				state = DeadAssignmentsElimination_VariableLiveness_MAYBE_LIVE
			}
		case rhino.Token_HOOK:
			// If current node is the condition, check each following
			// branch, otherwise it is a conditional branch and the
			// other branch can be ignored.
			hook := ast.MustHook(parent)
			if hook.Cond().Is(n) {
				state = d.checkHookBranchReadBeforeKillByID(hook.Then().Node(), hook.Else().Node(), varID)
			}
		default:
			for sibling := n.GetNext(); !sibling.IsNil(); sibling = sibling.GetNext() {
				state = d.isVariableReadBeforeKillByID(sibling, varID)
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
		n = parent
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
	n rhino.Node, variable string) DeadAssignmentsElimination_VariableLiveness {
	id := rhino.RhinoStringPool_LookupID(variable)
	if id == 0 {
		return DeadAssignmentsElimination_VariableLiveness_MAYBE_LIVE
	}
	return d.isVariableReadBeforeKillByID(n, id)
}

func (d *DeadAssignmentsElimination) isVariableReadBeforeKillByID(
	n rhino.Node, varID uint32) DeadAssignmentsElimination_VariableLiveness {
	if ControlFlowGraph_IsEnteringNewCfgNode(n) { // Not a FUNCTION
		return DeadAssignmentsElimination_VariableLiveness_MAYBE_LIVE
	}

	if n.IsName() && n.GetStringID() == varID {
		parent := n.GetParent()
		if NodeUtil_IsNameDeclOrSimpleAssignLhs(n, parent) {
			assign, isAssign := ast.AsAssign(parent)
			jutil.CheckStateMsg1(isAssign, "%s", parent)
			// The expression to which the assignment is made is evaluated before
			// the RHS is evaluated (normal left to right evaluation) but the KILL
			// occurs after the RHS is evaluated.
			state := d.isVariableReadBeforeKillByID(assign.Value().Node(), varID)
			if state == DeadAssignmentsElimination_VariableLiveness_READ {
				return state
			}
			return DeadAssignmentsElimination_VariableLiveness_KILL
		}
		return DeadAssignmentsElimination_VariableLiveness_READ
	}

	switch n.GetToken() {
	case rhino.Token_OPTCHAIN_GETPROP:
		// Go-only: OPTCHAIN_GETPROP has only one child (Object; property is its string),
		// so visiting Object once matches Java without O(2^k) re-traversal on chains.
		return d.isVariableReadBeforeKillByID(ast.MustOptChainGetProp(n).Object().Node(), varID)
	case rhino.Token_OR, rhino.Token_AND, rhino.Token_COALESCE:
		bin := ast.MustBinary(n)
		v1 := d.isVariableReadBeforeKillByID(bin.Left().Node(), varID)
		v2 := d.isVariableReadBeforeKillByID(bin.Right().Node(), varID)
		// With a AND/OR/COALESCE the first branch always runs, but the second
		// may not.
		if v1 != DeadAssignmentsElimination_VariableLiveness_MAYBE_LIVE {
			return v1
		} else if v2 == DeadAssignmentsElimination_VariableLiveness_READ {
			return DeadAssignmentsElimination_VariableLiveness_READ
		}
		return DeadAssignmentsElimination_VariableLiveness_MAYBE_LIVE
	case rhino.Token_OPTCHAIN_GETELEM:
		ge := ast.MustOptChainGetElem(n)
		v1 := d.isVariableReadBeforeKillByID(ge.Object().Node(), varID)
		v2 := d.isVariableReadBeforeKillByID(ge.Index().Node(), varID)
		if v1 != DeadAssignmentsElimination_VariableLiveness_MAYBE_LIVE {
			return v1
		} else if v2 == DeadAssignmentsElimination_VariableLiveness_READ {
			return DeadAssignmentsElimination_VariableLiveness_READ
		}
		return DeadAssignmentsElimination_VariableLiveness_MAYBE_LIVE
	case rhino.Token_OPTCHAIN_CALL:
		call := ast.MustOptChainCall(n)
		v1 := d.isVariableReadBeforeKillByID(call.Callee().Node(), varID)
		if v1 != DeadAssignmentsElimination_VariableLiveness_MAYBE_LIVE {
			return v1
		}
		for arg, ok := call.FirstArg(); ok; arg, ok = arg.Next() {
			state := d.isVariableReadBeforeKillByID(arg.Node(), varID)
			if state != DeadAssignmentsElimination_VariableLiveness_MAYBE_LIVE {
				if state == DeadAssignmentsElimination_VariableLiveness_READ {
					return DeadAssignmentsElimination_VariableLiveness_READ
				}
				return DeadAssignmentsElimination_VariableLiveness_MAYBE_LIVE
			}
		}
		return DeadAssignmentsElimination_VariableLiveness_MAYBE_LIVE
	case rhino.Token_HOOK:
		hook := ast.MustHook(n)
		first := d.isVariableReadBeforeKillByID(hook.Cond().Node(), varID)
		if first != DeadAssignmentsElimination_VariableLiveness_MAYBE_LIVE {
			return first
		}
		return d.checkHookBranchReadBeforeKillByID(hook.Then().Node(), hook.Else().Node(), varID)
	default:
		// Expressions are evaluated left-right, depth first.
		for child := n.GetFirstChild(); !child.IsNil(); child = child.GetNext() {
			state := d.isVariableReadBeforeKillByID(child, varID)
			if state != DeadAssignmentsElimination_VariableLiveness_MAYBE_LIVE {
				return state
			}
		}
	}

	return DeadAssignmentsElimination_VariableLiveness_MAYBE_LIVE
}

// Java: DeadAssignmentsElimination.checkHookBranchReadBeforeKill(Node, Node, String)
func (d *DeadAssignmentsElimination) CheckHookBranchReadBeforeKill(
	trueCase rhino.Node, falseCase rhino.Node, variable string) DeadAssignmentsElimination_VariableLiveness {
	id := rhino.RhinoStringPool_LookupID(variable)
	if id == 0 {
		return DeadAssignmentsElimination_VariableLiveness_MAYBE_LIVE
	}
	return d.checkHookBranchReadBeforeKillByID(trueCase, falseCase, id)
}

func (d *DeadAssignmentsElimination) checkHookBranchReadBeforeKillByID(
	trueCase rhino.Node, falseCase rhino.Node, varID uint32) DeadAssignmentsElimination_VariableLiveness {
	v1 := d.isVariableReadBeforeKillByID(trueCase, varID)
	v2 := d.isVariableReadBeforeKillByID(falseCase, varID)
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
