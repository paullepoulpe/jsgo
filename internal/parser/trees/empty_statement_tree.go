// Port of com/google/javascript/jscomp/parsing/parser/trees/EmptyStatementTree.java
// @ closure-compiler v20260415

package trees

type EmptyStatementTree struct {
	ParseTree
}

// Go-only: language-adapter
func NewEmptyStatementTree(location SourceRange) *EmptyStatementTree {
	return &EmptyStatementTree{
		ParseTree: ParseTree{Type: ParseTreeTypeEmptyStatement, Location: location},
	}
}
