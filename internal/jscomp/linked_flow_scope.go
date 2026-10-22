// Port of com/google/javascript/jscomp/LinkedFlowScope.java @ closure-compiler v20261015

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
	"github.com/paullepoulpe/jsgo/internal/base"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

type linkedFlowScopeArena struct {
	scopeBuf       [64]LinkedFlowScope
	scopeInlineLen int
	scopeChunks    []*[256]LinkedFlowScope
	scopeChunkIdx  int
	scopeOffset    int

	overlayBuf       [32]LinkedFlowScope_OverlayScope
	overlayInlineLen int
	overlayChunks    []*[256]LinkedFlowScope_OverlayScope
	overlayChunkIdx  int
	overlayOffset    int

	slotBuf       [64]LinkedFlowScope_OverlaySlot
	slotInlineLen int
	slotChunks    []*[256]LinkedFlowScope_OverlaySlot
	slotChunkIdx  int
	slotOffset    int

	scopeMapArena rhino.HamtPMapArena[*TypedScope, *LinkedFlowScope_OverlayScope]
	slotMapArena  rhino.HamtPMapArena[string, *LinkedFlowScope_OverlaySlot]
}

func (a *linkedFlowScopeArena) scopeMaps() *rhino.HamtPMapArena[*TypedScope, *LinkedFlowScope_OverlayScope] {
	if a == nil {
		return nil
	}
	return &a.scopeMapArena
}

func (a *linkedFlowScopeArena) slotMaps() *rhino.HamtPMapArena[string, *LinkedFlowScope_OverlaySlot] {
	if a == nil {
		return nil
	}
	return &a.slotMapArena
}

func (a *linkedFlowScopeArena) reset() {
	a.scopeMapArena.Reset()
	a.slotMapArena.Reset()
	if a.scopeInlineLen > 0 {
		clear(a.scopeBuf[:a.scopeInlineLen])
		a.scopeInlineLen = 0
	}
	if a.scopeChunkIdx > 0 || a.scopeOffset > 0 {
		for i := 0; i < a.scopeChunkIdx; i++ {
			clear(a.scopeChunks[i][:])
		}
		if a.scopeOffset > 0 && a.scopeChunkIdx < len(a.scopeChunks) {
			clear(a.scopeChunks[a.scopeChunkIdx][:a.scopeOffset])
		}
		a.scopeChunkIdx = 0
		a.scopeOffset = 0
	}

	if a.overlayInlineLen > 0 {
		clear(a.overlayBuf[:a.overlayInlineLen])
		a.overlayInlineLen = 0
	}
	if a.overlayChunkIdx > 0 || a.overlayOffset > 0 {
		for i := 0; i < a.overlayChunkIdx; i++ {
			clear(a.overlayChunks[i][:])
		}
		if a.overlayOffset > 0 && a.overlayChunkIdx < len(a.overlayChunks) {
			clear(a.overlayChunks[a.overlayChunkIdx][:a.overlayOffset])
		}
		a.overlayChunkIdx = 0
		a.overlayOffset = 0
	}

	if a.slotInlineLen > 0 {
		clear(a.slotBuf[:a.slotInlineLen])
		a.slotInlineLen = 0
	}
	if a.slotChunkIdx > 0 || a.slotOffset > 0 {
		for i := 0; i < a.slotChunkIdx; i++ {
			clear(a.slotChunks[i][:])
		}
		if a.slotOffset > 0 && a.slotChunkIdx < len(a.slotChunks) {
			clear(a.slotChunks[a.slotChunkIdx][:a.slotOffset])
		}
		a.slotChunkIdx = 0
		a.slotOffset = 0
	}
}

func (a *linkedFlowScopeArena) trimLargeChunks() {
	a.reset()
	const maxRetained = 8
	if len(a.scopeChunks) > maxRetained {
		clear(a.scopeChunks[maxRetained:])
		a.scopeChunks = a.scopeChunks[:maxRetained]
	}
	if len(a.overlayChunks) > maxRetained {
		clear(a.overlayChunks[maxRetained:])
		a.overlayChunks = a.overlayChunks[:maxRetained]
	}
	if len(a.slotChunks) > maxRetained {
		clear(a.slotChunks[maxRetained:])
		a.slotChunks = a.slotChunks[:maxRetained]
	}
	a.scopeMapArena.TrimLargeChunks()
	a.slotMapArena.TrimLargeChunks()
}

func (a *linkedFlowScopeArena) allocScope() *LinkedFlowScope {
	if a.scopeInlineLen < len(a.scopeBuf) {
		s := &a.scopeBuf[a.scopeInlineLen]
		a.scopeInlineLen++
		return s
	}
	return a.allocScopeChunk()
}

//go:noinline
func (a *linkedFlowScopeArena) allocScopeChunk() *LinkedFlowScope {
	if a.scopeOffset == 256 {
		a.scopeChunkIdx++
		a.scopeOffset = 0
	}
	if a.scopeChunkIdx == len(a.scopeChunks) {
		a.scopeChunks = append(a.scopeChunks, new([256]LinkedFlowScope))
	}
	s := &a.scopeChunks[a.scopeChunkIdx][a.scopeOffset]
	a.scopeOffset++
	return s
}

