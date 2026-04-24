// Port of com/google/javascript/jscomp/parsing/parser/trees/ForAwaitOfStatementTree.java
// @ closure-compiler v20260415

package trees

// ForAwaitOfStatementTree represents a for-await-of statement (ES2018).
type ForAwaitOfStatementTree struct {
	ParseTree
	Initializer *ParseTree
	Collection  *ParseTree
	Body        *ParseTree
}

func NewForAwaitOfStatementTree(
	location SourceRange,
	initializer *ParseTree,
	collection *ParseTree,
	body *ParseTree,
) *ForAwaitOfStatementTree {
	return &ForAwaitOfStatementTree{
		ParseTree:   NewParseTree(ParseTreeTypeForAwaitOfStatement, location),
		Initializer: initializer,
		Collection:  collection,
		Body:        body,
	}
}
