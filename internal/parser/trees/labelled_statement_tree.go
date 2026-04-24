// Port of com/google/javascript/jscomp/parsing/parser/trees/LabelledStatementTree.java
// @ closure-compiler v20260415

package trees

type LabelledStatementTree struct {
	Name      *IdentifierToken
	Statement ParseTree
}

func NewLabelledStatementTree(location SourceRange, name *IdentifierToken, statement ParseTree) *LabelledStatementTree {
	return &LabelledStatementTree{
		Name:      name,
		Statement: statement,
	}
}
