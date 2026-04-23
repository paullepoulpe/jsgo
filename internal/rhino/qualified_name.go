// Port of com/google/javascript/rhino/QualifiedName.java @ closure-compiler v20260415
//
// Abstraction over a qualified name. Unifies Node-based qualified names and
// string-based names, allowing lazy parsing of strings and representing a
// pre-parsed qualified name without the overhead of a whole Node.
// A qualified name is a linked list of components, starting from the outermost
// property access and ending with the root (a simple name with no owner).
//
// Java abstract class → Go interface.
// Java private inner classes StringListQname and GetpropQname → unexported structs.
// Java static final class NodeQname → exported struct (was package-private in Java).

package rhino

import (
	"fmt"
	"strings"
)

// qnameThis and qnameSuper mirror the Java static finals
// RhinoStringPool.addOrGet("this") / ("super"). In Go there is no string
// interning so we use plain constants and == comparison.
const qnameThis = "this"
const qnameSuper = "super"

// QualifiedName mirrors the Java abstract class QualifiedName.
// All method names match their Java counterparts exactly.
type QualifiedName interface {
	// GetOwner returns the qualified name of the owner, or nil for simple
	// names. For "foo.bar.baz" this returns an object representing "foo.bar".
	GetOwner() QualifiedName

	// GetComponent returns the outermost term of this qualified name, or the
	// entire name for simple names. For "foo.bar.baz" this returns "baz".
	GetComponent() string

	// IsSimple returns true if this is a simple (single-component) name.
	IsSimple() bool

	// appendTo appends the joined qualified name to the given builder.
	appendTo(sb *strings.Builder)

	// Matches checks whether the given node matches this name.
	Matches(n *Node) bool

	// GetRoot returns the root of this name, e.g. "foo" from "foo.bar.baz".
	GetRoot() string

	// Components returns the components as a slice of strings, starting at the
	// root. For "foo.bar.baz" this returns ["foo", "bar", "baz"].
	Components() []string

	// Join returns the qualified name as a dot-joined string.
	Join() string

	// GetComponentCount returns the number of components in this name.
	GetComponentCount() int

	// Getprop returns a new qualified name with this name as the owner and the
	// given string as the property name.
	Getprop(propertyName string) QualifiedName
}

// qualifiedNameBase embeds the shared concrete methods so each concrete type
// only needs to implement the four abstract primitives.
type qualifiedNameBase struct {
	self QualifiedName // back-pointer set by each constructor
}

func (b *qualifiedNameBase) GetRoot() string {
	name := b.self
	for !name.IsSimple() {
		name = name.GetOwner()
	}
	return name.GetComponent()
}

func (b *qualifiedNameBase) Components() []string {
	var out []string
	buildComponents(b.self, &out)
	return out
}

func buildComponents(q QualifiedName, out *[]string) {
	owner := q.GetOwner()
	if owner != nil {
		buildComponents(owner, out)
	}
	*out = append(*out, q.GetComponent())
}

func (b *qualifiedNameBase) Join() string {
	var sb strings.Builder
	b.self.appendTo(&sb)
	return sb.String()
}

func (b *qualifiedNameBase) Getprop(propertyName string) QualifiedName {
	return newGetpropQname(b.self, propertyName)
}

// ── Of ──────────────────────────────────────────────────────────────────── //

// Of ports QualifiedName.of(String). It parses a dot-separated name string
// into a StringListQname.
func Of(s string) QualifiedName {
	terms := strings.Split(s, ".")
	return newStringListQname(terms, len(terms))
}

// ── StringListQname ──────────────────────────────────────────────────────── //

// stringListQname ports the private static inner class StringListQname.
type stringListQname struct {
	qualifiedNameBase
	terms []string
	size  int
}

func newStringListQname(terms []string, size int) *stringListQname {
	q := &stringListQname{terms: terms, size: size}
	q.qualifiedNameBase.self = q
	return q
}

func (q *stringListQname) GetOwner() QualifiedName {
	if q.size > 1 {
		return newStringListQname(q.terms, q.size-1)
	}
	return nil
}

func (q *stringListQname) GetComponent() string {
	return q.terms[q.size-1]
}

func (q *stringListQname) IsSimple() bool {
	return q.size == 1
}

func (q *stringListQname) appendTo(sb *strings.Builder) {
	for i := 0; i < q.size; i++ {
		if i > 0 {
			sb.WriteByte('.')
		}
		sb.WriteString(q.terms[i])
	}
}

// Components overrides the base to use a direct slice rather than recursive
// construction, matching the Java override in StringListQname.
func (q *stringListQname) Components() []string {
	return q.terms[:q.size]
}

func (q *stringListQname) GetComponentCount() int {
	return q.size
}

func (q *stringListQname) Matches(n *Node) bool {
	pos := q.size - 1
	for pos > 0 && n.token == GETPROP {
		if n.str != q.terms[pos] {
			return false
		}
		pos--
		n = n.FirstChild()
	}
	if pos > 0 {
		return false
	}
	term := q.terms[0]
	switch n.token {
	case NAME, MEMBER_FUNCTION_DEF:
		return n.str == term
	case THIS:
		return term == qnameThis
	case SUPER:
		return term == qnameSuper
	default:
		return false
	}
}

