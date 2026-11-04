// Port of com/google/javascript/rhino/Node.java @ closure-compiler v20261028

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
	"fmt"
	"math/big"

	"github.com/paullepoulpe/jsgo/internal/base"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/strpool"
)

// Node implements the root of the intermediate representation.
//
// Go-only: language-adapter — Java's private subclasses NumberNode,
// BigIntNode, StringNode and TemplateLiteralSubstringNode only add a payload
// and override isEquivalentTo/cloneNode, so they are folded into Node
// (PORTING.md §4): `kind` records which class a node is, the payload fields
// are only meaningful for that kind, and the accessors throw the
// ClassCastException Java's casts throw on any other kind.

// node_kind records which Java class a Node is (see the Node doc comment).
type node_kind uint8

const (
	node_kindNode node_kind = iota
	node_kindString
	node_kindNumber
	node_kindBigInt
	node_kindTemplateLiteralSubstring
)

const (
	kindMask            uint8 = 7
	kindFlagTypeIsColor uint8 = 1 << 3
	kindFlagNoBlockDecl uint8 = 1 << 4
	// kindFlagHasExtra: the cold property list may hold the value of a
	// property that is normally stored inline (e.g. an int value other than
	// 1 for a boolean property, or CHANGE_TIME of a non-Node kind).
	kindFlagHasExtra uint8 = 1 << 5
	// kindFlagOrigNameIsString: the original name of this string node is its
	// own string (node_payload.go).
	kindFlagOrigNameIsString uint8 = 1 << 6
	// kindFlagHasExtBox: the aux entry holds a *nodeExtBox (otherwise the
	// *JSDocInfo when propFlagHasJSDoc is set).
	kindFlagHasExtBox uint8 = 1 << 7
)

// nodeExtBox holds the rarely-set parts of a node, in its store's cold-data
// map (node_tables.go).
type nodeExtBox struct {
	jsDocInfo      *JSDocInfo         // canonical *JSDocInfo (valid iff propFlagHasJSDoc != 0 && jsDocInfo != nil)
	propListHead   *Node_PropListItem // remaining cold properties
	typeBeforeCast uint32             // node type handle of the canonical TYPE_BEFORE_CAST (0 iff not present as JSType)
	length         int32              // the length, when Node.length16 is 0xFFFF
	cookedID       rhinoStringID      // the cooked string of a TEMPLATELIT_STRING, when it is not the raw one
	coldProps      uint16             // the 15 rare boolean properties (logical bits 32..46)
	cookedNil      bool               // the TEMPLATELIT_STRING has no cooked string
}

// The SOURCE_FILE property of a node is an id into its NodeStore's sourceFiles
// table (or staticStore's when staticBit is set): each PutProp(SOURCE_FILE, file)
// takes a new id in that store, holding the property's list item (shared by
// every node that copies the id, as Java shares the item).
type sourceFileSlot struct {
	item Node_PropListItem
}

func lookupSourceFileSlot(id uint32) *sourceFileSlot {
	if id&staticBit != 0 {
		return staticStore.sourceFiles.at(id &^ staticBit)
	}
	return CurrentNodeStore().sourceFiles.at(id)
}

func sourceFilePropItem(id uint32) *Node_PropListItem {
	return &lookupSourceFileSlot(id).item
}

func (s *NodeStore) allocSourceFileID(file any) uint32 {
	id := s.sourceFiles.add(&s.mu, sourceFileSlot{
		item: Node_PropListItem{
			objectValue: file,
			propType:    int8(Node_Prop_SOURCE_FILE),
			isObject:    true,
			shared:      true,
		},
	})
	if s.static {
		return staticBit | id
	}
	return id
}

// allocSourceFileID returns a new source-file id for file in the current store.
func allocSourceFileID(file any) uint32 {
	return CurrentNodeStore().allocSourceFileID(file)
}

// jsDocAux is the node's JSDocInfo when it is stored directly in the aux
// entry (no ext box).
func (n Node) jsDocAux() *JSDocInfo {
	if a, ok := n.auxOf(); ok {
		return a.jsDoc
	}
	return nil
}

// extBoxOrNil is the node's ext box, or nil if the node has none. Read
// paths use it; mutators use extBox.
func (n Node) extBoxOrNil() *nodeExtBox {
	if (n.rec().kindFlags & kindFlagHasExtBox) == 0 {
		return nil
	}
	if a, ok := n.auxOf(); ok {
		return a.ext
	}
	return nil
}

