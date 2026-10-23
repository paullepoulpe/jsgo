// Port of com/google/javascript/jscomp/ReferenceCollector.java @ closure-compiler v20261015

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
	"slices"
	"sync"
	"sync/atomic"
	"unsafe"

	"github.com/paullepoulpe/jsgo/internal/base"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

func rcHashNodePtr(n *rhino.Node) uint32 {
	x := uint64(uintptr(unsafe.Pointer(n)) >> 4)
	x ^= x >> 33
	x *= 0xff51afd7ed558ccd
	x ^= x >> 33
	return uint32(x)
}

type basicBlockEntry struct {
	root  *rhino.Node
	block *BasicBlock
}

type scopedVarKey struct {
	scopeRoot *rhino.Node
	nameID    uint32
}

type rcVarSlabs struct {
	vars               []*Var
	cols               []*ReferenceCollection
	refSlab            []Reference
	refChunks          [][]Reference
	refChunkIdx        int
	refPtrSlab         []*Reference
	refPtrChunks       [][]*Reference
	refPtrChunkIdx     int
	colSlab            []ReferenceCollection
	colChunks          [][]ReferenceCollection
	colChunkIdx        int
	baseVarsLen        int
	baseRefChunkIdx    int
	baseRefSlabLen     int
	baseRefPtrChunkIdx int
	baseRefPtrSlabLen  int
	baseColChunkIdx    int
	baseColSlabLen     int
	pooled             bool
}

func (s *rcVarSlabs) addVarCol(v *Var, col *ReferenceCollection, minGrow int) {
	if len(s.vars) == cap(s.vars) {
		newCap := max(cap(s.vars)*2, minGrow)
		newVars := make([]*Var, len(s.vars), newCap)
		copy(newVars, s.vars)
		s.vars = newVars
		newCols := make([]*ReferenceCollection, len(s.cols), newCap)
		copy(newCols, s.cols)
		s.cols = newCols
	}
	s.vars = append(s.vars, v)
	s.cols = append(s.cols, col)
	v.nameHash = int32(len(s.vars))
}

func (s *rcVarSlabs) allocReference(n *rhino.Node, scope *Scope, basicBlock *BasicBlock, inputId *rhino.InputId) *Reference {
	if len(s.refSlab) == cap(s.refSlab) {
		if s.pooled {
			if s.refChunkIdx > 0 {
				s.refChunks[s.refChunkIdx-1] = s.refSlab
			}
			if s.refChunkIdx < len(s.refChunks) {
				s.refSlab = s.refChunks[s.refChunkIdx][:0]
				s.refChunkIdx++
			} else {
				newCap := 256
				if cap(s.refSlab) >= 256 {
					newCap = min(cap(s.refSlab)*2, 1024) // 1024 * 32B = 32 KiB == _MaxSmallSize
				}
				s.refSlab = make([]Reference, 0, newCap)
				s.refChunks = append(s.refChunks, s.refSlab)
				s.refChunkIdx++
			}
		} else {
			newCap := 256
			if cap(s.refSlab) >= 256 {
				newCap = min(cap(s.refSlab)*2, 1024)
			}
			s.refSlab = make([]Reference, 0, newCap)
		}
	}
	idx := len(s.refSlab)
	s.refSlab = s.refSlab[:idx+1]
	r := &s.refSlab[idx]
	r.nameNode = n
	r.basicBlock = basicBlock
	r.scope = scope
	r.inputId = inputId
	return r
}

func (s *rcVarSlabs) allocReferenceCollection() *ReferenceCollection {
	if len(s.colSlab) == cap(s.colSlab) {
		if s.pooled {
			if s.colChunkIdx > 0 {
				s.colChunks[s.colChunkIdx-1] = s.colSlab
			}
			if s.colChunkIdx < len(s.colChunks) {
				s.colSlab = s.colChunks[s.colChunkIdx][:0]
				s.colChunkIdx++
			} else {
				newCap := 64
				if cap(s.colSlab) >= 64 {
					newCap = min(cap(s.colSlab)*2, 1024) // 1024 * 32B = 32 KiB == _MaxSmallSize
				}
				s.colSlab = make([]ReferenceCollection, 0, newCap)
				s.colChunks = append(s.colChunks, s.colSlab)
				s.colChunkIdx++
			}
		} else {
			newCap := 64
			if cap(s.colSlab) >= 64 {
				newCap = min(cap(s.colSlab)*2, 1024)
			}
			s.colSlab = make([]ReferenceCollection, 0, newCap)
		}
	}
	idx := len(s.colSlab)
	s.colSlab = s.colSlab[:idx+1]
	col := &s.colSlab[idx]
	col.references.InitWithBuffer(s.allocRefPtrSlice(4))
	return col
}

func (s *rcVarSlabs) allocRefPtrSlice(n int) []*Reference {
	if len(s.refPtrSlab)+n > cap(s.refPtrSlab) {
		if s.pooled {
			if s.refPtrChunkIdx > 0 {
				s.refPtrChunks[s.refPtrChunkIdx-1] = s.refPtrSlab
			}
			if s.refPtrChunkIdx < len(s.refPtrChunks) {
				ch := s.refPtrChunks[s.refPtrChunkIdx]
				if cap(ch) < n {
					ch = make([]*Reference, 0, max(4096, n)) // 4096 * 8B = 32 KiB == _MaxSmallSize
					s.refPtrChunks[s.refPtrChunkIdx] = ch
				}
				s.refPtrSlab = ch[:0]
				s.refPtrChunkIdx++
			} else {
				s.refPtrSlab = make([]*Reference, 0, max(4096, n))
				s.refPtrChunks = append(s.refPtrChunks, s.refPtrSlab)
				s.refPtrChunkIdx++
			}
		} else {
			newCap := 256
			if cap(s.refPtrSlab) >= 256 {
				newCap = min(cap(s.refPtrSlab)*2, 4096)
			}
			if newCap < n {
				newCap = n
			}
			s.refPtrSlab = make([]*Reference, 0, newCap)
		}
	}
	idx := len(s.refPtrSlab)
	s.refPtrSlab = s.refPtrSlab[:idx+n]
	return s.refPtrSlab[idx : idx+n : idx+n]
}

func (s *rcVarSlabs) growReferenceCollection(col *ReferenceCollection, slice []*Reference) {
	oldCap := cap(slice)
	newCap := 16
	if oldCap >= 16 {
		newCap = oldCap * 2
	}
	newBuf := s.allocRefPtrSlice(newCap)
	col.references.InitWithBuffer(newBuf)
	col.references.AddAll(slice...)
}

func (s *rcVarSlabs) markBase() {
	s.baseVarsLen = len(s.vars)
	s.baseRefChunkIdx = s.refChunkIdx
	s.baseRefSlabLen = len(s.refSlab)
	s.baseRefPtrChunkIdx = s.refPtrChunkIdx
	s.baseRefPtrSlabLen = len(s.refPtrSlab)
	s.baseColChunkIdx = s.colChunkIdx
	s.baseColSlabLen = len(s.colSlab)
}

func (s *rcVarSlabs) clearBase() {
	s.baseVarsLen = 0
	s.baseRefChunkIdx = 0
	s.baseRefSlabLen = 0
	s.baseRefPtrChunkIdx = 0
	s.baseRefPtrSlabLen = 0
	s.baseColChunkIdx = 0
	s.baseColSlabLen = 0
}

