// Port of com/google/javascript/jscomp/parsing/parser/trees/ReturnStatementTree.java @ closure-compiler v20260415

package parser

// ReturnStatementTree is the parse tree for ParseTreeType RETURN_STATEMENT.
type ReturnStatementTree struct {
	ParseTreeBase

	Expression ParseTree
}

// NewReturnStatementTree constructs a ReturnStatementTree.
//
// Java: ReturnStatementTree(SourceRange, ParseTree)
func NewReturnStatementTree(location *SourceRange, expression ParseTree) *ReturnStatementTree {
	t := &ReturnStatementTree{}
	t.InitParseTreeBase(t, ParseTreeType_RETURN_STATEMENT, location)
	t.Expression = expression
	return t
}
