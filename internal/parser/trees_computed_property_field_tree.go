// Port of com/google/javascript/jscomp/parsing/parser/trees/ComputedPropertyFieldTree.java @ closure-compiler v20260415

package parser

// Tree for public comptuted property fields
type ComputedPropertyFieldTree struct {
	ParseTreeBase

	Property    ParseTree
	IsStatic    bool
	Initializer ParseTree
}

// NewComputedPropertyFieldTree constructs a ComputedPropertyFieldTree.
//
// Java: ComputedPropertyFieldTree(SourceRange, ParseTree, boolean, ParseTree)
func NewComputedPropertyFieldTree(location *SourceRange, property ParseTree, isStatic bool, initializer ParseTree) *ComputedPropertyFieldTree {
	t := &ComputedPropertyFieldTree{}
	t.InitParseTreeBase(t, ParseTreeType_COMPUTED_PROPERTY_FIELD, location)
	t.Property = property
	t.IsStatic = isStatic
	t.Initializer = initializer
	return t
}
