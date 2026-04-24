// Port of com/google/javascript/jscomp/parsing/parser/trees/BreakStatementTree.java
// @ closure-compiler v20260415

package trees

type BreakStatementTree struct {
	ParseTree
	Name *IdentifierToken
}

// Go-only: language-adapter
func NewBreakStatementTree(location SourceRange, name *IdentifierToken) *BreakStatementTree {
	return &BreakStatementTree{
		ParseTree: ParseTree{Type: ParseTreeTypeBreakStatement, Location: location},
		Name:      name,
	}
}

func (c *BreakStatementTree) HasLabel() bool {
	return c.Name != nil
}
