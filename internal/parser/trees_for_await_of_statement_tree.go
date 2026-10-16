// Port of com/google/javascript/jscomp/parsing/parser/trees/ForAwaitOfStatementTree.java @ closure-compiler v20260415

package parser

// Represents a for-await-of statement, which was added in ES2018
type ForAwaitOfStatementTree struct {
	ParseTreeBase

	Initializer ParseTree
	Collection  ParseTree
	Body        ParseTree
}

// NewForAwaitOfStatementTree constructs a ForAwaitOfStatementTree.
//
// Java: ForAwaitOfStatementTree(SourceRange, ParseTree, ParseTree, ParseTree)
func NewForAwaitOfStatementTree(location *SourceRange, initializer ParseTree, collection ParseTree, body ParseTree) *ForAwaitOfStatementTree {
	t := &ForAwaitOfStatementTree{}
	t.InitParseTreeBase(t, ParseTreeType_FOR_AWAIT_OF_STATEMENT, location)
	t.Initializer = initializer
	t.Collection = collection
	t.Body = body
	return t
}
