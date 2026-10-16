// Port of com/google/javascript/jscomp/parsing/parser/trees/ContinueStatementTree.java @ closure-compiler v20260415

package parser

// ContinueStatementTree is the parse tree for ParseTreeType CONTINUE_STATEMENT.
type ContinueStatementTree struct {
	ParseTreeBase

	Name *IdentifierToken
}

// NewContinueStatementTree constructs a ContinueStatementTree.
//
// Java: ContinueStatementTree(SourceRange, IdentifierToken)
func NewContinueStatementTree(location *SourceRange, name *IdentifierToken) *ContinueStatementTree {
	t := &ContinueStatementTree{}
	t.InitParseTreeBase(t, ParseTreeType_CONTINUE_STATEMENT, location)
	t.Name = name
	return t
}

// HasLabel reports whether the statement names a label.
//
// Java: ContinueStatementTree.hasLabel()
func (t *ContinueStatementTree) HasLabel() bool {
	return t.Name != nil
}
