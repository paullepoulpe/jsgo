// Port of com/google/javascript/jscomp/colors/Color.java @ closure-compiler v20261021
// (with the AutoValue-generated AutoValue_Color.java and $AutoValue_Color.java)

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
	"strings"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// Java: Color.singleBuilder()
func Color_SingleBuilder() *Color_Builder {
	return &Color_Builder{
		prototypes:                    emptyColorSlice,
		instanceColors:                emptyColorSlice,
		ownProperties:                 emptyStringSlice,
		unionElements:                 emptyColorSlice,
		hasInvalidating:               true,
		hasPropertiesKeepOriginalName: true,
		hasConstructor:                true,
		hasClosureAssert:              true,
	}
}

func dedupCloneColors(src []*Color) []*Color {
	if len(src) == 0 {
		return emptyColorSlice
	}
	if len(src) == 1 {
		if src[0] == nil {
			panic(jutil.NewNullPointerException())
		}
		return []*Color{src[0]}
	}
	if len(src) <= 16 {
		allUnique := true
		for i, c := range src {
			if c == nil {
				panic(jutil.NewNullPointerException())
			}
			for j := 0; j < i; j++ {
				if src[j] == c || src[j].Equals(c) {
					allUnique = false
					break
				}
			}
		}
		if allUnique {
			out := make([]*Color, len(src))
			copy(out, src)
			return out
		}
	}
	out := make([]*Color, 0, len(src))
	var byHash map[int32][]*Color
	for _, c := range src {
		if c == nil {
			panic(jutil.NewNullPointerException())
		}
		out = appendUniqueColor(out, c, &byHash)
	}
	return out
}

func dedupCloneStrings(src []string) []string {
	if len(src) == 0 {
		return emptyStringSlice
	}
	if len(src) == 1 {
		return []string{src[0]}
	}
	if len(src) <= 16 {
		allUnique := true
		for i, s := range src {
			for j := 0; j < i; j++ {
				if src[j] == s {
					allUnique = false
					break
				}
			}
		}
		if allUnique {
			out := make([]string, len(src))
			copy(out, src)
			return out
		}
	}
	out := make([]string, 0, len(src))
	var seen map[string]struct{}
	for _, s := range src {
		out = appendUniqueString(out, s, &seen)
	}
	return out
}

// Color_Builder is the builder for a singleton color. Should be passed to
// Color#createSingleton(SingletonColorFields) after building and before
// using.
//
// Go-only: language-adapter — the abstract Color.Builder and the generated
// $AutoValue_Color.Builder are one struct.
type Color_Builder struct {
	prototypes                    []*Color
	instanceColors                []*Color
	ownProperties                 []string
	unionElements                 []*Color
	slab                          *jutil.Slab[Color]
	id                            ColorId
	boxId                         ColorId
	hasId                         bool
	hasBoxId                      bool
	invalidating                  bool
	hasInvalidating               bool
	propertiesKeepOriginalName    bool
	hasPropertiesKeepOriginalName bool
	constructor                   bool
	hasConstructor                bool
	closureAssert                 bool
	hasClosureAssert              bool
	trustUnique                   bool
}

// SetSlab attaches an optional Color slab to b.
func (b *Color_Builder) SetSlab(slab *jutil.Slab[Color]) *Color_Builder {
	b.slab = slab
	return b
}

// SetTrustUniqueSlices marks slice inputs passed to Set* as already immutable
// and deduplicated (used by ColorPool_Builder).
func (b *Color_Builder) SetTrustUniqueSlices(trust bool) *Color_Builder {
	b.trustUnique = trust
	return b
}

// SetId sets the id. (The Java null check has no Go counterpart: ColorId is
// a value type.)
//
// Java: Color.Builder.setId(ColorId)
func (b *Color_Builder) SetId(id ColorId) *Color_Builder {
	b.id = id
	b.hasId = true
	return b
}

// Java: Color.Builder.setInvalidating(boolean)
func (b *Color_Builder) SetInvalidating(invalidating bool) *Color_Builder {
	b.invalidating = invalidating
	b.hasInvalidating = true
	return b
}

// Java: Color.Builder.setPropertiesKeepOriginalName(boolean)
func (b *Color_Builder) SetPropertiesKeepOriginalName(propertiesKeepOriginalName bool) *Color_Builder {
	b.propertiesKeepOriginalName = propertiesKeepOriginalName
	b.hasPropertiesKeepOriginalName = true
	return b
}

// Java: Color.Builder.setConstructor(boolean)
func (b *Color_Builder) SetConstructor(constructor bool) *Color_Builder {
	b.constructor = constructor
	b.hasConstructor = true
	return b
}

// Java: Color.Builder.setOwnProperties(ImmutableSet<String>)
func (b *Color_Builder) SetOwnProperties(ownProperties []string) *Color_Builder {
	if ownProperties == nil {
		panic(jutil.NewNullPointerExceptionMsg("Null ownProperties"))
	}
	if len(ownProperties) == 0 {
		b.ownProperties = emptyStringSlice
	} else if b.trustUnique {
		b.ownProperties = ownProperties
	} else {
		b.ownProperties = dedupCloneStrings(ownProperties)
	}
	return b
}

// Java: Color.Builder.setClosureAssert(boolean)
func (b *Color_Builder) SetClosureAssert(closureAssert bool) *Color_Builder {
	b.closureAssert = closureAssert
	b.hasClosureAssert = true
	return b
}

