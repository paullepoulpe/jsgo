// Port of com/google/javascript/rhino/StaticSlot.java
// @ closure-compiler v20260415

package rhino

// StaticSlot must be implemented by variables that can appear as members of
// a StaticScope. Mirrors com.google.javascript.rhino.StaticSlot.
type StaticSlot interface {
	// GetName returns the name of the slot.
	GetName() string

	// GetDeclaration returns the declaration of this symbol. May be nil.
	GetDeclaration() StaticRef

	// GetJSDocInfo returns the JSDoc for this slot.
	GetJSDocInfo() *JSDocInfo

	// GetScope returns the scope that owns this slot.
	GetScope() StaticScope
}
