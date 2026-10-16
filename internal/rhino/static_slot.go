// Port of com/google/javascript/rhino/StaticSlot.java @ closure-compiler v20260415

package rhino

// StaticSlot must be implemented by variables that can appear as members of
// a StaticScope.
type StaticSlot interface {
	// GetName gets the name of the slot.
	GetName() string

	// GetDeclaration gets the declaration of this symbol. May not exist.
	GetDeclaration() StaticRef

	// GetJSDocInfo gets the JSDoc for this slot.
	GetJSDocInfo() *JSDocInfo

	GetScope() StaticScope
}
