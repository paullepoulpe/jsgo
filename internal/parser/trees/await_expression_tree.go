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

func NewAwaitExpressionTree(location SourceRange, expression *ParseTree) *AwaitExpressionTree {
	return &AwaitExpressionTree{
		ParseTree:  NewParseTree(ParseTreeTypeAwaitExpression, location),
		Expression: expression,
	}
}
