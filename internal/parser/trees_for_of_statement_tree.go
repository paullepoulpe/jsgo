// Port of com/google/javascript/jscomp/parsing/parser/trees/ForOfStatementTree.java @ closure-compiler v20260415

package parser

// ForOfStatementTree is the parse tree for ParseTreeType FOR_OF_STATEMENT.
type ForOfStatementTree struct {
	ParseTreeBase

	Initializer ParseTree
	Collection  ParseTree
	Body        ParseTree
}

// NewForOfStatementTree constructs a ForOfStatementTree.
//
// Java: ForOfStatementTree(SourceRange, ParseTree, ParseTree, ParseTree)
func NewForOfStatementTree(location *SourceRange, initializer ParseTree, collection ParseTree, body ParseTree) *ForOfStatementTree {
	t := &ForOfStatementTree{}
	t.InitParseTreeBase(t, ParseTreeType_FOR_OF_STATEMENT, location)
	t.Initializer = initializer
	t.Collection = collection
	t.Body = body
	return t
}
