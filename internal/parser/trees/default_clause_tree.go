// Port of com/google/javascript/jscomp/parsing/parser/trees/DefaultClauseTree.java
// @ closure-compiler v20260415

package trees

type DefaultClauseTree struct {
	Statements []ParseTree
}

func NewDefaultClauseTree(location SourceRange, statements []ParseTree) *DefaultClauseTree {
	return &DefaultClauseTree{
		Statements: statements,
	}
}
