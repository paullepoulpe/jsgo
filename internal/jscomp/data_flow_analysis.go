// Port of com/google/javascript/jscomp/DataFlowAnalysis.java @ closure-compiler v20261021

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
	"github.com/paullepoulpe/jsgo/internal/graph"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// DataFlowAnalysis_MAX_STEPS_PER_NODE is the maximum number of steps per
// individual CFG node before we assume the analysis is divergent.
//
// TODO(b/196398705): This is way too high. Find traversal ordering
// heurisitc that reduces it.
const DataFlowAnalysis_MAX_STEPS_PER_NODE = 20000

// dfaReuseMaxNodes is the largest flow state slice or queue that an analysis
// keeps for its next CFG.
const dfaReuseMaxNodes = 1024

// forwardCfgWorklist is a reusable, non-generic forward worklist over AST CFG
// nodes (min-heap on Priority when the CFG assigns node priorities, else FIFO).
type forwardCfgWorklist struct {
	queue       []*graph.LinkedDiGraphNode[rhino.Node, ControlFlowGraph_Branch]
	queueHead   int
	usePriority bool
}

func (w *forwardCfgWorklist) reset(cfg *ControlFlowGraph[rhino.Node], numNodes int) {
	clear(w.queue)
	w.queue = w.queue[:0]
	w.queueHead = 0
	w.usePriority = cfg.hasNodePriorities
	if cap(w.queue) < numNodes {
		w.queue = make([]*graph.LinkedDiGraphNode[rhino.Node, ControlFlowGraph_Branch], 0, numNodes)
	}
}

func (w *forwardCfgWorklist) releaseLarge() {
	if cap(w.queue) > dfaReuseMaxNodes {
		w.queue = nil
		w.queueHead = 0
	}
}

func (w *forwardCfgWorklist) notEmpty() bool {
	return len(w.queue) > w.queueHead
}

func (w *forwardCfgWorklist) enqueue(node *graph.LinkedDiGraphNode[rhino.Node, ControlFlowGraph_Branch]) {
	if node == nil {
		panic(jutil.NewNullPointerException())
	}
	if node.DfaInQueue {
		return
	}
	node.DfaInQueue = true
	w.queue = append(w.queue, node)
	if w.usePriority {
		i := len(w.queue) - 1
		xPrio := node.Priority
		if xPrio < 0 {
			panic(jutil.NewIllegalStateExceptionMsg("priority not set"))
		}
		for i > 0 {
			parent := (i - 1) / 2
			p := w.queue[parent]
			if xPrio >= p.Priority {
				break
			}
			w.queue[i] = p
			i = parent
		}
		w.queue[i] = node
	}
}

func (w *forwardCfgWorklist) dequeue() *graph.LinkedDiGraphNode[rhino.Node, ControlFlowGraph_Branch] {
	var curNode *graph.LinkedDiGraphNode[rhino.Node, ControlFlowGraph_Branch]
	if w.usePriority {
		n := len(w.queue)
		curNode = w.queue[0]
		last := w.queue[n-1]
		w.queue[n-1] = nil
		w.queue = w.queue[:n-1]
		if n > 1 {
			m := n - 1
			xPrio := last.Priority
			i := 0
			for {
				pick := 2*i + 1
				if pick >= m {
					break
				}
				smaller := w.queue[pick]
				smallerPrio := smaller.Priority
				if sibling := pick + 1; sibling < m {
					sib := w.queue[sibling]
					if sibPrio := sib.Priority; smallerPrio > sibPrio {
						pick = sibling
						smaller = sib
						smallerPrio = sibPrio
					}
				}
				if xPrio <= smallerPrio {
					break
				}
				w.queue[i] = smaller
				i = pick
			}
			w.queue[i] = last
		}
	} else {
		curNode = w.queue[w.queueHead]
		w.queue[w.queueHead] = nil
		w.queueHead++
		if w.queueHead == len(w.queue) {
			w.queue = w.queue[:0]
			w.queueHead = 0
		}
	}
	curNode.DfaInQueue = false
	return curNode
}

// backwardCfgWorklist is a reusable, non-generic backward worklist over AST CFG
// nodes (max-heap on Priority when the CFG assigns node priorities, else FIFO).
type backwardCfgWorklist struct {
	queue       []*graph.LinkedDiGraphNode[rhino.Node, ControlFlowGraph_Branch]
	queueHead   int
	usePriority bool
}

