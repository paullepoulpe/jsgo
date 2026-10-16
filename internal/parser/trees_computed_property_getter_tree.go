// Port of com/google/javascript/jscomp/parsing/parser/trees/ComputedPropertyGetterTree.java @ closure-compiler v20260415

package parser

// ComputedPropertyGetterTree is the parse tree for ParseTreeType COMPUTED_PROPERTY_GETTER.
type ComputedPropertyGetterTree struct {
	ParseTreeBase

	Property ParseTree
	IsStatic bool
	Body     *BlockTree
}

// NewComputedPropertyGetterTree constructs a ComputedPropertyGetterTree.
//
// Java: ComputedPropertyGetterTree(SourceRange, ParseTree, boolean, BlockTree)
func NewComputedPropertyGetterTree(location *SourceRange, property ParseTree, isStatic bool, body *BlockTree) *ComputedPropertyGetterTree {
	t := &ComputedPropertyGetterTree{}
	t.InitParseTreeBase(t, ParseTreeType_COMPUTED_PROPERTY_GETTER, location)
	t.Property = property
	t.IsStatic = isStatic
	t.Body = body
	return t
}