// extBox is the ext box of a node that has one (kindFlagHasExtBox set).
func (n Node) extBox() *nodeExtBox {
	if box := n.extBoxOrNil(); box != nil {
		return box
	}
	box := &nodeExtBox{}
	n.setAux(nodeAux{ext: box})
	return box
}

func (n Node) ensureExtBox() *nodeExtBox {
	if (n.rec().kindFlags & kindFlagHasExtBox) != 0 {
		return n.extBox()
	}
	box := &nodeExtBox{}
	if (n.props() & propFlagHasJSDoc) != 0 {
		box.jsDocInfo = n.jsDocAux()
	}
	n.rec().kindFlags |= kindFlagHasExtBox
	n.setAux(nodeAux{ext: box})
	return box
}

// maybeCollapseExtBox drops the box when only the JSDocInfo (or nothing) is
// left in it: the aux entry then refers to the JSDocInfo directly.
func (n Node) maybeCollapseExtBox(box *nodeExtBox) {
	if box.typeBeforeCast != 0 || box.propListHead != nil || n.extBoxHoldsNodeFields(box) {
		return
	}
	n.rec().kindFlags &^= kindFlagHasExtBox
	if (n.props()&propFlagHasJSDoc) != 0 && box.jsDocInfo != nil {
		n.setAux(nodeAux{jsDoc: box.jsDocInfo})
	} else {
		n.setAux(nodeAux{})
	}
}

// extBoxHoldsNodeFields reports whether box holds a length, cooked string, or
// rare boolean property that does not fit in the Node record itself.
func (n Node) extBoxHoldsNodeFields(box *nodeExtBox) bool {
	return n.rec().length16 == 0xFFFF || box.cookedID != 0 || box.cookedNil || box.coldProps != 0
}

func (n Node) coldPropListHead() *Node_PropListItem {
	if box := n.extBoxOrNil(); box != nil {
		return box.propListHead
	}
	return nil
}

func (n Node) setColdPropListHead(head *Node_PropListItem) {
	if head == nil {
		if (n.rec().kindFlags & kindFlagHasExtBox) == 0 {
			return
		}
		box := n.extBox()
		box.propListHead = nil
		n.maybeCollapseExtBox(box)
		return
	}
	box := n.ensureExtBox()
	box.propListHead = head
}

func (n Node) typeBeforeCastBase() *JSTypeBase {
	if box := n.extBoxOrNil(); box != nil {
		return n.store().typeAt(box.typeBeforeCast)
	}
	return nil
}

func (n Node) copyTemplateColdProps(other Node) {
	if otherBox := other.extBoxOrNil(); otherBox != nil {
		if otherBox.propListHead != nil {
			otherBox.propListHead.markShared()
		}
		if otherBox.jsDocInfo != nil && otherBox.typeBeforeCast == 0 && otherBox.propListHead == nil && otherBox.coldProps == 0 && (n.rec().kindFlags&kindFlagHasExtBox) == 0 {
			n.setAux(nodeAux{jsDoc: otherBox.jsDocInfo})
			return
		}
		if otherBox.jsDocInfo != nil || otherBox.typeBeforeCast != 0 || otherBox.propListHead != nil || otherBox.coldProps != 0 {
			box := n.ensureExtBox()
			box.jsDocInfo = otherBox.jsDocInfo
			box.typeBeforeCast = n.store().translateTypeHandle(other.store(), otherBox.typeBeforeCast)
			box.propListHead = otherBox.propListHead
			box.coldProps = otherBox.coldProps
		}
		return
	}
	if (other.props() & propFlagHasJSDoc) != 0 {
		if (n.rec().kindFlags & kindFlagHasExtBox) == 0 {
			n.setAux(nodeAux{jsDoc: other.jsDocAux()})
		} else {
			n.ensureExtBox().jsDocInfo = other.jsDocAux()
		}
	}
}

func (n Node) nodeKind() node_kind {
	return node_kind(n.rec().kindFlags & kindMask)
}

var node_kindClassNames = [...]string{
	"com.google.javascript.rhino.Node",
	"com.google.javascript.rhino.Node$StringNode",
	"com.google.javascript.rhino.Node$NumberNode",
	"com.google.javascript.rhino.Node$BigIntNode",
	"com.google.javascript.rhino.Node$TemplateLiteralSubstringNode",
}

