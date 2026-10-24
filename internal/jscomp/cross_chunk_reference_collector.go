// Port of com/google/javascript/jscomp/CrossChunkReferenceCollector.java @ closure-compiler v20261021

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
	"iter"
	"slices"
	"sync"
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

type crossChunkArena struct {
	topLevelStatements []*CrossChunkReferenceCollector_TopLevelStatement
	stmtChunks         [][]CrossChunkReferenceCollector_TopLevelStatement
	stmtChunkIdx       int
	stmtRefsChunks     [][]*Reference
	stmtRefsChunkIdx   int
	stmtVarsChunks     [][]*Var
	stmtVarsChunkIdx   int
	symbolsByVarIndex  []*CrossChunkCodeMotion_GlobalSymbol
	symSlab            []CrossChunkCodeMotion_GlobalSymbol
	symChunks          [][]CrossChunkCodeMotion_GlobalSymbol
	symChunkIdx        int
	dsgSlab            []CrossChunkCodeMotion_DeclarationStatementGroup
	dsgChunks          [][]CrossChunkCodeMotion_DeclarationStatementGroup
	dsgChunkIdx        int
	cycleSlab          []CrossChunkCodeMotion_GlobalSymbolCycle
	cycleChunks        [][]CrossChunkCodeMotion_GlobalSymbolCycle
	cycleChunkIdx      int
	componentContents  []*CrossChunkCodeMotion_GlobalSymbol
	componentRoots     []*CrossChunkCodeMotion_GlobalSymbol
	cycles             []*CrossChunkCodeMotion_GlobalSymbolCycle
}

var cachedCCAArena atomic.Pointer[crossChunkArena]

var crossChunkArenaPool = sync.Pool{
	New: func() any {
		return &crossChunkArena{}
	},
}

func acquireCCArena() *crossChunkArena {
	if a := cachedCCAArena.Swap(nil); a != nil {
		return a
	}
	return crossChunkArenaPool.Get().(*crossChunkArena)
}

func releaseCCArena(a *crossChunkArena) {
	if cachedCCAArena.CompareAndSwap(nil, a) {
		return
	}
	crossChunkArenaPool.Put(a)
}

// CrossChunkReferenceCollector collects global variable references for use
// by CrossChunkCodeMotion.
type CrossChunkReferenceCollector struct {
	// Maps global variable name to the corresponding Var object.
	varsByName *jutil.LinkedHashMap[string, *Var]

	// Maps a given variable to a collection of references to that name. Note
	// that Var objects are not stable across multiple traversals (unlike scope
	// root or name).
	referenceMap *jutil.ObjectLinkedHashMap[*Var, *ReferenceCollection]

	vars           []*Var
	cols           []*ReferenceCollection
	varsBuf        [16]*Var
	colsBuf        [16]*ReferenceCollection
	lastVar        *Var
	lastCollection *ReferenceCollection
	maxVarIndex    int
	ccArena        *crossChunkArena
	poolable       bool

	// The stack of basic blocks and scopes the current traversal is in.
	blockStack    []*BasicBlock
	blockStackBuf [16]*BasicBlock

	// Go-only: performance — one collection's blocks and references die
	// together with it (see jutil.Slab).
	blocks  jutil.Slab[BasicBlock]
	refSlab jutil.Slab[Reference]

	// List of all top-level statements in the order they appear in the AST.
	topLevelStatements []*CrossChunkReferenceCollector_TopLevelStatement

	scopeCreator ScopeCreator

	// JavaScript compiler to use in traversing.
	compiler AbstractCompiler

	statementCounter       int
	topLevelStatementDraft *CrossChunkReferenceCollector_TopLevelStatementDraft
	draftVal               CrossChunkReferenceCollector_TopLevelStatementDraft
	draftRefsBuf           [16]*Reference
	draftVarsBuf           [16]*Var

	stmtSlab     []CrossChunkReferenceCollector_TopLevelStatement
	stmtBuf      [4]CrossChunkReferenceCollector_TopLevelStatement
	stmtRefsSlab []*Reference
	stmtVarsSlab []*Var
}

// NewCrossChunkReferenceCollector: constructor initializes block stack.
//
// Java: CrossChunkReferenceCollector(AbstractCompiler, ScopeCreator)
func NewCrossChunkReferenceCollector(compiler AbstractCompiler, creator ScopeCreator) *CrossChunkReferenceCollector {
	c := &CrossChunkReferenceCollector{
		compiler:     compiler,
		scopeCreator: creator,
		maxVarIndex:  -1,
	}
	c.vars = c.varsBuf[:0]
	c.cols = c.colsBuf[:0]
	c.blockStack = c.blockStackBuf[:0]
	c.draftVal.nonDeclarationReferences = c.draftRefsBuf[:0]
	c.draftVal.nonDeclarationVars = c.draftVarsBuf[:0]
	c.stmtSlab = c.stmtBuf[:0]
	return c
}

func (c *CrossChunkReferenceCollector) enableArenaPool() {
	c.poolable = true
	ca := c.getCCArena()
	if cap(ca.topLevelStatements) > 0 {
		c.topLevelStatements = ca.topLevelStatements[:0]
	}
}

func (c *CrossChunkReferenceCollector) getCCArena() *crossChunkArena {
	if c.ccArena == nil {
		c.ccArena = acquireCCArena()
	}
	return c.ccArena
}

