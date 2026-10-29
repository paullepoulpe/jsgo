// Port of com/google/javascript/jscomp/MustBeReachingVariableDef.java @ closure-compiler v20261021

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

var argumentsStringID = rhino.RhinoStringPool_InternID(Var_ARGUMENTS)

// mustDefBottom is the sentinel for BOTTOM in MustBeReachingVariableDef_MustDef.defs
// (multiple or conditional reaching definitions). A nil entry represents TOP.
var mustDefBottom MustBeReachingVariableDef_Definition

type mustDefNodeState struct {
	in        *MustBeReachingVariableDef_MustDef
	out       *MustBeReachingVariableDef_MustDef
	stepCount int
}

// MustBeReachingVariableDef computes must-be-reaching definition for all
// uses of each variable.
//
// A definition of A in `A = foo()` is a must-be-reaching definition of the
// use of A in `alert(A)` if all paths from entry node to the use pass
// through that definition and it is the last definition before the use.
//
// By definition, a must-be-reaching definition for a given use is always a
// single definition and it "dominates" that use (i.e. always must execute
// before that use).
type MustBeReachingVariableDef struct {
	cfg        *ControlFlowGraph[rhino.Node]
	worklist   cfgWorklist
	nodeStates []mustDefNodeState

	compiler     AbstractCompiler
	vars         []*Var
	idByNameID   map[uint32]int
	escapedBits  inlineBitSet
	paramBits    inlineBitSet
	catchBits    inlineBitSet
	hasParams    bool
	hasCatchVars bool

	emptyLattice MustBeReachingVariableDef_MustDef
	joinOp       MustBeReachingVariableDef_MustDefJoin

	flowInput  *MustBeReachingVariableDef_MustDef
	flowOutput *MustBeReachingVariableDef_MustDef

	// Go-only: performance — rewound for each CFG: allocMustDefEmpty and
	// allocDefinition re-initialize what they hand out and reuse its storage.
	mustDefs jutil.ReuseSlab[MustBeReachingVariableDef_MustDef]
	defs     jutil.ReuseSlab[MustBeReachingVariableDef_Definition]
}

// Java: MustBeReachingVariableDef(ControlFlowGraph<Node>, AbstractCompiler, Set<Var>, Map<String, Var>)
func NewMustBeReachingVariableDef(
	cfg *ControlFlowGraph[rhino.Node],
	compiler AbstractCompiler,
	escaped *jutil.ObjectLinkedHashSet[*Var],
	allVarsInFn *jutil.LinkedHashMap[string, *Var]) *MustBeReachingVariableDef {
	d := &MustBeReachingVariableDef{}
	d.Reset(cfg, compiler, escaped, allVarsInFn)
	return d
}

// Reset resets d for reuse on a new function scope while retaining slab and
// variable table capacities.
func (d *MustBeReachingVariableDef) Reset(
	cfg *ControlFlowGraph[rhino.Node],
	compiler AbstractCompiler,
	escaped *jutil.ObjectLinkedHashSet[*Var],
	allVarsInFn *jutil.LinkedHashMap[string, *Var]) {
	if cap(d.nodeStates) > dfaReuseMaxNodes {
		d.nodeStates = nil
	} else {
		clear(d.nodeStates)
		d.nodeStates = d.nodeStates[:0]
	}
	d.worklist.releaseLarge(dfaReuseMaxNodes)
	d.cfg = cfg
	d.compiler = compiler
	d.initVars(allVarsInFn, escaped)
	d.flowInput = nil
	d.flowOutput = nil
	d.mustDefs.Rewind()
	d.defs.Rewind()
}

func (d *MustBeReachingVariableDef) resetWithTable(
	cfg *ControlFlowGraph[rhino.Node],
	compiler AbstractCompiler,
	table *fnVarTable,
) {
	if cap(d.nodeStates) > dfaReuseMaxNodes {
		d.nodeStates = nil
	} else {
		clear(d.nodeStates)
		d.nodeStates = d.nodeStates[:0]
	}
	d.worklist.releaseLarge(dfaReuseMaxNodes)
	d.cfg = cfg
	d.compiler = compiler
	d.initVarsFromTable(table)
	d.flowInput = nil
	d.flowOutput = nil
	d.mustDefs.Rewind()
	d.defs.Rewind()
}

// GetCfg returns the control flow graph analyzed by MustBeReachingVariableDef.
//
// Java: DataFlowAnalysis.getCfg()
func (d *MustBeReachingVariableDef) GetCfg() *ControlFlowGraph[rhino.Node] {
	return d.cfg
}

func (d *MustBeReachingVariableDef) inStateAt(index uint32) *MustBeReachingVariableDef_MustDef {
	return d.nodeStates[index].in
}

