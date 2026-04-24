// Port of com/google/javascript/rhino/Node.java
// @ closure-compiler v20260415 (subset).
//
// Retains the Java API shape callers from other passes rely on —
// parent/sibling links, Detach/DetachChildren, newString/newNumber
// factories, a small Prop property bag. Child storage uses the same
// intrusive circular doubly-linked list as Java: Node.first points to
// the first child, Node.last to the last; siblings link via next/previous.

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
	// PropQuotedStringKey mirrors Java's Node.isQuotedStringKey() — set by
	// the parser on STRING_KEY nodes whose source key was a string literal
	// (e.g. `{"x": 1}`) rather than a bare identifier (e.g. `{x: 1}`).
	PropQuotedStringKey
	// PropGetterDef / PropSetterDef mirror Java's Node.isGetterDef() /
	// Node.isSetterDef(). Set by the parser on FUNCTION nodes that are
	// ES6 class getter / setter methods (`get foo(){}` / `set foo(v){}`).
	PropGetterDef
	PropSetterDef
	// PropSideEffectFlags stores an inverted bitmask of SideEffectFlags bits.
	// Mirrors Java's Node.Prop.SIDE_EFFECT_FLAGS. Only meaningful on CALL/NEW.
	PropSideEffectFlags
	// PropIsNamespace mirrors Java's Node.IS_NAMESPACE — set on NAME nodes
	// that introduce a goog.provide / goog.module namespace.
	PropIsNamespace
	// PropFreeCall mirrors Java's Node.FREE_CALL — set on CALL nodes that
	// have no explicit receiver ("free" calls, e.g. `f()` not `obj.f()`).
	PropFreeCall
	// PropDirectEval mirrors Java's Node.DIRECT_EVAL — set on direct eval
	// CALL nodes (i.e. the callee is literally `eval`).
	PropDirectEval
	// PropStartOfOptChain mirrors Java's Node.START_OF_OPT_CHAIN — set on
	// OPTCHAIN_CALL / OPTCHAIN_GETPROP / OPTCHAIN_GETELEM nodes that begin
	// an optional chain (i.e. immediately follow the `?.` operator).
	PropStartOfOptChain
	// PropSynthetic mirrors Java's Node.SYNTHETIC — set on synthesized BLOCK
	// nodes that do not correspond to source-level braces.
	PropSynthetic
	// PropIncrDecr mirrors Java's Node.INCRDECR_PROP — set on INC/DEC nodes
	// when the operator is postfix (e.g. x++ vs ++x).
	PropIncrDecr
	// PropShorthandProperty mirrors Java's Node.isShorthandProperty() — set on
	// STRING_KEY nodes in object literals that use shorthand syntax (e.g. {x}).
	PropShorthandProperty
	// PropIsParenthesized mirrors Java's Node.isParenthesized() — set on
	// expression nodes that are wrapped in explicit parentheses in source.
	PropIsParenthesized
	// PropIsPrivateIdentifier mirrors Java's Node.isPrivateIdentifier() — set on
	// GETPROP / NAME nodes that reference a private class property (#foo).
	PropIsPrivateIdentifier
	// PropTemplateLitRaw stores the raw string value on TEMPLATELIT_STRING nodes.
	// Mirrors the second string argument to Node.newTemplateLitString(cooked, raw).
	PropTemplateLitRaw
)