// Java: Color.Builder.setInstanceColors(ImmutableSet<Color>)
func (b *Color_Builder) SetInstanceColors(instanceColors []*Color) *Color_Builder {
	if instanceColors == nil {
		panic(jutil.NewNullPointerExceptionMsg("Null instanceColors"))
	}
	if len(instanceColors) == 0 {
		b.instanceColors = emptyColorSlice
	} else if b.trustUnique {
		b.instanceColors = instanceColors
	} else {
		b.instanceColors = dedupCloneColors(instanceColors)
	}
	return b
}

// Java: Color.Builder.setPrototypes(ImmutableSet<Color>)
func (b *Color_Builder) SetPrototypes(prototypes []*Color) *Color_Builder {
	if prototypes == nil {
		panic(jutil.NewNullPointerExceptionMsg("Null prototypes"))
	}
	if len(prototypes) == 0 {
		b.prototypes = emptyColorSlice
	} else if b.trustUnique {
		b.prototypes = prototypes
	} else {
		b.prototypes = dedupCloneColors(prototypes)
	}
	return b
}

// Java: Color.Builder.setUnionElements(ImmutableSet<Color>)
func (b *Color_Builder) SetUnionElements(unionElements []*Color) *Color_Builder {
	if unionElements == nil {
		panic(jutil.NewNullPointerExceptionMsg("Null unionElements"))
	}
	if len(unionElements) == 0 {
		b.unionElements = emptyColorSlice
	} else if b.trustUnique {
		b.unionElements = unionElements
	} else {
		b.unionElements = dedupCloneColors(unionElements)
	}
	return b
}

// SetBoxIdValue sets boxId directly by value without allocating a *ColorId.
func (b *Color_Builder) SetBoxIdValue(boxId ColorId) *Color_Builder {
	b.boxId = boxId
	b.hasBoxId = true
	return b
}

// ClearBoxId clears boxId.
func (b *Color_Builder) ClearBoxId() *Color_Builder {
	b.boxId = ColorId{}
	b.hasBoxId = false
	return b
}

// Java: Color.Builder.setPrototype(Color)
func (b *Color_Builder) SetPrototype(x *Color) *Color_Builder {
	if x == nil {
		b.prototypes = emptyColorSlice
		return b
	}
	b.prototypes = []*Color{x}
	return b
}

// Java: Color.Builder.setInstanceColor(Color)
func (b *Color_Builder) SetInstanceColor(x *Color) *Color_Builder {
	if x == nil {
		b.instanceColors = emptyColorSlice
		return b
	}
	b.instanceColors = []*Color{x}
	return b
}

// Java: $AutoValue_Color.Builder.buildInternal()
func (b *Color_Builder) BuildInternal() *Color {
	if !b.hasId ||
		!b.hasInvalidating ||
		!b.hasPropertiesKeepOriginalName ||
		!b.hasConstructor ||
		!b.hasClosureAssert ||
		b.prototypes == nil ||
		b.instanceColors == nil ||
		b.ownProperties == nil ||
		b.unionElements == nil {
		var missing strings.Builder
		if !b.hasId {
			missing.WriteString(" id")
		}
		if b.prototypes == nil {
			missing.WriteString(" prototypes")
		}
		if b.instanceColors == nil {
			missing.WriteString(" instanceColors")
		}
		if !b.hasInvalidating {
			missing.WriteString(" invalidating")
		}
		if !b.hasPropertiesKeepOriginalName {
			missing.WriteString(" propertiesKeepOriginalName")
		}
		if !b.hasConstructor {
			missing.WriteString(" constructor")
		}
		if b.ownProperties == nil {
			missing.WriteString(" ownProperties")
		}
		if !b.hasClosureAssert {
			missing.WriteString(" closureAssert")
		}
		if b.unionElements == nil {
			missing.WriteString(" unionElements")
		}
		panic(jutil.NewIllegalStateExceptionMsg("Missing required properties:" + missing.String()))
	}
	return newColor(
		b.slab,
		b.id,
		b.prototypes,
		b.instanceColors,
		b.invalidating,
		b.propertiesKeepOriginalName,
		b.constructor,
		b.ownProperties,
		b.boxId,
		b.hasBoxId,
		b.closureAssert,
		b.unionElements)
}

// Java: Color.Builder.build()
func (b *Color_Builder) Build() *Color {
	result := b.BuildInternal()
	jutil.CheckStateMsg1(len(result.unionElements) == 0, "%s", result)
	jutil.CheckStateMsg1(!result.hasBoxId, "%s", result)
	jutil.CheckStateMsg1(StandardColors_AXIOMATIC_COLORS[result.GetId()] == nil, "%s", result)
	return result
}

// Java: Color.Builder.buildUnion()
func (b *Color_Builder) BuildUnion() *Color {
	result := b.BuildInternal()
	jutil.CheckStateMsg1(len(result.unionElements) > 1, "%s", result)
	jutil.CheckStateMsg1(!result.hasBoxId, "%s", result)
	return result
}

// Java: Color.Builder.buildAxiomatic()
func (b *Color_Builder) BuildAxiomatic() *Color {
	jutil.CheckStateMsg(StandardColors_AXIOMATIC_COLORS == nil, "StandardColors are all defined")

	result := b.BuildInternal()
	jutil.CheckStateMsg1(len(result.unionElements) == 0, "%s", result)
	return result
}
