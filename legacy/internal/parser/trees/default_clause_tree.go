// Port of com/google/javascript/jscomp/parsing/parser/trees/DefaultClauseTree.java
// @ closure-compiler v20260415

package trees

type DefaultClauseTree struct {
	ParseTree
	Statements []*ParseTree
}

// Go-only: language-adapter
func NewDefaultClauseTree(location SourceRange, statements []*ParseTree) *DefaultClauseTree {
	return &DefaultClauseTree{
		ParseTree:  ParseTree{Type: ParseTreeTypeDefaultClause, Location: location},
		Statements: statements,
	}
}
