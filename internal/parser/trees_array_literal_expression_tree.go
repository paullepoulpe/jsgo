// Port of com/google/javascript/jscomp/parsing/parser/trees/ArrayLiteralExpressionTree.java @ closure-compiler v20260415

package parser

// ArrayLiteralExpressionTree is the parse tree for ParseTreeType ARRAY_LITERAL_EXPRESSION.
type ArrayLiteralExpressionTree struct {
	ParseTreeBase

	Elements         []ParseTree
	HasTrailingComma bool
}

// NewArrayLiteralExpressionTree constructs a ArrayLiteralExpressionTree.
//
// Java: ArrayLiteralExpressionTree(SourceRange, ImmutableList, boolean)
func NewArrayLiteralExpressionTree(location *SourceRange, elements []ParseTree, hasTrailingComma bool) *ArrayLiteralExpressionTree {
	t := &ArrayLiteralExpressionTree{}
	t.InitParseTreeBase(t, ParseTreeType_ARRAY_LITERAL_EXPRESSION, location)
	t.Elements = elements
	t.HasTrailingComma = hasTrailingComma
	return t
}
