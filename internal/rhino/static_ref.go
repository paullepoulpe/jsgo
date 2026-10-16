// Port of com/google/javascript/rhino/StaticRef.java @ closure-compiler v20260415

package rhino

// StaticRef tells us all the ways that a StaticSlot is used in a program.
type StaticRef interface {
	// GetSymbol returns the variable that this reference points to.
	GetSymbol() StaticSlot

	// GetNode returns the node where the reference lives.
	GetNode() *Node

	// GetSourceFile returns the source file where the reference lives.
	GetSourceFile() StaticSourceFile
}
