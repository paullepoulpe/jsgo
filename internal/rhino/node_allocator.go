package rhino

import "math/big"

// NodeAllocator owns one independent parse's allocation stream. Its zero value
// is ready to use. It must not be copied after first use. It must have a single
// owner; separate allocators can be used
// concurrently. Nodes do not retain the allocator, and chunks are never pooled
// or published globally, so discarding a speculative parse cannot couple its
// lifetime to another parse's live AST.
//
// Go-only: performance — explicit ownership for speculative native parsing.
type NodeAllocator struct {
	current []Node
	used    int
}

func (a *NodeAllocator) NewNode(token Token) *Node {
	if a.used == len(a.current) {
		size := len(a.current) * 2
		if size == 0 {
			size = 8
		} else if size > 128 {
			size = 128
		}
		a.current = make([]Node, size)
		a.used = 0
	}
	n := &a.current[a.used]
	a.used++
	n.token = token
	n.linenoCharno = -1
	return n
}

func (a *NodeAllocator) NewString(token Token, value string) *Node {
	n := a.NewNode(token)
	n.kind = node_kindString
	n.SetString(value)
	return n
}

func (a *NodeAllocator) NewNumber(value float64) *Node {
	n := a.NewNode(Token_NUMBER)
	n.kind = node_kindNumber
	n.SetDouble(value)
	return n
}

func (a *NodeAllocator) NewBigInt(value *big.Int) *Node {
	n := a.NewNode(Token_BIGINT)
	n.kind = node_kindBigInt
	n.SetBigInt(value)
	return n
}

func (a *NodeAllocator) NewTemplateLiteral(cooked *string, raw string) *Node {
	n := a.NewNode(Token_TEMPLATELIT_STRING)
	n.kind = node_kindTemplateLiteralSubstring
	if cooked != nil {
		value := RhinoStringPool_AddOrGet(*cooked)
		n.node_getExt().cooked = &value
	}
	n.str = rhinoStringPoolCell(raw)
	return n
}

// CloneNode has the same detached, shallow semantics as Node.CloneNode,
// including folded subclass payloads and shared immutable properties.
func (a *NodeAllocator) CloneNode(source *Node) *Node {
	var clone *Node
	switch source.kind {
	case node_kindNumber:
		clone = a.NewNumber(source.node_number())
	case node_kindBigInt:
		clone = a.NewBigInt(source.GetBigInt())
	case node_kindString:
		clone = a.NewNode(source.token)
		clone.kind = node_kindString
		clone.str = source.str
	case node_kindTemplateLiteralSubstring:
		clone = a.NewTemplateLiteral(source.node_cooked(), source.nodeStringValue())
	default:
		clone = a.NewNode(source.token)
	}
	Node_CopyBaseNodeFields(source, clone, false)
	return clone
}
