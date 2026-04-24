// Port of com/google/javascript/jscomp/parsing/parser/trees/ExpressionStatementTree.java
// @ closure-compiler v20260415

package trees

// ExpressionStatementTree represents an expression used as a statement.
type ExpressionStatementTree struct {
	ParseTree
	Expression *ParseTree
}

func NewExpressionStatementTree(location SourceRange, expression *ParseTree) *ExpressionStatementTree {
	return &ExpressionStatementTree{
		ParseTree:  NewParseTree(ParseTreeTypeExpressionStatement, location),
		Expression: expression,
	}
}