type Node struct {
	token Token
	str   string
	num   float64

	// Intrusive circular doubly-linked child list, mirroring Java's
	// Node.first / Node.last. Siblings are linked by next/previous.
	first *Node // first child (nil when no children)
	last  *Node // last child  (nil when no children)

	parent   *Node
	next     *Node
	previous *Node

	jsdoc  *JSDocInfo
	lineno int // 1-indexed, matches Node.getLineno()
	charno int // 0-indexed, matches Node.getCharno()

	// jstypeOrColor mirrors Java's Node.jstypeOrColor: the compiler-inferred
	// JSType (or Color in the colour-only pipeline). Stored as `any` to avoid
	// an import cycle (jstype imports rhino). Callers use GetJSType/SetJSType.
	jstypeOrColor any

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

// AddChildToBack appends c as the last child of n. Mirrors Java's
// Node.addChildToBack(Node). O(1).
func (n *Node) AddChildToBack(c *Node) {
	if c == nil {
		return
	}
	c.parent = n
	c.next = nil
	c.previous = n.last
	if n.last != nil {
		n.last.next = c
	} else {
		n.first = c
	}
	n.last = c
}

// AddChildToFront prepends c as the first child of n. Mirrors Java's
// Node.addChildToFront(Node). O(1).
func (n *Node) AddChildToFront(c *Node) {
	if c == nil {
		return
	}
	c.parent = n
	c.previous = nil
	c.next = n.first
	if n.first != nil {
		n.first.previous = c
	} else {
		n.last = c
	}
	n.first = c
}

// FirstChild mirrors Java's getFirstChild. O(1).
func (n *Node) FirstChild() *Node { return n.first }

// LastChild mirrors Java's getLastChild. O(1).
func (n *Node) LastChild() *Node { return n.last }

// Children builds and returns a slice of children by walking the linked
// list. The returned slice is a snapshot — callers may range over it
// freely. For hot paths prefer iterating with FirstChild/Next directly.
func (n *Node) Children() []*Node {
	if n.first == nil {
		return nil
	}
	var out []*Node
	for c := n.first; c != nil; c = c.next {
		out = append(out, c)
	}
	return out
}

// Detach removes n from its parent. Mirrors Java's Node.detach(). O(1).
// No-op when already detached.
func (n *Node) Detach() {
	if n.parent == nil {
		return
	}
	p := n.parent
	if n.previous != nil {
		n.previous.next = n.next
	} else {
		p.first = n.next
	}
	if n.next != nil {
		n.next.previous = n.previous
	} else {
		p.last = n.previous
	}
	n.parent, n.next, n.previous = nil, nil, nil
}

// IsFirstChildOf mirrors Node.isFirstChildOf — true iff n is the first
// child of parent. Matches Java's Node.isFirstChildOf.
func (n *Node) IsFirstChildOf(parent *Node) bool {
	if parent == nil {
		return false
	}
	return parent.FirstChild() == n
}

// IsQualifiedName mirrors Node.isQualifiedName() — true for NAME (non-
// empty), THIS, and GETPROP whose receiver is itself a qualified name.
func (n *Node) IsQualifiedName() bool {
	switch n.token {
	case NAME:
		return n.str != ""
	case THIS:
		return true
	case GETPROP:
		if fc := n.FirstChild(); fc != nil {
			return fc.IsQualifiedName()
		}
		return false
	}
	return false
}

// GetQualifiedName mirrors Node.getQualifiedName() — returns the dotted
// GetQualifiedName is provided by qualified_name.go (the canonical port).

// DetachChildren mirrors Java's Node.detachChildren — remove every child
// and null out its links.
func (n *Node) DetachChildren() {
	for c := n.first; c != nil; {
		next := c.next
		c.parent, c.next, c.previous = nil, nil, nil
		c = next
	}
	n.first, n.last = nil, nil
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
	return n.first != nil
}

// HasParent mirrors Java's Node.hasParent().
func (n *Node) HasParent() bool { return n.parent != nil }

// IsStringLit mirrors Java's Node.isStringLit() — true for STRING_LIT
// token (used in type-expression ASTs).
func (n *Node) IsStringLit() bool { return n.token == STRINGLIT }

// IsStringKey mirrors Java's Node.isStringKey() — true for STRING_KEY
// token (used in record-type ASTs).
func (n *Node) IsStringKey() bool { return n.token == STRING_KEY }

// IsQuotedStringKey mirrors Java's Node.isQuotedStringKey() — true when
// a STRING_KEY node was parsed from a quoted string literal (e.g. {"x": 1})
// rather than a bare identifier ({x: 1}). Set by the parser via
// PutBooleanProp(PropQuotedStringKey, true).
func (n *Node) IsQuotedStringKey() bool {
	return n.token == STRING_KEY && n.GetBooleanProp(PropQuotedStringKey)
}

// IsBlock mirrors Java's Node.isBlock().
func (n *Node) IsBlock() bool { return n.token == BLOCK }

// IsGetterDef mirrors Java's Node.isGetterDef() — true when this FUNCTION
// node is an ES6 class getter method (`get foo() {}`).
func (n *Node) IsGetterDef() bool { return n.GetBooleanProp(PropGetterDef) }

// IsSetterDef mirrors Java's Node.isSetterDef() — true when this FUNCTION
// node is an ES6 class setter method (`set foo(v) {}`).
func (n *Node) IsSetterDef() bool { return n.GetBooleanProp(PropSetterDef) }

// ReplaceWith mirrors Java's Node.replaceWith(Node). Swaps n out of its
// parent and inserts replacement in the same position. O(1). No-op when
// n has no parent (already detached).
func (n *Node) ReplaceWith(replacement *Node) {
	if n.parent == nil {
		return
	}
	p := n.parent
	replacement.parent = p
	replacement.next = n.next
	replacement.previous = n.previous
	if n.previous != nil {
		n.previous.next = replacement
	} else {
		p.first = replacement
	}
	if n.next != nil {
		n.next.previous = replacement
	} else {
		p.last = replacement
	}
	n.parent, n.next, n.previous = nil, nil, nil
}

// RemoveChildren detaches all children from n and returns them as a slice.
// Mirrors Java's Node.removeChildren() (which returns the first-child node;
// we return a slice so AddChildrenToBack can re-attach them).
func (n *Node) RemoveChildren() []*Node {
	var kids []*Node
	for c := n.first; c != nil; {
		next := c.next
		kids = append(kids, c)
		c.parent, c.next, c.previous = nil, nil, nil
		c = next
	}
	n.first, n.last = nil, nil
	return kids
}

// AddChildrenToBack appends a slice of children (result of RemoveChildren)
// to n. Mirrors Java's Node.addChildrenToBack(Node firstChild).
func (n *Node) AddChildrenToBack(kids []*Node) {
	for _, c := range kids {
		n.AddChildToBack(c)
	}
}

// CloneTree performs a deep copy of `n` and its entire subtree. Mirrors
// Java's Node.cloneTree(). Position and string/number values are preserved;
// jsdoc and props are NOT deep-copied (shallow reference, matching Java's
// cloneTree which does copy props but not JSDocInfo).
func (n *Node) CloneTree() *Node {
	if n == nil {
		return nil
	}
	clone := &Node{
		token:  n.token,
		str:    n.str,
		num:    n.num,
		lineno: n.lineno,
		charno: n.charno,
		jsdoc:  n.jsdoc,
	}
	if n.props != nil {
		clone.props = make(map[Prop]any, len(n.props))
		for k, v := range n.props {
			clone.props[k] = v
		}
	}
	for c := n.first; c != nil; c = c.next {
		clone.AddChildToBack(c.CloneTree())
	}
	return clone
}

// IsEquivalentTo performs a structural equality check on two subtrees.
// Mirrors Java's Node.isEquivalentTo(Node) — compares token, string/number
// values, and children recursively.
func (n *Node) IsEquivalentTo(other *Node) bool {
	if n == nil && other == nil {
		return true
	}
	if n == nil || other == nil {
		return false
	}
	if n.token != other.token || n.str != other.str || n.num != other.num {
		return false
	}
	nc := n.first
	oc := other.first
	for nc != nil && oc != nil {
		if !nc.IsEquivalentTo(oc) {
			return false
		}
		nc = nc.next
		oc = oc.next
	}
	// Both must be exhausted simultaneously.
	return nc == nil && oc == nil
}

// ---------------- property bag ------------------------------------- //

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

// ---------------- int-prop bag ---------------------------------------- //

// GetIntProp ports Node.getIntProp(Prop). Returns 0 when absent.
func (n *Node) GetIntProp(p Prop) int {
	if n.props == nil {
		return 0
	}
	v, ok := n.props[p]
	if !ok {
		return 0
	}
	i, _ := v.(int)
	return i
}

// PutIntProp ports Node.putIntProp(Prop, int).
func (n *Node) PutIntProp(p Prop, val int) {
	if val == 0 {
		if n.props != nil {
			delete(n.props, p)
		}
		return
	}
	if n.props == nil {
		n.props = map[Prop]any{}
	}
	n.props[p] = val
}

// ---------------- JSType storage -------------------------------------- //

// GetJSType ports Node.getJSType(). Returns the compiler-inferred type stored
// on this node, or nil when none has been set. The concrete value is a
// *jstype.JSType (or subtype thereof) — callers in the jscomp layer assert the
// concrete type at use-sites via a type-switch/cast.
func (n *Node) GetJSType() any { return n.jstypeOrColor }

// SetJSType ports Node.setJSType(JSType). Stores the compiler-inferred type.
// Passing nil clears the annotation.
func (n *Node) SetJSType(t any) { n.jstypeOrColor = t }

// ---------------- SideEffectFlags ------------------------------------- //

// SideEffectFlags mirrors the Java inner class Node.SideEffectFlags.
// It encodes which side effects a CALL or NEW node may produce.
// Bit layout (identical to Java):
//
//	bit 0 MUTATES_GLOBAL_STATE
//	bit 1 MUTATES_THIS
//	bit 2 MUTATES_ARGUMENTS
//	bit 3 THROWS
type SideEffectFlags struct{ value int }

const (
	SideEffectMutatesGlobalState = 1
	SideEffectMutatesThis        = 2
	SideEffectMutatesArguments   = 4
	SideEffectThrows             = 8

	sideEffectUsedBitsMask = (1 << 4) - 1 // 0b1111

	SideEffectNoSideEffects = 0
	SideEffectAllSideEffects = SideEffectMutatesGlobalState |
		SideEffectMutatesThis | SideEffectMutatesArguments | SideEffectThrows
)

// NewSideEffectFlags returns a SideEffectFlags with all bits set (Java default).
func NewSideEffectFlags() SideEffectFlags {
	return SideEffectFlags{value: SideEffectAllSideEffects}
}

// NewSideEffectFlagsValue wraps a raw int value.
func NewSideEffectFlagsValue(v int) SideEffectFlags { return SideEffectFlags{value: v} }

func (f *SideEffectFlags) ValueOf() int { return f.value }

// SetAllFlags sets all side-effect bits.
func (f *SideEffectFlags) SetAllFlags() *SideEffectFlags {
	f.value = SideEffectAllSideEffects
	return f
}

// ClearAllFlags clears all side-effect bits.
func (f *SideEffectFlags) ClearAllFlags() *SideEffectFlags {
	f.value = SideEffectNoSideEffects
	return f
}

// SetMutatesGlobalState sets the global-state, arguments, and this mutation bits
// (matching Java's semantics: mutating global state implies the others).
func (f *SideEffectFlags) SetMutatesGlobalState() *SideEffectFlags {
	f.value |= SideEffectMutatesGlobalState | SideEffectMutatesArguments | SideEffectMutatesThis
	return f
}

// SetThrows sets the throws bit.
func (f *SideEffectFlags) SetThrows() *SideEffectFlags {
	f.value |= SideEffectThrows
	return f
}

// SetMutatesThis sets the mutates-this bit.
func (f *SideEffectFlags) SetMutatesThis() *SideEffectFlags {
	f.value |= SideEffectMutatesThis
	return f
}

// SetMutatesArguments sets the mutates-arguments bit.
func (f *SideEffectFlags) SetMutatesArguments() *SideEffectFlags {
	f.value |= SideEffectMutatesArguments
	return f
}

// SetSideEffectFlags ports Node.setSideEffectFlags(int).
// Only meaningful on CALL / NEW / TAGGED_TEMPLATELIT nodes.
// The value is stored inverted so the default (0-int-prop) == all flags set.
func (n *Node) SetSideEffectFlags(flags int) {
	n.PutIntProp(PropSideEffectFlags, ^flags&sideEffectUsedBitsMask)
}

// SetSideEffectFlagsValue is a convenience overload accepting SideEffectFlags.
func (n *Node) SetSideEffectFlagsValue(flags SideEffectFlags) {
	n.SetSideEffectFlags(flags.value)
}

// GetSideEffectFlags ports Node.getSideEffectFlags().
// Returns the raw int bitmask; defaults to all-bits-set (ALL_SIDE_EFFECTS)
// when the prop has never been written.
func (n *Node) GetSideEffectFlags() int {
	return ^n.GetIntProp(PropSideEffectFlags) & sideEffectUsedBitsMask
}

// IsNoSideEffectsCall returns true iff the call has no side effects.
func (n *Node) IsNoSideEffectsCall() bool {
	return n.GetSideEffectFlags() == SideEffectNoSideEffects
}

// IsOnlyModifiesThisCall mirrors Java: returns true when the only side-effect
// (if any) is mutating the receiver.
func (n *Node) IsOnlyModifiesThisCall() bool {
	return n.GetSideEffectFlags()&^SideEffectMutatesThis == SideEffectNoSideEffects
}

// IsOnlyModifiesArgumentsCall mirrors Java.
func (n *Node) IsOnlyModifiesArgumentsCall() bool {
	return n.GetSideEffectFlags()&^SideEffectMutatesArguments == SideEffectNoSideEffects
}

// MayMutateArguments returns true iff the MUTATES_ARGUMENTS bit is set.
func (n *Node) MayMutateArguments() bool {
	return n.GetSideEffectFlags()&SideEffectMutatesArguments != 0
}

// MayMutateGlobalStateOrThrow returns true iff the global-state or throws bit is set.
func (n *Node) MayMutateGlobalStateOrThrow() bool {
	return n.GetSideEffectFlags()&(SideEffectMutatesGlobalState|SideEffectThrows) != 0
}

// ---------------- shallow clone --------------------------------------- //

// CloneNode ports Node.cloneNode() — returns a detached shallow copy of n
// (same token/str/num/pos/jsdoc/props) without any children.
// Distinct from CloneTree which recursively clones children.
func (n *Node) CloneNode() *Node {
	if n == nil {
		return nil
	}
	clone := &Node{
		token:         n.token,
		str:           n.str,
		num:           n.num,
		lineno:        n.lineno,
		charno:        n.charno,
		jsdoc:         n.jsdoc,
		jstypeOrColor: n.jstypeOrColor,
	}
	if n.props != nil {
		clone.props = make(map[Prop]any, len(n.props))
		for k, v := range n.props {
			clone.props[k] = v
		}
	}
	return clone
}

// ---------------- navigation helpers ---------------------------------- //

// GetGrandparent ports Node.getGrandparent(). Returns nil when the node has
// fewer than two ancestors.
func (n *Node) GetGrandparent() *Node {
	if n.parent == nil {
		return nil
	}
	return n.parent.parent
}

// GetSecondChild ports Node.getSecondChild(). Returns the second child, or nil
// when n has fewer than two children.
func (n *Node) GetSecondChild() *Node {
	if n.first == nil {
		return nil
	}
	return n.first.next
}

// GetLastSibling ports Node.getLastSibling(). Returns the last sibling of n
// (following the chain of Next pointers). Returns n itself when it has no
// next sibling.
func (n *Node) GetLastSibling() *Node {
	cur := n
	for cur.next != nil {
		cur = cur.next
	}
	return cur
}

// ---------------- source-ref helpers ---------------------------------- //

// Srcref ports Node.srcref(Node other) — copies lineno/charno from other.
// Returns n for fluent use (mirrors Java's @CanIgnoreReturnValue signature).
func (n *Node) Srcref(other *Node) *Node {
	if other == nil {
		return n
	}
	n.lineno = other.lineno
	n.charno = other.charno
	return n
}

// SrcrefTree ports Node.srcrefTree(Node other) — copies source position from
// other onto n and every node in its subtree.
func (n *Node) SrcrefTree(other *Node) *Node {
	n.Srcref(other)
	for c := n.first; c != nil; c = c.next {
		c.SrcrefTree(other)
	}
	return n
}

// SrcrefIfMissing ports Node.srcrefIfMissing(Node other) — copies lineno/charno
// only when the node's current position is unset (lineno == 0).
func (n *Node) SrcrefIfMissing(other *Node) *Node {
	if other == nil {
		return n
	}
	if n.lineno == 0 {
		n.lineno = other.lineno
		n.charno = other.charno
	}
	return n
}

// SrcrefTreeIfMissing ports Node.srcrefTreeIfMissing(Node other).
func (n *Node) SrcrefTreeIfMissing(other *Node) *Node {
	n.SrcrefIfMissing(other)
	for c := n.first; c != nil; c = c.next {
		c.SrcrefTreeIfMissing(other)
	}
	return n
}

// ---------------- token classifiers ----------------------------------- //

// MatchesName ports Node.matchesName(String). Returns true iff this is a
// non-empty NAME node whose string value equals name.
func (n *Node) MatchesName(name string) bool {
	if n.token != NAME {
		return false
	}
	return n.str != "" && n.str == name
}

// MatchesQualifiedNameStr ports Node.matchesQualifiedName(String). Returns true
// iff this node represents the dotted name encoded in qname (e.g. "a.b.c").
func (n *Node) MatchesQualifiedNameStr(qname string) bool {
	return n.matchesQualifiedNameAt(qname, len(qname))
}

// matchesQualifiedNameAt is the recursive helper for MatchesQualifiedNameStr.
// It mirrors Java's private Node.matchesQualifiedName(String, int).
func (n *Node) matchesQualifiedNameAt(qname string, end int) bool {
	start := lastIndexByte(qname, '.', end) + 1
	switch n.token {
	case NAME:
		name := n.str
		return start == 0 && name != "" && len(name) == end && qname[:end] == name
	case THIS:
		return start == 0 && end == 4 && qname[:4] == "this"
	case SUPER:
		return start == 0 && end == 5 && qname[:5] == "super"
	case GETPROP:
		prop := n.str
		return start > 1 &&
			len(prop) == end-start &&
			qname[start:end] == prop &&
			n.FirstChild() != nil &&
			n.FirstChild().matchesQualifiedNameAt(qname, start-1)
	default:
		return false
	}
}

// lastIndexByte returns the index of the last occurrence of b in s[:end],
// or -1 if not found. Equivalent to strings.LastIndex restricted to a byte.
func lastIndexByte(s string, b byte, end int) int {
	for i := end - 1; i >= 0; i-- {
		if s[i] == b {
			return i
		}
	}
	return -1
}

// IsOptChainCall mirrors Java's Node.isOptChainCall().
func (n *Node) IsOptChainCall() bool { return n.token == OPTCHAIN_CALL }

// IsOptChainGetElem mirrors Java's Node.isOptChainGetElem().
func (n *Node) IsOptChainGetElem() bool { return n.token == OPTCHAIN_GETELEM }

// IsOptChainGetProp mirrors Java's Node.isOptChainGetProp().
func (n *Node) IsOptChainGetProp() bool { return n.token == OPTCHAIN_GETPROP }

// IsOptionalChainStart mirrors Java's Node.isOptionalChainStart() — true when
// this OPTCHAIN node is the start of the chain (immediately after `?.`).
func (n *Node) IsOptionalChainStart() bool {
	return n.GetBooleanProp(PropStartOfOptChain)
}

// SetIsOptionalChainStart mirrors Java's Node.setIsOptionalChainStart(bool).
func (n *Node) SetIsOptionalChainStart(b bool) {
	n.PutBooleanProp(PropStartOfOptChain, b)
}

// IsFreeCall mirrors Java's Node.getBooleanProp(FREE_CALL). True when this
// CALL node has no explicit receiver (i.e. callee is not a GETPROP/GETELEM).
func (n *Node) IsFreeCall() bool { return n.GetBooleanProp(PropFreeCall) }

// SetFreeCall mirrors Java's Node.putBooleanProp(FREE_CALL, b).
func (n *Node) SetFreeCall(b bool) { n.PutBooleanProp(PropFreeCall, b) }

// IsNamespace mirrors Java's Node.getBooleanProp(IS_NAMESPACE).
func (n *Node) IsNamespace() bool { return n.GetBooleanProp(PropIsNamespace) }

// SetNamespace mirrors Java's Node.putBooleanProp(IS_NAMESPACE, b).
func (n *Node) SetNamespace(b bool) { n.PutBooleanProp(PropIsNamespace, b) }

// IsGetProp mirrors Java's Node.isGetprop() — true for GETPROP token.
func (n *Node) IsGetProp() bool { return n.token == GETPROP }

// IsGetElem mirrors Java's Node.isGetelem() — true for GETELEM token.
func (n *Node) IsGetElem() bool { return n.token == GETELEM }

// IsCall mirrors Java's Node.isCall() — true for CALL token.
func (n *Node) IsCall() bool { return n.token == CALL }

// IsNew mirrors Java's Node.isNew() — true for NEW token.
func (n *Node) IsNew() bool { return n.token == NEW }

// IsName mirrors Java's Node.isName() — true for NAME token.
func (n *Node) IsName() bool { return n.token == NAME }

// IsFunction mirrors Java's Node.isFunction() — true for FUNCTION token.
func (n *Node) IsFunction() bool { return n.token == FUNCTION }

// IsScript mirrors Java's Node.isScript() — true for SCRIPT token.
func (n *Node) IsScript() bool { return n.token == SCRIPT }

// IsReturn mirrors Java's Node.isReturn() — true for RETURN token.
func (n *Node) IsReturn() bool { return n.token == RETURN }

// IsSpread mirrors Java's Node.isSpread() — true for ITER_SPREAD token.
func (n *Node) IsSpread() bool { return n.token == ITER_SPREAD }

// IsComputedProp mirrors Java's Node.isComputedProp() — true for COMPUTED_PROP token.
func (n *Node) IsComputedProp() bool { return n.token == COMPUTED_PROP }

// IsIn mirrors Java's Node.isIn() — true for the IN token.
func (n *Node) IsIn() bool { return n.token == IN }

// IsPrivateIdentifier mirrors Java's Node.isPrivateIdentifier() — true when
// PropIsPrivateIdentifier is set (e.g. #foo property access).
func (n *Node) IsPrivateIdentifier() bool { return n.GetBooleanProp(PropIsPrivateIdentifier) }

// IsShorthandProperty mirrors Java's Node.isShorthandProperty() — true when
// PropShorthandProperty is set on a STRING_KEY in an object literal.
func (n *Node) IsShorthandProperty() bool { return n.GetBooleanProp(PropShorthandProperty) }

// SetShorthandProperty mirrors Java's Node.setShorthandProperty(bool).
func (n *Node) SetShorthandProperty(b bool) { n.PutBooleanProp(PropShorthandProperty, b) }

// IsParenthesized mirrors Java's Node.isParenthesized() — true when the
// expression was wrapped in explicit parentheses in source.
func (n *Node) IsParenthesized() bool { return n.GetBooleanProp(PropIsParenthesized) }

// SetIsParenthesized mirrors Java's Node.setIsParenthesized(bool).
func (n *Node) SetIsParenthesized(b bool) { n.PutBooleanProp(PropIsParenthesized, b) }

// SetStaticMember mirrors Java's Node.setStaticMember(bool) — marks
// a method/accessor as a static class member.
func (n *Node) SetStaticMember(b bool) { n.PutBooleanProp(PropStaticMember, b) }

// IsStaticMember mirrors Java's Node.isStaticMember().
func (n *Node) IsStaticMember() bool { return n.GetBooleanProp(PropStaticMember) }

// NewTemplateLitString creates a TEMPLATELIT_STRING node with cooked and raw
// values. Mirrors Node.newTemplateLitString(String, String) in Java.
func NewTemplateLitString(cooked, raw string) *Node {
	n := &Node{token: TEMPLATELIT_STRING, str: cooked}
	// Store the raw value as a second string via the props map.
	if n.props == nil {
		n.props = make(map[Prop]any)
	}
	n.props[PropTemplateLitRaw] = raw
	return n
}

// GetRawString returns the raw template literal string stored on a
// TEMPLATELIT_STRING node. Mirrors Node.getRawString().
func (n *Node) GetRawString() string {
	if v, ok := n.props[PropTemplateLitRaw]; ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

