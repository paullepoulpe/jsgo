// Port of com/google/javascript/jscomp/parsing/parser/trees/DefaultClauseTree.java
// @ closure-compiler v20260415

package trees

// DefaultClauseTree represents the default clause in a switch statement.
type DefaultClauseTree struct {
	ParseTree
	Statements []*ParseTree
}

func NewDefaultClauseTree(location SourceRange, statements []*ParseTree) *DefaultClauseTree {
	return &DefaultClauseTree{
		ParseTree:  NewParseTree(ParseTreeTypeDefaultClause, location),
		Statements: statements,
	}
}
