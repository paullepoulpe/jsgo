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
	// PropGeneratorFn mirrors Java's Node.isGeneratorFunction() — set on FUNCTION
	// nodes that are generator functions (i.e. function*).
	PropGeneratorFn
	// PropAsyncFn mirrors Java's Node.isAsyncFunction() — set on FUNCTION nodes
	// that are async functions.
	PropAsyncFn
	// PropArrowFnBody is not used; PropArrowFn already exists above.
	// PropExportAll mirrors Java's Node.EXPORT_ALL_FROM — set on EXPORT nodes that
	// are `export * from 'mod'`.
	PropExportAll
	// PropExportDefault mirrors Java's Node.EXPORT_DEFAULT — set on EXPORT nodes
	// that are `export default`.
	PropExportDefault
	// PropYieldAll mirrors Java's Node.YIELD_ALL — set on YIELD nodes that are
	// yield* (i.e. delegate generators).
	PropYieldAll
	// PropOptEs6Typed mirrors Java's Node.OPT_ES6_TYPED — set on FUNCTION
	// parameter nodes that are TypeScript optional parameters.
	PropOptEs6Typed
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
	length int // source length in characters, matches Node.getLength()

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
func (n *Node) SetString_(s string)      { n.str = s }
func (n *Node) SetString(s string)       { n.str = s }
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

// InsertAfter inserts n immediately after existing in the child list of
// existing's parent. Mirrors Java's Node.insertAfter(Node existing).
// Panics if existing has no parent.
func (n *Node) InsertAfter(existing *Node) {
	p := existing.parent
	n.parent = p
	n.previous = existing
	n.next = existing.next
	if existing.next != nil {
		existing.next.previous = n
	} else {
		p.last = n
	}
	existing.next = n
}

// InsertBefore inserts n immediately before existing in the child list of
// existing's parent. Mirrors Java's Node.insertBefore(Node existing).
// Panics if existing has no parent.
func (n *Node) InsertBefore(existing *Node) {
	p := existing.parent
	n.parent = p
	n.next = existing
	n.previous = existing.previous
	if existing.previous != nil {
		existing.previous.next = n
	} else {
		p.first = n
	}
	existing.previous = n
}

// HasXChildren mirrors Java's Node.hasXChildren(int). Returns true when n has
// exactly count children.
func (n *Node) HasXChildren(count int) bool {
	c := n.first
	for i := 0; i < count; i++ {
		if c == nil {
			return false
		}
		c = c.next
	}
	return c == nil
}

// HasOneChild mirrors Java's Node.hasOneChild().
func (n *Node) HasOneChild() bool {
	return n.first != nil && n.first == n.last
}

// HasTwoChildren mirrors Java's Node.hasTwoChildren().
func (n *Node) HasTwoChildren() bool {
	return n.first != nil && n.first.next != nil && n.first.next == n.last
}

// HasMoreThanOneChild mirrors Java's Node.hasMoreThanOneChild().
func (n *Node) HasMoreThanOneChild() bool {
	return n.first != nil && n.first.next != nil
}

// RemoveFirstChild detaches and returns the first child of n, or nil if n has
// no children. Mirrors Java's Node.removeFirstChild().
func (n *Node) RemoveFirstChild() *Node {
	if n.first == nil {
		return nil
	}
	c := n.first
	n.first = c.next
	if n.first != nil {
		n.first.previous = nil
	} else {
		n.last = nil
	}
	c.parent, c.next, c.previous = nil, nil, nil
	return c
}

// ChildCount returns the number of direct children. Mirrors Java's Node.getChildCount().
func (n *Node) ChildCount() int {
	count := 0
	for c := n.first; c != nil; c = c.next {
		count++
	}
	return count
}

// IsTry mirrors Java's Node.isTry().
func (n *Node) IsTry() bool { return n.token == TRY }

// IsClass mirrors Java's Node.isClass().
func (n *Node) IsClass() bool { return n.token == CLASS }

// IsImport mirrors Java's Node.isImport().
func (n *Node) IsImport() bool { return n.token == IMPORT }

// IsExport mirrors Java's Node.isExport().
func (n *Node) IsExport() bool { return n.token == EXPORT }

// IsImportSpec mirrors Java's Node.isImportSpec().
func (n *Node) IsImportSpec() bool { return n.token == IMPORT_SPEC }

