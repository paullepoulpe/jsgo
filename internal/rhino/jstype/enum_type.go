// Port of com/google/javascript/rhino/jstype/EnumType.java +
// com/google/javascript/rhino/jstype/EnumElementType.java
// @ closure-compiler v20260415 (subset).
//
// An `@enum` declaration produces two related types:
//
//   - EnumType           — the enum object itself. Has TypeName
//     `enum{Name}`. References to the enum binding (`var MyEnum
//     = …`) resolve to an EnumType value; aliases (`/** @enum */
//     var Other = MyEnum`) share the SAME EnumType pointer.
//   - EnumElementType    — the type of any element access
//     (`MyEnum.FOO`). Has TypeName `Name<Underlying>` where
//     Underlying is the `@enum {T}` parameter (or Number when
//     omitted). Subtype rule: an EnumElementType is a subtype of
//     its own Underlying, and of any other EnumElementType sharing
//     the same Enum pointer.
//
// The Java class carries more — propertyType lookup, equals / hash
// implementations, union handling. jsgo's port only records the
// surface TypeCheck / TypedScopeCreator currently consume.

package jstype

import "github.com/paullepoulpe/jsgo/internal/rhino"

// EnumType is the type of the `@enum` binding itself. Two enums
// compare equal iff they're literally the same Go pointer —
// aliasing (`/** @enum */ var Alias = Other`) reuses the source
// enum's pointer rather than allocating a fresh one.
type EnumType struct {
	// Name is the source-level identifier (e.g. "MyEnum" for
	// `@enum */ var MyEnum = …`). Matches Java's
	// EnumType.getName().
	Name string
	// Underlying is the primitive (or nominal) type each element
	// must be assignable to. Defaults to Number when `@enum`
	// carries no type argument.
	Underlying JSType
	// Element is the associated EnumElementType — the type of
	// `MyEnum.FOO`. Lazily allocated so the two can circular-
	// reference each other without surprising the IsSubtypeOf
	// walks.
	Element *EnumElementType
	// Members tracks the member names declared on the enum
	// literal (`{FOO: 1, BAR: 2}` → {"FOO", "BAR"}). Used by
	// TypeCheck.visitGetProp to fire JSC_INEXISTENT_ENUM_ELEMENT
	// when a property access names a non-existent element.
	Members map[string]bool

	// docInfo mirrors EnumType.getJSDocInfo() — the JSDoc from the
	// @enum declaration. Set by InferJsDocInfo.
	docInfo *rhino.JSDocInfo
}

// NewEnumType allocates a fresh EnumType + its paired
// EnumElementType. Underlying defaults to Number when nil —
// matches how Java's EnumType defaults the @enum type argument.
func NewEnumType(name string, underlying JSType) *EnumType {
	if underlying == nil {
		underlying = Number
	}
	e := &EnumType{
		Name:       name,
		Underlying: underlying,
		Members:    map[string]bool{},
	}
	e.Element = &EnumElementType{Enum: e}
	return e
}

func (e *EnumType) TypeName() string { return "enum{" + e.Name + "}" }

// IsSubtypeOf: an enum is a subtype only of itself, AllType or
// Unknown. No structural coercion — two enums aren't
// interchangeable even when they share an underlying primitive.
func (e *EnumType) IsSubtypeOf(other JSType) bool {
	if other == nil {
		return true
	}
	if isAll(other) || isUnknown(other) {
		return true
	}
	if o, ok := other.(*EnumType); ok {
		return e == o
	}
	return false
}

// EnumElementType is the type of `Enum.MEMBER`. It behaves as
// "any element of Enum"; elements share the Enum pointer with
// other accesses so `E.FOO` and `E.BAR` are the same static type.
type EnumElementType struct {
	Enum *EnumType
}

func (e *EnumElementType) TypeName() string {
	return e.Enum.Name + "<" + e.Enum.Underlying.TypeName() + ">"
}

// IsSubtypeOf: an EnumElementType is a subtype of (1) its own
// Underlying primitive, (2) another EnumElementType with the
// same Enum pointer, and (3) AllType / Unknown. Notably NOT a
// subtype of an arbitrary primitive — `E<number>` doesn't match
// `string` and vice-versa.
func (e *EnumElementType) IsSubtypeOf(other JSType) bool {
	if other == nil {
		return true
	}
	if isAll(other) || isUnknown(other) {
		return true
	}
	if oe, ok := other.(*EnumElementType); ok {
		return e.Enum == oe.Enum
	}
	// The underlying rule: `E<T>` is a subtype of T.
	return IsSubtype(e.Enum.Underlying, other)
}

// SetJSDocInfo mirrors EnumType.setJSDocInfo — attaches the declaring
// JSDoc to this enum type. Called by InferJsDocInfo.
func (e *EnumType) SetJSDocInfo(info *rhino.JSDocInfo) { e.docInfo = info }

// GetJSDocInfo mirrors EnumType.getJSDocInfo.
func (e *EnumType) GetJSDocInfo() *rhino.JSDocInfo { return e.docInfo }

// GetElementsType mirrors EnumType.getElementsType() — returns the
// EnumElementType paired with this enum, which represents the type of
// `MyEnum.FOO`. Analogous to EnumType.getElementsType() in Java.
func (e *EnumType) GetElementsType() *EnumElementType { return e.Element }