// IsPlainNode reports whether this node is a plain Node rather than one of its
// folded payload subclasses (StringNode, NumberNode, BigIntNode,
// TemplateLiteralSubstringNode).
//
// Go-only: language-adapter — Java code compares `node.getClass() != Node.class`
// in CodeGenerator; the folded subclasses have no separate Go type.
func (n Node) IsPlainNode() bool {
	return n.nodeKind() == node_kindNode
}

// IsStringNode reports whether this node carries a string payload (StringNode).
func (n Node) IsStringNode() bool {
	return n.nodeKind() == node_kindString
}

//go:noinline
func failNodeCastKind(got, want node_kind) {
	panic(jutil.NewClassCastException(node_kindClassNames[got], node_kindClassNames[want]))
}

// node_castKind is the Java cast `((XNode) this)`: it throws a
// ClassCastException unless n is of the given kind.
func (n Node) node_castKind(k node_kind) {
	if got := n.nodeKind(); got != k {
		failNodeCastKind(got, k)
	}
}

//go:noinline
func failNodeCast(value any, zero any) {
	panic(jutil.NewClassCastException(fmt.Sprintf("%T", value), fmt.Sprintf("%T", zero)))
}

// node_cast is the Java reference cast `(T) value` of a property value:
// null stays null (the zero T), other values must be of type T.
func node_cast[T any](value any) T {
	if value == nil {
		var zero T
		return zero
	}
	t, ok := value.(T)
	if !ok {
		var zero T
		failNodeCast(value, any(zero))
	}
	return t
}

// NumberNode(double)
//
// Java: Node.NumberNode(double)
func newNode_NumberNode(number float64) Node {
	n := allocNode()
	nr := n.rec()
	nr.token = Token_NUMBER
	nr.linenoCharno = -1
	nr.kindFlags = uint8(node_kindNumber)
	n.SetDouble(number)
	return n
}

// BigIntNode(BigInteger)
//
// Java: Node.BigIntNode(BigInteger)
func newNode_BigIntNode(bigint *big.Int) Node {
	n := allocNode()
	nr := n.rec()
	nr.token = Token_BIGINT
	nr.linenoCharno = -1
	nr.kindFlags = uint8(node_kindBigInt)
	n.SetBigInt(bigint)
	return n
}

// Only for cloneNode.
//
// Java: Node.StringNode(Token)
func newNode_StringNodeForClone(token Token) Node {
	n := allocNode()
	nr := n.rec()
	nr.token = token
	nr.linenoCharno = -1
	nr.kindFlags = uint8(node_kindString)
	return n
}

// Java: Node.StringNode(Token, String)
func newNode_StringNode(token Token, str string) Node {
	n := allocNode()
	nr := n.rec()
	nr.token = token
	nr.linenoCharno = -1
	nr.kindFlags = uint8(node_kindString)
	nr.payload = uint32(strpool.Intern(str))
	return n
}

// Java: Node.TemplateLiteralSubstringNode(String, String)
func newNode_TemplateLiteralSubstringNode(cooked base.OptString, raw string) Node {
	n := allocNode()
	nr := n.rec()
	nr.token = Token_TEMPLATELIT_STRING
	nr.linenoCharno = -1
	nr.kindFlags = uint8(node_kindTemplateLiteralSubstring)
	n.initTemplateLitPayload(rhinoStringIDFromOpt(cooked), strpool.Intern(raw))
	return n
}

// Java: Node(Token)
func NewNode(token Token) Node {
	n := allocNode()
	nr := n.rec()
	nr.token = token
	nr.linenoCharno = -1
	return n
}

// AllocNode returns a fresh zeroed node (to be initialized with one of the
// Init methods).
func AllocNode() Node {
	return allocNode()
}

// Init initializes a zero Node with token and default source position.
func (n Node) Init(token Token) Node {
	nr := n.rec()
	nr.token = token
	nr.linenoCharno = -1
	return n
}

// InitFromTemplate initializes a zero Node from a slab with token and template properties.
// Template nodes carry only inline properties (SOURCE_FILE and hot boolean flags).
func (n Node) InitFromTemplate(token Token, template Node) Node {
	nr, tr := n.rec(), template.rec()
	nr.token = token
	nr.linenoCharno = -1
	nr.sourceFile = tr.sourceFile
	nr.propFlags = tr.propFlags
	return n
}

