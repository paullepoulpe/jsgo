// Port of com/google/javascript/jscomp/parsing/parser/trees/WhileStatementTree.java
// @ closure-compiler v20260415

package trees

type WhileStatementTree struct {
	Condition ParseTree
	Body      ParseTree
}

func NewWhileStatementTree(location SourceRange, condition ParseTree, body ParseTree) *WhileStatementTree {
	return &WhileStatementTree{
		Condition: condition,
		Body:      body,
	}
}
