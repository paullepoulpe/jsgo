// Port of com/google/javascript/jscomp/parsing/parser/trees/ReturnStatementTree.java
// @ closure-compiler v20260415

package trees

type ReturnStatementTree struct {
	Expression ParseTree
}

func NewReturnStatementTree(location SourceRange, expression ParseTree) *ReturnStatementTree {
	return &ReturnStatementTree{
		Expression: expression,
	}
}
