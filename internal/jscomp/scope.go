// Port of com/google/javascript/jscomp/Scope.java @ closure-compiler v20261015

/*
 * Copyright 2004 The Closure Compiler Authors.
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
	"sync"
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

const (
	minScopeChunkSize    = 16
	maxScopeChunkSize    = 4096
	maxScopeChunks       = 128
	minVarChunkSize      = 32
	maxVarChunkSize      = 4096
	maxVarChunks         = 128
	minSlotChunkSize     = 128
	maxSlotChunkSize     = 16384
	maxSlotChunks        = 128
	minHashChunkSize     = 256
	maxHashChunkSize     = 32768
	maxHashChunks        = 128
	minImplicitChunkSize = 8
	maxImplicitChunkSize = 128
	maxImplicitChunks    = 32
)

func scopeNameHash(s string) uint32 {
	h := uint32(2166136261)
	for i := 0; i < len(s); i++ {
		h ^= uint32(s[i])
		h *= 16777619
	}
	return h ^ (h >> 16)
}

// Scope contains information about a variable scope in JavaScript. Scopes
// can be nested, a scope points back to its parent scope. A Scope contains
// information about variables defined in that scope.
//
// Most methods come from AbstractScopeBase[*Scope, *Var] and are typed:
// GetVar returns *Var, GetParent and GetClosestHoistScope return *Scope,
// GetVarIterable is an iter.Seq[*Var]. GetSlot, GetOwnSlot and
// GetParentScope keep rhino.StaticScope's signatures; see the note at the
// top of abstract_scope.go.
//
// See NodeTraversal.
type Scope struct {
	AbstractScopeBase[*Scope, *Var]

	parent   *Scope
	arena    *scopeArena
	varList  []*Var
	varHash  []int32
	varCount int32
}

// scopeArena slab-allocates Scope, Var, varList, and varHash slices for a
// SyntacticScopeCreator and its descendant scopes.
type scopeArena struct {
	scopeChunk       []Scope
	scopeUsed        int
	scopeChunks      [][]Scope
	scopeChunkIdx    int
	varChunk         []Var
	varUsed          int
	varChunks        [][]Var
	varChunkIdx      int
	slotChunk        []*Var
	slotUsed         int
	slotChunks       [][]*Var
	slotChunkIdx     int
	hashChunk        []int32
	hashUsed         int
	hashChunks       [][]int32
	hashChunkIdx     int
	implicitChunk    [][4]*Var
	implicitUsed     int
	implicitChunks   [][][4]*Var
	implicitChunkIdx int

	globalVarList []*Var
	globalVarHash []int32
	globalScope   *Scope
	pooled        bool
}

var cachedScopeArenas [4]atomic.Pointer[scopeArena]

var scopeArenaPool = sync.Pool{
	New: func() any {
		return &scopeArena{pooled: true}
	},
}

func acquireScopeArena() *scopeArena {
	for i := range cachedScopeArenas {
		if a := cachedScopeArenas[i].Swap(nil); a != nil {
			a.pooled = true
			return a
		}
	}
	a := scopeArenaPool.Get().(*scopeArena)
	a.pooled = true
	return a
}

func (a *scopeArena) release() {
	if a == nil || !a.pooled {
		return
	}
	if gs := a.globalScope; gs != nil {
		a.globalScope = nil
		if cap(gs.varList) <= 32768 {
			clear(gs.varList)
			a.globalVarList = gs.varList[:0]
		} else {
			a.globalVarList = nil
		}
		if len(gs.varHash) <= 65536 {
			clear(gs.varHash)
			a.globalVarHash = gs.varHash
		} else {
			a.globalVarHash = nil
		}
	}

	for i := 0; i < a.scopeChunkIdx && i < maxScopeChunks; i++ {
		if i == a.scopeChunkIdx-1 {
			clear(a.scopeChunks[i][:a.scopeUsed])
		} else {
			clear(a.scopeChunks[i])
		}
	}
	if len(a.scopeChunks) > maxScopeChunks {
		clear(a.scopeChunks[maxScopeChunks:])
		a.scopeChunks = a.scopeChunks[:maxScopeChunks]
	}
	a.scopeChunk = nil
	a.scopeUsed = 0
	a.scopeChunkIdx = 0

	for i := 0; i < a.varChunkIdx && i < maxVarChunks; i++ {
		if i == a.varChunkIdx-1 {
			clear(a.varChunks[i][:a.varUsed])
		} else {
			clear(a.varChunks[i])
		}
	}
	if len(a.varChunks) > maxVarChunks {
		clear(a.varChunks[maxVarChunks:])
		a.varChunks = a.varChunks[:maxVarChunks]
	}
	a.varChunk = nil
	a.varUsed = 0
	a.varChunkIdx = 0

	for i := 0; i < a.slotChunkIdx && i < maxSlotChunks; i++ {
		if i == a.slotChunkIdx-1 {
			clear(a.slotChunks[i][:a.slotUsed])
		} else {
			clear(a.slotChunks[i])
		}
	}
	if len(a.slotChunks) > maxSlotChunks {
		clear(a.slotChunks[maxSlotChunks:])
		a.slotChunks = a.slotChunks[:maxSlotChunks]
	}
	a.slotChunk = nil
	a.slotUsed = 0
	a.slotChunkIdx = 0

	for i := 0; i < a.hashChunkIdx && i < maxHashChunks; i++ {
		if i == a.hashChunkIdx-1 {
			clear(a.hashChunks[i][:a.hashUsed])
		} else {
			clear(a.hashChunks[i])
		}
	}
	if len(a.hashChunks) > maxHashChunks {
		clear(a.hashChunks[maxHashChunks:])
		a.hashChunks = a.hashChunks[:maxHashChunks]
	}
	a.hashChunk = nil
	a.hashUsed = 0
	a.hashChunkIdx = 0

	for i := 0; i < a.implicitChunkIdx && i < maxImplicitChunks; i++ {
		if i == a.implicitChunkIdx-1 {
			clear(a.implicitChunks[i][:a.implicitUsed])
		} else {
			clear(a.implicitChunks[i])
		}
	}
	if len(a.implicitChunks) > maxImplicitChunks {
		clear(a.implicitChunks[maxImplicitChunks:])
		a.implicitChunks = a.implicitChunks[:maxImplicitChunks]
	}
	a.implicitChunk = nil
	a.implicitUsed = 0
	a.implicitChunkIdx = 0

	for i := range cachedScopeArenas {
		if cachedScopeArenas[i].CompareAndSwap(nil, a) {
			return
		}
	}
	scopeArenaPool.Put(a)
}

func (a *scopeArena) allocScope() *Scope {
	if a.scopeUsed == len(a.scopeChunk) {
		if a.scopeChunkIdx < len(a.scopeChunks) {
			a.scopeChunk = a.scopeChunks[a.scopeChunkIdx]
			a.scopeChunkIdx++
			a.scopeUsed = 1
			return &a.scopeChunk[0]
		}
		n := 4 * len(a.scopeChunk)
		if n < minScopeChunkSize {
			n = minScopeChunkSize
		} else if n > maxScopeChunkSize {
			n = maxScopeChunkSize
		}
		ch := make([]Scope, n)
		a.scopeChunks = append(a.scopeChunks, ch)
		a.scopeChunkIdx = len(a.scopeChunks)
		a.scopeChunk = ch
		a.scopeUsed = 1
		return &ch[0]
	}
	s := &a.scopeChunk[a.scopeUsed]
	a.scopeUsed++
	return s
}

func (a *scopeArena) allocVar() *Var {
	if a.varUsed == len(a.varChunk) {
		if a.varChunkIdx < len(a.varChunks) {
			a.varChunk = a.varChunks[a.varChunkIdx]
			a.varChunkIdx++
			a.varUsed = 1
			return &a.varChunk[0]
		}
		n := 4 * len(a.varChunk)
		if n < minVarChunkSize {
			n = minVarChunkSize
		} else if n > maxVarChunkSize {
			n = maxVarChunkSize
		}
		ch := make([]Var, n)
		a.varChunks = append(a.varChunks, ch)
		a.varChunkIdx = len(a.varChunks)
		a.varChunk = ch
		a.varUsed = 1
		return &ch[0]
	}
	v := &a.varChunk[a.varUsed]
	a.varUsed++
	return v
}

func (a *scopeArena) allocVarList(length, capacity int) []*Var {
	if capacity <= 0 {
		return nil
	}
	if capacity > maxSlotChunkSize {
		return make([]*Var, length, capacity)
	}
	if len(a.slotChunk)-a.slotUsed < capacity {
		if a.slotChunkIdx < len(a.slotChunks) {
			ch := a.slotChunks[a.slotChunkIdx]
			if len(ch) < capacity {
				size := 4 * len(ch)
				if size < minSlotChunkSize {
					size = minSlotChunkSize
				}
				if size < capacity {
					size = capacity
				}
				if size > maxSlotChunkSize {
					size = maxSlotChunkSize
				}
				ch = make([]*Var, size)
				a.slotChunks[a.slotChunkIdx] = ch
			}
			a.slotChunkIdx++
			a.slotChunk = ch
			a.slotUsed = capacity
			return ch[:length:capacity]
		}
		size := 4 * len(a.slotChunk)
		if size < minSlotChunkSize {
			size = minSlotChunkSize
		} else if size > maxSlotChunkSize {
			size = maxSlotChunkSize
		}
		if size < capacity {
			size = capacity
		}
		ch := make([]*Var, size)
		a.slotChunks = append(a.slotChunks, ch)
		a.slotChunkIdx = len(a.slotChunks)
		a.slotChunk = ch
		a.slotUsed = capacity
		return ch[:length:capacity]
	}
	u := a.slotUsed
	a.slotUsed = u + capacity
	return a.slotChunk[u : u+length : u+capacity]
}

func (a *scopeArena) allocVarHash(size int) []int32 {
	if size <= 0 {
		return nil
	}
	if size > maxHashChunkSize {
		return make([]int32, size)
	}
	if len(a.hashChunk)-a.hashUsed < size {
		if a.hashChunkIdx < len(a.hashChunks) {
			ch := a.hashChunks[a.hashChunkIdx]
			if len(ch) < size {
				chunkSize := 4 * len(ch)
				if chunkSize < minHashChunkSize {
					chunkSize = minHashChunkSize
				}
				if chunkSize < size {
					chunkSize = size
				}
				if chunkSize > maxHashChunkSize {
					chunkSize = maxHashChunkSize
				}
				ch = make([]int32, chunkSize)
				a.hashChunks[a.hashChunkIdx] = ch
			}
			a.hashChunkIdx++
			a.hashChunk = ch
			a.hashUsed = size
			return ch[:size:size]
		}
		chunkSize := 4 * len(a.hashChunk)
		if chunkSize < minHashChunkSize {
			chunkSize = minHashChunkSize
		} else if chunkSize > maxHashChunkSize {
			chunkSize = maxHashChunkSize
		}
		if chunkSize < size {
			chunkSize = size
		}
		ch := make([]int32, chunkSize)
		a.hashChunks = append(a.hashChunks, ch)
		a.hashChunkIdx = len(a.hashChunks)
		a.hashChunk = ch
		a.hashUsed = size
		return ch[:size:size]
	}
	u := a.hashUsed
	a.hashUsed = u + size
	return a.hashChunk[u : u+size : u+size]
}

func (sc *Scope) initLargeVars(n int) {
	if a := sc.arena; a != nil {
		if sc.parent == nil && a.globalScope == nil {
			a.globalScope = sc
			listCap := 256
			for listCap < n {
				listCap *= 2
			}
			if cap(a.globalVarList) < listCap {
				a.globalVarList = make([]*Var, 0, listCap)
			}
			sc.varList = a.globalVarList[:0]
			a.globalVarList = nil

			hashLen := 512
			for hashLen < cap(sc.varList)*2 {
				hashLen *= 2
			}
			if len(a.globalVarHash) < hashLen {
				a.globalVarHash = make([]int32, hashLen)
			}
			sc.varHash = a.globalVarHash
			a.globalVarHash = nil
			return
		}
		sc.varList = a.allocVarList(0, n)
		hashLen := 32
		for hashLen < n*2 {
			hashLen *= 2
		}
		sc.varHash = a.allocVarHash(hashLen)
		return
	}
	sc.varList = make([]*Var, 0, n)
	hashLen := 32
	for hashLen < n*2 {
		hashLen *= 2
	}
	sc.varHash = make([]int32, hashLen)
}

func (a *scopeArena) allocImplicitVars() *[4]*Var {
	if a.implicitUsed == len(a.implicitChunk) {
		if a.implicitChunkIdx < len(a.implicitChunks) {
			a.implicitChunk = a.implicitChunks[a.implicitChunkIdx]
			a.implicitChunkIdx++
			a.implicitUsed = 1
			return &a.implicitChunk[0]
		}
		n := 4 * len(a.implicitChunk)
		if n < minImplicitChunkSize {
			n = minImplicitChunkSize
		} else if n > maxImplicitChunkSize {
			n = maxImplicitChunkSize
		}
		ch := make([][4]*Var, n)
		a.implicitChunks = append(a.implicitChunks, ch)
		a.implicitChunkIdx = len(a.implicitChunks)
		a.implicitChunk = ch
		a.implicitUsed = 1
		return &ch[0]
	}
	p := &a.implicitChunk[a.implicitUsed]
	a.implicitUsed++
	return p
}

// Java: Scope.createGlobalScope(Node)
func Scope_CreateGlobalScope(rootNode *rhino.Node) *Scope {
	return newScope(rootNode)
}

// Java: Scope.createChildScope(Scope, Node)
func Scope_CreateChildScope(parent *Scope, rootNode *rhino.Node) *Scope {
	return newScopeWithParent(parent, rootNode)
}

// Java: Scope(Scope, Node)
func newScopeWithParent(parent *Scope, rootNode *rhino.Node) *Scope {
	var s *Scope
	if parent != nil && parent.arena != nil {
		s = parent.arena.allocScope()
		s.arena = parent.arena
	} else {
		s = &Scope{}
	}
	initScopeWithParent(s, parent, rootNode)
	return s
}

func newScopeWithParentInArena(parent *Scope, rootNode *rhino.Node, arena *scopeArena) *Scope {
	s := arena.allocScope()
	s.arena = arena
	initScopeWithParent(s, parent, rootNode)
	return s
}

func initScopeWithParent(s *Scope, parent *Scope, rootNode *rhino.Node) {
	if rootNode == nil || parent == nil {
		panic(jutil.NewNullPointerException())
	}
	if !NodeUtil_CreatesScope(rootNode) {
		panic(jutil.NewIllegalArgumentExceptionMsg(jutil.LenientFormat("%s", rootNode)))
	}
	if rootNode == parent.rootNode {
		panic(jutil.NewIllegalArgumentExceptionMsg(
			jutil.LenientFormat("rootNode should not be the parent's root node: %s", rootNode)))
	}
	s.self = s
	s.rootNode = rootNode
	s.parent = parent
	s.depth = parent.depth + 1
}

// Java: Scope(Node)
func newScope(rootNode *rhino.Node) *Scope {
	s := &Scope{}
	initScope(s, rootNode)
	return s
}

func newScopeInArena(rootNode *rhino.Node, arena *scopeArena) *Scope {
	s := arena.allocScope()
	s.arena = arena
	initScope(s, rootNode)
	return s
}

func initScope(s *Scope, rootNode *rhino.Node) {
	if rootNode == nil {
		panic(jutil.NewNullPointerException())
	}
	if !NodeUtil_CreatesScope(rootNode) && !rootNode.IsScript() && !rootNode.IsRoot() {
		panic(jutil.NewIllegalArgumentExceptionMsg(jutil.LenientFormat("%s", rootNode)))
	}
	s.self = s
	s.rootNode = rootNode
}

// Java: Scope.untyped()
func (s *Scope) Untyped() *Scope {
	return s
}

// Java: Scope.getParent()
func (s *Scope) GetParent() *Scope {
	return s.parent
}

// Java: AbstractScope.getParentScope()
func (s *Scope) GetParentScope() rhino.StaticScope {
	if s.parent != nil {
		return s.parent
	}
	return nil
}

// Java: AbstractScope.isGlobal()
func (s *Scope) IsGlobal() bool {
	return s.parent == nil
}

// Java: AbstractScope.isLocal()
func (s *Scope) IsLocal() bool {
	return s.parent != nil
}

// Java: AbstractScope.getGlobalScope()
func (s *Scope) GetGlobalScope() *Scope {
	for s.parent != nil {
		s = s.parent
	}
	return s
}

// Java: AbstractScope.isHoistScope()
func (s *Scope) IsHoistScope() bool {
	return s.parent == nil ||
		s.rootNode.IsFunction() ||
		NodeUtil_IsFunctionBlock(s.rootNode) ||
		s.rootNode.IsModuleBody() ||
		NodeUtil_IsClassStaticBlock(s.rootNode)
}

// Java: AbstractScope.getClosestHoistScope()
func (s *Scope) GetClosestHoistScope() *Scope {
	for ; s != nil; s = s.parent {
		if s.IsHoistScope() {
			return s
		}
	}
	return nil
}

// Java: AbstractScope.getClosestContainerScope()
func (s *Scope) GetClosestContainerScope() *Scope {
	scope := s.GetClosestHoistScope()
	if NodeUtil_IsFunctionBlock(scope.rootNode) {
		scope = scope.parent
		if scope.IsBlockScope() {
			panic(jutil.NewIllegalStateExceptionMsg(jutil.LenientFormat("%s", scope)))
		}
	}
	return scope
}

func (s *Scope) lookupOwnVar(name string) *Var {
	if s.varCount == 0 {
		return nil
	}
	if s.varHash != nil {
		mask := uint32(len(s.varHash) - 1)
		slot := scopeNameHash(name) & mask
		for {
			idx := s.varHash[slot]
			if idx == 0 {
				return nil
			}
			if v := s.varList[idx-1]; v.name == name {
				return v
			}
			slot = (slot + 1) & mask
		}
	}
	for _, v := range s.varList {
		if v != nil && v.name == name {
			return v
		}
	}
	return nil
}

// Java: AbstractScope.getOwnSlot(String) (final)
func (s *Scope) GetOwnSlot(name string) rhino.StaticSlot {
	if v := s.getOwnSlotV(name); v != nil {
		return v
	}
	return nil
}

func (s *Scope) getOwnSlotV(name string) *Var {
	if v := s.lookupOwnVar(name); v != nil {
		return v
	}
	if imp := AbstractScope_ImplicitVar_Of(name); imp != nil {
		return s.GetOwnImplicitSlot(imp)
	}
	return nil
}

// Java: AbstractScope.hasOwnSlot(String)
func (s *Scope) HasOwnSlot(name string) bool {
	if s.lookupOwnVar(name) != nil {
		return true
	}
	if imp := AbstractScope_ImplicitVar_Of(name); imp != nil {
		return s.HasOwnImplicitSlot(imp)
	}
	return false
}

// Java: AbstractScope.hasSlot(String)
func (s *Scope) HasSlot(name string) bool {
	implicit := AbstractScope_ImplicitVar_Of(name)
	var h uint32
	var hashed bool
	for ; s != nil; s = s.parent {
		if s.varCount > 0 {
			if s.varHash != nil {
				if !hashed {
					h = scopeNameHash(name)
					hashed = true
				}
				mask := uint32(len(s.varHash) - 1)
				slot := h & mask
				for {
					idx := s.varHash[slot]
					if idx == 0 {
						break
					}
					if s.varList[idx-1].name == name {
						return true
					}
					slot = (slot + 1) & mask
				}
			} else {
				for _, v := range s.varList {
					if v != nil && v.name == name {
						return true
					}
				}
			}
		}
		if implicit != nil && implicit.IsMadeByScope(s) {
			return true
		}
	}
	return false
}

func (s *Scope) inheritedInput() (*CompilerInput, bool) {
	for anc := s.parent; anc != nil && anc.parent != nil; anc = anc.parent {
		if anc.varCount > 0 {
			for _, v := range anc.varList {
				if v != nil && !v.IsImplicitGoogNamespace() {
					return v.GetInput(), true
				}
			}
		}
	}
	return nil, false
}

// GetVar uses the concrete parent links of Java's final Scope class. The
// generic base must support TypedScope's virtual methods; untyped scopes can
// avoid those dictionary calls and classify implicit names once per lookup.
func (s *Scope) GetVar(name string) *Var {
	implicit := AbstractScope_ImplicitVar_Of(name)
	var h uint32
	var hashed bool
	for {
		if s.varCount > 0 {
			if s.varHash != nil {
				if !hashed {
					h = scopeNameHash(name)
					hashed = true
				}
				mask := uint32(len(s.varHash) - 1)
				slot := h & mask
				for {
					idx := s.varHash[slot]
					if idx == 0 {
						break
					}
					if v := s.varList[idx-1]; v.name == name {
						return v
					}
					slot = (slot + 1) & mask
				}
			} else {
				for _, v := range s.varList {
					if v != nil && v.name == name {
						return v
					}
				}
			}
		}
		if implicit != nil {
			if v := s.GetOwnImplicitSlot(implicit); v != nil {
				return v
			}
		}
		s = s.parent
		if s == nil {
			return nil
		}
	}
}

// Java: AbstractScope.getVarCount()
func (s *Scope) GetVarCount() int {
	return int(s.varCount)
}

// Java: AbstractScope.getVarIterable()
func (s *Scope) GetVarIterable() iter.Seq[*Var] {
	return func(yield func(*Var) bool) {
		for i := 0; i < len(s.varList); i++ {
			if v := s.varList[i]; v != nil {
				if !yield(v) {
					return
				}
			}
		}
	}
}

// Java: AbstractScope.getAllSymbols()
func (s *Scope) GetAllSymbols() iter.Seq[*Var] {
	return s.GetVarIterable()
}

// Java: AbstractScope.clearVarsInternal()
func (s *Scope) ClearVarsInternal() {
	if s.varCount > 0 || len(s.varList) > 0 {
		clear(s.varList)
		s.varList = s.varList[:0]
		if s.varHash != nil {
			clear(s.varHash)
		}
		s.varCount = 0
	}
}

// Java: AbstractScope.undeclare(V)
func (s *Scope) Undeclare(var_ *Var) {
	jutil.CheckState(var_.scope == s)
	current := s.lookupOwnVar(var_.name)
	if current == nil {
		panic(jutil.NewNullPointerException())
	}
	jutil.CheckState(current.Equals(var_))
	s.UndeclareInteral(var_)
}

// Java: AbstractScope.undeclareInteral(V)
func (s *Scope) UndeclareInteral(var_ *Var) {
	if !s.varsInitialized {
		panic(jutil.NewUnsupportedOperationException())
	}
	name := var_.name
	if s.varHash != nil {
		mask := uint32(len(s.varHash) - 1)
		i := scopeNameHash(name) & mask
		for {
			idx := s.varHash[i]
			if idx == 0 {
				return
			}
			if s.varList[idx-1].name == name {
				s.varList[idx-1] = nil
				s.varCount--
				break
			}
			i = (i + 1) & mask
		}
		j := i
		for {
			j = (j + 1) & mask
			idx := s.varHash[j]
			if idx == 0 {
				break
			}
			k := scopeNameHash(s.varList[idx-1].name) & mask
			if ((j - k) & mask) >= ((j - i) & mask) {
				s.varHash[i] = idx
				i = j
			}
		}
		s.varHash[i] = 0
		return
	}
	for i, v := range s.varList {
		if v != nil && v.name == name {
			s.varList[i] = nil
			s.varCount--
			return
		}
	}
}

func (s *Scope) appendVar(v *Var) {
	if len(s.varList) < cap(s.varList) {
		s.varList = append(s.varList, v)
		return
	}
	s.growVarList(v)
}

func (s *Scope) growVarList(v *Var) {
	oldLen := len(s.varList)
	newCap := cap(s.varList) * 2
	if newCap < 4 {
		newCap = 4
	}
	if a := s.arena; a != nil && s != a.globalScope && newCap <= maxSlotChunkSize {
		if oldLen > 0 && a.slotUsed >= oldLen && &s.varList[0] == &a.slotChunk[a.slotUsed-oldLen] {
			start := a.slotUsed - oldLen
			if len(a.slotChunk)-start >= newCap {
				a.slotUsed = start + newCap
				s.varList = a.slotChunk[start : start+oldLen+1 : start+newCap]
				s.varList[oldLen] = v
				return
			}
		}
		next := a.allocVarList(oldLen+1, newCap)
		copy(next, s.varList)
		next[oldLen] = v
		s.varList = next
		return
	}
	next := make([]*Var, oldLen+1, newCap)
	copy(next, s.varList)
	next[oldLen] = v
	s.varList = next
}

func (s *Scope) initVarHash(size int) {
	if a := s.arena; a != nil && s != a.globalScope && size <= maxHashChunkSize {
		s.varHash = a.allocVarHash(size)
	} else {
		s.varHash = make([]int32, size)
	}
	mask := uint32(size - 1)
	for i, v := range s.varList {
		if v == nil {
			continue
		}
		slot := scopeNameHash(v.name) & mask
		for s.varHash[slot] != 0 {
			slot = (slot + 1) & mask
		}
		s.varHash[slot] = int32(i + 1)
	}
}

func (s *Scope) growVarHash() {
	newLen := len(s.varHash) * 2
	if newLen < 32 {
		newLen = 32
	}
	if a := s.arena; a != nil && s != a.globalScope && newLen <= maxHashChunkSize {
		oldLen := len(s.varHash)
		if oldLen > 0 && a.hashUsed >= oldLen && &s.varHash[0] == &a.hashChunk[a.hashUsed-oldLen] {
			start := a.hashUsed - oldLen
			if len(a.hashChunk)-start >= newLen {
				clear(s.varHash)
				a.hashUsed = start + newLen
				s.varHash = a.hashChunk[start : start+newLen : start+newLen]
			} else {
				s.varHash = a.allocVarHash(newLen)
			}
		} else {
			s.varHash = a.allocVarHash(newLen)
		}
	} else {
		s.varHash = make([]int32, newLen)
	}
	mask := uint32(newLen - 1)
	for i, v := range s.varList {
		if v == nil {
			continue
		}
		slot := scopeNameHash(v.name) & mask
		for s.varHash[slot] != 0 {
			slot = (slot + 1) & mask
		}
		s.varHash[slot] = int32(i + 1)
	}
}

func (s *Scope) appendVarToInitialized(var_ *Var, probeSlot uint32) {
	s.appendVar(var_)
	s.varCount++
	if int(s.varCount)*2 > len(s.varHash) {
		s.growVarHash()
		return
	}
	s.varHash[probeSlot] = int32(len(s.varList))
}

// Java: AbstractScope.declareInternal(String, V)
func (s *Scope) DeclareInternal(name string, var_ *Var) {
	if s.IsFunctionBlockScope() && s.parent.HasOwnSlot(name) && !s.HasOwnSlot(name) && !s.IsBleedingFunctionName(name) {
		panic(jutil.NewIllegalStateExceptionMsg(jutil.LenientFormat("Illegal shadow: %s", var_.GetNode())))
	}
	s.varsInitialized = true
	if s.varHash != nil {
		h := scopeNameHash(name)
		mask := uint32(len(s.varHash) - 1)
		slot := h & mask
		for {
			idx := s.varHash[slot]
			if idx == 0 {
				break
			}
			if s.varList[idx-1].name == name {
				s.varList[idx-1] = var_
				return
			}
			slot = (slot + 1) & mask
		}
		s.appendVarToInitialized(var_, slot)
		return
	}
	for i, existing := range s.varList {
		if existing != nil && existing.name == name {
			s.varList[i] = var_
			return
		}
	}
	s.appendVar(var_)
	s.varCount++
	if len(s.varList) > 8 {
		s.initVarHash(32)
	}
}

// Java: AbstractScope.hasOwnImplicitSlot(ImplicitVar)
func (s *Scope) HasOwnImplicitSlot(name *AbstractScope_ImplicitVar) bool {
	return name != nil && name.IsMadeByScope(s)
}

// Java: AbstractScope.getOwnImplicitSlot(ImplicitVar)
func (s *Scope) GetOwnImplicitSlot(name *AbstractScope_ImplicitVar) *Var {
	if !s.HasOwnImplicitSlot(name) {
		return nil
	}
	if s.implicitVars == nil {
		if s.arena != nil {
			s.implicitVars = s.arena.allocImplicitVars()
		} else {
			s.implicitVars = new([4]*Var)
		}
	}
	if v := s.implicitVars[*name]; v != nil {
		return v
	}
	v := s.MakeImplicitVar(*name)
	if v != nil {
		s.implicitVars[*name] = v
	}
	return v
}

// Java: AbstractScope.getSlot(String) (final)
func (s *Scope) GetSlot(name string) rhino.StaticSlot {
	if v := s.GetVar(name); v != nil {
		return v
	}
	return nil
}

// Declare declares a variable.
//
// name is the name of the variable, nameNode the NAME node declaring the
// variable, input the input in which this variable is defined.
//
// Java: Scope.declare(String, Node, CompilerInput)
func (s *Scope) Declare(name string, nameNode *rhino.Node, input *CompilerInput) *Var {
	jutil.CheckArgument(name != "")
	// Make sure that it's declared only once
	jutil.CheckState(s.getOwnSlotV(name) == nil)
	var var_ *Var
	if s.arena != nil {
		var_ = s.arena.allocVar()
		initVar(var_, name, nameNode, s, int(s.varCount), input, nil)
	} else {
		var_ = NewVar(
			name,
			nameNode,
			s,
			int(s.varCount),
			input,
			/* implicitGoogNamespaceDefinition= */ nil)
	}
	s.DeclareInternal(name, var_)
	return var_
}