func (a *linkedFlowScopeArena) allocOverlay() *LinkedFlowScope_OverlayScope {
	if a.overlayInlineLen < len(a.overlayBuf) {
		o := &a.overlayBuf[a.overlayInlineLen]
		a.overlayInlineLen++
		return o
	}
	return a.allocOverlayChunk()
}

//go:noinline
func (a *linkedFlowScopeArena) allocOverlayChunk() *LinkedFlowScope_OverlayScope {
	if a.overlayOffset == 256 {
		a.overlayChunkIdx++
		a.overlayOffset = 0
	}
	if a.overlayChunkIdx == len(a.overlayChunks) {
		a.overlayChunks = append(a.overlayChunks, new([256]LinkedFlowScope_OverlayScope))
	}
	o := &a.overlayChunks[a.overlayChunkIdx][a.overlayOffset]
	a.overlayOffset++
	return o
}

func (a *linkedFlowScopeArena) allocSlot() *LinkedFlowScope_OverlaySlot {
	if a.slotInlineLen < len(a.slotBuf) {
		s := &a.slotBuf[a.slotInlineLen]
		a.slotInlineLen++
		return s
	}
	return a.allocSlotChunk()
}

//go:noinline
func (a *linkedFlowScopeArena) allocSlotChunk() *LinkedFlowScope_OverlaySlot {
	if a.slotOffset == 256 {
		a.slotChunkIdx++
		a.slotOffset = 0
	}
	if a.slotChunkIdx == len(a.slotChunks) {
		a.slotChunks = append(a.slotChunks, new([256]LinkedFlowScope_OverlaySlot))
	}
	s := &a.slotChunks[a.slotChunkIdx][a.slotOffset]
	a.slotOffset++
	return s
}

// LinkedFlowScope is a flow scope that tries to store as little symbol
// information as possible, instead delegating to its parents. Optimized for
// low memory use.
type LinkedFlowScope struct {
	arena         *linkedFlowScopeArena
	inputProvider CompilerInputProvider

	// Map from TypedScope to OverlayScope.
	scopes *rhino.HamtPMap[*TypedScope, *LinkedFlowScope_OverlayScope]

	functionScope *TypedScope

	// The TypedScope for the block that this flow scope is defined for.
	syntacticScope *TypedScope
}

// NewLinkedFlowScope creates a flow scope without a direct parent. This can
// happen in three cases: (1) the "bottom" scope for a CFG root, (2) a direct
// child of a parent at the maximum depth, or (3) a joined scope with more
// than one direct parent. The parent is non-null only in the second case.
//
// Java: LinkedFlowScope(CompilerInputProvider, PMap<TypedScope, OverlayScope>, TypedScope, TypedScope)
func NewLinkedFlowScope(
	inputProvider CompilerInputProvider,
	scopes rhino.PMap[*TypedScope, *LinkedFlowScope_OverlayScope],
	syntacticScope *TypedScope,
	functionScope *TypedScope) *LinkedFlowScope {
	return newLinkedFlowScope(
		inputProvider,
		scopes.(*rhino.HamtPMap[*TypedScope, *LinkedFlowScope_OverlayScope]),
		syntacticScope,
		functionScope)
}

func newLinkedFlowScope(
	inputProvider CompilerInputProvider,
	scopes *rhino.HamtPMap[*TypedScope, *LinkedFlowScope_OverlayScope],
	syntacticScope *TypedScope,
	functionScope *TypedScope) *LinkedFlowScope {
	return newLinkedFlowScopeWithArena(nil, inputProvider, scopes, syntacticScope, functionScope)
}

func newLinkedFlowScopeWithArena(
	arena *linkedFlowScopeArena,
	inputProvider CompilerInputProvider,
	scopes *rhino.HamtPMap[*TypedScope, *LinkedFlowScope_OverlayScope],
	syntacticScope *TypedScope,
	functionScope *TypedScope) *LinkedFlowScope {
	var s *LinkedFlowScope
	if arena != nil {
		s = arena.allocScope()
	} else {
		s = &LinkedFlowScope{}
	}
	s.arena = arena
	s.inputProvider = inputProvider
	s.scopes = scopes
	s.syntacticScope = syntacticScope
	s.functionScope = functionScope
	return s
}

// TrimScopes returns the scope map, trimmed to the common ancestor between
// this FlowScope's syntacticScope and the given scope. Any inferred types on
// variables in deeper scopes cannot be propagated past this point (since
// they're no longer in scope), and trimming them eagerly allows us to ignore
// these irrelevant types when checking equality and joining.
//
// Java: LinkedFlowScope.trimScopes(TypedScope)
func (l *LinkedFlowScope) TrimScopes(scope *TypedScope) rhino.PMap[*TypedScope, *LinkedFlowScope_OverlayScope] {
	return l.trimScopesMap(scope)
}

