// Port of com/google/javascript/jscomp/InvalidatingTypes.java @ closure-compiler v20261021

/*
 * Copyright 2017 The Closure Compiler Authors.
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

package jscomp

import (
	"iter"
	"maps"
	"slices"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

type invalidatingTypeEntry struct {
	typ         rhino.JSType
	firstLoc    rhino.Node
	extraLocs   []rhino.Node
	nextPlusOne uint32
	hasFirstLoc bool
}

// InvalidatingTypes keeps track of "invalidating types" that force
// type-based optimizations to back off.
//
// Specifically for InlineProperties, disambiguate.AmbiguateProperties, and
// disambiguate.DisambiguateProperties.
type InvalidatingTypes struct {
	entries []invalidatingTypeEntry
	buckets map[int32]uint32
}

// Java: InvalidatingTypes.isInvalidating(JSType)
func (it *InvalidatingTypes) IsInvalidating(type_ rhino.JSType) bool {
	if type_ == nil || type_.IsUnknownType() || type_.IsEmptyType() {
		return true
	}

	// A union type is invalidating if any one of its members is invalidating
	if type_.IsUnionType() {
		type_ = type_.RestrictByNotNullOrUndefined()
		if type_.IsUnionType() {
			for _, alt := range type_.GetUnionMembers() {
				if it.IsInvalidating(alt) {
					return true
				}
			}
			return false
		}
	}

	objType := type_.ToMaybeObjectType()
	if objType == nil {
		// TODO(b/174534994): why can't scalars be invalidating?
		return false
	}

	if objType.IsTemplatizedType() {
		objType = objType.ToMaybeTemplatizedType().GetReferencedType()
	}

	// Don't disambiguate properties on object types that are structurally compared or that
	// don't come from a literal class or function definition.
	if InvalidatingTypes_IsAmbiguousOrStructuralType(objType) {
		return true
	}

	for cur := it.buckets[objType.HashCode()]; cur != 0; cur = it.entries[cur-1].nextPlusOne {
		e := it.entries[cur-1].typ
		if e == objType || e.Equals(objType) {
			return true
		}
	}
	return false
}

// MismatchNodeSet is an insertion-ordered, deduplicated view of mismatch AST nodes.
type MismatchNodeSet []rhino.Node

// Contains reports whether n is in s.
func (s MismatchNodeSet) Contains(n rhino.Node) bool {
	return slices.Contains(s, n)
}

// All iterates AST nodes in insertion order.
func (s MismatchNodeSet) All() iter.Seq[rhino.Node] {
	return slices.Values([]rhino.Node(s))
}

// Size returns the number of AST nodes.
func (s MismatchNodeSet) Size() int { return len(s) }

// IsEmpty reports whether s has no AST nodes.
func (s MismatchNodeSet) IsEmpty() bool { return len(s) == 0 }

// ToSlice returns a defensive copy of the AST nodes (non-nil empty slice when empty).
func (s MismatchNodeSet) ToSlice() []rhino.Node {
	if len(s) == 0 {
		return []rhino.Node{}
	}
	return slices.Clone([]rhino.Node(s))
}

type mismatchLocEntry struct {
	typ         rhino.JSType
	locs        []rhino.Node
	nextPlusOne uint32
}

// MismatchLocations maps invalidating JSTypes to their mismatch AST nodes in insertion order.
type MismatchLocations struct {
	entries []mismatchLocEntry
	buckets map[int32]uint32
}

// KeySet iterates the invalidating JSTypes in insertion order.
func (ml *MismatchLocations) KeySet() iter.Seq[rhino.JSType] {
	return func(yield func(rhino.JSType) bool) {
		for i := range ml.entries {
			if !yield(ml.entries[i].typ) {
				return
			}
		}
	}
}

// Get returns the mismatch AST nodes for typ (empty when absent).
func (ml *MismatchLocations) Get(typ rhino.JSType) MismatchNodeSet {
	if typ == nil || len(ml.entries) == 0 {
		return nil
	}
	for cur := ml.buckets[typ.HashCode()]; cur != 0; cur = ml.entries[cur-1].nextPlusOne {
		e := &ml.entries[cur-1]
		if e.typ == typ || e.typ.Equals(typ) {
			return MismatchNodeSet(e.locs)
		}
	}
	return nil
}

// InverseMismatchLocations maps mismatch AST nodes to their involved invalidating JSTypes.
type InverseMismatchLocations struct {
	order  []rhino.Node
	byNode map[rhino.Node][]rhino.JSType
}

// InverseMismatchLocationsMap adapts InverseMismatchLocations for ordered iteration.
type InverseMismatchLocationsMap struct {
	inv *InverseMismatchLocations
}

// All iterates (node, involved types) pairs in insertion order.
func (m InverseMismatchLocationsMap) All() iter.Seq2[rhino.Node, *jutil.ObjectLinkedHashSet[rhino.JSType]] {
	return func(yield func(rhino.Node, *jutil.ObjectLinkedHashSet[rhino.JSType]) bool) {
		if m.inv == nil {
			return
		}
		for _, loc := range m.inv.order {
			types := m.inv.byNode[loc]
			s := jutil.NewObjectLinkedHashSet[rhino.JSType]()
			for _, t := range types {
				s.Add(t)
			}
			if !yield(loc, s) {
				return
			}
		}
	}
}

// AsMap returns an ordered map view of the node-to-types index.
func (iml *InverseMismatchLocations) AsMap() InverseMismatchLocationsMap {
	return InverseMismatchLocationsMap{inv: iml}
}

// Inverse returns a node-to-types view preserving insertion order.
func (ml *MismatchLocations) Inverse() *InverseMismatchLocations {
	if len(ml.entries) == 0 {
		return &InverseMismatchLocations{}
	}
	byNode := make(map[rhino.Node][]rhino.JSType, len(ml.entries))
	order := make([]rhino.Node, 0, len(ml.entries))
	for i := range ml.entries {
		typ := ml.entries[i].typ
		for _, loc := range ml.entries[i].locs {
			existing, ok := byNode[loc]
			if !ok {
				order = append(order, loc)
			} else {
				dup := false
				for _, t := range existing {
					if t == typ || t.Equals(typ) {
						dup = true
						break
					}
				}
				if dup {
					continue
				}
			}
			byNode[loc] = append(existing, typ)
		}
	}
	return &InverseMismatchLocations{order: order, byNode: byNode}
}

// Java: InvalidatingTypes.getMismatchLocations()
func (it *InvalidatingTypes) GetMismatchLocations() *MismatchLocations {
	count := 0
	for i := range it.entries {
		if it.entries[i].hasFirstLoc {
			count++
		}
	}
	if count == 0 {
		return &MismatchLocations{}
	}
	entries := make([]mismatchLocEntry, 0, count)
	buckets := make(map[int32]uint32, count)
	for i := range it.entries {
		e := &it.entries[i]
		if !e.hasFirstLoc {
			continue
		}
		locs := make([]rhino.Node, 1+len(e.extraLocs))
		locs[0] = e.firstLoc
		copy(locs[1:], e.extraLocs)
		h := e.typ.HashCode()
		idx := uint32(len(entries))
		entries = append(entries, mismatchLocEntry{
			typ:         e.typ,
			locs:        locs,
			nextPlusOne: buckets[h],
		})
		buckets[h] = idx + 1
	}
	return &MismatchLocations{
		entries: entries,
		buckets: buckets,
	}
}

// InvalidatingTypes_IsAmbiguousOrStructuralType returns true if any of the
// following hold:
//   - this type obeys structural subtyping rules, as opposed to nominal
//     subtyping
//   - this type is some JSDoc-only or anonymous type like a mixin, as
//     opposed to a class or function literal
//
// Java: InvalidatingTypes.isAmbiguousOrStructuralType(ObjectType)
func InvalidatingTypes_IsAmbiguousOrStructuralType(type_ rhino.ObjectType) bool {
	if type_.IsEnumType() {
		// enum types are created via object literals, which are normally structural, but Closure
		// special-cases them to behave as if nominal.
		return false
	} else if type_.IsEnumElementType() {
		primitive := type_.ToMaybeEnumElementType().GetPrimitiveType().ToMaybeObjectType()
		// Treat an Enum<Foo> identically to a Foo
		return primitive == nil || InvalidatingTypes_IsAmbiguousOrStructuralType(primitive)
	} else if type_.IsFunctionType() {
		return !type_.IsNominalConstructorOrInterface() ||
			type_.ToMaybeFunctionType().IsAmbiguousConstructor()
	} else if type_.IsFunctionPrototypeType() {
		ownerFunction := type_.GetOwnerFunction()
		return ownerFunction == nil ||
			!ownerFunction.IsNominalConstructorOrInterface() ||
			ownerFunction.IsAmbiguousConstructor() ||
			ownerFunction.IsStructuralInterface()
	} else if type_.IsInstanceType() {
		ctor := type_.GetConstructor()
		return ctor == nil || ctor.IsAmbiguousConstructor() || ctor.IsStructuralInterface()
	}

	return true
}

// InvalidatingTypes_Builder is InvalidatingTypes.Builder.
type InvalidatingTypes_Builder struct {
	registry *rhino.JSTypeRegistry
	entries  []invalidatingTypeEntry
	buckets  map[int32]uint32
	consumed bool
}

// Java: InvalidatingTypes.Builder(JSTypeRegistry)
func NewInvalidatingTypes_Builder(registry *rhino.JSTypeRegistry) *InvalidatingTypes_Builder {
	return &InvalidatingTypes_Builder{
		registry: registry,
		entries:  make([]invalidatingTypeEntry, 0, 8),
		buckets:  make(map[int32]uint32, 8),
	}
}

func (b *InvalidatingTypes_Builder) ensureMutable() {
	if b.consumed {
		cloned := make([]invalidatingTypeEntry, len(b.entries))
		copy(cloned, b.entries)
		for i := range cloned {
			if len(cloned[i].extraLocs) > 0 {
				cloned[i].extraLocs = slices.Clone(cloned[i].extraLocs)
			}
		}
		b.entries = cloned
		b.buckets = maps.Clone(b.buckets)
		b.consumed = false
	}
}

func (b *InvalidatingTypes_Builder) putTypeLocation(typ rhino.JSType, location rhino.Node) {
	b.ensureMutable()
	h := typ.HashCode()
	for cur := b.buckets[h]; cur != 0; cur = b.entries[cur-1].nextPlusOne {
		e := &b.entries[cur-1]
		if e.typ == typ || e.typ.Equals(typ) {
			if !e.hasFirstLoc {
				e.firstLoc = location
				e.hasFirstLoc = true
				return
			}
			if e.firstLoc == location {
				return
			}
			for _, existing := range e.extraLocs {
				if existing == location {
					return
				}
			}
			e.extraLocs = append(e.extraLocs, location)
			return
		}
	}
	idx := uint32(len(b.entries))
	b.entries = append(b.entries, invalidatingTypeEntry{
		typ:         typ,
		firstLoc:    location,
		nextPlusOne: b.buckets[h],
		hasFirstLoc: true,
	})
	b.buckets[h] = idx + 1
}

// Java: InvalidatingTypes.Builder.build()
func (b *InvalidatingTypes_Builder) Build() *InvalidatingTypes {
	for _, t := range invalidatingTypes_ALWAYS_INVALIDATING_TYPES {
		b.putTypeLocation(b.registry.GetNativeType(t), invalidatingTypes_ALWAYS_INVALIDATING_LOCATION)
	}
	b.consumed = true
	return &InvalidatingTypes{
		entries: b.entries,
		buckets: b.buckets,
	}
}

// Java: InvalidatingTypes.Builder.addAllTypeMismatches(Iterable<TypeMismatch>)
func (b *InvalidatingTypes_Builder) AddAllTypeMismatches(mismatches []TypeMismatch) *InvalidatingTypes_Builder {
	for _, mismatch := range mismatches {
		b.AddTypeWithReason(mismatch.Found(), mismatch.Location())
		b.AddTypeWithReason(mismatch.Required(), mismatch.Location())
	}

	return b
}

// AddTypeWithReason invalidates the given type, so that no properties on it
// will be inlined or renamed.
//
// Java: InvalidatingTypes.Builder.addTypeWithReason(JSType, Node)
func (b *InvalidatingTypes_Builder) AddTypeWithReason(type_ rhino.JSType, location rhino.Node) {
	type_ = type_.RestrictByNotNullOrUndefined()

	if type_.IsUnionType() {
		for _, alt := range type_.GetUnionMembers() {
			b.AddTypeWithReason(alt, location)
		}
		return
	}

	if type_.IsEnumElementType() {
		b.AddTypeWithReason(type_.GetEnumeratedTypeOfEnumElement(), location)
		return
	}

	objType := type_.ToMaybeObjectType()
	if objType == nil {
		return
	}

	b.RecordTypeWithReason(objType, location)
	b.RecordTypeWithReason(objType.GetImplicitPrototype(), location)

	if objType.IsConstructor() {
		// TODO(b/142431852): This should never be null but it is possible.
		// Case: `function(new:T)`, `T = number`.
		instanceType := objType.ToMaybeFunctionType().GetInstanceType()
		b.RecordTypeWithReason(instanceType, location)
	} else if objType.IsInstanceType() {
		b.RecordTypeWithReason(objType.GetConstructor(), location)
	}
}

// Java: InvalidatingTypes.Builder.recordTypeWithReason(ObjectType, Node)
func (b *InvalidatingTypes_Builder) RecordTypeWithReason(type_ rhino.ObjectType, location rhino.Node) {
	if type_ == nil {
		return
	}
	if type_.IsTemplatizedType() {
		type_ = type_.ToMaybeTemplatizedType().GetReferencedType()
	}
	if InvalidatingTypes_IsAmbiguousOrStructuralType(type_) {
		// This type is inherently invalidating. Putting it in the map wastes memory.
		// This also fixes a performance regression: previously we saw ~4k structural types
		// hash to the same bucket in the "typeToLocation" map, causing >100 seconds spent in
		// hash map lookups for some builds.
		return
	}

	b.putTypeLocation(type_, location)
}

// Java: InvalidatingTypes.ALWAYS_INVALIDATING_TYPES
var invalidatingTypes_ALWAYS_INVALIDATING_TYPES = []rhino.JSTypeNative{
	rhino.JSTypeNative_FUNCTION_FUNCTION_TYPE,
	rhino.JSTypeNative_FUNCTION_TYPE,
	rhino.JSTypeNative_FUNCTION_PROTOTYPE,
	rhino.JSTypeNative_FUNCTION_INSTANCE_PROTOTYPE,
	rhino.JSTypeNative_OBJECT_TYPE,
	rhino.JSTypeNative_OBJECT_PROTOTYPE,
	rhino.JSTypeNative_OBJECT_FUNCTION_TYPE,
}

// Java: InvalidatingTypes.ALWAYS_INVALIDATING_LOCATION
var invalidatingTypes_ALWAYS_INVALIDATING_LOCATION = rhino.IR_Name("alwaysInvalidatingLocation").
	SetStaticSourceFile(
		SourceFile_FromCode("InvalidatingTypes_alwaysInvalidatingLocation", ""))
