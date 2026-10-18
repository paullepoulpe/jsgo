package rhino

import "github.com/paullepoulpe/jsgo/internal/jutil"

// Node_CloneTreeInSlab clones a tree into a caller-owned allocation group.
// This is for immutable long-lived parse templates: sharing the ordinary
// allocation chunk with a short-lived compilation can retain its whole AST.
// Type expressions in JSDoc are copied separately by the template owner.
//
// Go-only: performance — allocation ownership only; the detached clone has
// the same fields as CloneTreeWithTypeExprs(false).
func Node_CloneTreeInSlab(n *Node, slab *jutil.Slab[Node]) *Node {
	clone := slab.New()
	*clone = *n
	clone.parent, clone.previous, clone.next, clone.first = nil, nil, nil, nil
	if n.ext != nil {
		ext := *n.ext
		clone.ext = &ext
	}
	switch n.kind {
	case node_kindNumber:
		clone.token = Token_NUMBER
	case node_kindBigInt:
		clone.token = Token_BIGINT
	case node_kindTemplateLiteralSubstring:
		clone.token = Token_TEMPLATELIT_STRING
	}
	for child := n.first; child != nil; child = child.next {
		clone.AddChildToBack(Node_CloneTreeInSlab(child, slab))
	}
	if shadow := n.GetClosureUnawareShadow(); shadow != nil {
		clone.SetClosureUnawareShadow(Node_CloneTreeInSlab(shadow, slab))
	}
	return clone
}