func (l *LinkedFlowScope) trimScopesMap(scope *TypedScope) *rhino.HamtPMap[*TypedScope, *LinkedFlowScope_OverlayScope] {
	thisScope := l.syntacticScope
	thatScope := scope
	if thisScope == thatScope || l.scopes.IsEmpty() {
		return l.scopes
	}
	thisDepth := thisScope.GetDepth()
	thatDepth := thatScope.GetDepth()
	result := l.scopes
	for thatDepth > thisDepth {
		thatScope = thatScope.GetParent()
		thatDepth--
	}
	scopeMaps := l.arena.scopeMaps()
	for thisDepth > thatDepth {
		result = result.MinusMapWithArena(scopeMaps, thisScope)
		if result.IsEmpty() {
			return result
		}
		thisScope = thisScope.GetParent()
		thisDepth--
	}
	for thisScope != thatScope && thisScope != nil && thatScope != nil {
		result = result.MinusMapWithArena(scopeMaps, thisScope)
		if result.IsEmpty() {
			return result
		}
		thisScope = thisScope.GetParent()
		thatScope = thatScope.GetParent()
	}
	return result
}

// FlowsFromBottom: whether this flows from a bottom scope.
//
// Java: LinkedFlowScope.flowsFromBottom()
func (l *LinkedFlowScope) FlowsFromBottom() bool {
	return l.functionScope.IsBottom()
}

// LinkedFlowScope_CreateEntryLattice creates an entry lattice for the flow.
//
// Java: LinkedFlowScope.createEntryLattice(CompilerInputProvider, TypedScope)
func LinkedFlowScope_CreateEntryLattice(
	inputProvider CompilerInputProvider, scope *TypedScope) *LinkedFlowScope {
	return newLinkedFlowScope(
		inputProvider, LinkedFlowScope_EMPTY_SCOPES, scope, scope)
}

// Java: LinkedFlowScope.inferSlotType(String, JSType)
func (l *LinkedFlowScope) InferSlotType(symbol string, type_ rhino.JSType) FlowScope {
	scope := l.GetOverlayScopeForName(symbol, true)
	newScope := scope.Infer(symbol, type_)
	// Aggressively remove empty scopes to maintain a reasonable equivalence.
	var newScopes *rhino.HamtPMap[*TypedScope, *LinkedFlowScope_OverlayScope]
	if !newScope.slots.IsEmpty() {
		newScopes = l.scopes.PlusMapWithArena(l.arena.scopeMaps(), scope.scope, newScope)
	} else {
		newScopes = l.scopes.MinusMapWithArena(l.arena.scopeMaps(), scope.scope)
	}
	if newScopes != l.scopes {
		return newLinkedFlowScopeWithArena(l.arena, l.inputProvider, newScopes, l.syntacticScope, l.functionScope)
	}
	return l
}

// Java: LinkedFlowScope.inferQualifiedSlot(Node, String, JSType, JSType, boolean)
func (l *LinkedFlowScope) InferQualifiedSlot(
	node *rhino.Node, symbol string, bottomType rhino.JSType, inferredType rhino.JSType, declared bool) FlowScope {
	if l.functionScope.IsGlobal() {
		// Do not infer qualified names on the global scope.  Ideally these would be
		// added to the scope by TypedScopeCreator, but if they are not, adding them
		// here causes scaling problems (large projects can have tens of thousands of
		// undeclared qualified names in the global scope) with no real benefit.
		return l
	}
	v := l.syntacticScope.GetVar(symbol)
	if v == nil && !l.functionScope.IsBottom() {
		// NOTE(sdh): Qualified names are declared on scopes lazily via this method.
		// The difficulty is that it's not always clear which scope they need to be
		// defined on.  In particular, syntacticScope is wrong because it is often a
		// nested block scope that is ignored when branches are joined; functionScope
		// is also wrong because it could lead to ambiguity if the same root name is
		// declared in multiple different blocks.  Instead, the qualified name is declared
		// on the scope that owns the root, when possible. When the root is undeclared, the qualified
		// name is declared in the global scope, as only global variables can be undeclared.
		rootVar := l.syntacticScope.GetVar(LinkedFlowScope_GetRootOfQualifiedName(symbol))
		var rootScope *TypedScope
		if rootVar != nil {
			rootScope = rootVar.GetScope().(*TypedScope)
		} else {
			rootScope = l.syntacticScope.GetGlobalScope()
		}
		v =
			rootScope.Declare(
				symbol,
				node,
				bottomType,
				l.inputProvider.GetInput(NodeUtil_GetInputId(node)),
				!declared)
	}

	var declaredType rhino.JSType
	if v != nil {
		declaredType = v.GetType()
	}
	if v != nil {
		if !v.IsTypeInferred() {
			// Use the inferred type over the declared type only if the
			// inferred type is a strict subtype of the declared type.
			if declaredType == nil ||
				!inferredType.IsSubtypeOf(declaredType) ||
				declaredType.IsSubtypeOf(inferredType) ||
				inferredType.Equals(declaredType) {
				return l
			}
		} else if declaredType != nil && !inferredType.IsSubtypeOf(declaredType) {
			// If this inferred type is incompatible with another type previously
			// inferred and stored on the scope, then update the scope.
			v.SetType(v.GetType().GetLeastSupertype(inferredType))
		}
	}
	return l.InferSlotType(symbol, inferredType)
}

