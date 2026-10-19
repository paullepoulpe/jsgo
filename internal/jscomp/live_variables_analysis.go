// Port of com/google/javascript/jscomp/LiveVariablesAnalysis.java @ closure-compiler v20261015

/*
 * Copyright 2017 The Closure Compiler Authors.
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

// LiveVariablesAnalysis computes the "liveness" of all local variables. A
// variable is "live" at a point of a program if the value it is currently
// holding might be read later. Otherwise, the variable is considered "dead"
// if we know for sure that it will no longer be read. Dead variables are
// candidates for dead assignment elimination and variable name sharing. The
// worst case safe assumption is to assume that all variables are live. In
// that case, we will have no opportunity for optimizations. This is
// especially the case within a TRY block when an assignment is not
// guaranteed to take place. We bail out by assuming that all variables are
// live.
//
// Due to the possibility of inner functions and closures, certain "local"
// variables can escape the function. These variables will be considered as
// global and they can be retrieved with GetEscapedLocals().
type LiveVariablesAnalysis struct {
	DataFlowAnalysisBase[*rhino.Node, *LiveVariablesAnalysis_LiveVariableLattice]

	// The scope of the function that we are analyzing.
	jsScope *Scope

	// The scope of the body of the function that we are analyzing.
	jsScopeChild *Scope
	escaped      *jutil.ObjectLinkedHashSet[*Var]

	// Maps the variable name to it's position in this jsScope were we to
	// combine the function and function body scopes. The Integer represents
	// the equivalent of the variable index property within a scope
	scopeVariables *jutil.LinkedHashMap[string, int]

	// obtain variables in the order in which they appear in the code
	orderedVars []*Var

	allVarsInFn *jutil.LinkedHashMap[string, *Var]

	escapedBits jutil.BitSet
	genScratch  jutil.BitSet
	killScratch jutil.BitSet
}

const LiveVariablesAnalysis_MAX_VARIABLES_TO_ANALYZE = 100

// LiveVariablesAnalysis_LiveVariableJoinOp is the inner class
// LiveVariableJoinOp.
type LiveVariablesAnalysis_LiveVariableJoinOp struct {
	result *LiveVariablesAnalysis_LiveVariableLattice
}

// Java: LiveVariablesAnalysis.LiveVariableJoinOp()
func NewLiveVariablesAnalysis_LiveVariableJoinOp(outer *LiveVariablesAnalysis) *LiveVariablesAnalysis_LiveVariableJoinOp {
	return &LiveVariablesAnalysis_LiveVariableJoinOp{
		result: NewLiveVariablesAnalysis_LiveVariableLattice(len(outer.orderedVars)),
	}
}

// Java: LiveVariablesAnalysis.LiveVariableJoinOp.joinFlow(LiveVariableLattice)
func (j *LiveVariablesAnalysis_LiveVariableJoinOp) JoinFlow(x *LiveVariablesAnalysis_LiveVariableLattice) {
	j.result.liveSet.Or(x.liveSet)
}

// Java: LiveVariablesAnalysis.LiveVariableJoinOp.finish()
func (j *LiveVariablesAnalysis_LiveVariableJoinOp) Finish() *LiveVariablesAnalysis_LiveVariableLattice {
	return j.result
}

// LiveVariablesAnalysis_LiveVariableLattice is the lattice that stores the
// liveness of all local variables at a given point in the program. The whole
// lattice is the power set of all local variables and a variable is live if
// it is in the set.
type LiveVariablesAnalysis_LiveVariableLattice struct {
	liveSet  *jutil.BitSet
	liveBits jutil.BitSet
}

// numVars is the number of all local variables.
//
// Java: LiveVariablesAnalysis.LiveVariableLattice(int)
func NewLiveVariablesAnalysis_LiveVariableLattice(numVars int) *LiveVariablesAnalysis_LiveVariableLattice {
	l := &LiveVariablesAnalysis_LiveVariableLattice{}
	l.liveBits.InitWithSize(numVars)
	l.liveSet = &l.liveBits
	return l
}

// Java: LiveVariablesAnalysis.LiveVariableLattice(LiveVariableLattice)
func NewLiveVariablesAnalysis_LiveVariableLatticeOf(other *LiveVariablesAnalysis_LiveVariableLattice) *LiveVariablesAnalysis_LiveVariableLattice {
	jutil.CheckNotNull(other)
	l := &LiveVariablesAnalysis_LiveVariableLattice{}
	other.liveSet.CloneInto(&l.liveBits)
	l.liveSet = &l.liveBits
	return l
}

// Java: LiveVariablesAnalysis.LiveVariableLattice.equals(Object)
func (l *LiveVariablesAnalysis_LiveVariableLattice) Equals(other any) bool {
	jutil.CheckNotNull(other)
	liveVariableLattice, ok := other.(*LiveVariablesAnalysis_LiveVariableLattice)
	return ok && l.liveSet.Equals(liveVariableLattice.liveSet)
}

// IsLive: there is only a version of this function with index since
// var.index will return the wrong one. Use an instantiation of
// LiveVariablesAnalysis and getVarIndex(var) to get the right index.
//
// Java: LiveVariablesAnalysis.LiveVariableLattice.isLive(int)
func (l *LiveVariablesAnalysis_LiveVariableLattice) IsLive(index int) bool {
	return l.liveSet.Get(index)
}

// Java: LiveVariablesAnalysis.LiveVariableLattice.toString()
func (l *LiveVariablesAnalysis_LiveVariableLattice) String() string {
	return l.liveSet.String()
}

// NextSetBit returns the index of the first bit that is set to true that
// occurs on or after the specified starting index.
//
// Java: LiveVariablesAnalysis.LiveVariableLattice.nextSetBit(int)
func (l *LiveVariablesAnalysis_LiveVariableLattice) NextSetBit(fromIndex int) int {
	return l.liveSet.NextSetBit(fromIndex)
}

// Java: LiveVariablesAnalysis.LiveVariableLattice.hashCode()
func (l *LiveVariablesAnalysis_LiveVariableLattice) HashCode() int32 {
	return l.liveSet.HashCode()
}

// NewLiveVariablesAnalysis is Live Variables Analysis using the ES6 scope
// creator. This analysis should only be done on function where jsScope is
// the function scope. If we call LiveVariablesAnalysis from the function
// scope of our pass, we can pass a null value for the JsScopeChild, but if
// we call it from the function block scope, then JsScopeChild will be the
// function block scope.
//
// We call from the function scope when the pass requires us to traverse
// nodes beginning at the function parameters, and it from the function block
// scope when we are ignoring function parameters.
//
// jsScope is the function scope; jsScopeChild null or function block scope;
// scopeCreator Es6 Scope creator; allVarsDeclaredInFunction mapping of names
// to vars of everything reachable in a function.
//
// Java: LiveVariablesAnalysis(ControlFlowGraph<Node>, Scope, Scope, AbstractCompiler, ScopeCreator, AllVarsDeclaredInFunction)
func NewLiveVariablesAnalysis(
	cfg ControlFlowGraph[*rhino.Node],
	jsScope *Scope,
	jsScopeChild *Scope,
	compiler AbstractCompiler,
	scopeCreator ScopeCreator,
	allVarsDeclaredInFunction *NodeUtil_AllVarsDeclaredInFunction) *LiveVariablesAnalysis {
	l := &LiveVariablesAnalysis{}
	l.InitDataFlowAnalysisBase(l, cfg)
	jutil.CheckStateMsg1(jsScope.IsFunctionScope(), "%s", jsScope)

	l.jsScope = jsScope
	l.jsScopeChild = jsScopeChild
	l.escaped = jutil.NewObjectLinkedHashSet[*Var]()
	l.orderedVars = allVarsDeclaredInFunction.GetAllVariablesInOrder()
	l.allVarsInFn = allVarsDeclaredInFunction.GetAllVariables()
	l.scopeVariables = jutil.NewLinkedHashMapWithCapacity[string, int](len(l.orderedVars))

	DataFlowAnalysis_ComputeEscaped(jsScope, l.escaped, compiler, scopeCreator, l.allVarsInFn)

	l.AddScopeVariables()
	return l
}

// AddScopeVariables: parameters belong to the function scope, but variables
// defined in the function body belong to the function body scope. Assign a
// unique index to each variable, regardless of which scope it's in.
//
// Java: LiveVariablesAnalysis.addScopeVariables()
func (l *LiveVariablesAnalysis) AddScopeVariables() {
	l.escapedBits.InitWithSize(len(l.orderedVars))
	num := 0
	for _, v := range l.orderedVars {
		name := v.GetName()
		l.scopeVariables.Put(name, num)
		localScope, _ := v.GetScope().(*Scope)
		var local bool
		if localScope.IsFunctionBlockScope() {
			local = LiveVariablesAnalysis_IsDeclaredInFunctionBlockOrParameter(localScope, name)
		} else if localScope == l.jsScope && l.jsScopeChild != nil {
			local = LiveVariablesAnalysis_IsDeclaredInFunctionBlockOrParameter(l.jsScopeChild, name)
		} else {
			local = localScope.HasOwnSlot(name)
		}
		if !local || l.escaped.Contains(v) {
			l.escapedBits.Set(num)
		}
		num++
	}
}

// Java: LiveVariablesAnalysis.getEscapedLocals()
func (l *LiveVariablesAnalysis) GetEscapedLocals() *jutil.ObjectLinkedHashSet[*Var] {
	return l.escaped
}

// Java: LiveVariablesAnalysis.getAllVariables()
func (l *LiveVariablesAnalysis) GetAllVariables() *jutil.LinkedHashMap[string, *Var] {
	return l.allVarsInFn
}

// Java: LiveVariablesAnalysis.getAllVariablesInOrder()
func (l *LiveVariablesAnalysis) GetAllVariablesInOrder() []*Var {
	return l.orderedVars
}

// Java: LiveVariablesAnalysis.getVarIndex(String)
func (l *LiveVariablesAnalysis) GetVarIndex(var_ string) int {
	index, ok := l.scopeVariables.Lookup(var_)
	if !ok {
		// Unboxing the null Integer that Map.get returns.
		panic(jutil.NewNullPointerException())
	}
	return index
}

// Java: LiveVariablesAnalysis.isForward()
func (l *LiveVariablesAnalysis) IsForward() bool {
	return false
}

// Java: LiveVariablesAnalysis.createEntryLattice()
func (l *LiveVariablesAnalysis) CreateEntryLattice() *LiveVariablesAnalysis_LiveVariableLattice {
	return NewLiveVariablesAnalysis_LiveVariableLattice(len(l.orderedVars))
}

// Java: LiveVariablesAnalysis.createInitialEstimateLattice()
func (l *LiveVariablesAnalysis) CreateInitialEstimateLattice() *LiveVariablesAnalysis_LiveVariableLattice {
	return NewLiveVariablesAnalysis_LiveVariableLattice(len(l.orderedVars))
}

// Java: LiveVariablesAnalysis.createFlowJoiner()
func (l *LiveVariablesAnalysis) CreateFlowJoiner() DataFlowAnalysis_FlowJoiner[*LiveVariablesAnalysis_LiveVariableLattice] {
	return NewLiveVariablesAnalysis_LiveVariableJoinOp(l)
}

// Java: LiveVariablesAnalysis.flowThrough(Node, LiveVariableLattice)
func (l *LiveVariablesAnalysis) FlowThrough(node *rhino.Node, input *LiveVariablesAnalysis_LiveVariableLattice) *LiveVariablesAnalysis_LiveVariableLattice {
	size := input.liveSet.Size()
	l.genScratch.InitWithSize(size)
	l.killScratch.InitWithSize(size)

	// Make kills conditional if the node can end abruptly by an exception.
	conditional := false
	edgeList := l.GetCfg().GetOutEdges(node)
	for _, edge := range edgeList {
		if ControlFlowGraph_Branch_ON_EX == edge.GetValue() {
			conditional = true
		}
	}
	l.ComputeGenKill(node, &l.genScratch, &l.killScratch, conditional)
	result := NewLiveVariablesAnalysis_LiveVariableLatticeOf(input)
	// L_in = L_out - Kill + Gen
	result.liveSet.AndNot(&l.killScratch)
	result.liveSet.Or(&l.genScratch)
	return result
}

// ComputeGenKill computes the GEN and KILL set.
//
// n is the root node. gen: local variables that are live because of the
// instruction at n will be added to this set. kill: local variables that are
// killed because of the instruction at n will be added to this set.
// conditional: true if any assignments encountered are conditionally
// executed. These assignments might not kill a variable.
//
// Java: LiveVariablesAnalysis.computeGenKill(Node, BitSet, BitSet, boolean)
func (l *LiveVariablesAnalysis) ComputeGenKill(n *rhino.Node, gen *jutil.BitSet, kill *jutil.BitSet, conditional bool) {

	switch n.GetToken() {
	case rhino.Token_SCRIPT, rhino.Token_ROOT, rhino.Token_FUNCTION, rhino.Token_BLOCK:
		return
	case rhino.Token_WHILE, rhino.Token_DO, rhino.Token_IF, rhino.Token_FOR:
		l.ComputeGenKill(NodeUtil_GetConditionExpression(n), gen, kill, conditional)
		return
	case rhino.Token_FOR_OF, rhino.Token_FOR_AWAIT_OF, rhino.Token_FOR_IN:
		{
			// for (x in y) {...}
			lhs := n.GetFirstChild()
			if NodeUtil_IsNameDeclaration(lhs) {
				// for (var x in y) {...}
				lhs = lhs.GetLastChild()
			}

			// Note that the LHS may never be assigned to or evaluated, like in:
			//   for (x in []) {}
			// so should not be killed.
			l.ComputeGenKill(lhs, gen, kill, conditional)

			// rhs is executed only once so we don't go into it every loop.
			return
		}
	case rhino.Token_LET, rhino.Token_CONST, rhino.Token_VAR:
		for c := n.GetFirstChild(); c != nil; c = c.GetNext() {
			if c.IsName() {
				if c.HasChildren() {
					l.ComputeGenKill(c.GetFirstChild(), gen, kill, conditional)
					if !conditional {
						l.AddToSetIfLocal(c, kill)
					}
				}
			} else {
				jutil.CheckStateMsg1(c.IsDestructuringLhs(), "%s", c)
				if !conditional {
					NodeUtil_VisitLhsNodesInNode(c, func(lhsNode *rhino.Node) { l.AddToSetIfLocal(lhsNode, kill) })
				}
				l.ComputeGenKill(c.GetFirstChild(), gen, kill, conditional)
				l.ComputeGenKill(c.GetSecondChild(), gen, kill, conditional)
			}
		}
		return
	case rhino.Token_AND, rhino.Token_OR, rhino.Token_COALESCE, rhino.Token_OPTCHAIN_GETELEM, rhino.Token_OPTCHAIN_GETPROP:
		l.ComputeGenKill(n.GetFirstChild(), gen, kill, conditional)
		// May short circuit.
		l.ComputeGenKill(n.GetLastChild(), gen, kill, true)
		return
	case rhino.Token_OPTCHAIN_CALL:
		l.ComputeGenKill(n.GetFirstChild(), gen, kill, conditional)
		// Unlike OPTCHAIN_GETPROP and OPTCHAIN_GETELEM, the OPTCHAIN_CALLs can have multiple
		// children on rhs which get executed conditionally
		for c := n.GetSecondChild(); c != nil; c = c.GetNext() {
			l.ComputeGenKill(c, gen, kill, true)
		}
		return
	case rhino.Token_HOOK:
		l.ComputeGenKill(n.GetFirstChild(), gen, kill, conditional)
		// Assume both sides are conditional.
		l.ComputeGenKill(n.GetSecondChild(), gen, kill, true)
		l.ComputeGenKill(n.GetLastChild(), gen, kill, true)
		return
	case rhino.Token_NAME:
		if n.GetString() == "arguments" {
			l.MarkAllParametersEscaped()
		} else if !NodeUtil_IsLhsByDestructuring(n) {
			// Only add names in destructuring patterns if they're not lvalues.
			// e.g. "x" in "const {foo = x} = obj;"
			l.AddToSetIfLocal(n, gen)
		}
		return
	default:
		if NodeUtil_IsAssignmentOp(n) && n.GetFirstChild().IsName() {
			lhs := n.GetFirstChild()
			if !conditional {
				l.AddToSetIfLocal(lhs, kill)
			}
			if !n.IsAssign() {
				// assignments such as a += 1 reads a.
				l.AddToSetIfLocal(lhs, gen)
			}
			l.ComputeGenKill(lhs.GetNext(), gen, kill, conditional)
		} else if n.IsAssign() && n.GetFirstChild().IsDestructuringPattern() {
			if !conditional {
				NodeUtil_VisitLhsNodesInNode(
					n,
					func(child *rhino.Node) {
						if child.IsName() {
							l.AddToSetIfLocal(child, kill)
						}
					})
			}
			l.ComputeGenKill(n.GetFirstChild(), gen, kill, conditional)
			l.ComputeGenKill(n.GetSecondChild(), gen, kill, conditional)
		} else {
			for c := n.GetFirstChild(); c != nil; c = c.GetNext() {
				l.ComputeGenKill(c, gen, kill, conditional)
			}
		}
		return
	}
}

// Java: LiveVariablesAnalysis.addToSetIfLocal(Node, BitSet)
func (l *LiveVariablesAnalysis) AddToSetIfLocal(node *rhino.Node, set *jutil.BitSet) {
	jutil.CheckStateMsg1(node.IsName(), "%s", node)
	idx, ok := l.scopeVariables.Lookup(node.GetString())
	if !ok || l.escapedBits.Get(idx) {
		return
	}
	set.Set(idx)
}

// Java: LiveVariablesAnalysis.isDeclaredInFunctionBlockOrParameter(Scope, String)
func LiveVariablesAnalysis_IsDeclaredInFunctionBlockOrParameter(scope *Scope, name string) bool {
	// In ES6, we create a separate container scope above the function block scope to handle
	// default parameters. Since nothing in the function block scope is allowed to shadow
	// the variables in the function scope, we treat the two scopes as one in this method.
	jutil.CheckState(scope.IsFunctionBlockScope())
	return scope.HasOwnSlot(name) || scope.GetParent().HasOwnSlot(name)
}

// MarkAllParametersEscaped gives up computing liveness of formal parameters
// by putting all the simple parameters in the escaped set.
//
// This only applies to simple parameters, that is NAMEs, because other
// parameter syntaxes never need to be escaped in this way. The known
// applications of this method are for uses of `arguments`, and for IE8. In a
// function with non-simple paremeters, `arguments` is not parameter-mapped,
// and so referencing it doesn't escape paremeters. IE8 just doess't support
// non-simple parameters.
//
// We could actaully continue tracking simple parameters if any parameter is
// non-simple, but it wasn't worth the complexity or cost to do so.
//
// See https://tc39.github.io/ecma262/#sec-functiondeclarationinstantiation
//
// Java: LiveVariablesAnalysis.markAllParametersEscaped()
func (l *LiveVariablesAnalysis) MarkAllParametersEscaped() {
	paramList := NodeUtil_GetFunctionParameters(l.jsScope.GetRootNode())
	for param := paramList.GetFirstChild(); param != nil; param = param.GetNext() {
		if param.IsName() {
			name := param.GetString()
			l.escaped.Add(l.jsScope.GetVar(name))
			if idx, ok := l.scopeVariables.Lookup(name); ok {
				l.escapedBits.Set(idx)
			}
		}
	}
}
