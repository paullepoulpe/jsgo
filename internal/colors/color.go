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
	"strconv"
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
// struct here (PORTING.md §4, AutoValue → plain struct). Set fields are stored
// as immutable deduplicated slices preserving insertion order, with pointer
// fields grouped first to minimize GC scan span.
type Color struct {
	prototypes     []*Color
	instanceColors []*Color
	ownProperties  []string
	unionElements  []*Color

	// Memoized value of subtractNullOrVoid() (AutoValue_Color).
	subtractNullOrVoid atomic.Pointer[Color]

	id    ColorId
	boxId ColorId

	// Go-only: performance — the memoized hashCode(). All fields are immutable
	// once built, so the value never changes. Bit 32 indicates the hash has
	// been computed.
	hashCode atomic.Uint64

	// Go-only: the color's handle in rhino's node color table (owner id in
	// the high 32 bits, handle in the low 32), set by rhino.Node.SetColor.
	nodeHandle atomic.Uint64

	hasBoxId                   bool
	invalidating               bool
	propertiesKeepOriginalName bool
	constructor                bool
	closureAssert              bool
}

// LoadNodeHandle returns the owner id and handle stored by StoreNodeHandle.
func (c *Color) LoadNodeHandle() (ownerID uint32, handle uint32) {
	v := c.nodeHandle.Load()
	return uint32(v >> 32), uint32(v)
}

// StoreNodeHandle records the handle the node table of ownerID gave c.
func (c *Color) StoreNodeHandle(ownerID uint32, handle uint32) {
	c.nodeHandle.Store(uint64(ownerID)<<32 | uint64(handle))
}

// Java: Color.getId()
func (c *Color) GetId() ColorId { return c.id }

// GetPrototypes returns, given `function Foo() {}` or `class Foo {}`, the
// color of Foo.prototype. Empty otherwise.
//
// Java: Color.getPrototypes()
func (c *Color) GetPrototypes() []*Color { return c.prototypes }

// Java: Color.getInstanceColors()
func (c *Color) GetInstanceColors() []*Color { return c.instanceColors }

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
func (c *Color) GetOwnProperties() []string { return c.ownProperties }

// HasOwnProperty reports whether prop is declared on c.
func (c *Color) HasOwnProperty(prop string) bool {
	for _, p := range c.ownProperties {
		if p == prop {
			return true
		}
	}
	return false
}

// Java: Color.getBoxId()
func (c *Color) GetBoxId() (ColorId, bool) {
	return c.boxId, c.hasBoxId
}

// IsClosureAssert reports whether this type is some Closure assertion
// function removable by Closure-specific optimizations.
//
// Java: Color.isClosureAssert()
func (c *Color) IsClosureAssert() bool { return c.closureAssert }

// Java: Color.getUnionElements()
func (c *Color) GetUnionElements() []*Color { return c.unionElements }

var (
	emptyColorSlice  = []*Color{}
	emptyStringSlice = []string{}
)

// EmptyColorSlice returns the shared non-nil empty []*Color sentinel.
func EmptyColorSlice() []*Color { return emptyColorSlice }

// EmptyStringSlice returns the shared non-nil empty []string sentinel.
func EmptyStringSlice() []string { return emptyStringSlice }

// Java: Color.createUnion(Set<Color>)
func Color_CreateUnion(elements *jutil.ObjectLinkedHashSet[*Color]) *Color {
	return Color_CreateUnionSlice(elements.ToSlice())
}

// Color_CreateUnionSlice builds a union Color from a slice of Colors.
func Color_CreateUnionSlice(elements []*Color) *Color {
	var builder Color_Builder
	return Color_CreateUnionSliceWithBuilder(elements, &builder, nil, nil)
}