// IsExportSpec mirrors Java's Node.isExportSpec().
func (n *Node) IsExportSpec() bool { return n.token == EXPORT_SPEC }

// IsObjectPattern mirrors Java's Node.isObjectPattern().
func (n *Node) IsObjectPattern() bool { return n.token == OBJECT_PATTERN }

// IsArrayPattern mirrors Java's Node.isArrayPattern().
func (n *Node) IsArrayPattern() bool { return n.token == ARRAY_PATTERN }

// IsDestructuringLhs mirrors Java's Node.isDestructuringLhs().
func (n *Node) IsDestructuringLhs() bool { return n.token == DESTRUCTURING_LHS }

// IsParamList mirrors Java's Node.isParamList().
func (n *Node) IsParamList() bool { return n.token == PARAM_LIST }

// IsLabel mirrors Java's Node.isLabel().
func (n *Node) IsLabel() bool { return n.token == LABEL }

// IsVanillaFor mirrors Java's Node.isVanillaFor() — true for a plain FOR node
// (not FOR_IN or FOR_OF).
func (n *Node) IsVanillaFor() bool { return n.token == FOR }

// IsModuleBody mirrors Java's Node.isModuleBody().
func (n *Node) IsModuleBody() bool { return n.token == MODULE_BODY }

// IsRest mirrors Java's Node.isRest().
func (n *Node) IsRest() bool { return n.token == ITER_REST || n.token == OBJECT_REST }

// IsExprResult mirrors Java's Node.isExprResult().
func (n *Node) IsExprResult() bool { return n.token == EXPR_RESULT }

// IsVar mirrors Java's Node.isVar().
func (n *Node) IsVar() bool { return n.token == VAR }

// IsConst mirrors Java's Node.isConst().
func (n *Node) IsConst() bool { return n.token == CONST }

// IsLet mirrors Java's Node.isLet().
func (n *Node) IsLet() bool { return n.token == LET }

// IsCase mirrors Java's Node.isCase().
func (n *Node) IsCase() bool { return n.token == CASE }

// IsDefaultCase mirrors Java's Node.isDefaultCase().
func (n *Node) IsDefaultCase() bool { return n.token == DEFAULT_CASE }

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

// IsGeneratorFunction mirrors Java's Node.isGeneratorFunction() — true when
// PropGeneratorFn is set on a FUNCTION node.
func (n *Node) IsGeneratorFunction() bool { return n.GetBooleanProp(PropGeneratorFn) }

// SetIsGeneratorFunction mirrors Java's Node.setIsGeneratorFunction(bool).
func (n *Node) SetIsGeneratorFunction(b bool) { n.PutBooleanProp(PropGeneratorFn, b) }

// IsArrowFunction mirrors Java's Node.isArrowFunction() — true when
// PropArrowFn is set on a FUNCTION node.
func (n *Node) IsArrowFunction() bool { return n.GetBooleanProp(PropArrowFn) }

// SetIsArrowFunction mirrors Java's Node.setIsArrowFunction(bool).
func (n *Node) SetIsArrowFunction(b bool) { n.PutBooleanProp(PropArrowFn, b) }

// IsAsyncFunction mirrors Java's Node.isAsyncFunction() — true when
// PropAsyncFn is set on a FUNCTION node.
func (n *Node) IsAsyncFunction() bool { return n.GetBooleanProp(PropAsyncFn) }

// SetIsAsyncFunction mirrors Java's Node.setIsAsyncFunction(bool).
func (n *Node) SetIsAsyncFunction(b bool) { n.PutBooleanProp(PropAsyncFn, b) }

// IsYieldAll mirrors Java's Node.isYieldAll() — true when PropYieldAll is
// set on a YIELD node (i.e. yield*).
func (n *Node) IsYieldAll() bool { return n.GetBooleanProp(PropYieldAll) }

// SetYieldAll mirrors Java's Node.setYieldAll(bool).
func (n *Node) SetYieldAll(b bool) { n.PutBooleanProp(PropYieldAll, b) }

// IsExportAll mirrors Java's Node.getBooleanProp(EXPORT_ALL_FROM).
func (n *Node) IsExportAll() bool { return n.GetBooleanProp(PropExportAll) }