// InitWithChildFromTemplate initializes a zero Node from a slab with 1 child and template properties.
func (n Node) InitWithChildFromTemplate(token Token, child Node, template Node) Node {
	nr, tr, cr := n.rec(), template.rec(), child.rec()
	if (cr.parent|cr.next|cr.previous) != 0 ||
		tr.propFlags&propWordJSDoc != 0 ||
		tr.kindFlags >= kindFlagHasExtBox ||
		(uint32(n)^uint32(child))&staticBit != 0 {
		return n.initWithChildFromTemplateSlow(token, child, template)
	}
	nr.token = token
	nr.linenoCharno = -1
	nr.sourceFile = tr.sourceFile
	nr.propFlags = tr.propFlags
	nr.first = uint32(child)
	cr.previous = uint32(child)
	cr.parent = uint32(n)
	return n
}

//go:noinline
func (n Node) initWithChildFromTemplateSlow(token Token, child Node, template Node) Node {
	nr := n.rec()
	nr.token = token
	nr.linenoCharno = -1
	n.setPropsAndSourceFile(template)
	if template.hasAux() {
		n.copyTemplateColdProps(template)
	}
	nr.first = n.ref(child)
	child.CheckDetached()
	cr := child.rec()
	cr.previous = child.ref(child)
	cr.parent = child.ref(n)
	return n
}

// InitWith2ChildrenFromTemplate initializes a zero Node from a slab with 2 children and template properties.
func (n Node) InitWith2ChildrenFromTemplate(token Token, left Node, right Node, template Node) Node {
	nr, tr, lr := n.rec(), template.rec(), left.rec()
	if (lr.parent|lr.next|lr.previous) != 0 || right == 0 || left == right ||
		tr.propFlags&propWordJSDoc != 0 || tr.kindFlags >= kindFlagHasExtBox {
		return n.initWith2ChildrenFromTemplateSlow(token, left, right, template)
	}
	rr := right.rec()
	if (rr.parent|rr.next|rr.previous) != 0 ||
		((uint32(n)^uint32(left))|(uint32(n)^uint32(right)))&staticBit != 0 {
		return n.initWith2ChildrenFromTemplateSlow(token, left, right, template)
	}
	nr.token = token
	nr.linenoCharno = -1
	nr.sourceFile = tr.sourceFile
	nr.propFlags = tr.propFlags
	nr.first = uint32(left)
	lr.next = uint32(right)
	lr.previous = uint32(right)
	lr.parent = uint32(n)
	rr.previous = uint32(left)
	rr.parent = uint32(n)
	return n
}

//go:noinline
func (n Node) initWith2ChildrenFromTemplateSlow(token Token, left Node, right Node, template Node) Node {
	nr := n.rec()
	nr.token = token
	nr.linenoCharno = -1
	n.setPropsAndSourceFile(template)
	if template.hasAux() {
		n.copyTemplateColdProps(template)
	}
	nr.first = n.ref(left)

	left.CheckDetached()
	lr := left.rec()
	lr.next = left.ref(right)
	lr.previous = left.ref(right)
	lr.parent = left.ref(n)

	right.CheckDetached()
	rr := right.rec()
	rr.previous = right.ref(left)
	rr.parent = right.ref(n)
	return n
}

// InitWith3ChildrenFromTemplate initializes a zero Node from a slab with 3 children and template properties.
func (n Node) InitWith3ChildrenFromTemplate(token Token, left Node, mid Node, right Node, template Node) Node {
	nr, tr, lr := n.rec(), template.rec(), left.rec()
	if (lr.parent|lr.next|lr.previous) != 0 || mid == 0 || right == 0 ||
		left == mid || left == right || mid == right ||
		tr.propFlags&propWordJSDoc != 0 || tr.kindFlags >= kindFlagHasExtBox {
		return n.initWith3ChildrenFromTemplateSlow(token, left, mid, right, template)
	}
	mr, rr := mid.rec(), right.rec()
	if (mr.parent|mr.next|mr.previous|rr.parent|rr.next|rr.previous) != 0 ||
		((uint32(n)^uint32(left))|(uint32(n)^uint32(mid))|(uint32(n)^uint32(right)))&staticBit != 0 {
		return n.initWith3ChildrenFromTemplateSlow(token, left, mid, right, template)
	}
	nr.token = token
	nr.linenoCharno = -1
	nr.sourceFile = tr.sourceFile
	nr.propFlags = tr.propFlags
	nr.first = uint32(left)
	lr.next = uint32(mid)
	lr.previous = uint32(right)
	lr.parent = uint32(n)
	mr.next = uint32(right)
	mr.previous = uint32(left)
	mr.parent = uint32(n)
	rr.previous = uint32(mid)
	rr.parent = uint32(n)
	return n
}

