// Port of com/google/javascript/jscomp/colors/Color.java @ closure-compiler v20261015
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
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// Color is a simplified version of a Closure or TS type for use by
// optimizations.
//
// Go-only: language-adapter — Color is an @AutoValue: the abstract class,
// the generated $AutoValue_Color (fields, equals, hashCode, toString,
// Builder) and AutoValue_Color (the @Memoized subtractNullOrVoid) are one
// struct here (PORTING.md §4, AutoValue → plain struct). Sets of Colors are
// ObjectLinkedHashSets because Color has value equality.
type Color struct {
	id             ColorId
	prototypes     *jutil.ObjectLinkedHashSet[*Color]
	instanceColors *jutil.ObjectLinkedHashSet[*Color]
	ownProperties  *jutil.LinkedHashSet[string]
	boxId          *ColorId
	unionElements  *jutil.ObjectLinkedHashSet[*Color]

	// Memoized value of subtractNullOrVoid() (AutoValue_Color).
	subtractNullOrVoid atomic.Pointer[Color]

	// Go-only: performance — the memoized hashCode(). Java's is cheap
	// because every ImmutableSet field caches its own hash code (and every
	// String its hash); the Go sets recompute theirs on each call, which made
	// hashing a Color with many own properties dominate the hash lookups of
	// DisambiguateProperties. All fields are immutable once built, so the
	// value never changes. Bit 32 indicates the hash has been computed.
	hashCode atomic.Uint64

	invalidating               bool
	propertiesKeepOriginalName bool
	constructor                bool
	closureAssert              bool
}

// Java: Color.getId()
func (c *Color) GetId() ColorId { return c.id }

// GetPrototypes returns, given `function Foo() {}` or `class Foo {}`, the
// color of Foo.prototype. null otherwise.
//
// Java: Color.getPrototypes()
func (c *Color) GetPrototypes() *jutil.ObjectLinkedHashSet[*Color] { return c.prototypes }

// Java: Color.getInstanceColors()
func (c *Color) GetInstanceColors() *jutil.ObjectLinkedHashSet[*Color] { return c.instanceColors }

// Java: Color.isInvalidating()
func (c *Color) IsInvalidating() bool { return c.invalidating }

// Java: Color.getPropertiesKeepOriginalName()
func (c *Color) GetPropertiesKeepOriginalName() bool { return c.propertiesKeepOriginalName }

// Java: Color.isConstructor()
func (c *Color) IsConstructor() bool { return c.constructor }

// GetOwnProperties returns the property names 'declared' on an object (as
// opposed to being conceptually inherited from some supertype).
//
// Java: Color.getOwnProperties()
func (c *Color) GetOwnProperties() *jutil.LinkedHashSet[string] { return c.ownProperties }

// Java: Color.getBoxId()
func (c *Color) GetBoxId() *ColorId { return c.boxId }

// IsClosureAssert reports whether this type is some Closure assertion
// function removable by Closure-specific optimizations.
//
// Java: Color.isClosureAssert()
func (c *Color) IsClosureAssert() bool { return c.closureAssert }

// Java: Color.getUnionElements()
func (c *Color) GetUnionElements() *jutil.ObjectLinkedHashSet[*Color] { return c.unionElements }

var (
	colorEmptyColorSet  = jutil.NewObjectLinkedHashSet[*Color]()
	colorEmptyStringSet = jutil.NewLinkedHashSet[string]()
)

func colorImmutableSetOf() *jutil.ObjectLinkedHashSet[*Color] {
	return colorEmptyColorSet
}

// Java: Color.singleBuilder()
func Color_SingleBuilder() *Color_Builder {
	return NewColor_Builder().
		SetClosureAssert(false).
		SetConstructor(false).
		SetInstanceColors(colorEmptyColorSet).
		SetInvalidating(false).
		SetOwnProperties(colorEmptyStringSet).
		SetPropertiesKeepOriginalName(false).
		SetPrototypes(colorEmptyColorSet).
		SetUnionElements(colorEmptyColorSet)
}