// Java: LinkedFlowScope.getTypeOfThis()
func (l *LinkedFlowScope) GetTypeOfThis() rhino.JSType {
	return l.syntacticScope.GetTypeOfThis()
}

// Java: LinkedFlowScope.getRootNode()
func (l *LinkedFlowScope) GetRootNode() *rhino.Node {
	return l.syntacticScope.GetRootNode()
}

// Java: LinkedFlowScope.getParentScope()
func (l *LinkedFlowScope) GetParentScope() rhino.StaticScope {
	panic(jutil.NewUnsupportedOperationException())
}

// GetSlot gets the slot for the given symbol (a StaticTypedSlot, or nil).
//
// Java: LinkedFlowScope.getSlot(String)
func (l *LinkedFlowScope) GetSlot(name string) rhino.StaticSlot {
	var_ := l.syntacticScope.GetVar(name)
	var scope *LinkedFlowScope_OverlayScope
	if var_ == nil {
		scope = l.GetOverlayScopeForName(name, false)
	} else {
		scope = l.GetOverlayScopeForScope(var_.GetScope().(*TypedScope), false)
	}

	if scope != nil {
		if slot := scope.GetSlot(name); slot != nil {
			return slot
		}
		return nil
	}
	if var_ == nil {
		return nil
	}
	return var_
}

// Java: LinkedFlowScope.getRootOfQualifiedName(String)
func LinkedFlowScope_GetRootOfQualifiedName(name string) string {
	index := jutil.StringIndexOfChar(name, '.')
	if index < 0 {
		return name
	}
	return jutil.StringSubstring(name, 0, index)
}

// GetOverlayScopeForName returns the overlay scope corresponding to this
// qualified name
//
// create says whether to create a new OverlayScope if one does not already
// exist.
//
// Java: LinkedFlowScope.getOverlayScopeForName(String, boolean)
func (l *LinkedFlowScope) GetOverlayScopeForName(name string, create bool) *LinkedFlowScope_OverlayScope {
	rootVar := l.syntacticScope.GetVar(LinkedFlowScope_GetRootOfQualifiedName(name))
	var scope *TypedScope
	if rootVar != nil {
		scope = rootVar.GetScope().(*TypedScope)
	}
	if scope == nil {
		scope = l.functionScope
	}
	return l.GetOverlayScopeForScope(scope, create)
}

// GetOverlayScopeForScope returns the overlay scope corresponding to this
// syntactic scope.
//
// Use instead of getOverlayScopeForName(String, boolean) if you already know
// the correct scope in order to avoid a variable lookup.
//
// scope is the syntactic scope; create says whether to create a new
// OverlayScope if one does not already exist.
//
// Java: LinkedFlowScope.getOverlayScopeForScope(TypedScope, boolean)
func (l *LinkedFlowScope) GetOverlayScopeForScope(scope *TypedScope, create bool) *LinkedFlowScope_OverlayScope {
	overlay := l.scopes.Get(scope)
	if overlay == nil && create {
		overlay = newLinkedFlowScope_OverlayScopeWithArena(l.arena, scope, LinkedFlowScope_EMPTY_SLOTS)
	}
	return overlay
}

// Java: LinkedFlowScope.getOwnSlot(String)
func (l *LinkedFlowScope) GetOwnSlot(name string) rhino.StaticSlot {
	panic(jutil.NewUnsupportedOperationException())
}

// GetTopmostScopeOfEventualDeclaration is StaticScope's default method,
// which LinkedFlowScope does not override.
//
// Java: StaticScope.getTopmostScopeOfEventualDeclaration(String) (default method)
func (l *LinkedFlowScope) GetTopmostScopeOfEventualDeclaration(name string) rhino.StaticScope {
	return rhino.StaticScope_GetTopmostScopeOfEventualDeclaration(l, name)
}

// LookupQualifiedName is StaticTypedScope's default method, which
// LinkedFlowScope does not override.
//
// Java: StaticTypedScope.lookupQualifiedName(QualifiedName) (default method)
func (l *LinkedFlowScope) LookupQualifiedName(qname rhino.QualifiedName) rhino.JSType {
	return rhino.StaticTypedScope_LookupQualifiedName(l, qname)
}

// Java: LinkedFlowScope.withSyntacticScope(StaticTypedScope)
func (l *LinkedFlowScope) WithSyntacticScope(scope rhino.StaticTypedScope) FlowScope {
	typedScope := scope.(*TypedScope)
	if typedScope != l.syntacticScope {
		return newLinkedFlowScopeWithArena(l.arena, l.inputProvider, l.trimScopesMap(typedScope), typedScope, l.functionScope)
	}
	return l
}

