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
	"encoding/binary"
	"slices"
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
func (c *Color) GetBoxId() *ColorId {
	if !c.hasBoxId {
		return nil
	}
	return &c.boxId
}

// IsClosureAssert reports whether this type is some Closure assertion
// function removable by Closure-specific optimizations.
//
// Java: Color.isClosureAssert()
func (c *Color) IsClosureAssert() bool { return c.closureAssert }

// Java: Color.getUnionElements()
func (c *Color) GetUnionElements() []*Color { return c.unionElements }

var (
	colorEmptyColorSet = jutil.NewObjectLinkedHashSet[*Color]()
	emptyColorSlice    = []*Color{}
	emptyStringSlice   = []string{}
)

// EmptyColorSlice returns the shared non-nil empty []*Color sentinel.
func EmptyColorSlice() []*Color { return emptyColorSlice }

// EmptyStringSlice returns the shared non-nil empty []string sentinel.
func EmptyStringSlice() []string { return emptyStringSlice }

// Java: Color.singleBuilder()
func Color_SingleBuilder() *Color_Builder {
	return &Color_Builder{
		prototypes:     emptyColorSlice,
		instanceColors: emptyColorSlice,
		ownProperties:  emptyStringSlice,
		unionElements:  emptyColorSlice,
		set_0:          0x0f,
	}
}

