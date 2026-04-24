// Port of com/google/javascript/jscomp/parsing/parser/trees/ForAwaitOfStatementTree.java
// @ closure-compiler v20260415

package trees

// ForAwaitOfStatementTree represents `for await (lhs of rhs) body`
// (ES2018 async iteration).
type ForAwaitOfStatementTree struct {
	ParseTree
	Initializer *ParseTree
	Collection  *ParseTree
	Body        *ParseTree
}

// NewForAwaitOfStatementTree constructs a ForAwaitOfStatementTree.
func NewForAwaitOfStatementTree(location SourceRange, initializer, collection, body *ParseTree) *ForAwaitOfStatementTree {
	return &ForAwaitOfStatementTree{
		ParseTree:   ParseTree{Type: ParseTreeTypeForAwaitOfStatement, Location: location},
		Initializer: initializer,
		Collection:  collection,
		Body:        body,
	}
}
