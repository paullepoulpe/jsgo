// Port of com/google/javascript/jscomp/parsing/parser/trees/EmptyStatementTree.java
// @ closure-compiler v20260415

package trees

// EmptyStatementTree represents an empty statement (a lone semicolon).
type EmptyStatementTree struct {
	ParseTree
}

func NewEmptyStatementTree(location SourceRange) *EmptyStatementTree {
	return &EmptyStatementTree{
		ParseTree: NewParseTree(ParseTreeTypeEmptyStatement, location),
	}
}