// Color_CreateUnion2 builds a union Color from two non-null Colors (matching
// Color.createUnion(ImmutableSet.of(e1, e2))).
func Color_CreateUnion2(e1, e2 *Color) *Color {
	if e1 == nil || e2 == nil {
		panic(jutil.NewNullPointerException())
	}
	if e1 == e2 || e1.Equals(e2) {
		return e1
	}
	buf := [2]*Color{e1, e2}
	var builder Color_Builder
	return Color_CreateUnionSliceWithBuilder(buf[:], &builder, nil, nil)
}

type dedupSet[T comparable] struct {
	seen   map[T]struct{}
	byHash map[int32][]*Color
}

func elemEqual[T comparable](a, b T) bool {
	if a == b {
		return true
	}
	if _, isColor := any((*T)(nil)).(**Color); isColor {
		return any(a).(*Color).Equals(any(b).(*Color))
	}
	return false
}

func (s *dedupSet[T]) initialized() bool {
	return s.seen != nil || s.byHash != nil
}

func (s *dedupSet[T]) init(capHint int) {
	if _, isColor := any((*T)(nil)).(**Color); isColor {
		s.byHash = make(map[int32][]*Color, capHint)
	} else {
		s.seen = make(map[T]struct{}, capHint)
	}
}

func (s *dedupSet[T]) add(v T) bool {
	if _, isColor := any((*T)(nil)).(**Color); isColor {
		c := any(v).(*Color)
		h := c.HashCode()
		bucket := s.byHash[h]
		for _, existing := range bucket {
			if existing == c || existing.Equals(c) {
				return false
			}
		}
		s.byHash[h] = append(bucket, c)
		return true
	}
	if _, ok := s.seen[v]; ok {
		return false
	}
	s.seen[v] = struct{}{}
	return true
}

func appendUnique[T comparable](dst []T, v T, seen *dedupSet[T]) []T {
	if seen.initialized() {
		if seen.add(v) {
			return append(dst, v)
		}
		return dst
	}
	for _, existing := range dst {
		if elemEqual(existing, v) {
			return dst
		}
	}
	dst = append(dst, v)
	if len(dst) > 16 {
		seen.init(len(dst) * 2)
		for _, existing := range dst {
			seen.add(existing)
		}
	}
	return dst
}

func colorIdUnionFromColors(elems []*Color) ColorId {
	var stackBuf [32]int64
	var ids []int64
	if len(elems) <= len(stackBuf) {
		ids = stackBuf[:len(elems)]
	} else {
		ids = make([]int64, len(elems))
	}
	for i, e := range elems {
		ids[i] = e.id.rightAligned
	}
	return unionRightAlignedIds(ids)
}

