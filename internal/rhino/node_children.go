// Port of com/google/javascript/rhino/Node.java @ closure-compiler v20261021

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
	"iter"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// Java: Node.hasChildren()
func (n Node) HasChildren() bool {
	return n.rec().first != 0
}

// Java: Node.getOnlyChild()
func (n Node) GetOnlyChild() Node {
	jutil.CheckState(n.HasOneChild())
	return n.firstNode()
}

// Java: Node.getFirstChild()
func (n Node) GetFirstChild() Node {
	return n.firstNode()
}

// GetFirstFirstChild gets the first child of the first child. This method
// assumes that the first child exists.
//
// Java: Node.getFirstFirstChild()
func (n Node) GetFirstFirstChild() Node {
	return n.firstNode().firstNode()
}

// Java: Node.getSecondChild()
func (n Node) GetSecondChild() Node {
	return n.firstNode().nextNode()
}

// Java: Node.getLastChild()
func (n Node) GetLastChild() Node {
	// The first child's previous is the last child (a circular list).
	if first := n.firstNode(); first != NilNode {
		return first.prevNode()
	}
	return NilNode
}

// Java: Node.getNext()
func (n Node) GetNext() Node {
	return n.nextNode()
}

// Java: Node.getPrevious()
func (n Node) GetPrevious() Node {
	if uint32(n) == n.parentNode().rec().first {
		return NilNode
	}
	return n.prevNode()
}

// GetChildAtIndex gets the ith child, note that this is O(N) where N is the
// number of children.
//
// Java: Node.getChildAtIndex(int)
func (n Node) GetChildAtIndex(i int) Node {
	c := n.firstNode()
	for i > 0 {
		c = c.nextNode()
		i--
	}
	return c
}

// GetIndexOfChild gets the index of a child, note that this is O(N) where N
// is the number of children.
//
// Java: Node.getIndexOfChild(Node)
func (n Node) GetIndexOfChild(child Node) int {
	c := n.firstNode()
	i := 0
	for !c.IsNil() {
		if child == c {
			return i
		}

		c = c.nextNode()
		i++
	}
	return -1
}

// Java: Node.addChildToFront(Node)
func (n Node) AddChildToFront(child Node) {
	jutil.CheckArgument(child.rec().parent == 0)
	jutil.CheckArgument(child.rec().next == 0)
	jutil.CheckArgument(child.rec().previous == 0)
	n.rec().kindFlags &^= kindFlagNoBlockDecl
	child.rec().parent = child.ref(n)
	child.rec().next = child.ref(n.firstNode())
	if n.rec().first == 0 {
		// NOTE: child.next remains null
		child.rec().previous = child.ref(child)
	} else {
		last := n.firstNode().prevNode()
		// NOTE: last.next remains null
		child.rec().previous = child.ref(last)
		child.rec().next = child.ref(n.firstNode())
		n.firstNode().rec().previous = n.firstNode().ref(child)
	}
	n.rec().first = n.ref(child)
}

// Java: Node.addChildToBack(Node)
func (n Node) AddChildToBack(child Node) {
	jutil.CheckArgumentMsg3(
		child.rec().parent == 0,
		"Cannot add already-owned child node.\nChild: %s\nExisting parent: %s\nNew parent: %s",
		child,
		child.parentNode(),
		n)

	jutil.CheckArgument(child.rec().next == 0)
	jutil.CheckArgument(child.rec().previous == 0)
	n.rec().kindFlags &^= kindFlagNoBlockDecl

	if n.rec().first == 0 {
		// NOTE: child.next remains null
		child.rec().previous = child.ref(child)
		n.rec().first = n.ref(child)
	} else {
		last := n.firstNode().prevNode()
		last.rec().next = last.ref(child)
		// NOTE: child.next remains null
		child.rec().previous = child.ref(last)
		n.firstNode().rec().previous = n.firstNode().ref(child)
	}

	child.rec().parent = child.ref(n)
}

