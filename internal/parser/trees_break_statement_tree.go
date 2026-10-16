// Port of com/google/javascript/jscomp/parsing/parser/trees/BreakStatementTree.java @ closure-compiler v20260415

package parser

// BreakStatementTree is the parse tree for ParseTreeType BREAK_STATEMENT.
type BreakStatementTree struct {
	ParseTreeBase

	Name *IdentifierToken
}

// NewBreakStatementTree constructs a BreakStatementTree.
//
// Java: BreakStatementTree(SourceRange, IdentifierToken)
func NewBreakStatementTree(location *SourceRange, name *IdentifierToken) *BreakStatementTree {
	t := &BreakStatementTree{}
	t.InitParseTreeBase(t, ParseTreeType_BREAK_STATEMENT, location)
	t.Name = name
	return t
}

// HasLabel reports whether the statement names a label.
//
// Java: BreakStatementTree.hasLabel()
func (t *BreakStatementTree) HasLabel() bool {
	return t.Name != nil
}
