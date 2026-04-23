// Port of com/google/javascript/rhino/Node.java
// @ closure-compiler v20260415 (subset).
//
// Retains the Java API shape callers from other passes rely on —
// parent/sibling links, Detach/DetachChildren, newString/newNumber
// factories, a small Prop property bag — while keeping the child
// storage as a plain slice for readability. The Java side uses an
// intrusive singly-linked list (Node.first + child.next), but the
// slice + explicit parent/next/previous maintenance below produces
// the same observable API. Bits still not ported: cloneNode,
// SideEffectFlags, the full ~50-entry Prop enum, JSTypeNative
// annotation slot.

package rhino

// Prop mirrors Node.Prop in Java. Only entries the checker currently
// consults are listed; new props extend this enum as callers need them.
type Prop int

const (
	PropJSDocInfo Prop = iota
	PropStaticMember
	PropIsConstantName
	PropArrowFn
	// PropAddedBlock mirrors Node.ADDED_BLOCK (the IS_ADDED_BLOCK prop
	// bit in Java). Set by the parser when a BLOCK was synthesized
	// implicitly from a single statement body (i.e. without explicit
	// braces). CheckSuspiciousCode uses it to detect `if (x);`.
	PropAddedBlock
	PropDeclaredConstantVar // set by InferConsts for @const / const-keyword vars
	PropInferredConstantVar // set by InferConsts for single-assignment vars
)

type Node struct {
	token    Token
	str      string
	num      float64
	children []*Node // slice view; Java uses Node.first linked-list

	parent   *Node
	next     *Node
	previous *Node

	jsdoc  *JSDocInfo
	lineno int // 1-indexed, matches Node.getLineno()
	charno int // 0-indexed, matches Node.getCharno()

	props map[Prop]any
}

// ---------------- factories ---------------------------------------- //

// New is the token-only constructor, matching `new Node(Token)` in Java.
func New(tok Token) *Node { return &Node{token: tok} }

// NewString ports `Node.newString(Token, String)`.
func NewString(tok Token, s string) *Node { return &Node{token: tok, str: s} }

// NewNumber ports `Node.newNumber(double)`.
func NewNumber(v float64) *Node { return &Node{token: NUMBER, num: v} }

// ---------------- basic accessors (Java getters) ------------------- //

func (n *Node) Token() Token             { return n.token }
func (n *Node) String_() string          { return n.str }
func (n *Node) Number() float64          { return n.num }
func (n *Node) JSDocInfo() *JSDocInfo    { return n.jsdoc }
func (n *Node) SetJSDocInfo(i *JSDocInfo) { n.jsdoc = i }
func (n *Node) Lineno() int              { return n.lineno }
func (n *Node) Charno() int              { return n.charno }
func (n *Node) SetPos(line, col int)     { n.lineno, n.charno = line, col }

// Parent mirrors Java's getParent(). Nil for detached / root nodes.
func (n *Node) Parent() *Node { return n.parent }

// Next / Previous mirror getNext / getPrevious.
func (n *Node) Next() *Node     { return n.next }
func (n *Node) Previous() *Node { return n.previous }

// ---------------- tree building ------------------------------------ //

// AddChildToBack appends `c` and keeps parent/next/previous in sync.
// Matches the observable behaviour of Java's addChildToBack.
func (n *Node) AddChildToBack(c *Node) {
	if c == nil {
		return
	}
	c.parent = n
	if len(n.children) > 0 {
		last := n.children[len(n.children)-1]
		last.next = c
		c.previous = last
	} else {
		c.previous = nil
	}
	c.next = nil
	n.children = append(n.children, c)
}

// AddChildToFront ports Node.addChildToFront.
func (n *Node) AddChildToFront(c *Node) {
	if c == nil {
		return
	}
	c.parent = n
	c.previous = nil
	if len(n.children) > 0 {
		first := n.children[0]
		first.previous = c
		c.next = first
	} else {
		c.next = nil
	}
	n.children = append([]*Node{c}, n.children...)
}

// FirstChild mirrors Java's getFirstChild.
func (n *Node) FirstChild() *Node {
	if len(n.children) == 0 {
		return nil
	}
	return n.children[0]
}

// Children exposes the child slice — not in Java (which iterates
// first.next), but jsgo callers prefer range loops. Treat the returned
// slice as read-only.
func (n *Node) Children() []*Node { return n.children }

// Detach removes `n` from its parent. Mirrors Node.detach (Java returns
// `this`; Go doesn't need the fluent API so we return nothing). No-op
// when already detached.
func (n *Node) Detach() {
	if n.parent == nil {
		return
	}
	p := n.parent
	for i, c := range p.children {
		if c == n {
			p.children = append(p.children[:i], p.children[i+1:]...)
			break
		}
	}
	if n.previous != nil {
		n.previous.next = n.next
	}
	if n.next != nil {
		n.next.previous = n.previous
	}
	n.parent, n.next, n.previous = nil, nil, nil
}

// DetachChildren mirrors Node.detachChildren — remove every child and
// null out its links.
func (n *Node) DetachChildren() {
	for _, c := range n.children {
		c.parent, c.next, c.previous = nil, nil, nil
	}
	n.children = nil
}

// ---------------- property bag ------------------------------------- //

// IsAddedBlock mirrors Java's Node.isAddedBlock() — true when this
// BLOCK was synthesized implicitly (no explicit braces). The parser
// sets PropAddedBlock on BLOCK nodes it creates for single-statement
// control-flow bodies. CheckSuspiciousCode uses this to detect
// a misplaced semicolon like `if (x);`.
func (n *Node) IsAddedBlock() bool {
	return n.GetBooleanProp(PropAddedBlock)
}

// HasChildren mirrors Java's Node.hasChildren(). Returns true when
// the node has at least one child.
func (n *Node) HasChildren() bool {
	return len(n.children) > 0
}

// GetBooleanProp ports Node.getBooleanProp. Absent props are false.
func (n *Node) GetBooleanProp(p Prop) bool {
	if n.props == nil {
		return false
	}
	v, ok := n.props[p]
	if !ok {
		return false
	}
	b, _ := v.(bool)
	return b
}

// PutBooleanProp ports Node.putBooleanProp.
func (n *Node) PutBooleanProp(p Prop, b bool) {
	if !b {
		delete(n.props, p)
		return
	}
	if n.props == nil {
		n.props = map[Prop]any{}
	}
	n.props[p] = true
}

// IsDeclaredConstantVar ports Node.isDeclaredConstantVar(). Returns true
// when the name node was annotated @const or declared with the const keyword.
func (n *Node) IsDeclaredConstantVar() bool {
	return n.GetBooleanProp(PropDeclaredConstantVar)
}

// SetDeclaredConstantVar ports Node.setDeclaredConstantVar(boolean).
func (n *Node) SetDeclaredConstantVar(b bool) {
	n.PutBooleanProp(PropDeclaredConstantVar, b)
}

// IsInferredConstantVar ports Node.isInferredConstantVar(). Returns true
// when InferConsts determined the variable is assigned exactly once.
func (n *Node) IsInferredConstantVar() bool {
	return n.GetBooleanProp(PropInferredConstantVar)
}

// SetInferredConstantVar ports Node.setInferredConstantVar(boolean).
func (n *Node) SetInferredConstantVar(b bool) {
	n.PutBooleanProp(PropInferredConstantVar, b)
}

// HasChildren mirrors Node.hasChildren() — true when the node has at
// least one child. Used by InferConsts to detect initialised let bindings.
func (n *Node) HasChildren() bool {
	return len(n.children) > 0
}