// Java: Color.createUnion(Set<Color>)
func Color_CreateUnion(elements *jutil.ObjectLinkedHashSet[*Color]) *Color {
	n := elements.Size()
	switch n {
	case 0:
		panic(jutil.NewIllegalStateException())
	case 1:
		for e := range elements.All() {
			return e
		}
	}
	var stackBuf [8]*Color
	var slice []*Color
	if n <= len(stackBuf) {
		slice = stackBuf[:0]
		for e := range elements.All() {
			slice = append(slice, e)
		}
	} else {
		slice = elements.ToSlice()
	}
	var builder Color_Builder
	return Color_CreateUnionSliceWithBuilder(slice, &builder, nil, nil)
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

func appendUniqueColor(dst []*Color, c *Color, byHash *map[int32][]*Color) []*Color {
	if *byHash != nil {
		h := c.HashCode()
		bucket := (*byHash)[h]
		for _, existing := range bucket {
			if existing == c || existing.Equals(c) {
				return dst
			}
		}
		(*byHash)[h] = append(bucket, c)
		return append(dst, c)
	}
	for _, existing := range dst {
		if existing == c || existing.Equals(c) {
			return dst
		}
	}
	dst = append(dst, c)
	if len(dst) > 16 {
		m := make(map[int32][]*Color, len(dst)*2)
		for _, existing := range dst {
			eh := existing.HashCode()
			m[eh] = append(m[eh], existing)
		}
		*byHash = m
	}
	return dst
}

func appendUniqueString(dst []string, s string, seen *map[string]struct{}) []string {
	if *seen != nil {
		if _, ok := (*seen)[s]; ok {
			return dst
		}
		(*seen)[s] = struct{}{}
		return append(dst, s)
	}
	for _, existing := range dst {
		if existing == s {
			return dst
		}
	}
	dst = append(dst, s)
	if len(dst) > 16 {
		m := make(map[string]struct{}, len(dst)*2)
		for _, existing := range dst {
			m[existing] = struct{}{}
		}
		*seen = m
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
	slices.Sort(ids)
	ids = slices.Compact(ids)
	n := len(ids)
	if n == 1 {
		return NewColorId(ids[0])
	}
	if n == 2 {
		var buf [16]byte
		binary.LittleEndian.PutUint64(buf[0:8], uint64(ids[0]))
		binary.LittleEndian.PutUint64(buf[8:16], uint64(ids[1]))
		return NewColorId(int64(jutil.FarmHashFingerprint64(buf[:])))
	}
	var stackBytes [256]byte
	var raw []byte
	if n*8 <= len(stackBytes) {
		raw = stackBytes[:n*8]
	} else {
		raw = make([]byte, n*8)
	}
	for i := 0; i < n; i++ {
		binary.LittleEndian.PutUint64(raw[i*8:], uint64(ids[i]))
	}
	return NewColorId(int64(jutil.FarmHashFingerprint64(raw)))
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
	var flatByHash map[int32][]*Color

	var instStack [8]*Color
	var mergedInst []*Color
	var singleInst []*Color
	var instByHash map[int32][]*Color
	spilledInst := false

	var protoStack [8]*Color
	var mergedProtos []*Color
	var singleProtos []*Color
	var protoByHash map[int32][]*Color
	spilledProtos := false

	var propStack [8]string
	var mergedProps []string
	var singleProps []string
	var propSeen map[string]struct{}
	spilledProps := false

	isClosureAssert := true
	isConstructor := true
	isInvalidating := false
	propertiesKeepOriginalName := false

	for _, element := range elements {
		if element.IsUnion() {
			for _, nestedElement := range element.unionElements {
				flatElems = appendUniqueColor(flatElems, nestedElement, &flatByHash)
			}
		} else {
			flatElems = appendUniqueColor(flatElems, element, &flatByHash)
		}

		if ic := element.instanceColors; len(ic) > 0 {
			if spilledInst {
				for _, c := range ic {
					mergedInst = appendUniqueColor(mergedInst, c, &instByHash)
				}
			} else if singleInst == nil || &singleInst[0] == &ic[0] {
				singleInst = ic
			} else {
				spilledInst = true
				mergedInst = instStack[:0]
				for _, c := range singleInst {
					mergedInst = appendUniqueColor(mergedInst, c, &instByHash)
				}
				for _, c := range ic {
					mergedInst = appendUniqueColor(mergedInst, c, &instByHash)
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
					mergedProps = appendUniqueString(mergedProps, p, &propSeen)
				}
			} else if singleProps == nil || &singleProps[0] == &op[0] {
				singleProps = op
			} else {
				spilledProps = true
				mergedProps = propStack[:0]
				for _, p := range singleProps {
					mergedProps = appendUniqueString(mergedProps, p, &propSeen)
				}
				for _, p := range op {
					mergedProps = appendUniqueString(mergedProps, p, &propSeen)
				}
				singleProps = nil
			}
		}

		propertiesKeepOriginalName = propertiesKeepOriginalName || element.propertiesKeepOriginalName

		if pr := element.prototypes; len(pr) > 0 {
			if spilledProtos {
				for _, c := range pr {
					mergedProtos = appendUniqueColor(mergedProtos, c, &protoByHash)
				}
			} else if singleProtos == nil || &singleProtos[0] == &pr[0] {
				singleProtos = pr
			} else {
				spilledProtos = true
				mergedProtos = protoStack[:0]
				for _, c := range singleProtos {
					mergedProtos = appendUniqueColor(mergedProtos, c, &protoByHash)
				}
				for _, c := range pr {
					mergedProtos = appendUniqueColor(mergedProtos, c, &protoByHash)
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
	prototypes                 []*Color
	instanceColors             []*Color
	ownProperties              []string
	unionElements              []*Color
	slab                       *jutil.Slab[Color]
	id                         ColorId
	boxId                      ColorId
	hasBoxId                   bool
	invalidating               bool
	propertiesKeepOriginalName bool
	constructor                bool
	closureAssert              bool
	trustUnique                bool
	set_0                      byte
}

// Java: $AutoValue_Color.Builder()
func NewColor_Builder() *Color_Builder { return &Color_Builder{} }

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
	b.set_0 |= 8
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

// Java: Color.Builder.setBoxId(ColorId)
func (b *Color_Builder) SetBoxId(boxId *ColorId) *Color_Builder {
	if boxId == nil {
		b.boxId = ColorId{}
		b.hasBoxId = false
	} else {
		b.boxId = *boxId
		b.hasBoxId = true
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
		b.hasBoxId,
		b.closureAssert,
		b.unionElements)
}

// Java: Color.Builder.build()
func (b *Color_Builder) Build() *Color {
	result := b.BuildInternal()
	jutil.CheckStateMsg1(len(result.unionElements) == 0, "%s", result)
	jutil.CheckStateMsg1(!result.hasBoxId, "%s", result)
	jutil.CheckStateMsg1(!StandardColors_AXIOMATIC_COLORS.ContainsKey(result.GetId()), "%s", result)
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