// IsExportDefault mirrors Java's Node.getBooleanProp(EXPORT_DEFAULT).
func (n *Node) IsExportDefault() bool { return n.GetBooleanProp(PropExportDefault) }

// IsDestructuringPattern mirrors Java's Node.isDestructuringPattern() — true
// when this node is ARRAY_PATTERN or OBJECT_PATTERN.
func (n *Node) IsDestructuringPattern() bool {
	return n.token == ARRAY_PATTERN || n.token == OBJECT_PATTERN
}

// IsMemberFunctionDef mirrors Java's Node.isMemberFunctionDef().
func (n *Node) IsMemberFunctionDef() bool { return n.token == MEMBER_FUNCTION_DEF }

// IsMemberFieldDef mirrors Java's Node.isMemberFieldDef().
func (n *Node) IsMemberFieldDef() bool { return n.token == MEMBER_FIELD_DEF }

// IsClassMembers mirrors Java's Node.isClassMembers().
func (n *Node) IsClassMembers() bool { return n.token == CLASS_MEMBERS }

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

// ---------------- source length ------------------------------------------ //

// GetLength ports Node.getLength() — the source length in characters.
// Zero when not set (default).
func (n *Node) GetLength() int { return n.length }

// SetLength ports Node.setLength(int).
func (n *Node) SetLength(l int) { n.length = l }

// HasLength returns true when the node has a non-zero source length.
// Convenience wrapper used by tests that mirror Java's assertNode.hasLength().
func (n *Node) HasLength() bool { return n.length > 0 }

// ---------------- additional source-info ---------------------------------- //

// GetFirstFirstChild mirrors Java's Node.getFirstFirstChild(). Returns the
// first child of the first child, or nil when n has fewer than two levels.
func (n *Node) GetFirstFirstChild() *Node {
	if n.first == nil {
		return nil
	}
	return n.first.first
}

// GetOnlyChild mirrors Java's Node.getOnlyChild() — returns the sole child
// when n has exactly one child. Panics if n has zero or more than one child.
// Callers that are not certain should use FirstChild().
func (n *Node) GetOnlyChild() *Node {
	if n.first == nil || n.first != n.last {
		panic("Node.GetOnlyChild: node does not have exactly one child")
	}
	return n.first
}

// IsDescendantOf mirrors Java's Node.isDescendantOf(Node) — true when n is
// a descendant (not necessarily direct child) of the given ancestor.
func (n *Node) IsDescendantOf(ancestor *Node) bool {
	cur := n.parent
	for cur != nil {
		if cur == ancestor {
			return true
		}
		cur = cur.parent
	}
	return false
}

// IsOnlyChildOf mirrors Java's Node.isOnlyChildOf — true when n is the only
// child of possibleParent.
func (n *Node) IsOnlyChildOf(possibleParent *Node) bool {
	if possibleParent == nil {
		return false
	}
	return possibleParent.first == n && possibleParent.last == n
}

// IsSecondChildOf mirrors Java's Node.isSecondChildOf.
func (n *Node) IsSecondChildOf(possibleParent *Node) bool {
	if possibleParent == nil {
		return false
	}
	second := possibleParent.GetSecondChild()
	return second != nil && second == n
}

// GetChildAtIndex mirrors Java's Node.getChildAtIndex(int). Returns the child
// at 0-based index i, or nil if out of bounds. O(n).
func (n *Node) GetChildAtIndex(i int) *Node {
	j := 0
	for c := n.first; c != nil; c = c.next {
		if j == i {
			return c
		}
		j++
	}
	return nil
}

// GetIndexOfChild mirrors Java's Node.getIndexOfChild(Node). Returns the
// 0-based index of child in n's child list, or -1 if not found. O(n).
func (n *Node) GetIndexOfChild(child *Node) int {
	i := 0
	for c := n.first; c != nil; c = c.next {
		if c == child {
			return i
		}
		i++
	}
	return -1
}

// ---------------- bulk AST token predicates (1:1 Java port) -------------- //
// Each mirrors the identically-named final boolean isXxx() in Node.java.

// IsAdd mirrors Java's Node.isAdd().
func (n *Node) IsAdd() bool { return n.token == ADD }

