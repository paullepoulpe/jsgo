// Port of com/google/javascript/jscomp/parsing/parser/trees/ComputedPropertyMethodTree.java @ closure-compiler v20260415

package parser

// ComputedPropertyMethodTree is the parse tree for ParseTreeType COMPUTED_PROPERTY_METHOD.
type ComputedPropertyMethodTree struct {
	ParseTreeBase

	Property ParseTree
	Method   ParseTree
}

// NewComputedPropertyMethodTree constructs a ComputedPropertyMethodTree.
//
// Java: ComputedPropertyMethodTree(SourceRange, ParseTree, ParseTree)
func NewComputedPropertyMethodTree(location *SourceRange, property ParseTree, method ParseTree) *ComputedPropertyMethodTree {
	t := &ComputedPropertyMethodTree{}
	t.InitParseTreeBase(t, ParseTreeType_COMPUTED_PROPERTY_METHOD, location)
	t.Property = property
	t.Method = method
	return t
}