// Color_CreateUnionSliceWithBuilder builds a union Color reusing builder and
// optional slice allocators for flattened slices.
func Color_CreateUnionSliceWithBuilder(
	elements []*Color,
	builder *Color_Builder,
	allocColorSlice func(n int) []*Color,
	allocStringSlice func(n int) []string,
) *Color {
	if len(elements) == 0 {
		panic(jutil.NewIllegalStateException())
	}
	first := elements[0]
	if first == nil {
		panic(jutil.NewNullPointerException())
	}
	if len(elements) == 1 {
		return first
	}
	allSame := true
	for i := 1; i < len(elements); i++ {
		ei := elements[i]
		if ei == nil {
			panic(jutil.NewNullPointerException())
		}
		if allSame && ei != first && !ei.Equals(first) {
			allSame = false
		}
	}
	if allSame {
		return first
	}

	var flatStack [8]*Color
	flatElems := flatStack[:0]
	var flatSeen dedupSet[*Color]

	var instStack [8]*Color
	var mergedInst []*Color
	var singleInst []*Color
	var instSeen dedupSet[*Color]
	spilledInst := false

	var protoStack [8]*Color
	var mergedProtos []*Color
	var singleProtos []*Color
	var protoSeen dedupSet[*Color]
	spilledProtos := false

	var propStack [8]string
	var mergedProps []string
	var singleProps []string
	var propSeen dedupSet[string]
	spilledProps := false

	isClosureAssert := true
	isConstructor := true
	isInvalidating := false
	propertiesKeepOriginalName := false

	for _, element := range elements {
		if element.IsUnion() {
			for _, nestedElement := range element.unionElements {
				flatElems = appendUnique(flatElems, nestedElement, &flatSeen)
			}
		} else {
			flatElems = appendUnique(flatElems, element, &flatSeen)
		}

		if ic := element.instanceColors; len(ic) > 0 {
			if spilledInst {
				for _, c := range ic {
					mergedInst = appendUnique(mergedInst, c, &instSeen)
				}
			} else if singleInst == nil || &singleInst[0] == &ic[0] {
				singleInst = ic
			} else {
				spilledInst = true
				mergedInst = instStack[:0]
				for _, c := range singleInst {
					mergedInst = appendUnique(mergedInst, c, &instSeen)
				}
				for _, c := range ic {
					mergedInst = appendUnique(mergedInst, c, &instSeen)
				}
				singleInst = nil
			}
		}

		isClosureAssert = isClosureAssert && element.closureAssert
		isConstructor = isConstructor && element.constructor
		isInvalidating = isInvalidating || element.invalidating

		if op := element.ownProperties; len(op) > 0 {
			if spilledProps {
				for _, p := range op {
					mergedProps = appendUnique(mergedProps, p, &propSeen)
				}
			} else if singleProps == nil || &singleProps[0] == &op[0] {
				singleProps = op
			} else {
				spilledProps = true
				mergedProps = propStack[:0]
				for _, p := range singleProps {
					mergedProps = appendUnique(mergedProps, p, &propSeen)
				}
				for _, p := range op {
					mergedProps = appendUnique(mergedProps, p, &propSeen)
				}
				singleProps = nil
			}
		}

		propertiesKeepOriginalName = propertiesKeepOriginalName || element.propertiesKeepOriginalName

		if pr := element.prototypes; len(pr) > 0 {
			if spilledProtos {
				for _, c := range pr {
					mergedProtos = appendUnique(mergedProtos, c, &protoSeen)
				}
			} else if singleProtos == nil || &singleProtos[0] == &pr[0] {
				singleProtos = pr
			} else {
				spilledProtos = true
				mergedProtos = protoStack[:0]
				for _, c := range singleProtos {
					mergedProtos = appendUnique(mergedProtos, c, &protoSeen)
				}
				for _, c := range pr {
					mergedProtos = appendUnique(mergedProtos, c, &protoSeen)
				}
				singleProtos = nil
			}
		}
	}

	var newElements []*Color
	if allocColorSlice != nil {
		newElements = allocColorSlice(len(flatElems))
		copy(newElements, flatElems)
	} else {
		newElements = make([]*Color, len(flatElems))
		copy(newElements, flatElems)
	}

	var instanceColors []*Color
	if spilledInst {
		if allocColorSlice != nil {
			instanceColors = allocColorSlice(len(mergedInst))
		} else {
			instanceColors = make([]*Color, len(mergedInst))
		}
		copy(instanceColors, mergedInst)
	} else if singleInst != nil {
		instanceColors = singleInst
	} else {
		instanceColors = emptyColorSlice
	}

	var prototypes []*Color
	if spilledProtos {
		if allocColorSlice != nil {
			prototypes = allocColorSlice(len(mergedProtos))
		} else {
			prototypes = make([]*Color, len(mergedProtos))
		}
		copy(prototypes, mergedProtos)
	} else if singleProtos != nil {
		prototypes = singleProtos
	} else {
		prototypes = emptyColorSlice
	}

	var ownProperties []string
	if spilledProps {
		if allocStringSlice != nil {
			ownProperties = allocStringSlice(len(mergedProps))
		} else {
			ownProperties = make([]string, len(mergedProps))
		}
		copy(ownProperties, mergedProps)
	} else if singleProps != nil {
		ownProperties = singleProps
	} else {
		ownProperties = emptyStringSlice
	}

	prevTrust := builder.trustUnique
	builder.trustUnique = true
	res := builder.
		ClearBoxId().
		SetClosureAssert(isClosureAssert).
		SetConstructor(isConstructor).
		SetId(colorIdUnionFromColors(newElements)).
		SetInstanceColors(instanceColors).
		SetInvalidating(isInvalidating).
		SetOwnProperties(ownProperties).
		SetPropertiesKeepOriginalName(propertiesKeepOriginalName).
		SetPrototypes(prototypes).
		SetUnionElements(newElements).
		BuildUnion()
	builder.trustUnique = prevTrust
	return res
}

