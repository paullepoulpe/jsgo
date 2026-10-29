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
	"strings"

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
	sups := r.GetDisambiguationSupertypesSlice(x)
	out := jutil.NewObjectLinkedHashSet[*Color]()
	for _, s := range sups {
		out.Add(s)
	}
	return out
}

// ColorMismatchLocSet is an insertion-ordered, deduplicated view of mismatch location strings.
type ColorMismatchLocSet []string

// Contains reports whether loc is in s.
func (s ColorMismatchLocSet) Contains(loc string) bool {
	return slices.Contains(s, loc)
}

// ToSlice returns a defensive copy of the locations (non-nil empty slice when empty).
func (s ColorMismatchLocSet) ToSlice() []string {
	if len(s) == 0 {
		return []string{}
	}
	return slices.Clone([]string(s))
}

// All iterates locations in insertion order.
func (s ColorMismatchLocSet) All() iter.Seq[string] {
	return slices.Values([]string(s))
}

// Size returns the number of locations.
func (s ColorMismatchLocSet) Size() int { return len(s) }

// IsEmpty reports whether s has no locations.
func (s ColorMismatchLocSet) IsEmpty() bool { return len(s) == 0 }

// ColorMismatchLocations maps ColorIds to mismatch location sourceref strings in insertion order.
type ColorMismatchLocations struct {
	locs  map[ColorId][]string
	order []ColorId
	size  int
}

// NewColorMismatchLocations returns an empty ColorMismatchLocations.
func NewColorMismatchLocations() *ColorMismatchLocations {
	return &ColorMismatchLocations{}
}

// Put adds (id, loc) if not already present.
func (cml *ColorMismatchLocations) Put(id ColorId, loc string) bool {
	if cml.locs == nil {
		cml.locs = make(map[ColorId][]string)
	}
	existing, ok := cml.locs[id]
	if !ok {
		cml.order = append(cml.order, id)
	} else if slices.Contains(existing, loc) {
		return false
	}
	cml.locs[id] = append(existing, loc)
	cml.size++
	return true
}

// Get returns the mismatch location strings for id (empty set when absent).
func (cml *ColorMismatchLocations) Get(id ColorId) ColorMismatchLocSet {
	return ColorMismatchLocSet(cml.locs[id])
}

// Size returns the total number of (ColorId, location) pairs.
func (cml *ColorMismatchLocations) Size() int { return cml.size }

// Entries iterates (ColorId, location) pairs in insertion order.
func (cml *ColorMismatchLocations) Entries() iter.Seq2[ColorId, string] {
	return func(yield func(ColorId, string) bool) {
		for _, id := range cml.order {
			for _, loc := range cml.locs[id] {
				if !yield(id, loc) {
					return
				}
			}
		}
	}
}

func (cml *ColorMismatchLocations) String() string {
	if len(cml.order) == 0 {
		return "{}"
	}
	var sb strings.Builder
	sb.WriteByte('{')
	for i, id := range cml.order {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(id.String())
		sb.WriteString("=[")
		for j, loc := range cml.locs[id] {
			if j > 0 {
				sb.WriteString(", ")
			}
			sb.WriteString(loc)
		}
		sb.WriteByte(']')
	}
	sb.WriteByte('}')
	return sb.String()
}

// InverseColorMismatchLocations maps mismatch location strings to involved ColorIds.
type InverseColorMismatchLocations struct {
	locs  map[string][]ColorId
	order []string
}

// AsMapObject exposes the underlying map as a JsonObject for LogsGson serialization.
func (inv *InverseColorMismatchLocations) AsMapObject() any {
	obj := jutil.NewJsonObject()
	for _, loc := range inv.order {
		arr := jutil.NewJsonArray()
		for _, id := range inv.locs[loc] {
			arr.Add(jutil.NewJsonPrimitiveString(id.ToLogsGson()))
		}
		obj.Add(loc, arr)
	}
	return obj
}

// Inverse returns a location-to-ColorIds index preserving insertion order.
func (cml *ColorMismatchLocations) Inverse() *InverseColorMismatchLocations {
	if len(cml.order) == 0 {
		return &InverseColorMismatchLocations{}
	}
	locs := make(map[string][]ColorId, len(cml.order))
	order := make([]string, 0, len(cml.order))
	for _, id := range cml.order {
		for _, loc := range cml.locs[id] {
			ids, ok := locs[loc]
			if !ok {
				order = append(order, loc)
			} else if slices.Contains(ids, id) {
				continue
			}
			locs[loc] = append(ids, id)
		}
	}
	return &InverseColorMismatchLocations{locs: locs, order: order}
}

// GetMismatchLocationsForDebugging returns an index (mismatch sourceref
// string) => (involved color ids).
//
// This index is only intended for debugging. It may be empty or incomplete
// during production compilations.
//
// Java: ColorRegistry.getMismatchLocationsForDebugging()
func (r *ColorRegistry) GetMismatchLocationsForDebugging() *ColorMismatchLocations {
	if len(r.mismatchOrder) == 0 {
		return &ColorMismatchLocations{}
	}
	locs := cloneMismatchLocations(r.mismatchLocations)
	order := slices.Clone(r.mismatchOrder)
	size := 0
	for _, id := range order {
		size += len(locs[id])
	}
	return &ColorMismatchLocations{
		locs:  locs,
		order: order,
		size:  size,
	}
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
		if existing.equalsColor(supertype) {
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

// ColorIdSet is an immutable insertion-ordered set of ColorIds.
type ColorIdSet []ColorId

// Contains reports whether id is in s.
func (s ColorIdSet) Contains(id ColorId) bool { return slices.Contains(s, id) }

// Size returns the number of ColorIds in s.
func (s ColorIdSet) Size() int { return len(s) }

// All iterates ColorIds in insertion order.
func (s ColorIdSet) All() iter.Seq[ColorId] { return slices.Values([]ColorId(s)) }

// Java: ColorRegistry.REQUIRED_IDS
var ColorRegistry_REQUIRED_IDS = ColorIdSet{
	StandardColors_BIGINT_OBJECT_ID,
	StandardColors_BOOLEAN_OBJECT_ID,
	StandardColors_NUMBER_OBJECT_ID,
	StandardColors_STRING_OBJECT_ID,
	StandardColors_SYMBOL_OBJECT_ID,
	StandardColors_ARRAY_ID,
	StandardColors_READONLY_ARRAY_ID,
	StandardColors_ARGUMENTS_ID,
	StandardColors_ASYNC_ITERATOR_ITERABLE_ID,
	StandardColors_GENERATOR_ID,
	StandardColors_I_TEMPLATE_ARRAY_ID,
	StandardColors_ITERATOR_ID,
	StandardColors_PROMISE_ID,
}
