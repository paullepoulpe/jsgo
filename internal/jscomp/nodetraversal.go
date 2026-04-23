// Port of com/google/javascript/jscomp/NodeTraversal.java
// @ closure-compiler v20260415 (widened subset).
//
// Ported shapes:
//   - Callback interface (ShouldTraverse + Visit)
//   - ScopedCallback interface (Callback + EnterScope + ExitScope)
//   - AbstractPostOrderCallback (always traverses; subclass supplies Visit)
//   - AbstractPreOrderCallback  (always traverses; no-op Visit)
//   - AbstractScopedCallback    (ScopedCallback; always traverses; no-op Enter/Exit)
//   - NodeTraversal struct with scope-aware enter/exit hooks
//
// Not yet ported (no caller): AbstractCfgCallback, AbstractShallowCallback,
// AbstractShallowStatementCallback, AbstractModuleCallback,
// AbstractChangedScopeCallback, Builder, traverseRoots, traverseAtScope,
// traverseScopeRoots, traverseWithScope, GetScopeDepth, GetScopeRoot,
// GetClosestHoistScopeRoot, GetClosestContainerScope, GetClosestHoistScope.

package jscomp

import "github.com/paullepoulpe/jsgo/internal/rhino"

// Callback mirrors the Java NodeTraversal.Callback interface:
// ShouldTraverse runs in pre-order and decides whether to descend;
// Visit runs in post-order after children.
type Callback interface {
	ShouldTraverse(t *NodeTraversal, n, parent *rhino.Node) bool
	Visit(t *NodeTraversal, n, parent *rhino.Node)
}

// ScopedCallback mirrors the Java NodeTraversal.ScopedCallback
// interface: a Callback that is also notified when the traversal
// enters or exits a scope (function body, script root, etc.).
// EnterScope fires immediately after the new scope is pushed;
// ExitScope fires immediately before the scope is popped.
type ScopedCallback interface {
	Callback
	EnterScope(t *NodeTraversal)
	ExitScope(t *NodeTraversal)
}

// AbstractPostOrderCallback mirrors Java's
// NodeTraversal.AbstractPostOrderCallback: ShouldTraverse always
// returns true so every node is visited in post-order.  Concrete
// implementations embed this type and supply only Visit.
type AbstractPostOrderCallback struct{}

func (AbstractPostOrderCallback) ShouldTraverse(_ *NodeTraversal, _, _ *rhino.Node) bool {
	return true
}

// AbstractPreOrderCallback mirrors Java's
// NodeTraversal.AbstractPreOrderCallback: ShouldTraverse always
// returns true and Visit is a no-op.  Concrete implementations embed
// this type and supply only ShouldTraverse.
type AbstractPreOrderCallback struct{}

func (AbstractPreOrderCallback) ShouldTraverse(_ *NodeTraversal, _, _ *rhino.Node) bool {
	return true
}

func (AbstractPreOrderCallback) Visit(_ *NodeTraversal, _, _ *rhino.Node) {}

// AbstractScopedCallback mirrors Java's
// NodeTraversal.AbstractScopedCallback: a ScopedCallback that always
// traverses (ShouldTraverse returns true) and provides no-op
// EnterScope / ExitScope default implementations.  Concrete types
// embed this and override only the methods they care about.
type AbstractScopedCallback struct{}

func (AbstractScopedCallback) ShouldTraverse(_ *NodeTraversal, _, _ *rhino.Node) bool {
	return true
}

func (AbstractScopedCallback) EnterScope(_ *NodeTraversal) {}
func (AbstractScopedCallback) ExitScope(_ *NodeTraversal)  {}

// NodeTraversal drives the walk. Java's class tracks the current
// Compiler, scope stack, source file, input id, change-time flags,
// etc.; we carry the scope stack (TypedScope) because TypeCheck needs
// to route GETPROP / CALL / ASSIGN lookups through the nearest
// enclosing function or class scope.
//
// The scopeCallback field is set when the supplied Callback also
// implements ScopedCallback; EnterScope / ExitScope are then fired
// around every scope push / pop, matching Java's pushScope /
// popScope pattern.
type NodeTraversal struct {
	callback      Callback
	scopeCallback ScopedCallback // non-nil iff callback implements ScopedCallback
	creator       *TypedScopeCreator
	global        *TypedScope
	scopes        []*TypedScope
}

