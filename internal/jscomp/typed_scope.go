// Port of com/google/javascript/jscomp/TypedScope.java @ closure-compiler v20261021

/*
 * Copyright 2015 The Closure Compiler Authors.
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
	"strings"
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

var nextTypedScopeID atomic.Uint32

func allocTypedScopeID() uint32 {
	id := nextTypedScopeID.Add(1)
	if id == 0 {
		id = nextTypedScopeID.Add(1)
	}
	return id
}

// TypedScope contains information about variables and their types. Scopes
// can be nested; a scope points back to its parent scope.
//
// TypedScope is also used as a lattice element for flow-sensitive type
// inference. As a lattice element, a scope is viewed as a map from names to
// types. A name not in the map is considered to have the bottom type. The
// join of two maps m1 and m2 is the map of the union of names with
// JSType.getLeastSupertype to meet the m1 type and m2 type.
//
// TypedScopes also have a concept of a set of 'reserved' names. These are
// names that will be declared as actual TypedVars later in the compilation.
// Looking up any of these reserved names will always return null, even if
// the name is available in a parent scope.
//
// See NodeTraversal and DataFlowAnalysis.
//
// It is a rhino.StaticTypedScope: GetSlot, GetOwnSlot and GetParentScope
// keep rhino.StaticScope's signatures (assert `.(*TypedVar)`); see the note
// at the top of abstract_scope.go.
type TypedScope struct {
	scopeCommon
	scopeID uint32

	implicitVars *[4]*TypedVar
	parent       *TypedScope
	lookupParent *TypedScope

	// Go-only: performance — the Java `module` and `reservedNames` fields,
	// which only the scope builder reads; dropped by ValidateCompletelyBuilt
	// (nil afterwards, and for scopes built without either).
	build *typedScopeBuildState

	// Go-only: performance — the Java LinkedHashMap<String, TypedVar> is vars,
	// in declaration order (a var is keyed by its name), plus varIndex once
	// there are more than typedScopeSmallVars vars, plus varBloom, a filter
	// that lets lookups skip scopes without the name.
	vars     []*TypedVar
	varIndex *[]uint32 // a pointer: most scopes have none, and it saves 16 bytes in each
	varBloom uint64
}

var (
	_ AbstractScope          = (*TypedScope)(nil)
	_ rhino.StaticTypedScope = (*TypedScope)(nil)
)

// typedScopeBuildState is the part of a TypedScope that TypedScopeCreator
// needs while it builds the scope.
type typedScopeBuildState struct {
	module *Module

	// Will shrink over time; nil/empty slice is the ImmutableSet.of() sentinel.
	reservedNames []string
}

func (s *TypedScope) reservedNames() []string {
	if s.build == nil {
		return nil
	}
	return s.build.reservedNames
}

// Equals implements jutil.Equaler using pointer identity (matching Java Object.equals).
func (s *TypedScope) Equals(other any) bool {
	that, ok := other.(*TypedScope)
	return ok && s == that
}

// HashCode implements jutil.Hasher using scopeID (matching Java Object.hashCode).
func (s *TypedScope) HashCode() int32 {
	if s == nil {
		return 0
	}
	h := s.scopeID * 0x9e3779b1
	return int32(h ^ (h >> 16))
}

// Java: AbstractScope.checkChildScope(S)
func (s *TypedScope) CheckChildScope(parent *TypedScope) {
	if parent == nil {
		panic(jutil.NewNullPointerException())
	}
	s.checkChildScopeRoot(parent.rootNode)
}

// Java: TypedScope(TypedScope, Node)
func NewTypedScope(parent *TypedScope, rootNode rhino.Node) *TypedScope {
	return newTypedScopeOwnedReservedNames(parent, rootNode, nil, nil)
}

// NewTypedScopeWithReservedNames creates a new global or local TypedScope.
//
// Java: TypedScope(TypedScope, Node, Set<String>, Module)
func NewTypedScopeWithReservedNames(
	parent *TypedScope, rootNode rhino.Node, reservedNames *jutil.LinkedHashSet[string], module *Module) *TypedScope {
	var copied []string
	if !jutil.CheckNotNull(reservedNames).IsEmpty() { // Java: NPE on a null set
		copied = reservedNames.ToSlice()
	}
	return newTypedScopeOwnedReservedNames(parent, rootNode, copied, module)
}

func newTypedScopeOwnedReservedNames(
	parent *TypedScope, rootNode rhino.Node, reservedNames []string, module *Module) *TypedScope {
	s := &TypedScope{}
	s.initScopeCommon(rootNode)
	s.scopeID = allocTypedScopeID()
	if len(reservedNames) > 0 {
		s.vars = make([]*TypedVar, 0, len(reservedNames))
	}
	s.CheckChildScope(parent)
	s.parent = parent
	s.depth = parent.depth + 1
	if len(reservedNames) > 0 || module != nil {
		s.build = &typedScopeBuildState{module: module, reservedNames: reservedNames}
	}
	s.initLookupChain(parent)
	return s
}

func (s *TypedScope) initLookupChain(parent *TypedScope) {
	if parent == nil {
		s.hasImplicitVars = false
		s.lookupParent = nil
		return
	}
	switch s.rootNode.GetToken() {
	case rhino.Token_FUNCTION:
		s.hasImplicitVars = NodeUtil_IsNonArrowFunction(s.rootNode)
	case rhino.Token_MEMBER_FIELD_DEF, rhino.Token_COMPUTED_FIELD_DEF:
		s.hasImplicitVars = true
	case rhino.Token_BLOCK:
		s.hasImplicitVars = s.IsStaticBlockScope()
	default:
		s.hasImplicitVars = false
	}
	if parent.scanDone && len(parent.vars) == 0 && len(parent.reservedNames()) == 0 && !parent.hasImplicitVars && !parent.IsHoistScope() {
		s.lookupParent = parent.lookupParent
	} else {
		s.lookupParent = parent
	}
}

// Java: TypedScope(Node, boolean)
func newTypedScopeWithIsBottom(rootNode rhino.Node, isBottom bool) *TypedScope {
	s := &TypedScope{}
	s.initScopeCommon(rootNode)
	s.scopeID = allocTypedScopeID()
	s.CheckRootScope()
	s.isBottom = isBottom
	s.initLookupChain(nil)
	return s
}

func newGlobalTypedScopeWithCapacity(rootNode rhino.Node, varCap int) *TypedScope {
	s := &TypedScope{}
	s.initScopeCommon(rootNode)
	s.scopeID = allocTypedScopeID()
	if varCap > 0 {
		s.vars = make([]*TypedVar, 0, varCap)
	}
	s.CheckRootScope()
	s.initLookupChain(nil)
	return s
}

// Java: TypedScope.createGlobalScope(Node)
func TypedScope_CreateGlobalScope(rootNode rhino.Node) *TypedScope {
	return newTypedScopeWithIsBottom(rootNode, false)
}

// Java: TypedScope.createLatticeBottom(Node)
func TypedScope_CreateLatticeBottom(rootNode rhino.Node) *TypedScope {
	return newTypedScopeWithIsBottom(rootNode, true)
}

// Java: TypedScope.typed()
func (s *TypedScope) Typed() *TypedScope {
	return s
}

// ValidateCompletelyBuilt asserts that all reserved names in the scope have
// been actually declared.
//
// Java: TypedScope.validateCompletelyBuilt()
func (s *TypedScope) ValidateCompletelyBuilt() {
	if reserved := s.reservedNames(); len(reserved) > 0 {
		jutil.CheckStateMsg2(
			false,
			"Expected %s to have no reserved names, found: %s. This probably indicates a bug in"+
				" TypedScopeCreator where it is failing to declare a variable.",
			s,
			"["+strings.Join(reserved, ", ")+"]")
	}
	s.build = nil
}

// IsBottom reports whether this is the bottom of the lattice.
//
// Java: TypedScope.isBottom()
func (s *TypedScope) IsBottom() bool {
	return s.isBottom
}

// Java: TypedScope.getModule()
func (s *TypedScope) GetModule() *Module {
	if s.build == nil {
		return nil
	}
	return s.build.module
}

// Java: TypedScope.getParent()
func (s *TypedScope) GetParent() *TypedScope {
	return s.parent
}

// Java: AbstractScope.getParentScope()
func (s *TypedScope) GetParentScope() rhino.StaticScope {
	if s.parent != nil {
		return s.parent
	}
	return nil
}

// Java: AbstractScope.isGlobal()
func (s *TypedScope) IsGlobal() bool {
	return s.parent == nil
}

// Java: AbstractScope.isLocal()
func (s *TypedScope) IsLocal() bool {
	return s.parent != nil
}

// GetTypeOfThis gets the type of `this` in the current scope.
//
// Java: TypedScope.getTypeOfThis()
func (s *TypedScope) GetTypeOfThis() rhino.JSType {
	root := s.rootNode
	if s.parent == nil {
		if t := rhino.ObjectType_Cast(root.GetJSType()); t != nil {
			return t
		}
		return nil
	} else if NodeUtil_IsNonArrowFunction(root) {
		nodeType := root.GetJSType()
		if nodeType != nil && nodeType.IsFunctionType() {
			return nodeType.ToMaybeFunctionType().GetTypeOfThis()
		}
		// Executed when the current scope has not been typechecked.
		return nil
	} else if s.IsStaticBlockScope() {
		return s.parent.rootNode.GetJSType()
	} else if s.IsMemberFieldDefScope() || s.IsComputedFieldDefRhsScope() {
		classType := s.parent.rootNode.GetJSType()
		if root.IsStaticMember() {
			return classType
		}
		if t := classType.AssertFunctionType().GetInstanceType(); t != nil {
			return t
		}
		return nil
	}
	return s.parent.GetTypeOfThis()
}

// Java: TypedScope.declare(String, Node, JSType, CompilerInput, boolean)
func (s *TypedScope) Declare(name string, nameNode rhino.Node,
	type_ rhino.JSType, input *CompilerInput, inferred bool) *TypedVar {
	var_, _ := s.declareWithReserved(name, nameNode, type_, input, inferred)
	return var_
}

func (s *TypedScope) declareWithReserved(name string, nameNode rhino.Node,
	type_ rhino.JSType, input *CompilerInput, inferred bool) (*TypedVar, bool) {
	jutil.CheckState(name != "") // Java: name != null && !name.isEmpty()
	wasReserved := false
	if reserved := s.reservedNames(); len(reserved) > 0 {
		for i, rn := range reserved {
			if rn == name {
				wasReserved = true
				copy(reserved[i:], reserved[i+1:])
				reserved[len(reserved)-1] = ""
				s.build.reservedNames = reserved[:len(reserved)-1]
				break
			}
		}
	}
	var_ := NewTypedVar(inferred, name, nameNode, type_, s, len(s.vars), input)
	s.DeclareInternal(name, var_)
	return var_, wasReserved
}

// Java: AbstractScope.undeclare(V)
func (s *TypedScope) Undeclare(var_ *TypedVar) {
	jutil.CheckState(var_.scope == s)
	current := s.getOwnTypedSlot(var_.name)
	if current == nil {
		panic(jutil.NewNullPointerException())
	}
	jutil.CheckState(current.Equals(var_))
	s.UndeclareInteral(var_)
}

// Java: AbstractScope.undeclareInteral(V)
func (s *TypedScope) UndeclareInteral(var_ *TypedVar) {
	if !s.varsInitialized {
		panic(jutil.NewUnsupportedOperationException())
	}
	name := var_.name
	if i := s.findVar(name, scopeNameHash(name)); i >= 0 {
		s.vars = slices.Delete(s.vars, i, i+1)
		if s.varIndex != nil {
			s.rebuildVarIndex()
		}
	}
}

// Java: AbstractScope.declareInternal(String, V)
func (s *TypedScope) DeclareInternal(name string, var_ *TypedVar) {
	if s.IsFunctionBlockScope() && s.parent.HasOwnSlot(name) && !s.HasOwnSlot(name) && !s.IsBleedingFunctionName(name) {
		panic(jutil.NewIllegalStateExceptionMsg(jutil.LenientFormat("Illegal shadow: %s", var_.GetNode())))
	}
	jutil.CheckState(name == var_.name)
	h := scopeNameHash(name)
	s.varBloom |= scopeBloomFromHash(h)
	s.varsInitialized = true
	if i := s.findVar(name, h); i >= 0 {
		s.vars[i] = var_
		return
	}
	s.vars = append(s.vars, var_)
	if s.varIndex != nil && len(s.vars)*2 <= len(*s.varIndex) {
		s.indexVar(len(s.vars)-1, h)
	} else if len(s.vars) > typedScopeSmallVars {
		s.rebuildVarIndex()
	}
}

const (
	typedScopeSmallVars    = 8
	typedVarIndexTagMask   = 0xFF000000
	typedVarIndexPosMask   = 0x00FFFFFF
	typedVarIndexMinLength = 32
)

// findVar returns the position in vars of the var named name, whose
// scopeNameHash is h, or -1.
func (s *TypedScope) findVar(name string, h uint32) int {
	if s.varIndex == nil {
		for i, v := range s.vars {
			if v.name == name {
				return i
			}
		}
		return -1
	}
	index := *s.varIndex
	mask := uint32(len(index) - 1)
	for slot := h & mask; ; slot = (slot + 1) & mask {
		e := index[slot]
		if e == 0 {
			return -1
		}
		if e&typedVarIndexTagMask == h&typedVarIndexTagMask {
			if i := int(e&typedVarIndexPosMask) - 1; s.vars[i].name == name {
				return i
			}
		}
	}
}

// indexVar adds vars[i], whose name has scopeNameHash h, to varIndex.
func (s *TypedScope) indexVar(i int, h uint32) {
	index := *s.varIndex
	mask := uint32(len(index) - 1)
	slot := h & mask
	for index[slot] != 0 {
		slot = (slot + 1) & mask
	}
	index[slot] = h&typedVarIndexTagMask | uint32(i+1)
}

func (s *TypedScope) rebuildVarIndex() {
	if len(s.vars) <= typedScopeSmallVars {
		s.varIndex = nil
		return
	}
	n := typedVarIndexMinLength
	for n < 4*len(s.vars) {
		n *= 2
	}
	index := make([]uint32, n)
	s.varIndex = &index
	for i, v := range s.vars {
		s.indexVar(i, scopeNameHash(v.name))
	}
}

// Java: AbstractScope.clearVarsInternal()
func (s *TypedScope) ClearVarsInternal() {
	s.varBloom = 0
	clear(s.vars)
	s.vars = s.vars[:0]
	s.varIndex = nil
}

// Java: AbstractScope.getVarIterable()
func (s *TypedScope) GetVarIterable() iter.Seq[*TypedVar] {
	return func(yield func(*TypedVar) bool) {
		for i := 0; i < len(s.vars); i++ {
			if !yield(s.vars[i]) {
				return
			}
		}
	}
}

// Java: AbstractScope.getAllSymbols()
func (s *TypedScope) GetAllSymbols() iter.Seq[*TypedVar] {
	return s.GetVarIterable()
}

// Java: AbstractScope.getVarCount()
func (s *TypedScope) GetVarCount() int {
	return len(s.vars)
}

func (s *TypedScope) getOwnTypedSlot(name string) *TypedVar {
	return s.ownVar(name, scopeNameHash(name))
}

// ownVar is getOwnTypedSlot for a name whose scopeNameHash is h.
func (s *TypedScope) ownVar(name string, h uint32) *TypedVar {
	if bloom := scopeBloomFromHash(h); s.varBloom&bloom != bloom {
		return nil
	}
	if i := s.findVar(name, h); i >= 0 {
		return s.vars[i]
	}
	return nil
}

func (s *TypedScope) hasOwnImplicitVar(imp AbstractScope_ImplicitVar) bool {
	return imp != AbstractScope_ImplicitVar_EXPORTS && imp.isMadeByRootNode(s.rootNode)
}

func (s *TypedScope) getOwnImplicitVar(imp AbstractScope_ImplicitVar) *TypedVar {
	if !s.hasOwnImplicitVar(imp) {
		return nil
	}
	if s.implicitVars == nil {
		s.implicitVars = new([4]*TypedVar)
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

func (s *TypedScope) getOwnSlotV(name string) *TypedVar {
	if v := s.getOwnTypedSlot(name); v != nil {
		return v
	}
	if imp, ok := AbstractScope_ImplicitVar_Lookup(name); ok {
		return s.getOwnImplicitVar(imp)
	}
	return nil
}

func (s *TypedScope) GetOwnSlot(name string) rhino.StaticSlot {
	if v := s.getOwnSlotV(name); v != nil {
		return v
	}
	return nil
}

func (s *TypedScope) HasOwnSlot(name string) bool {
	if s.getOwnTypedSlot(name) != nil {
		return true
	}
	if imp, ok := AbstractScope_ImplicitVar_Lookup(name); ok {
		return s.hasOwnImplicitVar(imp)
	}
	return false
}

// Java: TypedScope.makeImplicitVar(ImplicitVar)
func (s *TypedScope) MakeImplicitVar(var_ AbstractScope_ImplicitVar) *TypedVar {
	if s.parent == nil || var_ == AbstractScope_ImplicitVar_EXPORTS {
		return nil
	}
	return NewTypedVar(false, var_.name(), rhino.NilNode, s.GetImplicitVarType(var_), s, -1, nil)
}

// Java: TypedScope.hasOwnImplicitSlot(ImplicitVar)
func (s *TypedScope) HasOwnImplicitSlot(name *AbstractScope_ImplicitVar) bool {
	return name != nil && s.hasOwnImplicitVar(*name)
}

// Java: AbstractScope.getOwnImplicitSlot(ImplicitVar)
func (s *TypedScope) GetOwnImplicitSlot(name *AbstractScope_ImplicitVar) *TypedVar {
	if name == nil {
		return nil
	}
	return s.getOwnImplicitVar(*name)
}

// Java: TypedScope.getImplicitVarType(ImplicitVar)
func (s *TypedScope) GetImplicitVarType(var_ AbstractScope_ImplicitVar) rhino.JSType {
	switch var_ {
	case AbstractScope_ImplicitVar_ARGUMENTS:
		globalArgs := s.GetGlobalScope().GetVar(Var_ARGUMENTS)
		if globalArgs != nil && globalArgs.IsExtern() {
			return globalArgs.GetType()
		}
		return nil
	case AbstractScope_ImplicitVar_THIS:
		return s.GetTypeOfThis()
	case AbstractScope_ImplicitVar_SUPER:
		receiverType := rhino.ObjectType_Cast(s.GetTypeOfThis())
		if receiverType == nil {
			return nil
		} else if receiverType.IsInstanceType() {
			superclassCtor := receiverType.GetSuperClassConstructor()
			if superclassCtor == nil {
				return nil
			}
			if t := superclassCtor.GetInstanceType(); t != nil {
				return t
			}
			return nil
		} else {
			if t := receiverType.GetImplicitPrototype(); t != nil {
				return t
			}
			return nil
		}
	case AbstractScope_ImplicitVar_EXPORTS:
		panic(jutil.NewAssertionError("TypedScopes should not contain an implicit 'exports'"))
	}

	panic(jutil.NewAssertionErrorNoMessage())
}

// GetDeclarativelyUnboundVarsWithoutTypes returns the variables in this
// scope that have been declared with 'var' and not yet declared with a known
// type.
//
// Java: TypedScope.getDeclarativelyUnboundVarsWithoutTypes()
func (s *TypedScope) GetDeclarativelyUnboundVarsWithoutTypes() iter.Seq[*TypedVar] {
	return func(yield func(*TypedVar) bool) {
		for _, v := range s.vars {
			if s.IsDeclarativelyUnboundVarWithoutType(v) && !yield(v) {
				return
			}
		}
	}
}

// Java: TypedScope.isDeclarativelyUnboundVarWithoutType(TypedVar)
func (s *TypedScope) IsDeclarativelyUnboundVarWithoutType(var_ *TypedVar) bool {
	return !var_.GetParentNode().IsNil() &&
		var_.GetType() == nil &&
		// TODO(bradfordcsmith): update this for destructuring
		var_.GetParentNode().IsVar() && // NOTE: explicitly excludes let/const
		!var_.IsExtern()
}

var (
	typedScopeSymbolID    = rhino.RhinoStringPool_InternID("Symbol")
	typedScopeSymbolHash  = scopeNameHash("Symbol")
	typedScopeSymbolBloom = scopeBloomFromHash(typedScopeSymbolHash)
)

// lookupRootVar is AbstractScope.getVar, for a name without dots.
func (s *TypedScope) lookupRootVar(name string) *TypedVar {
	h := scopeNameHash(name)
	imp, isImplicit := AbstractScope_ImplicitVar_Lookup(name)
	return s.lookupRootVarPrehashed(name, h, scopeBloomFromHash(h), imp, isImplicit)
}

func (s *TypedScope) lookupRootVarPrehashed(
	name string, h uint32, bloom uint64, imp AbstractScope_ImplicitVar, isImplicit bool,
) *TypedVar {
	for cur := s; cur != nil; cur = cur.lookupParent {
		if cur.varBloom&bloom == bloom {
			if i := cur.findVar(name, h); i >= 0 {
				return cur.vars[i]
			}
		}
		if isImplicit && cur.hasImplicitVars {
			if v := cur.getOwnImplicitVar(imp); v != nil {
				return v
			}
		}
	}
	return nil
}

// GetVarByID returns the slot for an interned Rhino string ID.
func (s *TypedScope) GetVarByID(id uint32) *TypedVar {
	if id == typedScopeSymbolID {
		return s.lookupRootVarPrehashed("Symbol", typedScopeSymbolHash, typedScopeSymbolBloom, 0, false)
	}
	name := rhino.RhinoStringPool_StringFromID(id)
	if strings.IndexByte(name, '.') >= 0 {
		return s.GetVar(name)
	}
	h := scopeNameHash(name)
	imp, isImplicit := AbstractScope_ImplicitVar_LookupID(id)
	return s.lookupRootVarPrehashed(name, h, scopeBloomFromHash(h), imp, isImplicit)
}

// GetVar returns the slot for name, considering shadowing of qualified names.
//
// Java: TypedScope.getVar(String)
func (s *TypedScope) GetVar(name string) *TypedVar {
	dot := strings.IndexByte(name, '.')
	if dot < 0 {
		return s.lookupRootVar(name)
	}
	ownSlot := s.getOwnSlotV(name)
	if ownSlot != nil {
		// Micro-optimization: variables declared directly in this scope cannot have been shadowed.
		return ownSlot
	} else if s.parent == nil {
		return nil
	}
	rootVar := s.lookupRootVar(name[:dot])
	if rootVar == nil {
		// Default to the global scope because externs may have qualified names with undeclared roots.
		return s.parent.GetGlobalScope().getOwnSlotV(name)
	}
	// Qualified names 'a.b.c' are declared in the same scope as 'a', never a child scope, which
	// is why calling `getOwnSlot` is sufficient.
	return rootVar.scope.getOwnSlotV(name)
}

// Java: TypedScope.getTypeThroughNamespace(String)
func (s *TypedScope) GetTypeThroughNamespace(moduleId string) rhino.JSType {
	split := strings.LastIndexByte(moduleId, '.')
	if split >= 0 {
		parentName := moduleId[:split]
		prop := moduleId[split+1:]
		parentType := s.GetTypeThroughNamespace(parentName)
		if parentType == nil || parentType.ToMaybeObjectType() == nil {
			return nil
		}
		return parentType.AssertObjectType().GetPropertyType(prop)
	}
	var_ := s.GetVar(moduleId)
	if var_ != nil {
		return var_.GetType()
	}
	return nil
}

// Java: TypedScope.getTopmostScopeOfEventualDeclaration(String)
func (s *TypedScope) GetTopmostScopeOfEventualDeclaration(name string) rhino.StaticScope {
	for cur := s; cur != nil; cur = cur.parent {
		if cur.getOwnSlotV(name) != nil || slices.Contains(cur.reservedNames(), name) {
			return cur
		}
	}
	return nil
}

// LookupQualifiedName is the rhino.StaticTypedScope default method, which
// TypedScope inherits.
//
// Java: StaticTypedScope.lookupQualifiedName(QualifiedName) (default method)
func (s *TypedScope) LookupQualifiedName(qname rhino.QualifiedName) rhino.JSType {
	return rhino.StaticTypedScope_LookupQualifiedName(s, qname)
}

// Java: AbstractScope.contains(S)
func (s *TypedScope) Contains(other *TypedScope) bool {
	for cur := jutil.CheckNotNull(other); cur != nil; cur = cur.parent {
		if cur == s {
			return true
		}
	}
	return false
}

// Java: AbstractScope.getGlobalScope()
func (s *TypedScope) GetGlobalScope() *TypedScope {
	for s.parent != nil {
		s = s.parent
	}
	return s
}

// Java: AbstractScope.hasSlot(String)
func (s *TypedScope) HasSlot(name string) bool {
	h := scopeNameHash(name)
	bloom := scopeBloomFromHash(h)
	imp, isImplicit := AbstractScope_ImplicitVar_Lookup(name)
	for cur := s; cur != nil; cur = cur.lookupParent {
		if cur.varBloom&bloom == bloom && cur.findVar(name, h) >= 0 {
			return true
		}
		if isImplicit && cur.hasImplicitVars && cur.hasOwnImplicitVar(imp) {
			return true
		}
	}
	return false
}

func (s *TypedScope) getSlotV(name string) *TypedVar { return s.GetVar(name) }

func (s *TypedScope) GetSlot(name string) rhino.StaticSlot {
	if v := s.GetVar(name); v != nil {
		return v
	}
	return nil
}

// Java: AbstractScope.getArgumentsVar()
func (s *TypedScope) GetArgumentsVar() *TypedVar {
	for cur := s; cur != nil; cur = cur.lookupParent {
		if cur.hasImplicitVars {
			if v := cur.getOwnImplicitVar(AbstractScope_ImplicitVar_ARGUMENTS); v != nil {
				return v
			}
		}
	}
	return nil
}

// Java: AbstractScope.canDeclare(String)
func (s *TypedScope) CanDeclare(name string) bool {
	return !s.HasOwnSlot(name) &&
		(!s.IsFunctionBlockScope() ||
			!s.parent.HasOwnSlot(name) ||
			s.IsBleedingFunctionName(name))
}

// Java: AbstractScope.isBleedingFunctionName(String)
func (s *TypedScope) IsBleedingFunctionName(name string) bool {
	var_ := s.GetVar(name)
	return var_ != nil &&
		!var_.nameNode.IsNil() &&
		var_.nameNode.GetParent().IsFunction()
}

// Java: AbstractScope.getAllAccessibleVariables()
func (s *TypedScope) GetAllAccessibleVariables() iter.Seq[*TypedVar] {
	return func(yield func(*TypedVar) bool) {
		var seen map[string]struct{}
		for cur := s; cur != nil; cur = cur.parent {
			for _, v := range cur.vars {
				if _, dup := seen[v.name]; dup {
					continue
				}
				if cur.parent != nil || seen != nil {
					if seen == nil {
						seen = make(map[string]struct{}, len(cur.vars))
					}
					seen[v.name] = struct{}{}
				}
				if !yield(v) {
					return
				}
			}
		}
	}
}

// Java: AbstractScope.isHoistScope()
func (s *TypedScope) IsHoistScope() bool {
	return s.parent == nil ||
		s.rootNode.IsFunction() ||
		NodeUtil_IsFunctionBlock(s.rootNode) ||
		s.rootNode.IsModuleBody() ||
		NodeUtil_IsClassStaticBlock(s.rootNode)
}

// Java: AbstractScope.getClosestHoistScope()
func (s *TypedScope) GetClosestHoistScope() *TypedScope {
	for cur := s; cur != nil; cur = cur.parent {
		if cur.IsHoistScope() {
			return cur
		}
	}
	return nil
}

// Java: AbstractScope.getClosestCfgRootScope()
func (s *TypedScope) GetClosestCfgRootScope() *TypedScope {
	cur := s
	for !cur.IsCfgRootScope() {
		cur = cur.parent
	}
	return cur
}

// Java: AbstractScope.getClosestContainerScope()
func (s *TypedScope) GetClosestContainerScope() *TypedScope {
	scope := s.GetClosestHoistScope()
	if NodeUtil_IsFunctionBlock(scope.rootNode) {
		scope = scope.parent
		if scope.IsBlockScope() {
			panic(jutil.NewIllegalStateExceptionMsg(jutil.LenientFormat("%s", scope)))
		}
	}
	return scope
}

// Java: AbstractScope.getCommonParent(S)
func (s *TypedScope) GetCommonParent(other *TypedScope) *TypedScope {
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
func (s *TypedScope) HasSameContainerScope(other *TypedScope) bool {
	return s == other || s.GetClosestContainerScope() == other.GetClosestContainerScope()
}

// Java: AbstractScope.getScopeOfThis()
func (s *TypedScope) GetScopeOfThis() *TypedScope {
	scope := s.GetClosestContainerScope()
	for scope.parent != nil && !NodeUtil_IsNonArrowFunction(scope.rootNode) {
		scope = scope.parent.GetClosestContainerScope()
	}
	return scope
}
