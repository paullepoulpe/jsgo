// Port of com/google/javascript/jscomp/parsing/parser/trees/ComprehensionIfTree.java @ closure-compiler v20260415

package parser

// ComprehensionIfTree is the parse tree for ParseTreeType COMPREHENSION_IF.
type ComprehensionIfTree struct {
	ParseTreeBase

	Expression ParseTree
}

// NewComprehensionIfTree constructs a ComprehensionIfTree.
//
// Java: ComprehensionIfTree(SourceRange, ParseTree)
func NewComprehensionIfTree(location *SourceRange, expression ParseTree) *ComprehensionIfTree {
	t := &ComprehensionIfTree{}
	t.InitParseTreeBase(t, ParseTreeType_COMPREHENSION_IF, location)
	t.Expression = expression
	return t
}
