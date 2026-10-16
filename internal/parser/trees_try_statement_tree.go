// Port of com/google/javascript/jscomp/parsing/parser/trees/TryStatementTree.java @ closure-compiler v20260415

package parser

// TryStatementTree is the parse tree for ParseTreeType TRY_STATEMENT.
type TryStatementTree struct {
	ParseTreeBase

	Body         ParseTree
	CatchBlock   ParseTree
	FinallyBlock ParseTree
}

// NewTryStatementTree constructs a TryStatementTree.
//
// Java: TryStatementTree(SourceRange, ParseTree, ParseTree, ParseTree)
func NewTryStatementTree(location *SourceRange, body ParseTree, catchBlock ParseTree, finallyBlock ParseTree) *TryStatementTree {
	t := &TryStatementTree{}
	t.InitParseTreeBase(t, ParseTreeType_TRY_STATEMENT, location)
	t.Body = body
	t.CatchBlock = catchBlock
	t.FinallyBlock = finallyBlock
	return t
}