func (w *backwardCfgWorklist) reset(cfg *ControlFlowGraph[rhino.Node], numNodes int) {
	clear(w.queue)
	w.queue = w.queue[:0]
	w.queueHead = 0
	w.usePriority = cfg.hasNodePriorities
	if cap(w.queue) < numNodes {
		w.queue = make([]*graph.LinkedDiGraphNode[rhino.Node, ControlFlowGraph_Branch], 0, numNodes)
	}
}

func (w *backwardCfgWorklist) releaseLarge() {
	if cap(w.queue) > dfaReuseMaxNodes {
		w.queue = nil
		w.queueHead = 0
	}
}

func (w *backwardCfgWorklist) notEmpty() bool {
	return len(w.queue) > w.queueHead
}

func (w *backwardCfgWorklist) enqueue(node *graph.LinkedDiGraphNode[rhino.Node, ControlFlowGraph_Branch]) {
	if node == nil {
		panic(jutil.NewNullPointerException())
	}
	if node.DfaInQueue {
		return
	}
	node.DfaInQueue = true
	w.queue = append(w.queue, node)
	if w.usePriority {
		i := len(w.queue) - 1
		xPrio := node.Priority
		if xPrio < 0 {
			panic(jutil.NewIllegalStateExceptionMsg("priority not set"))
		}
		for i > 0 {
			parent := (i - 1) / 2
			p := w.queue[parent]
			if xPrio <= p.Priority {
				break
			}
			w.queue[i] = p
			i = parent
		}
		w.queue[i] = node
	}
}

func (w *backwardCfgWorklist) dequeue() *graph.LinkedDiGraphNode[rhino.Node, ControlFlowGraph_Branch] {
	var curNode *graph.LinkedDiGraphNode[rhino.Node, ControlFlowGraph_Branch]
	if w.usePriority {
		n := len(w.queue)
		curNode = w.queue[0]
		last := w.queue[n-1]
		w.queue[n-1] = nil
		w.queue = w.queue[:n-1]
		if n > 1 {
			m := n - 1
			xPrio := last.Priority
			i := 0
			for {
				pick := 2*i + 1
				if pick >= m {
					break
				}
				larger := w.queue[pick]
				largerPrio := larger.Priority
				if sibling := pick + 1; sibling < m {
					sib := w.queue[sibling]
					if sibPrio := sib.Priority; largerPrio < sibPrio {
						pick = sibling
						larger = sib
						largerPrio = sibPrio
					}
				}
				if xPrio >= largerPrio {
					break
				}
				w.queue[i] = larger
				i = pick
			}
			w.queue[i] = last
		}
	} else {
		curNode = w.queue[w.queueHead]
		w.queue[w.queueHead] = nil
		w.queueHead++
		if w.queueHead == len(w.queue) {
			w.queue = w.queue[:0]
			w.queueHead = 0
		}
	}
	curNode.DfaInQueue = false
	return curNode
}

// dataFlowAnalysis_linearFlowStateOf is the unchecked cast
// `LinearFlowState<L> state = node.getAnnotation()`.
//
// Go-only: language-adapter — Java's generic getAnnotation() casts at the
// call site; a null annotation stays null.
func dataFlowAnalysis_linearFlowStateOf[L any](a graph.Annotatable) *DataFlowAnalysis_LinearFlowState[L] {
	annotation := a.GetAnnotation()
	if annotation == nil {
		return nil
	}
	return annotation.(*DataFlowAnalysis_LinearFlowState[L])
}

// DataFlowAnalysis_LinearFlowState holds the in and out states of a node.
type DataFlowAnalysis_LinearFlowState[L any] struct {
	stepCount int
	in        L
	out       L
}

// Java: DataFlowAnalysis.LinearFlowState(L, L)
func newDataFlowAnalysis_LinearFlowState[L any](in L, out L) *DataFlowAnalysis_LinearFlowState[L] {
	jutil.CheckNotNull(in)
	jutil.CheckNotNull(out)
	return &DataFlowAnalysis_LinearFlowState[L]{stepCount: 0, in: in, out: out}
}

// Java: DataFlowAnalysis.LinearFlowState.getStepCount()
func (s *DataFlowAnalysis_LinearFlowState[L]) GetStepCount() int {
	return s.stepCount
}

// Java: DataFlowAnalysis.LinearFlowState.getIn()
func (s *DataFlowAnalysis_LinearFlowState[L]) GetIn() L {
	return s.in
}

// Java: DataFlowAnalysis.LinearFlowState.setIn(L)
func (s *DataFlowAnalysis_LinearFlowState[L]) SetIn(in L) {
	jutil.CheckNotNull(in)
	s.in = in
}

// Java: DataFlowAnalysis.LinearFlowState.getOut()
func (s *DataFlowAnalysis_LinearFlowState[L]) GetOut() L {
	return s.out
}

