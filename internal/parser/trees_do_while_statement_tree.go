// Port of com/google/javascript/jscomp/parsing/parser/trees/DoWhileStatementTree.java @ closure-compiler v20260415

package parser

// DoWhileStatementTree is the parse tree for ParseTreeType DO_WHILE_STATEMENT.
type DoWhileStatementTree struct {
	ParseTreeBase

	Body      ParseTree
	Condition ParseTree
}

// NewDoWhileStatementTree constructs a DoWhileStatementTree.
//
// Java: DoWhileStatementTree(SourceRange, ParseTree, ParseTree)
func NewDoWhileStatementTree(location *SourceRange, body ParseTree, condition ParseTree) *DoWhileStatementTree {
	t := &DoWhileStatementTree{}
	t.InitParseTreeBase(t, ParseTreeType_DO_WHILE_STATEMENT, location)
	t.Body = body
	t.Condition = condition
	return t
}
