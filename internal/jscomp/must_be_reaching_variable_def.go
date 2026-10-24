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
	"iter"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

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
	DataFlowAnalysisBase[*rhino.Node, *MustBeReachingVariableDef_MustDef]

	// The scope of the function that we are analyzing.
	compiler AbstractCompiler
	// Go-only: language-adapter — Var overrides equals/hashCode (ScopedName),
	// hence the Object set.
	escaped     *jutil.ObjectLinkedHashSet[*Var]
	allVarsInFn *jutil.LinkedHashMap[string, *Var]
	// Go-only: performance — local variable IDs shared by lattice states.
	layout        *mustDefLayout
	layoutStorage mustDefLayout
	idByName      map[string]uint16
	escapedMask   uint64
	paramMask     uint64
	catchMask     uint64
	escapedWords  [maxMustDefLayoutWords]uint64
	paramWords    [maxMustDefLayoutWords]uint64
	catchWords    [maxMustDefLayoutWords]uint64
	hasParams     bool
	hasCatchVars  bool

	emptyLattice MustBeReachingVariableDef_MustDef
	joinOp       MustBeReachingVariableDef_MustDefJoin
	depCallback  mustBeReachingVariableDef_dependenceCallback

	flowInput  *MustBeReachingVariableDef_MustDef
	flowOutput *MustBeReachingVariableDef_MustDef

	// Go-only: performance — rewound for each CFG: allocMustDef(Raw) and
	// allocDefinition re-initialize what they hand out and reuse its storage.
	mustDefs reuseSlab[MustBeReachingVariableDef_MustDef]
	defs     reuseSlab[MustBeReachingVariableDef_Definition]
}

func (d *MustBeReachingVariableDef) initLayout(vars *jutil.LinkedHashMap[string, *Var]) {
	d.escapedMask = 0
	d.paramMask = 0
	d.catchMask = 0
	d.hasParams = false
	d.hasCatchVars = false
	if vars == nil || vars.Size() > maxMustDefLayoutVars {
		d.layout = nil
		return
	}
	n := vars.Size()
	if n == 0 {
		d.layout = emptyMustDefLayout
		return
	}
	l := &d.layoutStorage
	clear(l.vars)
	if cap(l.vars) < n {
		l.vars = make([]*Var, 0, max(n, 16))
	} else {
		l.vars = l.vars[:0]
	}
	if l.ids == nil {
		l.ids = make(map[*Var]uint16, max(n, 16))
	} else {
		clear(l.ids)
	}
	l.numWords = (n + 63) >> 6
	l.wide = n > 64
	useMap := n > 16
	if useMap {
		if d.idByName == nil {
			d.idByName = make(map[string]uint16, max(n, 32))
		} else {
			clear(d.idByName)
		}
	}
	if !l.wide {
		var escapedMask, paramMask, catchMask uint64
		for v := range vars.Values() {
			id := uint16(len(l.vars))
			bit := uint64(1) << id
			l.ids[v] = id
			l.vars = append(l.vars, v)
			if useMap {
				d.idByName[v.GetName()] = id
			}
			if d.escaped != nil && d.escaped.Contains(v) {
				escapedMask |= bit
			}
			if v.IsParam() {
				paramMask |= bit
			}
			if sc, ok := v.GetScope().(*Scope); ok && sc != nil && sc.IsCatchScope() {
				catchMask |= bit
			}
		}
		d.escapedMask = escapedMask
		d.paramMask = paramMask
		d.catchMask = catchMask
		d.hasParams = paramMask != 0
		d.hasCatchVars = catchMask != 0
	} else {
		clear(d.escapedWords[:l.numWords])
		clear(d.paramWords[:l.numWords])
		clear(d.catchWords[:l.numWords])
		for v := range vars.Values() {
			id := uint16(len(l.vars))
			w := id >> 6
			bit := uint64(1) << (id & 63)
			l.ids[v] = id
			l.vars = append(l.vars, v)
			if useMap {
				d.idByName[v.GetName()] = id
			}
			if d.escaped != nil && d.escaped.Contains(v) {
				d.escapedWords[w] |= bit
			}
			if v.IsParam() {
				d.paramWords[w] |= bit
				d.hasParams = true
			}
			if sc, ok := v.GetScope().(*Scope); ok && sc != nil && sc.IsCatchScope() {
				d.catchWords[w] |= bit
				d.hasCatchVars = true
			}
		}
	}
	d.layout = l
}

