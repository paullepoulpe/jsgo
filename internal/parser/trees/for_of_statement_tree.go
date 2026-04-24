// Port of com/google/javascript/jscomp/parsing/parser/trees/ForOfStatementTree.java
// @ closure-compiler v20260415

package trees

// ForOfStatementTree represents a for-of loop.
type ForOfStatementTree struct {
	ParseTree
	Initializer *ParseTree
	Collection  *ParseTree
	Body        *ParseTree
}

func NewForOfStatementTree(
	location SourceRange,
	initializer *ParseTree,
	collection *ParseTree,
	body *ParseTree,
) *ForOfStatementTree {
	return &ForOfStatementTree{
		ParseTree:   NewParseTree(ParseTreeTypeForOfStatement, location),
		Initializer: initializer,
		Collection:  collection,
		Body:        body,
	}
}
