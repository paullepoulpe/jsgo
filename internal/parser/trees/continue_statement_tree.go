// Port of com/google/javascript/jscomp/parsing/parser/trees/ContinueStatementTree.java
// @ closure-compiler v20260415

package trees

// ContinueStatementTree represents a continue statement, with an optional label.
type ContinueStatementTree struct {
	ParseTree
	// Name is nil when the continue has no label.
	Name *IdentifierToken
}

func NewContinueStatementTree(location SourceRange, name *IdentifierToken) *ContinueStatementTree {
	return &ContinueStatementTree{
		ParseTree: NewParseTree(ParseTreeTypeContinueStatement, location),
		Name:      name,
	}
}

func (t *ContinueStatementTree) HasLabel() bool {
	return t.Name != nil
}