// IsSub mirrors Java's Node.isSub().
func (n *Node) IsSub() bool { return n.token == SUB }

// IsAnd mirrors Java's Node.isAnd().
func (n *Node) IsAnd() bool { return n.token == AND }

// IsAssignAnd mirrors Java's Node.isAssignAnd().
func (n *Node) IsAssignAnd() bool { return n.token == ASSIGN_AND }

// IsArrayLit mirrors Java's Node.isArrayLit().
func (n *Node) IsArrayLit() bool { return n.token == ARRAYLIT }

// IsAssign mirrors Java's Node.isAssign().
func (n *Node) IsAssign() bool { return n.token == ASSIGN }

// IsAssignAdd mirrors Java's Node.isAssignAdd().
func (n *Node) IsAssignAdd() bool { return n.token == ASSIGN_ADD }

// IsNormalBlock mirrors Java's Node.isNormalBlock() which returns isBlock().
func (n *Node) IsNormalBlock() bool { return n.token == BLOCK }

// IsRoot mirrors Java's Node.isRoot().
func (n *Node) IsRoot() bool { return n.token == ROOT }

// IsAwait mirrors Java's Node.isAwait().
func (n *Node) IsAwait() bool { return n.token == AWAIT }

// IsBitNot mirrors Java's Node.isBitNot().
func (n *Node) IsBitNot() bool { return n.token == BITNOT }

// IsBreak mirrors Java's Node.isBreak().
func (n *Node) IsBreak() bool { return n.token == BREAK }

// IsCast mirrors Java's Node.isCast().
func (n *Node) IsCast() bool { return n.token == CAST }

// IsCatch mirrors Java's Node.isCatch().
func (n *Node) IsCatch() bool { return n.token == CATCH }

// IsComma mirrors Java's Node.isComma().
func (n *Node) IsComma() bool { return n.token == COMMA }

// IsContinue mirrors Java's Node.isContinue().
func (n *Node) IsContinue() bool { return n.token == CONTINUE }

// IsDebugger mirrors Java's Node.isDebugger().
func (n *Node) IsDebugger() bool { return n.token == DEBUGGER }

// IsDec mirrors Java's Node.isDec().
func (n *Node) IsDec() bool { return n.token == DEC }

// IsDefaultValue mirrors Java's Node.isDefaultValue().
func (n *Node) IsDefaultValue() bool { return n.token == DEFAULT_VALUE }

// IsDelProp mirrors Java's Node.isDelProp().
func (n *Node) IsDelProp() bool { return n.token == DELPROP }

// IsDo mirrors Java's Node.isDo().
func (n *Node) IsDo() bool { return n.token == DO }

// IsEmpty mirrors Java's Node.isEmpty().
func (n *Node) IsEmpty() bool { return n.token == EMPTY }

// IsExponent mirrors Java's Node.isExponent().
func (n *Node) IsExponent() bool { return n.token == EXPONENT }

// IsAssignExponent mirrors Java's Node.isAssignExponent().
func (n *Node) IsAssignExponent() bool { return n.token == ASSIGN_EXPONENT }

// [duplicate of first declaration removed]

// IsExportSpecs mirrors Java's Node.isExportSpecs().
func (n *Node) IsExportSpecs() bool { return n.token == EXPORT_SPECS }

// [duplicate of first declaration removed]

// IsFalse mirrors Java's Node.isFalse().
func (n *Node) IsFalse() bool { return n.token == FALSE }

// [duplicate of first declaration removed]

// IsForIn mirrors Java's Node.isForIn().
func (n *Node) IsForIn() bool { return n.token == FOR_IN }

// IsForOf mirrors Java's Node.isForOf().
func (n *Node) IsForOf() bool { return n.token == FOR_OF }

// IsForAwaitOf mirrors Java's Node.isForAwaitOf().
func (n *Node) IsForAwaitOf() bool { return n.token == FOR_AWAIT_OF }

// IsHook mirrors Java's Node.isHook() — true for the HOOK (ternary) token.
func (n *Node) IsHook() bool { return n.token == HOOK }

// IsIf mirrors Java's Node.isIf().
func (n *Node) IsIf() bool { return n.token == IF }

// [duplicate of first declaration removed]

// IsImportMeta mirrors Java's Node.isImportMeta().
func (n *Node) IsImportMeta() bool { return n.token == IMPORT_META }

