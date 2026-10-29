// Port of com/google/javascript/jscomp/disambiguate/ColorGraphNodeFactory.java @ closure-compiler v20261021

/*
 * Copyright 2020 The Closure Compiler Authors.
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
	"slices"

	"github.com/paullepoulpe/jsgo/internal/colors"
	"github.com/paullepoulpe/jsgo/internal/jutil"
)

type colorGraphNodeEntry struct {
	key         *colors.Color
	node        *ColorGraphNode
	nextPlusOne uint32
}

// colorGraphNodeIndex maps Colors (by structural Equals/HashCode) to ColorGraphNodes
// in deterministic insertion order using a flat slice + open hash bucket table.
type colorGraphNodeIndex struct {
	entries []colorGraphNodeEntry
	buckets map[int32]uint32
}

func (idx *colorGraphNodeIndex) init(capHint int) {
	idx.entries = make([]colorGraphNodeEntry, 0, capHint)
	idx.buckets = make(map[int32]uint32, capHint)
}

func (idx *colorGraphNodeIndex) Lookup(key *colors.Color) (*ColorGraphNode, bool) {
	if key == nil || len(idx.entries) == 0 {
		return nil, false
	}
	for cur := idx.buckets[key.HashCode()]; cur != 0; cur = idx.entries[cur-1].nextPlusOne {
		e := &idx.entries[cur-1]
		if e.key == key || e.key.Equals(key) {
			return e.node, true
		}
	}
	return nil, false
}

func (idx *colorGraphNodeIndex) Put(key *colors.Color, node *ColorGraphNode) {
	if idx.buckets == nil {
		idx.buckets = make(map[int32]uint32, 16)
	}
	h := key.HashCode()
	for cur := idx.buckets[h]; cur != 0; cur = idx.entries[cur-1].nextPlusOne {
		e := &idx.entries[cur-1]
		if e.key == key || e.key.Equals(key) {
			e.node = node
			return
		}
	}
	pos := uint32(len(idx.entries))
	idx.entries = append(idx.entries, colorGraphNodeEntry{
		key:         key,
		node:        node,
		nextPlusOne: idx.buckets[h],
	})
	idx.buckets[h] = pos + 1
}

func (idx *colorGraphNodeIndex) Size() int {
	return len(idx.entries)
}

func (idx *colorGraphNodeIndex) Values() iter.Seq[*ColorGraphNode] {
	return func(yield func(*ColorGraphNode) bool) {
		for i := range idx.entries {
			if !yield(idx.entries[i].node) {
				return
			}
		}
	}
}

// ColorGraphNodeSet is an insertion-ordered set of ColorGraphNodes.
type ColorGraphNodeSet struct {
	order   []*ColorGraphNode
	members map[*ColorGraphNode]struct{}
}

// Contains reports whether n is in s.
func (s ColorGraphNodeSet) Contains(n *ColorGraphNode) bool {
	_, ok := s.members[n]
	return ok
}

// All iterates the nodes in insertion order.
func (s ColorGraphNodeSet) All() iter.Seq[*ColorGraphNode] {
	return slices.Values(s.order)
}

// Size returns the number of distinct nodes.
func (s ColorGraphNodeSet) Size() int { return len(s.order) }

// ToSlice returns a defensive copy of the nodes in insertion order.
func (s ColorGraphNodeSet) ToSlice() []*ColorGraphNode {
	if len(s.order) == 0 {
		return []*ColorGraphNode{}
	}
	return slices.Clone(s.order)
}

// ColorGraphNodeFactory is a factory and cache for ColorGraphNode
// instances.
type ColorGraphNodeFactory struct {
	// Colors are keyed with their equals()/hashCode().
	typeIndex colorGraphNodeIndex
	registry  *colors.ColorRegistry
	nodeSlab  jutil.Slab[ColorGraphNode]
}

// Java: ColorGraphNodeFactory(LinkedHashMap<Color, ColorGraphNode>, ColorRegistry)
func NewColorGraphNodeFactory(registry *colors.ColorRegistry) *ColorGraphNodeFactory {
	f := &ColorGraphNodeFactory{registry: registry}
	f.typeIndex.init(64)
	return f
}

func (f *ColorGraphNodeFactory) allocNode(single *colors.Color, index int) *ColorGraphNode {
	jutil.CheckArgument(index >= 0)
	jutil.CheckNotNullPtr(single)
	n := f.nodeSlab.New()
	n.index = index
	n.color = single
	return n
}

// Java: ColorGraphNodeFactory.createFactory(ColorRegistry)
func ColorGraphNodeFactory_CreateFactory(colorRegistry *colors.ColorRegistry) *ColorGraphNodeFactory {
	f := NewColorGraphNodeFactory(colorRegistry)
	unknownColorNode := f.allocNode(colors.StandardColors_UNKNOWN, 0)
	f.typeIndex.Put(colors.StandardColors_UNKNOWN, unknownColorNode)
	return f
}

// CreateNode returns the ColorGraphNode known by this factory for type.
//
// For a given type and factory, this method will always return the same
// result. The results are cached.
//
// Java: ColorGraphNodeFactory.createNode(Color)
func (f *ColorGraphNodeFactory) CreateNode(color *colors.Color) *ColorGraphNode {
	key := f.simplifyColor(color)
	// typeIndex.computeIfAbsent(key, this::newColorGraphNode)
	if node, ok := f.typeIndex.Lookup(key); ok && node != nil {
		return node
	}
	node := f.newColorGraphNode(key)
	f.typeIndex.Put(key, node)
	return node
}

// Java: ColorGraphNodeFactory.getAllKnownTypes()
func (f *ColorGraphNodeFactory) GetAllKnownTypes() ColorGraphNodeSet {
	n := f.typeIndex.Size()
	members := make(map[*ColorGraphNode]struct{}, n)
	order := make([]*ColorGraphNode, 0, n)
	for v := range f.typeIndex.Values() {
		if _, exists := members[v]; !exists {
			members[v] = struct{}{}
			order = append(order, v)
		}
	}
	return ColorGraphNodeSet{order: order, members: members}
}

// Java: ColorGraphNodeFactory.newColorGraphNode(Color)
func (f *ColorGraphNodeFactory) newColorGraphNode(key *colors.Color) *ColorGraphNode {
	return f.allocNode(key, f.typeIndex.Size())
}

// simplifyColor merges different colors with the same ambiguation-behavior
// into one.
//
// Java: ColorGraphNodeFactory.simplifyColor(Color)
func (f *ColorGraphNodeFactory) simplifyColor(type_ *colors.Color) *colors.Color {
	if type_ == nil {
		return colors.StandardColors_UNKNOWN
	}

	if type_.IsUnion() {
		// First remove null/void, then recursively simplify any primitive
		// components
		type_ = type_.SubtractNullOrVoid()
		if type_.IsUnion() {
			elems := type_.GetUnionElements()
			var stackBuf [8]*colors.Color
			var simplified []*colors.Color
			if len(elems) <= len(stackBuf) {
				simplified = stackBuf[:0]
			} else {
				simplified = make([]*colors.Color, 0, len(elems))
			}
			changed := false
			for _, e := range elems {
				se := f.simplifyColor(e)
				if se != e {
					changed = true
				}
				simplified = append(simplified, se)
			}
			if !changed {
				return type_
			}
			return colors.Color_CreateUnionSlice(simplified)
		}
		return f.simplifyColor(type_)
	} else if type_.GetBoxId() != nil {
		return f.registry.Get(*type_.GetBoxId())
	} else if type_.Equals(colors.StandardColors_NULL_OR_VOID) {
		return colors.StandardColors_UNKNOWN
	} else {
		return type_
	}
}