// Java: Color.createUnion(Set<Color>)
func Color_CreateUnion(elements *jutil.ObjectLinkedHashSet[*Color]) *Color {
	var builder Color_Builder
	return Color_CreateUnionWithBuilder(elements, &builder, nil)
}

// Color_CreateUnionWithBuilder builds a union Color reusing builder (which may
// carry a Color slab) and optional allocSet for the flattened unionElements set.
func Color_CreateUnionWithBuilder(
	elements *jutil.ObjectLinkedHashSet[*Color],
	builder *Color_Builder,
	allocSet func(capHint int) *jutil.ObjectLinkedHashSet[*Color],
) *Color {
	elemSize := elements.Size()
	switch elemSize {
	case 0:
		panic(jutil.NewIllegalStateException())
	case 1:
		for e := range elements.All() {
			return e
		}
	default:
	}

	var instanceColors *jutil.ObjectLinkedHashSet[*Color]
	var singleInstanceColors *jutil.ObjectLinkedHashSet[*Color]
	var prototypes *jutil.ObjectLinkedHashSet[*Color]
	var singlePrototypes *jutil.ObjectLinkedHashSet[*Color]
	var ownProperties *jutil.LinkedHashSet[string]
	var singleOwnProps *jutil.LinkedHashSet[string]

	var newElements *jutil.ObjectLinkedHashSet[*Color]
	if allocSet != nil {
		newElements = allocSet(elemSize)
	} else {
		newElements = jutil.NewObjectLinkedHashSet[*Color]()
	}
	var idsBuf [8]jutil.LhmSlot[ColorId, struct{}]
	var ids jutil.LinkedHashSet[ColorId]
	ids.InitWithBuffer(idsBuf[:])

	isClosureAssert := true
	isConstructor := true
	isInvalidating := false
	propertiesKeepOriginalName := false

	for element := range elements.All() {
		if element.IsUnion() {
			for nestedElement := range element.GetUnionElements().All() {
				newElements.Add(nestedElement)
				ids.Add(nestedElement.GetId())
			}
		} else {
			newElements.Add(element)
			ids.Add(element.GetId())
		}

		if ic := element.GetInstanceColors(); !ic.IsEmpty() {
			if instanceColors != nil {
				instanceColors.AddAll(ic.All())
			} else if singleInstanceColors == nil || singleInstanceColors == ic {
				singleInstanceColors = ic
			} else {
				if allocSet != nil {
					instanceColors = allocSet(singleInstanceColors.Size() + ic.Size())
				} else {
					instanceColors = jutil.NewObjectLinkedHashSet[*Color]()
				}
				instanceColors.AddAll(singleInstanceColors.All())
				instanceColors.AddAll(ic.All())
				singleInstanceColors = nil
			}
		}

		isClosureAssert = isClosureAssert && element.IsClosureAssert()
		isConstructor = isConstructor && element.IsConstructor()
		isInvalidating = isInvalidating || element.IsInvalidating()

		if op := element.GetOwnProperties(); !op.IsEmpty() {
			if ownProperties != nil {
				for p := range op.All() {
					ownProperties.Add(p)
				}
			} else if singleOwnProps == nil || singleOwnProps == op {
				singleOwnProps = op
			} else {
				ownProperties = jutil.NewLinkedHashSet[string]()
				for p := range singleOwnProps.All() {
					ownProperties.Add(p)
				}
				for p := range op.All() {
					ownProperties.Add(p)
				}
				singleOwnProps = nil
			}
		}

		propertiesKeepOriginalName = propertiesKeepOriginalName || element.GetPropertiesKeepOriginalName()

		if pr := element.GetPrototypes(); !pr.IsEmpty() {
			if prototypes != nil {
				prototypes.AddAll(pr.All())
			} else if singlePrototypes == nil || singlePrototypes == pr {
				singlePrototypes = pr
			} else {
				if allocSet != nil {
					prototypes = allocSet(singlePrototypes.Size() + pr.Size())
				} else {
					prototypes = jutil.NewObjectLinkedHashSet[*Color]()
				}
				prototypes.AddAll(singlePrototypes.All())
				prototypes.AddAll(pr.All())
				singlePrototypes = nil
			}
		}
	}

	if instanceColors == nil {
		if singleInstanceColors != nil {
			instanceColors = singleInstanceColors
		} else {
			instanceColors = colorEmptyColorSet
		}
	}
	if prototypes == nil {
		if singlePrototypes != nil {
			prototypes = singlePrototypes
		} else {
			prototypes = colorEmptyColorSet
		}
	}
	if ownProperties == nil {
		if singleOwnProps != nil {
			ownProperties = singleOwnProps
		} else {
			ownProperties = colorEmptyStringSet
		}
	}

	return builder.
		SetBoxId(nil).
		SetClosureAssert(isClosureAssert).
		SetConstructor(isConstructor).
		SetId(ColorId_Union(&ids)).
		SetInstanceColors(instanceColors).
		SetInvalidating(isInvalidating).
		SetOwnProperties(ownProperties).
		SetPropertiesKeepOriginalName(propertiesKeepOriginalName).
		SetPrototypes(prototypes).
		SetUnionElements(newElements).
		BuildUnion()
}

