// Port of com/google/javascript/jscomp/colors/ColorRegistry.java @ closure-compiler v20261015

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

import "github.com/paullepoulpe/jsgo/internal/jutil"

// ColorRegistry stores "out of band" information about the set of Colors in
// use by a compilation.
//
// This includes any information that isn't specifically can't be set on a
// single Color object,
type ColorRegistry struct {
	nativeColors                        *jutil.LinkedHashMap[ColorId, *Color]
	colorToDisambiguationSupertypeGraph *jutil.SetMultimap[*Color, *Color]
	mismatchLocations                   *jutil.SetMultimap[ColorId, string]
}

type colorRegistryEdgeSet struct {
	set jutil.ObjectLinkedHashSet[*Color]
	buf [2]jutil.OlhmSlot[*Color, struct{}]
}

// Java: ColorRegistry(Builder)
func NewColorRegistry(builder *ColorRegistry_Builder) *ColorRegistry {
	r := &ColorRegistry{}
	if !builder.consumed {
		builder.consumed = true
		r.nativeColors = builder.nativeColors
		r.colorToDisambiguationSupertypeGraph = builder.colorToDisambiguationSupertypeGraph
		r.mismatchLocations = builder.mismatchLocations
	} else {
		r.nativeColors = jutil.NewLinkedHashMapFrom(builder.nativeColors)
		r.colorToDisambiguationSupertypeGraph = jutil.SetMultimapCopyOf(builder.colorToDisambiguationSupertypeGraph)
		r.mismatchLocations = jutil.SetMultimapCopyOf(builder.mismatchLocations)
	}

	ok := r.nativeColors.Size() == ColorRegistry_REQUIRED_IDS.Size()
	if ok {
		for id := range ColorRegistry_REQUIRED_IDS.All() {
			if !r.nativeColors.ContainsKey(id) {
				ok = false
				break
			}
		}
	}
	jutil.CheckState(ok)
	return r
}

// Java: ColorRegistry.get(ColorId)
func (r *ColorRegistry) Get(id ColorId) *Color {
	// checkNotNull(id): ColorId is a Go value type and cannot be null.
	c := r.nativeColors.Get(id)
	if c == nil {
		panic(jutil.NewNullPointerExceptionMsg("Missing Color for " + id.String()))
	}
	return c
}

// GetDisambiguationSupertypes returns the colors directly above `x` in the
// subtyping graph for the purposes of property (dis)ambiguation.
//
// Java: ColorRegistry.getDisambiguationSupertypes(Color)
func (r *ColorRegistry) GetDisambiguationSupertypes(x *Color) *jutil.ObjectLinkedHashSet[*Color] {
	return r.colorToDisambiguationSupertypeGraph.GetOrEmpty(x, colorEmptyColorSet)
}

// GetMismatchLocationsForDebugging returns an index (mismatch sourceref
// string) => (involved color ids).
//
// This index is only intended for debugging. It may be empty or incomplete
// during production compilations.
//
// Java: ColorRegistry.getMismatchLocationsForDebugging()
func (r *ColorRegistry) GetMismatchLocationsForDebugging() *jutil.SetMultimap[ColorId, string] {
	return r.mismatchLocations
}

// Java: ColorRegistry.builder()
func ColorRegistry_Builder_() *ColorRegistry_Builder {
	return NewColorRegistry_Builder()
}

// ColorRegistry_Builder builds a ColorRegistry.
type ColorRegistry_Builder struct {
	nativeColors                        *jutil.LinkedHashMap[ColorId, *Color]
	colorToDisambiguationSupertypeGraph *jutil.SetMultimap[*Color, *Color]
	mismatchLocations                   *jutil.SetMultimap[ColorId, string]
	edgeSetSlab                         jutil.Slab[colorRegistryEdgeSet]
	allocEdgeSetFn                      func() *jutil.ObjectLinkedHashSet[*Color]
	consumed                            bool
}

// Java: ColorRegistry.Builder()
func NewColorRegistry_Builder() *ColorRegistry_Builder {
	b := &ColorRegistry_Builder{
		nativeColors:                        jutil.NewLinkedHashMap[ColorId, *Color](),
		colorToDisambiguationSupertypeGraph: jutil.NewSetMultimap[*Color, *Color](),
		mismatchLocations:                   jutil.NewSetMultimap[ColorId, string](),
		edgeSetSlab:                         jutil.NewFixedSlab[colorRegistryEdgeSet](256), // 16 KiB <= 32 KiB
	}
	b.allocEdgeSetFn = b.allocEdgeSet
	return b
}

func (b *ColorRegistry_Builder) allocEdgeSet() *jutil.ObjectLinkedHashSet[*Color] {
	es := b.edgeSetSlab.New()
	es.set.InitWithBuffer(es.buf[:])
	return &es.set
}

func (b *ColorRegistry_Builder) ensureMutable() {
	if b.consumed {
		b.nativeColors = jutil.NewLinkedHashMapFrom(b.nativeColors)
		b.colorToDisambiguationSupertypeGraph = jutil.SetMultimapCopyOf(b.colorToDisambiguationSupertypeGraph)
		b.mismatchLocations = jutil.SetMultimapCopyOf(b.mismatchLocations)
		b.consumed = false
	}
}

// Java: ColorRegistry.Builder.setNativeColor(Color)
func (b *ColorRegistry_Builder) SetNativeColor(x *Color) *ColorRegistry_Builder {
	b.ensureMutable()
	jutil.CheckStateMsg1(ColorRegistry_REQUIRED_IDS.Contains(x.GetId()), "%s", x)
	b.nativeColors.Put(x.GetId(), x)
	return b
}

// Java: ColorRegistry.Builder.addDisambiguationEdge(Color, Color)
func (b *ColorRegistry_Builder) AddDisambiguationEdge(subtype *Color, supertype *Color) *ColorRegistry_Builder {
	b.ensureMutable()
	b.colorToDisambiguationSupertypeGraph.PutWithSetAlloc(subtype, supertype, b.allocEdgeSetFn)
	return b
}

// Java: ColorRegistry.Builder.addMismatchLocation(ColorId, String)
func (b *ColorRegistry_Builder) AddMismatchLocation(id ColorId, location string) *ColorRegistry_Builder {
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
