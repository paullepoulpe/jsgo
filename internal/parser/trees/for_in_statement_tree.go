// Port of com/google/javascript/jscomp/parsing/parser/trees/ForInStatementTree.java
// @ closure-compiler v20260415

package trees

type ForInStatementTree struct {
	ParseTree
	Initializer *ParseTree
	Collection  *ParseTree
	Body        *ParseTree
}

// Go-only: language-adapter
func NewForInStatementTree(location SourceRange, initializer *ParseTree, collection *ParseTree, body *ParseTree) *ForInStatementTree {
	return &ForInStatementTree{
		ParseTree:   ParseTree{Type: ParseTreeTypeForInStatement, Location: location},
		Initializer: initializer,
		Collection:  collection,
		Body:        body,
	}
}
