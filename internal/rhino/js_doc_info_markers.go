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

// JSDocInfo_Marker defines a struct for containing the parsing information
// for this JSDocInfo. For each annotation found in the JsDoc, a marker will
// be created indicating the annotation itself, the name of the annotation
// (if any; for example, a @param has a name, but a @return does not), the
// textual description found on that annotation and, if applicable, the type
// declaration. All this information is only collected if documentation
// collection is turned on.
//
// Go-only: performance — embeds the four SourcePosition fields inline by value
// so markers live in a flat []JSDocInfo_Marker slice without per-marker or
// per-position heap allocations.
type JSDocInfo_Marker struct {
	annotation  SourcePosition[string]
	nameNode    SourcePosition[Node]
	description SourcePosition[string]
	type_       SourcePosition[Node]
}

// Java: JSDocInfo.Marker()
func NewJSDocInfo_Marker() *JSDocInfo_Marker { return &JSDocInfo_Marker{} }

// GetAnnotation gets the position information for the annotation name.
// (e.g., "param")
//
// Java: JSDocInfo.Marker.getAnnotation()
func (m *JSDocInfo_Marker) GetAnnotation() *SourcePosition[string] {
	if !m.annotation.valid {
		return nil
	}
	return &m.annotation
}

// GetNameNode gets the position information for the name found in a @param
// tag.
//
// Java: JSDocInfo.Marker.getNameNode()
func (m *JSDocInfo_Marker) GetNameNode() *SourcePosition[Node] {
	if !m.nameNode.valid {
		return nil
	}
	return &m.nameNode
}

// GetDescription gets the position information for the description found in
// a block tag.
//
// Java: JSDocInfo.Marker.getDescription()
func (m *JSDocInfo_Marker) GetDescription() *SourcePosition[string] {
	if !m.description.valid {
		return nil
	}
	return &m.description
}

// GetType gets the position information for the type expression found in
// some block tags, like "@param" and "@return".
//
// Java: JSDocInfo.Marker.getType()
func (m *JSDocInfo_Marker) GetType() *SourcePosition[Node] {
	if !m.type_.valid {
		return nil
	}
	return &m.type_
}

func stringPosEquivalent(a, b *SourcePosition[string]) bool {
	if a.valid != b.valid {
		return false
	}
	return !a.valid || (a.IsSamePositionAs(b) && a.item == b.item)
}

func nodePosEquivalent(a, b *SourcePosition[Node]) bool {
	if a.valid != b.valid {
		return false
	}
	if !a.valid {
		return true
	}
	if !a.IsSamePositionAs(b) || a.brackets != b.brackets || a.item.IsNil() != b.item.IsNil() {
		return false
	}
	return a.item.IsNil() || a.item.IsEquivalentTo(b.item)
}

// Java: JSDocInfo.Marker.isEquivalentTo(Marker)
func (m *JSDocInfo_Marker) IsEquivalentTo(that *JSDocInfo_Marker) bool {
	if m == nil || that == nil {
		return m == that
	}
	return stringPosEquivalent(&m.annotation, &that.annotation) &&
		nodePosEquivalent(&m.nameNode, &that.nameNode) &&
		stringPosEquivalent(&m.description, &that.description) &&
		nodePosEquivalent(&m.type_, &that.type_)
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
	info      *JSDocInfo
	typeBits  uint64
	curList   []*JSTypeExpression
	curParams []JSDocParam
	subIdx    int
	cur       *JSTypeExpression
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
		if it.curParams != nil {
			for it.subIdx < len(it.curParams) {
				e := it.curParams[it.subIdx].Type
				it.subIdx++
				if e != nil {
					it.cur = e
					return true
				}
			}
			it.curParams = nil
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
			it.curParams = jsdocGet[[]JSDocParam](val)
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
