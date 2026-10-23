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
	"slices"
	"unsafe"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

func scopeNameEq(a, b string) bool {
	return len(a) == len(b) && (unsafe.StringData(a) == unsafe.StringData(b) || (len(a) > 0 && a[0] == b[0] && a == b))
}

func scopeNameHash(s string) uint32 {
	n := len(s)
	if n == 0 {
		return 0
	}
	p := unsafe.Pointer(unsafe.StringData(s))
	var h uint64
	if n >= 8 {
		v0 := *(*uint64)(p)
		v1 := *(*uint64)(unsafe.Add(p, n-8))
		if n <= 16 {
			h = (v0*0x9e3779b97f4a7c15 ^ v1*0xbf58476d1ce4e5b9) + uint64(n)*0x94d049bb133111eb
		} else if n <= 32 {
			v2 := *(*uint64)(unsafe.Add(p, 8))
			v3 := *(*uint64)(unsafe.Add(p, n-16))
			h = (v0*0x9e3779b97f4a7c15 ^ v1*0xbf58476d1ce4e5b9) +
				(v2*0x94d049bb133111eb ^ v3*0x2545f4914f6cdd1d) + uint64(n)
		} else {
			h = (v0*0x9e3779b97f4a7c15 ^ v1*0xbf58476d1ce4e5b9) + uint64(n)*0x94d049bb133111eb
			for off := 8; off+8 < n; off += 8 {
				h = (h ^ *(*uint64)(unsafe.Add(p, off))) * 0x9e3779b97f4a7c15
			}
		}
	} else if n >= 4 {
		v0 := uint64(*(*uint32)(p))
		v1 := uint64(*(*uint32)(unsafe.Add(p, n-4)))
		h = ((v0 << 32) | v1) ^ (uint64(n) * 0x9e3779b97f4a7c15)
	} else {
		b0 := uint32(*(*byte)(p))
		b1 := uint32(*(*byte)(unsafe.Add(p, n>>1)))
		b2 := uint32(*(*byte)(unsafe.Add(p, n-1)))
		u := (b0 | (b1 << 8) | (b2 << 16) | (uint32(n) << 24)) * 0xcc9e2d51
		return u ^ (u >> 16)
	}
	h ^= h >> 33
	h *= 0xff51afd7ed558ccd
	return uint32(h ^ (h >> 32))
}

func scopeBloomFromHash(h uint32) uint64 {
	return (uint64(1) << (h & 63)) | (uint64(1) << ((h >> 6) & 63))
}

func scopeBloomMask(s string) uint64 {
	return scopeBloomFromHash(scopeNameHash(s))
}

func scopeIDHash(id uint32) uint32 {
	h := id * 0x9e3779b1
	h ^= h >> 16
	h *= 0x85ebca6b
	return h ^ (h >> 13)
}

func scopeBloomFromID(id uint32) uint64 {
	return scopeBloomFromHash(scopeIDHash(id))
}

// scopeSmallVars is the most vars a Scope looks up by scanning varList; a
// larger scope builds varIndex.
const scopeSmallVars = 8

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
// Go-only: performance — the Java LinkedHashMap<String, Var> is varList, in
// declaration order, plus varIndex (name ID to varList index) once the scope
// has more than scopeSmallVars vars, plus varBloom, a two-bit-per-name filter
// that lets lookups skip scopes without the name. lookupParent is the
// closest ancestor worth searching: scanned scopes without vars or implicit
// vars are skipped.
//
// See NodeTraversal.
type Scope struct {
	AbstractScopeBase[*Scope, *Var]

	parent       *Scope
	lookupParent *Scope
	varList      []*Var
	varIndex     map[uint32]int32
	varBloom     uint64
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
	s := &Scope{parent: parent}
	s.rootNode = rootNode
	s.depth = parent.depth + 1
	s.initLookupChain(parent)
	return s
}

// Java: Scope(Node)
func newScope(rootNode *rhino.Node) *Scope {
	if rootNode == nil {
		panic(jutil.NewNullPointerException())
	}
	if !NodeUtil_CreatesScope(rootNode) && !rootNode.IsScript() && !rootNode.IsRoot() {
		panic(jutil.NewIllegalArgumentExceptionMsg(jutil.LenientFormat("%s", rootNode)))
	}
	s := &Scope{}
	s.rootNode = rootNode
	s.initLookupChain(nil)
	return s
}

