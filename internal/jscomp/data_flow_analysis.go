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
	"github.com/paullepoulpe/jsgo/internal/ast"
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

// cfgWorklist is a reusable, non-generic worklist over AST CFG nodes (min-heap
// on Priority when forward is true and the CFG assigns node priorities, max-heap
// when forward is false, else FIFO).
type cfgWorklist struct {
	queue       []*graph.LinkedDiGraphNode[rhino.Node, ControlFlowGraph_Branch]
	queueHead   int
	usePriority bool
	forward     bool
}

func (w *cfgWorklist) reset(cfg *ControlFlowGraph[rhino.Node], numNodes int, forward bool) {
	clear(w.queue)
	w.queue = w.queue[:0]
	w.queueHead = 0
	w.usePriority = cfg.hasNodePriorities
	w.forward = forward
	if cap(w.queue) < numNodes {
		w.queue = make([]*graph.LinkedDiGraphNode[rhino.Node, ControlFlowGraph_Branch], 0, numNodes)
	}
}

func (w *cfgWorklist) releaseLarge(maxCap int) {
	if cap(w.queue) > maxCap {
		w.queue = nil
		w.queueHead = 0
	}
}

func (w *cfgWorklist) notEmpty() bool {
	return len(w.queue) > w.queueHead
}

func (w *cfgWorklist) enqueue(node *graph.LinkedDiGraphNode[rhino.Node, ControlFlowGraph_Branch]) {
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
		if w.forward {
			for i > 0 {
				parent := (i - 1) / 2
				p := w.queue[parent]
				if xPrio >= p.Priority {
					break
				}
				w.queue[i] = p
				i = parent
			}
		} else {
			for i > 0 {
				parent := (i - 1) / 2
				p := w.queue[parent]
				if xPrio <= p.Priority {
					break
				}
				w.queue[i] = p
				i = parent
			}
		}
		w.queue[i] = node
	}
}

