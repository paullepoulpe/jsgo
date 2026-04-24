// Port of com/google/javascript/jscomp/parsing/parser/trees/BreakStatementTree.java
// @ closure-compiler v20260415

package trees

// BreakStatementTree represents a break statement, with an optional label.
type BreakStatementTree struct {
	ParseTree
	// Name is nil when the break has no label.
	Name *IdentifierToken
}

func NewBreakStatementTree(location SourceRange, name *IdentifierToken) *BreakStatementTree {
	return &BreakStatementTree{
		ParseTree: NewParseTree(ParseTreeTypeBreakStatement, location),
		Name:      name,
	}
}

func (t *BreakStatementTree) HasLabel() bool {
	return t.Name != nil
}
