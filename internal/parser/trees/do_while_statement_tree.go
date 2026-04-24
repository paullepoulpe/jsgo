// Port of com/google/javascript/jscomp/parsing/parser/trees/DoWhileStatementTree.java
// @ closure-compiler v20260415

package trees

type DoWhileStatementTree struct {
	Body      ParseTree
	Condition ParseTree
}

func NewDoWhileStatementTree(location SourceRange, body ParseTree, condition ParseTree) *DoWhileStatementTree {
	return &DoWhileStatementTree{
		Body:      body,
		Condition: condition,
	}
}