func resetRCSlabChunks[T any](chunks [][]T, slab []T, chunkIdx, baseChunkIdx, baseSlabLen, maxChunks int) ([][]T, []T, int) {
	if chunkIdx > 0 {
		chunks[chunkIdx-1] = slab
	}
	if baseChunkIdx == 0 {
		for i := 0; i < chunkIdx && i < maxChunks; i++ {
			ch := chunks[i]
			clear(ch)
			chunks[i] = ch[:0]
		}
		if len(chunks) > maxChunks {
			clear(chunks[maxChunks:])
			chunks = chunks[:maxChunks]
		}
		return chunks, nil, 0
	}
	b := baseChunkIdx - 1
	chB := chunks[b]
	clear(chB[baseSlabLen:])
	chB = chB[:baseSlabLen]
	chunks[b] = chB
	for i := b + 1; i < chunkIdx; i++ {
		ch := chunks[i]
		clear(ch)
		chunks[i] = ch[:0]
	}
	return chunks, chB, baseChunkIdx
}

func (s *rcVarSlabs) resetToBase(maxRefChunks, maxRefPtrChunks, maxColChunks, maxVarsCap int) {
	if s.baseVarsLen == 0 && s.baseRefChunkIdx == 0 && s.baseRefPtrChunkIdx == 0 && s.baseColChunkIdx == 0 {
		s.reset(maxRefChunks, maxRefPtrChunks, maxColChunks, maxVarsCap)
		return
	}
	for _, kv := range s.vars[s.baseVarsLen:] {
		kv.nameHash = 0
	}
	clear(s.vars[s.baseVarsLen:])
	s.vars = s.vars[:s.baseVarsLen]
	clear(s.cols[s.baseVarsLen:])
	s.cols = s.cols[:s.baseVarsLen]
	s.refChunks, s.refSlab, s.refChunkIdx = resetRCSlabChunks(s.refChunks, s.refSlab, s.refChunkIdx, s.baseRefChunkIdx, s.baseRefSlabLen, maxRefChunks)
	s.refPtrChunks, s.refPtrSlab, s.refPtrChunkIdx = resetRCSlabChunks(s.refPtrChunks, s.refPtrSlab, s.refPtrChunkIdx, s.baseRefPtrChunkIdx, s.baseRefPtrSlabLen, maxRefPtrChunks)
	s.colChunks, s.colSlab, s.colChunkIdx = resetRCSlabChunks(s.colChunks, s.colSlab, s.colChunkIdx, s.baseColChunkIdx, s.baseColSlabLen, maxColChunks)
}

func (s *rcVarSlabs) reset(maxRefChunks, maxRefPtrChunks, maxColChunks, maxVarsCap int) {
	s.clearBase()
	for _, kv := range s.vars {
		kv.nameHash = 0
	}
	if cap(s.vars) <= maxVarsCap {
		clear(s.vars)
		s.vars = s.vars[:0]
		clear(s.cols)
		s.cols = s.cols[:0]
	} else {
		s.vars = nil
		s.cols = nil
	}
	s.refChunks, s.refSlab, s.refChunkIdx = resetRCSlabChunks(s.refChunks, s.refSlab, s.refChunkIdx, 0, 0, maxRefChunks)
	s.refPtrChunks, s.refPtrSlab, s.refPtrChunkIdx = resetRCSlabChunks(s.refPtrChunks, s.refPtrSlab, s.refPtrChunkIdx, 0, 0, maxRefPtrChunks)
	s.colChunks, s.colSlab, s.colChunkIdx = resetRCSlabChunks(s.colChunks, s.colSlab, s.colChunkIdx, 0, 0, maxColChunks)
}

type referenceCollectorArena struct {
	global        rcVarSlabs
	local         rcVarSlabs
	blkChunks     [][]BasicBlock
	blkChunkIdx   int
	hoistedTable  []*rhino.Node
	hoistedList   []*rhino.Node
	argumentsRefs map[*rhino.Node]*ReferenceCollection
}

var cachedRCArenas [4]atomic.Pointer[referenceCollectorArena]

var referenceCollectorArenaPool = sync.Pool{
	New: func() any {
		return &referenceCollectorArena{
			global: rcVarSlabs{pooled: true},
			local:  rcVarSlabs{pooled: true},
		}
	},
}

func acquireRCArena() *referenceCollectorArena {
	bestIdx := -1
	bestLen := -1
	for i := range cachedRCArenas {
		if a := cachedRCArenas[i].Load(); a != nil {
			if n := len(a.global.refChunks) + len(a.blkChunks); n > bestLen {
				bestLen = n
				bestIdx = i
			}
		}
	}
	if bestIdx >= 0 {
		if a := cachedRCArenas[bestIdx].Swap(nil); a != nil {
			a.global.pooled = true
			a.local.pooled = true
			return a
		}
	}
	for i := range cachedRCArenas {
		if a := cachedRCArenas[i].Swap(nil); a != nil {
			a.global.pooled = true
			a.local.pooled = true
			return a
		}
	}
	a := referenceCollectorArenaPool.Get().(*referenceCollectorArena)
	a.global.pooled = true
	a.local.pooled = true
	return a
}

func (a *referenceCollectorArena) trimSmall() {
	a.global.reset(64, 16, 32, 4096)
	a.local.reset(64, 16, 32, 4096)
	if len(a.blkChunks) > 32 {
		clear(a.blkChunks[32:])
		a.blkChunks = a.blkChunks[:32]
	}
}

func releaseRCArena(a *referenceCollectorArena) {
	for i := 0; i < 2; i++ {
		if cachedRCArenas[i].CompareAndSwap(nil, a) {
			return
		}
	}
	if old := cachedRCArenas[0].Load(); old != nil && len(a.global.refChunks)+len(a.blkChunks) > len(old.global.refChunks)+len(old.blkChunks) {
		if cachedRCArenas[0].CompareAndSwap(old, a) {
			if cachedRCArenas[1].CompareAndSwap(nil, old) {
				return
			}
			old.trimSmall()
			for i := 2; i < len(cachedRCArenas); i++ {
				if cachedRCArenas[i].CompareAndSwap(nil, old) {
					return
				}
			}
			referenceCollectorArenaPool.Put(old)
			return
		}
	}
	a.trimSmall()
	for i := 2; i < len(cachedRCArenas); i++ {
		if cachedRCArenas[i].CompareAndSwap(nil, a) {
			return
		}
	}
	referenceCollectorArenaPool.Put(a)
}