// IsImportStar mirrors Java's Node.isImportStar().
func (n *Node) IsImportStar() bool { return n.token == IMPORT_STAR }

// [duplicate of first declaration removed]

// IsImportSpecs mirrors Java's Node.isImportSpecs().
func (n *Node) IsImportSpecs() bool { return n.token == IMPORT_SPECS }

// IsInc mirrors Java's Node.isInc().
func (n *Node) IsInc() bool { return n.token == INC }

// IsInstanceOf mirrors Java's Node.isInstanceOf().
func (n *Node) IsInstanceOf() bool { return n.token == INSTANCEOF }

// [duplicate of first declaration removed]

// IsLabelName mirrors Java's Node.isLabelName().
func (n *Node) IsLabelName() bool { return n.token == LABEL_NAME }

// [duplicate of first declaration removed]

// IsMemberVariableDef mirrors Java's Node.isMemberVariableDef().
// In Go this token is represented as MEMBER_FIELD_DEF (no separate MEMBER_VARIABLE_DEF).
func (n *Node) IsMemberVariableDef() bool { return n.token == MEMBER_FIELD_DEF }

// IsComputedFieldDef mirrors Java's Node.isComputedFieldDef().
func (n *Node) IsComputedFieldDef() bool { return n.token == COMPUTED_PROP }

// [duplicate of first declaration removed]

// IsNE mirrors Java's Node.isNE().
func (n *Node) IsNE() bool { return n.token == NE }

// IsSHNE mirrors Java's Node.isSHNE() — strict not-equal (!==).
func (n *Node) IsSHNE() bool { return n.token == SHNE }

// IsEQ mirrors Java's Node.isEQ() — loose equality (==).
func (n *Node) IsEQ() bool { return n.token == EQ }

// IsSHEQ mirrors Java's Node.isSHEQ() — strict equality (===).
func (n *Node) IsSHEQ() bool { return n.token == SHEQ }

// IsNeg mirrors Java's Node.isNeg() — unary negation.
func (n *Node) IsNeg() bool { return n.token == NEG }

// IsNot mirrors Java's Node.isNot() — logical NOT (!).
func (n *Node) IsNot() bool { return n.token == NOT }

// IsNull mirrors Java's Node.isNull().
func (n *Node) IsNull() bool { return n.token == NULL }

// IsNullishCoalesce mirrors Java's Node.isNullishCoalesce() — the ?? operator.
func (n *Node) IsNullishCoalesce() bool { return n.token == COALESCE }

// IsAssignNullishCoalesce mirrors Java's Node.isAssignNullishCoalesce() — ??=.
func (n *Node) IsAssignNullishCoalesce() bool { return n.token == ASSIGN_COALESCE }

// IsNumber mirrors Java's Node.isNumber().
func (n *Node) IsNumber() bool { return n.token == NUMBER }

// IsObjectLit mirrors Java's Node.isObjectLit().
func (n *Node) IsObjectLit() bool { return n.token == OBJECTLIT }

// [duplicate of first declaration removed]

// IsOr mirrors Java's Node.isOr() — logical OR (||).
func (n *Node) IsOr() bool { return n.token == OR }

// IsAssignOr mirrors Java's Node.isAssignOr() — ||=.
func (n *Node) IsAssignOr() bool { return n.token == ASSIGN_OR }

// [duplicate of first declaration removed]

// IsRegExp mirrors Java's Node.isRegExp().
func (n *Node) IsRegExp() bool { return n.token == REGEXP }

// [duplicate of first declaration removed]

// IsObjectRest mirrors Java's Node.isObjectRest().
func (n *Node) IsObjectRest() bool { return n.token == OBJECT_REST }

// IsSuper mirrors Java's Node.isSuper().
func (n *Node) IsSuper() bool { return n.token == SUPER }

// IsSwitch mirrors Java's Node.isSwitch().
func (n *Node) IsSwitch() bool { return n.token == SWITCH }

// IsSwitchBody mirrors Java's Node.isSwitchBody().
func (n *Node) IsSwitchBody() bool { return n.token == SWITCH_BODY }

// IsTaggedTemplateLit mirrors Java's Node.isTaggedTemplateLit().
func (n *Node) IsTaggedTemplateLit() bool { return n.token == TAGGED_TEMPLATELIT }

