// Port of com/google/javascript/jscomp/MaybeReachingVariableUse.java @ closure-compiler v20261015

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

package jscomp

import (
	"iter"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// MaybeReachingVariableUse computes "may be" reaching use for all
// definitions of each variable.
//
// A use of A in `alert(A)` is a "may be" reaching use of the definition of A
// at `A = foo()` if at least one path from the definition node to the end
// node reaches that use and it is the last definition before the use on
// that path.
//
// Example:
//
//	D1: var A = foo();
//	U1: alert(A);
//	    if(....) {
//	       D2: A = bar();
//	       U2: alert(A);
//	    }
//	U3: alert(A);
//
// Here, MaybeReachingUses[D1] = {U1, U3} and MaybeReachingUses[D2]={U2, U3}.
// The use U3 is not guaranteed to use def D1: this is a "may-be" analysis.
//
// This pass is a backwards-analysis pass, i.e. it traverses the CFG nodes
// bottom-up, `MAX_STEPS` times or till a fixed point solution is reached. At
// each `cfgNode`, it:
//
//  1. Creates a new output lattice element to store the set of upward
//     exposed variable uses at `cfgNode`.
//  2. Propagates the set of existing upwards exposed variable uses at
//     `cfgNode` from the input lattice to output lattice
//  3. Adds new exposed uses of a variable to the upward exposed set in the
//     output lattice
//  4. Removes killed(unconditionally redefined) variables from upward
//     exposed set in the output lattice
type MaybeReachingVariableUse struct {
	DataFlowAnalysisBase[*rhino.Node, *MaybeReachingVariableUse_ReachingUses]

	// The scope of the function that we are analyzing.
	//
	// Go-only: language-adapter — Var overrides equals/hashCode (ScopedName),
	// hence the Object set.
	escaped     *jutil.ObjectLinkedHashSet[*Var]
	allVarsInFn *jutil.LinkedHashMap[string, *Var]

	emptyLattice MaybeReachingVariableUse_ReachingUses
	joinOp       MaybeReachingVariableUse_ReachingUsesJoinOp
	latticeSlab  []MaybeReachingVariableUse_ReachingUses
	latticeBuf   [8]MaybeReachingVariableUse_ReachingUses
}

// Java: MaybeReachingVariableUse(ControlFlowGraph<Node>, Set<Var>, Map<String, Var>)
func NewMaybeReachingVariableUse(
	cfg ControlFlowGraph[*rhino.Node],
	escaped *jutil.ObjectLinkedHashSet[*Var],
	allVarsInFn *jutil.LinkedHashMap[string, *Var]) *MaybeReachingVariableUse {
	u := &MaybeReachingVariableUse{}
	u.InitDataFlowAnalysisBase(u, cfg)
	u.escaped = escaped
	u.allVarsInFn = allVarsInFn
	u.emptyLattice.mayUsePMap = rhino.HamtPMap_Empty[*Var, rhino.PMap[*rhino.Node, *rhino.Node]]()
	u.latticeSlab = u.latticeBuf[:0:len(u.latticeBuf)]
	return u
}

func (u *MaybeReachingVariableUse) allocLattice(m rhino.PMap[*Var, rhino.PMap[*rhino.Node, *rhino.Node]]) *MaybeReachingVariableUse_ReachingUses {
	if len(u.latticeSlab) == cap(u.latticeSlab) {
		nextCap := min(max(cap(u.latticeSlab)*2, 16), 256)
		u.latticeSlab = make([]MaybeReachingVariableUse_ReachingUses, 0, nextCap)
	}
	idx := len(u.latticeSlab)
	u.latticeSlab = u.latticeSlab[:idx+1]
	lat := &u.latticeSlab[idx]
	lat.mayUsePMap = m
	return lat
}

// MaybeReachingVariableUse_ReachingUses is the may use definition lattice
// representation. It captures a product lattice for each local
// (non-escaped) variable. The sub-lattice is a n + 2 power set element
// lattice with all the Nodes in the program, TOP and BOTTOM. This is better
// explained with an example:
//
// Consider: A sub-lattice element representing the variable A represented
// by { N_4, N_5} where N_x is a Node in the program. This implies at that
// particular point in the program the content of A is "upward exposed" at
// point N_4 and N_5.
//
// Example:
//
//	A = 1;
//	...
//	N_3:
//	N_4: print(A);
//	N_5: y = A;
//	N_6: A = 1;
//	N_7: print(A);
//
// At N_3, reads of A in {N_4, N_5} are said to be upward exposed.
type MaybeReachingVariableUse_ReachingUses struct {
	// Maps variables to all their uses that are upward exposed at the
	// current cfgNode.
	mayUsePMap rhino.PMap[*Var, rhino.PMap[*rhino.Node, *rhino.Node]]
}

// Java: MaybeReachingVariableUse.ReachingUses()
func NewMaybeReachingVariableUse_ReachingUses() *MaybeReachingVariableUse_ReachingUses {
	return &MaybeReachingVariableUse_ReachingUses{
		mayUsePMap: rhino.HamtPMap_Empty[*Var, rhino.PMap[*rhino.Node, *rhino.Node]](),
	}
}

// NewMaybeReachingVariableUse_ReachingUsesCopy is the copy constructor.
//
// other: the constructed object is a replicated copy of this element.
//
// Java: MaybeReachingVariableUse.ReachingUses(ReachingUses)
func NewMaybeReachingVariableUse_ReachingUsesCopy(other *MaybeReachingVariableUse_ReachingUses) *MaybeReachingVariableUse_ReachingUses {
	return &MaybeReachingVariableUse_ReachingUses{
		mayUsePMap: other.mayUsePMap,
	}
}

// Java: MaybeReachingVariableUse.ReachingUses.get(Var)
func (r *MaybeReachingVariableUse_ReachingUses) Get(v *Var) iter.Seq[*rhino.Node] {
	uses := r.mayUsePMap.Get(v)
	if uses != nil {
		return uses.Keys()
	}
	return func(func(*rhino.Node) bool) {}
}

// Java: MaybeReachingVariableUse.ReachingUses.removeAll(Var)
func (r *MaybeReachingVariableUse_ReachingUses) RemoveAll(v *Var) {
	r.mayUsePMap = r.mayUsePMap.Minus(v)
}

// Java: MaybeReachingVariableUse.ReachingUses.put(Var, Node)
func (r *MaybeReachingVariableUse_ReachingUses) Put(v *Var, n *rhino.Node) {
	values := r.mayUsePMap.Get(v)
	if values == nil {
		values = rhino.HamtPMap_Empty[*rhino.Node, *rhino.Node]()
	}
	newValues := values.Plus(n, n)
	if newValues != values {
		r.mayUsePMap = r.mayUsePMap.Plus(v, newValues)
	}
}

func maybeReachingUse_reconcileNodes(node *rhino.Node, unused1 *rhino.Node, unused2 *rhino.Node) *rhino.Node {
	return node
}

func maybeReachingUse_reconcileVars(var_ *Var, thisVal rhino.PMap[*rhino.Node, *rhino.Node], thatVal rhino.PMap[*rhino.Node, *rhino.Node]) rhino.PMap[*rhino.Node, *rhino.Node] {
	if thisVal == nil {
		return thatVal
	}
	if thatVal == nil || thisVal == thatVal {
		return thisVal
	}
	// The PMap as a set with the key and value pair have the same value.
	return thisVal.Reconcile(thatVal, maybeReachingUse_reconcileNodes)
}

// Java: MaybeReachingVariableUse.ReachingUses.join(ReachingUses)
func (r *MaybeReachingVariableUse_ReachingUses) Join(other *MaybeReachingVariableUse_ReachingUses) {
	empty := rhino.HamtPMap_Empty[*Var, rhino.PMap[*rhino.Node, *rhino.Node]]()
	if r.mayUsePMap == other.mayUsePMap || other.mayUsePMap == empty {
		return
	}
	if r.mayUsePMap == empty {
		r.mayUsePMap = other.mayUsePMap
		return
	}
	r.mayUsePMap = r.mayUsePMap.Reconcile(other.mayUsePMap, maybeReachingUse_reconcileVars)
}

// Java: MaybeReachingVariableUse.ReachingUses.equals(Object)
func (r *MaybeReachingVariableUse_ReachingUses) Equals(other any) bool {
	reachingUses, ok := other.(*MaybeReachingVariableUse_ReachingUses)
	return ok && reachingUses != nil &&
		(r == reachingUses || r.mayUsePMap == reachingUses.mayUsePMap ||
			reachingUses.mayUsePMap.Equivalent(r.mayUsePMap, MaybeReachingVariableUse_ReachingUses_EqualMaps))
}

// Java: MaybeReachingVariableUse.ReachingUses.equalMaps(PMap<Node, Node>, PMap<Node, Node>)
func MaybeReachingVariableUse_ReachingUses_EqualMaps(map1 rhino.PMap[*rhino.Node, *rhino.Node], map2 rhino.PMap[*rhino.Node, *rhino.Node]) bool {
	return map1 == map2 || map1.Equivalent(map2, MaybeReachingVariableUse_ReachingUses_EqualNodes)
}

// Java: MaybeReachingVariableUse.ReachingUses.equalNodes(Node, Node)
func MaybeReachingVariableUse_ReachingUses_EqualNodes(n1 *rhino.Node, n2 *rhino.Node) bool {
	return n1 == n2
}

// Java: MaybeReachingVariableUse.ReachingUses.hashCode()
func (r *MaybeReachingVariableUse_ReachingUses) HashCode() int32 {
	panic(jutil.NewUnsupportedOperationExceptionMsg("the hashcode of this object is not stable"))
}

// MaybeReachingVariableUse_ReachingUsesJoinOp: the join is a simple union
// because of the "may be" nature of the analysis.
//
// Consider: A = 1; if (x) { A = 2 }; alert(A);
//
// The read of A "may be" exposed to A = 1 in the beginning.
type MaybeReachingVariableUse_ReachingUsesJoinOp struct {
	result *MaybeReachingVariableUse_ReachingUses
}

// Java: MaybeReachingVariableUse.ReachingUsesJoinOp()
func NewMaybeReachingVariableUse_ReachingUsesJoinOp() *MaybeReachingVariableUse_ReachingUsesJoinOp {
	return &MaybeReachingVariableUse_ReachingUsesJoinOp{result: NewMaybeReachingVariableUse_ReachingUses()}
}

// Java: MaybeReachingVariableUse.ReachingUsesJoinOp.joinFlow(ReachingUses)
func (j *MaybeReachingVariableUse_ReachingUsesJoinOp) JoinFlow(uses *MaybeReachingVariableUse_ReachingUses) {
	j.result.Join(uses)
}

// Java: MaybeReachingVariableUse.ReachingUsesJoinOp.finish()
func (j *MaybeReachingVariableUse_ReachingUsesJoinOp) Finish() *MaybeReachingVariableUse_ReachingUses {
	return j.result
}

// Java: MaybeReachingVariableUse.isForward()
func (u *MaybeReachingVariableUse) IsForward() bool {
	return false
}

// Java: MaybeReachingVariableUse.createEntryLattice()
func (u *MaybeReachingVariableUse) CreateEntryLattice() *MaybeReachingVariableUse_ReachingUses {
	return &u.emptyLattice
}

// Java: MaybeReachingVariableUse.createInitialEstimateLattice()
func (u *MaybeReachingVariableUse) CreateInitialEstimateLattice() *MaybeReachingVariableUse_ReachingUses {
	return &u.emptyLattice
}

// Java: MaybeReachingVariableUse.createFlowJoiner()
func (u *MaybeReachingVariableUse) CreateFlowJoiner() DataFlowAnalysis_FlowJoiner[*MaybeReachingVariableUse_ReachingUses] {
	u.joinOp.result = u.allocLattice(u.emptyLattice.mayUsePMap)
	return &u.joinOp
}

// FlowThrough computes the new LatticeElement for a given node given its
// LatticeElement from previous iteration.
//
// input: backward dataflow analyses compute their LatticeElement bottom-up
// (i.e. LinearFlowState.out to LinearFlowState.in). See
// DataFlowAnalysis#flow(DiGraphNode). Here param `input` is the readonly
// input LinearFlowState.out that was constructed as `LinearFlowState.in` in
// the previous iteration, or the initial lattice element if this is the
// first iteration.
//
// Java: MaybeReachingVariableUse.flowThrough(Node, ReachingUses)
func (u *MaybeReachingVariableUse) FlowThrough(n *rhino.Node, input *MaybeReachingVariableUse_ReachingUses) *MaybeReachingVariableUse_ReachingUses {
	temp := MaybeReachingVariableUse_ReachingUses{mayUsePMap: input.mayUsePMap}

	// If there's an ON_EX edge, this cfgNode may or may not get executed.
	// We can express this concisely by just pretending this happens in
	// a conditional.
	conditional := u.HasExceptionHandler(n)
	u.ComputeMayUse(n, n, &temp, conditional)
	if temp.mayUsePMap == input.mayUsePMap {
		return input
	}
	return u.allocLattice(temp.mayUsePMap)
}

// Java: MaybeReachingVariableUse.hasExceptionHandler(Node)
func (u *MaybeReachingVariableUse) HasExceptionHandler(cfgNode *rhino.Node) bool {
	branchEdges := u.GetCfg().GetOutEdges(cfgNode)
	for _, edge := range branchEdges {
		if edge.GetValue() == ControlFlowGraph_Branch_ON_EX {
			return true
		}
	}
	return false
}

// ComputeMayUse updates, given a cfgNode, the output LatticeElement at that
// node by finding and storing all variables and their uses that are upward
// exposed at the cfgNode.
//
// n is the explorer node which searches for variables; cfgNode the CFG node
// for which the upward exposed variables are being searched; conditional
// whether n is only conditionally evaluated given that cfgNode is
// evaluated. Do not remove conditionally redefined variables from the
// reaching uses set.
//
// Java: MaybeReachingVariableUse.computeMayUse(Node, Node, ReachingUses, boolean)
func (u *MaybeReachingVariableUse) ComputeMayUse(n *rhino.Node, cfgNode *rhino.Node, output *MaybeReachingVariableUse_ReachingUses, conditional bool) {
	switch n.GetToken() {
	case rhino.Token_BLOCK, rhino.Token_ROOT, rhino.Token_FUNCTION:
		return
	case rhino.Token_NAME:
		if NodeUtil_IsLhsByDestructuring(n) {
			if !conditional {
				u.RemoveFromUseIfLocal(n.GetString(), output)
			}
		} else {
			u.AddToUseIfLocal(n.GetString(), cfgNode, output)
		}
		return
	case rhino.Token_WHILE, rhino.Token_DO, rhino.Token_IF, rhino.Token_FOR:
		condExpr := NodeUtil_GetConditionExpression(n)
		u.ComputeMayUse(condExpr, cfgNode, output, conditional)
		return
	case rhino.Token_FOR_IN, rhino.Token_FOR_OF, rhino.Token_FOR_AWAIT_OF:
		// for(x in y) {...}
		lhs := n.GetFirstChild()
		rhs := lhs.GetNext()
		if NodeUtil_IsNameDeclaration(lhs) {
			lhs = lhs.GetLastChild() // for(var x in y) {...}
			if lhs.IsDestructuringLhs() {
				lhs = lhs.GetFirstChild() // for (let [x] of obj) {...}
			}
		}
		if lhs.IsDestructuringPattern() {
			u.ComputeMayUse(lhs, cfgNode, output, true)
		}
		u.ComputeMayUse(rhs, cfgNode, output, conditional)
		return
	case rhino.Token_AND, rhino.Token_OR, rhino.Token_COALESCE, rhino.Token_OPTCHAIN_GETPROP, rhino.Token_OPTCHAIN_GETELEM:
		u.ComputeMayUse(n.GetLastChild(), cfgNode, output /* conditional= */, true)
		u.ComputeMayUse(n.GetFirstChild(), cfgNode, output, conditional)
		return
	case rhino.Token_OPTCHAIN_CALL:
		// As args are evaluated in AST order, we traverse in reverse AST order for backward
		// dataflow analysis.
		for c := n.GetLastChild(); c != n.GetFirstChild(); c = c.GetPrevious() {
			u.ComputeMayUse(c, cfgNode, output /* conditional= */, true)
		}
		u.ComputeMayUse(n.GetFirstChild(), cfgNode, output, conditional)
		return
	case rhino.Token_HOOK:
		u.ComputeMayUse(n.GetLastChild(), cfgNode, output /* conditional= */, true)
		u.ComputeMayUse(n.GetSecondChild(), cfgNode, output /* conditional= */, true)
		u.ComputeMayUse(n.GetFirstChild(), cfgNode, output, conditional)
		return
	case rhino.Token_VAR, rhino.Token_LET, rhino.Token_CONST:
		varName := n.GetFirstChild()
		jutil.CheckStateMsg1(n.HasChildren(), "AST should be normalized (%s)", n)

		if varName.IsDestructuringLhs() {
			// Note: since destructuring is evaluated in reverse AST order, we traverse the first
			// child before the second in order to do our backwards data flow analysis.
			u.ComputeMayUse(varName.GetFirstChild(), cfgNode, output, conditional)
			u.ComputeMayUse(varName.GetSecondChild(), cfgNode, output, conditional)
		} else if varName.HasChildren() {
			u.ComputeMayUse(varName.GetFirstChild(), cfgNode, output, conditional)
			if !conditional {
				u.RemoveFromUseIfLocal(varName.GetString(), output)
			}
		} // else var name declaration with no initial value
		return
	case rhino.Token_DEFAULT_VALUE:
		if n.GetFirstChild().IsDestructuringPattern() {
			u.ComputeMayUse(n.GetFirstChild(), cfgNode, output, conditional)
			u.ComputeMayUse(n.GetSecondChild(), cfgNode, output, true)
		} else if n.GetFirstChild().IsName() {
			// assigning to the name occurs after evaluating the default value
			if !conditional {
				u.RemoveFromUseIfLocal(n.GetFirstChild().GetString(), output)
			}
			u.ComputeMayUse(n.GetSecondChild(), cfgNode, output, true)
		} else {
			u.ComputeMayUse(n.GetSecondChild(), cfgNode, output, true)
			u.ComputeMayUse(n.GetFirstChild(), cfgNode, output, conditional)
		}
	default:
		if NodeUtil_IsAssignmentOp(n) && n.GetFirstChild().IsName() {
			jutil.CheckState(!NodeUtil_IsLogicalAssignmentOp(n))
			name := n.GetFirstChild()
			if !conditional {
				u.RemoveFromUseIfLocal(name.GetString(), output)
			}

			// In case of a += "Hello". There is a read of a.
			if !n.IsAssign() {
				u.AddToUseIfLocal(name.GetString(), cfgNode, output)
			}

			u.ComputeMayUse(name.GetNext(), cfgNode, output, conditional)
		} else if n.IsAssign() && n.GetFirstChild().IsDestructuringPattern() {
			// Note: the rhs of destructuring is evaluated before the lhs
			u.ComputeMayUse(n.GetFirstChild(), cfgNode, output, conditional)
			u.ComputeMayUse(n.GetSecondChild(), cfgNode, output, conditional)
		} else {
			/*
			 * We want to traverse in reverse order because we want the LAST
			 * definition in the sub-tree.
			 */
			for c := n.GetLastChild(); c != nil; c = c.GetPrevious() {
				u.ComputeMayUse(c, cfgNode, output, conditional)
			}
		}
	}
}

// AddToUseIfLocal sets the variable for the given name to the node value in
// the upward exposed lattice. Do nothing if the variable name is one of the
// escaped variable.
//
// Java: MaybeReachingVariableUse.addToUseIfLocal(String, Node, ReachingUses)
func (u *MaybeReachingVariableUse) AddToUseIfLocal(name string, node *rhino.Node, use *MaybeReachingVariableUse_ReachingUses) {
	var_ := u.allVarsInFn.Get(name)
	if var_ == nil {
		return
	}
	if !u.escaped.Contains(var_) {
		use.Put(var_, node)
	}
}

// RemoveFromUseIfLocal removes the variable for the given name from the node
// value in the upward exposed lattice. Do nothing if the variable name is
// one of the escaped variable.
//
// Java: MaybeReachingVariableUse.removeFromUseIfLocal(String, ReachingUses)
func (u *MaybeReachingVariableUse) RemoveFromUseIfLocal(name string, use *MaybeReachingVariableUse_ReachingUses) {
	var_ := u.allVarsInFn.Get(name)
	if var_ == nil {
		return
	}
	if !u.escaped.Contains(var_) {
		use.RemoveAll(var_)
	}
}

// GetUses gets a list of nodes that may be using the value assigned to name
// in defNode. defNode must be one of the control flow graph nodes.
//
// name is the name of the variable. It can only be names of local variable
// that are not function parameters, escaped variables or variables declared
// in catch. defNode is the control flow graph node that may assign a value
// to name. Returns the list of upward exposed uses of the variable name at
// defNode.
//
// Java: MaybeReachingVariableUse.getUses(String, Node)
func (u *MaybeReachingVariableUse) GetUses(name string, defNode *rhino.Node) iter.Seq[*rhino.Node] {
	n := u.GetCfg().GetNode(defNode)
	jutil.CheckNotNull(n)
	state := n.GetAnnotation().(*DataFlowAnalysis_LinearFlowState[*MaybeReachingVariableUse_ReachingUses])
	return state.GetOut().Get(u.allVarsInFn.Get(name))
}
