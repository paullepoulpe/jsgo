// Port of com/google/javascript/jscomp/InvalidatingTypes.java @ closure-compiler v20261015

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
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// InvalidatingTypes keeps track of "invalidating types" that force
// type-based optimizations to back off.
//
// Specifically for InlineProperties, disambiguate.AmbiguateProperties, and
// disambiguate.DisambiguateProperties.
type InvalidatingTypes struct {
	// Go-only: language-adapter — an ImmutableSetMultimap keyed by JSType
	// (JSType.equals/hashCode) with Node values (identity).
	typeToLocation *jutil.SetMultimap[rhino.JSType, *rhino.Node]
}

// Java: InvalidatingTypes(ImmutableSetMultimap<JSType, Node>)
func newInvalidatingTypes(typeToLocation *jutil.SetMultimap[rhino.JSType, *rhino.Node]) *InvalidatingTypes {
	t := &InvalidatingTypes{}
	t.typeToLocation = typeToLocation
	return t
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

	return it.typeToLocation.ContainsKey(objType) ||
		// Don't disambiguate properties on object types that are structurally compared or that
		// don't come from a literal class or function definition
		InvalidatingTypes_IsAmbiguousOrStructuralType(objType)
}

// Java: InvalidatingTypes.getMismatchLocations()
func (it *InvalidatingTypes) GetMismatchLocations() *jutil.SetMultimap[rhino.JSType, *rhino.Node] {
	return it.typeToLocation
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

	// Go-only: language-adapter — ImmutableSetMultimap.Builder: keys and
	// values in insertion order, duplicate pairs dropped.
	typeToLocation *jutil.SetMultimap[rhino.JSType, *rhino.Node]
}

// Java: InvalidatingTypes.Builder(JSTypeRegistry)
func NewInvalidatingTypes_Builder(registry *rhino.JSTypeRegistry) *InvalidatingTypes_Builder {
	b := &InvalidatingTypes_Builder{}
	b.typeToLocation = jutil.NewSetMultimap[rhino.JSType, *rhino.Node]()
	b.registry = registry
	return b
}

// Java: InvalidatingTypes.Builder.build()
func (b *InvalidatingTypes_Builder) Build() *InvalidatingTypes {
	for _, t := range invalidatingTypes_ALWAYS_INVALIDATING_TYPES {
		b.typeToLocation.Put(b.registry.GetNativeType(t), invalidatingTypes_ALWAYS_INVALIDATING_LOCATION)
	}

	return newInvalidatingTypes(jutil.SetMultimapCopyOf(b.typeToLocation))
}

// Java: InvalidatingTypes.Builder.addAllTypeMismatches(Iterable<TypeMismatch>)
func (b *InvalidatingTypes_Builder) AddAllTypeMismatches(mismatches []*TypeMismatch) *InvalidatingTypes_Builder {
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
func (b *InvalidatingTypes_Builder) AddTypeWithReason(type_ rhino.JSType, location *rhino.Node) {
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
func (b *InvalidatingTypes_Builder) RecordTypeWithReason(type_ rhino.ObjectType, location *rhino.Node) {
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

	b.typeToLocation.Put(type_, location)
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
