// Port of com/google/javascript/jscomp/parsing/parser/trees/DoWhileStatementTree.java
// @ closure-compiler v20260415

package trees

type DoWhileStatementTree struct {
	ParseTree
	Body      *ParseTree
	Condition *ParseTree
}

// Go-only: language-adapter
func NewDoWhileStatementTree(location SourceRange, body *ParseTree, condition *ParseTree) *DoWhileStatementTree {
	return &DoWhileStatementTree{
		ParseTree: ParseTree{Type: ParseTreeTypeDoWhileStatement, Location: location},
		Body:      body,
		Condition: condition,
	}
}
