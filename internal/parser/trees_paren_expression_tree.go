// Port of com/google/javascript/jscomp/parsing/parser/trees/ParenExpressionTree.java @ closure-compiler v20260415

package parser

// ParenExpressionTree is the parse tree for ParseTreeType PAREN_EXPRESSION.
type ParenExpressionTree struct {
	ParseTreeBase

	Expression ParseTree
}

// NewParenExpressionTree constructs a ParenExpressionTree.
//
// Java: ParenExpressionTree(SourceRange, ParseTree)
func NewParenExpressionTree(location *SourceRange, expression ParseTree) *ParenExpressionTree {
	t := &ParenExpressionTree{}
	t.InitParseTreeBase(t, ParseTreeType_PAREN_EXPRESSION, location)
	t.Expression = expression
	return t
}
