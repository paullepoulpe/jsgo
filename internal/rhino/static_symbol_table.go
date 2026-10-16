// Port of com/google/javascript/rhino/StaticSymbolTable.java @ closure-compiler v20260415

package rhino

import "iter"

// StaticSymbolTable looks up references by the symbols that they refer to.
type StaticSymbolTable[S StaticSlot, R StaticRef] interface {
	// GetReferences returns the references that point to the given symbol.
	GetReferences(symbol S) iter.Seq[R]

	// GetScope returns the scope for a given symbol.
	GetScope(symbol S) StaticScope

	// GetAllSymbols returns all variables in this symbol table.
	GetAllSymbols() iter.Seq[S]
}