// ReferenceCollector is a helper class for passes that want to access all
// information about where a variable is referenced and declared at once and
// then make a decision as to how it should be handled, possibly inlining,
// reordering, or generating warnings. Callers do this by providing Behavior
// and then calling Process.
//
// Java: final class ReferenceCollector implements CompilerPass,
// StaticSymbolTable<Var, Reference>.
//
// Go-only: language-adapter — Java narrows StaticSymbolTable's getScope(Var)
// (to Scope) and getReferences(Var) (Iterable<Reference> to
// ReferenceCollection, which callers need and an iter.Seq cannot give back).
// Go has no covariant overrides, so GetScope and GetReferences keep the
// Java-declared types and AsStaticSymbolTable returns the
// rhino.StaticSymbolTable[*Var, *Reference] view of this collector.
type ReferenceCollector struct {
	// Maps a given variable to a collection of references to that name. Note that
	// Var objects are not stable across multiple traversals (unlike scope root or
	// name).
	//
	// Go-only: language-adapter — Var overrides equals/hashCode (ScopedName),
	// hence the Object map.
	referenceMap *jutil.ObjectLinkedHashMap[*Var, *ReferenceCollection]

	defaultSlabs     rcVarSlabs
	globalSlabs      *rcVarSlabs
	localSlabs       *rcVarSlabs
	varsBuf          [16]*Var
	colsBuf          [16]*ReferenceCollection
	refByKey         map[scopedVarKey]*ReferenceCollection
	argumentsRefs    map[*rhino.Node]*ReferenceCollection
	lastVar          *Var
	lastCollection   *ReferenceCollection
	lastRejectedVar       *Var
	activeTraversal       *NodeTraversal
	arena                 *referenceCollectorArena
	poolable              bool
	reclaimLocals         bool
	canReclaimLocals      func() bool
	localScopeDepth       int
	globalRefInLocal      bool
	globalRefInNonFnLocal bool
	baseBlkChunkIdx       int
	baseBlkSlabLen        int
	dummyFnScopes         [2]Scope
	lastDummyFn           *rhino.Node
	dummyFnIdx            uint8

	// The stack of basic blocks and scopes the current traversal is in.
	//
	// Go-only: language-adapter — the ArrayDeque, first element first.
	blockStack    []basicBlockEntry
	blockStackBuf [64]basicBlockEntry

	// Source of behavior at various points in the traversal.
	behavior ReferenceCollector_Behavior

	scopeCreator ScopeCreator

	// JavaScript compiler to use in traversing.
	compiler AbstractCompiler

	// Only collect references for filtered variables.
	varFilter func(*Var) bool

	callback            *ReferenceCollector_CollectorCallback
	callbackVal         ReferenceCollector_CollectorCallback
	referenceMapWrapper ReferenceCollector_ReferenceMapWrapper

	collectedHoistedFunctions *jutil.LinkedHashSet[*rhino.Node]
	hoistedTable              []*rhino.Node
	hoistedList               []*rhino.Node
	hoistedBuf                [16]*rhino.Node
	hoistedCount              uint8
	lastHoisted               *rhino.Node

	narrowScope *Scope

	refBuf    [8]Reference
	refPtrBuf [64]*Reference
	colBuf    [2]ReferenceCollection
	blockSlab []BasicBlock
	blockBuf  [8]BasicBlock
}

// NewReferenceCollector: constructor initializes block stack.
//
// Java: ReferenceCollector(AbstractCompiler, Behavior, ScopeCreator)
func NewReferenceCollector(compiler AbstractCompiler, behavior ReferenceCollector_Behavior, creator ScopeCreator) *ReferenceCollector {
	return NewReferenceCollectorWithVarFilter(compiler, behavior, creator, nil)
}

// NewReferenceCollectorWithVarFilter: constructor only collects references
// that match the given variable.
//
// The test for Var equality uses reference equality, so it's necessary to
// inject a scope when you traverse.
//
// Java: ReferenceCollector(AbstractCompiler, Behavior, ScopeCreator, Predicate<Var>)
func NewReferenceCollectorWithVarFilter(
	compiler AbstractCompiler,
	behavior ReferenceCollector_Behavior,
	creator ScopeCreator,
	varFilter func(*Var) bool) *ReferenceCollector {
	c := &ReferenceCollector{
		compiler:     compiler,
		behavior:     behavior,
		scopeCreator: creator,
		varFilter:    varFilter,
	}
	c.defaultSlabs.vars = c.varsBuf[:0]
	c.defaultSlabs.cols = c.colsBuf[:0]
	c.defaultSlabs.refSlab = c.refBuf[:0]
	c.defaultSlabs.refPtrSlab = c.refPtrBuf[:0]
	c.defaultSlabs.colSlab = c.colBuf[:0]
	c.globalSlabs = &c.defaultSlabs
	c.localSlabs = &c.defaultSlabs
	c.blockStack = c.blockStackBuf[:0]
	c.blockSlab = c.blockBuf[:0]
	c.callbackVal.c = c
	c.callback = &c.callbackVal
	c.referenceMapWrapper.c = c
	return c
}

func (c *ReferenceCollector) enableArenaPool() {
	c.poolable = true
	a := c.getArena()
	c.globalSlabs = &a.global
	c.localSlabs = &a.global
	if sc, ok := c.scopeCreator.(*SyntacticScopeCreator); ok && sc != nil {
		sc.EnableArenaPool()
	}
}

func (c *ReferenceCollector) enableLocalReclamation(canReclaim func() bool) {
	c.enableArenaPool()
	c.reclaimLocals = true
	c.canReclaimLocals = canReclaim
	c.localSlabs = &c.arena.local
}

func (c *ReferenceCollector) getArena() *referenceCollectorArena {
	if c.arena == nil {
		c.arena = acquireRCArena()
	}
	return c.arena
}

func (c *ReferenceCollector) slabsForVar(v *Var) *rcVarSlabs {
	if c.reclaimLocals && v.scope != nil && v.scope.parent != nil {
		return c.localSlabs
	}
	return c.globalSlabs
}

var rcDummyFnNode = rhino.NewNode(rhino.Token_FUNCTION)

func (c *ReferenceCollector) getDummyFnScope(t *NodeTraversal) *Scope {
	gs := t.scopes[0].untypedScope
	if gs == nil {
		gs = t.getScopeSlow(0)
	}
	if t.currentFunction != c.lastDummyFn {
		c.lastDummyFn = t.currentFunction
		c.dummyFnIdx ^= 1
	}
	s := &c.dummyFnScopes[c.dummyFnIdx]
	s.rootNode = rcDummyFnNode
	s.parent = gs
	s.lookupParent = gs
	s.depth = 1
	s.scanDone = true
	return s
}

func (c *ReferenceCollector) maybeReclaimLocals() {
	if c.narrowScope != nil || c.referenceMap != nil {
		return
	}
	if c.canReclaimLocals != nil && !c.canReclaimLocals() {
		c.localSlabs.markBase()
		if sc, ok := c.scopeCreator.(*SyntacticScopeCreator); ok && sc != nil && sc.arena != nil {
			sc.arena.markGlobalBase()
		}
		c.globalRefInNonFnLocal = false
		return
	}
	c.localSlabs.resetToBase(512, 128, 256, 65536)
	if c.lastVar != nil && c.lastVar.scope != nil && c.lastVar.scope.parent != nil {
		c.lastVar = nil
		c.lastCollection = nil
	}
	if c.lastRejectedVar != nil && c.lastRejectedVar.scope != nil && c.lastRejectedVar.scope.parent != nil {
		c.lastRejectedVar = nil
	}
	if len(c.argumentsRefs) > 0 {
		clear(c.argumentsRefs)
	}
	if c.refByKey != nil {
		clear(c.refByKey)
		c.refByKey = nil
	}
	if sc, ok := c.scopeCreator.(*SyntacticScopeCreator); ok && sc != nil {
		if !c.globalRefInNonFnLocal {
			sc.ResetLocalArena()
		} else if sc.arena != nil {
			sc.arena.markGlobalBase()
			c.globalRefInNonFnLocal = false
		}
	}
	if !c.globalRefInLocal && c.arena != nil && c.baseBlkChunkIdx > 0 {
		a := c.arena
		if a.blkChunkIdx > 0 {
			a.blkChunks[a.blkChunkIdx-1] = c.blockSlab
		}
		b := c.baseBlkChunkIdx - 1
		if a.blkChunkIdx-1 == b {
			clear(c.blockSlab[c.baseBlkSlabLen:])
			c.blockSlab = c.blockSlab[:c.baseBlkSlabLen]
		} else if a.blkChunkIdx-1 > b {
			chB := a.blkChunks[b]
			clear(chB[c.baseBlkSlabLen:])
			a.blkChunks[b] = chB[:c.baseBlkSlabLen]
			for i := b + 1; i < a.blkChunkIdx; i++ {
				ch := a.blkChunks[i]
				clear(ch)
				a.blkChunks[i] = ch[:0]
			}
			c.blockSlab = a.blkChunks[b]
			a.blkChunkIdx = c.baseBlkChunkIdx
		}
	}
}

