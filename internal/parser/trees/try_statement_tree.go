// Port of com/google/javascript/jscomp/parsing/parser/trees/TryStatementTree.java
// @ closure-compiler v20260415

package trees

type TryStatementTree struct {
	Body         ParseTree
	CatchBlock   ParseTree
	FinallyBlock ParseTree
}

func NewTryStatementTree(location SourceRange, body ParseTree, catchBlock ParseTree, finallyBlock ParseTree) *TryStatementTree {
	return &TryStatementTree{
		Body:         body,
		CatchBlock:   catchBlock,
		FinallyBlock: finallyBlock,
	}
}