// Join merges two MustDef lattices directly without interface dispatch.
//
// Java: DataFlowAnalysis.join(MustDef, MustDef)
func (d *MustBeReachingVariableDef) Join(latticeA, latticeB *MustBeReachingVariableDef_MustDef) *MustBeReachingVariableDef_MustDef {
	d.joinOp.result = d.allocMustDefEmpty()
	d.joinOp.JoinFlow(latticeA)
	d.joinOp.JoinFlow(latticeB)
	return d.joinOp.result
}

// Analyze computes a forward linear fixed-point flow solution over d.cfg.
//
// Java: DataFlowAnalysis.analyze()
func (d *MustBeReachingVariableDef) Analyze() {
	cfg := d.cfg
	entry := cfg.entry
	implicitReturn := cfg.implicitReturn
	cfgNodes := &cfg.Nodes

	d.worklist.reset(cfg, cfgNodes.Size(), true)
	nodeCap := cfg.NodeIndexCap()
	if cap(d.nodeStates) < nodeCap {
		d.nodeStates = make([]mustDefNodeState, nodeCap)
	} else {
		clear(d.nodeStates[:nodeCap])
		d.nodeStates = d.nodeStates[:nodeCap]
	}

	initial := jutil.CheckNotNullPtr(d.CreateInitialEstimateLattice())
	for node := range cfgNodes.Values() {
		node.DfaInQueue = false
		d.nodeStates[node.Index] = mustDefNodeState{in: initial, out: initial}
		if node != implicitReturn {
			d.worklist.enqueue(node)
		}
	}

	for d.worklist.notEmpty() {
		curNode := d.worklist.dequeue()
		curState := &d.nodeStates[curNode.Index]
		stepCount := curState.stepCount
		curState.stepCount++
		if stepCount > DataFlowAnalysis_MAX_STEPS_PER_NODE {
			panic(jutil.NewIllegalStateExceptionMsg("Dataflow analysis appears to diverge around: " + jutil.ToJavaString(curNode)))
		}

		d.joinNodeInputs(curNode, curState, entry)
		outBefore := curState.out
		outAfter := jutil.CheckNotNullPtr(d.FlowThrough(curNode.Value, curState.in))
		curState.out = outAfter
		if !mustDefEquals(outBefore, outAfter) {
			for _, edge := range curNode.OutEdgeList.Slice() {
				nextNode := edge.DestNode
				if !nextNode.DfaInQueue && nextNode != implicitReturn {
					d.worklist.enqueue(nextNode)
				}
			}
		}
	}

	if implicitReturn != nil {
		d.joinNodeInputs(implicitReturn, &d.nodeStates[implicitReturn.Index], entry)
	}
}

func (d *MustBeReachingVariableDef) joinNodeInputs(
	node *graph.LinkedDiGraphNode[rhino.Node, ControlFlowGraph_Branch],
	state *mustDefNodeState,
	entry *graph.LinkedDiGraphNode[rhino.Node, ControlFlowGraph_Branch],
) {
	if node == entry {
		state.in = jutil.CheckNotNullPtr(d.CreateEntryLattice())
		return
	}
	inEdges := node.InEdgeList.Slice()
	switch len(inEdges) {
	case 0:
		return
	case 1:
		state.in = jutil.CheckNotNullPtr(d.nodeStates[inEdges[0].SourceNode.Index].out)
	default:
		d.joinOp.result = d.allocMustDefEmpty()
		for _, inEdge := range inEdges {
			d.joinOp.JoinFlow(d.nodeStates[inEdge.SourceNode.Index].out)
		}
		state.in = jutil.CheckNotNullPtr(d.joinOp.result)
	}
}

func (d *MustBeReachingVariableDef) resetVarStorage(n int) bool {
	releaseLarge := n <= 256
	d.escapedBits.clear(releaseLarge)
	d.paramBits.clear(releaseLarge)
	d.catchBits.clear(releaseLarge)
	d.hasParams = false
	d.hasCatchVars = false

	clear(d.vars)
	if n == 0 {
		if cap(d.vars) > 1024 {
			d.vars = nil
		} else {
			d.vars = d.vars[:0]
		}
		if len(d.idByNameID) > 1024 {
			d.idByNameID = nil
		} else if len(d.idByNameID) > 0 {
			clear(d.idByNameID)
		}
		return false
	}
	if cap(d.vars) < n || (cap(d.vars) > 1024 && releaseLarge) {
		d.vars = make([]*Var, 0, max(n, 16))
	} else {
		d.vars = d.vars[:0]
	}
	useMap := n > 16
	if useMap {
		if d.idByNameID == nil || (len(d.idByNameID) > 1024 && releaseLarge) {
			d.idByNameID = make(map[uint32]int, max(n, 32))
		} else {
			clear(d.idByNameID)
		}
	}
	return useMap
}

