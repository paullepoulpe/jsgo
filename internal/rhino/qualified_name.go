// Port of com/google/javascript/rhino/QualifiedName.java @ closure-compiler v20260415

/*
 *
 * ***** BEGIN LICENSE BLOCK *****
 * Version: MPL 1.1
 *
 * The contents of this file are subject to the Mozilla Public License Version
 * 1.1 (the "License"); you may not use this file except in compliance with
 * the License. You may obtain a copy of the License at
 * http://www.mozilla.org/MPL/
 *
 * Software distributed under the License is distributed on an "AS IS" basis,
 * WITHOUT WARRANTY OF ANY KIND, either express or implied. See the License
 * for the specific language governing rights and limitations under the
 * License.
 *
 * The Original Code is Rhino code, released
 * May 6, 1999.
 *
 * The Initial Developer of the Original Code is
 * Netscape Communications Corporation.
 * Portions created by the Initial Developer are Copyright (C) 1997-1999
 * the Initial Developer. All Rights Reserved.
 *
 * Contributor(s):
 *   Norris Boyd
 *   Roger Lawrence
 *   Mike McCabe
 *
 * This Go translation is distributed under the MPL 1.1 only, as the
 * original's dual licence allows; the other licence it offers is not
 * used.
 *
 * ***** END LICENSE BLOCK ***** */

package rhino