//go:noinline
func (n Node) initWith3ChildrenFromTemplateSlow(token Token, left Node, mid Node, right Node, template Node) Node {
	nr := n.rec()
	nr.token = token
	nr.linenoCharno = -1
	n.setPropsAndSourceFile(template)
	if template.hasAux() {
		n.copyTemplateColdProps(template)
	}
	nr.first = n.ref(left)

	left.CheckDetached()
	lr := left.rec()
	lr.next = left.ref(mid)
	lr.previous = left.ref(right)
	lr.parent = left.ref(n)

	mid.CheckDetached()
	mr := mid.rec()
	mr.next = mid.ref(right)
	mr.previous = mid.ref(left)
	mr.parent = mid.ref(n)

	right.CheckDetached()
	rr := right.rec()
	rr.previous = right.ref(mid)
	rr.parent = right.ref(n)
	return n
}

// InitStringFromTemplate initializes a zero Node from a slab as a StringNode with template properties.
func (n Node) InitStringFromTemplate(token Token, str string, template Node) Node {
	nr, tr := n.rec(), template.rec()
	if tr.propFlags&propWordJSDoc != 0 || tr.kindFlags >= kindFlagHasExtBox {
		panic("rhino.Node: template node has cold properties")
	}
	nr.token = token
	nr.linenoCharno = -1
	nr.kindFlags = uint8(node_kindString)
	nr.sourceFile = tr.sourceFile
	nr.propFlags = tr.propFlags
	nr.payload = uint32(strpool.Intern(str))
	return n
}

// InitStringIDFromTemplate is InitStringFromTemplate for a string already
// interned (its pool id). Template nodes carry only inline properties.
func (n Node) InitStringIDFromTemplate(token Token, id rhinoStringID, template Node) Node {
	nr, tr := n.rec(), template.rec()
	nr.token = token
	nr.linenoCharno = -1
	nr.kindFlags = uint8(node_kindString)
	nr.sourceFile = tr.sourceFile
	nr.propFlags = tr.propFlags
	nr.payload = uint32(id)
	return n
}

// SetStringID sets the interned string payload on a freshly initialized Node.
//
// Java: Node.newString(Token, RhinoStringPool.LazyInternedStringList, int), Node.StringNode(Token, RhinoStringPool.LazyInternedStringList, int)
func (n Node) SetStringID(id rhinoStringID) Node {
	nr := n.rec()
	nr.kindFlags |= uint8(node_kindString)
	nr.payload = uint32(id)
	return n
}

// SetNumberRaw sets the float64 payload on a freshly initialized Node.
func (n Node) SetNumberRaw(number float64) Node {
	nr := n.rec()
	nr.kindFlags |= uint8(node_kindNumber)
	nr.payload = numberPayload(number)
	return n
}

// InitTemplateLitStringFromStringPool initializes a zero Node from a slab as a TemplateLiteralSubstringNode from a LazyInternedStringList.
//
// Java: Node.newTemplateLitString(RhinoStringPool.LazyInternedStringList, int, int), Node.TemplateLiteralSubstringNode(RhinoStringPool.LazyInternedStringList, int, int)
func (n Node) InitTemplateLitStringFromStringPool(
	stringPool *RhinoStringPool_LazyInternedStringList,
	cookedOffsetOrNegativeOne int,
	rawOffset int,
	template Node,
) Node {
	nr, tr := n.rec(), template.rec()
	if tr.propFlags&propWordJSDoc != 0 || tr.kindFlags >= kindFlagHasExtBox {
		panic("rhino.Node: template node has cold properties")
	}
	nr.token = Token_TEMPLATELIT_STRING
	nr.linenoCharno = -1
	nr.kindFlags = uint8(node_kindTemplateLiteralSubstring)
	nr.sourceFile = tr.sourceFile
	nr.propFlags = tr.propFlags
	var cookedID rhinoStringID
	if cookedOffsetOrNegativeOne != -1 {
		cookedID = stringPool.GetID(cookedOffsetOrNegativeOne)
	}
	rawID := stringPool.GetID(rawOffset)
	n.initTemplateLitPayload(cookedID, rawID)
	return n
}