// newColor is the AutoValue constructor; the generated null checks of its
// reference fields are Builder.buildInternal's (the only caller).
//
// Java: $AutoValue_Color(ColorId, ImmutableSet<Color>, ImmutableSet<Color>, boolean, boolean, boolean, ImmutableSet<String>, ColorId, boolean, ImmutableSet<Color>)
func newColor(
	slab *jutil.Slab[Color],
	id ColorId,
	prototypes *jutil.ObjectLinkedHashSet[*Color],
	instanceColors *jutil.ObjectLinkedHashSet[*Color],
	invalidating bool,
	propertiesKeepOriginalName bool,
	constructor bool,
	ownProperties *jutil.LinkedHashSet[string],
	boxId *ColorId,
	closureAssert bool,
	unionElements *jutil.ObjectLinkedHashSet[*Color],
) *Color {
	if prototypes == nil {
		panic(jutil.NewNullPointerExceptionMsg("Null prototypes"))
	}
	if instanceColors == nil {
		panic(jutil.NewNullPointerExceptionMsg("Null instanceColors"))
	}
	if ownProperties == nil {
		panic(jutil.NewNullPointerExceptionMsg("Null ownProperties"))
	}
	if unionElements == nil {
		panic(jutil.NewNullPointerExceptionMsg("Null unionElements"))
	}
	var c *Color
	if slab != nil {
		c = slab.New()
	} else {
		c = &Color{}
	}
	c.id = id
	c.prototypes = prototypes
	c.instanceColors = instanceColors
	c.invalidating = invalidating
	c.propertiesKeepOriginalName = propertiesKeepOriginalName
	c.constructor = constructor
	c.ownProperties = ownProperties
	c.boxId = boxId
	c.closureAssert = closureAssert
	c.unionElements = unionElements
	return c
}

// IsPrimitive reports whether this corresponds to a single JavaScript
// primitive like number or symbol.
//
// Note that the boxed versions of primitives (String, Number, etc.) are
// /not/ considered "primitive" by this method.
//
// Java: Color.isPrimitive()
func (c *Color) IsPrimitive() bool {
	/*
	 * Avoid the design headache about whether unions are primitives. The union *color* isn't
	 * primitive, but the *value* held by a union reference may be.
	 */
	jutil.CheckStateMsg1(!c.IsUnion(), "%s", c)
	return StandardColors_PRIMITIVE_COLORS.ContainsKey(c.GetId())
}

// Java: Color.isUnion()
func (c *Color) IsUnion() bool {
	// Single element sets are banned in the builder.
	return !c.GetUnionElements().IsEmpty()
}

