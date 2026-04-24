// Port of com/google/javascript/jscomp/parsing/parser/trees/BreakStatementTree.java
// @ closure-compiler v20260415

package trees

type BreakStatementTree struct {
	Name *IdentifierToken
}

func NewBreakStatementTree(location SourceRange, name *IdentifierToken) *BreakStatementTree {
	return &BreakStatementTree{
		Name: name,
	}
}

func (c *BreakStatementTree) HasLabel() bool {
	return c.Name != nil
}
