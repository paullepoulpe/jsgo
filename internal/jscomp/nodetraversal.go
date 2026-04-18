// Port of com/google/javascript/jscomp/NodeTraversal.java
// @ closure-compiler v20260415 (subset).
//
// Java's NodeTraversal is ~1400 lines with: Callback + ScopedCallback
// interfaces, AbstractPreOrderCallback / AbstractPostOrderCallback /
// AbstractCfgCallback convenience bases, a Builder, per-scope /
// per-root entry points (traverse, traverseRoots, traverseScopeRoots,
// traverseFunctionOutOfBand), and rich scope-stack management. This
// file carries the Callback contract + the single-root depth-first
// driver TypeCheck uses today. Scope tracking, pre-order hooks, and
// Builder-driven configuration arrive with later callers.

package jscomp

import "github.com/paullepoulpe/jsgo/internal/rhino"

// Callback mirrors the Java NodeTraversal.Callback interface:
// ShouldTraverse runs in pre-order and decides whether to descend;
// Visit runs in post-order after children.
type Callback interface {
	ShouldTraverse(t *NodeTraversal, n, parent *rhino.Node) bool
	Visit(t *NodeTraversal, n, parent *rhino.Node)
}

// NodeTraversal drives the walk. Java's class tracks the current
// Compiler, scope stack, source file, input id, change-time flags,
// etc.; we carry the scope stack (TypedScope) because TypeCheck needs
// to route GETPROP / CALL / ASSIGN lookups through the nearest
// enclosing function or class scope.
type NodeTraversal struct {
	callback Callback
	creator  *TypedScopeCreator
	global   *TypedScope
	scopes   []*TypedScope
}

// NewNodeTraversal builds a traversal driven by the given callback,
// starting from the provided global scope. The scope creator is
// consulted when we descend into a FUNCTION (or class method) whose
// body needs its local scope active.
func NewNodeTraversal(cb Callback, creator *TypedScopeCreator, global *TypedScope) *NodeTraversal {
	return &NodeTraversal{callback: cb, creator: creator, global: global, scopes: []*TypedScope{global}}
}

// Traverse walks `root` depth-first, invoking ShouldTraverse before
// descent and Visit after children. Mirrors the external shape of
// Java's NodeTraversal.traverse(Compiler, Node, Callback).
func (t *NodeTraversal) Traverse(root *rhino.Node) {
	t.traverse(root, nil)
}

func (t *NodeTraversal) traverse(n, parent *rhino.Node) {
	if n == nil {
		return
	}
	if !t.callback.ShouldTraverse(t, n, parent) {
		return
	}
	// Enter a function body's local scope, if we have one. Mirrors
	// ScopedCallback.enterScope at the Java level, specialised to
	// TypeCheck's needs.
	pushed := false
	if t.creator != nil {
		if local, ok := t.creator.Functions[n]; ok && local != nil {
			t.scopes = append(t.scopes, local)
			pushed = true
		}
	}
	for _, c := range n.Children() {
		t.traverse(c, n)
	}
	if pushed {
		t.scopes = t.scopes[:len(t.scopes)-1]
	}
	t.callback.Visit(t, n, parent)
}

// CurrentScope returns the innermost active TypedScope (set of
// declarations reachable via name lookup at the current traversal
// position). Matches Java's NodeTraversal.getScope signature.
func (t *NodeTraversal) CurrentScope() *TypedScope {
	if len(t.scopes) == 0 {
		return t.global
	}
	return t.scopes[len(t.scopes)-1]
}
