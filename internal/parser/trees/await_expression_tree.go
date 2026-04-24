// Port of com/google/javascript/jscomp/parsing/parser/trees/AwaitExpressionTree.java
// @ closure-compiler v20260415

package trees

// AwaitExpressionTree represents an await expression.
//
//	AwaitExpression : await UnaryExpression
type AwaitExpressionTree struct {
	ParseTree
	Expression *ParseTree
}

// NewAwaitExpressionTree constructs an AwaitExpressionTree.
// Go-only: language-adapter
func NewAwaitExpressionTree(location SourceRange, expression *ParseTree) *AwaitExpressionTree {
	return &AwaitExpressionTree{
		ParseTree:  ParseTree{Type: ParseTreeTypeAwaitExpression, Location: location},
		Expression: expression,
	}
}
