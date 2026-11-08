// Port of com/google/javascript/jscomp/colors/ColorRegistry.java @ closure-compiler v20261104

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
	disambiguationByHash     map[int32][]*Color
	mismatchLocations        jutil.OrderedSetMultimap[ColorId, string]
}

// Java: ColorRegistry(Builder)
func NewColorRegistry(builder *ColorRegistryBuilder) *ColorRegistry {
	r := &ColorRegistry{}
	if !builder.consumed {
		builder.consumed = true
		r.nativeColors = builder.nativeColors
		r.disambiguationSupertypes = builder.disambiguationSupertypes
		r.disambiguationByHash = builder.disambiguationByHash
		r.mismatchLocations = builder.mismatchLocations
	} else {
		r.nativeColors = maps.Clone(builder.nativeColors)
		r.disambiguationSupertypes = cloneSliceMap(builder.disambiguationSupertypes)
		r.disambiguationByHash = cloneSliceMap(builder.disambiguationByHash)
		r.mismatchLocations = builder.mismatchLocations.Clone()
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

func cloneSliceMap[K comparable, V any](src map[K][]V) map[K][]V {
	if len(src) == 0 {
		return nil
	}
	dst := make(map[K][]V, len(src))
	for k, vs := range src { // order-independent: copying distinct keys into a new map
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
	if sups, ok := r.disambiguationSupertypes[x]; ok {
		return sups
	}
	for _, k := range r.disambiguationByHash[x.HashCode()] {
		if k.equalsColor(x) {
			return r.disambiguationSupertypes[k]
		}
	}
	return nil
}

// ColorMismatchLocSet is an insertion-ordered, deduplicated view of mismatch location strings.
type ColorMismatchLocSet = jutil.OrderedSetView[string]

// ColorMismatchLocations maps ColorIds to mismatch location sourceref strings in insertion order.
type ColorMismatchLocations struct {
	jutil.OrderedSetMultimap[ColorId, string]
}

// NewColorMismatchLocations returns an empty ColorMismatchLocations.
func NewColorMismatchLocations() *ColorMismatchLocations {
	return &ColorMismatchLocations{}
}

// InverseColorMismatchLocations maps mismatch location strings to involved ColorIds.
type InverseColorMismatchLocations struct {
	jutil.OrderedSetMultimap[string, ColorId]
}

// AsMapObject exposes the underlying map as a JsonObject for LogsGson serialization.
func (inv *InverseColorMismatchLocations) AsMapObject() any {
	obj := jutil.NewJsonObject()
	for loc, ids := range inv.AsMap() {
		arr := jutil.NewJsonArray()
		for _, id := range ids {
			arr.Add(jutil.NewJsonPrimitiveString(id.ToLogsGson()))
		}
		obj.Add(loc, arr)
	}
	return obj
}

// Inverse returns a location-to-ColorIds index preserving insertion order.
func (cml *ColorMismatchLocations) Inverse() *InverseColorMismatchLocations {
	return &InverseColorMismatchLocations{OrderedSetMultimap: *cml.OrderedSetMultimap.Inverse()}
}

// GetMismatchLocationsForDebugging returns an index (mismatch sourceref
// string) => (involved color ids).
//
// This index is only intended for debugging. It may be empty or incomplete
// during production compilations.
//
// Java: ColorRegistry.getMismatchLocationsForDebugging()
func (r *ColorRegistry) GetMismatchLocationsForDebugging() *ColorMismatchLocations {
	return &ColorMismatchLocations{OrderedSetMultimap: r.mismatchLocations.Clone()}
}

// ColorRegistryBuilder builds a ColorRegistry.
type ColorRegistryBuilder struct {
	nativeColors             map[ColorId]*Color
	disambiguationSupertypes map[*Color][]*Color
	disambiguationByHash     map[int32][]*Color
	mismatchLocations        jutil.OrderedSetMultimap[ColorId, string]
	consumed                 bool
}

// NewColorRegistryBuilder returns a new ColorRegistryBuilder.
func NewColorRegistryBuilder() *ColorRegistryBuilder {
	return &ColorRegistryBuilder{
		nativeColors: make(map[ColorId]*Color, ColorRegistry_REQUIRED_IDS.Size()),
	}
}

func (b *ColorRegistryBuilder) ensureMutable() {
	if b.consumed {
		b.nativeColors = maps.Clone(b.nativeColors)
		b.disambiguationSupertypes = cloneSliceMap(b.disambiguationSupertypes)
		b.disambiguationByHash = cloneSliceMap(b.disambiguationByHash)
		b.mismatchLocations = b.mismatchLocations.Clone()
		b.consumed = false
	}
}

// Java: ColorRegistry.Builder.setNativeColor(Color)
func (b *ColorRegistryBuilder) SetNativeColor(x *Color) *ColorRegistryBuilder {
	b.ensureMutable()
	jutil.CheckStateMsg1(ColorRegistry_REQUIRED_IDS.Contains(x.GetId()), "%s", x)
	b.nativeColors[x.GetId()] = x
	return b
}

// Java: ColorRegistry.Builder.addDisambiguationEdge(Color, Color)
func (b *ColorRegistryBuilder) AddDisambiguationEdge(subtype *Color, supertype *Color) *ColorRegistryBuilder {
	jutil.CheckNotNull(subtype)
	jutil.CheckNotNull(supertype)
	b.ensureMutable()
	if b.disambiguationSupertypes == nil {
		b.disambiguationSupertypes = make(map[*Color][]*Color)
		b.disambiguationByHash = make(map[int32][]*Color)
	}
	key := subtype
	if _, ok := b.disambiguationSupertypes[key]; !ok {
		h := subtype.HashCode()
		found := false
		for _, existingKey := range b.disambiguationByHash[h] {
			if existingKey.equalsColor(subtype) {
				key = existingKey
				found = true
				break
			}
		}
		if !found {
			b.disambiguationByHash[h] = append(b.disambiguationByHash[h], key)
		}
	}
	sups := b.disambiguationSupertypes[key]
	for _, existing := range sups {
		if existing.equalsColor(supertype) {
			return b
		}
	}
	b.disambiguationSupertypes[key] = append(sups, supertype)
	return b
}

// Java: ColorRegistry.Builder.addMismatchLocation(ColorId, String)
func (b *ColorRegistryBuilder) AddMismatchLocation(id ColorId, location string) *ColorRegistryBuilder {
	b.ensureMutable()
	b.mismatchLocations.Put(id, location)
	return b
}

// SetDefaultNativeColorsForTesting sets defaults for native Colors.
//
// Only for use in testing. In real compilations, certain native colors have
// fields that vary from compilation-to-compilation (like whether the
// "Number" object is invalidating), so should use setNativeColor instead.
//
// Java: ColorRegistry.Builder.setDefaultNativeColorsForTesting()
func (b *ColorRegistryBuilder) SetDefaultNativeColorsForTesting() *ColorRegistryBuilder {
	for id := range ColorRegistry_REQUIRED_IDS.All() {
		b.SetNativeColor(Color_SingleBuilder().SetId(id).Build())
	}
	return b
}

// Java: ColorRegistry.Builder.build()
func (b *ColorRegistryBuilder) Build() *ColorRegistry {
	return NewColorRegistry(b)
}

// ColorIdSet is an immutable insertion-ordered set of ColorIds.
type ColorIdSet = jutil.OrderedSetView[ColorId]

// Java: ColorRegistry.REQUIRED_IDS
var ColorRegistry_REQUIRED_IDS = StandardColors_STANDARD_OBJECT_IDS