// GetDeclarationScope returns the syntactic scope (a *TypedScope).
//
// Java: LinkedFlowScope.getDeclarationScope()
func (l *LinkedFlowScope) GetDeclarationScope() rhino.StaticTypedScope {
	if l.syntacticScope == nil {
		return nil
	}
	return l.syntacticScope
}

// LinkedFlowScope_FlowScopeJoinOp joins the two FlowScopes.
type LinkedFlowScope_FlowScopeJoinOp struct {
	// May be nil.
	result        *LinkedFlowScope
	inputProvider CompilerInputProvider
}

// Java: LinkedFlowScope.FlowScopeJoinOp(CompilerInputProvider)
func NewLinkedFlowScope_FlowScopeJoinOp(inputProvider CompilerInputProvider) *LinkedFlowScope_FlowScopeJoinOp {
	j := &LinkedFlowScope_FlowScopeJoinOp{}
	j.result = nil
	j.inputProvider = inputProvider
	return j
}

// JoinFlow: NOTE(sdh): When joining flow scopes with different syntactic
// scopes, we do not attempt to recover the correct syntactic scope. This is
// okay because joins only occur in two situations: (1) performed by the
// DataFlowAnalysis class automatically between CFG nodes, and (2) requested
// manually while traversing a single expression within a CFG node. The
// syntactic scope is always set at the beginning of flowing through a CFG
// node. In the case of (1), the join result's syntactic scope is immediately
// replaced with the correct one when we flow through the next node. In the
// case of (2), both inputs will always have the same syntactic scope. So
// simply propagating either input's scope is perfectly fine.
//
// Java: LinkedFlowScope.FlowScopeJoinOp.joinFlow(FlowScope)
func (j *LinkedFlowScope_FlowScopeJoinOp) JoinFlow(input FlowScope) {
	// To join the two scopes, we have to
	linkedInput := input.(*LinkedFlowScope)
	if j.result == nil {
		j.result = linkedInput
		return
	} else if j.result.scopes == linkedInput.scopes &&
		j.result.functionScope == linkedInput.functionScope {
		return
	} else if j.result.FlowsFromBottom() && j.result.scopes.IsEmpty() {
		j.result = linkedInput
		return
	} else if linkedInput.FlowsFromBottom() && linkedInput.scopes.IsEmpty() {
		return
	}

	// NOTE: it would be nice to put 'null' as the syntactic scope if they're not
	// equal, but this is not currently feasible.  For joins that occur within a
	// single CFG node's flow, it's irrelevant, but for joins between separate
	// CFG nodes, there is *one* place where the syntactic scope is actually used:
	// when joining more than two scopes, the first two scopes are joined, and
	// then the join result is joined with the third.  When joining, we look up
	// the types (and existence) of vars in one scope in the other; so when a var
	// from the third scope (say, a local) is missing from the join result, it
	// looks through the syntactic scope before realizing  this.  A quick fix
	// might be to just check that the scope is non-null before trying to join;
	// a better long-term fix would be to improve how we do joins to avoid
	// excessive map entry creation: find a common ancestor, etc.  One
	// interesting consequence of the current approach is that we may end up
	// adding irrelevant block-local variables to the joined scope unnecessarily.
	common := LinkedFlowScope_GetCommonParentDeclarationScope(j.result, linkedInput)
	var functionScope *TypedScope
	if j.result.FlowsFromBottom() {
		functionScope = linkedInput.functionScope
	} else {
		functionScope = j.result.functionScope
	}
	joinedScopes := linkedFlowScope_JoinMap(j.result, linkedInput, common)
	if joinedScopes == j.result.scopes &&
		common == j.result.syntacticScope &&
		functionScope == j.result.functionScope {
		return
	}
	if joinedScopes == linkedInput.scopes &&
		common == linkedInput.syntacticScope &&
		functionScope == linkedInput.functionScope {
		j.result = linkedInput
		return
	}
	arena := j.result.arena
	if arena == nil {
		arena = linkedInput.arena
	}
	j.result =
		newLinkedFlowScopeWithArena(
			arena,
			j.inputProvider,
			joinedScopes,
			common,
			functionScope)
}

// Java: LinkedFlowScope.FlowScopeJoinOp.finish()
func (j *LinkedFlowScope_FlowScopeJoinOp) Finish() FlowScope {
	if j.result == nil {
		return nil
	}
	return j.result
}

// Java: LinkedFlowScope.getCommonParentDeclarationScope(LinkedFlowScope, LinkedFlowScope)
func LinkedFlowScope_GetCommonParentDeclarationScope(left *LinkedFlowScope, right *LinkedFlowScope) *TypedScope {
	if left.FlowsFromBottom() {
		return right.syntacticScope
	} else if right.FlowsFromBottom() {
		return left.syntacticScope
	}
	return left.syntacticScope.GetCommonParent(right.syntacticScope)
}