func (c *CrossChunkReferenceCollector) release() {
	for _, kv := range c.vars {
		kv.nameHash = 0
	}
	if ca := c.ccArena; ca != nil {
		c.ccArena = nil
		if cap(c.topLevelStatements) > 0 {
			if cap(c.topLevelStatements) <= 524288 {
				clear(c.topLevelStatements)
				ca.topLevelStatements = c.topLevelStatements[:0]
			} else {
				ca.topLevelStatements = nil
			}
		}
		c.topLevelStatements = nil
		const (
			maxRetainedStmtChunks     = 256
			maxRetainedStmtRefsChunks = 128
			maxRetainedSymChunks      = 128
			maxRetainedDsgChunks      = 128
			maxRetainedCycleChunks    = 64
		)
		if ca.stmtChunkIdx > 0 {
			ca.stmtChunks[ca.stmtChunkIdx-1] = c.stmtSlab
		}
		for i := 0; i < ca.stmtChunkIdx && i < maxRetainedStmtChunks; i++ {
			ch := ca.stmtChunks[i]
			clear(ch)
			ca.stmtChunks[i] = ch[:0]
		}
		if len(ca.stmtChunks) > maxRetainedStmtChunks {
			clear(ca.stmtChunks[maxRetainedStmtChunks:])
			ca.stmtChunks = ca.stmtChunks[:maxRetainedStmtChunks]
		}
		if ca.stmtRefsChunkIdx > 0 {
			ca.stmtRefsChunks[ca.stmtRefsChunkIdx-1] = c.stmtRefsSlab
		}
		for i := 0; i < ca.stmtRefsChunkIdx && i < maxRetainedStmtRefsChunks; i++ {
			ch := ca.stmtRefsChunks[i]
			clear(ch)
			ca.stmtRefsChunks[i] = ch[:0]
		}
		if len(ca.stmtRefsChunks) > maxRetainedStmtRefsChunks {
			clear(ca.stmtRefsChunks[maxRetainedStmtRefsChunks:])
			ca.stmtRefsChunks = ca.stmtRefsChunks[:maxRetainedStmtRefsChunks]
		}
		if ca.stmtVarsChunkIdx > 0 {
			ca.stmtVarsChunks[ca.stmtVarsChunkIdx-1] = c.stmtVarsSlab
		}
		for i := 0; i < ca.stmtVarsChunkIdx && i < maxRetainedStmtRefsChunks; i++ {
			ch := ca.stmtVarsChunks[i]
			clear(ch)
			ca.stmtVarsChunks[i] = ch[:0]
		}
		if len(ca.stmtVarsChunks) > maxRetainedStmtRefsChunks {
			clear(ca.stmtVarsChunks[maxRetainedStmtRefsChunks:])
			ca.stmtVarsChunks = ca.stmtVarsChunks[:maxRetainedStmtRefsChunks]
		}
		if ca.symChunkIdx > 0 {
			ca.symChunks[ca.symChunkIdx-1] = ca.symSlab
		}
		for i := 0; i < ca.symChunkIdx && i < maxRetainedSymChunks; i++ {
			ch := ca.symChunks[i]
			clear(ch)
			ca.symChunks[i] = ch[:0]
		}
		if len(ca.symChunks) > maxRetainedSymChunks {
			clear(ca.symChunks[maxRetainedSymChunks:])
			ca.symChunks = ca.symChunks[:maxRetainedSymChunks]
		}
		ca.symSlab = nil
		if ca.dsgChunkIdx > 0 {
			ca.dsgChunks[ca.dsgChunkIdx-1] = ca.dsgSlab
		}
		for i := 0; i < ca.dsgChunkIdx && i < maxRetainedDsgChunks; i++ {
			ch := ca.dsgChunks[i]
			clear(ch)
			ca.dsgChunks[i] = ch[:0]
		}
		if len(ca.dsgChunks) > maxRetainedDsgChunks {
			clear(ca.dsgChunks[maxRetainedDsgChunks:])
			ca.dsgChunks = ca.dsgChunks[:maxRetainedDsgChunks]
		}
		ca.dsgSlab = nil
		if ca.cycleChunkIdx > 0 {
			ca.cycleChunks[ca.cycleChunkIdx-1] = ca.cycleSlab
		}
		for i := 0; i < ca.cycleChunkIdx && i < maxRetainedCycleChunks; i++ {
			ch := ca.cycleChunks[i]
			clear(ch)
			ca.cycleChunks[i] = ch[:0]
		}
		if len(ca.cycleChunks) > maxRetainedCycleChunks {
			clear(ca.cycleChunks[maxRetainedCycleChunks:])
			ca.cycleChunks = ca.cycleChunks[:maxRetainedCycleChunks]
		}
		ca.cycleSlab = nil
		if len(ca.symbolsByVarIndex) > 0 {
			if cap(ca.symbolsByVarIndex) <= 262144 {
				clear(ca.symbolsByVarIndex)
				ca.symbolsByVarIndex = ca.symbolsByVarIndex[:0]
			} else {
				ca.symbolsByVarIndex = nil
			}
		}
		if cap(ca.componentContents) > 262144 {
			ca.componentContents = nil
		}
		if cap(ca.componentRoots) > 262144 {
			ca.componentRoots = nil
		}
		if len(ca.cycles) > 0 {
			if cap(ca.cycles) <= 262144 {
				clear(ca.cycles)
				ca.cycles = ca.cycles[:0]
			} else {
				ca.cycles = nil
			}
		}
		ca.stmtChunkIdx = 0
		ca.stmtRefsChunkIdx = 0
		ca.stmtVarsChunkIdx = 0
		ca.symChunkIdx = 0
		ca.dsgChunkIdx = 0
		ca.cycleChunkIdx = 0
		releaseCCArena(ca)
	}
}

