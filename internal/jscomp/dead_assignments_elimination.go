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
	"github.com/paullepoulpe/jsgo/internal/ast"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// DeadAssignmentsElimination removes local variable assignments that are
// useless based on information from LiveVariablesAnalysis. If there is an
// assignment to variable x and x is dead after this assignment, we know that
// the current content of x will not be read and this assignment is useless.
type DeadAssignmentsElimination struct {
	cfgScopeStack

	compiler     AbstractCompiler
	scopeCreator *SyntacticScopeCreator
	liveness     *LiveVariablesAnalysis

	livenessStorage LiveVariablesAnalysis
	varsFinder      lvaAllVarsFinder

	deferChanges bool
	changeCount  int
}

// Java: DeadAssignmentsElimination(AbstractCompiler)
func NewDeadAssignmentsElimination(compiler AbstractCompiler) *DeadAssignmentsElimination {
	d := &DeadAssignmentsElimination{compiler: compiler}
	d.cfgs = d.cfgsBuf[:0]
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
	if allVarsDeclaredInFunction.table.isEmpty() {
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
	d.tryRemoveDeadAssignments(cfg)
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

	changeCounts := make([]int, len(candidates))
	runParallelWorkersIndexed(
		len(candidates),
		numThreads,
		func(_ int) *DeadAssignmentsElimination {
			wd := NewDeadAssignmentsElimination(d.compiler)
			wd.deferChanges = true
			wd.scopeCreator = NewSyntacticScopeCreator(d.compiler)
			return wd
		},
		func(wd **DeadAssignmentsElimination, idx int) {
			w := *wd
			w.changeCount = 0
			w.eliminateFunction(candidates[idx], dummyGlobal)
			changeCounts[idx] = w.changeCount
		},
		func(_ int, wd **DeadAssignmentsElimination) {
			(*wd).scopeCreator = nil
		},
	)
	for i, count := range changeCounts {
		for range count {
			d.compiler.ReportChangeToChangeScope(candidates[i])
		}
	}
}

func (d *DeadAssignmentsElimination) tryRemoveDeadAssignments(cfg *ControlFlowGraph[rhino.Node]) {
	for cfgNode := range cfg.Nodes.Values() {
		state := dataFlowAnalysis_linearFlowStateOf[*LiveVariablesAnalysis_LiveVariableLattice](cfgNode)
		n := cfgNode.Value
		if n.IsNil() {
			continue
		}
		switch n.GetToken() {
		case rhino.Token_IF, rhino.Token_WHILE, rhino.Token_DO, rhino.Token_FOR:
			cond := NodeUtil_GetConditionExpression(n)
			d.tryRemoveAssignment(cond, cond, state)
			continue
		case rhino.Token_FOR_IN, rhino.Token_FOR_OF, rhino.Token_FOR_AWAIT_OF:
			continue
		case rhino.Token_SWITCH:
			disc := ast.MustSwitch(n).Discriminant().Node()
			d.tryRemoveAssignment(disc, disc, state)
			continue
		case rhino.Token_CASE:
			val := ast.MustCase(n).Value().Node()
			d.tryRemoveAssignment(val, val, state)
			continue
		case rhino.Token_RETURN:
			if val, ok := ast.MustReturn(n).Value(); ok {
				vn := val.Node()
				d.tryRemoveAssignment(vn, vn, state)
			}
			continue
			// TODO(user): case VAR: Remove var a=1;a=2;.....
		default:
		}

		d.tryRemoveAssignment(n, n, state)
	}
}

// tryRemoveAssignment determines if any local variables are dead after the
// instruction n and are assigned within the subtree of n, removing those
// assignments if there are any.
func (d *DeadAssignmentsElimination) tryRemoveAssignment(
	n rhino.Node,
	exprRoot rhino.Node,
	state *DataFlowAnalysis_LinearFlowState[*LiveVariablesAnalysis_LiveVariableLattice],
) {
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
			d.tryRemoveAssignment(rhs, exprRoot, state)
			rhs = NodeUtil_GetRValueOfLValue(lhs)
		}

		// Multiple declarations should be processed from right-to-left to ensure side-effects
		// are run in the correct order.
		if isDeclarationNode && !lhs.GetNext().IsNil() {
			d.tryRemoveAssignment(lhs.GetNext(), exprRoot, state)
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
		lhsID := lhs.GetStringID()
		index, ok := d.liveness.lookupVarIndexByID(lhsID)
		if !ok || d.liveness.isEscapedIndex(index) {
			return
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
				d.tryRemoveAssignment(c, exprRoot, state)
			}
			c = next
		}
		return
	}
}

// isVariableStillLiveWithinExpressionByID returns true if there exists a read
// of varID before a write to varID that is on the right side of n within
// exprRoot.
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

// isVariableReadBeforeKillByID returns READ if the first reference of varID in
// n is a read, KILL if it is an assignment, or MAYBE_LIVE otherwise.
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
