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
	"github.com/paullepoulpe/jsgo/internal/colors"
	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// ColorGraphNodeFactory is a factory and cache for ColorGraphNode
// instances.
type ColorGraphNodeFactory struct {
	// Colors are keyed with their equals()/hashCode().
	typeIndex *jutil.ObjectLinkedHashMap[*colors.Color, *ColorGraphNode]
	registry  *colors.ColorRegistry
	nodeSlab  jutil.Slab[ColorGraphNode]
}

// Java: ColorGraphNodeFactory(LinkedHashMap<Color, ColorGraphNode>, ColorRegistry)
func NewColorGraphNodeFactory(
	initialTypeIndex *jutil.ObjectLinkedHashMap[*colors.Color, *ColorGraphNode], registry *colors.ColorRegistry) *ColorGraphNodeFactory {
	return &ColorGraphNodeFactory{typeIndex: initialTypeIndex, registry: registry}
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
	typeIndex := jutil.NewObjectLinkedHashMap[*colors.Color, *ColorGraphNode]()
	f := NewColorGraphNodeFactory(typeIndex, colorRegistry)
	unknownColorNode := f.allocNode(colors.StandardColors_UNKNOWN, 0)
	typeIndex.Put(colors.StandardColors_UNKNOWN, unknownColorNode)
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
func (f *ColorGraphNodeFactory) GetAllKnownTypes() *jutil.LinkedHashSet[*ColorGraphNode] {
	out := jutil.NewLinkedHashSet[*ColorGraphNode]()
	for v := range f.typeIndex.Values() {
		out.Add(v)
	}
	return out
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