func (d *MustBeReachingVariableDef) appendInitVar(v *Var, isEscaped bool, useMap bool) {
	id := len(d.vars)
	d.vars = append(d.vars, v)
	if useMap {
		d.idByNameID[v.nameID] = id
	}
	if isEscaped {
		d.escapedBits.set(id)
	}
	if v.IsParam() {
		d.paramBits.set(id)
		d.hasParams = true
	}
	if sc, ok := v.GetScope().(*Scope); ok && sc != nil && sc.IsCatchScope() {
		d.catchBits.set(id)
		d.hasCatchVars = true
	}
}

func (d *MustBeReachingVariableDef) initVars(
	vars *jutil.LinkedHashMap[string, *Var],
	escaped *jutil.ObjectLinkedHashSet[*Var]) {
	n := 0
	if vars != nil {
		n = vars.Size()
	}
	useMap := d.resetVarStorage(n)
	if n == 0 {
		return
	}
	for v := range vars.Values() {
		d.appendInitVar(v, escaped != nil && escaped.Contains(v), useMap)
	}
}

func (d *MustBeReachingVariableDef) initVarsFromTable(table *fnVarTable) {
	n := 0
	if table != nil {
		n = table.uniqueCount
	}
	useMap := d.resetVarStorage(n)
	if n == 0 {
		return
	}
	if !table.hasDuplicateNames {
		for i, v := range table.vars {
			d.appendInitVar(v, table.escaped.has(i), useMap)
		}
		return
	}
	for _, v := range table.vars {
		if useMap {
			if _, dup := d.idByNameID[v.nameID]; dup {
				continue
			}
		} else {
			dup := false
			for _, prev := range d.vars {
				if prev.nameID == v.nameID {
					dup = true
					break
				}
			}
			if dup {
				continue
			}
		}
		canonIdx, _ := table.lookupIndexByID(v.nameID)
		d.appendInitVar(table.vars[canonIdx], table.escaped.has(canonIdx), useMap)
	}
}

func (d *MustBeReachingVariableDef) lookupVarIDByStringID(nameID uint32) (int, bool) {
	if nameID == 0 {
		return 0, false
	}
	vars := d.vars
	if len(vars) <= 16 {
		for i, v := range vars {
			if v.nameID == nameID {
				return i, true
			}
		}
		return 0, false
	}
	id, ok := d.idByNameID[nameID]
	return id, ok
}

func (d *MustBeReachingVariableDef) lookupVarID(name string) (int, bool) {
	return d.lookupVarIDByStringID(rhino.RhinoStringPool_LookupID(name))
}

// ResolveVarID returns the dense variable ID for v if v matches an entry in
// this function's variable table by ScopedName equality.
func (d *MustBeReachingVariableDef) ResolveVarID(v *Var) (int, bool) {
	if v == nil {
		return -1, false
	}
	id, ok := d.lookupVarIDByStringID(v.nameID)
	if !ok {
		return -1, false
	}
	canon := d.vars[id]
	if canon != v && !canon.Equals(v) {
		return -1, false
	}
	return id, true
}

func (d *MustBeReachingVariableDef) hasDependentNameID(def *MustBeReachingVariableDef_Definition, nameID uint32) bool {
	if def == nil || def.deps.IsEmpty() {
		return false
	}
	id, ok := d.lookupVarIDByStringID(nameID)
	return ok && def.deps.has(id)
}

func (d *MustBeReachingVariableDef) allocMustDefEmpty() *MustBeReachingVariableDef_MustDef {
	m := d.mustDefs.Alloc()
	if len(m.defs) > 0 {
		clear(m.defs)
		if cap(m.defs) > 512 && len(d.vars) <= 128 {
			m.defs = nil
		} else {
			m.defs = m.defs[:0]
		}
	}
	return m
}

func materializeMustDef(m *MustBeReachingVariableDef_MustDef, n int) []*MustBeReachingVariableDef_Definition {
	if len(m.defs) == n {
		return m.defs
	}
	if len(m.defs) > n {
		clear(m.defs[n:])
	}
	if cap(m.defs) >= n && (cap(m.defs) <= 512 || n > 128) {
		defs := m.defs[:n]
		clear(defs)
		m.defs = defs
		return defs
	}
	if n <= len(m.defsBuf) {
		defs := m.defsBuf[:n]
		clear(defs)
		m.defs = defs
		return defs
	}
	defs := make([]*MustBeReachingVariableDef_Definition, n, max(n, 32))
	m.defs = defs
	return defs
}

