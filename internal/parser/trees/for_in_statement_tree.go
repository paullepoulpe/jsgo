// Port of com/google/javascript/jscomp/parsing/parser/trees/ForInStatementTree.java
// @ closure-compiler v20260415

package trees

// ForInStatementTree represents a for-in loop.
type ForInStatementTree struct {
	ParseTree
	Initializer *ParseTree
	Collection  *ParseTree
	Body        *ParseTree
}

func NewForInStatementTree(
	location SourceRange,
	initializer *ParseTree,
	collection *ParseTree,
	body *ParseTree,
) *ForInStatementTree {
	return &ForInStatementTree{
		ParseTree:   NewParseTree(ParseTreeTypeForInStatement, location),
		Initializer: initializer,
		Collection:  collection,
		Body:        body,
	}
}