func (c *ReferenceCollector) release() {
	if sc, ok := c.scopeCreator.(*SyntacticScopeCreator); ok && sc != nil {
		defer sc.ReleaseArena()
	}
	for _, kv := range c.defaultSlabs.vars {
		kv.nameHash = 0
	}
	a := c.arena
	if a == nil {
		return
	}
	c.arena = nil
	a.global.reset(4096, 512, 1024, 1048576)
	a.local.reset(512, 128, 256, 65536)
	c.globalSlabs = &c.defaultSlabs
	c.localSlabs = &c.defaultSlabs
	if len(c.hoistedTable) > 0 && len(c.hoistedList) > 0 {
		if len(c.hoistedTable) <= 524288 {
			if len(c.hoistedList)*16 < len(c.hoistedTable) {
				mask := uint32(len(c.hoistedTable) - 1)
				for _, hn := range c.hoistedList {
					slot := rcHashNodePtr(hn) & mask
					for c.hoistedTable[slot] != nil {
						c.hoistedTable[slot] = nil
						slot = (slot + 1) & mask
					}
				}
			} else {
				clear(c.hoistedTable)
			}
			a.hoistedTable = c.hoistedTable
		} else {
			a.hoistedTable = nil
		}
	}
	c.hoistedTable = nil
	if cap(c.hoistedList) > 0 {
		if cap(c.hoistedList) <= 262144 {
			clear(c.hoistedList)
			a.hoistedList = c.hoistedList[:0]
		} else {
			a.hoistedList = nil
		}
	}
	c.hoistedList = nil
	clear(c.hoistedBuf[:c.hoistedCount])
	c.hoistedCount = 0
	c.lastHoisted = nil
	if len(c.argumentsRefs) > 0 {
		if len(c.argumentsRefs) <= 16384 {
			clear(c.argumentsRefs)
		} else {
			a.argumentsRefs = nil
		}
	}
	c.argumentsRefs = nil
	const maxRetainedBlkChunks = 640
	if a.blkChunkIdx > 0 {
		a.blkChunks[a.blkChunkIdx-1] = c.blockSlab
	}
	for i := 0; i < a.blkChunkIdx && i < maxRetainedBlkChunks; i++ {
		ch := a.blkChunks[i]
		clear(ch)
		a.blkChunks[i] = ch[:0]
	}
	if len(a.blkChunks) > maxRetainedBlkChunks {
		clear(a.blkChunks[maxRetainedBlkChunks:])
		a.blkChunks = a.blkChunks[:maxRetainedBlkChunks]
	}
	a.blkChunkIdx = 0
	releaseRCArena(a)
}

func (c *ReferenceCollector) allocBasicBlock(parent *BasicBlock, root *rhino.Node) *BasicBlock {
	if len(c.blockSlab) == cap(c.blockSlab) {
		if c.poolable {
			a := c.getArena()
			if a.blkChunkIdx > 0 {
				a.blkChunks[a.blkChunkIdx-1] = c.blockSlab
			}
			if a.blkChunkIdx < len(a.blkChunks) {
				c.blockSlab = a.blkChunks[a.blkChunkIdx][:0]
				a.blkChunkIdx++
			} else {
				newCap := 64
				if cap(c.blockSlab) >= 64 {
					newCap = min(cap(c.blockSlab)*2, 1024) // 1024 * 24B = 24 KiB <= 32 KiB
				}
				c.blockSlab = make([]BasicBlock, 0, newCap)
				a.blkChunks = append(a.blkChunks, c.blockSlab)
				a.blkChunkIdx++
			}
		} else {
			newCap := 64
			if cap(c.blockSlab) >= 64 {
				newCap = min(cap(c.blockSlab)*2, 1024)
			}
			c.blockSlab = make([]BasicBlock, 0, newCap)
		}
	}
	idx := len(c.blockSlab)
	c.blockSlab = c.blockSlab[:idx+1]
	return initBasicBlock(&c.blockSlab[idx], parent, root)
}

func (c *ReferenceCollector) initHoistedTable() {
	if c.poolable {
		a := c.getArena()
		if len(a.hoistedTable) == 0 {
			a.hoistedTable = make([]*rhino.Node, 2048)
		}
		c.hoistedTable = a.hoistedTable
		if cap(a.hoistedList) > int(c.hoistedCount) {
			c.hoistedList = append(a.hoistedList[:0], c.hoistedBuf[:c.hoistedCount]...)
		} else {
			c.hoistedList = make([]*rhino.Node, c.hoistedCount, 1024)
			copy(c.hoistedList, c.hoistedBuf[:c.hoistedCount])
		}
	} else {
		c.hoistedTable = make([]*rhino.Node, 64)
		c.hoistedList = make([]*rhino.Node, c.hoistedCount, 64)
		copy(c.hoistedList, c.hoistedBuf[:c.hoistedCount])
	}
	table := c.hoistedTable
	mask := uint32(len(table) - 1)
	for _, hn := range c.hoistedList {
		slot := rcHashNodePtr(hn) & mask
		for table[slot] != nil {
			slot = (slot + 1) & mask
		}
		table[slot] = hn
	}
}

func (c *ReferenceCollector) growHoistedTable() {
	newCap := len(c.hoistedTable) * 2
	newTable := make([]*rhino.Node, newCap)
	mask := uint32(newCap - 1)
	for _, hn := range c.hoistedList {
		slot := rcHashNodePtr(hn) & mask
		for newTable[slot] != nil {
			slot = (slot + 1) & mask
		}
		newTable[slot] = hn
	}
	c.hoistedTable = newTable
	if c.arena != nil {
		c.arena.hoistedTable = newTable
	}
}

func (c *ReferenceCollector) addCollectedHoistedFunction(n *rhino.Node) bool {
	if n == c.lastHoisted && n != nil {
		return false
	}
	if c.hoistedTable == nil {
		for i := uint8(0); i < c.hoistedCount; i++ {
			if c.hoistedBuf[i] == n {
				c.lastHoisted = n
				return false
			}
		}
		if c.hoistedCount < uint8(len(c.hoistedBuf)) {
			c.hoistedBuf[c.hoistedCount] = n
			c.hoistedCount++
			c.lastHoisted = n
			return true
		}
		c.initHoistedTable()
	}
	table := c.hoistedTable
	mask := uint32(len(table) - 1)
	h := rcHashNodePtr(n)
	for slot := h & mask; ; slot = (slot + 1) & mask {
		v := table[slot]
		if v == n {
			c.lastHoisted = n
			return false
		}
		if v == nil {
			if (len(c.hoistedList)+1)*2 > len(table) {
				c.growHoistedTable()
				table = c.hoistedTable
				mask = uint32(len(table) - 1)
				slot = h & mask
				for table[slot] != nil {
					slot = (slot + 1) & mask
				}
			}
			table[slot] = n
			c.hoistedList = append(c.hoistedList, n)
			c.lastHoisted = n
			return true
		}
	}
}

