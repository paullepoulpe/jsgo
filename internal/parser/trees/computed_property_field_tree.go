// Port of com/google/javascript/jscomp/parsing/parser/trees/ComputedPropertyFieldTree.java
// @ closure-compiler v20260415

package trees

// ComputedPropertyFieldTree represents a class field with a computed
// property name, e.g. `class { [key] = expr; }` or `static [key];`.
type ComputedPropertyFieldTree struct {
	ParseTree
	Property    *ParseTree
	IsStatic    bool
	Initializer *ParseTree // nil when the field has no initialiser
}

// NewComputedPropertyFieldTree constructs a ComputedPropertyFieldTree.
func NewComputedPropertyFieldTree(location SourceRange, property *ParseTree, isStatic bool, initializer *ParseTree) *ComputedPropertyFieldTree {
	return &ComputedPropertyFieldTree{
		ParseTree:   ParseTree{Type: ParseTreeTypeComputedPropertyField, Location: location},
		Property:    property,
		IsStatic:    isStatic,
		Initializer: initializer,
	}
}