func (c *CrossChunkReferenceCollector) allocStmtRefs(n int) []*Reference {
	if cap(c.stmtRefsSlab)-len(c.stmtRefsSlab) < n {
		if c.poolable {
			ca := c.getCCArena()
			if ca.stmtRefsChunkIdx > 0 {
				ca.stmtRefsChunks[ca.stmtRefsChunkIdx-1] = c.stmtRefsSlab
			}
			if ca.stmtRefsChunkIdx < len(ca.stmtRefsChunks) {
				ch := ca.stmtRefsChunks[ca.stmtRefsChunkIdx]
				if cap(ch) < n {
					ch = make([]*Reference, 0, max(n, 4096))
					ca.stmtRefsChunks[ca.stmtRefsChunkIdx] = ch
				}
				c.stmtRefsSlab = ch[:0]
				ca.stmtRefsChunkIdx++
			} else {
				newCap := 256
				if cap(c.stmtRefsSlab) >= 256 {
					newCap = min(cap(c.stmtRefsSlab)*2, 4096)
				}
				if newCap < n {
					newCap = n
				}
				c.stmtRefsSlab = make([]*Reference, 0, newCap)
				ca.stmtRefsChunks = append(ca.stmtRefsChunks, c.stmtRefsSlab)
				ca.stmtRefsChunkIdx++
			}
		} else {
			newCap := 256
			if cap(c.stmtRefsSlab) >= 256 {
				newCap = min(cap(c.stmtRefsSlab)*2, 4096)
			}
			if newCap < n {
				newCap = n
			}
			c.stmtRefsSlab = make([]*Reference, 0, newCap)
		}
	}
	off := len(c.stmtRefsSlab)
	c.stmtRefsSlab = c.stmtRefsSlab[:off+n]
	return c.stmtRefsSlab[off : off+n : off+n]
}

func (c *CrossChunkReferenceCollector) allocStmtVars(n int) []*Var {
	if cap(c.stmtVarsSlab)-len(c.stmtVarsSlab) < n {
		if c.poolable {
			ca := c.getCCArena()
			if ca.stmtVarsChunkIdx > 0 {
				ca.stmtVarsChunks[ca.stmtVarsChunkIdx-1] = c.stmtVarsSlab
			}
			if ca.stmtVarsChunkIdx < len(ca.stmtVarsChunks) {
				ch := ca.stmtVarsChunks[ca.stmtVarsChunkIdx]
				if cap(ch) < n {
					ch = make([]*Var, 0, max(n, 4096))
					ca.stmtVarsChunks[ca.stmtVarsChunkIdx] = ch
				}
				c.stmtVarsSlab = ch[:0]
				ca.stmtVarsChunkIdx++
			} else {
				newCap := 256
				if cap(c.stmtVarsSlab) >= 256 {
					newCap = min(cap(c.stmtVarsSlab)*2, 4096)
				}
				if newCap < n {
					newCap = n
				}
				c.stmtVarsSlab = make([]*Var, 0, newCap)
				ca.stmtVarsChunks = append(ca.stmtVarsChunks, c.stmtVarsSlab)
				ca.stmtVarsChunkIdx++
			}
		} else {
			newCap := 256
			if cap(c.stmtVarsSlab) >= 256 {
				newCap = min(cap(c.stmtVarsSlab)*2, 4096)
			}
			if newCap < n {
				newCap = n
			}
			c.stmtVarsSlab = make([]*Var, 0, newCap)
		}
	}
	off := len(c.stmtVarsSlab)
	c.stmtVarsSlab = c.stmtVarsSlab[:off+n]
	return c.stmtVarsSlab[off : off+n : off+n]
}

func (c *CrossChunkReferenceCollector) allocTopLevelStatement(draft *CrossChunkReferenceCollector_TopLevelStatementDraft) *CrossChunkReferenceCollector_TopLevelStatement {
	if len(c.stmtSlab) == cap(c.stmtSlab) {
		if c.poolable {
			ca := c.getCCArena()
			if ca.stmtChunkIdx > 0 {
				ca.stmtChunks[ca.stmtChunkIdx-1] = c.stmtSlab
			}
			if ca.stmtChunkIdx < len(ca.stmtChunks) {
				c.stmtSlab = ca.stmtChunks[ca.stmtChunkIdx][:0]
				ca.stmtChunkIdx++
			} else {
				newCap := 64
				if cap(c.stmtSlab) >= 64 {
					newCap = min(cap(c.stmtSlab)*2, 2048)
				}
				c.stmtSlab = make([]CrossChunkReferenceCollector_TopLevelStatement, 0, newCap)
				ca.stmtChunks = append(ca.stmtChunks, c.stmtSlab)
				ca.stmtChunkIdx++
			}
		} else {
			newCap := 64
			if cap(c.stmtSlab) >= 64 {
				newCap = min(cap(c.stmtSlab)*2, 2048)
			}
			c.stmtSlab = make([]CrossChunkReferenceCollector_TopLevelStatement, 0, newCap)
		}
	}
	idx := len(c.stmtSlab)
	c.stmtSlab = c.stmtSlab[:idx+1]
	s := &c.stmtSlab[idx]
	var refs []*Reference
	var vars []*Var
	if n := len(draft.nonDeclarationReferences); n > 0 {
		refs = c.allocStmtRefs(n)
		copy(refs, draft.nonDeclarationReferences)
		if len(draft.nonDeclarationVars) == n {
			vars = c.allocStmtVars(n)
			copy(vars, draft.nonDeclarationVars)
		}
	}
	s.outer = c
	s.originalOrder = draft.originalOrder
	s.chunk = draft.chunk
	s.statementNode = draft.statementNode
	s.nonDeclarationReferences = refs
	s.nonDeclarationVars = vars
	s.declaredNameReference = draft.declaredNameReference
	s.declaredVar = draft.declaredVar
	s.declaredValueNode = draft.declaredValueNode
	return s
}

