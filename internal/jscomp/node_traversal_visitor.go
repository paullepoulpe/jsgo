package jscomp

import (
	"github.com/paullepoulpe/jsgo/internal/ast"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// Go-only: traversal with a typed visitor (internal/ast).

// NodeTraversal_TraverseVisitor traverses root like NodeTraversal_Traverse
// with a post-order callback: after the children of a node, it calls v's
// handler for the node's kind, if v has one. The traversal skips the nodes of
// the kinds v has no handler for, so a visitor costs no more than a callback
// that switches on the token.
func NodeTraversal_TraverseVisitor(compiler AbstractCompiler, root rhino.Node, v *ast.Visitor[*NodeTraversal]) {
	NodeTraversal_Traverse(compiler, root, &nodeTraversal_visitorCallback{v: v, kinds: v.Tokens()})
}

// NodeTraversal_TraverseRootsVisitor traverses externs and root like
// NodeTraversal_TraverseRoots with a post-order ast.Visitor.
func NodeTraversal_TraverseRootsVisitor(compiler AbstractCompiler, externs rhino.Node, root rhino.Node, v *ast.Visitor[*NodeTraversal]) {
	NodeTraversal_TraverseRoots(compiler, &nodeTraversal_visitorCallback{v: v, kinds: v.Tokens()}, externs, root)
}

// nodeTraversal_visitorCallback adapts an ast.Visitor to the callback
// interface.
type nodeTraversal_visitorCallback struct {
	v     *ast.Visitor[*NodeTraversal]
	kinds ast.TokenSet // the tokens v has a handler for
}

func (c *nodeTraversal_visitorCallback) visitTokens() *ast.TokenSet { return &c.kinds }

func (c *nodeTraversal_visitorCallback) ShouldTraverse(*NodeTraversal, rhino.Node, rhino.Node) bool {
	return true
}

func (c *nodeTraversal_visitorCallback) Visit(t *NodeTraversal, n rhino.Node, _ rhino.Node) {
	c.v.Visit(t, n)
}

func (c *nodeTraversal_visitorCallback) nodeTraversalNoShouldTraverse() {}