func (w *cfgWorklist) dequeue() *graph.LinkedDiGraphNode[rhino.Node, ControlFlowGraph_Branch] {
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
			if w.forward {
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
			} else {
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

type forwardCfgWorklist struct {
	cfgWorklist
}

func (w *forwardCfgWorklist) reset(cfg *ControlFlowGraph[rhino.Node], numNodes int) {
	w.cfgWorklist.reset(cfg, numNodes, true)
}

func (w *forwardCfgWorklist) releaseLarge() {
	w.cfgWorklist.releaseLarge(dfaReuseMaxNodes)
}

const fnVarTableInlineCap = 16

type varIDTable struct {
	inlineIDs [fnVarTableInlineCap]uint32
	idToIndex map[uint32]int
	count     int
}

func (t *varIDTable) reset() {
	if len(t.idToIndex) > dfaReuseMaxNodes {
		t.idToIndex = nil
	} else if len(t.idToIndex) > 0 {
		clear(t.idToIndex)
	}
	t.count = 0
}

func (t *varIDTable) initFromVars(vars []*Var) (int, bool) {
	n := len(vars)
	t.count = n
	if n == 0 {
		t.reset()
		return 0, false
	}
	if n <= fnVarTableInlineCap {
		unique := n
		for i, v := range vars {
			id := v.nameID
			t.inlineIDs[i] = id
			for j := 0; j < i; j++ {
				if t.inlineIDs[j] == id {
					unique--
					break
				}
			}
		}
		return unique, unique < n
	}
	if t.idToIndex == nil || len(t.idToIndex) > dfaReuseMaxNodes {
		t.idToIndex = make(map[uint32]int, max(n, 32))
	} else {
		clear(t.idToIndex)
	}
	for i, v := range vars {
		t.idToIndex[v.nameID] = i
	}
	unique := len(t.idToIndex)
	return unique, unique < n
}

func (t *varIDTable) initFromIDs(ids []uint32) {
	n := len(ids)
	t.count = n
	if n == 0 {
		t.reset()
		return
	}
	if n <= fnVarTableInlineCap {
		copy(t.inlineIDs[:n], ids)
		return
	}
	if t.idToIndex == nil || len(t.idToIndex) > dfaReuseMaxNodes {
		t.idToIndex = make(map[uint32]int, max(n, 32))
	} else {
		clear(t.idToIndex)
	}
	for i, id := range ids {
		t.idToIndex[id] = i
	}
}

func (t *varIDTable) lookup(id uint32) (int, bool) {
	if id == 0 || t.count == 0 {
		return 0, false
	}
	if t.count <= fnVarTableInlineCap {
		for i := t.count - 1; i >= 0; i-- {
			if t.inlineIDs[i] == id {
				return i, true
			}
		}
		return 0, false
	}
	idx, ok := t.idToIndex[id]
	return idx, ok
}

type fnVarTable struct {
	vars              []*Var
	ids               varIDTable
	escaped           inlineBitSet
	uniqueCount       int
	hasDuplicateNames bool
}

func (t *fnVarTable) reset() {
	clear(t.vars)
	if cap(t.vars) > dfaReuseMaxNodes {
		t.vars = nil
	} else {
		t.vars = t.vars[:0]
	}
	t.ids.reset()
	t.escaped.clear(true)
	t.uniqueCount = 0
	t.hasDuplicateNames = false
}

func (t *fnVarTable) isEmpty() bool {
	return len(t.vars) == 0
}

func (t *fnVarTable) appendScopeVars(currentScope *Scope) {
	n := currentScope.GetVarCount()
	if n == 0 {
		return
	}
	if need := len(t.vars) + n; need > cap(t.vars) {
		newCap := max(need, 2*cap(t.vars), 16)
		buf := make([]*Var, len(t.vars), newCap)
		copy(buf, t.vars)
		t.vars = buf
	}
	for _, v := range currentScope.varList {
		if v != nil {
			t.vars = append(t.vars, v)
		}
	}
}

func (t *fnVarTable) finalize() {
	t.uniqueCount, t.hasDuplicateNames = t.ids.initFromVars(t.vars)
}

func (t *fnVarTable) lookupIndexByID(id uint32) (int, bool) {
	return t.ids.lookup(id)
}

func (t *fnVarTable) lookupVarByID(id uint32) *Var {
	idx, ok := t.ids.lookup(id)
	if !ok {
		return nil
	}
	return t.vars[idx]
}

func (t *fnVarTable) resolveExactVarIndex(v *Var) (int, bool) {
	if v == nil {
		return -1, false
	}
	idx, ok := t.ids.lookup(v.nameID)
	if !ok {
		return -1, false
	}
	if t.vars[idx] == v || t.vars[idx].Equals(v) {
		return idx, true
	}
	if t.hasDuplicateNames {
		for i, cand := range t.vars {
			if cand == v || cand.Equals(v) {
				return i, true
			}
		}
	}
	return -1, false
}

// DataFlowAnalysis_LinearFlowState holds the in and out states of a node.
type DataFlowAnalysis_LinearFlowState[L any] struct {
	stepCount int
	in        L
	out       L
}

// Java: DataFlowAnalysis.LinearFlowState.getIn()
func (s *DataFlowAnalysis_LinearFlowState[L]) GetIn() L {
	return s.in
}

// Java: DataFlowAnalysis.LinearFlowState.getOut()
func (s *DataFlowAnalysis_LinearFlowState[L]) GetOut() L {
	return s.out
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

func computeEscapedInTable(
	jsScope *Scope,
	table *fnVarTable,
	compiler AbstractCompiler,
	scopeCreator ScopeCreator,
) {
	jutil.CheckArgument(jsScope.IsFunctionScope())
	if table.isEmpty() {
		return
	}

	hasNestedScopeVars := false

	// TODO (simranarora) catch variables should not be considered escaped in
	// ES6. Getting rid of the catch check is causing breakages however
	codingConvention := compiler.GetCodingConvention()
	for i, var_ := range table.vars {
		if table.hasDuplicateNames {
			if canonIdx, _ := table.lookupIndexByID(var_.nameID); canonIdx != i {
				continue
			}
		}
		if var_.scope != jsScope && (var_.scope == nil || !var_.scope.IsFunctionBlockScope()) {
			hasNestedScopeVars = true
		}
		if var_.GetParentNode().IsCatch() ||
			codingConvention.IsExported(var_.GetName()) {
			table.escaped.Set(i)
		}
	}

	fnRoot := jsScope.GetRootNode()
	fnNameNode := fnRoot.GetFirstChild()
	fnNameID := fnNameNode.GetStringID()
	fnNameShadowed := false
	if fnNameID > rhino.RhinoStringPool_EmptyID && jsScope.hasOwnSlotByID(fnNameID) {
		if v := table.lookupVarByID(fnNameID); v == nil || v.scope != jsScope {
			fnNameShadowed = true
		}
	}

	if !hasNestedScopeVars && !fnNameShadowed && table.escaped.Cardinality() >= table.uniqueCount {
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

	hasDuplicateVarNames := table.hasDuplicateNames || hasNestedScopeVars || fnNameShadowed || argsHasFnOrClass
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

	finder := fnVarEscapedFinder{
		jsScope:              jsScope,
		nonEscapingDepth:     nonEscapingDepth,
		table:                table,
		skipArgs:             skipArgs,
		fnBody:               filterFnBody,
		firstEscapingStmt:    firstEscapingStmt,
		hasDuplicateVarNames: hasDuplicateVarNames,
	}

	b := NodeTraversal_Builder{
		compiler:     compiler,
		callback:     &finder,
		scopeCreator: scopeCreator,
	}
	b.TraverseAtScope(jsScope)
}

var dfaEscapedFinderVisitTokens = func() ast.TokenSet {
	var s ast.TokenSet
	s.Add(rhino.Token_NAME)
	return s
}()

var (
	_ nodeTraversal_visitTokens    = (*fnVarEscapedFinder)(nil)
	_ nodeTraversal_onlyNameLeaves = (*fnVarEscapedFinder)(nil)
)

// fnVarEscapedFinder is the anonymous AbstractPostOrderCallback of
// computeEscaped.
type fnVarEscapedFinder struct {
	jsScope              *Scope
	nonEscapingDepth     int
	table                *fnVarTable
	skipArgs             rhino.Node
	fnBody               rhino.Node
	firstEscapingStmt    rhino.Node
	reachedFirstEscaping bool
	hasDuplicateVarNames bool
}

func (f *fnVarEscapedFinder) ShouldTraverse(t *NodeTraversal, n rhino.Node, parent rhino.Node) bool {
	if !f.hasDuplicateVarNames && f.table.escaped.Cardinality() >= f.table.uniqueCount {
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

func (f *fnVarEscapedFinder) isInsideEscapingScope(t *NodeTraversal) bool {
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

func (*fnVarEscapedFinder) visitTokens() *ast.TokenSet { return &dfaEscapedFinderVisitTokens }
func (*fnVarEscapedFinder) onlyNameLeaves()            {}

// Java: DataFlowAnalysis.computeEscaped(...).new AbstractPostOrderCallback() {...}.visit(NodeTraversal, Node, Node)
func (f *fnVarEscapedFinder) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	if !n.IsName() || parent.IsFunction() || len(t.scopes)-1 <= f.nonEscapingDepth || !f.isInsideEscapingScope(t) {
		return
	}

	id := n.GetStringID()
	if id <= rhino.RhinoStringPool_EmptyID {
		return
	}
	candIdx, ok := f.table.lookupIndexByID(id)
	if !ok {
		return
	}
	if !f.hasDuplicateVarNames && f.table.escaped.Get(candIdx) {
		return
	}

	var_ := t.LookupVarByID(id)
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
	if idx, ok := f.table.resolveExactVarIndex(var_); ok {
		f.table.escaped.Set(idx)
	}
}