func (c *CrossChunkReferenceCollector) growVarsCols() {
	newCap := max(cap(c.vars)*2, 4096)
	newVars := make([]*Var, len(c.vars), newCap)
	copy(newVars, c.vars)
	c.vars = newVars
	newCols := make([]*ReferenceCollection, len(c.cols), newCap)
	copy(newCols, c.cols)
	c.cols = newCols
}

func (c *CrossChunkReferenceCollector) getExistingReferences(v *Var) *ReferenceCollection {
	if idx := int(v.nameHash) - 1; idx >= 0 && idx < len(c.vars) && c.vars[idx] == v {
		return c.cols[idx]
	}
	return nil
}

func (c *CrossChunkReferenceCollector) createReferencesForNewVar(v *Var) *ReferenceCollection {
	if idx := v.GetIndex(); idx > c.maxVarIndex {
		c.maxVarIndex = idx
	}
	if c.varsByName != nil {
		name := v.GetName()
		if c.varsByName.ContainsKey(name) {
			jutil.CheckState(jutil.ObjectsEquals(c.varsByName.Get(name), v))
		} else {
			c.varsByName.Put(name, v)
		}
	}
	col := &ReferenceCollection{}
	if len(c.vars) == cap(c.vars) {
		c.growVarsCols()
	}
	c.vars = append(c.vars, v)
	c.cols = append(c.cols, col)
	if c.referenceMap != nil {
		c.referenceMap.Put(v, col)
	}
	v.nameHash = int32(len(c.vars))
	c.lastVar = v
	c.lastCollection = col
	return col
}

// Process is a convenience method for running this pass over a tree with
// this class as a callback.
//
// Java: CrossChunkReferenceCollector.process(Node, Node)
func (c *CrossChunkReferenceCollector) Process(externs *rhino.Node, root *rhino.Node) {
	jutil.CheckStateMsg(len(c.topLevelStatements) == 0 && c.statementCounter == 0, "process() called more than once")
	b := NodeTraversal_Builder{
		compiler:     c.compiler,
		callback:     c,
		scopeCreator: c.scopeCreator,
	}
	b.TraverseRoots(externs, root)
}

// ProcessRoot is the one-argument overload process(Node).
//
// Go-only: language-adapter — the two-argument overload keeps the plain
// name because it implements CompilerPass.process(Node, Node).
//
// Java: CrossChunkReferenceCollector.process(Node)
func (c *CrossChunkReferenceCollector) ProcessRoot(root *rhino.Node) {
	jutil.CheckStateMsg(len(c.topLevelStatements) == 0 && c.statementCounter == 0, "process() called more than once")
	b := NodeTraversal_Builder{
		compiler:     c.compiler,
		callback:     c,
		scopeCreator: c.scopeCreator,
	}
	b.Traverse(root)
}

// GetAllSymbols gets the variables that were referenced in this callback.
//
// Java: CrossChunkReferenceCollector.getAllSymbols()
func (c *CrossChunkReferenceCollector) GetAllSymbols() iter.Seq[*Var] {
	if c.referenceMap != nil {
		return c.referenceMap.Keys()
	}
	return slices.Values(c.vars)
}

// GetReferences gets the reference collection for the given variable.
//
// Java: CrossChunkReferenceCollector.getReferences(Var)
func (c *CrossChunkReferenceCollector) GetReferences(v *Var) *ReferenceCollection {
	if v == nil {
		if c.referenceMap != nil {
			return c.referenceMap.Get(nil)
		}
		return nil
	}
	if v == c.lastVar {
		return c.lastCollection
	}
	if col := c.getExistingReferences(v); col != nil {
		c.lastVar = v
		c.lastCollection = col
		return col
	}
	for i, kv := range c.vars {
		if kv.nameID == v.nameID && kv.scope.GetRootNode() == v.scope.GetRootNode() {
			col := c.cols[i]
			c.lastVar = v
			c.lastCollection = col
			return col
		}
	}
	if c.referenceMap != nil {
		return c.referenceMap.Get(v)
	}
	return nil
}

// Java: CrossChunkReferenceCollector.getGlobalVariableNamesMap()
func (c *CrossChunkReferenceCollector) GetGlobalVariableNamesMap() *jutil.LinkedHashMap[string, *Var] {
	if c.varsByName != nil {
		return jutil.NewLinkedHashMapFrom(c.varsByName)
	}
	out := jutil.NewLinkedHashMapWithCapacity[string, *Var](len(c.vars))
	for _, v := range c.vars {
		out.Put(v.GetName(), v)
	}
	return out
}

