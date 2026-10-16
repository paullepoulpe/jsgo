// Port of com/google/javascript/rhino/jstype/StaticTypedRef.java @ closure-compiler v20260415

package rhino

// StaticTypedRef tells us all the ways that a StaticTypedSlot is used in a
// program.
//
// Go-only: language-adapter — getSymbol() (StaticTypedSlot in Java) keeps
// StaticRef's signature; see StaticTypedScope.
type StaticTypedRef interface {
	StaticRef
}
