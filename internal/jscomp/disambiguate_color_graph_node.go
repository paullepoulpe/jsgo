// Port of com/google/javascript/jscomp/disambiguate/ColorGraphNode.java @ closure-compiler v20261021

/*
 * Copyright 2008 The Closure Compiler Authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package jscomp

import (
	"iter"

	"github.com/paullepoulpe/jsgo/internal/colors"
	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// ColorGraphNode is a struct representing a Color for use in ambiguation.
//
// Each instance pairs a Color with additional information computed by/for
// optimizations.
//
// Note: this design now depends on the implementation of Color to preserve
// invariants about recursive Colors. The node factory can't be depended on
// for post-processing Colors without losing type safety.
type ColorGraphNode struct {
	color *colors.Color

	props       []*PropertyClustering
	propsInline [4]*PropertyClustering
	propSet     map[*PropertyClustering]struct{}

	subtypeIndices *jutil.BitSet

	index            int
	reflexiveSubtype bool
}

// ColorGraphNode_PropAssociation lists the reasons a property name became
// associated with a type.
//
// This information is only used for debugging. It doesn't affect the
// behaviour of the pass.
type ColorGraphNode_PropAssociation int32

const (
	// ColorGraphNode_PropAssociation_AST: because of an access in the AST
	// (e.g. foo.prop)
	ColorGraphNode_PropAssociation_AST ColorGraphNode_PropAssociation = iota
	// ColorGraphNode_PropAssociation_TYPE_SYSTEM: because the type system
	// recorded such an association, without any reason visible in the AST.
	// Usually this means the relevant AST segment has been optimized away
	// by an earlier pass.
	ColorGraphNode_PropAssociation_TYPE_SYSTEM
	// ColorGraphNode_PropAssociation_SUPERTYPE: because it was inherited
	// from a supertype in the type graph
	ColorGraphNode_PropAssociation_SUPERTYPE
)

var colorGraphNode_PropAssociation_names = [...]string{"AST", "TYPE_SYSTEM", "SUPERTYPE"}

// Name is Enum.name().
func (a ColorGraphNode_PropAssociation) Name() string { return colorGraphNode_PropAssociation_names[a] }

// String is Enum.toString().
func (a ColorGraphNode_PropAssociation) String() string { return a.Name() }

// Ordinal is Enum.ordinal().
func (a ColorGraphNode_PropAssociation) Ordinal() int { return int(a) }

// Java: ColorGraphNode.create(Color, int)
func ColorGraphNode_Create(single *colors.Color, index int) *ColorGraphNode {
	jutil.CheckArgument(index >= 0)
	return newColorGraphNode(single, index)
}

// Java: ColorGraphNode.createForTesting(int)
func ColorGraphNode_CreateForTesting(index int) *ColorGraphNode {
	jutil.CheckArgument(index < 0) // All test nodes have negative indexs to differentiate them.
	return newColorGraphNode(colors.StandardColors_UNKNOWN, index)
}

// Java: ColorGraphNode.createForTesting(Color, int)
func ColorGraphNode_CreateForTestingWithColor(single *colors.Color, index int) *ColorGraphNode {
	jutil.CheckArgument(index < 0) // All test nodes have negative indexs to differentiate them.
	return newColorGraphNode(single, index)
}

// Java: ColorGraphNode(Color, int)
func newColorGraphNode(single *colors.Color, index int) *ColorGraphNode {
	jutil.CheckNotNullPtr(single)
	return &ColorGraphNode{
		index: index,
		color: single,
	}
}

// Java: ColorGraphNode.getColor()
func (n *ColorGraphNode) GetColor() *colors.Color {
	return jutil.CheckNotNullPtr(n.color)
}

// GetIndex is an ID used to efficiently construct a unique name for any
// cluster this node becomes the represenative of.
//
// Java: ColorGraphNode.getIndex()
func (n *ColorGraphNode) GetIndex() int {
	return n.index
}

// Props returns the properties associated with n in insertion order.
func (n *ColorGraphNode) Props() []*PropertyClustering {
	return n.props
}

// HasProp reports whether prop is associated with n.
func (n *ColorGraphNode) HasProp(prop *PropertyClustering) bool {
	if n.propSet != nil {
		_, ok := n.propSet[prop]
		return ok
	}
	for _, p := range n.props {
		if p == prop {
			return true
		}
	}
	return false
}

// AddProp associates prop with n if not already present, returning true if it
// was newly added.
func (n *ColorGraphNode) AddProp(prop *PropertyClustering) bool {
	if n.propSet != nil {
		if _, ok := n.propSet[prop]; ok {
			return false
		}
		n.propSet[prop] = struct{}{}
		n.props = append(n.props, prop)
		return true
	}
	for _, p := range n.props {
		if p == prop {
			return false
		}
	}
	if n.props == nil {
		n.props = n.propsInline[:0]
	}
	n.props = append(n.props, prop)
	if len(n.props) > 16 {
		m := make(map[*PropertyClustering]struct{}, 32)
		for _, p := range n.props {
			m[p] = struct{}{}
		}
		n.propSet = m
	}
	return true
}

// ColorGraphNodePropsView adapts ColorGraphNode's flat property list to the
// map-like API expected by tests.
type ColorGraphNodePropsView struct {
	n *ColorGraphNode
}

func (v ColorGraphNodePropsView) Size() int { return len(v.n.props) }

func (v ColorGraphNodePropsView) IsEmpty() bool { return len(v.n.props) == 0 }

func (v ColorGraphNodePropsView) ContainsKey(prop *PropertyClustering) bool {
	return v.n.HasProp(prop)
}

func (v ColorGraphNodePropsView) Put(prop *PropertyClustering, _ ColorGraphNode_PropAssociation) {
	v.n.AddProp(prop)
}

func (v ColorGraphNodePropsView) PutIfAbsent(prop *PropertyClustering, _ ColorGraphNode_PropAssociation) {
	v.n.AddProp(prop)
}

func (v ColorGraphNodePropsView) Keys() iter.Seq[*PropertyClustering] {
	return func(yield func(*PropertyClustering) bool) {
		for _, p := range v.n.props {
			if !yield(p) {
				return
			}
		}
	}
}

// GetAssociatedProps is the set of properties that might be accessed
// from this type.
//
// Java: ColorGraphNode.getAssociatedProps()
func (n *ColorGraphNode) GetAssociatedProps() ColorGraphNodePropsView {
	return ColorGraphNodePropsView{n: n}
}

func (n *ColorGraphNode) AssociatedPropsSize() int {
	return len(n.props)
}

// GetSubtypeIndices returns the IDs of other ColorGraphNodes that have been
// found to be subtypes of this type.
//
// Java: ColorGraphNode.getSubtypeIndices()
func (n *ColorGraphNode) GetSubtypeIndices() *jutil.BitSet {
	if n.subtypeIndices == nil {
		n.subtypeIndices = jutil.NewBitSet()
		if n.reflexiveSubtype && n.index >= 0 {
			n.subtypeIndices.Set(n.index)
		}
	}
	return n.subtypeIndices
}

// String is for debugging only (@DoNotCall).
//
// Java: ColorGraphNode.toString()
func (n *ColorGraphNode) String() string {
	// Just report a few important properties of color instead of the whole
	// thing. The string generated here will become the label of the `.dot`
	// graph node if this graph is logged for debugging. Including the
	// entire color, recursively includes all the colors it extends also,
	// making the labels unmanageably long to display and repeating
	// information available from other nodes.
	return jutil.MoreObjects_ToStringHelper("ColorGraphNode").
		Add("index", n.index).
		Add("color.id", n.color.GetId()).
		Add("color.ownProperties", jutil.CollectionToString(n.color.GetOwnProperties())).
		String()
}