// AddChildrenToFront adds all children to the front of this node.
//
// children is the first of a list of sibling nodes who have no parent.
// NOTE: Usually you would get this argument from a removeChildren() call. A
// single detached node will not work because its sibling pointers will not
// be correctly initialized.
//
// Java: Node.addChildrenToFront(Node)
func (n Node) AddChildrenToFront(children Node) {
	if children.IsNil() {
		return // removeChildren() returns null when there are none
	}
	// NOTE: If there is only one sibling, its previous pointer must point to
	// itself. Null indicates a fully detached node.
	CheckNotNullMsg1(children.prevNode(), "%s", children)
	n.rec().kindFlags &^= kindFlagNoBlockDecl
	for child := children; !child.IsNil(); child = child.nextNode() {
		jutil.CheckArgument(child.rec().parent == 0)
		child.rec().parent = child.ref(n)
	}

	lastSib := children.prevNode()
	if n.rec().first != 0 {
		last := n.firstNode().prevNode()
		// NOTE: last.next remains null
		children.rec().previous = children.ref(last)
		lastSib.rec().next = lastSib.ref(n.firstNode())
		n.firstNode().rec().previous = n.firstNode().ref(lastSib)
	}
	n.rec().first = n.ref(children)
}

// Java: Node.addChildrenToBack(Node)
func (n Node) AddChildrenToBack(children Node) {
	n.rec().kindFlags &^= kindFlagNoBlockDecl
	n.AddChildrenAfter(children, n.GetLastChild())
}

// Java: Node.insertAfter(Node)
func (n Node) InsertAfter(existing Node) {
	existing.CheckAttached()
	n.CheckDetached()

	existingParent := existing.parentNode()
	existingNext := existing.nextNode()
	existingParent.rec().kindFlags &^= kindFlagNoBlockDecl

	n.rec().parent = n.ref(existingParent)

	existing.rec().next = existing.ref(n)
	n.rec().previous = n.ref(existing)

	if existingNext.IsNil() {
		existingParent.firstNode().rec().previous = existingParent.firstNode().ref(n)
		// this.next remains null
	} else {
		existingNext.rec().previous = existingNext.ref(n)
		n.rec().next = n.ref(existingNext)
	}
}

// Java: Node.insertBefore(Node)
func (n Node) InsertBefore(existing Node) {
	existing.CheckAttached()
	n.CheckDetached()

	existingParent := existing.parentNode()
	existingPrevious := existing.prevNode()
	existingParent.rec().kindFlags &^= kindFlagNoBlockDecl

	n.rec().parent = n.ref(existingParent)

	n.rec().next = n.ref(existing)
	existing.rec().previous = existing.ref(n)

	n.rec().previous = n.ref(existingPrevious)
	if existingPrevious.rec().next == 0 {
		existingParent.rec().first = existingParent.ref(n)
		// existingPrevious.next remains null
	} else {
		// existingParent.first remains existing
		existingPrevious.rec().next = existingPrevious.ref(n)
	}
}

// AddChildrenAfter adds all children after 'node'. If 'node' is null, add
// them to the front of this node.
//
// children is the first of a list of sibling nodes who have no parent.
// NOTE: Usually you would get this argument from a removeChildren() call. A
// single detached node will not work because its sibling pointers will not
// be correctly initialized.
//
// Java: Node.addChildrenAfter(Node, Node)
func (n Node) AddChildrenAfter(children Node, node Node) {
	if children.IsNil() {
		return // removeChildren() returns null when there are none
	}
	jutil.CheckArgument(node.IsNil() || node.parentNode() == n)
	// NOTE: If there is only one sibling, its previous pointer must point to
	// itself. Null indicates a fully detached node.
	CheckNotNullMsg1(children.prevNode(), "%s", children)
	if node.IsNil() {
		n.AddChildrenToFront(children)
		return
	}
	n.rec().kindFlags &^= kindFlagNoBlockDecl

	for child := children; !child.IsNil(); child = child.nextNode() {
		jutil.CheckArgument(child.rec().parent == 0)
		child.rec().parent = child.ref(n)
	}

	lastSibling := children.prevNode()
	nodeAfter := node.nextNode()
	lastSibling.rec().next = lastSibling.ref(nodeAfter)
	if nodeAfter.IsNil() {
		n.firstNode().rec().previous = n.firstNode().ref(lastSibling)
	} else {
		nodeAfter.rec().previous = nodeAfter.ref(lastSibling)
	}
	node.rec().next = node.ref(children)
	children.rec().previous = children.ref(node)
}

