// Port of com/google/javascript/jscomp/colors/Color.java @ closure-compiler v20261028
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
	"math/bits"
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
	if e1.equalsColor(e2) {
		return e1
	}
	buf := [2]*Color{e1, e2}
	var builder Color_Builder
	return Color_CreateUnionSliceWithBuilder(buf[:], &builder, nil, nil)
}

type setOps[T comparable] interface {
	checkNonNull(v T)
	equal(a, b T) bool
	hash(v T) int32
}

type colorSetOps struct{}

func (colorSetOps) checkNonNull(c *Color) {
	if c == nil {
		panic(jutil.NewNullPointerException())
	}
}
func (colorSetOps) equal(a, b *Color) bool { return a.equalsColor(b) }
func (colorSetOps) hash(c *Color) int32    { return c.HashCode() }

type stringSetOps struct{}

func (stringSetOps) checkNonNull(string)    {}
func (stringSetOps) equal(a, b string) bool { return a == b }
func (stringSetOps) hash(s string) int32    { return jutil.StringHashCode(s) }

func rebuildIndexTable[T comparable, Ops setOps[T]](items []T, size int) []int32 {
	var ops Ops
	slots := make([]int32, size)
	mask := uint32(size - 1)
	for i, item := range items {
		for slot := uint32(ops.hash(item)) & mask; ; slot = (slot + 1) & mask {
			if slots[slot] == 0 {
				slots[slot] = int32(i + 1)
				break
			}
		}
	}
	return slots
}