// Visit, for each node, updates the block stack and reference collection as
// appropriate.
//
// Java: CrossChunkReferenceCollector.visit(NodeTraversal, Node, Node)
func (c *CrossChunkReferenceCollector) Visit(t *NodeTraversal, n *rhino.Node, parent *rhino.Node) {
	if c.topLevelStatementDraft != nil {
		if n == c.topLevelStatementDraft.statementNode {
			c.topLevelStatements = append(c.topLevelStatements, c.allocTopLevelStatement(c.topLevelStatementDraft))
			c.topLevelStatementDraft = nil
		} else if n.IsName() {
			id := n.GetStringID()
			if id > rhino.RhinoStringPool_EmptyID {
				scope := t.resolveLookupScope(len(t.scopes) - 1)
				v := scope.GetVarByID(id)

				if v != nil && v.IsGlobal() {
					var col *ReferenceCollection
					if v == c.lastVar {
						col = c.lastCollection
					} else {
						col = c.getExistingReferences(v)
						if col == nil {
							if c.compiler.GetCodingConvention().IsExportedWithLocal(v.GetName() /* local= */, false) {
								goto checkBoundary
							}
							col = c.createReferencesForNewVar(v)
						} else {
							c.lastVar = v
							c.lastCollection = col
						}
					}
					reference := c.refSlab.New()
					*reference = Reference{nameNode: n, basicBlock: crossChunkReferenceCollector_peek(c.blockStack), scope: scope, inputId: t.GetInput().GetInputId()}
					if n == c.topLevelStatementDraft.declaredNameNode {
						c.topLevelStatementDraft.declaredNameReference = reference
						c.topLevelStatementDraft.declaredVar = v
					} else {
						c.topLevelStatementDraft.nonDeclarationReferences = append(c.topLevelStatementDraft.nonDeclarationReferences, reference)
						c.topLevelStatementDraft.nonDeclarationVars = append(c.topLevelStatementDraft.nonDeclarationVars, v)
					}
					col.references.Add(reference)
				}
			}
		}
	}
checkBoundary:
	if CrossChunkReferenceCollector_IsBlockBoundary(n, parent) {
		crossChunkReferenceCollector_pop(&c.blockStack)
	}
}

// EnterScope updates block stack and invokes any additional behavior.
//
// Java: CrossChunkReferenceCollector.enterScope(NodeTraversal)
func (c *CrossChunkReferenceCollector) EnterScope(t *NodeTraversal) {
	n := t.GetScopeRoot()
	var parent *BasicBlock
	if len(c.blockStack) != 0 {
		parent = crossChunkReferenceCollector_peek(c.blockStack)
	}
	// Don't add all ES6 scope roots to blockStack, only those that are also scopes according to
	// the ES5 scoping rules. Other nodes that ought to be considered the root of a BasicBlock
	// are added in shouldTraverse() and removed in visit().
	if t.IsHoistScope() {
		c.blockStack = append(c.blockStack, initBasicBlock(c.blocks.New(), parent, n))
	}
}

// ExitScope updates block stack and invokes any additional behavior.
//
// Java: CrossChunkReferenceCollector.exitScope(NodeTraversal)
func (c *CrossChunkReferenceCollector) ExitScope(t *NodeTraversal) {
	if t.IsHoistScope() {
		crossChunkReferenceCollector_pop(&c.blockStack)
	}
}

// Java: CrossChunkReferenceCollector.shouldTraverse(NodeTraversal, Node, Node)
func (c *CrossChunkReferenceCollector) ShouldTraverse(t *NodeTraversal, n *rhino.Node, parent *rhino.Node) bool {
	if parent != nil && NodeUtil_IsTopLevel(parent) {
		jutil.CheckStateMsg1(c.topLevelStatementDraft == nil, "%s", n)
		c.topLevelStatementDraft = c.InitializeDraftStatement(t.GetChunk(), n)
	}
	// If node is a new basic block, put on basic block stack
	if CrossChunkReferenceCollector_IsBlockBoundary(n, parent) {
		c.blockStack = append(c.blockStack, initBasicBlock(c.blocks.New(), crossChunkReferenceCollector_peek(c.blockStack), n))
	}
	return true
}

