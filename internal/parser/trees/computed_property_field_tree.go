// Port of com/google/javascript/jscomp/parsing/parser/trees/ComputedPropertyFieldTree.java
// @ closure-compiler v20260415

package trees

// ComputedPropertyFieldTree represents a public computed property field in a class.
type ComputedPropertyFieldTree struct {
	ParseTree
	Property *ParseTree
	IsStatic bool
	// Initializer is nil when the field has no initializer.
	Initializer *ParseTree
}

func NewComputedPropertyFieldTree(
	location SourceRange,
	property *ParseTree,
	isStatic bool,
	initializer *ParseTree,
) *ComputedPropertyFieldTree {
	return &ComputedPropertyFieldTree{
		ParseTree:   NewParseTree(ParseTreeTypeComputedPropertyField, location),
		Property:    property,
		IsStatic:    isStatic,
		Initializer: initializer,
	}
}