// ReplaceWith swaps `replacement` and its subtree into the position of
// `this`.
//
// Java: Node.replaceWith(Node)
func (n Node) ReplaceWith(replacement Node) {
	n.CheckAttached()
	replacement.CheckDetached()

	existingParent := n.parentNode()
	existingNext := n.nextNode()
	existingPrevious := n.prevNode()
	existingParent.rec().kindFlags &^= kindFlagNoBlockDecl

	// Copy over important information.
	// TODO(nickreid): Stop doing this. It's totally unexpected.
	replacement.SrcrefIfMissing(n)

	// The sequence below also has to work when `this` is an only child, ehich
	// can cause many of the variables to point to the same object.

	n.rec().parent = 0
	replacement.rec().parent = replacement.ref(existingParent)

	n.rec().previous = 0
	replacement.rec().previous = replacement.ref(existingPrevious)
	if existingPrevious.rec().next == 0 {
		existingParent.rec().first = existingParent.ref(replacement)
		// existingPrevious.next remains null
	} else {
		// existingParent.first is unchanged;
		existingPrevious.rec().next = existingPrevious.ref(replacement)
	}

	if existingNext.IsNil() {
		// this.next remains null;
		existingParent.firstNode().rec().previous = existingParent.firstNode().ref(replacement)
		// replacement.next remains null
	} else {
		n.rec().next = 0
		existingNext.rec().previous = existingNext.ref(replacement)
		replacement.rec().next = replacement.ref(existingNext)
	}
}

// Detach removes this node from its parent, but retains its subtree.
//
// Java: Node.detach()
func (n Node) Detach() Node {
	n.CheckAttached()

	existingParent := n.parentNode()
	existingNext := n.nextNode()
	existingPrevious := n.prevNode()
	existingParent.rec().kindFlags &^= kindFlagNoBlockDecl

	// The sequence below also has to work when `this` is an only child or has
	// a single sibling, which can cause many of the variables to point to the
	// same object.

	n.rec().parent = 0

	if existingNext.IsNil() {
		// this.next remains null;
		existingParent.firstNode().rec().previous = existingParent.firstNode().ref(existingPrevious)
	} else {
		n.rec().next = 0
		existingNext.rec().previous = existingNext.ref(existingPrevious)
	}

	n.rec().previous = 0
	if existingPrevious.rec().next == 0 {
		existingParent.rec().first = existingParent.ref(existingNext)
		// existingPrevious.next remains null
	} else {
		// existingParent.first is unchanged;
		existingPrevious.rec().next = existingPrevious.ref(existingNext)
	}

	return n
}

// Java: Node.CheckAttached()
func (n Node) CheckAttached() {
	jutil.CheckStateMsg1(n.rec().parent != 0, "Has no parent: %s", n)
}

// Java: Node.CheckDetached()
func (n Node) CheckDetached() {
	jutil.CheckStateMsg1(n.rec().parent == 0, "Has parent: %s", n)
	jutil.CheckStateMsg1(n.rec().next == 0, "Has next: %s", n)
	jutil.CheckStateMsg1(n.rec().previous == 0, "Has previous: %s", n)
}

// RemoveFirstChild removes the first child of Node. Equivalent to:
// node.removeChild(node.getFirstChild());
//
// Java: Node.removeFirstChild()
func (n Node) RemoveFirstChild() Node {
	child := n.firstNode()
	if !child.IsNil() {
		child.Detach()
	}
	return child
}

// RemoveChildren removes all children, but leaves them linked to each other.
// Returns the first child node.
//
// Java: Node.removeChildren()
func (n Node) RemoveChildren() Node {
	n.rec().kindFlags &^= kindFlagNoBlockDecl
	children := n.firstNode()
	for child := n.firstNode(); !child.IsNil(); child = child.nextNode() {
		child.rec().parent = 0
	}
	n.rec().first = 0
	return children
}

// DetachChildren removes all children from this node and isolates the
// children from each other.
//
// Java: Node.detachChildren()
func (n Node) DetachChildren() {
	n.rec().kindFlags &^= kindFlagNoBlockDecl
	for child := n.firstNode(); !child.IsNil(); {
		nextChild := child.nextNode()
		child.rec().parent = 0
		child.rec().next = 0
		child.rec().previous = 0
		child = nextChild
	}
	n.rec().first = 0
}

// ==========================================================================
// Iteration

// Children returns an iterable object that iterates over this node's
// children. The iterator does not support the optional operation
// Iterator#remove().
//
// Deprecated: see above.
//
// Java: Node.children()
func (n Node) Children() iter.Seq[Node] {
	start := n.firstNode()
	return func(yield func(Node) bool) {
		for current := start; !current.IsNil(); {
			c := current
			current = current.GetNext()
			if !yield(c) {
				return
			}
		}
	}
}

// Java: Node.getParent()
func (n Node) GetParent() Node {
	return n.parentNode()
}

// Java: Node.hasParent()
func (n Node) HasParent() bool {
	return n.rec().parent != 0
}

// Java: Node.getGrandparent()
func (n Node) GetGrandparent() Node {
	if n.rec().parent == 0 {
		return NilNode
	}
	return n.parentNode().parentNode()
}