// InitBigIntFromTemplate initializes a zero Node from a slab as a BigIntNode with template properties.
func (n Node) InitBigIntFromTemplate(bigint *big.Int, template Node) Node {
	nr, tr := n.rec(), template.rec()
	if tr.propFlags&propWordJSDoc != 0 || tr.kindFlags >= kindFlagHasExtBox {
		panic("rhino.Node: template node has cold properties")
	}
	nr.token = Token_BIGINT
	nr.linenoCharno = -1
	nr.kindFlags = uint8(node_kindBigInt)
	nr.sourceFile = tr.sourceFile
	nr.propFlags = tr.propFlags
	n.SetBigInt(bigint)
	return n
}

// InitTemplateLitStringFromTemplate initializes a zero Node from a slab as a TemplateLiteralSubstringNode with template properties.
func (n Node) InitTemplateLitStringFromTemplate(cooked base.OptString, raw string, template Node) Node {
	nr, tr := n.rec(), template.rec()
	if tr.propFlags&propWordJSDoc != 0 || tr.kindFlags >= kindFlagHasExtBox {
		panic("rhino.Node: template node has cold properties")
	}
	nr.token = Token_TEMPLATELIT_STRING
	nr.linenoCharno = -1
	nr.kindFlags = uint8(node_kindTemplateLiteralSubstring)
	nr.sourceFile = tr.sourceFile
	nr.propFlags = tr.propFlags
	n.initTemplateLitPayload(rhinoStringIDFromOpt(cooked), strpool.Intern(raw))
	return n
}

// Java: Node(Token, Node)
func NewNodeWithChild(token Token, child Node) Node {
	n := allocNode()
	nr, cr := n.rec(), child.rec()
	if (cr.parent|cr.next|cr.previous) != 0 || (uint32(n)^uint32(child))&staticBit != 0 {
		return n.initNewNodeWithChildSlow(token, child)
	}
	nr.token = token
	nr.linenoCharno = -1
	nr.first = uint32(child)
	cr.previous = uint32(child)
	cr.parent = uint32(n)
	return n
}

//go:noinline
func (n Node) initNewNodeWithChildSlow(token Token, child Node) Node {
	nr := n.rec()
	nr.token = token
	nr.linenoCharno = -1
	nr.first = n.ref(child)
	child.CheckDetached()
	cr := child.rec()
	cr.previous = child.ref(child)
	cr.parent = child.ref(n)
	return n
}

// Java: Node(Token, Node, Node)
func NewNodeWith2Children(token Token, left Node, right Node) Node {
	n := allocNode()
	nr, lr := n.rec(), left.rec()
	if (lr.parent|lr.next|lr.previous) != 0 || right == 0 || left == right {
		return n.initNewNodeWith2ChildrenSlow(token, left, right)
	}
	rr := right.rec()
	if (rr.parent|rr.next|rr.previous) != 0 ||
		((uint32(n)^uint32(left))|(uint32(n)^uint32(right)))&staticBit != 0 {
		return n.initNewNodeWith2ChildrenSlow(token, left, right)
	}
	nr.token = token
	nr.linenoCharno = -1
	nr.first = uint32(left)
	lr.next = uint32(right)
	lr.previous = uint32(right)
	lr.parent = uint32(n)
	rr.previous = uint32(left)
	rr.parent = uint32(n)
	return n
}

//go:noinline
func (n Node) initNewNodeWith2ChildrenSlow(token Token, left Node, right Node) Node {
	nr := n.rec()
	nr.token = token
	nr.linenoCharno = -1
	nr.first = n.ref(left)

	left.CheckDetached()
	lr := left.rec()
	lr.next = left.ref(right)
	lr.previous = left.ref(right)
	lr.parent = left.ref(n)

	right.CheckDetached()
	rr := right.rec()
	rr.previous = right.ref(left)
	rr.parent = right.ref(n)
	return n
}

