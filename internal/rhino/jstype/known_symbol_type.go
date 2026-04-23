// Port of com/google/javascript/rhino/jstype/KnownSymbolType.java
// @ closure-compiler v20260415
//
// KnownSymbolType represents a concrete, well-known JavaScript symbol —
// Symbol.iterator, Symbol.asyncIterator, etc. Each instance is
// allocated once and compared by pointer, matching Java's class whose
// instances are unique within a JSTypeRegistry. A KnownSymbolType is a
// subtype of the general SymbolType but not of any other
// KnownSymbolType, mirroring Java's rule that "unique KnownSymbolTypes
// are not subtypes of one another".

package jstype

// KnownSymbolType is the type of a specific well-known symbol (e.g.
// Symbol.iterator). It is analogous to Java's KnownSymbolType which
// extends SymbolType. Identity is pointer-based: two *KnownSymbolType
// values are the same type only when they are the same pointer,
// mirroring Java's reference-equality semantics inside the registry.
type KnownSymbolType struct {
	// Name is the display name of the symbol (e.g. "Symbol.iterator").
	// Matches Java's KnownSymbolType.getDisplayName() return value.
	Name string
}

// NewKnownSymbolType allocates a fresh KnownSymbolType with the given
// name. Callers are responsible for interning the result — the registry
// typically stores exactly one instance per well-known symbol name.
func NewKnownSymbolType(name string) *KnownSymbolType {
	return &KnownSymbolType{Name: name}
}

// TypeName mirrors Java KnownSymbolType.getDisplayName(). Returns the
// name supplied at construction (e.g. "Symbol.iterator").
func (k *KnownSymbolType) TypeName() string { return k.Name }

// IsSubtypeOf mirrors Java JSType.isSubtypeOf for KnownSymbolType. A
// KnownSymbolType is a subtype of:
//   - the general SymbolType (any SymbolType value, including another
//     KnownSymbolType, satisfies this per Java's class hierarchy)
//   - itself (pointer identity)
//   - AllType
//   - UnknownType
//
// It is NOT a subtype of a different KnownSymbolType.
func (k *KnownSymbolType) IsSubtypeOf(o JSType) bool {
	if _, ok := o.(SymbolType); ok {
		// General symbol — every known symbol is a subtype of symbol.
		return true
	}
	if other, ok := o.(*KnownSymbolType); ok {
		// Same pointer → same type.
		return k == other
	}
	return isAll(o) || isUnknown(o)
}
