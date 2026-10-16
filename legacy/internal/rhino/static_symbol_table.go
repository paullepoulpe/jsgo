// Port of com/google/javascript/rhino/StaticSymbolTable.java
// @ closure-compiler v20260415

package rhino

// StaticSymbolTable allows looking up references by the symbols they refer to.
// Mirrors com.google.javascript.rhino.StaticSymbolTable<S, R>.
//
// Java uses generics bounded by StaticSlot / StaticRef. Go does not support
// generic interface methods in the same way, so this interface is expressed
// over the base interface types. Callers that need concrete types perform
// their own type assertions.
type StaticSymbolTable interface {
	// GetReferences returns the references that point to the given symbol.
	GetReferences(symbol StaticSlot) []StaticRef

	// GetScope returns the scope for a given symbol.
	GetScope(symbol StaticSlot) StaticScope

	// GetAllSymbols returns all variables in this symbol table.
	GetAllSymbols() []StaticSlot
}
