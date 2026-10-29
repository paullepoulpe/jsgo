// Port of com/google/javascript/jscomp/colors/ColorRegistry.java @ closure-compiler v20261021

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

package colors

import (
	"iter"
	"maps"
	"slices"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// ColorRegistry stores "out of band" information about the set of Colors in
// use by a compilation.
//
// This includes any information that isn't specifically can't be set on a
// single Color object,
type ColorRegistry struct {
	nativeColors             map[ColorId]*Color
	disambiguationSupertypes map[*Color][]*Color
	mismatchLocations        map[ColorId][]string
	mismatchOrder            []ColorId
}

// Java: ColorRegistry(Builder)
func NewColorRegistry(builder *ColorRegistry_Builder) *ColorRegistry {
	r := &ColorRegistry{}
	if !builder.consumed {
		builder.consumed = true
		r.nativeColors = builder.nativeColors
		r.disambiguationSupertypes = builder.disambiguationSupertypes
		r.mismatchLocations = builder.mismatchLocations
		r.mismatchOrder = builder.mismatchOrder
	} else {
		r.nativeColors = maps.Clone(builder.nativeColors)
		r.disambiguationSupertypes = cloneDisambiguationSupertypes(builder.disambiguationSupertypes)
		r.mismatchLocations = cloneMismatchLocations(builder.mismatchLocations)
		r.mismatchOrder = slices.Clone(builder.mismatchOrder)
	}

	ok := len(r.nativeColors) == ColorRegistry_REQUIRED_IDS.Size()
	if ok {
		for id := range ColorRegistry_REQUIRED_IDS.All() {
			if _, has := r.nativeColors[id]; !has {
				ok = false
				break
			}
		}
	}
	jutil.CheckState(ok)
	return r
}

func cloneDisambiguationSupertypes(src map[*Color][]*Color) map[*Color][]*Color {
	if len(src) == 0 {
		return nil
	}
	dst := make(map[*Color][]*Color, len(src))
	for k, vs := range src { // order-independent: copying distinct *Color keys into a new map
		dst[k] = slices.Clone(vs)
	}
	return dst
}

func cloneMismatchLocations(src map[ColorId][]string) map[ColorId][]string {
	if len(src) == 0 {
		return nil
	}
	dst := make(map[ColorId][]string, len(src))
	for k, vs := range src { // order-independent: copying distinct ColorId keys into a new map
		dst[k] = slices.Clone(vs)
	}
	return dst
}

// Java: ColorRegistry.get(ColorId)
func (r *ColorRegistry) Get(id ColorId) *Color {
	// checkNotNull(id): ColorId is a Go value type and cannot be null.
	c := r.nativeColors[id]
	if c == nil {
		panic(jutil.NewNullPointerExceptionMsg("Missing Color for " + id.String()))
	}
	return c
}

// GetDisambiguationSupertypesSlice returns the colors directly above `x` in the
// subtyping graph as a read-only slice in insertion order without allocation.
func (r *ColorRegistry) GetDisambiguationSupertypesSlice(x *Color) []*Color {
	jutil.CheckNotNull(x)
	return r.disambiguationSupertypes[x]
}

// GetDisambiguationSupertypes returns the colors directly above `x` in the
// subtyping graph for the purposes of property (dis)ambiguation.
//
// Java: ColorRegistry.getDisambiguationSupertypes(Color)
func (r *ColorRegistry) GetDisambiguationSupertypes(x *Color) *jutil.ObjectLinkedHashSet[*Color] {
	jutil.CheckNotNull(x)
	sups := r.disambiguationSupertypes[x]
	if len(sups) == 0 {
		return colorEmptyColorSet
	}
	out := jutil.NewObjectLinkedHashSet[*Color]()
	for _, s := range sups {
		out.Add(s)
	}
	return out
}

// ColorMismatchLocations maps ColorIds to mismatch location sourceref strings in insertion order.
type ColorMismatchLocations struct {
	m    *jutil.ObjectLinkedHashMap[ColorId, *jutil.LinkedHashSet[string]]
	size int
}

// NewColorMismatchLocations returns an empty ColorMismatchLocations.
func NewColorMismatchLocations() *ColorMismatchLocations {
	return &ColorMismatchLocations{
		m: jutil.NewObjectLinkedHashMap[ColorId, *jutil.LinkedHashSet[string]](),
	}
}

// Put adds (id, loc) if not already present.
func (cml *ColorMismatchLocations) Put(id ColorId, loc string) bool {
	s, ok := cml.m.Lookup(id)
	if !ok {
		s = jutil.NewLinkedHashSet[string]()
		cml.m.Put(id, s)
	}
	if s.Add(loc) {
		cml.size++
		return true
	}
	return false
}

// Get returns the mismatch location strings for id (empty set when absent).
func (cml *ColorMismatchLocations) Get(id ColorId) *jutil.LinkedHashSet[string] {
	if s, ok := cml.m.Lookup(id); ok {
		return s
	}
	return jutil.NewLinkedHashSet[string]()
}

// Size returns the total number of (ColorId, location) pairs.
func (cml *ColorMismatchLocations) Size() int { return cml.size }

// Entries iterates (ColorId, location) pairs in insertion order.
func (cml *ColorMismatchLocations) Entries() iter.Seq2[ColorId, string] {
	return func(yield func(ColorId, string) bool) {
		for id, locs := range cml.m.All() {
			for loc := range locs.All() {
				if !yield(id, loc) {
					return
				}
			}
		}
	}
}

func (cml *ColorMismatchLocations) String() string { return cml.m.String() }

// InverseColorMismatchLocations maps mismatch location strings to involved ColorIds.
type InverseColorMismatchLocations struct {
	m *jutil.ObjectLinkedHashMap[string, *jutil.ObjectLinkedHashSet[ColorId]]
}

// AsMapObject exposes the underlying map for LogsGson serialization.
func (inv *InverseColorMismatchLocations) AsMapObject() any { return inv.m }

// Inverse returns a location-to-ColorIds index preserving insertion order.
func (cml *ColorMismatchLocations) Inverse() *InverseColorMismatchLocations {
	inv := jutil.NewObjectLinkedHashMap[string, *jutil.ObjectLinkedHashSet[ColorId]]()
	for id, locs := range cml.m.All() {
		for loc := range locs.All() {
			s, ok := inv.Lookup(loc)
			if !ok {
				s = jutil.NewObjectLinkedHashSet[ColorId]()
				inv.Put(loc, s)
			}
			s.Add(id)
		}
	}
	return &InverseColorMismatchLocations{m: inv}
}

// GetMismatchLocationsForDebugging returns an index (mismatch sourceref
// string) => (involved color ids).
//
// This index is only intended for debugging. It may be empty or incomplete
// during production compilations.
//
// Java: ColorRegistry.getMismatchLocationsForDebugging()
func (r *ColorRegistry) GetMismatchLocationsForDebugging() *ColorMismatchLocations {
	out := NewColorMismatchLocations()
	for _, id := range r.mismatchOrder {
		for _, loc := range r.mismatchLocations[id] {
			out.Put(id, loc)
		}
	}
	return out
}

// Java: ColorRegistry.builder()
func ColorRegistry_Builder_() *ColorRegistry_Builder {
	return NewColorRegistry_Builder()
}

// ColorRegistry_Builder builds a ColorRegistry.
type ColorRegistry_Builder struct {
	nativeColors             map[ColorId]*Color
	disambiguationSupertypes map[*Color][]*Color
	mismatchLocations        map[ColorId][]string
	mismatchOrder            []ColorId
	consumed                 bool
}

// Java: ColorRegistry.Builder()
func NewColorRegistry_Builder() *ColorRegistry_Builder {
	return &ColorRegistry_Builder{
		nativeColors: make(map[ColorId]*Color, ColorRegistry_REQUIRED_IDS.Size()),
	}
}

func (b *ColorRegistry_Builder) ensureMutable() {
	if b.consumed {
		b.nativeColors = maps.Clone(b.nativeColors)
		b.disambiguationSupertypes = cloneDisambiguationSupertypes(b.disambiguationSupertypes)
		b.mismatchLocations = cloneMismatchLocations(b.mismatchLocations)
		b.mismatchOrder = slices.Clone(b.mismatchOrder)
		b.consumed = false
	}
}

// Java: ColorRegistry.Builder.setNativeColor(Color)
func (b *ColorRegistry_Builder) SetNativeColor(x *Color) *ColorRegistry_Builder {
	b.ensureMutable()
	jutil.CheckStateMsg1(ColorRegistry_REQUIRED_IDS.Contains(x.GetId()), "%s", x)
	b.nativeColors[x.GetId()] = x
	return b
}

// Java: ColorRegistry.Builder.addDisambiguationEdge(Color, Color)
func (b *ColorRegistry_Builder) AddDisambiguationEdge(subtype *Color, supertype *Color) *ColorRegistry_Builder {
	jutil.CheckNotNull(subtype)
	jutil.CheckNotNull(supertype)
	b.ensureMutable()
	if b.disambiguationSupertypes == nil {
		b.disambiguationSupertypes = make(map[*Color][]*Color)
	}
	sups := b.disambiguationSupertypes[subtype]
	for _, existing := range sups {
		if existing == supertype || existing.Equals(supertype) {
			return b
		}
	}
	b.disambiguationSupertypes[subtype] = append(sups, supertype)
	return b
}

// Java: ColorRegistry.Builder.addMismatchLocation(ColorId, String)
func (b *ColorRegistry_Builder) AddMismatchLocation(id ColorId, location string) *ColorRegistry_Builder {
	b.ensureMutable()
	if b.mismatchLocations == nil {
		b.mismatchLocations = make(map[ColorId][]string)
	}
	locs, exists := b.mismatchLocations[id]
	if !exists {
		b.mismatchOrder = append(b.mismatchOrder, id)
	} else {
		for _, existing := range locs {
			if existing == location {
				return b
			}
		}
	}
	b.mismatchLocations[id] = append(locs, location)
	return b
}

// SetDefaultNativeColorsForTesting sets defaults for native Colors.
//
// Only for use in testing. In real compilations, certain native colors have
// fields that vary from compilation-to-compilation (like whether the
// "Number" object is invalidating), so should use setNativeColor instead.
//
// Java: ColorRegistry.Builder.setDefaultNativeColorsForTesting()
func (b *ColorRegistry_Builder) SetDefaultNativeColorsForTesting() *ColorRegistry_Builder {
	for id := range ColorRegistry_REQUIRED_IDS.All() {
		b.SetNativeColor(Color_SingleBuilder().SetId(id).Build())
	}
	return b
}

// Java: ColorRegistry.Builder.build()
func (b *ColorRegistry_Builder) Build() *ColorRegistry {
	return NewColorRegistry(b)
}

// Java: ColorRegistry.REQUIRED_IDS
var ColorRegistry_REQUIRED_IDS = StandardColors_STANDARD_OBJECT_IDS