func (c *ReferenceCollector) hasCollectedHoistedFunction(n *rhino.Node) bool {
	if n == c.lastHoisted {
		return n != nil
	}
	if c.hoistedTable == nil {
		for i := uint8(0); i < c.hoistedCount; i++ {
			if c.hoistedBuf[i] == n {
				c.lastHoisted = n
				return true
			}
		}
		return false
	}
	table := c.hoistedTable
	mask := uint32(len(table) - 1)
	for slot := rcHashNodePtr(n) & mask; ; slot = (slot + 1) & mask {
		v := table[slot]
		if v == n {
			c.lastHoisted = n
			return true
		}
		if v == nil {
			return false
		}
	}
}

func (c *ReferenceCollector) recordNewVarCol(v *Var, col *ReferenceCollection) {
	if v.nameID == implicitArgumentsID {
		if c.argumentsRefs == nil {
			if c.poolable {
				a := c.getArena()
				if a.argumentsRefs == nil {
					a.argumentsRefs = make(map[*rhino.Node]*ReferenceCollection, 16)
				}
				c.argumentsRefs = a.argumentsRefs
			} else {
				c.argumentsRefs = make(map[*rhino.Node]*ReferenceCollection, 4)
			}
		}
		c.argumentsRefs[v.scope.GetRootNode()] = col
	}
	if c.refByKey != nil {
		c.refByKey[scopedVarKey{scopeRoot: v.scope.GetRootNode(), nameID: v.nameID}] = col
	}
	if c.referenceMap != nil {
		c.referenceMap.Put(v, col)
	}
}

func (c *ReferenceCollector) getOrCreateReferences(v *Var) (*ReferenceCollection, *rcVarSlabs) {
	slabs := c.slabsForVar(v)
	if v == c.lastVar && v != nil {
		return c.lastCollection, slabs
	}
	if idx := int(v.nameHash) - 1; idx >= 0 && idx < len(slabs.vars) && slabs.vars[idx] == v {
		col := slabs.cols[idx]
		c.lastVar = v
		c.lastCollection = col
		return col, slabs
	}
	if !c.poolable && len(slabs.vars) <= len(c.varsBuf) {
		for i, kv := range slabs.vars {
			if kv == v || (kv.nameID == v.nameID && kv.scope.GetRootNode() == v.scope.GetRootNode()) {
				col := slabs.cols[i]
				v.nameHash = int32(i + 1)
				c.lastVar = v
				c.lastCollection = col
				return col, slabs
			}
		}
	}
	col := slabs.allocReferenceCollection()
	col.varScope = v.scope
	minGrow := 4096
	if c.reclaimLocals && slabs == c.localSlabs {
		minGrow = 256
	}
	slabs.addVarCol(v, col, minGrow)
	c.recordNewVarCol(v, col)
	c.lastVar = v
	c.lastCollection = col
	return col, slabs
}

func (c *ReferenceCollector) getReferencesSlow(v *Var) *ReferenceCollection {
	slabs := c.slabsForVar(v)
	if s := v.scope; s != nil {
		if s.implicitVars != nil {
			if imp := AbstractScope_ImplicitVar_OfID(v.nameID); imp != nil {
				if canon := s.implicitVars[*imp]; canon != nil && canon != v {
					if idx := int(canon.nameHash) - 1; idx >= 0 && idx < len(slabs.vars) && slabs.vars[idx] == canon {
						col := slabs.cols[idx]
						v.nameHash = canon.nameHash
						c.lastVar = v
						c.lastCollection = col
						return col
					}
				}
			}
		}
		if t := c.activeTraversal; t != nil {
			if d := int(s.depth); d < len(t.scopes) && t.scopes[d].scope == s {
				return nil
			}
		}
		if s == c.narrowScope {
			return nil
		}
	}
	return c.getReferencesByScopedNameID(v.scope.GetRootNode(), v.nameID)
}

func (c *ReferenceCollector) getReferencesByScopedNameID(scopeRoot *rhino.Node, nameID uint32) *ReferenceCollection {
	if nameID == implicitArgumentsID {
		if c.argumentsRefs != nil {
			return c.argumentsRefs[scopeRoot]
		}
		if c.referenceMap == nil {
			return nil
		}
	}
	gVars := c.globalSlabs.vars
	gCols := c.globalSlabs.cols
	var lVars []*Var
	var lCols []*ReferenceCollection
	if c.reclaimLocals {
		lVars = c.localSlabs.vars
		lCols = c.localSlabs.cols
	}
	if len(gVars)+len(lVars) <= len(c.varsBuf) {
		for i, kv := range gVars {
			if kv.nameID == nameID && kv.scope.GetRootNode() == scopeRoot {
				return gCols[i]
			}
		}
		for i, kv := range lVars {
			if kv.nameID == nameID && kv.scope.GetRootNode() == scopeRoot {
				return lCols[i]
			}
		}
	} else {
		if c.refByKey == nil {
			c.refByKey = make(map[scopedVarKey]*ReferenceCollection, len(gVars)+len(lVars))
			for i, kv := range gVars {
				c.refByKey[scopedVarKey{scopeRoot: kv.scope.GetRootNode(), nameID: kv.nameID}] = gCols[i]
			}
			for i, kv := range lVars {
				c.refByKey[scopedVarKey{scopeRoot: kv.scope.GetRootNode(), nameID: kv.nameID}] = lCols[i]
			}
		}
		if col := c.refByKey[scopedVarKey{scopeRoot: scopeRoot, nameID: nameID}]; col != nil {
			return col
		}
	}
	if c.referenceMap != nil {
		for k, col := range c.referenceMap.All() {
			if k != nil && k.nameID == nameID && k.scope.GetRootNode() == scopeRoot {
				return col
			}
		}
	}
	return nil
}

func (c *ReferenceCollector) ensureRawReferenceMap() *jutil.ObjectLinkedHashMap[*Var, *ReferenceCollection] {
	if c.referenceMap == nil {
		c.referenceMap = jutil.NewObjectLinkedHashMap[*Var, *ReferenceCollection]()
		for i, v := range c.globalSlabs.vars {
			v.nameHash = 0
			c.referenceMap.Put(v, c.globalSlabs.cols[i])
		}
		if c.reclaimLocals {
			for i, v := range c.localSlabs.vars {
				v.nameHash = 0
				c.referenceMap.Put(v, c.localSlabs.cols[i])
			}
		}
	}
	return c.referenceMap
}

func (c *ReferenceCollector) builderVal() NodeTraversal_Builder {
	return NodeTraversal_Builder{
		compiler:                                       c.compiler,
		callback:                                       c.callback,
		scopeCreator:                                   c.scopeCreator,
		obeyDestructuringAndDefaultValueExecutionOrder: true,
	}
}

