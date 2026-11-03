// Port of com/google/javascript/jscomp/DeadAssignmentsElimination.java @ closure-compiler v20261028

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
		state := d.liveness.FlowStateAt(cfgNode.Index)
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
	} else if n.IsComputedFieldDef() {
		if key := n.GetFirstChild(); !key.IsNil() && !isCfgSubtreeBoundary(key, false, false) {
			d.tryRemoveAssignment(key, exprRoot, state)
		}
		if n.IsStaticMember() {
			if val := n.GetSecondChild(); !val.IsNil() && !isCfgSubtreeBoundary(val, false, false) {
				d.tryRemoveAssignment(val, exprRoot, state)
			}
		}
		return
	} else {
		for c := n.GetFirstChild(); !c.IsNil(); {
			next := c.GetNext()
			if !isCfgSubtreeBoundary(c, false, false) {
				d.tryRemoveAssignment(c, exprRoot, state)
			}
			c = next
		}
		return
	}
}

type daePeerEvalCtx struct {
	d     *DeadAssignmentsElimination
	varID uint32
}

func daeEvalPeerReadBeforeKill(ctx daePeerEvalCtx, expr rhino.Node) cfgCondEvalState {
	return cfgCondEvalState(ctx.d.isVariableReadBeforeKillByID(expr, ctx.varID))
}

// isVariableStillLiveWithinExpressionByID returns true if there exists a read
// of varID before a write to varID that evaluates after n within exprRoot.
func (d *DeadAssignmentsElimination) isVariableStillLiveWithinExpressionByID(
	n rhino.Node, exprRoot rhino.Node, varID uint32) bool {
	peerCtx := daePeerEvalCtx{d: d, varID: varID}
	for n != exprRoot {
		parent := n.GetParent()
		if (parent.IsComputedProp() || parent.IsComputedFieldDef()) && parent.GetParent().IsClassMembers() {
			if parent == exprRoot {
				if parent.IsComputedFieldDef() && parent.IsStaticMember() && n == parent.GetFirstChild() {
					if val := parent.GetSecondChild(); !val.IsNil() {
						return d.isVariableReadBeforeKillByID(val, varID) == DeadAssignmentsElimination_VariableLiveness_READ
					}
				}
				return false
			}
			classMembers := parent.GetParent()
			inPass2 := n != parent.GetFirstChild()
			if !inPass2 || parent.IsStaticMember() {
				switch DeadAssignmentsElimination_VariableLiveness(
					walkClassMembersAfterState(classMembers, parent, inPass2, peerCtx, daeEvalPeerReadBeforeKill),
				) {
				case DeadAssignmentsElimination_VariableLiveness_READ:
					return true
				case DeadAssignmentsElimination_VariableLiveness_KILL:
					return false
				}
			}
			n = classMembers
			continue
		}
		if parent.IsClassMembers() {
			if n.IsMemberFieldDef() && n.IsStaticMember() {
				switch DeadAssignmentsElimination_VariableLiveness(
					walkClassMembersAfterState(parent, n, true, peerCtx, daeEvalPeerReadBeforeKill),
				) {
				case DeadAssignmentsElimination_VariableLiveness_READ:
					return true
				case DeadAssignmentsElimination_VariableLiveness_KILL:
					return false
				}
			}
			n = parent
			continue
		}

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

const (
	_ = uint(DeadAssignmentsElimination_VariableLiveness_MAYBE_LIVE) - uint(cfgCondEvalContinue)
	_ = uint(cfgCondEvalContinue) - uint(DeadAssignmentsElimination_VariableLiveness_MAYBE_LIVE)
	_ = uint(DeadAssignmentsElimination_VariableLiveness_READ) - uint(cfgCondEvalRead)
	_ = uint(cfgCondEvalRead) - uint(DeadAssignmentsElimination_VariableLiveness_READ)
	_ = uint(DeadAssignmentsElimination_VariableLiveness_KILL) - uint(cfgCondEvalKill)
	_ = uint(cfgCondEvalKill) - uint(DeadAssignmentsElimination_VariableLiveness_KILL)
)

// isVariableReadBeforeKillByID returns READ if the first reference of varID in
// n is a read, KILL if it is an assignment, or MAYBE_LIVE otherwise.
func (d *DeadAssignmentsElimination) isVariableReadBeforeKillByID(
	n rhino.Node, varID uint32) DeadAssignmentsElimination_VariableLiveness {
	return DeadAssignmentsElimination_VariableLiveness(
		walkCfgCondEvalForward(n, false, varID, daeVisitReadBeforeKill),
	)
}

func daeVisitReadBeforeKill(varID uint32, n rhino.Node, _ bool) (cfgCondEvalState, bool) {
	if ControlFlowGraph_IsEnteringNewCfgNode(n) { // Not a FUNCTION
		return cfgCondEvalContinue, true
	}

	if n.IsName() && n.GetStringID() == varID {
		parent := n.GetParent()
		if NodeUtil_IsNameDeclOrSimpleAssignLhs(n, parent) {
			assign, isAssign := ast.AsAssign(parent)
			jutil.CheckStateMsg1(isAssign, "%s", parent)
			// The expression to which the assignment is made is evaluated before
			// the RHS is evaluated (normal left to right evaluation) but the KILL
			// occurs after the RHS is evaluated.
			if walkCfgCondEvalForward(assign.Value().Node(), false, varID, daeVisitReadBeforeKill) == cfgCondEvalRead {
				return cfgCondEvalRead, true
			}
			return cfgCondEvalKill, true
		}
		return cfgCondEvalRead, true
	}

	return cfgCondEvalContinue, false
}

func (d *DeadAssignmentsElimination) checkHookBranchReadBeforeKillByID(
	trueCase rhino.Node, falseCase rhino.Node, varID uint32) DeadAssignmentsElimination_VariableLiveness {
	// With a hook it is unknown which branch will run, so
	// we must be conservative. A read by either is a READ, and
	// a KILL is only considered if both KILL.
	return DeadAssignmentsElimination_VariableLiveness(
		joinCfgCondEvalBranches(
			walkCfgCondEvalForward(trueCase, true, varID, daeVisitReadBeforeKill),
			walkCfgCondEvalForward(falseCase, true, varID, daeVisitReadBeforeKill),
		),
	)
}