func cloneMustDefInto(src, dst *MustBeReachingVariableDef_MustDef) {
	n := len(src.defs)
	if n == 0 {
		if len(dst.defs) > 0 {
			clear(dst.defs)
			if cap(dst.defs) > 512 {
				dst.defs = nil
			} else {
				dst.defs = dst.defs[:0]
			}
		}
		return
	}
	if len(dst.defs) > n {
		clear(dst.defs[n:])
	}
	var dstDefs []*MustBeReachingVariableDef_Definition
	if cap(dst.defs) >= n && (cap(dst.defs) <= 512 || n > 128) {
		dstDefs = dst.defs[:n]
	} else if n <= len(dst.defsBuf) {
		dstDefs = dst.defsBuf[:n]
	} else {
		dstDefs = make([]*MustBeReachingVariableDef_Definition, n, max(n, 32))
	}
	copy(dstDefs, src.defs)
	dst.defs = dstDefs
}

func (d *MustBeReachingVariableDef) allocDefinition(node rhino.Node) *MustBeReachingVariableDef_Definition {
	def := d.defs.Alloc()
	def.node = node
	def.unknownDependencies = false
	def.deps.clear(len(d.vars) <= 256)
	return def
}

func (d *MustBeReachingVariableDef) ensureMutableDefs(def *MustBeReachingVariableDef_MustDef) []*MustBeReachingVariableDef_Definition {
	if def != nil && (def == d.flowInput || def == &d.emptyLattice) {
		if d.flowOutput == d.flowInput {
			cloned := d.mustDefs.Alloc()
			cloneMustDefInto(def, cloned)
			d.flowOutput = cloned
		}
		def = d.flowOutput
	}
	return materializeMustDef(def, len(d.vars))
}

// MustBeReachingVariableDef_Definition is an abstraction of a local variable
// definition. It represents the node which a local variable is defined as
// well as a set of other local variables that this definition reads from.
// For example N: a = b + foo.bar(c). The definition node will be N, the
// depending set would be {b,c}.
type MustBeReachingVariableDef_Definition struct {
	node                rhino.Node
	deps                inlineBitSet
	unknownDependencies bool
}

// Java: MustBeReachingVariableDef.Definition(Node)
func NewMustBeReachingVariableDef_Definition(node rhino.Node) *MustBeReachingVariableDef_Definition {
	return &MustBeReachingVariableDef_Definition{node: node}
}

// HasDependencies reports whether this definition reads any local variables.
func (d *MustBeReachingVariableDef_Definition) HasDependencies() bool {
	return !d.deps.IsEmpty()
}

// DependsOnVarID reports whether this definition reads local variable id.
func (d *MustBeReachingVariableDef_Definition) DependsOnVarID(id int) bool {
	return d.deps.has(id)
}

// ContainsAllDependencies reports whether this definition reads every local
// variable read by other.
func (d *MustBeReachingVariableDef_Definition) ContainsAllDependencies(other *MustBeReachingVariableDef_Definition) bool {
	return d.deps.ContainsAll(&other.deps)
}

// Java: MustBeReachingVariableDef.Definition.equals(Object)
func (d *MustBeReachingVariableDef_Definition) Equals(other any) bool {
	otherDef, ok := other.(*MustBeReachingVariableDef_Definition)
	if !ok || otherDef == nil {
		return false
	}
	if d == &mustDefBottom || otherDef == &mustDefBottom {
		return d == otherDef
	}
	// If the var has the same definition node we can assume they have the
	// same depends set.
	return otherDef.node == d.node
}

// Java: MustBeReachingVariableDef.Definition.toString()
func (d *MustBeReachingVariableDef_Definition) String() string {
	return "Definition@" + jutil.ToJavaString(d.node)
}

// Java: MustBeReachingVariableDef.Definition.hashCode()
func (d *MustBeReachingVariableDef_Definition) HashCode() int32 {
	if d == &mustDefBottom {
		return -1
	}
	return jutil.ObjectHashCode(d.node)
}

// MustBeReachingVariableDef_MustDef is the must reaching definition lattice
// representation. It captures a product lattice for each local
// (non-escaped) variable. The sub-lattice is a n + 2 element lattice with
// all the Definition in the program, TOP and BOTTOM.
//
// Since this is a Must-Define analysis, BOTTOM represents the case where
// there might be more than one reaching definition for the variable.
//
// (TOP) / | | \ N1 N2 N3 ....Nn \ | | / (BOTTOM)
type MustBeReachingVariableDef_MustDef struct {
	// Indexed by dense variable ID:
	//   nil           -> TOP (unreached / uninitialized)
	//   &mustDefBottom -> BOTTOM (ambiguous / conditional / escaped)
	//   *Definition   -> single reaching definition
	// When len(defs) == 0, every variable is implicitly TOP.
	defs    []*MustBeReachingVariableDef_Definition
	defsBuf [16]*MustBeReachingVariableDef_Definition
}

