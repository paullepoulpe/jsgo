// Port of com/google/javascript/rhino/jstype/SymbolType.java
// @ closure-compiler v20260415
//
// SymbolType represents the primitive `symbol` type in Closure's type
// system. It mirrors the Java class of the same name, which extends
// ValueType. A single package-level sentinel (Symbol) is interned —
// matching the Java JSTypeRegistry's singleton pattern for primitive
// types. KnownSymbolType (well-known symbols such as Symbol.iterator)
// extends this type and lives in known_symbol_type.go.

package jstype

// SymbolType is the type of primitive symbol values (e.g. the result
// of `Symbol("description")`). It is the Go equivalent of Java's
// SymbolType extends ValueType. Two symbol values that are both plain
// SymbolType are the same type; a KnownSymbolType is a subtype of
// SymbolType but not of another KnownSymbolType.
type SymbolType struct{}

// TypeName mirrors Java SymbolType.getDisplayName() → "symbol".
func (SymbolType) TypeName() string { return "symbol" }

// IsSubtypeOf mirrors Java JSType.isSubtypeOf. A plain SymbolType is a
// subtype of:
//   - any other SymbolType (including KnownSymbolType, which embeds one)
//   - AllType
//   - UnknownType
func (SymbolType) IsSubtypeOf(o JSType) bool {
	if isSymbol(o) {
		return true
	}
	return isAll(o) || isUnknown(o)
}

// Symbol is the canonical interned SymbolType instance. Mirrors Java's
// JSTypeRegistry returning the same SYMBOL_TYPE singleton.
var Symbol JSType = SymbolType{}

// isSymbol reports whether t is a SymbolType or a *KnownSymbolType.
// Used by both IsSubtypeOf implementations to share the "am I a symbol?"
// check without creating a circular dependency between the two files.
func isSymbol(t JSType) bool {
	if _, ok := t.(SymbolType); ok {
		return true
	}
	if _, ok := t.(*KnownSymbolType); ok {
		return true
	}
	return false
}