func functionScopeVarSlotCount(fn *rhino.Node) int {
	fnNameNode := fn.GetFirstChild()
	n := 0
	if fnNameNode.GetString() != "" && NodeUtil_IsFunctionExpression(fn) {
		n = 1
	}
	for c := fnNameNode.GetNext().GetFirstChild(); c != nil; c = c.GetNext() {
		switch c.GetToken() {
		case rhino.Token_NAME:
			n++
		case rhino.Token_DEFAULT_VALUE, rhino.Token_ITER_REST:
			if c.GetFirstChild().IsName() {
				n++
			} else {
				n += 4
			}
		default:
			n += 4
		}
	}
	if n < 1 {
		return 1
	}
	if n > 64 {
		return 64
	}
	return n
}

const maxPreallocScopeVars = 512

func countLhsVars(n *rhino.Node) int {
	if n.HasOneChild() && n.GetFirstChild().IsName() {
		return 1
	}
	c := 0
	for child := n.GetFirstChild(); child != nil; child = child.GetNext() {
		if child.IsName() {
			c++
		} else {
			c += 2
		}
	}
	if c == 0 {
		return 1
	}
	return c
}

func countScopeVars(n *rhino.Node, isHoist bool, inBlockScope bool, root *rhino.Node, count int, steps *int) int {
	if count >= maxPreallocScopeVars {
		return maxPreallocScopeVars
	}
	*steps++
	if *steps > 256 {
		if count < 16 {
			return 32
		}
		if count*2 > maxPreallocScopeVars {
			return maxPreallocScopeVars
		}
		return count * 2
	}
	switch n.GetToken() {
	case rhino.Token_VAR, rhino.Token_IMPORT:
		if isHoist {
			count += countLhsVars(n)
		}
		return count
	case rhino.Token_LET, rhino.Token_CONST:
		if inBlockScope {
			count += countLhsVars(n)
		}
		return count
	case rhino.Token_EXPORT:
		if first := n.GetFirstChild(); first != nil {
			return countScopeVars(first, isHoist, inBlockScope, root, count, steps)
		}
		return count
	case rhino.Token_FUNCTION, rhino.Token_CLASS:
		if inBlockScope && !NodeUtil_IsFunctionExpression(n) && !NodeUtil_IsClassExpression(n) {
			count++
		}
		return count
	case rhino.Token_CATCH:
		if inBlockScope {
			count++
		}
		return countScopeVars(n.GetSecondChild(), isHoist, inBlockScope, root, count, steps)
	case rhino.Token_MODULE_BODY:
		if n != root {
			return count + 1
		}
	case rhino.Token_BLOCK:
		if n != root && scanVarsBlockCreatesScope(n) {
			if !isHoist {
				return count
			}
			inBlockScope = false
		}
	case rhino.Token_FOR, rhino.Token_FOR_IN, rhino.Token_FOR_OF, rhino.Token_FOR_AWAIT_OF, rhino.Token_SWITCH_BODY:
		if n != root {
			if !isHoist {
				return count
			}
			inBlockScope = false
		}
	case rhino.Token_SCRIPT, rhino.Token_ROOT,
		rhino.Token_DO, rhino.Token_WHILE, rhino.Token_WITH, rhino.Token_IF,
		rhino.Token_LABEL, rhino.Token_TRY, rhino.Token_SWITCH, rhino.Token_CASE, rhino.Token_DEFAULT_CASE:
	default:
		return count
	}
	for child := n.GetFirstChild(); child != nil; child = child.GetNext() {
		if scanVarsTokenKind[uint8(child.GetToken())] == 1 {
			count = countScopeVars(child, isHoist, inBlockScope, root, count, steps)
			if count >= maxPreallocScopeVars || *steps > 256 {
				return count
			}
		}
	}
	return count
}