// Process is a convenience method for running this pass over a tree with
// this class as a callback.
//
// Java: ReferenceCollector.process(Node, Node)
func (c *ReferenceCollector) Process(externs *rhino.Node, root *rhino.Node) {
	b := c.builderVal()
	b.TraverseRoots(externs, root)
}

// ProcessRoot is the one-argument process overload.
//
// Go-only: language-adapter — the name: CompilerPass fixes the plain name
// to the two-argument overload.
//
// Java: ReferenceCollector.process(Node)
func (c *ReferenceCollector) ProcessRoot(root *rhino.Node) {
	b := c.builderVal()
	b.Traverse(root)
}

// ProcessScope targets reference collection to a particular scope.
//
// Java: ReferenceCollector.processScope(Scope)
func (c *ReferenceCollector) ProcessScope(scope *Scope) {
	shouldAddToBlockStack := !scope.IsHoistScope()
	c.narrowScope = scope
	if shouldAddToBlockStack {
		c.PushNewBlock(scope.GetRootNode())
	}
	b := c.builderVal()
	b.TraverseAtScope(scope)
	if shouldAddToBlockStack {
		c.PopLastBlock(scope.GetRootNode())
	}
	c.narrowScope = nil
}

// CreateTraversalBuilder: same as process but only runs on a part of AST
// associated to one script.
//
// Java: ReferenceCollector.createTraversalBuilder()
func (c *ReferenceCollector) CreateTraversalBuilder() *NodeTraversal_Builder {
	return NodeTraversal_Builder_().
		SetCompiler(c.compiler).
		SetCallback(c.callback).
		SetScopeCreator(c.scopeCreator).
		SetObeyDestructuringAndDefaultValueExecutionOrder(true)
}

// GetAllSymbols gets the variables that were referenced in this callback.
//
// Java: ReferenceCollector.getAllSymbols()
func (c *ReferenceCollector) GetAllSymbols() iter.Seq[*Var] {
	if c.referenceMap != nil {
		return c.referenceMap.Keys()
	}
	if !c.reclaimLocals || len(c.localSlabs.vars) == 0 {
		return slices.Values(c.globalSlabs.vars)
	}
	return func(yield func(*Var) bool) {
		for _, v := range c.globalSlabs.vars {
			if !yield(v) {
				return
			}
		}
		for _, v := range c.localSlabs.vars {
			if !yield(v) {
				return
			}
		}
	}
}

// GetScope narrows StaticSymbolTable.getScope; see ReferenceCollector.
//
// Java: ReferenceCollector.getScope(Var)
func (c *ReferenceCollector) GetScope(var_ *Var) *Scope {
	return var_.getScopeS()
}

// GetReferences gets the reference collection for the given variable. It
// narrows StaticSymbolTable.getReferences; see ReferenceCollector.
//
// Java: ReferenceCollector.getReferences(Var)
func (c *ReferenceCollector) GetReferences(v *Var) *ReferenceCollection {
	if v == nil {
		if c.referenceMap != nil {
			return c.referenceMap.Get(nil)
		}
		return nil
	}
	if v == c.lastVar {
		return c.lastCollection
	}
	slabs := c.slabsForVar(v)
	if idx := int(v.nameHash) - 1; idx >= 0 && idx < len(slabs.vars) && slabs.vars[idx] == v {
		col := slabs.cols[idx]
		c.lastVar = v
		c.lastCollection = col
		return col
	}
	return c.getReferencesSlow(v)
}

// AsStaticSymbolTable returns this collector as the Java
// StaticSymbolTable<Var, Reference> it implements (for
// SymbolTable.addSymbolsFrom).
//
// Go-only: language-adapter — see ReferenceCollector.
func (c *ReferenceCollector) AsStaticSymbolTable() rhino.StaticSymbolTable[*Var, *Reference] {
	return referenceCollector_staticSymbolTable{c}
}

// referenceCollector_staticSymbolTable is the
// rhino.StaticSymbolTable[*Var, *Reference] view of a ReferenceCollector:
// the same three methods, with StaticSymbolTable's Go result types.
//
// Go-only: language-adapter — see ReferenceCollector.
type referenceCollector_staticSymbolTable struct {
	c *ReferenceCollector
}

// Java: ReferenceCollector.getReferences(Var)
func (t referenceCollector_staticSymbolTable) GetReferences(symbol *Var) iter.Seq[*Reference] {
	refs := t.c.GetReferences(symbol)
	if refs == nil {
		return nil
	}
	return refs.All()
}

// Java: ReferenceCollector.getScope(Var)
func (t referenceCollector_staticSymbolTable) GetScope(symbol *Var) rhino.StaticScope {
	if s := t.c.GetScope(symbol); s != nil {
		return s
	}
	return nil
}

// Java: ReferenceCollector.getAllSymbols()
func (t referenceCollector_staticSymbolTable) GetAllSymbols() iter.Seq[*Var] {
	return t.c.GetAllSymbols()
}

func rcIsHoistedFunctionDeclaration(fnNode *rhino.Node) bool {
	parent := fnNode.GetParent()
	if parent == nil {
		return false
	}
	switch parent.GetToken() {
	case rhino.Token_SCRIPT, rhino.Token_MODULE_BODY, rhino.Token_EXPORT:
		first := fnNode.GetFirstChild()
		return first != nil && first.IsName() && first.GetString() != ""
	case rhino.Token_BLOCK, rhino.Token_LABEL, rhino.Token_DECLARE, rhino.Token_NAMESPACE_ELEMENTS:
		if gp := parent.GetParent(); gp != nil && gp.IsFunction() {
			first := fnNode.GetFirstChild()
			return first != nil && first.IsName() && first.GetString() != ""
		}
	}
	return false
}

// Java: ReferenceCollector.maybeJumpToHoistedFunction(Var, NodeTraversal)
func (c *ReferenceCollector) MaybeJumpToHoistedFunction(var_ *Var, t *NodeTraversal) {
	nameNode := var_.nameNode
	if nameNode == nil {
		return
	}
	fnNode := nameNode.GetParent()
	if fnNode == nil || !fnNode.IsFunction() || fnNode == c.lastHoisted {
		return
	}
	c.maybeJumpToHoistedFunctionSlow(var_, fnNode, t)
}

