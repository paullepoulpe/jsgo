// Port of com/google/javascript/jscomp/parsing/parser/trees/ExpressionStatementTree.java @ closure-compiler v20260415

package parser

// ExpressionStatementTree is the parse tree for ParseTreeType EXPRESSION_STATEMENT.
type ExpressionStatementTree struct {
	ParseTreeBase

	Expression ParseTree
}

// NewExpressionStatementTree constructs a ExpressionStatementTree.
//
// Java: ExpressionStatementTree(SourceRange, ParseTree)
func NewExpressionStatementTree(location *SourceRange, expression ParseTree) *ExpressionStatementTree {
	t := &ExpressionStatementTree{}
	t.InitParseTreeBase(t, ParseTreeType_EXPRESSION_STATEMENT, location)
	t.Expression = expression
	return t
}
