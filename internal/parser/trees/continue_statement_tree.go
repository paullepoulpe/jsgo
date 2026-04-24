// Port of com/google/javascript/jscomp/parsing/parser/trees/ContinueStatementTree.java
// @ closure-compiler v20260415

package trees

type ContinueStatementTree struct {
	Name *IdentifierToken
}

func NewContinueStatementTree(location SourceRange, name *IdentifierToken) *ContinueStatementTree {
	return &ContinueStatementTree{
		Name: name,
	}
}

func (c *ContinueStatementTree) HasLabel() bool {
	return c.Name != nil
}
