// Port of com/google/javascript/jscomp/parsing/parser/trees/ForInStatementTree.java @ closure-compiler v20260415

package parser

// ForInStatementTree is the parse tree for ParseTreeType FOR_IN_STATEMENT.
type ForInStatementTree struct {
	ParseTreeBase

	Initializer ParseTree
	Collection  ParseTree
	Body        ParseTree
}

// NewForInStatementTree constructs a ForInStatementTree.
//
// Java: ForInStatementTree(SourceRange, ParseTree, ParseTree, ParseTree)
func NewForInStatementTree(location *SourceRange, initializer ParseTree, collection ParseTree, body ParseTree) *ForInStatementTree {
	t := &ForInStatementTree{}
	t.InitParseTreeBase(t, ParseTreeType_FOR_IN_STATEMENT, location)
	t.Initializer = initializer
	t.Collection = collection
	t.Body = body
	return t
}