// ── GetpropQname ─────────────────────────────────────────────────────────── //

// getpropQname ports the private static inner class GetpropQname.
type getpropQname struct {
	qualifiedNameBase
	owner QualifiedName
	prop  string
}

func newGetpropQname(owner QualifiedName, prop string) *getpropQname {
	q := &getpropQname{owner: owner, prop: prop}
	q.qualifiedNameBase.self = q
	return q
}

func (q *getpropQname) GetOwner() QualifiedName { return q.owner }
func (q *getpropQname) GetComponent() string    { return q.prop }
func (q *getpropQname) IsSimple() bool          { return false }

func (q *getpropQname) appendTo(sb *strings.Builder) {
	q.owner.appendTo(sb)
	sb.WriteByte('.')
	sb.WriteString(q.prop)
}

func (q *getpropQname) GetComponentCount() int {
	return q.owner.GetComponentCount() + 1
}

func (q *getpropQname) Matches(n *Node) bool {
	return n.token == GETPROP &&
		n.str == q.prop &&
		q.owner.Matches(n.FirstChild())
}

// ── NodeQname ────────────────────────────────────────────────────────────── //

// NodeQname ports the package-private static final class NodeQname.
// It wraps a *Node that is assumed to be a valid qualified name.
type NodeQname struct {
	qualifiedNameBase
	node *Node
}

// NewNodeQname constructs a NodeQname from a *Node.
// Mirrors the NodeQname(Node n) constructor in Java.
func NewNodeQname(n *Node) *NodeQname {
	q := &NodeQname{node: n}
	q.qualifiedNameBase.self = q
	return q
}

func (q *NodeQname) GetOwner() QualifiedName {
	if q.node.token == GETPROP {
		return NewNodeQname(q.node.FirstChild())
	}
	return nil
}

func (q *NodeQname) GetComponent() string {
	switch q.node.token {
	case THIS:
		return qnameThis
	case SUPER:
		return qnameSuper
	case NAME, GETPROP, MEMBER_FUNCTION_DEF:
		return q.node.str
	default:
		panic(fmt.Sprintf("Not a qualified name: %v", q.node.token))
	}
}

func (q *NodeQname) IsSimple() bool {
	return q.node.token != GETPROP
}

func (q *NodeQname) appendTo(sb *strings.Builder) {
	sb.WriteString(q.Join())
}

// Join overrides the base Join to delegate to GetQualifiedName, matching the
// Java NodeQname override.
func (q *NodeQname) Join() string {
	s := q.node.GetQualifiedName()
	if s == "" {
		return ""
	}
	return s
}

func (q *NodeQname) GetComponentCount() int {
	count := 1
	current := q.node
	for current.token == GETPROP {
		count++
		current = current.FirstChild()
	}
	return count
}

func (q *NodeQname) Matches(n *Node) bool {
	return n.MatchesQualifiedName(q.node)
}

// ── Node helpers added for QualifiedName ────────────────────────────────── //

// GetQualifiedName ports Node.getQualifiedName() from Node.java.
// Returns "" when the node is not a valid qualified name (Java returns null).
func (n *Node) GetQualifiedName() string {
	switch n.token {
	case NAME:
		if n.str == "" {
			return ""
		}
		return n.str
	case GETPROP:
		return getQualifiedNameForGetProp(n)
	case THIS:
		return qnameThis
	case SUPER:
		return qnameSuper
	default:
		return ""
	}
}

// getQualifiedNameForGetProp is the recursive helper for GetQualifiedName on
// GETPROP nodes (mirrors private Node.getQualifiedNameForGetProp in Java).
func getQualifiedNameForGetProp(n *Node) string {
	propName := n.str
	first := n.FirstChild()
	if first == nil {
		return ""
	}
	var left string
	if first.token == GETPROP {
		left = getQualifiedNameForGetProp(first)
	} else {
		left = first.GetQualifiedName()
	}
	if left == "" {
		return ""
	}
	return left + "." + propName
}

// MatchesQualifiedName ports Node.matchesQualifiedName(Node) from Node.java.
func (n *Node) MatchesQualifiedName(other *Node) bool {
	if other.token != n.token {
		return false
	}
	switch n.token {
	case NAME:
		return n.str != "" && other.str != "" && n.str == other.str
	case THIS, SUPER:
		return true
	case GETPROP:
		return n.str == other.str &&
			n.FirstChild().MatchesQualifiedName(other.FirstChild())
	default:
		return false
	}
}

// GetQualifiedNameObject ports Node.getQualifiedNameObject() from Node.java.
// Returns nil when the node is not a valid qualified name.
func (n *Node) GetQualifiedNameObject() QualifiedName {
	if n.isQualifiedName() {
		return NewNodeQname(n)
	}
	return nil
}

// isQualifiedName ports Node.isQualifiedName() — private helper.
func (n *Node) isQualifiedName() bool {
	switch n.token {
	case NAME:
		return n.str != ""
	case THIS, SUPER:
		return true
	case GETPROP:
		return n.FirstChild().isQualifiedName()
	default:
		return false
	}
}