// Java: CrossChunkReferenceCollector.initializeDraftStatement(JSChunk, Node)
func (c *CrossChunkReferenceCollector) InitializeDraftStatement(chunk *JSChunk, statementNode *rhino.Node) *CrossChunkReferenceCollector_TopLevelStatementDraft {
	draft := &c.draftVal
	refs := draft.nonDeclarationReferences[:0]
	vars := draft.nonDeclarationVars[:0]
	*draft = CrossChunkReferenceCollector_TopLevelStatementDraft{
		originalOrder:            c.statementCounter,
		chunk:                    chunk,
		statementNode:            statementNode,
		nonDeclarationReferences: refs,
		nonDeclarationVars:       vars,
	}
	c.statementCounter++
	// Determine whether this statement declares a name or not.
	// If so, save its name node and value node, if any.
	if NodeUtil_IsNameDeclaration(statementNode) {
		// variable declaration
		draft.declaredNameNode = statementNode.GetFirstChild()
		draft.declaredValueNode = statementNode.GetFirstFirstChild()
	} else if statementNode.IsClass() {
		draft.declaredNameNode = statementNode.GetFirstChild()
		draft.declaredValueNode = statementNode
	} else if statementNode.IsFunction() {
		// function declaration
		draft.declaredNameNode = statementNode.GetFirstChild()
		draft.declaredValueNode = statementNode
	} else if statementNode.IsExprResult() {
		expr := jutil.CheckNotNull(statementNode.GetFirstChild())
		if expr.IsAssign() {
			lhs := jutil.CheckNotNull(expr.GetFirstChild())
			rhs := jutil.CheckNotNull(expr.GetSecondChild())
			if lhs.IsName() {
				// `varName = value;`
				draft.declaredNameNode = lhs
				draft.declaredValueNode = rhs
			} else if lhs.IsGetProp() || CrossChunkReferenceCollector_IsPotentiallyMoveableGetElem(lhs) {
				nameNode := jutil.CheckNotNull(lhs.GetFirstChild())
				for nameNode.IsGetProp() {
					nameNode = jutil.CheckNotNull(nameNode.GetFirstChild())
				}
				if nameNode.IsName() {
					// `varName.some.property = value;`
					// OR
					// `varName[someProperty] = value;`
					draft.declaredNameNode = nameNode
					draft.declaredValueNode = rhs
				}
			}
		} else if expr.IsCall() {
			var nameNode *rhino.Node
			var valueNode *rhino.Node
			relationship :=
				c.compiler.GetCodingConvention().GetClassesDefinedByCall(expr)
			if relationship != nil {
				// Check for $jscomp.inherits(SubC, SuperC), goog.inherits(Sub, SuperC), etc.
				declaredName := relationship.subclassName
				for callArg := expr.GetSecondChild(); callArg != nil; callArg = callArg.GetNext() {
					// We're assuming that the child class must be an argument to the function that
					// establishes its inheritance, which is true for `goog.inherits()` and
					// `$jscomp.inherits()`
					// TODO(bradfordcsmith): handle cases like `goog.inherits(x.ChildClass, SuperClass)`
					if callArg.IsName() && declaredName == callArg.GetString() {
						nameNode = callArg
						break
					}
				}
			} else if NodeUtil_IsObjectDefinePropertiesDefinition(expr) {
				// Check for $jscomp$global.Object.defineProperties.
				targetObject := expr.GetSecondChild()

				// Get the global var being referenced in the first parameter of Object.defineProperties.
				// Can be 'Foo' or 'Foo.prototype'.
				if targetObject.IsName() {
					nameNode = targetObject
				} else {
					nameNode = targetObject.GetFirstChild()
				}

				// Second parameter of defineProperties is the value node.
				valueNode = targetObject.GetNext()
			}

			if nameNode != nil {
				draft.declaredNameNode = nameNode
				draft.declaredValueNode = valueNode
			}
		}
	}
	return draft
}

// Java: CrossChunkReferenceCollector.isPotentiallyMoveableGetElem(Node)
func CrossChunkReferenceCollector_IsPotentiallyMoveableGetElem(maybeGetElem *rhino.Node) bool {
	if !maybeGetElem.IsGetElem() {
		return false
	}

	// We'll only consider a GETELEM as potentially moveable if the property being indexed is on
	// `Symbol` (e.g. `Symbol.iterator`). We don't generally want to consider GETELEM nodes as
	// moveable as registries (ex. `GlobalRegistry['Foo'] = Foo`) will create a reference from `Foo`
	// to `GlobalRegistry`. As other pieces of code reference `GlobalRegistry` this will ultimately
	// cause a large cluster of symbols that would have to be considered as moving together, but
	// likely will no to be able to.
	// Therefore, we we'll only consider a subset of cases, like `Symbol.iterator`, as moveable as
	// these large reference clusters are unlikely in this case. This should allow classes that
	// define these properties to potentialy be moved to later chunks.
	property := maybeGetElem.GetLastChild()
	if !property.IsGetProp() {
		return false
	}
	name := property.GetFirstChild()
	return name.IsName() && name.GetString() == "Symbol"
}

// Java: CrossChunkReferenceCollector.pop(List<T>)
func crossChunkReferenceCollector_pop[T any](list *[]T) T {
	l := *list
	if len(l) == 0 {
		// ArrayList.remove(-1)
		panic(jutil.NewIndexOutOfBoundsException("Index -1 out of bounds for length 0"))
	}
	last := l[len(l)-1]
	var zero T
	l[len(l)-1] = zero
	*list = l[:len(l)-1]
	return last
}

// Java: CrossChunkReferenceCollector.peek(List<T>)
func crossChunkReferenceCollector_peek[T any](list []T) T {
	if len(list) == 0 {
		// Iterables.getLast of an empty list.
		panic(jutil.NewNoSuchElementException())
	}
	return list[len(list)-1]
}

var crossChunkParentBoundaryKind = [256]uint8{
	rhino.Token_DO:           1,
	rhino.Token_FOR:          1,
	rhino.Token_FOR_IN:       1,
	rhino.Token_FOR_OF:       1,
	rhino.Token_FOR_AWAIT_OF: 1,
	rhino.Token_TRY:          1,
	rhino.Token_WHILE:        1,
	rhino.Token_WITH:         1,
	rhino.Token_CLASS:        1,
	rhino.Token_AND:          2,
	rhino.Token_HOOK:         2,
	rhino.Token_IF:           2,
	rhino.Token_OR:           2,
	rhino.Token_SWITCH:       2,
	rhino.Token_COALESCE:     2,
}

// CrossChunkReferenceCollector_IsBlockBoundary returns true if this node
// marks the start of a new basic block.
//
// Java: CrossChunkReferenceCollector.isBlockBoundary(Node, Node)
func CrossChunkReferenceCollector_IsBlockBoundary(n *rhino.Node, parent *rhino.Node) bool {
	if parent != nil {
		switch crossChunkParentBoundaryKind[uint8(parent.GetToken())] {
		case 1:
			return true
		case 2:
			// The first child of a conditional is not a boundary,
			// but all the rest of the children are.
			return n != parent.GetFirstChild()
		}
	}

	return n.IsCase()
}