func (d *MustBeReachingVariableDef) lookupLocalVar(name string) (uint16, *Var, bool) {
	vars := d.layout.vars
	if len(vars) <= 16 {
		id := rhino.RhinoStringPool_LookupID(name)
		if id == 0 {
			return 0, nil, false
		}
		for i, v := range vars {
			if v.nameID == id {
				return uint16(i), v, true
			}
		}
		return 0, nil, false
	}
	id, ok := d.idByName[name]
	if !ok {
		return 0, nil, false
	}
	return id, vars[id], true
}

func resetMustDefMap(m *mustDefMap, layout *mustDefLayout) {
	m.layout = layout
	m.present = 0
	if layout != nil {
		if layout.wide {
			clear(m.presentWords[:layout.numWords])
		}
		n := len(layout.vars)
		if n > 0 {
			if cap(m.values) >= n && cap(m.order) >= n {
				clear(m.values[:n])
				m.values = m.values[:n]
				m.order = m.order[:0]
			} else if n <= len(m.valuesBuf) {
				clear(m.valuesBuf[:n])
				m.values = m.valuesBuf[:n]
				m.order = m.orderBuf[:0:len(m.orderBuf)]
			} else {
				c := max(n, 64)
				m.values = make([]*MustBeReachingVariableDef_Definition, n, c)
				m.order = make([]uint16, 0, c)
			}
		} else {
			if len(m.values) > 0 {
				clear(m.values)
				m.values = m.values[:0]
			}
			m.order = m.order[:0]
		}
		return
	}
	if len(m.values) > 0 {
		clear(m.values)
		m.values = m.values[:0]
	}
	m.order = m.order[:0]
	if m.fallback == nil {
		m.fallback = jutil.NewObjectLinkedHashMap[*Var, *MustBeReachingVariableDef_Definition]()
	} else {
		m.fallback.Clear()
	}
}

func cloneMustDefMapInto(src, dst *mustDefMap) {
	if src.layout != nil {
		dst.layout = src.layout
		if !src.layout.wide {
			dst.present = src.present
		} else {
			copy(dst.presentWords[:src.layout.numWords], src.presentWords[:src.layout.numWords])
		}
		if len(src.values) == 0 {
			if len(dst.values) > 0 {
				clear(dst.values)
				dst.values = dst.values[:0]
			}
			dst.order = dst.order[:0]
			return
		}
		if cap(dst.values) < len(src.values) {
			if len(src.values) <= len(dst.valuesBuf) {
				dst.values = dst.valuesBuf[:len(src.values)]
			} else {
				dst.values = make([]*MustBeReachingVariableDef_Definition, len(src.values), max(len(src.values), 64))
			}
		} else {
			dst.values = dst.values[:len(src.values)]
		}
		copy(dst.values, src.values)
		if cap(dst.order) < len(src.values) {
			if len(src.values) <= len(dst.orderBuf) {
				dst.order = dst.orderBuf[:len(src.order):len(dst.orderBuf)]
			} else {
				dst.order = make([]uint16, len(src.order), max(len(src.values), 64))
			}
		} else {
			dst.order = dst.order[:len(src.order)]
		}
		copy(dst.order, src.order)
		return
	}
	dst.layout = nil
	dst.present = 0
	if len(dst.values) > 0 {
		clear(dst.values)
		dst.values = dst.values[:0]
	}
	dst.order = dst.order[:0]
	if dst.fallback == nil {
		dst.fallback = jutil.NewObjectLinkedHashMap[*Var, *MustBeReachingVariableDef_Definition]()
	} else {
		dst.fallback.Clear()
	}
	for v, d := range src.fallback.All() {
		dst.fallback.Put(v, d)
	}
}

// Java: MustBeReachingVariableDef(ControlFlowGraph<Node>, AbstractCompiler, Set<Var>, Map<String, Var>)
func NewMustBeReachingVariableDef(
	cfg ControlFlowGraph[*rhino.Node],
	compiler AbstractCompiler,
	escaped *jutil.ObjectLinkedHashSet[*Var],
	allVarsInFn *jutil.LinkedHashMap[string, *Var]) *MustBeReachingVariableDef {
	d := &MustBeReachingVariableDef{}
	d.Reset(cfg, compiler, escaped, allVarsInFn)
	return d
}