func (s *Scope) initLookupChain(parent *Scope) {
	if parent == nil {
		s.hasImplicitVars = true
		s.lookupParent = nil
		return
	}
	switch s.rootNode.GetToken() {
	case rhino.Token_FUNCTION, rhino.Token_CLASS, rhino.Token_MODULE_BODY, rhino.Token_MEMBER_FIELD_DEF, rhino.Token_COMPUTED_FIELD_DEF:
		s.hasImplicitVars = true
	case rhino.Token_BLOCK:
		s.hasImplicitVars = s.IsStaticBlockScope() || NodeUtil_IsBundledGoogModuleScopeRoot(s.rootNode)
	default:
		s.hasImplicitVars = false
	}
	if parent.scanDone && len(parent.varList) == 0 && !parent.hasImplicitVars && !NodeUtil_IsFunctionBlock(parent.rootNode) {
		s.lookupParent = parent.lookupParent
	} else {
		s.lookupParent = parent
	}
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

func (s *Scope) isBleedingFunctionNameID(id uint32) bool {
	var_ := s.GetVarByID(id)
	return var_ != nil &&
		var_.GetNode() != nil &&
		var_.GetNode().GetParent().IsFunction()
}

// ownIndex returns the varList index of the var named id declared in this
// scope (implicit vars excluded), or -1.
func (s *Scope) ownIndex(id uint32, bloom uint64) int {
	if s.varBloom&bloom != bloom {
		return -1
	}
	if s.varIndex != nil {
		if i, ok := s.varIndex[id]; ok {
			return int(i)
		}
		return -1
	}
	for i, v := range s.varList {
		if v.nameID == id {
			return i
		}
	}
	return -1
}

func (s *Scope) ownVar(id uint32, bloom uint64) *Var {
	if i := s.ownIndex(id, bloom); i >= 0 {
		return s.varList[i]
	}
	return nil
}

func (s *Scope) lookupOwnVarByID(id uint32) *Var {
	if len(s.varList) == 0 || id == 0 {
		return nil
	}
	return s.ownVar(id, scopeBloomFromID(id))
}

// hasOwnSlotByID is hasOwnSlot for an interned name.
func (s *Scope) hasOwnSlotByID(id uint32) bool {
	if s.lookupOwnVarByID(id) != nil {
		return true
	}
	if s.hasImplicitVars {
		if imp := AbstractScope_ImplicitVar_OfID(id); imp != nil {
			return s.HasOwnImplicitSlot(imp)
		}
	}
	return false
}

// Java: AbstractScope.getOwnSlot(String) (final)
func (s *Scope) GetOwnSlot(name string) rhino.StaticSlot {
	if v := s.getOwnSlotV(name); v != nil {
		return v
	}
	return nil
}

func (s *Scope) getOwnSlotV(name string) *Var {
	id := rhino.RhinoStringPool_LookupID(name)
	if id == 0 {
		return nil
	}
	if v := s.lookupOwnVarByID(id); v != nil {
		return v
	}
	if s.hasImplicitVars {
		if imp := AbstractScope_ImplicitVar_OfID(id); imp != nil {
			return s.GetOwnImplicitSlot(imp)
		}
	}
	return nil
}

// Java: AbstractScope.hasOwnSlot(String)
func (s *Scope) HasOwnSlot(name string) bool {
	id := rhino.RhinoStringPool_LookupID(name)
	return id != 0 && s.hasOwnSlotByID(id)
}

// Java: AbstractScope.canDeclare(String)
func (s *Scope) CanDeclare(name string) bool {
	id := rhino.RhinoStringPool_LookupID(name)
	if id == 0 {
		return true
	}
	if s.hasOwnSlotByID(id) {
		return false
	}
	if s.IsFunctionBlockScope() && s.parent.hasOwnSlotByID(id) {
		return s.isBleedingFunctionNameID(id)
	}
	return true
}

// Java: AbstractScope.hasSlot(String)
func (s *Scope) HasSlot(name string) bool {
	id := rhino.RhinoStringPool_LookupID(name)
	return id != 0 && s.HasSlotByID(id)
}

func (s *Scope) HasSlotByID(id uint32) bool {
	bloom := scopeBloomFromID(id)
	implicit := AbstractScope_ImplicitVar_OfID(id)
	for cur := s; cur != nil; cur = cur.lookupParent {
		if cur.ownIndex(id, bloom) >= 0 {
			return true
		}
		if implicit != nil && cur.hasImplicitVars && implicit.IsMadeByScope(cur) {
			return true
		}
	}
	return false
}

func (s *Scope) inheritedInput() (*CompilerInput, bool) {
	for anc := s.parent; anc != nil && anc.parent != nil; anc = anc.parent {
		for _, v := range anc.varList {
			if !v.IsImplicitGoogNamespace() {
				return v.GetInput(), true
			}
		}
	}
	return nil, false
}

// GetVar uses the concrete parent links of Java's final Scope class. The
// generic base must support TypedScope's virtual methods; untyped scopes can
// avoid those dictionary calls and classify implicit names once per lookup.
func (s *Scope) GetVar(name string) *Var {
	id := rhino.RhinoStringPool_LookupID(name)
	if id == 0 {
		return nil
	}
	return s.GetVarByID(id)
}

func (s *Scope) GetVarByNode(n *rhino.Node) *Var {
	return s.GetVarByID(n.GetStringID())
}

func (s *Scope) GetVarByID(id uint32) *Var {
	if id == 0 {
		return nil
	}
	bloom := scopeBloomFromID(id)
	implicit := AbstractScope_ImplicitVar_OfID(id)
	for cur := s; cur != nil; cur = cur.lookupParent {
		if v := cur.ownVar(id, bloom); v != nil {
			return v
		}
		if implicit != nil && cur.hasImplicitVars {
			if v := cur.GetOwnImplicitSlot(implicit); v != nil {
				return v
			}
		}
	}
	return nil
}

// Java: AbstractScope.getVarCount()
func (s *Scope) GetVarCount() int {
	return len(s.varList)
}

// Java: AbstractScope.getVarIterable()
func (s *Scope) GetVarIterable() iter.Seq[*Var] {
	return func(yield func(*Var) bool) {
		for i := 0; i < len(s.varList); i++ {
			if !yield(s.varList[i]) {
				return
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
	clear(s.varList)
	s.varList = s.varList[:0]
	s.varIndex = nil
	s.varBloom = 0
}

// Java: AbstractScope.undeclare(V)
func (s *Scope) Undeclare(var_ *Var) {
	jutil.CheckState(var_.scope == s)
	current := s.lookupOwnVarByID(var_.nameID)
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
	id := var_.nameID
	i := s.ownIndex(id, scopeBloomFromID(id))
	if i < 0 {
		return
	}
	s.varList = slices.Delete(s.varList, i, i+1)
	if s.varIndex != nil {
		delete(s.varIndex, id)
		for j := i; j < len(s.varList); j++ {
			s.varIndex[s.varList[j].nameID] = int32(j)
		}
	}
	if len(s.varList) == 0 {
		s.varBloom = 0
	}
}

// addVar appends var_, whose name is not declared in this scope.
func (s *Scope) addVar(var_ *Var) {
	s.varsInitialized = true
	s.varBloom |= scopeBloomFromID(var_.nameID)
	s.varList = append(s.varList, var_)
	if s.varIndex != nil {
		s.varIndex[var_.nameID] = int32(len(s.varList) - 1)
	} else if len(s.varList) > scopeSmallVars {
		s.varIndex = make(map[uint32]int32, len(s.varList))
		for i, v := range s.varList {
			s.varIndex[v.nameID] = int32(i)
		}
	}
}

// Java: AbstractScope.declareInternal(String, V)
func (s *Scope) DeclareInternal(name string, var_ *Var) {
	id := rhino.RhinoStringPool_InternID(name)
	if s.IsFunctionBlockScope() && s.parent.hasOwnSlotByID(id) && !s.hasOwnSlotByID(id) && !s.isBleedingFunctionNameID(id) {
		panic(jutil.NewIllegalStateExceptionMsg(jutil.LenientFormat("Illegal shadow: %s", var_.GetNode())))
	}
	if i := s.ownIndex(id, scopeBloomFromID(id)); i >= 0 {
		s.varsInitialized = true
		s.varList[i] = var_
		return
	}
	s.addVar(var_)
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
		s.implicitVars = new([4]*Var)
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
	id := rhino.RhinoStringPool_InternID(name)
	// Make sure that it's declared only once
	if s.getOwnSlotByID(id) != nil {
		panic(jutil.NewIllegalStateException())
	}
	var_ := newVarWithID(id, nameNode, s, len(s.varList), input, nil)
	if s.IsFunctionBlockScope() && s.parent.hasOwnSlotByID(id) && !s.isBleedingFunctionNameID(id) {
		panic(jutil.NewIllegalStateExceptionMsg(jutil.LenientFormat("Illegal shadow: %s", var_.GetNode())))
	}
	s.addVar(var_)
	return var_
}

// getOwnSlotByID is getOwnSlot for an interned name.
func (s *Scope) getOwnSlotByID(id uint32) *Var {
	if v := s.lookupOwnVarByID(id); v != nil {
		return v
	}
	if s.hasImplicitVars {
		if imp := AbstractScope_ImplicitVar_OfID(id); imp != nil {
			return s.GetOwnImplicitSlot(imp)
		}
	}
	return nil
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
	id := rhino.RhinoStringPool_InternID(name)
	var_ := s.getOwnSlotByID(id)
	if var_ == nil {
		var_ = newVarWithID(id, nil, s, -1, nil, definition)
		s.addVar(var_)
	} else if var_.IsImplicitGoogNamespace() {
		var_.AddImplicitGoogNamespaceDefinition(definition)
	}
	return var_
}

// Java: Scope.makeImplicitVar(ImplicitVar)
func (s *Scope) MakeImplicitVar(var_ AbstractScope_ImplicitVar) *Var {
	return newVarWithID(var_.nameID(), nil, s, -1, nil, nil)
}
