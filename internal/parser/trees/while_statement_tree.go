// Port of com/google/javascript/jscomp/parsing/parser/trees/WhileStatementTree.java
// @ closure-compiler v20260415

package trees

type WhileStatementTree struct {
	ParseTree
	Condition *ParseTree
	Body      *ParseTree
}

// Go-only: language-adapter
func NewWhileStatementTree(location SourceRange, condition *ParseTree, body *ParseTree) *WhileStatementTree {
	return &WhileStatementTree{
		ParseTree: ParseTree{Type: ParseTreeTypeWhileStatement, Location: location},
		Condition: condition,
		Body:      body,
	}
}
