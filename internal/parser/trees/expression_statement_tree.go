// Port of com/google/javascript/jscomp/parsing/parser/trees/ExpressionStatementTree.java
// @ closure-compiler v20260415

package trees

type ExpressionStatementTree struct {
	ParseTree
	Expression *ParseTree
}

func NewExpressionStatementTree(location SourceRange, expression *ParseTree) *ExpressionStatementTree {
	return &ExpressionStatementTree{
		ParseTree:  ParseTree{Type: ParseTreeTypeExpressionStatement, Location: location},
		Expression: expression,
	}
}