// Java: DataFlowAnalysis.LinearFlowState.setOut(L)
func (s *DataFlowAnalysis_LinearFlowState[L]) SetOut(out L) {
	jutil.CheckNotNull(out)
	s.out = out
}

// Java: DataFlowAnalysis.LinearFlowState.toString()
func (s *DataFlowAnalysis_LinearFlowState[L]) String() string {
	return jutil.StringFormat("IN: %s OUT: %s", s.in, s.out)
}

func dfaHasNestedFunctionOrClass(n rhino.Node) bool {
	for c := n.GetFirstChild(); !c.IsNil(); c = c.GetNext() {
		tok := c.GetToken()
		if tok == rhino.Token_FUNCTION || tok == rhino.Token_CLASS {
			return true
		}
		if dfaHasNestedFunctionOrClass(c) {
			return true
		}
	}
	return false
}

// DataFlowAnalysis_ComputeEscaped computes the set of escaped variables.
// When a variable is escaped in a dataflow analysis, it can be referenced
// outside of the code that we are analyzing. A variable is escaped if any
// of the following is true:
//
// 1. Exported variables as they can be needed after the script terminates.
// 2. Names of named functions because in JavaScript, function foo(){} does
// not kill foo in the dataflow.
//
// jsScope must be a function scope.
//
// Java: DataFlowAnalysis.computeEscaped(Scope, Set<Var>, AbstractCompiler, ScopeCreator, Map<String, Var>)
func DataFlowAnalysis_ComputeEscaped(
	jsScope *Scope,
	escaped *jutil.ObjectLinkedHashSet[*Var],
	compiler AbstractCompiler,
	scopeCreator ScopeCreator,
	allVarsInFn *jutil.LinkedHashMap[string, *Var]) {

	// Optimize for Class Static Blocks.
	jutil.CheckArgument(jsScope.IsFunctionScope())
	if allVarsInFn != nil && allVarsInFn.IsEmpty() {
		return
	}

	escapedStartedEmpty := escaped.IsEmpty()
	hasNestedScopeVars := false

	// TODO (simranarora) catch variables should not be considered escaped in
	// ES6. Getting rid of the catch check is causing breakages however
	codingConvention := compiler.GetCodingConvention()
	if allVarsInFn != nil {
		for var_ := range allVarsInFn.Values() {
			if var_.scope != jsScope && (var_.scope == nil || !var_.scope.IsFunctionBlockScope()) {
				hasNestedScopeVars = true
			}
			if var_.GetParentNode().IsCatch() ||
				codingConvention.IsExported(var_.GetName()) {
				escaped.Add(var_)
			}
		}
	}

	fnRoot := jsScope.GetRootNode()
	fnNameNode := fnRoot.GetFirstChild()
	fnName := fnNameNode.GetString()
	fnNameShadowed := false
	if fnName != "" && allVarsInFn != nil && jsScope.HasOwnSlot(fnName) {
		if v := allVarsInFn.Get(fnName); v == nil || v.scope != jsScope {
			fnNameShadowed = true
		}
	}

	if allVarsInFn != nil && !hasNestedScopeVars && !fnNameShadowed && escapedStartedEmpty && escaped.Size() >= allVarsInFn.Size() {
		return
	}

	args := fnNameNode.GetNext()
	fnBody := args.GetNext()
	argsHasFnOrClass := dfaHasNestedFunctionOrClass(args)

	var firstEscapingStmt rhino.Node
	if !argsHasFnOrClass && fnBody.IsBlock() {
		for stmt := fnBody.GetFirstChild(); !stmt.IsNil(); stmt = stmt.GetNext() {
			if stmt.IsFunction() || stmt.IsClass() || dfaHasNestedFunctionOrClass(stmt) {
				firstEscapingStmt = stmt
				break
			}
		}
		if firstEscapingStmt.IsNil() {
			return
		}
	} else if !argsHasFnOrClass && !dfaHasNestedFunctionOrClass(fnBody) {
		return
	}

	hasDuplicateVarNames := allVarsInFn == nil || hasNestedScopeVars || fnNameShadowed || argsHasFnOrClass
	nonEscapingDepth := jsScope.GetDepth()
	var skipArgs rhino.Node
	var filterFnBody rhino.Node
	if !argsHasFnOrClass {
		skipArgs = args
		if fnBody.IsBlock() {
			nonEscapingDepth = jsScope.GetDepth() + 1
			filterFnBody = fnBody
		}
	}

	finder := dataFlowAnalysis_escapedFinder{
		jsScope:              jsScope,
		nonEscapingDepth:     nonEscapingDepth,
		escaped:              escaped,
		allVarsInFn:          allVarsInFn,
		skipArgs:             skipArgs,
		fnBody:               filterFnBody,
		firstEscapingStmt:    firstEscapingStmt,
		hasDuplicateVarNames: hasDuplicateVarNames,
		escapedStartedEmpty:  escapedStartedEmpty,
	}

	b := NodeTraversal_Builder{
		compiler:     compiler,
		callback:     &finder,
		scopeCreator: scopeCreator,
	}
	b.TraverseAtScope(jsScope)
}