// newColor is the AutoValue constructor; the generated null checks of its
// reference fields are Builder.buildInternal's (the only caller).
//
// Java: $AutoValue_Color(ColorId, ImmutableSet<Color>, ImmutableSet<Color>, boolean, boolean, boolean, ImmutableSet<String>, ColorId, boolean, ImmutableSet<Color>)
func newColor(
	slab *jutil.Slab[Color],
	id ColorId,
	prototypes []*Color,
	instanceColors []*Color,
	invalidating bool,
	propertiesKeepOriginalName bool,
	constructor bool,
	ownProperties []string,
	boxId ColorId,
	hasBoxId bool,
	closureAssert bool,
	unionElements []*Color,
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
	c.hasBoxId = hasBoxId
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
	return isPrimitiveColorId(c.GetId())
}

// Java: Color.isUnion()
func (c *Color) IsUnion() bool {
	// Single element sets are banned in the builder.
	return len(c.unionElements) != 0
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

	elems := c.unionElements
	foundIdx := -1
	for i, e := range elems {
		if e == StandardColors_NULL_OR_VOID || e.Equals(StandardColors_NULL_OR_VOID) {
			foundIdx = i
			break
		}
	}
	if foundIdx < 0 {
		return c
	}
	if len(elems) == 2 {
		return elems[1-foundIdx]
	}

	var stackBuf [8]*Color
	var filtered []*Color
	if len(elems)-1 <= len(stackBuf) {
		filtered = stackBuf[:0]
	} else {
		filtered = make([]*Color, 0, len(elems)-1)
	}
	for i, e := range elems {
		if i != foundIdx {
			filtered = append(filtered, e)
		}
	}
	return Color_CreateUnionSlice(filtered)
}

// Java: $AutoValue_Color.toString()
func (c *Color) String() string {
	var sb strings.Builder
	sb.WriteString("Color{id=")
	sb.WriteString(c.id.String())
	sb.WriteString(", prototypes=")
	sb.WriteString(jutil.CollectionToString(c.prototypes))
	sb.WriteString(", instanceColors=")
	sb.WriteString(jutil.CollectionToString(c.instanceColors))
	sb.WriteString(", invalidating=")
	sb.WriteString(strconv.FormatBool(c.invalidating))
	sb.WriteString(", propertiesKeepOriginalName=")
	sb.WriteString(strconv.FormatBool(c.propertiesKeepOriginalName))
	sb.WriteString(", constructor=")
	sb.WriteString(strconv.FormatBool(c.constructor))
	sb.WriteString(", ownProperties=")
	sb.WriteString(jutil.CollectionToString(c.ownProperties))
	if !c.hasBoxId {
		sb.WriteString(", boxId=null, closureAssert=")
	} else {
		sb.WriteString(", boxId=")
		sb.WriteString(c.boxId.String())
		sb.WriteString(", closureAssert=")
	}
	sb.WriteString(strconv.FormatBool(c.closureAssert))
	sb.WriteString(", unionElements=")
	sb.WriteString(jutil.CollectionToString(c.unionElements))
	sb.WriteString("}")
	return sb.String()
}

