// Port of com/google/javascript/jscomp/parsing/parser/trees/WithStatementTree.java
// @ closure-compiler v20260415

package trees

type WithStatementTree struct {
	Expression ParseTree
	Body       ParseTree
}

func NewWithStatementTree(location SourceRange, expression ParseTree, body ParseTree) *WithStatementTree {
	return &WithStatementTree{
		Expression: expression,
		Body:       body,
	}
}