func (c *ReferenceCollector) maybeJumpToHoistedFunctionSlow(var_ *Var, fnNode *rhino.Node, t *NodeTraversal) {
	if c.hasCollectedHoistedFunction(fnNode) || !rcIsHoistedFunctionDeclaration(fnNode) {
		return
	}
	varScope := var_.getScopeS()

	// If we're only traversing a narrow scope, do not try to climb outside.
	if c.narrowScope != nil && c.narrowScope.GetDepth() > varScope.GetDepth() {
		return
	}

	/*
	 * Replace the block stack with a new one matching the hoist position.
	 *
	 * <p>This algorithm only works because we know hoisted functions cannot be inside loops. It
	 * will have to change if we ever do general function continuations.
	 *
	 * <p>This is tricky to compute because of the weird traverseAtScope call for
	 * AggressiveInlineAliases.
	 */
	oldBlockStack := c.blockStack
	if varScope.IsGlobal() {
		if len(oldBlockStack) == 0 {
			panic(jutil.NewNoSuchElementException())
		}
		if oldBlockStack[0].block == nil {
			c.materializeBlock(0)
		}
		c.blockStack = append(oldBlockStack[len(oldBlockStack):len(oldBlockStack)], oldBlockStack[0])
	} else {
		rootNode := varScope.GetRootNode()
		end := len(oldBlockStack)
		for i := range oldBlockStack {
			if oldBlockStack[i].root == rootNode {
				end = i + 1
				break
			}
		}
		if end > 0 && oldBlockStack[end-1].block == nil {
			c.materializeBlock(end - 1)
		}
		c.blockStack = append(oldBlockStack[len(oldBlockStack):len(oldBlockStack)], oldBlockStack[:end]...)
	}

	/*
	 * Record the function declaration reference explicitly because it is a bleeding name. The
	 * reference must be recorded with the containing scope, and traverseAtScope skips bleeding
	 * function NAME nodes
	 */
	c.AddReference(var_, fnNode.GetFirstChild(), t)
	b := c.builderVal()
	b.TraverseAtScope(c.scopeCreator.CreateScope(fnNode, varScope))

	clear(c.blockStack)
	c.blockStack = oldBlockStack
}

// referenceCollector_getFirst is ArrayDeque.getFirst().
//
// Go-only: language-adapter — see ReferenceCollector.blockStack.
func referenceCollector_getFirst(stack []*BasicBlock) *BasicBlock {
	if len(stack) == 0 {
		panic(jutil.NewNoSuchElementException())
	}
	return stack[0]
}

func (c *ReferenceCollector) getCurrentBlock() *BasicBlock {
	idx := len(c.blockStack) - 1
	if b := c.blockStack[idx].block; b != nil {
		return b
	}
	return c.materializeBlock(idx)
}

func (c *ReferenceCollector) materializeBlock(idx int) *BasicBlock {
	i := idx
	for i > 0 && c.blockStack[i-1].block == nil {
		i--
	}
	var parent *BasicBlock
	if i > 0 {
		parent = c.blockStack[i-1].block
	}
	for ; i <= idx; i++ {
		parent = c.allocBasicBlock(parent, c.blockStack[i].root)
		c.blockStack[i].block = parent
	}
	return parent
}

// ReferenceCollector_CollectorCallback is the traversal callback.
type ReferenceCollector_CollectorCallback struct {
	c *ReferenceCollector
}

// ShouldTraverse updates block stack.
//
// Java: ReferenceCollector.CollectorCallback.shouldTraverse(NodeTraversal, Node, Node)
func (cb *ReferenceCollector_CollectorCallback) ShouldTraverse(nodeTraversal *NodeTraversal, n *rhino.Node, parent *rhino.Node) bool {
	if n.GetFirstChild() == nil {
		switch n.GetToken() {
		case rhino.Token_NAME, rhino.Token_IMPORT_STAR, rhino.Token_BLOCK, rhino.Token_SWITCH_BODY:
			if referenceCollector_isBlockBoundary(n, parent) {
				cb.c.PushNewBlock(n)
			}
			return true
		default:
			return false
		}
	}

	// We automatically traverse a hoisted function body when that function
	// is first referenced, so that the reference lists are in the right order.
	//
	// TODO(nicksantos): Maybe generalize this to a continuation mechanism
	// like in RemoveUnusedCode.
	if n.IsFunction() && rcIsHoistedFunctionDeclaration(n) {
		if !cb.c.addCollectedHoistedFunction(n) {
			return false
		}
	}

	// If node is a new basic block, put on basic block stack
	if referenceCollector_isBlockBoundary(n, parent) {
		cb.c.PushNewBlock(n)
	}

	return true
}

// EnterScope updates block stack and invokes any additional behavior.
//
// Java: ReferenceCollector.CollectorCallback.enterScope(NodeTraversal)
func (cb *ReferenceCollector_CollectorCallback) EnterScope(t *NodeTraversal) {
	if !t.InGlobalScope() {
		if cb.c.localScopeDepth == 0 && cb.c.reclaimLocals {
			cb.c.globalRefInLocal = false
			if cb.c.arena != nil {
				cb.c.baseBlkChunkIdx = cb.c.arena.blkChunkIdx
				cb.c.baseBlkSlabLen = len(cb.c.blockSlab)
			}
		}
		cb.c.localScopeDepth++
	}
	// Don't add all ES6 scope roots to blockStack, only those that are also scopes according to
	// the ES5 scoping rules. Other nodes that ought to be considered the root of a BasicBlock
	// are added in shouldTraverse() or processScope() and removed in visit().
	if root := t.GetScopeRoot(); nodeTraversal_isHoistScopeRootNode(root) {
		cb.c.PushNewBlock(root)
	}
}

// ExitScope updates block stack and invokes any additional behavior.
//
// Java: ReferenceCollector.CollectorCallback.exitScope(NodeTraversal)
func (cb *ReferenceCollector_CollectorCallback) ExitScope(t *NodeTraversal) {
	if root := t.GetScopeRoot(); nodeTraversal_isHoistScopeRootNode(root) {
		cb.c.PopLastBlock(root)
	}
	prevTraversal := cb.c.activeTraversal
	cb.c.activeTraversal = t
	cb.c.behavior.AfterExitScope(t, &cb.c.referenceMapWrapper)
	cb.c.activeTraversal = prevTraversal
	if !t.InGlobalScope() {
		cb.c.localScopeDepth--
		if cb.c.localScopeDepth == 0 && cb.c.reclaimLocals {
			cb.c.maybeReclaimLocals()
		}
	}
}

// Visit: for each node, update the block stack and reference collection as
// appropriate.
//
// Java: ReferenceCollector.CollectorCallback.visit(NodeTraversal, Node, Node)
func (cb *ReferenceCollector_CollectorCallback) Visit(t *NodeTraversal, n *rhino.Node, parent *rhino.Node) {
	if tok := n.GetToken(); tok == rhino.Token_NAME || tok == rhino.Token_IMPORT_STAR {
		if parent != nil {
			if pt := parent.GetToken(); (pt == rhino.Token_IMPORT_SPEC && n != parent.GetLastChild()) ||
				(pt == rhino.Token_EXPORT_SPEC && n != parent.GetFirstChild()) {
				// The n in `import {n as x}` or `export {x as n}` are not references, even though
				// they are represented in the AST as NAME nodes.
				return
			}
		}

		id := n.GetStringID()
		var v *Var
		var scope *Scope
		if id == implicitArgumentsID {
			scope = t.GetScope()
			v = scope.GetVarByID(id)
		} else {
			v = t.LookupVarByID(id)
		}
		if v != nil {
			cb.c.addReferenceWithScope(v, n, t, scope)
			if v.nameNode != n {
				cb.c.MaybeJumpToHoistedFunction(v, t)
			}
		}
	}

	if referenceCollector_isBlockBoundary(n, parent) {
		cb.c.PopLastBlock(n)
	}
}

// Java: ReferenceCollector.pushNewBlock(Node)
func (c *ReferenceCollector) PushNewBlock(root *rhino.Node) {
	c.blockStack = append(c.blockStack, basicBlockEntry{root: root})
}

