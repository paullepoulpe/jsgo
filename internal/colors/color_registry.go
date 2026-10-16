// Port of com/google/javascript/jscomp/colors/ColorRegistry.java @ closure-compiler v20260415

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

// Java: ColorRegistry(Builder)
func NewColorRegistry(builder *ColorRegistry_Builder) *ColorRegistry {
	r := &ColorRegistry{}
	r.nativeColors = jutil.NewLinkedHashMapFrom(builder.nativeColors)
	r.colorToDisambiguationSupertypeGraph = jutil.SetMultimapCopyOf(builder.colorToDisambiguationSupertypeGraph)
	r.mismatchLocations = jutil.SetMultimapCopyOf(builder.mismatchLocations)

	jutil.CheckState(jutil.NewLinkedHashSet(r.nativeColors.KeySlice()...).Equals(ColorRegistry_REQUIRED_IDS))
	return r
}

// Java: ColorRegistry.get(ColorId)
func (r *ColorRegistry) Get(id ColorId) *Color {
	// checkNotNull(id): ColorId is a Go value type and cannot be null.
	return jutil.CheckNotNullMsg(r.nativeColors.Get(id), "Missing Color for %s", id)
}

// GetDisambiguationSupertypes returns the colors directly above `x` in the
// subtyping graph for the purposes of property (dis)ambiguation.
//
// Java: ColorRegistry.getDisambiguationSupertypes(Color)
func (r *ColorRegistry) GetDisambiguationSupertypes(x *Color) *jutil.ObjectLinkedHashSet[*Color] {
	return r.colorToDisambiguationSupertypeGraph.Get(x)
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
}

// Java: ColorRegistry.Builder()
func NewColorRegistry_Builder() *ColorRegistry_Builder {
	return &ColorRegistry_Builder{
		nativeColors:                        jutil.NewLinkedHashMap[ColorId, *Color](),
		colorToDisambiguationSupertypeGraph: jutil.NewSetMultimap[*Color, *Color](),
		mismatchLocations:                   jutil.NewSetMultimap[ColorId, string](),
	}
}

// Java: ColorRegistry.Builder.setNativeColor(Color)
func (b *ColorRegistry_Builder) SetNativeColor(x *Color) *ColorRegistry_Builder {
	jutil.CheckStateMsg(ColorRegistry_REQUIRED_IDS.Contains(x.GetId()), "%s", x)
	b.nativeColors.Put(x.GetId(), x)
	return b
}

// Java: ColorRegistry.Builder.addDisambiguationEdge(Color, Color)
func (b *ColorRegistry_Builder) AddDisambiguationEdge(subtype *Color, supertype *Color) *ColorRegistry_Builder {
	b.colorToDisambiguationSupertypeGraph.Put(subtype, supertype)
	return b
}

// Java: ColorRegistry.Builder.addMismatchLocation(ColorId, String)
func (b *ColorRegistry_Builder) AddMismatchLocation(id ColorId, location string) *ColorRegistry_Builder {
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
