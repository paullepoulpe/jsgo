// Port of com/google/javascript/rhino/jstype/EnumElementType.java @ closure-compiler v20260415
//
// The type of individual elements of an enum type (see EnumType).
//
// In Java, EnumElementType extends ObjectType and wraps a primitiveType
// field plus a back-pointer to the owning EnumType. In this Go port the
// struct is defined in enum_type.go (where EnumType lives so the two can
// mutually reference each other at allocation time). This file adds the
// faithful 1:1 method set that downstream TypeCheck enum-mismatch
// diagnostics will consume.
//
// Field mapping (Java → Go):
//   primitiveType        → e.Enum.Underlying  (the @enum {T} parameter)
//   primitiveObjectType  → not tracked (no ObjectType cast needed at Tier-0/1)
//   name                 → e.Enum.Name
//   enumType             → e.Enum

package jstype

// GetPrimitiveType returns the underlying primitive (or nominal) type
// that every element of this enum must be assignable to. Mirrors Java's
// EnumElementType.getPrimitiveType().
//
// For example, given `/** @enum {number} */ var E = {A: 1}`, the
// primitive type of E's EnumElementType is Number.
func (e *EnumElementType) GetPrimitiveType() JSType {
	return e.Enum.Underlying
}

// GetEnumType returns the owning EnumType (the enum collection object).
// Mirrors Java's EnumElementType.getEnumType().
func (e *EnumElementType) GetEnumType() *EnumType {
	return e.Enum
}

// GetReferenceName returns the source-level name of this enum element
// type (e.g. "MyEnum" for `/** @enum */ var MyEnum = …`). Mirrors Java's
// EnumElementType.getReferenceName().
func (e *EnumElementType) GetReferenceName() string {
	return e.Enum.Name
}

// ToMaybeEnumElementType implements JSType.toMaybeEnumElementType — only
// EnumElementType returns non-nil; all other types return nil. Mirrors
// Java's EnumElementType.toMaybeEnumElementType().
func (e *EnumElementType) ToMaybeEnumElementType() *EnumElementType {
	return e
}

// GetGreatestSubtype returns the infimum of an EnumElementType and
// another type, or nil if the infimum is empty. Mirrors Java's static
// EnumElementType.getGreatestSubtype(EnumElementType, JSType).
//
// From the Java comment:
//
//	This can be a little bit weird. For example, suppose you have an enum
//	of {(string|number)}, and you want the greatest subtype of the enum and
//	a {number}.
//
//	The infimum is non-empty. But at the same time, we don't really have a
//	name for this infimum. It's equivalent to "elements of this enum that
//	are numbers".
//
//	The best we can do is make up a new type. This is similar to what we do
//	in UnionType#meet.
//
// In the Go port we check whether the underlying primitive type is a
// subtype of `that` (or vice-versa). If the two are disjoint we return
// nil; otherwise we return a fresh EnumElementType with the narrowed
// underlying type. This approximation is sufficient for the Tier-0..3
// diagnostic surface.
func GetGreatestSubtype(element *EnumElementType, that JSType) *EnumElementType {
	// Nil or unknown target: the element type trivially participates.
	if that == nil || isUnknown(that) || isAll(that) {
		return element
	}
	// Check whether the underlying primitive intersects `that`. We use a
	// simple compatibility test: if the underlying is a subtype of `that`
	// we use the underlying as the narrowed type; if `that` is a subtype
	// of the underlying we use `that`. Otherwise the intersection is empty.
	//
	// TODO(nickreid): Align with Java's full greatestSubtype lattice once
	// UnionType.meet is ported (needed for enum-of-union narrowing).
	underlying := element.Enum.Underlying
	var meetPrimitive JSType
	switch {
	case IsSubtype(underlying, that):
		// underlying ⊆ that → infimum is underlying.
		meetPrimitive = underlying
	case IsSubtype(that, underlying):
		// that ⊆ underlying → infimum is that.
		meetPrimitive = that
	default:
		// Disjoint — empty intersection.
		return nil
	}
	// Allocate a fresh (anonymous) EnumType backed by the narrowed
	// primitive. The name is inherited so diagnostics stay meaningful.
	narrowed := NewEnumType(element.Enum.Name, meetPrimitive)
	return narrowed.Element
}
