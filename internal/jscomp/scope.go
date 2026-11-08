// Port of com/google/javascript/jscomp/Scope.java @ closure-compiler v20261104

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

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
	"github.com/paullepoulpe/jsgo/internal/strpool"
)

// Scope contains information about a variable scope in JavaScript. Scopes
// can be nested, a scope points back to its parent scope. A Scope contains
// information about variables defined in that scope.
//
// GetVar returns *Var, GetParent and GetClosestHoistScope return *Scope,
// GetVarIterable is an iter.Seq[*Var]. GetSlot, GetOwnSlot and
// GetParentScope keep rhino.StaticScope's signatures; see the note at the
// top of abstract_scope.go.
//
// Go-only: performance — the Java Map<String, Var> is varList, in
// declaration order, plus scopeVarIndex once the scope has more than
// scopeSmallVars vars. lookupParent is the closest ancestor worth searching:
// scanned scopes without vars or implicit vars are skipped.
//
// See NodeTraversal.
type Scope struct {
	scopeCommon
	inferFrameIdx int32

	implicitVars *[4]*Var
	parent       *Scope
	lookupParent *Scope
	varList      []*Var
	scopeVarIndex
}

var _ AbstractScope = (*Scope)(nil)

func varNameID(v *Var) uint32 { return v.nameID }
func varIDHash(v *Var) uint32 { return scopeIDHash(v.nameID) }

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
	s := &Scope{}
	s.initScopeCommon(rootNode)
	checkChildScope(&s.scopeCommon, parent)
	return initScopeWithParent(s, parent, rootNode)
}

// initScopeWithParent is newScopeWithParent on a zero Scope (one a
// SyntacticScopeCreator's slab handed out).
func initScopeWithParent(s *Scope, parent *Scope, rootNode rhino.Node) *Scope {
	s.inferFrameIdx = -1
	s.parent = parent
	s.rootNode = rootNode
	s.depth = parent.depth + 1
	s.initLookupChain(parent)
	return s
}

