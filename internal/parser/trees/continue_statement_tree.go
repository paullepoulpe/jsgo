// Port of com/google/javascript/jscomp/parsing/parser/trees/ContinueStatementTree.java
// @ closure-compiler v20260415

package trees

type ContinueStatementTree struct {
	ParseTree
	Name *IdentifierToken
}

func NewContinueStatementTree(location SourceRange, name *IdentifierToken) *ContinueStatementTree {
	return &ContinueStatementTree{
		ParseTree: ParseTree{Type: ParseTreeTypeContinueStatement, Location: location},
		Name:      name,
	}
}

func (c *ContinueStatementTree) HasLabel() bool {
	return c.Name != nil
}