// SubtractNullOrVoid is memoized (AutoValue_Color.subtractNullOrVoid).
//
// Java: Color.subtractNullOrVoid()
func (c *Color) SubtractNullOrVoid() *Color {
	if v := c.subtractNullOrVoid.Load(); v != nil {
		return v
	}
	v := c.subtractNullOrVoidImpl()
	if v == nil {
		panic(jutil.NewNullPointerExceptionMsg("subtractNullOrVoid() cannot return null"))
	}
	if c.subtractNullOrVoid.CompareAndSwap(nil, v) {
		return v
	}
	return c.subtractNullOrVoid.Load()
}

// subtractNullOrVoidImpl is the body of Color.subtractNullOrVoid() that
// AutoValue_Color memoizes (its super.subtractNullOrVoid()).
func (c *Color) subtractNullOrVoidImpl() *Color {
	/* Avoid defining NULL_OR_VOID.subtract(NULL_OR_VOID) */
	jutil.CheckState(c.IsUnion())

	if !c.GetUnionElements().Contains(StandardColors_NULL_OR_VOID) {
		return c
	}

	elements := jutil.NewObjectLinkedHashSet(c.GetUnionElements().ToSlice()...)
	elements.Remove(StandardColors_NULL_OR_VOID)
	return Color_CreateUnion(elements)
}

// Java: $AutoValue_Color.toString()
func (c *Color) String() string {
	var sb strings.Builder
	sb.WriteString("Color{")
	sb.WriteString("id=" + c.id.String() + ", ")
	sb.WriteString("prototypes=" + c.prototypes.String() + ", ")
	sb.WriteString("instanceColors=" + c.instanceColors.String() + ", ")
	sb.WriteString("invalidating=" + jutil.ToJavaString(c.invalidating) + ", ")
	sb.WriteString("propertiesKeepOriginalName=" + jutil.ToJavaString(c.propertiesKeepOriginalName) + ", ")
	sb.WriteString("constructor=" + jutil.ToJavaString(c.constructor) + ", ")
	sb.WriteString("ownProperties=" + c.ownProperties.String() + ", ")
	if c.boxId == nil {
		sb.WriteString("boxId=null, ")
	} else {
		sb.WriteString("boxId=" + c.boxId.String() + ", ")
	}
	sb.WriteString("closureAssert=" + jutil.ToJavaString(c.closureAssert) + ", ")
	sb.WriteString("unionElements=" + c.unionElements.String())
	sb.WriteString("}")
	return sb.String()
}

// Java: $AutoValue_Color.equals(Object)
func (c *Color) Equals(o any) bool {
	if that, ok := o.(*Color); ok && that == c {
		return true
	}
	if that, ok := o.(*Color); ok && that != nil {
		return c.id.Equals(that.GetId()) &&
			c.prototypes.Equals(that.GetPrototypes()) &&
			c.instanceColors.Equals(that.GetInstanceColors()) &&
			c.invalidating == that.IsInvalidating() &&
			c.propertiesKeepOriginalName == that.GetPropertiesKeepOriginalName() &&
			c.constructor == that.IsConstructor() &&
			c.ownProperties.Equals(that.GetOwnProperties()) &&
			(func() bool {
				if c.boxId == nil {
					return that.GetBoxId() == nil
				}
				return that.GetBoxId() != nil && c.boxId.Equals(*that.GetBoxId())
			})() &&
			c.closureAssert == that.IsClosureAssert() &&
			c.unionElements.Equals(that.GetUnionElements())
	}
	return false
}

func colorBoolHash(b bool) int32 {
	if b {
		return 1231
	}
	return 1237
}

// Java: $AutoValue_Color.hashCode()
func (c *Color) HashCode() int32 {
	if v := c.hashCode.Load(); v != 0 {
		return int32(uint32(v))
	}
	h := c.computeHashCode()
	c.hashCode.Store((1 << 32) | uint64(uint32(h)))
	return h
}