// dataFlowAnalysis_escapedFinder is the anonymous AbstractPostOrderCallback
// of computeEscaped.
type dataFlowAnalysis_escapedFinder struct {
	jsScope              *Scope
	nonEscapingDepth     int
	escaped              *jutil.ObjectLinkedHashSet[*Var]
	allVarsInFn          *jutil.LinkedHashMap[string, *Var]
	skipArgs             rhino.Node
	fnBody               rhino.Node
	firstEscapingStmt    rhino.Node
	reachedFirstEscaping bool
	hasDuplicateVarNames bool
	escapedStartedEmpty  bool
}

func (f *dataFlowAnalysis_escapedFinder) ShouldTraverse(t *NodeTraversal, n rhino.Node, parent rhino.Node) bool {
	if !f.hasDuplicateVarNames && f.escapedStartedEmpty && f.escaped.Size() >= f.allVarsInFn.Size() {
		return false
	}
	if n == f.skipArgs {
		return false
	}
	if !parent.IsNil() && parent == f.fnBody {
		if !f.reachedFirstEscaping {
			if n != f.firstEscapingStmt {
				return false
			}
			f.reachedFirstEscaping = true
			return true
		}
		return n.IsFunction() || n.IsClass() || dfaHasNestedFunctionOrClass(n)
	}
	return true
}

func (f *dataFlowAnalysis_escapedFinder) isInsideEscapingScope(t *NodeTraversal) bool {
	for i := len(t.scopes) - 1; i > f.nonEscapingDepth; i-- {
		root := t.scopes[i].node
		switch root.GetToken() {
		case rhino.Token_FUNCTION, rhino.Token_SCRIPT, rhino.Token_MODULE_BODY, rhino.Token_ROOT:
			return true
		case rhino.Token_BLOCK:
			if NodeUtil_IsClassStaticBlock(root) {
				return true
			}
		case rhino.Token_MEMBER_FIELD_DEF, rhino.Token_COMPUTED_FIELD_DEF:
			if !root.IsStaticMember() {
				return true
			}
		}
	}
	return false
}

func (*dataFlowAnalysis_escapedFinder) onlyNameLeaves() {}

// Java: DataFlowAnalysis.computeEscaped(...).new AbstractPostOrderCallback() {...}.visit(NodeTraversal, Node, Node)
func (f *dataFlowAnalysis_escapedFinder) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {

	if !n.IsName() || parent.IsFunction() || len(t.scopes)-1 <= f.nonEscapingDepth || !f.isInsideEscapingScope(t) {
		return
	}

	name := n.GetString()
	if f.allVarsInFn != nil {
		targetVar := f.allVarsInFn.Get(name)
		if targetVar == nil {
			return
		}
		if !f.hasDuplicateVarNames && f.escaped.Contains(targetVar) {
			return
		}
	}

	var_ := t.LookupVar(name)
	if var_ == nil {
		return
	}

	// AbstractVar.getScope() returns S (here Scope); see abstract_scope.go.
	varScope, _ := var_.GetScope().(*Scope)
	variableCfgScope := varScope.GetClosestCfgRootScope()

	// If the current examined scope and var decleration scope are different
	// then var is 'not escaped'.
	if variableCfgScope != f.jsScope {
		return
	}

	// Variable is referenced outside of declaration scope it is 'escaped'.
	// Non-static class field initializers are evaluated at instance creation time
	// (like closures), so variables referenced in non-static class fields must also be
	// marked as escaped.
	f.escaped.Add(var_)
}

// Java: DataFlowAnalysis.isReferencedInNonStaticClassField(Scope, Scope)
func DataFlowAnalysis_IsReferencedInNonStaticClassField(refScope *Scope, targetScope *Scope) bool {
	for s := refScope; s != nil && s != targetScope; s = s.GetParent() {
		root := s.GetRootNode()
		if (root.IsMemberFieldDef() || root.IsComputedFieldDef()) && !root.IsStaticMember() {
			return true
		}
	}
	return false
}
