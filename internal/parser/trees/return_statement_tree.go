// Port of com/google/javascript/jscomp/parsing/parser/trees/ReturnStatementTree.java
// @ closure-compiler v20260415

package trees

type ReturnStatementTree struct {
	ParseTree
	Expression *ParseTree
}

// Go-only: language-adapter
func NewReturnStatementTree(location SourceRange, expression *ParseTree) *ReturnStatementTree {
	return &ReturnStatementTree{
		ParseTree:  ParseTree{Type: ParseTreeTypeReturnStatement, Location: location},
		Expression: expression,
	}
}