// computeHashCode is the body of $AutoValue_Color.hashCode(), which
// HashCode memoizes (see Color.hashCode).
func (c *Color) computeHashCode() int32 {
	var h int32 = 1
	h *= 1000003
	h ^= c.id.HashCode()
	h *= 1000003
	h ^= c.prototypes.HashCode()
	h *= 1000003
	h ^= c.instanceColors.HashCode()
	h *= 1000003
	h ^= colorBoolHash(c.invalidating)
	h *= 1000003
	h ^= colorBoolHash(c.propertiesKeepOriginalName)
	h *= 1000003
	h ^= colorBoolHash(c.constructor)
	h *= 1000003
	h ^= c.ownProperties.HashCode()
	h *= 1000003
	if c.boxId == nil {
		h ^= 0
	} else {
		h ^= c.boxId.HashCode()
	}
	h *= 1000003
	h ^= colorBoolHash(c.closureAssert)
	h *= 1000003
	h ^= c.unionElements.HashCode()
	return h
}

// Color_Builder is the builder for a singleton color. Should be passed to
// Color#createSingleton(SingletonColorFields) after building and before
// using.
//
// Go-only: language-adapter — the abstract Color.Builder and the generated
// $AutoValue_Color.Builder are one struct.
type Color_Builder struct {
	id                         ColorId
	prototypes                 *jutil.ObjectLinkedHashSet[*Color]
	instanceColors             *jutil.ObjectLinkedHashSet[*Color]
	invalidating               bool
	propertiesKeepOriginalName bool
	constructor                bool
	ownProperties              *jutil.LinkedHashSet[string]
	boxId                      *ColorId
	closureAssert              bool
	unionElements              *jutil.ObjectLinkedHashSet[*Color]
	slab                       *jutil.Slab[Color]
	set_0                      byte
}

// Java: $AutoValue_Color.Builder()
func NewColor_Builder() *Color_Builder { return &Color_Builder{} }

// SetSlab attaches an optional Color slab to b.
func (b *Color_Builder) SetSlab(slab *jutil.Slab[Color]) *Color_Builder {
	b.slab = slab
	return b
}

// SetId sets the id. (The Java null check has no Go counterpart: ColorId is
// a value type.)
//
// Java: Color.Builder.setId(ColorId)
func (b *Color_Builder) SetId(id ColorId) *Color_Builder {
	b.id = id
	b.set_0 |= 16
	return b
}

// Java: Color.Builder.setInvalidating(boolean)
func (b *Color_Builder) SetInvalidating(invalidating bool) *Color_Builder {
	b.invalidating = invalidating
	b.set_0 |= 1
	return b
}

// Java: Color.Builder.setPropertiesKeepOriginalName(boolean)
func (b *Color_Builder) SetPropertiesKeepOriginalName(propertiesKeepOriginalName bool) *Color_Builder {
	b.propertiesKeepOriginalName = propertiesKeepOriginalName
	b.set_0 |= 2
	return b
}

// Java: Color.Builder.setConstructor(boolean)
func (b *Color_Builder) SetConstructor(constructor bool) *Color_Builder {
	b.constructor = constructor
	b.set_0 |= 4
	return b
}

// Java: Color.Builder.setOwnProperties(ImmutableSet<String>)
func (b *Color_Builder) SetOwnProperties(ownProperties *jutil.LinkedHashSet[string]) *Color_Builder {
	if ownProperties == nil {
		panic(jutil.NewNullPointerExceptionMsg("Null ownProperties"))
	}
	b.ownProperties = ownProperties
	return b
}

// Java: Color.Builder.setClosureAssert(boolean)
func (b *Color_Builder) SetClosureAssert(closureAssert bool) *Color_Builder {
	b.closureAssert = closureAssert
	b.set_0 |= 8
	return b
}

// Java: Color.Builder.setInstanceColors(ImmutableSet<Color>)
func (b *Color_Builder) SetInstanceColors(instanceColors *jutil.ObjectLinkedHashSet[*Color]) *Color_Builder {
	if instanceColors == nil {
		panic(jutil.NewNullPointerExceptionMsg("Null instanceColors"))
	}
	b.instanceColors = instanceColors
	return b
}

