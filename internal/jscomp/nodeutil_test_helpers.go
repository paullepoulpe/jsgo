// Port of test/com/google/javascript/jscomp/NodeUtilTest.java helper methods
// @ closure-compiler v20260415

package jscomp

import (
	"testing"

	"github.com/paullepoulpe/jsgo/internal/parser"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// nuParse parses a JS string and returns the SCRIPT-rooted AST.
// Mirrors NodeUtilTest.parse(String).
func nuParse(t *testing.T, js string) *rhino.Node {
	t.Helper()
	n, err := parser.Parse(js)
	if err != nil {
		t.Fatalf("nuParse(%q): %v", js, err)
	}
	return n
}

// nuParseExpr parses an expression by wrapping it in parens, then returns
// the inner expression node. Mirrors NodeUtilTest.parseExpr(String).
func nuParseExpr(t *testing.T, js string) *rhino.Node {
	t.Helper()
	script := nuParse(t, "("+js+");")
	// script → EXPR_RESULT → expr
	exprResult := script.FirstChild()
	if exprResult == nil {
		t.Fatalf("nuParseExpr(%q): no first child of script", js)
	}
	expr := exprResult.FirstChild()
	if expr == nil {
		t.Fatalf("nuParseExpr(%q): no first child of EXPR_RESULT", js)
	}
	return expr
}

// nuGetNode performs DFS over root searching for a descendant with token.
// Mirrors NodeUtilTest.getNode(Node, Token).
func nuGetNode(t *testing.T, root *rhino.Node, token rhino.Token) *rhino.Node {
	t.Helper()
	n := nuGetNodeOrNull(root, token)
	if n == nil {
		t.Fatalf("nuGetNode: no %v node found in tree", token)
	}
	return n
}

// nuGetNodeOrNull performs DFS, returning nil if not found.
func nuGetNodeOrNull(root *rhino.Node, token rhino.Token) *rhino.Node {
	for c := root.FirstChild(); c != nil; c = c.Next() {
		if c.Token() == token {
			return c
		}
		if found := nuGetNodeOrNull(c, token); found != nil {
			return found
		}
	}
	return nil
}

// nuParseFirst parses JS and returns the first node with the given token in DFS order.
// Mirrors NodeUtilTest.parseFirst(Token, String).
func nuParseFirst(t *testing.T, token rhino.Token, js string) *rhino.Node {
	t.Helper()
	root := nuParse(t, js)
	if token == rhino.SCRIPT {
		return root
	}
	return nuGetNode(t, root, token)
}

// nuParseSecond parses JS and returns the second node with the given token.
// Mirrors NodeUtilTest.parseSecond(Token, String).
func nuParseSecond(t *testing.T, token rhino.Token, js string) *rhino.Node {
	t.Helper()
	first := nuParseFirst(t, token, js)
	return nuGetNode(t, first, token)
}
