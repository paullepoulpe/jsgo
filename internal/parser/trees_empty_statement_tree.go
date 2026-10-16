// Port of com/google/javascript/jscomp/parsing/parser/trees/EmptyStatementTree.java @ closure-compiler v20260415

package parser

// EmptyStatementTree is the parse tree for ParseTreeType EMPTY_STATEMENT.
type EmptyStatementTree struct {
	ParseTreeBase
}

// NewEmptyStatementTree constructs a EmptyStatementTree.
//
// Java: EmptyStatementTree(SourceRange)
func NewEmptyStatementTree(location *SourceRange) *EmptyStatementTree {
	t := &EmptyStatementTree{}
	t.InitParseTreeBase(t, ParseTreeType_EMPTY_STATEMENT, location)
	return t
}
