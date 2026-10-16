// Port of com/google/javascript/jscomp/TypedScope.java @ closure-compiler v20260415

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
	"strings"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

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
	depth  int
	module *Module

	// Whether this is a bottom scope for the purposes of type inference.
	isBottom bool

	// Will shrink over time. Set to "ImmutableSet.of()" once empty (so don't try removing
	// elements from an empty set)
	//
	// Go-only: language-adapter — nil is the ImmutableSet.of() sentinel.
	reservedNames *jutil.LinkedHashSet[string]

	// Scope.java contains an arguments field.
	// We haven't added it here because it's unused by the passes that need typed scopes.
}

// Java: TypedScope(TypedScope, Node)
func NewTypedScope(parent *TypedScope, rootNode *rhino.Node) *TypedScope {
	return NewTypedScopeWithReservedNames(parent, rootNode, jutil.NewLinkedHashSet[string](), nil)
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
	s := &TypedScope{}
	s.InitAbstractScopeBase(s, rootNode)
	s.CheckChildScope(parent)
	s.parent = parent
	s.depth = parent.depth + 1
	s.isBottom = false
	if jutil.CheckNotNull(reservedNames).IsEmpty() { // Java: NPE on a null set
		s.reservedNames = nil
	} else {
		s.reservedNames = jutil.NewLinkedHashSet[string](reservedNames.ToSlice()...)
	}
	s.module = module
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
	s.CheckRootScope()
	s.parent = nil
	s.depth = 0
	s.isBottom = isBottom
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
	jutil.CheckStateMsg2(
		s.reservedNames.IsEmpty(),
		"Expected %s to have no reserved names, found: %s. This probably indicates a bug in"+
			" TypedScopeCreator where it is failing to declare a variable.",
		s,
		s.reservedNames) // non-empty, hence non-nil, when printed
	// let a (16-bit) VM garbage collect 64 bytes per TypedScope. (ImmutableSet.of() returns a
	// singleton)
	s.reservedNames = nil
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

// Java: TypedScope.getDepth()
func (s *TypedScope) GetDepth() int {
	return s.depth
}

// Java: TypedScope.getParent()
func (s *TypedScope) GetParent() *TypedScope {
	return s.parent
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
	jutil.CheckState(name != "") // Java: name != null && !name.isEmpty()
	if !s.reservedNames.IsEmpty() {
		// (reservedNames only contains simple, non-qualified names; so it's completely normal to
		// declare qualified names that were not 'reserved')
		s.reservedNames.Remove(name)
	}
	var_ := NewTypedVar(inferred, name, nameNode, type_, s, s.GetVarCount(), input)
	s.DeclareInternal(name, var_)
	return var_
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
	vars := s.GetVarIterable()
	return func(yield func(*TypedVar) bool) {
		for v := range vars {
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
// Java: TypedScope.getVar(String)
func (s *TypedScope) GetVar(name string) *TypedVar {
	ownSlot := s.getOwnSlotV(name)
	if ownSlot != nil {
		// Micro-optimization: variables declared directly in this scope cannot have been shadowed.
		return ownSlot
	} else if s.GetParent() == nil {
		return nil
	}
	// Find the root name and its slot.
	// Go-only: language-adapter — '.' is ASCII, so its byte index splits the
	// string exactly where Java's UTF-16 index does.
	dot := strings.IndexByte(name, '.')
	rootName := name
	if dot >= 0 {
		rootName = name[:dot]
	}
	rootVar := s.AbstractScopeBase.GetVar(rootName) // Use the superclass method to skip string checks.

	if dot < 0 {
		return rootVar
	} else if rootVar == nil {
		// Default to the global scope because externs may have qualified names with undeclared roots.
		return s.GetParent().GetGlobalScope().getOwnSlotV(name)
	} else {
		// Qualified names 'a.b.c' are declared in the same scope as 'a', never a child scope, which
		// is why calling `getOwnSlot` is sufficient.
		return rootVar.getScopeS().getOwnSlotV(name)
	}
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
	if s.getOwnSlotV(name) != nil || s.reservedNames.Contains(name) {
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
