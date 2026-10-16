// Port of com/google/javascript/jscomp/parsing/parser/trees/UnaryExpressionTree.java @ closure-compiler v20260415

package parser

// UnaryExpressionTree is the parse tree for ParseTreeType UNARY_EXPRESSION.
type UnaryExpressionTree struct {
	ParseTreeBase

	Operator Token
	Operand  ParseTree
}

// NewUnaryExpressionTree constructs a UnaryExpressionTree.
//
// Java: UnaryExpressionTree(SourceRange, Token, ParseTree)
func NewUnaryExpressionTree(location *SourceRange, operator Token, operand ParseTree) *UnaryExpressionTree {
	t := &UnaryExpressionTree{}
	t.InitParseTreeBase(t, ParseTreeType_UNARY_EXPRESSION, location)
	t.Operator = operator
	t.Operand = operand
	return t
}
