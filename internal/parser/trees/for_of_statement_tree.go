// Port of com/google/javascript/jscomp/parsing/parser/trees/ForOfStatementTree.java
// @ closure-compiler v20260415

package trees

type ForOfStatementTree struct {
	ParseTree
	Initializer *ParseTree
	Collection  *ParseTree
	Body        *ParseTree
}

// Go-only: language-adapter
func NewForOfStatementTree(location SourceRange, initializer *ParseTree, collection *ParseTree, body *ParseTree) *ForOfStatementTree {
	return &ForOfStatementTree{
		ParseTree:   ParseTree{Type: ParseTreeTypeForOfStatement, Location: location},
		Initializer: initializer,
		Collection:  collection,
		Body:        body,
	}
}
