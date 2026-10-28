// Port of com/google/javascript/jscomp/AbstractScope.java @ closure-compiler v20261021

/*
 * Copyright 2018 The Closure Compiler Authors.
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
)

// How the generic scope hierarchy is mapped (read this before porting code
// that uses scopes or vars).
//
// Java has AbstractScope<S extends AbstractScope<S, V>, V extends
// AbstractVar<S, V>> with two subclasses, Scope (S=Scope, V=Var) and
// TypedScope (S=TypedScope, V=TypedVar), and AbstractVar<S, V> with Var and
// TypedVar. Go-only: language-adapter — Go has neither wildcards nor
// covariant overrides, so:
//
//   - AbstractScope (this interface) is the wildcard type AbstractScope<?, ?>
//     (what NodeTraversal, ScopeCreator and MemoizedScopeCreator handle). It
//     declares every AbstractScope method whose signature mentions neither S
//     nor V.
//   - AbstractScopeBase[S, V] is the generic base struct with the fields and
//     every implementation. The methods whose signature mentions S or V
//     (getParent, getVar, getGlobalScope, getClosestHoistScope,
//     getVarIterable, ...) exist only there and on the concrete scopes, where
//     they are typed: (*Scope).GetVar returns *Var, (*TypedScope).GetParent
//     returns *TypedScope, (*Scope).GetVarIterable is an iter.Seq[*Var].
//   - Java narrows three rhino.StaticScope methods to S/V: getParentScope(),
//     getSlot(String) and getOwnSlot(String). They keep rhino.StaticScope's
//     Go signatures (as rhino.StaticTypedScope does), so that *Scope is a
//     rhino.StaticScope and *TypedScope a rhino.StaticTypedScope; they never
//     return a typed nil. Callers that need the Java type assert, with the
//     comma-ok form so that Java null stays nil:
//       v, _ := scope.GetOwnSlot(name).(*Var)
//       tv, _ := typedScope.GetSlot(name).(*TypedVar)
//     The same holds for AbstractVar's getScope() (rhino.StaticSlot),
//     getDeclaration() and getSymbol(): `v.GetScope().(*Scope)`.
//   - abstractScopeBound / abstractVarBound are the Go form of the bounds
//     `S extends AbstractScope<S, V>` / `V extends AbstractVar<S, V>`.
//
// Virtual calls from the base go through self (an S): getDepth, getParent,
// getVar (TypedScope overrides it), hasOwnImplicitSlot (TypedScope overrides
// it), makeImplicitVar, typed()/untyped().

// AbstractScope contains information about a variable scope in JavaScript.
// Scopes can be nested, a scope points back to its parent scope. A Scope
// contains information about variables defined in that scope.
//
// ES 2015 introduces new scoping rules, which adds some complexity to this
// class. In particular, scopes fall into two mutually exclusive categories:
// block and container. Block scopes are all scopes whose roots are blocks,
// as well as any control structures whose optional blocks are omitted. These
// scopes did not exist at all prior to ES 2015. Container scopes comprise
// function scopes, global scopes, and module scopes, and (aside from modules,
// which didn't exist in ES5) corresponds to the ES5 scope rules. This
// corresponds roughly to one container scope per CFG root (but not exactly,
// due to SCRIPT-level CFGs).
//
// All container scopes, plus the outermost block scope within a function
// (i.e. the function block scope) are considered hoist scopes. All functions
// thus have two hoist scopes: the function scope and the function block
// scope. Hoist scopes are relevant because "var" declarations are hoisted to
// the closest hoist scope, as opposed to "let" and "const" which always apply
// to the specific scope in which they occur.
//
// Note that every function actually has two distinct hoist scopes: a
// container scope on the FUNCTION node, and a block-scope on the top-level
// BLOCK in the function (the "function block"). Local variables are declared
// on the function block, while parameters and optionally the function name
// (if it bleeds, i.e. from a named function expression) are declared on the
// container scope. This is required so that default parameter initializers
// can refer to names from outside the function that could possibly be
// shadowed in the function block. But these scopes are not fully independent
// of one another, since the language does not allow a top-level local
// variable to shadow a parameter name - so in some situations these scopes
// must be treated as a single scope.
//
// See NodeTraversal.
type AbstractScope interface {
	rhino.StaticScope

	// GetDepth is the depth of the scope. The global scope has depth 0.
	//
	// Java: AbstractScope.getDepth()
	GetDepth() int

	// Java: AbstractScope.toString()
	String() string

	// Java: AbstractScope.untyped()
	Untyped() *Scope

	// Java: AbstractScope.typed()
	Typed() *TypedScope

	// Java: AbstractScope.clearVarsInternal()
	ClearVarsInternal()

	// HasOwnImplicitSlot returns true iff this scope implies a slot with the
	// given name.
	//
	// Java: AbstractScope.hasOwnImplicitSlot(ImplicitVar)
	HasOwnImplicitSlot(name *AbstractScope_ImplicitVar) bool

	// Java: AbstractScope.hasOwnSlot(String)
	HasOwnSlot(name string) bool

	// Java: AbstractScope.hasSlot(String)
	HasSlot(name string) bool

	// Java: AbstractScope.canDeclare(String)
	CanDeclare(name string) bool

	// Java: AbstractScope.getVarCount()
	GetVarCount() int

	// Java: AbstractScope.isGlobal()
	IsGlobal() bool

	// Java: AbstractScope.isLocal()
	IsLocal() bool

	// Java: AbstractScope.isBlockScope()
	IsBlockScope() bool

	// Java: AbstractScope.isStaticBlockScope()
	IsStaticBlockScope() bool

	// Java: AbstractScope.isFunctionBlockScope()
	IsFunctionBlockScope() bool

	// Java: AbstractScope.isFunctionScope()
	IsFunctionScope() bool

	// Java: AbstractScope.isModuleScope()
	IsModuleScope() bool

	// Java: AbstractScope.isMemberFieldDefScope()
	IsMemberFieldDefScope() bool

	// Java: AbstractScope.isComputedFieldDefRhsScope()
	IsComputedFieldDefRhsScope() bool

	// Java: AbstractScope.isCatchScope()
	IsCatchScope() bool

	// Java: AbstractScope.isCfgRootScope()
	IsCfgRootScope() bool

	// Java: AbstractScope.isHoistScope()
	IsHoistScope() bool

	// Java: AbstractScope.checkRootScope()
	CheckRootScope()
}

// abstractScopeBound is the bound `S extends AbstractScope<S, V>`: the
// methods AbstractScopeBase and AbstractVarBase call on an S other than
// their own self. All but GetParent, GetDepth and MakeImplicitVar are
// promoted from AbstractScopeBase.
//
// Go-only: language-adapter — the bound of a Java type parameter; see the
// note at the top of this file. It may not mention Scope, TypedScope, their
// vars or AbstractScope (whose Untyped() returns *Scope): Scope embeds
// AbstractScopeBase[*Scope, *Var], and a constraint that refers back to its
// own instantiations is an invalid recursive type in Go. Where Java passes an
// S as an AbstractScope<?, ?>, the Go code converts it with
// abstractScope_of.
type abstractScopeBound[S any, V any] interface {
	comparable
	rhino.StaticScope

	// Java: AbstractScope.getDepth()
	GetDepth() int

	// Java: AbstractScope.hasOwnImplicitSlot(ImplicitVar)
	HasOwnImplicitSlot(name *AbstractScope_ImplicitVar) bool

	// Java: AbstractScope.hasOwnSlot(String)
	HasOwnSlot(name string) bool

	// Java: AbstractScope.canDeclare(String)
	CanDeclare(name string) bool

	// Java: AbstractScope.isGlobal()
	IsGlobal() bool

	// Java: AbstractScope.isLocal()
	IsLocal() bool

	// Java: AbstractScope.isBlockScope()
	IsBlockScope() bool

	// Java: AbstractScope.isStaticBlockScope()
	IsStaticBlockScope() bool

	// Java: AbstractScope.isFunctionBlockScope()
	IsFunctionBlockScope() bool

	// Java: AbstractScope.isFunctionScope()
	IsFunctionScope() bool

	// Java: AbstractScope.isModuleScope()
	IsModuleScope() bool

	// Java: AbstractScope.isCfgRootScope()
	IsCfgRootScope() bool

	// Java: AbstractScope.isHoistScope()
	IsHoistScope() bool

	// Java: AbstractScope.getParent()
	GetParent() S

	// Java: AbstractScope.getVar(String)
	GetVar(name string) V

	// Java: AbstractScope.makeImplicitVar(ImplicitVar)
	MakeImplicitVar(type_ AbstractScope_ImplicitVar) V

	// Java: AbstractScope.getOwnImplicitSlot(ImplicitVar)
	GetOwnImplicitSlot(name *AbstractScope_ImplicitVar) V

	// Java: AbstractScope.getVarIterable()
	GetVarIterable() iter.Seq[V]

	// Java: AbstractScope.getClosestContainerScope()
	GetClosestContainerScope() S

	getOwnSlotV(name string) V
}

// AbstractScopeBase holds the fields and the implementations of
// AbstractScope; Scope and TypedScope embed it and call
// InitAbstractScopeBase with themselves.
type AbstractScopeBase[S abstractScopeBound[S, V], V abstractVarBound[S, V]] struct {
	// Go-only: language-adapter — the EnumMap<ImplicitVar, V>, indexed by
	// ordinal (a zero V is an absent key, as EnumMap.computeIfAbsent never
	// stores null); nil is the ImmutableMap.of() sentinel.
	implicitVars *[4]V

	rootNode rhino.Node

	depth           int32
	varsInitialized bool
	isBottom        bool
	hasImplicitVars bool
	scanDone        bool
	scopeID         uint32
}

// Java: AbstractScope.getDepth()
func (b *AbstractScopeBase[S, V]) GetDepth() int {
	return int(b.depth)
}

func abstractScope_of[S abstractScopeBound[S, V], V abstractVarBound[S, V]](s S) AbstractScope {
	var null S
	if s == null {
		return nil
	}
	return any(s).(AbstractScope)
}

// InitAbstractScopeBase is the AbstractScope(Node) constructor; self is the
// concrete scope.
//
// Java: AbstractScope(Node)
func (b *AbstractScopeBase[S, V]) InitAbstractScopeBase(_ S, rootNode rhino.Node) {
	b.varsInitialized = false
	b.isBottom = false
	b.hasImplicitVars = false
	b.scanDone = false
	b.depth = 0
	b.implicitVars = nil
	b.rootNode = rhino.CheckNotNull(rootNode)
}

// String is final in Java.
//
// Java: AbstractScope.toString()
func (b *AbstractScopeBase[S, V]) String() string {
	return "Scope@" + jutil.ToJavaString(b.rootNode)
}

// Java: AbstractScope.untyped()
func (b *AbstractScopeBase[S, V]) Untyped() *Scope {
	panic(jutil.NewIllegalStateExceptionMsg("untyped() called, but not an untyped scope."))
}

// Java: AbstractScope.typed()
func (b *AbstractScopeBase[S, V]) Typed() *TypedScope {
	panic(jutil.NewIllegalStateExceptionMsg("typed() called, but not a typed scope."))
}

// Java: AbstractScope.contains(S)
func (b *AbstractScopeBase[S, V]) contains(self S, other S) bool {
	var null S
	s := jutil.CheckNotNull(other)
	for s != null {
		if s == self {
			return true
		}
		s = s.GetParent()
	}
	return false
}

// GetRootNode gets the container node of the scope. This is typically the
// FUNCTION node or the global BLOCK/SCRIPT node.
//
// Java: AbstractScope.getRootNode()
func (b *AbstractScopeBase[S, V]) GetRootNode() rhino.Node {
	return b.rootNode
}

// GetGlobalScope walks up the tree to find the global scope.
//
// Java: AbstractScope.getGlobalScope()
func (b *AbstractScopeBase[S, V]) getGlobalScope(self S) S {
	var null S
	result := self
	for result.GetParent() != null {
		result = result.GetParent()
	}
	return result
}

// HasSlot returns true if a variable is declared in this or any parent
// scope.
//
// Java: AbstractScope.hasSlot(String)
func (b *AbstractScopeBase[S, V]) hasSlot(self S, name string) bool {
	var null S
	for scope := self; scope != null; scope = scope.GetParent() {
		if scope.HasOwnSlot(name) {
			return true
		}
	}
	return false
}

// Java: AbstractScope.getOwnImplicitSlot(ImplicitVar)
func (b *AbstractScopeBase[S, V]) getOwnImplicitSlot(self S, name *AbstractScope_ImplicitVar) V {
	if !self.HasOwnImplicitSlot(name) {
		var null V
		return null
	}

	// For memory savings, only initialize the map once it needs to add its first element
	if b.implicitVars == nil {
		b.implicitVars = new([4]V)
	}

	var null V
	if v := b.implicitVars[*name]; v != null {
		return v
	}
	v := self.MakeImplicitVar(*name)
	if v != null {
		b.implicitVars[*name] = v
	}
	return v
}

// Java: AbstractScope.getArgumentsVar()
func (b *AbstractScopeBase[S, V]) getArgumentsVar(self S) V {
	var nullS S
	var null V
	for scope := self; scope != nullS; scope = scope.GetParent() {
		if arguments := scope.GetOwnImplicitSlot(AbstractScope_ImplicitVar_ARGUMENTS.ptr()); arguments != null {
			return arguments
		}
	}
	return null
}

// Java: AbstractScope.canDeclare(String)
func (b *AbstractScopeBase[S, V]) canDeclare(self S, name string) bool {
	return !self.HasOwnSlot(name) &&
		(!self.IsFunctionBlockScope() ||
			!self.GetParent().HasOwnSlot(name) ||
			b.isBleedingFunctionName(self, name))
}

// Java: AbstractScope.isBleedingFunctionName(String)
func (b *AbstractScopeBase[S, V]) isBleedingFunctionName(self S, name string) bool {
	var null V
	var_ := self.GetVar(name)
	return var_ != null &&
		!var_.GetNode().IsNil() &&
		var_.GetNode().GetParent().IsFunction()
}

// Java: AbstractScope.getAllAccessibleVariables()
func (b *AbstractScopeBase[S, V]) getAllAccessibleVariables(self S) iter.Seq[V] {
	accessibleVars := jutil.NewLinkedHashMap[string, V]()
	var null S
	s := self

	for s != null {
		for v := range s.GetVarIterable() {
			accessibleVars.PutIfAbsent(v.GetName(), v)
		}
		s = s.GetParent()
	}
	return accessibleVars.Values()
}

// Java: AbstractScope.isBlockScope()
func (b *AbstractScopeBase[S, V]) IsBlockScope() bool {
	return NodeUtil_CreatesBlockScope(b.rootNode)
}

// Java: AbstractScope.isStaticBlockScope()
func (b *AbstractScopeBase[S, V]) IsStaticBlockScope() bool {
	return NodeUtil_IsClassStaticBlock(b.GetRootNode())
}

// Java: AbstractScope.isFunctionBlockScope()
func (b *AbstractScopeBase[S, V]) IsFunctionBlockScope() bool {
	return NodeUtil_IsFunctionBlock(b.GetRootNode())
}

// Java: AbstractScope.isFunctionScope()
func (b *AbstractScopeBase[S, V]) IsFunctionScope() bool {
	return b.GetRootNode().IsFunction()
}

// Java: AbstractScope.isModuleScope()
func (b *AbstractScopeBase[S, V]) IsModuleScope() bool {
	return b.GetRootNode().IsModuleBody()
}

// Java: AbstractScope.isMemberFieldDefScope()
func (b *AbstractScopeBase[S, V]) IsMemberFieldDefScope() bool {
	return b.GetRootNode().IsMemberFieldDef()
}

// Java: AbstractScope.isComputedFieldDefRhsScope()
func (b *AbstractScopeBase[S, V]) IsComputedFieldDefRhsScope() bool {
	return b.GetRootNode().IsComputedFieldDef()
}

// Java: AbstractScope.isCatchScope()
func (b *AbstractScopeBase[S, V]) IsCatchScope() bool {
	return b.GetRootNode().IsBlock() &&
		b.GetRootNode().HasOneChild() &&
		b.GetRootNode().GetFirstChild().IsCatch()
}

// Java: AbstractScope.isCfgRootScope()
func (b *AbstractScopeBase[S, V]) IsCfgRootScope() bool {
	return NodeUtil_IsValidCfgRoot(b.rootNode)
}

// Java: AbstractScope.isHoistScope()
func (b *AbstractScopeBase[S, V]) isHoistScope(self S) bool {
	return self.IsFunctionScope() ||
		self.IsFunctionBlockScope() ||
		self.IsGlobal() ||
		self.IsModuleScope() ||
		self.IsStaticBlockScope()
}

// Java: AbstractScope.getClosestHoistScope()
func (b *AbstractScopeBase[S, V]) getClosestHoistScope(self S) S {
	var null S
	current := self
	for current != null {
		if current.IsHoistScope() {
			return current
		}
		current = current.GetParent()
	}
	return null
}

// Java: AbstractScope.getClosestCfgRootScope()
func (b *AbstractScopeBase[S, V]) getClosestCfgRootScope(self S) S {
	current := self
	for !current.IsCfgRootScope() {
		current = current.GetParent()
	}
	return current
}

// Java: AbstractScope.getClosestContainerScope()
func (b *AbstractScopeBase[S, V]) getClosestContainerScope(self S) S {
	scope := b.getClosestHoistScope(self)
	if scope.IsFunctionBlockScope() {
		scope = scope.GetParent()
		if scope.IsBlockScope() {
			panic(jutil.NewIllegalStateExceptionMsg(jutil.LenientFormat("%s", scope)))
		}
	}
	return scope
}

// CheckChildScope performs simple validity checks on when constructing a
// child scope.
//
// Java: AbstractScope.checkChildScope(S)
func (b *AbstractScopeBase[S, V]) CheckChildScope(parent S) {
	var null S
	if parent == null {
		panic(jutil.NewNullPointerException())
	}
	if !NodeUtil_CreatesScope(b.rootNode) {
		panic(jutil.NewIllegalArgumentExceptionMsg(jutil.LenientFormat("%s", b.rootNode)))
	}
	if b.rootNode == parent.GetRootNode() {
		panic(jutil.NewIllegalArgumentExceptionMsg(
			jutil.LenientFormat("rootNode should not be the parent's root node: %s", b.rootNode)))
	}
}

// CheckRootScope performs simple validity checks on when constructing a
// root scope.
//
// Java: AbstractScope.checkRootScope()
func (b *AbstractScopeBase[S, V]) CheckRootScope() {
	// TODO(tbreisacher): Can we tighten this to just NodeUtil.createsScope?
	if !NodeUtil_CreatesScope(b.rootNode) && !b.rootNode.IsScript() && !b.rootNode.IsRoot() {
		panic(jutil.NewIllegalArgumentExceptionMsg(jutil.LenientFormat("%s", b.rootNode)))
	}
}

// Java: AbstractScope.getCommonParent(S)
func (b *AbstractScopeBase[S, V]) getCommonParent(self S, other S) S {
	var null S
	left := self
	right := other
	for left != null && right != null && left != right {
		leftDepth := left.GetDepth()
		rightDepth := right.GetDepth()
		if leftDepth >= rightDepth {
			left = left.GetParent()
		}
		if leftDepth <= rightDepth {
			right = right.GetParent()
		}
	}

	jutil.CheckState(left != null && left == right)
	return left
}

// Java: AbstractScope.hasSameContainerScope(S)
func (b *AbstractScopeBase[S, V]) hasSameContainerScope(self S, other S) bool {
	return self == other || self.GetClosestContainerScope() == other.GetClosestContainerScope()
}

// Java: AbstractScope.getScopeOfThis()
func (b *AbstractScopeBase[S, V]) getScopeOfThis(self S) S {
	scope := self.GetClosestContainerScope()
	for !scope.IsGlobal() && !NodeUtil_IsNonArrowFunction(scope.GetRootNode()) {
		scope = scope.GetParent().GetClosestContainerScope()
	}
	return scope
}

// AbstractScope_ImplicitVar is the three implicit var types, which are
// defined implicitly (at least) in every vanilla function scope without
// actually being declared.
//
// Java: AbstractScope.ImplicitVar
type AbstractScope_ImplicitVar int32

const (
	AbstractScope_ImplicitVar_ARGUMENTS AbstractScope_ImplicitVar = iota
	AbstractScope_ImplicitVar_EXPORTS
	AbstractScope_ImplicitVar_SUPER
	// TODO(sdh): Expand THIS.isMadeByScope to check super.isMadeByScope(scope) || scope.isGlobal()
	// Currently this causes a number of problems (see b/74980936), but could eventually lead to
	// better type information.  We might also want to restrict this so that module-root scopes
	// explicitly *don't* have access to the global this, though I think this is more than just
	// returning false in isMadeByScope - rather, getVar() needs to stop checking and immediately
	// return null.
	AbstractScope_ImplicitVar_THIS
)

var abstractScope_ImplicitVar_names = [...]string{"ARGUMENTS", "EXPORTS", "SUPER", "THIS"}

// abstractScope_ImplicitVar_name is the field ImplicitVar.name.
var abstractScope_ImplicitVar_name = [...]string{"arguments", "exports", "super", "this"}

var (
	implicitArgumentsID = rhino.RhinoStringPool_InternID("arguments")
	implicitExportsID   = rhino.RhinoStringPool_InternID("exports")
	implicitSuperID     = rhino.RhinoStringPool_InternID("super")
	implicitThisID      = rhino.RhinoStringPool_InternID("this")

	abstractScope_ImplicitVar_id = [...]uint32{
		implicitArgumentsID,
		implicitExportsID,
		implicitSuperID,
		implicitThisID,
	}
)

// Go-only: language-adapter — a Java ImplicitVar reference may be null; such
// values are *AbstractScope_ImplicitVar here (nil is null), always pointing
// into this table so that no call allocates.
var abstractScope_ImplicitVar_ptrs = [...]AbstractScope_ImplicitVar{
	AbstractScope_ImplicitVar_ARGUMENTS, AbstractScope_ImplicitVar_EXPORTS,
	AbstractScope_ImplicitVar_SUPER, AbstractScope_ImplicitVar_THIS,
}

// String is Enum.name().
func (v AbstractScope_ImplicitVar) String() string { return abstractScope_ImplicitVar_names[v] }

// Ordinal is Enum.ordinal().
func (v AbstractScope_ImplicitVar) Ordinal() int { return int(v) }

// name is the field ImplicitVar.name.
func (v AbstractScope_ImplicitVar) name() string { return abstractScope_ImplicitVar_name[v] }

func (v AbstractScope_ImplicitVar) nameID() uint32 { return abstractScope_ImplicitVar_id[v] }

// ptr returns the canonical non-null *AbstractScope_ImplicitVar for v.
func (v AbstractScope_ImplicitVar) ptr() *AbstractScope_ImplicitVar {
	return &abstractScope_ImplicitVar_ptrs[v]
}

// IsMadeByScope reports whether this kind of implicit variable is
// created/owned by the given scope.
//
// Java: AbstractScope.ImplicitVar.isMadeByScope(AbstractScope<?, ?>)
func (v AbstractScope_ImplicitVar) IsMadeByScope(scope AbstractScope) bool {
	switch v {
	case AbstractScope_ImplicitVar_EXPORTS:
		return scope.IsModuleScope() &&
			scope.GetRootNode().GetParent().GetBooleanProp(rhino.Node_GOOG_MODULE)
	case AbstractScope_ImplicitVar_SUPER, AbstractScope_ImplicitVar_THIS:
		return scope.IsStaticBlockScope() ||
			NodeUtil_IsNonArrowFunction(scope.GetRootNode()) ||
			scope.IsMemberFieldDefScope() ||
			scope.IsComputedFieldDefRhsScope()
	case AbstractScope_ImplicitVar_ARGUMENTS:
		return NodeUtil_IsNonArrowFunction(scope.GetRootNode())
	}
	panic(jutil.NewIllegalStateException())
}

// Java: AbstractScope.ImplicitVar.of(String)
func AbstractScope_ImplicitVar_Of(name string) *AbstractScope_ImplicitVar {
	switch len(name) {
	case 4:
		if name == "this" {
			return &abstractScope_ImplicitVar_ptrs[AbstractScope_ImplicitVar_THIS]
		}
	case 5:
		if name == "super" {
			return &abstractScope_ImplicitVar_ptrs[AbstractScope_ImplicitVar_SUPER]
		}
	case 7:
		if name == "exports" {
			return &abstractScope_ImplicitVar_ptrs[AbstractScope_ImplicitVar_EXPORTS]
		}
	case 9:
		if name == "arguments" {
			return &abstractScope_ImplicitVar_ptrs[AbstractScope_ImplicitVar_ARGUMENTS]
		}
	}
	return nil
}

func AbstractScope_ImplicitVar_OfID(id uint32) *AbstractScope_ImplicitVar {
	switch id {
	case implicitArgumentsID:
		return &abstractScope_ImplicitVar_ptrs[AbstractScope_ImplicitVar_ARGUMENTS]
	case implicitThisID:
		return &abstractScope_ImplicitVar_ptrs[AbstractScope_ImplicitVar_THIS]
	case implicitExportsID:
		return &abstractScope_ImplicitVar_ptrs[AbstractScope_ImplicitVar_EXPORTS]
	case implicitSuperID:
		return &abstractScope_ImplicitVar_ptrs[AbstractScope_ImplicitVar_SUPER]
	}
	return nil
}
