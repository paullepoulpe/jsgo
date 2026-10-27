// Port of com/google/javascript/rhino/JSDocInfo.java @ closure-compiler v20261021

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
 *   Bob Jervis
 *   Google Inc.
 *
 * This Go translation is distributed under the MPL 1.1 only, as the
 * original's dual licence allows; the other licence it offers is not
 * used.
 *
 * ***** END LICENSE BLOCK ***** */

package rhino

import (
	"iter"
	"math/bits"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// JSDocInfo_ComparableSourcePosition is a SourcePosition that can be
// compared for equivalence.
//
// Java: JSDocInfo.ComparableSourcePosition<T>
type JSDocInfo_ComparableSourcePosition[T any] interface {
	SourcePosition[T]
	IsEquivalentTo(that SourcePosition[T]) bool
}

// JSDocInfo_StringPosition is a piece of information (found in a marker)
// which contains a position with a string.
//
// Go-only: language-adapter — StringPosition has a subclass
// (TrimmedStringPosition), so it is an interface plus
// JSDocInfo_StringPositionBase; the item type String is Go string (a
// position's item is always set right after construction).
type JSDocInfo_StringPosition interface {
	JSDocInfo_ComparableSourcePosition[string]
}

// JSDocInfo_StringPositionBase is `new StringPosition()`.
type JSDocInfo_StringPositionBase struct {
	SourcePositionBase[string]
}

// NewJSDocInfo_StringPosition is the implicit StringPosition constructor.
//
// Java: JSDocInfo.StringPosition()
func NewJSDocInfo_StringPosition() *JSDocInfo_StringPositionBase {
	return &JSDocInfo_StringPositionBase{}
}

// Java: JSDocInfo.StringPosition.isEquivalentTo(SourcePosition<String>)
func (p *JSDocInfo_StringPositionBase) IsEquivalentTo(that SourcePosition[string]) bool {
	return !jutil.IsNil(that) && p.IsSamePositionAs(that) && jutil.ObjectsEquals(p.GetItem(), that.GetItem())
}

// JSDocInfo_TrimmedStringPosition is a piece of information (found in a
// marker) which contains a position with a string that has no leading or
// trailing whitespace.
type JSDocInfo_TrimmedStringPosition struct {
	JSDocInfo_StringPositionBase
}

// Java: JSDocInfo.TrimmedStringPosition()
func NewJSDocInfo_TrimmedStringPosition() *JSDocInfo_TrimmedStringPosition {
	return &JSDocInfo_TrimmedStringPosition{}
}

// Java: JSDocInfo.TrimmedStringPosition.setItem(String)
func (p *JSDocInfo_TrimmedStringPosition) SetItem(item string) {
	jutil.CheckArgumentMsg(
		jutil.StringCharAt(item, 0) != ' ' && jutil.StringCharAt(item, jutil.StringLength(item)-1) != ' ',
		"String has leading or trailing whitespace")
	p.JSDocInfo_StringPositionBase.SetItem(item)
}

// JSDocInfo_NamePosition is a piece of information (found in a marker) which
// contains a position with a name node.
type JSDocInfo_NamePosition struct {
	SourcePositionBase[Node]
}

// Java: JSDocInfo.NamePosition()
func NewJSDocInfo_NamePosition() *JSDocInfo_NamePosition {
	return &JSDocInfo_NamePosition{}
}

// Java: JSDocInfo.NamePosition.isEquivalentTo(SourcePosition<Node>)
func (p *JSDocInfo_NamePosition) IsEquivalentTo(that SourcePosition[Node]) bool {
	if jutil.IsNil(that) ||
		!p.IsSamePositionAs(that) ||
		(p.GetItem().IsNil()) != (that.GetItem().IsNil()) {
		return false
	}
	return p.GetItem().IsNil() || p.GetItem().IsEquivalentTo(that.GetItem())
}

// JSDocInfo_TypePosition is a piece of information (found in a marker) which
// contains a position with a type expression syntax tree.
type JSDocInfo_TypePosition struct {
	SourcePositionBase[Node]
	brackets bool
}

// Java: JSDocInfo.TypePosition()
func NewJSDocInfo_TypePosition() *JSDocInfo_TypePosition {
	return &JSDocInfo_TypePosition{brackets: false}
}

// HasBrackets returns whether the type has curly braces around it.
//
// Java: JSDocInfo.TypePosition.hasBrackets()
func (p *JSDocInfo_TypePosition) HasBrackets() bool {
	return p.brackets
}

// Java: JSDocInfo.TypePosition.setHasBrackets(boolean)
func (p *JSDocInfo_TypePosition) SetHasBrackets(newVal bool) {
	p.brackets = newVal
}

// Java: JSDocInfo.TypePosition.isEquivalentTo(SourcePosition<Node>)
func (p *JSDocInfo_TypePosition) IsEquivalentTo(that SourcePosition[Node]) bool {
	typePosition, ok := that.(*JSDocInfo_TypePosition)
	if !ok || typePosition == nil ||
		!p.IsSamePositionAs(that) ||
		p.brackets != typePosition.brackets ||
		(p.GetItem().IsNil()) != (that.GetItem().IsNil()) {
		return false
	}
	return p.GetItem().IsNil() || p.GetItem().IsEquivalentTo(that.GetItem())
}

// JSDocInfo_Marker defines a class for containing the parsing information
// for this JSDocInfo. For each annotation found in the JsDoc, a marker will
// be created indicating the annotation itself, the name of the annotation
// (if any; for example, a @param has a name, but a @return does not), the
// textual description found on that annotation and, if applicable, the type
// declaration. All this information is only collected if documentation
// collection is turned on.
type JSDocInfo_Marker struct {
	annotation  *JSDocInfo_TrimmedStringPosition
	nameNode    *JSDocInfo_NamePosition
	description JSDocInfo_StringPosition
	type_       *JSDocInfo_TypePosition
}

// Java: JSDocInfo.Marker()
func NewJSDocInfo_Marker() *JSDocInfo_Marker { return &JSDocInfo_Marker{} }

// GetAnnotation gets the position information for the annotation name.
// (e.g., "param")
//
// Java: JSDocInfo.Marker.getAnnotation()
func (m *JSDocInfo_Marker) GetAnnotation() JSDocInfo_StringPosition {
	if m.annotation == nil {
		return nil
	}
	return m.annotation
}

// Java: JSDocInfo.Marker.setAnnotation(TrimmedStringPosition)
func (m *JSDocInfo_Marker) SetAnnotation(p *JSDocInfo_TrimmedStringPosition) {
	m.annotation = p
}

// GetNameNode gets the position information for the name found in a @param
// tag.
//
// Java: JSDocInfo.Marker.getNameNode()
func (m *JSDocInfo_Marker) GetNameNode() *JSDocInfo_NamePosition {
	return m.nameNode
}

// Java: JSDocInfo.Marker.setNameNode(NamePosition)
func (m *JSDocInfo_Marker) SetNameNode(p *JSDocInfo_NamePosition) {
	m.nameNode = p
}

// GetDescription gets the position information for the description found in
// a block tag.
//
// Java: JSDocInfo.Marker.getDescription()
func (m *JSDocInfo_Marker) GetDescription() JSDocInfo_StringPosition {
	return m.description
}

// Java: JSDocInfo.Marker.setDescription(StringPosition)
func (m *JSDocInfo_Marker) SetDescription(p JSDocInfo_StringPosition) {
	if jutil.IsNil(p) {
		m.description = nil
	} else {
		m.description = p
	}
}

// GetType gets the position information for the type expression found in
// some block tags, like "@param" and "@return".
//
// Java: JSDocInfo.Marker.getType()
func (m *JSDocInfo_Marker) GetType() *JSDocInfo_TypePosition {
	return m.type_
}

// Java: JSDocInfo.Marker.setType(TypePosition)
func (m *JSDocInfo_Marker) SetType(p *JSDocInfo_TypePosition) {
	m.type_ = p
}

// Java: JSDocInfo.Marker.isEquivalentTo(Marker)
func (m *JSDocInfo_Marker) IsEquivalentTo(that *JSDocInfo_Marker) bool {
	return JSDocInfo_Marker_AreEquivalent[string](m.annotation, that.annotation) &&
		JSDocInfo_Marker_AreEquivalent[Node](m.nameNode, that.nameNode) &&
		JSDocInfo_Marker_AreEquivalent[string](m.description, that.description) &&
		JSDocInfo_Marker_AreEquivalent[Node](m.type_, that.type_)
}

// JSDocInfo_Marker_AreEquivalent compares two nullable positions.
//
// Go-only: language-adapter — a nil pointer converted to the interface
// parameter is Java null (jutil.IsNil).
//
// Java: JSDocInfo.Marker.areEquivalent(ComparableSourcePosition<T>, ComparableSourcePosition<T>)
func JSDocInfo_Marker_AreEquivalent[T any](p1 JSDocInfo_ComparableSourcePosition[T], p2 JSDocInfo_ComparableSourcePosition[T]) bool {
	p1Null, p2Null := jutil.IsNil(p1), jutil.IsNil(p2)
	if p2Null {
		p2 = nil
	}
	return p1Null == p2Null && (p1Null || p1.IsEquivalentTo(p2))
}

// String is java.lang.Object.toString() (Marker does not override it).
//
// Go-only: language-adapter — Java prints an identity hash code, which
// differs from run to run; this one is derived from the object's address.
func (m *JSDocInfo_Marker) String() string {
	return jutil.ObjectToString("com.google.javascript.rhino.JSDocInfo$Marker", m)
}

// This includes:
//
//   - base type
//   - @extends
//   - @implements
//   - @lend
//   - @param
//   - @return
//   - @template
//   - @this
//   - @throws
//   - @type
//
// Any future type specific JSDoc should make sure to add the appropriate
// nodes here.
//
// JSDocTypeExprIterator is a stack-allocated, zero-closure cursor over all
// non-nil JSTypeExpressions in a JSDocInfo.
type JSDocTypeExprIterator struct {
	info     *JSDocInfo
	typeBits uint64
	curList  []*JSTypeExpression
	curMap   *jutil.LinkedHashMap[string, *JSTypeExpression]
	subIdx   int
	cur      *JSTypeExpression
}

// TypeExprIter returns a stack-allocated JSDocTypeExprIterator for i.
func (i *JSDocInfo) TypeExprIter() JSDocTypeExprIterator {
	if i == nil {
		return JSDocTypeExprIterator{}
	}
	return JSDocTypeExprIterator{
		info:     i,
		typeBits: uint64(i.propertyKeysBitset & jsDocInfo_typePropertiesMask),
	}
}

// Expr returns the current JSTypeExpression after Next returns true.
func (it *JSDocTypeExprIterator) Expr() *JSTypeExpression { return it.cur }

// Node returns the root Node of the current JSTypeExpression after Next returns true.
func (it *JSDocTypeExprIterator) Node() Node { return it.cur.GetRoot() }

// Next advances to the next non-nil JSTypeExpression and reports whether one exists.
func (it *JSDocTypeExprIterator) Next() bool {
	for {
		if it.curList != nil {
			for it.subIdx < len(it.curList) {
				e := it.curList[it.subIdx]
				it.subIdx++
				if e != nil {
					it.cur = e
					return true
				}
			}
			it.curList = nil
		}
		if it.curMap != nil {
			for {
				v, nextSlot, ok := it.curMap.NextSlotValue(it.subIdx)
				if !ok {
					break
				}
				it.subIdx = nextSlot
				if v != nil {
					it.cur = v
					return true
				}
			}
			it.curMap = nil
		}
		if it.typeBits == 0 {
			it.cur = nil
			return false
		}
		low := bits.TrailingZeros64(it.typeBits)
		propMask := uint64(1) << low
		it.typeBits &^= propMask
		index := bits.OnesCount64(uint64(it.info.propertyKeysBitset) & (propMask - 1))
		val := it.info.GetPropertyValueByIndex(index)
		switch JSDocInfo_Property_values[low].kind {
		case jsDocPropKindType:
			if e := jsdocGet[*JSTypeExpression](val); e != nil {
				it.cur = e
				return true
			}
		case jsDocPropKindTypeList:
			it.curList = jsdocGet[[]*JSTypeExpression](val)
			it.subIdx = 0
		case jsDocPropKindTypeMap:
			it.curMap = jsdocGet[*jutil.LinkedHashMap[string, *JSTypeExpression]](val)
			it.subIdx = 0
		}
	}
}

// TypeExpressions iterates over all non-nil JSTypeExpressions in this JSDocInfo
// without allocating intermediate slices.
func (i *JSDocInfo) TypeExpressions() iter.Seq[*JSTypeExpression] {
	return func(yield func(*JSTypeExpression) bool) {
		for it := i.TypeExprIter(); it.Next(); {
			if !yield(it.Expr()) {
				return
			}
		}
	}
}

// TypeNodes iterates over all type AST root nodes in this JSDocInfo without
// allocating intermediate slices.
func (i *JSDocInfo) TypeNodes() iter.Seq[Node] {
	return func(yield func(Node) bool) {
		for it := i.TypeExprIter(); it.Next(); {
			if !yield(it.Node()) {
				return
			}
		}
	}
}

// HasTypeNodes reports whether this JSDocInfo contains at least one type node,
// without allocating slices or closures.
func (i *JSDocInfo) HasTypeNodes() bool {
	it := i.TypeExprIter()
	return it.Next()
}

// Java: JSDocInfo.getTypeExpressions()
func (i *JSDocInfo) GetTypeExpressions() []*JSTypeExpression {
	builder := []*JSTypeExpression{}
	for it := i.TypeExprIter(); it.Next(); {
		builder = append(builder, it.Expr())
	}
	return builder
}

// GetTypeNodes returns a collection of all type nodes that are a part of
// this JSDocInfo.
//
// This includes:
//
//   - @extends
//   - @implements
//   - @lend
//   - @param
//   - @return
//   - @template
//   - @this
//   - @throws
//   - @type
//
// Any future type specific JSDoc should make sure to add the appropriate
// nodes here.
//
// Java: JSDocInfo.getTypeNodes()
func (i *JSDocInfo) GetTypeNodes() []Node {
	nodes := []Node{}

	for it := i.TypeExprIter(); it.Next(); {
		nodes = append(nodes, it.Node())
	}

	return nodes
}