// Java: LinkedFlowScope.equals(Object)
func (l *LinkedFlowScope) Equals(other any) bool {
	that, ok := other.(*LinkedFlowScope)
	if !ok || that == nil {
		return false
	}
	if l == that {
		return true
	}

	// If two flow scopes are in the same function, then they could have
	// two possible function scopes: the real one and the BOTTOM scope.
	// If they have different function scopes, we *should* iterate through all
	// the variables in each scope and compare. However, 99.9% of the time,
	// they're not equal. And the other .1% of the time, we can pretend
	// they're equal--this just means that data flow analysis will have
	// to propagate the entry lattice a little bit further than it
	// really needs to. Everything will still come out ok.
	return l.functionScope == that.functionScope &&
		l.scopes.EquivalentMap(that.scopes, LinkedFlowScope_EqualScopes)
}

// Java: LinkedFlowScope.equalScopes(OverlayScope, OverlayScope)
func LinkedFlowScope_EqualScopes(left *LinkedFlowScope_OverlayScope, right *LinkedFlowScope_OverlayScope) bool {
	if left == right {
		return true
	}
	return left.slots.EquivalentMap(right.slots, LinkedFlowScope_EqualSlots)
}

// LinkedFlowScope_EqualSlots determines whether two slots are meaningfully
// different for the purposes of data flow analysis.
//
// Java: LinkedFlowScope.equalSlots(StaticTypedSlot, StaticTypedSlot)
func LinkedFlowScope_EqualSlots(slotA *LinkedFlowScope_OverlaySlot, slotB *LinkedFlowScope_OverlaySlot) bool {
	return slotA == slotB || !slotA.GetType().DiffersFrom(slotB.GetType())
}

// Java: LinkedFlowScope.hashCode()
func (l *LinkedFlowScope) HashCode() int32 {
	panic(jutil.NewUnsupportedOperationException())
}

// Java: LinkedFlowScope.join(LinkedFlowScope, LinkedFlowScope, TypedScope)
func LinkedFlowScope_Join(
	linkedA *LinkedFlowScope, linkedB *LinkedFlowScope, commonParent *TypedScope) rhino.PMap[*TypedScope, *LinkedFlowScope_OverlayScope] {
	return linkedFlowScope_JoinMap(linkedA, linkedB, commonParent)
}