func appendUnique[T comparable, Ops setOps[T]](dst []T, v T, table *[]int32) []T {
	var ops Ops
	if slots := *table; slots != nil {
		mask := uint32(len(slots) - 1)
		for slot := uint32(ops.hash(v)) & mask; ; slot = (slot + 1) & mask {
			idx := slots[slot]
			if idx == 0 {
				dst = append(dst, v)
				if len(dst)*2 > len(slots) {
					*table = rebuildIndexTable[T, Ops](dst, len(slots)*2)
				} else {
					slots[slot] = int32(len(dst))
				}
				return dst
			}
			if ops.equal(dst[idx-1], v) {
				return dst
			}
		}
	}
	for _, existing := range dst {
		if ops.equal(existing, v) {
			return dst
		}
	}
	dst = append(dst, v)
	if len(dst) > 16 {
		*table = rebuildIndexTable[T, Ops](dst, 64)
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

func allocCopySlice[T any](src []T, alloc func(int) []T) []T {
	var out []T
	if alloc != nil {
		out = alloc(len(src))
	} else {
		out = make([]T, len(src))
	}
	copy(out, src)
	return out
}

// setAccum deduplicates elements in insertion order using an inline [8]T stack
// buffer before spilling to a linear slice (<= 16) and open-addressed index table.
type setAccum[T comparable, Ops setOps[T]] struct {
	inline [8]T
	spill  []T
	table  []int32
	n      int
}

func (a *setAccum[T, Ops]) add(v T) {
	if a.spill != nil {
		a.spill = appendUnique[T, Ops](a.spill, v, &a.table)
		return
	}
	var ops Ops
	for _, existing := range a.inline[:a.n] {
		if ops.equal(existing, v) {
			return
		}
	}
	if a.n < len(a.inline) {
		a.inline[a.n] = v
		a.n++
		return
	}
	a.spill = make([]T, len(a.inline)+1, len(a.inline)*2)
	copy(a.spill, a.inline[:])
	a.spill[len(a.inline)] = v
}

func (a *setAccum[T, Ops]) copyOut(alloc func(int) []T) []T {
	if a.spill != nil {
		return allocCopySlice(a.spill, alloc)
	}
	return allocCopySlice(a.inline[:a.n], alloc)
}

// sliceUnionAccum merges slices across union members, reusing a single member's
// slice directly until two distinct non-empty slices appear.
type sliceUnionAccum[T comparable, Ops setOps[T]] struct {
	single  []T
	merged  setAccum[T, Ops]
	spilled bool
}

func (a *sliceUnionAccum[T, Ops]) addSlice(s []T) {
	if len(s) == 0 {
		return
	}
	if a.spilled {
		for _, v := range s {
			a.merged.add(v)
		}
		return
	}
	if a.single == nil || &a.single[0] == &s[0] {
		a.single = s
		return
	}
	a.spilled = true
	for _, v := range a.single {
		a.merged.add(v)
	}
	for _, v := range s {
		a.merged.add(v)
	}
	a.single = nil
}

func (a *sliceUnionAccum[T, Ops]) finish(empty []T, alloc func(int) []T) []T {
	if a.spilled {
		return a.merged.copyOut(alloc)
	}
	if a.single != nil {
		return a.single
	}
	return empty
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
		if allSame && !ei.equalsColor(first) {
			allSame = false
		}
	}
	if allSame {
		return first
	}

	var flatElems setAccum[*Color, colorSetOps]
	var inst, protos sliceUnionAccum[*Color, colorSetOps]
	var props sliceUnionAccum[string, stringSetOps]

	isClosureAssert := true
	isConstructor := true
	isInvalidating := false
	propertiesKeepOriginalName := false

	for _, element := range elements {
		if element.IsUnion() {
			for _, nestedElement := range element.unionElements {
				flatElems.add(nestedElement)
			}
		} else {
			flatElems.add(element)
		}
		inst.addSlice(element.instanceColors)
		isClosureAssert = isClosureAssert && element.closureAssert
		isConstructor = isConstructor && element.constructor
		isInvalidating = isInvalidating || element.invalidating
		props.addSlice(element.ownProperties)
		propertiesKeepOriginalName = propertiesKeepOriginalName || element.propertiesKeepOriginalName
		protos.addSlice(element.prototypes)
	}

	newElements := flatElems.copyOut(allocColorSlice)
	instanceColors := inst.finish(emptyColorSlice, allocColorSlice)
	prototypes := protos.finish(emptyColorSlice, allocColorSlice)
	ownProperties := props.finish(emptyStringSlice, allocStringSlice)

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

// newColor constructs a Color from validated Color_Builder fields.
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
		if e.equalsColor(StandardColors_NULL_OR_VOID) {
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

func sliceSetEquals[T comparable, Ops setOps[T]](a, b []T) bool {
	if len(a) != len(b) {
		return false
	}
	if len(a) == 0 || &a[0] == &b[0] {
		return true
	}
	var ops Ops
	if len(a) <= 16 {
		for i, x := range a {
			if ops.equal(x, b[i]) {
				continue
			}
			found := false
			for j, y := range b {
				if j != i && ops.equal(x, y) {
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
	var slots []int32
	var mask uint32
	for i, x := range a {
		if ops.equal(x, b[i]) {
			continue
		}
		if slots == nil {
			tableSize := 1 << bits.Len(uint(len(b)*2-1))
			slots = rebuildIndexTable[T, Ops](b, tableSize)
			mask = uint32(tableSize - 1)
		}
		found := false
		for slot := uint32(ops.hash(x)) & mask; ; slot = (slot + 1) & mask {
			idx := slots[slot]
			if idx == 0 {
				break
			}
			if ops.equal(x, b[idx-1]) {
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

func (c *Color) equalsColor(that *Color) bool {
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
		sliceSetEquals[*Color, colorSetOps](c.prototypes, that.prototypes) &&
		sliceSetEquals[*Color, colorSetOps](c.instanceColors, that.instanceColors) &&
		sliceSetEquals[string, stringSetOps](c.ownProperties, that.ownProperties) &&
		sliceSetEquals[*Color, colorSetOps](c.unionElements, that.unionElements)
}

// Java: $AutoValue_Color.equals(Object)
func (c *Color) Equals(o any) bool {
	that, ok := o.(*Color)
	return ok && c.equalsColor(that)
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
