// Port of com/google/javascript/rhino/jstype/StaticTypedSlot.java @ closure-compiler v20260415

package rhino

// StaticTypedSlot must be implemented by variables that can appear as
// members of a StaticTypedScope.
//
// Go-only: language-adapter — getDeclaration() (StaticTypedRef in Java) and
// getScope() (StaticTypedScope in Java) keep StaticSlot's signatures; see
// StaticTypedScope.
type StaticTypedSlot interface {
	StaticSlot

	// GetType returns the type information, if any, for this slot, or null
	// if no type is declared for it.
	GetType() JSType

	// IsTypeInferred returns whether the type has been inferred (as opposed
	// to declared).
	IsTypeInferred() bool
}
