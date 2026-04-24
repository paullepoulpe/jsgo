// Port of com/google/javascript/jscomp/NodeTraversal.java
// @ closure-compiler v20260415 (widened subset).
//
// Ported shapes:
//   - Callback interface (ShouldTraverse + Visit)
//   - ScopedCallback interface (Callback + EnterScope + ExitScope)
//   - ScopeCreator interface (createScope 2-pass protocol)
//   - AbstractPostOrderCallback (always traverses; subclass supplies Visit)
//   - AbstractPreOrderCallback  (always traverses; no-op Visit)
//   - AbstractScopedCallback    (ScopedCallback; always traverses; no-op Enter/Exit)
//   - NodeTraversal struct with scope-aware enter/exit hooks
//   - TraverseWithScope — traverse under a pre-built scope
//   - TraverseRoots     — multi-root batch entry
//   - TraverseAtScope   — scope-rooted traversal with ancestor scopes pushed quietly
//   - TraverseFunctionOutOfBand — traverse a function body in isolation
//
// Not yet ported (no caller): AbstractCfgCallback, AbstractShallowCallback,
// AbstractShallowStatementCallback, AbstractModuleCallback,
// AbstractChangedScopeCallback, Builder, traverseScopeRoots,
// GetScopeDepth, GetScopeRoot, GetClosestHoistScopeRoot,
// GetClosestContainerScope, GetClosestHoistScope.

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

func (AbstractPreOrderCallback) Visit(t *NodeTraversal, n, parent *rhino.Node) {
	// Default no-op: AbstractPreOrderCallback processes nodes in ShouldTraverse only.
	_ = t
}

// AbstractScopedCallback mirrors Java's
// NodeTraversal.AbstractScopedCallback: a ScopedCallback that always
// traverses (ShouldTraverse returns true) and provides no-op
// EnterScope / ExitScope default implementations.  Concrete types
// embed this and override only the methods they care about.
type AbstractScopedCallback struct{}

func (AbstractScopedCallback) ShouldTraverse(_ *NodeTraversal, _, _ *rhino.Node) bool {
	return true
}

func (AbstractScopedCallback) EnterScope(t *NodeTraversal) {
	// Default no-op: AbstractScopedCallback subclasses override as needed.
	_ = t
}
func (AbstractScopedCallback) ExitScope(t *NodeTraversal) {
	// Default no-op: AbstractScopedCallback subclasses override as needed.
	_ = t
}

// ScopeCreator mirrors the Java ScopeCreator interface.  It is called
// in pre-order (when a scope boundary node is first entered) to
// create the scope for that node; the scope is then active for all
// descendants until the traversal exits the boundary, at which point
// ExitScope fires on any ScopedCallback (post-order).  This 2-pass
// protocol matches Java's lazy getAbstractScope / popScope sequence.
//
// CreateScope receives the scope-boundary node and its parent scope
// and returns the new child scope.  ExitScope is called with the same
// node immediately before the scope is popped — implementors that
// need to do post-processing (e.g. finalise inferred return types)
// may hook here.  A nil ExitScope implementation is valid; the
// interface method is still called, so a no-op pointer receiver is
// fine.
type ScopeCreatorIface interface {
	// CreateScope is called in pre-order when the traversal pushes a
	// new scope boundary (function body, block-scope node, etc.).
	// parent is the enclosing TypedScope (never nil).
	CreateScopeForNode(n *rhino.Node, parent *TypedScope) *TypedScope

	// ExitScopeForNode is called in post-order, just before the scope
	// pushed by CreateScopeForNode is removed from the stack.
	ExitScopeForNode(n *rhino.Node, scope *TypedScope)
}

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

// Go-only: language-adapter — Java uses NodeTraversal.Builder; Go uses a
// constructor function since there is no builder pattern here.
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

// Go-only: language-adapter — Java's traverseBranch is private; Go exposes
// it as TraverseInto for callers that need to drive child traversal order.
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

// TraverseWithScope walks `root` under the provided pre-built `scope`,
// replacing the current innermost scope for the duration of the walk.
// EnterScope fires for `scope` before the walk; ExitScope fires after.
// This mirrors Java's private NodeTraversal.traverseWithScope(Node,
// AbstractScope), which pushes the given scope and traverses `root`
// with parent==nil (i.e. treating `root` as the traversal root).
//
// Use this when you hold a pre-built global or module scope and want
// to traverse a subtree under it without rebuilding the scope from
// scratch.
func (t *NodeTraversal) TraverseWithScope(root *rhino.Node, scope *TypedScope) {
	t.pushScope(scope)
	t.traverse(root, nil)
	t.popScope()
}

