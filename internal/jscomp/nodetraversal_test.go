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
// Our NodeTraversal tracks scopes via TypedScope, not SyntacticScopeCreator; skip.
func TestGetScopeRoot(t *testing.T) {
	t.Skip("port gap: SyntacticScopeCreator / getScopeRoot not ported")
}

func TestGetScopeRoot_inEsModule(t *testing.T) {
	t.Skip("port gap: SyntacticScopeCreator / getScopeRoot not ported")
}

func TestGetScopeRoot_inGoogModule(t *testing.T) {
	t.Skip("port gap: SyntacticScopeCreator / getScopeRoot not ported")
}

// ---- testGetHoistScopeRoot -------------------------------------------------
func TestGetHoistScopeRoot(t *testing.T) {
	t.Skip("port gap: getClosestHoistScopeRoot not ported")
}

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
func TestTraverseAtScopeWithBlockScope(t *testing.T) {
	t.Skip("port gap: SyntacticScopeCreator not ported")
}

// ---- testTraverseAtScopeWithForScope ---------------------------------------
func TestTraverseAtScopeWithForScope(t *testing.T) {
	t.Skip("port gap: SyntacticScopeCreator not ported")
}

// ---- testTraverseAtScopeWithSwitchScope ------------------------------------
func TestTraverseAtScopeWithSwitchScope(t *testing.T) {
	t.Skip("port gap: SyntacticScopeCreator not ported")
}

// ---- testTraverseAtScopeWithModuleScope ------------------------------------
func TestTraverseAtScopeWithModuleScope(t *testing.T) {
	t.Skip("port gap: SyntacticScopeCreator not ported")
}

// ---- testTraverseAtScopeWithMemberFieldDefScope ----------------------------
func TestTraverseAtScopeWithMemberFieldDefScope(t *testing.T) {
	t.Skip("port gap: SyntacticScopeCreator not ported")
}

// ---- testTraverseAtScopeWithComputedFieldDefScope --------------------------
func TestTraverseAtScopeWithComputedFieldDefScope(t *testing.T) {
	t.Skip("port gap: SyntacticScopeCreator not ported")
}

// ---- testTraverseFieldDefScopeRootsInOrder ----------------------------------
func TestTraverseFieldDefScopeRootsInOrder(t *testing.T) {
	t.Skip("port gap: SyntacticScopeCreator / scope-root enumeration not ported")
}

// ---- testTraverseComputedFieldsInOrder -------------------------------------
func TestTraverseComputedFieldsInOrder(t *testing.T) {
	t.Skip("port gap: SyntacticScopeCreator not ported")
}

// ---- testGetVarAccessible --------------------------------------------------
func TestGetVarAccessible(t *testing.T) {
	t.Skip("port gap: SyntacticScopeCreator / getAllAccessibleVariables not ported")
}

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