// Java: CrossChunkReferenceCollector.addReferenceToCollection(Var, Reference)
func (c *CrossChunkReferenceCollector) AddReferenceToCollection(v *Var, reference *Reference) {
	var col *ReferenceCollection
	if v == c.lastVar && v != nil {
		col = c.lastCollection
	} else {
		col = c.getExistingReferences(v)
		if col == nil {
			col = c.createReferencesForNewVar(v)
		} else {
			c.lastVar = v
			c.lastCollection = col
		}
	}
	col.Add(reference)
}

// Java: CrossChunkReferenceCollector.getTopLevelStatements()
func (c *CrossChunkReferenceCollector) GetTopLevelStatements() []*CrossChunkReferenceCollector_TopLevelStatement {
	return c.topLevelStatements
}

// CanMoveValue determines whether the given value is eligible to be moved
// across chunks.
//
// Java: CrossChunkReferenceCollector.canMoveValue(Scope, Node)
func (c *CrossChunkReferenceCollector) CanMoveValue(scope *Scope, valueNode *rhino.Node) bool {
	// the value is only movable if it's
	// - nothing,
	// - a constant literal,
	// - a function
	// - class without computed properties, or whose computed properties expressions are moveable
	// - an array/object literal of movable values
	// - a template literal with movable substitutions
	// - a function stub generated by CrossChunkMethodMotion.
	if valueNode == nil || NodeUtil_IsLiteralValue(valueNode, true) || valueNode.IsFunction() {
		return true
	}

	// Consider the node eligible for move if the current statement or expression is annotated
	// with @pureOrBreakMyCode.
	jsDocInfo := NodeUtil_GetBestJSDocInfo(valueNode)
	if jsDocInfo != nil && jsDocInfo.IsPureOrBreakMyCode() {
		return true
	}

	switch valueNode.GetToken() {
	case rhino.Token_CLASS:
		{
			// NOTE: Why aren't class extends expressions checked here?
			classMembers := valueNode.GetLastChild()
			for member := classMembers.GetFirstChild(); member != nil; member = member.GetNext() {
				if member.IsComputedProp() {
					keyExpr := member.GetFirstChild()
					method := member.GetLastChild()
					jutil.CheckStateMsg1(method.IsFunction(), "%s", method)
					if !c.CanMoveValue(scope, keyExpr) {
						return false
					}
				} else if member.IsComputedFieldDef() {
					keyExpr := member.GetFirstChild()
					if !c.CanMoveValue(scope, keyExpr) {
						return false
					}
				} else if member.IsMemberFieldDef() {
					if member.IsStaticMember() && !c.CanMoveValue(scope, member.GetFirstChild()) {
						// Unlike non-static fields,the RHS of a static field is executed at
						// class definition time. So, we must check canMoveValue for static fields.
						return false
					}
				} else if member.IsBlock() {
					// TODO(bradfordcsmith): Ideally could move these in some cases, fix later
					return false
				} else {
					jutil.CheckStateMsg1(member.IsMemberFunctionDef() || NodeUtil_IsGetOrSetKey(member), "%s", member)
				}
			}
			return true
		}
	case rhino.Token_CALL:
		{
			// In general it is not safe to move function calls, but we carve out an exception
			// for the special stub method calls used for CrossChunkMethodMotion.
			// Case: `JSCompiler_stubMethod(x)`
			functionName := jutil.CheckNotNull(valueNode.GetFirstChild())
			return functionName.IsName() &&
				functionName.GetString() == CrossChunkMethodMotion_STUB_METHOD_NAME
		}
	case rhino.Token_ARRAYLIT:
		// Movable if all of the array values are movable.
		for child := valueNode.GetFirstChild(); child != nil; child = child.GetNext() {
			if !c.CanMoveValue(scope, child) {
				return false
			}
		}

		return true
	case rhino.Token_OBJECTLIT:
		// Movable if all of the keys and values are movable.
		for child := valueNode.GetFirstChild(); child != nil; child = child.GetNext() {
			switch child.GetToken() {
			case rhino.Token_MEMBER_FUNCTION_DEF, rhino.Token_GETTER_DEF, rhino.Token_SETTER_DEF:

			case rhino.Token_COMPUTED_PROP:
				if !c.CanMoveValue(scope, child.GetFirstChild()) ||
					!c.CanMoveValue(scope, child.GetLastChild()) {
					return false
				}

			case rhino.Token_STRING_KEY, rhino.Token_OBJECT_SPREAD:
				if !c.CanMoveValue(scope, child.GetOnlyChild()) {
					return false
				}

			default:
				panic(jutil.NewIllegalStateExceptionMsg(
					"Unexpected child of OBJECTLIT: " + child.ToStringTree()))
			}
		}

		return true
	case rhino.Token_NAME:
		{
			// If the value is guaranteed to never be changed after
			// this reference, then we can move it.
			v := scope.GetVar(valueNode.GetString())
			if v != nil && v.IsGlobal() {
				refCollection := c.GetReferences(v)
				if refCollection != nil &&
					refCollection.IsWellDefined() &&
					refCollection.IsAssignedOnceInLifetime() {
					return true
				}
			}
		}
	case rhino.Token_TEMPLATELIT:
		// A template literal is movable if all of the substitutions it contains are movable.
		for child := valueNode.GetFirstChild(); child != nil; child = child.GetNext() {
			if child.IsTemplateLitSub() {
				if !c.CanMoveValue(scope, child.GetFirstChild()) {
					return false
				}
			} else {
				jutil.CheckStateMsg1(child.IsTemplateLitString(), "%s", child)
			}
		}
		return true
	case rhino.Token_GETPROP:
		{
			// `Symbol` is a built-in JavaScript object that is guaranteed to be unique. Therefore it
			// is movable. Any properties on `Symbol` are also movable (e.g `Symbol.iterator`).
			child := valueNode.GetFirstChild()
			if child != nil && child.IsName() {
				symbol := child.GetString()
				if symbol == "Symbol" {
					return true
				}
			}
			return false
		}
	default:
	}

	return false
}