import (
	"strings"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// QualifiedName is an abstraction over a qualified name. Unifies Node-based
// qualified names and string-based names, allowing to lazily parse strings
// and represent a pre-parsed qualified name without the overhead of a whole
// Node. Essentially, a qualified name is a linked list of components,
// starting from the outermost property access and ending with the root of
// the name, which is a simple name with no owner.
type QualifiedName interface {
	// GetOwner returns the qualified name of the owner, or null for simple
	// names. For the name "foo.bar.baz", this returns an object representing
	// "foo.bar".
	GetOwner() QualifiedName

	// GetComponent returns outer-most term of this qualified name, or the
	// entire name for simple names. For the name "foo.bar.baz", this returns
	// "baz".
	GetComponent() string

	// IsSimple returns true if this is a simple name.
	IsSimple() bool

	// AppendTo appends the joined qualified name to the given StringBuilder.
	AppendTo(sb *strings.Builder)

	// Matches checks whether the given node matches this name.
	Matches(n *Node) bool

	GetRoot() string
	Components() []string
	BuildComponents(builder *[]string)
	Join() string
	GetComponentCount() int
	Getprop(propertyName string) QualifiedName

	GetQualifiedNameBase() *QualifiedNameBase
}

// QualifiedNameBase holds QualifiedName's concrete methods; self dispatches
// the abstract ones.
type QualifiedNameBase struct {
	self QualifiedName
}

// All subclasses must be defined in this file.
//
// Java: QualifiedName()
func (q *QualifiedNameBase) InitQualifiedNameBase(self QualifiedName) {
	q.self = self
}

// GetQualifiedNameBase gives access to the base through the interface.
func (q *QualifiedNameBase) GetQualifiedNameBase() *QualifiedNameBase { return q }

// Java: QualifiedName.of(String)
func QualifiedName_Of(string_ string) QualifiedName {
	lastIndex := 0
	var index int
	var builder []string
	for {
		index = jutil.StringIndexOfCharFrom(string_, '.', lastIndex)
		end := index
		if index < 0 {
			end = jutil.StringLength(string_)
		}
		term := jutil.StringSubstring(string_, lastIndex, end)
		builder = append(builder, RhinoStringPool_AddOrGet(term))
		lastIndex = index + 1
		if index < 0 {
			break
		}
	}
	terms := builder
	return NewQualifiedName_StringListQname(terms, len(terms))
}

// GetRoot returns the root of this name, e.g. "foo" from "foo.bar.baz".
//
// Java: QualifiedName.getRoot()
func (q *QualifiedNameBase) GetRoot() string {
	name := q.self
	for !name.IsSimple() {
		name = name.GetOwner()
	}
	return name.GetComponent()
}

// Components returns the components of this name as an iterable of strings,
// starting at the root. For the qualified name foo.bar.baz, this returns
// ["foo", "bar", "baz"].
//
// Java: QualifiedName.components()
func (q *QualifiedNameBase) Components() []string {
	var components []string
	q.self.BuildComponents(&components)
	return components
}

// Java: QualifiedName.buildComponents(ImmutableList.Builder<String>)
func (q *QualifiedNameBase) BuildComponents(builder *[]string) {
	owner := q.self.GetOwner()
	if owner != nil {
		owner.BuildComponents(builder)
	}
	*builder = append(*builder, q.self.GetComponent())
}

// Join returns the qualified name as a string.
//
// Java: QualifiedName.join()
func (q *QualifiedNameBase) Join() string {
	var sb strings.Builder
	q.self.AppendTo(&sb)
	return sb.String()
}

// Getprop returns a new qualified name object with `this` name as the owner
// and the given string as the property name.
//
// Java: QualifiedName.getprop(String)
func (q *QualifiedNameBase) Getprop(propertyName string) QualifiedName {
	return NewQualifiedName_GetpropQname(q.self, propertyName)
}

// QualifiedName_StringListQname is a qualified name based on a list of
// string terms.
type QualifiedName_StringListQname struct {
	QualifiedNameBase
	terms []string
	size  int
}

// Java: QualifiedName.StringListQname(ImmutableList<String>, int)
func NewQualifiedName_StringListQname(terms []string, size int) *QualifiedName_StringListQname {
	q := &QualifiedName_StringListQname{}
	q.InitQualifiedNameBase(q)
	q.terms = terms
	q.size = size
	return q
}

// Java: QualifiedName.StringListQname.getOwner()
func (q *QualifiedName_StringListQname) GetOwner() QualifiedName {
	if q.size > 1 {
		return NewQualifiedName_StringListQname(q.terms, q.size-1)
	}
	return nil
}

// Java: QualifiedName.StringListQname.getComponent()
func (q *QualifiedName_StringListQname) GetComponent() string {
	return q.terms[q.size-1]
}

// Java: QualifiedName.StringListQname.isSimple()
func (q *QualifiedName_StringListQname) IsSimple() bool {
	return q.size == 1
}

// Java: QualifiedName.StringListQname.appendTo(StringBuilder)
func (q *QualifiedName_StringListQname) AppendTo(sb *strings.Builder) {
	for i := 0; i < q.size; i++ {
		if i > 0 {
			sb.WriteByte('.')
		}
		sb.WriteString(q.terms[i])
	}
}

// Java: QualifiedName.StringListQname.components()
func (q *QualifiedName_StringListQname) Components() []string {
	return q.terms[0:q.size]
}

// Java: QualifiedName.StringListQname.getComponentCount()
func (q *QualifiedName_StringListQname) GetComponentCount() int {
	return q.size
}

// Java: QualifiedName.StringListQname.matches(Node)
func (q *QualifiedName_StringListQname) Matches(n *Node) bool {
	pos := q.size - 1
	for pos > 0 && n.IsGetProp() {
		if !RhinoStringPool_UncheckedEquals(n.GetString(), q.terms[pos]) {
			return false
		}
		pos--
		n = n.GetFirstChild()
	}
	if pos > 0 {
		return false
	}

	term := q.terms[0]
	switch n.GetToken() {
	case Token_NAME, Token_MEMBER_FUNCTION_DEF:
		return RhinoStringPool_UncheckedEquals(term, n.GetString())
	case Token_THIS:
		return RhinoStringPool_UncheckedEquals(term, QualifiedName_THIS)
	case Token_SUPER:
		return RhinoStringPool_UncheckedEquals(term, QualifiedName_SUPER)
	default:
		return false
	}
}

// QualifiedName_GetpropQname is a qualified name built with an extra
// property access on an existing qualified name.
type QualifiedName_GetpropQname struct {
	QualifiedNameBase
	owner QualifiedName
	prop  string
}

// Java: QualifiedName.GetpropQname(QualifiedName, String)
func NewQualifiedName_GetpropQname(owner QualifiedName, prop string) *QualifiedName_GetpropQname {
	q := &QualifiedName_GetpropQname{}
	q.InitQualifiedNameBase(q)
	q.owner = owner
	q.prop = RhinoStringPool_AddOrGet(prop)
	return q
}

// Java: QualifiedName.GetpropQname.getOwner()
func (q *QualifiedName_GetpropQname) GetOwner() QualifiedName {
	return q.owner
}

// Java: QualifiedName.GetpropQname.getComponent()
func (q *QualifiedName_GetpropQname) GetComponent() string {
	return q.prop
}

// Java: QualifiedName.GetpropQname.isSimple()
func (q *QualifiedName_GetpropQname) IsSimple() bool {
	return false
}

// Java: QualifiedName.GetpropQname.appendTo(StringBuilder)
func (q *QualifiedName_GetpropQname) AppendTo(sb *strings.Builder) {
	q.owner.AppendTo(sb)
	sb.WriteByte('.')
	sb.WriteString(q.prop)
}

// Java: QualifiedName.GetpropQname.matches(Node)
func (q *QualifiedName_GetpropQname) Matches(n *Node) bool {
	return n.IsGetProp() &&
		RhinoStringPool_UncheckedEquals(n.GetString(), q.prop) &&
		q.owner.Matches(n.GetFirstChild())
}

// Java: QualifiedName.GetpropQname.getComponentCount()
func (q *QualifiedName_GetpropQname) GetComponentCount() int {
	return q.owner.GetComponentCount() + 1
}

// QualifiedName_NodeQname is a qualified name from a node. The precondition
// that the node is actually a qualified name is not actually checked here,
// though it will throw IllegalStateException eventually if it is not.
type QualifiedName_NodeQname struct {
	QualifiedNameBase
	node *Node
}

// Java: QualifiedName.NodeQname(Node)
func NewQualifiedName_NodeQname(n *Node) *QualifiedName_NodeQname {
	q := &QualifiedName_NodeQname{}
	q.InitQualifiedNameBase(q)
	q.node = n
	return q
}

// Java: QualifiedName.NodeQname.getOwner()
func (q *QualifiedName_NodeQname) GetOwner() QualifiedName {
	if q.node.IsGetProp() {
		return NewQualifiedName_NodeQname(q.node.GetFirstChild())
	}
	return nil
}

// Java: QualifiedName.NodeQname.getComponent()
func (q *QualifiedName_NodeQname) GetComponent() string {
	switch q.node.GetToken() {
	case Token_THIS:
		return QualifiedName_THIS
	case Token_SUPER:
		return QualifiedName_SUPER
	case Token_NAME, Token_GETPROP, Token_MEMBER_FUNCTION_DEF:
		return q.node.GetString()
	default:
		panic(jutil.NewIllegalStateExceptionMsg("Not a qualified name: " + q.node.String()))
	}
}

// Java: QualifiedName.NodeQname.isSimple()
func (q *QualifiedName_NodeQname) IsSimple() bool {
	return !q.node.IsGetProp()
}

// Java: QualifiedName.NodeQname.appendTo(StringBuilder)
func (q *QualifiedName_NodeQname) AppendTo(sb *strings.Builder) {
	sb.WriteString(q.Join())
}

// Join is the node's qualified name.
//
// Java: QualifiedName.NodeQname.join()
func (q *QualifiedName_NodeQname) Join() string {
	// Go-only: language-adapter — Java returns null when the node is not a
	// qualified name (which callers of getQualifiedNameObject exclude); here
	// that prints as "null", as Java's appendTo does.
	return jutil.ToJavaString(q.node.GetQualifiedName())
}

// Java: QualifiedName.NodeQname.matches(Node)
func (q *QualifiedName_NodeQname) Matches(n *Node) bool {
	return n.MatchesQualifiedNameNode(q.node)
}

// Java: QualifiedName.NodeQname.getComponentCount()
func (q *QualifiedName_NodeQname) GetComponentCount() int {
	count := 1
	current := q.node
	for current.IsGetProp() {
		count++
		current = current.GetFirstChild()
	}
	return count
}

// Java: QualifiedName.THIS
var QualifiedName_THIS = RhinoStringPool_AddOrGet("this")

// Java: QualifiedName.SUPER
var QualifiedName_SUPER = RhinoStringPool_AddOrGet("super")