func linkedFlowScope_JoinMap(
	linkedA *LinkedFlowScope, linkedB *LinkedFlowScope, commonParent *TypedScope) *rhino.HamtPMap[*TypedScope, *LinkedFlowScope_OverlayScope] {
	scopesA := linkedA.trimScopesMap(commonParent)
	scopesB := linkedB.trimScopesMap(commonParent)
	if scopesA == scopesB && linkedA.FlowsFromBottom() == linkedB.FlowsFromBottom() {
		return scopesA
	}
	arena := linkedA.arena
	if arena == nil {
		arena = linkedB.arena
	}
	return scopesA.
		ReconcileMapWithArena(
			arena.scopeMaps(),
			scopesB,
			func(scopeKey *TypedScope, scopeA *LinkedFlowScope_OverlayScope, scopeB *LinkedFlowScope_OverlayScope) *LinkedFlowScope_OverlayScope {
				if scopeA == scopeB && linkedA.FlowsFromBottom() == linkedB.FlowsFromBottom() {
					return scopeA
				}
				var slotsA, slotsB *rhino.HamtPMap[string, *LinkedFlowScope_OverlaySlot]
				if scopeA != nil {
					slotsA = scopeA.slots
				} else {
					slotsA = LinkedFlowScope_EMPTY_SLOTS
				}
				if scopeB != nil {
					slotsB = scopeB.slots
				} else {
					slotsB = LinkedFlowScope_EMPTY_SLOTS
				}
				// TODO(sdh): Simplify this logic: we want the best non-bottom scope we can get,
				// for the purpose of (a) passing to the joined OverlayScope constructor, and
				// (b) joining types only present in one scope.
				var typedScopeA *TypedScope
				if !linkedA.FlowsFromBottom() {
					if scopeA != nil {
						typedScopeA = scopeA.scope
					} else {
						typedScopeA = scopeB.scope
					}
				}
				var typedScopeB *TypedScope
				if !linkedB.FlowsFromBottom() {
					if scopeB != nil {
						typedScopeB = scopeB.scope
					} else {
						typedScopeB = scopeA.scope
					}
				}
				bestScope := typedScopeA
				if bestScope == nil {
					bestScope = typedScopeB
				}
				if bestScope == nil {
					if scopeA != nil {
						bestScope = scopeA.scope
					} else {
						bestScope = scopeB.scope
					}
				}
				var reconciledSlots *rhino.HamtPMap[string, *LinkedFlowScope_OverlaySlot]
				if slotsA == slotsB && typedScopeA == typedScopeB {
					reconciledSlots = slotsA
				} else {
					reconciledSlots = slotsA.ReconcileMapWithArena(
						arena.slotMaps(),
						slotsB,
						func(slotKey string, slotA *LinkedFlowScope_OverlaySlot, slotB *LinkedFlowScope_OverlaySlot) *LinkedFlowScope_OverlaySlot {
							// There are 5 different join cases:
							// 1) The type is present in joinedScopeA, not in joinedScopeB,
							//    and not in functionScope. Just use the one in A.
							// 2) The type is present in joinedScopeB, not in joinedScopeA,
							//    and not in functionScope. Just use the one in B.
							// 3) The type is present in functionScope and joinedScopeA, but
							//    not in joinedScopeB. Join the two types.
							// 4) The type is present in functionScope and joinedScopeB, but
							//    not in joinedScopeA. Join the two types.
							// 5) The type is present in joinedScopeA and joinedScopeB. Join
							//    the two types.
							var name string
							if slotA != nil {
								name = slotA.GetName()
							} else {
								name = slotB.GetName()
							}
							if slotB == nil || slotB.GetType() == nil {
								var fnSlot *TypedVar
								if typedScopeB != nil {
									fnSlot, _ = typedScopeB.GetSlot(name).(*TypedVar)
								}
								var fnSlotType rhino.JSType
								if fnSlot != nil {
									fnSlotType = fnSlot.GetType()
								}
								if fnSlotType == nil || base.JSCompObjects_Identical(fnSlotType, slotA.GetType()) {
									// Case #1
									return slotA
								} else {
									// Case #3
									joinedType := slotA.GetType().GetLeastSupertype(fnSlotType)
									if base.JSCompObjects_Identical(joinedType, slotA.GetType()) {
										return slotA
									}
									return newLinkedFlowScope_OverlaySlotWithArena(arena, name, joinedType)
								}
							} else if slotA == nil || slotA.GetType() == nil {
								var fnSlot *TypedVar
								if typedScopeA != nil {
									fnSlot, _ = typedScopeA.GetSlot(name).(*TypedVar)
								}
								var fnSlotType rhino.JSType
								if fnSlot != nil {
									fnSlotType = fnSlot.GetType()
								}
								if fnSlotType == nil || base.JSCompObjects_Identical(fnSlotType, slotB.GetType()) {
									// Case #2
									return slotB
								} else {
									// Case #4
									joinedType := slotB.GetType().GetLeastSupertype(fnSlotType)
									if base.JSCompObjects_Identical(joinedType, slotB.GetType()) {
										return slotB
									}
									return newLinkedFlowScope_OverlaySlotWithArena(arena, name, joinedType)
								}
							}
							// Case #5
							if base.JSCompObjects_Identical(slotA.GetType(), slotB.GetType()) {
								return slotA
							}
							joinedType := slotA.GetType().GetLeastSupertype(slotB.GetType())
							if base.JSCompObjects_Identical(joinedType, slotA.GetType()) {
								return slotA
							}
							if base.JSCompObjects_Identical(joinedType, slotB.GetType()) {
								return slotB
							}
							return newLinkedFlowScope_OverlaySlotWithArena(arena, name, joinedType)
						})
				}
				if scopeA != nil && scopeA.scope == bestScope && reconciledSlots == scopeA.slots {
					return scopeA
				}
				if scopeB != nil && scopeB.scope == bestScope && reconciledSlots == scopeB.slots {
					return scopeB
				}
				return newLinkedFlowScope_OverlayScopeWithArena(arena, bestScope, reconciledSlots)
			})
}

// LinkedFlowScope_OverlayScope is the overlay of one syntactic scope.
type LinkedFlowScope_OverlayScope struct {
	arena *linkedFlowScopeArena
	scope *TypedScope
	slots *rhino.HamtPMap[string, *LinkedFlowScope_OverlaySlot]
}

// Equals implements jutil.Equaler using pointer identity (matching Java Object.equals).
func (o *LinkedFlowScope_OverlayScope) Equals(other any) bool {
	that, ok := other.(*LinkedFlowScope_OverlayScope)
	return ok && o == that
}

// Java: LinkedFlowScope.OverlayScope(TypedScope)
func NewLinkedFlowScope_OverlayScope(scope *TypedScope) *LinkedFlowScope_OverlayScope {
	return newLinkedFlowScope_OverlayScopeWithSlots(scope, LinkedFlowScope_EMPTY_SLOTS)
}

// Java: LinkedFlowScope.OverlayScope(TypedScope, PMap<String, OverlaySlot>)
func NewLinkedFlowScope_OverlayScopeWithSlots(
	scope *TypedScope, slots rhino.PMap[string, *LinkedFlowScope_OverlaySlot]) *LinkedFlowScope_OverlayScope {
	return newLinkedFlowScope_OverlayScopeWithSlots(
		scope, slots.(*rhino.HamtPMap[string, *LinkedFlowScope_OverlaySlot]))
}