// Reset resets d for reuse on a new function scope while retaining slab and
// layout capacities.
func (d *MustBeReachingVariableDef) Reset(
	cfg ControlFlowGraph[*rhino.Node],
	compiler AbstractCompiler,
	escaped *jutil.ObjectLinkedHashSet[*Var],
	allVarsInFn *jutil.LinkedHashMap[string, *Var]) {
	d.InitDataFlowAnalysisBase(d, cfg)
	d.compiler = compiler
	d.escaped = escaped
	d.allVarsInFn = allVarsInFn
	d.initLayout(allVarsInFn)
	resetMustDefMap(&d.emptyLattice.reachingDef, d.layout)
	d.depCallback.outer = d
	d.depCallback.def = nil
	d.flowInput = nil
	d.flowOutput = nil
	d.mustDefs.rewind()
	d.defs.rewind()
}

func (d *MustBeReachingVariableDef) allocMustDefRaw() *MustBeReachingVariableDef_MustDef {
	return d.mustDefs.alloc()
}

func (d *MustBeReachingVariableDef) allocMustDef() *MustBeReachingVariableDef_MustDef {
	m := d.allocMustDefRaw()
	resetMustDefMap(&m.reachingDef, d.layout)
	return m
}

func (d *MustBeReachingVariableDef) allocDefinition(node *rhino.Node) *MustBeReachingVariableDef_Definition {
	def := d.defs.alloc()
	def.node = node
	def.unknownDependencies = false
	def.depMask = 0
	if d.layout != nil && d.layout.wide {
		clear(def.depWords[:d.layout.numWords])
	}
	def.dependsSet.Clear()
	def.depends = &def.dependsSet
	return def
}

func (d *MustBeReachingVariableDef) ensureMutable(def *MustBeReachingVariableDef_MustDef) *MustBeReachingVariableDef_MustDef {
	if def != nil && def == d.flowInput {
		if d.flowOutput == d.flowInput {
			cloned := d.allocMustDefRaw()
			cloneMustDefMapInto(&d.flowInput.reachingDef, &cloned.reachingDef)
			d.flowOutput = cloned
		}
		return d.flowOutput
	}
	return def
}

// MustBeReachingVariableDef_Definition is an abstraction of a local variable
// definition. It represents the node which a local variable is defined as
// well as a set of other local variables that this definition reads from.
// For example N: a = b + foo.bar(c). The definition node will be N, the
// depending set would be {b,c}.
type MustBeReachingVariableDef_Definition struct {
	node *rhino.Node
	// Go-only: language-adapter — Var overrides equals/hashCode (ScopedName),
	// hence the Object set.
	depends             *jutil.ObjectLinkedHashSet[*Var]
	dependsSet          jutil.ObjectLinkedHashSet[*Var]
	depMask             uint64
	depWords            [maxMustDefLayoutWords]uint64
	unknownDependencies bool
}

// Java: MustBeReachingVariableDef.Definition(Node)
func NewMustBeReachingVariableDef_Definition(node *rhino.Node) *MustBeReachingVariableDef_Definition {
	d := &MustBeReachingVariableDef_Definition{
		node:                node,
		unknownDependencies: false,
	}
	d.depends = &d.dependsSet
	return d
}