// GetAncestor gets the ancestor node relative to this.
//
// level: 0 = this, 1 = the parent, etc.
//
// Java: Node.getAncestor(int)
func (n Node) GetAncestor(level int) Node {
	jutil.CheckArgument(level >= 0)
	node := n
	for !node.IsNil() {
		cont := level > 0
		level--
		if !cont {
			break
		}
		node = node.GetParent()
	}
	return node
}

// IsDescendantOf reports whether this Node is the same as `node` or a
// descendant of `node`.
//
// Java: Node.isDescendantOf(Node)
func (n Node) IsDescendantOf(node Node) bool {
	for c := n; !c.IsNil(); c = c.parentNode() {
		if c == node {
			return true
		}
	}
	return false
}

// Java: Node.isOnlyChildOf(Node)
func (n Node) IsOnlyChildOf(possibleParent Node) bool {
	return possibleParent == n.GetParent() && n.GetPrevious().IsNil() && n.GetNext().IsNil()
}

// Java: Node.isFirstChildOf(Node)
func (n Node) IsFirstChildOf(possibleParent Node) bool {
	return possibleParent == n.GetParent() && n.GetPrevious().IsNil()
}

// Java: Node.isSecondChildOf(Node)
func (n Node) IsSecondChildOf(possibleParent Node) bool {
	previousNode := n.GetPrevious()

	return !previousNode.IsNil() && previousNode.IsFirstChildOf(possibleParent)
}

// GetAncestors iterates all of the node's ancestors excluding itself.
//
// Java: Node.getAncestors()
func (n Node) GetAncestors() *Node_AncestorIterable {
	return NewNode_AncestorIterable(n.GetParent())
}

// Node_AncestorIterable is an iterator to go up the ancestor tree.
type Node_AncestorIterable struct {
	cur Node
}

// NewNode_AncestorIterable starts at `cur`.
//
// Java: Node.AncestorIterable(Node)
func NewNode_AncestorIterable(cur Node) *Node_AncestorIterable {
	return &Node_AncestorIterable{cur: cur}
}

// All is the iterable's iterator(). Like the Java iterator, it advances the
// iterable's own cursor, so an AncestorIterable can be traversed only once.
//
// Java: Node.AncestorIterable.iterator()
func (a *Node_AncestorIterable) All() iter.Seq[Node] {
	return func(yield func(Node) bool) {
		for !a.cur.IsNil() {
			c := a.cur
			a.cur = a.cur.GetParent()
			if !yield(c) {
				return
			}
		}
	}
}

// HasOneChild checks for one child more efficiently than by iterating over
// all the children as is done with Node.getChildCount().
//
// Java: Node.hasOneChild()
func (n Node) HasOneChild() bool {
	return n.rec().first != 0 && n.firstNode().rec().next == 0
}

// HasTwoChildren checks for two children more efficiently than
// `getChildCount() == 2`.
//
// Java: Node.hasTwoChildren()
func (n Node) HasTwoChildren() bool {
	return n.rec().first != 0 && n.firstNode().rec().next != 0 && n.firstNode().nextNode() == n.GetLastChild()
}

// HasZeroOrOneChild checks for zero or one child more efficiently than by
// iterating over all the children as is done with Node.getChildCount().
//
// Java: Node.hasZeroOrOneChild()
func (n Node) HasZeroOrOneChild() bool {
	return n.firstNode() == n.GetLastChild()
}

// HasMoreThanOneChild checks for more than one child more efficiently than
// by iterating over all the children as is done with Node.getChildCount().
//
// Java: Node.hasMoreThanOneChild()
func (n Node) HasMoreThanOneChild() bool {
	return n.rec().first != 0 && n.firstNode().rec().next != 0
}

// HasXChildren checks for has exactly the number of specified children.
//
// Java: Node.hasXChildren(int)
func (n Node) HasXChildren(x int) bool {
	c := 0
	for c2 := n.firstNode(); !c2.IsNil() && c <= x; c2 = c2.nextNode() {
		c++
	}
	return c == x
}

// Java: Node.getChildCount()
func (n Node) GetChildCount() int {
	c := 0
	for c2 := n.firstNode(); !c2.IsNil(); c2 = c2.nextNode() {
		c++
	}
	return c
}

// HasChild is intended for testing and verification only.
//
// Java: Node.hasChild(Node)
func (n Node) HasChild(child Node) bool {
	for c := n.firstNode(); !c.IsNil(); c = c.nextNode() {
		if child == c {
			return true
		}
	}
	return false
}