func newLinkedFlowScope_OverlayScopeWithSlots(
	scope *TypedScope, slots *rhino.HamtPMap[string, *LinkedFlowScope_OverlaySlot]) *LinkedFlowScope_OverlayScope {
	return newLinkedFlowScope_OverlayScopeWithArena(nil, scope, slots)
}

func newLinkedFlowScope_OverlayScopeWithArena(
	arena *linkedFlowScopeArena,
	scope *TypedScope,
	slots *rhino.HamtPMap[string, *LinkedFlowScope_OverlaySlot]) *LinkedFlowScope_OverlayScope {
	jutil.CheckNotNull(scope)
	var o *LinkedFlowScope_OverlayScope
	if arena != nil {
		o = arena.allocOverlay()
	} else {
		o = &LinkedFlowScope_OverlayScope{}
	}
	o.arena = arena
	o.scope = scope
	o.slots = slots
	return o
}

// Java: LinkedFlowScope.OverlayScope.infer(String, JSType)
func (o *LinkedFlowScope_OverlayScope) Infer(name string, type_ rhino.JSType) *LinkedFlowScope_OverlayScope {
	// TODO(sdh): variants that do or don't clobber properties (i.e. look up and modify instead)
	slot := o.slots.Get(name)
	if slot != nil && base.JSCompObjects_Identical(type_, slot.type_) {
		return o
	}
	return newLinkedFlowScope_OverlayScopeWithArena(o.arena, o.scope, o.slots.PlusMapWithArena(o.arena.slotMaps(), name, newLinkedFlowScope_OverlaySlotWithArena(o.arena, name, type_)))
}

// GetSlot returns the overlay's slot for name, else the scope's (nil if
// there is none).
//
// Java: LinkedFlowScope.OverlayScope.getSlot(String)
func (o *LinkedFlowScope_OverlayScope) GetSlot(name string) rhino.StaticTypedSlot {
	slot := o.slots.Get(name)
	if slot != nil {
		return slot
	}
	if s := o.scope.GetSlot(name); s != nil {
		return s.(rhino.StaticTypedSlot)
	}
	return nil
}

// LinkedFlowScope_OverlaySlot is an inferred slot of an OverlayScope.
type LinkedFlowScope_OverlaySlot struct {
	// TODO(sdh): add a final PMap<String, OverlaySlot> for properties
	name  string
	type_ rhino.JSType
}

// Equals implements jutil.Equaler using pointer identity (matching Java Object.equals).
func (s *LinkedFlowScope_OverlaySlot) Equals(other any) bool {
	that, ok := other.(*LinkedFlowScope_OverlaySlot)
	return ok && s == that
}

// Java: LinkedFlowScope.OverlaySlot(String, JSType)
func NewLinkedFlowScope_OverlaySlot(name string, type_ rhino.JSType) *LinkedFlowScope_OverlaySlot {
	return &LinkedFlowScope_OverlaySlot{
		name:  name,
		type_: type_,
	}
}

func newLinkedFlowScope_OverlaySlotWithArena(
	arena *linkedFlowScopeArena, name string, type_ rhino.JSType) *LinkedFlowScope_OverlaySlot {
	var s *LinkedFlowScope_OverlaySlot
	if arena != nil {
		s = arena.allocSlot()
	} else {
		s = &LinkedFlowScope_OverlaySlot{}
	}
	s.name = name
	s.type_ = type_
	return s
}

// Java: LinkedFlowScope.OverlaySlot.getName()
func (s *LinkedFlowScope_OverlaySlot) GetName() string {
	return s.name
}

// Java: LinkedFlowScope.OverlaySlot.getType()
func (s *LinkedFlowScope_OverlaySlot) GetType() rhino.JSType {
	return s.type_
}

// Java: LinkedFlowScope.OverlaySlot.isTypeInferred()
func (s *LinkedFlowScope_OverlaySlot) IsTypeInferred() bool {
	return true
}

// Java: LinkedFlowScope.OverlaySlot.getDeclaration()
func (s *LinkedFlowScope_OverlaySlot) GetDeclaration() rhino.StaticRef {
	return nil
}

// Java: LinkedFlowScope.OverlaySlot.getJSDocInfo()
func (s *LinkedFlowScope_OverlaySlot) GetJSDocInfo() *rhino.JSDocInfo {
	return nil
}

// Java: LinkedFlowScope.OverlaySlot.getScope()
func (s *LinkedFlowScope_OverlaySlot) GetScope() rhino.StaticScope {
	panic(jutil.NewUnsupportedOperationException())
}

var LinkedFlowScope_EMPTY_SCOPES *rhino.HamtPMap[*TypedScope, *LinkedFlowScope_OverlayScope] = rhino.HamtPMap_Empty[*TypedScope, *LinkedFlowScope_OverlayScope]()
var LinkedFlowScope_EMPTY_SLOTS *rhino.HamtPMap[string, *LinkedFlowScope_OverlaySlot] = rhino.HamtPMap_Empty[string, *LinkedFlowScope_OverlaySlot]()