// Java: MustBeReachingVariableDef.Definition.equals(Object)
func (d *MustBeReachingVariableDef_Definition) Equals(other any) bool {
	otherDef, ok := other.(*MustBeReachingVariableDef_Definition)
	if !ok || otherDef == nil {
		return false
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
	// TODO(user): Use bit vector instead of maps might get better
	// performance. Change it after this is tested to be fully functional.

	// When a Var "A" = "TOP", "A" does not exist in reachingDef's keySet.
	// When a Var "A" = Node N, "A" maps to that node.
	// When a Var "A" = "BOTTOM", "A" maps to null.
	//
	// Go-only: language-adapter — Var overrides equals/hashCode (ScopedName),
	// hence the Object map.
	reachingDef mustDefMap
}

// Java: MustBeReachingVariableDef.MustDef()
func NewMustBeReachingVariableDef_MustDef() *MustBeReachingVariableDef_MustDef {
	return &MustBeReachingVariableDef_MustDef{
		reachingDef: newMustDefMap(nil),
	}
}

// Java: MustBeReachingVariableDef.MustDef(Collection<Var>)
func NewMustBeReachingVariableDef_MustDefWithVars(vars iter.Seq[*Var]) *MustBeReachingVariableDef_MustDef {
	m := NewMustBeReachingVariableDef_MustDef()
	for var_ := range vars {
		m.reachingDef.Put(var_, NewMustBeReachingVariableDef_Definition(var_.GetScope().(*Scope).GetRootNode()))
	}
	return m
}

// NewMustBeReachingVariableDef_MustDefCopy is the copy constructor.
//
// other: the constructed object is a replicated copy of this element.
//
// Java: MustBeReachingVariableDef.MustDef(MustDef)
func NewMustBeReachingVariableDef_MustDefCopy(other *MustBeReachingVariableDef_MustDef) *MustBeReachingVariableDef_MustDef {
	m := &MustBeReachingVariableDef_MustDef{}
	other.reachingDef.cloneInto(&m.reachingDef)
	return m
}

// Java: MustBeReachingVariableDef.MustDef.equals(Object)
func (m *MustBeReachingVariableDef_MustDef) Equals(other any) bool {
	mustDef, ok := other.(*MustBeReachingVariableDef_MustDef)
	return ok && mustDef != nil && (m == mustDef || mustDef.reachingDef.equals(&m.reachingDef))
}

// Java: MustBeReachingVariableDef.MustDef.hashCode()
func (m *MustBeReachingVariableDef_MustDef) HashCode() int32 {
	return m.reachingDef.HashCode()
}

// MustBeReachingVariableDef_MustDefJoin is the FlowJoiner of this analysis.
type MustBeReachingVariableDef_MustDefJoin struct {
	result    *MustBeReachingVariableDef_MustDef
	resultMap *mustDefMap
}

// Java: MustBeReachingVariableDef.MustDefJoin()
func NewMustBeReachingVariableDef_MustDefJoin() *MustBeReachingVariableDef_MustDefJoin {
	j := &MustBeReachingVariableDef_MustDefJoin{result: NewMustBeReachingVariableDef_MustDef()}
	j.resultMap = &j.result.reachingDef
	return j
}

// Java: MustBeReachingVariableDef.MustDefJoin.joinFlow(MustDef)
func (j *MustBeReachingVariableDef_MustDefJoin) JoinFlow(input *MustBeReachingVariableDef_MustDef) {
	// Go-only: performance — shared IDs avoid hashing every key in a join.
	dst, src := j.resultMap, &input.reachingDef
	if dst.layout != nil && dst.layout == src.layout {
		if !dst.layout.wide {
			for _, id := range src.order {
				def := src.values[id]
				if def != nil && dst.present&(uint64(1)<<id) != 0 {
					previous := dst.values[id]
					if previous != nil && previous.node == def.node {
						continue
					}
					def = nil
				}
				dst.putID(id, def)
			}
		} else {
			for _, id := range src.order {
				def := src.values[id]
				if def != nil && dst.presentWords[id>>6]&(uint64(1)<<(id&63)) != 0 {
					previous := dst.values[id]
					if previous != nil && previous.node == def.node {
						continue
					}
					def = nil
				}
				dst.putID(id, def)
			}
		}
		return
	}
	for var_, def := range input.reachingDef.All() {
		j.MergeVarDef(var_, def)
	}
}

// Java: MustBeReachingVariableDef.MustDefJoin.mergeVarDef(Var, Definition)
func (j *MustBeReachingVariableDef_MustDefJoin) MergeVarDef(var_ *Var, def *MustBeReachingVariableDef_Definition) {
	var resultDef *MustBeReachingVariableDef_Definition
	if def == nil {
		resultDef = nil
	} else if !j.resultMap.ContainsKey(var_) {
		resultDef = def
	} else if def.Equals(j.resultMap.Get(var_)) {
		return
	} else {
		resultDef = nil
	}
	j.resultMap.Put(var_, resultDef)
}

// Java: MustBeReachingVariableDef.MustDefJoin.finish()
func (j *MustBeReachingVariableDef_MustDefJoin) Finish() *MustBeReachingVariableDef_MustDef {
	return j.result
}

// Java: MustBeReachingVariableDef.isForward()
func (d *MustBeReachingVariableDef) IsForward() bool {
	return true
}

func (d *MustBeReachingVariableDef) EqualsLattice(a, b *MustBeReachingVariableDef_MustDef) bool {
	if a == b {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	return b.reachingDef.equals(&a.reachingDef)
}

// Java: MustBeReachingVariableDef.createEntryLattice()
func (d *MustBeReachingVariableDef) CreateEntryLattice() *MustBeReachingVariableDef_MustDef {
	if d.allVarsInFn.IsEmpty() && d.layout != nil {
		return &d.emptyLattice
	}
	m := d.allocMustDef()
	if d.layout != nil {
		for id, var_ := range d.layout.vars {
			m.reachingDef.putID(uint16(id), d.allocDefinition(var_.GetScope().(*Scope).GetRootNode()))
		}
		return m
	}
	for var_ := range d.allVarsInFn.Values() {
		m.reachingDef.Put(var_, d.allocDefinition(var_.GetScope().(*Scope).GetRootNode()))
	}
	return m
}

// Java: MustBeReachingVariableDef.createInitialEstimateLattice()
func (d *MustBeReachingVariableDef) CreateInitialEstimateLattice() *MustBeReachingVariableDef_MustDef {
	if d.layout != nil {
		return &d.emptyLattice
	}
	return d.allocMustDef()
}

// Java: MustBeReachingVariableDef.createFlowJoiner()
func (d *MustBeReachingVariableDef) CreateFlowJoiner() DataFlowAnalysis_FlowJoiner[*MustBeReachingVariableDef_MustDef] {
	m := d.allocMustDef()
	d.joinOp.result = m
	d.joinOp.resultMap = &m.reachingDef
	return &d.joinOp
}

// Java: MustBeReachingVariableDef.flowThrough(Node, MustDef)
func (d *MustBeReachingVariableDef) FlowThrough(n *rhino.Node, input *MustBeReachingVariableDef_MustDef) *MustBeReachingVariableDef_MustDef {
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
func (d *MustBeReachingVariableDef) ComputeMustDef(n *rhino.Node, cfgNode *rhino.Node, output *MustBeReachingVariableDef_MustDef, conditional bool) {
	switch n.GetToken() {
	case rhino.Token_BLOCK, rhino.Token_ROOT, rhino.Token_FUNCTION:
		return
	case rhino.Token_WHILE, rhino.Token_DO, rhino.Token_IF, rhino.Token_FOR:
		d.ComputeMustDef(NodeUtil_GetConditionExpression(n), cfgNode, output, conditional)
		return
	case rhino.Token_FOR_IN, rhino.Token_FOR_OF, rhino.Token_FOR_AWAIT_OF:
		// for(x in y) {...}
		lhs := n.GetFirstChild()
		rhs := lhs.GetNext()
		if NodeUtil_IsNameDeclaration(lhs) {
			lhs = lhs.GetLastChild() // for(var x in y) {...}
		}
		if lhs.IsName() {
			// We pass null for the def node because an enhanced for-loop conditionally assigns
			// elements from the collection at runtime (or doesn't execute if empty), so there is no
			// unconditional reaching definition from the loop header.
			d.AddToDefIfLocal(lhs.GetString(), nil, rhs, output)
		} else if lhs.IsDestructuringLhs() {
			lhs = lhs.GetFirstChild()
		}
		if lhs.IsDestructuringPattern() {
			d.ComputeMustDef(lhs, cfgNode, output, true)
		}
		return
	case rhino.Token_OPTCHAIN_GETPROP:
		d.ComputeMustDef(n.GetFirstChild(), cfgNode, output, conditional)
		return
	case rhino.Token_AND, rhino.Token_OR, rhino.Token_COALESCE, rhino.Token_OPTCHAIN_GETELEM:
		d.ComputeMustDef(n.GetFirstChild(), cfgNode, output, conditional)
		d.ComputeMustDef(n.GetLastChild(), cfgNode, output, true)
		return
	case rhino.Token_OPTCHAIN_CALL:
		d.ComputeMustDef(n.GetFirstChild(), cfgNode, output, conditional)
		for c := n.GetSecondChild(); c != nil; c = c.GetNext() {
			d.ComputeMustDef(c, cfgNode, output, true)
		}
		return
	case rhino.Token_HOOK:
		d.ComputeMustDef(n.GetFirstChild(), cfgNode, output, conditional)
		d.ComputeMustDef(n.GetSecondChild(), cfgNode, output, true)
		d.ComputeMustDef(n.GetLastChild(), cfgNode, output, true)
		return
	case rhino.Token_LET, rhino.Token_CONST, rhino.Token_VAR:
		for c := n.GetFirstChild(); c != nil; c = c.GetNext() {
			if c.HasChildren() {
				if c.IsName() {
					d.ComputeMustDef(c.GetFirstChild(), cfgNode, output, conditional)
					var defNode *rhino.Node
					if !conditional {
						defNode = cfgNode
					}
					d.AddToDefIfLocal(c.GetString(), defNode, c.GetFirstChild(), output)
				} else {
					jutil.CheckStateMsg1(c.IsDestructuringLhs(), "%s", c)
					d.ComputeMustDef(c.GetSecondChild(), cfgNode, output, conditional)
					d.ComputeMustDef(c.GetFirstChild(), cfgNode, output, conditional)
				}
			}
		}
		return
	case rhino.Token_DEFAULT_VALUE:
		if n.GetFirstChild().IsDestructuringPattern() {
			d.ComputeMustDef(n.GetSecondChild(), cfgNode, output, true)
			d.ComputeMustDef(n.GetFirstChild(), cfgNode, output, conditional)
		} else if n.GetFirstChild().IsName() {
			d.ComputeMustDef(n.GetSecondChild(), cfgNode, output, true)
			var defNode *rhino.Node
			if !conditional {
				defNode = cfgNode
			}
			d.AddToDefIfLocal(n.GetFirstChild().GetString(), defNode, nil, output)
		} else {
			d.ComputeMustDef(n.GetFirstChild(), cfgNode, output, conditional)
			d.ComputeMustDef(n.GetSecondChild(), cfgNode, output, true)
		}
	case rhino.Token_NAME:
		if NodeUtil_IsLhsByDestructuring(n) {
			var defNode *rhino.Node
			if !conditional {
				defNode = cfgNode
			}
			d.AddToDefIfLocal(n.GetString(), defNode, nil, output)
		} else if "arguments" == n.GetString() {
			d.EscapeParameters(output)
		}
		return
	default:
		if NodeUtil_IsAssignmentOp(n) {
			if n.GetFirstChild().IsName() {
				name := n.GetFirstChild()
				d.ComputeMustDef(name.GetNext(), cfgNode, output, conditional)
				var defNode *rhino.Node
				if !conditional {
					defNode = cfgNode
				}
				d.AddToDefIfLocal(name.GetString(), defNode, n.GetLastChild(), output)
				return
			} else if NodeUtil_IsNormalGet(n.GetFirstChild()) {
				// Treat all assignments to arguments as redefining the
				// parameters itself.
				obj := n.GetFirstFirstChild()
				if obj.IsName() && "arguments" == obj.GetString() {
					// TODO(user): More accuracy can be introduced
					// i.e. We know exactly what arguments[x] is if x is a constant
					// number.
					d.EscapeParameters(output)
				}
			} else if n.GetFirstChild().IsDestructuringPattern() {
				d.ComputeMustDef(n.GetSecondChild(), cfgNode, output, conditional)
				d.ComputeMustDef(n.GetFirstChild(), cfgNode, output, conditional)
				return
			}
		}

		// DEC and INC actually defines the variable.
		if n.IsDec() || n.IsInc() {
			target := n.GetFirstChild()
			if target.IsName() {
				var defNode *rhino.Node
				if !conditional {
					defNode = cfgNode
				}
				d.AddToDefIfLocal(target.GetString(), defNode, nil, output)
				return
			}
		}

		for c := n.GetFirstChild(); c != nil; c = c.GetNext() {
			d.ComputeMustDef(c, cfgNode, output, conditional)
		}
	}
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
	name string, node *rhino.Node, rValue *rhino.Node, def *MustBeReachingVariableDef_MustDef) {
	if d.layout != nil && def.reachingDef.layout == d.layout {
		id, _, ok := d.lookupLocalVar(name)
		if !ok {
			return
		}
		def = d.ensureMutable(def)
		m := &def.reachingDef
		if !d.layout.wide {
			bit := uint64(1) << id
			for _, otherID := range m.order {
				if otherDef := m.values[otherID]; otherDef != nil && otherDef.depMask&bit != 0 {
					m.values[otherID] = nil
				}
			}
			if d.escapedMask&bit == 0 {
				if node == nil {
					m.putID(id, nil)
				} else {
					definition := d.allocDefinition(node)
					if rValue != nil {
						d.ComputeDependence(definition, rValue)
					}
					m.putID(id, definition)
				}
			}
		} else {
			w := id >> 6
			bit := uint64(1) << (id & 63)
			for _, otherID := range m.order {
				if otherDef := m.values[otherID]; otherDef != nil && otherDef.depWords[w]&bit != 0 {
					m.values[otherID] = nil
				}
			}
			if d.escapedWords[w]&bit == 0 {
				if node == nil {
					m.putID(id, nil)
				} else {
					definition := d.allocDefinition(node)
					if rValue != nil {
						d.ComputeDependence(definition, rValue)
					}
					m.putID(id, definition)
				}
			}
		}
		return
	}

	var_ := d.allVarsInFn.Get(name)

	// var might be null if the variable is defined in the externs
	if var_ == nil {
		return
	}

	def = d.ensureMutable(def)
	def.reachingDef.invalidateDependent(var_)

	if !d.escaped.Contains(var_) {
		if node == nil {
			def.reachingDef.Put(var_, nil)
		} else {
			definition := d.allocDefinition(node)
			if rValue != nil {
				d.ComputeDependence(definition, rValue)
			}
			def.reachingDef.Put(var_, definition)
		}
	}
}

// Java: MustBeReachingVariableDef.escapeParameters(MustDef)
func (d *MustBeReachingVariableDef) EscapeParameters(output *MustBeReachingVariableDef_MustDef) {
	output = d.ensureMutable(output)
	if d.layout != nil && output.reachingDef.layout == d.layout {
		if d.hasParams {
			m := &output.reachingDef
			if !d.layout.wide {
				for id := range d.layout.vars {
					if d.paramMask&(uint64(1)<<id) != 0 {
						m.putID(uint16(id), nil)
					}
				}
				for _, id := range m.order {
					if value := m.values[id]; value != nil && value.depMask&d.paramMask != 0 {
						m.values[id] = nil
					}
				}
			} else {
				for id := range d.layout.vars {
					if d.paramWords[id>>6]&(uint64(1)<<(id&63)) != 0 {
						m.putID(uint16(id), nil)
					}
				}
				nWords := d.layout.numWords
				for _, id := range m.order {
					if value := m.values[id]; value != nil {
						for w := 0; w < nWords; w++ {
							if value.depWords[w]&d.paramWords[w] != 0 {
								m.values[id] = nil
								break
							}
						}
					}
				}
			}
		}
		return
	}
	for v := range d.allVarsInFn.Values() {
		if MustBeReachingVariableDef_IsParameter(v) {
			// Assume we no longer know where the parameter comes from
			// anymore.
			output.reachingDef.Put(v, nil)
		}
	}

	// Also, assume we no longer know anything that depends on a parameter.
	output.reachingDef.invalidateParameters()
}

// Java: MustBeReachingVariableDef.isParameter(Var)
func MustBeReachingVariableDef_IsParameter(v *Var) bool {
	return v.IsParam()
}

func (d *MustBeReachingVariableDef) computeDependenceWalk(def *MustBeReachingVariableDef_Definition, n, parent *rhino.Node) {
	if parent != nil && ControlFlowGraph_IsEnteringNewCfgNode(n) {
		return
	}
	if n.GetToken() == rhino.Token_CLASS_MEMBERS {
		for child := n.GetFirstChild(); child != nil; child = child.GetNext() {
			if child.IsComputedProp() || child.IsComputedFieldDef() {
				d.computeDependenceWalk(def, child.GetFirstChild(), child)
			}
		}
		for child := n.GetFirstChild(); child != nil; child = child.GetNext() {
			switch child.GetToken() {
			case rhino.Token_COMPUTED_PROP, rhino.Token_COMPUTED_FIELD_DEF:
				if !ControlFlowGraph_IsEnteringNewCfgNode(child) && child.HasTwoChildren() {
					d.computeDependenceWalk(def, child.GetLastChild(), child)
				}
			default:
				d.computeDependenceWalk(def, child, n)
			}
		}
		return
	}
	for c := n.GetFirstChild(); c != nil; c = c.GetNext() {
		d.computeDependenceWalk(def, c, n)
	}
	if n.IsName() {
		name := n.GetString()
		if d.layout != nil {
			id, dep, ok := d.lookupLocalVar(name)
			if !ok {
				def.unknownDependencies = true
			} else if !d.layout.wide {
				bit := uint64(1) << id
				if def.depMask&bit == 0 {
					def.depMask |= bit
					def.depends.Add(dep)
				}
			} else {
				w := id >> 6
				bit := uint64(1) << (id & 63)
				if def.depWords[w]&bit == 0 {
					def.depWords[w] |= bit
					def.depMask = 1
					def.depends.Add(dep)
				}
			}
		} else {
			dep := d.allVarsInFn.Get(name)
			if dep == nil {
				def.unknownDependencies = true
			} else {
				def.depends.Add(dep)
			}
		}
	}
}

// ComputeDependence computes all the local variables that rValue reads from
// and store that in the def's depends set.
//
// Java: MustBeReachingVariableDef.computeDependence(Definition, Node)
func (d *MustBeReachingVariableDef) ComputeDependence(def *MustBeReachingVariableDef_Definition, rValue *rhino.Node) {
	jutil.CheckNotNull(d.compiler)
	if rValue == nil {
		d.depCallback.def = def
		traverseCfgNode(d.compiler, rValue, &d.depCallback)
		d.depCallback.def = nil
		return
	}
	d.computeDependenceWalk(def, rValue, nil)
}

// mustBeReachingVariableDef_dependenceCallback is the anonymous
// AbstractCfgNodeTraversalCallback of computeDependence.
type mustBeReachingVariableDef_dependenceCallback struct {
	ControlFlowGraph_AbstractCfgNodeTraversalCallbackBase
	outer *MustBeReachingVariableDef
	def   *MustBeReachingVariableDef_Definition
}

// Java: MustBeReachingVariableDef.computeDependence(...).new AbstractCfgNodeTraversalCallback() {...}.visit(NodeTraversal, Node, Node)
func (c *mustBeReachingVariableDef_dependenceCallback) Visit(t *NodeTraversal, n *rhino.Node, parent *rhino.Node) {
	if n.IsName() {
		dep := c.outer.allVarsInFn.Get(n.GetString())
		if dep == nil {
			c.def.unknownDependencies = true
		} else {
			c.def.depends.Add(dep)
		}
	}
}

func (d *MustBeReachingVariableDef) getDefFromMustDef(name string, inState *MustBeReachingVariableDef_MustDef) *MustBeReachingVariableDef_Definition {
	m := &inState.reachingDef
	if d.layout != nil && m.layout == d.layout {
		id, _, ok := d.lookupLocalVar(name)
		if !ok || !m.hasID(id) {
			return nil
		}
		return m.values[id]
	}
	return m.Get(d.allVarsInFn.Get(name))
}

// GetDef gets the must-be-reaching definition of a given use node.
//
// name is the name of the variable. It can only be names of local variable
// that are not function parameters, escaped variables or variables declared
// in catch. useNode is the location of the use where the definition
// reaches.
//
// Java: MustBeReachingVariableDef.getDef(String, Node)
func (d *MustBeReachingVariableDef) GetDef(name string, useNode *rhino.Node) *MustBeReachingVariableDef_Definition {
	n := d.GetCfg().GetNode(useNode)
	jutil.CheckArgument(n != nil)
	state := n.GetAnnotation().(*DataFlowAnalysis_LinearFlowState[*MustBeReachingVariableDef_MustDef])
	return d.getDefFromMustDef(name, state.GetIn())
}

// Java: MustBeReachingVariableDef.getDefNode(String, Node)
func (d *MustBeReachingVariableDef) GetDefNode(name string, useNode *rhino.Node) *rhino.Node {
	def := d.GetDef(name, useNode)
	if def == nil {
		return nil
	}
	return def.node
}

// Java: MustBeReachingVariableDef.dependsOnOuterScopeVars(Definition)
func (d *MustBeReachingVariableDef) DependsOnOuterScopeVars(def *MustBeReachingVariableDef_Definition) bool {
	if def.unknownDependencies {
		return true
	}
	if d.layout != nil && (def.depMask != 0 || def.depends.IsEmpty()) {
		if !d.hasCatchVars {
			return false
		}
		if !d.layout.wide {
			return def.depMask&d.catchMask != 0
		}
		for w := 0; w < d.layout.numWords; w++ {
			if def.depWords[w]&d.catchWords[w] != 0 {
				return true
			}
		}
		return false
	}

	for s := range def.depends.All() {
		// Don't inline try catch
		if s.GetScope().(*Scope).IsCatchScope() {
			return true
		}
	}
	return false
}
