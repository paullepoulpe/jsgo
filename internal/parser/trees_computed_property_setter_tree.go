// Port of com/google/javascript/jscomp/parsing/parser/trees/ComputedPropertySetterTree.java @ closure-compiler v20260415

package parser

// ComputedPropertySetterTree is the parse tree for ParseTreeType COMPUTED_PROPERTY_SETTER.
type ComputedPropertySetterTree struct {
	ParseTreeBase

	Property  ParseTree
	Parameter *FormalParameterListTree
	IsStatic  bool
	Body      *BlockTree
}

// NewComputedPropertySetterTree constructs a ComputedPropertySetterTree.
//
// Java: ComputedPropertySetterTree(SourceRange, ParseTree, boolean, FormalParameterListTree, BlockTree)
func NewComputedPropertySetterTree(location *SourceRange, property ParseTree, isStatic bool, parameter *FormalParameterListTree, body *BlockTree) *ComputedPropertySetterTree {
	t := &ComputedPropertySetterTree{}
	t.InitParseTreeBase(t, ParseTreeType_COMPUTED_PROPERTY_SETTER, location)
	t.Property = property
	t.IsStatic = isStatic
	t.Parameter = parameter
	t.Body = body
	return t
}