// CrossChunkReferenceCollector_TopLevelStatement represents a top-level
// statement and the references to global names it contains.
type CrossChunkReferenceCollector_TopLevelStatement struct {
	outer *CrossChunkReferenceCollector

	// 0-based index indicating original order of this statement in the source.
	originalOrder int

	chunk                    *JSChunk
	statementNode            *rhino.Node
	nonDeclarationReferences []*Reference
	nonDeclarationVars       []*Var
	declaredNameReference    *Reference
	declaredVar              *Var
	declaredValueNode        *rhino.Node
}

// Java: CrossChunkReferenceCollector.TopLevelStatement(TopLevelStatementDraft)
func NewCrossChunkReferenceCollector_TopLevelStatement(outer *CrossChunkReferenceCollector, draft *CrossChunkReferenceCollector_TopLevelStatementDraft) *CrossChunkReferenceCollector_TopLevelStatement {
	return &CrossChunkReferenceCollector_TopLevelStatement{
		outer:                    outer,
		originalOrder:            draft.originalOrder,
		chunk:                    draft.chunk,
		statementNode:            draft.statementNode,
		nonDeclarationReferences: draft.nonDeclarationReferences,
		nonDeclarationVars:       draft.nonDeclarationVars,
		declaredNameReference:    draft.declaredNameReference,
		declaredVar:              draft.declaredVar,
		declaredValueNode:        draft.declaredValueNode,
	}
}

// Java: CrossChunkReferenceCollector.TopLevelStatement.getOriginalOrder()
func (s *CrossChunkReferenceCollector_TopLevelStatement) GetOriginalOrder() int {
	return s.originalOrder
}

// Java: CrossChunkReferenceCollector.TopLevelStatement.getChunk()
func (s *CrossChunkReferenceCollector_TopLevelStatement) GetChunk() *JSChunk {
	return s.chunk
}

// Java: CrossChunkReferenceCollector.TopLevelStatement.getStatementNode()
func (s *CrossChunkReferenceCollector_TopLevelStatement) GetStatementNode() *rhino.Node {
	return s.statementNode
}

// Java: CrossChunkReferenceCollector.TopLevelStatement.getNonDeclarationReferences()
func (s *CrossChunkReferenceCollector_TopLevelStatement) GetNonDeclarationReferences() []*Reference {
	return s.nonDeclarationReferences
}

// Java: CrossChunkReferenceCollector.TopLevelStatement.isDeclarationStatement()
func (s *CrossChunkReferenceCollector_TopLevelStatement) IsDeclarationStatement() bool {
	return s.declaredNameReference != nil
}

// Java: CrossChunkReferenceCollector.TopLevelStatement.getDeclaredNameReference()
func (s *CrossChunkReferenceCollector_TopLevelStatement) GetDeclaredNameReference() *Reference {
	return jutil.CheckNotNull(s.declaredNameReference)
}

// Java: CrossChunkReferenceCollector.TopLevelStatement.getDeclaredValueNode()
func (s *CrossChunkReferenceCollector_TopLevelStatement) GetDeclaredValueNode() *rhino.Node {
	return s.declaredValueNode
}

// Java: CrossChunkReferenceCollector.TopLevelStatement.isMovableDeclaration()
func (s *CrossChunkReferenceCollector_TopLevelStatement) IsMovableDeclaration() bool {
	return s.IsDeclarationStatement() &&
		s.outer.CanMoveValue(s.declaredNameReference.GetScope(), s.declaredValueNode)
}

// CrossChunkReferenceCollector_TopLevelStatementDraft holds statement info
// temporarily while the statement is being traversed.
type CrossChunkReferenceCollector_TopLevelStatementDraft struct {
	// 0-based index indicating original order of this statement in the source.
	originalOrder int

	chunk                    *JSChunk
	statementNode            *rhino.Node
	nonDeclarationReferences []*Reference
	nonDeclarationVars       []*Var
	declaredValueNode        *rhino.Node
	declaredNameNode         *rhino.Node
	declaredNameReference    *Reference
	declaredVar              *Var
}

// Java: CrossChunkReferenceCollector.TopLevelStatementDraft(int, JSChunk, Node)
func NewCrossChunkReferenceCollector_TopLevelStatementDraft(originalOrder int, chunk *JSChunk, statementNode *rhino.Node) *CrossChunkReferenceCollector_TopLevelStatementDraft {
	return &CrossChunkReferenceCollector_TopLevelStatementDraft{
		originalOrder:            originalOrder,
		chunk:                    chunk,
		statementNode:            statementNode,
		nonDeclarationReferences: []*Reference{},
	}
}