// IsTemplateLit mirrors Java's Node.isTemplateLit().
func (n *Node) IsTemplateLit() bool { return n.token == TEMPLATELIT }

// IsTemplateLitString mirrors Java's Node.isTemplateLitString().
func (n *Node) IsTemplateLitString() bool { return n.token == TEMPLATELIT_STRING }

// IsTemplateLitSub mirrors Java's Node.isTemplateLitSub().
func (n *Node) IsTemplateLitSub() bool { return n.token == TEMPLATELIT_SUB }

// IsThis mirrors Java's Node.isThis().
func (n *Node) IsThis() bool { return n.token == THIS }

// IsThrow mirrors Java's Node.isThrow().
func (n *Node) IsThrow() bool { return n.token == THROW }

// IsTrue mirrors Java's Node.isTrue().
func (n *Node) IsTrue() bool { return n.token == TRUE }

// [duplicate of first declaration removed]

// IsTypeOf mirrors Java's Node.isTypeOf().
func (n *Node) IsTypeOf() bool { return n.token == TYPEOF }

// [duplicate of first declaration removed]

// IsVoid mirrors Java's Node.isVoid().
func (n *Node) IsVoid() bool { return n.token == VOID }

// IsWhile mirrors Java's Node.isWhile().
func (n *Node) IsWhile() bool { return n.token == WHILE }

// IsWith mirrors Java's Node.isWith().
func (n *Node) IsWith() bool { return n.token == WITH }

// IsYield mirrors Java's Node.isYield().
func (n *Node) IsYield() bool { return n.token == YIELD }

// IsDeclare mirrors Java's Node.isDeclare() — TypeScript declare modifier.
func (n *Node) IsDeclare() bool { return n.token == DECLARE }

// IsString mirrors Java's Node.isString() — same as isStringLit().
func (n *Node) IsString() bool { return n.token == STRINGLIT }

// IsGetter mirrors Java's Node.isGetterDef() — same as IsGetterDef.
// Provided for callers using the Java name directly.
func (n *Node) IsGetter() bool { return n.GetBooleanProp(PropGetterDef) }

// IsSetter mirrors Java's Node.isSetterDef() — same as IsSetterDef.
func (n *Node) IsSetter() bool { return n.GetBooleanProp(PropSetterDef) }

// IsSyntheticBlock mirrors Java's Node.isSyntheticBlock() — true when the
// PropSynthetic flag is set on a BLOCK node.
func (n *Node) IsSyntheticBlock() bool {
	return n.token == BLOCK && n.GetBooleanProp(PropSynthetic)
}

// SetSyntheticBlock mirrors Java's Node.setIsSyntheticBlock(boolean).
func (n *Node) SetSyntheticBlock(b bool) { n.PutBooleanProp(PropSynthetic, b) }

// IsUnscopedQualifiedName mirrors Java's Node.isUnscopedQualifiedName().
// True for NAME (non-empty) and GETPROP whose receiver is an unscoped qualified name.
func (n *Node) IsUnscopedQualifiedName() bool {
	switch n.token {
	case NAME:
		return n.str != ""
	case GETPROP:
		if fc := n.FirstChild(); fc != nil {
			return fc.IsUnscopedQualifiedName()
		}
		return false
	}
	return false
}

// IsValidAssignmentTarget mirrors Java's Node.isValidAssignmentTarget().
func (n *Node) IsValidAssignmentTarget() bool {
	switch n.token {
	case NAME, GETPROP, GETELEM, ARRAY_PATTERN, OBJECT_PATTERN:
		return true
	}
	return false
}

// GetDouble mirrors Java's Node.getDouble() — alias for Number().
func (n *Node) GetDouble() float64 { return n.num }

// SetDouble mirrors Java's Node.setDouble(double) — alias for setting num.
func (n *Node) SetDouble(v float64) { n.num = v }

// GetString mirrors Java's Node.getString() — alias for String_().
func (n *Node) GetString() string { return n.str }

// SetString_ already exists above as an alias. GetString is the Java name.

// SetToken mirrors Java's Node.setToken(Token).
func (n *Node) SetToken(tok Token) { n.token = tok }

// [duplicate of first declaration removed]

