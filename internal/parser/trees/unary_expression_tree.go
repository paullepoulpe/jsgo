// Port of com/google/javascript/jscomp/parsing/parser/trees/UnaryExpressionTree.java
// @ closure-compiler v20260415

package trees

// UnaryExpressionTree represents a unary operator expression.
type UnaryExpressionTree struct {
	ParseTree
	Operator *Token
	Operand  *ParseTree
}

func NewUnaryExpressionTree(location SourceRange, operator *Token, operand *ParseTree) *UnaryExpressionTree {
	return &UnaryExpressionTree{
		ParseTree: NewParseTree(ParseTreeTypeUnaryExpression, location),
		Operator:  operator,
		Operand:   operand,
	}
}
