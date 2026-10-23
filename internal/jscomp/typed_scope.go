// Port of com/google/javascript/jscomp/TypedScope.java @ closure-compiler v20261015

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

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

type typedScopeArena struct {
	scopes   []TypedScope
	vars     []TypedVar
	varSlots []jutil.LhmSlot[string, *TypedVar]
	names    []string
}

func (a *typedScopeArena) allocScope() *TypedScope {
	if len(a.scopes) == cap(a.scopes) {
		newCap := cap(a.scopes) * 2
		if newCap < 64 {
			newCap = 64
		} else if newCap > 1024 {
			newCap = 1024
		}
		a.scopes = make([]TypedScope, 0, newCap)
	}
	idx := len(a.scopes)
	a.scopes = a.scopes[:idx+1]
	return &a.scopes[idx]
}

func (a *typedScopeArena) allocVar(
	inferred bool,
	name string,
	nameNode *rhino.Node,
	type_ rhino.JSType,
	scope *TypedScope,
	index int,
	input *CompilerInput) *TypedVar {
	if len(a.vars) == cap(a.vars) {
		newCap := cap(a.vars) * 2
		if newCap < 128 {
			newCap = 128
		} else if newCap > 2048 {
			newCap = 2048
		}
		a.vars = make([]TypedVar, 0, newCap)
	}
	idx := len(a.vars)
	a.vars = a.vars[:idx+1]
	v := &a.vars[idx]
	v.InitAbstractVarBase(name, nameNode, scope, index, input, nil /* implicitGoogNamespaceDefinition */)
	if nameNode != nil {
		jutil.CheckArgumentMsg1(
			typedVar_NAME_NODE_TYPES.Contains(nameNode.GetToken()),
			"Invalid name node token %s",
			nameNode.GetToken())
	}
	v.type_ = type_
	if inferred {
		v.indexAndFlags |= varFlagTypeInferred
	}
	return v
}

func (a *typedScopeArena) allocVarSlots(n int) []jutil.LhmSlot[string, *TypedVar] {
	if n > 256 {
		return make([]jutil.LhmSlot[string, *TypedVar], n)
	}
	if len(a.varSlots)+n > cap(a.varSlots) {
		newCap := cap(a.varSlots) * 2
		if newCap < 256 {
			newCap = 256
		} else if newCap > 4096 {
			newCap = 4096
		}
		a.varSlots = make([]jutil.LhmSlot[string, *TypedVar], 0, newCap)
	}
	start := len(a.varSlots)
	a.varSlots = a.varSlots[:start+n]
	return a.varSlots[start : start+n : start+n]
}

