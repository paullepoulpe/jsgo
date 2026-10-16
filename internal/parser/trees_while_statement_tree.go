// Port of com/google/javascript/jscomp/parsing/parser/trees/WhileStatementTree.java @ closure-compiler v20260415

package parser

// WhileStatementTree is the parse tree for ParseTreeType WHILE_STATEMENT.
type WhileStatementTree struct {
	ParseTreeBase

	Condition ParseTree
	Body      ParseTree
}

// NewWhileStatementTree constructs a WhileStatementTree.
//
// Java: WhileStatementTree(SourceRange, ParseTree, ParseTree)
func NewWhileStatementTree(location *SourceRange, condition ParseTree, body ParseTree) *WhileStatementTree {
	t := &WhileStatementTree{}
	t.InitParseTreeBase(t, ParseTreeType_WHILE_STATEMENT, location)
	t.Condition = condition
	t.Body = body
	return t
}
