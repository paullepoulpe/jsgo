package jscomp

import (
	"github.com/paullepoulpe/jsgo/internal/ast"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// Go-only: traversal with a typed visitor (internal/ast).

// NodeTraversal_TraverseVisitor traverses root like NodeTraversal_Traverse
// with a post-order callback: after the children of a node, it calls v's
// handler for the node's kind, if v has one.
func NodeTraversal_TraverseVisitor(compiler AbstractCompiler, root rhino.Node, v *ast.Visitor[*NodeTraversal]) {
	NodeTraversal_Traverse(compiler, root, &nodeTraversal_visitorCallback{v})
}

// nodeTraversal_visitorCallback adapts an ast.Visitor to the callback
// interface.
type nodeTraversal_visitorCallback struct {
	v *ast.Visitor[*NodeTraversal]
}

func (c *nodeTraversal_visitorCallback) ShouldTraverse(*NodeTraversal, rhino.Node, rhino.Node) bool {
	return true
}

func (c *nodeTraversal_visitorCallback) Visit(t *NodeTraversal, n rhino.Node, _ rhino.Node) {
	c.v.Visit(t, n)
}

func (c *nodeTraversal_visitorCallback) nodeTraversalNoShouldTraverse() {}