// TraverseRoots walks each root in `roots` under the shared global
// scope held by the traversal.  EnterScope fires once before the
// first root and ExitScope fires once after the last root — matching
// Java's private NodeTraversal.traverseRoots which calls
// pushScope(scopeRoot) once and then traverseBranch for each root in
// sequence.
//
// The public static Java overload is
// NodeTraversal.traverseRoots(compiler, cb, externs, root); this Go
// method is the equivalent instance-level entry point.
func (t *NodeTraversal) TraverseRoots(roots ...*rhino.Node) {
	if len(roots) == 0 {
		return
	}
	// The global scope is already the bottom of t.scopes; fire
	// EnterScope once for it (mirrors Java's single pushScope call).
	if t.scopeCallback != nil {
		t.scopeCallback.EnterScope(t)
	}
	for _, root := range roots {
		t.traverse(root, nil)
	}
	if t.scopeCallback != nil {
		t.scopeCallback.ExitScope(t)
	}
}

// TraverseAtScope begins a traversal rooted at `scopeRoot`, with all
// ancestor scopes in `scope`'s parent chain pushed silently (no
// EnterScope for ancestors) and EnterScope / ExitScope firing only
// for `scope` itself.
//
// This mirrors Java's private NodeTraversal.traverseAtScope(AbstractScope)
// which calls s.getRootNode() to obtain the root, then walks
// s.getParent() to collect ancestors, pushes them with quietly=true,
// and finally pushes the target scope normally before traversing the
// root node.
//
// Because Go's TypedScope does not carry its root AST node (unlike
// Java's AbstractScope), the caller must supply `scopeRoot` explicitly.
func (t *NodeTraversal) TraverseAtScope(scopeRoot *rhino.Node, scope *TypedScope) {
	// Collect ancestor scopes so we can push outermost-first.
	var ancestors []*TypedScope
	for p := scope.parent; p != nil; p = p.parent {
		ancestors = append(ancestors, p)
	}
	// Push ancestors quietly (no EnterScope callback), outermost-first.
	for i := len(ancestors) - 1; i >= 0; i-- {
		t.scopes = append(t.scopes, ancestors[i])
	}
	// Push the target scope with normal EnterScope notification.
	t.pushScope(scope)
	// Traverse from the scope's root node.
	t.traverse(scopeRoot, nil)
	// Pop target scope with normal ExitScope notification.
	t.popScope()
	// Pop the silently-pushed ancestor scopes (no ExitScope callbacks).
	t.scopes = t.scopes[:len(t.scopes)-len(ancestors)]
}

// GetScopeRoot returns the root AST node of the current innermost scope.
// Mirrors Java's NodeTraversal.getScopeRoot() which returns
// t.getScope().getRootNode(). In Go the TypedScope carries its root node
// through the AbstractScopeBase embedded struct.
func (t *NodeTraversal) GetScopeRoot() *rhino.Node {
	return t.CurrentScope().GetRootNode()
}

// GetClosestHoistScopeRoot returns the root AST node of the nearest hoist
// scope at or above the current scope. Mirrors
// NodeTraversal.getClosestHoistScopeRoot() in Java.
func (t *NodeTraversal) GetClosestHoistScopeRoot() *rhino.Node {
	for i := len(t.scopes) - 1; i >= 0; i-- {
		s := t.scopes[i]
		if s.IsHoistScope() {
			return s.GetRootNode()
		}
	}
	return t.global.GetRootNode()
}

// Go-only: language-adapter — Java's traverseFunction is private; Go exposes
// it for TypeCheck's out-of-band function re-traversal.
func (t *NodeTraversal) TraverseFunctionOutOfBand(root *rhino.Node, scope *TypedScope) {
	if !t.callback.ShouldTraverse(t, root, nil) {
		return
	}
	t.pushScope(scope)
	for _, c := range root.Children() {
		t.traverse(c, root)
	}
	t.popScope()
	t.callback.Visit(t, root, nil)
}
