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

type scopedVarKey struct {
	scopeRoot *rhino.Node
	name      string
}

type referenceCollectorArena struct {
	vars           []*Var
	cols           []*ReferenceCollection
	refChunks      [][]Reference
	refChunkIdx    int
	refPtrChunks   [][]*Reference
	refPtrChunkIdx int
	colChunks      [][]ReferenceCollection
	colChunkIdx    int
	blkChunks      [][]BasicBlock
	blkChunkIdx    int
	hoistedTable   []*rhino.Node
	hoistedList    []*rhino.Node
	argumentsRefs  map[*rhino.Node]*ReferenceCollection
}

var cachedRCArenas [2]atomic.Pointer[referenceCollectorArena]

var referenceCollectorArenaPool = sync.Pool{
	New: func() any {
		return &referenceCollectorArena{}
	},
}

func acquireRCArena() *referenceCollectorArena {
	for i := range cachedRCArenas {
		if a := cachedRCArenas[i].Swap(nil); a != nil {
			return a
		}
	}
	return referenceCollectorArenaPool.Get().(*referenceCollectorArena)
}

func releaseRCArena(a *referenceCollectorArena) {
	for i := range cachedRCArenas {
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

	vars            []*Var
	cols            []*ReferenceCollection
	varsBuf         [16]*Var
	colsBuf         [16]*ReferenceCollection
	refByKey        map[scopedVarKey]*ReferenceCollection
	argumentsRefs   map[*rhino.Node]*ReferenceCollection
	lastVar         *Var
	lastCollection  *ReferenceCollection
	lastRejectedVar *Var
	activeTraversal *NodeTraversal
	arena           *referenceCollectorArena
	poolable        bool

	// The stack of basic blocks and scopes the current traversal is in.
	//
	// Go-only: language-adapter — the ArrayDeque, first element first.
	blockStack    []*BasicBlock
	blockStackBuf [64]*BasicBlock

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

	refSlab    []Reference
	refBuf     [8]Reference
	refPtrSlab []*Reference
	refPtrBuf  [64]*Reference
	colSlab    []ReferenceCollection
	colBuf     [2]ReferenceCollection
	blockSlab  []BasicBlock
	blockBuf   [8]BasicBlock
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
	c.vars = c.varsBuf[:0]
	c.cols = c.colsBuf[:0]
	c.blockStack = c.blockStackBuf[:0]
	c.refSlab = c.refBuf[:0]
	c.refPtrSlab = c.refPtrBuf[:0]
	c.colSlab = c.colBuf[:0]
	c.blockSlab = c.blockBuf[:0]
	c.callbackVal.c = c
	c.callback = &c.callbackVal
	c.referenceMapWrapper.c = c
	return c
}

func (c *ReferenceCollector) enableArenaPool() {
	c.poolable = true
	if sc, ok := c.scopeCreator.(*SyntacticScopeCreator); ok && sc != nil {
		sc.EnableArenaPool()
	}
}

func (c *ReferenceCollector) getArena() *referenceCollectorArena {
	if c.arena == nil {
		c.arena = acquireRCArena()
	}
	return c.arena
}

func (c *ReferenceCollector) release() {
	if sc, ok := c.scopeCreator.(*SyntacticScopeCreator); ok && sc != nil {
		defer sc.ReleaseArena()
	}
	for _, kv := range c.vars {
		kv.nameHash = 0
	}
	a := c.arena
	if a == nil {
		return
	}
	c.arena = nil
	if cap(c.vars) > len(c.varsBuf) {
		if cap(c.vars) <= 1048576 {
			clear(c.vars)
			a.vars = c.vars[:0]
			clear(c.cols)
			a.cols = c.cols[:0]
		} else {
			a.vars = nil
			a.cols = nil
		}
	}
	c.vars = nil
	c.cols = nil
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
	const (
		maxRetainedRefChunks    = 256
		maxRetainedRefPtrChunks = 256
		maxRetainedColChunks    = 128
		maxRetainedBlkChunks    = 128
	)
	if a.refChunkIdx > 0 {
		a.refChunks[a.refChunkIdx-1] = c.refSlab
	}
	for i := 0; i < a.refChunkIdx && i < maxRetainedRefChunks; i++ {
		ch := a.refChunks[i]
		clear(ch)
		a.refChunks[i] = ch[:0]
	}
	if len(a.refChunks) > maxRetainedRefChunks {
		clear(a.refChunks[maxRetainedRefChunks:])
		a.refChunks = a.refChunks[:maxRetainedRefChunks]
	}
	if a.refPtrChunkIdx > 0 {
		a.refPtrChunks[a.refPtrChunkIdx-1] = c.refPtrSlab
	}
	for i := 0; i < a.refPtrChunkIdx && i < maxRetainedRefPtrChunks; i++ {
		ch := a.refPtrChunks[i]
		clear(ch)
		a.refPtrChunks[i] = ch[:0]
	}
	if len(a.refPtrChunks) > maxRetainedRefPtrChunks {
		clear(a.refPtrChunks[maxRetainedRefPtrChunks:])
		a.refPtrChunks = a.refPtrChunks[:maxRetainedRefPtrChunks]
	}
	if a.colChunkIdx > 0 {
		a.colChunks[a.colChunkIdx-1] = c.colSlab
	}
	for i := 0; i < a.colChunkIdx && i < maxRetainedColChunks; i++ {
		ch := a.colChunks[i]
		for j := range ch {
			col := &ch[j]
			if cap(col.refsList.Slice()) == len(col.refsBuf) {
				col.refsList.Clear()
			} else {
				col.refsList.InitWithBuffer(col.refsBuf[:])
			}
		}
		a.colChunks[i] = ch[:0]
	}
	if len(a.colChunks) > maxRetainedColChunks {
		clear(a.colChunks[maxRetainedColChunks:])
		a.colChunks = a.colChunks[:maxRetainedColChunks]
	}
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
	a.refChunkIdx = 0
	a.refPtrChunkIdx = 0
	a.colChunkIdx = 0
	a.blkChunkIdx = 0
	releaseRCArena(a)
}

func (c *ReferenceCollector) allocReference(n *rhino.Node, scope *Scope, basicBlock *BasicBlock, inputId *rhino.InputId) *Reference {
	if len(c.refSlab) == cap(c.refSlab) {
		if c.poolable {
			a := c.getArena()
			if a.refChunkIdx > 0 {
				a.refChunks[a.refChunkIdx-1] = c.refSlab
			}
			if a.refChunkIdx < len(a.refChunks) {
				c.refSlab = a.refChunks[a.refChunkIdx][:0]
				a.refChunkIdx++
			} else {
				newCap := 256
				if cap(c.refSlab) >= 256 {
					newCap = min(cap(c.refSlab)*2, 8192)
				}
				c.refSlab = make([]Reference, 0, newCap)
				a.refChunks = append(a.refChunks, c.refSlab)
				a.refChunkIdx++
			}
		} else {
			newCap := 256
			if cap(c.refSlab) >= 256 {
				newCap = min(cap(c.refSlab)*2, 8192)
			}
			c.refSlab = make([]Reference, 0, newCap)
		}
	}
	idx := len(c.refSlab)
	c.refSlab = c.refSlab[:idx+1]
	r := &c.refSlab[idx]
	r.nameNode = n
	r.basicBlock = basicBlock
	r.scope = scope
	r.inputId = inputId
	return r
}

func (c *ReferenceCollector) allocReferenceCollection() *ReferenceCollection {
	if len(c.colSlab) == cap(c.colSlab) {
		if c.poolable {
			a := c.getArena()
			if a.colChunkIdx > 0 {
				a.colChunks[a.colChunkIdx-1] = c.colSlab
			}
			if a.colChunkIdx < len(a.colChunks) {
				c.colSlab = a.colChunks[a.colChunkIdx][:0]
				a.colChunkIdx++
			} else {
				newCap := 64
				if cap(c.colSlab) >= 64 {
					newCap = min(cap(c.colSlab)*2, 4096)
				}
				c.colSlab = make([]ReferenceCollection, 0, newCap)
				a.colChunks = append(a.colChunks, c.colSlab)
				a.colChunkIdx++
			}
		} else {
			newCap := 64
			if cap(c.colSlab) >= 64 {
				newCap = min(cap(c.colSlab)*2, 4096)
			}
			c.colSlab = make([]ReferenceCollection, 0, newCap)
		}
	}
	idx := len(c.colSlab)
	c.colSlab = c.colSlab[:idx+1]
	col := &c.colSlab[idx]
	if col.references != nil {
		return col
	}
	return initReferenceCollection(col)
}

func (c *ReferenceCollector) allocRefPtrSlice(n int) []*Reference {
	if len(c.refPtrSlab)+n > cap(c.refPtrSlab) {
		if c.poolable {
			a := c.getArena()
			if a.refPtrChunkIdx > 0 {
				a.refPtrChunks[a.refPtrChunkIdx-1] = c.refPtrSlab
			}
			if a.refPtrChunkIdx < len(a.refPtrChunks) {
				ch := a.refPtrChunks[a.refPtrChunkIdx]
				if cap(ch) < n {
					ch = make([]*Reference, 0, max(8192, n))
					a.refPtrChunks[a.refPtrChunkIdx] = ch
				}
				c.refPtrSlab = ch[:0]
				a.refPtrChunkIdx++
			} else {
				c.refPtrSlab = make([]*Reference, 0, max(8192, n))
				a.refPtrChunks = append(a.refPtrChunks, c.refPtrSlab)
				a.refPtrChunkIdx++
			}
		} else {
			newCap := 256
			if cap(c.refPtrSlab) >= 256 {
				newCap = min(cap(c.refPtrSlab)*2, 4096)
			}
			if newCap < n {
				newCap = n
			}
			c.refPtrSlab = make([]*Reference, 0, newCap)
		}
	}
	idx := len(c.refPtrSlab)
	c.refPtrSlab = c.refPtrSlab[:idx+n]
	return c.refPtrSlab[idx : idx+n : idx+n]
}

func (c *ReferenceCollector) growReferenceCollection(col *ReferenceCollection, s []*Reference) {
	oldCap := cap(s)
	newCap := 16
	if oldCap >= 16 {
		newCap = oldCap * 2
	}
	newBuf := c.allocRefPtrSlice(newCap)
	col.refsList.InitWithBuffer(newBuf)
	col.refsList.AddAll(s...)
	if oldCap == len(col.refsBuf) {
		clear(s)
	}
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
					newCap = min(cap(c.blockSlab)*2, 4096)
				}
				c.blockSlab = make([]BasicBlock, 0, newCap)
				a.blkChunks = append(a.blkChunks, c.blockSlab)
				a.blkChunkIdx++
			}
		} else {
			newCap := 64
			if cap(c.blockSlab) >= 64 {
				newCap = min(cap(c.blockSlab)*2, 4096)
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

func (c *ReferenceCollector) growVarsCols() {
	if c.poolable && cap(c.vars) == len(c.varsBuf) {
		a := c.getArena()
		if cap(a.vars) > len(c.vars) {
			c.vars = append(a.vars[:0], c.vars...)
			c.cols = append(a.cols[:0], c.cols...)
			return
		}
	}
	newCap := max(cap(c.vars)*2, 4096)
	newVars := make([]*Var, len(c.vars), newCap)
	copy(newVars, c.vars)
	c.vars = newVars
	newCols := make([]*ReferenceCollection, len(c.cols), newCap)
	copy(newCols, c.cols)
	c.cols = newCols
	if c.arena != nil {
		c.arena.vars = newVars[:0]
		c.arena.cols = newCols[:0]
	}
}

func (c *ReferenceCollector) recordNewVarCol(v *Var, col *ReferenceCollection) {
	if v.name == Var_ARGUMENTS {
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
		c.refByKey[scopedVarKey{scopeRoot: v.scope.GetRootNode(), name: v.name}] = col
	}
	if c.referenceMap != nil {
		c.referenceMap.Put(v, col)
	}
}

func (c *ReferenceCollector) getOrCreateReferences(v *Var) *ReferenceCollection {
	if v == c.lastVar && v != nil {
		return c.lastCollection
	}
	if idx := int(v.nameHash) - 1; idx >= 0 && idx < len(c.vars) && c.vars[idx] == v {
		col := c.cols[idx]
		c.lastVar = v
		c.lastCollection = col
		return col
	}
	if len(c.vars) <= len(c.varsBuf) {
		for i, kv := range c.vars {
			if kv == v || (kv.name == v.name && kv.scope.GetRootNode() == v.scope.GetRootNode()) {
				col := c.cols[i]
				v.nameHash = int32(i + 1)
				c.lastVar = v
				c.lastCollection = col
				return col
			}
		}
	}
	col := c.allocReferenceCollection()
	if len(c.vars) == cap(c.vars) {
		c.growVarsCols()
	}
	c.vars = append(c.vars, v)
	c.cols = append(c.cols, col)
	c.recordNewVarCol(v, col)
	v.nameHash = int32(len(c.vars))
	c.lastVar = v
	c.lastCollection = col
	return col
}

func (c *ReferenceCollector) getReferencesSlow(v *Var) *ReferenceCollection {
	if s := v.scope; s != nil {
		if s.implicitVars != nil {
			if imp := AbstractScope_ImplicitVar_Of(v.name); imp != nil {
				if canon := s.implicitVars[*imp]; canon != nil && canon != v {
					if idx := int(canon.nameHash) - 1; idx >= 0 && idx < len(c.vars) && c.vars[idx] == canon {
						col := c.cols[idx]
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
	return c.getReferencesByScopedName(v.scope.GetRootNode(), v.name)
}

func (c *ReferenceCollector) getReferencesByScopedName(scopeRoot *rhino.Node, name string) *ReferenceCollection {
	if name == Var_ARGUMENTS {
		if c.argumentsRefs != nil {
			return c.argumentsRefs[scopeRoot]
		}
		if c.referenceMap == nil {
			return nil
		}
	}
	if len(c.vars) <= len(c.varsBuf) {
		for i, kv := range c.vars {
			if kv.name == name && kv.scope.GetRootNode() == scopeRoot {
				return c.cols[i]
			}
		}
	} else {
		if c.refByKey == nil {
			c.refByKey = make(map[scopedVarKey]*ReferenceCollection, len(c.vars))
			for i, kv := range c.vars {
				c.refByKey[scopedVarKey{scopeRoot: kv.scope.GetRootNode(), name: kv.name}] = c.cols[i]
			}
		}
		if col := c.refByKey[scopedVarKey{scopeRoot: scopeRoot, name: name}]; col != nil {
			return col
		}
	}
	if c.referenceMap != nil {
		for k, col := range c.referenceMap.All() {
			if k != nil && k.name == name && k.scope.GetRootNode() == scopeRoot {
				return col
			}
		}
	}
	return nil
}

func (c *ReferenceCollector) ensureRawReferenceMap() *jutil.ObjectLinkedHashMap[*Var, *ReferenceCollection] {
	if c.referenceMap == nil {
		c.referenceMap = jutil.NewObjectLinkedHashMap[*Var, *ReferenceCollection]()
		for i, v := range c.vars {
			v.nameHash = 0
			c.referenceMap.Put(v, c.cols[i])
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
	return slices.Values(c.vars)
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
	if idx := int(v.nameHash) - 1; idx >= 0 && idx < len(c.vars) && c.vars[idx] == v {
		col := c.cols[idx]
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
	c.blockStack = oldBlockStack[len(oldBlockStack):len(oldBlockStack)]

	if varScope.IsGlobal() {
		c.blockStack = append(c.blockStack, referenceCollector_getFirst(oldBlockStack))
	} else {
		rootNode := varScope.GetRootNode()
		for _, b := range oldBlockStack {
			c.blockStack = append(c.blockStack, b)
			if b.GetRoot() == rootNode {
				break
			}
		}
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

// ReferenceCollector_CollectorCallback is the traversal callback.
type ReferenceCollector_CollectorCallback struct {
	c *ReferenceCollector
}

// ShouldTraverse updates block stack.
//
// Java: ReferenceCollector.CollectorCallback.shouldTraverse(NodeTraversal, Node, Node)
func (cb *ReferenceCollector_CollectorCallback) ShouldTraverse(nodeTraversal *NodeTraversal, n *rhino.Node, parent *rhino.Node) bool {
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

		scope := t.GetScope()
		v := scope.GetVar(n.GetString())
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
	var peekLast *BasicBlock
	if len(c.blockStack) > 0 {
		peekLast = c.blockStack[len(c.blockStack)-1]
	}
	c.blockStack = append(c.blockStack, c.allocBasicBlock(peekLast, root))
}

// Java: ReferenceCollector.popLastBlock(Node)
func (c *ReferenceCollector) PopLastBlock(root *rhino.Node) {
	if len(c.blockStack) == 0 {
		panic(jutil.NewNoSuchElementException()) // ArrayDeque.removeLast()
	}
	last := c.blockStack[len(c.blockStack)-1]
	c.blockStack[len(c.blockStack)-1] = nil
	c.blockStack = c.blockStack[:len(c.blockStack)-1]
	// Verfiy that the stack is unwound correctly.
	jutil.CheckStateMsg1(base.JSCompObjects_Identical(root, last.GetRoot()), "%s", root)
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
	if v == c.lastVar && v != nil {
		collection = c.lastCollection
	} else {
		if c.varFilter != nil && (v == c.lastRejectedVar || !c.varFilter(v)) {
			c.lastRejectedVar = v
			return
		}
		collection = c.getOrCreateReferences(v)
	}

	// Add this particular reference
	if len(c.blockStack) == 0 {
		panic(jutil.NewNoSuchElementException()) // ArrayDeque.getLast()
	}
	if scope == nil {
		scope = t.GetScope()
	}
	inputId := t.inputId
	if t.compilerInput == nil {
		inputId = t.GetInput().GetInputId()
	}
	ref := c.allocReference(n, scope, c.blockStack[len(c.blockStack)-1], inputId)
	s := collection.refsList.Slice()
	if len(s) == cap(s) {
		c.growReferenceCollection(collection, s)
	}
	collection.refsList.Add(ref)
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
		return w.c.getReferencesByScopedName(scopeRoot, name)
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
