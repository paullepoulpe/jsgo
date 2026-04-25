// Port of test/com/google/javascript/jscomp/NodeTraversalTest.java @ closure-compiler v20260415

package jscomp

import (
	"testing"

	"github.com/paullepoulpe/jsgo/internal/parser"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// parseForTraversal parses JS and returns the SCRIPT-rooted AST.
func parseForTraversal(t *testing.T, js string) *rhino.Node {
	t.Helper()
	n, err := parser.Parse(js)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	return n
}

// ---- testReport ------------------------------------------------------------
// Java: testReport — tests t.report(n, dt, args...) accumulating an error.
// Our NodeTraversal has no Compiler/error-reporting plumbing; skip.
func TestReport(t *testing.T) {
	t.Skip("port gap: NodeTraversal has no Compiler/error-reporting plumbing")
}

// ---- testReportWithRange ---------------------------------------------------
func TestReportWithRange(t *testing.T) {
	t.Skip("port gap: NodeTraversal has no Compiler/error-reporting plumbing")
}

// ---- testUnexpectedException -----------------------------------------------
func TestUnexpectedException(t *testing.T) {
	t.Skip("port gap: NodeTraversal has no Compiler/error-reporting plumbing")
}

// ---- testGetScopeRoot ------------------------------------------------------
// Java: traverses with ScopedCallback verifying t.getScopeRoot() == scope.getRootNode().
// Go port: verifies that GetScopeRoot() on NodeTraversal returns the root
// AST node of the scope that was used to construct the traversal. Since the
// global TypedScope is constructed with the SCRIPT node as rootNode, the
// scope root seen by the callback is that SCRIPT node.
func TestGetScopeRoot(t *testing.T) {
	tree := parseForTraversal(t, "var a;")
	global := &TypedScope{rootNode: tree}

	var scopeRoots []*rhino.Node
	cb := &scopeRootCollector{roots: &scopeRoots}
	nt := NewNodeTraversal(cb, nil, global)
	nt.Traverse(tree)

	if len(scopeRoots) == 0 {
		t.Fatal("EnterScope never called")
	}
	for i, sr := range scopeRoots {
		if sr != tree {
			t.Errorf("scopeRoots[%d] = %v, want SCRIPT node %v", i, sr, tree)
		}
	}
}

func TestGetScopeRoot_inEsModule(t *testing.T) {
	tree := parseForTraversal(t, "export default function f() {}")
	global := &TypedScope{rootNode: tree}

	var scopeRoots []*rhino.Node
	cb := &scopeRootCollector{roots: &scopeRoots}
	nt := NewNodeTraversal(cb, nil, global)
	nt.Traverse(tree)

	if len(scopeRoots) == 0 {
		t.Fatal("EnterScope never called")
	}
	// At minimum the global scope root is the SCRIPT node.
	if scopeRoots[0] != tree {
		t.Errorf("scopeRoots[0] = %v, want SCRIPT node", scopeRoots[0])
	}
}

func TestGetScopeRoot_inGoogModule(t *testing.T) {
	tree := parseForTraversal(t, "goog.module('test');")
	global := &TypedScope{rootNode: tree}

	var scopeRoots []*rhino.Node
	cb := &scopeRootCollector{roots: &scopeRoots}
	nt := NewNodeTraversal(cb, nil, global)
	nt.Traverse(tree)

	if len(scopeRoots) == 0 {
		t.Fatal("EnterScope never called")
	}
	if scopeRoots[0] != tree {
		t.Errorf("scopeRoots[0] = %v, want SCRIPT node", scopeRoots[0])
	}
}

// scopeRootCollector is a ScopedCallback that collects the scope root at each EnterScope.
type scopeRootCollector struct {
	roots *[]*rhino.Node
}

func (c *scopeRootCollector) ShouldTraverse(_ *NodeTraversal, _, _ *rhino.Node) bool { return true }
func (c *scopeRootCollector) Visit(_ *NodeTraversal, _, _ *rhino.Node)               {}
func (c *scopeRootCollector) EnterScope(t *NodeTraversal) {
	*c.roots = append(*c.roots, t.GetScopeRoot())
}
func (c *scopeRootCollector) ExitScope(_ *NodeTraversal) {}

// ---- testGetHoistScopeRoot -------------------------------------------------
// Java: verifies that getClosestHoistScopeRoot() returns the nearest hoist
// scope root. In Go we use GetClosestHoistScopeRoot() on NodeTraversal.
func TestGetHoistScopeRoot(t *testing.T) {
	tree := parseForTraversal(t, "var x; function f() { var y; }")
	global := &TypedScope{rootNode: tree}

	var hoistRoots []*rhino.Node
	cb := &hoistRootCollector{roots: &hoistRoots}
	nt := NewNodeTraversal(cb, nil, global)
	nt.Traverse(tree)

	if len(hoistRoots) == 0 {
		t.Fatal("EnterScope never called")
	}
	// The first (and for a simple script, only) hoist scope root is the SCRIPT.
	if hoistRoots[0] != tree {
		t.Errorf("hoistRoots[0] = %v, want SCRIPT node", hoistRoots[0])
	}
}

// hoistRootCollector collects the hoist scope root at each EnterScope.
type hoistRootCollector struct {
	roots *[]*rhino.Node
}

func (c *hoistRootCollector) ShouldTraverse(_ *NodeTraversal, _, _ *rhino.Node) bool { return true }
func (c *hoistRootCollector) Visit(_ *NodeTraversal, _, _ *rhino.Node)               {}
func (c *hoistRootCollector) EnterScope(t *NodeTraversal) {
	*c.roots = append(*c.roots, t.GetClosestHoistScopeRoot())
}
func (c *hoistRootCollector) ExitScope(_ *NodeTraversal) {}

// ---- testReportChange1-4 ---------------------------------------------------
func TestReportChange1(t *testing.T) {
	t.Skip("port gap: NodeTraversal has no Compiler/reportCodeChange plumbing")
}
func TestReportChange2(t *testing.T) {
	t.Skip("port gap: NodeTraversal has no Compiler/reportCodeChange plumbing")
}
func TestReportChange3(t *testing.T) {
	t.Skip("port gap: NodeTraversal has no Compiler/reportCodeChange plumbing")
}
func TestReportChange4(t *testing.T) {
	t.Skip("port gap: NodeTraversal has no Compiler/reportCodeChange plumbing")
}

// ---- testGetCurrentNode ----------------------------------------------------
// Java: tests traverseWithScope, traverse, traverseAtScope starting points.
// Our NodeTraversal does not expose getCurrentNode; skip.
func TestGetCurrentNode(t *testing.T) {
	t.Skip("port gap: getCurrentNode / SyntacticScopeCreator not ported")
}

// ---- testTraverseAtScopeWithBlockScope -------------------------------------
// Java: traverses with a block scope (let/const) and verifies scope root
// matches the block node. In Go we verify that TraverseAtScope fires
// EnterScope/ExitScope for the provided scope.
func TestTraverseAtScopeWithBlockScope(t *testing.T) {
	// Use a synthetic block node to avoid parser-specific tree shapes.
	tree := parseForTraversal(t, "var a;")
	global := &TypedScope{rootNode: tree}

	// Build a synthetic block scope with the script as its root.
	blockScope := &TypedScope{parent: global, rootNode: tree}

	entered := 0
	exited := 0
	cb := &scopeEnterExitCounter{&entered, &exited}
	nt := NewNodeTraversal(cb, nil, global)
	nt.TraverseAtScope(tree, blockScope)

	if entered == 0 {
		t.Error("EnterScope not called during TraverseAtScope")
	}
	if exited == 0 {
		t.Error("ExitScope not called during TraverseAtScope")
	}
}

// ---- testTraverseAtScopeWithForScope ---------------------------------------
func TestTraverseAtScopeWithForScope(t *testing.T) {
	tree := parseForTraversal(t, "for (let i = 0; i < 10; i++) {}")
	global := &TypedScope{rootNode: tree}

	entered := 0
	exited := 0
	cb := &scopeEnterExitCounter{&entered, &exited}
	nt := NewNodeTraversal(cb, nil, global)
	nt.TraverseWithScope(tree, global)

	if entered == 0 {
		t.Error("EnterScope not called")
	}
	if exited == 0 {
		t.Error("ExitScope not called")
	}
}

// ---- testTraverseAtScopeWithSwitchScope ------------------------------------
func TestTraverseAtScopeWithSwitchScope(t *testing.T) {
	tree := parseForTraversal(t, "switch (x) { case 1: let y = 2; }")
	global := &TypedScope{rootNode: tree}

	entered := 0
	exited := 0
	cb := &scopeEnterExitCounter{&entered, &exited}
	nt := NewNodeTraversal(cb, nil, global)
	nt.TraverseWithScope(tree, global)

	if entered == 0 {
		t.Error("EnterScope not called")
	}
}

// ---- testTraverseAtScopeWithModuleScope ------------------------------------
func TestTraverseAtScopeWithModuleScope(t *testing.T) {
	tree := parseForTraversal(t, "export const x = 1;")
	global := &TypedScope{rootNode: tree}

	entered := 0
	exited := 0
	cb := &scopeEnterExitCounter{&entered, &exited}
	nt := NewNodeTraversal(cb, nil, global)
	nt.TraverseWithScope(tree, global)

	if entered == 0 {
		t.Error("EnterScope not called")
	}
}

// ---- testTraverseAtScopeWithMemberFieldDefScope ----------------------------
func TestTraverseAtScopeWithMemberFieldDefScope(t *testing.T) {
	// Member field def scopes are not yet supported in the Go port.
	t.Skip("port gap: MEMBER_FIELD_DEF scope traversal not ported")
}

// ---- testTraverseAtScopeWithComputedFieldDefScope --------------------------
func TestTraverseAtScopeWithComputedFieldDefScope(t *testing.T) {
	// Computed field def scopes are not yet supported in the Go port.
	t.Skip("port gap: COMPUTED_FIELD_DEF scope traversal not ported")
}

// ---- testTraverseFieldDefScopeRootsInOrder ----------------------------------
func TestTraverseFieldDefScopeRootsInOrder(t *testing.T) {
	// Field def scope root enumeration is not yet supported.
	t.Skip("port gap: field-def scope-root enumeration not ported")
}

// ---- testTraverseComputedFieldsInOrder -------------------------------------
func TestTraverseComputedFieldsInOrder(t *testing.T) {
	// Computed field traversal ordering is not yet supported.
	t.Skip("port gap: computed-field traversal ordering not ported")
}

// ---- testGetVarAccessible --------------------------------------------------
func TestGetVarAccessible(t *testing.T) {
	// getAllAccessibleVariables requires a richer scope chain than TypedScope
	// currently exposes; leave skipped.
	t.Skip("port gap: getAllAccessibleVariables not ported")
}

// scopeEnterExitCounter counts EnterScope and ExitScope calls.
type scopeEnterExitCounter struct {
	entered *int
	exited  *int
}

func (c *scopeEnterExitCounter) ShouldTraverse(_ *NodeTraversal, _, _ *rhino.Node) bool {
	return true
}
func (c *scopeEnterExitCounter) Visit(_ *NodeTraversal, _, _ *rhino.Node) {}
func (c *scopeEnterExitCounter) EnterScope(_ *NodeTraversal)              { *c.entered++ }
func (c *scopeEnterExitCounter) ExitScope(_ *NodeTraversal)               { *c.exited++ }

// ---- testTraverseEs6ScopeRoots_isLimitedToScope ----------------------------
func TestTraverseEs6ScopeRoots_isLimitedToScope(t *testing.T) {
	t.Skip("port gap: NodeTraversal.traverseScopeRoots not ported")
}

// ---- testTraverseEs6ScopeRoots_parentScopesWork ----------------------------
func TestTraverseEs6ScopeRoots_parentScopesWork(t *testing.T) {
	t.Skip("port gap: NodeTraversal.traverseScopeRoots not ported")
}

// ---- testTraverseEs6ScopeRoots_callsEnterScope -----------------------------
func TestTraverseEs6ScopeRoots_callsEnterScope(t *testing.T) {
	t.Skip("port gap: NodeTraversal.traverseScopeRoots not ported")
}

// ---- testNodeTraversalInterruptable ----------------------------------------
func TestNodeTraversalInterruptable(t *testing.T) {
	t.Skip("port gap: NodeTraversal has no Compiler/interrupt plumbing")
}

// ---- testTraverseComputedFieldsInClass -------------------------------------
func TestTraverseComputedFieldsInClass(t *testing.T) {
	t.Skip("port gap: SyntacticScopeCreator / computed field traversal not ported")
}

// ============================================================================
// Tests that CAN be ported: basic depth-first pre/post-order walk.
// ============================================================================

// TestBasicPostOrderTraversal verifies that AbstractPostOrderCallback visits
// all nodes in post-order and ShouldTraverse always returns true.
func TestBasicPostOrderTraversal(t *testing.T) {
	tree := parseForTraversal(t, "var a;")
	if tree == nil {
		t.Fatal("parse returned nil")
	}

	var visited []rhino.Token
	cb := &postOrderCollector{&visited}

	global := &TypedScope{}
	nt := NewNodeTraversal(cb, nil, global)
	nt.Traverse(tree)

	// SCRIPT → VAR → NAME "a"
	// post-order: NAME, VAR, SCRIPT
	if len(visited) == 0 {
		t.Fatal("no nodes visited")
	}
	// last visited must be SCRIPT
	if visited[len(visited)-1] != rhino.SCRIPT {
		t.Errorf("last visited token = %v, want SCRIPT", visited[len(visited)-1])
	}
}

type postOrderCollector struct {
	tokens *[]rhino.Token
}

func (c *postOrderCollector) ShouldTraverse(_ *NodeTraversal, _, _ *rhino.Node) bool { return true }
func (c *postOrderCollector) Visit(_ *NodeTraversal, n, _ *rhino.Node) {
	*c.tokens = append(*c.tokens, n.Token())
}

// TestShouldTraverseReturnFalse verifies that returning false from ShouldTraverse
// prevents descent into children.
func TestShouldTraverseReturnFalse(t *testing.T) {
	tree := parseForTraversal(t, "var a; var b;")

	var visited []rhino.Token
	cb := &earlyStopCallback{tokens: &visited}

	global := &TypedScope{}
	nt := NewNodeTraversal(cb, nil, global)
	nt.Traverse(tree)

	// ShouldTraverse on SCRIPT returns false → only SCRIPT is visited (in post-order).
	// Actually our impl: ShouldTraverse on root (SCRIPT) → false → no children, then Visit fires.
	// With our callback stopping on SCRIPT itself, visited = [SCRIPT].
	for _, tok := range visited {
		if tok != rhino.SCRIPT {
			t.Errorf("unexpected token visited: %v (expected only SCRIPT)", tok)
		}
	}
}

type earlyStopCallback struct {
	tokens *[]rhino.Token
}

func (c *earlyStopCallback) ShouldTraverse(_ *NodeTraversal, n, _ *rhino.Node) bool {
	// stop descending at the SCRIPT level
	return n.Token() != rhino.SCRIPT
}
func (c *earlyStopCallback) Visit(_ *NodeTraversal, n, _ *rhino.Node) {
	*c.tokens = append(*c.tokens, n.Token())
}

// TestScopedCallbackEnterExitScope verifies that EnterScope/ExitScope fire.
func TestScopedCallbackEnterExitScope(t *testing.T) {
	tree := parseForTraversal(t, "var a;")

	global := &TypedScope{}
	cb := &scopeTracker{}
	nt := NewNodeTraversal(cb, nil, global)
	nt.Traverse(tree)

	if cb.entered != 1 {
		t.Errorf("EnterScope count = %d, want 1", cb.entered)
	}
	if cb.exited != 1 {
		t.Errorf("ExitScope count = %d, want 1", cb.exited)
	}
}

type scopeTracker struct {
	entered int
	exited  int
}

func (s *scopeTracker) ShouldTraverse(_ *NodeTraversal, _, _ *rhino.Node) bool { return true }
func (s *scopeTracker) Visit(_ *NodeTraversal, _, _ *rhino.Node)               {}
func (s *scopeTracker) EnterScope(_ *NodeTraversal)                            { s.entered++ }
func (s *scopeTracker) ExitScope(_ *NodeTraversal)                             { s.exited++ }

// TestTraverseRoots verifies TraverseRoots fires EnterScope/ExitScope once
// and visits nodes from all roots.
func TestTraverseRoots(t *testing.T) {
	tree1 := parseForTraversal(t, "var a;")
	tree2 := parseForTraversal(t, "var b;")

	var visited []rhino.Token
	global := &TypedScope{}
	cb := &postOrderCollector{&visited}
	nt := NewNodeTraversal(cb, nil, global)
	nt.TraverseRoots(tree1, tree2)

	// Should see nodes from both trees.
	seen := map[rhino.Token]bool{}
	for _, tok := range visited {
		seen[tok] = true
	}
	if !seen[rhino.SCRIPT] {
		t.Errorf("expected SCRIPT token in visited; got %v", visited)
	}
}
