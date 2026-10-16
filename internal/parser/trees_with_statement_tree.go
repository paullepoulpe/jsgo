// Port of com/google/javascript/jscomp/parsing/parser/trees/WithStatementTree.java @ closure-compiler v20260415

package parser

// WithStatementTree is the parse tree for ParseTreeType WITH_STATEMENT.
type WithStatementTree struct {
	ParseTreeBase

	Expression ParseTree
	Body       ParseTree
}

// NewWithStatementTree constructs a WithStatementTree.
//
// Java: WithStatementTree(SourceRange, ParseTree, ParseTree)
func NewWithStatementTree(location *SourceRange, expression ParseTree, body ParseTree) *WithStatementTree {
	t := &WithStatementTree{}
	t.InitParseTreeBase(t, ParseTreeType_WITH_STATEMENT, location)
	t.Expression = expression
	t.Body = body
	return t
}
