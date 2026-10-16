// Port of com/google/javascript/jscomp/parsing/parser/trees/AwaitExpressionTree.java @ closure-compiler v20260415

package parser

// Represents an await expression.
//
//	AwaitExpression : await UnaryExpression
type AwaitExpressionTree struct {
	ParseTreeBase

	Expression ParseTree
}

// NewAwaitExpressionTree constructs a AwaitExpressionTree.
//
// Java: AwaitExpressionTree(SourceRange, ParseTree)
func NewAwaitExpressionTree(location *SourceRange, expression ParseTree) *AwaitExpressionTree {
	t := &AwaitExpressionTree{}
	t.InitParseTreeBase(t, ParseTreeType_AWAIT_EXPRESSION, location)
	t.Expression = expression
	return t
}