// Java: Node(Token, Node, Node, Node)
func NewNodeWith3Children(token Token, left Node, mid Node, right Node) Node {
	n := allocNode()
	nr, lr := n.rec(), left.rec()
	if (lr.parent|lr.next|lr.previous) != 0 || mid == 0 || right == 0 ||
		left == mid || left == right || mid == right {
		return n.initNewNodeWith3ChildrenSlow(token, left, mid, right)
	}
	mr, rr := mid.rec(), right.rec()
	if (mr.parent|mr.next|mr.previous|rr.parent|rr.next|rr.previous) != 0 ||
		((uint32(n)^uint32(left))|(uint32(n)^uint32(mid))|(uint32(n)^uint32(right)))&staticBit != 0 {
		return n.initNewNodeWith3ChildrenSlow(token, left, mid, right)
	}
	nr.token = token
	nr.linenoCharno = -1
	nr.first = uint32(left)
	lr.next = uint32(mid)
	lr.previous = uint32(right)
	lr.parent = uint32(n)
	mr.next = uint32(right)
	mr.previous = uint32(left)
	mr.parent = uint32(n)
	rr.previous = uint32(mid)
	rr.parent = uint32(n)
	return n
}

//go:noinline
func (n Node) initNewNodeWith3ChildrenSlow(token Token, left Node, mid Node, right Node) Node {
	nr := n.rec()
	nr.token = token
	nr.linenoCharno = -1
	nr.first = n.ref(left)

	left.CheckDetached()
	lr := left.rec()
	lr.next = left.ref(mid)
	lr.previous = left.ref(right)
	lr.parent = left.ref(n)

	mid.CheckDetached()
	mr := mid.rec()
	mr.next = mid.ref(right)
	mr.previous = mid.ref(left)
	mr.parent = mid.ref(n)

	right.CheckDetached()
	rr := right.rec()
	rr.previous = right.ref(mid)
	rr.parent = right.ref(n)
	return n
}

// Java: Node.newNumber(double)
func Node_NewNumber(number float64) Node {
	return newNode_NumberNode(number)
}

// Java: Node.newBigInt(BigInteger)
func Node_NewBigInt(bigint *big.Int) Node {
	return newNode_BigIntNode(bigint)
}

// Java: Node.newString(String)
func Node_NewString(str string) Node {
	return newNode_StringNode(Token_STRINGLIT, str)
}

// Java: Node.newString(Token, String)
func Node_NewStringWithToken(token Token, str string) Node {
	return newNode_StringNode(token, str)
}

// Java: Node.newTemplateLitString(String, String)
func Node_NewTemplateLitString(cooked base.OptString, raw string) Node {
	return newNode_TemplateLiteralSubstringNode(cooked, raw)
}

// Java: Node.getToken()
func (n Node) GetToken() Token {
	return n.rec().token
}

// Java: Node.setToken(Token)
func (n Node) SetToken(token Token) {
	nr := n.rec()
	nr.token = token
	nr.kindFlags &^= kindFlagNoBlockDecl
	if nr.parent != 0 {
		Node(nr.parent).rec().kindFlags &^= kindFlagNoBlockDecl
	}
}

// Java: Node.setClosureUnawareShadow(Node)
func (n Node) SetClosureUnawareShadow(shadowRoot Node) {
	nr := n.rec()
	jutil.CheckStateMsg(nr.first == 0, "Cannot set shadow root on a node with children")
	jutil.CheckStateMsg(nr.token == Token_NAME, "Only NAME nodes can be used as shadows")
	n.PutProp(Node_Prop_CLOSURE_UNAWARE_SHADOW, shadowRoot)
}

// Java: Node.getClosureUnawareShadow()
func (n Node) GetClosureUnawareShadow() Node {
	nr := n.rec()
	if nr.token != Token_NAME || (nr.kindFlags&kindFlagHasExtBox) == 0 {
		// As checked in setClosureUnawareShadow, only NAME nodes can have shadow
		// content. We short-circuit the getProp call here because that can be
		// very expensive for nodes that have long proplists (and we know it
		// wouldn't find this prop anyways).
		return NilNode
	}
	return node_cast[Node](n.GetProp(Node_Prop_CLOSURE_UNAWARE_SHADOW))
}

