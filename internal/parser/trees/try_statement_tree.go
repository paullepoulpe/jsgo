// Port of com/google/javascript/jscomp/parsing/parser/trees/TryStatementTree.java
// @ closure-compiler v20260415

package trees

type TryStatementTree struct {
	ParseTree
	Body         *ParseTree
	CatchBlock   *ParseTree
	FinallyBlock *ParseTree
}

func NewTryStatementTree(location SourceRange, body *ParseTree, catchBlock *ParseTree, finallyBlock *ParseTree) *TryStatementTree {
	return &TryStatementTree{
		ParseTree:    ParseTree{Type: ParseTreeTypeTryStatement, Location: location},
		Body:         body,
		CatchBlock:   catchBlock,
		FinallyBlock: finallyBlock,
	}
}