func mustDefEquals(a, b *MustBeReachingVariableDef_MustDef) bool {
	if a == b {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	aDefs, bDefs := a.defs, b.defs
	if len(aDefs) == 0 {
		for _, x := range bDefs {
			if x != nil {
				return false
			}
		}
		return true
	}
	if len(bDefs) == 0 {
		for _, x := range aDefs {
			if x != nil {
				return false
			}
		}
		return true
	}
	if len(aDefs) != len(bDefs) {
		return false
	}
	for i, da := range aDefs {
		db := bDefs[i]
		if da == db {
			continue
		}
		if da == nil || db == nil || da == &mustDefBottom || db == &mustDefBottom || da.node != db.node {
			return false
		}
	}
	return true
}

// Java: MustBeReachingVariableDef.MustDef.equals(Object)
func (m *MustBeReachingVariableDef_MustDef) Equals(other any) bool {
	mustDef, ok := other.(*MustBeReachingVariableDef_MustDef)
	return ok && mustDefEquals(m, mustDef)
}

// Java: MustBeReachingVariableDef.MustDef.hashCode()
func (m *MustBeReachingVariableDef_MustDef) HashCode() int32 {
	var h int32
	for i, def := range m.defs {
		if def != nil {
			h = 31*h + int32(i)*97 + def.HashCode()
		}
	}
	return h
}

// MustBeReachingVariableDef_MustDefJoin is the FlowJoiner of this analysis.
type MustBeReachingVariableDef_MustDefJoin struct {
	result *MustBeReachingVariableDef_MustDef
}

// Java: MustBeReachingVariableDef.MustDefJoin.joinFlow(MustDef)
func (j *MustBeReachingVariableDef_MustDefJoin) JoinFlow(input *MustBeReachingVariableDef_MustDef) {
	srcDefs := input.defs
	if len(srcDefs) == 0 {
		return
	}
	dst := j.result
	if len(dst.defs) == 0 {
		cloneMustDefInto(input, dst)
		return
	}
	dstDefs := dst.defs
	for i, s := range srcDefs {
		if s == nil {
			continue
		}
		cur := dstDefs[i]
		if cur == s {
			continue
		}
		if cur == nil {
			dstDefs[i] = s
		} else if cur == &mustDefBottom || s == &mustDefBottom || cur.node != s.node {
			dstDefs[i] = &mustDefBottom
		}
	}
}

// Java: MustBeReachingVariableDef.MustDefJoin.finish()
func (j *MustBeReachingVariableDef_MustDefJoin) Finish() *MustBeReachingVariableDef_MustDef {
	return j.result
}

// Java: MustBeReachingVariableDef.isForward()
func (d *MustBeReachingVariableDef) IsForward() bool {
	return true
}

// Java: MustBeReachingVariableDef.createEntryLattice()
func (d *MustBeReachingVariableDef) CreateEntryLattice() *MustBeReachingVariableDef_MustDef {
	if len(d.vars) == 0 {
		return &d.emptyLattice
	}
	m := d.allocMustDefEmpty()
	defs := materializeMustDef(m, len(d.vars))
	for id, var_ := range d.vars {
		defs[id] = d.allocDefinition(var_.GetScope().(*Scope).GetRootNode())
	}
	return m
}

// Java: MustBeReachingVariableDef.createInitialEstimateLattice()
func (d *MustBeReachingVariableDef) CreateInitialEstimateLattice() *MustBeReachingVariableDef_MustDef {
	return &d.emptyLattice
}

// Java: MustBeReachingVariableDef.flowThrough(Node, MustDef)
func (d *MustBeReachingVariableDef) FlowThrough(n rhino.Node, input *MustBeReachingVariableDef_MustDef) *MustBeReachingVariableDef_MustDef {
	d.flowInput = input
	d.flowOutput = input
	// TODO(user): This must know about ON_EX edges but it should handle
	// it better than what we did in liveness. Because we are in a forward mode,
	// we can used the branched forward analysis.
	d.ComputeMustDef(n, n, input, false)
	output := d.flowOutput
	d.flowInput = nil
	d.flowOutput = nil
	return output
}

// ComputeMustDef: n is the node in question; cfgNode the node to add;
// conditional true if the definition is not always executed.
//
// Java: MustBeReachingVariableDef.computeMustDef(Node, Node, MustDef, boolean)
func (d *MustBeReachingVariableDef) ComputeMustDef(n rhino.Node, cfgNode rhino.Node, output *MustBeReachingVariableDef_MustDef, conditional bool) {
	switch n.GetToken() {
	case rhino.Token_BLOCK, rhino.Token_ROOT, rhino.Token_FUNCTION:
		return
	case rhino.Token_WHILE, rhino.Token_DO, rhino.Token_IF, rhino.Token_FOR:
		d.ComputeMustDef(NodeUtil_GetConditionExpression(n), cfgNode, output, conditional)
		return
	case rhino.Token_FOR_IN, rhino.Token_FOR_OF, rhino.Token_FOR_AWAIT_OF:
		// for(x in y) {...}
		loop := ast.MustEnhancedFor(n)
		rhs := loop.Collection().Node()
		var lhs rhino.Node
		if decl, ok := loop.Target().AsNameDeclaration(); ok {
			if last, hasLast := decl.LastDeclarator(); hasLast {
				lhs = last.Node() // for(var x in y) {...}
			}
		} else {
			lhs = loop.Target().Node()
		}
		if name, ok := ast.AsName(lhs); ok {
			// We pass null for the def node because an enhanced for-loop conditionally assigns
			// elements from the collection at runtime (or doesn't execute if empty), so there is no
			// unconditional reaching definition from the loop header.
			d.addToDefIfLocalID(name.Node().GetStringID(), rhino.NilNode, rhs, output)
		} else if dlhs, ok := ast.AsDestructuringLhs(lhs); ok {
			lhs = dlhs.Pattern().Node()
		}
		if _, ok := ast.AsPattern(lhs); ok {
			d.ComputeMustDef(lhs, cfgNode, output, true)
		}
		return
	case rhino.Token_OPTCHAIN_GETPROP:
		d.ComputeMustDef(ast.MustOptChainGetProp(n).Object().Node(), cfgNode, output, conditional)
		return
	case rhino.Token_AND, rhino.Token_OR, rhino.Token_COALESCE:
		bin := ast.MustBinary(n)
		d.ComputeMustDef(bin.Left().Node(), cfgNode, output, conditional)
		d.ComputeMustDef(bin.Right().Node(), cfgNode, output, true)
		return
	case rhino.Token_OPTCHAIN_GETELEM:
		ge := ast.MustOptChainGetElem(n)
		d.ComputeMustDef(ge.Object().Node(), cfgNode, output, conditional)
		d.ComputeMustDef(ge.Index().Node(), cfgNode, output, true)
		return
	case rhino.Token_OPTCHAIN_CALL:
		call := ast.MustOptChainCall(n)
		d.ComputeMustDef(call.Callee().Node(), cfgNode, output, conditional)
		for arg, ok := call.FirstArg(); ok; arg, ok = arg.Next() {
			d.ComputeMustDef(arg.Node(), cfgNode, output, true)
		}
		return
	case rhino.Token_HOOK:
		hook := ast.MustHook(n)
		d.ComputeMustDef(hook.Cond().Node(), cfgNode, output, conditional)
		d.ComputeMustDef(hook.Then().Node(), cfgNode, output, true)
		d.ComputeMustDef(hook.Else().Node(), cfgNode, output, true)
		return
	case rhino.Token_LET, rhino.Token_CONST, rhino.Token_VAR:
		decl := ast.MustNameDeclaration(n)
		for declItem, ok := decl.FirstDeclarator(); ok; declItem, ok = declItem.Next() {
			if name, isName := declItem.AsName(); isName {
				if init, hasInit := name.Init(); hasInit {
					d.ComputeMustDef(init.Node(), cfgNode, output, conditional)
					var defNode rhino.Node
					if !conditional {
						defNode = cfgNode
					}
					d.addToDefIfLocalID(name.Node().GetStringID(), defNode, init.Node(), output)
				}
			} else {
				dlhs := declItem.MustDestructuringLhs()
				if init, hasInit := dlhs.Init(); hasInit {
					d.ComputeMustDef(init.Node(), cfgNode, output, conditional)
				}
				d.ComputeMustDef(dlhs.Pattern().Node(), cfgNode, output, conditional)
			}
		}
		return
	case rhino.Token_DEFAULT_VALUE:
		dv := ast.MustDefaultValue(n)
		target := dv.Target()
		value := dv.Value().Node()
		if pat, ok := target.AsPattern(); ok {
			d.ComputeMustDef(value, cfgNode, output, true)
			d.ComputeMustDef(pat.Node(), cfgNode, output, conditional)
		} else if name, ok := target.AsName(); ok {
			d.ComputeMustDef(value, cfgNode, output, true)
			var defNode rhino.Node
			if !conditional {
				defNode = cfgNode
			}
			d.addToDefIfLocalID(name.Node().GetStringID(), defNode, rhino.NilNode, output)
		} else {
			d.ComputeMustDef(target.Node(), cfgNode, output, conditional)
			d.ComputeMustDef(value, cfgNode, output, true)
		}
	case rhino.Token_NAME:
		strID := n.GetStringID()
		if NodeUtil_IsLhsByDestructuring(n) {
			var defNode rhino.Node
			if !conditional {
				defNode = cfgNode
			}
			d.addToDefIfLocalID(strID, defNode, rhino.NilNode, output)
		} else if strID == argumentsStringID {
			d.EscapeParameters(output)
		}
		return
	default:
		if assign, ok := ast.AsAssignment(n); ok {
			target := assign.Target()
			value := assign.Value().Node()
			if name, isName := target.AsName(); isName {
				d.ComputeMustDef(value, cfgNode, output, conditional)
				var defNode rhino.Node
				if !conditional {
					defNode = cfgNode
				}
				d.addToDefIfLocalID(name.Node().GetStringID(), defNode, value, output)
				return
			} else if get, isGet := target.AsNormalGet(); isGet {
				// Treat all assignments to arguments as redefining the
				// parameters itself.
				if obj, isObjName := get.Object().AsName(); isObjName && obj.Node().GetStringID() == argumentsStringID {
					// TODO(user): More accuracy can be introduced
					// i.e. We know exactly what arguments[x] is if x is a constant
					// number.
					d.EscapeParameters(output)
				}
			} else if pat, isPat := target.AsPattern(); isPat {
				d.ComputeMustDef(value, cfgNode, output, conditional)
				d.ComputeMustDef(pat.Node(), cfgNode, output, conditional)
				return
			}
		}

		// DEC and INC actually defines the variable.
		if incDec, ok := ast.AsIncDec(n); ok {
			if target, isName := incDec.Target().AsName(); isName {
				var defNode rhino.Node
				if !conditional {
					defNode = cfgNode
				}
				d.addToDefIfLocalID(target.Node().GetStringID(), defNode, rhino.NilNode, output)
				return
			}
		}

		for c := n.GetFirstChild(); !c.IsNil(); c = c.GetNext() {
			d.ComputeMustDef(c, cfgNode, output, conditional)
		}
	}
}

func hasDependent(defs []*MustBeReachingVariableDef_Definition, id int) bool {
	if len(defs) == 0 {
		return false
	}
	if id < 64 {
		bit := uint64(1) << id
		for _, otherDef := range defs {
			if otherDef != nil && otherDef.deps.w0&bit != 0 {
				return true
			}
		}
		return false
	}
	if id < 128 {
		bit := uint64(1) << (id - 64)
		for _, otherDef := range defs {
			if otherDef != nil && otherDef.deps.w1&bit != 0 {
				return true
			}
		}
		return false
	}
	wi := (id - 128) >> 6
	bit := uint64(1) << ((id - 128) & 63)
	for _, otherDef := range defs {
		if otherDef != nil && wi < len(otherDef.deps.ext) && otherDef.deps.ext[wi]&bit != 0 {
			return true
		}
	}
	return false
}

func invalidateDependents(defs []*MustBeReachingVariableDef_Definition, id int) {
	if id < 64 {
		bit := uint64(1) << id
		for i, otherDef := range defs {
			if otherDef != nil && otherDef.deps.w0&bit != 0 {
				defs[i] = &mustDefBottom
			}
		}
		return
	}
	if id < 128 {
		bit := uint64(1) << (id - 64)
		for i, otherDef := range defs {
			if otherDef != nil && otherDef.deps.w1&bit != 0 {
				defs[i] = &mustDefBottom
			}
		}
		return
	}
	wi := (id - 128) >> 6
	bit := uint64(1) << ((id - 128) & 63)
	for i, otherDef := range defs {
		if otherDef != nil && wi < len(otherDef.deps.ext) && otherDef.deps.ext[wi]&bit != 0 {
			defs[i] = &mustDefBottom
		}
	}
}

func (d *MustBeReachingVariableDef) addToDefIfLocalID(
	nameID uint32, node rhino.Node, rValue rhino.Node, def *MustBeReachingVariableDef_MustDef) {
	id, ok := d.lookupVarIDByStringID(nameID)
	if !ok {
		return
	}
	activeDefs := def.defs
	if def == d.flowInput && d.flowOutput != d.flowInput {
		activeDefs = d.flowOutput.defs
	}
	if d.escapedBits.has(id) {
		if !hasDependent(activeDefs, id) {
			return
		}
		defs := d.ensureMutableDefs(def)
		invalidateDependents(defs, id)
		return
	}
	defs := d.ensureMutableDefs(def)
	invalidateDependents(defs, id)
	if node.IsNil() {
		defs[id] = &mustDefBottom
		return
	}
	definition := d.allocDefinition(node)
	if !rValue.IsNil() {
		d.ComputeDependence(definition, rValue)
	}
	defs[id] = definition
}

// AddToDefIfLocal sets the variable lattice for the given name to the node
// value in the def lattice. Do nothing if the variable name is one of the
// escaped variable.
//
// node is the CFG node where the definition should be record to; nil if
// this is a conditional define.
//
// Java: MustBeReachingVariableDef.addToDefIfLocal(String, Node, Node, MustDef)
func (d *MustBeReachingVariableDef) AddToDefIfLocal(
	name string, node rhino.Node, rValue rhino.Node, def *MustBeReachingVariableDef_MustDef) {
	d.addToDefIfLocalID(rhino.RhinoStringPool_LookupID(name), node, rValue, def)
}

// Java: MustBeReachingVariableDef.escapeParameters(MustDef)
func (d *MustBeReachingVariableDef) EscapeParameters(output *MustBeReachingVariableDef_MustDef) {
	if !d.hasParams {
		return
	}
	defs := d.ensureMutableDefs(output)
	d.paramBits.ForEach(func(id int) {
		defs[id] = &mustDefBottom
	})
	for i, val := range defs {
		if val != nil && val.deps.Intersects(&d.paramBits) {
			defs[i] = &mustDefBottom
		}
	}
}

func (d *MustBeReachingVariableDef) recordDependenceID(def *MustBeReachingVariableDef_Definition, nameID uint32) {
	id, ok := d.lookupVarIDByStringID(nameID)
	if !ok {
		def.unknownDependencies = true
		return
	}
	def.deps.set(id)
}

func (d *MustBeReachingVariableDef) computeDependenceWalk(def *MustBeReachingVariableDef_Definition, n, parent rhino.Node) {
	if !parent.IsNil() && ControlFlowGraph_IsEnteringNewCfgNode(n) {
		return
	}
	if n.GetToken() == rhino.Token_CLASS_MEMBERS {
		for child := n.GetFirstChild(); !child.IsNil(); child = child.GetNext() {
			if cm, ok := ast.AsComputedMember(child); ok {
				d.computeDependenceWalk(def, cm.Key().Node(), child)
			}
		}
		for child := n.GetFirstChild(); !child.IsNil(); child = child.GetNext() {
			if cm, ok := ast.AsComputedMember(child); ok {
				if !ControlFlowGraph_IsEnteringNewCfgNode(child) {
					if val, hasVal := cm.Value(); hasVal {
						d.computeDependenceWalk(def, val.Node(), child)
					}
				}
			} else {
				d.computeDependenceWalk(def, child, n)
			}
		}
		return
	}
	for c := n.GetFirstChild(); !c.IsNil(); c = c.GetNext() {
		d.computeDependenceWalk(def, c, n)
	}
	if n.IsName() {
		d.recordDependenceID(def, n.GetStringID())
	}
}

// ComputeDependence computes all the local variables that rValue reads from
// and store that in the def's depends set.
//
// Java: MustBeReachingVariableDef.computeDependence(Definition, Node)
func (d *MustBeReachingVariableDef) ComputeDependence(def *MustBeReachingVariableDef_Definition, rValue rhino.Node) {
	jutil.CheckNotNull(d.compiler)
	if rValue.IsNil() {
		throwNilCfgNodeTraversal(d.compiler)
		return
	}
	d.computeDependenceWalk(def, rValue, rhino.NilNode)
}

func (d *MustBeReachingVariableDef) getDefFromMustDefByID(nameID uint32, inState *MustBeReachingVariableDef_MustDef) *MustBeReachingVariableDef_Definition {
	if inState == nil || len(inState.defs) == 0 {
		return nil
	}
	id, ok := d.lookupVarIDByStringID(nameID)
	if !ok {
		return nil
	}
	def := inState.defs[id]
	if def == &mustDefBottom {
		return nil
	}
	return def
}

func (d *MustBeReachingVariableDef) getDefFromMustDef(name string, inState *MustBeReachingVariableDef_MustDef) *MustBeReachingVariableDef_Definition {
	return d.getDefFromMustDefByID(rhino.RhinoStringPool_LookupID(name), inState)
}

func (d *MustBeReachingVariableDef) getDefByID(nameID uint32, useNode rhino.Node) *MustBeReachingVariableDef_Definition {
	n := d.cfg.Nodes.Get(useNode)
	jutil.CheckArgument(n != nil)
	return d.getDefFromMustDefByID(nameID, d.nodeStates[n.Index].in)
}

// GetDef gets the must-be-reaching definition of a given use node.
//
// name is the name of the variable. It can only be names of local variable
// that are not function parameters, escaped variables or variables declared
// in catch. useNode is the location of the use where the definition
// reaches.
//
// Java: MustBeReachingVariableDef.getDef(String, Node)
func (d *MustBeReachingVariableDef) GetDef(name string, useNode rhino.Node) *MustBeReachingVariableDef_Definition {
	return d.getDefByID(rhino.RhinoStringPool_LookupID(name), useNode)
}

// Java: MustBeReachingVariableDef.getDefNode(String, Node)
func (d *MustBeReachingVariableDef) GetDefNode(name string, useNode rhino.Node) rhino.Node {
	def := d.GetDef(name, useNode)
	if def == nil {
		return rhino.NilNode
	}
	return def.node
}

// Java: MustBeReachingVariableDef.dependsOnOuterScopeVars(Definition)
func (d *MustBeReachingVariableDef) DependsOnOuterScopeVars(def *MustBeReachingVariableDef_Definition) bool {
	if def.unknownDependencies {
		return true
	}
	if !d.hasCatchVars {
		return false
	}
	return def.deps.Intersects(&d.catchBits)
}