// Java: Scope(Node)
func newScope(rootNode rhino.Node) *Scope {
	s := &Scope{inferFrameIdx: -1}
	s.initScopeCommon(rootNode)
	s.CheckRootScope()
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

func (s *Scope) parentScope() *Scope       { return s.parent }
func (s *Scope) lookupParentScope() *Scope { return s.lookupParent }

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

// Java: AbstractScope.getGlobalScope()
func (s *Scope) GetGlobalScope() *Scope {
	return scopeGlobal(s)
}

// Java: AbstractScope.getClosestHoistScope()
func (s *Scope) GetClosestHoistScope() *Scope {
	return scopeClosestHoist(s)
}

// Java: AbstractScope.getClosestContainerScope()
func (s *Scope) GetClosestContainerScope() *Scope {
	return scopeClosestContainer(s)
}

func (s *Scope) isBleedingFunctionNameID(id uint32) bool {
	var_ := s.GetVarByID(id)
	return var_ != nil && isBleedingFunctionNode(var_.nameNode)
}

func (s *Scope) ownIndex(id, h uint32) int {
	if s.varIndex == nil {
		for i, v := range s.varList {
			if v.nameID == id {
				return i
			}
		}
		return -1
	}
	table := *s.varIndex
	mask := uint32(len(table) - 1)
	for slot := h & mask; ; slot = (slot + 1) & mask {
		e := table[slot]
		if e == 0 {
			return -1
		}
		if idx := int(e&scopeIndexPosMask) - 1; s.varList[idx].nameID == id {
			return idx
		}
	}
}

func (s *Scope) ownVar(id, h uint32, bloom uint64) *Var {
	if s.varBloom&bloom != bloom {
		return nil
	}
	if i := s.ownIndex(id, h); i >= 0 {
		return s.varList[i]
	}
	return nil
}

func (s *Scope) lookupOwnVarByID(id uint32) *Var {
	if len(s.varList) == 0 || id == 0 {
		return nil
	}
	h := scopeIDHash(id)
	return s.ownVar(id, h, scopeBloomFromHash(h))
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
	id := strpool.Lookup(name)
	if id == 0 {
		return nil
	}
	return s.getOwnSlotByID(uint32(id))
}

// Java: AbstractScope.hasOwnSlot(String)
func (s *Scope) HasOwnSlot(name string) bool {
	id := strpool.Lookup(name)
	return id != 0 && s.hasOwnSlotByID(uint32(id))
}

// Java: AbstractScope.hasSlot(String)
func (s *Scope) HasSlot(name string) bool {
	id := strpool.Lookup(name)
	return id != 0 && s.HasSlotByID(uint32(id))
}

func (s *Scope) HasSlotByID(id uint32) bool {
	h := scopeIDHash(id)
	bloom := scopeBloomFromHash(h)
	imp, isImplicit := AbstractScope_ImplicitVar_LookupID(id)
	for cur := s; cur != nil; cur = cur.lookupParent {
		if cur.ownVar(id, h, bloom) != nil || (isImplicit && cur.hasImplicitVars && cur.hasOwnImplicitVar(imp)) {
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
	id := strpool.Lookup(name)
	if id == 0 {
		return nil
	}
	return s.GetVarByID(uint32(id))
}

func (s *Scope) GetVarByNode(n rhino.Node) *Var {
	return s.GetVarByID(n.GetStringID())
}

func (s *Scope) GetVarByID(id uint32) *Var {
	if id == 0 {
		return nil
	}
	h := scopeIDHash(id)
	bloom := scopeBloomFromHash(h)
	imp, isImplicit := AbstractScope_ImplicitVar_LookupID(id)
	for cur := s; cur != nil; cur = cur.lookupParent {
		if v := cur.ownVar(id, h, bloom); v != nil {
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
	return scopeVarSeq(&s.varList)
}

// Java: AbstractScope.getAllSymbols()
func (s *Scope) GetAllSymbols() iter.Seq[*Var] {
	return s.GetVarIterable()
}

// Java: AbstractScope.undeclare(V)
func (s *Scope) Undeclare(var_ *Var) {
	jutil.CheckState(var_.scope == s)
	h := scopeIDHash(var_.nameID)
	pos := -1
	if bloom := scopeBloomFromHash(h); s.varBloom&bloom == bloom {
		pos = s.ownIndex(var_.nameID, h)
	}
	scopeUndeclare(var_, &s.scopeCommon, &s.scopeVarIndex, &s.varList, pos, varIDHash)
}

// addVar appends var_, whose name is not declared in this scope.
func (s *Scope) addVar(var_ *Var) {
	h := scopeIDHash(var_.nameID)
	scopeAppendVar(&s.scopeCommon, &s.scopeVarIndex, &s.varList, var_, h, scopeBloomFromHash(h), varIDHash)
}

func (s *Scope) hasOwnImplicitVar(imp AbstractScope_ImplicitVar) bool {
	return imp.isMadeByRootNode(s.rootNode)
}

func (s *Scope) getOwnImplicitVar(imp AbstractScope_ImplicitVar) *Var {
	if !s.hasOwnImplicitVar(imp) {
		return nil
	}
	return scopeGetImplicitVar(s, &s.implicitVars, imp, (*Scope).MakeImplicitVar)
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
	id := uint32(strpool.Intern(name))
	// Make sure that it's declared only once
	if s.getOwnSlotByID(id) != nil {
		panic(jutil.NewIllegalStateException())
	}
	var_ := newVarWithID(id, nameNode, s, len(s.varList), input, rhino.NilNode)
	if s.IsFunctionBlockScope() && s.parent.hasOwnSlotByID(id) && !s.isBleedingFunctionNameID(id) {
		failIllegalShadow(var_.nameNode)
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
	id := uint32(strpool.Intern(name))
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

// Java: AbstractScope.contains(S)
func (s *Scope) Contains(other *Scope) bool {
	return scopeContains(s, other)
}

// Java: AbstractScope.getArgumentsVar()
func (s *Scope) GetArgumentsVar() *Var {
	return scopeArgumentsVar(s, (*Scope).getOwnImplicitVar)
}

// Java: AbstractScope.getAllAccessibleVariables()
func (s *Scope) GetAllAccessibleVariables() iter.Seq[*Var] {
	return scopeAllAccessibleVars(
		s,
		func(cur *Scope) []*Var { return cur.varList },
		varNameID,
	)
}

// Java: AbstractScope.getClosestCfgRootScope()
func (s *Scope) GetClosestCfgRootScope() *Scope {
	return scopeClosestCfgRoot(s)
}

// Java: AbstractScope.getCommonParent(S)
func (s *Scope) GetCommonParent(other *Scope) *Scope {
	return scopeCommonParent(s, other)
}

// Java: AbstractScope.hasSameContainerScope(S)
func (s *Scope) HasSameContainerScope(other *Scope) bool {
	return scopeHasSameContainer(s, other)
}

// Java: AbstractScope.getScopeOfThis()
func (s *Scope) GetScopeOfThis() *Scope {
	return scopeOfThis(s)
}

func (s *Scope) GetTopmostScopeOfEventualDeclaration(name string) rhino.StaticScope {
	return rhino.StaticScope_GetTopmostScopeOfEventualDeclaration(s, name)
}