func (a *typedScopeArena) allocNames(n int) []string {
	if n > 256 {
		return make([]string, 0, n)
	}
	if len(a.names)+n > cap(a.names) {
		newCap := cap(a.names) * 2
		if newCap < 256 {
			newCap = 256
		} else if newCap > 4096 {
			newCap = 4096
		}
		a.names = make([]string, 0, newCap)
	}
	start := len(a.names)
	a.names = a.names[:start+n]
	return a.names[start : start : start+n]
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
// See NodeTraversal and DataFlowAnalysis. Several methods in this class,
// such as isBlockScope throw an exception when called. The reason for this
// is that we want to shadow methods from the parent class, to avoid calling
// them accidentally.
//
// It is a rhino.StaticTypedScope: GetSlot, GetOwnSlot and GetParentScope
// keep rhino.StaticScope's signatures (assert `.(*TypedVar)`); see the note
// at the top of abstract_scope.go.
type TypedScope struct {
	AbstractScopeBase[*TypedScope, *TypedVar]

	parent *TypedScope
	module *Module
	arena  *typedScopeArena

	// Will shrink over time. Set to "ImmutableSet.of()" once empty (so don't try removing
	// elements from an empty set)
	//
	// Go-only: language-adapter — nil/empty slice is the ImmutableSet.of() sentinel.
	reservedNames []string

	// Go-only: performance — held by value so declaring the first variable
	// in a scope does not allocate a separate LinkedHashMap header; placed
	// before varBloom so trailing scalar fields lie past GC ptrdata.
	vars jutil.LinkedHashMap[string, *TypedVar]

	varBloom uint64

	// Scope.java contains an arguments field.
	// We haven't added it here because it's unused by the passes that need typed scopes.
}

// Equals implements jutil.Equaler using pointer identity (matching Java Object.equals).
func (s *TypedScope) Equals(other any) bool {
	that, ok := other.(*TypedScope)
	return ok && s == that
}

// HashCode implements jutil.Hasher using pointer identity (matching Java Object.hashCode).
func (s *TypedScope) HashCode() int32 {
	if s == nil {
		return 0
	}
	return jutil.IdentityHashCodeOfPointer(s)
}

// Java: TypedScope(TypedScope, Node)
func NewTypedScope(parent *TypedScope, rootNode *rhino.Node) *TypedScope {
	return newTypedScopeOwnedReservedNames(parent, rootNode, nil, nil, nil)
}

// NewTypedScopeWithReservedNames creates a new global or local TypedScope.
//
// parent is the containing scope, unless this is the global scope. rootNode
// can be any node that creates a scope. reservedNames is the set of simple
// names that will eventually be declared in this scope.
//
// Java: TypedScope(TypedScope, Node, Set<String>, Module)
func NewTypedScopeWithReservedNames(
	parent *TypedScope, rootNode *rhino.Node, reservedNames *jutil.LinkedHashSet[string], module *Module) *TypedScope {
	var copied []string
	if !jutil.CheckNotNull(reservedNames).IsEmpty() { // Java: NPE on a null set
		copied = reservedNames.ToSlice()
	}
	return newTypedScopeOwnedReservedNames(parent, rootNode, copied, module, nil)
}

func newTypedScopeOwnedReservedNames(
	parent *TypedScope, rootNode *rhino.Node, reservedNames []string, module *Module, arena *typedScopeArena) *TypedScope {
	var s *TypedScope
	if arena != nil {
		s = arena.allocScope()
	} else {
		s = &TypedScope{}
	}
	s.InitAbstractScopeBase(s, rootNode)
	if arena != nil {
		slotCap := len(reservedNames)
		if module != nil {
			slotCap += 4
		}
		if slotCap > 0 {
			s.vars.InitWithBuffer(arena.allocVarSlots(slotCap))
		} else {
			s.vars.InitWithBuffer(nil)
		}
	} else {
		s.vars.InitWithBuffer(nil)
	}
	s.CheckChildScope(parent)
	s.parent = parent
	s.depth = parent.depth + 1
	s.isBottom = false
	if len(reservedNames) > 0 {
		s.reservedNames = reservedNames
	}
	s.module = module
	s.arena = arena
	return s
}

// newTypedScopeWithIsBottom creates a empty Scope (bottom of the lattice).
//
// rootNode is typically a FUNCTION node or the global ROOT. isBottom is
// whether this is the bottom of a lattice. Otherwise, it must be a global
// scope.
//
// Java: TypedScope(Node, boolean)
func newTypedScopeWithIsBottom(rootNode *rhino.Node, isBottom bool) *TypedScope {
	s := &TypedScope{}
	s.InitAbstractScopeBase(s, rootNode)
	s.vars.InitWithBuffer(nil)
	s.CheckRootScope()
	s.parent = nil
	s.depth = 0
	s.isBottom = isBottom
	s.reservedNames = nil
	s.module = nil
	return s
}

func newGlobalTypedScopeWithCapacity(rootNode *rhino.Node, varCap int) *TypedScope {
	s := &TypedScope{}
	s.InitAbstractScopeBase(s, rootNode)
	if varCap > 0 {
		s.vars = *jutil.NewLinkedHashMapWithCapacity[string, *TypedVar](varCap)
	} else {
		s.vars.InitWithBuffer(nil)
	}
	s.CheckRootScope()
	s.parent = nil
	s.depth = 0
	s.isBottom = false
	s.reservedNames = nil
	s.module = nil
	return s
}

// Java: TypedScope.createGlobalScope(Node)
func TypedScope_CreateGlobalScope(rootNode *rhino.Node) *TypedScope {
	return newTypedScopeWithIsBottom(rootNode, false)
}

// Java: TypedScope.createLatticeBottom(Node)
func TypedScope_CreateLatticeBottom(rootNode *rhino.Node) *TypedScope {
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
	if len(s.reservedNames) > 0 {
		jutil.CheckStateMsg2(
			false,
			"Expected %s to have no reserved names, found: %s. This probably indicates a bug in"+
				" TypedScopeCreator where it is failing to declare a variable.",
			s,
			"["+strings.Join(s.reservedNames, ", ")+"]")
	}
	// let a (16-bit) VM garbage collect 64 bytes per TypedScope. (ImmutableSet.of() returns a
	// singleton)
	s.reservedNames = nil
	s.arena = nil
}

// IsBottom reports whether this is the bottom of the lattice.
//
// Java: TypedScope.isBottom()
func (s *TypedScope) IsBottom() bool {
	return s.isBottom
}

// Java: TypedScope.getModule()
func (s *TypedScope) GetModule() *Module {
	return s.module
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
	root := s.GetRootNode()
	if s.IsGlobal() {
		if t := rhino.ObjectType_Cast(root.GetJSType()); t != nil {
			return t
		}
		return nil
	} else if NodeUtil_IsNonArrowFunction(root) {
		nodeType := root.GetJSType()
		if nodeType != nil && nodeType.IsFunctionType() {
			return nodeType.ToMaybeFunctionType().GetTypeOfThis()
		} else {
			// Executed when the current scope has not been typechecked.
			return nil
		}
	} else if s.IsStaticBlockScope() {
		return s.GetParent().GetRootNode().GetJSType()
	} else if s.IsMemberFieldDefScope() || s.IsComputedFieldDefRhsScope() {
		classType := s.GetParent().GetRootNode().GetJSType()
		if root.IsStaticMember() {
			return classType
		} else {
			if t := classType.AssertFunctionType().GetInstanceType(); t != nil {
				return t
			}
			return nil
		}
	} else {
		return s.GetParent().GetTypeOfThis()
	}
}

// Java: TypedScope.declare(String, Node, JSType, CompilerInput, boolean)
func (s *TypedScope) Declare(name string, nameNode *rhino.Node,
	type_ rhino.JSType, input *CompilerInput, inferred bool) *TypedVar {
	var_, _ := s.declareWithReserved(name, nameNode, type_, input, inferred)
	return var_
}

func (s *TypedScope) declareWithReserved(name string, nameNode *rhino.Node,
	type_ rhino.JSType, input *CompilerInput, inferred bool) (*TypedVar, bool) {
	jutil.CheckState(name != "") // Java: name != null && !name.isEmpty()
	wasReserved := false
	if len(s.reservedNames) > 0 {
		// (reservedNames only contains simple, non-qualified names; so it's completely normal to
		// declare qualified names that were not 'reserved')
		for i, rn := range s.reservedNames {
			if rn == name {
				wasReserved = true
				copy(s.reservedNames[i:], s.reservedNames[i+1:])
				s.reservedNames[len(s.reservedNames)-1] = ""
				s.reservedNames = s.reservedNames[:len(s.reservedNames)-1]
				break
			}
		}
	}
	var var_ *TypedVar
	if s.arena != nil {
		var_ = s.arena.allocVar(inferred, name, nameNode, type_, s, s.GetVarCount(), input)
	} else {
		var_ = NewTypedVar(inferred, name, nameNode, type_, s, s.GetVarCount(), input)
	}
	s.DeclareInternal(name, var_)
	return var_, wasReserved
}

// Java: AbstractScope.undeclare(V)
func (s *TypedScope) Undeclare(var_ *TypedVar) {
	jutil.CheckState(var_.getScopeS() == s)
	current := s.vars.Get(var_.GetName())
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
	s.vars.Remove(var_.GetName())
}

// Java: AbstractScope.declareInternal(String, V)
func (s *TypedScope) DeclareInternal(name string, var_ *TypedVar) {
	if s.IsFunctionBlockScope() && s.parent.HasOwnSlot(name) && !s.HasOwnSlot(name) && !s.IsBleedingFunctionName(name) {
		panic(jutil.NewIllegalStateExceptionMsg(jutil.LenientFormat("Illegal shadow: %s", var_.GetNode())))
	}
	s.varBloom |= scopeBloomMask(name)
	s.varsInitialized = true
	s.vars.Put(name, var_)
}

// Java: AbstractScope.clearVarsInternal()
func (s *TypedScope) ClearVarsInternal() {
	s.varBloom = 0
	if !s.vars.IsEmpty() {
		s.vars.Clear()
	}
}

// Java: AbstractScope.getVarIterable()
func (s *TypedScope) GetVarIterable() iter.Seq[*TypedVar] {
	return s.vars.Values()
}

// Java: AbstractScope.getAllSymbols()
func (s *TypedScope) GetAllSymbols() iter.Seq[*TypedVar] {
	return s.vars.Values()
}

// Java: AbstractScope.getVarCount()
func (s *TypedScope) GetVarCount() int {
	return s.vars.Size()
}

func (s *TypedScope) getOwnTypedSlot(name string) *TypedVar {
	bloom := scopeBloomMask(name)
	if (s.varBloom&bloom) != bloom || s.vars.IsEmpty() {
		return nil
	}
	return s.vars.Get(name)
}

func (s *TypedScope) getOwnSlotV(name string) *TypedVar {
	if v := s.getOwnTypedSlot(name); v != nil {
		return v
	}
	if imp := AbstractScope_ImplicitVar_Of(name); imp != nil {
		return s.GetOwnImplicitSlot(imp)
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
	if imp := AbstractScope_ImplicitVar_Of(name); imp != nil {
		return s.HasOwnImplicitSlot(imp)
	}
	return false
}

// Java: TypedScope.makeImplicitVar(ImplicitVar)
func (s *TypedScope) MakeImplicitVar(var_ AbstractScope_ImplicitVar) *TypedVar {
	if s.IsGlobal() {
		// TODO(sdh): This is incorrect for 'global this', but since that's currently not handled
		// by this code, it's okay to bail out now until we find the root cause.  See b/74980936.
		return nil
	} else if var_ == AbstractScope_ImplicitVar_EXPORTS {
		// Instead of using the implicit 'exports' var, we want to pretend that the var is actually
		// declared.
		return nil
	}
	if s.arena != nil {
		return s.arena.allocVar(false, var_.name(), nil, s.GetImplicitVarType(var_), s, -1, nil)
	}
	return NewTypedVar(false, var_.name(), nil, s.GetImplicitVarType(var_), s, -1, nil)
}

// Java: TypedScope.hasOwnImplicitSlot(ImplicitVar)
func (s *TypedScope) HasOwnImplicitSlot(name *AbstractScope_ImplicitVar) bool {
	return name != nil && *name != AbstractScope_ImplicitVar_EXPORTS && name.IsMadeByScope(s)
}

// Java: TypedScope.getImplicitVarType(ImplicitVar)
func (s *TypedScope) GetImplicitVarType(var_ AbstractScope_ImplicitVar) rhino.JSType {
	switch var_ {
	case AbstractScope_ImplicitVar_ARGUMENTS:
		// Look for an extern named "arguments" and use its type if available.
		// TODO(sdh): consider looking for "Arguments" ctor rather than "arguments" var: this could
		// allow deleting the variable, which doesn't really belong in externs in the first place.
		globalArgs := s.GetGlobalScope().GetVar(Var_ARGUMENTS)
		if globalArgs != nil && globalArgs.IsExtern() {
			return globalArgs.GetType()
		}
		return nil
	case AbstractScope_ImplicitVar_THIS:
		return s.GetTypeOfThis()
	case AbstractScope_ImplicitVar_SUPER:
		// Inside a constructor, `super` may have two different types. Calls to `super()` use the
		// super-ctor type, while property accesses use the super-instance type. This logic always
		// returns the latter case.
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
// These variables can safely be set to undefined (rather than unknown) at
// the start of type inference, and will be reset to the correct type when
// analyzing the first assignment to them. Parameters and externs are
// excluded because they are not initialized in the function body, and
// lexically-bound variables (let and const) are excluded because they are
// initialized when inferring the LET/CONST node, which is guaranteed to
// occur before any use, since they are not hoisted.
//
// Java: TypedScope.getDeclarativelyUnboundVarsWithoutTypes()
func (s *TypedScope) GetDeclarativelyUnboundVarsWithoutTypes() iter.Seq[*TypedVar] {
	return func(yield func(*TypedVar) bool) {
		for v := range s.vars.Values() {
			if s.IsDeclarativelyUnboundVarWithoutType(v) && !yield(v) {
				return
			}
		}
	}
}

// Java: TypedScope.isDeclarativelyUnboundVarWithoutType(TypedVar)
func (s *TypedScope) IsDeclarativelyUnboundVarWithoutType(var_ *TypedVar) bool {
	return var_.GetParentNode() != nil &&
		var_.GetType() == nil &&
		// TODO(bradfordcsmith): update this for destructuring
		var_.GetParentNode().IsVar() && // NOTE: explicitly excludes let/const
		!var_.IsExtern()
}

// GetVar returns the slot for name, considering shadowing of qualified
// names.
//
// The superclass method does not handle shadowing.
//
// Lookup of qualified names (i.e. names with dots) executes against scopes
// in the following precedence:
//
//  1. This Scope.
//  2. The first ancestor Scope, if any, that declares the root of the
//     qualified name.
//  3. The global Scope.
//
// An example of where this is necessary: say the global scope contains "a"
// and "a.b" and a function scope contains "a". When looking up "a.b" in the
// function scope, AbstractScope::getVar returns "a.b". This method returns
// null because the global "a" is shadowed.
//
func (s *TypedScope) lookupRootVar(name string) *TypedVar {
	bloom := scopeBloomMask(name)
	implicit := AbstractScope_ImplicitVar_Of(name)
	if implicit == nil {
		for cur := s; cur != nil; cur = cur.parent {
			if (cur.varBloom&bloom) == bloom && !cur.vars.IsEmpty() {
				if v := cur.vars.Get(name); v != nil {
					return v
				}
			}
		}
		return nil
	}
	for cur := s; cur != nil; cur = cur.parent {
		if (cur.varBloom&bloom) == bloom && !cur.vars.IsEmpty() {
			if v := cur.vars.Get(name); v != nil {
				return v
			}
		}
		if v := cur.GetOwnImplicitSlot(implicit); v != nil {
			return v
		}
	}
	return nil
}

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
	return rootVar.getScopeS().getOwnSlotV(name)
}

// Java: TypedScope.getTypeThroughNamespace(String)
func (s *TypedScope) GetTypeThroughNamespace(moduleId string) rhino.JSType {
	// Go-only: language-adapter — '.' is ASCII, so byte indices split the
	// string exactly where Java's UTF-16 indices do.
	split := strings.LastIndexByte(moduleId, '.')
	if split >= 0 {
		parentName := moduleId[:split]
		prop := moduleId[split+1:]
		parentType := s.GetTypeThroughNamespace(parentName)
		if parentType == nil || parentType.ToMaybeObjectType() == nil {
			return nil
		}
		return parentType.AssertObjectType().GetPropertyType(prop)
	} else {
		var_ := s.getSlotV(moduleId)
		if var_ != nil {
			return var_.GetType()
		}
		return nil
	}
}

// Java: TypedScope.getTopmostScopeOfEventualDeclaration(String)
func (s *TypedScope) GetTopmostScopeOfEventualDeclaration(name string) rhino.StaticScope {
	if s.getOwnSlotV(name) != nil || slices.Contains(s.reservedNames, name) {
		return s
	} else if s.GetParent() == nil {
		return nil
	} else {
		// Recurse on the parent because it, too, may be incomplete.
		return s.GetParent().GetTopmostScopeOfEventualDeclaration(name)
	}
}

// LookupQualifiedName is the rhino.StaticTypedScope default method, which
// TypedScope inherits.
//
// Java: StaticTypedScope.lookupQualifiedName(QualifiedName) (default method)
func (s *TypedScope) LookupQualifiedName(qname rhino.QualifiedName) rhino.JSType {
	return rhino.StaticTypedScope_LookupQualifiedName(s, qname)
}
