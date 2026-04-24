// Port of com/google/javascript/jscomp/parsing/parser/trees/UnaryExpressionTree.java
// @ closure-compiler v20260415

package trees

// UnaryExpressionTree represents a unary operator expression, e.g. !a, -b, typeof c.
type UnaryExpressionTree struct {
	ParseTree
	Operator *Token
	Operand  *ParseTree
}

// NewUnaryExpressionTree constructs a UnaryExpressionTree.
func NewUnaryExpressionTree(location SourceRange, operator *Token, operand *ParseTree) *UnaryExpressionTree {
	return &UnaryExpressionTree{
		ParseTree: ParseTree{Type: ParseTreeTypeUnaryExpression, Location: location},
		Operator:  operator,
		Operand:   operand,
	}
}