func largeScopeVarSlotCount(sc *Scope) int {
	root := sc.rootNode
	if root.IsFunction() {
		n := functionScopeVarSlotCount(root)
		if n < 12 {
			return 12
		}
		return n
	}
	isHoist := sc.IsHoistScope()
	steps := 0
	n := countScopeVars(root, isHoist, true, root, 0, &steps)
	if n < 16 {
		return 16
	}
	return n
}

// DeclareImplicitGoogNamespaceIfAbsent declares an implicit goog.provide or
// goog.module namespace in this scope.
//
// name is the name of the variable, definition the STRINGLIT node holding
// the full namespace.
//
// Java: Scope.declareImplicitGoogNamespaceIfAbsent(String, Node)
func (s *Scope) DeclareImplicitGoogNamespaceIfAbsent(name string, definition *rhino.Node) *Var {
	jutil.CheckArgument(name != "")
	jutil.CheckStateMsg1(s.IsGlobal(), "Cannot declare implicit goog namespace in local scope %s", s)
	// Allow redeclarations of provides, since they are implicit and don't have a single
	// declaration site.
	var_ := s.getOwnSlotV(name)
	if var_ == nil {
		if s.arena != nil {
			var_ = s.arena.allocVar()
			initVar(var_, name, nil, s, -1, nil, definition)
		} else {
			var_ = Var_CreateImplicitGoogNamespace(name, s, definition)
		}
		s.DeclareInternal(name, var_)
	} else if var_.IsImplicitGoogNamespace() {
		var_.AddImplicitGoogNamespaceDefinition(definition)
	}
	return var_
}

// Java: Scope.makeImplicitVar(ImplicitVar)
func (s *Scope) MakeImplicitVar(var_ AbstractScope_ImplicitVar) *Var {
	if s.arena != nil {
		v := s.arena.allocVar()
		initVar(v, var_.name(), nil, s, -1, nil, nil)
		return v
	}
	return NewVar(
		var_.name(),
		/* nameNode= */ nil,
		s,
		/* index= */ -1,
		/* input= */ nil,
		/* implicitGoogNamespaceDefinition= */ nil)
}