// Java: ReferenceCollector.popLastBlock(Node)
func (c *ReferenceCollector) PopLastBlock(root *rhino.Node) {
	if len(c.blockStack) == 0 {
		panic(jutil.NewNoSuchElementException()) // ArrayDeque.removeLast()
	}
	idx := len(c.blockStack) - 1
	lastRoot := c.blockStack[idx].root
	c.blockStack[idx] = basicBlockEntry{}
	c.blockStack = c.blockStack[:idx]
	// Verfiy that the stack is unwound correctly.
	if lastRoot != root {
		jutil.CheckStateMsg1(base.JSCompObjects_Identical(root, lastRoot), "%s", root)
	}
}

var refCollectorParentBoundaryKind = [256]uint8{
	rhino.Token_DO:               1,
	rhino.Token_FOR:              1,
	rhino.Token_FOR_IN:           1,
	rhino.Token_FOR_OF:           1,
	rhino.Token_FOR_AWAIT_OF:     1,
	rhino.Token_TRY:              1,
	rhino.Token_WHILE:            1,
	rhino.Token_WITH:             1,
	rhino.Token_CLASS:            1,
	rhino.Token_SWITCH_BODY:      1,
	rhino.Token_AND:              2,
	rhino.Token_HOOK:             2,
	rhino.Token_IF:               2,
	rhino.Token_OR:               2,
	rhino.Token_COALESCE:         2,
	rhino.Token_OPTCHAIN_GETPROP: 2,
	rhino.Token_OPTCHAIN_GETELEM: 2,
	rhino.Token_OPTCHAIN_CALL:    2,
	rhino.Token_DEFAULT_VALUE:    2,
}

// referenceCollector_isBlockBoundary returns true if this node marks the
// start of a new basic block
//
// Java: ReferenceCollector.isBlockBoundary(Node, Node)
func referenceCollector_isBlockBoundary(n *rhino.Node, parent *rhino.Node) bool {
	if parent != nil {
		switch refCollectorParentBoundaryKind[uint8(parent.GetToken())] {
		case 1:
			return true
		case 2:
			return n != parent.GetFirstChild()
		}
	}

	return n.IsCase()
}

// Java: ReferenceCollector.addReference(Var, Node, NodeTraversal)
func (c *ReferenceCollector) AddReference(v *Var, n *rhino.Node, t *NodeTraversal) {
	c.addReferenceWithScope(v, n, t, nil)
}

func (c *ReferenceCollector) addReferenceWithScope(v *Var, n *rhino.Node, t *NodeTraversal, scope *Scope) {
	var collection *ReferenceCollection
	var slabs *rcVarSlabs
	if v == c.lastVar && v != nil {
		collection = c.lastCollection
		slabs = c.slabsForVar(v)
	} else {
		if c.varFilter != nil && (v == c.lastRejectedVar || !c.varFilter(v)) {
			c.lastRejectedVar = v
			return
		}
		collection, slabs = c.getOrCreateReferences(v)
	}

	// Add this particular reference
	if len(c.blockStack) == 0 {
		panic(jutil.NewNoSuchElementException()) // ArrayDeque.getLast()
	}
	if slabs == c.globalSlabs && c.localScopeDepth > 0 {
		c.globalRefInLocal = true
		if c.reclaimLocals && scope == nil && t.currentFunction != nil && !t.GetScopeRoot().IsMemberFieldDef() {
			scope = c.getDummyFnScope(t)
		} else {
			c.globalRefInNonFnLocal = true
		}
	}
	if scope == nil {
		scope = t.GetScope()
	}
	inputId := t.inputId
	if t.compilerInput == nil {
		inputId = t.GetInput().GetInputId()
	}
	ref := slabs.allocReference(n, scope, c.getCurrentBlock(), inputId)
	s := collection.references.Slice()
	if len(s) == cap(s) {
		slabs.growReferenceCollection(collection, s)
	}
	collection.references.Add(ref)
}

// ReferenceCollector_ReferenceMapWrapper is the ReferenceMap handed to
// Behavior.afterExitScope.
type ReferenceCollector_ReferenceMapWrapper struct {
	c            *ReferenceCollector
	referenceMap *jutil.ObjectLinkedHashMap[*Var, *ReferenceCollection]
}

// Java: ReferenceCollector.ReferenceMapWrapper(Map<Var, ReferenceCollection>)
func NewReferenceCollector_ReferenceMapWrapper(
	referenceMap *jutil.ObjectLinkedHashMap[*Var, *ReferenceCollection]) *ReferenceCollector_ReferenceMapWrapper {
	return &ReferenceCollector_ReferenceMapWrapper{referenceMap: referenceMap}
}

// Java: ReferenceCollector.ReferenceMapWrapper.getReferences(Var)
func (w *ReferenceCollector_ReferenceMapWrapper) GetReferences(var_ *Var) *ReferenceCollection {
	if w.c != nil {
		return w.c.GetReferences(var_)
	}
	return w.referenceMap.Get(var_)
}

func (w *ReferenceCollector_ReferenceMapWrapper) getReferencesByScopedName(scopeRoot *rhino.Node, name string) *ReferenceCollection {
	if w.c != nil {
		id := rhino.RhinoStringPool_LookupID(name)
		if id == 0 {
			return nil
		}
		return w.c.getReferencesByScopedNameID(scopeRoot, id)
	}
	return nil
}

// Java: ReferenceCollector.ReferenceMapWrapper.getRawReferenceMap()
func (w *ReferenceCollector_ReferenceMapWrapper) GetRawReferenceMap() *jutil.ObjectLinkedHashMap[*Var, *ReferenceCollection] {
	if w.c != nil {
		return w.c.ensureRawReferenceMap()
	}
	return w.referenceMap
}

// Java: ReferenceCollector.ReferenceMapWrapper.toString()
func (w *ReferenceCollector_ReferenceMapWrapper) String() string {
	return w.GetRawReferenceMap().String()
}

// ReferenceCollector_Behavior is the way for callers to add specific
// behavior during traversal that utilizes the built-up reference
// information.
type ReferenceCollector_Behavior interface {
	// AfterExitScope is called after we finish with a scope.
	//
	// Java: ReferenceCollector.Behavior.afterExitScope(NodeTraversal, ReferenceMap)
	AfterExitScope(t *NodeTraversal, referenceMap ReferenceMap)
}

// ReferenceCollector_BehaviorFunc is a Behavior given as a function (a Java
// lambda for Behavior).
//
// Go-only: language-adapter — the adapter for lambdas implementing the
// single-method interface Behavior (PORTING.md §4).
type ReferenceCollector_BehaviorFunc func(t *NodeTraversal, referenceMap ReferenceMap)

// Java: ReferenceCollector.Behavior.afterExitScope(NodeTraversal, ReferenceMap)
func (f ReferenceCollector_BehaviorFunc) AfterExitScope(t *NodeTraversal, referenceMap ReferenceMap) {
	f(t, referenceMap)
}

var ReferenceCollector_DO_NOTHING_BEHAVIOR ReferenceCollector_Behavior = &referenceCollector_doNothingBehavior{}

// referenceCollector_doNothingBehavior is the anonymous class of
// DO_NOTHING_BEHAVIOR.
type referenceCollector_doNothingBehavior struct {
	_ byte // Go-only: language-adapter — distinct identity
}

// Java: ReferenceCollector.DO_NOTHING_BEHAVIOR.afterExitScope(NodeTraversal, ReferenceMap)
func (*referenceCollector_doNothingBehavior) AfterExitScope(t *NodeTraversal, referenceMap ReferenceMap) {
}
