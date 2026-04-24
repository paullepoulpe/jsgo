// Port of com/google/javascript/jscomp/parsing/parser/trees/ReturnStatementTree.java
// @ closure-compiler v20260415

package trees

// ReturnStatementTree represents a return statement with an optional expression.
type ReturnStatementTree struct {
	ParseTree
	// Expression is nil for bare return statements.
	Expression *ParseTree
}

func NewReturnStatementTree(location SourceRange, expression *ParseTree) *ReturnStatementTree {
	return &ReturnStatementTree{
		ParseTree:  NewParseTree(ParseTreeTypeReturnStatement, location),
		Expression: expression,
	}
}