// Java: Color.Builder.setPrototypes(ImmutableSet<Color>)
func (b *Color_Builder) SetPrototypes(prototypes *jutil.ObjectLinkedHashSet[*Color]) *Color_Builder {
	if prototypes == nil {
		panic(jutil.NewNullPointerExceptionMsg("Null prototypes"))
	}
	b.prototypes = prototypes
	return b
}

// Java: Color.Builder.setUnionElements(ImmutableSet<Color>)
func (b *Color_Builder) SetUnionElements(unionElements *jutil.ObjectLinkedHashSet[*Color]) *Color_Builder {
	if unionElements == nil {
		panic(jutil.NewNullPointerExceptionMsg("Null unionElements"))
	}
	b.unionElements = unionElements
	return b
}

// Java: Color.Builder.setBoxId(ColorId)
func (b *Color_Builder) SetBoxId(boxId *ColorId) *Color_Builder {
	b.boxId = boxId
	return b
}

// Java: Color.Builder.setPrototype(Color)
func (b *Color_Builder) SetPrototype(x *Color) *Color_Builder {
	if x == nil {
		return b.SetPrototypes(colorImmutableSetOf())
	}
	return b.SetPrototypes(jutil.NewObjectLinkedHashSet(x))
}

// Java: Color.Builder.setInstanceColor(Color)
func (b *Color_Builder) SetInstanceColor(x *Color) *Color_Builder {
	if x == nil {
		return b.SetInstanceColors(colorImmutableSetOf())
	}
	return b.SetInstanceColors(jutil.NewObjectLinkedHashSet(x))
}

// Java: $AutoValue_Color.Builder.buildInternal()
func (b *Color_Builder) BuildInternal() *Color {
	if b.set_0 != 0x1f ||
		b.prototypes == nil ||
		b.instanceColors == nil ||
		b.ownProperties == nil ||
		b.unionElements == nil {
		var missing strings.Builder
		if (b.set_0 & 16) == 0 {
			missing.WriteString(" id")
		}
		if b.prototypes == nil {
			missing.WriteString(" prototypes")
		}
		if b.instanceColors == nil {
			missing.WriteString(" instanceColors")
		}
		if (b.set_0 & 1) == 0 {
			missing.WriteString(" invalidating")
		}
		if (b.set_0 & 2) == 0 {
			missing.WriteString(" propertiesKeepOriginalName")
		}
		if (b.set_0 & 4) == 0 {
			missing.WriteString(" constructor")
		}
		if b.ownProperties == nil {
			missing.WriteString(" ownProperties")
		}
		if (b.set_0 & 8) == 0 {
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
		b.closureAssert,
		b.unionElements)
}

// Java: Color.Builder.build()
func (b *Color_Builder) Build() *Color {
	result := b.BuildInternal()
	jutil.CheckStateMsg1(result.GetUnionElements().IsEmpty(), "%s", result)
	jutil.CheckStateMsg1(result.GetBoxId() == nil, "%s", result)
	jutil.CheckStateMsg1(!StandardColors_AXIOMATIC_COLORS.ContainsKey(result.GetId()), "%s", result)
	return result
}

// Java: Color.Builder.buildUnion()
func (b *Color_Builder) BuildUnion() *Color {
	result := b.BuildInternal()
	jutil.CheckStateMsg1(result.GetUnionElements().Size() > 1, "%s", result)
	jutil.CheckStateMsg1(result.GetBoxId() == nil, "%s", result)
	return result
}

// Java: Color.Builder.buildAxiomatic()
func (b *Color_Builder) BuildAxiomatic() *Color {
	jutil.CheckStateMsg(StandardColors_AXIOMATIC_COLORS == nil, "StandardColors are all defined")

	result := b.BuildInternal()
	jutil.CheckStateMsg1(result.GetUnionElements().IsEmpty(), "%s", result)
	return result
}