// NewNodeTraversal builds a traversal driven by the given callback,
// starting from the provided global scope.  The scope creator is
// consulted when we descend into a FUNCTION (or class method) whose
// body needs its local scope active.
//
// If cb also satisfies ScopedCallback its EnterScope / ExitScope
// methods are called whenever a scope is pushed or popped.
func NewNodeTraversal(cb Callback, creator *TypedScopeCreator, global *TypedScope) *NodeTraversal {
	t := &NodeTraversal{
		callback: cb,
		creator:  creator,
		global:   global,
		scopes:   []*TypedScope{global},
	}
	if sc, ok := cb.(ScopedCallback); ok {
		t.scopeCallback = sc
	}
	return t
}

// Traverse walks `root` depth-first, invoking ShouldTraverse before
// descent and Visit after children.  Mirrors the external shape of
// Java's NodeTraversal.traverse(Compiler, Node, Callback).
//
// The global scope is treated as the outermost scope entry: if a
// ScopedCallback is present, EnterScope fires for the global scope
// before the walk begins and ExitScope fires after.
func (t *NodeTraversal) Traverse(root *rhino.Node) {
	// Fire EnterScope for the implicit global scope that was pushed
	// in NewNodeTraversal, matching Java's traverse() which calls
	// pushScope(root) before traverseBranch.
	if t.scopeCallback != nil {
		t.scopeCallback.EnterScope(t)
	}
	t.traverse(root, nil)
	if t.scopeCallback != nil {
		t.scopeCallback.ExitScope(t)
	}
}

// TraverseInto exposes the recursive walker so callbacks can drive
// custom child orderings — e.g. TypeCheck visits an IF's condition,
// then descends into the then-branch under one refinement frame and
// the else-branch under another.  Mirrors NodeTraversal#traverse on
// the Java side.
//
// TraverseInto does not fire EnterScope / ExitScope for the node
// itself; those fire only at scope-boundary nodes (FUNCTION bodies)
// encountered during the walk.
func (t *NodeTraversal) TraverseInto(n, parent *rhino.Node) {
	t.traverse(n, parent)
}

// pushScope activates `s` as the innermost scope and fires
// EnterScope on any ScopedCallback.  Mirrors Java's pushScope(Node).
func (t *NodeTraversal) pushScope(s *TypedScope) {
	t.scopes = append(t.scopes, s)
	if t.scopeCallback != nil {
		t.scopeCallback.EnterScope(t)
	}
}

// popScope removes the innermost scope and fires ExitScope on any
// ScopedCallback.  Mirrors Java's popScope().
func (t *NodeTraversal) popScope() {
	if t.scopeCallback != nil {
		t.scopeCallback.ExitScope(t)
	}
	t.scopes = t.scopes[:len(t.scopes)-1]
}

func (t *NodeTraversal) traverse(n, parent *rhino.Node) {
	if n == nil {
		return
	}
	if !t.callback.ShouldTraverse(t, n, parent) {
		return
	}
	// Enter a function body's local scope, if we have one.  Mirrors
	// Java's handleFunction → pushScope pattern: the scope is pushed
	// before children are traversed, EnterScope fires, children run,
	// ExitScope fires, then scope is popped.
	pushed := false
	if t.creator != nil {
		if local, ok := t.creator.Functions[n]; ok && local != nil {
			t.pushScope(local)
			pushed = true
		}
	}
	for _, c := range n.Children() {
		t.traverse(c, n)
	}
	if pushed {
		t.popScope()
	}
	t.callback.Visit(t, n, parent)
}

// CurrentScope returns the innermost active TypedScope (set of
// declarations reachable via name lookup at the current traversal
// position).  Matches Java's NodeTraversal.getScope signature.
func (t *NodeTraversal) CurrentScope() *TypedScope {
	if len(t.scopes) == 0 {
		return t.global
	}
	return t.scopes[len(t.scopes)-1]
}
