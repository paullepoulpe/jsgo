// Port of com/google/javascript/jscomp/parsing/parser/trees/ComputedPropertyDefinitionTree.java @ closure-compiler v20260415

package parser

// Represents a computed property in an object literal.
type ComputedPropertyDefinitionTree struct {
	ParseTreeBase

	Property ParseTree
	Value    ParseTree
}

// NewComputedPropertyDefinitionTree constructs a ComputedPropertyDefinitionTree.
//
// Java: ComputedPropertyDefinitionTree(SourceRange, ParseTree, ParseTree)
func NewComputedPropertyDefinitionTree(location *SourceRange, property ParseTree, value ParseTree) *ComputedPropertyDefinitionTree {
	t := &ComputedPropertyDefinitionTree{}
	t.InitParseTreeBase(t, ParseTreeType_COMPUTED_PROPERTY_DEFINITION, location)
	t.Property = property
	t.Value = value
	return t
}