func colorSliceSetEquals(a, b []*Color) bool {
	if len(a) != len(b) {
		return false
	}
	if len(a) == 0 || &a[0] == &b[0] {
		return true
	}
	if len(a) <= 16 {
		for i, x := range a {
			bi := b[i]
			if x == bi || x.Equals(bi) {
				continue
			}
			found := false
			for j, y := range b {
				if j != i && (x == y || x.Equals(y)) {
					found = true
					break
				}
			}
			if !found {
				return false
			}
		}
		return true
	}
	var byHash map[int32][]*Color
	for i, x := range a {
		bi := b[i]
		if x == bi || x.Equals(bi) {
			continue
		}
		if byHash == nil {
			byHash = make(map[int32][]*Color, len(b)*2)
			for _, y := range b {
				yh := y.HashCode()
				byHash[yh] = append(byHash[yh], y)
			}
		}
		found := false
		for _, y := range byHash[x.HashCode()] {
			if x == y || x.Equals(y) {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}

func stringSliceSetEquals(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	if len(a) == 0 || &a[0] == &b[0] {
		return true
	}
	if len(a) <= 16 {
		for i, x := range a {
			if x == b[i] {
				continue
			}
			found := false
			for j, y := range b {
				if j != i && x == y {
					found = true
					break
				}
			}
			if !found {
				return false
			}
		}
		return true
	}
	var seen map[string]struct{}
	for i, x := range a {
		if x == b[i] {
			continue
		}
		if seen == nil {
			seen = make(map[string]struct{}, len(b)*2)
			for _, y := range b {
				seen[y] = struct{}{}
			}
		}
		if _, ok := seen[x]; !ok {
			return false
		}
	}
	return true
}

// Java: $AutoValue_Color.equals(Object)
func (c *Color) Equals(o any) bool {
	that, ok := o.(*Color)
	if !ok {
		return false
	}
	if that == c {
		return true
	}
	if c == nil || that == nil {
		return false
	}
	return c.id == that.id &&
		c.invalidating == that.invalidating &&
		c.propertiesKeepOriginalName == that.propertiesKeepOriginalName &&
		c.constructor == that.constructor &&
		c.hasBoxId == that.hasBoxId &&
		(!c.hasBoxId || c.boxId == that.boxId) &&
		c.closureAssert == that.closureAssert &&
		len(c.prototypes) == len(that.prototypes) &&
		len(c.instanceColors) == len(that.instanceColors) &&
		len(c.ownProperties) == len(that.ownProperties) &&
		len(c.unionElements) == len(that.unionElements) &&
		colorSliceSetEquals(c.prototypes, that.prototypes) &&
		colorSliceSetEquals(c.instanceColors, that.instanceColors) &&
		stringSliceSetEquals(c.ownProperties, that.ownProperties) &&
		colorSliceSetEquals(c.unionElements, that.unionElements)
}

func colorBoolHash(b bool) int32 {
	if b {
		return 1231
	}
	return 1237
}

func colorSliceSetHash(cs []*Color) int32 {
	var h int32
	for _, c := range cs {
		if c != nil {
			h += c.HashCode()
		}
	}
	return h
}

func stringSliceSetHash(ss []string) int32 {
	var h int32
	for _, s := range ss {
		h += jutil.StringHashCode(s)
	}
	return h
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
	h ^= colorSliceSetHash(c.prototypes)
	h *= 1000003
	h ^= colorSliceSetHash(c.instanceColors)
	h *= 1000003
	h ^= colorBoolHash(c.invalidating)
	h *= 1000003
	h ^= colorBoolHash(c.propertiesKeepOriginalName)
	h *= 1000003
	h ^= colorBoolHash(c.constructor)
	h *= 1000003
	h ^= stringSliceSetHash(c.ownProperties)
	h *= 1000003
	if c.hasBoxId {
		h ^= c.boxId.HashCode()
	}
	h *= 1000003
	h ^= colorBoolHash(c.closureAssert)
	h *= 1000003
	h ^= colorSliceSetHash(c.unionElements)
	return h
}
