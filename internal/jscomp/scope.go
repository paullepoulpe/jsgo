// Port of com/google/javascript/jscomp/Scope.java @ closure-compiler v20261021

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
	"hash/maphash"
	"iter"
	"slices"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

var scopeHashSeed = maphash.MakeSeed()

func scopeNameHash(s string) uint32 {
	if s == "" {
		return 0
	}
	return uint32(maphash.String(scopeHashSeed, s))
}

func scopeBloomFromHash(h uint32) uint64 {
	return (uint64(1) << (h & 63)) | (uint64(1) << ((h >> 6) & 63))
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
// GetVar returns *Var, GetParent and GetClosestHoistScope return *Scope,
// GetVarIterable is an iter.Seq[*Var]. GetSlot, GetOwnSlot and
// GetParentScope keep rhino.StaticScope's signatures; see the note at the
// top of abstract_scope.go.
//
// Go-only: performance — the Java LinkedHashMap<String, Var> is varList, in
// declaration order, plus varIndex (name ID to varList index, stored behind a
// pointer to save 16 bytes on small scopes) once the scope has more than
// scopeSmallVars vars, plus varBloom, a two-bit-per-name filter that lets
// lookups skip scopes without the name. lookupParent is the closest ancestor
// worth searching: scanned scopes without vars or implicit vars are skipped.
//
// See NodeTraversal.
type Scope struct {
	scopeCommon

	implicitVars *[4]*Var
	parent       *Scope
	lookupParent *Scope
	varList      []*Var
	varIndex     *[]int32 // see ownIndex; pointer saves 16 bytes on small scopes
	varBloom     uint64
}

var _ AbstractScope = (*Scope)(nil)

// Java: Scope.createGlobalScope(Node)
func Scope_CreateGlobalScope(rootNode rhino.Node) *Scope {
	return newScope(rootNode)
}

// Java: Scope.createChildScope(Scope, Node)
func Scope_CreateChildScope(parent *Scope, rootNode rhino.Node) *Scope {
	return newScopeWithParent(parent, rootNode)
}

// Java: Scope(Scope, Node)
func newScopeWithParent(parent *Scope, rootNode rhino.Node) *Scope {
	if rootNode.IsNil() || parent == nil {
		panic(jutil.NewNullPointerException())
	}
	if !NodeUtil_CreatesScope(rootNode) {
		panic(jutil.NewIllegalArgumentExceptionMsg(jutil.LenientFormat("%s", rootNode)))
	}
	if rootNode == parent.rootNode {
		panic(jutil.NewIllegalArgumentExceptionMsg(
			jutil.LenientFormat("rootNode should not be the parent's root node: %s", rootNode)))
	}
	return initScopeWithParent(new(Scope), parent, rootNode)
}

// initScopeWithParent is newScopeWithParent on a zero Scope (one a
// SyntacticScopeCreator's slab handed out).
func initScopeWithParent(s *Scope, parent *Scope, rootNode rhino.Node) *Scope {
	s.parent = parent
	s.rootNode = rootNode
	s.depth = parent.depth + 1
	s.initLookupChain(parent)
	return s
}

// Java: Scope(Node)
func newScope(rootNode rhino.Node) *Scope {
	if rootNode.IsNil() {
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

// Java: AbstractScope.checkChildScope(S)
func (s *Scope) CheckChildScope(parent *Scope) {
	if parent == nil {
		panic(jutil.NewNullPointerException())
	}
	s.checkChildScopeRoot(parent.rootNode)
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
		!var_.nameNode.IsNil() &&
		var_.nameNode.GetParent().IsFunction()
}

// ownIndex returns the varList index of the var named id declared in this
// scope (implicit vars excluded), or -1.
func (s *Scope) ownIndex(id uint32, bloom uint64) int {
	if s.varBloom&bloom != bloom {
		return -1
	}
	if s.varIndex != nil {
		// Go-only: performance — an open-addressing table of varList
		// positions + 1 (0 = empty), pointer-free and at most half full.
		index := *s.varIndex
		mask := uint32(len(index) - 1)
		for slot := scopeIDHash(id) & mask; ; slot = (slot + 1) & mask {
			e := index[slot]
			if e == 0 {
				return -1
			}
			if s.varList[e-1].nameID == id {
				return int(e - 1)
			}
		}
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
		if imp, ok := AbstractScope_ImplicitVar_LookupID(id); ok {
			return s.hasOwnImplicitVar(imp)
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
	return s.getOwnSlotByID(id)
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
	imp, isImplicit := AbstractScope_ImplicitVar_LookupID(id)
	for cur := s; cur != nil; cur = cur.lookupParent {
		if cur.ownIndex(id, bloom) >= 0 {
			return true
		}
		if isImplicit && cur.hasImplicitVars && imp.isMadeByRootNode(cur.rootNode) {
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

// GetVar looks up name in s and its ancestors.
func (s *Scope) GetVar(name string) *Var {
	id := rhino.RhinoStringPool_LookupID(name)
	if id == 0 {
		return nil
	}
	return s.GetVarByID(id)
}

func (s *Scope) GetVarByNode(n rhino.Node) *Var {
	return s.GetVarByID(n.GetStringID())
}

func (s *Scope) GetVarByID(id uint32) *Var {
	if id == 0 {
		return nil
	}
	bloom := scopeBloomFromID(id)
	imp, isImplicit := AbstractScope_ImplicitVar_LookupID(id)
	for cur := s; cur != nil; cur = cur.lookupParent {
		if v := cur.ownVar(id, bloom); v != nil {
			return v
		}
		if isImplicit && cur.hasImplicitVars {
			if v := cur.getOwnImplicitVar(imp); v != nil {
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
		s.rebuildVarIndex()
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
	if s.varIndex != nil && len(s.varList)*2 <= len(*s.varIndex) {
		s.indexVar(len(s.varList) - 1)
	} else if len(s.varList) > scopeSmallVars {
		s.rebuildVarIndex()
	}
}

// indexVar adds varList[i] to varIndex.
func (s *Scope) indexVar(i int) {
	index := *s.varIndex
	mask := uint32(len(index) - 1)
	slot := scopeIDHash(s.varList[i].nameID) & mask
	for index[slot] != 0 {
		slot = (slot + 1) & mask
	}
	index[slot] = int32(i + 1)
}

func (s *Scope) rebuildVarIndex() {
	if len(s.varList) <= scopeSmallVars {
		s.varIndex = nil
		return
	}
	n := 32
	for n < 4*len(s.varList) {
		n *= 2
	}
	index := make([]int32, n)
	s.varIndex = &index
	for i := range s.varList {
		s.indexVar(i)
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

func (s *Scope) hasOwnImplicitVar(imp AbstractScope_ImplicitVar) bool {
	return imp.isMadeByRootNode(s.rootNode)
}

func (s *Scope) getOwnImplicitVar(imp AbstractScope_ImplicitVar) *Var {
	if !imp.isMadeByRootNode(s.rootNode) {
		return nil
	}
	if s.implicitVars == nil {
		s.implicitVars = new([4]*Var)
	}
	if v := s.implicitVars[imp]; v != nil {
		return v
	}
	v := s.MakeImplicitVar(imp)
	if v != nil {
		s.implicitVars[imp] = v
	}
	return v
}

// Java: AbstractScope.hasOwnImplicitSlot(ImplicitVar)
func (s *Scope) HasOwnImplicitSlot(name *AbstractScope_ImplicitVar) bool {
	return name != nil && s.hasOwnImplicitVar(*name)
}

// Java: AbstractScope.getOwnImplicitSlot(ImplicitVar)
func (s *Scope) GetOwnImplicitSlot(name *AbstractScope_ImplicitVar) *Var {
	if name == nil {
		return nil
	}
	return s.getOwnImplicitVar(*name)
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
func (s *Scope) Declare(name string, nameNode rhino.Node, input *CompilerInput) *Var {
	jutil.CheckArgument(name != "")
	id := rhino.RhinoStringPool_InternID(name)
	// Make sure that it's declared only once
	if s.getOwnSlotByID(id) != nil {
		panic(jutil.NewIllegalStateException())
	}
	var_ := newVarWithID(id, nameNode, s, len(s.varList), input, rhino.NilNode)
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
		if imp, ok := AbstractScope_ImplicitVar_LookupID(id); ok {
			return s.getOwnImplicitVar(imp)
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
func (s *Scope) DeclareImplicitGoogNamespaceIfAbsent(name string, definition rhino.Node) *Var {
	jutil.CheckArgument(name != "")
	jutil.CheckStateMsg1(s.IsGlobal(), "Cannot declare implicit goog namespace in local scope %s", s)
	// Allow redeclarations of provides, since they are implicit and don't have a single
	// declaration site.
	id := rhino.RhinoStringPool_InternID(name)
	var_ := s.getOwnSlotByID(id)
	if var_ == nil {
		var_ = newVarWithID(id, rhino.NilNode, s, -1, nil, definition)
		s.addVar(var_)
	} else if var_.IsImplicitGoogNamespace() {
		var_.AddImplicitGoogNamespaceDefinition(definition)
	}
	return var_
}

// Java: Scope.makeImplicitVar(ImplicitVar)
func (s *Scope) MakeImplicitVar(var_ AbstractScope_ImplicitVar) *Var {
	return newVarWithID(var_.nameID(), rhino.NilNode, s, -1, nil, rhino.NilNode)
}

func (s *Scope) ThisScope() *Scope { return s }

// Java: AbstractScope.contains(S)
func (s *Scope) Contains(other *Scope) bool {
	for cur := jutil.CheckNotNull(other); cur != nil; cur = cur.parent {
		if cur == s {
			return true
		}
	}
	return false
}

// Java: AbstractScope.getArgumentsVar()
func (s *Scope) GetArgumentsVar() *Var {
	for cur := s; cur != nil; cur = cur.parent {
		if v := cur.getOwnImplicitVar(AbstractScope_ImplicitVar_ARGUMENTS); v != nil {
			return v
		}
	}
	return nil
}

// Java: AbstractScope.isBleedingFunctionName(String)
func (s *Scope) IsBleedingFunctionName(name string) bool {
	id := rhino.RhinoStringPool_LookupID(name)
	return id != 0 && s.isBleedingFunctionNameID(id)
}

// Java: AbstractScope.getAllAccessibleVariables()
func (s *Scope) GetAllAccessibleVariables() iter.Seq[*Var] {
	return func(yield func(*Var) bool) {
		var seen map[uint32]struct{}
		for cur := s; cur != nil; cur = cur.parent {
			for _, v := range cur.varList {
				if _, dup := seen[v.nameID]; dup {
					continue
				}
				if cur.parent != nil || seen != nil {
					if seen == nil {
						seen = make(map[uint32]struct{}, len(cur.varList))
					}
					seen[v.nameID] = struct{}{}
				}
				if !yield(v) {
					return
				}
			}
		}
	}
}

// Java: AbstractScope.getClosestCfgRootScope()
func (s *Scope) GetClosestCfgRootScope() *Scope {
	cur := s
	for !cur.IsCfgRootScope() {
		cur = cur.parent
	}
	return cur
}

// Java: AbstractScope.getCommonParent(S)
func (s *Scope) GetCommonParent(other *Scope) *Scope {
	left, right := s, other
	for left != nil && right != nil && left != right {
		ld, rd := left.depth, right.depth
		if ld >= rd {
			left = left.parent
		}
		if ld <= rd {
			right = right.parent
		}
	}
	jutil.CheckState(left != nil && left == right)
	return left
}

// Java: AbstractScope.hasSameContainerScope(S)
func (s *Scope) HasSameContainerScope(other *Scope) bool {
	return s == other || s.GetClosestContainerScope() == other.GetClosestContainerScope()
}

// Java: AbstractScope.getScopeOfThis()
func (s *Scope) GetScopeOfThis() *Scope {
	scope := s.GetClosestContainerScope()
	for scope.parent != nil && !NodeUtil_IsNonArrowFunction(scope.rootNode) {
		scope = scope.parent.GetClosestContainerScope()
	}
	return scope
}

func (s *Scope) getSlotV(name string) *Var { return s.GetVar(name) }

func (s *Scope) GetTopmostScopeOfEventualDeclaration(name string) rhino.StaticScope {
	return rhino.StaticScope_GetTopmostScopeOfEventualDeclaration(s, name)
}

func (s *Scope) renameOwnVar(oldName, newName string) {
	if s == nil {
		return
	}
	oldID := rhino.RhinoStringPool_LookupID(oldName)
	if oldID == 0 {
		return
	}
	idx := s.ownIndex(oldID, scopeBloomFromID(oldID))
	if idx < 0 {
		return
	}
	newID := rhino.RhinoStringPool_InternID(newName)
	s.varList[idx].nameID = newID
	var bloom uint64
	for _, v := range s.varList {
		bloom |= scopeBloomFromID(v.nameID)
	}
	s.varBloom = bloom
	if s.varIndex != nil {
		s.rebuildVarIndex()
	}
}