// Java: Node.setIsInClosureUnawareSubtree(boolean)
func (n Node) SetIsInClosureUnawareSubtree(value bool) {
	nr := n.rec()
	hasOnlySourceFile := n.props() == 0 &&
		n.typeBeforeCastBase() == nil &&
		(node_kind(nr.kindFlags&kindMask) != node_kindNode || nr.payload == 0) &&
		nr.sourceFile != 0 && n.coldPropListHead() == nil
	if !hasOnlySourceFile {
		panic(jutil.NewIllegalArgumentExceptionMsg(
			jutil.StringFormat(
				"Cannot set IS_IN_CLOSURE_UNAWARE_SUBTREE property on a node which proplist is not a"+
					" single SOURCE_FILE prop: %s %s",
				n, n.GetPropListDebugString())))
	}
	n.PutBooleanProp(Node_Prop_IS_IN_CLOSURE_UNAWARE_SUBTREE, value)
}

// Java: Node.getIsInClosureUnawareSubtree()
func (n Node) GetIsInClosureUnawareSubtree() bool {
	return n.rec().propFlags&propWord(propFlag_IS_IN_CLOSURE_UNAWARE_SUBTREE) != 0
}

// nodeData is the record of an AST node (fields as declared in Java, plus
// the folded subclass payloads); a Node is a handle to one (node_store.go).
//
// Go-only: performance. 36 bytes, no Go pointers, 4-byte aligned: the
// tree links are node ids; the leaf payload is one 32-bit value whose
// meaning depends on the kind (a string id, a number, a change time...;
// node_payload.go); the source file is a 32-bit id (allocSourceFileID); the
// 25 common boolean properties, JSDoc bit, CONSTANT_VAR_FLAGS and
// SIDE_EFFECT_FLAGS share the 32-bit propFlags word directly, and the 15
// ultra-rare boolean properties live in nodeExtBox.coldProps.
// What most nodes do not have is kept beside the record (node_tables.go):
// the JSType or Color handle and the original name in per-block columns,
// the cold data (JSDocInfo, remaining properties, overflow length, a cooked
// template string that differs from the raw one: about 1% of nodes) in a
// map of the store, looked up only when the record's flags say the node has
// some.
type nodeData struct {
	next     nodeID // next sibling, a linked list (0 = none)
	first    nodeID // first element of a linked list of children
	parent   nodeID // parent node
	previous nodeID // previous sibling, a circular linked list

	payload      uint32   // one value whose meaning depends on the kind (node_payload.go)
	sourceFile   uint32   // the SOURCE_FILE property: an id from allocSourceFileID (0 = none)
	propFlags    propWord // the property bits (Node.props): bool props, JSDoc, ConstVar, SideEffect
	linenoCharno int32    // line (upper 20 bits) and column (low 12 bits)
	length16     uint16   // code length (0..65534 inline; 65535: in the nodeExtBox)
	token        Token    // token kind (uint8)
	kindFlags    uint8    // bits 0..2: node_kind; bit 3: typeIsColor; bit 4: noBlockDecls; bit 5: hasExtra; bit 6: origNameIsString; bit 7: hasExtBox
}

func (n Node) HasNoBlockScopeDecls() bool {
	return n.rec().kindFlags&kindFlagNoBlockDecl != 0
}

func (n Node) SetNoBlockScopeDecls() {
	n.rec().kindFlags |= kindFlagNoBlockDecl
}

// Preconditions.checkNotNull for a Node. The generic jutil versions take an
// any, which would put the 16-byte handle on the heap on every call; these
// do the same check without allocating.

// CheckNotNull is Preconditions.checkNotNull(T) for a Node.
func CheckNotNull(n Node) Node {
	if n == 0 {
		panic(jutil.NewNullPointerException())
	}
	return n
}

// CheckNotNullMsg is Preconditions.checkNotNull(T, String, Object...) for a Node.
func CheckNotNullMsg(n Node, template string, args ...any) Node {
	if n == 0 {
		failNotNull(template, args...)
	}
	return n
}

// CheckNotNullMsg1 is Preconditions.checkNotNull(T, String, Object) for a Node.
func CheckNotNullMsg1[A any](n Node, template string, p1 A) Node {
	if n == 0 {
		failNotNull(template, p1)
	}
	return n
}

// CheckNotNullMsg2 is Preconditions.checkNotNull(T, String, Object, Object) for a Node.
func CheckNotNullMsg2[A, B any](n Node, template string, p1 A, p2 B) Node {
	if n == 0 {
		failNotNull(template, p1, p2)
	}
	return n
}

//go:noinline
func failNotNull(template string, args ...any) {
	panic(jutil.NewNullPointerExceptionMsg(jutil.LenientFormat(template, args...)))
}
